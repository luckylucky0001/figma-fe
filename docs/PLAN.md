# Plan：Figma → 代码的组件复用与视觉还原（痛点 4 验证 Demo）

## 1. 目标

用一个可复现的 Demo 证明下面这套做法有效，并且能直接推广：

> 梳理 Figma ↔ 代码组件映射并放进 Agent 的 Context → 开发前先识别要复用的组件、确认方案 → 开发后检查实际用了哪些组件，并拿截图和 Figma 对比 → 把偏差交给 Agent 修正 → 沉淀成规则。

交付物有三个：

| 交付物 | 说明 |
|---|---|
| **Skill** `figma-to-code` | 规定 Agent 的工作流程：先 Plan，等确认，再实现、验证、修正、沉淀 |
| **CLI** `figmafe`（Go） | 给 Skill 提供可执行、可量化的工具：sync / plan / audit / diff / report |
| **Demo** | 组件库、一个典型页面、A/B/C 对照实验、结论报告 |

说服力靠**数据对照**，不靠展示一个效果好的样例。

---

## 2. 实验设计

### 2.1 分组（消融实验，确认每一部分各自贡献了多少）

| 组 | Agent 能拿到的东西 | 验证的假设 |
|---|---|---|
| **A 基线** | 任务描述、Figma 导出图、节点 JSON、组件库源码（现状：组件在仓库里，但没人告诉它） | 现状复现 |
| **B 加 Context** | A 的全部，再加 component-map 和 catalog（组件位置、示例、约束） | 只补 Context 能提升多少 |
| **C 完整 Skill** | B 的全部，再加 Plan 闸门和 audit/diff 修正闭环（最多 3 轮） | 完整流程的效果 |

- 每组跑 3 次，用同一个模型、同一个任务描述，各自在全新的 worktree 里跑，互相隔离
- 三组拿到的设计信息完全一样，唯一的变量是 Context 和流程
- 实验里的 Plan 确认环节：把 Plan 落盘后自动批准，同时记录 Plan 本身的准确率（Plan 识别出的组件和 Figma 实际实例的吻合度）

### 2.2 指标（全部由 CLI 自动计算，以 Figma 为准，不需要参考实现）

| 指标 | 计算方式 | 达标线 |
|---|---|---|
| 组件复用率 | 被正确实现成映射组件的 Figma INSTANCE 数 ÷ 已映射的 INSTANCE 总数 | ≥ 95% |
| 违规数 | 原生 `<button>`/`<input>`/`<select>`/`<table>`、硬编码 hex/rgb、不在 token 内的 px | 0 |
| 像素偏差 | 截图与 Figma PNG 不一致的像素占比（带抗锯齿容差） | ≤ 3% |
| 节点级偏差 | 按 `data-figma-id` 逐个节点比较尺寸、padding、gap、圆角、字号、字重、颜色，超出容差的条目数 | 0 条严重偏差 |
| Plan 准确率（仅 C 组） | Plan 列出的组件与 Figma 实例的 F1 | ≥ 0.95 |
| 修正轮数（仅 C 组） | 达标前一共跑了几轮 audit/diff | 记录即可 |

### 2.3 预期结论的形式

```
            复用率   违规   像素偏差   节点偏差
A 基线       ~50%     10+     10%+      20+
B +Context   ~85%     3~5     5~8%      8~12
C 完整Skill  ≥95%      0      ≤3%       0
```

（以上是预期量级，最终以实测为准，如实写进报告。）

---

## 3. 仓库结构

```
figma-fe/
├── CLAUDE.md                        # 只写项目基本信息，不包含 Skill 内容（保证 A 组干净）
├── apps/demo/                       # Vite + React + TS
│   ├── src/tokens.css               # 颜色 / 间距 / 圆角 / 字号 / 阴影
│   ├── src/components/              # 组件库（见 4.1）
│   ├── src/pages/Orders/            # 目标页面（各组产出放在 experiments/ 里，跑分时挂载）
│   └── src/mock/orders.ts           # 页面数据，文案和 Figma 完全一致
├── design/
│   ├── figma/orders.nodes.json      # Figma REST API /v1/files/:key/nodes 的原始格式
│   ├── figma/orders@2x.png          # Figma /v1/images 导出（scale=2）
│   └── component-map.yaml           # ★ Figma 组件 ↔ 代码组件 映射
├── cli/                             # Go：figmafe
│   ├── cmd/figmafe/
│   └── internal/{figma,mapping,plan,audit,shot,diff,report}/
├── .claude/skills/figma-to-code/
│   ├── SKILL.md
│   └── references/
│       ├── catalog.md               # 由 `figmafe catalog` 生成
│       └── rules.md                 # 沉淀下来的规则（会持续更新）
├── experiments/
│   ├── prompts/{A,B,C}.md
│   └── runs/{A1..C3}/{Orders.tsx, plan.md, shot.png, diff.png, report.json}
└── docs/
    ├── PLAN.md
    └── report.md                    # 最终报告（另外发布成可分享的页面）
```

---

## 4. 各部分详细设计

### 4.1 组件库与典型页面

**页面：订单管理（Orders）**，1440 × 1024，基本覆盖中后台的常见元素：

- Sidebar（导航，带选中态）
- PageHeader（标题、面包屑、主操作 Button）
- StatCard × 4（指标数值、环比 Badge）
- FilterBar（Input 搜索框、Select × 2、Button 次级按钮）
- Tabs（全部 / 待发货 / 已发货 / 已完成）
- Table（订单行，状态列用 Badge，客户列用 Avatar）
- Pagination

**组件**：Button、Input、Select、Badge、Avatar、Card、StatCard、Tabs、Table、Pagination、PageHeader、Sidebar。

- 所有组件只用 token，并且把 `data-*` 属性透传到根节点（让 `data-figma-id` 能挂上去）
- 故意设置几个「容易用错」的点，模拟真实项目，例如 Badge 有 `tone` 语义（success/warning/…），不能直接传颜色；StatCard 内置了 padding 24，Agent 自己用 div 拼就容易写成 16

### 4.2 Figma 设计数据

- **首选**：你提供真实 Figma 文件（file key 和 `FIGMA_TOKEN`），通过 `figmafe sync` 拉取
- **兜底**：按 Figma REST API 的真实结构手工构造 `nodes.json`（FRAME / INSTANCE / TEXT，含 `componentId`、`componentProperties`、`absoluteBoundingBox`、`padding*`、`itemSpacing`、`cornerRadius`、`fills`、`style`），连同导出图一起入库。CLI 读这份数据和读 API 的代码路径完全一样
- 为了避免「自己出题、自己答题」：**先定稿设计数据，再写组件库**，实验过程中设计数据不再改动

### 4.3 component-map.yaml（核心 Context）

```yaml
version: 1
components:
  - figma: { name: "Button", key: "btn_key" }
    code:  { import: "@/components/Button", export: "Button" }
    variants:                         # Figma variant → props
      "Type=Primary":   'variant="primary"'
      "Type=Secondary": 'variant="secondary"'
      "Size=Large":     'size="lg"'
      "State=Disabled": "disabled"
    usage: '<Button variant="primary" size="md" icon={<Plus/>}>新建订单</Button>'
    constraints:
      - 一个页面最多 1 个 primary 按钮
      - 禁止用原生 <button> 加 className 自己拼
  - figma: { name: "Badge" }
    code:  { import: "@/components/Badge", export: "Badge" }
    variants: { "Status=Success": 'tone="success"', ... }
    constraints:
      - 状态色只能通过 tone 表达，禁止传 color/style
ignore:                               # 纯装饰节点，不需要映射
  - "Divider/*"
```

### 4.4 CLI `figmafe`（Go）

| 命令 | 输入 → 输出 | 实现要点 |
|---|---|---|
| `sync` | file key / node id → `nodes.json` 和 `@2x.png` | Figma REST：`GET /v1/files/:key/nodes`、`GET /v1/images/:key?scale=2`；支持 `--offline` 读本地 |
| `catalog` | component-map 和组件源码 → `catalog.md` | 解析组件的 Props 类型，生成「位置 / 示例 / 约束 / props 表」，交给 Skill 引用 |
| `plan <frame>` | nodes.json 和 map → `plan.md`/`plan.json` | 遍历 INSTANCE，按 key/name 匹配 map，结合 variant 推导 props；输出「Figma 节点 → 组件 → 建议 JSX」清单，未映射的列为待决项 |
| `audit <file.tsx>` | TSX 和 plan → 违规列表 | tree-sitter TSX 解析 import 和 JSX；检查组件覆盖、原生元素、硬编码颜色和尺寸、constraints 中能静态检查的规则 |
| `diff --url` | 页面 URL 和 nodes.json / png → `diff.png` 和节点偏差 | chromedp 驱动 Chromium（viewport 1440，DPR 2）；像素对比用 Go 实现 pixelmatch（YIQ 加抗锯齿检测）；节点级：执行 JS 读取 `[data-figma-id]` 的 `getBoundingClientRect` 和 `getComputedStyle`，和 Figma 节点属性逐项比较；差异像素聚类后归到对应 Figma 节点 |
| `report` | 以上全部 → `report.json`/`report.md` | 汇总指标；不达标时返回非零退出码（可接 CI） |

`diff` 的输出要能直接交给 Agent 修正：

```
✗ [12:88] StatCard「本月订单」 padding: 16px → 期望 24px
✗ [12:91] Badge「已发货」      未使用 <Badge>，检测到 <span style="color:#16a34a">
✗ [14:02] Table 行高           48px → 期望 56px
  pixel mismatch 7.8%（阈值 3%），最集中区域：StatCard 行（见 diff.png）
```

工程要求：每个 internal 包都有单元测试（table-driven）；`go vet` 通过；`Makefile` 提供 `make build / test`。

### 4.5 Skill `figma-to-code`

**触发**：用户要求按 Figma 实现、还原页面或组件。

**流程（SKILL.md 写死，不允许跳步）**：

1. **加载 Context**：读 `references/catalog.md` 和 `rules.md`；拿到 Figma 链接或 frame 后执行 `figmafe sync`
2. **Plan（闸门）**：执行 `figmafe plan` → 整理成「页面结构树、每个区块用哪个组件、props、未映射项及处理方案（扩展现有组件 / 新建 / 用基础元素）」→ **停下来等用户确认**
3. **实现**：只用 Plan 里确认的组件和 token；每个对应 Figma 实例的根节点必须加 `data-figma-id`；文案和数据取自 mock
4. **验证闭环**：`figmafe audit` 加 `figmafe diff` → 按报告逐条修正 → 重复，直到 report 达标或满 3 轮；满 3 轮仍不达标就向用户汇报剩下的偏差和原因，不硬凑
5. **沉淀**：本次发现的新问题（例如「Table 默认行高与设计不符」）以规则形式写进 `rules.md` 或 map 的 constraints，并附上来源页面

**SKILL.md 本身**：保持精简（≤ 150 行），细节放在 references 里按需读取；description 写清楚触发条件。

---

## 5. 里程碑

| # | 内容 | 完成标准 |
|---|---|---|
| M1 | Figma 设计数据定稿（真实文件或 fixture） | nodes.json 和 png 入库 |
| M2 | token、组件库、mock 数据、component-map | 组件 storybook 风格的预览页能跑起来 |
| M3 | CLI 六个命令和单元测试 | `make test` 全绿；用一份手写的「标准答案页面」跑 report 全部达标（证明 CLI 准确） |
| M4 | Skill 和 catalog/rules | 人工走一遍 C 流程跑通 |
| M5 | A/B/C × 3 实验 | 9 份产出和报告入库 |
| M6 | 结论报告 | `docs/report.md`，加上可分享的页面（指标表、三组截图和 diff 图并排、规则清单、推广步骤） |

M3 里的「标准答案页面」只用来校准 CLI，不参与实验。

---

## 6. 风险与对策

| 风险 | 对策 |
|---|---|
| 用 fixture 时设计数据是自己造的，说服力打折 | 优先用真实 Figma；用 fixture 时先定稿再写组件，并公开数据 |
| 字体和渲染差异导致像素噪声 | 本地打包同一种字体（Inter），截图前等字体加载完；pixelmatch 开启抗锯齿容差；阈值经 M3 校准 |
| LLM 输出波动 | 每组 3 次，报告给出均值和区间 |
| A 组被仓库里的线索「污染」 | CLAUDE.md 不提 Skill；A/B 的 worktree 删掉 `.claude/skills` 和 `design/component-map.yaml` |
| `data-figma-id` 依赖 Agent 自觉加 | audit 检查覆盖率；漏挂的节点退回到按文本内容和位置匹配 |

---

## 7. 推广路径（报告结尾部分）

1. 每个业务组件库维护一份 component-map（由组件 owner 负责）
2. Skill 和 CLI 放进各前端仓库，`figmafe report` 接入 PR CI
3. 每做完一个页面，把 `rules.md` 新增的规则回收到共享规则库

---

## 8. 需要你确认

1. **Figma 来源**：有真实 Figma 文件和 token 吗？没有就用 fixture。
2. **前端栈**：默认 React + TS + Vite 加自建组件库；如果真实项目用 antd、MUI 或 Vue，我对齐那个。
3. **典型页面**：默认「订单管理」列表页；你们有更典型的页面可以换。
4. **实验规模**：默认 A/B/C × 3 共 9 次；想省时间可以只做 A 和 C。
