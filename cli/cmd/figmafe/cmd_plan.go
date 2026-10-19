package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/luckylucky0001/figma-fe/cli/internal/config"
	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
	"github.com/luckylucky0001/figma-fe/cli/internal/mapping"
	"github.com/luckylucky0001/figma-fe/cli/internal/plan"
)

type pageArgs struct {
	page, platform, root string
	quiet                bool
}

func pageFlags(name string, args []string) (*pageArgs, *flag.FlagSet) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	a := &pageArgs{}
	fs.StringVar(&a.page, "page", "", "Figma page name (figmafe.yaml figma.pages)")
	fs.StringVar(&a.platform, "platform", "", "h5|android")
	fs.StringVar(&a.root, "root", "", "repo root")
	fs.BoolVar(&a.quiet, "q", false, "quiet")
	parse(fs, args)
	return a, fs
}

func buildPlan(cfg *config.Config, page, platform string) (*plan.Plan, *mapping.Map, *figma.Export, error) {
	dir, err := cfg.PageDir(page)
	if err != nil {
		return nil, nil, nil, err
	}
	exp, err := figma.LoadExport(dir)
	if err != nil {
		return nil, nil, nil, err
	}
	m, err := mapping.Load(cfg.P(cfg.Map))
	if err != nil {
		return nil, nil, nil, err
	}
	tok, err := plan.LoadTokens(cfg.P(cfg.Figma.Variables))
	if err != nil {
		return nil, nil, nil, err
	}
	return plan.Build(exp, m, tok, page, platform), m, exp, nil
}

func runPlan(args []string) error {
	a, _ := pageFlags("plan", args)
	if a.page == "" || (a.platform != "h5" && a.platform != "android") {
		return fmt.Errorf("usage: figmafe plan --page <Page> --platform h5|android")
	}
	cfg, err := load(a.root)
	if err != nil {
		return err
	}
	p, m, _, err := buildPlan(cfg, a.page, a.platform)
	if err != nil {
		return err
	}
	out := cfg.OutDir(a.platform, a.page)
	md := p.Markdown(m)
	if err := os.WriteFile(filepath.Join(out, "plan.md"), []byte(md), 0o644); err != nil {
		return err
	}
	writeJSON(filepath.Join(out, "plan.json"), p)
	if !a.quiet {
		fmt.Print(md)
	}
	fmt.Printf("\n✓ → %s\n", relTo(cfg, filepath.Join(out, "plan.md")))
	return nil
}
