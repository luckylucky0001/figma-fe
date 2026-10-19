// Package audit statically checks an implemented page: which mapped components it actually uses
// versus the plan, and bad practices (raw elements, direct Material, hard-coded colors, self-built look-alikes).
package audit

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/luckylucky0001/figma-fe/cli/internal/inventory"
	"github.com/luckylucky0001/figma-fe/cli/internal/plan"
)

type Finding struct {
	Rule     string `json:"rule"`
	Severity string `json:"severity"` // error counts as violation, warn does not
	File     string `json:"file"`
	Line     int    `json:"line"`
	Text     string `json:"text"`
	Hint     string `json:"hint,omitempty"`
}

type Usage struct {
	Component string `json:"component"`
	Expected  int    `json:"expected"`
	Actual    int    `json:"actual"`
}

type Result struct {
	Platform   string    `json:"platform"`
	Files      []string  `json:"files"`
	Usage      []Usage   `json:"usage"`
	ReuseRate  float64   `json:"reuseRate"`
	FigmaIDs   int       `json:"figmaIds"`
	Violations int       `json:"violations"`
	Findings   []Finding `json:"findings"`
}

var (
	jsxTagRe   = regexp.MustCompile(`<([A-Z]\w*)[\s/>]`)
	rawH5Re    = regexp.MustCompile(`<(button|input|select|textarea|img)\b`)
	hexRe      = regexp.MustCompile(`#[0-9a-fA-F]{6}\b|#[0-9a-fA-F]{3}\b|rgba?\([^)]*\)`)
	pxRe       = regexp.MustCompile(`\b\d+(\.\d+)?px\b`)
	figmaIDH5  = regexp.MustCompile(`data-figma-id=`)
	ktCallRe   = regexp.MustCompile(`\b([A-Z]\w*)\s*\(`)
	ktM3Import = regexp.MustCompile(`import\s+androidx\.compose\.material3\.(\w+)`)
	ktColorRe  = regexp.MustCompile(`Color\(0x[0-9A-Fa-f]{6,8}\)|Color\(\s*\d+\s*,`)
	ktDpRe     = regexp.MustCompile(`\b\d+(\.\d+)?\.(dp|sp)\b`)
	figmaIDKt  = regexp.MustCompile(`figmaId\(`)
	legacyRe   = regexp.MustCompile(`\bLegacy\w+`)
	classRe    = regexp.MustCompile(`className=["{][^"}]*\b(btn|button|tag|badge|price|cell|card|navbar|nav-bar|notice|divider|action-?bar)\b`)
	cssSelRe   = regexp.MustCompile(`\.([\w-]*(btn|button|tag|badge|price|cell|card|navbar|notice|divider|actionbar)[\w-]*)\s*\{`)
	ktLocalFun = regexp.MustCompile(`@Composable\s+(?:private\s+|internal\s+)?fun\s+(\w*(Button|Tag|Badge|Price|Cell|Card|TopBar|NavBar|Banner|Notice|Divider|BottomBar|ListItem)\w*)\s*\(`)
)

// Run audits all source files under dir for the given platform against the plan.
// dups maps business look-alike components (e.g. StatusTag) to the design-system component they duplicate.
func Run(platform, repoRoot, dir string, p *plan.Plan, inv *inventory.Inventory, dups map[string]string) (*Result, error) {
	res := &Result{Platform: platform}
	exts := []string{".tsx", ".ts", ".css"}
	if platform == "android" {
		exts = []string{".kt"}
	}
	var files []string
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			for _, e := range exts {
				if strings.HasSuffix(path, e) {
					files = append(files, path)
				}
			}
		}
		return nil
	})
	sort.Strings(files)
	known := map[string]bool{}
	for _, c := range inv.Components {
		if c.Kind != "business" {
			known[c.Name] = true
		}
	}
	dupRe := map[string]*regexp.Regexp{}
	for b := range dups {
		if platform == "h5" {
			dupRe[b] = regexp.MustCompile(`<` + b + `[\s/>]`)
		} else {
			dupRe[b] = regexp.MustCompile(`\b` + b + `\s*\(`)
		}
	}
	actual := map[string]int{}
	add := func(f Finding) {
		res.Findings = append(res.Findings, f)
		if f.Severity == "error" {
			res.Violations++
		}
	}
	for _, f := range files {
		b, _ := os.ReadFile(f)
		src := string(b)
		relf, _ := filepath.Rel(repoRoot, f)
		res.Files = append(res.Files, relf)
		lines := strings.Split(src, "\n")
		at := func(idx int) (int, string) {
			ln := strings.Count(src[:idx], "\n")
			return ln + 1, strings.TrimSpace(lines[ln])
		}
		for b, re := range dupRe {
			for _, m := range re.FindAllStringIndex(src, -1) {
				if strings.Contains(src[max(0, m[0]-6):m[0]], "fun ") {
					continue // the page defines it itself; lookalike rule covers that
				}
				ln, t := at(m[0])
				add(Finding{"business-duplicate", "error", relf, ln, t, "用了业务重复组件 " + b + "，应改用设计系统组件 " + dups[b]})
			}
		}
		if platform == "h5" {
			for _, m := range jsxTagRe.FindAllStringSubmatchIndex(src, -1) {
				name := src[m[2]:m[3]]
				if known[name] {
					actual[name]++
				}
			}
			for _, m := range rawH5Re.FindAllStringIndex(src, -1) {
				ln, t := at(m[0])
				add(Finding{"raw-element", "error", relf, ln, t, "改用组件库组件（Button 等）"})
			}
			for _, m := range hexRe.FindAllStringIndex(src, -1) {
				ln, t := at(m[0])
				add(Finding{"hardcoded-color", "error", relf, ln, t, "改用 var(--ds-color-*) token"})
			}
			for _, m := range classRe.FindAllStringSubmatchIndex(src, -1) {
				ln, t := at(m[0])
				add(Finding{"lookalike", "warn", relf, ln, t, "疑似自己拼的组件（类名 " + src[m[2]:m[3]] + "）：组件库里有对应组件时应直接复用"})
			}
			for _, m := range cssSelRe.FindAllStringSubmatchIndex(src, -1) {
				ln, t := at(m[0])
				add(Finding{"lookalike", "warn", relf, ln, t, "疑似自造组件样式 ." + src[m[2]:m[3]]})
			}
			if n := len(pxRe.FindAllStringIndex(src, -1)); n > 0 {
				add(Finding{"raw-number", "warn", relf, 0, "", "出现 " + itoa(n) + " 处裸 px 数值；间距 / 圆角 / 字号优先用 token"})
			}
			res.FigmaIDs += len(figmaIDH5.FindAllStringIndex(src, -1))
		} else {
			for _, m := range ktCallRe.FindAllStringSubmatchIndex(src, -1) {
				name := src[m[2]:m[3]]
				if known[name] {
					actual[name]++
				}
			}
			for _, m := range ktM3Import.FindAllStringSubmatchIndex(src, -1) {
				name := src[m[2]:m[3]]
				if inventory.M3Components[name] {
					ln, t := at(m[0])
					add(Finding{"direct-material3", "error", relf, ln, t, "页面不能直接用 material3." + name + "，改用 com.demo.ds 组件"})
				}
			}
			for _, m := range ktColorRe.FindAllStringIndex(src, -1) {
				ln, t := at(m[0])
				add(Finding{"hardcoded-color", "error", relf, ln, t, "改用 DsColors.*"})
			}
			for _, m := range legacyRe.FindAllStringIndex(src, -1) {
				ln, t := at(m[0])
				add(Finding{"legacy-view", "error", relf, ln, t, "新页面不能用 XML View 存量组件 Legacy*"})
			}
			for _, m := range ktLocalFun.FindAllStringSubmatchIndex(src, -1) {
				ln, t := at(m[0])
				add(Finding{"lookalike", "warn", relf, ln, t, "页面里自定义了 " + src[m[2]:m[3]] + "：组件库里有对应组件时应直接复用"})
			}
			if n := len(ktDpRe.FindAllStringIndex(src, -1)); n > 0 {
				add(Finding{"raw-number", "warn", relf, 0, "", "出现 " + itoa(n) + " 处裸 dp/sp 数值；优先用 DsSpacing / DsTypography"})
			}
			res.FigmaIDs += len(figmaIDKt.FindAllStringIndex(src, -1))
		}
	}
	exp, hit := 0, 0
	comps := make([]string, 0, len(p.Counts))
	for c := range p.Counts {
		comps = append(comps, c)
	}
	sort.Strings(comps)
	for _, c := range comps {
		e := p.Counts[c]
		a := actual[c]
		res.Usage = append(res.Usage, Usage{c, e, a})
		exp += e
		hit += min(e, a)
		if a < e {
			add(Finding{"missing-reuse", "warn", "", 0, c, "方案要求 " + itoa(e) + " 处 " + c + "，代码里只有 " + itoa(a) + " 处"})
		}
	}
	if exp > 0 {
		res.ReuseRate = float64(hit) / float64(exp)
	}
	return res, nil
}

func itoa(n int) string { return strconv.Itoa(n) }
