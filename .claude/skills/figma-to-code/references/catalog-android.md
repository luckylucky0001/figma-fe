# 组件手册：Android（Compose）

> 由 `figmafe catalog` 根据资产扫描结果（inventory）和 component-map 自动生成，不要手改。
> 规则：Figma 里出现的组件，**必须**用下表对应的代码组件；表里没有的才允许自己写，并在方案里列为待决项。

## 速查

| Figma 组件 | 代码组件 | 业务代码中已用次数 |
|---|---|---|
| ActionBar | `DsBottomBar` | 1 |
| Button | `DsButton` | 14 |
| Card | `DsCard` | 7 |
| Cell | `DsListItem` | 5 |
| CountDown | **无** —— 见待决项 | - |
| Divider | `DsDivider` | 6 |
| NavBar | `DsTopBar` | 2 |
| NoticeBar | `DsBanner` | 3 |
| Price | `DsPriceText` | 4 |
| ProductCard | `DsGoodsCard` | 3 |
| Tag | `DsBadge` | 2 |

## 组件详情

### ActionBar → `DsBottomBar`

吸底操作栏，按钮右对齐。源码：`android/ds/src/main/kotlin/com/demo/ds/DsBottomBar.kt`。

Bottom action bar (Figma: `ActionBar`). Height 56, bg-card, 1dp top border (color.border), horizontal padding 16, children End-aligned with 8dp spacing, vertically centered. It does not position itself — the screen places it at the bottom.

```
import com.demo.ds.DsBottomBar
```

参数：`modifier: Modifier = Modifier`、`content: RowScope.() -> Unit`

示例：`DsBottomBar { DsButton(text = "联系客服", onClick = {}, style = DsButtonStyle.Default) }`

- 约束：吸底操作区一律用 ActionBar；页面内容区底部要留出 56 的空间
- 注意：DsBottomBar 不负责定位，页面需要把它放在底部（Box + Alignment.BottomCenter）

### Button → `DsButton`

按钮。源码：`android/ds/src/main/kotlin/com/demo/ds/DsButton.kt`。

Button (Figma: `Button`, variants `Type=Primary|Default|Text`, `Size=Large|Medium|Small`). - Primary: bg primary, white text - Default: bg-card, 1dp border, text-primary - Text: no background / border, primary text Large 44 / pad 24 / 16·600 / r22; Medium 32 / pad 16 / 14·500 / r16; Small 24 / pad 12 / 12·500 / r12.

```
import com.demo.ds.DsButton
```

| Figma 变体 | 代码 |
|---|---|
| Size=Large | `size = DsButtonSize.Large` |
| Size=Medium | `size = DsButtonSize.Medium` |
| Size=Small | `size = DsButtonSize.Small` |
| Type=Primary | `style = DsButtonStyle.Primary` |
| Type=Default | `style = DsButtonStyle.Default` |
| Type=Text | `style = DsButtonStyle.Text` |

文本：`text` ← 图层 {Label}

参数：`text: String`、`onClick: () -> Unit`、`modifier: Modifier = Modifier`、`style: DsButtonStyle = DsButtonStyle.Primary`、`size: DsButtonSize = DsButtonSize.Medium`

示例：`DsButton(text = "提醒发货", onClick = { }, style = DsButtonStyle.Primary, size = DsButtonSize.Medium)`

- 约束：一屏最多一个 Primary 按钮
- 约束：H5 禁止原生 <button>；Android 禁止 material3.Button / OutlinedButton
- 注意：不要用 LegacyButton（XML View 存量组件），不要用 material3.Button

### Card → `DsCard`

白底圆角卡片，可带标题。源码：`android/ds/src/main/kotlin/com/demo/ds/DsCard.kt`。

Card container (Figma: `Card`, variant `Title=On|Off`). bg-card, radius.lg (12), padding 16. Optional title (font.subtitle) followed by a 12dp gap; children stacked vertically with no spacing.

```
import com.demo.ds.DsCard
```

| Figma 变体 | 代码 |
|---|---|
| Title=On | 默认 |
| Title=Off | 默认 |

文本：`title` ← 图层 {CardTitle}

参数：`modifier: Modifier = Modifier`、`title: String? = null`、`content: ColumnScope.() -> Unit`

示例：`DsCard(title = "商品信息") { … }`

- 约束：Card 已内置 16 内边距和 12 圆角，页面不要再给 Card 加 padding / 背景 / 圆角
- 约束：卡片之间的间距用 space.md (12)，由页面容器的 gap 提供

### Cell → `DsListItem`

单元格：标题/副标题 + 右侧值 + 箭头。源码：`android/ds/src/main/kotlin/com/demo/ds/DsListItem.kt`。

Key/value list row (Figma: `Cell`, variants `Arrow=On|Off`, `Label=On|Off`). Vertical padding 12, no horizontal padding (provided by the Card), min height 44. Left: title (font.body, text-primary) + optional subtitle (font.caption, text-tertiary, 2dp gap). Right: trailing value (font.body, text-secondary) + optional 16×16 chevron (4dp gap).

```
import com.demo.ds.DsListItem
```

| Figma 变体 | 代码 |
|---|---|
| Arrow=On | `showChevron = true` |
| Arrow=Off | 默认 |
| Label=On | 默认 |
| Label=Off | 默认 |

文本：`subtitle` ← 图层 {Label}，`title` ← 图层 {Title}，`trailingText` ← 图层 {Value}

参数：`title: String`、`modifier: Modifier = Modifier`、`subtitle: String? = null`、`trailingText: String? = null`、`showChevron: Boolean = false`

示例：`DsListItem(title = "运费", trailingText = "免运费")`

- 约束：左标题 / 右值的行一律用 Cell，不要自己用 flex/Row + 两个文本拼
- 注意：不要用 LegacyCellView（XML View 存量组件）

### Divider → `DsDivider`

1px 分割线。源码：`android/ds/src/main/kotlin/com/demo/ds/DsDivider.kt`。

1dp horizontal divider in color.border, no margins (Figma: `Divider`).

```
import com.demo.ds.DsDivider
```

参数：`modifier: Modifier = Modifier`

示例：`DsDivider()`

- 约束：不要用 border-bottom / material3.Divider 代替

### NavBar → `DsTopBar`

顶部导航栏，标题居中 + 返回。源码：`android/ds/src/main/kotlin/com/demo/ds/DsTopBar.kt`。

Top navigation bar (Figma: `NavBar`). Height 44, full width, bg-card, optional 24×24 back chevron 12dp from the left, centered title (font.title).

```
import com.demo.ds.DsTopBar
```

文本：`title` ← 图层 {Title}

参数：`title: String`、`modifier: Modifier = Modifier`、`onBack: (() -> Unit)? = null`

示例：`DsTopBar(title = "订单详情", onBack = { })`

- 约束：页面顶部只放一个 NavBar；不要自己用 div/Row 拼标题栏

### NoticeBar → `DsBanner`

通知条。源码：`android/ds/src/main/kotlin/com/demo/ds/DsBanner.kt`。

Notice banner (Figma: `NoticeBar`). Height 36, horizontal padding 12, radius.md (8), font.caption, single line, vertically centered. Info = primary on primary-light, Warning = warning on warning-light.

```
import com.demo.ds.DsBanner
```

| Figma 变体 | 代码 |
|---|---|
| Tone=Info | `tone = DsBannerTone.Info` |
| Tone=Warning | `tone = DsBannerTone.Warning` |

文本：`text` ← 图层 {Text}

参数：`text: String`、`modifier: Modifier = Modifier`、`tone: DsBannerTone = DsBannerTone.Info`

示例：`DsBanner(text = "商家已接单", tone = DsBannerTone.Warning)`


### Price → `DsPriceText`

价格：¥ + 整数 + 小数。源码：`android/ds/src/main/kotlin/com/demo/ds/DsPriceText.kt`。

Price text (Figma: `Price`): "¥" + integer + "." + two decimals in a single baseline-aligned AnnotatedString. Integer weight 700, "¥" and decimals 600. Line height = 1.2 × integer size.

```
import com.demo.ds.DsPriceText
```

| Figma 变体 | 代码 |
|---|---|
| Size=Large | `size = DsPriceSize.Large` |
| Size=Medium | `size = DsPriceSize.Medium` |
| Size=Small | `size = DsPriceSize.Small` |
| Tone=Primary | `tone = DsPriceTone.Primary` |
| Tone=Default | `tone = DsPriceTone.Default` |

文本：`amount` ← 图层 {Integer}{Decimal}|num

参数：`amount: Double`、`modifier: Modifier = Modifier`、`size: DsPriceSize = DsPriceSize.Medium`、`tone: DsPriceTone = DsPriceTone.Primary`

示例：`DsPriceText(amount = 269.0, size = DsPriceSize.Large, tone = DsPriceTone.Primary)`

- 约束：金额一律用 Price，不要拼「¥」+ 数字文本（字号层级和基线对齐会错）
- 约束：例外：Cell 右侧的值（如「¥299.00」「-¥30.00」）在设计稿里是普通文本，用 Cell 的 value 传字符串

### ProductCard → `DsGoodsCard`

商品卡：缩略图 + 标题 + 规格 + 价格 + 数量。源码：`android/ds/src/main/kotlin/com/demo/ds/DsGoodsCard.kt`。

Goods row (Figma: `ProductCard`). 80×80 placeholder thumbnail (radius.md, bg-placeholder), 12dp gap, then an 80dp-tall column: title (font.body, ≤2 lines, ellipsis), optional spec (font.caption, text-tertiary, 4dp below), and a bottom row pinned to the bottom with `DsPriceText(Medium, Default)` on the left and "x{count}" (font.caption, text-tertiary) on the right.

```
import com.demo.ds.DsGoodsCard
```

文本：`count` ← 图层 {Count}|int，`price` ← 图层 {Integer}{Decimal}|num，`spec` ← 图层 {Spec}，`title` ← 图层 {Title}

参数：`title: String`、`price: Double`、`modifier: Modifier = Modifier`、`spec: String? = null`、`count: Int = 1`

示例：`DsGoodsCard(title = "…", price = 299.0, spec = "米白色 / M", count = 1)`

- 约束：商品行一律用 ProductCard，内部的缩略图 / 标题 / 规格 / 价格 / 数量不要拆开自己拼

### Tag → `DsBadge`

状态标签。源码：`android/ds/src/main/kotlin/com/demo/ds/DsBadge.kt`。

Status badge (Figma: `Tag`). Height 20, horizontal padding 6, radius.sm (4), font 12 / 20 / 500. Only the Light look is implemented: background `{tone}-light`, text `{tone}`. (Figma's `Variant=Solid` has no Android counterpart yet.)

```
import com.demo.ds.DsBadge
```

| Figma 变体 | 代码 |
|---|---|
| Tone=Primary | `tone = DsBadgeTone.Primary` |
| Tone=Success | `tone = DsBadgeTone.Success` |
| Tone=Warning | `tone = DsBadgeTone.Warning` |
| Tone=Danger | `tone = DsBadgeTone.Danger` |
| Variant=Light | 默认 |
| Variant=Solid | **不支持** |

文本：`text` ← 图层 {Label}

参数：`text: String`、`modifier: Modifier = Modifier`、`tone: DsBadgeTone = DsBadgeTone.Primary`

示例：`DsBadge(text = "待发货", tone = DsBadgeTone.Warning)`

- 约束：状态色只能通过 type/tone 表达，禁止自己写背景色和文字色
- 注意：DsBadge 只有 Light 变体；Solid 需求已登记给组件库 owner

## 不对应 Figma 组件的辅助组件

- `DsText` Public text helper for screens (pages may only use `com.demo.ds.*`). Defaults to font.body / text-primary. Intentionally NOT tagged with [dsComponent]: it is a typography primitive, not a Figma component. — `android/ds/src/main/kotlin/com/demo/ds/DsText.kt`
- `LegacyButton` Legacy button (Figma: `Button`, `Type=Primary|Default|Text`, `Size=Large|Medium|Small`). - primary: bg ds_color_primary, white text - default: bg ds_color_bg_card, 1dp ds_color_border stroke, ds_color_text_primary text - text: no background, ds_color_primary text - large 44dp / pad 24 / 16sp bold / r22; medium 32dp / pad 16 / 14sp / r16; small 24dp / pad 12 / 12sp / r12 <com.demo.dsview.LegacyButton android:layout_width="wrap_content" android:layout_height="wrap_content" android:text="去支付" app:btnType="primary" app:btnSize="medium" />（XML View 存量组件，新页面不要用） — `android/ds-view/src/main/java/com/demo/dsview/LegacyButton.kt`
- `LegacyCellView` Legacy key/value row (Figma: `Cell`). Vertical padding 12dp, min height 44dp, title on the left (14sp text_primary), value on the right (14sp text_secondary), optional 16dp chevron 4dp after the value. The Figma `Label` (subtitle) variant is not supported by this legacy view. <com.demo.dsview.LegacyCellView android:layout_width="match_parent" android:layout_height="wrap_content" app:cellTitle="订单编号" app:cellValue="DD20261018093021" app:cellShowArrow="false" />（XML View 存量组件，新页面不要用） — `android/ds-view/src/main/java/com/demo/dsview/LegacyCellView.kt`
- `LegacyTagView` Legacy status tag (Figma: `Tag`, Light variant only). Height 20dp, horizontal padding 6dp, radius 4dp, 12sp medium text. Background `{tone}_light`, text `{tone}`. <com.demo.dsview.LegacyTagView android:layout_width="wrap_content" android:layout_height="20dp" android:text="待付款" app:tagTone="primary" />（XML View 存量组件，新页面不要用） — `android/ds-view/src/main/java/com/demo/dsview/LegacyTagView.kt`

## 待决项（Figma 有、代码没有）

- **CountDown**（all）：两端暂无组件。本期按普通文本实现（font.caption + color.danger，静态文案），已登记为组件需求；不要在页面里自造倒计时组件

## Token（页面级文本 / 容器只能用这些）

| Figma 变量 | 代码 |
|---|---|
| color.bg-card #FFFFFF | `DsColors.BgCard` |
| color.bg-page #F5F6F7 | `DsColors.BgPage` |
| color.bg-placeholder #F2F3F5 | `DsColors.BgPlaceholder` |
| color.border #E5E6EB | `DsColors.Border` |
| color.danger #F53F3F | `DsColors.Danger` |
| color.danger-light #FFECE8 | `DsColors.DangerLight` |
| color.primary #FF4D1A | `DsColors.Primary` |
| color.primary-light #FFF1EB | `DsColors.PrimaryLight` |
| color.success #00B578 | `DsColors.Success` |
| color.success-light #E6F8F1 | `DsColors.SuccessLight` |
| color.text-primary #1F2329 | `DsColors.TextPrimary` |
| color.text-secondary #646A73 | `DsColors.TextSecondary` |
| color.text-tertiary #8F959E | `DsColors.TextTertiary` |
| color.warning #FF8F1F | `DsColors.Warning` |
| color.warning-light #FFF3E8 | `DsColors.WarningLight` |
| font.body（14/20/400）| `DsTypography.Body` |
| font.caption（12/18/400）| `DsTypography.Caption` |
| font.subtitle（15/22/600）| `DsTypography.Subtitle` |
| font.title（17/24/600）| `DsTypography.Title` |
| space.xs（4）| `DsSpacing.Xs` |
| space.sm（8）| `DsSpacing.Sm` |
| space.md（12）| `DsSpacing.Md` |
| space.lg（16）| `DsSpacing.Lg` |
| space.xl（24）| `DsSpacing.Xl` |

页面级文本用 `DsText(text, style = DsTypography.X, color = DsColors.Y)`；挂 Figma id 用 `Modifier.figmaId("1:23")`。
