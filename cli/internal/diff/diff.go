// Package diff compares a rendered page with the Figma frame: pixel mismatch (pixelmatch-style YIQ
// distance) plus node-level checks (text position/size/style, component bounds), producing
// structured deviations the Agent can act on.
package diff

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/luckylucky0001/figma-fe/cli/internal/capture"
	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
	"github.com/luckylucky0001/figma-fe/cli/internal/plan"
)

type Deviation struct {
	Severity string  `json:"severity"` // error | warn
	Kind     string  `json:"kind"`     // missing-text | missing-component | not-ds | position | size | font-size | font-weight | color
	FigmaID  string  `json:"figmaId"`
	Target   string  `json:"target"` // human label
	Expected string  `json:"expected"`
	Actual   string  `json:"actual"`
	Delta    float64 `json:"delta,omitempty"`
}

type Region struct {
	FigmaID  string  `json:"figmaId"`
	Target   string  `json:"target"`
	Mismatch float64 `json:"mismatch"`
}

type Result struct {
	Platform      string      `json:"platform"`
	PixelMismatch float64     `json:"pixelMismatch"`
	SizeOK        bool        `json:"sizeOk"`
	Texts         int         `json:"texts"`
	TextsMatched  int         `json:"textsMatched"`
	Instances     int         `json:"instances"`
	InstancesByDS int         `json:"instancesByDs"` // rendered by the mapped design-system component (runtime-verified reuse)
	Deviations    []Deviation `json:"deviations"`
	Regions       []Region    `json:"regions"`
	Errors, Warns int
}

// Pixel compares two PNGs, writes a diff PNG, and returns mismatch ratio plus the mismatch mask.
func Pixel(designPath, screenPath, diffPath string) (float64, bool, *image.Gray, error) {
	a, err := readPNG(designPath)
	if err != nil {
		return 0, false, nil, err
	}
	b, err := readPNG(screenPath)
	if err != nil {
		return 0, false, nil, err
	}
	ab, bb := a.Bounds(), b.Bounds()
	sizeOK := ab.Dx() == bb.Dx() && ab.Dy() == bb.Dy()
	w, h := min(ab.Dx(), bb.Dx()), max(ab.Dy(), bb.Dy())
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	mask := image.NewGray(image.Rect(0, 0, w, h))
	const maxDelta = 35215.0
	thr := 0.1 * 0.1 * maxDelta // pixelmatch default threshold 0.1
	bad := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var ca, cb color.Color = color.White, color.White
			if y < ab.Dy() {
				ca = a.At(ab.Min.X+x, ab.Min.Y+y)
			}
			if y < bb.Dy() {
				cb = b.At(bb.Min.X+x, bb.Min.Y+y)
			}
			d := yiq(ca, cb)
			if d > thr {
				bad++
				out.Set(x, y, color.RGBA{255, 0, 0, 255})
				mask.SetGray(x, y, color.Gray{255})
			} else {
				r, g, bl, _ := ca.RGBA()
				l := uint8(0.1*float64(r>>8)*0.299+0.1*float64(g>>8)*0.587+0.1*float64(bl>>8)*0.114) + 229
				out.Set(x, y, color.RGBA{l, l, l, 255})
			}
		}
	}
	f, err := os.Create(diffPath)
	if err != nil {
		return 0, sizeOK, nil, err
	}
	defer f.Close()
	if err := png.Encode(f, out); err != nil {
		return 0, sizeOK, nil, err
	}
	return float64(bad) / float64(w*h), sizeOK, mask, nil
}

func readPNG(p string) (image.Image, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

// yiq is pixelmatch's perceptual colour delta (alpha-blended onto white).
func yiq(a, b color.Color) float64 {
	r1, g1, b1 := blend(a)
	r2, g2, b2 := blend(b)
	y := (r1-r2)*0.29889531 + (g1-g2)*0.58662247 + (b1-b2)*0.11448223
	i := (r1-r2)*0.59597799 - (g1-g2)*0.27417610 - (b1-b2)*0.32180189
	q := (r1-r2)*0.21147017 - (g1-g2)*0.52261711 + (b1-b2)*0.31114694
	return 0.5053*y*y + 0.299*i*i + 0.1957*q*q
}

func blend(c color.Color) (float64, float64, float64) {
	r, g, b, a := c.RGBA()
	af := float64(a) / 65535
	f := func(v uint32) float64 { return 255 + (float64(v>>8)-255)*af }
	if a == 0 {
		return 255, 255, 255
	}
	return f(r), f(g), f(b)
}

// Nodes runs node-level checks of the Figma frame against rendered nodes.
func Nodes(exp *figma.Export, p *plan.Plan, rn *capture.Nodes, tol float64, mask *image.Gray, dsName func(string) string) *Result {
	res := &Result{Platform: p.Platform}
	add := func(d Deviation) {
		res.Deviations = append(res.Deviations, d)
		if d.Severity == "error" {
			res.Errors++
		} else {
			res.Warns++
		}
	}
	used := map[int]bool{}
	norm := func(s string) string { return strings.Join(strings.Fields(s), "") }
	near := func(cands []int, b figma.Box) int {
		best, bd := -1, math.MaxFloat64
		for _, i := range cands {
			n := rn.Nodes[i]
			d := math.Hypot(n.X-b.X, n.Y-b.Y)
			if d < bd {
				best, bd = i, d
			}
		}
		return best
	}
	// 1) component instances
	for _, e := range p.Expected {
		if e.Component == "" {
			continue
		}
		res.Instances++
		label := fmt.Sprintf("%s「%s」", e.Set, trunc(e.Text, 12))
		idx := -1
		for i, n := range rn.Nodes {
			if n.FigmaID != nil && *n.FigmaID == e.FigmaID {
				idx = i
				break
			}
		}
		want := dsName(e.Component)
		if idx < 0 {
			var cands []int
			for i, n := range rn.Nodes {
				if !used[i] && n.DS != nil && *n.DS == want {
					cands = append(cands, i)
				}
			}
			idx = near(cands, e.Box)
			if idx >= 0 && math.Hypot(rn.Nodes[idx].X-e.Box.X, rn.Nodes[idx].Y-e.Box.Y) > 60 {
				idx = -1
			}
		}
		if idx < 0 {
			add(Deviation{"error", "missing-component", e.FigmaID, label, "由 " + e.Component + " 渲染", "页面里没有找到对应的组件节点", 0})
			continue
		}
		used[idx] = true
		n := rn.Nodes[idx]
		if n.DS == nil || *n.DS != want {
			got := "非组件库节点"
			if n.DS != nil {
				got = *n.DS
			}
			add(Deviation{"error", "not-ds", e.FigmaID, label, want, got, 0})
		} else {
			res.InstancesByDS++
		}
		checkBox(add, e.FigmaID, label, e.Box, n, tol)
	}
	// 2) texts
	var texts []*figma.Node
	figma.Walk(exp.Frame, func(n *figma.Node) bool {
		if n.Type == "TEXT" && strings.TrimSpace(n.Characters) != "" {
			texts = append(texts, n)
		}
		return true
	})
	usedT := map[int]bool{}
	for _, t := range texts {
		res.Texts++
		b := exp.Rel(t)
		want := norm(t.Characters)
		// Prefer plain text nodes; component roots also carry their joined text (e.g. DsTopBar "订单详情").
		var cands, compCands []int
		for i, n := range rn.Nodes {
			if !usedT[i] && n.Text != nil && norm(*n.Text) == want {
				if n.DS == nil {
					cands = append(cands, i)
				} else {
					compCands = append(compCands, i)
				}
			}
		}
		if len(cands) == 0 {
			cands = compCands
		}
		idx := near(cands, b)
		label := fmt.Sprintf("文本「%s」", trunc(t.Characters, 14))
		if idx < 0 {
			// rendered as part of a bigger text run (e.g. Price "¥299.00" on Android)?
			if inst := t.NearestInstance(); inst != nil {
				joined := ""
				for _, c := range inst.Children {
					if c.Type == "TEXT" {
						joined += c.Characters
					}
				}
				found := false
				for _, n := range rn.Nodes {
					if n.Text != nil && strings.Contains(norm(*n.Text), norm(joined)) && joined != "" {
						found = true
						break
					}
				}
				if found {
					res.TextsMatched++
					continue
				}
			}
			add(Deviation{"error", "missing-text", t.ID, label, t.Characters, "页面里没有这段文本", 0})
			continue
		}
		usedT[idx] = true
		res.TextsMatched++
		n := rn.Nodes[idx]
		// Text boxes differ in width between Figma (often fixed / fill width) and code (hug content or
		// padded parent), so x is alignment-agnostic: left, centre or right edge may match. y uses the centre.
		dx := math.Min(math.Abs(n.X-b.X), math.Min(math.Abs((n.X+n.W/2)-(b.X+b.Width/2)), math.Abs((n.X+n.W)-(b.X+b.Width))))
		dy := math.Abs((n.Y + n.H/2) - (b.Y + b.Height/2))
		if d := math.Max(dx, dy); d > tol {
			add(Deviation{sev(d), "position", t.ID, label, fmt.Sprintf("(%.0f, %.0f)", b.X, b.Y), fmt.Sprintf("(%.0f, %.0f)", n.X, n.Y), round1(d)})
		}
		if t.Style != nil && n.FontSize > 0 {
			if math.Abs(n.FontSize-t.Style.FontSize) > 0.5 {
				add(Deviation{"error", "font-size", t.ID, label, fmt.Sprintf("%gpx", t.Style.FontSize), fmt.Sprintf("%gpx", n.FontSize), math.Abs(n.FontSize - t.Style.FontSize)})
			}
			if n.FontWeight > 0 && math.Abs(n.FontWeight-t.Style.FontWeight) >= 100 {
				add(Deviation{"warn", "font-weight", t.ID, label, fmt.Sprint(t.Style.FontWeight), fmt.Sprint(n.FontWeight), 0})
			}
		} else if t.Style != nil && n.H > 0 && math.Abs(n.H-t.Style.LineHeightPx) > tol+1 {
			add(Deviation{"warn", "line-height", t.ID, label, fmt.Sprintf("行高 %g", t.Style.LineHeightPx), fmt.Sprintf("高 %.1f", n.H), round1(math.Abs(n.H - t.Style.LineHeightPx))})
		}
		if c, ok := t.Fill(); ok && n.Color != "" {
			if dist := colorDist(c, n.Color); dist > 12 {
				add(Deviation{"error", "color", t.ID, label, c.Hex(), cssToHex(n.Color), round1(dist)})
			}
		}
	}
	// 3) mismatch regions
	if mask != nil {
		for _, e := range p.Expected {
			if e.Set == "Card" || e.Set == "ActionBar" { // containers: their children are reported individually
				continue
			}
			r := ratioIn(mask, e.Box)
			if r > 0.08 {
				res.Regions = append(res.Regions, Region{e.FigmaID, fmt.Sprintf("%s「%s」", e.Set, trunc(e.Text, 12)), round3(r)})
			}
		}
		sort.Slice(res.Regions, func(i, j int) bool { return res.Regions[i].Mismatch > res.Regions[j].Mismatch })
		if len(res.Regions) > 8 {
			res.Regions = res.Regions[:8]
		}
	}
	sort.SliceStable(res.Deviations, func(i, j int) bool {
		return res.Deviations[i].Severity == "error" && res.Deviations[j].Severity != "error"
	})
	return res
}

func checkBox(add func(Deviation), id, label string, want figma.Box, n capture.Node, tol float64) {
	dx, dy := n.X-want.X, n.Y-want.Y
	if d := math.Max(math.Abs(dx), math.Abs(dy)); d > tol {
		add(Deviation{sev(d), "position", id, label, fmt.Sprintf("(%.0f, %.0f)", want.X, want.Y), fmt.Sprintf("(%.0f, %.0f)", n.X, n.Y), round1(d)})
	}
	dw, dh := n.W-want.Width, n.H-want.Height
	if d := math.Max(math.Abs(dw), math.Abs(dh)); d > tol {
		add(Deviation{sev(d), "size", id, label, fmt.Sprintf("%.0f×%.0f", want.Width, want.Height), fmt.Sprintf("%.0f×%.0f", n.W, n.H), round1(d)})
	}
}

func sev(d float64) string {
	if d > 4 {
		return "error"
	}
	return "warn"
}

func ratioIn(mask *image.Gray, b figma.Box) float64 {
	x0, y0 := int(b.X*2), int(b.Y*2)
	x1, y1 := int((b.X+b.Width)*2), int((b.Y+b.Height)*2)
	mb := mask.Bounds()
	x0, y0, x1, y1 = max(x0, 0), max(y0, 0), min(x1, mb.Dx()), min(y1, mb.Dy())
	if x1 <= x0 || y1 <= y0 {
		return 0
	}
	bad := 0
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if mask.GrayAt(x, y).Y > 0 {
				bad++
			}
		}
	}
	return float64(bad) / float64((x1-x0)*(y1-y0))
}

func colorDist(c figma.Color, css string) float64 {
	var r, g, b float64
	if _, err := fmt.Sscanf(strings.ReplaceAll(css, " ", ""), "rgb(%g,%g,%g)", &r, &g, &b); err != nil {
		if _, err := fmt.Sscanf(strings.ReplaceAll(css, " ", ""), "rgba(%g,%g,%g", &r, &g, &b); err != nil {
			return 0
		}
	}
	return math.Sqrt(math.Pow(c.R*255-r, 2) + math.Pow(c.G*255-g, 2) + math.Pow(c.B*255-b, 2))
}

func cssToHex(css string) string {
	var r, g, b float64
	if _, err := fmt.Sscanf(strings.ReplaceAll(css, " ", ""), "rgb(%g,%g,%g)", &r, &g, &b); err != nil {
		return css
	}
	return fmt.Sprintf("#%02X%02X%02X", int(r), int(g), int(b))
}

func trunc(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n]) + "…"
	}
	return s
}

func round1(f float64) float64 { return math.Round(f*10) / 10 }
func round3(f float64) float64 { return math.Round(f*1000) / 1000 }
