
## 页面复用检查：7 个文件（app/src/main/java/org/schabi/newpipe/orderdetail/OrderDetailFragment.kt、app/src/main/java/org/schabi/newpipe/orderdetail/OrderDetailMock.kt …）  ❌ FAIL

- 复用的现成组件（3）：@layout/order_detail_cell（业务组件）、ImageView（组件库）、TextView（组件库）
- 自造轮子（0）
- 有组件却没用（0）
- 疑似手写已有组件（2）
  - ✗ app/src/main/res/layout/fragment_order_detail.xml:18 页面手写了「navbar」（.nav_bar），已有 @layout/toolbar_layout（业务组件 include）（项目里已用 6 次）
  - ✗ app/src/main/res/layout/fragment_order_detail.xml:90 页面手写了「tag」（.status_tag），已有 Chip（组件库 Material Components）（项目里已用 1 次）
- 页面私有调色板：页面里自己定义了 13 个颜色变量（应使用项目 / 组件库已有的 token）
- 写死颜色（1）
  - · app/src/main/res/values/colors_order_detail.xml:9 #FFFFFF：新定义的颜色 order_detail_surface 与已有 token @color/light_license_background_color 同值

---

# 存量复用体检：/home/user/exp/Rnp-A-1（android）

技术栈：XML View / Jetpack Compose / Material Components / Compose Foundation / Compose Material3 · 扫描文件 653

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（638 / 638） |
| 重复概念（同一 UI 概念有多个实现） | 4 |
| 重复结构（可抽成组件的复制粘贴结构） | 8 组 |
| 硬编码颜色 | 102 处（其中 32 处已有同值 token） |
| 硬编码尺寸 | 1100 处 |

## 0. 优先处理（按影响面排序）

1. 32 处写死的颜色已有同值 token，可以直接替换（见第 6 节）
2. 「empty（View）」有 2 种实现（共用 20 次）→ 统一到 `@layout/error_panel`
3. 「navbar（Compose）」有 2 种实现（共用 13 次）→ 统一到 `TopAppBar`
4. 「navbar（View）」有 3 种实现（共用 12 次）→ 统一到 `@layout/toolbar_layout`
5. 5 个元素的结构被复制了 5 次（app/src/main/res/layout-large-land/fragment_video_detail.xml:416 等）→ 抽成组件

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Android SDK widgets | 223 | ImageView×65、TextView×44、RecyclerView×26、ImageButton×22、ProgressBar×22、Button×15、CardView×7、RadioButton×5 |
| Material Components | 31 | ShapeableImageView×16、AppBarLayout×5、FloatingActionButton×3、TabLayout×3、Chip×1、NavigationView×1、ChipGroup×1、TextInputLayout×1 |
| Compose Material3 | 17 | TextButton×5、ListItem×4、TopAppBar×3、AlertDialog×2、Switch×1、RadioButton×1、Tab×1 |
| Compose Foundation | 5 | LazyColumn×4、HorizontalPager×1 |

**组件库里有、但本项目不能直接用的组件**（match 不会推荐它们）：

- `BottomNavigationView`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃
- `CircularProgressIndicator`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃
- `Divider`：已废弃，改用 HorizontalDivider
- `ExtendedFloatingActionButton`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃
- `LinearProgressIndicator`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃
- `MaterialButtonToggleGroup`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃
- `MaterialButton`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃
- `MaterialCardView`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃
- `MaterialDivider`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃
- `MaterialSwitch`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃
- `SearchBar`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃
- `SearchView`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃
- `Slider`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃

**项目惯例写法**（没有组件、但项目里固定这么写；新设计稿里没有对应组件时沿用它，不要另造）：

- divider：1px View 细线，14 处，例 `<View android:layout_width="match_parent" android:layout_height="1dp" android:background="?attr/separator_color" />` app/src/main/res/layout/dialog_feed_group_reorder.xml:19

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `ImageView` | 组件库 | image | 65 | `<ImageView` app/src/main/res/layout/comment_replies_header.xml:54 |
| `TextView` | 组件库 |  | 44 | `<TextView` app/src/main/res/layout/fragment_order_detail.xml:35 |
| `RecyclerView` | 组件库 | list | 26 | `<androidx.recyclerview.widget.RecyclerView` app/src/main/res/layout/activity_player_queue_control.xml:30 |
| `ImageButton` | 组件库 | button | 22 | `<ImageButton` app/src/main/res/layout/activity_player_queue_control.xml:165 |
| `ProgressBar` | 组件库 | loading | 22 | `<ProgressBar` app/src/main/res/layout/activity_player_queue_control.xml:226 |
| `ShapeableImageView` | 组件库 | image | 16 | `<com.google.android.material.imageview.ShapeableImageView` app/src/main/res/layout/comment_replies_header.xml:8 |
| `Button` | 组件库 | button | 15 | `<Button` app/src/main/res/layout/activity_error.xml:118 |
| `CardView` | 组件库 | card | 7 | `<androidx.cardview.widget.CardView xmlns:android="http://schemas.android.com/apk/res/android"` app/src/main/res/layout/feed_group_add_new_grid_item.xml:2 |
| `Toolbar` | 组件库 | navbar | 5 | `<androidx.appcompat.widget.Toolbar` app/src/main/res/layout/activity_player_queue_control.xml:18 |
| `TextButton` | 组件库 | button | 5 | `TextButton(onClick = { showDialog = false }) {` shared/src/commonMain/kotlin/net/newpipe/app/composable/ListPreference.kt:104 |
| `AppBarLayout` | 组件库 |  | 5 | `<com.google.android.material.appbar.AppBarLayout` app/src/main/res/layout/activity_player_queue_control.xml:11 |
| `RadioButton` | 组件库 | radio | 5 | `<RadioButton` app/src/main/res/layout/download_dialog.xml:43 |
| `ListItem` | 组件库 | cell | 4 | `ListItem(` shared/src/commonMain/kotlin/net/newpipe/app/composable/ListPreference.kt:68 |
| `LazyColumn` | 组件库 | list | 4 | `LazyColumn(` shared/src/commonMain/kotlin/net/newpipe/app/screen/about/AboutPage.kt:78 |
| `ViewPager` | 组件库 | carousel | 4 | `<androidx.viewpager.widget.ViewPager` app/src/main/res/layout/fragment_channel.xml:156 |
| `Spinner` | 组件库 | loading | 4 | `<Spinner` app/src/main/res/layout/download_dialog.xml:63 |
| `TopAppBar` | 组件库 | navbar | 3 | `fun TopAppBar(` shared/src/commonMain/kotlin/net/newpipe/app/composable/TopAppBar.kt:37 |
| `TabLayout` | 组件库 | tabs | 3 | `<com.google.android.material.tabs.TabLayout` app/src/main/res/layout/fragment_channel.xml:145 |
| `CheckBox` | 组件库 | checkbox | 3 | `<CheckBox` app/src/main/res/layout/dialog_playback_parameter.xml:503 |
| `FloatingActionButton` | 组件库 | button | 3 | `<com.google.android.material.floatingactionbutton.FloatingActionButton` app/src/main/res/layout/fragment_choose_tabs.xml:35 |
| `AlertDialog` | 组件库 | dialog | 2 | `AlertDialog(` shared/src/commonMain/kotlin/net/newpipe/app/composable/ListPreference.kt:77 |
| `Tab` | 组件库 | tabs | 1 | `Tab(` shared/src/commonMain/kotlin/net/newpipe/app/screen/about/AboutScreen.kt:87 |
| `SeekBar` | 组件库 |  | 1 | `<SeekBar` app/src/main/res/layout/download_dialog.xml:137 |
| `Switch` | 组件库 | switch | 1 | `Switch(` shared/src/commonMain/kotlin/net/newpipe/app/composable/SwitchPreference.kt:64 |
| `HorizontalPager` | 组件库 | carousel | 1 | `HorizontalPager(modifier = Modifier.fillMaxSize(), state = pagerState) { page ->` shared/src/commonMain/kotlin/net/newpipe/app/screen/about/AboutScreen.kt:104 |
| `ChipGroup` | 组件库 | tag | 1 | `<com.google.android.material.chip.ChipGroup` app/src/main/res/layout/item_metadata_tags.xml:22 |
| `Chip` | 组件库 | tag | 1 | `<com.google.android.material.chip.Chip xmlns:android="http://schemas.android.com/apk/res/android"` app/src/main/res/layout/chip.xml:4 |
| `NavigationView` | 组件库 |  | 1 | `<com.google.android.material.navigation.NavigationView xmlns:android="http://schemas.android.com/apk/res/android"` app/src/main/res/layout/drawer_layout.xml:2 |
| `TextInputLayout` | 组件库 |  | 1 | `<com.google.android.material.textfield.TextInputLayout` app/src/main/res/layout/dialog_feed_group_create.xml:34 |
| `RadioButton` | 组件库 | radio | 1 | `RadioButton(selected = entry.value == selectedValue, onClick = null)` shared/src/commonMain/kotlin/net/newpipe/app/composable/ListPreference.kt:97 |
| `SwitchPreference` | Composable |  | 25 | `SwitchPreference(` shared/src/androidMain/kotlin/net/newpipe/app/screen/ExoPlayerSettingsSections.kt:112 |
| `ListPreference` | Composable |  | 20 | `ListPreference(` shared/src/androidMain/kotlin/net/newpipe/app/screen/ExoPlayerSettingsSections.kt:102 |
| `PreferenceRow` | Composable |  | 12 | `PreferenceRow(` shared/src/commonMain/kotlin/net/newpipe/app/composable/PreferenceRow.kt:62 |
| `TopAppBar` | Composable | navbar | 10 | `TopAppBar(` shared/src/commonMain/kotlin/net/newpipe/app/composable/TopAppBar.kt:45 |
| `LicenseDialog` | Composable | dialog | 3 | `LicenseDialog(` shared/src/commonMain/kotlin/net/newpipe/app/screen/about/LicenseDialog.kt:92 |
| `LinkListItem` | Composable | cell | 3 | `LinkListItem(link = link)` shared/src/commonMain/kotlin/net/newpipe/app/composable/about/LinkListItem.kt:64 |
| `LibraryListItem` | Composable | cell | 3 | `LibraryListItem(library = library)` shared/src/commonMain/kotlin/net/newpipe/app/composable/about/LibraryListItem.kt:72 |
| `PreferenceCategoryTitle` | Composable |  | 3 | `PreferenceCategoryTitle(title = "Player")` shared/src/commonMain/kotlin/net/newpipe/app/composable/PreferenceCategoryTitle.kt:35 |
| `NavDisplay` | Composable |  | 2 | `NavDisplay(` shared/src/commonMain/kotlin/net/newpipe/app/App.kt:37 |
| `ExoPlayerSettingsSections` | Composable | card | 1 | `override fun Render() = ExoPlayerSettingsSections()` shared/src/androidMain/kotlin/net/newpipe/app/platform/AndroidPlayerSettingsSections.kt:17 |
| `CircleClipTapView` | 自定义 View |  | 1 | `<org.schabi.newpipe.views.player.CircleClipTapView` app/src/main/res/layout/player_fast_seek_overlay.xml:8 |
| `NewPipeTextView` | 自定义 View |  | 261 | `<org.schabi.newpipe.views.NewPipeTextView` app/src/main/res/layout/activity_error.xml:29 |
| `Render` | Composable |  | 0 | shared/src/commonMain/kotlin/net/newpipe/app/platform/PlayerSettingsSections.kt |
| `LicensePage` | Composable |  | 1 | `Page.LICENSE -> LicensePage()` shared/src/commonMain/kotlin/net/newpipe/app/screen/about/AboutScreen.kt:52 |
| `ScrollableTabLayout` | 自定义 View |  | 1 | `<org.schabi.newpipe.views.ScrollableTabLayout` app/src/main/res/layout/fragment_main.xml:8 |
| `PlayerFastSeekOverlay` | 自定义 View |  | 2 | `<org.schabi.newpipe.views.player.PlayerFastSeekOverlay` app/src/main/res/layout/player.xml:824 |
| `SecondsView` | 自定义 View |  | 1 | `<org.schabi.newpipe.views.player.SecondsView` app/src/main/res/layout/player_fast_seek_overlay.xml:15 |
| `CustomCollapsingToolbarLayout` | 自定义 View |  | 3 | `<org.schabi.newpipe.views.CustomCollapsingToolbarLayout` app/src/main/res/layout/fragment_channel.xml:17 |
| `NewPipeEditText` | 自定义 View | input | 6 | `<org.schabi.newpipe.views.NewPipeEditText` app/src/main/res/layout/activity_error.xml:112 |
| `NewPipeRecyclerView` | 自定义 View | list | 2 | `<org.schabi.newpipe.views.NewPipeRecyclerView` app/src/main/res/layout/fragment_comments.xml:7 |
| `ExpandableSurfaceView` | 自定义 View |  | 1 | `<org.schabi.newpipe.views.ExpandableSurfaceView` app/src/main/res/layout/player.xml:11 |
| `AboutPage` | Composable |  | 1 | `Page.ABOUT -> AboutPage()` shared/src/commonMain/kotlin/net/newpipe/app/screen/about/AboutScreen.kt:51 |
| `FocusAwareSeekBar` | 自定义 View |  | 0 | app/src/main/java/org/schabi/newpipe/views/FocusAwareSeekBar.java |
| `@layout/error_panel` | 业务组件 | empty | 14 | include |
| `@layout/list_empty_view` | 业务组件 | empty | 6 | include |
| `@layout/order_detail_cell` | 业务组件 | cell | 6 | include |
| `@layout/toolbar_layout` | 业务组件 | navbar | 6 | include |
| `@layout/settings_notification_action` | 业务组件 |  | 5 | include |
| `@layout/playlist_control` | 业务组件 |  | 3 | include |
| `@layout/toolbar_search_layout` | 业务组件 |  | 2 | include |
| … | | | | 共 63 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| navbar（View） | @layout/toolbar_layout（业务组件，用 6 次）<br>Toolbar（组件库，用 5 次）<br>@layout/settings_category_header_title（业务组件，用 1 次） | `@layout/toolbar_layout` |
| empty（View） | @layout/error_panel（业务组件，用 14 次）<br>@layout/list_empty_view（业务组件，用 6 次） | `@layout/error_panel` |
| navbar（Compose） | TopAppBar（组件库，用 3 次）<br>TopAppBar（Composable，用 10 次） | `TopAppBar` |
| card（Compose） | Card（Compose Material3，项目未使用）<br>ExoPlayerSettingsSections（Composable，用 1 次） | `Card` |

## 5. 重复结构（复制粘贴的 UI 片段，可抽成组件）

- **7 个元素 × 4 处**：app/src/main/res/layout-land/list_stream_card_item.xml:2、app/src/main/res/layout/list_stream_card_item.xml:2、app/src/main/res/layout/list_stream_grid_item.xml:2、app/src/main/res/layout/list_stream_item.xml:2

```
<androidx.constraintlayout.widget.ConstraintLayout>
  <ImageView>
  <org.schabi.newpipe.views.NewPipeTextView>
  <org.schabi.newpipe.views.NewPipeTextView>
  <org.schabi.newpipe.views.NewPipeTextView>
  <org.schabi.newpipe.views.NewPipeTextView>
  <org.schabi.newpipe.views.AnimatedProgressBar>
```
- **5 个元素 × 5 处**：app/src/main/res/layout-large-land/fragment_video_detail.xml:416、app/src/main/res/layout-large-land/fragment_video_detail.xml:492、app/src/main/res/layout/dialog_playback_parameter.xml:406、app/src/main/res/layout/fragment_video_detail.xml:403、app/src/main/res/layout/fragment_video_detail.xml:477

```
<LinearLayout>
  <org.schabi.newpipe.views.NewPipeTextView>
  <org.schabi.newpipe.views.NewPipeTextView>
  <org.schabi.newpipe.views.NewPipeTextView>
  <org.schabi.newpipe.views.NewPipeTextView>
```
- **8 个元素 × 3 处**：app/src/main/res/layout/fragment_channel_videos.xml:2、app/src/main/res/layout/fragment_playlist.xml:2、app/src/main/res/layout/fragment_related_items.xml:2

```
<RelativeLayout>
  <androidx.recyclerview.widget.RecyclerView>
  <ProgressBar>
  <LinearLayout>
    <org.schabi.newpipe.views.NewPipeTextView>
    <org.schabi.newpipe.views.NewPipeTextView>
  <include>
  <View>
```
- **5 个元素 × 4 处**：app/src/main/res/layout/list_playlist_grid_item.xml:2、app/src/main/res/layout/list_playlist_item.xml:2、app/src/main/res/layout/list_playlist_mini_item.xml:2、app/src/main/res/layout/stream_quality_item.xml:2
- **8 个元素 × 2 处**：app/src/main/res/layout/fragment_comments.xml:2、app/src/main/res/layout/fragment_kiosk.xml:2
- **5 个元素 × 3 处**：app/src/main/res/layout/select_channel_fragment.xml:2、app/src/main/res/layout/select_feed_group_fragment.xml:2、app/src/main/res/layout/select_playlist_fragment.xml:2
- **7 个元素 × 2 处**：app/src/main/res/layout/list_stream_playlist_grid_item.xml:2、app/src/main/res/layout/list_stream_playlist_item.xml:2
- **5 个元素 × 2 处**：app/src/main/res/layout/fragment_bookmarks.xml:2、app/src/main/res/layout/fragment_subscription.xml:2

## 6. 硬编码颜色热点

| 颜色 | 次数 | 已有 token | 出现在 |
|---|---|---|---|
| #FFFFFF | 15 | `@color/light_license_background_color` | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt、shared/src/commonMain/kotlin/net/newpipe/app/theme/ServiceTheme.kt |
| #FFDAD6 | 6 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #F1DEDC | 4 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #1A1110 | 3 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #FFF8F7 | 3 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #231918 | 3 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #534342 | 3 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #000000 | 2 | `@color/black_background_color` | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #FFB3AC | 2 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #23000000 | 2 | — | app/src/main/res/layout/fragment_video_detail.xml、app/src/main/res/layout-large-land/fragment_video_detail.xml |
| #D8C2BF | 2 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #392E2D | 2 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
