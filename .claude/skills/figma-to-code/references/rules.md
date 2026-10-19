# 沉淀规则（每完成一个页面追加；格式：现象 → 规则 → 来源）

> 由 C 组 Agent 在 OrderDetail（H5 / Android）验证闭环中发现、人工整理后合并。已经通过修工具解决的问题（文本对齐口径、NavBar 标题误匹配、Android 像素阈值）不再列为规则，见 docs/REPORT.md「工具校准」。

## 通用
- Card 自带 16 内边距和 12 圆角，页面再包一层 padding 会让内容整体错位 → 页面不给 Card 加 padding / 背景 / 圆角 → 组件规范
- 左标题 / 右值的行很容易被写成 flex + 两个 span → 一律用 Cell（H5）/ DsListItem（Android），值通过 value / trailingText 传字符串 → 组件规范
- 金额拼文本会丢失「¥ / 整数 / 小数」三级字号 → 用 Price / DsPriceText；Cell 右侧的金额值是普通文本例外 → 组件规范
- Figma 有、组件库没有的组件（如 CountDown）→ 严格按 component-map 的 gaps 方案实现（本期为静态文本），不要在页面里自造带逻辑的组件 → OrderDetail
- 组件内部的偏差（尺寸、内部间距）→ 页面不要补偿、不要拆开组件自己拼；在交付回复里列出，报给组件库 owner（见 docs/component-issues.md）→ OrderDetail
- 没有对应 token 的尺寸（ActionBar 预留 56、Price 所在行高 32）可以保留裸值，verify 只记为警告 → OrderDetail

## H5
- Tag 的 CSS 自带 `align-self: flex-start`，放进 `align-items: center` 的横向行里会贴顶 2px → 页面给 Tag 挂一个只做 `align-self: center` 的布局类，类名不要含 tag / badge 等字样（否则命中 audit/lookalike）→ OrderDetail（h5）
- 页面 CSS 注释里的 px 也会被 audit/raw-number 计数 → 注释里别写数值单位 → OrderDetail（h5）

## Android
- 页面里的私有分区 Composable 命名为 `XxxCard` 会触发 audit/lookalike → 分区命名为 `XxxSection`，卡片外壳仍用 DsCard → OrderDetail（android）
- DsPriceText 的节点高度是字体自然高度（Large 32），不是设计里的 26 → 不要在页面给 Price 强行设高度；放 Price 的行用设计行高（TotalRow 32）+ `CenterVertically` → OrderDetail（android）
- 页面根用 `Box(fillMaxSize)`，内容 `Column` 底部留 56.dp，`DsBottomBar(Modifier.align(Alignment.BottomCenter))` → OrderDetail（android）
