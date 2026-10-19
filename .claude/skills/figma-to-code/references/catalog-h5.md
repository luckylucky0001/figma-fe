# 组件手册：H5（React）

> 由 `figmafe catalog` 根据资产扫描结果（inventory）和 component-map 自动生成，不要手改。
> 规则：Figma 里出现的组件，**必须**用下表对应的代码组件；表里没有的才允许自己写，并在方案里列为待决项。

## 速查

| Figma 组件 | 代码组件 | 业务代码中已用次数 |
|---|---|---|
| ActionBar | `ActionBar` | 1 |
| Button | `Button` | 15 |
| Card | `Card` | 7 |
| Cell | `Cell` | 5 |
| CountDown | **无** —— 见待决项 | - |
| Divider | `Divider` | 6 |
| NavBar | `NavBar` | 2 |
| NoticeBar | `NoticeBar` | 3 |
| Price | `Price` | 4 |
| ProductCard | `ProductCard` | 3 |
| Tag | `Tag` | 3 |

## 组件详情

### ActionBar → `ActionBar`

吸底操作栏，按钮右对齐。源码：`h5/src/ds/ActionBar/index.tsx`。

Usually Buttons; right-aligned with an 8px gap. */ children?: ReactNode; } /** ActionBar — bar fixed to the bottom of the viewport, 56px high, white background, 1px top border, 16px horizontal padding; children are right-aligned with an 8px gap. Pages should reserve 56px of bottom padding so content is not covered.

```
import { ActionBar } from '@ds';
```

参数：`children: ReactNode`

示例：`<ActionBar><Button>联系客服</Button><Button type="primary">提醒发货</Button></ActionBar>`

- 约束：吸底操作区一律用 ActionBar；页面内容区底部要留出 56 的空间

### Button → `Button`

按钮。源码：`h5/src/ds/Button/index.tsx`。

Visual type. Default `'default'`. */ type?: 'primary' | 'default' | 'text'; /** Large h44 · Medium h32 · Small h24. Default `'medium'`. */ size?: 'large' | 'medium' | 'small'; /** Stretch to the full width of the container. */ block?: boolean; onClick?: MouseEventHandler<HTMLButtonElement>; children?: ReactNode; } /** Button — pill button. Primary (brand fill, white text), Default (white fill, 1px border), Text (no fill/border, brand text). Sizes Large 44 / Medium 32 / Small 24 px high.

```
import { Button } from '@ds';
```

| Figma 变体 | 代码 |
|---|---|
| Size=Large | `size="large"` |
| Size=Medium | `size="medium"` |
| Size=Small | `size="small"` |
| Type=Primary | `type="primary"` |
| Type=Default | `type="default"` |
| Type=Text | `type="text"` |

文本：`children` ← 图层 {Label}

参数：`type: 'primary' | 'default' | 'text'`、`size: 'large' | 'medium' | 'small'`、`block: boolean`、`onClick: MouseEventHandler<HTMLButtonElement>`、`children: ReactNode`

示例：`<Button type="primary" size="medium">提醒发货</Button>`

- 约束：一屏最多一个 Primary 按钮
- 约束：H5 禁止原生 <button>；Android 禁止 material3.Button / OutlinedButton

### Card → `Card`

白底圆角卡片，可带标题。源码：`h5/src/ds/Card/index.tsx`。

Optional card title (font.subtitle), 12px above the content. */ title?: string; children?: ReactNode; } /** Card — white container with 12px radius and 16px padding. Children stack vertically with no gap; an optional title sits 12px above them.

```
import { Card } from '@ds';
```

| Figma 变体 | 代码 |
|---|---|
| Title=On | 默认 |
| Title=Off | 默认 |

文本：`title` ← 图层 {CardTitle}

参数：`title: string`、`children: ReactNode`

示例：`<Card title="商品信息">…</Card>`

- 约束：Card 已内置 16 内边距和 12 圆角，页面不要再给 Card 加 padding / 背景 / 圆角
- 约束：卡片之间的间距用 space.md (12)，由页面容器的 gap 提供

### Cell → `Cell`

单元格：标题/副标题 + 右侧值 + 箭头。源码：`h5/src/ds/Cell/index.tsx`。

Left title (font.body). */ title: string; /** Optional sub-title under the title (font.caption, tertiary). */ label?: string; /** Right-aligned value (font.body, secondary). */ value?: ReactNode; /** Show a 16×16 right chevron after the value. */ arrow?: boolean; } /** Cell — list row: title (+ optional label) on the left, value (+ optional arrow) on the right. 12px vertical padding, no horizontal padding (the Card provides it), min height 44.

```
import { Cell } from '@ds';
```

| Figma 变体 | 代码 |
|---|---|
| Arrow=On | `arrow` |
| Arrow=Off | 默认 |
| Label=On | 默认 |
| Label=Off | 默认 |

文本：`label` ← 图层 {Label}，`title` ← 图层 {Title}，`value` ← 图层 {Value}

参数：`title: string`、`label: string`、`value: ReactNode`、`arrow: boolean`

示例：`<Cell title="运费" value="免运费" />`

- 约束：左标题 / 右值的行一律用 Cell，不要自己用 flex/Row + 两个文本拼

### Divider → `Divider`

1px 分割线。源码：`h5/src/ds/Divider/index.tsx`。

Divider has no own props; only passthrough attributes (data-*, className, style). */ // eslint-disable-next-line @typescript-eslint/no-empty-interface export interface DividerProps extends Omit<HTMLAttributes<HTMLDivElement>, 'children'> {} /** Divider — 1px horizontal line in color.border, no margin.

```
import { Divider } from '@ds';
```

示例：`<Divider />`

- 约束：不要用 border-bottom / material3.Divider 代替

### NavBar → `NavBar`

顶部导航栏，标题居中 + 返回。源码：`h5/src/ds/NavBar/index.tsx`。

Centered title text (font.title). */ title: string; /** Called when the back chevron is tapped. */ onBack?: () => void; } /** NavBar — top navigation bar, 44px high, white background, back chevron on the left (24×24, 12px from the left edge) and a centered title.

```
import { NavBar } from '@ds';
```

文本：`title` ← 图层 {Title}

参数：`title: string`、`onBack: () => void`

示例：`<NavBar title="订单详情" onBack={() => history.back()} />`

- 约束：页面顶部只放一个 NavBar；不要自己用 div/Row 拼标题栏

### NoticeBar → `NoticeBar`

通知条。源码：`h5/src/ds/NoticeBar/index.tsx`。

`'info'` = primary-light bg / primary text; `'warning'` = warning-light bg / warning text. Default `'info'`. */ tone?: 'info' | 'warning'; children?: ReactNode; } /** NoticeBar — single-line notice, 36px high, 12px horizontal padding, 8px radius, font.caption. Overflowing text is truncated with an ellipsis.

```
import { NoticeBar } from '@ds';
```

| Figma 变体 | 代码 |
|---|---|
| Tone=Info | `tone="info"` |
| Tone=Warning | `tone="warning"` |

文本：`children` ← 图层 {Text}

参数：`tone: 'info' | 'warning'`、`children: ReactNode`

示例：`<NoticeBar tone="warning">商家已接单</NoticeBar>`


### Price → `Price`

价格：¥ + 整数 + 小数。源码：`h5/src/ds/Price/index.tsx`。

Amount in yuan, e.g. `128.5` renders as ¥128.50. */ value: number; /** Large: ¥14 / 22 / 14 · Medium: ¥12 / 18 / 12 · Small: ¥12 / 14 / 12. Default `'medium'`. */ size?: 'large' | 'medium' | 'small'; /** `'primary'` uses color.primary, `'default'` uses text-primary. Default `'primary'`. */ tone?: 'primary' | 'default'; } /** Price — "¥" + integer + "." + two decimals, baseline aligned. Integer is 700 weight, symbol and decimals are 600.

```
import { Price } from '@ds';
```

| Figma 变体 | 代码 |
|---|---|
| Size=Large | `size="large"` |
| Size=Medium | `size="medium"` |
| Size=Small | `size="small"` |
| Tone=Primary | `tone="primary"` |
| Tone=Default | `tone="default"` |

文本：`value` ← 图层 {Integer}{Decimal}|num

参数：`value: number`、`size: 'large' | 'medium' | 'small'`、`tone: 'primary' | 'default'`

示例：`<Price value={269} size="large" tone="primary" />`

- 约束：金额一律用 Price，不要拼「¥」+ 数字文本（字号层级和基线对齐会错）
- 约束：例外：Cell 右侧的值（如「¥299.00」「-¥30.00」）在设计稿里是普通文本，用 Cell 的 value 传字符串

### ProductCard → `ProductCard`

商品卡：缩略图 + 标题 + 规格 + 价格 + 数量。源码：`h5/src/ds/ProductCard/index.tsx`。

Product name, clamped to two lines. */ title: string; /** Spec / SKU description (font.caption, tertiary). */ spec?: string; /** Unit price, rendered as `Price size="medium" tone="default"`. */ price: number; /** Quantity, rendered as "x{count}". */ count?: number; /** Thumbnail URL; a placeholder background shows when omitted. */ thumb?: string; } /** ProductCard — 80×80 thumbnail + info column (title, spec, bottom row with price and quantity).

```
import { ProductCard } from '@ds';
```

文本：`count` ← 图层 {Count}|num，`price` ← 图层 {Integer}{Decimal}|num，`spec` ← 图层 {Spec}，`title` ← 图层 {Title}

参数：`title: string`、`spec: string`、`price: number`、`count: number`、`thumb: string`

示例：`<ProductCard title="…" spec="米白色 / M" price={299} count={1} />`

- 约束：商品行一律用 ProductCard，内部的缩略图 / 标题 / 规格 / 价格 / 数量不要拆开自己拼

### Tag → `Tag`

状态标签。源码：`h5/src/ds/Tag/index.tsx`。

Tone. Default `'primary'`. */ type?: 'primary' | 'success' | 'warning' | 'danger'; /** `true` = Light variant (tinted bg, toned text); `false` = Solid (toned bg, white text). Default `true`. */ light?: boolean; children?: ReactNode; } /** Tag — small status label, 20px high, 6px horizontal padding, 4px radius, 12/20/500 text.

```
import { Tag } from '@ds';
```

| Figma 变体 | 代码 |
|---|---|
| Tone=Primary | `type="primary"` |
| Tone=Success | `type="success"` |
| Tone=Warning | `type="warning"` |
| Tone=Danger | `type="danger"` |
| Variant=Light | 默认 |
| Variant=Solid | `light={false}` |

文本：`children` ← 图层 {Label}

参数：`type: 'primary' | 'success' | 'warning' | 'danger'`、`light: boolean`、`children: ReactNode`

示例：`<Tag type="warning">待发货</Tag>`

- 约束：状态色只能通过 type/tone 表达，禁止自己写背景色和文字色

## 不对应 Figma 组件的辅助组件


## 待决项（Figma 有、代码没有）

- **CountDown**（all）：两端暂无组件。本期按普通文本实现（font.caption + color.danger，静态文案），已登记为组件需求；不要在页面里自造倒计时组件

## Token（页面级文本 / 容器只能用这些）

| Figma 变量 | 代码 |
|---|---|
| color.bg-card #FFFFFF | `var(--ds-color-bg-card)` |
| color.bg-page #F5F6F7 | `var(--ds-color-bg-page)` |
| color.bg-placeholder #F2F3F5 | `var(--ds-color-bg-placeholder)` |
| color.border #E5E6EB | `var(--ds-color-border)` |
| color.danger #F53F3F | `var(--ds-color-danger)` |
| color.danger-light #FFECE8 | `var(--ds-color-danger-light)` |
| color.primary #FF4D1A | `var(--ds-color-primary)` |
| color.primary-light #FFF1EB | `var(--ds-color-primary-light)` |
| color.success #00B578 | `var(--ds-color-success)` |
| color.success-light #E6F8F1 | `var(--ds-color-success-light)` |
| color.text-primary #1F2329 | `var(--ds-color-text-primary)` |
| color.text-secondary #646A73 | `var(--ds-color-text-secondary)` |
| color.text-tertiary #8F959E | `var(--ds-color-text-tertiary)` |
| color.warning #FF8F1F | `var(--ds-color-warning)` |
| color.warning-light #FFF3E8 | `var(--ds-color-warning-light)` |
| font.body（14/20/400）| `var(--ds-font-body-size) / var(--ds-font-body-line-height) / var(--ds-font-body-weight)` |
| font.caption（12/18/400）| `var(--ds-font-caption-size) / var(--ds-font-caption-line-height) / var(--ds-font-caption-weight)` |
| font.subtitle（15/22/600）| `var(--ds-font-subtitle-size) / var(--ds-font-subtitle-line-height) / var(--ds-font-subtitle-weight)` |
| font.title（17/24/600）| `var(--ds-font-title-size) / var(--ds-font-title-line-height) / var(--ds-font-title-weight)` |
| space.xs（4）| `var(--ds-space-xs)` |
| space.sm（8）| `var(--ds-space-sm)` |
| space.md（12）| `var(--ds-space-md)` |
| space.lg（16）| `var(--ds-space-lg)` |
| space.xl（24）| `var(--ds-space-xl)` |

页面级文本用普通元素 + CSS 变量；组件和关键节点上挂 `data-figma-id="1:23"`（组件会透传到根节点）。
