# 真实存量项目实验：newbee-mall-vue3-app（Vue 3 + Vant）× 新设计稿「订单详情」

任务：按新 Figma 重做订单详情页（`src/views/OrderDetailV2.vue`）。A = 只给任务；C = Skill（`figmafe reuse` 体检 → `figmafe match` 复用清单 → 实现 → `reuse --focus` 自检）。统一用同一版 `figmafe reuse --focus` 打分。

| run | 结果 | 复用现成组件 | 自造轮子 | 有组件没用 | 手写已有组件（确定/疑似） | 页面私有调色板 | 写死颜色 | 新建组件 |
|---|---|---|---|---|---|---|---|---|
| Rnb-A-1 | FAIL | 6 | 0 | 1 | 2 / 1 | 10 个变量 | 2 | — |
| Rnb-A-2 | FAIL | 6 | 0 | 1 | 3 / 1 | 10 个变量 | 7 | — |
| Rnb-C-1 | PASS | 12 | 0 | 0 | 0 / 0 | 0 个变量 | 0 | Price.vue（组件目录） |
| Rnb-C-2 | PASS | 12 | 0 | 0 | 0 / 0 | 0 个变量 | 0 | Price.vue（组件目录） |

复用组件明细：

- **Rnb-A-1**：Button（组件库）、Cell（组件库）、CountDown（组件库）、NavBar（组件库）、NoticeBar（组件库）、Tag（组件库）
- **Rnb-A-2**：Button（组件库）、Cell（组件库）、CountDown（组件库）、NavBar（组件库）、NoticeBar（组件库）、Tag（组件库）
- **Rnb-C-1**：ActionBar（组件库）、Button（组件库）、Card（组件库）、CellGroup（组件库）、Cell（组件库）、CountDown（组件库）、Divider（组件库）、Image（组件库）、NoticeBar（组件库）、Price（业务组件）、SimpleHeader（业务组件）、Tag（组件库）
- **Rnb-C-2**：ActionBar（组件库）、Button（组件库）、Card（组件库）、CellGroup（组件库）、Cell（组件库）、CountDown（组件库）、Divider（组件库）、Image（组件库）、NoticeBar（组件库）、Price（业务组件）、SimpleHeader（业务组件）、Tag（组件库）
