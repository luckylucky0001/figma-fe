package inventory

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Business components live outside the design system (e.g. h5/src/components, android app widgets).
// They are assets too: scanning them reveals duplicates of design-system components that pages keep
// copying ("参照退款详情页的写法"), which is a common root cause of poor reuse.

var (
	tsBizFunRe = regexp.MustCompile(`export\s+(?:function|const)\s+([A-Z]\w*)\s*(?:=\s*)?\(`)
	ktBizFunRe = regexp.MustCompile(`@Composable\s+(?:public\s+)?fun\s+([A-Z]\w*)\s*\(`)
)

// ScanBusinessH5 scans React components exported from dirs (relative to h5Root) and counts their usage in srcRoot.
func ScanBusinessH5(repoRoot, h5Root string, dirs []string, inv *Inventory) {
	var comps []Component
	for _, d := range dirs {
		dir := filepath.Join(h5Root, d)
		for _, f := range walkFiles(dir, []string{".tsx"}, nil) {
			src := readFile(f)
			for _, loc := range tsBizFunRe.FindAllStringSubmatchIndex(src, -1) {
				name := src[loc[2]:loc[3]]
				c := Component{Name: name, Kind: "business", File: rel(repoRoot, f), Import: "import { " + name + " } from '<" + d + ">';"}
				c.Description, c.Example = docBefore(src, loc[0])
				open := loc[1] - 1
				if end := matchClose(src, open); end > open {
					c.Props = parseInlineTSProps(src[open+1 : end])
				}
				comps = append(comps, c)
			}
		}
	}
	countUsage(filepath.Join(h5Root, "src"), []string{".tsx"}, comps, func(src, name string) int {
		if !regexp.MustCompile(`import\s*\{[^}]*\b` + name + `\b[^}]*\}\s*from`).MatchString(src) {
			return 0
		}
		return len(regexp.MustCompile(`<`+name+`[\s/>]`).FindAllStringIndex(src, -1))
	}, repoRoot)
	inv.Components = append(inv.Components, comps...)
}

// ScanBusinessAndroid scans public composables in dirs (relative to aRoot) and counts calls under appDir.
func ScanBusinessAndroid(repoRoot, aRoot string, dirs []string, appDir string, inv *Inventory) {
	var comps []Component
	for _, d := range dirs {
		dir := filepath.Join(aRoot, d)
		for _, f := range walkFiles(dir, []string{".kt"}, nil) {
			src := readFile(f)
			pkg := regexp.MustCompile(`package\s+([\w.]+)`).FindStringSubmatch(src)
			for _, loc := range ktBizFunRe.FindAllStringSubmatchIndex(src, -1) {
				name := src[loc[2]:loc[3]]
				pre := src[max(0, loc[0]-20):loc[0]]
				if strings.Contains(pre, "private") || strings.HasSuffix(name, "Preview") {
					continue
				}
				imp := name
				if pkg != nil {
					imp = pkg[1] + "." + name
				}
				c := Component{Name: name, Kind: "business", File: rel(repoRoot, f), Import: "import " + imp}
				c.Description, _ = docBefore(src, loc[0])
				open := loc[1] - 1
				if end := matchClose(src, open); end > open {
					for _, p := range splitTop(src[open+1:end], ',') {
						parts := strings.SplitN(strings.TrimSpace(p), ":", 2)
						if len(parts) == 2 {
							c.Props = append(c.Props, Prop{Name: strings.TrimSpace(parts[0]), Type: strings.TrimSpace(strings.SplitN(parts[1], "=", 2)[0])})
						}
					}
				}
				comps = append(comps, c)
			}
		}
	}
	countUsage(filepath.Join(aRoot, appDir), []string{".kt"}, comps, func(src, name string) int {
		return len(regexp.MustCompile(`\b`+name+`\s*\(`).FindAllStringIndex(src, -1)) - len(regexp.MustCompile(`fun\s+`+name+`\s*\(`).FindAllStringIndex(src, -1))
	}, repoRoot)
	inv.Components = append(inv.Components, comps...)
}

func countUsage(root string, exts []string, comps []Component, count func(src, name string) int, repoRoot string) {
	for _, f := range walkFiles(root, exts, nil) {
		b, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		src := string(b)
		for i := range comps {
			if n := count(src, comps[i].Name); n > 0 {
				comps[i].Usage += n
				comps[i].UsageFiles = append(comps[i].UsageFiles, rel(repoRoot, f))
			}
		}
	}
}

// parseInlineTSProps parses a destructured parameter with an inline type literal: ({ a, b }: { a: T; b?: U }).
func parseInlineTSProps(param string) []Prop {
	i := strings.Index(param, "}:")
	if i < 0 {
		return nil
	}
	rest := strings.TrimSpace(param[i+2:])
	if !strings.HasPrefix(rest, "{") {
		return nil
	}
	end := matchClose(rest, 0)
	if end < 0 {
		return nil
	}
	var props []Prop
	for _, part := range splitTop(rest[1:end], ';') {
		m := h5PropRe.FindStringSubmatch(strings.TrimSpace(part))
		if m == nil {
			continue
		}
		p := Prop{Name: m[1], Optional: m[2] == "?", Type: strings.TrimSpace(m[3])}
		if litUnionRe.MatchString(p.Type) {
			for _, v := range litRe.FindAllStringSubmatch(p.Type, -1) {
				p.Values = append(p.Values, v[1])
			}
		}
		props = append(props, p)
	}
	return props
}
