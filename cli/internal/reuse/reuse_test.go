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
