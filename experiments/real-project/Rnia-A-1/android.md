
## 页面复用检查：app/src/main/res/values/strings.xml,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/icon/NiaIcons.kt,app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailMockData.kt,app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/ActionBar.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/Card.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/Cell.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/CountDown.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/NavBar.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/NoticeBar.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/PillButton.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/Price.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/ProductCard.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/StatusTag.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt  ❌ FAIL

- 复用的现成组件（2）：DynamicAsyncImage（Composable）、HorizontalDivider（组件库）
- 自造轮子（6）
  - ✗ core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/Cell.kt：页面里自造了 NiaCell，已有 ListItem（组件库 Compose Material3）
  - ✗ core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/PillButton.kt：页面里自造了 NiaPillButton，已有 NiaTextButton（Composable core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt）
  - ✗ core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/Card.kt：页面里自造了 NiaDivider，已有 HorizontalDivider（组件库 Compose Material3）
  - ✗ core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/Card.kt：页面里自造了 NiaCard，已有 Card（组件库 Compose Material3）
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/StatusTag.kt：页面里自造了 NiaStatusTag，已有 Badge（组件库 Compose Material3）（组件库标准组件，项目还没用过，请确认）
  - ✗ core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/NavBar.kt：页面里自造了 NiaNavBar，已有 NiaTopAppBar（Composable core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/TopAppBar.kt）
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

技术栈：Jetpack Compose / Compose Material3 · 扫描文件 258

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（195 / 195） |
| 重复概念（同一 UI 概念有多个实现） | 14 |
| 重复结构（可抽成组件的复制粘贴结构） | 0 组 |
| 硬编码颜色 | 80 处（其中 5 处已有同值 token） |

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Compose Material3 | 21 | HorizontalDivider×3、TextButton×3、CircularProgressIndicator×2、Card×2、TextField×1、OutlinedButton×1、AlertDialog×1、FilterChip×1 |

**组件库里有、但本项目不能直接用的组件**（match 不会推荐它们）：

- `Divider`：已废弃，改用 HorizontalDivider

**设计系统封装**（组件库组件 → 项目的封装；新页面一律用右边的封装）：`Button` → `NiaButton`、`CenterAlignedTopAppBar` → `NiaTopAppBar`、`FilterChip` → `NiaFilterChip`、`HorizontalDivider` → `NiaDivider`、`NavigationBar` → `NiaNavigationBar`、`OutlinedButton` → `NiaOutlinedButton`、`TabRow` → `NiaTabRow`、`Tab` → `NiaTab`、`TextButton` → `NiaTextButton`

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `TextButton` | 组件库 | button | 3 | `TextButton(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:197 |
| `HorizontalDivider` | 组件库 | divider | 3 | `HorizontalDivider(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/Card.kt:76 |
| `CircularProgressIndicator` | 组件库 | loading | 2 | `CircularProgressIndicator(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/DynamicAsyncImage.kt:72 |
| `Card` | 组件库 | card | 2 | `Card(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:113 |
| `TextField` | 组件库 | input | 1 | `TextField(` feature/search/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/search/impl/SearchScreen.kt:484 |
| `Button` | 组件库 | button | 1 | `Button(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:58 |
| `ListItem` | 组件库 | cell | 1 | `ListItem(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:55 |
| `AlertDialog` | 组件库 | dialog | 1 | `AlertDialog(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:106 |
| `Tab` | 组件库 | tabs | 1 | `Tab(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Tabs.kt:55 |
| `TabRow` | 组件库 | tabs | 1 | `TabRow(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Tabs.kt:88 |
| `OutlinedButton` | 组件库 | button | 1 | `OutlinedButton(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:124 |
| `RadioButton` | 组件库 | switch | 1 | `RadioButton(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:242 |
| `FilterChip` | 组件库 | tag | 1 | `FilterChip(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Chip.kt:52 |
| `NavigationBar` | 组件库 | tabbar | 1 | `NavigationBar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Navigation.kt:101 |
| `CenterAlignedTopAppBar` | 组件库 | navbar | 1 | `CenterAlignedTopAppBar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/TopAppBar.kt:52 |
| `NiaCell` | Composable | cell | 11 | `NiaCell(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:168 |
| `NiaTextButton` | Composable | button | 10 | `NiaTextButton(onClick = {}) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:87 |
| `NiaIconToggleButton` | Composable | switch | 9 | `NiaIconToggleButton(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:204 |
| `NiaBackground` | Composable |  | 8 | `NiaBackground(modifier = modifier) {` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:98 |
| `NiaButton` | Composable | button | 8 | `NiaButton(onClick = {}) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:81 |
| `NiaPillButton` | Composable | button | 7 | `NiaPillButton(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:117 |
| `InterestsItem` | Composable | cell | 7 | `InterestsItem(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:127 |
| `SettingsDialogThemeChooserRow` | Composable | dialog | 7 | `SettingsDialogThemeChooserRow(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:168 |
| `TrackScreenViewEvent` | Composable |  | 6 | `TrackScreenViewEvent(screenName = "Saved")` feature/bookmarks/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/bookmarks/impl/BookmarksScreen.kt:149 |
| `NiaLoadingWheel` | Composable | loading | 6 | `NiaLoadingWheel(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/LoadingWheel.kt:146 |
| `NiaOutlinedButton` | Composable | button | 6 | `NiaOutlinedButton(onClick = {}) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:84 |
| `NiaFilterChip` | Composable | tag | 6 | `NiaFilterChip(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:175 |
| `NiaViewToggleButton` | Composable | switch | 5 | `NiaViewToggleButton(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:277 |
| `NiaTopicTag` | Composable | tag | 5 | `NiaTopicTag(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:302 |
| `DynamicAsyncImage` | Composable | image | 4 | `DynamicAsyncImage(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/commerce/ProductCard.kt:71 |
| `NiaGradientBackground` | Composable |  | 4 | `NiaGradientBackground(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:99 |
| `NiaDivider` | Composable | divider | 4 | `NiaDivider()` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:167 |
| `TrackScrollJank` | Composable |  | 4 | `TrackScrollJank(scrollableState = scrollableState, stateName = "bookmarks:grid")` feature/bookmarks/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/bookmarks/impl/BookmarksScreen.kt:172 |
| `SettingsDialog` | Composable | dialog | 4 | `SettingsDialog(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:153 |
| `NiaPrice` | Composable | price | 4 | `NiaPrice(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:220 |
| `NewsResourceCardExpanded` | Composable | card | 3 | `NewsResourceCardExpanded(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsFeed.kt:66 |
| `BookmarkButton` | Composable | button | 3 | `BookmarkButton(isBookmarked, onToggleBookmark)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:152 |
| `NiaCard` | Composable | card | 3 | `NiaCard(title = stringResource(R.string.order_detail_goods_section)) {` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:181 |
| `NiaStatusTag` | Composable | tag | 3 | `NiaStatusTag(text = order.statusLabel, tone = NiaStatusTagTone.Warning)` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:144 |
| `NiaCountDown` | Composable | countdown | 2 | `NiaCountDown(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:155 |
| `NiaTabRow` | Composable | tabs | 2 | `NiaTabRow(selectedTabIndex = selectedTabIndex) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:324 |
| `EmptySearchResultBody` | Composable | input | 2 | `EmptySearchResultBody(` feature/search/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/search/impl/SearchScreen.kt:173 |
| `InterestsDetailPlaceholder` | Composable | empty | 2 | `InterestsDetailPlaceholder()` feature/interests/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/interests/impl/InterestsDetailPlaceholder.kt:70 |
| `NiaOverlayLoadingWheel` | Composable | loading | 2 | `NiaOverlayLoadingWheel(contentDesc = "LoadingWheel")` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/LoadingWheel.kt:167 |
| `Scrollbar` | Composable |  | 2 | `Scrollbar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/scrollbar/AppScrollbars.kt:84 |
| `NiaNoticeBar` | Composable | notice | 2 | `NiaNoticeBar(text = order.notice, tone = NiaNoticeBarTone.Warning)` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:111 |
| `NiaNavBar` | Composable | navbar | 2 | `NiaNavBar(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:99 |
| `NiaProductCard` | Composable | goods | 2 | `NiaProductCard(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:182 |
| `NiaTopAppBar` | Composable | navbar | 2 | `NiaTopAppBar(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:229 |
| `NiaTab` | Composable | tabs | 2 | `NiaTab(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:326 |
| `NewsResourceShortDescription` | Composable |  | 1 | `NewsResourceShortDescription(userNewsResource.content)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:166 |
| `NewsResourceTitle` | Composable |  | 1 | `NewsResourceTitle(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:137 |
| `NewsResourceTopics` | Composable |  | 1 | `NewsResourceTopics(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:168 |
| `TrackDisposableJank` | Composable |  | 1 | `TrackDisposableJank(navigationState.currentKey) { metricsHolder ->` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaAppState.kt:115 |
| `NiaCatalog` | Composable |  | 1 | `setContent { NiaCatalog() }` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/NiaCatalogActivity.kt:31 |
| `TrackJank` | Composable |  | 1 | `TrackJank(scrollableState) { metricsHolder ->` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/JankStatsExtensions.kt:82 |
| `DebugVisibleWindowInsets` | Composable |  | 1 | `DebugVisibleWindowInsets()` app/src/testDemo/kotlin/com/google/samples/apps/nowinandroid/ui/SnackbarInsetsScreenshotTests.kt:267 |
| `OrderDetailContent` | Composable |  | 1 | `OrderDetailContent(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:72 |
| `NewsResourceMetaData` | Composable |  | 1 | `NewsResourceMetaData(userNewsResource.publishDate, userNewsResource.type)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:163 |
| `NiaNavigationSuiteScaffold` | Composable | tabbar | 1 | `NiaNavigationSuiteScaffold(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:162 |
| … | | | | 共 87 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| button | NiaTextButton（Composable，用 10 次）<br>NiaButton（Composable，用 8 次）<br>NiaPillButton（Composable，用 7 次）<br>NiaOutlinedButton（Composable，用 6 次）<br>IconButtonPreviewUnchecked（Composable，用 0 次）<br>TextButton（组件库，用 3 次）<br>Button（组件库，用 1 次）<br>OutlinedButton（组件库，用 1 次）<br>BookmarkButton（Composable，用 3 次） | `NiaTextButton` |
| loading | NiaLoadingWheel（Composable，用 6 次）<br>NiaOverlayLoadingWheel（Composable，用 2 次）<br>CircularProgressIndicator（组件库，用 2 次）<br>ForYouScreenPopulatedAndLoading（Composable，用 0 次）<br>TopicScreenLoading（Composable，用 0 次）<br>InterestsScreenLoading（Composable，用 0 次）<br>ForYouScreenLoading（Composable，用 0 次） | `NiaLoadingWheel` |
| switch | NiaIconToggleButton（Composable，用 9 次）<br>NiaViewToggleButton（Composable，用 5 次）<br>ViewTogglePreviewCompact（Composable，用 0 次）<br>ViewTogglePreviewExpanded（Composable，用 0 次）<br>RadioButton（组件库，用 1 次）<br>Switch（Compose Material3，项目未使用） | `NiaIconToggleButton` |
| tag | NiaFilterChip（Composable，用 6 次）<br>NiaTopicTag（Composable，用 5 次）<br>NiaStatusTag（Composable，用 3 次）<br>FilterChip（组件库，用 1 次） | `NiaFilterChip` |
| tabs | NiaTab（Composable，用 2 次）<br>Tab（组件库，用 1 次）<br>TabRow（组件库，用 1 次）<br>TopicsTabContent（Composable，用 1 次） | `NiaTab` |
| tabbar | NiaNavigationSuiteScaffold（Composable，用 1 次）<br>NiaNavigationRail（Composable，用 0 次）<br>NiaNavigationBar（Composable，用 0 次）<br>NavigationBar（组件库，用 1 次） | `NiaNavigationSuiteScaffold` |
| cell | NiaCell（Composable，用 11 次）<br>ListItem（组件库，用 1 次）<br>InterestsItem（Composable，用 7 次） | `NiaCell` |
| card | NiaCard（Composable，用 3 次）<br>Card（组件库，用 2 次）<br>NewsResourceCardExpanded（Composable，用 3 次） | `NiaCard` |
| divider | NiaDivider（Composable，用 4 次）<br>HorizontalDivider（组件库，用 3 次）<br>Divider（Compose Material3，项目未使用） | `NiaDivider` |
| navbar | NiaNavBar（Composable，用 2 次）<br>NiaTopAppBar（Composable，用 2 次）<br>CenterAlignedTopAppBar（组件库，用 1 次） | `NiaNavBar` |
| dialog | AlertDialog（组件库，用 1 次）<br>SettingsDialog（Composable，用 4 次） | `AlertDialog` |
| image | DynamicAsyncImage（Composable，用 4 次）<br>NewsResourceHeaderImage（Composable，用 1 次） | `DynamicAsyncImage` |
| input | TextField（组件库，用 1 次）<br>EmptySearchResultBody（Composable，用 2 次） | `TextField` |
| empty | InterestsDetailPlaceholder（Composable，用 2 次）<br>InterestsScreenEmpty（Composable，用 0 次） | `InterestsDetailPlaceholder` |

## 4. 有组件却没用（2 处，含绕过设计系统封装直接用组件库）

- feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:117 `HorizontalDivider()` → 直接用了 HorizontalDivider（Compose Material3），项目设计系统已封装为 NiaDivider，应使用封装
- feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:137 `HorizontalDivider(Modifier.padding(top = 8.dp))` → 直接用了 HorizontalDivider（Compose Material3），项目设计系统已封装为 NiaDivider，应使用封装

## 6. 硬编码颜色热点

| 颜色 | 次数 | 已有 token | 出现在 |
|---|---|---|---|
| #FFFFFF | 2 | `@color/ic_launcher_background_tint` | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt |
| #0D1F12 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #4F6352 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #A23F16 | 1 | `@color/ic_launcher_foreground_tint` | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #702776 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #D0C3CC | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #690005 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #A2CED9 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #E8F3FF | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceTheme.kt |
| #ECDFE0 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #00210B | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #DDE5DB | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
