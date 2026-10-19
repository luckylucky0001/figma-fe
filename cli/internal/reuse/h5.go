package reuse

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	vueTemplateRe = regexp.MustCompile(`(?s)<template[^>]*>(.*)</template>`)
	vueStyleRe    = regexp.MustCompile(`(?s)<style[^>]*>(.*?)</style>`)
	jsxTagRe2     = regexp.MustCompile(`<([A-Za-z][\w.-]*)[\s/>]`)
	importRe      = regexp.MustCompile(`import\s+(?:(\w+)\s*,?\s*)?(?:\{([^}]*)\})?\s*from\s*['"]([^'"]+)['"]`)
	reactCompRe   = regexp.MustCompile(`(?m)^(?:export\s+)?(?:default\s+)?(?:function\s+([A-Z]\w*)\s*\(|const\s+([A-Z]\w*)\s*(?::\s*[\w.<>]+)?\s*=\s*(?:React\.)?(?:memo|forwardRef)?\(?\s*(?:\([^)]*\)|\w+)\s*=>)`)
	hexColorRe    = regexp.MustCompile(`#(?:[0-9a-fA-F]{6}|[0-9a-fA-F]{3})\b`)
	cssVarDefRe   = regexp.MustCompile(`(--[\w-]+)\s*:\s*(#[0-9a-fA-F]{3,6})\b`)
	preVarDefRe   = regexp.MustCompile(`([@$][\w-]+)\s*:\s*(#[0-9a-fA-F]{3,6})\b`)
	definePropsRe = regexp.MustCompile(`(?s)defineProps(?:<\s*\{(.*?)\}\s*>)?\(\s*(\{.*?\}|\[.*?\])?\s*\)`)
	propsObjRe    = regexp.MustCompile(`(?s)props\s*:\s*(\{.*?\}|\[.*?\])\s*,?\s*\n`)
	fixedRe       = regexp.MustCompile(`position:\s*fixed`)
	fixedBottomRe = regexp.MustCompile(`(?s)position:\s*fixed[^}]*?bottom:\s*0|bottom:\s*0[^}]*?position:\s*fixed`)
	propNameRe    = regexp.MustCompile(`(?:^|[{,;\n])\s*['"]?(\w+)['"]?\??\s*:`)
)

var htmlTags = map[string]bool{}

func init() {
	for _, t := range strings.Fields("div span p a ul ol li img button input select option textarea form label section header footer nav main article aside h1 h2 h3 h4 h5 h6 table tr td th tbody thead i b strong em small br hr svg path view text") {
		htmlTags[t] = true
	}
}

func kebabToPascal(s string) string {
	var b strings.Builder
	for _, p := range strings.Split(s, "-") {
		if p != "" {
			b.WriteString(strings.ToUpper(p[:1]) + p[1:])
		}
	}
	return b.String()
}

type srcFile struct {
	path, rel, kind string // kind: vue | jsx | style | js
	text, markup    string
}

// ScanH5 inventories a Vue / React mobile project at root.
func ScanH5(root string) *Report {
	r := &Report{Platform: "h5", Root: root}
	ign := ignoreList(root)
	deps := map[string]bool{}
	depVer := map[string]string{}
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && (skipDir(info.Name()) || ignored(ign, relPath(root, p))) {
			return filepath.SkipDir
		}
		if info.Name() == "package.json" {
			var pj struct{ Dependencies, DevDependencies map[string]string }
			if b, e := os.ReadFile(p); e == nil && json.Unmarshal(b, &pj) == nil {
				for k, v := range pj.Dependencies {
					deps[k] = true
					depVer[k] = cleanVersion(v)
				}
				for k := range pj.DevDependencies {
					deps[k] = true
				}
			}
		}
		return nil
	})
	for _, s := range []string{"vue", "react", "react-dom", "@tarojs/taro", "uni-app"} {
		if deps[s] {
			r.Stack = append(r.Stack, s)
		}
	}
	var files []srcFile
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if skipDir(info.Name()) || ignored(ign, relPath(root, p)) {
				return filepath.SkipDir
			}
			return nil
		}
		ext := filepath.Ext(p)
		kind := map[string]string{".vue": "vue", ".tsx": "jsx", ".jsx": "jsx", ".css": "style", ".less": "style", ".scss": "style", ".sass": "style", ".ts": "js", ".js": "js"}[ext]
		if kind == "" || strings.HasSuffix(p, ".d.ts") || strings.Contains(p, ".min.") || info.Size() > 400_000 {
			return nil
		}
		b, _ := os.ReadFile(p)
		f := srcFile{path: p, rel: relPath(root, p), kind: kind, text: string(b)}
		switch kind {
		case "vue":
			f.text = stripComments(f.text)
			if m := vueTemplateRe.FindStringSubmatch(f.text); m != nil {
				f.markup = m[1]
			}
		case "jsx":
			f.markup = f.text
		}
		files = append(files, f)
		return nil
	})
	r.Metrics.Files = len(files)

	// UI libraries in use (declared dependency, or imported / tag-prefixed in code).
	type libHit struct {
		lib   uiLib
		names map[string]bool
	}
	var libs []*libHit
	for _, l := range uiLibs {
		if l.Platform != "h5" {
			continue
		}
		on := false
		for _, p := range l.Packages {
			on = on || deps[p]
		}
		if !on {
			for _, f := range files {
				if l.TagPrefix != "" && strings.Contains(f.markup, "<"+l.TagPrefix) || strings.Contains(f.text, "from '"+l.Packages[0]+"'") || strings.Contains(f.text, `from "`+l.Packages[0]+`"`) {
					on = true
					break
				}
			}
		}
		if on {
			h := &libHit{lib: l, names: map[string]bool{}}
			for _, c := range l.Components {
				h.names[c] = true
			}
			libs = append(libs, h)
			r.Stack = append(r.Stack, l.Name)
		}
	}

	assets := map[string]*Asset{} // key kind|name
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
	versions := map[string]string{}
	for _, h := range libs {
		for c := range h.names {
			get("ui-lib", c, h.lib.Name)
		}
		for _, p := range h.lib.Packages {
			if v := depVer[p]; v != "" {
				versions[h.lib.Name] = v
			}
		}
	}
	markCompat(assets, versions, "", "")
	// Project components.
	projByName := map[string]*Asset{}
	fileAsset := map[string]*Asset{} // defining file → asset, to resolve aliased imports
	for _, f := range files {
		kind := classify(f.rel)
		switch f.kind {
		case "vue":
			name := strings.TrimSuffix(filepath.Base(f.path), ".vue")
			if strings.EqualFold(name, "index") {
				name = filepath.Base(filepath.Dir(f.path))
			}
			name = kebabToPascal(name)
			if kind == "page" || name == "App" {
				continue
			}
			a := get(kind, name, f.rel)
			a.Props = vueProps(f.text)
			a.Enums = unionEnums(f.text)
			if fixedBottomRe.MatchString(f.text) && fixedRe.MatchString(f.text) && (a.Concept == "navbar" || a.Concept == "") {
				a.Concept = "tabs" // a bar pinned to the bottom is a tab bar / action bar, whatever its name says
				if regexp.MustCompile(`(?i)button|btn|submit`).MatchString(f.markup) {
					a.Concept = "actionbar"
				}
			}
			projByName[name] = a
			fileAsset[f.rel] = a
		case "jsx":
			for _, m := range reactCompRe.FindAllStringSubmatch(f.text, -1) {
				name := m[1] + m[2]
				k := kind
				if k == "page" {
					k = "page-local"
				}
				a := get(k, name, f.rel)
				a.Props = tsProps(f.text)
				if a.Enums = unionEnums(f.text); len(a.Enums) > 0 {
					for p := range a.Enums {
						if !contains(a.Props, p) {
							a.Props = append(a.Props, p)
						}
					}
				}
				projByName[name] = a
			}
		}
	}

	// Usage: tags in templates / JSX.
	colorTokens := map[string]string{}
	for _, f := range files {
		src := f.text
		for _, m := range cssVarDefRe.FindAllStringSubmatch(src, -1) {
			colorTokens[strings.ToUpper(expandHex(m[2]))] = "var(" + m[1] + ")"
		}
		for _, m := range preVarDefRe.FindAllStringSubmatch(src, -1) {
			if _, ok := colorTokens[strings.ToUpper(expandHex(m[2]))]; !ok {
				colorTokens[strings.ToUpper(expandHex(m[2]))] = m[1]
			}
		}
	}
	rawCount := 0
	colorHits := map[string]*Hotspot{}
	for _, f := range files {
		imported := map[string]string{} // local name → source
		for _, m := range importRe.FindAllStringSubmatch(f.text, -1) {
			if m[1] != "" {
				imported[m[1]] = m[3]
			}
			for _, n := range strings.Split(m[2], ",") {
				n = strings.TrimSpace(n)
				if i := strings.Index(n, " as "); i >= 0 {
					n = strings.TrimSpace(n[i+4:])
				}
				if n != "" {
					imported[n] = m[3]
				}
			}
		}
		// aliased imports: import sHeader from '@/components/SimpleHeader.vue' → <s-header>
		alias := map[string]*Asset{}
		for _, m := range importRe.FindAllStringSubmatch(f.text, -1) {
			if m[1] == "" {
				continue
			}
			base := strings.TrimSuffix(filepath.Base(m[3]), ".vue")
			for rel, a := range fileAsset {
				if strings.TrimSuffix(filepath.Base(rel), ".vue") == base || strings.EqualFold(filepath.Base(filepath.Dir(rel)), base) && strings.HasSuffix(rel, "index.vue") {
					alias[strings.ToLower(m[1])] = a
				}
			}
		}
		lines := strings.Split(f.text, "\n")
		// at returns the line number and the tag snippet starting at idx (up to its closing '>' on that line).
		at := func(idx int) (int, string) {
			ln := strings.Count(f.text[:idx], "\n")
			start := strings.LastIndex(f.text[:idx], "\n") + 1
			rest := f.text[idx:]
			if e := strings.IndexByte(rest, '\n'); e >= 0 {
				rest = rest[:e]
			}
			if e := strings.Index(rest, "</"); e > 0 {
				if g := strings.IndexByte(rest[e:], '>'); g >= 0 {
					rest = rest[:e+g+1]
				}
			} else if e := strings.IndexByte(rest, '>'); e >= 0 {
				rest = rest[:e+1]
			}
			_ = start
			t := strings.TrimSpace(rest)
			if len(t) > 140 {
				t = t[:140] + "…"
			}
			if t == "" {
				t = strings.TrimSpace(lines[ln])
			}
			return ln + 1, t
		}
		if f.markup != "" {
			off := strings.Index(f.text, f.markup)
			for _, m := range jsxTagRe2.FindAllStringSubmatchIndex(f.markup, -1) {
				tag := f.markup[m[2]:m[3]]
				var a *Asset
				for _, h := range libs {
					if h.lib.TagPrefix != "" && strings.HasPrefix(tag, h.lib.TagPrefix) {
						if n := kebabToPascal(strings.TrimPrefix(tag, h.lib.TagPrefix)); h.names[n] {
							a = assets["ui-lib|"+n]
						}
					}
					pascal := kebabToPascal(tag)
					if a == nil && h.names[pascal] {
						if src, ok := imported[pascal]; ok && contains(h.lib.Import, src) {
							a = assets["ui-lib|"+pascal]
						}
					}
					// compound members: <List.Item>, <Form.Item> → ListItem, FormItem
					if i := strings.IndexByte(tag, '.'); a == nil && i > 0 && h.names[strings.ReplaceAll(tag, ".", "")] {
						if src, ok := imported[tag[:i]]; ok && contains(h.lib.Import, src) {
							a = assets["ui-lib|"+strings.ReplaceAll(tag, ".", "")]
						}
					}
				}
				if a == nil {
					if pa, ok := alias[strings.ToLower(strings.ReplaceAll(tag, "-", ""))]; ok {
						a = pa
					} else if pa, ok := projByName[kebabToPascal(tag)]; ok && (f.kind == "vue" || tag == kebabToPascal(tag)) {
						a = pa
					}
				}
				if a != nil {
					ln, t := at(off + m[0])
					a.Usage++
					if !contains(a.UsageFiles, f.rel) {
						a.UsageFiles = append(a.UsageFiles, f.rel)
					}
					if a.Example == "" || len(t) > len(a.Example) && len(a.Example) < 30 {
						a.Example, a.ExampleAt = t, f.rel+":"+itoa(ln)
					}
					r.Metrics.ComponentUsages++
					continue
				}
				if concept, ok := rawH5[tag]; ok {
					rawCount++
					ln, t := at(off + m[0])
					if alt := bestAsset(assets, concept); alt != nil {
						r.Findings = append(r.Findings, Finding{"missed-reuse", softSeverity(concept), f.rel, ln, t, "原生 <" + tag + ">：项目里已有 " + alt.Name + "（" + kindLabel(alt.Kind) + "）"})
					}
				}
			}
		}
		// colours: style files, <style> blocks and inline styles in markup
		styleText := ""
		switch f.kind {
		case "style":
			styleText = f.text
		case "vue":
			for _, m := range vueStyleRe.FindAllStringSubmatch(f.text, -1) {
				styleText += m[1] + "\n"
			}
			styleText += f.markup
		case "jsx":
			styleText = f.text
		}
		for _, loc := range hexColorRe.FindAllStringIndex(styleText, -1) {
			before := styleText[max(0, loc[0]-40):loc[0]]
			if cssVarDefRe.MatchString(before+styleText[loc[0]:loc[1]]) || preVarDefRe.MatchString(before+styleText[loc[0]:loc[1]]) {
				continue // a token definition, not a usage
			}
			v := strings.ToUpper(expandHex(styleText[loc[0]:loc[1]]))
			h := colorHits[v]
			if h == nil {
				h = &Hotspot{Value: v, Token: colorTokens[v]}
				colorHits[v] = h
			}
			h.Count++
			if len(h.Sample) < 3 && !contains(h.Sample, f.rel) {
				h.Sample = append(h.Sample, f.rel)
			}
		}
	}
	r.Metrics.RawPrimitives = rawCount
	r.ColorTokens = colorTokens
	r.Conventions = conventionsH5(files)
	finish(r, assets, colorHits)

	markup := map[string]string{}
	for _, f := range files {
		if f.markup != "" {
			markup[f.rel] = f.markup
		}
	}
	r.Clones = findClones(markup, func(s string) []*mnode { return filterJSX(parseMarkup(s)) }, 8)
	r.Metrics.CloneGroups = len(r.Clones)
	return r
}

// filterJSX drops pseudo-tags produced by TypeScript generics (useState<string>), keeping html + components.
func filterJSX(ns []*mnode) []*mnode {
	var out []*mnode
	for _, n := range ns {
		n.children = filterJSX(n.children)
		if htmlTags[n.tag] || strings.Contains(n.tag, "-") || n.tag[0] >= 'A' && n.tag[0] <= 'Z' {
			out = append(out, n)
		} else {
			out = append(out, n.children...)
		}
	}
	return out
}

func vueProps(src string) []string {
	var body string
	if m := definePropsRe.FindStringSubmatch(src); m != nil {
		body = m[1] + m[2]
	} else if m := propsObjRe.FindStringSubmatch(src); m != nil {
		body = m[1]
	}
	var out []string
	if strings.HasPrefix(strings.TrimSpace(body), "[") {
		for _, x := range regexp.MustCompile(`['"](\w+)['"]`).FindAllStringSubmatch(body, -1) {
			out = append(out, x[1])
		}
		return out
	}
	for _, x := range propNameRe.FindAllStringSubmatch(body, -1) {
		if x[1] != "type" && x[1] != "default" && x[1] != "required" && !contains(out, x[1]) {
			out = append(out, x[1])
		}
	}
	return out
}

// classify decides the asset kind from the path of the defining file.
func classify(rel string) string {
	p := "/" + strings.ToLower(filepath.ToSlash(rel))
	switch {
	case regexp.MustCompile(`/(ds|design-?system|ui|base|basic|common-ui)/`).MatchString(p):
		return "design-system"
	case regexp.MustCompile(`/(components?|widgets?|common|shared)/`).MatchString(p):
		return "business"
	case regexp.MustCompile(`/(views?|pages?|screens?|routes?)/`).MatchString(p):
		return "page"
	}
	return "business"
}

var appShell = map[string]bool{"App": true, "Router": true, "Routes": true, "Root": true, "Main": true, "Index": true, "AppRouter": true, "RouterView": true}

func bestAsset(assets map[string]*Asset, concept string) *Asset {
	var best *Asset
	for _, a := range assets {
		if a.Concept == concept && a.Kind != "page-local" && usable(a) && (best == nil || rank(*a) > rank(*best) || rank(*a) == rank(*best) && len(a.Name) < len(best.Name)) {
			best = a
		}
	}
	return best
}

// finish converts maps to sorted slices and fills metrics shared by both platforms.
func finish(r *Report, assets map[string]*Asset, colors map[string]*Hotspot) {
	for k, a := range assets {
		if a.Kind != "ui-lib" && appShell[a.Name] {
			delete(assets, k) // app shell / routing, not a UI component anyone would reuse
		}
	}
	for _, a := range assets {
		if a.Kind == "ui-lib" && a.Usage == 0 {
			continue // keep the report focused: unused library components are listed via r.Libraries
		}
		r.Assets = append(r.Assets, *a)
	}
	libUse := map[string]*LibUse{}
	for _, a := range assets {
		if a.Kind != "ui-lib" {
			continue
		}
		lu := libUse[a.Source]
		if lu == nil {
			lu = &LibUse{Name: a.Source}
			libUse[a.Source] = lu
		}
		lu.Usage += a.Usage
	}
	for _, lu := range libUse {
		var used []Asset
		for _, a := range r.Assets {
			if a.Kind == "ui-lib" && a.Source == lu.Name {
				used = append(used, a)
			}
		}
		sort.Slice(used, func(i, j int) bool { return used[i].Usage > used[j].Usage })
		for i, a := range used {
			if i < 8 {
				lu.Top = append(lu.Top, a.Name+"×"+itoa(a.Usage))
			}
		}
		r.Libraries = append(r.Libraries, *lu)
	}
	sort.Slice(r.Libraries, func(i, j int) bool { return r.Libraries[i].Usage > r.Libraries[j].Usage })
	sort.Slice(r.Assets, func(i, j int) bool {
		if r.Assets[i].Kind != r.Assets[j].Kind {
			return rank(Asset{Kind: r.Assets[i].Kind}) > rank(Asset{Kind: r.Assets[j].Kind})
		}
		return r.Assets[i].Usage > r.Assets[j].Usage
	})
	var libAll []Asset
	for _, a := range assets {
		if a.Kind == "ui-lib" {
			libAll = append(libAll, *a)
		}
	}
	r.LibAll = libAll
	r.Clusters = clusters(r.Assets, libAll)
	r.Metrics.DuplicateConcepts = len(r.Clusters)
	for _, h := range colors {
		r.Colors = append(r.Colors, *h)
		r.Metrics.HardcodedColors += h.Count
		if h.Token != "" {
			r.Metrics.TokenizableColors += h.Count
		}
	}
	sort.Slice(r.Colors, func(i, j int) bool { return r.Colors[i].Count > r.Colors[j].Count })
	if t := r.Metrics.ComponentUsages + r.Metrics.RawPrimitives; t > 0 {
		r.Metrics.Componentization = float64(r.Metrics.ComponentUsages) / float64(t)
	}
}

func expandHex(h string) string {
	h = strings.TrimPrefix(h, "#")
	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	return "#" + h
}

// ignoreList reads <root>/.figmafeignore: one path prefix per line (relative to root), # comments.
func ignoreList(root string) []string {
	b, err := os.ReadFile(filepath.Join(root, ".figmafeignore"))
	if err != nil {
		return nil
	}
	var out []string
	for _, l := range strings.Split(string(b), "\n") {
		l = strings.TrimSpace(l)
		if l != "" && !strings.HasPrefix(l, "#") {
			out = append(out, strings.TrimSuffix(filepath.ToSlash(l), "/"))
		}
	}
	return out
}

func ignored(ign []string, rel string) bool {
	for _, p := range ign {
		if rel == p || strings.HasPrefix(rel, p+"/") {
			return true
		}
	}
	return false
}

func skipDir(n string) bool {
	return n == "test" || n == "tests" || n == "androidTest" || n == "__tests__" || n == "e2e" || n == "testFixtures" || n == "screenshots" || n == "node_modules" || n == "dist" || n == "build" || n == ".git" || n == ".gradle" || n == ".idea" || n == "coverage" || n == ".figmafe" || n == "unpackage" || n == "miniprogram_npm"
}

func relPath(root, p string) string {
	r, err := filepath.Rel(root, p)
	if err != nil {
		return p
	}
	return filepath.ToSlash(r)
}

func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s || strings.HasPrefix(s, x) && strings.HasSuffix(x, ".") {
			return true
		}
	}
	return false
}
