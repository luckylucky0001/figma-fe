package reuse

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
	"github.com/luckylucky0001/figma-fe/cli/internal/mapping"
)

// Match answers "a new Figma page arrived — what can we reuse from what we already have?" by
// matching every Figma component (instance) against the project's inventoried assets.

type Candidate struct {
	Name      string   `json:"name"`
	Kind      string   `json:"kind"`
	Source    string   `json:"source"`
	Score     float64  `json:"score"`
	Usage     int      `json:"usage"`
	Example   string   `json:"example,omitempty"`
	ExampleAt string   `json:"exampleAt,omitempty"`
	Why       string   `json:"why"`
	Props     []string `json:"props,omitempty"` // how to configure it for the Figma variants
}

type MatchItem struct {
	Figma      string      `json:"figma"`
	Instances  int         `json:"instances"`
	NodeIDs    []string    `json:"nodeIds"`
	Variants   []string    `json:"variants,omitempty"`
	Verdict    string      `json:"verdict"`            // reuse | confirm | new
	Decision   string      `json:"decision,omitempty"` // confirmed gap decision
	Candidates []Candidate `json:"candidates"`
}

type MatchResult struct {
	Platform string      `json:"platform"`
	Project  string      `json:"project"`
	Page     string      `json:"page"`
	Items    []MatchItem `json:"items"`
	Reuse    int         `json:"reuse"`
	Confirm  int         `json:"confirm"`
	New      int         `json:"new"`
}

// Target is the technology of the page being built: "" (any), "xml" (Android View layouts) or
// "compose". Assets that cannot be used from that technology are pushed down.
var Target string

// InferTarget guesses the page technology from what the project's screens are mostly built with:
// component usages from XML layouts / Views vs composables. "" when it is too close to call.
func InferTarget(r *Report) string {
	var views, compose int
	for _, a := range r.Assets {
		switch {
		case a.Kind == "composable" || isComposeLib(a.Source):
			compose += a.Usage
		case a.Kind == "custom-view" || a.Kind == "view-helper" || a.Kind == "ui-lib":
			views += a.Usage
		}
	}
	switch {
	case views > 2*compose:
		return "xml"
	case compose > 2*views:
		return "compose"
	}
	return ""
}

// Confirmed is a project's confirmed component map (design/component-map.yaml or
// .figmafe/component-map.yaml). Confirmed entries win over inference: the team decided once.
var Confirmed *mapping.Map

// Match scores project assets for each Figma component set used on the page.
func Match(exp *figma.Export, r *Report) *MatchResult {
	res := &MatchResult{Platform: r.Platform, Project: r.Root, Page: exp.Frame.Name}
	type agg struct {
		ids      []string
		variants map[string]bool
		layers   map[string]bool
	}
	sets := map[string]*agg{}
	var order []string
	figma.Walk(exp.Frame, func(n *figma.Node) bool {
		if n.Type == "INSTANCE" {
			name := exp.SetName(n)
			// slash-named components without a variant set ("Button/Primary/Large"): the first segment
			// is the component, the rest are its variants
			var pathVariant string
			if i := strings.Index(name, "/"); i > 0 && len(n.Variant()) == 0 {
				pathVariant = strings.TrimSpace(name[i+1:])
				name = strings.TrimSpace(name[:i])
			}
			a := sets[name]
			if a == nil {
				a = &agg{variants: map[string]bool{}, layers: map[string]bool{}}
				sets[name] = a
				order = append(order, name)
			}
			a.ids = append(a.ids, n.ID)
			for k, v := range n.Variant() {
				a.variants[k+"="+v] = true
			}
			if pathVariant != "" {
				a.variants["Variant="+pathVariant] = true
			}
			for _, c := range n.Children {
				if c.Type == "TEXT" {
					a.layers[strings.ToLower(c.Name)] = true
				}
			}
		}
		return true
	})
	// Detached / never-componentized layers whose names express a UI concept ("按钮", "Tag 状态")
	// are reuse candidates too: designers often draw them by hand.
	figma.Walk(exp.Frame, func(n *figma.Node) bool {
		if n.Type == "INSTANCE" {
			return false
		}
		if (n.Type == "FRAME" || n.Type == "GROUP" || n.Type == "COMPONENT") && n != exp.Frame && len(n.Children) > 0 && ConceptOf(n.Name) != "" {
			name := n.Name + "（未组件化）"
			a := sets[name]
			if a == nil {
				a = &agg{variants: map[string]bool{}, layers: map[string]bool{}}
				sets[name] = a
				order = append(order, name)
			}
			a.ids = append(a.ids, n.ID)
		}
		return true // keep walking: a list frame contains item / tag / button frames
	})
	pool := append([]Asset{}, r.Assets...)
	seen := map[string]bool{}
	for _, a := range r.Assets {
		seen[a.Kind+"|"+a.Source+"|"+a.Name] = true
	}
	for _, a := range r.LibAll {
		if !seen[a.Kind+"|"+a.Source+"|"+a.Name] && a.Source != "Android SDK widgets" {
			pool = append(pool, a)
		}
	}
	setFamilyPrefixes(pool)
	established = map[string]string{}
	best := map[string]int{}
	for _, a := range r.Assets {
		if a.Kind != "ui-lib" && a.Kind != "page-local" && a.Concept != "" && a.Usage >= 3 && a.Usage > best[a.Concept] {
			best[a.Concept], established[a.Concept] = a.Usage, a.Name
		}
	}
	libShare = map[string]float64{}
	total := 0
	for _, l := range r.Libraries {
		if !builtinLib(l.Name) { // platform primitives (<view>, TextView) would drown every real toolkit
			total += l.Usage
		}
	}
	for _, l := range r.Libraries {
		if total > 0 {
			libShare[l.Name] = float64(l.Usage) / float64(total)
		}
	}
	familyBase = familyBases(pool)
	libChild = map[string]bool{}
	names := map[string]bool{}
	for _, a := range pool {
		names[a.Source+"|"+a.Name] = true
	}
	for _, a := range pool {
		for _, suf := range []string{"s", "Group", "Row"} {
			if a.Concept != "" && a.Concept != "button" && names[a.Source+"|"+a.Name+suf] {
				libChild[a.Source+"|"+a.Name] = true // Tab in Tabs, Radio in RadioGroup, NiaTab in NiaTabRow
			}
		}
	}
	for _, name := range order {
		a := sets[name]
		item := MatchItem{Figma: name, Instances: len(a.ids), NodeIDs: a.ids}
		for v := range a.variants {
			item.Variants = append(item.Variants, v)
		}
		sort.Strings(item.Variants)
		if Confirmed != nil {
			if e := Confirmed.Entry(name); e != nil {
				if im := e.Impl(r.Platform); im != nil {
					c := Candidate{Name: im.Component, Kind: "confirmed", Source: im.Import, Score: 1, Example: im.Usage, Why: "已确认映射（component-map.yaml）"}
					for _, as := range r.Assets {
						if as.Name == im.Component {
							c.Kind, c.Usage = as.Kind, as.Usage
							if c.Example == "" {
								c.Example, c.ExampleAt = as.Example, as.ExampleAt
							}
						}
					}
					item.Candidates = []Candidate{c}
					item.Verdict = "reuse"
					res.Reuse++
					res.Items = append(res.Items, item)
					continue
				}
			}
			if g := Confirmed.GapFor(name, r.Platform); g != nil {
				item.Verdict = "new"
				item.Decision = g.Decision
				item.Candidates = nil
				res.New++
				res.Items = append(res.Items, item)
				continue
			}
		}
		fc := ConceptOf(strings.TrimSuffix(name, "（未组件化）"))
		for _, as := range pool {
			if as.Unavailable {
				continue
			}
			sc, why := score(name, fc, a.variants, as)
			if sc > 0 && !usable(&as) {
				sc -= 0.3
				why += "；⚠ " + as.Caveat
			}
			if sc < 0.5 {
				continue
			}
			ex, at := as.Example, as.ExampleAt
			if ex == "" && as.Kind == "ui-lib" {
				ex, at = libSnippet(as), "（"+as.Source+" 标准用法，项目未用过）"
			}
			hints := PropHints(item.Variants, as)
			if st, n := topStyle(as.Styles); n >= 3 {
				hints = append(hints, fmt.Sprintf(`项目惯用 style="@style/%s"（%d 处）`, st, n))
			}
			if habits := habitList(as.Attrs, as.Usage); len(habits) > 0 {
				hints = append(hints, "项目惯用写法："+strings.Join(habits, " "))
			}
			item.Candidates = append(item.Candidates, Candidate{as.Name, as.Kind, as.Source, round2(sc), as.Usage, ex, at, why, hints})
		}
		sort.SliceStable(item.Candidates, func(i, j int) bool {
			ci, cj := item.Candidates[i], item.Candidates[j]
			if ci.Score != cj.Score {
				return ci.Score > cj.Score
			}
			if ci.Usage != cj.Usage {
				return ci.Usage > cj.Usage // tie: what this project already uses keeps it consistent
			}
			return len(ci.Name) < len(cj.Name) // generic before specialised (List before AddressList)
		})
		// no component fits well: the project's established way of drawing it (1px View, hairline class)
		if len(item.Candidates) == 0 || item.Candidates[0].Score < 0.7 {
			for _, cv := range r.Conventions {
				if cv.Concept == fc && fc != "" {
					sc := 0.72 + math.Min(0.2, 0.05*math.Log2(float64(cv.Usage)+1))
					item.Candidates = append([]Candidate{{cv.Name, cv.Kind, cv.Source, round2(sc), cv.Usage, cv.Example, cv.ExampleAt, "没有对应组件，项目里惯用这种写法（" + itoa(cv.Usage) + " 处）", nil}}, item.Candidates...)
					break
				}
			}
		}
		// still nothing: a generic container the project uses that can draw it through attributes
		// (QMUILinearLayout's qmui_bottomDividerHeight) — a capability, so it needs confirmation
		if fc == "divider" && (len(item.Candidates) == 0 || item.Candidates[0].Score < 0.7) {
			var best *Asset
			var attrs []string
			for i := range pool {
				as := &pool[i]
				// among established containers (≥3 uses) the plainest one: QMUIFrameLayout over QMUIPriorityLinearLayout
				if as.Concept != "" || as.Usage == 0 || best != nil && !(as.Usage >= 3 && (best.Usage < 3 || len(as.Name) < len(best.Name)) || best.Usage < 3 && as.Usage > best.Usage) {
					continue
				}
				var da []string
				for _, p := range as.Props {
					if strings.Contains(strings.ToLower(p), "divider") {
						da = append(da, p)
					}
				}
				if len(da) > 0 {
					best, attrs = as, da
				}
			}
			if best != nil {
				if len(attrs) > 2 {
					attrs = attrs[:2]
				}
				item.Candidates = append([]Candidate{{best.Name, best.Kind, best.Source, 0.75, best.Usage, best.Example, best.ExampleAt, "没有分割线组件；这个容器自带分割线属性（" + strings.Join(attrs, "、") + " …），项目里已用 " + itoa(best.Usage) + " 次", nil}}, item.Candidates...)
			}
		}
		if len(item.Candidates) > 4 {
			item.Candidates = item.Candidates[:4]
		}
		switch {
		case len(item.Candidates) > 0 && item.Candidates[0].Score >= 0.8:
			item.Verdict = "reuse"
			res.Reuse++
		case len(item.Candidates) > 0:
			item.Verdict = "confirm"
			res.Confirm++
		default:
			item.Verdict = "new"
			res.New++
		}
		res.Items = append(res.Items, item)
	}
	return res
}

func score(figmaName, figmaConcept string, variants map[string]bool, a Asset) (float64, string) {
	figmaName = strings.TrimSuffix(figmaName, "（未组件化）")
	fn, an := normName(figmaName), normName(a.Name)
	if hasCJK(figmaName) {
		fn = "" // Chinese layer names match on concept only
	}
	var s float64
	var why []string
	switch {
	case fn != "" && fn == an && figmaConcept != "" && a.Concept != "" && figmaConcept != a.Concept:
		s, why = 0.6, append(why, "同名但用途不同（"+a.Concept+"），需确认")
	case fn != "" && fn == an:
		s, why = 0.85, append(why, "同名")
	case figmaConcept != "" && figmaConcept == a.Concept:
		s, why = 0.72, append(why, "同为 "+figmaConcept)
	case strings.Contains(an, fn) && len(fn) >= 4 && (figmaConcept == "" || a.Concept == "" || a.Concept == figmaConcept) || strings.Contains(fn, an) && len(an) >= 4:
		// MaterialRadioButton contains "button" but is a radio: a more specific name of another concept
		// is not a match (the other way round, ProductCard → Card, is)
		s, why = 0.6, append(why, "名称包含")
	default:
		return 0, ""
	}
	// specialised components (AddressList, ContactCard, CouponCell) carry domain words the design does not
	if s < 0.9 && a.Concept != "" && a.Kind != "convention" {
		residual := specialWords(a.Name, a.Concept)
		joined := strings.Join(residual, "")
		lowerFigma := strings.ToLower(figmaName)
		specialised := len(joined) >= 3 && !strings.Contains(lowerFigma, joined)
		if a.Kind != "ui-lib" {
			// a project's StatusTagView / InterestsItem is its own vocabulary; only a multi-word domain
			// qualifier (NewsResourceCardExpanded) marks a component made for one feature
			specialised = specialised && len(residual) >= 2
		}
		if specialised {
			s -= 0.12
			why = append(why, "专用组件（"+joined+"）")
		}
	}
	// <include layout="@layout/activity_article_float_btn"> is a slice of one screen, not a component
	if strings.HasPrefix(a.Name, "@layout/activity_") || strings.HasPrefix(a.Name, "@layout/fragment_") {
		s -= 0.15
		why = append(why, "页面布局片段")
	}
	if libChild[a.Source+"|"+a.Name] && fn != an {
		s -= 0.1
		why = append(why, "是容器组件（"+a.Name+"s / Group / Row）的子项")
	}
	// a feature component: several required business parameters (InterestsItem needs following,
	// topicImageUrl, onFollowButtonClick); design-system components are exempt by definition
	if dr := domainRequired(a.Required); len(dr) >= 3 && !inDesignSystem(a.Source) || hasModelParam(a.Required) {
		s -= 0.12
		why = append(why, "业务专用（必填 "+strings.Join(dr, "、")+"）")
	}
	// a composable / view only ever used inside the file that defines it is a page's helper
	if (a.Kind == "composable" || a.Kind == "custom-view" || a.Kind == "view-helper") && len(a.UsageFiles) > 0 && !inDesignSystem(a.Source) {
		local := true
		for _, f := range a.UsageFiles {
			local = local && f == a.Source
		}
		if local {
			s -= 0.15
			why = append(why, "只在定义它的文件里用过（页面内部组件）")
		}
	}
	if _, primitive := rawAndroid[a.Name]; primitive && a.Source == "Android SDK widgets" {
		s -= 0.15 // Button / EditText / Toolbar are what projects wrap: the project's own component wins
		why = append(why, "系统原生控件")
	}
	// the project already has its own way for this concept (mix-list-cell × 6): an unused library
	// component would be the inconsistency the page check flags
	if e := established[a.Concept]; e != "" && a.Kind == "ui-lib" && a.Usage == 0 {
		s -= 0.08
		why = append(why, "项目里已惯用 "+e)
	}
	// a library the project barely uses (one stray <van-button> in a uni-ui project) is not its toolkit
	if share, ok := libShare[a.Source]; ok && a.Kind == "ui-lib" && a.Usage == 0 && share < 0.05 {
		s -= 0.15
		why = append(why, "项目几乎不用 "+a.Source)
	}
	// made for one job the design does not ask for (PasswordInputView, ConfirmButton, VideoCard)
	if a.Kind != "ui-lib" {
		for _, w := range specialWords(a.Name, a.Concept) {
			if purposeWords[w] && !strings.Contains(strings.ToLower(figmaName), w) {
				s -= 0.1
				why = append(why, "专门用途（"+w+"）")
				break
			}
		}
	}
		// the name links to the concept only through a weak word (swipeactions_row, InterestsItem): less sure
	if a.Kind != "ui-lib" && a.Concept != "" && weakOnly(a.Name, a.Concept) && fn != an {
		s -= 0.08
		why = append(why, "名称只以 row / item 之类的泛称对上")
	}
	// a selectable component (FilterChip: selected / onSelectedChange) for a design element that has
	// no selected state is the wrong interaction
	if (selectable(a.Props) || selectableName(a.Name)) && !hasAxis(variants, "selected", "checked", "active", "state", "on") {
		s -= 0.06
		why = append(why, "可选中组件，设计稿里是静态展示")
	}
	// the library's canonical component of the concept (Cell for cell, Tag for tag) beats look-alikes
	for _, c := range concepts {
		if c.Name == a.Concept && (an == c.Words[0] || familyBase[a.Name]) && fn != an {
			s += 0.04
			why = append(why, "标准组件")
		}
	}
	// variant axes ↔ props
	if len(variants) > 0 && len(a.Props) > 0 {
		hit := 0
		axes := map[string]bool{}
		for v := range variants {
			axes[strings.ToLower(strings.SplitN(v, "=", 2)[0])] = true
		}
		for ax := range axes {
			for _, p := range a.Props {
				if strings.Contains(strings.ToLower(p), ax) || ax == "type" && p == "variant" || ax == "tone" && (p == "type" || p == "color" || p == "theme") {
					hit++
					break
				}
			}
		}
		if hit > 0 {
			s += 0.05 * float64(hit)
			why = append(why, fmt.Sprintf("变体轴匹配 %d/%d", hit, len(axes)))
		}
	}
	// proven in this project / preferred kinds
	if a.Usage > 0 {
		s += math.Min(0.25, 0.05*math.Log2(float64(a.Usage)+1)) // proven in this codebase, more so when widely used
		why = append(why, fmt.Sprintf("项目里已用 %d 次", a.Usage))
	}
	if a.WrappedBy != "" {
		s -= 0.15
		why = append(why, "项目设计系统已封装为 "+a.WrappedBy+"，应使用封装")
	}
	switch {
	case a.Kind == "design-system" || (a.Kind == "composable" || a.Kind == "custom-view") && inDesignSystem(a.Source):
		s += 0.05
		why = append(why, "设计系统组件")
	case a.Kind == "page-local":
		s -= 0.2
		why = append(why, "页面私有，需先抽出")
	}
	compose := a.Kind == "composable" || isComposeLib(a.Source)
	switch {
	case Target == "xml" && compose:
		s = math.Min(s-0.25, 0.79) // never a direct reuse: an XML page cannot host it
		why = append(why, "Composable，XML 页面不能直接用")
	case Target == "compose" && (a.Kind == "custom-view" || a.Kind == "view-helper" || a.Kind == "ui-lib" && !compose && r2android(a)):
		s = math.Min(s-0.2, 0.79)
		why = append(why, "View 组件，Compose 页面需 AndroidView 包装")
	case Target == "compose" && a.Kind == "composable":
		s += 0.03
	}
	if strings.Contains(strings.ToLower(a.Name), "legacy") || strings.Contains(strings.ToLower(a.Name), "deprecated") || strings.HasPrefix(a.Name, "Old") {
		s -= 0.1
		why = append(why, "名称表明是老组件")
	}
	if s > 1 {
		s = 1
	}
	return s, strings.Join(why, "，")
}

func normName(s string) string {
	w := splitWords(s)
	if len(w) > 1 && famPrefixes[w[0]] {
		w = w[1:] // JetsnackCard is the project's Card
	}
	j := strings.Join(w, "")
	for _, p := range []string{"van", "nut", "ds", "legacy", "qmui", "material", "app"} {
		j = strings.TrimPrefix(j, p)
	}
	if t := strings.TrimSuffix(strings.TrimSuffix(j, "view"), "component"); t != "" {
		return t
	}
	return j // "View" itself stays "view"
}

func round2(f float64) float64 { return float64(int(f*100+0.5)) / 100 }

// Markdown renders the reuse list for a new Figma page.
func (m *MatchResult) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# 新设计稿复用清单：%s → %s（%s）\n\n", m.Page, m.Project, m.Platform)
	fmt.Fprintf(&b, "Figma 组件 %d 个：**可直接复用 %d** · 需确认 %d · 需新建 %d\n\n", len(m.Items), m.Reuse, m.Confirm, m.New)
	b.WriteString("| Figma 组件 | 实例 | 结论 | 推荐复用 | 依据 | 项目里的真实用法 |\n|---|---|---|---|---|---|\n")
	for _, it := range m.Items {
		verdict := map[string]string{"reuse": "✅ 直接复用", "confirm": "❓ 需确认", "new": "🆕 需新建"}[it.Verdict]
		if len(it.Candidates) == 0 {
			why := "项目和组件库里都没有"
			if it.Decision != "" {
				why = "已确认的缺口决定：" + it.Decision
			}
			fmt.Fprintf(&b, "| %s | %d | %s | — | %s | — |\n", it.Figma, it.Instances, verdict, why)
			continue
		}
		c := it.Candidates[0]
		ex := "—"
		if c.Example != "" {
			ex = "`" + strings.ReplaceAll(c.Example, "|", "\\|") + "` " + c.ExampleAt
		}
		fmt.Fprintf(&b, "| %s | %d | %s | `%s`（%s %s） | %.2f：%s | %s |\n", it.Figma, it.Instances, verdict, c.Name, kindLabel(c.Kind), c.Source, c.Score, c.Why, ex)
	}
	var cfg []string
	for _, it := range m.Items {
		if it.Verdict != "new" && len(it.Candidates) > 0 && len(it.Candidates[0].Props) > 0 {
			cfg = append(cfg, "- **"+it.Figma+"** → `"+it.Candidates[0].Name+"`：\n  - "+strings.Join(it.Candidates[0].Props, "\n  - "))
		}
	}
	if len(cfg) > 0 {
		b.WriteString("\n## 变体 → 属性（推荐组件怎么配置）\n\n" + strings.Join(cfg, "\n") + "\n")
	}
	b.WriteString("\n## 候选明细\n\n")
	for _, it := range m.Items {
		if len(it.Candidates) < 2 {
			continue
		}
		var cs []string
		for _, c := range it.Candidates {
			cs = append(cs, fmt.Sprintf("`%s`（%s，%.2f，用 %d 次）", c.Name, kindLabel(c.Kind), c.Score, c.Usage))
		}
		fmt.Fprintf(&b, "- **%s**：%s\n", it.Figma, strings.Join(cs, "、"))
	}
	return b.String()
}

// libSnippet renders the canonical tag / call for a library component the project has not used yet.
func libSnippet(a Asset) string {
	for _, l := range uiLibs {
		if l.Name != a.Source {
			continue
		}
		switch {
		case l.TagPrefix != "":
			return "<" + l.TagPrefix + strings.Join(splitWords(a.Name), "-") + " … />"
		case len(l.Import) == 0: // built-in tags (uni-app <switch>, <checkbox>) need no import
			return "<" + strings.ToLower(strings.Join(splitWords(a.Name), "-")) + " … />"
		case l.Platform == "h5" && compoundOf(l, a.Name) != "":
			p := compoundOf(l, a.Name)
			return "import { " + p + " } from '" + l.Import[0] + "'; <" + p + "." + strings.TrimPrefix(a.Name, p) + " … />"
		case l.Platform == "h5":
			return "import { " + a.Name + " } from '" + l.Import[0] + "'; <" + a.Name + " … />"
		case strings.HasPrefix(l.Import[0], "androidx.compose"):
			return a.Name + "(…)  // import " + l.Import[0] + a.Name
		default:
			return "<" + l.Import[0] + "…" + a.Name + " … />"
		}
	}
	return ""
}

// Draft turns the match result into a component-map draft for this project (figma → component),
// so onboarding a brownfield project starts from a generated mapping instead of a blank file.
func (m *MatchResult) Draft() *mapping.Map {
	d := &mapping.Map{Version: 1}
	for _, it := range m.Items {
		if it.Verdict == "new" || len(it.Candidates) == 0 {
			d.Gaps = append(d.Gaps, mapping.Gap{Figma: it.Figma, Platform: m.Platform, Decision: "TODO：项目和组件库里都没有，确认后新建组件（建议放进组件库而不是页面）"})
			continue
		}
		c := it.Candidates[0]
		im := &mapping.Impl{Component: c.Name, Usage: c.Example}
		im.Import = importHint(c)
		var alts []string
		for _, x := range it.Candidates[1:] {
			alts = append(alts, x.Name)
		}
		im.Notes = fmt.Sprintf("[%s %.2f] %s；来源 %s", it.Verdict, c.Score, c.Why, c.Source)
		if len(alts) > 0 {
			im.Notes += "；备选 " + strings.Join(alts, " / ")
		}
		e := mapping.Entry{Figma: it.Figma}
		if m.Platform == "android" {
			e.Android = im
		} else {
			e.H5 = im
		}
		if len(c.Props) > 0 {
			e.Constraints = append(e.Constraints, "变体 → 属性建议（确认后写进 props）："+strings.Join(c.Props, "；"))
		} else if len(it.Variants) > 0 {
			e.Constraints = append(e.Constraints, "TODO：补齐变体 → props 映射："+strings.Join(it.Variants, "、"))
		}
		d.Components = append(d.Components, e)
	}
	return d
}

func importHint(c Candidate) string {
	if c.Kind != "ui-lib" {
		return c.Source // defining file
	}
	for _, l := range uiLibs {
		if l.Name == c.Source && len(l.Import) > 0 {
			if l.Platform == "h5" {
				return "import { " + c.Name + " } from '" + l.Import[0] + "';"
			}
			return l.Import[0] + "*." + c.Name
		}
	}
	return c.Source
}

func r2android(a Asset) bool {
	for _, l := range uiLibs {
		if l.Name == a.Source {
			return l.Platform == "android"
		}
	}
	return false
}

// styleWords are orientation / visual-style words, not a specialisation (HorizontalDivider,
// OutlinedButton, FilledTonalButton are still the generic component).
// purposeWords: a component named for one job (login, video, confirm) is not the generic one.
var purposeWords = map[string]bool{"password": true, "confirm": true, "dismiss": true, "cancel": true, "delete": true, "login": true, "logout": true,
	"share": true, "download": true, "upload": true, "video": true, "audio": true, "music": true, "play": true, "pause": true, "retry": true, "float": true, "link": true, "back": true, "verify": true, "captcha": true}

var styleWords = []string{"horizontal", "vertical", "outlined", "filled", "elevated", "tonal", "default", "basic", "simple", "standard", "normal",
	"layout", "view", "widget", "component", "container", "wrapper", "round", "common", "custom", "base", "text", "priority"}

// familyPrefix: a word that prefixes a project family of components (Nia in NiaButton, NiaTab …),
// set from the pool in Match so it is not mistaken for a domain word.
func familyPrefix(w string) bool { return famPrefixes[w] }

var famPrefixes = map[string]bool{}

var libChild = map[string]bool{}

var libShare = map[string]float64{}

// established: concept → the project's own component used ≥3 times for it (set per Match)
var established = map[string]string{}

// builtinLib: the platform's own element set (Android SDK widgets, uni-app built-in tags), not a toolkit.
func builtinLib(name string) bool {
	if name == "Android SDK widgets" {
		return true
	}
	for _, l := range uiLibs {
		if l.Name == name {
			return len(l.Import) == 0
		}
	}
	return false
}

// familyBase marks a project family's base component: NiaButton among NiaTextButton /
// NiaOutlinedButton, DsCell among DsCellGroup … (prefix + the concept's own word).
var familyBase map[string]bool

func familyBases(pool []Asset) map[string]bool {
	out := map[string]bool{}
	for _, a := range pool {
		if a.Concept == "" {
			continue
		}
		var word string
		for _, c := range concepts {
			if c.Name == a.Concept {
				word = c.Words[0]
			}
		}
		n := normName(a.Name)
		if !strings.HasSuffix(n, word) || n == word {
			continue
		}
		prefix := strings.TrimSuffix(n, word)
		kin := 0
		for _, b := range pool {
			if bn := normName(b.Name); b.Concept == a.Concept && bn != n && strings.HasPrefix(bn, prefix) {
				kin++
			}
		}
		if kin >= 2 {
			out[a.Name] = true
		}
	}
	// a library's only / shortest component of a concept is its canonical one (antd-mobile ListItem for cell)
	short := map[string]string{}
	for _, a := range pool {
		if a.Kind != "ui-lib" || a.Concept == "" {
			continue
		}
		k := a.Source + "|" + a.Concept
		if cur, ok := short[k]; !ok || len(a.Name) < len(cur) {
			short[k] = a.Name
		}
	}
	for _, n := range short {
		out[n] = true
	}
	return out
}

// compoundOf: ListItem → List when the library exposes it as <List.Item> (antd-mobile style).
func compoundOf(l uiLib, name string) string {
	for _, suf := range []string{"Item", "Panel"} {
		if p := strings.TrimSuffix(name, suf); p != name && l.TagPrefix == "" && l.Name == "antd-mobile" {
			for _, c := range l.Components {
				if c == p {
					return p
				}
			}
		}
	}
	return ""
}

// specialWords: the words of a component name beyond its concept, vendor / family prefix and
// neutral style words — StatusTag → [status], NewsResourceCardExpanded → [news resource expanded].
func specialWords(name, concept string) []string {
	ws := splitWords(name)
	for len(ws) > 1 && (vendorPrefix[ws[0]] || projectPrefix(ws) || familyPrefix(ws[0])) {
		ws = ws[1:]
	}
	var conceptWords []string
	for _, c := range concepts {
		if c.Name == concept {
			conceptWords = c.Words
		}
	}
	// drop the concept's own words (longest n-gram first: top+app+bar) and neutral style words
	var residual []string
	for i := 0; i < len(ws); {
		n := 0
		for k := min(3, len(ws)-i); k >= 1 && n == 0; k-- {
			g := strings.Join(ws[i:i+k], "")
			if contains(conceptWords, g) || contains(conceptWords, strings.TrimSuffix(g, "s")) || k == 1 && contains(styleWords, g) {
				n = k
			}
		}
		if n == 0 {
			residual = append(residual, ws[i])
			n = 1
		}
		i += n
	}
	return residual
}

// setFamilyPrefixes learns project family prefixes (Nia, Ds, QD …) from the asset pool.
func setFamilyPrefixes(pool []Asset) {
	famPrefixes = map[string]bool{}
	firstWords := map[string]int{}
	for _, a := range pool {
		if ws := splitWords(a.Name); len(ws) > 1 && a.Kind != "ui-lib" {
			firstWords[ws[0]]++
		}
	}
	for w, n := range firstWords {
		if n >= 4 && ConceptOf(w) == "" && !purposeWords[w] { // Video* in Seal is a domain, not a family
			famPrefixes[w] = true
		}
	}
}

func selectable(props []string) bool {
	for _, p := range props {
		switch normProp(p) {
		case "selected", "checked", "isselected", "onselectedchange", "oncheckedchange", "onselect":
			return true
		}
	}
	return false
}

func hasAxis(variants map[string]bool, names ...string) bool {
	for v := range variants {
		axis := strings.ToLower(strings.SplitN(v, "=", 2)[0])
		for _, n := range names {
			if axis == n {
				return true
			}
		}
	}
	return false
}

var genericParams = map[string]bool{"modifier": true, "onclick": true, "content": true, "text": true, "title": true, "label": true,
	"icon": true, "enabled": true, "name": true, "value": true, "description": true, "subtitle": true, "image": true, "imageurl": true}

func domainRequired(req []string) []string {
	var out []string
	for _, p := range req {
		if !genericParams[normProp(p)] && !strings.Contains(p, ":") {
			out = append(out, p)
		}
	}
	return out
}

// hasModelParam: a required parameter typed as an app data model ("snack:Snack").
func hasModelParam(req []string) bool {
	for _, p := range req {
		if strings.Contains(p, ":") {
			return true
		}
	}
	return false
}

// weakOnly: the name reaches the concept only through an ambiguous structural word (row, item,
// panel …); "header" for a top bar is a fair name and does not count.
func weakOnly(name, concept string) bool {
	ws := splitWords(name)
	ambiguous := false
	for _, w := range ws {
		switch w {
		case "row", "rows", "item", "items", "panel", "box", "line":
			ambiguous = true
		}
	}
	if !ambiguous {
		return false
	}
	for _, c := range concepts {
		if c.Name != concept {
			continue
		}
		for i := range ws {
			for k := 1; k <= 3 && i+k <= len(ws); k++ {
				g := strings.Join(ws[i:i+k], "")
				for _, w := range c.Words {
					if (g == w || strings.TrimSuffix(g, "s") == w) && !weak[w] {
						return false
					}
				}
			}
		}
	}
	return true
}

// selectableName: FilterChip / ChoiceChip / ToggleTag … are interactive by name.
func selectableName(name string) bool {
	for _, w := range splitWords(name) {
		switch w {
		case "filter", "choice", "toggle", "selectable", "select":
			return true
		}
	}
	return false
}

// FeatureSummary rolls the reuse lists of a feature's pages into one plan, grouped by what to use:
// every component the feature needs, which Figma elements map to it, on how many pages.
func FeatureSummary(ms []*MatchResult) string {
	type row struct {
		use     string
		figma   []string
		pages   map[string]bool
		inst    int
		verdict string
	}
	rows := map[string]*row{}
	var order []string
	rank := map[string]int{"reuse": 0, "confirm": 1, "new": 2}
	for _, m := range ms {
		for _, it := range m.Items {
			name := strings.TrimSuffix(it.Figma, "（未组件化）")
			use, key := "需新建："+name, "new|"+strings.ToLower(name)
			if len(it.Candidates) > 0 && it.Verdict != "new" {
				use, key = it.Candidates[0].Name, "use|"+it.Candidates[0].Source+"|"+it.Candidates[0].Name
			}
			r := rows[key]
			if r == nil {
				r = &row{use: use, pages: map[string]bool{}, verdict: it.Verdict}
				rows[key] = r
				order = append(order, key)
			}
			if !contains(r.figma, name) {
				r.figma = append(r.figma, name)
			}
			r.pages[m.Page] = true
			r.inst += it.Instances
			if rank[it.Verdict] < rank[r.verdict] {
				r.verdict = it.Verdict // the same component confidently matched on one page is reusable everywhere
			}
		}
	}
	sort.SliceStable(order, func(i, j int) bool {
		a, b := rows[order[i]], rows[order[j]]
		if rank[a.verdict] != rank[b.verdict] {
			return rank[a.verdict] < rank[b.verdict]
		}
		return a.inst > b.inst
	})
	label := map[string]string{"reuse": "✅ 直接复用", "confirm": "❓ 需确认", "new": "🆕 需新建"}
	var b strings.Builder
	var pages []string
	for _, m := range ms {
		pages = append(pages, m.Page)
	}
	fmt.Fprintf(&b, "# 整个需求的复用清单（%d 个页面：%s）\n\n| 用什么 | 结论 | 对应的 Figma 元素 | 页面数 | 实例 |\n|---|---|---|---|---|\n", len(ms), strings.Join(pages, "、"))
	n := map[string]int{}
	for _, k := range order {
		r := rows[k]
		n[r.verdict]++
		fmt.Fprintf(&b, "| `%s` | %s | %s | %d | %d |\n", r.use, label[r.verdict], strings.Join(r.figma, "、"), len(r.pages), r.inst)
	}
	fmt.Fprintf(&b, "\n这个需求共要用到 %d 个组件：直接复用 %d · 需确认 %d · 需新建 %d。需新建的组件只做一次、放进组件库，所有页面共用。\n", len(order), n["reuse"], n["confirm"], n["new"])
	return b.String()
}

// topStyle: the style the project applies most often with a component.
func topStyle(m map[string]int) (string, int) {
	best, n := "", 0
	for k, v := range m {
		if v > n || v == n && k < best {
			best, n = k, v
		}
	}
	return best, n
}

// habitList: props passed in at least 3 usages and 40% of them — the project's way of using it.
func habitList(attrs map[string]int, usage int) []string {
	best := map[string]string{} // prop → its most used value (size="small" vs size="large")
	for k, n := range attrs {
		if n < 3 || float64(n) < 0.4*float64(usage) {
			continue
		}
		prop := strings.SplitN(k, "=", 2)[0]
		if cur, ok := best[prop]; !ok || n > attrs[cur] || n == attrs[cur] && k < cur {
			best[prop] = k
		}
	}
	var out []string
	for _, k := range best {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool {
		if attrs[out[i]] != attrs[out[j]] {
			return attrs[out[i]] > attrs[out[j]]
		}
		return out[i] < out[j]
	})
	if len(out) > 4 {
		out = out[:4]
	}
	return out
}
