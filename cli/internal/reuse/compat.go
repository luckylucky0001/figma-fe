package reuse

import (
	"regexp"
	"strconv"
	"strings"
)

// Compatibility of library components with the project as it is: a component the installed
// library version does not have cannot be reused (it will not compile), and a Material
// Components widget inflated under an AppCompat theme crashes at runtime. Both were real
// mistakes an agent made on a brownfield project, so match / page checks must know them.

// libMinVersion: component → first library version that ships it. Only versions we are sure of.
var libMinVersion = map[string]map[string]string{
	"Material Components": {
		"MaterialDivider": "1.5.0", "MaterialSwitch": "1.7.0", "ShapeableImageView": "1.2.0", "Slider": "1.2.0",
		"CircularProgressIndicator": "1.3.0", "LinearProgressIndicator": "1.3.0", "MaterialButtonToggleGroup": "1.1.0",
		"ExtendedFloatingActionButton": "1.1.0", "SearchBar": "1.8.0", "SearchView": "1.8.0",
	},
	// antd-mobile 5 is a rewrite: these exist only from 5.0
	"antd-mobile": {
		"Swiper": "5.0.0", "Input": "5.0.0", "TextArea": "5.0.0", "SpinLoading": "5.0.0", "DotLoading": "5.0.0", "Space": "5.0.0",
		"CapsuleTabs": "5.0.0", "JumboTabs": "5.0.0", "Ellipsis": "5.0.0", "FloatingPanel": "5.0.0", "Collapse": "5.0.0", "ErrorBlock": "5.0.0",
		"InfiniteScroll": "5.0.0", "AutoCenter": "5.0.0", "Divider": "5.0.0", "Image": "5.0.0", "CheckList": "5.0.0", "Selector": "5.0.0", "Popup": "5.0.0", "CenterPopup": "5.0.0",
	},
	// components added in Vant 4 (projects on Vant 2 / 3 do not have them)
	"Vant": {
		"BackTop": "4.0.0", "DatePicker": "4.0.0", "TimePicker": "4.0.0", "PickerGroup": "4.0.0", "FloatingBubble": "4.0.0",
		"FloatingPanel": "4.0.0", "Highlight": "4.0.0", "RollingText": "4.0.0", "TextEllipsis": "4.0.0", "Watermark": "4.0.0",
		"Signature": "4.0.0", "Barrage": "4.0.0",
		"ActionBar": "3.0.0", "ActionBarButton": "3.0.0", "ActionBarIcon": "3.0.0", // GoodsAction* before Vant 3
	},
}

// libDeprecated: component → its replacement in the same library.
var libDeprecated = map[string]map[string]string{
	"Compose Material3": {"Divider": "HorizontalDivider", "SwipeToDismiss": "SwipeToDismissBox"},
}

// libRemoved: component → [first version without it, its successor] (renamed / removed).
var libRemoved = map[string]map[string][2]string{
	"antd-mobile": {"InputItem": {"5.0.0", "Input"}, "TextareaItem": {"5.0.0", "TextArea"}, "Carousel": {"5.0.0", "Swiper"},
		"ActivityIndicator": {"5.0.0", "SpinLoading / DotLoading"}, "WhiteSpace": {"5.0.0", "Space"}, "WingBlank": {"5.0.0", "Space"},
		"ListView": {"5.0.0", "List + InfiniteScroll"}, "Flex": {"5.0.0", "Space / Grid"}, "Accordion": {"5.0.0", "Collapse"}, "Range": {"5.0.0", "Slider"}, "Menu": {"5.0.0", "CascaderView"}},
	"Vant": {"GoodsAction": {"3.0.0", "ActionBar"}, "GoodsActionButton": {"3.0.0", "ActionBarButton"}, "GoodsActionIcon": {"3.0.0", "ActionBarIcon"}},
}

// material3Only: Material Components widgets that need a Theme.Material3 ancestor.
var material3Only = map[string]bool{"SearchBar": true, "SearchView": true, "MaterialSwitch": true}

// materialThemed: Material Components widgets that enforce a Theme.MaterialComponents / Material3
// ancestor ("The style on this component requires your app theme to be Theme.MaterialComponents").
var materialThemed = map[string]bool{
	"MaterialButton": true, "Chip": true, "ChipGroup": true, "MaterialCardView": true, "TextInputLayout": true,
	"BottomNavigationView": true, "Slider": true, "MaterialButtonToggleGroup": true, "ExtendedFloatingActionButton": true,
	"CircularProgressIndicator": true, "LinearProgressIndicator": true, "MaterialSwitch": true,
	"SearchBar": true, "SearchView": true,
	"MaterialDivider": true, // colours come from ?attr/colorOnSurface / colorOutline (Material theme attributes)
}

var (
	materialDepRe   = regexp.MustCompile(`com\.google\.android\.material:material:(\d+\.\d+\.\d+)`)
	materialTomlRe  = regexp.MustCompile(`(?m)^\s*material\s*=\s*"(\d+\.\d+\.\d+)"`)
	styleDefRe      = regexp.MustCompile(`<style\s+name="([\w.]+)"(?:\s+parent="(?:@style/|@android:style/)?([\w.]*)")?`)
	manifestThemeRe = regexp.MustCompile(`<application[^>]*android:theme="@style/([\w.]+)"`)
)

// versionLess compares dotted numeric versions ("1.4.0" < "1.5.0"); unknown means not less.
func versionLess(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < len(pa) || i < len(pb); i++ {
		var x, y int
		if i < len(pa) {
			x, _ = strconv.Atoi(pa[i])
		}
		if i < len(pb) {
			y, _ = strconv.Atoi(pb[i])
		}
		if x != y {
			return x < y
		}
	}
	return false
}

// cleanVersion turns an npm range ("^4.8.0", "~2.12.1") into its base version.
func cleanVersion(v string) string {
	v = strings.TrimLeft(strings.TrimSpace(v), "^~>=v ")
	if i := strings.IndexAny(v, " -|"); i > 0 {
		v = v[:i]
	}
	if strings.Count(v, ".") == 0 && v != "" {
		v += ".0.0"
	}
	return v
}

func materialVersion(texts ...string) string {
	for _, t := range texts {
		if m := materialDepRe.FindStringSubmatch(t); m != nil {
			return m[1]
		}
		if m := materialTomlRe.FindStringSubmatch(t); m != nil {
			return m[1]
		}
	}
	return ""
}

// appThemeFamily resolves the application theme's parent chain: "material" (MaterialComponents /
// Material3, incl. Bridge themes), "appcompat", or "" when it cannot be resolved.
func appThemeFamily(manifests, values []string) (family, theme string) {
	for _, m := range manifests {
		if mm := manifestThemeRe.FindStringSubmatch(m); mm != nil {
			theme = mm[1]
			break
		}
	}
	if theme == "" {
		return "", ""
	}
	parent := map[string]string{}
	for _, v := range values {
		for _, m := range styleDefRe.FindAllStringSubmatch(v, -1) {
			p := m[2]
			if p == "" && strings.Contains(m[1], ".") {
				p = m[1][:strings.LastIndex(m[1], ".")] // implicit parent: QMUI.NoActionBar → QMUI
			}
			if _, ok := parent[m[1]]; !ok || p != "" {
				parent[m[1]] = p
			}
		}
	}
	for cur, i := theme, 0; cur != "" && i < 30; i++ {
		switch {
		case strings.HasPrefix(cur, "Theme.Material3"):
			return "material3", theme
		case strings.HasPrefix(cur, "Theme.MaterialComponents"):
			return "material", theme
		case strings.HasPrefix(cur, "Theme.AppCompat"):
			return "appcompat", theme
		}
		cur = parent[cur]
	}
	return "", theme
}

// markCompat annotates library assets that the project cannot use as-is.
func markCompat(assets map[string]*Asset, versions map[string]string, themeFamily, theme string) {
	for _, a := range assets {
		if a.Kind != "ui-lib" {
			continue
		}
		if min := libMinVersion[a.Source][a.Name]; min != "" && versionLess(versions[a.Source], min) {
			a.Unavailable = true
			a.Caveat = "需要 " + a.Source + " ≥ " + min + "，项目是 " + versions[a.Source] + "，用了会编译失败"
			continue
		}
		if rm, ok := libRemoved[a.Source][a.Name]; ok && versions[a.Source] != "" && !versionLess(versions[a.Source], rm[0]) {
			a.Unavailable = true
			a.Caveat = a.Source + " " + rm[0] + " 起已改名为 " + rm[1]
			continue
		}
		if repl := libDeprecated[a.Source][a.Name]; repl != "" {
			a.Caveat = "已废弃，改用 " + repl
			continue
		}
		if a.Source == "Material Components" && materialThemed[a.Name] && themeFamily == "appcompat" {
			a.Caveat = "需要 Theme.MaterialComponents 主题，项目主题 " + theme + " 继承 AppCompat，直接用会崩溃"
		} else if a.Source == "Material Components" && material3Only[a.Name] && themeFamily == "material" {
			a.Caveat = "是 Material 3 组件，需要 Theme.Material3 主题，项目主题 " + theme + " 是 MaterialComponents"
		}
	}
}

// usable: a component that can be recommended without a compatibility caveat. A themed widget the
// project already uses evidently works (e.g. a local theme overlay), so usage clears the caveat.
func usable(a *Asset) bool {
	return !a.Unavailable && (a.Caveat == "" || a.Usage > 0)
}
