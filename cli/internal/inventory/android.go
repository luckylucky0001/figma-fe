package inventory

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ktFunRe     = regexp.MustCompile(`@Composable\s+(?:public\s+)?fun\s+([A-Z]\w*)\s*\(|(?:public\s+)?@Composable\s+fun\s+([A-Z]\w*)\s*\(`)
	ktEnumRe    = regexp.MustCompile(`enum\s+class\s+(\w+)`)
	ktObjRe     = regexp.MustCompile(`object\s+(Ds\w+)\s*\{`)
	ktColorRe   = regexp.MustCompile(`val\s+(\w+)\s*=\s*Color\(0x([0-9A-Fa-f]{8})\)`)
	ktDpRe      = regexp.MustCompile(`val\s+(\w+)\s*=\s*([\d.]+)\.dp`)
	ktM3Import  = regexp.MustCompile(`import\s+androidx\.compose\.material3\.(\w+)`)
	ktColorLit  = regexp.MustCompile(`Color\(0x[0-9A-Fa-f]{6,8}\)`)
	viewClassRe = regexp.MustCompile(`class\s+(\w+)\s*(?:@JvmOverloads\s*)?(?:constructor\s*)?\(`)
	styleableRe = regexp.MustCompile(`(?s)<declare-styleable\s+name="(\w+)"\s*>(.*?)</declare-styleable>`)
	attrRe      = regexp.MustCompile(`(?s)<attr\s+name="(\w+)"\s+format="([\w|]+)"\s*(/>|>(.*?)</attr>)`)
	enumValRe   = regexp.MustCompile(`<(?:enum|flag)\s+name="(\w+)"`)
	xmlColorRe  = regexp.MustCompile(`<color\s+name="ds_(\w+)">\s*(#[0-9A-Fa-f]{6,8})\s*</color>`)
	xmlDimenRe  = regexp.MustCompile(`<dimen\s+name="ds_(\w+)">\s*([\d.]+)dp\s*</dimen>`)
	xmlRawRe    = regexp.MustCompile(`<(Button|TextView|ImageView|EditText)\b`)
	xmlHardRe   = regexp.MustCompile(`android:(textColor|background|textSize)="(#[0-9A-Fa-f]{6,8}|[\d.]+sp)"`)
)

// material3 components that must not be used directly in screens (text/icon primitives are allowed).
var M3Components = map[string]bool{
	"Button": true, "OutlinedButton": true, "TextButton": true, "FilledTonalButton": true, "ElevatedButton": true,
	"Card": true, "OutlinedCard": true, "ElevatedCard": true, "Divider": true, "HorizontalDivider": true,
	"Surface": true, "TopAppBar": true, "CenterAlignedTopAppBar": true, "BottomAppBar": true,
	"AssistChip": true, "SuggestionChip": true, "FilterChip": true, "Badge": true, "ListItem": true,
	"NavigationBar": true, "Scaffold": true,
}

// ScanAndroid scans Compose components under dsRel, XML View components under dsViewRel, and usage under appRel.
func ScanAndroid(repoRoot, aRoot, dsRel, dsViewRel, appRel string) (*Inventory, error) {
	inv := &Inventory{Schema: "figmafe-inventory/1", Platform: "android", Root: rel(repoRoot, aRoot)}
	dsDir := filepath.Join(aRoot, dsRel)
	enums := map[string][]string{}
	type fn struct {
		c   Component
		raw string
	}
	var fns []fn
	for _, f := range walkFiles(dsDir, []string{".kt"}, nil) {
		src := readFile(f)
		for _, loc := range ktEnumRe.FindAllStringSubmatchIndex(src, -1) {
			name := src[loc[2]:loc[3]]
			i := loc[1]
			for i < len(src) && src[i] != '{' && src[i] != '(' {
				i++
			}
			if i < len(src) && src[i] == '(' {
				i = matchClose(src, i) + 1
				for i < len(src) && src[i] != '{' {
					i++
				}
			}
			if i >= len(src) {
				continue
			}
			end := matchClose(src, i)
			body := src[i+1 : end]
			if k := strings.IndexByte(body, ';'); k >= 0 {
				body = body[:k]
			}
			for _, e := range splitTop(body, ',') {
				e = strings.TrimSpace(e)
				if m := regexp.MustCompile(`^(\w+)`).FindString(e); m != "" {
					enums[name] = append(enums[name], m)
				}
			}
		}
		for _, loc := range ktFunRe.FindAllStringSubmatchIndex(src, -1) {
			name := ""
			if loc[2] >= 0 {
				name = src[loc[2]:loc[3]]
			} else {
				name = src[loc[4]:loc[5]]
			}
			pre := src[max(0, loc[0]-20):loc[0]]
			if strings.HasSuffix(name, "Preview") || strings.Contains(pre, "private") || strings.Contains(pre, "internal") || !strings.HasPrefix(name, "Ds") {
				continue
			}
			open := loc[1] - 1
			close := matchClose(src, open)
			c := Component{Name: name, Kind: "compose", File: rel(repoRoot, f), Import: "import com.demo.ds." + name}
			c.Description, c.Example = docBefore(src, loc[0])
			fns = append(fns, fn{c, src[open+1 : close]})
		}
	}
	for _, f := range fns {
		c := f.c
		for _, p := range splitTop(f.raw, ',') {
			p = strings.TrimSpace(regexp.MustCompile(`@\w+\s*`).ReplaceAllString(p, ""))
			if p == "" {
				continue
			}
			pr := Prop{}
			if k := strings.Index(p, "="); k >= 0 && !strings.Contains(p[:k], "->") {
				pr.Default = strings.TrimSpace(p[k+1:])
				pr.Optional = true
				p = p[:k]
			}
			parts := strings.SplitN(p, ":", 2)
			if len(parts) != 2 {
				continue
			}
			pr.Name, pr.Type = strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			if v, ok := enums[pr.Type]; ok {
				pr.Values = v
			}
			if pr.Type == "Boolean" {
				pr.Values = []string{"true", "false"}
			}
			c.Props = append(c.Props, pr)
		}
		inv.Components = append(inv.Components, c)
	}

	// XML View components (legacy).
	if dsViewRel != "" {
		vroot := filepath.Join(aRoot, dsViewRel)
		attrs := map[string][]Prop{}
		for _, f := range walkFiles(filepath.Join(vroot, "res"), []string{".xml"}, nil) {
			src := readFile(f)
			for _, m := range styleableRe.FindAllStringSubmatch(src, -1) {
				for _, a := range attrRe.FindAllStringSubmatch(m[2], -1) {
					p := Prop{Name: a[1], Type: a[2], Optional: true}
					for _, ev := range enumValRe.FindAllStringSubmatch(a[4], -1) {
						p.Values = append(p.Values, ev[1])
					}
					if a[2] == "boolean" {
						p.Values = []string{"true", "false"}
					}
					attrs[m[1]] = append(attrs[m[1]], p)
				}
			}
			for _, m := range xmlColorRe.FindAllStringSubmatch(src, -1) {
				inv.Tokens = append(inv.Tokens, Token{Name: strings.Replace(strings.ReplaceAll(m[1], "_", "-"), "color-", "color.", 1), Value: m[2], Source: rel(repoRoot, f)})
			}
			for _, m := range xmlDimenRe.FindAllStringSubmatch(src, -1) {
				n := strings.ReplaceAll(m[1], "_", "-")
				if i := strings.IndexByte(n, '-'); i > 0 {
					n = n[:i] + "." + n[i+1:]
				}
				inv.Tokens = append(inv.Tokens, Token{Name: n + "@xml", Value: m[2] + "dp", Source: rel(repoRoot, f)})
			}
		}
		for _, f := range walkFiles(filepath.Join(vroot, "java"), []string{".kt", ".java"}, nil) {
			src := readFile(f)
			for _, loc := range viewClassRe.FindAllStringSubmatchIndex(src, -1) {
				name := src[loc[2]:loc[3]]
				pkg := regexp.MustCompile(`package\s+([\w.]+)`).FindStringSubmatch(src)
				fq := name
				if pkg != nil {
					fq = pkg[1] + "." + name
				}
				c := Component{Name: name, Kind: "view", File: rel(repoRoot, f), Import: "<" + fq + " />", Props: attrs[name]}
				c.Description, _ = docBefore(src, loc[0])
				inv.Components = append(inv.Components, c)
			}
		}
		for _, f := range walkFiles(filepath.Join(vroot, "res", "layout"), []string{".xml"}, nil) {
			src := readFile(f)
			for i := range inv.Components {
				c := &inv.Components[i]
				if c.Kind != "view" {
					continue
				}
				if n := len(regexp.MustCompile(`<[\w.]*\b`+c.Name+`\b`).FindAllStringIndex(src, -1)); n > 0 {
					c.Usage += n
					c.UsageFiles = append(c.UsageFiles, rel(repoRoot, f))
				}
			}
			for _, loc := range xmlRawRe.FindAllStringIndex(src, -1) {
				inv.Baseline = append(inv.Baseline, Finding{Rule: "raw-view", File: rel(repoRoot, f), Line: lineOf(src, loc[0]), Text: src[loc[0]:loc[1]]})
			}
			for _, loc := range xmlHardRe.FindAllStringIndex(src, -1) {
				inv.Baseline = append(inv.Baseline, Finding{Rule: "hardcoded-value", File: rel(repoRoot, f), Line: lineOf(src, loc[0]), Text: src[loc[0]:loc[1]]})
			}
		}
	}

	// Compose usage + baseline in app code.
	appDir := filepath.Join(aRoot, appRel)
	for _, f := range walkFiles(appDir, []string{".kt"}, nil) {
		src := readFile(f)
		for i := range inv.Components {
			c := &inv.Components[i]
			if c.Kind != "compose" {
				continue
			}
			if n := len(regexp.MustCompile(`\b`+c.Name+`\s*\(`).FindAllStringIndex(src, -1)); n > 0 {
				c.Usage += n
				c.UsageFiles = append(c.UsageFiles, rel(repoRoot, f))
			}
		}
		for _, m := range ktM3Import.FindAllStringSubmatchIndex(src, -1) {
			name := src[m[2]:m[3]]
			if M3Components[name] {
				inv.Baseline = append(inv.Baseline, Finding{Rule: "direct-material3", File: rel(repoRoot, f), Line: lineOf(src, m[0]), Text: "material3." + name})
			}
		}
		for _, loc := range ktColorLit.FindAllStringIndex(src, -1) {
			inv.Baseline = append(inv.Baseline, Finding{Rule: "hardcoded-color", File: rel(repoRoot, f), Line: lineOf(src, loc[0]), Text: src[loc[0]:loc[1]]})
		}
	}

	// Compose tokens.
	for _, f := range walkFiles(dsDir, []string{".kt"}, nil) {
		src := readFile(f)
		for _, loc := range ktObjRe.FindAllStringSubmatchIndex(src, -1) {
			obj := src[loc[2]:loc[3]]
			end := matchClose(src, loc[1]-1)
			body := src[loc[1]:end]
			prefix := map[string]string{"DsColors": "color", "DsSpacing": "space", "DsRadius": "radius"}[obj]
			if prefix == "" {
				continue
			}
			for _, m := range ktColorRe.FindAllStringSubmatch(body, -1) {
				inv.Tokens = append(inv.Tokens, Token{Name: prefix + "." + kebab(m[1]), Value: "#" + strings.ToUpper(m[2][2:]), Source: rel(repoRoot, f)})
			}
			for _, m := range ktDpRe.FindAllStringSubmatch(body, -1) {
				inv.Tokens = append(inv.Tokens, Token{Name: prefix + "." + kebab(m[1]), Value: m[2] + "dp", Source: rel(repoRoot, f)})
			}
		}
	}
	return inv, nil
}
