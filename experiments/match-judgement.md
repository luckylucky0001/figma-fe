# figmafe match 准确率

## 1. demo 仓库（有人工确认过的标准答案，自动评估）

运行 `python3 experiments/match-accuracy.py`：H5 **11/11**，Android（`--target compose`）**11/11**。

> 注意：demo 仓库同时被用来校准打分规则，结果偏乐观。下面两个真实开源项目的结果更有参考价值。

## 2. 真实项目（人工逐条判断）

Figma：订单详情（11 个组件）。判断标准：推荐的组件是否就是这个项目里实现该元素时应该用的那个。

### newbee-mall-vue3-app（Vue 3 + Vant）

| Figma | 推荐 | 结论 | 人工判断 |
|---|---|---|---|
| NavBar | SimpleHeader（业务组件，项目里已用 11 次） | 直接复用 | ✓ 项目里的顶部导航就是它。名字叫 NavBar 的那个组件其实是底部标签栏，工具已经通过样式判断出来，没有选它 |
| NoticeBar | Vant NoticeBar | 直接复用 | ✓ |
| Card | Vant CellGroup（inset） | 需确认 | ✓ Vant 里没有通用卡片，用 inset 圆角分组是常见做法。标为需确认是合理的 |
| Tag | Vant Tag | 直接复用 | ✓ |
| CountDown | Vant CountDown | 直接复用 | ✓ 正是「不造轮子」的典型例子 |
| Divider | Vant Divider | 直接复用 | ✓ |
| Cell | Vant Cell | 直接复用 | ✓ |
| ProductCard | Vant Card（项目里已用 2 次） | 直接复用 | ✓ van-card 就是商品卡 |
| Price | — | 需新建 | ✓ Vant 和项目里都没有独立的价格组件 |
| ActionBar | Vant ActionBar（项目里已用 1 次） | 直接复用 | ✓ |
| Button | Vant Button（项目里已用 14 次） | 直接复用 | ✓ |

**11/11**

### QMUI_Android（XML View，未传 `--target`，工具按项目主流写法推断为 xml）

| Figma | 推荐 | 结论 | 人工判断 |
|---|---|---|---|
| NavBar | QMUITopBarLayout（View，项目里已用 45 次） | 直接复用 | ✓ 两个用 Skill 的 agent 最后都选了它 |
| NoticeBar | — | 需新建 | ✓ QMUI 没有通知条组件 |
| Card | QMUIGroupListView（项目里已用 25 次） | 直接复用 | ✓ QMUI 惯用分组列表来做白底分区 |
| Tag | — | 需新建 | ✓ 组件库里的 Chip 在本项目的 AppCompat 主题下会崩溃，工具已排除它 |
| CountDown | — | 需新建 | ✓ |
| Divider | QMUIPriorityLinearLayout 自带的分割线属性（qmui_bottomDividerHeight 等） | 需确认 | △ 思路对（两个 agent 都用了 QMUI 容器自带的分割线），但更常用的是 QMUILinearLayout，推荐的容器偏冷门 |
| Cell | QMUICommonListItemView（项目里已用 48 次） | 直接复用 | ✓ |
| ProductCard | — | 需新建 | ✓ MaterialCardView 同样有主题崩溃问题，已排除 |
| Price | — | 需新建 | ✓ |
| ActionBar | — | 需新建 | ✓ |
| Button | QMUIRoundButton（项目里已用 33 次） | 直接复用 | ✓ |

**10/11 完全正确，1/11 可用但不是最优。** 上一版是 8/11 + 3 个「可用」，其中 MaterialDivider（项目的 Material 1.4.0 里没有）和 Chip / MaterialCardView（AppCompat 主题下会崩溃）其实是错的推荐，是用 Skill 的 agent 实现时发现并反馈的。现在体检会检查组件库版本和 App 主题，报告里列出「组件库里有、但本项目不能直接用的组件」，match 不再推荐它们。

### Now in Android（Jetpack Compose + Material3 + 自有设计系统 `Nia*`，Google 官方示例）

| Figma | 推荐 | 结论 | 人工判断 |
|---|---|---|---|
| NavBar | NiaTopAppBar（设计系统） | 直接复用 | ✓ |
| NoticeBar | — | 需新建 | ✓ |
| Card | Material3 Card（项目里直接用了 2 次，设计系统没有封装） | 直接复用 | ✓ |
| Tag | NiaFilterChip（设计系统） | 直接复用 | △ 可用。NiaFilterChip 是可选中的筛选标签，纯展示的状态标签用 NiaTopicTag 也说得通 |
| CountDown | — | 需新建 | ✓ |
| Divider | Material3 HorizontalDivider（项目里已用 2 次） | 直接复用 | ✓ 已废弃的 `Divider` 被排除 |
| Cell | InterestsItem（core/ui，项目里已用 7 次） | 直接复用 | △ 它是「兴趣」专用的列表项，通用列表项用 Material3 ListItem 更合适 |
| ProductCard | Material3 Card | 需确认 | ✓ |
| Price | — | 需新建 | ✓ |
| ActionBar | — | 需新建 | ✓ |
| Button | NiaButton（设计系统） | 直接复用 | ✓ 同一家族里有 NiaButton / NiaTextButton / NiaOutlinedButton，选了基础款 |

**9/11 完全正确，2/11 可用但不是最优。** 体检还找出了设计系统的全部 8 个封装关系（Button → NiaButton、TabRow → NiaTabRow 等），并确认业务代码里没有绕过封装直接用 Material3 的地方（用 grep 核对过）。

### React + antd-mobile v5（两个开源移动端模板：lhw-dev/react-h5-template、xmllein/react-mobile）

两个项目页面都不多，推荐主要来自 antd-mobile 组件目录，结论一致：

| Figma | 推荐 | 结论 | 人工判断 |
|---|---|---|---|
| NavBar / NoticeBar / Card / Tag / Divider / Button | antd-mobile 同名组件 | 直接复用 | ✓ |
| Cell | antd-mobile `List.Item` | 用过的项目：直接复用；没用过的：需确认 | ✓ antd-mobile 的列表项写作 `<List.Item>`，工具能识别这种点号写法，示例也给的是 `<List.Item>` |
| ProductCard | antd-mobile Card | 需确认 | ✓ 只能当容器 |
| CountDown / Price / ActionBar | — | 需新建 | ✓ antd-mobile 没有这三个组件 |

**每个项目 11/11。** 这两个是模板项目，代码量小，主要验证的是 React + antd-mobile 这条链路能跑通，以及组件目录是否正确。

### 没组件化、图层是中文名的设计稿（`design/figma-export/holdings`，券商「我的持仓」页）

设计稿里全是 FRAME / TEXT，没有组件实例，图层名是「导航栏」「涨跌标签」「收盘倒计时」这样的中文。工具按中文名推断 UI 概念再匹配，结果标「（未组件化）」：

| 项目 | 12 个元素的推荐 | 人工判断 |
|---|---|---|
| newbee（Vue + Vant） | 导航栏 → SimpleHeader、标签页 → Tabs、持仓列表项 → Cell、涨跌标签 → Tag、收盘倒计时 → CountDown、底部操作栏 → ActionBar、买入 / 卖出 / 更多 → Button … | 12/12 ✓ |
| react-h5-template（React + antd-mobile） | 导航栏 → NavBar、标签页 → Tabs、持仓列表项 → List.Item、收盘倒计时 / 底部操作栏 → 需新建 … | 12/12 ✓ |
| QMUI（XML） | 导航栏 → QMUITopBarLayout、标签页 → QMUITabSegment、持仓列表 → RecyclerView、持仓列表项 → QMUICommonListItemView … | 12/12 ✓ |
| Now in Android（Compose） | 导航栏 → NiaTopAppBar、资产卡片 → Card、标签页 → NiaTabRow、按钮 → NiaButton … | 9/12 ✓，3 △：持仓列表判成需新建（Compose 里列表是 LazyColumn，不算组件）、持仓列表项推荐了兴趣专用的 InterestsItem、涨跌标签推荐了可选中的 NiaFilterChip |

### 小结

- 订单详情（组件化的设计稿）：五个真实项目（覆盖 Vue+Vant、React+antd-mobile、XML+QMUI、Compose+Material3+自有设计系统）共 55 项，完全正确 52 项，可用但不是最优 3 项。
- 我的持仓（没组件化、中文图层名）：四个真实项目共 48 项，完全正确 45 项，可用但不是最优 3 项。
- 两份设计稿都没有出现推荐错误组件的情况。拿不准的情况，工具给出的都是「需确认」。
- 存量项目里「组件库有这个组件」不等于「项目能用」。版本、主题、废弃 API、设计系统封装都会影响结论，这几项现在都会检查。
- 项目惯例：没有对应组件、但项目里反复出现的写法（1px View 细线、Vant 的 `van-hairline--*` 类、容器自带的分割线属性）会作为「项目惯例」推荐，不会直接判成「需新建」。
- 还没覆盖的：「借用」写法（例如用按钮组件来充当标签）。这类可以在 `component-map.yaml` 里人工确认一次，之后 match 会优先采用确认过的映射。
