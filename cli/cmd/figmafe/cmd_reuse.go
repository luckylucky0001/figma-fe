package main

import (
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
	pos := parse(fs, args)
	root := "."
	if len(pos) > 0 {
		root = pos[0]
	}
	root, _ = filepath.Abs(root)
	if *out == "" {
		*out = filepath.Join(root, ".figmafe", "reuse")
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
		fmt.Printf("\n✓ [%s] %d 个资产 · 组件化率 %.0f%% · 重复概念 %d · 重复结构 %d 组 · 硬编码颜色 %d → %s\n",
			p, len(r.Assets), m.Componentization*100, m.DuplicateConcepts, m.CloneGroups, m.HardcodedColors, filepath.Join(*out, p+".md"))
	}
	return nil
}

// runMatch: a new Figma page + an existing project → which existing components can be reused.
func runMatch(args []string) error {
	fs := flag.NewFlagSet("match", flag.ExitOnError)
	figmaDir := fs.String("figma", "", "Figma page export dir (export.json)")
	target := fs.String("target", "", "android page technology: xml|compose (ranks usable components first)")
	pf := fs.String("platform", "auto", "h5|android|auto")
	out := fs.String("o", "", "output dir (default <project>/.figmafe/match)")
	quiet := fs.Bool("q", false, "quiet")
	pos := parse(fs, args)
	if *figmaDir == "" || len(pos) != 1 {
		return fmt.Errorf("usage: figmafe match --figma <export dir> <project dir> [--platform h5|android]")
	}
	exp, err := figma.LoadExport(*figmaDir)
	if err != nil {
		return err
	}
	root, _ := filepath.Abs(pos[0])
	if *out == "" {
		*out = filepath.Join(root, ".figmafe", "match")
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
		m := reuse.Match(exp, r)
		md := m.Markdown()
		_ = os.WriteFile(filepath.Join(*out, p+".md"), []byte(md), 0o644)
		writeJSON(filepath.Join(*out, p+".json"), m)
		if b, err := yaml.Marshal(m.Draft()); err == nil {
			hdr := "# component-map 初稿（figmafe match 自动生成）：人工确认推荐组件、补齐变体 → props 后，合并进 design/component-map.yaml\n"
			_ = os.WriteFile(filepath.Join(*out, "component-map.draft."+p+".yaml"), append([]byte(hdr), b...), 0o644)
		}
		if !*quiet {
			fmt.Print(md)
		}
		fmt.Printf("\n✓ [%s] Figma 组件 %d 个：可直接复用 %d · 需确认 %d · 需新建 %d → %s\n", p, len(m.Items), m.Reuse, m.Confirm, m.New, filepath.Join(*out, p+".md"))
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
		if cm, err := mapping.Load(filepath.Join(root, cand)); err == nil {
			return cm
		}
	}
	return nil
}
