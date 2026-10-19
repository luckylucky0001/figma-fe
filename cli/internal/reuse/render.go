package reuse

import (
	"fmt"
	"sort"
	"strings"
)

// Markdown renders the brownfield reuse check report.
func (r *Report) Markdown() string {
	var b strings.Builder
	m := r.Metrics
	fmt.Fprintf(&b, "# 存量复用体检：%s（%s）\n\n", r.Root, r.Platform)
	fmt.Fprintf(&b, "技术栈：%s · 扫描文件 %d\n\n", strings.Join(r.Stack, " / "), m.Files)
	b.WriteString("| 指标 | 值 |\n|---|---|\n")
	fmt.Fprintf(&b, "| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | %.0f%%（%d / %d） |\n", m.Componentization*100, m.ComponentUsages, m.ComponentUsages+m.RawPrimitives)
	fmt.Fprintf(&b, "| 重复概念（同一 UI 概念有多个实现） | %d |\n", m.DuplicateConcepts)
	fmt.Fprintf(&b, "| 重复结构（可抽成组件的复制粘贴结构） | %d 组 |\n", m.CloneGroups)
	fmt.Fprintf(&b, "| 硬编码颜色 | %d 处（其中 %d 处已有同值 token） |\n", m.HardcodedColors, m.TokenizableColors)
	if m.HardcodedSizes > 0 {
		fmt.Fprintf(&b, "| 硬编码尺寸 | %d 处 |\n", m.HardcodedSizes)
	}

	if len(r.Libraries) > 0 {
		b.WriteString("\n## 1. 三方组件库使用\n\n| 组件库 | 使用次数 | 最常用 |\n|---|---|---|\n")
		for _, l := range r.Libraries {
			fmt.Fprintf(&b, "| %s | %d | %s |\n", l.Name, l.Usage, strings.Join(l.Top, "、"))
		}
		var bad []string
		for i := range r.LibAll {
			if a := &r.LibAll[i]; !usable(a) {
				bad = append(bad, "- `"+a.Name+"`："+a.Caveat)
			}
		}
		if len(bad) > 0 {
			sort.Strings(bad)
			b.WriteString("\n**组件库里有、但本项目不能直接用的组件**（match 不会推荐它们）：\n\n" + strings.Join(bad, "\n") + "\n")
		}
	}

	if len(r.Conventions) > 0 {
		b.WriteString("\n**项目惯例写法**（没有组件、但项目里固定这么写；新设计稿里没有对应组件时沿用它，不要另造）：\n\n")
		for _, c := range r.Conventions {
			fmt.Fprintf(&b, "- %s：%s，%d 处，例 `%s` %s\n", c.Concept, c.Name, c.Usage, strings.ReplaceAll(c.Example, "`", "'"), c.ExampleAt)
		}
	}

	b.WriteString("\n## 2. 组件资产（可复用的现成组件）\n\n| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |\n|---|---|---|---|---|\n")
	n := 0
	for _, a := range r.Assets {
		if a.Kind == "page-local" && a.Usage < 2 {
			continue
		}
		if n >= 60 {
			fmt.Fprintf(&b, "| … | | | | 共 %d 个，见 reuse.json |\n", len(r.Assets))
			break
		}
		ex := a.Example
		if ex != "" {
			ex = "`" + strings.ReplaceAll(ex, "|", "\\|") + "` " + a.ExampleAt
		} else {
			ex = a.Source
		}
		fmt.Fprintf(&b, "| `%s` | %s | %s | %d | %s |\n", a.Name, kindLabel(a.Kind), a.Concept, a.Usage, ex)
		n++
	}

	if len(r.Clusters) > 0 {
		b.WriteString("\n## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）\n\n| 概念 | 实现 | 推荐统一为 |\n|---|---|---|\n")
		for _, c := range r.Clusters {
			fmt.Fprintf(&b, "| %s | %s | `%s` |\n", c.Concept, strings.Join(c.Members, "<br>"), c.Recommend)
		}
	}

	var missed []Finding
	for _, f := range r.Findings {
		if f.Rule == "missed-reuse" {
			missed = append(missed, f)
		}
	}
	if len(missed) > 0 {
		fmt.Fprintf(&b, "\n## 4. 有组件却没用（%d 处）\n\n", len(missed))
		for i, f := range missed {
			if i >= 25 {
				fmt.Fprintf(&b, "- … 还有 %d 处\n", len(missed)-i)
				break
			}
			fmt.Fprintf(&b, "- %s:%d `%s` → %s\n", f.File, f.Line, strings.ReplaceAll(f.Text, "`", "'"), f.Hint)
		}
	}

	if len(r.Clones) > 0 {
		fmt.Fprintf(&b, "\n## 5. 重复结构（复制粘贴的 UI 片段，可抽成组件）\n\n")
		for i, c := range r.Clones {
			if i >= 10 {
				fmt.Fprintf(&b, "- … 共 %d 组\n", len(r.Clones))
				break
			}
			fmt.Fprintf(&b, "- **%d 个元素 × %d 处**：%s\n", c.Elements, len(c.Locations), strings.Join(c.Locations, "、"))
			if i < 3 {
				fmt.Fprintf(&b, "\n```\n%s```\n", c.Sample)
			}
		}
	}

	if len(r.Colors) > 0 {
		b.WriteString("\n## 6. 硬编码颜色热点\n\n| 颜色 | 次数 | 已有 token | 出现在 |\n|---|---|---|---|\n")
		for i, h := range r.Colors {
			if i >= 12 {
				break
			}
			fmt.Fprintf(&b, "| %s | %d | %s | %s |\n", h.Value, h.Count, orDash(h.Token), strings.Join(h.Sample, "、"))
		}
	}
	return b.String()
}

func orDash(s string) string {
	if s == "" {
		return "—"
	}
	return "`" + s + "`"
}
