package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
	"github.com/luckylucky0001/figma-fe/cli/internal/mapping"
	"github.com/luckylucky0001/figma-fe/cli/internal/reuse"
)

// runReuse is the brownfield reuse check. It needs no figmafe.yaml and no Figma: point it at any
// Vue / React H5 project or Android (XML View / Compose) project.
func runReuse(args []string) error {
	fs := flag.NewFlagSet("reuse", flag.ExitOnError)
	pf := fs.String("platform", "auto", "h5|android|auto")
	out := fs.String("o", "", "output dir (default <root>/.figmafe/reuse)")
	quiet := fs.Bool("q", false, "quiet")
	focus := fs.String("focus", "", "only judge this page dir/file (relative to project); the whole project is the asset pool")
	changed := fs.String("changed", "", "PR gate: judge files changed vs this git ref (e.g. origin/main)")
	baseline := fs.String("baseline", "", "earlier <platform>.json report: print what improved / got worse")
	noRegress := fs.Bool("no-regress", false, "with --baseline: exit 3 when any project-wide metric got worse")
	pos := parse(fs, args)
	root := "."
	if len(pos) > 0 {
		root = pos[0]
	}
	root, _ = filepath.Abs(root)
	if *out == "" {
		*out = filepath.Join(root, ".figmafe", "reuse")
		ignoreWorkDir(root)
	}
	_ = os.MkdirAll(*out, 0o755)
	plats := []string{*pf}
	if *pf == "auto" {
		plats = reuse.Detect(root)
		if len(plats) == 0 {
			return fmt.Errorf("%s: no H5 (package.json) or Android (build.gradle) project found", root)
		}
	}
	var lines map[string]map[int]bool
	if *changed != "" {
		files, err := changedUIFiles(root, *changed)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Println("✓ 没有改动 UI 文件，跳过页面复用检查")
			return nil
		}
		*focus = strings.Join(files, ",")
		lines = addedLines(root, *changed)
		fmt.Printf("· 检查改动的 %d 个 UI 文件（已有文件只看新增的行）\n", len(files))
	}
	regressed := false
	defer func() {
		if regressed {
			os.Exit(3)
		}
	}()
	for _, p := range plats {
		var r *reuse.Report
		if p == "h5" {
			r = reuse.ScanH5(root)
		} else {
			r = reuse.ScanAndroid(root)
		}
		md := r.Markdown()
		reuse.Confirmed = loadConfirmed(root)
		var pc *reuse.PageCheck
		if *focus != "" {
			pc = r.CheckFocusLines(*focus, lines)
			md = pc.Markdown() + "\n---\n\n" + md
			writeJSON(filepath.Join(*out, p+".page.json"), pc)
		}
		if *baseline != "" {
			var old reuse.Report
			if b, err := os.ReadFile(*baseline); err != nil || json.Unmarshal(b, &old) != nil {
				return fmt.Errorf("--baseline %s: not a figmafe reuse report", *baseline)
			}
			text, worse := reuse.CompareMetrics(&old, r)
			md = text + "\n---\n\n" + md
			fmt.Print(text)
			if worse && *noRegress {
				regressed = true
			}
		}
		_ = os.WriteFile(filepath.Join(*out, p+".md"), []byte(md), 0o644)
		writeJSON(filepath.Join(*out, p+".json"), r)
		if !*quiet {
			fmt.Print(md)
		} else if pc != nil {
			fmt.Print(pc.Markdown())
		}
		if pc != nil && !pc.Pass {
			defer func() { os.Exit(3) }()
		}
		m := r.Metrics
		wraps, blocked := 0, 0
		for i := range r.LibAll {
			if r.LibAll[i].WrappedBy != "" {
				wraps++
			}
			if r.LibAll[i].Unavailable || r.LibAll[i].Caveat != "" && r.LibAll[i].Usage == 0 {
				blocked++
			}
		}
		extra := ""
		if wraps > 0 {
			extra += fmt.Sprintf(" · 设计系统封装 %d", wraps)
		}
		if blocked > 0 {
			extra += fmt.Sprintf(" · 不能直接用的库组件 %d", blocked)
		}
		if len(r.Conventions) > 0 {
			extra += fmt.Sprintf(" · 项目惯例 %d", len(r.Conventions))
		}
		fmt.Printf("\n✓ [%s] %d 个资产 · 组件化率 %.0f%% · 重复概念 %d · 重复结构 %d 组 · 硬编码颜色 %d%s → %s\n",
			p, len(r.Assets), m.Componentization*100, m.DuplicateConcepts, m.CloneGroups, m.HardcodedColors, extra, filepath.Join(*out, p+".md"))
	}
	return nil
}

// runMatch: a new Figma page + an existing project → which existing components can be reused.
func runMatch(args []string) error {
	fs := flag.NewFlagSet("match", flag.ExitOnError)
	figmaDir := fs.String("figma", "", "Figma page export dir (export.json), or a dir of page export dirs (a whole feature)")
	target := fs.String("target", "", "android page technology: xml|compose (ranks usable components first)")
	pf := fs.String("platform", "auto", "h5|android|auto")
	out := fs.String("o", "", "output dir (default <project>/.figmafe/match)")
	quiet := fs.Bool("q", false, "quiet")
	pos := parse(fs, args)
	if *figmaDir == "" || len(pos) != 1 {
		return fmt.Errorf("usage: figmafe match --figma <export dir | dir of exports> <project dir> [--platform h5|android]")
	}
	// one page, or every page export under the dir (a feature's pages)
	pages := []string{*figmaDir}
	multi := false
	if _, err := os.Stat(filepath.Join(*figmaDir, "export.json")); err != nil {
		pages = nil
		entries, _ := os.ReadDir(*figmaDir)
		for _, e := range entries {
			if _, err := os.Stat(filepath.Join(*figmaDir, e.Name(), "export.json")); e.IsDir() && err == nil {
				pages = append(pages, filepath.Join(*figmaDir, e.Name()))
			}
		}
		if len(pages) == 0 {
			return fmt.Errorf("%s: no export.json (and no page export dirs inside)", *figmaDir)
		}
		multi = true
	}
	root, _ := filepath.Abs(pos[0])
	if *out == "" {
		*out = filepath.Join(root, ".figmafe", "match")
		ignoreWorkDir(root)
	}
	_ = os.MkdirAll(*out, 0o755)
	plats := []string{*pf}
	if *pf == "auto" {
		plats = reuse.Detect(root)
	}
	for _, p := range plats {
		var r *reuse.Report
		if p == "h5" {
			r = reuse.ScanH5(root)
		} else {
			r = reuse.ScanAndroid(root)
		}
		reuse.Target = *target
		if p == "android" && *target == "" {
			if reuse.Target = reuse.InferTarget(r); reuse.Target != "" {
				fmt.Printf("· 未指定 --target，按项目主流写法推断新页面用 %s（不对请显式传 --target xml|compose）\n", reuse.Target)
			}
		}
		if reuse.Confirmed = loadConfirmed(root); reuse.Confirmed != nil {
			fmt.Println("· 使用已确认映射 component-map.yaml")
		}
		var results []*reuse.MatchResult
		for _, page := range pages {
			exp, err := figma.LoadExport(page)
			if err != nil {
				return err
			}
			dir := *out
			if multi {
				dir = filepath.Join(*out, filepath.Base(page))
				_ = os.MkdirAll(dir, 0o755)
			}
			m := reuse.Match(exp, r)
			results = append(results, m)
			md := m.Markdown()
			if exp.Images.Frame != "" {
				if b, err := os.ReadFile(filepath.Join(page, exp.Images.Frame)); err == nil {
					img := p + ".frame.png"
					_ = os.WriteFile(filepath.Join(dir, img), b, 0o644)
					_ = os.WriteFile(filepath.Join(dir, p+".html"), []byte(m.HTML(exp, img)), 0o644)
				}
			}
			_ = os.WriteFile(filepath.Join(dir, p+".md"), []byte(md), 0o644)
			writeJSON(filepath.Join(dir, p+".json"), m)
			if b, err := yaml.Marshal(m.Draft()); err == nil {
				hdr := "# component-map 初稿（figmafe match 自动生成）：人工确认推荐组件、补齐变体 → props 后，合并进 design/component-map.yaml\n"
				_ = os.WriteFile(filepath.Join(dir, "component-map.draft."+p+".yaml"), append([]byte(hdr), b...), 0o644)
			}
			if !*quiet && !multi {
				fmt.Print(md)
			}
			fmt.Printf("\n✓ [%s] %s：Figma 组件 %d 个：可直接复用 %d · 需确认 %d · 需新建 %d → %s\n", p, m.Page, len(m.Items), m.Reuse, m.Confirm, m.New, filepath.Join(dir, p+".md"))
		}
		if multi {
			sum := reuse.FeatureSummary(results)
			_ = os.WriteFile(filepath.Join(*out, p+".summary.md"), []byte(sum), 0o644)
			if !*quiet {
				fmt.Print(sum)
			}
			fmt.Printf("✓ [%s] %d 个页面的汇总 → %s\n", p, len(pages), filepath.Join(*out, p+".summary.md"))
		}
	}
	return nil
}

// changedUIFiles lists added / modified UI source files vs ref (git diff --diff-filter=AM).
func changedUIFiles(root, ref string) ([]string, error) {
	out, err := exec.Command("git", "-C", root, "diff", "--name-only", "--relative", "--diff-filter=AM", ref, "--").Output()
	if err != nil {
		return nil, fmt.Errorf("git diff %s: %w", ref, err)
	}
	untracked, _ := exec.Command("git", "-C", root, "ls-files", "--others", "--exclude-standard").Output()
	var files []string
	seen := map[string]bool{}
	for _, f := range strings.Fields(string(out) + "\n" + string(untracked)) {
		ext := filepath.Ext(f)
		ui := ext == ".vue" || ext == ".tsx" || ext == ".jsx" || ext == ".kt" || ext == ".java" ||
			ext == ".xml" && (strings.Contains("/"+f, "/res/layout") || strings.Contains("/"+f, "/res/values"))
		if ui && !seen[f] {
			seen[f] = true
			files = append(files, f)
		}
	}
	return files, nil
}

var hunkRe = regexp.MustCompile(`(?m)^\+\+\+ b/(.+)$|^@@ -\d+(?:,\d+)? \+(\d+)(?:,(\d+))? @@`)

// addedLines: for files modified (not added) since ref, the line numbers this change added, so the
// gate judges the change and not the file's history. Files missing from the map are judged whole.
func addedLines(root, ref string) map[string]map[int]bool {
	out, err := exec.Command("git", "-C", root, "diff", "-U0", "--relative", "--diff-filter=M", ref, "--").Output()
	if err != nil {
		return nil
	}
	res := map[string]map[int]bool{}
	cur := ""
	for _, m := range hunkRe.FindAllStringSubmatch(string(out), -1) {
		if m[1] != "" {
			cur = m[1]
			res[cur] = map[int]bool{}
			continue
		}
		start, _ := strconv.Atoi(m[2])
		n := 1
		if m[3] != "" {
			n, _ = strconv.Atoi(m[3])
		}
		for i := 0; i < n; i++ {
			res[cur][start+i] = true
		}
	}
	return res
}

// loadConfirmed reads the project's confirmed component map, if it has one.
func loadConfirmed(root string) *mapping.Map {
	for _, cand := range []string{"design/component-map.yaml", ".figmafe/component-map.yaml", "component-map.yaml"} {
		p := filepath.Join(root, cand)
		if _, err := os.Stat(p); err != nil {
			continue
		}
		cm, err := mapping.Load(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠ %s 解析失败，本次忽略其中的已确认映射和决定：%v\n", cand, err)
			return nil
		}
		return cm
	}
	return nil
}

// ignoreWorkDir keeps figmafe's reports out of the project's git status (.figmafe/.gitignore = "*").
func ignoreWorkDir(root string) {
	dir := filepath.Join(root, ".figmafe")
	_ = os.MkdirAll(dir, 0o755)
	if _, err := os.Stat(filepath.Join(dir, ".gitignore")); err != nil {
		_ = os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("# figmafe reports (generated)\n*\n"), 0o644)
	}
}
