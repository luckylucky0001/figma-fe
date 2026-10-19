# 报告：Figma → 代码的组件复用与设计还原度（痛点 4 验证）

> PDF 版（3 页）：[`docs/report.pdf`](report.pdf)。数据：[`experiments/summary.md`](../experiments/summary.md)。操作教程：[`docs/tutorial/README.md`](tutorial/README.md)。

## 结论

1. **流程能跑通，也能量化。** 「盘点资产 → 建立映射 → 复用方案（闸门）→ 实现 → 审计 + 截图对比 → 修正 → 沉淀规则」在 H5（React）和 Android（Compose，兼容 compileSdk 32 / minSdk 23）两端都已经实现成命令。其中 `figmafe verify` 一条命令输出复用率、违规数、像素差异和节点级偏差，退出码可以直接给 CI 当门禁。
2. **实测结果**：一共 13 次对照运行。C 组（完整 Skill）的 H5 运行 3/3 次达标。Android 端剩下的 2 处严重偏差都出在组件库内部（DsPriceText 高度、DsGoodsCard 规格行），A、C 两组完全一样。也就是说，verify 能区分「页面写错了」和「组件本身不符合设计」。
3. **需要如实说明的一点**：在这个 demo 仓库里，A 组基线的复用率也是 100%。只给截图、组件库换成 npm 包、模型换成 Haiku，结果都一样。所以本 demo **证明不了「Skill 提升了复用率」**。两组的差距体现在规范和可追溯性上：
   - A 组 **7/8** 次在页面里自造了 CountDown，其中 6 次还带计时逻辑；C 组 **0/5** 次，都按缺口方案实现成静态文本，并登记了组件需求。
   - A 组 **2/8** 次给节点挂了 Figma 节点 id；C 组 **5/5** 次。挂了 id，每处偏差都能定位到具体节点。
   - token 漂移、变体缺失、两套资产并存、存量页面里的坏写法，这些资产层面的问题只有 C 组的 scan / map check 系统地查了出来。
   - 代价：C 组平均耗时是 A 组的 1.5 倍，token 是 1.3 倍。
4. **下一步**：拿真实项目量一次现状。真实项目里复用率低的原因（仓库大、组件命名和 Figma 对不上、新旧两套组件并存、设计稿没有组件化）在这个 demo 里复现不出来，只能到真实项目上量。做法见下文「接入真实项目」。

## 实验设计

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

## 结果

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

## 资产层面的发现

完整清单见 [`docs/component-issues.md`](component-issues.md)：token 漂移、变体缺失、CountDown 缺口、Compose 和 XML View 两套资产并存、5 处组件内部偏差、存量 XML 页面里的坏写法。

## 工具校准（实验中发现并已修复）

- **文本位置比对改为对齐方式无关**：左对齐、居中、右对齐，只要有一种对上就算一致。原因是 Figma 的文本框常常是定宽，代码里多是随内容收缩，原来的比法会误报。这个问题由 C 组 Agent 发现。
- **文本匹配优先找纯文本节点**：之前会匹配到组件的根节点（比如 DsTopBar 本身也带着标题文字），导致位置判断错误。由 C 组 Android Agent 发现。
- **Android 像素阈值从 4% 调到 5%**：页面完全正确时，Skia 和 Chromium 字形渲染不同带来的本底差异约为 4.0%。
- **catalog 只取紧挨着组件声明的那段 JSDoc / KDoc**：之前会把前面属性的注释也混进来。由 C-haiku 组的 Agent 发现。

## 局限

- **Figma 数据不是真实文件**：本环境既没有真实的 Figma 文件，也访问不了 Figma（`api.figma.com` 被网络策略拒绝）。设计数据是用 `figmafe fixture` 从一份独立手写的设计源生成的，格式和插件导出一致；插件代码本身没有在真实 Figma 里跑过。
- **组件库是自建的小库**，不是你们的真实资产；样本量也小，每种条件只跑了 1～2 次。
- **Android 不是真机渲染**：用的是 Compose Desktop（JVM 上的 Compose，API 一致），因为本环境访问不了 `dl.google.com`，装不了 Android SDK。真实项目把渲染器换成 Roborazzi 或 Paparazzi 即可。

## 接入真实项目

1. 用插件导出组件库和一个典型页面。
2. 改 `figmafe.yaml`，指向真实的 H5 组件库（源码目录，或者 npm 包的 `dist/types`）和 Android 组件库（Compose 或 XML View）。
3. 跑 `scan` → `map suggest` → 由组件 owner 确认 → `map check`。
4. 拿这个页面跑 A / C 对照，各 3 次，得到真实的 before / after 数据。
5. 把 `verify` 接进 PR CI，`rules.md` 随着页面增多持续积累。
