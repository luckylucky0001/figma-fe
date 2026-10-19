package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/luckylucky0001/figma-fe/cli/internal/config"
	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
	"github.com/luckylucky0001/figma-fe/cli/internal/inventory"
	"github.com/luckylucky0001/figma-fe/cli/internal/mapping"
)

func inventoryPath(cfg *config.Config, pf string) string {
	return cfg.P(cfg.Out, "inventory-"+pf+".json")
}

func scanPlatform(cfg *config.Config, pf string) (*inventory.Inventory, error) {
	switch pf {
	case "h5":
		return inventory.ScanH5(cfg.Root, cfg.P(cfg.H5.Root), cfg.H5.DS)
	case "android":
		return inventory.ScanAndroid(cfg.Root, cfg.P(cfg.Android.Root), cfg.Android.DS, cfg.Android.DSView, cfg.Android.App)
	}
	return nil, fmt.Errorf("unknown platform %q (h5|android)", pf)
}

func platforms(p string) []string {
	if p == "" || p == "all" {
		return []string{"h5", "android"}
	}
	return []string{p}
}

func runScan(args []string) error {
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	pf := fs.String("platform", "all", "h5|android|all")
	root := fs.String("root", "", "repo root")
	parse(fs, args)
	cfg, err := load(*root)
	if err != nil {
		return err
	}
	for _, p := range platforms(*pf) {
		inv, err := scanPlatform(cfg, p)
		if err != nil {
			return err
		}
		if err := inv.Save(inventoryPath(cfg, p)); err != nil {
			return err
		}
		used := 0
		for _, c := range inv.Components {
			if c.Usage > 0 {
				used++
			}
		}
		fmt.Printf("✓ [%s] %d components (%d used in business code), %d tokens, %d baseline findings → %s\n",
			p, len(inv.Components), used, len(inv.Tokens), len(inv.Baseline), relTo(cfg, inventoryPath(cfg, p)))
		for _, c := range inv.Components {
			var axes []string
			for _, pr := range c.Props {
				if len(pr.Values) > 0 {
					axes = append(axes, pr.Name+"="+strings.Join(pr.Values, "|"))
				}
			}
			fmt.Printf("    %-16s %-8s usage=%-3d %s\n", c.Name, c.Kind, c.Usage, strings.Join(axes, "  "))
		}
	}
	return nil
}

func loadInventories(cfg *config.Config, fresh bool) (map[string]*inventory.Inventory, error) {
	out := map[string]*inventory.Inventory{}
	for _, p := range []string{"h5", "android"} {
		var inv *inventory.Inventory
		var err error
		if !fresh {
			inv, err = inventory.Load(inventoryPath(cfg, p))
		}
		if fresh || err != nil {
			inv, err = scanPlatform(cfg, p)
			if err != nil {
				return nil, err
			}
			_ = inv.Save(inventoryPath(cfg, p))
		}
		out[p] = inv
	}
	return out, nil
}

func relTo(cfg *config.Config, p string) string {
	r, err := filepath.Rel(cfg.Root, p)
	if err != nil {
		return p
	}
	return r
}

func runMap(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: figmafe map suggest|check")
	}
	sub := args[0]
	fs := flag.NewFlagSet("map", flag.ExitOnError)
	root := fs.String("root", "", "repo root")
	parse(fs, args[1:])
	cfg, err := load(*root)
	if err != nil {
		return err
	}
	lib, err := figma.LoadLibrary(cfg.P(cfg.Figma.Library))
	if err != nil {
		return err
	}
	invs, err := loadInventories(cfg, true)
	if err != nil {
		return err
	}
	switch sub {
	case "suggest":
		sugg := mapping.Suggest(lib, invs)
		out := cfg.P(cfg.Out, "component-map.suggested.yaml")
		b, _ := yaml.Marshal(map[string]any{"generatedBy": "figmafe map suggest", "note": "候选映射，需人工确认后写入 design/component-map.yaml", "suggestions": sugg})
		if err := os.WriteFile(out, b, 0o644); err != nil {
			return err
		}
		fmt.Printf("%-12s │ %-26s │ %s\n", "Figma", "H5", "Android")
		fmt.Println(strings.Repeat("─", 13) + "┼" + strings.Repeat("─", 28) + "┼" + strings.Repeat("─", 44))
		for _, s := range sugg {
			cell := func(pf string) string {
				cs := s.Options[pf]
				if len(cs) == 0 {
					return "✗ 无候选"
				}
				mark := map[string]string{"auto": "✓", "review": "?"}[s.Status[pf]]
				out := fmt.Sprintf("%s %s %.2f", mark, cs[0].Component, cs[0].Score)
				if len(cs) > 1 && cs[1].Score >= 0.85 {
					out += fmt.Sprintf(" / %s %.2f", cs[1].Component, cs[1].Score) // several assets compete: human decides
				}
				return out
			}
			fmt.Printf("%-12s │ %-26s │ %s\n", s.Figma, cell("h5"), cell("android"))
		}
		fmt.Printf("\n✓ → %s  (✓ 高置信  ? 需人工确认  ✗ 缺口)\n", relTo(cfg, out))
	case "check":
		m, err := mapping.Load(cfg.P(cfg.Map))
		if err != nil {
			return err
		}
		res, err := mapping.Check(m, lib, invs, cfg.P(cfg.Figma.Variables))
		if err != nil {
			return err
		}
		md := renderCheck(res)
		out := cfg.P(cfg.Out, "map-check.md")
		_ = os.WriteFile(out, []byte(md), 0o644)
		writeJSON(cfg.P(cfg.Out, "map-check.json"), res)
		fmt.Print(md)
		fmt.Printf("\n✓ → %s\n", relTo(cfg, out))
	default:
		return fmt.Errorf("unknown map subcommand %q", sub)
	}
	return nil
}

func renderCheck(res *mapping.CheckResult) string {
	var b strings.Builder
	b.WriteString("# 映射检查报告（缺口 + token 漂移）\n\n")
	for _, pf := range []string{"h5", "android"} {
		c := res.Coverage[pf]
		fmt.Fprintf(&b, "- **%s**：Figma 组件覆盖 %d / %d\n", pf, c.Mapped, c.Total)
	}
	groups := map[string][]mapping.Issue{}
	for _, i := range res.Issues {
		groups[i.Kind] = append(groups[i.Kind], i)
	}
	titles := []struct{ k, t string }{
		{"missing-impl", "Figma 有、代码没有"},
		{"unknown-component", "映射指向不存在的组件"},
		{"variant-unsupported", "变体在该端不支持"},
		{"variant-unmapped", "变体没有 props 映射"},
		{"token-drift", "Token 漂移（Figma Variables ≠ 代码）"},
		{"token-missing", "Token 缺失"},
		{"code-only", "代码有、Figma 没有"},
	}
	for _, t := range titles {
		is := groups[t.k]
		if len(is) == 0 {
			continue
		}
		sort.Slice(is, func(a, c int) bool { return is[a].Platform+is[a].Figma < is[c].Platform+is[c].Figma })
		fmt.Fprintf(&b, "\n## %s（%d）\n\n", t.t, len(is))
		for _, i := range is {
			name := i.Figma
			if name != "" {
				name = "`" + name + "` "
			}
			fmt.Fprintf(&b, "- [%s][%s] %s%s\n", i.Severity, i.Platform, name, i.Detail)
		}
	}
	return b.String()
}
