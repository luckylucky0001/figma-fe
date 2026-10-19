// Package fixture turns an annotated design HTML into a Figma-format export
// (same schema as the figma-plugin exporter). It exists only because no real Figma
// file is reachable from this environment; see docs/CONTRACTS.md §6.
package fixture

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/luckylucky0001/figma-fe/cli/internal/capture"
	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
)

// extractJS walks [data-fig] elements and emits Figma-like nodes (coordinates relative to #frame).
const extractJS = `(() => {
  const frame = document.querySelector('[data-fig="FRAME"]');
  const fr = frame.getBoundingClientRect();
  let seq = 2;
  const color = (s) => {
    const m = s && s.match(/rgba?\(([^)]+)\)/);
    if (!m) return null;
    const p = m[1].split(',').map(x => parseFloat(x));
    const a = p.length > 3 ? p[3] : 1;
    if (a === 0) return null;
    return {r: p[0]/255, g: p[1]/255, b: p[2]/255, a};
  };
  const conv = (el) => {
    const cs = getComputedStyle(el);
    const r = el.getBoundingClientRect();
    const type = el.dataset.fig;
    const n = { id: type === 'FRAME' && el === frame ? '1:1' : '1:' + (seq++), name: el.dataset.figName || el.tagName.toLowerCase(), type,
      absoluteBoundingBox: { x: r.left - fr.left, y: r.top - fr.top, width: r.width, height: r.height } };
    const bg = color(cs.backgroundColor);
    if (type === 'TEXT') {
      n.characters = el.textContent.trim();
      const fg = color(cs.color);
      if (fg) n.fills = [{type: 'SOLID', color: fg}];
      n.style = { fontFamily: cs.fontFamily.split(',')[0].replace(/['"]/g, '').trim(), fontWeight: parseFloat(cs.fontWeight),
        fontSize: parseFloat(cs.fontSize), lineHeightPx: parseFloat(cs.lineHeight) || parseFloat(cs.fontSize) * 1.2 };
    } else {
      if (bg) n.fills = [{type: 'SOLID', color: bg}];
      const bw = parseFloat(cs.borderTopWidth);
      if (bw > 0 && cs.borderTopStyle !== 'none') { const bc = color(cs.borderTopColor); if (bc) { n.strokes = [{type: 'SOLID', color: bc}]; n.strokeWeight = bw; } }
      const rad = parseFloat(cs.borderTopLeftRadius); if (rad) n.cornerRadius = rad;
      ['Left','Right','Top','Bottom'].forEach(s => { const v = parseFloat(cs['padding' + s]); if (v) n['padding' + s] = v; });
      if (cs.display.includes('flex')) { n.layoutMode = cs.flexDirection.startsWith('column') ? 'VERTICAL' : 'HORIZONTAL'; const g = parseFloat(cs.gap); if (g) n.itemSpacing = g; }
    }
    if (type === 'INSTANCE') { n._set = el.dataset.figSet; n._props = el.dataset.figProps ? JSON.parse(el.dataset.figProps) : {}; }
    const kids = [];
    const visit = (p) => { for (const c of p.children) { if (c.dataset && c.dataset.fig) kids.push(conv(c)); else visit(c); } };
    if (type !== 'TEXT') visit(el);
    if (kids.length) n.children = kids;
    return n;
  };
  return JSON.stringify({ frame: conv(frame), width: fr.width, height: fr.height });
})()`

type rawNode struct {
	figma.Node
	Set      string            `json:"_set"`
	Props    map[string]string `json:"_props"`
	Children []*rawNode        `json:"children"`
}

// Build renders htmlPath and writes export.json + frame@2x.png into outDir.
func Build(htmlPath, libraryPath, outDir string) error {
	lib, err := figma.LoadLibrary(libraryPath)
	if err != nil {
		return err
	}
	abs, err := filepath.Abs(htmlPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	b, err := capture.NewBrowser(60 * time.Second)
	if err != nil {
		return err
	}
	defer b.Close()
	// First pass at a tall viewport to read the frame size, then re-open at the exact size.
	if err := b.Open("file://"+abs, 375, 2000, 2); err != nil {
		return err
	}
	var raw string
	if err := b.Eval(extractJS, &raw); err != nil {
		return fmt.Errorf("extract: %w", err)
	}
	var res struct {
		Frame  *rawNode `json:"frame"`
		Width  float64  `json:"width"`
		Height float64  `json:"height"`
	}
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		return err
	}
	w, h := int(res.Width), int(res.Height)
	if err := b.Open("file://"+abs, w, h, 2); err != nil {
		return err
	}
	if err := b.Screenshot(filepath.Join(outDir, "frame@2x.png"), w, h, 2); err != nil {
		return err
	}

	exp := figma.Export{Schema: "figmafe-export/1", Components: map[string]figma.Component{}, ComponentSets: map[string]figma.ComponentSet{}}
	exp.Source.FileKey = "fixture"
	exp.Source.FileName = "Demo DS / OrderDetail (fixture from " + filepath.Base(htmlPath) + ")"
	exp.Source.ExportedAt = time.Now().UTC().Format(time.RFC3339)
	exp.Source.Exporter = "figmafe fixture"
	exp.Images.Frame = "frame@2x.png"
	exp.Images.Scale = 2

	setByName := map[string]string{}
	for id, s := range lib.ComponentSets {
		setByName[s.Name] = id
	}
	var convert func(r *rawNode) (*figma.Node, error)
	convert = func(r *rawNode) (*figma.Node, error) {
		n := r.Node
		n.Children = nil
		if r.Type == "INSTANCE" {
			setID, ok := setByName[r.Set]
			if !ok {
				return nil, fmt.Errorf("node %s: component set %q not in library", r.ID, r.Set)
			}
			set := lib.ComponentSets[setID]
			exp.ComponentSets[setID] = set
			// Fill unspecified variant axes with their first value, like Figma does.
			props := map[string]figma.ComponentProperty{}
			var parts []string
			axes := make([]string, 0, len(set.VariantProperties))
			for ax := range set.VariantProperties {
				axes = append(axes, ax)
			}
			sort.Strings(axes)
			for _, ax := range axes {
				v, ok := r.Props[ax]
				if !ok {
					v = set.VariantProperties[ax][0]
				}
				props[ax] = figma.ComponentProperty{Type: "VARIANT", Value: v}
				parts = append(parts, ax+"="+v)
			}
			compName := strings.Join(parts, ", ")
			if compName == "" {
				compName = set.Name
			}
			compID := setID + "/" + compName
			exp.Components[compID] = figma.Component{Key: set.Key + "/" + compName, Name: compName, ComponentSetID: setID}
			n.ComponentID = compID
			if len(props) > 0 {
				n.ComponentProperties = props
			}
		}
		for _, c := range r.Children {
			cn, err := convert(c)
			if err != nil {
				return nil, err
			}
			n.Children = append(n.Children, cn)
		}
		return &n, nil
	}
	frame, err := convert(res.Frame)
	if err != nil {
		return err
	}
	exp.Frame = frame
	out, err := json.MarshalIndent(exp, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, "export.json"), out, 0o644)
}
