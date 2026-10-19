
## 页面复用检查：core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/TopAppBar.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/icon/NiaIcons.kt,app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/ActionBar.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Countdown.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/ListItem.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Price.kt,core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Status.kt,core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/ProductCard.kt  ❌ FAIL

- 复用的现成组件（6）：Card（组件库）、CenterAlignedTopAppBar（组件库）、DynamicAsyncImage（Composable）、HorizontalDivider（组件库）、NiaButton（Composable）、NiaOutlinedButton（Composable）
- 自造轮子（2）
  - ✗ core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/ListItem.kt：页面里自造了 NiaListItem，已有 ListItem（组件库 Compose Material3）
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Status.kt：页面里自造了 NiaStatusTag，已有 Badge（组件库 Compose Material3）（组件库标准组件，项目还没用过，请确认）
- 有组件却没用（0）
- 疑似手写已有组件（0）
- 写死颜色（0）

---

# 存量复用体检：/home/user/exp/Rnia-C-1（android）

技术栈：Jetpack Compose / Compose Material3 · 扫描文件 252

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（181 / 181） |
| 重复概念（同一 UI 概念有多个实现） | 13 |
| 重复结构（可抽成组件的复制粘贴结构） | 0 组 |
| 硬编码颜色 | 62 处（其中 3 处已有同值 token） |

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Compose Material3 | 25 | HorizontalDivider×5、TextButton×3、Card×3、CircularProgressIndicator×2、CenterAlignedTopAppBar×2、RadioButton×1、NavigationBar×1、FilterChip×1 |

**组件库里有、但本项目不能直接用的组件**（match 不会推荐它们）：

- `Divider`：已废弃，改用 HorizontalDivider

**设计系统封装**（组件库组件 → 项目的封装；新页面一律用右边的封装）：`Button` → `NiaButton`、`CenterAlignedTopAppBar` → `NiaTopAppBar`、`FilterChip` → `NiaFilterChip`、`NavigationBar` → `NiaNavigationBar`、`OutlinedButton` → `NiaOutlinedButton`、`TabRow` → `NiaTabRow`、`Tab` → `NiaTab`、`TextButton` → `NiaTextButton`

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `HorizontalDivider` | 组件库 | divider | 5 | `HorizontalDivider(modifier = Modifier.padding(top = 12.dp))` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:144 |
| `TextButton` | 组件库 | button | 3 | `TextButton(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:197 |
| `Card` | 组件库 | card | 3 | `Card(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:204 |
| `CenterAlignedTopAppBar` | 组件库 | navbar | 2 | `CenterAlignedTopAppBar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/TopAppBar.kt:52 |
| `CircularProgressIndicator` | 组件库 | loading | 2 | `CircularProgressIndicator(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/DynamicAsyncImage.kt:72 |
| `AlertDialog` | 组件库 | dialog | 1 | `AlertDialog(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:106 |
| `ListItem` | 组件库 | cell | 1 | `ListItem(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:55 |
| `Button` | 组件库 | button | 1 | `Button(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:58 |
| `OutlinedButton` | 组件库 | button | 1 | `OutlinedButton(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:124 |
| `NavigationBar` | 组件库 | tabbar | 1 | `NavigationBar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Navigation.kt:101 |
| `Tab` | 组件库 | tabs | 1 | `Tab(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Tabs.kt:55 |
| `FilterChip` | 组件库 | tag | 1 | `FilterChip(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Chip.kt:52 |
| `TextField` | 组件库 | input | 1 | `TextField(` feature/search/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/search/impl/SearchScreen.kt:484 |
| `RadioButton` | 组件库 | switch | 1 | `RadioButton(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:242 |
| `TabRow` | 组件库 | tabs | 1 | `TabRow(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Tabs.kt:88 |
| `NiaButton` | Composable | button | 10 | `NiaButton(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:96 |
| `NiaTextButton` | Composable | button | 10 | `NiaTextButton(onClick = {}) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:87 |
| `NiaIconToggleButton` | Composable | switch | 9 | `NiaIconToggleButton(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:204 |
| `NiaOutlinedButton` | Composable | button | 9 | `NiaOutlinedButton(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:86 |
| `NiaBackground` | Composable |  | 8 | `NiaBackground(modifier = modifier) {` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:98 |
| `SettingsDialogThemeChooserRow` | Composable | dialog | 7 | `SettingsDialogThemeChooserRow(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:168 |
| `InterestsItem` | Composable | cell | 7 | `InterestsItem(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:127 |
| `TrackScreenViewEvent` | Composable |  | 6 | `TrackScreenViewEvent(screenName = "Saved")` feature/bookmarks/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/bookmarks/impl/BookmarksScreen.kt:149 |
| `NiaLoadingWheel` | Composable | loading | 6 | `NiaLoadingWheel(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/LoadingWheel.kt:146 |
| `NiaFilterChip` | Composable | tag | 6 | `NiaFilterChip(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:175 |
| `NiaViewToggleButton` | Composable | switch | 5 | `NiaViewToggleButton(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:277 |
| `NiaTopicTag` | Composable | tag | 5 | `NiaTopicTag(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:302 |
| `NiaListItem` | Composable | cell | 5 | `NiaListItem(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:145 |
| `DynamicAsyncImage` | Composable | image | 4 | `DynamicAsyncImage(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:114 |
| `TrackScrollJank` | Composable |  | 4 | `TrackScrollJank(scrollableState = scrollableState, stateName = "bookmarks:grid")` feature/bookmarks/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/bookmarks/impl/BookmarksScreen.kt:172 |
| `NiaPriceText` | Composable | price | 4 | `NiaPriceText(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:180 |
| `NiaGradientBackground` | Composable |  | 4 | `NiaGradientBackground(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:99 |
| `SettingsDialog` | Composable | dialog | 4 | `SettingsDialog(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:153 |
| `NewsResourceCardExpanded` | Composable | card | 3 | `NewsResourceCardExpanded(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsFeed.kt:66 |
| `NiaTopAppBar` | Composable | navbar | 3 | `NiaTopAppBar(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:229 |
| `BookmarkButton` | Composable | button | 3 | `BookmarkButton(isBookmarked, onToggleBookmark)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:152 |
| `NiaStatusTag` | Composable | tag | 2 | `NiaStatusTag(text = order.statusTag, tone = NiaStatusTone.Warning)` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:128 |
| `InterestsDetailPlaceholder` | Composable | empty | 2 | `InterestsDetailPlaceholder()` feature/interests/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/interests/impl/InterestsDetailPlaceholder.kt:70 |
| `ProductCard` | Composable | goods | 2 | `ProductCard(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:157 |
| `EmptySearchResultBody` | Composable | input | 2 | `EmptySearchResultBody(` feature/search/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/search/impl/SearchScreen.kt:173 |
| `NiaTab` | Composable | tabs | 2 | `NiaTab(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:326 |
| `NiaTabRow` | Composable | tabs | 2 | `NiaTabRow(selectedTabIndex = selectedTabIndex) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:324 |
| `Scrollbar` | Composable |  | 2 | `Scrollbar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/scrollbar/AppScrollbars.kt:84 |
| `NiaCountdownText` | Composable | countdown | 2 | `NiaCountdownText(remaining = order.remaining, prefix = "剩余 ")` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:136 |
| `NiaNoticeBar` | Composable | notice | 2 | `NiaNoticeBar(text = order.notice, tone = NiaStatusTone.Warning)` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:112 |
| `NiaOverlayLoadingWheel` | Composable | loading | 2 | `NiaOverlayLoadingWheel(contentDesc = "LoadingWheel")` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/LoadingWheel.kt:167 |
| `NewsResourceShortDescription` | Composable |  | 1 | `NewsResourceShortDescription(userNewsResource.content)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:166 |
| `NewsResourceTitle` | Composable |  | 1 | `NewsResourceTitle(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:137 |
| `NiaNavigationRailItem` | Composable | tabbar | 1 | `NiaNavigationRailItem(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Navigation.kt:324 |
| `NewsResourceMetaData` | Composable |  | 1 | `NewsResourceMetaData(userNewsResource.publishDate, userNewsResource.type)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:163 |
| `TopicsTabContent` | Composable | tabs | 1 | `TopicsTabContent(` feature/interests/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/interests/impl/InterestsScreen.kt:78 |
| `TrackJank` | Composable |  | 1 | `TrackJank(scrollableState) { metricsHolder ->` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/JankStatsExtensions.kt:82 |
| `NiaNavigationSuiteScaffold` | Composable | tabbar | 1 | `NiaNavigationSuiteScaffold(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:162 |
| `NiaCatalog` | Composable |  | 1 | `setContent { NiaCatalog() }` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/NiaCatalogActivity.kt:31 |
| `TopicIcon` | Composable | icon | 1 | `TopicIcon(` feature/foryou/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/foryou/impl/ForYouScreen.kt:400 |
| `NewsResourceTopics` | Composable |  | 1 | `NewsResourceTopics(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:168 |
| `NewsResourceHeaderImage` | Composable | image | 1 | `NewsResourceHeaderImage(userNewsResource.headerImageUrl)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:128 |
| `NotificationDot` | Composable |  | 1 | `NotificationDot(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:157 |
| `DebugVisibleWindowInsets` | Composable |  | 1 | `DebugVisibleWindowInsets()` app/src/testDemo/kotlin/com/google/samples/apps/nowinandroid/ui/SnackbarInsetsScreenshotTests.kt:267 |
| `TrackDisposableJank` | Composable |  | 1 | `TrackDisposableJank(navigationState.currentKey) { metricsHolder ->` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaAppState.kt:115 |
| … | | | | 共 82 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| button | NiaButton（Composable，用 10 次）<br>NiaTextButton（Composable，用 10 次）<br>NiaOutlinedButton（Composable，用 9 次）<br>IconButtonPreviewUnchecked（Composable，用 0 次）<br>TextButton（组件库，用 3 次）<br>Button（组件库，用 1 次）<br>OutlinedButton（组件库，用 1 次）<br>BookmarkButton（Composable，用 3 次） | `NiaButton` |
| loading | NiaLoadingWheel（Composable，用 6 次）<br>NiaOverlayLoadingWheel（Composable，用 2 次）<br>CircularProgressIndicator（组件库，用 2 次）<br>TopicScreenLoading（Composable，用 0 次）<br>ForYouScreenPopulatedAndLoading（Composable，用 0 次）<br>InterestsScreenLoading（Composable，用 0 次）<br>ForYouScreenLoading（Composable，用 0 次） | `NiaLoadingWheel` |
| switch | NiaIconToggleButton（Composable，用 9 次）<br>NiaViewToggleButton（Composable，用 5 次）<br>ViewTogglePreviewExpanded（Composable，用 0 次）<br>ViewTogglePreviewCompact（Composable，用 0 次）<br>RadioButton（组件库，用 1 次）<br>Switch（Compose Material3，项目未使用） | `NiaIconToggleButton` |
| tag | NiaFilterChip（Composable，用 6 次）<br>NiaTopicTag（Composable，用 5 次）<br>NiaStatusTag（Composable，用 2 次）<br>FilterChip（组件库，用 1 次） | `NiaFilterChip` |
| tabs | NiaTab（Composable，用 2 次）<br>Tab（组件库，用 1 次）<br>TabRow（组件库，用 1 次）<br>TopicsTabContent（Composable，用 1 次） | `NiaTab` |
| tabbar | NiaNavigationSuiteScaffold（Composable，用 1 次）<br>NiaNavigationRail（Composable，用 0 次）<br>NiaNavigationBar（Composable，用 0 次）<br>NavigationBar（组件库，用 1 次） | `NiaNavigationSuiteScaffold` |
| cell | NiaListItem（Composable，用 5 次）<br>ListItem（组件库，用 1 次）<br>InterestsItem（Composable，用 7 次） | `NiaListItem` |
| card | Card（组件库，用 3 次）<br>NewsResourceCardExpanded（Composable，用 3 次） | `Card` |
| dialog | AlertDialog（组件库，用 1 次）<br>SettingsDialog（Composable，用 4 次） | `AlertDialog` |
| image | DynamicAsyncImage（Composable，用 4 次）<br>NewsResourceHeaderImage（Composable，用 1 次） | `DynamicAsyncImage` |
| navbar | NiaTopAppBar（Composable，用 3 次）<br>CenterAlignedTopAppBar（组件库，用 2 次） | `NiaTopAppBar` |
| input | TextField（组件库，用 1 次）<br>EmptySearchResultBody（Composable，用 2 次） | `TextField` |
| empty | InterestsDetailPlaceholder（Composable，用 2 次）<br>InterestsScreenEmpty（Composable，用 0 次） | `InterestsDetailPlaceholder` |

## 6. 硬编码颜色热点

| 颜色 | 次数 | 已有 token | 出现在 |
|---|---|---|---|
| #006780 | 1 | `@color/ic_launcher_foreground_tint` | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #D3E8D3 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #0EE37C | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #FFB59B | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #7F747C | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #A2CED9 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #690005 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #93000A | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #FFB4AB | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #001F28 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #727971 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #5B1A00 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
