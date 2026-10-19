# 真实存量项目实验：拿到新设计稿「订单详情」，在开源项目里实现

A = 只给任务和 Figma 导出；C = 加上 Skill（`reuse` 体检 → `match` 复用清单 → 实现 → 页面复用检查）。所有 run 用同一版 `figmafe` 打分，由 `experiments/real-summary.py` 生成。

## newbee-mall-vue3-app（Vue 3 + Vant）

新建 `src/views/OrderDetailV2.vue`；打分：`figmafe reuse <项目> --focus src/views/OrderDetailV2.vue`

| run | 组 | 结果 | 复用现成组件 | 自造轮子（确定/疑似） | 有组件没用 | 手写已有组件（确定/疑似） | 新增颜色变量 | 写死颜色 | 新建组件 |
|---|---|---|---|---|---|---|---|---|---|
| Rnb-A-1 | A 只给任务（Opus） | ❌ FAIL | 6 | 0 / 0 | 0 | 3 / 1 | 10 | 2 | — |
| Rnb-A-2 | A 只给任务（Opus） | ❌ FAIL | 6 | 0 / 0 | 0 | 3 / 1 | 10 | 7 | — |
| Rnb-Ahk-1 | A 只给任务（Haiku） | ❌ FAIL | 0 | 0 / 0 | 4 | 5 / 5 | 0 | 28 | — |
| Rnb-C-1 | C 用 Skill（Opus） | ✅ PASS | 12 | 0 / 0 | 0 | 0 / 0 | 0 | 0 | Price.vue |
| Rnb-C-2 | C 用 Skill（Opus） | ✅ PASS | 12 | 0 / 0 | 0 | 0 / 0 | 0 | 0 | Price.vue |
| Rnb-Chk-1 | C 用 Skill（Haiku） | ✅ PASS | 10 | 0 / 0 | 0 | 0 / 0 | 0 | 26 | Price.vue |

## QMUI_Android（XML View + QMUI）

新建 Fragment + layout；打分：`figmafe reuse <项目> --changed HEAD`（新文件整份判，已有文件只看新增行）

| run | 组 | 结果 | 复用现成组件 | 自造轮子（确定/疑似） | 有组件没用 | 手写已有组件（确定/疑似） | 新增颜色变量 | 写死颜色 | 新建组件 |
|---|---|---|---|---|---|---|---|---|---|
| Rqm-A-1 | A 只给任务（Opus） | ❌ FAIL | 8 | 0 / 0 | 0 | 1 / 0 | 0 | 57 | — |
| Rqm-A-2 | A 只给任务（Opus） | ❌ FAIL | 9 | 0 / 0 | 0 | 1 / 0 | 11 | 4 | — |
| Rqm-Ahk-1 | A 只给任务（Haiku） | ❌ FAIL | 6 | 0 / 0 | 0 | 1 / 0 | 11 | 1 | — |
| Rqm-C-1 | C 用 Skill（Opus） | ✅ PASS | 9 | 0 / 0 | 0 | 0 / 0 | 10 | 0 | QDNoticeBarView、QDCountDownView、QDPriceView、QDActionBar |
| Rqm-C-2 | C 用 Skill（Opus） | ✅ PASS | 8 | 0 / 0 | 0 | 0 / 0 | 8 | 0 | QDNoticeBarView、QDCountDownView、QDPriceView、QDActionBarLayout、QDTagView、QDProductCardView |
| Rqm-Chk-1 | C 用 Skill（Haiku） | ✅ PASS | 8 | 0 / 0 | 0 | 0 / 0 | 10 | 0 | 无（Tag / Price / CountDown / 商品行用 QMUICommonListItemView 的标题、描述、accessory 和 span 实现，决定写进 component-map.yaml） |

## Now in Android（Jetpack Compose + Material3 + 自有设计系统 Nia*）

新建 `OrderDetailScreen.kt`，通用组件放哪由 agent 决定；打分：`figmafe reuse <项目> --changed HEAD`

| run | 组 | 结果 | 复用现成组件 | 自造轮子（确定/疑似） | 有组件没用 | 手写已有组件（确定/疑似） | 新增颜色变量 | 写死颜色 | 新建组件 |
|---|---|---|---|---|---|---|---|---|---|
| Rnia-A-1 | A 只给任务（Opus） | ❌ FAIL | 2 | 5 / 1 | 0 | 0 / 0 | 0 | 18 | 12 个：NiaNavBar、NiaCard、NiaDivider、NiaCell、NiaPillButton、NiaStatusTag、NiaPrice、NiaCountDown、NiaNoticeBar、NiaActionBar、NiaProductCard + CommerceTheme |
| Rnia-A-2 | A 只给任务（Opus） | ❌ FAIL | 2 | 5 / 1 | 0 | 0 / 0 | 0 | 18 | 11 个 Commerce*（NavBar、NoticeBar、Card、Divider、Cell、Tag、Price、ProductCard、CountDown、Button、ActionBar）+ CommerceTheme |
| Rnia-C-1 | C 用 Skill（Opus） | ❌ FAIL | 6 | 1 / 1 | 0 | 0 / 0 | 0 | 0 | NiaNoticeBar、NiaPriceText、NiaCountdownText、NiaActionBar、NiaListItem、NiaStatusTag、ProductCard |
| Rnia-C-2 | C 用 Skill（Opus） | ✅ PASS | 8 | 0 / 0 | 0 | 0 / 0 | 0 | 33 | NiaStatusTag、NiaCell、NiaNoticeBar、NiaCountDown、NiaPriceText、NiaBottomActionBar、ProductItem（每个都在 component-map.yaml 记录了决定和理由） |

## 说明

- newbee：A 组页面里手写了商品卡、分割线、底部操作栏，还自建了 10 个颜色变量；C 组直接用 Vant 的 Card / Divider / ActionBar 和项目已有的 SimpleHeader，缺的 Price 放进组件目录。
- Haiku 不带 Skill 时一个现成组件都没用（0 个），带 Skill 后复用 10 个、检查通过：Skill 把「先查再写」变成了固定流程，模型能力弱时收益更大。
- QMUI + Haiku：不带 Skill 时没用任何列表组件（FAIL）；带 Skill 时复用 QMUIGroupListView / QMUICommonListItemView ×7 / QMUIRoundButton，检查 PASS。它还指出 Skill 的 catalog / workflow 文件讲的是 demo 组件库，在别的项目里会误导，已拆分说明。
- QMUI 的 Haiku 基线（Rqm-Ahk-1）没用任何列表组件，7 行「标题 + 值」都是 LinearLayout + 两个 TextView 手写的，而且没有能说明用途的 id / class。按名字的检查看不出来，于是加了按结构识别（横排、两个文字、可选箭头）的检查，现在能抓到；两次 C 组没有误报。
- QMUI：A 组地址行用 LinearLayout 手写（项目里有 QMUICommonListItemView，已用 48 次），C 组用 QMUIGroupListView + QMUICommonListItemView。
- QMUI 的两个 C run 都发现 match 推荐了项目用不了的组件（MaterialDivider 需要 Material 1.5，项目是 1.4；Chip / MaterialCardView 在 AppCompat 主题下会崩溃），自行绕开并写进了报告。这两类问题已经做进工具：现在体检会列出「组件库里有、但本项目不能直接用的组件」，match 不再推荐它们。
- Now in Android：A 组没有先查项目已有组件，在 core/designsystem 下新建了 12 个组件，其中 NiaNavBar / NiaCard / NiaDivider / NiaPillButton / NiaCell 和已有的 NiaTopAppBar / Card / HorizontalDivider / NiaButton / ListItem 重复，另外写了一套 18 个写死颜色的主题。C 组复用了 NiaTopAppBar、Card、HorizontalDivider、NiaButton、NiaOutlinedButton，颜色全部映射到 MaterialTheme；它按语义拒绝了 match 的两个 △ 推荐（InterestsItem、NiaFilterChip），新建的 NiaListItem 被 PR 闸门判为和 Material3 ListItem 重复，属于需要人定的项，agent 在交付里如实说明了理由。
- Rnia-C-2 按 Skill 新规则把「为什么不用推荐组件 / 为什么新建」写进了 component-map.yaml，PR 闸门认可这些决定后通过；它还发现了工具的一个 bug（component-map 的 props 写成 `Type=Primary: …` 扁平格式时整份文件被静默忽略），已修复：两种写法都支持，解析失败会明确报错。
- C 组 agent 会把这次学到的经验追加到 rules.md（例如「直接复用也要核对组件语义」），这正是方案里「沉淀规则」那一步；其中有价值的已合并进 Skill。
- 新增颜色变量：设计稿里的颜色项目里大多没有同值 token，C 组把它们集中加到 colors.xml；检查只把「与已有 token 同值却另起名字」的算问题。
