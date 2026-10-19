package reuse

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
)

func TestConceptOf(t *testing.T) {
	cases := map[string]string{
		"StatusTag": "tag", "van-tag": "tag", "LegacyTagView": "tag", "OrderStatusChip": "tag",
		"InfoRow": "cell", "QMUICommonListItemView": "cell", "ListItem": "cell",
		"PageHeader": "navbar", "QMUITopBar": "navbar", "SimpleHeader": "navbar",
		"RadioButton": "switch", "QMUIRoundButton": "button", "SubmitBtn": "button",
		"SwipeBackLayout": "", "QMUICheckBox": "switch", "Swiper": "carousel",
		"CountDown": "countdown", "GoodsItem": "goods", "SectionCard": "card",
		"按钮/主按钮": "button", "订单状态标签": "tag", "底部操作栏": "actionbar", "倒计时": "countdown", "商品卡片": "goods", "Tag / 状态": "tag",
	}
	for in, want := range cases {
		if got := ConceptOf(in); got != want {
			t.Errorf("ConceptOf(%q)=%q want %q", in, got, want)
		}
	}
}

func TestClonesCollapseSiblingRuns(t *testing.T) {
	a := `<div><span>a</span><span>b</span><img/><p><i></i><b></b></p></div>`
	b := `<div><span>x</span><span>y</span><span>z</span><img/><p><i></i><b></b></p></div>`
	cs := findClones(map[string]string{"a.vue": a, "b.vue": b}, func(s string) []*mnode { return filterJSX(parseMarkup(s)) }, 5)
	if len(cs) != 1 || len(cs[0].Locations) != 2 {
		t.Fatalf("clones=%+v", cs)
	}
}

func write(t *testing.T, root, rel, content string) {
	p := filepath.Join(root, rel)
	_ = os.MkdirAll(filepath.Dir(p), 0o755)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestScanVueVantProject(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"^3.4.0","vant":"^4.8.0"}}`)
	write(t, root, "src/components/NavBar.vue", `<template><div class="nav"><span>‹</span><b>{{title}}</b></div></template><script setup>defineProps({ title: String })</script>`)
	write(t, root, "src/views/Order.vue", `<template><nav-bar title="订单"/><van-cell title="运费" value="免运费"/><button class="btn">去支付</button></template><style>.btn{color:#1baeae}</style>`)
	write(t, root, "src/styles/var.less", `@primary: #1baeae;`)
	r := ScanH5(root)
	if !strings.Contains(strings.Join(r.Stack, ","), "Vant") {
		t.Errorf("stack=%v", r.Stack)
	}
	got := map[string]Asset{}
	for _, a := range r.Assets {
		got[a.Kind+"|"+a.Name] = a
	}
	if a := got["ui-lib|Cell"]; a.Usage != 1 || !strings.Contains(a.Example, "van-cell") {
		t.Errorf("van-cell usage=%+v", a)
	}
	if a := got["business|NavBar"]; a.Usage != 1 || len(a.Props) != 1 {
		t.Errorf("business NavBar=%+v", a)
	}
	if r.Metrics.RawPrimitives != 1 || len(r.Findings) != 1 {
		t.Errorf("raw=%d findings=%+v", r.Metrics.RawPrimitives, r.Findings)
	}
	if len(r.Colors) == 0 || r.Colors[0].Token != "@primary" {
		t.Errorf("colors=%+v", r.Colors)
	}
	// business NavBar duplicates Vant NavBar (unused): the "own wheel" case
	found := false
	for _, c := range r.Clusters {
		found = found || c.Concept == "navbar" && strings.Contains(strings.Join(c.Members, ","), "项目未使用")
	}
	if !found {
		t.Errorf("clusters=%+v", r.Clusters)
	}
}

func TestScanAndroidAndMatch(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/build.gradle", `dependencies { implementation 'com.google.android.material:material:1.6.0' }`)
	write(t, root, "app/src/main/java/com/x/StatusTagView.kt", `package com.x
class StatusTagView @JvmOverloads constructor(c: Context) : AppCompatTextView(c) {}`)
	write(t, root, "app/src/main/res/values/attrs.xml", `<resources><declare-styleable name="StatusTagView"><attr name="tagTone" format="enum"><enum name="warning" value="1"/></attr></declare-styleable></resources>`)
	write(t, root, "app/src/main/res/values/colors.xml", `<resources><color name="brand">#FF4D1A</color></resources>`)
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><com.x.StatusTagView android:id="@+id/t"/><TextView android:textColor="#FF4D1A"/><com.google.android.material.button.MaterialButton/></LinearLayout>`)
	r := ScanAndroid(root)
	got := map[string]Asset{}
	for _, a := range r.Assets {
		got[a.Name] = a
	}
	if a := got["StatusTagView"]; a.Kind != "custom-view" || a.Usage != 1 || len(a.Props) != 1 {
		t.Errorf("custom view=%+v", a)
	}
	if got["MaterialButton"].Usage != 1 {
		t.Errorf("material=%+v", got["MaterialButton"])
	}
	if len(r.Colors) != 1 || r.Colors[0].Token != "@color/brand" {
		t.Errorf("colors=%+v", r.Colors)
	}
	exp := &figma.Export{
		Frame:         &figma.Node{Type: "FRAME", Children: []*figma.Node{{ID: "1:2", Type: "INSTANCE", ComponentID: "c1"}, {ID: "1:3", Type: "INSTANCE", ComponentID: "c2"}}},
		Components:    map[string]figma.Component{"c1": {ComponentSetID: "s1"}, "c2": {ComponentSetID: "s2"}},
		ComponentSets: map[string]figma.ComponentSet{"s1": {Name: "Tag"}, "s2": {Name: "CountDown"}},
	}
	m := Match(exp, r)
	if m.Items[0].Verdict == "new" || m.Items[0].Candidates[0].Name != "StatusTagView" {
		t.Errorf("Tag → %+v", m.Items[0])
	}
	if m.Items[1].Verdict != "new" {
		t.Errorf("CountDown should need a new component here: %+v", m.Items[1])
	}
}

func TestCheckFocusHandBuilt(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"^3.4.0","vant":"^4.8.0"}}`)
	write(t, root, "src/views/Order.vue", `<template><van-card title="x"/><van-action-bar/></template>`)
	write(t, root, "src/views/New.vue", `<template><div class="goods-card">x</div><van-button>ok</van-button></template>
<style lang="less">
@brand: #FF4D1A;
.goods-card { color: @brand; }
.tip { color: #1BAEAE; }
</style>`)
	write(t, root, "src/styles/var.less", `@primary: #1BAEAE;`)
	pc := ScanH5(root).CheckFocus("src/views/New.vue")
	if pc.Pass || len(pc.HandBuilt) != 1 || pc.HandBuilt[0].Severity != "error" || !strings.Contains(pc.HandBuilt[0].Hint, "Card") {
		t.Errorf("handBuilt=%+v pass=%v", pc.HandBuilt, pc.Pass)
	}
	if pc.LocalTokens != 1 || len(pc.Colors) != 1 || !strings.Contains(pc.Colors[0].Hint, "@primary") {
		t.Errorf("localTokens=%d colors=%+v", pc.LocalTokens, pc.Colors)
	}
}

func TestVueAliasAndFixedBottom(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"^3.4.0","vant":"^4.8.0"}}`)
	write(t, root, "src/components/NavBar.vue", "<template><div class=\"nav-bar\"><i>首页</i><i>我的</i></div></template>\n<style>.nav-bar { position: fixed;\n left: 0;\n bottom: 0; }</style>")
	write(t, root, "src/components/SimpleHeader.vue", `<template><header><b>{{name}}</b></header></template>`)
	write(t, root, "src/views/A.vue", "<template><s-header :name=\"'详情'\"></s-header><nav-bar/></template>\n<script>import sHeader from '@/components/SimpleHeader.vue'\nimport navBar from '@/components/NavBar.vue'</script>")
	r := ScanH5(root)
	got := map[string]Asset{}
	for _, a := range r.Assets {
		got[a.Name] = a
	}
	if got["SimpleHeader"].Usage != 1 {
		t.Errorf("aliased SimpleHeader usage=%d", got["SimpleHeader"].Usage)
	}
	if got["NavBar"].Concept != "tabs" {
		t.Errorf("fixed-bottom NavBar concept=%q want tabs (it is a bottom tab bar)", got["NavBar"].Concept)
	}
}

func TestLibrariesDetection(t *testing.T) {
	cases := []struct {
		name, pkg, file, src, lib, comp string
	}{
		{"nutui", `{"dependencies":{"vue":"3","@nutui/nutui":"4"}}`, "src/views/A.vue", `<template><nut-button type="primary">去支付</nut-button><nut-price :price="269"/></template>`, "NutUI", "Price"},
		{"antd-mobile", `{"dependencies":{"react":"18","antd-mobile":"5"}}`, "src/pages/A.tsx", "import { Button, NoticeBar } from 'antd-mobile';\nexport default function A() { return <div><NoticeBar content=\"x\"/><Button color=\"primary\">ok</Button></div> }", "antd-mobile", "NoticeBar"},
		{"tdesign", `{"dependencies":{"vue":"3","tdesign-mobile-vue":"1"}}`, "src/views/A.vue", `<template><t-cell title="运费" note="免运费"/><t-count-down :time="1000"/></template>`, "TDesign Mobile", "CountDown"},
	}
	for _, c := range cases {
		root := t.TempDir()
		write(t, root, "package.json", c.pkg)
		write(t, root, c.file, c.src)
		r := ScanH5(root)
		found := false
		for _, a := range r.Assets {
			if a.Kind == "ui-lib" && a.Source == c.lib && a.Name == c.comp && a.Usage == 1 {
				found = true
			}
		}
		if !found {
			t.Errorf("%s: %s/%s not detected; assets=%+v", c.name, c.lib, c.comp, r.Assets)
		}
	}
}

func TestComposeProject(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/build.gradle.kts", `dependencies { implementation("androidx.compose.material3:material3") }`)
	write(t, root, "app/src/main/java/x/ui/StatusBadge.kt", "package x.ui\n@Composable\nfun StatusBadge(text: String, tone: Tone = Tone.Warning) {}\n@Composable\nprivate fun Inner() {}")
	write(t, root, "app/src/main/java/x/OrderScreen.kt", "import androidx.compose.material3.Button\n@Composable\nfun OrderScreen() { StatusBadge(\"待发货\"); Button(onClick = {}) { Text(\"ok\") } }")
	r := ScanAndroid(root)
	got := map[string]Asset{}
	for _, a := range r.Assets {
		got[a.Kind+"|"+a.Name] = a
	}
	if a := got["composable|StatusBadge"]; a.Usage != 1 || a.Concept != "tag" {
		t.Errorf("StatusBadge=%+v", a)
	}
	if _, ok := got["composable|Inner"]; ok {
		t.Error("private composables are not reusable assets")
	}
	if got["ui-lib|Button"].Usage != 1 {
		t.Errorf("material3 Button=%+v", got["ui-lib|Button"])
	}
}

func TestUViewUniApp(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"3","uview-plus":"3"}}`)
	write(t, root, "pages/order/detail.vue", `<template><u-navbar title="订单详情"/><u-count-down :time="1000"/><u-tag text="待发货" type="warning"/></template>`)
	r := ScanH5(root)
	n := 0
	for _, a := range r.Assets {
		if a.Source == "uView" && a.Usage == 1 {
			n++
		}
	}
	if n != 3 {
		t.Errorf("uView assets used=%d assets=%+v", n, r.Assets)
	}
}

func TestCheckFocusAndroidLayout(t *testing.T) {
	root := t.TempDir()
	write(t, root, "lib/src/main/java/com/x/CellView.java", "package com.x;\npublic class CellView extends LinearLayout {}")
	write(t, root, "lib/src/main/java/com/x/TagView.java", "package com.x;\npublic class TagView extends TextView {}")
	write(t, root, "lib/src/main/java/com/x/RoundButton.java", "package com.x;\npublic class RoundButton extends Button {}")
	write(t, root, "lib/src/main/java/com/x/compose/Cell.kt", "package com.x.compose\n@Composable\nfun Cell(title: String) {}")
	write(t, root, "app/src/main/res/values/colors.xml", `<resources><color name="white">#FFFFFF</color></resources>`)
	write(t, root, "app/src/main/res/layout/old.xml", `<LinearLayout><com.x.CellView/><com.x.CellView/><com.x.CellView/></LinearLayout>`)
	write(t, root, "app/src/main/res/layout/page.xml", `<LinearLayout>
  <com.x.RoundButton android:id="@+id/status_tag"/>
  <LinearLayout android:id="@+id/address_cell"/>
  <ImageView android:src="@drawable/a"/>
</LinearLayout>`)
	write(t, root, "app/src/main/res/values/page_colors.xml", `<resources>
<color name="page_bg">#FFFFFF</color>
<color name="page_red">#F53F3F</color>
</resources>`)
	r := ScanAndroid(root)
	pc := r.CheckFocus("app/src/main/res/layout/page.xml,app/src/main/res/values/page_colors.xml")
	// the id on <com.x.RoundButton> is a deliberate component choice, not a hand-built tag; the one on a bare LinearLayout is, and the
	// suggestion must be the View (an XML layout cannot host the composable)
	if len(pc.HandBuilt) != 1 || !strings.Contains(pc.HandBuilt[0].Text, "address_cell") || !strings.Contains(pc.HandBuilt[0].Hint, "CellView") {
		t.Errorf("handBuilt=%+v", pc.HandBuilt)
		for _, a := range r.Assets {
			t.Logf("%s %s %s %d %v", a.Name, a.Kind, a.Concept, a.Usage, a.UsageFiles)
		}
	}
	for _, f := range pc.MissedReuse {
		if f.Severity != "info" {
			t.Errorf("raw ImageView must not fail the page: %+v", f)
		}
	}
	if pc.LocalTokens != 2 || len(pc.Colors) != 1 || !strings.Contains(pc.Colors[0].Hint, "@color/white") {
		t.Errorf("localTokens=%d colors=%+v", pc.LocalTokens, pc.Colors)
	}
}

func TestCompatVersionAndTheme(t *testing.T) {
	root := t.TempDir()
	write(t, root, "buildSrc/Dep.kt", `const val material = "com.google.android.material:material:1.4.0"`)
	write(t, root, "app/build.gradle", `dependencies { implementation(Dep.material) }`)
	write(t, root, "app/src/main/AndroidManifest.xml", `<manifest><application android:theme="@style/AppTheme"/></manifest>`)
	write(t, root, "app/src/main/res/values/themes.xml", `<resources>
<style name="AppTheme" parent="Base.NoActionBar"/>
<style name="Base"  parent="Theme.AppCompat.Light"/>
<style name="Base.NoActionBar"/>
</resources>`)
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><TextView/><com.google.android.material.button.MaterialButton/></LinearLayout>`)
	r := ScanAndroid(root)
	bad := map[string]string{}
	for _, a := range r.LibAll {
		if !usable(&a) {
			bad[a.Name] = a.Caveat
		}
	}
	if !strings.Contains(bad["MaterialDivider"], "1.5.0") || !strings.Contains(bad["Chip"], "AppCompat") {
		t.Errorf("caveats=%v", bad)
	}
	if _, ok := bad["MaterialButton"]; ok {
		t.Errorf("MaterialButton is used by the project, so it evidently works: %v", bad["MaterialButton"])
	}
	exp := &figma.Export{
		Frame:         &figma.Node{Type: "FRAME", Children: []*figma.Node{{ID: "1:2", Type: "INSTANCE", ComponentID: "c1"}, {ID: "1:3", Type: "INSTANCE", ComponentID: "c2"}}},
		Components:    map[string]figma.Component{"c1": {ComponentSetID: "s1"}, "c2": {ComponentSetID: "s2"}},
		ComponentSets: map[string]figma.ComponentSet{"s1": {Name: "Divider"}, "s2": {Name: "Tag"}},
	}
	for _, it := range Match(exp, r).Items {
		for _, c := range it.Candidates {
			if c.Name == "MaterialDivider" || c.Name == "Chip" && c.Score >= 0.8 {
				t.Errorf("%s → %+v", it.Figma, c)
			}
		}
	}
}

func TestConventionsAndVantVersion(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"^2.6.0","vant":"^2.12.0"}}`)
	for i, n := range []string{"A", "B", "C"} {
		write(t, root, "src/views/"+n+".vue", `<template><div class="head van-hairline--bottom">`+itoa(i)+`</div><van-button>x</van-button></template>`)
	}
	r := ScanH5(root)
	if len(r.Conventions) != 1 || r.Conventions[0].Usage != 3 {
		t.Fatalf("conventions=%+v", r.Conventions)
	}
	for _, a := range r.LibAll {
		if a.Name == "DatePicker" && usable(&a) {
			t.Errorf("Vant 2 has no DatePicker: %+v", a)
		}
	}
}
