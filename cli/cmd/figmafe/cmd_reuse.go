package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	for _, p := range plats {
		var r *reuse.Report
		if p == "h5" {
			r = reuse.ScanH5(root)
		} else {
			r = reuse.ScanAndroid(root)
		}
		md := r.Markdown()
		var pc *reuse.PageCheck
		if *focus != "" {
			pc = r.CheckFocus(*focus)
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
		reuse.Confirmed = nil
		for _, cand := range []string{"design/component-map.yaml", ".figmafe/component-map.yaml", "component-map.yaml"} {
			if cm, err := mapping.Load(filepath.Join(root, cand)); err == nil {
				reuse.Confirmed = cm
				fmt.Printf("· 使用已确认映射 %s\n", cand)
				break
			}
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
