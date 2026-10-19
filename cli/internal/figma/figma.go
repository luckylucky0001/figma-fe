// Package figma holds the subset of the Figma REST node model that figmafe uses.
// Field names follow GET /v1/files/:key/nodes so plugin exports and REST exports are interchangeable.
package figma

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type Color struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
	A float64 `json:"a"`
}

// Hex returns #RRGGBB (alpha ignored).
func (c Color) Hex() string {
	return fmt.Sprintf("#%02X%02X%02X", to255(c.R), to255(c.G), to255(c.B))
}

func to255(v float64) int { return int(math.Round(v * 255)) }

type Paint struct {
	Type  string `json:"type"`
	Color *Color `json:"color,omitempty"`
}

type Box struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type TypeStyle struct {
	FontFamily   string  `json:"fontFamily,omitempty"`
	FontWeight   float64 `json:"fontWeight,omitempty"`
	FontSize     float64 `json:"fontSize,omitempty"`
	LineHeightPx float64 `json:"lineHeightPx,omitempty"`
}

type ComponentProperty struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type Node struct {
	ID                  string                       `json:"id"`
	Name                string                       `json:"name"`
	Type                string                       `json:"type"`
	Children            []*Node                      `json:"children,omitempty"`
	AbsoluteBoundingBox *Box                         `json:"absoluteBoundingBox,omitempty"`
	Fills               []Paint                      `json:"fills,omitempty"`
	Strokes             []Paint                      `json:"strokes,omitempty"`
	StrokeWeight        float64                      `json:"strokeWeight,omitempty"`
	CornerRadius        float64                      `json:"cornerRadius,omitempty"`
	PaddingLeft         float64                      `json:"paddingLeft,omitempty"`
	PaddingRight        float64                      `json:"paddingRight,omitempty"`
	PaddingTop          float64                      `json:"paddingTop,omitempty"`
	PaddingBottom       float64                      `json:"paddingBottom,omitempty"`
	ItemSpacing         float64                      `json:"itemSpacing,omitempty"`
	LayoutMode          string                       `json:"layoutMode,omitempty"`
	Characters          string                       `json:"characters,omitempty"`
	Style               *TypeStyle                   `json:"style,omitempty"`
	ComponentID         string                       `json:"componentId,omitempty"`
	ComponentProperties map[string]ComponentProperty `json:"componentProperties,omitempty"`

	Parent *Node `json:"-"`
}

type Component struct {
	Key            string `json:"key"`
	Name           string `json:"name"`
	ComponentSetID string `json:"componentSetId"`
}

type ComponentSet struct {
	Key               string              `json:"key"`
	Name              string              `json:"name"`
	Description       string              `json:"description,omitempty"`
	VariantProperties map[string][]string `json:"variantProperties,omitempty"`
}

type Export struct {
	Schema string `json:"schema"`
	Source struct {
		FileKey    string `json:"fileKey"`
		FileName   string `json:"fileName"`
		ExportedAt string `json:"exportedAt"`
		Exporter   string `json:"exporter"`
	} `json:"source"`
	Frame         *Node                   `json:"frame"`
	Components    map[string]Component    `json:"components"`
	ComponentSets map[string]ComponentSet `json:"componentSets"`
	Images        struct {
		Frame string  `json:"frame"`
		Scale float64 `json:"scale"`
	} `json:"images"`

	Dir string `json:"-"`
}

// Library is the component-set listing exported from the Figma library file.
type Library struct {
	Schema        string                  `json:"schema"`
	FileName      string                  `json:"fileName"`
	ComponentSets map[string]ComponentSet `json:"componentSets"`
}

func LoadExport(dir string) (*Export, error) {
	b, err := os.ReadFile(filepath.Join(dir, "export.json"))
	if err != nil {
		return nil, err
	}
	var e Export
	if err := json.Unmarshal(b, &e); err != nil {
		return nil, fmt.Errorf("parse export.json: %w", err)
	}
	if e.Frame == nil {
		return nil, fmt.Errorf("%s: export has no frame", dir)
	}
	e.Dir = dir
	link(e.Frame, nil)
	return &e, nil
}

func LoadLibrary(path string) (*Library, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var l Library
	if err := json.Unmarshal(b, &l); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &l, nil
}

func link(n, parent *Node) {
	n.Parent = parent
	for _, c := range n.Children {
		link(c, n)
	}
}

// Walk visits nodes depth-first; returning false skips the node's children.
func Walk(n *Node, fn func(*Node) bool) {
	if !fn(n) {
		return
	}
	for _, c := range n.Children {
		Walk(c, fn)
	}
}

// SetName resolves the component-set name of an INSTANCE node.
func (e *Export) SetName(n *Node) string {
	c, ok := e.Components[n.ComponentID]
	if !ok {
		return ""
	}
	if s, ok := e.ComponentSets[c.ComponentSetID]; ok {
		return s.Name
	}
	return c.Name
}

// Variant returns the VARIANT property values of an instance, e.g. {"Type":"Primary"}.
func (n *Node) Variant() map[string]string {
	out := map[string]string{}
	for k, p := range n.ComponentProperties {
		out[k] = fmt.Sprint(p.Value)
	}
	return out
}

// NearestInstance returns the closest INSTANCE ancestor (excluding n).
func (n *Node) NearestInstance() *Node {
	for p := n.Parent; p != nil; p = p.Parent {
		if p.Type == "INSTANCE" {
			return p
		}
	}
	return nil
}

// Text concatenates the characters of all TEXT descendants.
func (n *Node) Text() string {
	var parts []string
	Walk(n, func(x *Node) bool {
		if x.Type == "TEXT" && x.Characters != "" {
			parts = append(parts, x.Characters)
		}
		return true
	})
	return strings.Join(parts, " ")
}

// Fill returns the first solid fill color, if any.
func (n *Node) Fill() (Color, bool) {
	for _, f := range n.Fills {
		if f.Type == "SOLID" && f.Color != nil {
			return *f.Color, true
		}
	}
	return Color{}, false
}

// Rel returns the bounding box relative to the frame origin.
func (e *Export) Rel(n *Node) Box {
	b := *n.AbsoluteBoundingBox
	if f := e.Frame.AbsoluteBoundingBox; f != nil {
		b.X -= f.X
		b.Y -= f.Y
	}
	return b
}

// ParseHex parses #RGB/#RRGGBB/#AARRGGBB-less strings into a Color.
func ParseHex(s string) (Color, bool) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(s) == 3 {
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	}
	if len(s) != 6 && len(s) != 8 {
		return Color{}, false
	}
	var r, g, b, a uint64 = 0, 0, 0, 255
	var err error
	parse := func(x string) uint64 {
		var v uint64
		_, e := fmt.Sscanf(x, "%02x", &v)
		if e != nil {
			err = e
		}
		return v
	}
	r, g, b = parse(s[0:2]), parse(s[2:4]), parse(s[4:6])
	if len(s) == 8 {
		a = parse(s[6:8])
	}
	if err != nil {
		return Color{}, false
	}
	return Color{float64(r) / 255, float64(g) / 255, float64(b) / 255, float64(a) / 255}, true
}
