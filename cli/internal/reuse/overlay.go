package reuse

import (
	"fmt"
	"html"
	"strings"

	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
)

// HTML renders the reuse list over the Figma frame: every component instance gets a box coloured by
// its verdict and labelled with the component to reuse, so the plan can be reviewed at a glance.
func (m *MatchResult) HTML(exp *figma.Export, img string) string {
	nodes := map[string]*figma.Node{}
	var walk func(n *figma.Node)
	walk = func(n *figma.Node) {
		nodes[n.ID] = n
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(exp.Frame)
	fb := exp.Frame.AbsoluteBoundingBox
	if fb == nil {
		fb = &figma.Box{Width: 375, Height: 812}
	}
	color := map[string]string{"reuse": "#16a34a", "confirm": "#d97706", "new": "#dc2626"}
	label := map[string]string{"reuse": "直接复用", "confirm": "需确认", "new": "需新建"}
	var boxes, rows strings.Builder
	// labels sit above their box; one that would cover an earlier label (Swiper and NavBar both start at
	// the frame's top-left) moves to the bottom edge inside its box
	type rect struct{ x, y, w, h float64 }
	var placed []rect
	free := func(r rect) bool {
		for _, p := range placed {
			if r.x < p.x+p.w && p.x < r.x+r.w && r.y < p.y+p.h && p.y < r.y+r.h {
				return false
			}
		}
		return true
	}
	for _, it := range m.Items {
		target := "需新建组件"
		if len(it.Candidates) > 0 && it.Verdict != "new" {
			target = it.Candidates[0].Name
		}
		for i, id := range it.NodeIDs {
			n := nodes[id]
			if n == nil || n.AbsoluteBoundingBox == nil {
				continue
			}
			b := n.AbsoluteBoundingBox
			tag := ""
			if i == 0 {
				text := strings.TrimSuffix(it.Figma, "（未组件化）") + " → " + target
				x, y := b.X-fb.X-2, b.Y-fb.Y
				w := float64(len([]rune(text)))*7 + 10
				pos := ""
				if x+w > fb.Width { // would run off the frame: align to the box's right edge instead
					x = b.X - fb.X + b.Width + 2 - w
					pos = ";left:auto;right:-2px"
				}
				if r := (rect{x, y - 18, w, 16}); free(r) || !free(rect{x, y + b.Height - 18, w, 16}) {
					placed = append(placed, r)
				} else {
					placed = append(placed, rect{x, y + b.Height - 18, w, 16})
					pos += ";top:auto;bottom:0"
				}
				tag = fmt.Sprintf(`<span style="background:%s%s">%s</span>`, color[it.Verdict], pos, html.EscapeString(text))
			}
			fmt.Fprintf(&boxes, `<div class="box" style="left:%.0fpx;top:%.0fpx;width:%.0fpx;height:%.0fpx;border-color:%s">%s</div>`,
				b.X-fb.X, b.Y-fb.Y, b.Width, b.Height, color[it.Verdict], tag)
		}
		why, ex := "项目和组件库里都没有", ""
		if len(it.Candidates) > 0 {
			c := it.Candidates[0]
			why = fmt.Sprintf("%.2f：%s", c.Score, c.Why)
			if c.Example != "" {
				ex = "<code>" + html.EscapeString(c.Example) + "</code>"
			}
			var ok []string // resolved mappings only; the unresolved notes stay in the .md for the agent
			for _, p := range c.Props {
				if strings.Contains(p, "→") || strings.HasPrefix(p, "项目惯用") {
					ok = append(ok, p)
				}
			}
			if len(ok) > 0 {
				ex += "<br><em>" + html.EscapeString(strings.Join(ok, "；")) + "</em>"
			}
		}
		fmt.Fprintf(&rows, `<tr><td>%s</td><td><b style="color:%s">%s</b></td><td><code>%s</code></td><td>%s</td><td>%s</td></tr>`,
			html.EscapeString(it.Figma), color[it.Verdict], label[it.Verdict], html.EscapeString(target), html.EscapeString(why), ex)
	}
	return fmt.Sprintf(`<!doctype html><html lang="zh"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>复用清单 · %s</title>
<style>
body{margin:0;font:13px/1.5 -apple-system,"PingFang SC","Noto Sans SC",sans-serif;color:#1f2329;background:#f5f6f7}
.wrap{display:flex;gap:24px;padding:24px;align-items:flex-start;flex-wrap:wrap}
.frame{position:relative;width:%.0fpx;flex:none;box-shadow:0 2px 12px rgba(0,0,0,.12);background:#fff}
.frame img{display:block;width:100%%}
.box{position:absolute;border:2px solid;border-radius:4px;box-sizing:border-box}
.box span{position:absolute;left:-2px;top:-18px;color:#fff;font-size:11px;line-height:16px;padding:0 5px;border-radius:3px;white-space:nowrap}
.side{flex:1;min-width:360px}
h1{font-size:18px;margin:0 0 4px}.sum{color:#646a73;margin:0 0 12px}
table{border-collapse:collapse;width:100%%;background:#fff}td,th{border:1px solid #e5e6eb;padding:6px 8px;text-align:left;vertical-align:top}
th{background:#f2f3f5}td:nth-child(-n+3){white-space:nowrap}td:nth-child(4){min-width:150px}td:last-child{min-width:280px}code{font-size:12px;word-break:break-all}em{color:#646a73;font-style:normal}
</style></head><body><div class="wrap">
<div class="frame"><img src="%s" alt="Figma frame">%s</div>
<div class="side"><h1>新设计稿复用清单：%s</h1>
<p class="sum">%s · %s · 共 %d 个组件：<b style="color:#16a34a">直接复用 %d</b> · <b style="color:#d97706">需确认 %d</b> · <b style="color:#dc2626">需新建 %d</b></p>
<table><tr><th>Figma</th><th>结论</th><th>推荐</th><th>依据</th><th>项目里的真实用法 / 怎么配置</th></tr>%s</table></div>
</div></body></html>
`, html.EscapeString(m.Page), fb.Width, html.EscapeString(img), boxes.String(), html.EscapeString(m.Page), html.EscapeString(m.Project), m.Platform,
		len(m.Items), m.Reuse, m.Confirm, m.New, rows.String())
}
