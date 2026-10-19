package mapping

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/luckylucky0001/figma-fe/cli/internal/figma"
	"github.com/luckylucky0001/figma-fe/cli/internal/inventory"
)

type Issue struct {
	Kind     string `json:"kind"` // missing-impl | unknown-component | variant-unsupported | variant-unmapped | code-only | token-drift | token-missing
	Severity string `json:"severity"`
	Figma    string `json:"figma,omitempty"`
	Platform string `json:"platform,omitempty"`
	Detail   string `json:"detail"`
}

type CheckResult struct {
	Issues   []Issue `json:"issues"`
	Coverage map[string]struct {
		Mapped int `json:"mapped"`
		Total  int `json:"total"`
	} `json:"coverage"`
}

// Check validates the confirmed map against the Figma library, inventories and Figma variables.
func Check(m *Map, lib *figma.Library, invs map[string]*inventory.Inventory, variablesPath string) (*CheckResult, error) {
	res := &CheckResult{Coverage: map[string]struct {
		Mapped int `json:"mapped"`
		Total  int `json:"total"`
	}{}}
	add := func(i Issue) { res.Issues = append(res.Issues, i) }
	var sets []figma.ComponentSet
	for _, s := range lib.ComponentSets {
		sets = append(sets, s)
	}
	sort.Slice(sets, func(i, j int) bool { return sets[i].Name < sets[j].Name })
	used := map[string]map[string]bool{"h5": {}, "android": {}}
	for _, pf := range []string{"h5", "android"} {
		cov := res.Coverage[pf]
		cov.Total = len(sets)
		for _, set := range sets {
			e := m.Entry(set.Name)
			var im *Impl
			if e != nil {
				im = e.Impl(pf)
			}
			if im == nil {
				g := m.GapFor(set.Name, pf)
				d := "Figma 组件在该端没有实现"
				if g != nil {
					d += "（已登记处理方案：" + g.Decision + "）"
				}
				add(Issue{Kind: "missing-impl", Severity: sev(g == nil), Figma: set.Name, Platform: pf, Detail: d})
				continue
			}
			cov.Mapped++
			used[pf][im.Component] = true
			inv := invs[pf]
			if inv != nil && inv.Find(im.Component) == nil {
				add(Issue{Kind: "unknown-component", Severity: "error", Figma: set.Name, Platform: pf, Detail: fmt.Sprintf("映射的组件 %s 不在资产清单里（改名或删除了？）", im.Component)})
			}
			for axis, vals := range set.VariantProperties {
				for _, v := range vals {
					frag, ok := im.Props[axis][v]
					switch {
					case !ok:
						add(Issue{Kind: "variant-unmapped", Severity: "warn", Figma: set.Name, Platform: pf, Detail: fmt.Sprintf("变体 %s=%s 没有写 props 映射", axis, v)})
					case frag == "-":
						add(Issue{Kind: "variant-unsupported", Severity: "warn", Figma: set.Name, Platform: pf, Detail: fmt.Sprintf("变体 %s=%s 在该端组件中不支持", axis, v)})
					}
				}
			}
		}
		res.Coverage[pf] = cov
		if inv := invs[pf]; inv != nil {
			for _, c := range inv.Components {
				if !used[pf][c.Name] {
					add(Issue{Kind: "code-only", Severity: "info", Platform: pf, Detail: fmt.Sprintf("%s（%s）没有对应的 Figma 组件：设计侧缺口、老组件或待废弃", c.Name, c.Kind)})
				}
			}
		}
	}
	if variablesPath != "" {
		vars, err := loadDTCGColors(variablesPath)
		if err != nil {
			return nil, err
		}
		names := make([]string, 0, len(vars))
		for k := range vars {
			names = append(names, k)
		}
		sort.Strings(names)
		for _, pf := range []string{"h5", "android"} {
			inv := invs[pf]
			if inv == nil {
				continue
			}
			byName := map[string][]inventory.Token{}
			for _, t := range inv.Tokens {
				byName[t.Name] = append(byName[t.Name], t)
			}
			for _, n := range names {
				toks := byName[n]
				if len(toks) == 0 {
					add(Issue{Kind: "token-missing", Severity: "warn", Platform: pf, Detail: fmt.Sprintf("Figma 变量 %s 在该端没有对应 token", n)})
					continue
				}
				for _, t := range toks {
					if !strings.EqualFold(strings.TrimSpace(t.Value), vars[n]) {
						add(Issue{Kind: "token-drift", Severity: "error", Platform: pf, Detail: fmt.Sprintf("%s：Figma %s ≠ 代码 %s（%s）", n, vars[n], strings.ToUpper(t.Value), t.Source)})
					}
				}
			}
		}
	}
	return res, nil
}

func sev(err bool) string {
	if err {
		return "error"
	}
	return "info"
}

// loadDTCGColors flattens the color group of a DTCG file into name → #RRGGBB.
func loadDTCGColors(path string) (map[string]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var root map[string]any
	if err := json.Unmarshal(b, &root); err != nil {
		return nil, err
	}
	out := map[string]string{}
	var walk func(prefix string, v map[string]any)
	walk = func(prefix string, v map[string]any) {
		for k, x := range v {
			if strings.HasPrefix(k, "$") {
				continue
			}
			obj, ok := x.(map[string]any)
			if !ok {
				continue
			}
			name := k
			if prefix != "" {
				name = prefix + "." + k
			}
			if val, ok := obj["$value"].(string); ok && strings.HasPrefix(val, "#") {
				out[name] = strings.ToUpper(val)
				continue
			}
			walk(name, obj)
		}
	}
	if c, ok := root["color"].(map[string]any); ok {
		walk("color", c)
	}
	return out, nil
}
