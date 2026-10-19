package diff

import (
	"image/color"
	"testing"
)

func TestYIQ(t *testing.T) {
	w := color.RGBA{255, 255, 255, 255}
	if d := yiq(w, w); d != 0 {
		t.Errorf("identical colours delta=%v", d)
	}
	thr := 0.1 * 0.1 * 35215.0
	if d := yiq(color.RGBA{0xFF, 0x8F, 0x1F, 255}, color.RGBA{0xFF, 0x88, 0x00, 255}); d <= 0 {
		t.Error("token drift #FF8F1F vs #FF8800 must differ")
	}
	if d := yiq(w, color.RGBA{0, 0, 0, 255}); d < thr {
		t.Error("black vs white must exceed threshold")
	}
}

func TestSeverity(t *testing.T) {
	if sev(3) != "warn" || sev(5) != "error" {
		t.Error("severity bands")
	}
}
