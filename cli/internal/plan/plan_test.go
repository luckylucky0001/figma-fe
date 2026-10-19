package plan

import (
	"testing"

	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
)

func TestFillTemplates(t *testing.T) {
	price := &figma.Node{Type: "INSTANCE", Children: []*figma.Node{
		{Type: "TEXT", Name: "Symbol", Characters: "¥"},
		{Type: "TEXT", Name: "Integer", Characters: "269"},
		{Type: "TEXT", Name: "Decimal", Characters: ".00"},
	}}
	if v, num := fill(price, "{Integer}{Decimal}|num"); v != "269.00" || !num {
		t.Errorf("price = %q %v", v, num)
	}
	cnt := &figma.Node{Children: []*figma.Node{{Type: "TEXT", Name: "Count", Characters: "x12"}}}
	if v, _ := fill(cnt, "{Count}|int"); v != "12" {
		t.Errorf("count = %q", v)
	}
	if v, _ := fill(cnt, "{Missing}"); v != "" {
		t.Errorf("missing layer should yield empty, got %q", v)
	}
}

func TestRef(t *testing.T) {
	cases := map[[2]string]string{
		{"h5", "color.text-primary"}:      "var(--ds-color-text-primary)",
		{"android", "color.text-primary"}: "DsColors.TextPrimary",
		{"android", "font.title"}:         "DsTypography.Title",
		{"android", "space.md"}:           "DsSpacing.Md",
	}
	for in, want := range cases {
		if got := Ref(in[0], in[1]); got != want {
			t.Errorf("Ref(%v)=%q want %q", in, got, want)
		}
	}
}
