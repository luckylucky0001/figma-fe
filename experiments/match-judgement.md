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

- **最终结果**：23 个项目（5 个调参 + 18 个留出）× 3 份设计稿共 759 项，**✓734 / △25 / ✗0**，逐项见 `real-match-eval.md`。下面几条是只有前两份设计稿、五个调参项目时的阶段性结论，留出项目第一次跑的错误见第 3 节。
- 订单详情（组件化的设计稿）：五个真实项目（覆盖 Vue+Vant、React+antd-mobile、XML+QMUI、Compose+Material3+自有设计系统）共 55 项，完全正确 52 项，可用但不是最优 3 项。
- 我的持仓（没组件化、中文图层名）：四个真实项目共 48 项，完全正确 45 项，可用但不是最优 3 项。
- 两份设计稿在这五个项目上都没有出现推荐错误组件的情况；但留出项目 Jetsnack 第一次跑出现了 3 个错误（见第 3 节），说明规则对没见过的项目结构仍会出错，修正后已加入回归。拿不准的情况，工具给出的都是「需确认」。
- 存量项目里「组件库有这个组件」不等于「项目能用」。版本、主题、废弃 API、设计系统封装都会影响结论，这几项现在都会检查。
- 项目惯例：没有对应组件、但项目里反复出现的写法（1px View 细线、Vant 的 `van-hairline--*` 类、容器自带的分割线属性）会作为「项目惯例」推荐，不会直接判成「需新建」。
- 还没覆盖的：「借用」写法（例如用按钮组件来充当标签）。这类可以在 `component-map.yaml` 里人工确认一次，之后 match 会优先采用确认过的映射。

## 3. 留出验证（调规则时没用过的项目）

上面的真实项目在迭代过程中也被用来发现问题、调整规则，所以结果偏乐观。为了看泛化能力，又拿两个之前没碰过的项目直接跑：

| 项目 | 第一次跑（没针对它调过） | 发现的问题 | 修正后 |
|---|---|---|---|
| Jetsnack（android/compose-samples，Compose + 自有 `Jetsnack*` 设计系统） | 订单详情 6 ✓ / 2 △ / **3 ✗** | ① 把桌面小部件（Glance）里的 `IconButton`、`ListItem` 当成了 App 组件；② `JetsnackCard` 输给了同名的 Material3 `Card`（没认出 `Jetsnack` 是项目前缀）；③ `JetsnackBottomBar` 其实是底部导航，被当成了操作栏；④ `ui/components` 目录没被认成设计系统 | 订单详情 10 ✓ / 1 △，我的持仓 11 ✓ / 1 △（Tag 推荐了可选中的 FilterChip，项目里没有纯展示标签） |
| vant-demo（youzan 官方，Vue 3 + TS / Vue 2 两套） | 两套各 23 项全部 ✓ | — | 同左 |
| litemall-vue（开源商城的移动端，Vue 2 + **Vant 2**） | 订单详情 10 ✓ / **1 ✗**，我的持仓 11 ✓ / **1 ✗** | 推荐了 `ActionBar`，但 Vant 2 里没有这个组件（Vant 3 才由 `GoodsAction` 改名而来）。同一个问题也让我在 vant-demo（Vue 2）的标准答案里写错了 | 组件库目录加入「改名 / 移除」信息，按项目实际安装的版本过滤：Vant 2 推荐 `GoodsAction`，Vant 3+ 推荐 `ActionBar`。之后又从一次实验里发现项目全局注册了 npm 包 `@chenfengyuan/vue-countdown`（5 个页面在用），工具原先看不到这种「单组件第三方包」，现在能识别，倒计时改推荐它。两份设计稿全部 ✓ |
| hello-uniapp（DCloud 官方 uni-app 示例，uni-ui + 内置组件） | 订单详情 7 ✓ / 1 △ / **3 ✗**，我的持仓 **12 ✗** | ① 没有 uni-ui 组件库目录；② `<view>` `<button>` 这类 uni-app 内置的小写标签被当成了同名的项目页面组件（View 被「用」了 1271 次）；③ 中文图层名和 `View` 规范化后都是空字符串，被判成「同名」，于是持仓页 12 项全推荐 View；④ 项目里偶尔出现一次的 Vant 组件压过了主力的 uni-ui；⑤ uni-ui / Vant 的 Badge 是角标而不是标签 | 加入 uni-ui 和 uni-app 内置组件目录、修掉空名字误判、项目几乎不用的组件库降权、Badge 单独成概念（带 text / tone 的 Badge 仍算标签）。订单详情 10 ✓ / 1 △，我的持仓 12 ✓ |
| taro-ui-demo（NervJS 官方，Taro + taro-ui，React 写在 .js 里） | 两份设计稿推荐的组件都对，但全部只到「需确认」，体检显示 0 个资产 | ① 没有 taro-ui 组件库目录；② React 组件写在普通 `.js` 文件里（不是 .jsx / .tsx），没被当成有标记的文件，用法一次都没统计到 | 加入 taro-ui 目录，`.js` 里有 JSX 且引用 react / @tarojs 时按 JSX 处理。资产 49 个、组件化率 100%，两份设计稿全部 ✓ |
| nutui-demo（京东官方，Vue 3 + NutUI 4，vite / taro 两套） | 23 项里 22 ✓ / **1 ✗** | 商品卡推荐了 `Sku`（规格选择弹窗），NutUI 的商品卡其实叫 `Card` | 补上 NutUI 的概念覆盖（Card = 商品卡，Sku 不是卡片）。全部 ✓，价格直接复用 NutUI 自带的 `Price` |
| Pokedex（skydoves，XML + Material Components 1.14，MaterialComponents 主题） | 两份设计稿第一次跑就全部 ✓ | — | 与 QMUI / NewPipe 相反，这里 MaterialButton / Chip / MaterialDivider 都能用，工具也照常推荐，说明版本和主题检查不是一刀切 |
| Sunflower（android 官方示例，Compose + Material3，没有自有设计系统） | 三份设计稿第一次跑就没有 ✗（31 ✓ / 2 △） | — | Material3 没有纯展示的标签组件，Tag 给的是「需确认」的 InputChip，算 △ |
| NewPipe（大型开源 App，XML View 为主，部分 Compose） | 订单详情 8 ✓ / 2 △ / **1 ✗**，我的持仓 5 ✓ / 5 △ / **2 ✗** | ① AppCompat `Toolbar` 等标准控件被当成「有组件却没用」（只因组件库目录里有 MaterialToolbar，而项目根本没用过），也就没算进项目资产，导航栏推荐成了 Compose 的 TopAppBar；② `TabLayout` 这类「概念 + Layout」的两词名字没识别出概念；③ `@layout/error_panel` 按 panel 当成了卡片；④ XML 页面里 Compose 的 Button 仍可能得到「直接复用」 | 订单详情 9 ✓ / 2 △，我的持仓 8 ✓ / 4 △，无 ✗。剩下的 △ 都是「项目里只有 Compose 版的按钮 / 列表行、XML 页面没有现成可用的」，工具给的是「需确认」 |

修正的都是通用规则（组件库版本、识别 Glance、学习项目前缀、按参数判断底部导航、`ui/components` 算设计系统、必填参数里有业务数据模型的算业务组件），不是针对 Jetsnack 的特判。之后 `python3 experiments/real-match-eval.py` 会同时跑调参集和留出集，作为回归测试。

### 第二批留出项目（报告初版之后加入）

写完报告初版后，又找了三个之前没碰过的项目，先写好标准答案再跑：

| 项目 | 第一次跑 | 发现的问题 → 通用修正 |
|---|---|---|
| vue3-vant-mobile（easy-temps，Vue 3 + Vant 4 模板） | 33 项全部 ✓ | — |
| mall-app-web（macrozheng，uni-app，自有组件为主，uni-ui 装了但几乎没用） | 订单详情 11 ✓；我的持仓 11 ✓ / **1 ✗**；买入下单 **程序崩溃（10 ✗）** | ① uni-app 内置组件没有 import 路径，生成示例代码时数组越界崩溃 → 内置标签直接写成 `<switch … />`，并加了覆盖所有组件库的测试；② `<view>` `<text>` 这类内置标签用了上千次，把 uni-ui 的「使用占比」压到 5% 以下，被当成「项目几乎不用的组件库」降权，标签页因此判成需新建 → 算占比时排除平台内置元素（和 Android SDK 控件同样处理） |
| AntennaPod（大型播客 App，XML + Material 1.14，Material3 主题） | 订单详情 10 ✓ / 1 △；我的持仓 11 ✓ / 1 △；买入下单 8 ✓ / 1 △ / **1 ✗** | 空状态推荐「需新建」，但项目的 6 个页面都通过 `EmptyViewHandler` 显示空状态。它不是 View 子类，只是一个加载布局的封装类，工具看不到 → 新增「界面封装类」资产：名字以 Handler / Helper / Controller 结尾、加载了布局、名字能对上 UI 概念的类，用法按代码里的调用统计 |

修 ② 时顺带加了一条通用规则：**项目自己的组件对某个概念已经用了 3 次以上时，项目从没用过的组件库组件要降权**。否则 mall-app-web 的单元格会推荐 uni-ui 的 ListItem，而不是项目里用了 6 次的 mix-list-cell；页面检查本来就会把这种情况报成「和项目惯用组件不一致」，match 应该和它一致。这条规则让两个老项目的推荐变了，我把它们的标准答案改了，理由如下（需要说明：这是看到结果之后改的）：

- NewPipe 的输入框：现在推荐 `NewPipeEditText`。它是项目对 AppCompatEditText 的封装（修 EMUI 上的分享菜单），布局里用了 6 次，原生 EditText 只有 1 次。原来的标准答案漏了它。
- litemall 的空状态：现在推荐项目自己的 `is-empty`（3 个页面在用）。Vant 的 Empty 仍算正确，但新页面引入 van-empty 会被页面检查报成「不一致」。

另外，mall-app-web 的商品卡在 uni-ui 不再被降权后，从「需新建」变成推荐 uni-ui 的 Card（需确认），记为 △。

### 第三批留出项目

| 项目 | 第一次跑 | 发现的问题 → 通用修正 |
|---|---|---|
| WanAndroid（goweii，XML，AppCompat 主题，国内常见写法） | 订单详情 8 ✓ / 1 △ / **2 ✗**；我的持仓 7 ✓ / 1 △ / **4 ✗**；买入下单 4 ✓ / 1 △ / **5 ✗** | ① 顶栏用的是三方库 `per.goweii.actionbarex` 的 ActionBarCommon（34 处），开关是 androidx 的 SwitchCompat，空状态是三方 MultiStateView（26 处）：组件库目录里没有的全限定名 View 原来一律被忽略 → 布局里出现的三方 View 按包名归为组件库资产（`AppCompat*` 仍算平台控件，已收录组件库的内部类不算）；② Android 的 ActionBar 指顶栏，H5（Vant）的 ActionBar 指底部操作栏 → Android 上 `ActionBar*` / `TitleBar*` 归为导航栏；③ 输入框推荐了专门输密码的 PasswordInputView，按钮推荐了文章页里的 `@layout/activity_article_float_btn` 和 FloatingActionButton → 名字带专门用途词（password / confirm / video / download …）的项目组件降权，`activity_*` / `fragment_*` 布局片段不算组件 |
| Seal（JunkFood02，Compose + Material3，自有 ui/component） | 订单详情 7 ✓ / 2 △ / **2 ✗**；我的持仓 7 ✓ / 1 △ / **4 ✗**；买入下单 8 ✓ / 1 △ / **1 ✗** | ① 卡片推荐了 VideoCard：项目里 Video* 组件多，被当成了「项目前缀」，于是 VideoCard 和 Card 算同名 → 专门用途词不能当项目前缀；② 按钮推荐了 ConfirmButton（对话框的确认按钮）/ ActionButton（必填 `Task.DownloadState`）→ 嵌套的业务类型也算业务参数；③ M3 Button 被判成「已被 LinkButton 封装」，但项目直接用 Button 16 次、LinkButton 只有 4 次 → 封装组件的使用次数不到库组件一半时，不算设计系统封装 |

| KunMinX Jetpack-MVVM-Best-Practice（XML + DataBinding，AppCompat 主题，依赖里有 Material 1.9） | 订单详情 8 ✓ / 1 △ / **2 ✗**；我的持仓 11 ✓ / 1 △；买入下单 6 ✓ / 3 △ / **1 ✗** | 卡片、分割线、搜索栏推荐了 MaterialCardView / MaterialDivider / SearchBar：工具已经知道它们在 AppCompat 主题下会崩溃，也降了分，但没有别的候选时它们仍排第一、结论是「需确认」，和体检报告里「match 不会推荐它们」矛盾 → 用不了的组件只作为警告列在可用候选之后，不再决定结论；没有可用候选时判为需新建 |
| jeecg-uniapp（JeecgBoot 的移动端，uni-app + wot-design-uni） | 订单详情 7 ✓ / 1 △ / **3 ✗**；我的持仓 8 ✓ / 1 △ / **3 ✗**；买入下单 5 ✓ / **5 ✗** | 项目的组件库是 wot-design-uni（`wd-*`，488 处），组件库目录里没有，整个没被识别 → 加入 wot-design-uni 目录（StatusTip = 空状态、InputNumber = 步进器）；之后分割线推荐了项目里的 JLine（其实是 echarts 折线图，靠「line」一个弱词沾上分割线），被「项目惯用组件」规则抬高 → 只靠弱词关联概念的组件不算项目惯用 |

### 实测：sunniejs/vue-h5-template（国内常用的 Vue 3 + Vant H5 模板）

这个模板同时装了 Vant、NutUI、Varlet 三套组件库，用 `.env` 里的 `VITE_UI_FRAMEWORK=vant` 选定一套（三套各有一个演示页 `src/views/ui/*.vue`）。第一次跑的结果：

- 复用清单 **10 ✓ / 23 ✗**：Vant 和 NutUI 的使用次数完全相同（都只在演示页用了 9 次），打分打平，按名字排序选中了 NutUI，于是同一个页面里一半推荐 Vant、一半推荐 NutUI，照做会把两套组件库混在一个页面里；
- 体检把 `src/views/ui/vant.vue` 这些演示页当成了「设计系统组件」（路径里的 `/ui/` 命中了设计系统目录规则），把 Vite 别名 `#ui-demo` 当成了三方组件库。

修正（通用规则）：装了多套 H5 组件库时先定一套「主用组件库」（优先看 `.env*` 里名字带 UI 的配置项，否则按使用次数），其他组件库的组件降权，主用库没有的组件（如 Price）才会给出别的库的候选并标「需确认」；目录分类按最外层目录判断（`views/ui/` 是页面）；只有 `package.json` 里依赖的包才算三方组件库。修正后 **32 ✓ / 1 △ / 0 ✗**。

这次修正让 hello-uniapp 的分割线从 Vant 的 Divider（项目只是偶尔用了 Vant）变成推荐项目自己的 CSS 类 `.uni-divider`（定义在项目的 `common/uni.css` 里）。我把它加进了标准答案，这是看到结果之后改的。

另外修正了 MaterialRadioButton 因名字包含「Button」被推荐给按钮的问题：候选名字比设计稿名字更具体、又属于别的概念时不算匹配（反过来 ProductCard → Card 仍然算）。修正后两个项目都没有 ✗，之前 18 个项目的结果不变（有一项 △ 变成 ✓）。

## 4. 第三份设计稿：买入下单（英文图层名、没组件化、新概念）

前两份设计稿覆盖的 UI 概念有限，于是又做了一份券商「买入下单」页（`design/figma-export/trade-order`）：图层叫 SearchBar / Tabs / Stepper / Input / Switch / Checkbox / Empty，但没有组件化。拿当时的工具直接在 14 个项目上跑（140 项），这是第一次的结果：

**123 ✓ / 5 △ / 12 ✗**，主要问题：

- 开关、复选框、单选原来是同一个概念，于是 QMUI / NewPipe 的「Switch」「Checkbox」都推荐成了 RadioButton，Now in Android 推荐成了切换按钮（6 处 ✗）；
- 搜索算在「输入框」里，Vant 项目推荐了 Field 而不是 Search；
- NutUI 的步进器叫 InputNumber、Jetsnack 的叫 QuantitySelector，没被认出来；
- Material 的 SearchBar / MaterialSwitch 是 Material 3 组件，在 MaterialComponents 主题的项目里不能直接用；
- NewPipe 的 SwitchPreference（设置框架里的行）被当成了开关组件。

修正后（同样是通用规则：拆开概念、补词表、Material 3 主题检查、排除 Preference）：**132 ✓ / 8 △ / 0 ✗**（加上另外两份设计稿，三份合计见 `real-match-eval.md`）。

