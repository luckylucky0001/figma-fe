# 验证报告：OrderDetail（android）

**❌ FAIL：节点级严重偏差 2 处 > 0**

| 指标 | 值 |
|---|---|
| 组件复用率（静态，按方案） | 100% |
| 组件复用率（运行时，由组件库渲染） | 100% |
| 违规（原生元素 / 直接用 Material / 硬编码颜色 / 老组件） | 0 |
| 警告（疑似自造组件 / 裸数值 / 少用组件） | 4 |
| 像素差异 | 4.00% |
| 节点级偏差（严重 / 轻微） | 2 / 6 |
| Figma 文本覆盖 | 100% |
| 挂了 Figma id 的节点 | 0 |

## 组件使用（方案 vs 实际）

| 组件 | 方案 | 实际 |
|---|---|---|
| `DsBadge` | 1 | 1 |
| `DsBanner` | 1 | 1 |
| `DsBottomBar` | 1 | 1 |
| `DsButton` | 3 | 3 |
| `DsCard` | 3 | 3 |
| `DsDivider` | 2 | 2 |
| `DsGoodsCard` | 1 | 1 |
| `DsListItem` | 7 | 7 |
| `DsPriceText` | 1 | 1 |
| `DsTopBar` | 1 | 1 |

## 需要修正（按严重程度）

- ✗ [diff/size] 1:41 Price「¥ 269 .00」：期望 68×26，实际 68×32
- ✗ [diff/position] 1:23 文本「米白色 / M」：期望 (120, 361)，实际 (120, 366)
- · [diff/position] 1:41 Price「¥ 269 .00」：期望 (279, 553)，实际 (279, 550)
- · [diff/position] 1:57 Button「联系客服」：期望 (75, 856)，实际 (79, 856)
- · [diff/line-height] 1:5 文本「商家已接单，预计 10月22…」：期望 行高 36，实际 高 18.0
- · [diff/position] 1:22 文本「轻薄羽绒服 女款 短款连帽 …」：期望 (120, 321)，实际 (120, 321)
- · [diff/line-height] 1:22 文本「轻薄羽绒服 女款 短款连帽 …」：期望 行高 20，实际 高 41.0
- · [diff/position] 1:58 文本「联系客服」：期望 (92, 862)，实际 (95, 862)
- · [audit/lookalike] android/app/src/main/kotlin/com/demo/app/screens/OrderDetail/OrderDetailScreen.kt:70 页面里自定义了 StatusCard：组件库里有对应组件时应直接复用
- · [audit/lookalike] android/app/src/main/kotlin/com/demo/app/screens/OrderDetail/OrderDetailScreen.kt:88 页面里自定义了 GoodsCard：组件库里有对应组件时应直接复用
- · [audit/lookalike] android/app/src/main/kotlin/com/demo/app/screens/OrderDetail/OrderDetailScreen.kt:109 页面里自定义了 InfoCard：组件库里有对应组件时应直接复用
- · [audit/raw-number] android/app/src/main/kotlin/com/demo/app/screens/OrderDetail/OrderDetailScreen.kt:0 出现 1 处裸 dp/sp 数值；优先用 DsSpacing / DsTypography

## 像素差异最集中的区域（见 diff.png）

- 1:11 CountDown「剩余 23:59:59」：29%
- 1:8 Tag「待发货」：24%
- 1:57 Button「联系客服」：16%
- 1:59 Button「修改地址」：14%
- 1:41 Price「¥ 269 .00」：13%
- 1:61 Button「提醒发货」：10%
- 1:20 ProductCard「轻薄羽绒服 女款 短款连…」：8%
