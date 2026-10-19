---
name: figma-to-code
description: 按 Figma 设计稿实现或还原 H5 / Android 页面时使用（「照 Figma 做」「还原设计稿」「实现这个页面」「按设计稿改 UI」「这个设计稿哪些组件能复用」）。强制流程：先盘点现有组件资产，再根据新设计稿找出可直接复用的组件、出方案并等人确认，最后实现并验证，避免重复造轮子。支持常见栈：Vue / React / uni-app（Vant、React Vant、NutUI、antd-mobile、TDesign、uView、Element Plus）与 Android XML View / Compose（Material、QMUI、项目自有设计系统）。
---

# Figma → 代码（H5 / Android）

目标：**最大化复用现有组件，实现与 Figma 一致**。不要凭截图手写 UI。
所有确定性的工作都交给 CLI `figmafe`（仓库根目录 `bin/figmafe`，没有就先 `cd cli && go build -o ../bin/figmafe ./cmd/figmafe`）。

## 先判断走哪条流程

- 仓库根目录**有 `figmafe.yaml`**（已经接入：映射、组件手册、截图比对都配好了）→ 按下面的完整流程做。
- **没有**（普通的存量项目）→ 按 `references/brownfield.md` 做：`figmafe reuse` 体检 → `figmafe match` 拿到新设计稿的复用清单（每个组件：用什么、怎么配置；`.html` 是叠在设计图上的可视版，方案闸门时给人看）→ 实现 → `figmafe reuse --focus`（或 `--changed <ref>`）做页面复用检查 → 不采纳推荐 / 新建组件的决定写进 `component-map.yaml`。

## 硬性规则

1. **先盘点资产再动手**：资产清单（`.figmafe/inventory-*.json`）不存在，或比组件库源码旧时，先跑 Phase 0。
2. **写代码前必须有复用方案并得到确认**（Phase 2 闸门）。用户明确说「直接做、不用确认」时，把方案贴出来后继续。
3. **只能用 catalog 里的组件和 token**：Figma 组件 → 对应代码组件；页面级文本和容器只能用 token。不能用原生 `<button>`/`<input>`，不能直接用 `material3.*` 组件，不能用 `Legacy*` XML 老组件，不能写死颜色。
4. **catalog 里没有的组件**：按 component-map 的 `gaps` 处理方案执行；没有登记的，在方案里列为待决项问人，不要在页面里自造组件。
5. **交付前 `figmafe verify` 必须 PASS**；修正最多 3 轮，3 轮后仍不达标就如实汇报剩下的偏差和原因，不要硬凑。
6. 每个对应 Figma 组件实例的节点都挂上 Figma 节点 id（H5 `data-figma-id`，Android `Modifier.figmaId`），verify 才能逐节点比对。

## 流程

### Phase 0：资产盘点（组件库有变化时）
```
bin/figmafe scan                 # 两端：组件 / props / 变体 / 示例 / 使用次数 / token / 存量问题（含业务层重复组件）
bin/figmafe map check            # 缺口、变体不支持、token 漂移、重复资产
bin/figmafe catalog              # 更新 references/catalog-*.md
bin/figmafe match --figma design/figma-export/<page> .   # 新设计稿里还有映射没覆盖的组件？看能复用什么
```
新增了组件、映射还没有的：`bin/figmafe map suggest` 给出候选 → 人确认后写入 `design/component-map.yaml`。

### Phase 1：读手册
读 `references/catalog-<platform>.md`（组件手册：Figma 组件 → 代码组件、变体 → props、示例、约束、token）和 `references/rules.md`（沉淀下来的经验）。平台细节见 `references/workflow-h5.md` / `references/workflow-android.md`。

### Phase 2：复用方案（闸门）
```
bin/figmafe plan --page <Page> --platform h5|android
```
输出 `.figmafe/<platform>/<Page>/plan.md`：每个 Figma 节点 → 组件 → 建议代码（含 props、文本、figma id），以及待决项和约束。
把**组件清单 + 待决项 + 你的实现要点**贴给用户确认。方案以 plan.md 为准，不要自己改组件选择。

### Phase 3：实现
按 plan.md 写代码：组件按建议代码使用，页面级文本 / 容器按方案里的 token 写，布局间距用 token。

### Phase 4：验证闭环（≤3 轮）
```
bin/figmafe verify --page <Page> --platform h5|android     # 退出码 0=PASS 3=未达标 4=渲染失败
```
读 `.figmafe/<platform>/<Page>/report.md`：
- `audit/*` 是静态问题（用错组件、原生元素、写死颜色等），按提示改。
- `diff/*` 是节点级偏差（缺文本、位置 / 尺寸 / 字号 / 颜色不对），对照 figma id 改。
- 需要看图时读 `screen.png`（实现）、`diff.png`（红色为差异像素），设计图在 `design/figma-export/<page>/frame@2x.png`。
改完再跑 verify，直到 PASS。

### Phase 5：沉淀
本次遇到的、下次还会遇到的问题（组件用法坑、映射缺口、设计与组件不一致），写进 `references/rules.md`（一条一行：现象 → 规则 → 来源页面）。组件或映射本身的问题，告诉用户需要改 component-map 或组件库。

## 交付时回复
复用方案摘要、verify 最终指标（复用率 / 违规 / 像素差异 / 节点偏差）、剩下的偏差及原因、新增的规则。
