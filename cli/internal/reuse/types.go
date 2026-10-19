// Package reuse is the brownfield ("存量项目") reuse check: it works on any common H5 / Android
// project without Figma. It inventories component assets (UI library, design system, business,
// page-local, custom views, composables), then finds reuse problems:
//   - duplicate assets: several components implementing the same UI concept
//   - missed reuse: raw primitives / hand-built UI where a component exists
//   - structural clones: the same markup subtree copy-pasted across files (extraction candidates)
//   - hard-coded values: colours / sizes that already exist as tokens
package reuse

import (
	"sort"
	"strings"
)

// Asset is one reusable UI building block found in the project or its UI library.
type Asset struct {
	Name       string   `json:"name"`
	Kind       string   `json:"kind"`   // ui-lib | design-system | business | page-local | custom-view | composable
	Source     string   `json:"source"` // library name or defining file
	Concept    string   `json:"concept,omitempty"`
	Props      []string `json:"props,omitempty"`
	Usage      int      `json:"usage"`
	UsageFiles []string `json:"usageFiles,omitempty"`
	Example    string   `json:"example,omitempty"` // a real usage snippet from the project
	ExampleAt  string   `json:"exampleAt,omitempty"`
}

type Finding struct {
	Rule     string `json:"rule"`
	Severity string `json:"severity"` // error | warn | info
	File     string `json:"file"`
	Line     int    `json:"line"`
	Text     string `json:"text"`
	Hint     string `json:"hint"`
}

type Cluster struct {
	Concept    string   `json:"concept"`
	Members    []string `json:"members"` // "Name(kind, usage)"
	Recommend  string   `json:"recommend"`
	TotalUsage int      `json:"totalUsage"`
}

type Clone struct {
	Signature string   `json:"signature"`
	Elements  int      `json:"elements"`
	Locations []string `json:"locations"` // file:line
	Sample    string   `json:"sample"`
}

type Hotspot struct {
	Value  string   `json:"value"`
	Count  int      `json:"count"`
	Token  string   `json:"token,omitempty"` // existing token with the same value
	Sample []string `json:"sample,omitempty"`
}

type Metrics struct {
	Files             int     `json:"files"`
	ComponentUsages   int     `json:"componentUsages"`
	RawPrimitives     int     `json:"rawPrimitives"`    // primitives that have a component equivalent
	Componentization  float64 `json:"componentization"` // componentUsages / (componentUsages + rawPrimitives)
	DuplicateConcepts int     `json:"duplicateConcepts"`
	CloneGroups       int     `json:"cloneGroups"`
	HardcodedColors   int     `json:"hardcodedColors"`
	TokenizableColors int     `json:"tokenizableColors"`
	HardcodedSizes    int     `json:"hardcodedSizes"`
}

type Report struct {
	Platform    string            `json:"platform"`
	Root        string            `json:"root"`
	Stack       []string          `json:"stack"`
	Libraries   []LibUse          `json:"libraries"`
	Assets      []Asset           `json:"assets"`
	Clusters    []Cluster         `json:"clusters"`
	Findings    []Finding         `json:"findings"`
	Clones      []Clone           `json:"clones"`
	Colors      []Hotspot         `json:"colors"`
	LibAll      []Asset           `json:"-"` // every component of the detected UI libraries (used or not), for matching
	ColorTokens map[string]string `json:"-"` // #RRGGBB → token reference
	Metrics     Metrics           `json:"metrics"`
}

type LibUse struct {
	Name  string   `json:"name"`
	Usage int      `json:"usage"`
	Top   []string `json:"top"` // most used components
}

// concepts: canonical UI concept → name words that indicate it. Word-level match on CamelCase /
// kebab-case names, so StatusTag, van-tag, LegacyTagView, OrderStatusChip all map to "tag".
var concepts = []struct {
	Name  string
	Words []string
}{
	{"button", []string{"button", "btn"}},
	{"tag", []string{"tag", "badge", "chip", "label", "pill"}},
	{"cell", []string{"cell", "listitem", "row", "item", "inforow", "commonlistitemview"}},
	{"card", []string{"card", "panel", "section", "cellgroup", "grouplistview"}},
	{"navbar", []string{"navbar", "topbar", "appbar", "titlebar", "toolbar", "header"}},
	{"price", []string{"price", "amount", "money", "pricetext"}},
	{"dialog", []string{"dialog", "modal", "alert", "confirm"}},
	{"toast", []string{"toast", "snackbar", "notify", "tips"}},
	{"popup", []string{"popup", "sheet", "actionsheet", "bottomsheet", "drawer", "popover", "floatingpanel"}},
	{"tabs", []string{"tabs", "tab", "tabbar", "segment", "segmented"}},
	{"input", []string{"input", "field", "edittext", "textfield", "search", "searchbar", "textarea"}},
	{"carousel", []string{"swipe", "swiper", "carousel", "banner", "viewpager"}},
	{"icon", []string{"icon"}},
	{"image", []string{"image", "img", "avatar", "imageview"}},
	{"empty", []string{"empty", "emptyview", "nodata", "placeholder"}},
	{"loading", []string{"loading", "spinner", "progress", "skeleton"}},
	{"list", []string{"list", "recyclerview", "listview", "lazycolumn"}},
	{"divider", []string{"divider", "separator", "line"}},
	{"notice", []string{"noticebar", "notice"}},
	{"switch", []string{"switch", "toggle", "checkbox", "radio", "radiobutton", "togglebutton", "checkboxgroup", "radiogroup"}},
	{"stepper", []string{"stepper", "counter"}},
	{"countdown", []string{"countdown", "timer"}},
	{"goods", []string{"goods", "product", "sku", "goodscard", "productcard"}},
	{"actionbar", []string{"actionbar", "bottombar", "submitbar", "bottomactionbar", "footerbar"}},
}

// noise words that never decide a concept on their own
var weak = map[string]bool{"item": true, "row": true, "box": true, "field": true, "line": true, "label": true, "section": true, "header": true, "footer": true, "list": true}

// splitWords splits CamelCase, kebab-case, snake_case and dotted names into lower-case words.
func splitWords(s string) []string {
	if i := strings.LastIndex(s, "."); i >= 0 {
		s = s[i+1:]
	}
	var out []string
	cur := []rune{}
	flush := func() {
		if len(cur) > 0 {
			out = append(out, strings.ToLower(string(cur)))
			cur = cur[:0]
		}
	}
	rs := []rune(s)
	for i, r := range rs {
		switch {
		case r == '-' || r == '_' || r == ' ':
			flush()
		case r >= 'A' && r <= 'Z' && len(cur) > 0 && (i+1 < len(rs) && rs[i+1] >= 'a' && rs[i+1] <= 'z' || rs[i-1] >= 'a' && rs[i-1] <= 'z'):
			flush()
			cur = append(cur, r)
		default:
			cur = append(cur, r)
		}
	}
	flush()
	return out
}

var vendorPrefix = map[string]bool{"qmui": true, "van": true, "nut": true, "ds": true, "legacy": true, "app": true, "common": true, "base": true, "my": true, "custom": true, "ui": true, "md": true, "t": true}

// infrastructure suffixes: layout helpers, managers, adapters … are not UI components of a concept
var infraSuffix = map[string]bool{"layout": true, "helper": true, "manager": true, "util": true, "utils": true, "controller": true, "background": true, "behavior": true, "decoration": true, "adapter": true, "holder": true, "mask": true, "root": true, "wrapper": true, "container": true, "delegate": true, "provider": true, "state": true}

// Chinese layer / component names are common in Figma files of Chinese teams. Longest keyword wins.
var zhConcepts = []struct{ kw, concept string }{
	{"标签页", "tabs"}, {"选项卡", "tabs"}, {"导航栏", "navbar"}, {"标题栏", "navbar"}, {"顶栏", "navbar"},
	{"操作栏", "actionbar"}, {"底部栏", "actionbar"}, {"提交栏", "actionbar"}, {"吸底", "actionbar"},
	{"通知栏", "notice"}, {"公告", "notice"}, {"提示条", "notice"}, {"倒计时", "countdown"},
	{"单元格", "cell"}, {"列表项", "cell"}, {"信息行", "cell"}, {"分割线", "divider"}, {"分隔线", "divider"},
	{"商品卡", "goods"}, {"商品", "goods"}, {"价格", "price"}, {"金额", "price"}, {"按钮", "button"},
	{"标签", "tag"}, {"徽标", "tag"}, {"角标", "tag"}, {"卡片", "card"}, {"弹窗", "dialog"}, {"对话框", "dialog"},
	{"轻提示", "toast"}, {"吐司", "toast"}, {"弹出层", "popup"}, {"面板", "popup"}, {"输入框", "input"}, {"搜索", "input"},
	{"图片", "image"}, {"头像", "image"}, {"空状态", "empty"}, {"加载", "loading"}, {"开关", "switch"}, {"复选框", "switch"},
	{"单选", "switch"}, {"步进器", "stepper"}, {"轮播", "carousel"}, {"图标", "icon"}, {"列表", "list"},
}

func hasCJK(s string) bool {
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
	}
	return false
}

// ConceptOf returns the UI concept a component name expresses ("" if none). The last strong word
// wins (StatusTag → tag, TagList → list is weak so → tag).
func ConceptOf(name string) string {
	if hasCJK(name) {
		best, bl := "", 0
		for _, z := range zhConcepts {
			if strings.Contains(name, z.kw) && len(z.kw) > bl {
				best, bl = z.concept, len(z.kw)
			}
		}
		if best != "" {
			return best
		}
		// mixed names like "Tag / 状态" fall through to the latin words
	}
	words := splitWords(name)
	for len(words) > 1 && vendorPrefix[words[0]] {
		words = words[1:]
	}
	joined := strings.Join(words, "")
	for _, c := range concepts {
		for _, w := range c.Words {
			if joined == w {
				return c.Name
			}
		}
	}
	if len(words) > 1 && infraSuffix[words[len(words)-1]] {
		return ""
	}
	for _, c := range concepts {
		for _, w := range c.Words {
			if joined == w {
				return c.Name
			}
		}
	}
	best, weakBest := "", ""
	for _, w := range words {
		w2 := strings.TrimSuffix(w, "s")
		for _, c := range concepts {
			for _, cw := range c.Words {
				if w == cw || w2 == cw {
					if weak[cw] {
						weakBest = c.Name
					} else {
						best = c.Name
					}
				}
			}
		}
	}
	if best != "" {
		return best
	}
	return weakBest
}

// clusters groups assets of the same concept; a concept implemented by ≥2 assets is a duplicate.
func clusters(assets []Asset, libAll []Asset) []Cluster {
	by := map[string][]Asset{}
	for _, a := range assets {
		if a.Concept == "" || a.Kind == "page-local" && a.Usage == 0 {
			continue
		}
		by[a.Concept] = append(by[a.Concept], a)
	}
	// An unused library component of the same concept is the classic "built our own wheel" case.
	for _, l := range libAll {
		if l.Usage > 0 || l.Concept == "" || l.Source == "Android SDK widgets" {
			continue
		}
		hasProject := false
		for _, a := range by[l.Concept] {
			hasProject = hasProject || a.Kind != "ui-lib"
		}
		if hasProject && strings.EqualFold(l.Name, conceptCanonical(l.Concept, by[l.Concept])) {
			by[l.Concept] = append(by[l.Concept], l)
		}
	}
	var out []Cluster
	for c, as := range by {
		as = collapseFamilies(as)
		if len(as) < 2 {
			continue
		}
		projects := 0
		for _, a := range as {
			if a.Kind != "ui-lib" {
				projects++
			}
		}
		if projects == 0 {
			continue // siblings inside one UI library (Tab / Tabs) are not duplicates
		}
		sort.Slice(as, func(i, j int) bool { return rank(as[i]) > rank(as[j]) })
		cl := Cluster{Concept: c}
		for _, a := range as {
			cl.Members = append(cl.Members, a.Name+"（"+kindLabel(a.Kind)+"，用 "+itoa(a.Usage)+" 次）")
			cl.TotalUsage += a.Usage
		}
		cl.Recommend = as[0].Name
		for _, a := range as {
			if a.Kind == "ui-lib" && a.Usage == 0 {
				cl.Members[indexOf(as, a)] = a.Name + "（" + a.Source + "，项目未使用）"
			}
		}
		out = append(out, cl)
	}
	sort.Slice(out, func(i, j int) bool {
		return len(out[i].Members) > len(out[j].Members) || len(out[i].Members) == len(out[j].Members) && out[i].TotalUsage > out[j].TotalUsage
	})
	return out
}

// rank prefers design-system / UI-library assets, then usage.
func rank(a Asset) int {
	w := map[string]int{"design-system": 3000, "ui-lib": 2000, "custom-view": 1000, "composable": 1000, "business": 500, "page-local": 0}[a.Kind]
	if a.Source == "Android SDK widgets" {
		w = 100 // framework widgets are the thing projects wrap, never the recommendation
	}
	return w + a.Usage
}

// collapseFamilies drops sub-parts of a component family (QMUIDialogAction, QMUIDialogTitle belong
// to QMUIDialog) so clusters list alternatives, not anatomy.
func collapseFamilies(as []Asset) []Asset {
	var out []Asset
	for _, a := range as {
		part := false
		for _, b := range as {
			if b.Name != a.Name && len(b.Name) >= 4 && strings.HasPrefix(a.Name, b.Name) && len(a.Name) > len(b.Name) && a.Name[len(b.Name)] >= 'A' && a.Name[len(b.Name)] <= 'Z' {
				part = true
				break
			}
		}
		if !part {
			out = append(out, a)
		}
	}
	return out
}

func kindLabel(k string) string {
	return map[string]string{"ui-lib": "组件库", "design-system": "设计系统", "business": "业务组件", "page-local": "页面内私有", "custom-view": "自定义 View", "composable": "Composable"}[k]
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func indexOf(as []Asset, a Asset) int {
	for i := range as {
		if as[i].Name == a.Name && as[i].Kind == a.Kind {
			return i
		}
	}
	return 0
}

// conceptCanonical returns the library component name considered canonical for a concept
// (exact concept word in the name: NavBar for navbar, Tag for tag …).
func conceptCanonical(concept string, _ []Asset) string {
	for _, c := range concepts {
		if c.Name == concept {
			return c.Words[0]
		}
	}
	return concept
}
