# 组件库待修清单（由 figmafe map check + verify 发现）

> 这些问题出在组件库或设计资产本身，页面层不应该补偿。每一条都能复现：`bin/figmafe map check`，或者对 OrderDetail 跑 `bin/figmafe verify`。

| # | 端 | 组件 / 资产 | 问题 | 证据 | 建议 |
|---|---|---|---|---|---|
| 1 | Android | `DsColors.Warning` | Token 漂移：代码 #FF8800，Figma `color.warning` #FF8F1F | map check `token-drift` | 改为 #FF8F1F（XML 资源 `colors.xml` 已经是正确值） |
| 2 | Android | `DsBadge` | 缺 `Variant=Solid` 变体 | map check `variant-unsupported` | 补 Solid，或在 Figma 标注 Android 不支持 |
| 3 | 两端 | CountDown | Figma 有，两端都没有组件 | map check `missing-impl`、plan 待决项 | 立项实现；本期按静态文本处理 |
| 4 | H5 | `Button` Primary | 1px 同色边框画在盒外，宽度比 Figma 多 2px（88→90） | verify `diff/position` 轻微偏差 | 用 `box-shadow: inset` 或 `border-color: transparent` 加内缩 |
| 5 | H5 | `Tag` | 根节点 `align-self: flex-start`，放进居中的横向行里会贴顶 2px | verify `diff/position` | 去掉 `align-self` |
| 6 | 两端 | `ProductCard` / `DsGoodsCard` | 两行标题按 2×20 行高渲染为 40，Figma 为 36.4；规格、价格、数量整体下移约 4px，Android 上价格还和规格重叠 | verify `diff/position`（1:23 等） | 和设计确认标题行高，统一按 18 或调整卡片高度 |
| 7 | Android | `DsPriceText` | 节点高 32，Figma 26（Large） | verify `diff/size` 1:41 | 给价格文本设固定行高 |
| 8 | Android | `DsButton` Default | 1dp 描边画在内部，宽 88（Figma 90） | verify `diff/position` 轻微偏差 | 描边放在 padding 外 |
| 9 | Android | `LegacyCellView` | XML 存量组件不支持 Figma `Label`（副标题）变体 | scan 资产清单 | 存量页面迁移到 DsListItem |
| 10 | Android | `activity_order_list_legacy.xml` | 存量页面里有原生 `<Button>`，并写死了 `#FF4D1A`、`14sp` | scan 存量问题（baseline） | 迁移到 LegacyButton / DsButton |
