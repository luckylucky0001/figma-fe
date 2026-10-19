# Plan v3：先盘资产，再识别复用 —— Figma → H5 / Android 双端组件复用

> v1 → v2 变化：
> 1. 目标端从单一 Web 改为 **H5 + Android 两端**，各自基于自己的组件库实现
> 2. 新增 **Phase 0：资产扫描**。先把两端「已有什么组件、什么 token、长什么样、用了多少次」扫描成结构化清单，之后才做映射和复用识别
> 3. 每一步都对应到业界的 best practice 出处（见第 2 节）
>
> v2 → v3 变化（只用免费方案）：
> 4. **不使用任何付费能力**：Figma 官方 MCP（免费档每月只有 6 次调用）、Code Connect（要 Organization / Enterprise）、Variables REST API（要 Enterprise）。它们只作为「思路出处」，不作为依赖
> 5. Figma 数据改为通过**自建的 Figma 导出插件**获取（Plugin API 免费，没有调用配额）；REST API 只作为可选通道，用一次缓存一次
> 6. 开源调研见 `docs/RESEARCH.md`：Mística / Spark / Backpack 等多端设计系统的工程实践，以及可以复用的开源工具

---

## 1. 目标与交付物

| 交付物 | 说明 |
|---|---|
| **Skill** `figma-to-code` | Agent 的强制流程：扫描资产 → 映射 → 页面复用方案（等人确认）→ 实现 → 审计和视觉对比 → 修正 → 沉淀规则；H5 和 Android 两端通用，按 `--platform` 切换 |
| **CLI** `figmafe`（Go） | 把每一步做成可执行、可量化的命令；两端的差异由各自的 adapter 处理 |
| **Demo** | H5 组件库、Android 组件库、同一个 Figma 页面、两端 A/B/C 对照实验、结论报告 |

---

## 2. Best practice 调研结论 → 设计决策

| 业界做法 | 出处 | 在本方案里怎么落地 |
|---|---|---|
| **先做 Interface Inventory**：建设计系统或做复用之前，先把现有 UI 部件全部盘点、分类 | Brad Frost「Interface Inventory」 | **Phase 0 资产扫描**，产出 `inventory.json` 和可视化目录 |
| **组件清单要做成结构化 Manifest**，Agent 读 manifest（接口、变体、示例），不去扫整个源码树，token 消耗少得多，准确度也更高 | Storybook MCP / Component Manifest | 扫描结果统一成 manifest 格式（名称、路径、props、变体、示例、预览图、使用次数） |
| **用使用量数据衡量组件采用率**（静态分析统计组件和 props 的使用情况） | react-scanner、Omlet、Productboard 实践 | 扫描时顺带统计每个组件在业务代码里的使用次数；Agent 优先用高频组件，同时作为复用率指标的基线 |
| **Figma 组件 ↔ 代码组件建立显式映射**；一个 Figma 组件可以按平台分别映射，每个平台一个 label（React / Compose / …） | Figma Code Connect 的思路（付费服务**不用**）；Mística 开源的 `*.figma.ts` / `*.mobile.figma.ts` 按端拆分映射 | `component-map.yaml` 每条记录带 `h5` 和 `android` 两个实现；能离线解析开源的 `*.figma.ts` 文件导入 |
| **生成代码前先查映射**，有映射就复用，不重新生成 | Figma MCP `get_code_connect_map`、`create_design_system_rules` 的思路（付费，**不用**）；Mística `mistica-react` Skill「先读 llms.md 再写代码」 | `figmafe catalog` 代替 `create_design_system_rules`，`figmafe plan` 代替 `get_code_connect_map`；Skill 的 Plan 闸门：`figmafe plan` 先给出复用清单，人确认后才写代码 |
| **先拿轻量结构（metadata），再按节点取细节**；token 用 variables 取，不从截图里猜 | Figma MCP 的分层取数思路（**不用**）；Framelink Figma-Context-MCP 的节点精简做法 | 自建 Figma 导出插件（免费）一次导出节点树、组件集、variables（DTCG）和 PNG @2x；`figmafe import` 精简后入库 |
| **设计 token 用统一格式管理，再生成各端产物** | W3C DTCG Format 2025.10（稳定版）、Style Dictionary | token 统一成 DTCG JSON，生成 H5 的 CSS 变量和 Android 的 Compose Theme / XML 资源；扫描时校验两端 token 和 Figma variables 是否一致 |
| **实现完要和 Figma 截图对照验证** | Figma mcp-server-guide 推荐流程的最后一步（只借鉴流程） | Phase 4 视觉对比：像素级加节点级 |
| **Android 用 JVM 截图测试**（不需要模拟器，支持 Compose Preview） | Roborazzi（兼容 Robolectric，支持 Compose 和 View） | Android 端用 Roborazzi 截图，以 Figma 导出图作为对比基准 |
| **节点级位置和尺寸取 semantics 树** | Compose `SemanticsNode`（size、positionInRoot）、`testTag` / `testTagsAsResourceId` | Android 用 `testTag = figma 节点 id` 导出 bounds；H5 用 `data-figma-id` 加 `getBoundingClientRect` |
| **用 lint 强制「只能用设计系统组件」**（禁止直接用 Material 或原生元素，并给出替换建议） | 社区自建 Compose lint、DuckDuckGo 设计系统 lint、Slack compose-lints | `figmafe audit` 做静态检查；长期建议沉淀为 Android Lint 和 ESLint 规则，接入 CI |
| **Skill 写法**：SKILL.md 保持精简，细节放 references 按需读取，确定性的步骤做成脚本，用有/无 Skill 的 A/B 评测验证效果 | Agent Skills 规范（progressive disclosure）、skill-creator 的评测方法 | SKILL.md ≤ 200 行；确定性步骤全部交给 CLI；用 A/B/C 实验证明效果 |

---

## 3. 整体流程

```
             ┌──────────────────────── Phase 0：资产扫描（Inventory）────────────────────────┐
H5 组件库 ──▶│ figmafe scan --platform h5      → 组件/props/变体/示例/预览图/使用次数/token  │
Android 库 ─▶│ figmafe scan --platform android → @Composable 或 View/参数/Preview/截图/token │
Figma 库  ──▶│ 导出插件 → figmafe import        → Figma 组件集/变体属性/variables            │
             └───────────────────────────────────────┬──────────────────────────────────────┘
                                                     ▼
             ┌──────────────────────── Phase 1：映射（Mapping）───────────────────────────────┐
             │ figmafe map suggest  名称 + 变体/props 结构 + 预览图视觉相似度 → 候选映射       │
             │ 人工确认 → component-map.yaml（每条都有 h5 和 android 两个实现）                 │
             │ figmafe map check    缺口报告：Figma 有但代码没有 / 代码有但 Figma 没有 / 两端不一致 │
             └───────────────────────────────────────┬──────────────────────────────────────┘
                                                     ▼  （Phase 0、1 是一次性的，组件库变化时增量更新）
             ┌──────────────────────── Phase 2：页面复用方案（每个页面）────────────────────────┐
Figma 页面 ─▶│ figmafe plan --platform h5|android → 节点树 → 组件 → props → 建议代码；待决项   │
             │ ⛔ 闸门：人确认方案后才进入下一步                                                │
             └───────────────────────────────────────┬──────────────────────────────────────┘
                                                     ▼
                                    Phase 3：Agent 按方案实现
                                                     ▼
             ┌──────────────────────── Phase 4：验证闭环（≤3 轮）──────────────────────────────┐
             │ figmafe audit  静态检查：实际用到的组件 vs 方案；直接用 Material/原生元素；硬编码值 │
             │ figmafe diff   截图 vs Figma：像素差异 + 节点级尺寸/间距/颜色/字号偏差          │
             │ figmafe report 结构化偏差清单 → Agent 修正 → 再跑                               │
             └───────────────────────────────────────┬──────────────────────────────────────┘
                                                     ▼
                                    Phase 5：沉淀 → rules.md 和 map 的 constraints
```

---

## 4. 各部分详细设计

### 4.1 Phase 0：资产扫描

**H5 adapter**（以 React + TS 为例；如果是 Vue + Vant 类组件库，换成 SFC 解析）

| 扫描项 | 方法 |
|---|---|
| 组件清单 | 从组件库入口（`index.ts`）沿 export 解析，得到组件名和文件路径 |
| Props / 变体 | tree-sitter 解析 Props 类型：联合字面量类型（如 `'primary' \| 'secondary'`）识别为变体轴 |
| 示例 | 读 `*.stories.tsx` 或 `demo/`；如果都没有，就从业务代码里取一处真实用法 |
| 预览图 | 用 chromedp 渲染每个组件的每个变体（移动端视口，例如 375 宽，DPR 2）并截图 |
| 使用次数 | 统计业务代码里每个组件和每个 props 组合的出现次数（类似 react-scanner） |
| Token | 解析 CSS 变量 / Less 变量 → 统一成 DTCG 格式 |

**Android adapter**（Compose 优先；兼容 XML 自定义 View）

| 扫描项 | 方法 |
|---|---|
| 组件清单 | tree-sitter-kotlin 扫描设计系统模块里的 `public @Composable` 函数；XML 体系则扫描 `View` 子类和 `attrs.xml` 里的 `declare-styleable` |
| 参数 / 变体 | 函数参数里的 enum / sealed class 识别为变体轴，`Modifier` 和 slot 参数单独标记 |
| 示例 | 读 `@Preview` 函数（也兼容 Showkase 的 `@ShowkaseComposable`） |
| 预览图 | 用 Roborazzi 给所有 `@Preview` 截图（JVM 上跑，不需要模拟器） |
| 使用次数 | 统计 app 模块里的调用次数；同时统计「绕开设计系统直接用 Material3」的次数（作为现状基线） |
| Token | 解析 `Color.kt` / `Theme.kt` / `colors.xml` / `dimens.xml` → DTCG |

**统一产物**：`assets/inventory.json`（manifest），外加一个人能看的可视化目录页（每个组件、每个变体、每个平台一张预览图）。

### 4.2 Phase 1：映射

`figmafe map suggest` 给每个 Figma 组件集，在每个平台上找候选实现，综合打分：

1. **名称相似度**（归一化后计算编辑距离，再加同义词表，例如 Tag ≈ Badge、Cell ≈ ListItem）
2. **变体结构相似度**（Figma 的 variant 属性和值，对比代码的 props / enum）
3. **视觉相似度**（Figma 组件导出图 vs 代码预览图，用感知哈希或 SSIM）

输出候选和置信度，**人工确认后**写进 `component-map.yaml`：

```yaml
- figma: { set: "Button", key: "…" }
  variants: [Type: [Primary, Secondary, Text], Size: [L, M, S], State: [Default, Disabled]]
  h5:
    import: "@ds/h5/Button"
    props: { "Type=Primary": 'type="primary"', "Size=L": 'size="large"', "State=Disabled": "disabled" }
    usage: '<Button type="primary" size="large">提交订单</Button>'
  android:
    call: "com.xx.ds.button.DsButton"
    props: { "Type=Primary": "style = DsButtonStyle.Primary", "Size=L": "size = DsButtonSize.Large", "State=Disabled": "enabled = false" }
    usage: 'DsButton(text = "提交订单", style = DsButtonStyle.Primary, onClick = {})'
  constraints:
    - 一屏最多一个 Primary
    - Android 禁止直接用 material3.Button；H5 禁止用原生 <button>
```

`figmafe map check` 输出**缺口报告**，这本身就是很有价值的产出：

- Figma 有、某一端没有的组件 → 需要补组件，或明确降级方案
- 代码有、Figma 没有的组件 → 设计侧缺口，或者已经废弃
- 两端变体不一致（例如 Android 少了 Size=S）
- token 值不一致（Figma variables vs H5 vs Android）

可选：如果组件库里已经有开源格式的 Code Connect 文件（`*.figma.ts`，例如 Mística），`figmafe map import-code-connect` 离线解析，作为映射的初始候选。全程不调用 Code Connect 服务。

### 4.3 Phase 2：页面复用方案（闸门）

`figmafe plan <frame> --platform h5|android` 的输出：

- 页面结构树：每个节点标注**复用哪个组件、props 是什么、对应的 Figma 节点 id**
- 布局容器和间距对应的 token（不允许出现裸数值）
- **待决项**：没有映射的节点，给出建议方案（扩展现有组件 / 新建组件 / 用基础布局拼）
- Skill 规定：**方案必须经人确认**；实验里自动批准，但记录方案准确率

### 4.4 Phase 4：验证闭环

| | H5 | Android |
|---|---|---|
| 静态审计 | tree-sitter TSX：实际用到的组件 vs 方案；原生 `<button>` / `<input>` 等；硬编码 hex / px | tree-sitter-kotlin：实际用到的组件 vs 方案；直接 import `androidx.compose.material3.*` 的组件；`Color(0x…)`、裸 `.dp` / `.sp` |
| 截图 | chromedp，视口和 DPR 与 Figma 导出对齐 | Roborazzi 截取页面 Composable，dp 和 density 与 Figma 导出对齐 |
| 节点级 | `[data-figma-id]` → `getBoundingClientRect` 和 `getComputedStyle` | `Modifier.testTag(figmaId)` → semantics 树（`useUnmergedTree`）导出 bounds 和颜色 |
| 像素对比 | Go 实现 pixelmatch（YIQ 加抗锯齿检测），输出 diff 图，差异区域归到对应的 Figma 节点 | 同左 |

输出给 Agent 的是**结构化偏差清单**：

```
[android] ✗ 12:88 OrderCard 内边距 12dp → 期望 16dp（token: space/lg）
[android] ✗ 12:91 状态标签   使用了 material3.AssistChip，应使用 DsTag(tone = Success)
[h5]      ✗ 14:02 列表行高   48px → 期望 56px
[h5]      pixel mismatch 6.1%（阈值 3%），主要集中在：OrderCard 区域
```

### 4.5 Skill 结构

```
.claude/skills/figma-to-code/
├── SKILL.md                     # 流程和闸门；≤200 行
└── references/
    ├── workflow-h5.md           # H5 端的命令、截图参数、常见坑
    ├── workflow-android.md      # Android 端的命令、Roborazzi、testTag 约定
    ├── catalog-h5.md            # 由 inventory 和 map 自动生成：组件 / 示例 / 约束
    ├── catalog-android.md       # 同上
    └── rules.md                 # 沉淀下来的规则（持续增长）
```

SKILL.md 的硬性规则：

1. 没有 inventory 和 component-map，或者它们比组件库源码旧 → 先跑 `scan` 和 `map check`，**不允许直接写页面**
2. 写页面之前必须产出 plan 并得到确认
3. 只能使用 catalog 里的组件和 token；需要新增组件时走待决项流程，不能就地自己拼
4. 交付前 `report` 必须达标，或者在 3 轮之后如实汇报剩下的偏差
5. 每完成一个页面，把新发现的问题写进 rules.md

### 4.6 CLI `figmafe`（Go）命令

| 命令 | 作用 |
|---|---|
| `import <export.zip>` | 读取自建 Figma 导出插件的产物（节点 JSON、组件集、variables 的 DTCG JSON、PNG @2x），精简后入库。**主通道，免费** |
| `sync <figma-url>`（可选） | Figma REST：`GET /v1/files` 和 `/v1/images` 各调用 1 次，结果缓存；免费档每月 6 次额度足够每月同步两三次。不调用 Variables API（要 Enterprise） |
| `scan --platform h5\|android` | Phase 0 资产扫描 → `inventory.json`、预览图、目录页 |
| `map suggest` / `map check` / `map export-code-connect` | Phase 1 映射 |
| `catalog` | 由 inventory 和 map 生成给 Skill 用的 catalog-*.md |
| `plan <frame> --platform` | Phase 2 页面复用方案 |
| `audit <path> --platform` | 静态审计 |
| `diff <frame> --platform` | 截图（H5 用 chromedp；Android 调用 Gradle Roborazzi 任务）、像素对比、节点对比 |
| `report` | 汇总指标；不达标时返回非零退出码（可接 CI） |

代码结构：`internal/figma`、`internal/inventory`、`internal/platform/{h5,android}`（实现同一个 `Adapter` 接口：Scan / Audit / Capture / Nodes）、`internal/mapping`、`internal/diff`、`internal/report`。各包都有 table-driven 单元测试。

---

## 5. Demo 与对照实验

### 5.1 Demo 素材

- **组件库**：最好直接用你们真实的 H5 组件库和 Android 组件库（哪怕是裁剪版）。拿不到的话，我搭两个小型库（各 10～12 个组件：Button、Tag、Cell/ListItem、Card、NavBar、Tabs、Input、Price、Stepper、Dialog、Toast、Avatar），并故意做出常见的不一致：命名不同（Tag vs Badge）、一端少某个变体、token 有偏差
- **Figma**：最好用真实文件。没有的话，我按 Figma REST 的真实结构手工构造 fixture（组件库文件和页面文件），**先定稿设计，再写组件**，避免自己出题自己答
- **典型页面**：移动端「订单详情」或「订单列表」页（NavBar、状态卡片、商品 Cell、价格、Tag、底部操作栏、Dialog）

### 5.2 实验分组（两端各做一次）

| 组 | Agent 拿到的东西 |
|---|---|
| A 现状 | Figma 数据、导出图、组件库源码（组件就在仓库里，但没有清单和映射） |
| B 加资产 Context | A，加上 inventory、component-map、catalog |
| C 完整 Skill | B，加上 Plan 闸门和 audit/diff 修正闭环 |

每组跑 3 次，两端合计 18 次（只做 A 和 C 的话是 12 次）。指标：组件复用率、违规数、像素偏差、节点级偏差、Plan 准确率、修正轮数。

B 组和 A 组对比，可以单独证明「先扫描资产」这一步的价值；C 组和 B 组对比，证明闭环的价值。

### 5.3 产出

`docs/report.md`，同时发布成一个可分享的页面，内容包括：

- 指标表
- 两端 A/B/C 截图和 diff 图并排
- 缺口报告
- 沉淀出的规则
- 推广步骤：每个组件库接入 `scan` 和 `map` → Skill 进各端仓库 → `report` 接入 PR CI

---

## 6. 里程碑

| # | 内容 | 完成标准 |
|---|---|---|
| M1 | 确定素材：真实组件库 / Figma 还是 fixture | 素材入库 |
| M2 | Figma 导出插件；CLI：import、scan（两个 adapter） | 两端 inventory 和目录页生成正确；单元测试通过 |
| M3 | CLI：map suggest/check 和 catalog | 映射建议的准确率可以人工核对；缺口报告输出 |
| M4 | CLI：plan / audit / diff / report | 用一份手写的标准实现跑出达标结果，以此校准阈值 |
| M5 | Skill | 人工把 C 流程两端各走通一遍 |
| M6 | 对照实验 | 所有运行的产出都入库到 `experiments/` |
| M7 | 报告 | 报告和分享页面完成 |

---

## 7. 风险与对策

| 风险 | 对策 |
|---|---|
| Android 截图需要 JDK、Android SDK 和 Gradle；当前云端环境可能没有，或下载受网络限制 | M1 先验证环境；不行的话，Android 截图在你本地或 CI 跑，CLI 只负责编排和对比 |
| Figma MCP、Code Connect、Variables API 都要付费 | **全部不用**；自建导出插件，加自有的 component-map 和 catalog |
| 免费 REST 额度只有每月 6 次 | 插件作为主通道；REST 只作为可选通道，并且缓存结果 |
| 字体和渲染差异带来像素噪声 | 两端用同一种字体；开启抗锯齿容差；在 M4 校准阈值 |
| 自己构造素材，说服力打折 | 优先用真实素材；用 fixture 时先定稿设计、数据公开 |

---

## 8. 需要你确认

1. **你说的两点，我的理解是**：同一份 Figma 设计，H5 端基于 H5 组件库实现，Android 端基于 Android 组件库实现，两端都要先盘点资产。还是说 **H5 组件本身就是参照 Android 组件做的**，要以 Android 组件为基准去对齐 H5？
2. **技术栈**：H5 是 React 还是 Vue（有没有基于 Vant / NutUI 之类的库）？Android 是 Compose 还是 XML View？
3. **素材**（详见 `docs/RESEARCH.md` 第 5 节）：
   - **方案 A（推荐）**：你们真实的 H5 和 Android 组件库（可以裁剪），加上你们的 Figma 页面，用免费插件导出
   - **方案 B**：Material 3：免费的 Figma M3 Kit，加 Compose Material3 和 material-web。设计稿是第三方出的，最客观
   - **方案 C**：仿 Spark 的结构自建两端组件库，再手工构造 Figma 数据
4. **实验规模**：两端 A/B/C 共 18 次，还是只做 A 和 C 共 12 次？
