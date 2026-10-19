package mapping

import (
	"testing"

	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
	"github.com/luckylucky0001/figma-fe/cli/internal/inventory"
)

func TestNameScore(t *testing.T) {
	cases := []struct {
		figma, code string
		min, max    float64
	}{
		{"Button", "Button", 1, 1},
		{"Button", "DsButton", 1, 1},
		{"NavBar", "DsTopBar", 0.9, 0.9}, // synonym
		{"Cell", "DsListItem", 0.9, 0.9}, // synonym
		{"Tag", "LegacyTagView", 1, 1},   // prefix/suffix stripped
		{"CountDown", "Button", 0, 0.45}, // unrelated
	}
	for _, c := range cases {
		got := nameScore(c.figma, c.code)
		if got < c.min || got > c.max {
			t.Errorf("nameScore(%q,%q)=%.2f want [%.2f,%.2f]", c.figma, c.code, got, c.min, c.max)
		}
	}
}

func TestVariantScore(t *testing.T) {
	set := figma.ComponentSet{Name: "Tag", VariantProperties: map[string][]string{
		"Tone":    {"Primary", "Success", "Warning", "Danger"},
		"Variant": {"Light", "Solid"},
	}}
	full := inventory.Component{Props: []inventory.Prop{
		{Name: "type", Values: []string{"primary", "success", "warning", "danger"}},
		{Name: "variant", Values: []string{"light", "solid"}},
	}}
	half := inventory.Component{Props: []inventory.Prop{{Name: "tone", Values: []string{"Primary", "Success", "Warning", "Danger"}}}}
	if s := variantScore(set, full); s != 1 {
		t.Errorf("full=%v want 1", s)
	}
	if s := variantScore(set, half); s != 0.5 {
		t.Errorf("half=%v want 0.5 (Solid variant missing)", s)
	}
	toggle := figma.ComponentSet{VariantProperties: map[string][]string{"Label": {"On", "Off"}}}
	if s := variantScore(toggle, inventory.Component{Props: []inventory.Prop{{Name: "label", Type: "string"}}}); s != 1 {
		t.Errorf("toggle axis should match optional prop of the same name, got %v", s)
	}
}

func TestSuggestFlagsMissing(t *testing.T) {
	lib := &figma.Library{ComponentSets: map[string]figma.ComponentSet{
		"1": {Name: "Button", VariantProperties: map[string][]string{"Type": {"Primary", "Default"}}},
		"2": {Name: "CountDown"},
	}}
	inv := &inventory.Inventory{Components: []inventory.Component{
		{Name: "DsButton", Props: []inventory.Prop{{Name: "style", Values: []string{"Primary", "Default"}}}},
	}}
	got := Suggest(lib, map[string]*inventory.Inventory{"android": inv})
	st := map[string]string{}
	for _, s := range got {
		st[s.Figma] = s.Status["android"]
	}
	if st["Button"] != "auto" || st["CountDown"] != "missing" {
		t.Fatalf("status=%v", st)
	}
}
