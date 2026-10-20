package reuse

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	xmlTagRe  = regexp.MustCompile(`<([A-Za-z][\w.]*)[\s/>]`)
	includeRe = regexp.MustCompile(`<include[^>]*layout="@layout/(\w+)"`)
	// modifiers in any order, generics (with bounds), and a base with or without "(" (secondary constructors)
	viewClassKtRe = regexp.MustCompile(`(?m)^\s*(?:(?:open|abstract|internal|public|private|sealed|data|inner)\s+)*class\s+([A-Z]\w*)\s*(?:<[^{]*?>)?[^{:]*?(?:\([^)]*\))?\s*:\s*([A-Z][\w.]*)\s*[({]`)
	// a class that inflates a layout for one UI concept (AntennaPod's EmptyViewHandler): the project's way
	// to show that concept, though it is not a View itself
	// a class declared with no supertype at all (Java: no extends; Kotlin: no ":" before the body)
	plainClassJvRe  = regexp.MustCompile(`(?m)^\s*(?:public\s+|abstract\s+|final\s+)*class\s+([A-Z]\w*)\s*(?:<[^>{]*>)?\s*(?:implements\s+[^{]*)?\{`)
	plainClassKtRe  = regexp.MustCompile(`(?m)^\s*(?:internal\s+|open\s+|abstract\s+|data\s+|sealed\s+)*class\s+([A-Z]\w*)\s*(?:<[^>{]*>)?\s*(?:\([^)]*\))?\s*\{`)
	helperClassRe   = regexp.MustCompile(`(?m)^\s*(?:public\s+|internal\s+|open\s+)?(?:final\s+)?class\s+([A-Z]\w*(?:Handler|Helper|Controller))\b`)
	viewClassJvRe   = regexp.MustCompile(`(?m)^\s*(?:(?:public|abstract|final)\s+)*class\s+([A-Z]\w*)\s*(?:<[^{]*?>)?\s+extends\s+([A-Z][\w.]*)`)
	bannerContentRe = regexp.MustCompile(`(?i)image|item|page|banner|url|list|data|count`)
	searchArgRe     = regexp.MustCompile(`(?i)\b(?:value|text|query)\s*=\s*[\w.]*(?:search|query|keyword)`)
	composCallRe    = regexp.MustCompile(`(?:^|[^\w.])([A-Z]\w*)\s*[({]`)
	composableRe    = regexp.MustCompile(`@Composable\s+(?:(?:public|internal)\s+)?fun\s+(?:<[^>]*>\s*)?(?:[\w.]+\.)?([A-Z]\w*)\s*\(`) // fun <T> X(, fun RowScope.X(
	styleableRe2    = regexp.MustCompile(`(?s)<declare-styleable\s+name="(\w+)"\s*>(.*?)</declare-styleable>`)
	attrNameRe      = regexp.MustCompile(`<attr\s+name="(\w+)"`)
	colorResRe      = regexp.MustCompile(`<color\s+name="(\w+)"[^>]*>\s*(#[0-9A-Fa-f]{3,8})\s*</color>`)
	// any colour attribute (textColorHint, backgroundTint, app:tint, cardBackgroundColor …), #rgb to #aarrggbb
	xmlHardColorRe = regexp.MustCompile(`\b(?:android|app):\w*(?:[Cc]olor|[Bb]ackground|[Tt]int|src)\w*="(#[0-9A-Fa-f]{3,8})"`)
	xmlHardSizeRe  = regexp.MustCompile(`android:(?:textSize|layout_width|layout_height|padding\w*|layout_margin\w*)="(\d+(?:\.\d+)?)(sp|dp|dip|px)"`)
	ktColorRe2     = regexp.MustCompile(`Color\(0x([0-9A-Fa-f]{8})\)|Color\.parseColor\("(#[0-9A-Fa-f]{6,8})"\)`)
)

// view base classes that make a class a reusable View component
var styleAttrRe = regexp.MustCompile(`\bstyle="@style/([\w.]+)"`)

var viewBases = regexp.MustCompile(`(View|Layout|Button|TextView|ImageView|EditText|Toolbar|Bar|Group|Dialog|Popup|Fragment)$`)

// layout containers are structure, not reusable components
var containers = map[string]bool{"LinearLayout": true, "FrameLayout": true, "RelativeLayout": true, "ConstraintLayout": true, "androidx.constraintlayout.widget.ConstraintLayout": true, "merge": true, "include": true, "layout": true, "data": true, "variable": true, "ScrollView": true, "androidx.core.widget.NestedScrollView": true, "androidx.coordinatorlayout.widget.CoordinatorLayout": true, "Space": true, "View": true, "ViewStub": true, "requestFocus": true, "tag": true, "import": true, "fragment": true, "androidx.fragment.app.FragmentContainerView": true,
	"DrawerLayout": true, "SwipeRefreshLayout": true, "SmartRefreshLayout": true, "MotionLayout": true, "Barrier": true, "Guideline": true, "Placeholder": true, "Group": true, "Flow": true, "Layer": true}

// h5FrameworkRe: a dependency that makes a package.json an H5 UI project.
var h5FrameworkRe = regexp.MustCompile(`"(vue|react|react-dom|preact|@dcloudio/uni-app|@tarojs/taro|nuxt|next|svelte|solid-js)"\s*:`)

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
			// a front-end project, not a tooling package.json (gkd: build scripts, a TS selector library)
			if b, err := os.ReadFile(p); err == nil && h5FrameworkRe.Match(b) {
				h5 = true
			}
		case "pages.json": // uni-app
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
	ign := ignoreList(root)
	type file struct{ path, rel, text string }
	var layouts, code, values []file
	var manifests []string
	gradle := ""
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
		case name == "AndroidManifest.xml":
			b, _ := os.ReadFile(p)
			manifests = append(manifests, string(b))
		case strings.HasSuffix(name, ".xml") && strings.Contains("/"+rel, "/res/layout"):
			b, _ := os.ReadFile(p)
			layouts = append(layouts, file{p, rel, string(b)})
		case strings.HasSuffix(name, ".xml") && strings.Contains("/"+rel, "/res/values"):
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
		if kind == "ui-lib" {
			k = "ui-lib|" + source + "|" + name // SDK Button and Compose Material3 Button are different things
		}
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
			if isComposeLib(l.Name) {
				continue // composables never appear as XML tags
			}
			if _, ok := libOf[c]; !ok || l.Name != "Android SDK widgets" {
				libOf[c] = a
			}
		}
	}
	// what the installed Material version / app theme allow
	verTexts := []string{gradle}
	for _, f := range code {
		if strings.Contains(f.text, "com.google.android.material:material:") {
			verTexts = append(verTexts, f.text)
		}
	}
	var valueTexts []string
	for _, f := range values {
		valueTexts = append(valueTexts, f.text)
	}
	family, theme := appThemeFamily(manifests, valueTexts)
	markCompat(assets, map[string]string{"Material Components": materialVersion(verTexts...)}, family, theme)
	// custom Views + attrs
	styleable := map[string][]string{}
	styleEnums := map[string]map[string][]string{}
	ktEnums := map[string][]string{}
	for _, f := range code {
		if strings.Contains(f.text, "enum class") {
			for k, v := range kotlinEnums(f.text) {
				ktEnums[k] = v
			}
		}
	}
	for _, f := range values {
		for _, m := range styleableRe2.FindAllStringSubmatch(f.text, -1) {
			styleEnums[m[1]] = styleableEnums(m[2])
			for _, a := range attrNameRe.FindAllStringSubmatch(m[2], -1) {
				styleable[m[1]] = append(styleable[m[1]], a[1])
			}
		}
	}
	custom := map[string]*Asset{}
	compos := map[string]*Asset{}
	type composDef struct {
		line int
		name string
	}
	composDefs := map[string][]composDef{} // file → composables in definition order (overloads included)
	composBody := map[string]string{}      // composable → its body (the first overload)
	for _, f := range code {
		for _, re := range []*regexp.Regexp{viewClassKtRe, viewClassJvRe} {
			for _, m := range re.FindAllStringSubmatch(f.text, -1) {
				base := m[2][strings.LastIndex(m[2], ".")+1:]
				if !viewBases.MatchString(base) || strings.HasSuffix(base, "Fragment") || strings.HasSuffix(m[1], "Activity") || strings.HasSuffix(m[1], "Adapter") || strings.HasSuffix(m[1], "Holder") {
					continue
				}
				a := get("custom-view", m[1], f.rel)
				a.Extends = base
				if i := strings.Index(f.text, "class "+m[1]); i >= 0 {
					a.Purpose = docPurpose(docBefore(f.text, i))
				}
				if lib := libOf[base]; lib != nil && lib.Source != "Android SDK widgets" && inDesignSystem(f.rel) && lib.WrappedBy == "" {
					lib.WrappedBy = m[1] // class DsButton : MaterialButton
				}
				a.Props = styleable[m[1]]
				a.Enums = styleEnums[m[1]]
				refineConcept(a)
				custom[m[1]] = a
			}
		}
		if strings.Contains(f.text, "R.layout.") {
			for _, m := range helperClassRe.FindAllStringSubmatch(f.text, -1) {
				if custom[m[1]] == nil && ConceptOf(m[1]) != "" {
					custom[m[1]] = get("view-helper", m[1], f.rel)
				}
			}
		}
		glance := strings.Contains(f.text, "import androidx.glance") // home-screen widget UI, a different toolkit
		for _, m := range composableRe.FindAllStringSubmatchIndex(f.text, -1) {
			if glance {
				break
			}
			name := f.text[m[2]:m[3]]
			pre := f.text[max(0, m[0]-12):m[0]]
			ann := annotationsBefore(f.text, m[0])
			if strings.Contains(pre, "private") || strings.Contains(name, "Preview") || strings.Contains(ann, "Preview") || strings.Contains(name, "Screen") ||
				strings.HasSuffix(name, "App") || strings.HasSuffix(name, "Theme") || strings.HasSuffix(name, "Content") || strings.HasSuffix(name, "Route") {
				continue // previews, screens and their sections are not reusable components
			}
			a := get("composable", name, f.rel)
			if end := matchParen(f.text, m[1]-1); end > 0 && composBody[name] == "" {
				composBody[name] = funcBody(f.text, end+1)
			}
			if a.Purpose == "" {
				a.Purpose = docPurpose(docBefore(f.text, m[0]))
			}
			composDefs[f.rel] = append(composDefs[f.rel], composDef{strings.Count(f.text[:m[0]], "\n") + 1, name})
			if end := matchParen(f.text, m[1]-1); end > m[1] {
				for _, prm := range splitTopLevel(f.text[m[1]:end]) {
					if kv := strings.SplitN(strings.TrimSpace(prm), ":", 2); len(kv) == 2 {
						pn := paramName(kv[0])
						a.Props = append(a.Props, pn)
						if !strings.Contains(kv[1], "=") {
							a.Required = append(a.Required, pn)
							if t := strings.TrimSuffix(strings.TrimSpace(kv[1]), "?"); domainType(t) {
								a.Required = append(a.Required, pn+":"+t)
							}
						}
						typ := strings.TrimSuffix(strings.TrimSpace(strings.SplitN(kv[1], "=", 2)[0]), "?")
						if vals := ktEnums[typ]; len(vals) > 0 {
							if a.Enums == nil {
								a.Enums = map[string][]string{}
							}
							a.Enums[pn] = vals
						}
					}
				}
			}
			refineConcept(a)
			compos[name] = a
		}
	}
	// project composables a composable calls in its body (GkSettingsDialog → GkDialog): built on it
	for name, body := range composBody {
		a := compos[name]
		for _, m := range composCallRe.FindAllStringSubmatch(body, -1) {
			if compos[m[1]] != nil && m[1] != name && !contains(a.BuiltOn, m[1]) {
				a.BuiltOn = append(a.BuiltOn, m[1])
			}
		}
	}
	// DownloadActionButton extends ItemActionButton, a plain project class with no supertype: a name
	// ending in Button does not make a View. Drop custom "views" whose base is such a class.
	declared := map[string]bool{} // only classes known to have no supertype: BaseItemView<T> extends FrameLayout stays a View
	for _, f := range code {
		for _, re := range []*regexp.Regexp{plainClassJvRe, plainClassKtRe} {
			for _, m := range re.FindAllStringSubmatch(f.text, -1) {
				declared[m[1]] = true
			}
		}
	}
	for changed := true; changed; {
		changed = false
		for name, a := range custom {
			if a.Kind == "custom-view" && declared[a.Extends] && (custom[a.Extends] == nil || custom[a.Extends].Kind != "custom-view") {
				delete(custom, name)
				delete(assets, "custom-view|"+name)
				changed = true
			}
		}
	}
	// the project's own packages (org.schabi.newpipe): views from them are never a third-party library
	ownPkgs := map[string]bool{}
	for _, f := range code {
		if m := pkgDeclRe.FindStringSubmatch(f.text); m != nil {
			ownPkgs[pkgRoot(m[1])] = true
		}
	}
	// colour tokens
	colorTokens := map[string]string{}
	for _, f := range values {
		for _, m := range colorResRe.FindAllStringSubmatch(f.text, -1) {
			v := androidHex(m[2])
			if _, ok := colorTokens[v]; !ok {
				colorTokens[v] = "@color/" + m[1]
			}
		}
	}
	colorHits := map[string]*Hotspot{}
	hit := func(v, rel string) {
		v = androidHex(v)
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
			line = clip(line, 140)
			a.Example, a.ExampleAt = line, rel+":"+itoa(ln)
		}
		r.Metrics.ComponentUsages++
	}
	type rawOcc struct {
		a                    *Asset
		simple, concept, rel string
		ln                   int
		line                 string
	}
	var raws []rawOcc
	layoutMarkup := map[string]string{}
	for i := range layouts {
		layouts[i].text = stripComments(layouts[i].text)
	}
	for _, f := range layouts {
		layoutMarkup[f.rel] = f.text
		lines := strings.Split(f.text, "\n")
		for _, m := range xmlTagRe.FindAllStringSubmatchIndex(f.text, -1) {
			tag := f.text[m[2]:m[3]]
			ln := strings.Count(f.text[:m[0]], "\n") + 1
			line := strings.TrimSpace(lines[ln-1])
			simple := tag[strings.LastIndex(tag, ".")+1:]
			// <androidx.appcompat.widget.AppCompatButton> is the platform Button (AndroidProject × 54)
			if base := strings.TrimPrefix(simple, "AppCompat"); base != simple && libOf[base] != nil && libOf[base].Source == "Android SDK widgets" {
				simple = base
			}
			if containers[tag] || containers[simple] {
				continue
			}
			style := ""
			if end := strings.IndexByte(f.text[m[0]:], '>'); end > 0 {
				if sm := styleAttrRe.FindStringSubmatch(f.text[m[0] : m[0]+end]); sm != nil {
					style = sm[1]
				}
			}
			noteStyle := func(a *Asset) {
				if style != "" {
					if a.Styles == nil {
						a.Styles = map[string]int{}
					}
					a.Styles[style]++
				}
			}
			if a, ok := custom[simple]; ok && (strings.Contains(tag, ".") || true) {
				noteStyle(a)
				use(a, f.rel, ln, line)
				continue
			}
			if a, ok := libOf[simple]; ok {
				if a.WrappedBy != "" && !inDesignSystem(f.rel) {
					r.Findings = append(r.Findings, Finding{"bypass", "warn", f.rel, ln, line, "直接用了 " + simple + "（" + a.Source + "），项目设计系统已封装为 " + a.WrappedBy + "，应使用封装"})
				}
				if a.Source == "Android SDK widgets" {
					if concept, raw := rawAndroid[simple]; raw {
						// judged after all usages are counted: only an alternative the project has
						// (its own view, or a library view it actually uses) makes this a missed reuse
						raws = append(raws, rawOcc{a, simple, concept, f.rel, ln, line})
						continue
					}
				}
				noteStyle(a)
				use(a, f.rel, ln, line)
				continue
			}
			// a fully-qualified view from a library the catalog does not know (WanAndroid's
			// per.goweii.actionbarex ActionBarCommon × 24, androidx SwitchCompat): still the project's toolkit
			if i := strings.LastIndex(tag, "."); i > 0 && !strings.HasPrefix(tag, "android.") && !strings.HasPrefix(simple, "AppCompat") && !catalogued(tag) && !ownPkgs[pkgRoot(tag[:i])] {
				src := thirdPartySource(tag[:i])
				a := get("ui-lib", simple, src)
				if a.Concept == "" || a.Concept == "actionbar" && strings.HasPrefix(simple, "ActionBar") {
					a.Concept = androidConcept(simple)
				}
				noteStyle(a)
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
		for _, m := range xmlHardSizeRe.FindAllStringSubmatch(f.text, -1) {
			if m[1] != "0" { // layout_width="0dp" is ConstraintLayout's match-constraint, not a size
				r.Metrics.HardcodedSizes++
			}
		}
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
		composRe = regexp.MustCompile(`\b(` + alt(compos) + `)\s*[({]`) // NiaButton(…), NiaButton { … }, DsCard(DsChip(…))
	}
	if len(custom) > 0 {
		customRe = regexp.MustCompile(`\b(` + alt(custom) + `)\b[(.]`)
	}
	inputCalls, searchCalls := map[string]int{}, map[string]int{}
	for _, f := range code {
		lines := strings.Split(f.text, "\n")
		if composRe != nil {
			for _, m := range composRe.FindAllStringSubmatchIndex(f.text, -1) {
				if m[2] > 0 && f.text[m[2]-1] == '.' || declaredAt(f.text, m[2]) {
					continue // obj.NiaButton(…) / fun NiaButton(
				}
				name := f.text[m[2]:m[3]]
				ln := strings.Count(f.text[:m[0]], "\n") + 1
				use(compos[name], f.rel, ln, strings.TrimSpace(lines[ln-1]))
				if compos[name].Concept == "input" && f.text[m[1]-1] == '(' {
					inputCalls[name]++
					if end := matchParen(f.text, m[1]-1); end > 0 && searchArgRe.MatchString(f.text[m[1]:end]) {
						searchCalls[name]++
					}
				}
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
	// an input whose calls mostly bind a search query (GkAppBarTextField(value = searchStr) × 7) is
	// the project's search box, whatever its name says
	for name, n := range inputCalls {
		if s := searchCalls[name]; s >= 2 && s*2 >= n {
			compos[name].Concept = "search"
			compos[name].ConceptBy = itoa(n) + " 处调用里有 " + itoa(s) + " 处传的是搜索词，是项目的搜索框"
		}
	}
	for _, o := range raws {
		// the alternative must be a generic component of the concept (QMUIRoundButton), not a
		// specialised one (FloatingActionButton is no replacement for a Button)
		alt := bestProjectAsset(assets, o.concept, "custom-view")
		// under a Material theme the inflater already turns <Button> / <CheckBox> / <RadioButton> into
		// MaterialButton / MaterialCheckBox / MaterialRadioButton: nothing was bypassed
		if alt != nil && alt.Source == "Material Components" && (family == "material" || family == "material3") && materialInflated[o.simple] {
			// the platform widget already is the Material one; the project's own component (RoundButton)
			// would still be the reuse, if it has one
			alt = nil
			for _, a := range assets {
				if a.Concept == o.concept && a.Kind != "ui-lib" && a.Kind != "page-local" && a.Usage > 0 && usable(a) && len(specialWords(a.Name, a.Concept)) == 0 &&
					(alt == nil || a.Usage > alt.Usage || a.Usage == alt.Usage && a.Name < alt.Name) {
					alt = a
				}
			}
		}
		if alt != nil && (alt.Kind != "ui-lib" || alt.Usage > 0) && len(specialWords(alt.Name, alt.Concept)) == 0 {
			r.Metrics.RawPrimitives++
			r.Findings = append(r.Findings, Finding{"missed-reuse", softSeverity(o.concept), o.rel, o.ln, o.line, "原生 " + o.simple + "：项目里已有 " + alt.Name + "（" + kindLabel(alt.Kind) + "）"})
			continue
		}
		use(o.a, o.rel, o.ln, o.line) // AppCompat Toolbar in a project without an alternative is just the standard widget
	}
	// Compose material3 / library calls
	for _, l := range libs {
		if !isComposeLib(l.Name) {
			continue
		}
		type call struct {
			f    file
			c    string
			at   int
			line string
		}
		var calls []call
		for _, f := range code {
			uses := false
			for _, p := range l.Packages {
				uses = uses || strings.Contains(f.text, p)
			}
			if !uses {
				continue
			}
			lines := strings.Split(f.text, "\n")
			for _, c := range l.Components {
				// the project defines a composable of the same name (Jetsnack's FilterChip): only an
				// explicit import makes a call the library one
				if compos[c] != nil && !strings.Contains(f.text, "import "+l.Import[0]+c) && !strings.Contains(f.text, "import "+l.Import[0]+"*") {
					continue
				}
				for _, m := range regexp.MustCompile(`\b`+c+`\s*[({]`).FindAllStringIndex(f.text, -1) {
					if m[0] > 0 && f.text[m[0]-1] == '.' || declaredAt(f.text, m[0]) {
						continue
					}
					ln := strings.Count(f.text[:m[0]], "\n") + 1
					calls = append(calls, call{f, c, ln, strings.TrimSpace(lines[ln-1])})
				}
			}
		}
		// a design-system file that calls a library composable wraps it (NiaButton → Button)
		for _, k := range calls {
			lib := assets["ui-lib|"+l.Name+"|"+k.c]
			if lib == nil || lib.WrappedBy != "" || !inDesignSystem(k.f.rel) {
				continue
			}
			// the composable whose body holds the call: the nearest one defined above it in that file
			// (HelpDialog calls AlertDialog; SealDialog, further down, is built on BasicAlertDialog)
			var best *Asset
			for _, d := range composDefs[k.f.rel] {
				if d.line > k.at {
					break
				}
				best = compos[d.name]
			}
			if best != nil && best.Concept != lib.Concept && !strings.HasSuffix(best.Name, k.c) {
				best = nil
			}
			// GkMenuItemRadioButton is a menu row holding a RadioButton and GkRuleCategoryTabs the tabs of one
			// feature, not the project's RadioButton / tabs: a wrapper adds at most its prefix and one word
			if best != nil && (strings.HasSuffix(best.Name, k.c) && len(splitWords(strings.TrimSuffix(best.Name, k.c))) > 2 ||
				!strings.HasSuffix(best.Name, k.c) && featureNamed(*best, lib.Concept)) {
				best = nil
			}
			if best != nil && best.Name != k.c {
				lib.WrappedBy = best.Name
			}
		}
		for _, k := range calls {
			lib := assets["ui-lib|"+l.Name+"|"+k.c]
			use(lib, k.f.rel, k.at, k.line)
			if lib != nil && lib.WrappedBy != "" && !inDesignSystem(k.f.rel) {
				r.Findings = append(r.Findings, Finding{"bypass", "warn", k.f.rel, k.at, k.line, "直接用了 " + k.c + "（" + l.Name + "），项目设计系统已封装为 " + lib.WrappedBy + "，应使用封装"})
			}
		}
	}
	r.ColorTokens = colorTokens
	r.Conventions = conventionsAndroid(layoutMarkup)
	// a "wrapper" the project calls far less than the library component itself (Seal: Button × 16,
	// LinkButton × 4) is a specialisation, not the house component
	dropped := map[string]bool{}
	for _, a := range assets {
		if w := a.WrappedBy; w != "" && a.Kind == "ui-lib" {
			if wa := compos[w]; wa != nil && wa.Usage*2 < a.Usage {
				a.WrappedBy, dropped[a.Name+"|"+w] = "", true
			} else if wa := custom[w]; wa != nil && wa.Usage*2 < a.Usage {
				a.WrappedBy, dropped[a.Name+"|"+w] = "", true
			}
		}
	}
	if len(dropped) > 0 { // the bypass findings raised for it during the scan go too
		kept := r.Findings[:0]
		for _, f := range r.Findings {
			i, j := strings.LastIndex(f.Hint, "封装为 "), strings.LastIndex(f.Hint, "，应使用封装")
			k, e := strings.Index(f.Hint, "直接用了 "), strings.Index(f.Hint, "（")
			if f.Rule == "bypass" && i >= 0 && j > i && k >= 0 && e > k && dropped[strings.TrimSpace(f.Hint[k+len("直接用了 "):e])+"|"+f.Hint[i+len("封装为 "):j]] {
				continue
			}
			kept = append(kept, f)
		}
		r.Findings = kept
	}
	if family == "material" || family == "material3" {
		for _, a := range assets {
			if a.Source == "Android SDK widgets" && materialInflated[a.Name] {
				a.Concept = "" // inflated as its Material version: the same component, not a second one
			}
		}
	}
	finish(r, assets, colorHits)
	r.Clones = findClones(layoutMarkup, func(s string) []*mnode { return dropContainers(parseMarkup(s)) }, 5)
	r.Metrics.CloneGroups = len(r.Clones)
	return r
}

// bestProjectAsset picks the replacement for a raw widget; prefer is the kind usable in that context
// (custom-view inside XML layouts, composable inside Compose code).
// softSeverity: a raw image is often fine (icons, plain pictures), so it is reported but never fails a page.
func softSeverity(concept string) string {
	if concept == "image" {
		return "info"
	}
	return "warn"
}

func bestProjectAsset(assets map[string]*Asset, concept, prefer string) *Asset {
	var best *Asset
	score := func(a *Asset) int {
		s := rank(*a)
		if a.Kind == prefer || prefer == "custom-view" && a.Kind == "ui-lib" && !isComposeLib(a.Source) {
			s += 5000
		}
		if prefer == "custom-view" && (a.Kind == "composable" || isComposeLib(a.Source)) {
			s -= 5000 // a composable cannot be dropped into an XML layout
		}
		return s
	}
	for _, a := range assets {
		if a.Concept == concept && usable(a) && (a.Kind == "custom-view" || a.Kind == "composable" || a.Kind == "design-system" || a.Kind == "business" || a.Kind == "ui-lib" && a.Source != "Android SDK widgets") && (best == nil || score(a) > score(best) || score(a) == score(best) && (len(a.Name) < len(best.Name) || len(a.Name) == len(best.Name) && a.Name+"|"+a.Source < best.Name+"|"+best.Source)) {
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

var commentRe = regexp.MustCompile(`(?s)<!--.*?-->`)

// stripComments blanks out markup comments but keeps line numbers.
func stripComments(s string) string {
	return commentRe.ReplaceAllStringFunc(s, func(c string) string {
		return strings.Repeat("\n", strings.Count(c, "\n"))
	})
}

// matchParen returns the index of the ')' closing the '(' at open.
// funcBody is the body of a Kotlin function whose parameter list ends right before from: the
// { … } block after an optional return type, or an expression body up to the next blank line.
func funcBody(s string, from int) string {
	for i := from; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth := 0
			for j := i; j < len(s); j++ {
				switch s[j] {
				case '{':
					depth++
				case '}':
					if depth--; depth == 0 {
						return s[i : j+1]
					}
				}
			}
			return s[i:]
		case '=':
			if e := strings.Index(s[i:], "\n\n"); e > 0 {
				return s[i : i+e]
			}
			return s[i:]
		}
	}
	return ""
}

func matchParen(s string, open int) int {
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// refineConcept resolves ambiguous names with the props: a "Banner" taking text / tone is a notice
// bar, one taking images / items is a carousel.
// thirdPartySource names an uncatalogued view library by its package (androidx.appcompat.widget →
// AppCompat, per.goweii.actionbarex.common → per.goweii.actionbarex).
func thirdPartySource(pkg string) string {
	if strings.HasPrefix(pkg, "androidx.appcompat") {
		return "AppCompat"
	}
	if strings.HasPrefix(pkg, "androidx.") {
		return "AndroidX"
	}
	if ps := strings.Split(pkg, "."); len(ps) > 3 {
		pkg = strings.Join(ps[:3], ".")
	}
	return pkg
}

var pkgDeclRe = regexp.MustCompile(`(?m)^\s*package\s+([\w.]+)`)

// pkgRoot: the first three segments of a package (org.schabi.newpipe.views → org.schabi.newpipe).
func pkgRoot(pkg string) string {
	if ps := strings.Split(pkg, "."); len(ps) > 3 {
		return strings.Join(ps[:3], ".")
	}
	return pkg
}

// catalogued: the class belongs to a library the catalog lists (QMUI's internal LineTypeView): only its
// listed components are public ones.
func catalogued(tag string) bool {
	for _, l := range uiLibs {
		for _, imp := range l.Import {
			if root := libRoot(imp); root != "" && strings.HasPrefix(tag, root+".") {
				return true
			}
		}
	}
	return false
}

// libRoot: a library's package up to its first distinctive segment (com.google.android.material,
// com.qmuiteam), so com.google.android.exoplayer2 is not mistaken for Material.
func libRoot(imp string) string {
	ps := strings.Split(strings.TrimSuffix(imp, "."), ".")
	for i, p := range ps {
		switch p {
		case "com", "org", "io", "net", "me", "cn", "google", "android", "androidx":
			continue
		}
		return strings.Join(ps[:i+1], ".")
	}
	return ""
}

// androidConcept: on Android an "action bar" is the top app bar (ActionBarCommon, ActionBarEx).
func androidConcept(name string) string {
	if strings.HasPrefix(name, "ActionBar") || strings.HasPrefix(name, "TitleBar") {
		return "navbar"
	}
	return ConceptOf(name)
}

func refineConcept(a *Asset) {
	refineBadge(a)
	// a "bottom bar" that takes tabs / routes is bottom navigation, not an action bar
	if a.Concept == "actionbar" {
		for _, p := range a.Props {
			switch normProp(p) {
			case "tabs", "items", "currentroute", "selectedtab", "selectedindex", "navigatetoroute", "destinations", "onnavigate":
				a.Concept = "tabbar"
				return
			}
		}
	}
	if a.Concept != "carousel" || !strings.Contains(strings.ToLower(a.Name), "banner") {
		return
	}
	for _, p := range a.Props {
		switch strings.ToLower(p) {
		case "text", "tone", "message", "content", "type", "title":
			a.Concept = "notice"
			return
		}
	}
	// a composable banner with nothing to page through (GkMatchingBanner(enabled, onEnable) is one
	// feature's prompt card) is no carousel, and not the generic notice bar either
	if a.Kind == "composable" {
		for _, p := range a.Props {
			if bannerContentRe.MatchString(p) {
				return
			}
		}
		a.Concept = ""
	}
}

// isComposeLib: a Jetpack Compose library catalog (composables, not Views).
// splitTopLevel splits a parameter list at commas outside (), <>, {}, [] and quotes:
// onPick: (id: String, index: Int) -> Unit stays one parameter.
func splitTopLevel(s string) []string {
	var out []string
	depth, q, start := 0, byte(0), 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case q != 0:
			if c == q && s[i-1] != '\\' {
				q = 0
			}
		case c == '"' || c == '\'':
			q = c
		case c == '(' || c == '<' || c == '{' || c == '[':
			depth++
		case c == ')' || c == '}' || c == ']' || c == '>' && (i == 0 || s[i-1] != '-'): // not the -> of a lambda type
			depth--
		case c == ',' && depth == 0:
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}

// paramName strips annotations and modifiers: "@StringRes labelRes" → labelRes, "vararg tags" → tags.
func paramName(s string) string {
	f := strings.Fields(s)
	for len(f) > 1 && (strings.HasPrefix(f[0], "@") || f[0] == "vararg" || f[0] == "noinline" || f[0] == "crossinline") {
		f = f[1:]
	}
	if len(f) == 0 {
		return ""
	}
	return strings.TrimPrefix(f[len(f)-1], "@")
}

// androidHex normalises an Android colour: #rgb / #argb expanded, upper case, opaque alpha dropped
// (#fff, #FFFFFFFF and #ffffff are one colour).
func androidHex(v string) string {
	h := strings.ToUpper(strings.TrimPrefix(v, "#"))
	if len(h) == 3 || len(h) == 4 {
		var b []byte
		for i := 0; i < len(h); i++ {
			b = append(b, h[i], h[i])
		}
		h = string(b)
	}
	if len(h) == 8 && strings.HasPrefix(h, "FF") {
		h = h[2:]
	}
	return "#" + h
}

// declaredAt: the name at idx is being declared (fun / class / object / interface X), not called.
func declaredAt(src string, idx int) bool {
	return declKwRe.MatchString(src[max(0, idx-16):idx])
}

var declKwRe = regexp.MustCompile(`\b(?:fun|class|object|interface)\s+$`)

func isComposeLib(source string) bool { return strings.HasPrefix(source, "Compose ") }

// domainType: a parameter type that is an app data model (Snack, Topic, Order) rather than UI
// plumbing (String, Color, Modifier, ImageVector, lambdas, generics).
func domainType(t string) bool {
	if i := strings.IndexByte(t, '.'); i > 0 && !strings.ContainsAny(t, "()<>-") && t[0] >= 'A' && t[0] <= 'Z' && t[i+1:] != "" && t[i+1] >= 'A' && t[i+1] <= 'Z' {
		return domainType(t[:i]) // Task.DownloadState is a model; Alignment.Horizontal is Compose's own
	}
	if t == "" || t[0] < 'A' || t[0] > 'Z' || strings.ContainsAny(t, "()<>-.") {
		return false
	}
	switch t {
	case "String", "Int", "Long", "Float", "Double", "Boolean", "Color", "Dp", "TextUnit", "Modifier", "ImageVector", "Painter",
		"Shape", "TextStyle", "Brush", "Alignment", "Arrangement", "PaddingValues", "AnnotatedString", "Unit", "Any", "Text",
		"MutableInteractionSource", "BorderStroke", "ButtonColors", "CardColors", "ImageBitmap", "Uri", "Bitmap", "LazyListState", "ScrollState":
		return false
	}
	return !strings.HasSuffix(t, "Colors") && !strings.HasSuffix(t, "State") && !strings.HasSuffix(t, "Style") && !strings.HasSuffix(t, "Tone") && !strings.HasSuffix(t, "Size") && !strings.HasSuffix(t, "Type") && !strings.HasSuffix(t, "Variant")
}

// annotationsBefore returns the annotation lines directly above a declaration (@Preview, @ThemePreviews …).
func annotationsBefore(src string, at int) string {
	lines := strings.Split(src[max(0, at-400):at], "\n")
	var ann []string
	for i := len(lines) - 1; i >= 0; i-- {
		l := strings.TrimSpace(lines[i])
		if l == "" && i == len(lines)-1 {
			continue
		}
		if !strings.HasPrefix(l, "@") {
			break
		}
		ann = append(ann, l)
	}
	return strings.Join(ann, " ")
}

// refineBadge: a "badge" that takes text / tone and no count is a status tag (StatusBadge, DsBadge).
func refineBadge(a *Asset) {
	if a.Concept != "badge" {
		return
	}
	text := false
	for _, p := range a.Props {
		switch normProp(p) {
		case "count", "max", "dot", "isdot", "number", "badgecount", "showzero":
			return
		case "text", "label", "tone", "type", "content", "title":
			text = true
		}
	}
	if text {
		a.Concept = "tag"
	}
}
