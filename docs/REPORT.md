# 报告：Figma → 代码的组件复用（痛点 4 验证）

> PDF 版（3 页）：[`docs/report.pdf`](report.pdf)。操作教程：[`docs/tutorial/README.md`](tutorial/README.md)。

## 结论

1. **对存量项目，先把「有什么」查清楚，才知道新设计稿能复用什么。** `figmafe reuse` 扫一遍现有 H5 或 Android 项目，得到组件资产（含真实用法）、项目惯例、重复实现、复制粘贴结构，以及「组件库里有、但本项目用不了」的组件（版本、改名 / 移除、主题依赖、废弃 API、设计系统封装）。`figmafe match` 再把新设计稿逐个元素对上去，给出「直接复用 / 需确认 / 需新建」，附项目里的真实写法和 Figma 变体 → props 的配置。
2. **复用推荐准确率**：21 个真实开源项目（Vue / React / uni-app / Taro + Vant / antd-mobile / NutUI / uni-ui / taro-ui / 自有组件；XML View / Compose + QMUI / Material / AppCompat / 三方 View 库 / 自有设计系统）× 3 份设计稿，共 693 项：**✓665 / △28 / ✗0**。其中 16 个是留出项目，第一次跑时出现过错误（Jetsnack 3 ✗、litemall 2 ✗、hello-uniapp 12 ✗、NewPipe 3 ✗；第二批 mall-app-web 一份设计稿直接崩溃、AntennaPod 1 ✗；第三批 WanAndroid 11 ✗、Seal 7 ✗、KunMinX 3 ✗），修的都是通用规则，之后作为回归测试。数据：[`experiments/real-match-eval.md`](../experiments/real-match-eval.md)、[`experiments/match-judgement.md`](../experiments/match-judgement.md)。
3. **Agent 实测**：5 个真实项目里按同一份设计稿实现新页面，用同一版 `figmafe reuse --focus / --changed` 打分。不用 Skill **0/11** 通过（Now in Android 里新建了 11～12 个和已有组件重复的组件；Haiku 在 newbee 里一个现成组件都没用），用 Skill **11/12** 通过（Haiku 0/4 → 4/4）。代价是多约 35% token、45% 时间。数据：[`experiments/real-project/summary.md`](../experiments/real-project/summary.md)。
4. **新页面闭环**：demo 仓库里 `figmafe verify` 一条命令输出复用率、违规、像素差异和节点级偏差，退出码可做 CI 门禁。但这个仓库小，基线 Agent 复用率也是 100%，证明不了复用率提升；复用率的差距要看第 3 条的真实项目实验。

## 存量项目的用法

| 步骤 | 命令 | 产出 |
|---|---|---|
| 体检 | `figmafe reuse <项目>` | 资产清单、优先处理项、不能用的库组件、项目惯例、指标基线 |
| 新设计稿 | `figmafe match <项目> --figma <导出目录>` | 复用清单 + 叠图 + component-map 初稿；多页面时有整体汇总 |
| 确认 | 编辑 `component-map.yaml` | 需确认 / 需新建的决定，之后 match 和检查都会采用 |
| 写完检查 | `figmafe reuse <项目> --changed origin/main` | 自造轮子、有组件没用、手写已有组件，带行号；未通过退出码 3 |
| 对照设计 | `figmafe compare --url <页面> --figma <导出目录>` | 截图 + 像素 / 文本位置偏差 |
| 不许变差 | `figmafe reuse <项目> --baseline main.json --no-regress` | 组件化率、写死颜色等指标只许变好（CI 示例：`docs/ci/`） |

## 附：demo 仓库的新页面实验

### 实验设计

- **页面**：订单详情，375×900，包含 21 个组件实例，外加 1 个两端都没有实现的组件（CountDown）。
- **组件库**：H5 `h5/src/ds`（React，10 个组件）；Android `android/ds`（Compose，10 个组件，加 1 个辅助组件 `DsText`）；`android/ds-view`（XML View 存量组件，只扫描源码，不编译）。
- **人为埋的问题**：Android 端 `color.warning` 与设计不一致（token 漂移）；DsBadge 缺少 Solid 变体；两端组件命名不一致（NavBar ↔ DsTopBar 等）；存在 Legacy 组件。
- **分组**：每次运行都在独立的工作区里完成，最后统一用同一版 `figmafe verify` 打分（`experiments/score.sh`）。

| 组 | Agent 拿到的东西 |
|---|---|
| A | Figma 设计图 + 节点 JSON + 组件库（就在仓库里） |
| A-pkg | 同 A，但组件库是一个只有编译产物和 `.d.ts` 的 npm 包 `@company/mobile-ui` |
| A-img | 只有设计图，没有节点 JSON |
| A-haiku | 同 A，模型换成 Haiku |
| C / C-pkg / C-haiku | 在对应 A 组的基础上，加上 Skill（资产清单 + 映射 + 组件手册 + 方案闸门 + verify 闭环 + 规则） |

### 结果

| run | 平台 | 结果 | 复用率（静态/运行时） | 违规 | 警告 | 像素差异 | 节点偏差（严重/轻微） | 文本覆盖 | Figma id | 自造 CountDown |
|---|---|---|---|---|---|---|---|---|---|---|
| A-android-1 | android | FAIL | 100% / 100% | 0 | 4 | 4.00% | 2 / 6 | 100% | 28 | 是 |
| A-h5-1 | h5 | PASS | 100% / 100% | 0 | 0 | 0.65% | 0 / 5 | 100% | 0 | 是 |
| A-h5-2 | h5 | PASS | 100% / 100% | 0 | 0 | 0.65% | 0 / 5 | 100% | 0 | 是 |
| Ahk-android-1 | android | FAIL | 100% / 100% | 0 | 4 | 4.00% | 2 / 6 | 100% | 27 | 否 |
| Ahk-h5-1 | h5 | PASS | 100% / 100% | 0 | 0 | 0.65% | 0 / 5 | 100% | 0 | 是 |
| Aimg-android-1 | android | FAIL | 100% / 100% | 0 | 4 | 4.00% | 2 / 6 | 100% | 0 | 是 |
| Aimg-h5-1 | h5 | PASS | 100% / 100% | 0 | 0 | 0.65% | 0 / 5 | 100% | 0 | 是 |
| Apkg-h5-1 | h5 | PASS | 100% / 100% | 0 | 0 | 0.65% | 0 / 5 | 100% | 0 | 是 |
| C-android-1 | android | FAIL | 100% / 100% | 0 | 1 | 4.00% | 2 / 6 | 100% | 27 | 否 |
| C-h5-1 | h5 | PASS | 100% / 100% | 0 | 1 | 0.65% | 0 / 5 | 100% | 27 | 否 |
| Chk-android-1 | android | FAIL | 100% / 100% | 0 | 1 | 4.00% | 2 / 6 | 100% | 27 | 否 |
| Chk-h5-1 | h5 | PASS | 100% / 100% | 0 | 1 | 0.65% | 0 / 5 | 100% | 24 | 否 |
| Cpkg-h5-1 | h5 | PASS | 100% / 100% | 0 | 0 | 0.65% | 0 / 5 | 100% | 27 | 否 |

每次运行的报告、截图、差异图和页面源码都在 `experiments/runs/<run>/`。耗时和 token 用量见 `experiments/cost.json`。

**verify 确实能抓出问题**：教程第 8 步在页面里故意放了 3 处常见错误（自己拼 Tag、用原生 `<button>`、手写 Cell 行），verify 报出了 3 处违规和 8 处严重偏差，每一条都带文件行号或 Figma 节点 id。照着报告改完，第二次 verify 就通过了。

### 资产层面的发现

完整清单见 [`docs/component-issues.md`](component-issues.md)：token 漂移、变体缺失、CountDown 缺口、Compose 和 XML View 两套资产并存、5 处组件内部偏差、存量 XML 页面里的坏写法。

### 工具校准（实验中发现并已修复）

- **文本位置比对改为对齐方式无关**：左对齐、居中、右对齐，只要有一种对上就算一致。原因是 Figma 的文本框常常是定宽，代码里多是随内容收缩，原来的比法会误报。这个问题由 C 组 Agent 发现。
- **文本匹配优先找纯文本节点**：之前会匹配到组件的根节点（比如 DsTopBar 本身也带着标题文字），导致位置判断错误。由 C 组 Android Agent 发现。
- **Android 像素阈值从 4% 调到 5%**：页面完全正确时，Skia 和 Chromium 字形渲染不同带来的本底差异约为 4.0%。
- **catalog 只取紧挨着组件声明的那段 JSDoc / KDoc**：之前会把前面属性的注释也混进来。由 C-haiku 组的 Agent 发现。

## 局限

- **没有真实 Figma 文件**，本环境也访问不了 Figma：三份设计稿由手写的设计源生成（`figmafe fixture`，格式和插件导出一致）。插件只用 mock 的 Plugin API 测过（`make plugin-test`），没在真实 Figma 里跑过。
- **没有编译真实项目**：存量项目实验只做静态检查，没有构建和运行这些开源项目；demo 的 Android 用 Compose Desktop 在 JVM 上渲染，不是真机（本环境装不了 Android SDK）。
- **实验有泄漏**：A 组 agent 也能看到本仓库的 Skill；Rlm-A-1 自己运行了 figmafe，已剔除。泄漏只会让 A 组更好，所以 A / C 差距是保守估计。样本量小（每种条件 1～2 次），只能说明方向。
- **准确率的判断标准是人定的**：✓ / △ / ✗ 由我逐条判断并写成回归用例；真实接入时应由组件 owner 复核一次。

## 接入真实项目

1. `figmafe reuse` 分别扫 H5 和 Android 仓库（只读），把体检报告给组件 owner 看，先处理「优先处理」里的项。
2. 用插件导出下一个要做的页面，`figmafe match` 出复用清单，组件 owner 确认「需确认」项，写进 `component-map.yaml`。
3. 开发时让 Agent 按 Skill 走；提交前 `reuse --changed origin/main`。
4. 把 `reuse --baseline --no-regress` 接进 CI；偏差和决定持续沉淀到 `rules.md` 和 `component-map.yaml`。
