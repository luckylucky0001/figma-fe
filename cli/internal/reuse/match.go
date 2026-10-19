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
	Name      string  `json:"name"`
	Kind      string  `json:"kind"`
	Source    string  `json:"source"`
	Score     float64 `json:"score"`
	Usage     int     `json:"usage"`
	Example   string  `json:"example,omitempty"`
	ExampleAt string  `json:"exampleAt,omitempty"`
	Why       string  `json:"why"`
}

type MatchItem struct {
	Figma      string      `json:"figma"`
	Instances  int         `json:"instances"`
	NodeIDs    []string    `json:"nodeIds"`
	Variants   []string    `json:"variants,omitempty"`
	Verdict    string      `json:"verdict"` // reuse | confirm | new
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
		seen[a.Kind+"|"+a.Name] = true
	}
	for _, a := range r.LibAll {
		if !seen[a.Kind+"|"+a.Name] && a.Source != "Android SDK widgets" {
			pool = append(pool, a)
		}
	}
	for _, name := range order {
		a := sets[name]
		item := MatchItem{Figma: name, Instances: len(a.ids), NodeIDs: a.ids}
		for v := range a.variants {
			item.Variants = append(item.Variants, v)
		}
		sort.Strings(item.Variants)
		fc := ConceptOf(strings.TrimSuffix(name, "（未组件化）"))
		for _, as := range pool {
			sc, why := score(name, fc, a.variants, as)
			if sc < 0.5 {
				continue
			}
			ex, at := as.Example, as.ExampleAt
			if ex == "" && as.Kind == "ui-lib" {
				ex, at = libSnippet(as), "（"+as.Source+" 标准用法，项目未用过）"
			}
			item.Candidates = append(item.Candidates, Candidate{as.Name, as.Kind, as.Source, round2(sc), as.Usage, ex, at, why})
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
	case fn == an && figmaConcept != "" && a.Concept != "" && figmaConcept != a.Concept:
		s, why = 0.6, append(why, "同名但用途不同（"+a.Concept+"），需确认")
	case fn == an:
		s, why = 0.9, append(why, "同名")
	case figmaConcept != "" && figmaConcept == a.Concept:
		s, why = 0.72, append(why, "同为 "+figmaConcept)
	case strings.Contains(an, fn) && len(fn) >= 4 || strings.Contains(fn, an) && len(an) >= 4:
		s, why = 0.6, append(why, "名称包含")
	default:
		return 0, ""
	}
	// specialised components (AddressList, ContactCard, CouponCell) carry domain words the design does not
	if s < 0.9 && a.Concept != "" && a.Kind == "ui-lib" {
		ws := splitWords(a.Name)
		for len(ws) > 1 && vendorPrefix[ws[0]] {
			ws = ws[1:]
		}
		residual := strings.Join(ws, "")
		for _, c := range concepts {
			if c.Name != a.Concept {
				continue
			}
			for _, w := range c.Words {
				if len(w) >= 3 {
					residual = strings.Replace(residual, w, "", 1)
				}
			}
		}
		residual = strings.TrimSuffix(residual, "s")
		if len(residual) >= 3 && !strings.Contains(strings.ToLower(figmaName), residual) {
			s -= 0.12
			why = append(why, "专用组件（"+residual+"）")
		}
	}
	// the library's canonical component of the concept (Cell for cell, Tag for tag) beats look-alikes
	for _, c := range concepts {
		if c.Name == a.Concept && an == c.Words[0] && fn != an {
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
		s += math.Min(0.2, 0.05*math.Log2(float64(a.Usage)+1)) // proven in this codebase, more so when widely used
		why = append(why, fmt.Sprintf("项目里已用 %d 次", a.Usage))
	}
	switch a.Kind {
	case "design-system":
		s += 0.05
	case "page-local":
		s -= 0.2
		why = append(why, "页面私有，需先抽出")
	}
	if s > 1 {
		s = 1
	}
	return s, strings.Join(why, "，")
}

func normName(s string) string {
	w := splitWords(s)
	j := strings.Join(w, "")
	for _, p := range []string{"van", "nut", "ds", "legacy", "qmui", "material", "app"} {
		j = strings.TrimPrefix(j, p)
	}
	return strings.TrimSuffix(strings.TrimSuffix(j, "view"), "component")
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
			fmt.Fprintf(&b, "| %s | %d | %s | — | 项目和组件库里都没有 | — |\n", it.Figma, it.Instances, verdict)
			continue
		}
		c := it.Candidates[0]
		ex := "—"
		if c.Example != "" {
			ex = "`" + strings.ReplaceAll(c.Example, "|", "\\|") + "` " + c.ExampleAt
		}
		fmt.Fprintf(&b, "| %s | %d | %s | `%s`（%s %s） | %.2f：%s | %s |\n", it.Figma, it.Instances, verdict, c.Name, kindLabel(c.Kind), c.Source, c.Score, c.Why, ex)
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
		if len(it.Variants) > 0 {
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
