package reuse

import (
	"regexp"
	"sort"
	"strings"
)

// Variant → prop mapping: knowing *which* component to reuse is half the answer; the other half is
// how to configure it for the design (Figma Type=Warning → type="warning", Arrow=On → is-link).
// Enum-like props come from the project's own declarations (TS / Vue union literals, Android
// styleable enums, Kotlin enum parameters) and from a small catalog of common library components.

var (
	unionPropRe = regexp.MustCompile(`(\w+)\??\s*:\s*((?:['"][\w-]+['"]\s*\|\s*)+['"][\w-]+['"])`)
	unionValRe  = regexp.MustCompile(`['"]([\w-]+)['"]`)
	enumAttrRe  = regexp.MustCompile(`(?s)<attr\s+name="(\w+)"[^>]*>(.*?)</attr>`)
	enumValRe   = regexp.MustCompile(`<(?:enum|flag)\s+name="(\w+)"`)
	ktEnumRe    = regexp.MustCompile(`enum\s+class\s+(\w+)[^{]*\{([^}]*)\}`)
	ktIdentRe   = regexp.MustCompile(`^[A-Z]\w*$`)
	parenRe     = regexp.MustCompile(`\([^()]*\)`)
	tsPropsRe   = regexp.MustCompile(`(?s)(?:interface|type)\s+\w*Props\b[^{]*\{(.*?)\n\}`)
	tsMemberRe  = regexp.MustCompile(`(?m)^\s*(?:readonly\s+)?(\w+)\??\s*:`)
)

// libEnums: enum props of common library components (from their docs).
var libEnums = map[string]map[string][]string{
	"Vant/Button":           {"type": {"primary", "success", "warning", "danger", "default"}, "size": {"large", "normal", "small", "mini"}},
	"Vant/Tag":              {"type": {"primary", "success", "warning", "danger"}, "size": {"medium", "large"}},
	"Vant/NoticeBar":        {"mode": {"closeable", "link"}},
	"Vant/Cell":             {"size": {"large", "normal"}},
	"antd-mobile/Button":    {"color": {"default", "primary", "success", "warning", "danger"}, "fill": {"solid", "outline", "none"}, "size": {"mini", "small", "middle", "large"}},
	"antd-mobile/Tag":       {"color": {"default", "primary", "success", "warning", "danger"}, "fill": {"solid", "outline"}},
	"antd-mobile/NoticeBar": {"color": {"default", "alert", "error", "info"}},
	"NutUI/Button":          {"type": {"default", "primary", "info", "success", "warning", "danger"}, "size": {"large", "normal", "small"}},
	"NutUI/Tag":             {"type": {"default", "primary", "success", "warning", "danger"}},
	"TDesign Mobile/Button": {"theme": {"default", "primary", "danger", "light"}, "size": {"small", "medium", "large"}},
	"TDesign Mobile/Tag":    {"theme": {"default", "primary", "warning", "danger", "success"}},
}

// libFlags: boolean / presence props of common library components.
var libFlags = map[string][]string{
	"Vant/Cell":            {"title", "value", "label", "is-link", "center", "border", "icon"},
	"Vant/Button":          {"plain", "round", "block", "disabled", "loading"},
	"Vant/Tag":             {"plain", "round", "mark"},
	"Vant/NavBar":          {"title", "left-arrow", "left-text", "right-text", "fixed"},
	"Vant/CellGroup":       {"title", "inset", "border"},
	"antd-mobile/ListItem": {"title", "description", "extra", "arrow", "clickable", "prefix"},
	"antd-mobile/Card":     {"title", "extra"},
	"antd-mobile/NavBar":   {"back", "right"},
}

// flag synonyms: Figma boolean axis word → prop words that implement it
var flagSyn = map[string][]string{
	"arrow": {"arrow", "islink", "showarrow", "rightarrow", "accessory", "chevron"},
	"label": {"label", "description", "subtitle", "desc", "detail"},
	"title": {"title", "header", "heading"},
	"icon":  {"icon", "lefticon", "prefix"},
	"close": {"closeable", "closable", "close"},
	"back":  {"back", "leftarrow", "showback"},
}

// value synonyms: a Figma variant value → equivalent enum values in code
var valSyn = map[string][]string{
	"default": {"default", "normal", "secondary"},
	"medium":  {"medium", "normal", "middle", "md", "default"},
	"small":   {"small", "sm", "mini"},
	"large":   {"large", "lg", "big"},
	"warning": {"warning", "warn", "alert"},
	"danger":  {"danger", "error", "critical"},
	"info":    {"info", "primary"},
	"primary": {"primary", "brand"},
	"success": {"success", "positive"},
}

func unionEnums(src string) map[string][]string {
	out := map[string][]string{}
	for _, m := range unionPropRe.FindAllStringSubmatch(src, -1) {
		for _, v := range unionValRe.FindAllStringSubmatch(m[2], -1) {
			out[m[1]] = append(out[m[1]], v[1])
		}
	}
	return out
}

func styleableEnums(body string) map[string][]string {
	out := map[string][]string{}
	for _, m := range enumAttrRe.FindAllStringSubmatch(body, -1) {
		for _, v := range enumValRe.FindAllStringSubmatch(m[2], -1) {
			out[m[1]] = append(out[m[1]], v[1])
		}
	}
	return out
}

func kotlinEnums(src string) map[string][]string {
	out := map[string][]string{}
	for _, m := range ktEnumRe.FindAllStringSubmatch(src, -1) {
		body := m[2]
		for parenRe.MatchString(body) {
			body = parenRe.ReplaceAllString(body, "") // Info(Color(0xFF…), …) → Info
		}
		if i := strings.IndexByte(body, ';'); i >= 0 {
			body = body[:i] // members after the constants
		}
		for _, v := range strings.FieldsFunc(body, func(r rune) bool { return r == ',' || r == '\n' || r == ' ' || r == '\t' }) {
			if ktIdentRe.MatchString(v) {
				out[m[1]] = append(out[m[1]], m[1]+"."+v)
			}
		}
	}
	return out
}

// tsProps: member names of a React component's Props interface / type.
func tsProps(src string) []string {
	var out []string
	for _, m := range tsPropsRe.FindAllStringSubmatch(src, -1) {
		for _, x := range tsMemberRe.FindAllStringSubmatch(m[1], -1) {
			if !contains(out, x[1]) {
				out = append(out, x[1])
			}
		}
	}
	return out
}

func normProp(s string) string {
	return strings.ToLower(strings.NewReplacer("-", "", "_", "").Replace(s))
}

// PropHints maps the Figma variants of an item onto the chosen component's props.
func PropHints(variants []string, a Asset) []string {
	enums := map[string][]string{}
	for k, v := range a.Enums {
		enums[k] = v
	}
	for k, v := range libEnums[a.Source+"/"+a.Name] {
		enums[k] = v
	}
	flags := append([]string{}, a.Props...)
	flags = append(flags, libFlags[a.Source+"/"+a.Name]...)
	axes := map[string][]string{}
	var order []string
	for _, v := range variants {
		kv := strings.SplitN(v, "=", 2)
		if len(kv) != 2 {
			continue
		}
		if _, ok := axes[kv[0]]; !ok {
			order = append(order, kv[0])
		}
		axes[kv[0]] = append(axes[kv[0]], kv[1])
	}
	sort.Strings(order)
	var out []string
	usedProp := map[string]bool{}
	for _, axis := range order {
		vals := axes[axis]
		if onOff(vals) {
			if p := flagProp(axis, flags); p != "" {
				out = append(out, axis+"=On → "+flagAssign(a, p)+"；"+axis+"=Off → 不传")
			} else {
				out = append(out, axis+"：组件没有对应属性，需确认（用插槽 / 自行组合）")
			}
			continue
		}
		// the enum prop that resolves most of this axis' values
		best, bestHits := "", 0
		var props []string
		for p := range enums {
			props = append(props, p)
		}
		sort.Strings(props)
		for _, p := range props {
			if usedProp[p] {
				continue
			}
			hits := 0
			for _, v := range vals {
				if resolve(v, enums[p]) != "" {
					hits++
				}
			}
			if normProp(p) == normProp(axis) {
				hits++ // same name breaks ties
			}
			if hits > bestHits {
				best, bestHits = p, hits
			}
		}
		if best == "" {
			out = append(out, axis+"："+strings.Join(vals, " / ")+" 没有对应的枚举属性，需确认（可能要用 color / class 等实现）")
			continue
		}
		usedProp[best] = true
		var parts []string
		for _, v := range vals {
			if r := resolve(v, enums[best]); r != "" {
				parts = append(parts, axis+"="+v+" → "+assign(a, best, r))
			} else {
				parts = append(parts, axis+"="+v+" → "+best+" 没有对应值，需确认")
			}
		}
		out = append(out, strings.Join(parts, "；"))
	}
	return out
}

func onOff(vals []string) bool {
	for _, v := range vals {
		if l := strings.ToLower(v); l != "on" && l != "off" && l != "true" && l != "false" && l != "yes" && l != "no" {
			return false
		}
	}
	return true
}

func flagProp(axis string, props []string) string {
	words := flagSyn[strings.ToLower(axis)]
	if words == nil {
		words = []string{normProp(axis)}
	}
	for _, w := range words {
		for _, p := range props {
			if normProp(p) == w {
				return p
			}
		}
	}
	for _, w := range words { // showChevron, hasArrow, withIcon
		for _, p := range props {
			if n := normProp(p); len(w) >= 4 && strings.HasSuffix(n, w) && len(n) <= len(w)+4 {
				return p
			}
		}
	}
	return ""
}

func resolve(v string, enum []string) string {
	l := strings.ToLower(v)
	cands := append([]string{l}, valSyn[l]...)
	for _, c := range cands {
		for _, e := range enum {
			if strings.ToLower(e[strings.LastIndex(e, ".")+1:]) == c {
				return e
			}
		}
	}
	return ""
}

// assign renders prop = value in the component's own syntax.
func assign(a Asset, prop, val string) string {
	switch {
	case a.Kind == "composable" || strings.Contains(val, "."):
		return prop + " = " + val // tone = DsBannerTone.Warning
	case a.Kind == "custom-view":
		return "app:" + prop + "=\"" + val + "\""
	}
	return prop + "=\"" + val + "\""
}

func flagAssign(a Asset, prop string) string {
	switch a.Kind {
	case "composable":
		return prop + " = true / 传入内容"
	case "custom-view":
		return "app:" + prop
	}
	return prop
}
