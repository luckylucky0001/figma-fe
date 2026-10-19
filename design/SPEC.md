# 设计规范（Figma 组件库的文字版）

> 这是 Figma 组件库 `Demo DS / Mobile` 的规范文字版，H5 和 Android 组件库都按它实现。
> 所有尺寸单位：Figma px = H5 CSS px = Android dp。导出图统一为 @2x。
> 字体：`Noto Sans SC`（安装方法见 `scripts/fetch-fonts.sh`）。

## 1. Tokens（Figma Variables，DTCG 格式见 `design/figma-export/order-detail/variables.tokens.json`）

| Token | 值 | 说明 |
|---|---|---|
| color.primary | #FF4D1A | 品牌色 |
| color.primary-light | #FFF1EB | |
| color.success | #00B578 | |
| color.success-light | #E6F8F1 | |
| color.warning | #FF8F1F | |
| color.warning-light | #FFF3E8 | |
| color.danger | #F53F3F | |
| color.danger-light | #FFECE8 | |
| color.text-primary | #1F2329 | |
| color.text-secondary | #646A73 | |
| color.text-tertiary | #8F959E | |
| color.bg-page | #F5F6F7 | 页面背景 |
| color.bg-card | #FFFFFF | |
| color.bg-placeholder | #F2F3F5 | 图片占位 |
| color.border | #E5E6EB | 分割线 / 描边 |
| space.xs / sm / md / lg / xl | 4 / 8 / 12 / 16 / 24 | |
| radius.sm / md / lg / full | 4 / 8 / 12 / 999 | |
| font.title | 17 / 24 / 600 | 字号 / 行高 / 字重 |
| font.subtitle | 15 / 22 / 600 | |
| font.body | 14 / 20 / 400 | |
| font.caption | 12 / 18 / 400 | |

## 2. 组件（Figma 组件集名称 → 规格）

### NavBar
- 高 44，背景 bg-card，宽度撑满
- 左侧返回图标：24×24，距左 12，垂直居中；图标是向左的 chevron，描边 2，颜色 text-primary
- 标题：font.title，颜色 text-primary，水平居中
- 变体：无

### Card
- 背景 bg-card，圆角 radius.lg (12)，内边距 16
- 可选标题：font.subtitle，text-primary；标题与内容间距 12
- 子元素纵向排列，默认无间距
- 变体：`Title=On|Off`

### Cell
- 纵向内边距 12，横向内边距 0（由 Card 提供），最小高度 44
- 左：标题 font.body text-primary；可选副标题（label）font.caption text-tertiary，与标题间距 2
- 右：值 font.body text-secondary，右对齐；可选箭头 16×16 向右 chevron（描边 1.5，text-tertiary），与值间距 4
- 左右区域垂直居中
- 变体：`Arrow=On|Off`、`Label=On|Off`

### Tag
- 高 20，横向内边距 6，圆角 radius.sm (4)，font 12 / 20 / 500
- 变体：`Tone=Primary|Success|Warning|Danger`，`Variant=Light|Solid`
  - Light：背景 {tone}-light，文字 {tone}
  - Solid：背景 {tone}，文字 #FFFFFF

### Price
- 结构：「¥」+ 整数部分 + 「.」+ 两位小数，基线对齐，字重 700（¥ 和小数 600）
- 变体：`Size=Large|Medium|Small`
  - Large：¥ 14，整数 22，小数 14
  - Medium：¥ 12，整数 18，小数 12
  - Small：¥ 12，整数 14，小数 12
- 变体：`Tone=Primary|Default`：Primary 用 color.primary，Default 用 text-primary

### Button
- 变体：`Type=Primary|Default|Text`，`Size=Large|Medium|Small`
  - Primary：背景 primary，文字 #FFFFFF
  - Default：背景 bg-card，1px 描边 border，文字 text-primary
  - Text：无背景无描边，文字 primary
  - Large：高 44，横向内边距 24，字 16 / 600，圆角 22
  - Medium：高 32，横向内边距 16，字 14 / 500，圆角 16
  - Small：高 24，横向内边距 12，字 12 / 500，圆角 12
- 文字水平垂直居中

### Divider
- 1px 横线，颜色 border，无外边距
- 变体：无

### ProductCard
- 横向排列：左缩略图 80×80，圆角 radius.md (8)，背景 bg-placeholder；间距 12
- 右侧纵向，高 80：
  - 标题 font.body text-primary，最多两行，超出省略
  - 规格 font.caption text-tertiary，与标题间距 4
  - 底部一行（贴底）：左 Price `Size=Medium Tone=Default`，右「x{数量}」font.caption text-tertiary
- 变体：无

### ActionBar
- 吸底，高 56，背景 bg-card，顶部 1px 描边 border
- 横向内边距 16；内容（按钮）右对齐，按钮间距 8，垂直居中
- 变体：无

### NoticeBar
- 高 36，横向内边距 12，圆角 radius.md (8)，font.caption
- 变体：`Tone=Info|Warning`：Info 背景 primary-light 文字 primary；Warning 背景 warning-light 文字 warning
- 文字单行，垂直居中

### CountDown（只有 Figma 有，两端都还没有实现，用来验证缺口报告）
- 文字「剩余 HH:MM:SS」，font.caption，颜色 danger
