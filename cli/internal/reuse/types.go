// Package reuse is the brownfield ("存量项目") reuse check: it works on any common H5 / Android
// project without Figma. It inventories component assets (UI library, design system, business,
// page-local, custom views, composables), then finds reuse problems:
//   - duplicate assets: several components implementing the same UI concept
//   - missed reuse: raw primitives / hand-built UI where a component exists
//   - structural clones: the same markup subtree copy-pasted across files (extraction candidates)
//   - hard-coded values: colours / sizes that already exist as tokens
package reuse

import (
	"regexp"
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
	// compatibility with the project as it is (see compat.go)
	Caveat      string              `json:"caveat,omitempty"`
	Unavailable bool                `json:"unavailable,omitempty"` // the installed library version does not ship it
	WrappedBy   string              `json:"wrappedBy,omitempty"`   // the design-system component that wraps this library component
	Enums       map[string][]string `json:"enums,omitempty"`       // enum-like props → allowed values
	Required    []string            `json:"required,omitempty"`    // parameters without a default (composables)
	Styles      map[string]int      `json:"styles,omitempty"`      // XML style="@style/…" applied with it, by count
	Attrs       map[string]int      `json:"attrs,omitempty"`
	Extends     string              `json:"extends,omitempty"` // Android custom view: its base class (SquareImageView → AppCompatImageView)       // H5 static props the project passes, by count (round, color="#1baeae")
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
	Conventions []Asset           `json:"conventions,omitempty"` // established non-component ways (see convention.go)
	LibAll      []Asset           `json:"-"`                     // every component of the detected UI libraries (used or not), for matching
	ColorTokens map[string]string `json:"-"`                     // #RRGGBB → token reference
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
	{"tag", []string{"tag", "chip", "label", "pill"}},
	{"badge", []string{"badge", "dot", "badgedrawable"}}, // a count / dot on an icon, not a tag
	{"cell", []string{"cell", "listitem", "row", "item", "inforow", "commonlistitemview"}},
	{"card", []string{"card", "panel", "section", "cellgroup", "grouplistview"}},
	{"navbar", []string{"navbar", "topbar", "appbar", "titlebar", "toolbar", "header", "topappbar"}},
	{"price", []string{"price", "amount", "money", "pricetext"}},
	{"dialog", []string{"dialog", "modal", "alert", "confirm"}},
	{"toast", []string{"toast", "snackbar", "notify", "tips", "tipdialog", "hud"}},
	{"popup", []string{"popup", "drawer", "popover", "floatingpanel", "floatlayout"}},
	{"sheet", []string{"sheet", "actionsheet", "bottomsheet", "sharesheet", "bottomsheetdialog"}},
	{"tabs", []string{"tabs", "tab", "segment", "segmented", "tabrow", "tabsegment"}},
	{"tabbar", []string{"tabbar", "navigationbar", "bottomnavigation", "bottomnavigationview", "bottomnav", "navigationrail", "navigationsuite", "navigationitem", "navigationbaritem"}},
	{"input", []string{"input", "field", "edittext", "textfield", "textarea", "inputitem", "easyinput"}},
	{"search", []string{"search", "searchbar", "searchview", "searchfield"}},
	{"carousel", []string{"swipe", "swiper", "carousel", "banner", "viewpager"}},
	{"icon", []string{"icon"}},
	{"image", []string{"image", "img", "avatar", "imageview"}},
	{"empty", []string{"empty", "emptyview", "nodata", "placeholder", "error", "errorpanel", "errorview", "errorblock", "errorstate", "multistateview", "stateview", "statelayout", "statusview", "multistatelayout"}},
	{"loading", []string{"loading", "spinner", "progress", "skeleton"}},
	{"list", []string{"list", "recyclerview", "listview", "lazycolumn"}},
	{"divider", []string{"divider", "separator", "line"}},
	{"notice", []string{"noticebar", "notice", "alertbar", "tipsbar"}},
	{"switch", []string{"switch", "toggle", "toggleswitch", "switchcompat", "materialswitch", "switchbutton", "switchview"}},
	{"chart", []string{"chart", "kline", "klinechart", "candlestick", "linechart", "barchart", "piechart", "timeshare", "timesharing", "minutechart", "trendchart", "sparkline"}},
	{"checkbox", []string{"checkbox", "checkboxgroup", "check", "materialcheckbox", "checklist", "datacheckbox"}},
	{"radio", []string{"radio", "radiobutton", "radiogroup", "materialradiobutton"}},
	{"stepper", []string{"stepper", "counter", "quantityselector", "numberbox", "inputnumber", "numberpicker", "numberstepper"}},
	{"countdown", []string{"countdown", "timer"}},
	{"goods", []string{"goods", "product", "sku", "goodscard", "productcard"}},
	{"actionbar", []string{"actionbar", "goodsaction", "goodsnav", "bottombar", "bottomactionbar", "footerbar"}},
}

// noise words that never decide a concept on their own
var weak = map[string]bool{"panel": true, "item": true, "row": true, "box": true, "field": true, "line": true, "label": true, "section": true, "header": true, "footer": true, "list": true}

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

var vendorPrefix = map[string]bool{"at": true, "vue": true, "react": true, "material": true, "qmui": true, "van": true, "nut": true, "ds": true, "legacy": true, "app": true, "common": true, "base": true, "my": true, "custom": true, "ui": true, "md": true, "t": true}

var viewSuffix = map[string]bool{"view": true, "layout": true, "widget": true, "component": true, "comp": true}

// projectPrefix: a short project abbreviation in front of a longer name (QDTagView, GJButton, XYCell).
func projectPrefix(words []string) bool {
	w := words[0]
	if len(words) < 2 || len(w) < 2 || len(w) > 3 {
		return false
	}
	for _, c := range concepts {
		for _, cw := range c.Words {
			if strings.HasPrefix(cw, w) {
				return false // btn, nav, tag, img … are words, not prefixes
			}
		}
	}
	return !strings.ContainsAny(w, "aeiou") // abbreviations rarely have vowels (qd, gj, xy)
}

// infrastructure suffixes: layout helpers, managers, adapters … are not UI components of a concept
var infraSuffix = map[string]bool{"layout": true, "helper": true, "manager": true, "util": true, "utils": true, "controller": true, "background": true, "behavior": true, "decoration": true, "adapter": true, "holder": true, "mask": true, "root": true, "wrapper": true, "container": true, "delegate": true, "provider": true, "state": true}

// Chinese layer / component names are common in Figma files of Chinese teams. Longest keyword wins.
var zhConcepts = []struct{ kw, concept string }{
	{"标签页", "tabs"}, {"标签栏", "tabbar"}, {"底部导航", "tabbar"}, {"底部菜单", "tabbar"}, {"选项卡", "tabs"}, {"导航栏", "navbar"}, {"标题栏", "navbar"}, {"顶栏", "navbar"},
	{"操作栏", "actionbar"}, {"底部栏", "actionbar"}, {"提交栏", "actionbar"}, {"吸底", "actionbar"},
	{"通知栏", "notice"}, {"公告", "notice"}, {"提示条", "notice"}, {"倒计时", "countdown"},
	{"单元格", "cell"}, {"列表项", "cell"}, {"信息行", "cell"}, {"分割线", "divider"}, {"分隔线", "divider"},
	{"商品卡", "goods"}, {"商品", "goods"}, {"价格", "price"}, {"金额", "price"}, {"按钮", "button"},
	{"标签", "tag"}, {"徽标", "badge"}, {"角标", "badge"}, {"卡片", "card"}, {"弹窗", "dialog"}, {"对话框", "dialog"},
	{"轻提示", "toast"}, {"吐司", "toast"}, {"弹出层", "popup"}, {"动作面板", "sheet"}, {"底部弹窗", "sheet"}, {"分享面板", "sheet"}, {"面板", "popup"}, {"输入框", "input"}, {"搜索", "search"}, {"K线", "chart"}, {"k线", "chart"}, {"分时", "chart"}, {"走势图", "chart"}, {"图表", "chart"},
	{"图片", "image"}, {"头像", "image"}, {"空状态", "empty"}, {"加载", "loading"}, {"开关", "switch"}, {"复选框", "checkbox"}, {"勾选", "checkbox"},
	{"单选", "radio"}, {"步进器", "stepper"}, {"轮播", "carousel"}, {"图标", "icon"}, {"列表", "list"},
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
	for _, w := range words {
		if w == "preference" || w == "preferences" {
			return "" // androidx.preference rows belong to the settings framework, not to page UI
		}
	}
	for len(words) > 1 && (vendorPrefix[words[0]] || projectPrefix(words)) {
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
	// a compound concept wrapped as a View: ProductCardView, ActionBarLayout, TabLayout, AppBarLayout
	if n := len(words); n >= 2 && viewSuffix[words[n-1]] {
		inner := strings.Join(words[:n-1], "")
		for _, c := range concepts {
			for _, w := range c.Words {
				if inner == w {
					return c.Name
				}
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
	for e := range words {
		// longest n-gram ending here: top+app+bar → topappbar, navigation+bar → navigationbar
		for n := 3; n >= 1; n-- {
			if e-n+1 < 0 {
				continue
			}
			w := strings.Join(words[e-n+1:e+1], "")
			w2 := strings.TrimSuffix(w, "s")
			hit := ""
			for _, c := range concepts {
				for _, cw := range c.Words {
					if w == cw || w2 == cw {
						hit = c.Name
						if n == 1 && weak[cw] {
							hit = "~" + c.Name
						}
					}
				}
			}
			if hit == "" {
				continue
			}
			if strings.HasPrefix(hit, "~") {
				weakBest = hit[1:]
			} else {
				best = hit
			}
			break
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
		// H5 feature components (SelectUserModal) and names tied to the concept only by a weak word
		// (JFlyLineMap, an echarts line chart) are not alternative implementations to converge
		if a.Concept == "" || a.Kind == "page-local" && a.Usage == 0 || featureOnly(a) || a.WrappedBy != "" ||
			a.Kind != "ui-lib" && techOf(a) == "" && (featureNamed(a, a.Concept) || weakOnly(a.Name, a.Concept)) {
			continue // a feature's own component (SnackItem(snack: Snack), a page helper) is not a duplicate to converge
		}
		by[a.Concept+techOf(a)] = append(by[a.Concept+techOf(a)], a)
	}
	// An unused library component of the same concept is the classic "built our own wheel" case.
	for _, l := range libAll {
		if l.Usage > 0 || l.Concept == "" || l.Source == "Android SDK widgets" {
			continue
		}
		k := l.Concept + techOf(l)
		hasProject := false
		for _, a := range by[k] {
			// a feature component (VueImageVerify, a captcha) is not a home-made Image
			hasProject = hasProject || a.Kind != "ui-lib" && !featureNamed(a, a.Concept)
		}
		if hasProject && strings.EqualFold(l.Name, conceptCanonical(l.Concept, by[k])) {
			by[k] = append(by[k], l)
		}
	}
	var out []Cluster
	setFamilyPrefixes(assets)
	for c, as := range by {
		as = collapseFamilies(dropSubclasses(as))
		if len(as) < 2 || oneFamily(as) {
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
		// same preference as match and the page check: established in this project first, then canonical
		sort.SliceStable(as, func(i, j int) bool { return preference(as[i], c) > preference(as[j], c) })
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
	if builtinLib(a.Source) {
		w = 100 // framework widgets (Android SDK, uni-app <image>) are what projects wrap, never the recommendation
	}
	if (a.Kind == "composable" || a.Kind == "custom-view") && inDesignSystem(a.Source) {
		w = 3000 // an Android design-system module (core/designsystem, ds, uikit)
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

var dsPathRe = regexp.MustCompile(`/(ds|design-?system|ui-?kit|core-?ui|base-?ui|widget-?lib|ui/components?)/`)

// inDesignSystem: the file lives in a design-system / UI-kit module or folder.
func inDesignSystem(rel string) bool {
	return dsPathRe.MatchString("/" + strings.ToLower(filepathToSlash(rel)))
}

func filepathToSlash(p string) string { return strings.ReplaceAll(p, "\\", "/") }

func kindLabel(k string) string {
	return map[string]string{"ui-lib": "组件库", "design-system": "设计系统", "business": "业务组件", "page-local": "页面内私有", "custom-view": "自定义 View", "view-helper": "界面封装类", "composable": "Composable", "convention": "项目惯例"}[k]
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

// featureOnly: a component made for one feature — it needs an app data model, or it is only used
// inside the file that defines it.
func featureOnly(a Asset) bool {
	if hasModelParam(a.Required) || len(domainRequired(a.Required)) >= 3 && !inDesignSystem(a.Source) {
		return true
	}
	if (a.Kind == "composable" || a.Kind == "custom-view") && len(a.UsageFiles) > 0 && !inDesignSystem(a.Source) {
		for _, f := range a.UsageFiles {
			if f != a.Source {
				return false
			}
		}
		return true
	}
	return false
}

// oneFamily: every member is a design-system variant of one family (NiaButton / NiaTextButton /
// NiaOutlinedButton) — deliberate styles, not duplicates.
func oneFamily(as []Asset) bool {
	prefix := ""
	for _, a := range as {
		ws := splitWords(a.Name)
		if a.Kind == "ui-lib" || len(ws) < 2 || !famPrefixes[ws[0]] || !inDesignSystem(a.Source) {
			return false
		}
		if prefix != "" && ws[0] != prefix {
			return false
		}
		prefix = ws[0]
	}
	return true
}

// preference orders alternatives of one concept: what the project has established (≥3 uses), then
// anything it uses, then the concept's canonical component (Cell, Tag, NavBar), then kind rank.
func preference(a Asset, concept string) int {
	p := rank(a)
	switch {
	case a.Usage >= 3:
		p += 3000
	case a.Usage > 0:
		p += 1000
	}
	for _, c := range concepts {
		if c.Name == concept && normName(a.Name) == c.Words[0] {
			p += 1500
		}
	}
	// a feature-specific member (@layout/feed_statistics_card: ≥2 feature words, as in match) is never the convergence
	// target over a generic one
	if a.Kind != "ui-lib" && !inDesignSystem(a.Source) && featureNamed(a, strings.TrimSuffix(strings.TrimSuffix(concept, "（View）"), "（Compose）")) {
		p -= 2500
	}
	return p
}

// featureNamed: the name says the component is made for one feature (@layout/feed_statistics_card,
// VueImageVerify): ≥2 words beyond its concept, or a purpose word. Judged against the concept the name
// itself says (NavBar.vue that is really a tab bar is still a generic name).
func featureNamed(a Asset, concept string) bool {
	if c := ConceptOf(a.Name); c != "" {
		concept = c
	}
	ws := specialWords(a.Name, concept)
	for _, w := range ws {
		if purposeWords[w] {
			return true
		}
	}
	return len(ws) >= 2
}

// dropSubclasses removes members built on another member: custom views extending its widget
// (SquareImageView extends AppCompatImageView, a member ImageView) and H5 components whose markup is
// that library component (GhostButton → <van-button plain>). Specialisations, not second implementations.
func dropSubclasses(as []Asset) []Asset {
	var out []Asset
	for _, a := range as {
		sub := false
		if a.Extends != "" {
			base := strings.TrimPrefix(a.Extends, "AppCompat")
			for _, b := range as {
				if (b.Name != a.Name || b.Kind != a.Kind) && (b.Name == a.Extends || b.Name == base) {
					sub = true
				}
			}
		}
		if !sub {
			out = append(out, a)
		}
	}
	return out
}

// techOf separates Android View and Compose implementations: they are not interchangeable, so a
// duplicate cluster never asks an XML page to converge on a composable. "" for H5.
func techOf(a Asset) string {
	switch {
	case a.Kind == "composable" || isComposeLib(a.Source):
		return "（Compose）"
	case a.Kind == "custom-view" || a.Kind == "view-helper" || strings.HasPrefix(a.Name, "@layout/") || a.Source == "Android SDK widgets" || a.Source == "Material Components" || a.Source == "QMUI" ||
		a.Source == "AppCompat" || a.Kind == "ui-lib" && strings.Count(a.Source, ".") >= 1 && !strings.Contains(a.Source, " "):
		return "（View）"
	}
	return ""
}
