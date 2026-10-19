package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/luckylucky0001/figma-fe/cli/internal/audit"
	"github.com/luckylucky0001/figma-fe/cli/internal/capture"
	"github.com/luckylucky0001/figma-fe/cli/internal/config"
	"github.com/luckylucky0001/figma-fe/cli/internal/diff"
	"github.com/luckylucky0001/figma-fe/cli/internal/mapping"
	"github.com/luckylucky0001/figma-fe/cli/internal/report"
)

func checkArgs(a *pageArgs, cmd string) error {
	if a.page == "" || (a.platform != "h5" && a.platform != "android") {
		return fmt.Errorf("usage: figmafe %s --page <Page> --platform h5|android", cmd)
	}
	return nil
}

func pageSrcDir(cfg *config.Config, platform, page string) string {
	if platform == "h5" {
		return cfg.P(cfg.H5.Root, cfg.H5.Pages, page)
	}
	return cfg.P(cfg.Android.Root, cfg.Android.Screens, page)
}

func doAudit(cfg *config.Config, a *pageArgs) (*audit.Result, error) {
	p, _, _, err := buildPlan(cfg, a.page, a.platform)
	if err != nil {
		return nil, err
	}
	invs, err := loadInventories(cfg, false)
	if err != nil {
		return nil, err
	}
	dir := pageSrcDir(cfg, a.platform, a.page)
	if _, err := os.Stat(dir); err != nil {
		return nil, fmt.Errorf("page source dir %s not found", relTo(cfg, dir))
	}
	var dups map[string]string
	if m, err := mapping.Load(cfg.P(cfg.Map)); err == nil {
		dups = mapping.Duplicates(m, invs[a.platform], a.platform)
	}
	res, err := audit.Run(a.platform, cfg.Root, dir, p, invs[a.platform], dups)
	if err != nil {
		return nil, err
	}
	writeJSON(filepath.Join(cfg.OutDir(a.platform, a.page), "audit.json"), res)
	return res, nil
}

func runAudit(args []string) error {
	a, _ := pageFlags("audit", args)
	if err := checkArgs(a, "audit"); err != nil {
		return err
	}
	cfg, err := load(a.root)
	if err != nil {
		return err
	}
	res, err := doAudit(cfg, a)
	if err != nil {
		return err
	}
	r := report.Build(a.page, a.platform, res, nil, cfg.Thresholds)
	fmt.Print(r.Markdown(40))
	return nil
}

func doCapture(cfg *config.Config, a *pageArgs, build bool) (string, error) {
	p, _, _, err := buildPlan(cfg, a.page, a.platform)
	if err != nil {
		return "", err
	}
	out := cfg.OutDir(a.platform, a.page)
	w, h := int(p.Width), int(p.Height)
	if a.platform == "h5" {
		err = capture.H5(cfg.P(cfg.H5.Root), cfg.H5.Build, cfg.H5.Dist, a.page, out, w, h, build)
	} else {
		err = capture.Android(cfg.P(cfg.Android.Root), cfg.Android.Build, cfg.Android.Run, a.page, out, w, h, build)
	}
	return out, err
}

func runCapture(args []string) error {
	a, _ := pageFlags("capture", args)
	if err := checkArgs(a, "capture"); err != nil {
		return err
	}
	cfg, err := load(a.root)
	if err != nil {
		return err
	}
	out, err := doCapture(cfg, a, true)
	if err != nil {
		return err
	}
	fmt.Printf("✓ → %s/screen.png, nodes.json\n", relTo(cfg, out))
	return nil
}

func doDiff(cfg *config.Config, a *pageArgs) (*diff.Result, error) {
	p, _, exp, err := buildPlan(cfg, a.page, a.platform)
	if err != nil {
		return nil, err
	}
	out := cfg.OutDir(a.platform, a.page)
	ratio, sizeOK, mask, err := diff.Pixel(filepath.Join(exp.Dir, exp.Images.Frame), filepath.Join(out, "screen.png"), filepath.Join(out, "diff.png"))
	if err != nil {
		return nil, err
	}
	nodes, err := capture.LoadNodes(filepath.Join(out, "nodes.json"))
	if err != nil {
		return nil, err
	}
	res := diff.Nodes(exp, p, nodes, cfg.Thresholds.PositionTol, mask, func(c string) string { return c })
	res.PixelMismatch, res.SizeOK = ratio, sizeOK
	writeJSON(filepath.Join(out, "diff.json"), res)
	return res, nil
}

func runDiff(args []string) error {
	a, _ := pageFlags("diff", args)
	if err := checkArgs(a, "diff"); err != nil {
		return err
	}
	cfg, err := load(a.root)
	if err != nil {
		return err
	}
	res, err := doDiff(cfg, a)
	if err != nil {
		return err
	}
	fmt.Printf("pixel mismatch %.2f%%, node errors %d, warnings %d\n", res.PixelMismatch*100, res.Errors, res.Warns)
	return nil
}

// runVerify is the Skill's loop step: audit + capture + diff + report. Exit 3 when below threshold.
func runVerify(args []string) error {
	a, _ := pageFlags("verify", args)
	if err := checkArgs(a, "verify"); err != nil {
		return err
	}
	cfg, err := load(a.root)
	if err != nil {
		return err
	}
	au, err := doAudit(cfg, a)
	if err != nil {
		return err
	}
	var d *diff.Result
	if _, err := doCapture(cfg, a, true); err != nil {
		fmt.Fprintln(os.Stderr, "⚠ capture failed:", err)
	} else if d, err = doDiff(cfg, a); err != nil {
		return err
	}
	r := report.Build(a.page, a.platform, au, d, cfg.Thresholds)
	out := cfg.OutDir(a.platform, a.page)
	md := r.Markdown(40)
	_ = os.WriteFile(filepath.Join(out, "report.md"), []byte(md), 0o644)
	writeJSON(filepath.Join(out, "report.json"), r)
	if !a.quiet {
		fmt.Print(md)
	}
	fmt.Printf("\n→ %s/{report.md,report.json,screen.png,diff.png}\n", relTo(cfg, out))
	if d == nil {
		return exitError{4, "capture failed: page did not render"}
	}
	if !r.Pass {
		return exitError{3, "below threshold — fix the items above and run verify again"}
	}
	return nil
}
