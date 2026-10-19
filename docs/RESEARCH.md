# 开源调研：可参考的组件库与工具

> 调研日期：2026-10-19。star 数、维护状态取自 GitHub 页面或第三方聚合站当时的快照，可能有偏差。

## 1. 结论先行

1. **没有找到一个开源设计系统同时满足「公开免费的 Figma 文件 + H5 组件库 + Android 原生组件库 + 仍在积极维护」**：
   - 结构最接近我们场景的是 Mística（Telefónica）、Spark（Adevinta / leboncoin）、Backpack（Skyscanner），都是「同一套设计 → Web + Android（+ iOS）」，但它们的 Figma 库是内部库，或者没有确认有公开的 Community 文件
   - 唯一「Figma 公开免费 + Android + Web 都是官方」的是 **Material 3**：Figma M3 Design Kit 免费，Compose Material3 是官方实现，Web 端的 `material-web` 已进入维护模式
2. **这些项目最值得参考的不是组件本身，而是它们的工程实践**（第 3 节），可以直接作为 best practice 的出处：
   - Mística：给 Agent 用的 `llms.md` 组件手册、官方 Claude Skill、`*.figma.ts` 映射文件
   - Spark：Figma → 三端 token 同步、lint 模块、截图测试模块、catalog app
   - Backpack：同一个库同时支持 View 和 Compose，自带 lint 模块
3. **工具层面基本都有开源实现可以复用或借鉴**（第 4 节），我们的 CLI 主要做「串起来 + 双端统一 + 结构化偏差输出」，不重复造轮子。

## 2. 付费边界（按用户要求：付费的一律不用）

| 能力 | 免费吗 | 结论 |
|---|---|---|
| Figma 官方 MCP（`get_code_connect_map`、`create_design_system_rules`、`get_design_context` 等） | Starter 计划和 View / Collab 席位：**每月 6 次工具调用**；Dev / Full 席位要付费（Professional 起，每天 200 次） | **不用** |
| Figma Code Connect（`get_code_connect_map` 依赖它） | 只有 **Organization / Enterprise** 计划，并且要 Full 或 Dev 席位 | **不用这个服务**；但 Code Connect 的 CLI 和 `*.figma.ts` 文件格式是 MIT 开源的，可以参考它的映射格式 |
| Figma REST API（个人 token） | 可以用免费 token，但读取 Starter 计划里的文件，Tier 1 接口（`GET /v1/files`、`/nodes`、`/images`）**每月只有 6 次** | 只作为可选通道，用一次缓存一次 |
| Figma Variables REST API | 只有 Enterprise 才能用 | **不用** |
| Figma Plugin API（在 Figma 桌面端跑插件） | **免费，没有调用配额**；开发中的插件在免费账号里也能运行 | ✅ **主通道**：自己写一个导出插件，导出节点 JSON、PNG 和 variables |
| 社区 MCP（Framelink `Figma-Context-MCP`） | 开源，但底层调用的是 REST API，同样受每月 6 次的限制 | 不当主通道，可以借鉴它「精简节点数据」的做法 |

**替代方案**：
- `get_code_connect_map` → 自己维护的 `component-map.yaml`，加上 `figmafe plan`
- `create_design_system_rules` → `figmafe catalog`，从扫描出的资产清单和映射生成 rules 和 catalog
- Figma 数据 → 自己写的导出插件，或开源的 figma-to-json 加 TokensBrücke

## 3. 可参考的开源设计系统（多端）

| 项目 | Web / H5 | Android | Figma | 状态 | 可借鉴点 |
|---|---|---|---|---|---|
| **Mística**（Telefónica） | `mistica-web`（React，MIT） | `mistica-android`（Compose 和 XML，Maven 上有 `mistica-compose`、`mistica-catalog-compose`） | `mistica-design` 仓库通过 Figma API 同步 variables；Figma 库应该是内部的，没找到公开的 Community 文件 | 活跃（2026-10 发布了 v18 beta） | ① 随 npm 包附带 **`doc/llms/llms.md`**，作为 Agent 的组件手册；② 官方 Claude Skill：`mistica-react`（「先读 llms.md 再写代码」）和 `mistica-figma`；③ `src/__figma__/` 里有 **55 个以上组件的 `*.figma.ts` / `*.mobile.figma.ts` 映射文件**，按端拆分；④ 仓库里有 `AGENTS.md` 和 `CLAUDE.md` |
| **Spark**（Adevinta / leboncoin） | `adevinta/spark`（React + Tailwind，MIT，有 Storybook 和 Chromatic） | `leboncoin/spark-android`（Compose，基于 Material3 覆盖 token，MIT） | `adevinta/spark-tokens`：**Figma ↔ Web / Android / iOS 三端 token 同步**，图标自动分发到各端 | 活跃 | ① `spark-lint` 模块（设计系统 lint）；② `spark-screenshot-testing` 模块；③ `catalog` app；④ 「基于 Material3 加覆盖 token」的模式，很接近很多公司 Android 组件库的做法 |
| **Backpack**（Skyscanner） | `backpack`（React，Apache-2.0）——**2026-08 已归档** | `backpack-android`（**同时支持 View 和 Compose**，Apache-2.0） | 仓库里有 `figma.config.json`（Code Connect）和 `token-sync` | Web 已归档，Android 仍在 | ① `backpack-lint` 模块；② 同一套组件同时提供 View 和 Compose 两种实现；③ `templates/component` 组件模板 |
| **Orbit**（Kiwi.com） | `orbit`（React，MIT） | `orbit-compose`——**2025-07 已归档** | 社区评价它的 Figma 文件不太规范 | Android 已归档 | 只作为参考 |
| **Material 3**（Google） | `material-web`（Lit Web Components，Apache-2.0，**维护模式**） | Compose Material3（官方） | **Figma M3 Design Kit 在 Community 免费提供**，Figma 官方 UI kit 还预置了 Code Connect 片段 | Android 活跃，Web 维护模式 | 唯一一套「设计稿公开免费、两端都有官方实现」的组合，适合拿来做**第三方出题**的 demo |

**国内 H5 组件库**（只有 H5，没有对应的 Android 原生库）：

| 项目 | 技术栈 | Figma |
|---|---|---|
| Vant（有赞） | Vue 2 / 3，小程序 | 没找到官方 Figma kit |
| NutUI（京东） | Vue 3 / React，Taro | 没找到 |
| antd-mobile（蚂蚁） | React | 只有第三方付费 kit |
| TDesign Mobile（腾讯） | Vue 3（LTS）/ React（Alpha） | **官方提供 Figma 设计资源**；没有 Android 原生库，只有 Flutter 版 |

## 4. 可复用 / 借鉴的工具

| 环节 | 开源项目 | 用法 |
|---|---|---|
| Figma 数据导出（免费） | `Livshitz/figma-to-json`（插件，导出节点 JSON） | 参考，或者在它的基础上写自己的导出插件 |
| | **TokensBrücke**（插件和 CLI，把 variables 导出成 **DTCG** JSON） | 免费拿到 variables（REST 接口要 Enterprise） |
| | Plugin API `exportAsync` | 插件内导出 PNG @2x |
| | `GLips/Figma-Context-MCP`（MIT） | 借鉴它把 Figma 原始 JSON 精简成「布局 + 样式」给 LLM 的做法 |
| Figma ↔ 代码映射格式 | `figma/code-connect`（CLI MIT） | 参考 `*.figma.ts` 的 props / variant 映射写法；有现成 Code Connect 文件的组件库（如 Mística），可以离线解析成我们的 map，不需要付费服务 |
| H5 资产扫描 | `moroshko/react-scanner` | 统计组件和 props 的使用量（React）；Vue 需要另写 SFC 解析 |
| | Storybook `@storybook/addon-mcp` 和 Component Manifest | 参考 manifest 字段设计（目前只支持 React） |
| Android 资产扫描 | `airbnb/Showkase`（KSP） | 参考它用 KSP 收集 `@Preview` / `@Composable` 生成组件目录的做法 |
| Android 截图 | `takahirom/roborazzi` | JVM 截图，支持 Compose、View 和 Preview，兼容 Robolectric |
| | `cashapp/paparazzi` | 备选（不兼容 Robolectric，不能交互） |
| 像素对比 | `mapbox/pixelmatch`（算法），`dmtrKovalenko/odiff`（快，Lost Pixel 用它做引擎） | Go 里实现 pixelmatch 算法，或直接调用 odiff |
| 规则强制 | `slackhq/compose-lints`、`mrmans0n/compose-rules`、DuckDuckGo 的设计系统 lint、Backpack / Spark 的 lint 模块 | 参考如何写「禁止直接用 Material、给出替换建议」的自定义 lint |
| Agent Skill | Mística `mistica-react` / `mistica-figma` | 「先读组件手册再写代码」的模式；SKILL.md 很短，细节外置 |
| | `figma/mcp-server-guide` 的 skills（`figma-implement-design` 等） | 只参考流程（依赖付费 MCP，不直接用） |
| | `senlindesign/claude2figma`、`scoobynko/claude-code-design-skills` | 「先搜索已有组件再新建」的规则写法 |

## 5. 对我们 Demo 的启示

1. **「组件手册（llms.md / catalog）+ 映射文件 + Skill」三件套已经是业界做法**（Mística 是现成的例子），我们的方案方向正确；我们额外加的是「先扫描资产」「双端统一映射」「截图和节点级对比闭环」，这些是现有项目都没做全的部分。
2. **Demo 素材有三种选择**（见 PLAN 第 8 节）：
   - **方案 A（推荐）**：真实的 H5 和 Android 组件库，加上你们的 Figma 页面（用免费插件导出）。效果最好，最能说服人
   - **方案 B**：Material 3：Figma M3 Kit（免费）+ Compose Material3 + material-web。设计稿是第三方出的，不存在「自己出题自己答」；缺点是 Web 端在维护模式，命名比较统一，映射难度偏低
   - **方案 C**：仿 Spark 的结构自建「基于 Material3 加覆盖 token」的 Android 库和一个 H5 库，故意制造命名差异和变体缺口，再手工构造 Figma 数据。可控，但说服力最弱

## Sources

- Figma 付费边界：[Code Connect (Help Center)](https://help.figma.com/hc/en-us/articles/23920389749655-Code-Connect)、[REST API rate limits](https://developers.figma.com/docs/rest-api/rate-limits)、[Figma MCP rate limits and billing](https://instagit.com/figma/mcp-server-guide/figma-mcp-rate-limits-and-billing-information/)、[Figma MCP pricing](https://toolradar.com/tools/figma-mcp/pricing)、[Figma forum: MCP rate limit](https://forum.figma.com/report-a-problem-6/mcp-rate-limit-not-acknowledging-subscription-51368)
- Mística：[mistica-web](https://github.com/Telefonica/mistica-web)、[mistica-android](https://github.com/Telefonica/mistica-android)、[mistica-design](https://github.com/Telefonica/mistica-design)、[mistica-compose (Maven)](https://mvnrepository.com/artifact/com.telefonica/mistica-compose/6.7.0)、[mistica-react skill](https://skills.sh/telefonica/mistica-web/mistica-react)
- Spark：[spark-android](https://github.com/leboncoin/spark-android/)、[adevinta/spark](https://github.com/adevinta/spark)、[spark-tokens](https://github.com/adevinta/spark-tokens)
- Backpack：[backpack](https://github.com/Skyscanner/backpack)、[backpack-android](https://github.com/Skyscanner/backpack-android)
- Orbit：[orbit](https://github.com/kiwicom/orbit)、[orbit-compose](https://github.com/kiwicom/orbit-compose)、[HN 讨论](https://news.ycombinator.com/item?id=27832680)
- Material：[material-web](https://github.com/material-components/material-web)、[M3 Figma Design Kit](https://material.io/blog/material-3-figma-design-kit)
- 国内：[Tencent/tdesign](https://github.com/Tencent/tdesign)、[掘金：Vant / NutUI / antd-mobile 对比](https://juejin.cn/post/7120411647742590990)
- 工具：[figma-to-json](https://github.com/Livshitz/figma-to-json)、[TokensBrücke](https://cdn.jsdelivr.net/npm/tokens-bruecke@3.15.1/README.md)、[figex](https://github.com/iodigital-com/figex)、[Figma-Context-MCP](https://jimmysong.io/en/ai/figma-context-mcp)、[Productboard + react-scanner](https://www.productboard.com/blog/how-we-measure-adoption-of-a-design-system-at-productboard/)、[Storybook MCP](https://storybook.js.org/docs/ai/mcp/overview)、[Showkase (KotlinConf)](https://resources.jetbrains.com/storage/products/kotlinconf-2023/Showkase_%20Automated%20UI%20Infrastructure%20in%20Jetpack%20Compose%20using%20KSP.pdf)、[Roborazzi (Telefónica)](https://hub.telefonica.com/en/engineering/testing-screenshot-roborazzi-systems-design)、[odiff](https://www.npmjs.com/package/odiff-bin)、[compose-lints](https://p.codekk.com/detail/Android/slackhq/compose-lints)、[DuckDuckGo DS lint](https://github.com/duckduckgo/Android/pull/9249)、[claude2figma](https://github.com/senlindesign/claude2figma)、[figma/mcp-server-guide](https://github.com/figma/mcp-server-guide)
