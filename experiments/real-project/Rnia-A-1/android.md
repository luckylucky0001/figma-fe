
## 页面复用检查：15 个文件（app/src/main/res/values/strings.xml、core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/icon/NiaIcons.kt …）  ❌ FAIL

- 复用的现成组件（2）：DynamicAsyncImage（Composable）、HorizontalDivider（组件库）
- 自造轮子（6）
  - ✗ core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/Cell.kt：新建了 NiaCell，已有 ListItem（组件库 Compose Material3）
  - ✗ core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/PillButton.kt：新建了 NiaPillButton，已有 NiaButton（Composable core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt）
  - ✗ core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/Card.kt：新建了 NiaDivider，已有 HorizontalDivider（组件库 Compose Material3）
  - ✗ core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/Card.kt：新建了 NiaCard，已有 Card（组件库 Compose Material3）
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/StatusTag.kt：新建了 NiaStatusTag；同类已有 NiaFilterChip（Composable），用途可能不同（filter），请确认
  - ✗ core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/NavBar.kt：新建了 NiaNavBar，已有 NiaTopAppBar（Composable core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/TopAppBar.kt）
- 有组件却没用（0）
- 疑似手写已有组件（0）
- 写死颜色（18）
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt:34 Color(0xFFFF4D1A)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt:35 Color(0xFFFFEDE8)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt:36 Color(0xFFFFFFFF)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt:37 Color(0xFF1F2329)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt:38 Color(0xFF646A73)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt:39 Color(0xFF8F959E)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt:40 Color(0xFFF5F6F7)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt:41 Color(0xFFFFFFFF)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt:42 Color(0xFFF2F3F5)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt:43 Color(0xFFE5E6EB)：写死颜色
  - … 共 18 处

---

# 存量复用体检：/home/user/exp/Rnia-A-1（android）

技术栈：Jetpack Compose / Compose Foundation / Compose Material3 · 扫描文件 258

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（197 / 197） |
| 重复概念（同一 UI 概念有多个实现） | 4 |
| 重复结构（可抽成组件的复制粘贴结构） | 0 组 |
| 硬编码颜色 | 80 处（其中 5 处已有同值 token） |

## 0. 优先处理（按影响面排序）

1. 「cell（Compose）」有 2 种实现（共用 12 次）→ 统一到 `NiaCell`
2. 「loading（Compose）」有 3 种实现（共用 10 次）→ 统一到 `NiaLoadingWheel`
3. 2 处HorizontalDivider（Compose Material3），项目设计系统已封装为 NiaDivider，应使用封装
4. 「card（Compose）」有 2 种实现（共用 5 次）→ 统一到 `NiaCard`
5. 5 处写死的颜色已有同值 token，可以直接替换（见第 6 节）

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Compose Material3 | 21 | HorizontalDivider×3、TextButton×3、Card×2、CircularProgressIndicator×2、CenterAlignedTopAppBar×1、ListItem×1、Tab×1、OutlinedButton×1 |
| Compose Foundation | 10 | LazyVerticalStaggeredGrid×5、LazyColumn×4、LazyHorizontalGrid×1 |

**组件库里有、但本项目不能直接用的组件**（match 不会推荐它们）：

- `Divider`：已废弃，改用 HorizontalDivider

**设计系统封装**（组件库组件 → 项目的封装；新页面一律用右边的封装）：`Button` → `NiaButton`、`CenterAlignedTopAppBar` → `NiaTopAppBar`、`FilterChip` → `NiaFilterChip`、`HorizontalDivider` → `NiaDivider`、`NavigationBar` → `NiaNavigationBar`、`OutlinedButton` → `NiaOutlinedButton`、`TabRow` → `NiaTabRow`、`Tab` → `NiaTab`、`TextButton` → `NiaTextButton`

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `LazyVerticalStaggeredGrid` | 组件库 |  | 5 | `LazyVerticalStaggeredGrid(columns = StaggeredGridCells.Adaptive(300.dp)) {` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsFeed.kt:129 |
| `LazyColumn` | 组件库 | list | 4 | `LazyColumn(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:67 |
| `TextButton` | 组件库 | button | 3 | `TextButton(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:197 |
| `HorizontalDivider` | 组件库 | divider | 3 | `HorizontalDivider(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/Card.kt:76 |
| `Card` | 组件库 | card | 2 | `Card(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:113 |
| `CircularProgressIndicator` | 组件库 | loading | 2 | `CircularProgressIndicator(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/DynamicAsyncImage.kt:72 |
| `LazyHorizontalGrid` | 组件库 |  | 1 | `LazyHorizontalGrid(` feature/foryou/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/foryou/impl/ForYouScreen.kt:333 |
| `Tab` | 组件库 | tabs | 1 | `Tab(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Tabs.kt:55 |
| `Button` | 组件库 | button | 1 | `Button(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:58 |
| `ListItem` | 组件库 | cell | 1 | `ListItem(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:55 |
| `TabRow` | 组件库 | tabs | 1 | `TabRow(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Tabs.kt:88 |
| `RadioButton` | 组件库 | radio | 1 | `RadioButton(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:242 |
| `TextField` | 组件库 | input | 1 | `TextField(` feature/search/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/search/impl/SearchScreen.kt:484 |
| `OutlinedButton` | 组件库 | button | 1 | `OutlinedButton(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:124 |
| `CenterAlignedTopAppBar` | 组件库 | navbar | 1 | `CenterAlignedTopAppBar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/TopAppBar.kt:52 |
| `FilterChip` | 组件库 | tag | 1 | `FilterChip(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Chip.kt:52 |
| `AlertDialog` | 组件库 | dialog | 1 | `AlertDialog(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:106 |
| `NavigationBar` | 组件库 | tabbar | 1 | `NavigationBar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Navigation.kt:101 |
| `NiaCell` | Composable | cell | 11 | `NiaCell(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:168 |
| `NiaTextButton` | Composable | button | 10 | `NiaTextButton(onClick = {}) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:87 |
| `NiaIconToggleButton` | Composable | button | 9 | `NiaIconToggleButton(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:204 |
| `NiaBackground` | Composable |  | 8 | `NiaBackground(modifier = modifier) {` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:98 |
| `NiaButton` | Composable | button | 8 | `NiaButton(onClick = {}) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:81 |
| `SettingsDialogThemeChooserRow` | Composable | dialog | 7 | `SettingsDialogThemeChooserRow(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:168 |
| `NiaPillButton` | Composable | button | 7 | `NiaPillButton(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:117 |
| `InterestsItem` | Composable | cell | 7 | `InterestsItem(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:127 |
| `NiaLoadingWheel` | Composable | loading | 6 | `NiaLoadingWheel(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/LoadingWheel.kt:146 |
| `NiaOutlinedButton` | Composable | button | 6 | `NiaOutlinedButton(onClick = {}) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:84 |
| `NiaFilterChip` | Composable | tag | 6 | `NiaFilterChip(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:175 |
| `NiaViewToggleButton` | Composable | button | 5 | `NiaViewToggleButton(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:277 |
| `NiaTopicTag` | Composable | tag | 5 | `NiaTopicTag(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:302 |
| `TrackScrollJank` | Composable |  | 4 | `TrackScrollJank(scrollableState = scrollableState, stateName = "bookmarks:grid")` feature/bookmarks/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/bookmarks/impl/BookmarksScreen.kt:172 |
| `NiaGradientBackground` | Composable |  | 4 | `NiaGradientBackground(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:99 |
| `NiaPrice` | Composable | price | 4 | `NiaPrice(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:220 |
| `NiaDivider` | Composable | divider | 4 | `NiaDivider()` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:167 |
| `SettingsDialog` | Composable | dialog | 4 | `SettingsDialog(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:153 |
| `DynamicAsyncImage` | Composable | image | 4 | `DynamicAsyncImage(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/ProductCard.kt:71 |
| `NewsResourceCardExpanded` | Composable | card | 3 | `NewsResourceCardExpanded(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsFeed.kt:66 |
| `NiaCard` | Composable | card | 3 | `NiaCard(title = stringResource(R.string.order_detail_goods_section)) {` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:181 |
| `NiaStatusTag` | Composable | tag | 3 | `NiaStatusTag(text = order.statusLabel, tone = NiaStatusTagTone.Warning)` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:144 |
| `BookmarkButton` | Composable | button | 3 | `BookmarkButton(isBookmarked, onToggleBookmark)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:152 |
| `NiaTab` | Composable | tabs | 2 | `NiaTab(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:326 |
| `InterestsDetailPlaceholder` | Composable | empty | 2 | `InterestsDetailPlaceholder()` feature/interests/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/interests/impl/InterestsDetailPlaceholder.kt:70 |
| `NiaNoticeBar` | Composable | notice | 2 | `NiaNoticeBar(text = order.notice, tone = NiaNoticeBarTone.Warning)` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:111 |
| `Scrollbar` | Composable |  | 2 | `Scrollbar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/scrollbar/AppScrollbars.kt:84 |
| `NiaProductCard` | Composable | goods | 2 | `NiaProductCard(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:182 |
| `NiaTopAppBar` | Composable | navbar | 2 | `NiaTopAppBar(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:229 |
| `NiaNavBar` | Composable | navbar | 2 | `NiaNavBar(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:99 |
| `NiaCountDown` | Composable | countdown | 2 | `NiaCountDown(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:155 |
| `EmptySearchResultBody` | Composable | search | 2 | `EmptySearchResultBody(` feature/search/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/search/impl/SearchScreen.kt:173 |
| `NiaTabRow` | Composable | tabs | 2 | `NiaTabRow(selectedTabIndex = selectedTabIndex) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:324 |
| `NiaOverlayLoadingWheel` | Composable | loading | 2 | `NiaOverlayLoadingWheel(contentDesc = "LoadingWheel")` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/LoadingWheel.kt:167 |
| `TrackDisposableJank` | Composable |  | 1 | `TrackDisposableJank(navigationState.currentKey) { metricsHolder ->` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaAppState.kt:115 |
| `NewsResourceTitle` | Composable |  | 1 | `NewsResourceTitle(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:137 |
| `NewsResourceHeaderImage` | Composable | image | 1 | `NewsResourceHeaderImage(userNewsResource.headerImageUrl)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:128 |
| `NiaNavigationRailItem` | Composable | tabbar | 1 | `NiaNavigationRailItem(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Navigation.kt:324 |
| `TopicIcon` | Composable | icon | 1 | `TopicIcon(` feature/foryou/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/foryou/impl/ForYouScreen.kt:400 |
| `NiaNavigationSuiteScaffold` | Composable | tabbar | 1 | `NiaNavigationSuiteScaffold(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:162 |
| `TrackJank` | Composable |  | 1 | `TrackJank(scrollableState) { metricsHolder ->` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/JankStatsExtensions.kt:82 |
| `NewsResourceMetaData` | Composable |  | 1 | `NewsResourceMetaData(userNewsResource.publishDate, userNewsResource.type)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:163 |
| … | | | | 共 68 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| loading（Compose） | NiaLoadingWheel（Composable，用 6 次）<br>NiaOverlayLoadingWheel（Composable，用 2 次）<br>CircularProgressIndicator（组件库，用 2 次） | `NiaLoadingWheel` |
| cell（Compose） | NiaCell（Composable，用 11 次）<br>ListItem（组件库，用 1 次） | `NiaCell` |
| card（Compose） | NiaCard（Composable，用 3 次）<br>Card（组件库，用 2 次） | `NiaCard` |
| divider（Compose） | NiaDivider（Composable，用 4 次）<br>Divider（Compose Material3，项目未使用） | `NiaDivider` |

## 4. 有组件却没用（2 处，含绕过设计系统封装直接用组件库）

- feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:117 `HorizontalDivider()` → 直接用了 HorizontalDivider（Compose Material3），项目设计系统已封装为 NiaDivider，应使用封装
- feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:137 `HorizontalDivider(Modifier.padding(top = 8.dp))` → 直接用了 HorizontalDivider（Compose Material3），项目设计系统已封装为 NiaDivider，应使用封装

## 6. 硬编码颜色热点

| 颜色 | 次数 | 已有 token | 出现在 |
|---|---|---|---|
| #FFFFFF | 2 | `@color/ic_launcher_background_tint` | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt |
| #BA1A1A | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #E8F3FF | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt |
| #B7CCB8 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #003919 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #414941 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #E8FFEA | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt |
| #B8EAFF | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #C1C9BF | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #380D00 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #FFA9FE | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #FFEDE8 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt |
