package reuse

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	xmlTagRe       = regexp.MustCompile(`<([A-Za-z][\w.]*)[\s/>]`)
	includeRe      = regexp.MustCompile(`<include[^>]*layout="@layout/(\w+)"`)
	viewClassKtRe  = regexp.MustCompile(`(?m)^\s*(?:open\s+|abstract\s+|internal\s+)?class\s+([A-Z]\w*)\b[^{:]*?(?:\([^)]*\))?\s*:\s*([A-Z][\w.]*)\s*\(`)
	viewClassJvRe  = regexp.MustCompile(`(?m)^\s*public\s+(?:abstract\s+)?class\s+([A-Z]\w*)\s+extends\s+([A-Z][\w.]*)`)
	composableRe   = regexp.MustCompile(`@Composable\s+(?:(?:public|internal)\s+)?fun\s+([A-Z]\w*)\s*\(`)
	styleableRe2   = regexp.MustCompile(`(?s)<declare-styleable\s+name="(\w+)"\s*>(.*?)</declare-styleable>`)
	attrNameRe     = regexp.MustCompile(`<attr\s+name="(\w+)"`)
	colorResRe     = regexp.MustCompile(`<color\s+name="(\w+)">\s*(#[0-9A-Fa-f]{6,8})\s*</color>`)
	xmlHardColorRe = regexp.MustCompile(`android:(?:textColor|background|tint|src)="(#[0-9A-Fa-f]{6,8})"`)
	xmlHardSizeRe  = regexp.MustCompile(`android:(?:textSize|layout_width|layout_height|padding\w*|layout_margin\w*)="(\d+(?:\.\d+)?)(sp|dp|dip|px)"`)
	ktColorRe2     = regexp.MustCompile(`Color\(0x([0-9A-Fa-f]{8})\)|Color\.parseColor\("(#[0-9A-Fa-f]{6,8})"\)`)
)

// view base classes that make a class a reusable View component
var viewBases = regexp.MustCompile(`(View|Layout|Button|TextView|ImageView|EditText|Toolbar|Bar|Group|Dialog|Popup|Fragment)$`)

// layout containers are structure, not reusable components
var containers = map[string]bool{"LinearLayout": true, "FrameLayout": true, "RelativeLayout": true, "ConstraintLayout": true, "androidx.constraintlayout.widget.ConstraintLayout": true, "merge": true, "include": true, "layout": true, "data": true, "variable": true, "ScrollView": true, "androidx.core.widget.NestedScrollView": true, "androidx.coordinatorlayout.widget.CoordinatorLayout": true, "Space": true, "View": true, "ViewStub": true, "requestFocus": true, "tag": true, "import": true, "fragment": true, "androidx.fragment.app.FragmentContainerView": true}

// Detect returns the platforms present under root.
func Detect(root string) []string {
	var h5, an bool
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && skipDir(info.Name()) {
			return filepath.SkipDir
		}
		switch info.Name() {
		case "package.json":
			h5 = true
		case "build.gradle", "build.gradle.kts", "AndroidManifest.xml":
			an = true
		}
		return nil
	})
	var out []string
	if h5 {
		out = append(out, "h5")
	}
	if an {
		out = append(out, "android")
	}
	return out
}

// ScanAndroid inventories an Android project: XML layouts, custom Views (Kotlin / Java), composables.
func ScanAndroid(root string) *Report {
	r := &Report{Platform: "android", Root: root}
	type file struct{ path, rel, text string }
	var layouts, code, values []file
	gradle := ""
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if skipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		rel := relPath(root, p)
		name := info.Name()
		if strings.HasSuffix(name, ".gradle") || strings.HasSuffix(name, ".gradle.kts") || name == "libs.versions.toml" {
			b, _ := os.ReadFile(p)
			gradle += string(b)
			return nil
		}
		if info.Size() > 400_000 {
			return nil
		}
		switch {
		case strings.HasSuffix(name, ".xml") && strings.Contains(rel, "/res/layout"):
			b, _ := os.ReadFile(p)
			layouts = append(layouts, file{p, rel, string(b)})
		case strings.HasSuffix(name, ".xml") && strings.Contains(rel, "/res/values"):
			b, _ := os.ReadFile(p)
			values = append(values, file{p, rel, string(b)})
		case strings.HasSuffix(name, ".kt") || strings.HasSuffix(name, ".java"):
			b, _ := os.ReadFile(p)
			code = append(code, file{p, rel, string(b)})
		}
		return nil
	})
	r.Metrics.Files = len(layouts) + len(code)
	allCode := gradle
	for _, f := range code {
		if len(allCode) < 3_000_000 {
			allCode += f.text[:min(len(f.text), 3000)]
		}
	}
	if len(layouts) > 0 {
		r.Stack = append(r.Stack, "XML View")
	}
	if strings.Contains(allCode, "@Composable") || strings.Contains(gradle, "compose") {
		r.Stack = append(r.Stack, "Jetpack Compose")
	}
	var libs []uiLib
	for _, l := range uiLibs {
		if l.Platform != "android" {
			continue
		}
		on := l.Name == "Android SDK widgets" && len(layouts) > 0
		for _, p := range l.Packages {
			on = on || p != "android" && strings.Contains(allCode, p)
		}
		if on {
			libs = append(libs, l)
			if l.Name != "Android SDK widgets" {
				r.Stack = append(r.Stack, l.Name)
			}
		}
	}
	assets := map[string]*Asset{}
	get := func(kind, name, source string) *Asset {
		k := kind + "|" + name
		if a, ok := assets[k]; ok {
			return a
		}
		a := &Asset{Name: name, Kind: kind, Source: source, Concept: ConceptOf(name)}
		if kind == "ui-lib" {
			if c, ok := libConcept[source+"/"+name]; ok {
				a.Concept = c
			}
		}
		assets[k] = a
		return a
	}
	libOf := map[string]*Asset{}
	for _, l := range libs {
		for _, c := range l.Components {
			a := get("ui-lib", c, l.Name)
			if _, ok := libOf[c]; !ok || l.Name != "Android SDK widgets" {
				libOf[c] = a
			}
		}
	}
	// custom Views + attrs
	styleable := map[string][]string{}
	for _, f := range values {
		for _, m := range styleableRe2.FindAllStringSubmatch(f.text, -1) {
			for _, a := range attrNameRe.FindAllStringSubmatch(m[2], -1) {
				styleable[m[1]] = append(styleable[m[1]], a[1])
			}
		}
	}
	custom := map[string]*Asset{}
	compos := map[string]*Asset{}
	for _, f := range code {
		for _, re := range []*regexp.Regexp{viewClassKtRe, viewClassJvRe} {
			for _, m := range re.FindAllStringSubmatch(f.text, -1) {
				base := m[2][strings.LastIndex(m[2], ".")+1:]
				if !viewBases.MatchString(base) || strings.HasSuffix(base, "Fragment") || strings.HasSuffix(m[1], "Activity") || strings.HasSuffix(m[1], "Adapter") || strings.HasSuffix(m[1], "Holder") {
					continue
				}
				a := get("custom-view", m[1], f.rel)
				a.Props = styleable[m[1]]
				custom[m[1]] = a
			}
		}
		for _, m := range composableRe.FindAllStringSubmatchIndex(f.text, -1) {
			name := f.text[m[2]:m[3]]
			pre := f.text[max(0, m[0]-12):m[0]]
			if strings.Contains(pre, "private") || strings.HasSuffix(name, "Preview") || strings.HasSuffix(name, "Screen") || strings.HasSuffix(name, "App") || strings.HasSuffix(name, "Theme") {
				continue
			}
			compos[name] = get("composable", name, f.rel)
		}
	}
	// colour tokens
	colorTokens := map[string]string{}
	for _, f := range values {
		for _, m := range colorResRe.FindAllStringSubmatch(f.text, -1) {
			v := strings.ToUpper(m[2])
			if len(v) == 9 && strings.HasPrefix(v, "#FF") {
				v = "#" + v[3:]
			}
			if _, ok := colorTokens[v]; !ok {
				colorTokens[v] = "@color/" + m[1]
			}
		}
	}
	colorHits := map[string]*Hotspot{}
	hit := func(v, rel string) {
		v = strings.ToUpper(v)
		if len(v) == 9 && strings.HasPrefix(v, "#FF") {
			v = "#" + v[3:]
		}
		h := colorHits[v]
		if h == nil {
			h = &Hotspot{Value: v, Token: colorTokens[v]}
			colorHits[v] = h
		}
		h.Count++
		if len(h.Sample) < 3 && !contains(h.Sample, rel) {
			h.Sample = append(h.Sample, rel)
		}
	}
	use := func(a *Asset, rel string, ln int, line string) {
		a.Usage++
		if !contains(a.UsageFiles, rel) {
			a.UsageFiles = append(a.UsageFiles, rel)
		}
		if a.Example == "" {
			if len(line) > 140 {
				line = line[:140] + "…"
			}
			a.Example, a.ExampleAt = line, rel+":"+itoa(ln)
		}
		r.Metrics.ComponentUsages++
	}
	layoutMarkup := map[string]string{}
	for _, f := range layouts {
		layoutMarkup[f.rel] = f.text
		lines := strings.Split(f.text, "\n")
		for _, m := range xmlTagRe.FindAllStringSubmatchIndex(f.text, -1) {
			tag := f.text[m[2]:m[3]]
			ln := strings.Count(f.text[:m[0]], "\n") + 1
			line := strings.TrimSpace(lines[ln-1])
			simple := tag[strings.LastIndex(tag, ".")+1:]
			if containers[tag] || containers[simple] {
				continue
			}
			if a, ok := custom[simple]; ok && (strings.Contains(tag, ".") || true) {
				use(a, f.rel, ln, line)
				continue
			}
			if a, ok := libOf[simple]; ok {
				if a.Source == "Android SDK widgets" {
					if concept, raw := rawAndroid[simple]; raw {
						if alt := bestProjectAsset(assets, concept); alt != nil {
							r.Metrics.RawPrimitives++
							r.Findings = append(r.Findings, Finding{"missed-reuse", "warn", f.rel, ln, line, "原生 " + simple + "：项目里已有 " + alt.Name + "（" + kindLabel(alt.Kind) + "）"})
							continue
						}
					}
				}
				use(a, f.rel, ln, line)
			}
		}
		for _, m := range includeRe.FindAllStringSubmatch(f.text, -1) {
			a := get("business", "@layout/"+m[1], "include")
			a.Concept = ConceptOf(m[1])
			a.Usage++
			if !contains(a.UsageFiles, f.rel) {
				a.UsageFiles = append(a.UsageFiles, f.rel)
			}
		}
		for _, m := range xmlHardColorRe.FindAllStringSubmatch(f.text, -1) {
			hit(m[1], f.rel)
		}
		r.Metrics.HardcodedSizes += len(xmlHardSizeRe.FindAllStringIndex(f.text, -1))
	}
	// one alternation regex per kind: a single pass per file instead of one regex per asset
	alt := func(names map[string]*Asset) string {
		var ns []string
		for n := range names {
			ns = append(ns, regexp.QuoteMeta(n))
		}
		sort.Slice(ns, func(i, j int) bool { return len(ns[i]) > len(ns[j]) })
		return strings.Join(ns, "|")
	}
	var composRe, customRe *regexp.Regexp
	if len(compos) > 0 {
		composRe = regexp.MustCompile(`(?:^|[^\w.])(` + alt(compos) + `)\s*\(`)
	}
	if len(custom) > 0 {
		customRe = regexp.MustCompile(`\b(` + alt(custom) + `)\b[(.]`)
	}
	for _, f := range code {
		lines := strings.Split(f.text, "\n")
		if composRe != nil {
			for _, m := range composRe.FindAllStringSubmatchIndex(f.text, -1) {
				if strings.Contains(f.text[max(0, m[0]-5):m[0]+1], "fun ") {
					continue
				}
				name := f.text[m[2]:m[3]]
				ln := strings.Count(f.text[:m[0]], "\n") + 1
				use(compos[name], f.rel, ln, strings.TrimSpace(lines[ln-1]))
			}
		}
		// programmatic custom view / library usage in code (e.g. QMUITipDialog.Builder, MaterialButton(context))
		if customRe != nil {
			for _, m := range customRe.FindAllStringSubmatchIndex(f.text, -1) {
				name := f.text[m[2]:m[3]]
				if strings.HasSuffix(f.rel, "/"+name+".kt") || strings.HasSuffix(f.rel, "/"+name+".java") {
					continue
				}
				ln := strings.Count(f.text[:m[0]], "\n") + 1
				use(custom[name], f.rel, ln, strings.TrimSpace(lines[ln-1]))
			}
		}
		for _, m := range ktColorRe2.FindAllStringSubmatch(f.text, -1) {
			if m[1] != "" {
				hit("#"+m[1], f.rel)
			} else {
				hit(m[2], f.rel)
			}
		}
	}
	// Compose material3 / library calls
	for _, l := range libs {
		if l.Name != "Compose Material3" {
			continue
		}
		for _, f := range code {
			if !strings.Contains(f.text, "androidx.compose.material3") {
				continue
			}
			lines := strings.Split(f.text, "\n")
			for _, c := range l.Components {
				for _, m := range regexp.MustCompile(`(?:^|[^\w.])`+c+`\s*\(`).FindAllStringIndex(f.text, -1) {
					ln := strings.Count(f.text[:m[0]], "\n") + 1
					use(assets["ui-lib|"+c], f.rel, ln, strings.TrimSpace(lines[ln-1]))
				}
			}
		}
	}
	r.ColorTokens = colorTokens
	finish(r, assets, colorHits)
	r.Clones = findClones(layoutMarkup, func(s string) []*mnode { return dropContainers(parseMarkup(s)) }, 5)
	r.Metrics.CloneGroups = len(r.Clones)
	return r
}

func bestProjectAsset(assets map[string]*Asset, concept string) *Asset {
	var best *Asset
	for _, a := range assets {
		if a.Concept == concept && (a.Kind == "custom-view" || a.Kind == "composable" || a.Kind == "design-system" || a.Kind == "business" || a.Kind == "ui-lib" && a.Source != "Android SDK widgets") && (best == nil || rank(*a) > rank(*best) || rank(*a) == rank(*best) && len(a.Name) < len(best.Name)) {
			best = a
		}
	}
	return best
}

// dropContainers keeps containers in the tree (they are structure) but ignores xml prolog-like nodes.
func dropContainers(ns []*mnode) []*mnode {
	var out []*mnode
	for _, n := range ns {
		if strings.HasPrefix(n.tag, "?") || n.tag == "data" || n.tag == "variable" || n.tag == "import" {
			continue
		}
		n.children = dropContainers(n.children)
		out = append(out, n)
	}
	return out
}
