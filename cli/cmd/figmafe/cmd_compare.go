package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/luckylucky0001/figma-fe/cli/internal/capture"
	"github.com/luckylucky0001/figma-fe/cli/internal/diff"
	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
)

// runCompare: the brownfield screenshot check. Any running H5 page (dev server, test environment,
// a static file) against a Figma export — no figmafe.yaml, no data-figma-id needed. Reports the
// pixel difference with a diff image, and which Figma texts are missing from the page.
func runCompare(args []string) error {
	fs := flag.NewFlagSet("compare", flag.ExitOnError)
	url := fs.String("url", "", "page to check (http(s)://… or a local .html file)")
	figmaDir := fs.String("figma", "", "Figma page export dir (export.json + frame@2x.png)")
	out := fs.String("o", ".figmafe/compare", "output dir")
	max := fs.Float64("max-diff", 0.05, "pixel difference allowed (0.05 = 5%)")
	wait := fs.Duration("wait", 500*time.Millisecond, "extra wait after load (data, animations)")
	parse(fs, args)
	if *url == "" || *figmaDir == "" {
		return fmt.Errorf("usage: figmafe compare --url <page> --figma <export dir> [-o dir] [--max-diff 0.05]")
	}
	exp, err := figma.LoadExport(*figmaDir)
	if err != nil {
		return err
	}
	fb := exp.Frame.AbsoluteBoundingBox
	if fb == nil {
		return fmt.Errorf("export has no frame size")
	}
	scale := exp.Images.Scale
	if scale == 0 {
		scale = 2
	}
	w, h := int(fb.Width+0.5), int(fb.Height+0.5)
	target := *url
	if !strings.Contains(target, "://") {
		abs, _ := filepath.Abs(target)
		target = "file://" + abs
	}
	_ = os.MkdirAll(*out, 0o755)
	b, err := capture.NewBrowser(60 * time.Second)
	if err != nil {
		return err
	}
	defer b.Close()
	if err := b.Open(target, w, h, scale); err != nil {
		return fmt.Errorf("open %s: %w", target, err)
	}
	time.Sleep(*wait)
	screen := filepath.Join(*out, "screen.png")
	if err := b.Screenshot(screen, w, h, scale); err != nil {
		return err
	}
	ratio, sameSize, _, err := diff.Pixel(filepath.Join(*figmaDir, exp.Images.Frame), screen, filepath.Join(*out, "diff.png"))
	if err != nil {
		return err
	}
	var text string
	_ = b.Eval(`document.body.innerText`, &text)
	flat := strings.Join(strings.Fields(text), "")
	var missing, all []string
	figma.Walk(exp.Frame, func(n *figma.Node) bool {
		if n.Type == "TEXT" && strings.TrimSpace(n.Characters) != "" {
			all = append(all, n.Characters)
			if !strings.Contains(flat, strings.Join(strings.Fields(n.Characters), "")) {
				missing = append(missing, n.Characters)
			}
		}
		return true
	})
	pass := ratio <= *max && len(missing) == 0
	var md strings.Builder
	verdict := "✅ PASS"
	if !pass {
		verdict = "❌ FAIL"
	}
	fmt.Fprintf(&md, "# 截图比对：%s  %s\n\n- 页面：%s（%d×%d @%gx）\n- 设计稿：%s\n", exp.Frame.Name, verdict, *url, w, h, scale, *figmaDir)
	fmt.Fprintf(&md, "- 像素差异：%.2f%%（阈值 %.0f%%）%s → diff.png（红色为不一致的像素）\n", ratio*100, *max*100, map[bool]string{false: "，截图尺寸和设计图不同，按重叠部分比较", true: ""}[sameSize])
	fmt.Fprintf(&md, "- Figma 文字覆盖：%d / %d\n", len(all)-len(missing), len(all))
	for i, m := range missing {
		if i == 15 {
			fmt.Fprintf(&md, "  - … 共 %d 处\n", len(missing))
			break
		}
		fmt.Fprintf(&md, "  - 页面上没有：「%s」\n", m)
	}
	md.WriteString("\n> 没有 data-figma-id 时只能做整页像素比对和文字覆盖；需要逐节点的位置 / 字号 / 颜色比对，按 docs/CONTRACTS.md 给组件挂上 Figma 节点 id 后用 verify。\n")
	_ = os.WriteFile(filepath.Join(*out, "report.md"), []byte(md.String()), 0o644)
	fmt.Print(md.String())
	if !pass {
		os.Exit(3)
	}
	return nil
}
