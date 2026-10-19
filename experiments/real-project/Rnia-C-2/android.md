
## 页面复用检查：15 个文件（app/src/main/res/values/strings.xml、core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt …）  ✅ PASS

- 复用的现成组件（8）：Button（组件库）、Card（组件库）、CenterAlignedTopAppBar（组件库）、DynamicAsyncImage（Composable）、HorizontalDivider（组件库）、NiaBackground（Composable）、OutlinedButton（组件库）、TextButton（组件库）
- 自造轮子（0）
- 有组件却没用（0）
- 疑似手写已有组件（0）
- 写死颜色（33）
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt:51 Color(0xFFFF4D1A)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt:53 Color(0xFFFFEDE8)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt:54 Color(0xFF3370FF)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt:55 Color(0xFFE1EAFF)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt:56 Color(0xFF00B42A)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt:57 Color(0xFFE8FFEA)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt:58 Color(0xFFFF8F1F)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt:59 Color(0xFFFFF3E8)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt:60 Color(0xFFF53F3F)：写死颜色
  - · core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt:61 Color(0xFFFFECE8)：写死颜色
  - … 共 33 处

---

# 存量复用体检：/home/user/exp/Rnia-C-2（android）

技术栈：Jetpack Compose / Compose Foundation / Compose Material3 · 扫描文件 253

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（186 / 186） |
| 重复概念（同一 UI 概念有多个实现） | 2 |
| 重复结构（可抽成组件的复制粘贴结构） | 0 组 |
| 硬编码颜色 | 95 处（其中 3 处已有同值 token） |

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Compose Material3 | 25 | HorizontalDivider×5、TextButton×3、Card×3、CircularProgressIndicator×2、CenterAlignedTopAppBar×2、AlertDialog×1、RadioButton×1、TextField×1 |
| Compose Foundation | 10 | LazyVerticalStaggeredGrid×5、LazyColumn×4、LazyHorizontalGrid×1 |

**组件库里有、但本项目不能直接用的组件**（match 不会推荐它们）：

- `Divider`：已废弃，改用 HorizontalDivider

**设计系统封装**（组件库组件 → 项目的封装；新页面一律用右边的封装）：`Button` → `NiaButton`、`CenterAlignedTopAppBar` → `NiaTopAppBar`、`FilterChip` → `NiaFilterChip`、`NavigationBar` → `NiaNavigationBar`、`OutlinedButton` → `NiaOutlinedButton`、`TabRow` → `NiaTabRow`、`Tab` → `NiaTab`、`TextButton` → `NiaTextButton`

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `LazyVerticalStaggeredGrid` | 组件库 |  | 5 | `LazyVerticalStaggeredGrid(columns = StaggeredGridCells.Adaptive(300.dp)) {` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsFeed.kt:129 |
| `HorizontalDivider` | 组件库 | divider | 5 | `HorizontalDivider(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:138 |
| `LazyColumn` | 组件库 | list | 4 | `LazyColumn(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:67 |
| `TextButton` | 组件库 | button | 3 | `TextButton(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:205 |
| `Card` | 组件库 | card | 3 | `Card(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:269 |
| `CenterAlignedTopAppBar` | 组件库 | navbar | 2 | `CenterAlignedTopAppBar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/TopAppBar.kt:52 |
| `CircularProgressIndicator` | 组件库 | loading | 2 | `CircularProgressIndicator(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/DynamicAsyncImage.kt:72 |
| `RadioButton` | 组件库 | radio | 1 | `RadioButton(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:242 |
| `AlertDialog` | 组件库 | dialog | 1 | `AlertDialog(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:106 |
| `ListItem` | 组件库 | cell | 1 | `ListItem(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:55 |
| `OutlinedButton` | 组件库 | button | 1 | `OutlinedButton(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:134 |
| `Button` | 组件库 | button | 1 | `Button(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Button.kt:64 |
| `NavigationBar` | 组件库 | tabbar | 1 | `NavigationBar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Navigation.kt:101 |
| `TabRow` | 组件库 | tabs | 1 | `TabRow(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Tabs.kt:88 |
| `Tab` | 组件库 | tabs | 1 | `Tab(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Tabs.kt:55 |
| `FilterChip` | 组件库 | tag | 1 | `FilterChip(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/Chip.kt:52 |
| `TextField` | 组件库 | input | 1 | `TextField(` feature/search/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/search/impl/SearchScreen.kt:484 |
| `LazyHorizontalGrid` | 组件库 |  | 1 | `LazyHorizontalGrid(` feature/foryou/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/foryou/impl/ForYouScreen.kt:333 |
| `NiaTextButton` | Composable | button | 10 | `NiaTextButton(onClick = {}) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:87 |
| `NiaButton` | Composable | button | 10 | `NiaButton(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:245 |
| `NiaIconToggleButton` | Composable | button | 9 | `NiaIconToggleButton(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:204 |
| `NiaCell` | Composable | cell | 9 | `NiaCell(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:143 |
| `NiaOutlinedButton` | Composable | button | 8 | `NiaOutlinedButton(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:227 |
| `NiaBackground` | Composable |  | 8 | `NiaBackground(modifier = modifier) {` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:98 |
| `InterestsItem` | Composable | cell | 7 | `InterestsItem(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:127 |
| `SettingsDialogThemeChooserRow` | Composable | dialog | 7 | `SettingsDialogThemeChooserRow(` feature/settings/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/settings/impl/SettingsDialog.kt:168 |
| `NiaFilterChip` | Composable | tag | 6 | `NiaFilterChip(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:175 |
| `NiaLoadingWheel` | Composable | loading | 6 | `NiaLoadingWheel(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/LoadingWheel.kt:146 |
| `NiaViewToggleButton` | Composable | button | 5 | `NiaViewToggleButton(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:277 |
| `NiaTopicTag` | Composable | tag | 5 | `NiaTopicTag(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:302 |
| `NiaGradientBackground` | Composable |  | 4 | `NiaGradientBackground(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:99 |
| `DynamicAsyncImage` | Composable | image | 4 | `DynamicAsyncImage(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/InterestsItem.kt:114 |
| `TrackScrollJank` | Composable |  | 4 | `TrackScrollJank(scrollableState = scrollableState, stateName = "bookmarks:grid")` feature/bookmarks/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/bookmarks/impl/BookmarksScreen.kt:172 |
| `SettingsDialog` | Composable | dialog | 4 | `SettingsDialog(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:153 |
| `NewsResourceCardExpanded` | Composable | card | 3 | `NewsResourceCardExpanded(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsFeed.kt:66 |
| `NiaPriceText` | Composable | price | 3 | `NiaPriceText(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:192 |
| `NiaTopAppBar` | Composable | navbar | 3 | `NiaTopAppBar(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:229 |
| `BookmarkButton` | Composable | button | 3 | `BookmarkButton(isBookmarked, onToggleBookmark)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:152 |
| `ProductItem` | Composable | goods | 2 | `ProductItem(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:156 |
| `NiaTabRow` | Composable | tabs | 2 | `NiaTabRow(selectedTabIndex = selectedTabIndex) {` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:324 |
| `NiaNoticeBar` | Composable | notice | 2 | `NiaNoticeBar(text = order.notice, tone = NiaNoticeTone.Warning)` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:100 |
| `NiaCountDown` | Composable | countdown | 2 | `NiaCountDown(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:127 |
| `InterestsDetailPlaceholder` | Composable | empty | 2 | `InterestsDetailPlaceholder()` feature/interests/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/interests/impl/InterestsDetailPlaceholder.kt:70 |
| `NiaOverlayLoadingWheel` | Composable | loading | 2 | `NiaOverlayLoadingWheel(contentDesc = "LoadingWheel")` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/LoadingWheel.kt:167 |
| `Scrollbar` | Composable |  | 2 | `Scrollbar(` core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/component/scrollbar/AppScrollbars.kt:84 |
| `NiaTab` | Composable | tabs | 2 | `NiaTab(` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/ui/Catalog.kt:326 |
| `EmptySearchResultBody` | Composable | search | 2 | `EmptySearchResultBody(` feature/search/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/search/impl/SearchScreen.kt:173 |
| `NiaStatusTag` | Composable | tag | 2 | `NiaStatusTag(text = order.statusTag, tone = NiaTagTone.Warning)` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/orderdetail/OrderDetailScreen.kt:118 |
| `NewsResourceMetaData` | Composable |  | 1 | `NewsResourceMetaData(userNewsResource.publishDate, userNewsResource.type)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:163 |
| `TrackDisposableJank` | Composable |  | 1 | `TrackDisposableJank(navigationState.currentKey) { metricsHolder ->` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaAppState.kt:115 |
| `TrackJank` | Composable |  | 1 | `TrackJank(scrollableState) { metricsHolder ->` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/JankStatsExtensions.kt:82 |
| `NewsResourceTitle` | Composable |  | 1 | `NewsResourceTitle(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:137 |
| `NewsResourceHeaderImage` | Composable | image | 1 | `NewsResourceHeaderImage(userNewsResource.headerImageUrl)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:128 |
| `NotificationDot` | Composable | badge | 1 | `NotificationDot(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:157 |
| `NewsResourceTopics` | Composable |  | 1 | `NewsResourceTopics(` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:168 |
| `DebugVisibleWindowInsets` | Composable |  | 1 | `DebugVisibleWindowInsets()` app/src/testDemo/kotlin/com/google/samples/apps/nowinandroid/ui/SnackbarInsetsScreenshotTests.kt:267 |
| `NewsResourceShortDescription` | Composable |  | 1 | `NewsResourceShortDescription(userNewsResource.content)` core/ui/src/main/kotlin/com/google/samples/apps/nowinandroid/core/ui/NewsResourceCard.kt:166 |
| `NiaCatalog` | Composable |  | 1 | `setContent { NiaCatalog() }` app-nia-catalog/src/main/kotlin/com/google/samples/apps/niacatalog/NiaCatalogActivity.kt:31 |
| `NiaNavigationSuiteScaffold` | Composable | tabbar | 1 | `NiaNavigationSuiteScaffold(` app/src/main/kotlin/com/google/samples/apps/nowinandroid/ui/NiaApp.kt:162 |
| `TopicIcon` | Composable | icon | 1 | `TopicIcon(` feature/foryou/impl/src/main/kotlin/com/google/samples/apps/nowinandroid/feature/foryou/impl/ForYouScreen.kt:400 |
| … | | | | 共 64 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| loading | NiaLoadingWheel（Composable，用 6 次）<br>NiaOverlayLoadingWheel（Composable，用 2 次）<br>CircularProgressIndicator（组件库，用 2 次） | `NiaLoadingWheel` |
| cell | NiaCell（Composable，用 9 次）<br>ListItem（组件库，用 1 次） | `NiaCell` |

## 6. 硬编码颜色热点

| 颜色 | 次数 | 已有 token | 出现在 |
|---|---|---|---|
| #1F2329 | 2 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt |
| #006780 | 1 | `@color/ic_launcher_foreground_tint` | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #4F6352 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #00210B | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #A23F16 | 1 | `@color/ic_launcher_foreground_tint` | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #560A5D | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #2B2F36 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt |
| #004D61 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #5B1A00 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #EDDEE8 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Color.kt |
| #F53F3F | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt |
| #373C43 | 1 | — | core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/CommerceColors.kt |
