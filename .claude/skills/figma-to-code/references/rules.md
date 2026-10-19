# 沉淀规则（每完成一个页面追加；格式：现象 → 规则 → 来源）

- Card 自带 16 内边距和 12 圆角，页面再包一层 padding 会让内容整体错位 → 页面不给 Card 加 padding / 背景 / 圆角 → 组件规范
- 左标题 / 右值的行很容易被写成 flex + 两个 span → 一律用 Cell（H5）/ DsListItem（Android），值通过 value / trailingText 传字符串 → 组件规范
- 金额拼文本会丢失「¥ / 整数 / 小数」三级字号 → 用 Price / DsPriceText；Cell 右侧的金额值是普通文本例外 → 组件规范
