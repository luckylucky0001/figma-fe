package inventory

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	h5ExportRe = regexp.MustCompile(`export\s*\{\s*(\w+)\s*\}\s*from\s*'\./([\w/]+)'`)
	h5PropRe   = regexp.MustCompile(`^\s*(\w+)(\?)?\s*:\s*(.+?);?\s*$`)
	litUnionRe = regexp.MustCompile(`^'[^']*'(\s*\|\s*'[^']*')*$`)
	litRe      = regexp.MustCompile(`'([^']*)'`)
	cssVarRe   = regexp.MustCompile(`--ds-([\w-]+)\s*:\s*([^;]+);`)
	h5RawRe    = regexp.MustCompile(`<(button|input|select|textarea|img)\b`)
)

// ScanH5 scans a React component library (dsDir) and its usage in the rest of srcRoot.
func ScanH5(repoRoot, h5Root, dsRel string) (*Inventory, error) {
	dsDir := filepath.Join(h5Root, dsRel)
	// Source library (index.ts + Name/index.tsx) or a published package's type declarations (index.d.ts).
	ext := ".tsx"
	idx := readFile(filepath.Join(dsDir, "index.ts"))
	if idx == "" {
		idx = readFile(filepath.Join(dsDir, "index.d.ts"))
		ext = ".d.ts"
	}
	if idx == "" {
		return nil, fmt.Errorf("%s/index.ts (or index.d.ts) not found", dsDir)
	}
	pkgName := packageName(dsDir)
	inv := &Inventory{Schema: "figmafe-inventory/1", Platform: "h5", Root: rel(repoRoot, h5Root)}
	for _, m := range h5ExportRe.FindAllStringSubmatch(idx, -1) {
		name, sub := m[1], m[2]
		file := filepath.Join(dsDir, sub, "index"+ext)
		src := readFile(file)
		if src == "" {
			file = filepath.Join(dsDir, sub+ext)
			src = readFile(file)
		}
		imp := "@ds"
		if pkgName != "" {
			imp = pkgName
		}
		c := Component{Name: name, Kind: "react", File: rel(repoRoot, file), Import: fmt.Sprintf("import { %s } from '%s';", name, imp)}
		c.Props = parseTSProps(src, name+"Props")
		if loc := regexp.MustCompile(`export\s+(declare\s+)?(function|const)\s+` + name + `\b`).FindStringIndex(src); loc != nil {
			c.Description, c.Example = docBefore(src, loc[0])
		}
		inv.Components = append(inv.Components, c)
	}
	srcRoot := filepath.Join(h5Root, "src")
	files := walkFiles(srcRoot, []string{".tsx", ".ts", ".css"}, func(p string) bool { return p == dsDir })
	for _, f := range files {
		src := readFile(f)
		if strings.HasSuffix(f, ".tsx") {
			for i := range inv.Components {
				c := &inv.Components[i]
				if !regexp.MustCompile(`import\s*\{[^}]*\b` + c.Name + `\b[^}]*\}\s*from`).MatchString(src) {
					continue
				}
				n := len(regexp.MustCompile(`<`+c.Name+`[\s/>]`).FindAllStringIndex(src, -1))
				if n > 0 {
					c.Usage += n
					c.UsageFiles = append(c.UsageFiles, rel(repoRoot, f))
				}
			}
		}
		for _, loc := range h5RawRe.FindAllStringIndex(src, -1) {
			inv.Baseline = append(inv.Baseline, Finding{Rule: "raw-element", File: rel(repoRoot, f), Line: lineOf(src, loc[0]), Text: src[loc[0]:loc[1]]})
		}
		for _, loc := range hexRe.FindAllStringIndex(src, -1) {
			inv.Baseline = append(inv.Baseline, Finding{Rule: "hardcoded-color", File: rel(repoRoot, f), Line: lineOf(src, loc[0]), Text: src[loc[0]:loc[1]]})
		}
	}
	tokFile := filepath.Join(dsDir, "tokens.css")
	if readFile(tokFile) == "" { // packaged library: tokens are bundled into the stylesheet next to the types
		if m, _ := filepath.Glob(filepath.Join(dsDir, "..", "*.css")); len(m) > 0 {
			tokFile = m[0]
		}
	}
	for _, m := range cssVarRe.FindAllStringSubmatch(readFile(tokFile), -1) {
		name := m[1]
		if i := strings.IndexByte(name, '-'); i > 0 {
			name = name[:i] + "." + name[i+1:]
		}
		inv.Tokens = append(inv.Tokens, Token{Name: name, Value: strings.TrimSpace(m[2]), Source: rel(repoRoot, tokFile)})
	}
	return inv, nil
}

func parseTSProps(src, iface string) []Prop {
	loc := regexp.MustCompile(`export\s+interface\s+` + iface + `\b[^{]*\{`).FindStringIndex(src)
	if loc == nil {
		return nil
	}
	end := matchClose(src, loc[1]-1)
	if end < 0 {
		return nil
	}
	body := src[loc[1]:end]
	var props []Prop
	doc := ""
	for _, line := range strings.Split(body, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "/**") {
			doc = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(t, "/**"), "*/"))
			continue
		}
		if strings.HasPrefix(t, "*") || strings.HasPrefix(t, "//") || t == "" {
			continue
		}
		m := h5PropRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		p := Prop{Name: m[1], Optional: m[2] == "?", Type: strings.TrimSpace(m[3]), Doc: doc}
		if litUnionRe.MatchString(p.Type) {
			for _, v := range litRe.FindAllStringSubmatch(p.Type, -1) {
				p.Values = append(p.Values, v[1])
			}
		}
		if p.Type == "boolean" {
			p.Values = []string{"true", "false"}
		}
		props = append(props, p)
		doc = ""
	}
	return props
}

// packageName returns the npm package name if dir sits inside node_modules/<pkg>.
func packageName(dir string) string {
	for d := dir; d != filepath.Dir(d); d = filepath.Dir(d) {
		if b := readFile(filepath.Join(d, "package.json")); b != "" && strings.Contains(filepath.Dir(d), "node_modules") {
			if m := regexp.MustCompile(`"name"\s*:\s*"([^"]+)"`).FindStringSubmatch(b); m != nil {
				return m[1]
			}
		}
	}
	return ""
}
