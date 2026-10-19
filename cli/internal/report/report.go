// Package report merges plan, audit and diff results into metrics, a pass/fail verdict and an
// Agent-readable list of fixes.
package report

import (
	"fmt"
	"strings"

	"github.com/luckylucky0001/figma-fe/cli/internal/audit"
	"github.com/luckylucky0001/figma-fe/cli/internal/config"
	"github.com/luckylucky0001/figma-fe/cli/internal/diff"
)

type Metrics struct {
	ReuseRate     float64 `json:"reuseRate"`     // static: mapped components used vs plan
	RuntimeReuse  float64 `json:"runtimeReuse"`  // rendered by the mapped DS component
	Violations    int     `json:"violations"`    // audit errors
	Warnings      int     `json:"warnings"`      // audit warnings
	PixelMismatch float64 `json:"pixelMismatch"` // ratio
	NodeErrors    int     `json:"nodeErrors"`    // diff errors
	NodeWarnings  int     `json:"nodeWarnings"`  // diff warnings
	TextCoverage  float64 `json:"textCoverage"`  // Figma texts found on page
	FigmaIDs      int     `json:"figmaIds"`
}

type Report struct {
	Page     string        `json:"page"`
	Platform string        `json:"platform"`
	Pass     bool          `json:"pass"`
	Failed   []string      `json:"failed"`
	Metrics  Metrics       `json:"metrics"`
	Audit    *audit.Result `json:"audit"`
	Diff     *diff.Result  `json:"diff"`
}

func Build(page, platform string, a *audit.Result, d *diff.Result, th config.Thresholds) *Report {
	r := &Report{Page: page, Platform: platform, Audit: a, Diff: d}
	m := &r.Metrics
	m.ReuseRate = a.ReuseRate
	m.Violations = a.Violations
	for _, f := range a.Findings {
		if f.Severity != "error" {
			m.Warnings++
		}
	}
	m.FigmaIDs = a.FigmaIDs
	if d != nil {
		if d.Instances > 0 {
			m.RuntimeReuse = float64(d.InstancesByDS) / float64(d.Instances)
		}
		m.PixelMismatch = d.PixelMismatch
		m.NodeErrors = d.Errors
		m.NodeWarnings = d.Warns
		if d.Texts > 0 {
			m.TextCoverage = float64(d.TextsMatched) / float64(d.Texts)
		}
	}
	pm := th.PixelMismatch[platform]
	if m.ReuseRate < th.Reuse {
		r.Failed = append(r.Failed, fmt.Sprintf("组件复用率 %.0f%% < %.0f%%", m.ReuseRate*100, th.Reuse*100))
	}
	if m.Violations > th.Violations {
		r.Failed = append(r.Failed, fmt.Sprintf("违规 %d 处 > %d", m.Violations, th.Violations))
	}
	if d != nil && m.PixelMismatch > pm {
		r.Failed = append(r.Failed, fmt.Sprintf("像素差异 %.2f%% > %.2f%%", m.PixelMismatch*100, pm*100))
	}
	if d != nil && m.NodeErrors > th.NodeErrors {
		r.Failed = append(r.Failed, fmt.Sprintf("节点级严重偏差 %d 处 > %d", m.NodeErrors, th.NodeErrors))
	}
	r.Pass = len(r.Failed) == 0
	return r
}

func (r *Report) Markdown(maxItems int) string {
	var b strings.Builder
	m := r.Metrics
	verdict := "✅ PASS"
	if !r.Pass {
		verdict = "❌ FAIL：" + strings.Join(r.Failed, "；")
	}
	fmt.Fprintf(&b, "# 验证报告：%s（%s）\n\n**%s**\n\n", r.Page, r.Platform, verdict)
	b.WriteString("| 指标 | 值 |\n|---|---|\n")
	fmt.Fprintf(&b, "| 组件复用率（静态，按方案） | %.0f%% |\n", m.ReuseRate*100)
	if r.Diff != nil {
		fmt.Fprintf(&b, "| 组件复用率（运行时，由组件库渲染） | %.0f%% |\n", m.RuntimeReuse*100)
	}
	fmt.Fprintf(&b, "| 违规（原生元素 / 直接用 Material / 硬编码颜色 / 老组件） | %d |\n", m.Violations)
	fmt.Fprintf(&b, "| 警告（疑似自造组件 / 裸数值 / 少用组件） | %d |\n", m.Warnings)
	if r.Diff != nil {
		fmt.Fprintf(&b, "| 像素差异 | %.2f%% |\n", m.PixelMismatch*100)
		fmt.Fprintf(&b, "| 节点级偏差（严重 / 轻微） | %d / %d |\n", m.NodeErrors, m.NodeWarnings)
		fmt.Fprintf(&b, "| Figma 文本覆盖 | %.0f%% |\n", m.TextCoverage*100)
	}
	fmt.Fprintf(&b, "| 挂了 Figma id 的节点 | %d |\n", m.FigmaIDs)

	b.WriteString("\n## 组件使用（方案 vs 实际）\n\n| 组件 | 方案 | 实际 |\n|---|---|---|\n")
	for _, u := range r.Audit.Usage {
		mark := ""
		if u.Actual < u.Expected {
			mark = " ⚠"
		}
		fmt.Fprintf(&b, "| `%s` | %d | %d%s |\n", u.Component, u.Expected, u.Actual, mark)
	}
	b.WriteString("\n## 需要修正（按严重程度）\n\n")
	n := 0
	for _, f := range r.Audit.Findings {
		if f.Severity != "error" {
			continue
		}
		fmt.Fprintf(&b, "- ✗ [audit/%s] %s:%d `%s` → %s\n", f.Rule, f.File, f.Line, trunc(f.Text, 80), f.Hint)
		n++
	}
	if r.Diff != nil {
		for _, d := range r.Diff.Deviations {
			if n >= maxItems {
				fmt.Fprintf(&b, "- … 还有 %d 条，见 report.json\n", len(r.Diff.Deviations)-n)
				break
			}
			mark := "✗"
			if d.Severity != "error" {
				mark = "·"
			}
			fmt.Fprintf(&b, "- %s [diff/%s] %s %s：期望 %s，实际 %s\n", mark, d.Kind, d.FigmaID, d.Target, d.Expected, d.Actual)
			n++
		}
	}
	for _, f := range r.Audit.Findings {
		if f.Severity == "error" || n >= maxItems+10 {
			continue
		}
		loc := ""
		if f.File != "" {
			loc = fmt.Sprintf("%s:%d ", f.File, f.Line)
		}
		fmt.Fprintf(&b, "- · [audit/%s] %s%s\n", f.Rule, loc, f.Hint)
		n++
	}
	if r.Diff != nil && len(r.Diff.Regions) > 0 {
		b.WriteString("\n## 像素差异最集中的区域（见 diff.png）\n\n")
		for _, rg := range r.Diff.Regions {
			fmt.Fprintf(&b, "- %s %s：%.0f%%\n", rg.FigmaID, rg.Target, rg.Mismatch*100)
		}
	}
	return b.String()
}

func trunc(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n]) + "…"
	}
	return s
}
