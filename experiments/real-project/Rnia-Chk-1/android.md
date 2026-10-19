
## 页面复用检查：8 个文件（core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt、core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/TopAppBar.kt …）  ✅ PASS

- 复用的现成组件（6）：Button（组件库）、CenterAlignedTopAppBar（组件库）、HorizontalDivider（组件库）、NiaBackground（Composable）、OutlinedButton（组件库）、TextButton（组件库）
- 自造轮子（1）
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/StatusTag.kt：新建了 NiaStatusTag；同类已有 NiaFilterChip（Composable），用途可能不同（filter），请确认
- 有组件却没用（0）
- 疑似手写已有组件（0）
- 写死颜色（11）
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/NiaCommerceColors.kt:28 Color(0xFFF5F6F7)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/NiaCommerceColors.kt:29 Color(0xFFFFFFFF)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/NiaCommerceColors.kt:30 Color(0xFF1F2329)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/NiaCommerceColors.kt:31 Color(0xFF646A73)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/NiaCommerceColors.kt:32 Color(0xFF8F959E)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/NiaCommerceColors.kt:33 Color(0xFFE5E6EB)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/NiaCommerceColors.kt:34 Color(0xFFF2F3F5)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/NiaCommerceColors.kt:35 Color(0xFFFF4D1A)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/NiaCommerceColors.kt:36 Color(0xFFF53F3F)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/NiaCommerceColors.kt:37 Color(0xFFFF8F1F)：写死颜色
  - … 共 11 处

---

# 存量复用体检：/home/user/exp/Rnia-Chk-1（android）

技术栈：Jetpack Compose / Compose Foundation / Compose Material3 · 扫描文件 250

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（156 / 156） |
| 重复概念（同一 UI 概念有多个实现） | 1 |
| 重复结构（可抽成组件的复制粘贴结构） | 0 组 |
| 硬编码颜色 | 73 处（其中 4 处已有同值 token） |

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Compose Material3 | 21 | HorizontalDivider×3、TextButton×3、CircularProgressIndicator×2、Card×2、CenterAlignedTopAppBar×1、Button×1、TextField×1、OutlinedButton×1 |
| Compose Foundation | 10 | LazyVerticalStaggeredGrid×5、LazyColumn×4、LazyHorizontalGrid×1 |

**组件库里有、但本项目不能直接用的组件**（match 不会推荐它们）：

- `Divider`：已废弃，改用 HorizontalDivider

**设计系统封装**（组件库组件 → 项目的封装；新页面一律用右边的封装）：`Button` → `NiaButton`、`CenterAlignedTopAppBar` → `NiaTopAppBar`、`FilterChip` → `NiaFilterChip`、`NavigationBar` → `NiaNavigationBar`、`OutlinedButton` → `NiaOutlinedButton`、`TabRow` → `NiaTabRow`、`Tab` → `NiaTab`、`TextButton` → `NiaTextButton`

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `LazyVerticalStaggeredGrid` | 组件库 |  | 5 | `LazyVerticalStaggeredGrid(columns = StaggeredGridCells.Adaptive(300.dp)) {` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsFeed.kt:129 |
| `LazyColumn` | 组件库 | list | 4 | `LazyColumn(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:67 |
| `TextButton` | 组件库 | button | 3 | `TextButton(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:204 |
| `HorizontalDivider` | 组件库 | divider | 3 | `HorizontalDivider(color = NiaCommerceColors.Divider)` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/ActionBar.kt:53 |
| `CircularProgressIndicator` | 组件库 | loading | 2 | `CircularProgressIndicator(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/DynamicAsyncImage.kt:72 |
| `Card` | 组件库 | card | 2 | `Card(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:113 |
| `TabRow` | 组件库 | tabs | 1 | `TabRow(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Tabs.kt:88 |
| `TextField` | 组件库 | input | 1 | `TextField(` feature/search/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/search/impl/SearchScreen.kt:484 |
| `Tab` | 组件库 | tabs | 1 | `Tab(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Tabs.kt:55 |
| `Button` | 组件库 | button | 1 | `Button(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:64 |
| `NavigationBar` | 组件库 | tabbar | 1 | `NavigationBar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Navigation.kt:101 |
| `LazyHorizontalGrid` | 组件库 |  | 1 | `LazyHorizontalGrid(` feature/foryou/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/foryou/impl/ForYouScreen.kt:333 |
| `CenterAlignedTopAppBar` | 组件库 | navbar | 1 | `CenterAlignedTopAppBar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/TopAppBar.kt:52 |
| `AlertDialog` | 组件库 | dialog | 1 | `AlertDialog(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:106 |
| `RadioButton` | 组件库 | radio | 1 | `RadioButton(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:242 |
| `FilterChip` | 组件库 | tag | 1 | `FilterChip(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Chip.kt:52 |
| `OutlinedButton` | 组件库 | button | 1 | `OutlinedButton(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:131 |
| `ListItem` | 组件库 | cell | 1 | `ListItem(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:55 |
| `NiaTextButton` | Composable | button | 10 | `NiaTextButton(onClick = {}) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:87 |
| `NiaIconToggleButton` | Composable | button | 9 | `NiaIconToggleButton(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:204 |
| `NiaBackground` | Composable |  | 8 | `NiaBackground(modifier = modifier) {` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:98 |
| `NiaButton` | Composable | button | 8 | `NiaButton(onClick = {}) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:81 |
| `SettingsDialogThemeChooserRow` | Composable | dialog | 7 | `SettingsDialogThemeChooserRow(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:168 |
| `InterestsItem` | Composable | cell | 7 | `InterestsItem(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:127 |
| `NiaLoadingWheel` | Composable | loading | 6 | `NiaLoadingWheel(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/LoadingWheel.kt:146 |
| `NiaOutlinedButton` | Composable | button | 6 | `NiaOutlinedButton(onClick = {}) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:84 |
| `NiaFilterChip` | Composable | tag | 6 | `NiaFilterChip(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:175 |
| `NiaViewToggleButton` | Composable | button | 5 | `NiaViewToggleButton(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:277 |
| `NiaTopicTag` | Composable | tag | 5 | `NiaTopicTag(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:302 |
| `NiaGradientBackground` | Composable |  | 4 | `NiaGradientBackground(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:99 |
| `SettingsDialog` | Composable | dialog | 4 | `SettingsDialog(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:153 |
| `TrackScrollJank` | Composable |  | 4 | `TrackScrollJank(scrollableState = scrollableState, stateName = "bookmarks:grid")` feature/bookmarks/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/bookmarks/impl/BookmarksScreen.kt:172 |
| `BookmarkButton` | Composable | button | 3 | `BookmarkButton(isBookmarked, onToggleBookmark)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:152 |
| `NewsResourceCardExpanded` | Composable | card | 3 | `NewsResourceCardExpanded(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsFeed.kt:66 |
| `DynamicAsyncImage` | Composable | image | 3 | `DynamicAsyncImage(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:114 |
| `EmptySearchResultBody` | Composable | search | 2 | `EmptySearchResultBody(` feature/search/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/search/impl/SearchScreen.kt:173 |
| `Scrollbar` | Composable |  | 2 | `Scrollbar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/scrollbar/AppScrollbars.kt:84 |
| `NiaTabRow` | Composable | tabs | 2 | `NiaTabRow(selectedTabIndex = selectedTabIndex) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:324 |
| `NiaTab` | Composable | tabs | 2 | `NiaTab(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:326 |
| `NiaTopAppBar` | Composable | navbar | 2 | `NiaTopAppBar(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:229 |
| `NiaOverlayLoadingWheel` | Composable | loading | 2 | `NiaOverlayLoadingWheel(contentDesc = "LoadingWheel")` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/LoadingWheel.kt:167 |
| `InterestsDetailPlaceholder` | Composable | empty | 2 | `InterestsDetailPlaceholder()` feature/interests/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/interests/impl/InterestsDetailPlaceholder.kt:70 |
| `NiaNavigationSuiteScaffold` | Composable | tabbar | 1 | `NiaNavigationSuiteScaffold(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:162 |
| `NiaCatalog` | Composable |  | 1 | `setContent { NiaCatalog() }` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/NiaCatalogActivity.kt:31 |
| `NewsResourceShortDescription` | Composable |  | 1 | `NewsResourceShortDescription(userNewsResource.content)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:166 |
| `TopicIcon` | Composable | icon | 1 | `TopicIcon(` feature/foryou/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/foryou/impl/ForYouScreen.kt:400 |
| `NiaNavigationRailItem` | Composable | tabbar | 1 | `NiaNavigationRailItem(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Navigation.kt:324 |
| `TrackJank` | Composable |  | 1 | `TrackJank(scrollableState) { metricsHolder ->` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/JankStatsExtensions.kt:82 |
| `NotificationDot` | Composable | badge | 1 | `NotificationDot(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:157 |
| `TrackDisposableJank` | Composable |  | 1 | `TrackDisposableJank(navigationState.currentKey) { metricsHolder ->` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaAppState.kt:115 |
| `NewsResourceMetaData` | Composable |  | 1 | `NewsResourceMetaData(userNewsResource.publishDate, userNewsResource.type)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:163 |
| `DebugVisibleWindowInsets` | Composable |  | 1 | `DebugVisibleWindowInsets()` app/src/testDemo/kotlin/com/google/samples/apps/nowinandroid/ui/SnackbarInsetsScreenshotTests.kt:267 |
| `NewsResourceTopics` | Composable |  | 1 | `NewsResourceTopics(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:168 |
| `NewsResourceTitle` | Composable |  | 1 | `NewsResourceTitle(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:137 |
| `NewsResourceHeaderImage` | Composable | image | 1 | `NewsResourceHeaderImage(userNewsResource.headerImageUrl)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:128 |
| `NiaStatusTag` | Composable | tag | 0 | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/StatusTag.kt |
| `NiaNoticeBar` | Composable | notice | 0 | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/NoticeBar.kt |
| `NiaPrice` | Composable | price | 0 | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Price.kt |
| `NiaActionBar` | Composable | actionbar | 0 | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/ActionBar.kt |
| `NiaNavigationRail` | Composable | tabbar | 0 | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Navigation.kt |
| … | | | | 共 61 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| loading | NiaLoadingWheel（Composable，用 6 次）<br>NiaOverlayLoadingWheel（Composable，用 2 次）<br>CircularProgressIndicator（组件库，用 2 次） | `NiaLoadingWheel` |

## 6. 硬编码颜色热点

| 颜色 | 次数 | 已有 token | 出现在 |
|---|---|---|---|
| #5DD5FC | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #ECDFE0 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #727971 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #4D444C | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #003544 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #36003C | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #560A5D | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #702776 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #362F30 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #646A73 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/NiaCommerceColors.kt |
| #006D36 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #5B1A00 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
