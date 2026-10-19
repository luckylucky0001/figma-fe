package reuse

import (
	"os"
	"path/filepath"
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
	Colors      []Finding `json:"colors"`      // hard-coded colours (with token if one exists)
	Pass        bool      `json:"pass"`
}

// CheckFocus evaluates files under focus (relative to the report root).
func (r *Report) CheckFocus(focus string) *PageCheck {
	focus = filepath.ToSlash(strings.TrimSuffix(focus, "/"))
	in := func(rel string) bool { return rel == focus || strings.HasPrefix(rel, focus+"/") }
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
	for _, h := range r.Colors {
		if h.Token != "" {
			tokens[h.Value] = h.Token
		}
	}
	for k, v := range r.ColorTokens {
		tokens[k] = v
	}
	_ = filepath.Walk(filepath.Join(r.Root, focus), func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		b, _ := os.ReadFile(p)
		src := string(b)
		rel := relPath(r.Root, p)
		for _, loc := range hexColorRe.FindAllStringIndex(src, -1) {
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
		return nil
	})
	pc.Pass = len(pc.Reinvented) == 0 && len(pc.MissedReuse) == 0
	return pc
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
