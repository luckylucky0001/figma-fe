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
	if *focus != "" && *changed == "" {
		// a mistyped path would otherwise check nothing and pass: in CI that silently disables the gate
		for _, f := range strings.Split(*focus, ",") {
			if f = strings.TrimSpace(f); f != "" {
				if _, err := os.Stat(filepath.Join(root, f)); err != nil {
					return fmt.Errorf("--focus %s: not found under %s", f, root)
				}
			}
		}
	}
	if *changed != "" {
		files, err := changedUIFiles(root, *changed)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Println("✓ 没有改动 UI 文件，跳过页面复用检查")
			return nil
		}
		*focus = strings.Join(files, reuse.FocusSep)
		lines = addedLines(root, *changed)
		fmt.Printf("· 检查改动的 %d 个 UI 文件（已有文件只看新增的行）\n", len(files))
	}
	regressed, failed := false, false // reported after every platform ran, so real errors are not swallowed
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
			var probe map[string]json.RawMessage
			b, err := os.ReadFile(*baseline)
			if err != nil || json.Unmarshal(b, &probe) != nil || probe["metrics"] == nil || json.Unmarshal(b, &old) != nil || old.Platform == "" {
				return fmt.Errorf("--baseline %s: not a figmafe reuse report (<platform>.json)", *baseline)
			}
			if old.Platform != p {
				fmt.Printf("· --baseline 是 %s 的报告，跳过 %s 的对比\n", old.Platform, p)
			} else {
				text, worse := reuse.CompareMetrics(&old, r)
				md = text + "\n---\n\n" + md
				fmt.Print(text)
				if worse && *noRegress {
					regressed = true
				}
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
			failed = true
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
	switch {
	case failed:
		return exitError{3, "✗ 页面复用检查未通过"}
	case regressed:
		return exitError{3, "✗ 和上次体检相比有指标变差（--no-regress）"}
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

// gitOut runs git with stable path output: no quoting of non-ASCII names, a/ b/ prefixes whatever
// the user's diff.noprefix / mnemonicPrefix, and renames reported as delete + add.
func gitOut(root string, args ...string) ([]byte, error) {
	base := []string{"-C", root, "-c", "core.quotePath=false", "-c", "diff.noprefix=false", "-c", "diff.mnemonicPrefix=false"}
	return exec.Command("git", append(base, args...)...).Output()
}

// nulList splits git's -z output.
func nulList(b []byte) []string {
	var out []string
	for _, f := range strings.Split(string(b), "\x00") {
		if f = strings.TrimSpace(f); f != "" {
			out = append(out, f)
		}
	}
	return out
}

func uiFile(f string) bool {
	ext := filepath.Ext(f)
	// style sheets too: a hard-coded colour added in page.scss is part of the change
	return ext == ".vue" || ext == ".nvue" || ext == ".tsx" || ext == ".jsx" || ext == ".kt" || ext == ".java" ||
		ext == ".css" || ext == ".scss" || ext == ".less" || ext == ".sass" ||
		ext == ".xml" && (strings.Contains("/"+f, "/res/layout") || strings.Contains("/"+f, "/res/values"))
}

// changedUIFiles lists added / modified UI source files vs ref (a renamed file counts as added).
func changedUIFiles(root, ref string) ([]string, error) {
	out, err := gitOut(root, "diff", "--name-only", "-z", "--relative", "--no-renames", "--diff-filter=AM", ref, "--")
	if err != nil {
		return nil, fmt.Errorf("git diff %s: %w", ref, err)
	}
	untracked, _ := gitOut(root, "ls-files", "-z", "--others", "--exclude-standard")
	var files []string
	seen := map[string]bool{}
	for _, f := range append(nulList(out), nulList(untracked)...) {
		if uiFile(f) && !seen[f] {
			seen[f] = true
			files = append(files, f)
		}
	}
	return files, nil
}

var hunkHeadRe = regexp.MustCompile(`^@@ -\d+(?:,\d+)? \+(\d+)(?:,(\d+))? @@`)

// addedLines: for files modified (not added) since ref, the line numbers this change added, so the
// gate judges the change and not the file's history. Files missing from the map are judged whole.
func addedLines(root, ref string) map[string]map[int]bool {
	names, err := gitOut(root, "diff", "--name-only", "-z", "--relative", "--no-renames", "--diff-filter=M", ref, "--")
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠ git diff %s 失败，改动过的文件将整份检查：%v\n", ref, err)
		return nil
	}
	res := map[string]map[int]bool{}
	for _, f := range nulList(names) {
		res[f] = map[int]bool{} // a mode-only change adds no lines
	}
	out, err := gitOut(root, "diff", "-U0", "--relative", "--no-renames", "--diff-filter=M", "--src-prefix=a/", "--dst-prefix=b/", ref, "--")
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠ git diff %s 失败，改动过的文件将整份检查：%v\n", ref, err)
		return nil
	}
	cur := ""
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "+++ ") {
			p := strings.TrimRight(strings.TrimPrefix(line, "+++ "), "\t")
			if strings.HasPrefix(p, `"`) {
				if u, err := strconv.Unquote(p); err == nil {
					p = u
				}
			}
			cur = ""
			if p != "/dev/null" {
				cur = strings.TrimPrefix(p, "b/")
				if res[cur] == nil {
					res[cur] = map[int]bool{}
				}
			}
			continue
		}
		m := hunkHeadRe.FindStringSubmatch(line)
		if m == nil || cur == "" {
			continue
		}
		start, _ := strconv.Atoi(m[1])
		n := 1
		if m[2] != "" {
			n, _ = strconv.Atoi(m[2])
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
