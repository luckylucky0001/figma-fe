package inventory

import "testing"

func TestParseTSProps(t *testing.T) {
	src := `export interface ButtonProps extends Omit<X, 'type'> {
  /** Visual type. */
  type?: 'primary' | 'default' | 'text';
  block?: boolean;
  onClick?: () => void;
}`
	ps := parseTSProps(src, "ButtonProps")
	if len(ps) != 3 {
		t.Fatalf("got %d props: %+v", len(ps), ps)
	}
	if ps[0].Name != "type" || len(ps[0].Values) != 3 || ps[0].Doc != "Visual type." {
		t.Errorf("type prop = %+v", ps[0])
	}
	if ps[1].Values[0] != "true" {
		t.Errorf("boolean should expose true/false: %+v", ps[1])
	}
}

func TestDocBefore(t *testing.T) {
	src := "/**\n * Button — pill button.\n *\n * @example\n * <Button type=\"primary\">去支付</Button>\n */\nexport function Button() {}"
	d, ex := docBefore(src, len(src)-len("export function Button() {}"))
	if d != "Button — pill button." || ex != "<Button type=\"primary\">去支付</Button>" {
		t.Errorf("desc=%q example=%q", d, ex)
	}
}

func TestSplitTopAndKebab(t *testing.T) {
	parts := splitTop(`text: String, onClick: () -> Unit, size: Pair<Int, Int> = Pair(1, 2)`, ',')
	if len(parts) != 3 {
		t.Errorf("splitTop=%q", parts)
	}
	if kebab("TextPrimary") != "text-primary" || kebab("BgPage") != "bg-page" {
		t.Error("kebab")
	}
}

func TestDocBeforeIgnoresPropDocs(t *testing.T) {
	src := "export interface P {\n  /** prop doc */\n  children?: ReactNode;\n}\n\n/**\n * ActionBar — fixed bar.\n */\nexport function ActionBar() {}"
	d, _ := docBefore(src, len(src)-len("export function ActionBar() {}"))
	if d != "ActionBar — fixed bar." {
		t.Errorf("desc=%q", d)
	}
}
