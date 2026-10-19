package reuse

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// PageCheck is the brownfield verification of one new page: the whole project is the asset pool,
// only the files under Focus are judged.
type PageCheck struct {
	Focus       string    `json:"focus"`
	Used        []string  `json:"used"`        // existing assets the page uses ("Name(kind)")
	Reinvented  []Finding `json:"reinvented"`  // components defined in the page that already exist
	MissedReuse []Finding `json:"missedReuse"` // raw primitives where a component exists
	HandBuilt   []Finding `json:"handBuilt"`   // markup that hand-builds a concept an existing component covers
	Colors      []Finding `json:"colors"`      // hard-coded colours (with token if one exists)
	LocalTokens int       `json:"localTokens"` // colour variables defined inside the page itself
	Pass        bool      `json:"pass"`
}

// CheckFocus evaluates files under focus (relative to the report root).
func (r *Report) CheckFocus(focus string) *PageCheck {
	var focuses []string // comma-separated: a page is often several files (Fragment + layout)
	for _, f := range strings.Split(focus, ",") {
		if f = filepath.ToSlash(strings.TrimSuffix(strings.TrimSpace(f), "/")); f != "" {
			focuses = append(focuses, f)
		}
	}
	in := func(rel string) bool {
		for _, f := range focuses {
			if rel == f || strings.HasPrefix(rel, f+"/") {
				return true
			}
		}
		return false
	}
	pc := &PageCheck{Focus: focus}
	pool := append([]Asset{}, r.Assets...)
	pool = append(pool, r.LibAll...)
	for _, a := range r.Assets {
		for _, f := range a.UsageFiles {
			if in(f) && !(in(a.Source) && a.Kind != "ui-lib") {
				pc.Used = append(pc.Used, a.Name+"（"+kindLabel(a.Kind)+"）")
				break
			}
		}
	}
	sort.Strings(pc.Used)
	for _, a := range r.Assets {
		if a.Kind == "ui-lib" || !in(a.Source) || a.Concept == "" {
			continue
		}
		var existing *Asset
		for i := range pool {
			b := &pool[i]
			if b.Concept == a.Concept && !in(b.Source) && b.Source != "Android SDK widgets" && (existing == nil || rank(*b) > rank(*existing)) {
				existing = b
			}
		}
		if existing != nil {
			pc.Reinvented = append(pc.Reinvented, Finding{"reinvented", "error", a.Source, 0, a.Name, "页面里自造了 " + a.Name + "，已有 " + existing.Name + "（" + kindLabel(existing.Kind) + " " + existing.Source + "）"})
		}
	}
	for _, f := range r.Findings {
		if f.Rule == "missed-reuse" && in(f.File) {
			pc.MissedReuse = append(pc.MissedReuse, f)
		}
	}
	tokens := map[string]string{}
	for k, v := range r.ColorTokens {
		tokens[k] = v
	}
	usedConcepts := map[string]bool{}
	for _, a := range r.Assets {
		for _, f := range a.UsageFiles {
			if in(f) {
				usedConcepts[a.Concept] = true
			}
		}
	}
	var focusSrc []struct{ rel, src string }
	for _, fo := range focuses {
		_ = filepath.Walk(filepath.Join(r.Root, fo), func(p string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				b, _ := os.ReadFile(p)
				focusSrc = append(focusSrc, struct{ rel, src string }{relPath(r.Root, p), string(b)})
			}
			return nil
		})
	}
	// colour variables defined in the page are not project tokens
	for _, f := range focusSrc {
		for _, re := range []*regexp.Regexp{cssVarDefRe, preVarDefRe} {
			for _, m := range re.FindAllStringSubmatch(f.src, -1) {
				pc.LocalTokens++
				v := strings.ToUpper(expandHex(m[2]))
				if t := tokens[v]; t == m[1] || t == "var("+m[1]+")" {
					delete(tokens, v)
				}
			}
		}
	}
	// hand-built concepts: class names / selectors that express a concept the page does not use a component for
	seenHB := map[string]bool{}
	for _, f := range focusSrc {
		for _, m := range classAttrRe.FindAllStringSubmatchIndex(f.src, -1) {
			for _, cls := range strings.Fields(f.src[m[2]:m[3]]) {
				checkHandBuilt(pc, r, pool, usedConcepts, seenHB, f.rel, strings.Count(f.src[:m[0]], "\n")+1, cls, in)
			}
		}
		for _, m := range androidIDRe.FindAllStringSubmatchIndex(f.src, -1) {
			checkHandBuilt(pc, r, pool, usedConcepts, seenHB, f.rel, strings.Count(f.src[:m[0]], "\n")+1, f.src[m[2]:m[3]], in)
		}
		for _, m := range cssSelRe2.FindAllStringSubmatchIndex(f.src, -1) {
			checkHandBuilt(pc, r, pool, usedConcepts, seenHB, f.rel, strings.Count(f.src[:m[0]], "\n")+1, f.src[m[2]:m[3]], in)
		}
	}
	for _, f := range focusSrc {
		src, rel := f.src, f.rel
		for _, loc := range hexColorRe.FindAllStringIndex(src, -1) {
			if lineHasVarDef(src, loc[0]) {
				continue
			}
			v := strings.ToUpper(expandHex(src[loc[0]:loc[1]]))
			ln := strings.Count(src[:loc[0]], "\n") + 1
			hint := "写死颜色"
			if t := tokens[v]; t != "" {
				hint += "，已有同值 token " + t
			}
			pc.Colors = append(pc.Colors, Finding{"hardcoded-color", "warn", rel, ln, v, hint})
		}
		for _, m := range ktColorRe2.FindAllStringSubmatchIndex(src, -1) {
			ln := strings.Count(src[:m[0]], "\n") + 1
			pc.Colors = append(pc.Colors, Finding{"hardcoded-color", "warn", rel, ln, src[m[0]:m[1]], "写死颜色"})
		}
	}
	pc.Pass = len(pc.Reinvented) == 0 && len(pc.MissedReuse) == 0
	for _, f := range pc.HandBuilt {
		if f.Severity == "error" {
			pc.Pass = false
		}
	}
	return pc
}

var (
	classAttrRe = regexp.MustCompile(`\bclass(?:Name)?="([^"]+)"`)
	cssSelRe2   = regexp.MustCompile(`(?m)^\s*\.([A-Za-z][\w-]*)[^{\n]*\{`)
	androidIDRe = regexp.MustCompile(`android:id="@\+id/(\w+)"`)
)

func lineHasVarDef(src string, idx int) bool {
	ls := strings.LastIndex(src[:idx], "\n") + 1
	le := strings.IndexByte(src[idx:], '\n')
	if le < 0 {
		le = len(src) - idx
	}
	line := src[ls : idx+le]
	return cssVarDefRe.MatchString(line) || preVarDefRe.MatchString(line)
}

func checkHandBuilt(pc *PageCheck, r *Report, pool []Asset, used, seen map[string]bool, rel string, ln int, cls string, in func(string) bool) {
	concept := ConceptOf(cls)
	if concept == "" || used[concept] || seen[concept] || concept == "list" || concept == "icon" || concept == "image" {
		return
	}
	// same preference as match: proven in this project, then the concept's canonical component
	pref := func(a *Asset) int {
		p := rank(*a)
		switch {
		case a.Usage >= 3:
			p += 3000 // the project's established choice
		case a.Usage > 0:
			p += 1000
		}
		for _, c := range concepts {
			if c.Name == concept && normName(a.Name) == c.Words[0] {
				p += 1500 // the concept's canonical component (Cell, Tag, NavBar)
			}
		}
		return p
	}
	var best *Asset
	for i := range pool {
		b := &pool[i]
		if b.Concept == concept && !in(b.Source) && b.Source != "Android SDK widgets" && (best == nil || pref(b) > pref(best) || pref(b) == pref(best) && len(b.Name) < len(best.Name)) {
			best = b
		}
	}
	if best == nil {
		return
	}
	seen[concept] = true
	sev, note := "warn", "（组件库标准组件，项目还没用过，请确认）"
	if best.Usage > 0 {
		sev, note = "error", "（项目里已用 "+itoa(best.Usage)+" 次）"
	}
	pc.HandBuilt = append(pc.HandBuilt, Finding{"hand-built", sev, rel, ln, "." + cls, "页面手写了「" + concept + "」（." + cls + "），已有 " + best.Name + "（" + kindLabel(best.Kind) + " " + best.Source + "）" + note})
}

func (pc *PageCheck) Markdown() string {
	var b strings.Builder
	verdict := "✅ PASS"
	if !pc.Pass {
		verdict = "❌ FAIL"
	}
	b.WriteString("\n## 页面复用检查：" + pc.Focus + "  " + verdict + "\n\n")
	b.WriteString("- 复用的现成组件（" + itoa(len(pc.Used)) + "）：" + strings.Join(pc.Used, "、") + "\n")
	b.WriteString("- 自造轮子（" + itoa(len(pc.Reinvented)) + "）\n")
	for _, f := range pc.Reinvented {
		b.WriteString("  - ✗ " + f.File + "：" + f.Hint + "\n")
	}
	b.WriteString("- 有组件却没用（" + itoa(len(pc.MissedReuse)) + "）\n")
	for _, f := range pc.MissedReuse {
		b.WriteString("  - ✗ " + f.File + ":" + itoa(f.Line) + " `" + f.Text + "` → " + f.Hint + "\n")
	}
	b.WriteString("- 疑似手写已有组件（" + itoa(len(pc.HandBuilt)) + "）\n")
	for _, f := range pc.HandBuilt {
		mark := "·"
		if f.Severity == "error" {
			mark = "✗"
		}
		b.WriteString("  - " + mark + " " + f.File + ":" + itoa(f.Line) + " " + f.Hint + "\n")
	}
	if pc.LocalTokens > 0 {
		b.WriteString("- 页面私有调色板：页面里自己定义了 " + itoa(pc.LocalTokens) + " 个颜色变量（应使用项目 / 组件库已有的 token）\n")
	}
	b.WriteString("- 写死颜色（" + itoa(len(pc.Colors)) + "）\n")
	for i, f := range pc.Colors {
		if i >= 10 {
			b.WriteString("  - … 共 " + itoa(len(pc.Colors)) + " 处\n")
			break
		}
		b.WriteString("  - · " + f.File + ":" + itoa(f.Line) + " " + f.Text + "：" + f.Hint + "\n")
	}
	return b.String()
}
