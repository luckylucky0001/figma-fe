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
	Focus        string    `json:"focus"`
	Used         []string  `json:"used"`                   // existing assets the page uses ("Name(kind)")
	Reinvented   []Finding `json:"reinvented"`             // components defined in the page that already exist
	MissedReuse  []Finding `json:"missedReuse"`            // raw primitives where a component exists
	HandBuilt    []Finding `json:"handBuilt"`              // markup that hand-builds a concept an existing component covers
	Inconsistent []Finding `json:"inconsistent,omitempty"` // a component new to the project where it has an established one
	Colors       []Finding `json:"colors"`                 // hard-coded colours (with token if one exists)
	LocalTokens  int       `json:"localTokens"`            // colour variables defined inside the page itself
	Pass         bool      `json:"pass"`
}

// CheckFocus evaluates files under focus (relative to the report root).
func (r *Report) CheckFocus(focus string) *PageCheck { return r.CheckFocusLines(focus, nil) }

// CheckFocusLines is CheckFocus restricted, per file, to the given lines (a PR's added lines);
// files not in lines are judged whole.
func (r *Report) CheckFocusLines(focus string, lines map[string]map[int]bool) *PageCheck {
	inLines := func(rel string, ln int) bool {
		set, ok := lines[rel]
		return !ok || set[ln]
	}
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
	setFamilyPrefixes(r.Assets)
	var pool []Asset
	for _, a := range append(append([]Asset{}, r.Assets...), r.LibAll...) {
		if usable(&a) {
			pool = append(pool, a)
		}
	}
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
		if _, modified := lines[a.Source]; modified {
			continue // the component predates this change; only newly added components are judged
		}
		// an existing component of the same concept that is generic or made for the same purpose
		// (StatusTag vs Tag / OrderStatusTag) is reinvention; one made for another purpose
		// (StatusTag vs FilterChip) only deserves a look
		if accepted(a.Name, r.Platform) {
			continue // a recorded decision in component-map.yaml settles it
		}
		mine := specialWords(a.Name, a.Concept)
		var existing, other *Asset
		for i := range pool {
			b := &pool[i]
			if b.Concept != a.Concept || in(b.Source) || b.Source == "Android SDK widgets" {
				continue
			}
			if subset(specialWords(b.Name, b.Concept), mine) {
				if existing == nil || preference(*b, a.Concept) > preference(*existing, a.Concept) {
					existing = b
				}
			} else if other == nil || preference(*b, a.Concept) > preference(*other, a.Concept) {
				other = b
			}
		}
		if existing == nil && other != nil {
			pc.Reinvented = append(pc.Reinvented, Finding{"reinvented", "warn", a.Source, 0, a.Name, "新建了 " + a.Name + "；同类已有 " + other.Name + "（" + kindLabel(other.Kind) + "），用途可能不同（" + strings.Join(specialWords(other.Name, other.Concept), " ") + "），请确认"})
		}
		if existing != nil {
			// an unused library component is a suggestion; a component the project already relies on is a must
			sev, note := "error", ""
			if existing.Kind == "ui-lib" && existing.Usage == 0 {
				sev, note = "warn", "（组件库标准组件，项目还没用过，请确认）"
			}
			pc.Reinvented = append(pc.Reinvented, Finding{"reinvented", sev, a.Source, 0, a.Name, "新建了 " + a.Name + "，已有 " + existing.Name + "（" + kindLabel(existing.Kind) + " " + existing.Source + "）" + note})
		}
	}
	for _, f := range r.Findings {
		if (f.Rule == "missed-reuse" || f.Rule == "bypass") && in(f.File) && inLines(f.File, f.Line) {
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
		for _, re := range []*regexp.Regexp{cssVarDefRe, preVarDefRe, xmlColorDefRe} {
			for _, m := range re.FindAllStringSubmatchIndex(f.src, -1) {
				if !inLines(f.rel, lineAt(f.src, m[0])) {
					continue
				}
				pc.LocalTokens++
				name := f.src[m[2]:m[3]]
				v := strings.ToUpper(expandHex(f.src[m[4]:m[5]]))
				t := tokens[v]
				if t == name || t == "var("+name+")" || t == "@color/"+name {
					delete(tokens, v)
				} else if t != "" {
					ln := strings.Count(f.src[:m[0]], "\n") + 1
					pc.Colors = append(pc.Colors, Finding{"hardcoded-color", "warn", f.rel, ln, v, "新定义的颜色 " + name + " 与已有 token " + t + " 同值"})
				}
			}
		}
	}
	// hand-built concepts: class names / selectors that express a concept the page does not use a component for
	seenHB := map[string]bool{}
	for _, f := range focusSrc {
		for _, m := range classAttrRe.FindAllStringSubmatchIndex(f.src, -1) {
			if onComponent(f.src, m[0]) {
				continue // <van-tag class="status-tag">: the class styles a component, it does not replace one
			}
			if !inLines(f.rel, lineAt(f.src, m[0])) {
				continue
			}
			for _, cls := range strings.Fields(f.src[m[2]:m[3]]) {
				checkHandBuilt(pc, r, pool, usedConcepts, seenHB, f.rel, strings.Count(f.src[:m[0]], "\n")+1, cls, in)
			}
		}
		for _, m := range androidIDRe.FindAllStringSubmatchIndex(f.src, -1) {
			if onComponent(f.src, m[0]) {
				continue // <QMUIRoundButton android:id="@+id/status_tag">: a deliberate component choice
			}
			if !inLines(f.rel, lineAt(f.src, m[0])) {
				continue
			}
			checkHandBuilt(pc, r, pool, usedConcepts, seenHB, f.rel, strings.Count(f.src[:m[0]], "\n")+1, f.src[m[2]:m[3]], in)
		}
		for _, m := range cssSelRe2.FindAllStringSubmatchIndex(f.src, -1) {
			if !inLines(f.rel, lineAt(f.src, m[0])) {
				continue
			}
			checkHandBuilt(pc, r, pool, usedConcepts, seenHB, f.rel, strings.Count(f.src[:m[0]], "\n")+1, f.src[m[2]:m[3]], in)
		}
	}
	// inconsistent choices: the page brings in a component the project has not used before for a
	// concept it already has an established component for (Vant CountDown where vue-countdown is used 5×)
	for _, a := range r.Assets {
		if a.Concept == "" || a.Usage == 0 || a.WrappedBy != "" {
			continue // a library component the design system wraps is used inside its wrapper on purpose
		}
		inFocus, outside := false, false
		for _, f := range a.UsageFiles {
			if in(f) && !inDesignSystem(f) {
				inFocus = true
			} else {
				outside = true
			}
		}
		if !inFocus || outside || in(a.Source) {
			continue
		}
		var est *Asset
		for i := range r.Assets {
			b := &r.Assets[i]
			if b.Concept != a.Concept || b.Name == a.Name && b.Source == a.Source || b.Usage < 3 || featureOnly(*b) {
				continue
			}
			viewSide := strings.HasSuffix(a.Source, ".xml") || a.Kind == "custom-view" || strings.HasPrefix(a.Name, "@layout/")
			composeSide := a.Kind == "composable" || isComposeLib(a.Source)
			if viewSide && (b.Kind == "composable" || isComposeLib(b.Source)) || composeSide && (b.Kind == "custom-view" || strings.HasPrefix(b.Name, "@layout/")) {
				continue // XML and Compose alternatives are not interchangeable
			}
			if est == nil || preference(*b, a.Concept) > preference(*est, a.Concept) {
				est = b
			}
		}
		if est != nil {
			pc.Inconsistent = append(pc.Inconsistent, Finding{"inconsistent", "warn", a.UsageFiles[0], 0, a.Name, "新引入了 " + a.Name + "（" + kindLabel(a.Kind) + " " + a.Source + "），项目里「" + a.Concept + "」惯用 " + est.Name + "（用 " + itoa(est.Usage) + " 次），请确认是否统一"})
		}
	}
	// hand-built list rows by structure (no telling id / class): XML rows of label + value TextViews
	for _, f := range focusSrc {
		xmlLayout := strings.HasSuffix(f.rel, ".xml") && strings.Contains("/"+f.rel, "/res/layout")
		markupFile := strings.HasSuffix(f.rel, ".vue") || strings.HasSuffix(f.rel, ".tsx") || strings.HasSuffix(f.rel, ".jsx")
		if (xmlLayout || markupFile) && !usedConcepts["cell"] && !seenHB["cell"] {
			isRow, minRows := textRow, 2
			if markupFile {
				isRow, minRows = htmlTextRow, 3 // markup is freer: ask for a clearer repetition
			}
			var rows []int
			var walk func(ns []*mnode, depth int)
			walk = func(ns []*mnode, depth int) {
				for _, n := range ns {
					// a layout whose root is the row is a reusable item layout (info_cell.xml), not a hand-built one
					if depth > 0 && isRow(n) && inLines(f.rel, n.line) {
						rows = append(rows, n.line)
						continue
					}
					walk(n.children, depth+1)
				}
			}
			walk(parseMarkup(f.src), 0)
			if len(rows) >= minRows {
				pref := func(a *Asset) int { return preference(*a, "cell") }
				var best *Asset
				for i := range pool {
					b := &pool[i]
					if b.Concept == "cell" && !in(b.Source) && b.Kind != "composable" && !isComposeLib(b.Source) && b.Source != "Android SDK widgets" && !featureOnly(*b) && (best == nil || pref(b) > pref(best)) {
						best = b
					}
				}
				if best != nil && best.Usage > 0 {
					seenHB["cell"] = true
					var at []string
					for i, l := range rows {
						if i == 4 {
							at = append(at, "…")
							break
						}
						at = append(at, itoa(l))
					}
					pc.HandBuilt = append(pc.HandBuilt, Finding{"hand-built", "error", f.rel, rows[0], "list rows", "页面里有 " + itoa(len(rows)) + " 行「标题 + 值」横排是手写的（第 " + strings.Join(at, "、") + " 行），已有 " + best.Name + "（" + kindLabel(best.Kind) + " " + best.Source + "）（项目里已用 " + itoa(best.Usage) + " 次）"})
				}
			}
		}
	}
	for _, f := range focusSrc {
		src, rel := f.src, f.rel
		for _, loc := range hexColorRe.FindAllStringIndex(src, -1) {
			if lineHasVarDef(src, loc[0]) || !inLines(rel, lineAt(src, loc[0])) {
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
			if !inLines(rel, ln) {
				continue
			}
			pc.Colors = append(pc.Colors, Finding{"hardcoded-color", "warn", rel, ln, src[m[0]:m[1]], "写死颜色"})
		}
	}
	pc.Pass = true
	for _, f := range pc.Reinvented {
		if f.Severity == "error" {
			pc.Pass = false
		}
	}
	for _, f := range pc.MissedReuse {
		if f.Severity != "info" {
			pc.Pass = false
		}
	}
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
	// <color name="order_red">#F53F3F</color> in a values file the page adds
	xmlColorDefRe = regexp.MustCompile(`<color\s+name="(\w+)"\s*>\s*(#[0-9a-fA-F]{3,8})\s*<`)
	androidIDRe   = regexp.MustCompile(`android:id="@\+id/(\w+)"`)
)

var openTagRe = regexp.MustCompile(`<([A-Za-z][\w.:-]*)[^<>]*$`)

// onComponent reports whether the attribute at idx sits on a component tag rather than on
// plain markup (div/view, LinearLayout, TextView …).
func onComponent(src string, idx int) bool {
	start := idx - 4000
	if start < 0 {
		start = 0
	}
	m := openTagRe.FindStringSubmatch(src[start:idx])
	if m == nil {
		return false
	}
	tag := m[1]
	if strings.Contains(tag, "-") || tag[0] >= 'A' && tag[0] <= 'Z' && !strings.Contains(tag, ".") && !sdkPlain[tag] && !containers[tag] {
		return true
	}
	simple := tag[strings.LastIndex(tag, ".")+1:]
	return strings.Contains(tag, ".") && !containers[tag] && !containers[simple] && !sdkPlain[simple]
}

// sdkPlain are framework widgets that carry no UI concept by themselves.
var sdkPlain = map[string]bool{"TextView": true, "ImageView": true, "View": true, "LinearLayout": true, "FrameLayout": true, "RelativeLayout": true, "ConstraintLayout": true, "ScrollView": true, "NestedScrollView": true, "CardView": true, "Space": true}

func lineHasVarDef(src string, idx int) bool {
	ls := strings.LastIndex(src[:idx], "\n") + 1
	le := strings.IndexByte(src[idx:], '\n')
	if le < 0 {
		le = len(src) - idx
	}
	line := src[ls : idx+le]
	return cssVarDefRe.MatchString(line) || preVarDefRe.MatchString(line) || xmlColorDefRe.MatchString(line)
}

func checkHandBuilt(pc *PageCheck, r *Report, pool []Asset, used, seen map[string]bool, rel string, ln int, cls string, in func(string) bool) {
	concept := ConceptOf(cls)
	if concept == "" || used[concept] || seen[concept] || concept == "list" || concept == "icon" || concept == "image" {
		return
	}
	pref := func(a *Asset) int { return preference(*a, concept) }
	xmlCtx := strings.HasSuffix(rel, ".xml") || strings.HasSuffix(rel, ".java")
	var best *Asset
	for i := range pool {
		b := &pool[i]
		if xmlCtx && (b.Kind == "composable" || isComposeLib(b.Source)) {
			continue
		}
		if b.Kind == "ui-lib" && len(specialWords(b.Name, b.Concept)) > 0 && len(splitWords(b.Name)) > 1 {
			continue // FloatingActionButton is no replacement for a hand-built button
		}
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
	focus := pc.Focus
	if fs := strings.Split(focus, ","); len(fs) > 3 {
		focus = itoa(len(fs)) + " 个文件（" + strings.Join(fs[:2], "、") + " …）"
	}
	b.WriteString("\n## 页面复用检查：" + focus + "  " + verdict + "\n\n")
	b.WriteString("- 复用的现成组件（" + itoa(len(pc.Used)) + "）：" + strings.Join(pc.Used, "、") + "\n")
	b.WriteString("- 自造轮子（" + itoa(len(pc.Reinvented)) + "）\n")
	for _, f := range pc.Reinvented {
		mark := "✗"
		if f.Severity != "error" {
			mark = "·"
		}
		b.WriteString("  - " + mark + " " + f.File + "：" + f.Hint + "\n")
	}
	b.WriteString("- 有组件却没用（" + itoa(len(pc.MissedReuse)) + "）\n")
	for _, f := range pc.MissedReuse {
		mark := "✗"
		if f.Severity == "info" {
			mark = "·"
		}
		b.WriteString("  - " + mark + " " + f.File + ":" + itoa(f.Line) + " `" + f.Text + "` → " + f.Hint + "\n")
	}
	b.WriteString("- 疑似手写已有组件（" + itoa(len(pc.HandBuilt)) + "）\n")
	for _, f := range pc.HandBuilt {
		mark := "·"
		if f.Severity == "error" {
			mark = "✗"
		}
		b.WriteString("  - " + mark + " " + f.File + ":" + itoa(f.Line) + " " + f.Hint + "\n")
	}
	if len(pc.Inconsistent) > 0 {
		b.WriteString("- 和项目惯用组件不一致（" + itoa(len(pc.Inconsistent)) + "，不影响通过）\n")
		for _, f := range pc.Inconsistent {
			b.WriteString("  - · " + f.File + "：" + f.Hint + "\n")
		}
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

func subset(a, b []string) bool {
	for _, x := range a {
		if !contains(b, x) {
			return false
		}
	}
	return true
}

// accepted: the team already decided this component should exist — it is the confirmed
// implementation of a Figma component, or named in a gap decision of component-map.yaml.
func accepted(name, platform string) bool {
	if Confirmed == nil {
		return false
	}
	for _, e := range Confirmed.Components {
		if im := e.Impl(platform); im != nil && im.Component == name {
			return true
		}
	}
	for _, g := range Confirmed.Gaps {
		if (g.Platform == platform || g.Platform == "all") && strings.Contains(g.Decision, name) {
			return true
		}
	}
	return false
}

// textRow: a horizontal layout whose content is a label and a value (TextViews, possibly a
// vertical title/subtitle stack) and at most an arrow / spacer — a list cell built by hand.
func textRow(n *mnode) bool {
	tag := n.tag[strings.LastIndex(n.tag, ".")+1:]
	if tag != "LinearLayout" && tag != "RelativeLayout" && tag != "ConstraintLayout" {
		return false
	}
	if tag == "LinearLayout" && strings.Contains(n.attrs, `orientation="vertical"`) {
		return false
	}
	texts := 0
	for _, c := range n.children {
		switch t := c.tag[strings.LastIndex(c.tag, ".")+1:]; {
		case t == "TextView" || t == "AppCompatTextView" || t == "MaterialTextView":
			texts++
		case t == "LinearLayout" && strings.Contains(c.attrs, `orientation="vertical"`):
			for _, g := range c.children {
				if g.tag[strings.LastIndex(g.tag, ".")+1:] != "TextView" {
					return false
				}
			}
			texts++
		case t == "ImageView" || t == "AppCompatImageView" || t == "View" || t == "Space":
		default:
			return false
		}
	}
	return texts >= 2 && len(n.children) <= 4
}

// htmlTextRow: a div / li / view holding just a label and a value as plain text elements.
func htmlTextRow(n *mnode) bool {
	switch strings.ToLower(n.tag) {
	case "div", "li", "view", "section", "p":
	default:
		return false
	}
	if len(n.children) < 2 || len(n.children) > 3 {
		return false
	}
	for _, c := range n.children {
		switch strings.ToLower(c.tag) {
		case "span", "div", "p", "label", "text", "em", "b", "strong", "i":
		default:
			return false
		}
		if len(c.children) > 0 {
			return false
		}
	}
	return true
}
