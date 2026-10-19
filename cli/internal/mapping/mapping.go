// Package mapping holds the Figma ↔ code component map (Phase 1): suggest candidates from the
// inventories, and check the confirmed map for gaps and token drift.
package mapping

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
	"github.com/luckylucky0001/figma-fe/cli/internal/inventory"
)

// Impl is one platform's implementation of a Figma component set.
type Impl struct {
	Component string                       `yaml:"component"`
	Import    string                       `yaml:"import,omitempty"`
	Props     map[string]map[string]string `yaml:"props,omitempty"` // axis → value → code fragment ("" = default, "-" = unsupported)
	Texts     map[string]string            `yaml:"texts,omitempty"` // prop → template over Figma TEXT layer names, e.g. "{Integer}{Decimal}|num"
	Usage     string                       `yaml:"usage,omitempty"`
	Notes     string                       `yaml:"notes,omitempty"`
}

type Entry struct {
	Figma       string   `yaml:"figma"`
	Slot        bool     `yaml:"slot,omitempty"` // children are composed by the page (Card, ActionBar)
	H5          *Impl    `yaml:"h5,omitempty"`
	Android     *Impl    `yaml:"android,omitempty"`
	Constraints []string `yaml:"constraints,omitempty"`
}

type Gap struct {
	Figma    string `yaml:"figma"`
	Platform string `yaml:"platform"`
	Decision string `yaml:"decision"`
}

type Map struct {
	Version    int     `yaml:"version"`
	Components []Entry `yaml:"components"`
	Gaps       []Gap   `yaml:"gaps,omitempty"`
}

func Load(path string) (*Map, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Map
	if err := yaml.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return &m, nil
}

func (m *Map) Entry(figmaSet string) *Entry {
	for i := range m.Components {
		if m.Components[i].Figma == figmaSet {
			return &m.Components[i]
		}
	}
	return nil
}

func (e *Entry) Impl(platform string) *Impl {
	if platform == "android" {
		return e.Android
	}
	return e.H5
}

// GapFor returns the recorded decision for an unmapped component, if any.
func (m *Map) GapFor(figmaSet, platform string) *Gap {
	for i := range m.Gaps {
		g := &m.Gaps[i]
		if g.Figma == figmaSet && (g.Platform == platform || g.Platform == "all") {
			return g
		}
	}
	return nil
}

// ---------- suggest ----------

var synonyms = [][]string{
	{"navbar", "topbar", "appbar", "header", "titlebar", "toolbar"},
	{"cell", "listitem", "row", "item", "field"},
	{"tag", "badge", "chip", "label", "pill"},
	{"price", "pricetext", "amount", "money"},
	{"productcard", "goodscard", "goodsitem", "skucard", "productitem", "commoditycard"},
	{"actionbar", "bottombar", "footerbar", "submitbar", "toolbarbottom"},
	{"noticebar", "banner", "notice", "alertbar", "tips"},
	{"card", "panel", "section", "container"},
	{"button", "btn"},
	{"divider", "separator", "line", "hr"},
	{"countdown", "timer", "countdowntimer"},
}

func norm(s string) string {
	s = strings.ToLower(s)
	for _, p := range []string{"ds", "legacy", "app", "base"} {
		s = strings.TrimPrefix(s, p)
	}
	s = strings.TrimSuffix(s, "view")
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}

func nameScore(figmaName, code string) float64 {
	a, b := norm(figmaName), norm(code)
	if a == b {
		return 1
	}
	for _, g := range synonyms {
		ia, ib := false, false
		for _, w := range g {
			ia = ia || a == w
			ib = ib || b == w
		}
		if ia && ib {
			return 0.9
		}
	}
	d := levenshtein(a, b)
	l := max(len(a), len(b))
	if l == 0 {
		return 0
	}
	return max(0, 1-float64(d)/float64(l)) * 0.8
}

func isToggle(vals []string) bool {
	return len(vals) == 2 && strings.EqualFold(vals[0], "On") && strings.EqualFold(vals[1], "Off")
}

func variantScore(set figma.ComponentSet, c inventory.Component) float64 {
	if len(set.VariantProperties) == 0 {
		return 1
	}
	matched := 0
	for axis, vals := range set.VariantProperties {
		best := 0.0
		if isToggle(vals) {
			for _, p := range c.Props {
				if strings.Contains(strings.ToLower(p.Name), strings.ToLower(axis)) || strings.Contains(strings.ToLower(axis), strings.ToLower(p.Name)) {
					best = 1
				}
			}
		}
		for _, p := range c.Props {
			if len(p.Values) == 0 {
				continue
			}
			hit := 0
			for _, v := range vals {
				for _, pv := range p.Values {
					if strings.EqualFold(v, pv) || (strings.EqualFold(v, "On") && pv == "true") || (strings.EqualFold(v, "Off") && pv == "false") {
						hit++
						break
					}
				}
			}
			best = max(best, float64(hit)/float64(len(vals)))
		}
		if best >= 0.5 {
			matched++
		}
	}
	return float64(matched) / float64(len(set.VariantProperties))
}

type Candidate struct {
	Component string  `yaml:"component"`
	Kind      string  `yaml:"kind"`
	Score     float64 `yaml:"score"`
	Name      float64 `yaml:"nameScore"`
	Variant   float64 `yaml:"variantScore"`
}

type Suggestion struct {
	Figma   string                 `yaml:"figma"`
	Status  map[string]string      `yaml:"status"` // platform → auto | review | missing
	Options map[string][]Candidate `yaml:"candidates"`
}

// Suggest ranks code components for each Figma component set, per platform.
func Suggest(lib *figma.Library, invs map[string]*inventory.Inventory) []Suggestion {
	var ids []string
	for id := range lib.ComponentSets {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return lib.ComponentSets[ids[i]].Name < lib.ComponentSets[ids[j]].Name })
	var out []Suggestion
	for _, id := range ids {
		set := lib.ComponentSets[id]
		s := Suggestion{Figma: set.Name, Status: map[string]string{}, Options: map[string][]Candidate{}}
		for _, pf := range []string{"h5", "android"} {
			inv := invs[pf]
			if inv == nil {
				continue
			}
			var cs []Candidate
			for _, c := range inv.Components {
				ns := nameScore(set.Name, c.Name)
				vs := variantScore(set, c)
				sc := 0.7*ns + 0.3*vs
				if len(set.VariantProperties) == 0 {
					sc = ns // no variant axes: the name is the only evidence
				}
				if sc >= 0.6 {
					cs = append(cs, Candidate{Component: c.Name, Kind: c.Kind, Score: round2(sc), Name: round2(ns), Variant: round2(vs)})
				}
			}
			sort.Slice(cs, func(i, j int) bool { return cs[i].Score > cs[j].Score })
			if len(cs) > 3 {
				cs = cs[:3]
			}
			switch {
			case len(cs) == 0:
				s.Status[pf] = "missing"
			case cs[0].Score >= 0.85:
				s.Status[pf] = "auto"
			default:
				s.Status[pf] = "review"
			}
			s.Options[pf] = cs
		}
		out = append(out, s)
	}
	return out
}

func round2(f float64) float64 { return float64(int(f*100+0.5)) / 100 }

func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	prev := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(ra); i++ {
		cur := make([]int, len(rb)+1)
		cur[0] = i
		for j := 1; j <= len(rb); j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev = cur
	}
	return prev[len(rb)]
}
