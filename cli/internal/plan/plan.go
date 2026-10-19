// Package plan turns a Figma page export + the component map into a per-platform reuse plan (Phase 2).
// The plan is the gate artefact: the Agent shows it, a human confirms, only then code is written.
package plan

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
	"github.com/luckylucky0001/figma-fe/cli/internal/mapping"
)

// Expected is one Figma instance the page must render with a mapped component.
type Expected struct {
	FigmaID   string            `json:"figmaId"`
	Set       string            `json:"set"`
	Component string            `json:"component,omitempty"` // empty = unmapped (gap)
	Variant   map[string]string `json:"variant,omitempty"`
	Code      string            `json:"code,omitempty"`
	Text      string            `json:"text,omitempty"`
	Box       figma.Box         `json:"box"`
}

type Plan struct {
	Page     string         `json:"page"`
	Platform string         `json:"platform"`
	Width    float64        `json:"width"`
	Height   float64        `json:"height"`
	Expected []Expected     `json:"expected"`
	Counts   map[string]int `json:"counts"` // component → expected count
	Gaps     []string       `json:"gaps"`
	Lines    []string       `json:"-"`
}

type Tokens struct {
	Colors map[string]string // #HEX → name
	Fonts  map[string]string // "size/lh/weight" → name
}

func LoadTokens(path string) (*Tokens, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	root := map[string]map[string]any{}
	for k, v := range raw {
		var g map[string]any
		if json.Unmarshal(v, &g) == nil {
			root[k] = g
		}
	}
	t := &Tokens{Colors: map[string]string{}, Fonts: map[string]string{}}
	for k, v := range root["color"] {
		if o, ok := v.(map[string]any); ok {
			if s, ok := o["$value"].(string); ok {
				t.Colors[strings.ToUpper(s)] = "color." + k
			}
		}
	}
	for k, v := range root["font"] {
		o, ok := v.(map[string]any)
		if !ok {
			continue
		}
		val, _ := o["$value"].(map[string]any)
		fs, _ := val["fontSize"].(map[string]any)
		size, _ := fs["value"].(float64)
		lh, _ := val["lineHeight"].(float64)
		w, _ := val["fontWeight"].(float64)
		t.Fonts[fontKey(size, math.Round(size*lh), w)] = "font." + k
	}
	return t, nil
}

func fontKey(size, lh, w float64) string { return fmt.Sprintf("%g/%g/%g", size, lh, w) }

// Ref renders a token name as platform code.
func Ref(platform, token string) string {
	if token == "" {
		return ""
	}
	parts := strings.SplitN(token, ".", 2)
	if platform == "h5" {
		if parts[0] == "font" {
			return "var(--ds-font-" + parts[1] + "-*)"
		}
		return "var(--ds-" + parts[0] + "-" + parts[1] + ")"
	}
	camel := func(s string) string {
		var b strings.Builder
		for _, w := range strings.Split(s, "-") {
			if w != "" {
				b.WriteString(strings.ToUpper(w[:1]) + w[1:])
			}
		}
		return b.String()
	}
	switch parts[0] {
	case "color":
		return "DsColors." + camel(parts[1])
	case "font":
		return "DsTypography." + camel(parts[1])
	case "space":
		return "DsSpacing." + camel(parts[1])
	case "radius":
		return "DsRadius." + camel(parts[1])
	}
	return token
}

var spaceTokens = map[float64]string{4: "space.xs", 8: "space.sm", 12: "space.md", 16: "space.lg", 24: "space.xl"}

// Build walks the frame and resolves each node against the map.
func Build(exp *figma.Export, m *mapping.Map, tok *Tokens, page, platform string) *Plan {
	p := &Plan{Page: page, Platform: platform, Counts: map[string]int{}}
	fb := exp.Frame.AbsoluteBoundingBox
	p.Width, p.Height = fb.Width, fb.Height
	gapSeen := map[string]bool{}
	var walk func(n *figma.Node, depth int, inComponent bool)
	walk = func(n *figma.Node, depth int, inComponent bool) {
		ind := strings.Repeat("  ", depth)
		b := exp.Rel(n)
		geo := fmt.Sprintf("(%g,%g %g×%g)", b.X, b.Y, b.Width, b.Height)
		switch n.Type {
		case "INSTANCE":
			set := exp.SetName(n)
			e := m.Entry(set)
			var im *mapping.Impl
			if e != nil {
				im = e.Impl(platform)
			}
			if im == nil {
				g := m.GapFor(set, platform)
				msg := "⚠ 未映射：" + set
				if g != nil {
					msg += " → 处理方案：" + g.Decision
				}
				p.Lines = append(p.Lines, fmt.Sprintf("%s- [%s] **%s** %s %s — 文本「%s」", ind, n.ID, set, geo, msg, n.Text()))
				if !gapSeen[set] {
					p.Gaps = append(p.Gaps, set)
					gapSeen[set] = true
				}
				p.Expected = append(p.Expected, Expected{FigmaID: n.ID, Set: set, Variant: n.Variant(), Text: n.Text(), Box: b})
				return
			}
			code := codeFor(platform, im, n)
			p.Lines = append(p.Lines, fmt.Sprintf("%s- [%s] **%s** → `%s` %s", ind, n.ID, set, code, geo))
			p.Expected = append(p.Expected, Expected{FigmaID: n.ID, Set: set, Component: im.Component, Variant: n.Variant(), Code: code, Text: n.Text(), Box: b})
			p.Counts[im.Component]++
			for _, w := range variantWarnings(im, n) {
				p.Lines = append(p.Lines, ind+"  - ⚠ "+w)
			}
			if e.Slot {
				used := templateLayers(im)
				for _, c := range n.Children {
					if c.Type == "TEXT" && used[c.Name] {
						continue // rendered by the component's own text prop (e.g. Card title)
					}
					walk(c, depth+1, false)
				}
			}
			return
		case "TEXT":
			if inComponent {
				return
			}
			st := n.Style
			var font, color string
			if st != nil {
				font = tok.Fonts[fontKey(st.FontSize, math.Round(st.LineHeightPx), st.FontWeight)]
			}
			if c, ok := n.Fill(); ok {
				color = tok.Colors[c.Hex()]
			}
			desc := "文本「" + n.Characters + "」"
			if font != "" || color != "" {
				desc += fmt.Sprintf(" 样式 %s + %s", orRaw(Ref(platform, font), st), orHex(Ref(platform, color), n))
			} else {
				desc += " ⚠ 样式不在 token 内"
			}
			p.Lines = append(p.Lines, fmt.Sprintf("%s- [%s] %s %s", ind, n.ID, desc, geo))
		case "RECTANGLE":
			p.Lines = append(p.Lines, fmt.Sprintf("%s- [%s] 矩形 %s %s", ind, n.ID, n.Name, geo))
		default: // FRAME / GROUP: layout container
			lay := ""
			if n.LayoutMode != "" {
				dir := map[string]string{"VERTICAL": "纵向", "HORIZONTAL": "横向"}[n.LayoutMode]
				lay = dir
				if n.ItemSpacing > 0 {
					lay += fmt.Sprintf(" 间距 %g(%s)", n.ItemSpacing, spaceTokens[n.ItemSpacing])
				}
			}
			if n.PaddingTop+n.PaddingLeft > 0 {
				lay += fmt.Sprintf(" 内边距 %g/%g", n.PaddingTop, n.PaddingLeft)
			}
			if c, ok := n.Fill(); ok {
				lay += " 背景 " + orHex(Ref(platform, tok.Colors[c.Hex()]), n)
			}
			if depth > 0 {
				p.Lines = append(p.Lines, fmt.Sprintf("%s- [%s] 容器 %s %s %s", ind, n.ID, n.Name, strings.TrimSpace(lay), geo))
			}
			for _, c := range n.Children {
				walk(c, depth+1, inComponent)
			}
		}
	}
	walk(exp.Frame, 0, false)
	return p
}

func orRaw(ref string, st *figma.TypeStyle) string {
	if ref != "" || st == nil {
		return ref
	}
	return fmt.Sprintf("⚠%g/%g/%g", st.FontSize, st.LineHeightPx, st.FontWeight)
}

func orHex(ref string, n *figma.Node) string {
	if ref != "" {
		return ref
	}
	if c, ok := n.Fill(); ok {
		return "⚠" + c.Hex()
	}
	return ""
}

var layerRe = regexp.MustCompile(`\{(\w+)\}`)

// templateLayers lists the TEXT layer names an Impl's templates consume.
func templateLayers(im *mapping.Impl) map[string]bool {
	out := map[string]bool{}
	for _, t := range im.Texts {
		for _, m := range layerRe.FindAllStringSubmatch(t, -1) {
			out[m[1]] = true
		}
	}
	return out
}

func layerText(n *figma.Node, layer string) string {
	t := ""
	figma.Walk(n, func(x *figma.Node) bool {
		if t == "" && x != n && x.Type == "TEXT" && x.Name == layer {
			t = x.Characters
		}
		// do not descend into nested instances that are composed separately
		return t == "" && (x == n || x.Type != "INSTANCE" || len(x.ComponentProperties) == 0 || true)
	})
	return t
}

// fill resolves "{A}{B}|num" against the instance's TEXT layers.
func fill(n *figma.Node, tpl string) (string, bool) {
	mod := ""
	if i := strings.LastIndex(tpl, "|"); i >= 0 {
		tpl, mod = tpl[:i], tpl[i+1:]
	}
	val := layerRe.ReplaceAllStringFunc(tpl, func(m string) string { return layerText(n, m[1:len(m)-1]) })
	if strings.TrimSpace(val) == "" {
		return "", false
	}
	switch mod {
	case "num":
		val = strings.TrimLeft(regexp.MustCompile(`[^0-9.]`).ReplaceAllString(val, ""), "")
		return val, true
	case "int":
		return regexp.MustCompile(`[^0-9]`).ReplaceAllString(val, ""), true
	}
	return val, false
}

var propOrder = []string{"title", "text", "amount", "value", "price"}

func sortedProps(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	rank := func(k string) int {
		for i, p := range propOrder {
			if p == k {
				return i
			}
		}
		return len(propOrder)
	}
	sort.Slice(keys, func(i, j int) bool {
		if rank(keys[i]) != rank(keys[j]) {
			return rank(keys[i]) < rank(keys[j])
		}
		return keys[i] < keys[j]
	})
	return keys
}

func codeFor(platform string, im *mapping.Impl, n *figma.Node) string {
	var args []string
	children := ""
	for _, prop := range sortedProps(im.Texts) {
		v, numeric := fill(n, im.Texts[prop])
		if v == "" {
			continue
		}
		if prop == "children" {
			children = v
			continue
		}
		switch {
		case platform == "h5" && numeric:
			args = append(args, fmt.Sprintf("%s={%s}", prop, v))
		case platform == "h5":
			args = append(args, fmt.Sprintf("%s=\"%s\"", prop, v))
		case numeric && strings.Contains(im.Texts[prop], "|num") && !strings.Contains(v, "."):
			args = append(args, fmt.Sprintf("%s = %s.0", prop, v))
		case numeric:
			args = append(args, fmt.Sprintf("%s = %s", prop, v))
		default:
			args = append(args, fmt.Sprintf("%s = \"%s\"", prop, v))
		}
	}
	keys := make([]string, 0, len(n.ComponentProperties))
	for k := range n.ComponentProperties {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := fmt.Sprint(n.ComponentProperties[k].Value)
		if frag := im.Props[k][v]; frag != "" && frag != "-" {
			args = append(args, frag)
		}
	}
	if platform == "h5" {
		s := "<" + im.Component
		if len(args) > 0 {
			s += " " + strings.Join(args, " ")
		}
		s += fmt.Sprintf(` data-figma-id="%s"`, n.ID)
		if children != "" {
			return s + ">" + children + "</" + im.Component + ">"
		}
		return s + " />"
	}
	args = append(args, fmt.Sprintf(`modifier = Modifier.figmaId("%s")`, n.ID))
	return im.Component + "(" + strings.Join(args, ", ") + ")"
}

func variantWarnings(im *mapping.Impl, n *figma.Node) []string {
	var out []string
	for k, p := range n.ComponentProperties {
		v := fmt.Sprint(p.Value)
		if frag, ok := im.Props[k][v]; ok && frag == "-" {
			out = append(out, fmt.Sprintf("变体 %s=%s 在 %s 不支持：需要和组件库 owner 确认降级方案", k, v, im.Component))
		}
	}
	return out
}

func (p *Plan) Markdown(m *mapping.Map) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# 复用方案：%s（%s）\n\n", p.Page, p.Platform)
	fmt.Fprintf(&b, "> 由 `figmafe plan` 生成。**闸门：请确认本方案后再写代码。**\n> 设计帧 %g×%g；节点 id 即 Figma 节点 id，实现时挂到对应组件上（H5 `data-figma-id`，Android `Modifier.figmaId`）。\n\n", p.Width, p.Height)
	b.WriteString("## 1. 组件清单\n\n| 组件 | 数量 |\n|---|---|\n")
	keys := make([]string, 0, len(p.Counts))
	for k := range p.Counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	total := 0
	for _, k := range keys {
		fmt.Fprintf(&b, "| `%s` | %d |\n", k, p.Counts[k])
		total += p.Counts[k]
	}
	fmt.Fprintf(&b, "| **合计** | **%d** |\n\n", total)
	if len(p.Gaps) > 0 {
		b.WriteString("## 2. 待决项（Figma 有、代码没有）\n\n")
		for _, g := range p.Gaps {
			d := "（未登记，需要和设计 / 组件库 owner 确认）"
			if gp := m.GapFor(g, p.Platform); gp != nil {
				d = gp.Decision
			}
			fmt.Fprintf(&b, "- **%s**：%s\n", g, d)
		}
		b.WriteString("\n")
	}
	b.WriteString("## 3. 页面结构（节点 → 组件 → 代码）\n\n")
	for _, l := range p.Lines {
		b.WriteString(l + "\n")
	}
	b.WriteString("\n## 4. 约束（来自 component-map）\n\n")
	seen := map[string]bool{}
	for _, e := range p.Expected {
		if seen[e.Set] {
			continue
		}
		seen[e.Set] = true
		if en := m.Entry(e.Set); en != nil {
			im := en.Impl(p.Platform)
			for _, c := range en.Constraints {
				fmt.Fprintf(&b, "- **%s**：%s\n", e.Set, c)
			}
			if im != nil && im.Notes != "" {
				fmt.Fprintf(&b, "- **%s**（%s）：%s\n", e.Set, p.Platform, im.Notes)
			}
		}
	}
	return b.String()
}
