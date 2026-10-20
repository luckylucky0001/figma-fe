package reuse

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/luckylucky0001/figma-fe/cli/internal/mapping"

	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
)

func TestConceptOf(t *testing.T) {
	cases := map[string]string{
		"StatusTag": "tag", "van-tag": "tag", "LegacyTagView": "tag", "OrderStatusChip": "tag",
		"InfoRow": "cell", "QMUICommonListItemView": "cell", "ListItem": "cell",
		"PageHeader": "navbar", "QMUITopBar": "navbar", "SimpleHeader": "navbar",
		"RadioButton": "radio", "QMUIRoundButton": "button", "SubmitBtn": "button",
		"SwipeBackLayout": "", "QMUICheckBox": "checkbox", "Swiper": "carousel",
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
	if got["NavBar"].Concept != "tabbar" {
		t.Errorf("fixed-bottom NavBar concept=%q want tabbar (it is a bottom tab bar)", got["NavBar"].Concept)
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

func TestReactCompoundMembers(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"react":"^18.2.0","antd-mobile":"^5.38.0"}}`)
	write(t, root, "src/App.tsx", `import { BrowserRouter } from 'react-router-dom'
export default function App() { return <BrowserRouter><Profile/></BrowserRouter> }`)
	write(t, root, "src/pages/Profile.tsx", `import { List, Form } from 'antd-mobile'
export default function Profile() {
  return <List><List.Item extra="x">a</List.Item><Form><Form.Item label="n">b</Form.Item></Form></List>
}`)
	r := ScanH5(root)
	got := map[string]Asset{}
	for _, a := range r.Assets {
		got[a.Name] = a
	}
	if got["ListItem"].Usage != 1 || got["ListItem"].Concept != "cell" {
		t.Errorf("ListItem=%+v", got["ListItem"])
	}
	if got["FormItem"].Concept != "" {
		t.Errorf("Form.Item is not a list cell: %+v", got["FormItem"])
	}
	if _, ok := got["App"]; ok {
		t.Errorf("app shell listed as a component")
	}
}

func TestPropHints(t *testing.T) {
	vant := Asset{Name: "Button", Kind: "ui-lib", Source: "Vant"}
	got := strings.Join(PropHints([]string{"Size=Medium", "Type=Primary", "Type=Default"}, vant), " | ")
	if !strings.Contains(got, `size="normal"`) || !strings.Contains(got, `type="primary"`) || !strings.Contains(got, `type="default"`) {
		t.Errorf("vant button: %s", got)
	}
	cell := Asset{Name: "Cell", Kind: "ui-lib", Source: "Vant"}
	if got := strings.Join(PropHints([]string{"Arrow=On", "Arrow=Off"}, cell), ""); !strings.Contains(got, "is-link") {
		t.Errorf("vant cell: %s", got)
	}
	enums := kotlinEnums("enum class Tone(val c: Color) { Info(Color(0xFF1677FF)), Warning(Color(0xFFFF8F1F));\n val x = 1 }")
	if strings.Join(enums["Tone"], ",") != "Tone.Info,Tone.Warning" {
		t.Errorf("kotlin enums: %v", enums)
	}
	banner := Asset{Name: "DsBanner", Kind: "composable", Enums: map[string][]string{"tone": enums["Tone"]}, Props: []string{"text", "showChevron"}}
	got = strings.Join(PropHints([]string{"Tone=Warning", "Arrow=On", "Arrow=Off"}, banner), " | ")
	if !strings.Contains(got, "tone = Tone.Warning") || !strings.Contains(got, "showChevron") {
		t.Errorf("compose: %s", got)
	}
	if p := tsProps("export interface CellProps extends X {\n  title: string;\n  arrow?: boolean;\n}\n"); strings.Join(p, ",") != "title,arrow" {
		t.Errorf("ts props: %v", p)
	}
}

func TestReinventionPurposeAndDecision(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"^3.4.0"}}`)
	write(t, root, "src/components/FilterChip.vue", `<template><span/></template>`)
	write(t, root, "src/components/BaseCell.vue", `<template><div/></template>`)
	write(t, root, "src/views/Old.vue", `<template><FilterChip/><FilterChip/><BaseCell/><BaseCell/></template>`)
	write(t, root, "src/components/StatusChip.vue", `<template><span/></template>`)
	write(t, root, "src/components/OrderCell.vue", `<template><div/></template>`)
	r := ScanH5(root)
	pc := r.CheckFocus("src/components/StatusChip.vue,src/components/OrderCell.vue")
	sev := map[string]string{}
	for _, f := range pc.Reinvented {
		sev[f.Text] = f.Severity
	}
	// StatusChip vs FilterChip: different purpose → only a look; OrderCell vs the generic BaseCell → reinvention
	if sev["StatusChip"] != "warn" || sev["OrderCell"] != "error" {
		t.Fatalf("reinvented=%+v", pc.Reinvented)
	}
	Confirmed = &mapping.Map{Gaps: []mapping.Gap{{Figma: "Cell", Platform: "all", Decision: "新建 OrderCell：订单行需要右侧金额对齐"}}}
	defer func() { Confirmed = nil }()
	pc = r.CheckFocus("src/components/StatusChip.vue,src/components/OrderCell.vue")
	for _, f := range pc.Reinvented {
		if f.Text == "OrderCell" {
			t.Errorf("a recorded decision should settle OrderCell: %+v", f)
		}
	}
}

func TestH5WrapperBypass(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"^3.4.0","vant":"^4.8.0"}}`)
	write(t, root, "src/components/base/BaseButton.vue", `<template><van-button :type="type" round><slot/></van-button></template>
<script setup lang="ts">defineProps<{ type?: 'primary' | 'default' }>()</script>`)
	write(t, root, "src/views/A.vue", `<template><BaseButton type="primary">ok</BaseButton><BaseButton>x</BaseButton><BaseButton>y</BaseButton></template>`)
	write(t, root, "src/views/New.vue", `<template><van-button type="primary">pay</van-button></template>`)
	r := ScanH5(root)
	var bypass []Finding
	for _, f := range r.Findings {
		if f.Rule == "bypass" {
			bypass = append(bypass, f)
		}
	}
	if len(bypass) != 1 || bypass[0].File != "src/views/New.vue" || !strings.Contains(bypass[0].Hint, "BaseButton") {
		t.Fatalf("bypass=%+v", bypass)
	}
	if pc := r.CheckFocus("src/views/New.vue"); pc.Pass {
		t.Errorf("bypassing the wrapper should fail the page: %+v", pc.MissedReuse)
	}
	exp := &figma.Export{
		Frame:         &figma.Node{Type: "FRAME", Children: []*figma.Node{{ID: "1:2", Type: "INSTANCE", ComponentID: "c1"}}},
		Components:    map[string]figma.Component{"c1": {ComponentSetID: "s1"}},
		ComponentSets: map[string]figma.ComponentSet{"s1": {Name: "Button"}},
	}
	if c := Match(exp, r).Items[0].Candidates[0]; c.Name != "BaseButton" {
		t.Errorf("Button → %+v", c)
	}
}

func TestHandBuiltRowsByStructure(t *testing.T) {
	root := t.TempDir()
	write(t, root, "lib/src/main/java/com/x/CellView.java", "package com.x;\npublic class CellView extends LinearLayout {}")
	write(t, root, "app/src/main/res/layout/old.xml", `<LinearLayout><com.x.CellView/><com.x.CellView/></LinearLayout>`)
	row := `<LinearLayout android:layout_width="match_parent" android:orientation="horizontal"><TextView android:text="运费"/><TextView android:text="免运费"/></LinearLayout>`
	write(t, root, "app/src/main/res/layout/page.xml", `<LinearLayout android:orientation="vertical">
`+row+`
`+row+`
<LinearLayout android:orientation="horizontal"><Button/><Button/></LinearLayout>
</LinearLayout>`)
	pc := ScanAndroid(root).CheckFocus("app/src/main/res/layout/page.xml")
	if len(pc.HandBuilt) != 1 || !strings.Contains(pc.HandBuilt[0].Hint, "2 行") || !strings.Contains(pc.HandBuilt[0].Hint, "CellView") || pc.Pass {
		t.Errorf("handBuilt=%+v", pc.HandBuilt)
	}
}

func TestHandBuiltRowsH5(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"^3.4.0","vant":"^4.8.0"}}`)
	write(t, root, "src/views/Old.vue", `<template><van-cell title="a"/><van-cell title="b"/><van-cell title="c"/></template>`)
	write(t, root, "src/views/New.vue", `<template><div>
  <div><span>运费</span><span>免运费</span></div>
  <div><span>优惠</span><span>-¥30</span></div>
  <div><span>实付</span><span>¥269</span></div>
</div></template>`)
	pc := ScanH5(root).CheckFocus("src/views/New.vue")
	if len(pc.HandBuilt) != 1 || !strings.Contains(pc.HandBuilt[0].Hint, "3 行") || !strings.Contains(pc.HandBuilt[0].Hint, "Cell") {
		t.Errorf("handBuilt=%+v", pc.HandBuilt)
	}
}

func TestSameNameAcrossLibraries(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/build.gradle", `dependencies { implementation "androidx.compose.material3:material3:1.2.0" }`)
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><Button/><Button/></LinearLayout>`)
	write(t, root, "app/src/main/java/x/S.kt", "import androidx.compose.material3.Button\n@Composable\nfun S() { Button(onClick = {}) {} }")
	r := ScanAndroid(root)
	use := map[string]int{}
	for _, a := range r.Assets {
		if a.Name == "Button" {
			use[a.Source] = a.Usage
		}
	}
	if use["Android SDK widgets"] != 2 || use["Compose Material3"] != 1 {
		t.Errorf("XML <Button> and Compose Button() must be counted separately: %v", use)
	}
}

func TestThirdPartyComponentPackage(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"^2.6.0","vant":"^2.12.0","@chenfengyuan/vue-countdown":"^1.1.2"}}`)
	write(t, root, "src/main.js", "import Vue from 'vue'\nimport VueCountdown from '@chenfengyuan/vue-countdown'\nVue.component(VueCountdown.name, VueCountdown)\n")
	write(t, root, "src/views/A.vue", `<template><div><countdown :time="60000"/></div></template>`)
	write(t, root, "src/views/B.vue", `<template><div><countdown :time="1000"/></div></template>`)
	r := ScanH5(root)
	for _, a := range r.Assets {
		if a.Name == "VueCountdown" {
			if a.Usage != 2 || a.Concept != "countdown" || a.Source != "@chenfengyuan/vue-countdown" {
				t.Errorf("%+v", a)
			}
			return
		}
	}
	t.Errorf("third-party countdown not found: %+v", r.Assets)
}

func TestInconsistentChoice(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"^3.4.0","vant":"^4.8.0"}}`)
	write(t, root, "src/components/SimpleHeader.vue", `<template><header/></template>`)
	for _, n := range []string{"A", "B", "C"} {
		write(t, root, "src/views/"+n+".vue", `<template><SimpleHeader/></template>`)
	}
	write(t, root, "src/views/New.vue", `<template><van-nav-bar title="x"/></template>`)
	pc := ScanH5(root).CheckFocus("src/views/New.vue")
	if len(pc.Inconsistent) != 1 || !strings.Contains(pc.Inconsistent[0].Hint, "SimpleHeader") || !pc.Pass {
		t.Errorf("inconsistent=%+v pass=%v", pc.Inconsistent, pc.Pass)
	}
}

func TestTaroJSXInPlainJS(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"@tarojs/taro":"^3.0.0","taro-ui":"^2.2.0","prop-types":"^15.0.0"}}`)
	write(t, root, "src/pages/index/index.js", "import Taro from '@tarojs/taro'\nimport PropTypes from 'prop-types'\nimport { AtButton, AtList, AtListItem } from 'taro-ui'\nexport default class Index extends Taro.Component {\n  render() {\n    return (<AtList><AtListItem title='a'/><AtButton type='primary'>ok</AtButton></AtList>)\n  }\n}\n")
	r := ScanH5(root)
	use := map[string]int{}
	for _, a := range r.Assets {
		use[a.Name] = a.Usage
	}
	if use["AtButton"] != 1 || use["AtListItem"] != 1 || len(r.Libraries) != 1 {
		t.Errorf("usage=%v libraries=%+v", use, r.Libraries)
	}
}

func TestSlashNamedComponents(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"^3.4.0","vant":"^4.8.0"}}`)
	write(t, root, "src/views/A.vue", `<template><van-button type="primary">x</van-button></template>`)
	exp := &figma.Export{
		Frame:         &figma.Node{Type: "FRAME", Children: []*figma.Node{{ID: "1:2", Type: "INSTANCE", ComponentID: "c1"}, {ID: "1:3", Type: "INSTANCE", ComponentID: "c2"}}},
		Components:    map[string]figma.Component{"c1": {ComponentSetID: "c1"}, "c2": {ComponentSetID: "c2"}},
		ComponentSets: map[string]figma.ComponentSet{"c1": {Name: "Button/Primary/Large"}, "c2": {Name: "Button / Secondary"}},
	}
	m := Match(exp, ScanH5(root))
	if len(m.Items) != 1 || m.Items[0].Figma != "Button" || m.Items[0].Instances != 2 || m.Items[0].Candidates[0].Name != "Button" {
		t.Errorf("items=%+v", m.Items)
	}
}

func TestAntdMobileV2(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"react":"^16.8.0","antd-mobile":"^2.3.4"}}`)
	write(t, root, "src/pages/A.jsx", "import { List, InputItem, Carousel } from 'antd-mobile'\nexport default function A() { return <List><InputItem placeholder='x'/><List.Item>a</List.Item><Carousel/></List> }")
	r := ScanH5(root)
	state := map[string]bool{}
	for _, a := range r.LibAll {
		state[a.Name] = usable(&a)
	}
	if !state["InputItem"] || !state["Carousel"] || state["Swiper"] || state["Divider"] {
		t.Errorf("antd-mobile 2: %v", state)
	}
}

func TestBrokerageNames(t *testing.T) {
	for name, want := range map[string]string{
		"GJTitleBar": "navbar", "GjStockListItem": "cell", "GJSwitchButton": "switch", "KLineView": "chart", "TimeSharingView": "chart",
		"分时图": "chart", "GjPriceTextView": "price", "GjStepperView": "stepper", "GjSearchView": "search", "SplitLine": "divider",
	} {
		if got := ConceptOf(name); got != want {
			t.Errorf("ConceptOf(%q)=%q want %q", name, got, want)
		}
	}
}

func TestFeatureSummary(t *testing.T) {
	a := &MatchResult{Page: "A", Items: []MatchItem{{Figma: "Button", Instances: 2, Verdict: "reuse", Candidates: []Candidate{{Name: "Button", Source: "Vant"}}}}}
	b := &MatchResult{Page: "B", Items: []MatchItem{{Figma: "买入按钮（未组件化）", Instances: 1, Verdict: "confirm", Candidates: []Candidate{{Name: "Button", Source: "Vant"}}},
		{Figma: "Price", Instances: 1, Verdict: "new"}}}
	s := FeatureSummary([]*MatchResult{a, b})
	if !strings.Contains(s, "| `Button` | ✅ 直接复用 | Button、买入按钮 | 2 | 3 |") || !strings.Contains(s, "需新建：Price") {
		t.Errorf("%s", s)
	}
}

func TestCompareMetrics(t *testing.T) {
	old := &Report{Metrics: Metrics{Componentization: 0.80, HardcodedColors: 10}}
	better := &Report{Metrics: Metrics{Componentization: 0.85, HardcodedColors: 10}}
	worse := &Report{Metrics: Metrics{Componentization: 0.85, HardcodedColors: 12}}
	if _, w := CompareMetrics(old, better); w {
		t.Error("an improvement flagged as regression")
	}
	if s, w := CompareMetrics(old, worse); !w || !strings.Contains(s, "❌ 变差") {
		t.Errorf("regression not flagged: %s", s)
	}
}

func TestHabitAttrs(t *testing.T) {
	got := strings.Join(habitAttrs(`<van-button round :color="'#1baeae'" type="primary" @click="go" :title="item.name" block>`), " ")
	if got != `round color="#1baeae" type="primary" block` {
		t.Errorf("habitAttrs=%q", got)
	}
	h := habitList(map[string]int{`size="small"`: 5, `size="large"`: 3, "round": 4, "plain": 1}, 8)
	if strings.Join(h, " ") != `size="small" round` {
		t.Errorf("habitList=%v", h)
	}
}

// every catalogued library component must render a snippet (uni-app built-ins have no import path)
func TestLibSnippetAllLibraries(t *testing.T) {
	for _, l := range uiLibs {
		for _, c := range l.Components {
			if s := libSnippet(Asset{Name: c, Source: l.Name, Kind: "ui-lib"}); s == "" {
				t.Errorf("%s/%s: empty snippet", l.Name, c)
			}
		}
	}
	if s := libSnippet(Asset{Name: "Switch", Source: "uni-app 内置组件"}); s != "<switch … />" {
		t.Errorf("uni-app built-in: %q", s)
	}
}

// AntennaPod shows every empty state through EmptyViewHandler (not a View): it is the project's empty
// component and must beat "needs a new component"
func TestViewHelperIsAnAsset(t *testing.T) {
	root := t.TempDir()
	write(t, root, "ui/common/src/main/java/a/EmptyViewHandler.java", `package a;
public class EmptyViewHandler {
    public EmptyViewHandler(Context c) { view = inflater.inflate(R.layout.empty_view_layout, null); }
}`)
	for _, p := range []string{"A", "B", "C"} {
		write(t, root, "app/src/main/java/a/"+p+"Fragment.java", "class "+p+"Fragment { void f() { emptyView = new EmptyViewHandler(getContext()); } }")
	}
	write(t, root, "app/src/main/java/a/Util.java", "public class LayoutHelper { int x = R.layout.foo; }") // no concept: not an asset
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><TextView/></LinearLayout>`)
	r := ScanAndroid(root)
	var got *Asset
	for i := range r.Assets {
		if r.Assets[i].Kind == "view-helper" {
			if got != nil {
				t.Fatalf("extra helper %s", r.Assets[i].Name)
			}
			got = &r.Assets[i]
		}
	}
	if got == nil || got.Name != "EmptyViewHandler" || got.Concept != "empty" || got.Usage != 3 {
		t.Fatalf("helper = %+v", got)
	}
	exp := &figma.Export{
		Frame:         &figma.Node{Type: "FRAME", Children: []*figma.Node{{ID: "1:2", Type: "INSTANCE", ComponentID: "c1"}}},
		Components:    map[string]figma.Component{"c1": {ComponentSetID: "s1"}},
		ComponentSets: map[string]figma.ComponentSet{"s1": {Name: "Empty"}},
	}
	if it := Match(exp, r).Items[0]; len(it.Candidates) == 0 || it.Candidates[0].Name != "EmptyViewHandler" {
		t.Errorf("Empty → %+v", it.Candidates)
	}
}

// AntennaPod: the manifest theme is a splash theme switched in code; every Theme.* style of the app is
// Material3, so <Button> / <CheckBox> are inflated as Material widgets and are not missed reuse
func TestMaterialInflatedWidgetsAndSplashTheme(t *testing.T) {
	root := t.TempDir()
	write(t, root, "gradle/libs.versions.toml", `material = { module = "com.google.android.material:material", version = "1.14.0" }`)
	write(t, root, "app/src/main/AndroidManifest.xml", `<manifest><application android:theme="@style/Theme.App.Splash"/></manifest>`)
	write(t, root, "app/src/main/res/values/styles.xml", `<resources>
<style name="Theme.App.Splash" parent="Theme.SplashScreen"><item name="android:statusBarColor">#000</item></style>
<style name="Theme.Base.App" parent="Theme.Material3.DynamicColors.Light"/>
<style name="Theme.App.Light" parent="Theme.Base.App"/>
</resources>`)
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><Button/><CheckBox/><com.google.android.material.button.MaterialButton/><com.google.android.material.checkbox.MaterialCheckBox/></LinearLayout>`)
	r := ScanAndroid(root)
	for _, f := range r.Findings {
		if f.Rule == "missed-reuse" {
			t.Errorf("Material theme inflates platform widgets: %+v", f)
		}
	}
	// without the Material theme the same layout bypasses MaterialCheckBox, which is a checkbox (not a switch)
	write(t, root, "app/src/main/res/values/styles.xml", `<resources><style name="Theme.App.Splash" parent="Theme.AppCompat.Light"/></resources>`)
	var hints []string
	for _, f := range ScanAndroid(root).Findings {
		if f.Rule == "missed-reuse" {
			hints = append(hints, f.Hint)
		}
	}
	if s := strings.Join(hints, "|"); !strings.Contains(s, "原生 CheckBox：项目里已有 MaterialCheckBox") || !strings.Contains(s, "原生 Button") {
		t.Errorf("hints=%v", hints)
	}
}

// SquareImageView extends AppCompatImageView: a specialisation, not a second image implementation
func TestSubclassIsNotADuplicate(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/src/main/java/a/SquareImageView.java", "package a;\npublic class SquareImageView extends AppCompatImageView {}")
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><a.SquareImageView/><a.SquareImageView/><a.SquareImageView/><ImageView/><ImageView/><ImageView/><ImageView/></LinearLayout>`)
	for _, c := range ScanAndroid(root).Clusters {
		if strings.HasPrefix(c.Concept, "image") {
			t.Errorf("cluster %+v", c)
		}
	}
}

// AntennaPod's DownloadActionButton extends ItemActionButton, an abstract class that is not a View
func TestNonViewBaseIsNotACustomView(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/src/main/java/a/ItemActionButton.java", "package a;\npublic abstract class ItemActionButton {\n}")
	write(t, root, "app/src/main/java/a/DownloadActionButton.java", "package a;\npublic class DownloadActionButton extends ItemActionButton {\n}")
	write(t, root, "app/src/main/java/a/PlayButton.java", "package a;\npublic class PlayButton extends AppCompatImageButton {\n}")
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><a.PlayButton/></LinearLayout>`)
	names := map[string]bool{}
	for _, a := range ScanAndroid(root).Assets {
		names[a.Name] = true
	}
	if names["DownloadActionButton"] || !names["PlayButton"] {
		t.Errorf("assets=%v", names)
	}
}

// mall-app-web imports its own copy of uni-number-box: easycom would map the tag to uni-ui, the import wins
func TestExplicitLocalImportBeatsLibraryPrefix(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"@dcloudio/uni-app":"3","@dcloudio/uni-ui":"1.5"}}`)
	write(t, root, "src/components/uni-number-box.vue", `<template><view class="box"><text>-</text><input/><text>+</text></view></template>`)
	write(t, root, "src/pages/cart/cart.vue", `<template><view><uni-number-box :value="1"></uni-number-box><uni-load-more/></view></template>
<script setup>
import uniNumberBox from '@/components/uni-number-box.vue'
</script>`)
	used := map[string]string{}
	for _, a := range ScanH5(root).Assets {
		if a.Usage > 0 {
			used[a.Name] = a.Kind
		}
	}
	if used["UniNumberBox"] != "business" || used["NumberBox"] != "" || used["LoadMore"] != "ui-lib" {
		t.Errorf("used=%v", used)
	}
}

// WanAndroid: its top bar is the third-party ActionBarCommon and its switches SwitchCompat — views from
// libraries the catalog does not know are still the project's toolkit; AppCompat* and a catalogued
// library's internals (QMUI LineTypeView) are not
func TestUncataloguedThirdPartyViews(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout>
<per.goweii.actionbarex.common.ActionBarCommon/><per.goweii.actionbarex.common.ActionBarCommon/>
<androidx.appcompat.widget.SwitchCompat/><androidx.appcompat.widget.AppCompatImageButton/>
<com.qmuiteam.qmui.type.view.LineTypeView/><com.kennyc.view.MultiStateView/></LinearLayout>`)
	got := map[string]string{}
	for _, a := range ScanAndroid(root).Assets {
		got[a.Name] = a.Source + "/" + a.Concept
	}
	want := map[string]string{"ActionBarCommon": "per.goweii.actionbarex/navbar", "SwitchCompat": "AppCompat/switch", "MultiStateView": "com.kennyc.view/empty"}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
	if _, ok := got["AppCompatImageButton"]; ok {
		t.Errorf("AppCompat* is a platform widget: %v", got)
	}
	if _, ok := got["LineTypeView"]; ok {
		t.Errorf("QMUI internals are not components: %v", got)
	}
}

// vue3-vant-mobile's GhostButton is <van-button plain>: a preset built on Vant's Button, not a duplicate
func TestH5PresetIsNotADuplicate(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"3","vant":"4.8.0"}}`)
	write(t, root, "src/components/GhostButton.vue", `<template><van-button size="small" plain type="primary"><slot /></van-button></template>`)
	write(t, root, "src/pages/a.vue", `<template><div><GhostButton/><GhostButton/><GhostButton/><van-button/><van-button/></div></template>`)
	for _, c := range ScanH5(root).Clusters {
		if c.Concept == "button" {
			t.Errorf("cluster %+v", c)
		}
	}
}

// unplugin-vue-components' VantResolver lets templates write <VanNavBar>: still Vant's NavBar
func TestPascalCaseLibraryTag(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"3","vant":"4.8.0"}}`)
	write(t, root, "src/components/NavBar.vue", `<template><VanNavBar :title="title" left-arrow /></template>`)
	write(t, root, "src/App.vue", `<template><div><nav-bar /><VanButton type="primary">ok</VanButton></div></template>`)
	used := map[string]int{}
	for _, a := range ScanH5(root).Assets {
		if a.Kind == "ui-lib" {
			used[a.Name] = a.Usage
		}
	}
	if used["NavBar"] != 1 || used["Button"] != 1 {
		t.Errorf("used=%v", used)
	}
}

// litemall registers { [Field.name]: Field } and writes <md-field>: the name option is the tag
func TestVueNameOptionAlias(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"2","vant":"2.12.0"}}`)
	write(t, root, "src/components/field/index.vue", "<template><div class=\"md_field\"><input/></div></template>\n<script>\nexport default {\n  name: 'md-field',\n}\n</script>")
	write(t, root, "src/views/login.vue", "<template><div><md-field/><md-field/></div></template>\n<script>\nimport Field from '@/components/field/';\nexport default { components: { [Field.name]: Field } }\n</script>")
	for _, a := range ScanH5(root).Assets {
		if a.Kind == "business" && a.Concept == "input" && a.Usage != 2 {
			t.Errorf("%s usage=%d", a.Name, a.Usage)
		}
	}
}

// KunMinX: AppCompat theme, Material 1.9 on the classpath. MaterialCardView would crash, so a Card with
// nothing else usable needs a new component; the warning stays listed for the agent
func TestUnusableCandidateNeverDecidesVerdict(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/build.gradle", `dependencies { api "com.google.android.material:material:1.9.0" }`)
	write(t, root, "app/src/main/AndroidManifest.xml", `<manifest><application android:theme="@style/AppTheme"/></manifest>`)
	write(t, root, "app/src/main/res/values/styles.xml", `<resources><style name="AppTheme" parent="Theme.AppCompat.Light.NoActionBar"/></resources>`)
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><TextView/><com.google.android.material.tabs.TabLayout/></LinearLayout>`)
	exp := &figma.Export{
		Frame:         &figma.Node{Type: "FRAME", Children: []*figma.Node{{ID: "1:2", Type: "INSTANCE", ComponentID: "c1"}}},
		Components:    map[string]figma.Component{"c1": {ComponentSetID: "s1"}},
		ComponentSets: map[string]figma.ComponentSet{"s1": {Name: "Card"}},
	}
	it := Match(exp, ScanAndroid(root)).Items[0]
	if it.Verdict != "new" {
		t.Errorf("verdict %s, candidates %+v", it.Verdict, it.Candidates)
	}
	for _, c := range it.Candidates {
		if c.Name == "MaterialCardView" && !strings.Contains(c.Why, "⚠") {
			t.Errorf("MaterialCardView listed without its warning: %+v", c)
		}
	}
}

// jeecg-uniapp is built on wot-design-uni: wd-* tags are its toolkit (wd-status-tip is the empty state)
func TestWotDesignUni(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"@dcloudio/uni-app":"3","wot-design-uni":"^1.9.1"}}`)
	write(t, root, "src/pages/a.vue", `<template><view><wd-cell title="a"/><wd-cell title="b"/><wd-status-tip image="content" tip="暂无内容"/><wd-input-number v-model="n"/></view></template>`)
	got := map[string]string{}
	for _, a := range ScanH5(root).Assets {
		if a.Source == "wot-design-uni" {
			got[a.Name] = a.Concept
		}
	}
	if got["Cell"] != "cell" || got["StatusTip"] != "empty" || got["InputNumber"] != "stepper" {
		t.Errorf("got=%v", got)
	}
}

// jeecg uses wd-input 38 times: a uni-app built-in <input> there is a missed reuse; without a toolkit
// input in real use (mall-app-web) it is just the platform's input
func TestUniBuiltinMissedReuse(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"@dcloudio/uni-app":"3","wot-design-uni":"^1.9.1"}}`)
	write(t, root, "src/pages/a.vue", `<template><view><wd-input/><wd-input/><wd-input/><input class="raw"/></view></template>`)
	n := 0
	for _, f := range ScanH5(root).Findings {
		if f.Rule == "missed-reuse" && strings.Contains(f.Hint, "uni-app 内置") {
			n++
		}
	}
	if n != 1 {
		t.Errorf("missed-reuse findings = %d", n)
	}
	write(t, root, "src/pages/a.vue", `<template><view><wd-input/><input class="raw"/></view></template>`)
	for _, f := range ScanH5(root).Findings {
		if f.Rule == "missed-reuse" {
			t.Errorf("toolkit input used once is not established: %+v", f)
		}
	}
}

// jeecg: echarts line charts (JFlyLineMap) reach "divider" only through "line", and SelectUserModal is
// a feature dialog: neither is an alternative implementation to converge (H5)
func TestFeatureAndWeakNamesStayOutOfClusters(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"3","vant":"4.8.0"}}`)
	write(t, root, "src/components/JFlyLineMap.vue", `<template><div class="chart"></div></template>`)
	write(t, root, "src/components/SelectUserModal.vue", `<template><div class="modal"></div></template>`)
	write(t, root, "src/components/SelectDeptModal.vue", `<template><div class="modal"></div></template>`)
	write(t, root, "src/pages/a.vue", `<template><div><JFlyLineMap/><JFlyLineMap/><JFlyLineMap/><SelectUserModal/><SelectUserModal/><SelectDeptModal/><van-dialog/></div></template>`)
	for _, c := range ScanH5(root).Clusters {
		if c.Concept == "divider" || c.Concept == "dialog" {
			t.Errorf("cluster %+v", c)
		}
	}
}

// WanAndroid's top bar is the third-party ActionBarCommon (× 24): a hand-built top bar must point to it,
// while a FloatingActionButton used a few times is still no replacement for a hand-built button
func TestHandBuiltPointsToEstablishedThirdPartyView(t *testing.T) {
	root := t.TempDir()
	var bars, fabs string
	for i := 0; i < 12; i++ {
		bars += `<per.goweii.actionbarex.common.ActionBarCommon/>`
	}
	for i := 0; i < 3; i++ {
		fabs += `<com.google.android.material.floatingactionbutton.FloatingActionButton/>`
	}
	write(t, root, "app/src/main/res/layout/old.xml", "<LinearLayout>"+bars+fabs+"</LinearLayout>")
	write(t, root, "app/src/main/res/layout/page.xml", `<LinearLayout>
<RelativeLayout android:id="@+id/top_bar"><ImageView android:id="@+id/iv_back"/><TextView android:id="@+id/tv_title"/></RelativeLayout>
<TextView android:id="@+id/btn_buy"/></LinearLayout>`)
	pc := ScanAndroid(root).CheckFocus("app/src/main/res/layout/page.xml")
	var hints []string
	for _, f := range pc.HandBuilt {
		hints = append(hints, f.Hint)
	}
	s := strings.Join(hints, "|")
	if !strings.Contains(s, "已有 ActionBarCommon") || strings.Contains(s, "FloatingActionButton") {
		t.Errorf("hints=%v", hints)
	}
}

// NewPipe's own org.schabi.newpipe.* views are not a third-party library; QMUILoadingView (named for
// nothing but loading) leads the loading cluster over the QMUIProgressBar variant
func TestOwnPackageAndGenericClusterLead(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/src/main/java/org/schabi/newpipe/App.kt", "package org.schabi.newpipe.app\n\nclass App")
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><org.schabi.newpipe.views.AnimatedProgressBar/><per.goweii.actionbarex.common.ActionBarCommon/></LinearLayout>`)
	for _, l := range ScanAndroid(root).Libraries {
		if l.Name == "org.schabi.newpipe" {
			t.Errorf("own package listed as a library: %+v", l)
		}
	}
	a := Asset{Name: "QMUILoadingView", Kind: "custom-view", Concept: "loading", Usage: 5}
	b := Asset{Name: "QMUIProgressBar", Kind: "custom-view", Concept: "loading", Usage: 6}
	if preference(a, "loading（View）") <= preference(b, "loading（View）") {
		t.Errorf("generic QMUILoadingView should lead QMUIProgressBar")
	}
}

// review: name: 'Button' must not turn raw <button> tags into uses of that component
func TestNameOptionDoesNotShadowHTMLTags(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"3"}}`)
	write(t, root, "src/components/Button.vue", "<template><button class=\"btn\"><slot/></button></template>\n<script>\nexport default {\n  name: 'Button',\n}\n</script>")
	write(t, root, "src/views/Order.vue", `<template><div><button>a</button><button>b</button></div></template>`)
	n := 0
	for _, f := range ScanH5(root).Findings {
		if f.Rule == "missed-reuse" {
			n++
		}
	}
	if n != 2 {
		t.Errorf("missed-reuse findings = %d, want 2", n)
	}
}

// review: a base view the regexes cannot read (generics, Kotlin secondary constructors) is still a View
func TestUnreadableBaseViewKeepsSubclass(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/src/main/java/a/BaseItemView.java", "package a;\npublic abstract class BaseItemView<T> extends FrameLayout {\n}")
	write(t, root, "app/src/main/java/a/GoodsItemView.java", "package a;\npublic class GoodsItemView extends BaseItemView<Goods> {\n}")
	write(t, root, "app/src/main/java/a/RoundView.kt", "package a\n\nopen class RoundView : FrameLayout {\n    constructor(c: Context) : super(c)\n}")
	write(t, root, "app/src/main/java/a/AvatarView.kt", "package a\n\nclass AvatarView(c: Context) : RoundView(c)")
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><a.GoodsItemView/><a.GoodsItemView/><a.AvatarView/></LinearLayout>`)
	names := map[string]bool{}
	for _, a := range ScanAndroid(root).Assets {
		names[a.Name] = true
	}
	if !names["GoodsItemView"] || !names["AvatarView"] {
		t.Errorf("assets=%v", names)
	}
}

// review: a self-closing <style/> before the splash style must not swallow its postSplashScreenTheme
func TestSplashThemeAfterSelfClosingStyle(t *testing.T) {
	values := []string{`<resources>
<style name="Theme.Other" parent="Theme.AppCompat.Light" />
<style name="Theme.App.Starting" parent="Theme.SplashScreen"><item name="postSplashScreenTheme">@style/Theme.App</item></style>
<style name="Theme.App" parent="Theme.MaterialComponents.DayNight"/>
</resources>`}
	if f, _ := appThemeFamily([]string{`<application android:theme="@style/Theme.App.Starting">`}, values); f != "material" {
		t.Errorf("family = %q", f)
	}
}

// review: a project input built on uni-app's <input> is not a missed reuse of itself
func TestUniComponentBuiltOnPrimitive(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"@dcloudio/uni-app":"3"}}`)
	write(t, root, "src/components/MyInput.vue", `<template><view class="my-input"><input class="inner"/></view></template>`)
	write(t, root, "src/pages/a.vue", `<template><view><my-input/><my-input/><my-input/></view></template>`)
	for _, f := range ScanH5(root).Findings {
		if f.Rule == "missed-reuse" {
			t.Errorf("%+v", f)
		}
	}
}

func TestReviewSmallFixes(t *testing.T) {
	if domainType("Alignment.Horizontal") || domainType("Arrangement.Vertical") || !domainType("Task.DownloadState") {
		t.Error("nested types: only a model's nested type is a model")
	}
	if catalogued("com.google.android.exoplayer2.ui.PlayerView") || !catalogued("com.google.android.material.internal.Foo") || !catalogued("com.qmuiteam.qmui.type.view.LineTypeView") {
		t.Error("catalogued uses the library's real package root")
	}
	root := t.TempDir()
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><androidx.constraintlayout.widget.Placeholder/><androidx.constraintlayout.helper.widget.Flow/></LinearLayout>`)
	for _, a := range ScanAndroid(root).Assets {
		if a.Name == "Placeholder" || a.Name == "Flow" {
			t.Errorf("layout helper became an asset: %+v", a)
		}
	}
}

// Seal: HelpDialog calls AlertDialog, SealDialog (same file, same name length) is built on BasicAlertDialog;
// NiaButton has overloads and the second call belongs to NiaOutlinedButton defined in between
func TestComposeWrapperIsTheEnclosingComposable(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/build.gradle", `dependencies { implementation("androidx.compose.material3:material3:1.2.0") }`)
	write(t, root, "core/designsystem/src/main/java/a/Dialogs.kt", `package a
import androidx.compose.material3.*

@Composable
fun HelpDialog(text: String) {
    AlertDialog(onDismissRequest = {}, text = { Text(text) })
}

@Composable
fun SealDialog(content: @Composable () -> Unit) {
    BasicAlertDialog(onDismissRequest = {}) { content() }
}
`)
	write(t, root, "core/designsystem/src/main/java/a/Buttons.kt", `package a
import androidx.compose.material3.*

@Composable
fun NiaButton(onClick: () -> Unit) {
    Button(onClick = onClick) {}
}

@Composable
fun NiaOutlinedButton(onClick: () -> Unit) {
    OutlinedButton(onClick = onClick) {}
}

@Composable
fun NiaButton(text: String, onClick: () -> Unit) {
    NiaButton(onClick)
}
`)
	write(t, root, "feature/src/main/java/b/Screen.kt", `package b
import a.*

@Composable
internal fun Screen() {
    HelpDialog("x")
    HelpDialog("y")
    NiaOutlinedButton {}
    NiaOutlinedButton {}
    NiaButton {}
    NiaButton {}
}
`)
	wrapped := map[string]string{}
	for _, a := range ScanAndroid(root).LibAll {
		if a.WrappedBy != "" {
			wrapped[a.Name] = a.WrappedBy
		}
	}
	if wrapped["AlertDialog"] == "SealDialog" || wrapped["Button"] != "NiaButton" || wrapped["OutlinedButton"] != "NiaOutlinedButton" {
		t.Errorf("wrapped=%v", wrapped)
	}
}

func TestReviewPlausibleFixes(t *testing.T) {
	if got := mdTags("原生 <input><button> 和 `<view>`"); got != "原生 `<input>``<button>` 和 `<view>`" {
		t.Errorf("mdTags = %q", got)
	}
	if javaPackage("chart.js") || javaPackage("@scope/x.y") || !javaPackage("per.goweii.actionbarex") || !javaPackage("com.lihang") {
		t.Error("javaPackage")
	}
	// Material theme: <Button> is MaterialButton already, but the project's own RoundButton is still the reuse
	root := t.TempDir()
	write(t, root, "app/src/main/AndroidManifest.xml", `<manifest><application android:theme="@style/AppTheme"/></manifest>`)
	write(t, root, "app/src/main/res/values/styles.xml", `<resources><style name="AppTheme" parent="Theme.MaterialComponents.Light"/></resources>`)
	write(t, root, "app/src/main/java/a/RoundButton.java", "package a;\npublic class RoundButton extends MaterialButton {\n}")
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><a.RoundButton/><a.RoundButton/><com.google.android.material.button.MaterialButton/><com.google.android.material.button.MaterialButton/><com.google.android.material.button.MaterialButton/><Button/></LinearLayout>`)
	var hints []string
	for _, f := range ScanAndroid(root).Findings {
		if f.Rule == "missed-reuse" {
			hints = append(hints, f.Hint)
		}
	}
	if len(hints) != 1 || !strings.Contains(hints[0], "RoundButton") {
		t.Errorf("hints=%v", hints)
	}
}

// review: a changed file whose name has a comma is one file; a user's --focus stays comma-separated
func TestFocusListKeepsCommas(t *testing.T) {
	if got := strings.Join(focusList("a,b.vue"+FocusSep+"c.vue"), "|"); got != "a,b.vue|c.vue" {
		t.Errorf("--changed list: %q", got)
	}
	if got := strings.Join(focusList("a.xml,B.kt"), "|"); got != "a.xml|B.kt" {
		t.Errorf("--focus list: %q", got)
	}
}

// review 3: H5 scanner fixes
func TestH5ScannerReviewFixes(t *testing.T) {
	// a native <select> is not the imported Select component; a JSX comment / type argument is not a usage
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"react":"18","vue":"3"}}`)
	write(t, root, "src/components/Select.vue", `<template><div class="select"><slot/></div></template>`)
	write(t, root, "src/views/A.vue", "<template><div><select></select><Select/></div></template>\n<script>\nimport Select from '@/components/Select.vue'\n</script>")
	write(t, root, "src/components/Button.tsx", "export function Button(props: Props) { return <button className=\"b\">{props.children}</button> }")
	write(t, root, "src/views/B.tsx", "import { Button } from '../components/Button'\nconst r = useRef<Button>(null)\n// <Button old/>\nexport function B() { return <div>{/* <Button/> */}<Button/></div> }")
	r := ScanH5(root)
	use := map[string]int{}
	for _, a := range r.Assets {
		use[a.Name] += a.Usage
	}
	if use["Select"] != 1 || use["Button"] != 1 {
		t.Errorf("usage = %v", use)
	}
	for _, f := range r.Findings {
		if f.Rule == "missed-reuse" && strings.HasSuffix(f.File, "Button.tsx") {
			t.Errorf("Button.tsx flagged for its own <button>: %+v", f)
		}
	}

	// a colour after a token definition is still a usage; border-bottom: 0 is not "pinned to the bottom"
	root = t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"3"}}`)
	write(t, root, "src/components/TopBar.vue", "<template><div class=\"top\">t</div></template>\n<style>\n$blue: #00f;\n.top { position: fixed; top: 0; border-bottom: 0; color: #123456; }\n</style>")
	write(t, root, "src/views/P.vue", "<template><div><top-bar/><p class=\"t\">长长的中文说明文字长长的中文说明文字长长的中文说明文字长长的中文说明文字长长的中文说明文字长长的中文说明文字长长的中文</p><p class=\"t\">x</p></div></template>")
	r = ScanH5(root)
	found := false
	for _, h := range r.Colors {
		found = found || h.Value == "#123456"
	}
	if !found {
		t.Errorf("#123456 after $blue: was not counted: %+v", r.Colors)
	}
	for _, a := range r.Assets {
		if a.Name == "TopBar" && a.Concept == "tabbar" {
			t.Error("border-bottom: 0 made a top bar a tab bar")
		}
	}
	if s := clip("长长的中文", 4); s != "长…" {
		t.Errorf("clip = %q", s)
	}
}

// review 3: copy-paste within one file is a clone; layout-land copies of one layout are not
func TestClonesWithinOneFile(t *testing.T) {
	block := `<div class="row"><span>a</span><span>b</span><em>c</em><i>d</i></div>`
	files := map[string]string{"Same.vue": "<div>" + block + "<p/>" + block + "</div>"}
	if c := findClones(files, parseMarkup, 5); len(c) != 1 {
		t.Errorf("clones in one file = %d", len(c))
	}
	files = map[string]string{"app/src/main/res/layout/a.xml": block, "app/src/main/res/layout-land/a.xml": block}
	if c := findClones(files, parseMarkup, 5); len(c) != 0 {
		t.Errorf("layout-land copy reported as clone: %+v", c)
	}
}

// review 3: with <script setup> first, clone and convention locations are lines of the file
func TestMarkupLineNumbersAreFileLines(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"3","vant":"4.8.0"}}`)
	block := `<div class="row"><span>a</span><p>b</p><ul><li>c</li><li>d</li></ul><a>e</a><img/><section>f</section></div>`
	src := "<script setup>\nconst a = 1\nconst b = 2\n</script>\n\n<template>\n<div>\n" + block + "\n<div class=\"van-hairline--top\"/><div class=\"van-hairline--top\"/><div class=\"van-hairline--top\"/>\n</div>\n</template>\n"
	write(t, root, "src/views/A.vue", src)
	write(t, root, "src/views/B.vue", "<template>\n<div>\n"+block+"\n</div>\n</template>\n")
	r := ScanH5(root)
	if len(r.Clones) == 0 || !strings.Contains(strings.Join(r.Clones[0].Locations, " "), "src/views/A.vue:8") {
		t.Errorf("clones = %+v", r.Clones)
	}
	if len(r.Conventions) == 0 || r.Conventions[0].ExampleAt != "src/views/A.vue:9" {
		t.Errorf("conventions = %+v", r.Conventions)
	}
}

// review 3: multi-line tags keep their attributes (habits); quoted "<" in attributes does not break parsing
func TestMultiLineTagsAndQuotedAngles(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"3","vant":"4.8.0"}}`)
	var tags string
	for i := 0; i < 4; i++ {
		tags += "<van-button\n  type=\"primary\"\n  round\n>ok</van-button>\n"
	}
	write(t, root, "src/views/A.vue", "<template>\n<div>\n"+tags+"</div>\n</template>\n")
	for _, a := range ScanH5(root).Assets {
		if a.Name == "Button" && a.Attrs[`type="primary"`] != 4 {
			t.Errorf("habit attrs = %v (example %q)", a.Attrs, a.Example)
		}
	}
	ns := parseMarkup(`<div><p v-if="a < b">x</p><span :show="n > 0"/></div>`)
	if len(ns) != 1 || len(ns[0].children) != 2 {
		t.Errorf("parsed %d roots / %+v", len(ns), ns)
	}
}

// review 3: every object-syntax prop is found (not just the first), nested keys are not props, and
// an inline array works too
func TestVuePropsParsing(t *testing.T) {
	obj := "export default {\n  props: {\n    title: { type: String, default: '' },\n    size: { type: String, validator: (v) => ['a', 'b'].includes(v) },\n    'is-link': Boolean,\n  },\n}"
	if got := strings.Join(vueProps(obj), ","); got != "title,size,is-link" {
		t.Errorf("object props = %q", got)
	}
	if got := strings.Join(vueProps("export default { props: ['label', 'value'] }"), ","); got != "label,value" {
		t.Errorf("inline array props = %q", got)
	}
}

// review 3: two Card.vue in different folders are two components; an import resolves to its own file
func TestSameNamedComponentsInTwoFolders(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"3"}}`)
	write(t, root, "src/components/goods/Card.vue", "<template><div class=\"goods\"/></template>\n<script>\nexport default { props: { price: Number } }\n</script>")
	write(t, root, "src/components/user/Card.vue", "<template><div class=\"user\"/></template>\n<script>\nexport default { props: { avatar: String } }\n</script>")
	write(t, root, "src/views/A.vue", "<template><div><Card/><Card/></div></template>\n<script>\nimport Card from '@/components/user/Card.vue'\n</script>")
	got := map[string]string{}
	for _, a := range ScanH5(root).Assets {
		if a.Name == "Card" {
			got[a.Source] = strings.Join(a.Props, ",") + "/" + itoa(a.Usage)
		}
	}
	if got["src/components/goods/Card.vue"] != "price/0" || got["src/components/user/Card.vue"] != "avatar/2" {
		t.Errorf("cards = %v", got)
	}
}

// review 4: Android scanner fixes
func TestAndroidScannerReviewFixes(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/src/main/java/a/Views.kt", `package a
internal open class RoundView : FrameLayout {
    constructor(c: Context) : super(c)
}
class TagView<T : Any>(c: Context) : TextView(c)
`)
	write(t, root, "app/src/main/java/a/Fancy.java", "package a;\npublic final class FancyButton extends Button {\n}\nclass PlainCard extends FrameLayout {\n}")
	write(t, root, "app/src/main/java/a/Ui.kt", `package a
@Composable
fun DsChip(@StringRes labelRes: Int, onPick: (id: String, index: Int) -> Unit, vararg tags: String, title: String = "a, b: c") {}
@Composable
fun DsCard(content: @Composable () -> Unit) {}
@Composable
fun <T> DsPicker(items: List<T>) {}
@Composable
fun RowScope.DsTab(selected: Boolean) {}
@Composable
internal fun Screen() { DsCard(DsChip(1, { _, _ -> })) ; DsPicker(listOf(1)); DsTab(true) }
`)
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><a.RoundView/><a.TagView/><a.FancyButton/><a.PlainCard/><View android:layout_width="0dp" android:layout_height="1dp"/></LinearLayout>`)
	r := ScanAndroid(root)
	got := map[string]*Asset{}
	for i := range r.Assets {
		got[r.Assets[i].Name] = &r.Assets[i]
	}
	for _, n := range []string{"RoundView", "TagView", "FancyButton", "PlainCard", "DsPicker", "DsTab"} {
		if got[n] == nil {
			t.Errorf("%s not found", n)
		}
	}
	if c := got["DsChip"]; c == nil || c.Usage != 1 || strings.Join(c.Props, ",") != "labelRes,onPick,tags,title" || strings.Join(c.Required, ",") != "labelRes,onPick,tags" {
		t.Errorf("DsChip = %+v", c)
	}
	if r.Metrics.HardcodedSizes != 1 {
		t.Errorf("hardcoded sizes = %d (0dp is not a size)", r.Metrics.HardcodedSizes)
	}
	// res/ at the project root (Eclipse / Ant layout)
	root = t.TempDir()
	write(t, root, "src/a/FancyButton.java", "package a;\npublic class FancyButton extends Button {\n}")
	write(t, root, "res/layout/a.xml", `<LinearLayout><a.FancyButton/></LinearLayout>`)
	write(t, root, "build.gradle", "")
	for _, a := range ScanAndroid(root).Assets {
		if a.Name == "FancyButton" && a.Usage != 1 {
			t.Errorf("root res/ layout not scanned: %+v", a)
		}
	}
}

// review 4: every colour attribute counts, short forms included, and matches its token
func TestAndroidColourAttributes(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/src/main/res/values/colors.xml", `<resources><color name="white" tools:ignore="UnusedResources">#FFFFFF</color></resources>`)
	write(t, root, "app/src/main/res/layout/a.xml", `<LinearLayout><TextView android:textColorHint="#fff" android:backgroundTint="#FF123456"/><androidx.cardview.widget.CardView app:cardBackgroundColor="#abc"/></LinearLayout>`)
	got := map[string]string{}
	for _, h := range ScanAndroid(root).Colors {
		got[h.Value] = h.Token
	}
	if tok, ok := got["#FFFFFF"]; !ok || tok != "@color/white" || got["#123456"] != "" || len(got) != 3 {
		t.Errorf("colours = %v", got)
	}
}

// a changed shared stylesheet: tokens added there are not a page palette, hard-coded colours still count
func TestStylesheetInPageCheck(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"3"}}`)
	write(t, root, "src/assets/scss/_var.scss", "$brand: #ff0000;\n$sub: #00ff00;\n.a { color: #123456; }\n")
	write(t, root, "src/A.vue", `<template><div class="a"/></template>`)
	pc := ScanH5(root).CheckFocus("src/assets/scss/_var.scss")
	if pc.LocalTokens != 0 || len(pc.Colors) != 1 {
		t.Errorf("localTokens=%d colors=%+v", pc.LocalTokens, pc.Colors)
	}
}

// sunniejs/vue-h5-template: Vant / NutUI demo pages under views/ui are pages, '#ui-demo' is a build alias,
// and with both suites installed the .env setting decides which one new pages use — never a mix
func TestPrimaryUISuite(t *testing.T) {
	if classify("src/views/ui/vant.vue") != "page" || classify("src/components/ui/Button.vue") != "design-system" || classify("src/ui/Button.vue") != "design-system" {
		t.Error("classify: the outermost directory decides")
	}
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"3","vant":"^4.10.0","@nutui/nutui":"^4.3.13"}}`)
	write(t, root, ".env", "# VITE_UI_FRAMEWORK - vant / nutui\nVITE_UI_FRAMEWORK=vant\n")
	write(t, root, "src/views/ui/vant.vue", `<template><div><VanButton/><VanCell/></div></template>`)
	write(t, root, "src/views/ui/nutui.vue", `<template><div><NutButton/><NutCell/></div></template>`)
	write(t, root, "src/views/ui/index.vue", "<template><SelectedUiDemo/></template>\n<script setup>\nimport SelectedUiDemo from '#ui-demo'\n</script>")
	write(t, root, "src/views/order.vue", `<template><div><button class="buy">买入</button></div></template>`)
	r := ScanH5(root)
	if r.PrimaryLib != "Vant" || !strings.Contains(r.PrimaryWhy, "VITE_UI_FRAMEWORK=vant") {
		t.Errorf("primary = %q (%s)", r.PrimaryLib, r.PrimaryWhy)
	}
	for _, l := range r.Libraries {
		if strings.HasPrefix(l.Name, "#") {
			t.Errorf("build alias listed as a library: %+v", l)
		}
	}
	for _, a := range r.Assets {
		if a.Kind == "design-system" {
			t.Errorf("demo page became a design-system asset: %+v", a)
		}
	}
	for _, f := range r.Findings {
		if f.Rule == "missed-reuse" && !strings.Contains(f.Hint, "Vant") {
			t.Errorf("raw <button> should point to Vant: %s", f.Hint)
		}
	}
	exp := &figma.Export{
		Frame:         &figma.Node{Type: "FRAME", Children: []*figma.Node{{ID: "1:2", Type: "INSTANCE", ComponentID: "c1"}, {ID: "1:3", Type: "INSTANCE", ComponentID: "c2"}}},
		Components:    map[string]figma.Component{"c1": {ComponentSetID: "s1"}, "c2": {ComponentSetID: "s2"}},
		ComponentSets: map[string]figma.ComponentSet{"s1": {Name: "Button"}, "s2": {Name: "Cell"}},
	}
	for _, it := range Match(exp, r).Items {
		if c := it.Candidates[0]; c.Source != "Vant" {
			t.Errorf("%s → %s (%s): mixes libraries", it.Figma, c.Name, c.Source)
		}
	}
}

// yudao-mall-uniapp: a ternary in a prop value (a top-level ':' after the key) crashed the parser
func TestVuePropsTernaryValue(t *testing.T) {
	src := "export default {\n  props: {\n    size: { type: String, default: isPad ? 'lg' : 'md' },\n    mode: isDark ? 'dark' : 'light',\n    title: String,\n  },\n}"
	if got := strings.Join(vueProps(src), ","); got != "size,mode,title" {
		t.Errorf("props = %q", got)
	}
}

// yudao: #fff is both $white (used everywhere) and a theme var: suggest the token the project uses
func TestMostReferencedTokenWins(t *testing.T) {
	root := t.TempDir()
	write(t, root, "package.json", `{"dependencies":{"vue":"3"}}`)
	write(t, root, "src/styles/_var.scss", "$white: #ffffff;\n:root { --ui-BG-Main-TC: #ffffff; }\n.a { color: $white; }\n.b { background: $white; }\n")
	write(t, root, "src/views/A.vue", "<template><div class=\"x\"/></template>\n<style>\n.x { color: #FFFFFF; }\n</style>")
	for _, h := range ScanH5(root).Colors {
		if h.Value == "#FFFFFF" && h.Token != "$white" {
			t.Errorf("token = %q", h.Token)
		}
	}
}

// complex-system test round: AppCompatButton is the platform Button; an icon button or a domain card
// is not "the" Button / Card; a tooling package.json does not make an H5 project
func TestComplexSystemsRound(t *testing.T) {
	root := t.TempDir()
	write(t, root, "app/src/main/java/a/SubmitButton.java", "package a;\npublic class SubmitButton extends AppCompatButton {\n}")
	var tags string
	for i := 0; i < 5; i++ {
		tags += `<androidx.appcompat.widget.AppCompatButton/>`
	}
	write(t, root, "app/src/main/res/layout/a.xml", "<LinearLayout>"+tags+"<a.SubmitButton/></LinearLayout>")
	write(t, root, "package.json", `{"devDependencies":{"husky":"9"}}`)
	r := ScanAndroid(root)
	for _, a := range r.Assets {
		if a.Name == "Button" && a.Source == "Android SDK widgets" && a.Usage != 5 {
			t.Errorf("AppCompatButton not counted as Button: %+v", a)
		}
	}
	if ps := Detect(root); len(ps) != 1 || ps[0] != "android" {
		t.Errorf("Detect = %v (tooling package.json is not H5)", ps)
	}
	if !purposeWords["icon"] || !purposeWords["live"] {
		t.Error("purpose words")
	}
}
