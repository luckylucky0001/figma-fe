package reuse

import (
	"regexp"
	"sort"
	"strings"
)

// Conventions are the project's established ways of drawing a concept without a component — a
// divider as a 1px View or as the qmui_bottomDividerHeight attribute of a QMUI layout, Vant's
// van-hairline classes. When a Figma element has no component, the convention is still "reuse,
// don't invent": match offers it before concluding that a new component is needed.

var (
	hairlineViewRe = regexp.MustCompile(`<View\b[^>]*?android:layout_height="(?:1px|0\.5dp|1dp|0\.5px|@dimen/\w*(?:divider|line|hairline)\w*)"[^>]*>`)
	dividerAttrRe  = regexp.MustCompile(`\b(\w+):(\w*[Dd]ivider\w*)="`)
	hairlineClsRe  = regexp.MustCompile(`(?:^|\s)(van-hairline(?:--[\w-]+)?|[\w-]*(?:divider|hairline|split-line)[\w-]*)`)
)

const minConvention = 3 // fewer uses than this is not a convention

type convAcc struct {
	label    string
	n        int
	ex, exAt string
}

func (c *convAcc) add(ex, at string) {
	c.n++
	if c.ex == "" {
		ex = clip(ex, 140)
		c.ex, c.exAt = ex, at
	}
}

func conventionsAndroid(layouts map[string]string) []Asset {
	acc := map[string]*convAcc{}
	get := func(key, label string) *convAcc {
		if acc[key] == nil {
			acc[key] = &convAcc{label: label}
		}
		return acc[key]
	}
	rels := make([]string, 0, len(layouts))
	for rel := range layouts {
		rels = append(rels, rel)
	}
	sort.Strings(rels) // the example shown is the first in path order, the same on every run
	for _, rel := range rels {
		text := layouts[rel]
		for _, m := range hairlineViewRe.FindAllStringIndex(text, -1) {
			get("view", "1px View 细线").add(oneLine(text[m[0]:m[1]]), rel+":"+itoa(lineAt(text, m[0])))
		}
		for _, m := range xmlTagRe.FindAllStringSubmatchIndex(text, -1) {
			end := strings.IndexByte(text[m[0]:], '>')
			if end < 0 {
				continue
			}
			open := text[m[0] : m[0]+end+1]
			tag := text[m[2]:m[3]]
			tag = tag[strings.LastIndex(tag, ".")+1:]
			if tag == "LinearLayout" || ConceptOf(tag) != "" {
				continue // showDividers / a tab bar's own divider are not a standalone divider look
			}
			for _, am := range dividerAttrRe.FindAllStringSubmatch(open, -1) {
				if am[1] == "android" || am[1] == "tools" {
					continue
				}
				// group the family (qmui_topDividerHeight, qmui_bottomDividerColor …) per component
				get(tag, tag+" 自带分割线（app:"+am[2]+" 等属性）").add(oneLine(open), rel+":"+itoa(lineAt(text, m[0])))
				break
			}
		}
	}
	return pickConventions(acc)
}

func conventionsH5(files []srcFile) []Asset {
	acc := map[string]*convAcc{}
	for _, f := range files {
		src := f.markup
		for _, m := range classAttrRe.FindAllStringSubmatchIndex(src, -1) {
			for _, cm := range hairlineClsRe.FindAllStringSubmatch(src[m[2]:m[3]], -1) {
				key, label := cm[1], "CSS 类 ."+cm[1]
				if strings.HasPrefix(key, "van-hairline") {
					key, label = "van-hairline", "Vant 细线类 van-hairline--*"
				}
				if acc[key] == nil {
					acc[key] = &convAcc{label: label}
				}
				ls := strings.LastIndex(src[:m[0]], "\n") + 1
				le := strings.IndexByte(src[m[0]:], '\n')
				if le < 0 {
					le = len(src) - m[0]
				}
				acc[key].add(strings.TrimSpace(src[ls:m[0]+le]), f.rel+":"+itoa(lineAt(src, m[0])+markupLine(f)))
			}
		}
	}
	return pickConventions(acc)
}

func pickConventions(acc map[string]*convAcc) []Asset {
	var out []Asset
	for _, c := range acc {
		if c.n >= minConvention {
			out = append(out, Asset{Name: c.label, Kind: "convention", Source: "项目惯例", Concept: "divider", Usage: c.n, Example: c.ex, ExampleAt: c.exAt})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Usage != out[j].Usage {
			return out[i].Usage > out[j].Usage
		}
		return out[i].Name < out[j].Name
	})
	return out
}

func lineAt(s string, idx int) int { return strings.Count(s[:idx], "\n") + 1 }

func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }
