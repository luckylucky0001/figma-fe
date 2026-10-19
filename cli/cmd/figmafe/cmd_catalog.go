package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
	"github.com/luckylucky0001/figma-fe/cli/internal/inventory"
	"github.com/luckylucky0001/figma-fe/cli/internal/mapping"
	"github.com/luckylucky0001/figma-fe/cli/internal/plan"
)

// runCatalog writes the Skill's per-platform component catalog (the Agent's "component manual"):
// Figma component → code component, import, props & variant mapping, example, constraints, tokens.
func runCatalog(args []string) error {
	fs := flag.NewFlagSet("catalog", flag.ExitOnError)
	root := fs.String("root", "", "repo root")
	outDir := fs.String("o", ".claude/skills/figma-to-code/references", "output dir")
	parse(fs, args)
	cfg, err := load(*root)
	if err != nil {
		return err
	}
	lib, err := figma.LoadLibrary(cfg.P(cfg.Figma.Library))
	if err != nil {
		return err
	}
	m, err := mapping.Load(cfg.P(cfg.Map))
	if err != nil {
		return err
	}
	invs, err := loadInventories(cfg, true)
	if err != nil {
		return err
	}
	tok, err := plan.LoadTokens(cfg.P(cfg.Figma.Variables))
	if err != nil {
		return err
	}
	_ = os.MkdirAll(cfg.P(*outDir), 0o755)
	for _, pf := range []string{"h5", "android"} {
		md := catalogMarkdown(pf, lib, m, invs[pf], tok)
		p := cfg.P(*outDir, "catalog-"+pf+".md")
		if err := os.WriteFile(p, []byte(md), 0o644); err != nil {
			return err
		}
		fmt.Printf("✓ → %s\n", relTo(cfg, p))
	}
	return nil
}

func catalogMarkdown(pf string, lib *figma.Library, m *mapping.Map, inv *inventory.Inventory, tok *plan.Tokens) string {
	var b strings.Builder
	title := map[string]string{"h5": "H5（React）", "android": "Android（Compose）"}[pf]
	fmt.Fprintf(&b, "# 组件手册：%s\n\n> 由 `figmafe catalog` 根据资产扫描结果（inventory）和 component-map 自动生成，不要手改。\n> 规则：Figma 里出现的组件，**必须**用下表对应的代码组件；表里没有的才允许自己写，并在方案里列为待决项。\n\n", title)
	var sets []figma.ComponentSet
	for _, s := range lib.ComponentSets {
		sets = append(sets, s)
	}
	sort.Slice(sets, func(i, j int) bool { return sets[i].Name < sets[j].Name })
	b.WriteString("## 速查\n\n| Figma 组件 | 代码组件 | 业务代码中已用次数 |\n|---|---|---|\n")
	for _, s := range sets {
		e := m.Entry(s.Name)
		if e == nil || e.Impl(pf) == nil {
			d := "**无** —— 见待决项"
			fmt.Fprintf(&b, "| %s | %s | - |\n", s.Name, d)
			continue
		}
		im := e.Impl(pf)
		u := 0
		if c := inv.Find(im.Component); c != nil {
			u = c.Usage
		}
		fmt.Fprintf(&b, "| %s | `%s` | %d |\n", s.Name, im.Component, u)
	}
	b.WriteString("\n## 组件详情\n")
	for _, s := range sets {
		e := m.Entry(s.Name)
		if e == nil || e.Impl(pf) == nil {
			continue
		}
		im := e.Impl(pf)
		c := inv.Find(im.Component)
		fmt.Fprintf(&b, "\n### %s → `%s`\n\n", s.Name, im.Component)
		if s.Description != "" {
			fmt.Fprintf(&b, "%s。", s.Description)
		}
		if c != nil {
			fmt.Fprintf(&b, "源码：`%s`。\n\n", c.File)
			if c.Description != "" {
				fmt.Fprintf(&b, "%s\n\n", c.Description)
			}
		}
		if im.Import != "" {
			fmt.Fprintf(&b, "```\n%s\n```\n\n", im.Import)
		}
		if len(s.VariantProperties) > 0 {
			b.WriteString("| Figma 变体 | 代码 |\n|---|---|\n")
			axes := make([]string, 0, len(s.VariantProperties))
			for a := range s.VariantProperties {
				axes = append(axes, a)
			}
			sort.Strings(axes)
			for _, a := range axes {
				for _, v := range s.VariantProperties[a] {
					frag, ok := im.Props[a][v]
					switch {
					case !ok:
						frag = "（未映射）"
					case frag == "":
						frag = "默认"
					case frag == "-":
						frag = "**不支持**"
					default:
						frag = "`" + frag + "`"
					}
					fmt.Fprintf(&b, "| %s=%s | %s |\n", a, v, frag)
				}
			}
			b.WriteString("\n")
		}
		if len(im.Texts) > 0 {
			var parts []string
			for k, v := range im.Texts {
				parts = append(parts, fmt.Sprintf("`%s` ← 图层 %s", k, v))
			}
			sort.Strings(parts)
			fmt.Fprintf(&b, "文本：%s\n\n", strings.Join(parts, "，"))
		}
		if c != nil && len(c.Props) > 0 {
			var ps []string
			for _, p := range c.Props {
				s := p.Name + ": " + p.Type
				if p.Default != "" {
					s += " = " + p.Default
				}
				ps = append(ps, "`"+s+"`")
			}
			fmt.Fprintf(&b, "参数：%s\n\n", strings.Join(ps, "、"))
		}
		if im.Usage != "" {
			fmt.Fprintf(&b, "示例：`%s`\n\n", im.Usage)
		}
		for _, k := range e.Constraints {
			fmt.Fprintf(&b, "- 约束：%s\n", k)
		}
		if im.Notes != "" {
			fmt.Fprintf(&b, "- 注意：%s\n", im.Notes)
		}
	}
	b.WriteString("\n## 不对应 Figma 组件的辅助组件\n\n")
	for _, c := range inv.Components {
		mapped := false
		for _, e := range m.Components {
			if im := e.Impl(pf); im != nil && im.Component == c.Name {
				mapped = true
			}
		}
		if mapped {
			continue
		}
		note := ""
		if c.Kind == "view" {
			note = "（XML View 存量组件，新页面不要用）"
		}
		fmt.Fprintf(&b, "- `%s` %s%s — `%s`\n", c.Name, c.Description, note, c.File)
	}
	if len(m.Gaps) > 0 {
		b.WriteString("\n## 待决项（Figma 有、代码没有）\n\n")
		for _, g := range m.Gaps {
			fmt.Fprintf(&b, "- **%s**（%s）：%s\n", g.Figma, g.Platform, g.Decision)
		}
	}
	b.WriteString("\n## Token（页面级文本 / 容器只能用这些）\n\n| Figma 变量 | 代码 |\n|---|---|\n")
	var cols []string
	for hex, name := range tok.Colors {
		cols = append(cols, name+"\t"+hex)
	}
	sort.Strings(cols)
	for _, c := range cols {
		p := strings.SplitN(c, "\t", 2)
		fmt.Fprintf(&b, "| %s %s | `%s` |\n", p[0], p[1], plan.Ref(pf, p[0]))
	}
	var fonts []string
	for k, name := range tok.Fonts {
		fonts = append(fonts, name+"\t"+k)
	}
	sort.Strings(fonts)
	for _, f := range fonts {
		p := strings.SplitN(f, "\t", 2)
		ref := plan.Ref(pf, p[0])
		if pf == "h5" {
			n := strings.TrimPrefix(p[0], "font.")
			ref = fmt.Sprintf("var(--ds-font-%s-size) / var(--ds-font-%s-line-height) / var(--ds-font-%s-weight)", n, n, n)
		}
		fmt.Fprintf(&b, "| %s（%s）| `%s` |\n", p[0], p[1], ref)
	}
	for _, s := range []string{"xs 4", "sm 8", "md 12", "lg 16", "xl 24"} {
		k := strings.Fields(s)
		fmt.Fprintf(&b, "| space.%s（%s）| `%s` |\n", k[0], k[1], plan.Ref(pf, "space."+k[0]))
	}
	if pf == "android" {
		b.WriteString("\n页面级文本用 `DsText(text, style = DsTypography.X, color = DsColors.Y)`；挂 Figma id 用 `Modifier.figmaId(\"1:23\")`。\n")
	} else {
		b.WriteString("\n页面级文本用普通元素 + CSS 变量；组件和关键节点上挂 `data-figma-id=\"1:23\"`（组件会透传到根节点）。\n")
	}
	return b.String()
}

var _ = filepath.Join
