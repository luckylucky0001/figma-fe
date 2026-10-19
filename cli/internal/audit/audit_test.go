package audit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/luckylucky0001/figma-fe/cli/internal/inventory"
	"github.com/luckylucky0001/figma-fe/cli/internal/plan"
)

func TestAuditAndroid(t *testing.T) {
	dir := t.TempDir()
	src := `package x
import androidx.compose.material3.Button
import androidx.compose.material3.Text
@Composable fun Page() {
  DsTopBar(title = "订单详情", modifier = Modifier.figmaId("1:2"))
  Button(onClick = {}) { Text("提醒发货") }
  Box(Modifier.background(Color(0xFFFF4D1A)).padding(12.dp))
  LegacyCellView()
}
@Composable private fun PriceRow() {}`
	_ = os.WriteFile(filepath.Join(dir, "Page.kt"), []byte(src), 0o644)
	inv := &inventory.Inventory{Components: []inventory.Component{{Name: "DsTopBar"}, {Name: "DsButton"}}}
	p := &plan.Plan{Counts: map[string]int{"DsTopBar": 1, "DsButton": 1}}
	r, err := Run("android", dir, dir, p, inv)
	if err != nil {
		t.Fatal(err)
	}
	if r.ReuseRate != 0.5 {
		t.Errorf("reuse=%v want 0.5", r.ReuseRate)
	}
	rules := map[string]int{}
	for _, f := range r.Findings {
		rules[f.Rule]++
	}
	for _, want := range []string{"direct-material3", "hardcoded-color", "legacy-view", "lookalike", "raw-number", "missing-reuse"} {
		if rules[want] == 0 {
			t.Errorf("missing finding %s: %v", want, rules)
		}
	}
	if rules["direct-material3"] != 1 { // Text is an allowed primitive
		t.Errorf("material3 findings=%d want 1", rules["direct-material3"])
	}
	if r.FigmaIDs != 1 || r.Violations != 3 {
		t.Errorf("figmaIds=%d violations=%d", r.FigmaIDs, r.Violations)
	}
}

func TestAuditH5(t *testing.T) {
	dir := t.TempDir()
	src := `import { Cell } from '@ds';
export default () => <div><Cell title="运费" data-figma-id="1:3" /><button style={{color:'#ff4d1a'}}>x</button><span className="my-tag">t</span></div>`
	_ = os.WriteFile(filepath.Join(dir, "index.tsx"), []byte(src), 0o644)
	inv := &inventory.Inventory{Components: []inventory.Component{{Name: "Cell"}, {Name: "Tag"}}}
	p := &plan.Plan{Counts: map[string]int{"Cell": 1, "Tag": 1}}
	r, _ := Run("h5", dir, dir, p, inv)
	if r.ReuseRate != 0.5 || r.Violations != 2 || r.FigmaIDs != 1 {
		t.Errorf("reuse=%v violations=%d ids=%d findings=%+v", r.ReuseRate, r.Violations, r.FigmaIDs, r.Findings)
	}
}
