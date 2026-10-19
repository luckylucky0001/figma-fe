# 真实存量项目实验：拿到新设计稿「订单详情」，在开源项目里实现

A = 只给任务和 Figma 导出；C = 加上 Skill（`reuse` 体检 → `match` 复用清单 → 实现 → 页面复用检查）。所有 run 用同一版 `figmafe` 打分，由 `experiments/real-summary.py` 生成。

## newbee-mall-vue3-app（Vue 3 + Vant）

新建 `src/views/OrderDetailV2.vue`；打分：`figmafe reuse <项目> --focus src/views/OrderDetailV2.vue`

| run | 组 | 结果 | 复用现成组件 | 自造轮子 | 有组件没用 | 手写已有组件（确定/疑似） | 新增颜色变量 | 写死颜色 | 新建组件 |
|---|---|---|---|---|---|---|---|---|---|
| Rnb-A-1 | A 只给任务（Opus） | ❌ FAIL | 6 | 0 | 0 | 2 / 1 | 10 | 2 | — |
| Rnb-A-2 | A 只给任务（Opus） | ❌ FAIL | 6 | 0 | 0 | 3 / 1 | 10 | 7 | — |
| Rnb-Ahk-1 | A 只给任务（Haiku） | ❌ FAIL | 0 | 0 | 4 | 5 / 5 | 0 | 28 | — |
| Rnb-C-1 | C 用 Skill（Opus） | ✅ PASS | 12 | 0 | 0 | 0 / 0 | 0 | 0 | Price.vue |
| Rnb-C-2 | C 用 Skill（Opus） | ✅ PASS | 12 | 0 | 0 | 0 / 0 | 0 | 0 | Price.vue |
| Rnb-Chk-1 | C 用 Skill（Haiku） | ✅ PASS | 10 | 0 | 0 | 0 / 0 | 0 | 26 | Price.vue |

## QMUI_Android（XML View + QMUI）

新建 Fragment + layout；打分：`figmafe reuse <项目> --changed HEAD`（新文件整份判，已有文件只看新增行）

| run | 组 | 结果 | 复用现成组件 | 自造轮子 | 有组件没用 | 手写已有组件（确定/疑似） | 新增颜色变量 | 写死颜色 | 新建组件 |
|---|---|---|---|---|---|---|---|---|---|
| Rqm-A-1 | A 只给任务（Opus） | ❌ FAIL | 7 | 0 | 0 | 1 / 0 | 0 | 57 | — |
| Rqm-A-2 | A 只给任务（Opus） | ❌ FAIL | 8 | 0 | 0 | 1 / 0 | 11 | 4 | — |
| Rqm-C-1 | C 用 Skill（Opus） | ✅ PASS | 9 | 0 | 0 | 0 / 0 | 10 | 0 | QDNoticeBarView、QDCountDownView、QDPriceView、QDActionBar |
| Rqm-C-2 | C 用 Skill（Opus） | ✅ PASS | 8 | 0 | 0 | 0 / 0 | 8 | 0 | QDNoticeBarView、QDCountDownView、QDPriceView、QDActionBarLayout、QDTagView、QDProductCardView |

## 说明

- newbee：A 组页面里手写了商品卡、分割线、底部操作栏，还自建了 10 个颜色变量；C 组直接用 Vant 的 Card / Divider / ActionBar 和项目已有的 SimpleHeader，缺的 Price 放进组件目录。
- Haiku 不带 Skill 时一个现成组件都没用（0 个），带 Skill 后复用 10 个、检查通过：Skill 把「先查再写」变成了固定流程，模型能力弱时收益更大。
- QMUI：A 组地址行用 LinearLayout 手写（项目里有 QMUICommonListItemView，已用 48 次），C 组用 QMUIGroupListView + QMUICommonListItemView。
- QMUI 的两个 C run 都发现 match 推荐了项目用不了的组件（MaterialDivider 需要 Material 1.5，项目是 1.4；Chip / MaterialCardView 在 AppCompat 主题下会崩溃），自行绕开并写进了报告。这两类问题已经做进工具：现在体检会列出「组件库里有、但本项目不能直接用的组件」，match 不再推荐它们。
- 新增颜色变量：设计稿里的颜色项目里大多没有同值 token，C 组把它们集中加到 colors.xml；检查只把「与已有 token 同值却另起名字」的算问题。
