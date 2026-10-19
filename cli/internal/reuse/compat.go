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
		"ExtendedFloatingActionButton": "1.1.0",
	},
	// components added in Vant 4 (projects on Vant 2 / 3 do not have them)
	"Vant": {
		"BackTop": "4.0.0", "DatePicker": "4.0.0", "TimePicker": "4.0.0", "PickerGroup": "4.0.0", "FloatingBubble": "4.0.0",
		"FloatingPanel": "4.0.0", "Highlight": "4.0.0", "RollingText": "4.0.0", "TextEllipsis": "4.0.0", "Watermark": "4.0.0",
		"Signature": "4.0.0", "Barrage": "4.0.0",
	},
}

// libDeprecated: component → its replacement in the same library.
var libDeprecated = map[string]map[string]string{
	"Compose Material3": {"Divider": "HorizontalDivider", "SwipeToDismiss": "SwipeToDismissBox"},
}

// materialThemed: Material Components widgets that enforce a Theme.MaterialComponents / Material3
// ancestor ("The style on this component requires your app theme to be Theme.MaterialComponents").
var materialThemed = map[string]bool{
	"MaterialButton": true, "Chip": true, "ChipGroup": true, "MaterialCardView": true, "TextInputLayout": true,
	"BottomNavigationView": true, "Slider": true, "MaterialButtonToggleGroup": true, "ExtendedFloatingActionButton": true,
	"CircularProgressIndicator": true, "LinearProgressIndicator": true, "MaterialSwitch": true,
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
		case strings.HasPrefix(cur, "Theme.MaterialComponents"), strings.HasPrefix(cur, "Theme.Material3"):
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
		if repl := libDeprecated[a.Source][a.Name]; repl != "" {
			a.Caveat = "已废弃，改用 " + repl
			continue
		}
		if a.Source == "Material Components" && materialThemed[a.Name] && themeFamily == "appcompat" {
			a.Caveat = "需要 Theme.MaterialComponents 主题，项目主题 " + theme + " 继承 AppCompat，直接用会崩溃"
		}
	}
}

// usable: a component that can be recommended without a compatibility caveat. A themed widget the
// project already uses evidently works (e.g. a local theme overlay), so usage clears the caveat.
func usable(a *Asset) bool {
	return !a.Unavailable && (a.Caveat == "" || a.Usage > 0)
}
