
## 页面复用检查：13 个文件（app/src/main/res/values-zh-rCN/strings.xml、app/src/main/res/values/attrs.xml …）  ✅ PASS

- 复用的现成组件（6）：@layout/info_cell（业务组件）、Button（组件库）、CardView（组件库）、ImageView（组件库）、NewPipeTextView（自定义 View）、ShapeableImageView（组件库）
- 自造轮子（0）
- 有组件却没用（0）
- 疑似手写已有组件（0）
- 页面私有调色板：页面里自己定义了 10 个颜色变量（应使用项目 / 组件库已有的 token）
- 写死颜色（0）

---

# 存量复用体检：/home/user/exp/Rnp-C-1（android）

技术栈：XML View / Jetpack Compose / Material Components / Compose Foundation / Compose Material3 · 扫描文件 656

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（640 / 640） |
| 重复概念（同一 UI 概念有多个实现） | 5 |
| 重复结构（可抽成组件的复制粘贴结构） | 8 组 |
| 硬编码颜色 | 102 处（其中 32 处已有同值 token） |
| 硬编码尺寸 | 1075 处 |

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Android SDK widgets | 206 | ImageView×63、RecyclerView×26、TextView×23、ProgressBar×22、ImageButton×22、Button×18、CardView×10、Toolbar×5 |
| Material Components | 32 | ShapeableImageView×17、AppBarLayout×5、TabLayout×3、FloatingActionButton×3、Chip×1、NavigationView×1、ChipGroup×1、TextInputLayout×1 |
| Compose Material3 | 17 | TextButton×5、ListItem×4、TopAppBar×3、AlertDialog×2、Switch×1、Tab×1、RadioButton×1 |
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
- `Slider`：需要 Theme.MaterialComponents 主题，项目主题 OpeningTheme 继承 AppCompat，直接用会崩溃

**项目惯例写法**（没有组件、但项目里固定这么写；新设计稿里没有对应组件时沿用它，不要另造）：

- divider：1px View 细线，16 处，例 `<View android:id="@+id/separator" android:layout_width="match_parent" android:layout_height="1dp" android:background="?attr/separator_color"…` app/src/main/res/layout/dialog_feed_group_create.xml:172

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `ImageView` | 组件库 | image | 63 | `<ImageView` app/src/main/res/layout/comment_replies_header.xml:54 |
| `RecyclerView` | 组件库 | list | 26 | `<androidx.recyclerview.widget.RecyclerView` app/src/main/res/layout/activity_player_queue_control.xml:30 |
| `TextView` | 组件库 |  | 23 | `<TextView xmlns:android="http://schemas.android.com/apk/res/android"` app/src/main/res/layout/instance_spinner_item.xml:2 |
| `ImageButton` | 组件库 | button | 22 | `<ImageButton` app/src/main/res/layout/activity_player_queue_control.xml:165 |
| `ProgressBar` | 组件库 | loading | 22 | `<ProgressBar` app/src/main/res/layout/activity_player_queue_control.xml:226 |
| `Button` | 组件库 | button | 18 | `<Button` app/src/main/res/layout/activity_error.xml:118 |
| `ShapeableImageView` | 组件库 | image | 17 | `<com.google.android.material.imageview.ShapeableImageView` app/src/main/res/layout/comment_replies_header.xml:8 |
| `CardView` | 组件库 | card | 10 | `<androidx.cardview.widget.CardView xmlns:android="http://schemas.android.com/apk/res/android"` app/src/main/res/layout/feed_group_add_new_grid_item.xml:2 |
| `TextButton` | 组件库 | button | 5 | `TextButton(onClick = { showDialog = false }) {` shared/src/commonMain/kotlin/net/newpipe/app/composable/ListPreference.kt:104 |
| `AppBarLayout` | 组件库 |  | 5 | `<com.google.android.material.appbar.AppBarLayout` app/src/main/res/layout/activity_player_queue_control.xml:11 |
| `Toolbar` | 组件库 | navbar | 5 | `<androidx.appcompat.widget.Toolbar` app/src/main/res/layout/activity_player_queue_control.xml:18 |
| `RadioButton` | 组件库 | switch | 5 | `<RadioButton` app/src/main/res/layout/download_dialog.xml:43 |
| `LazyColumn` | 组件库 | list | 4 | `LazyColumn(` shared/src/commonMain/kotlin/net/newpipe/app/screen/about/AboutPage.kt:78 |
| `ViewPager` | 组件库 | carousel | 4 | `<androidx.viewpager.widget.ViewPager` app/src/main/res/layout/fragment_channel.xml:156 |
| `Spinner` | 组件库 | loading | 4 | `<Spinner` app/src/main/res/layout/download_dialog.xml:63 |
| `ListItem` | 组件库 | cell | 4 | `ListItem(` shared/src/commonMain/kotlin/net/newpipe/app/composable/ListPreference.kt:68 |
| `TopAppBar` | 组件库 | navbar | 3 | `fun TopAppBar(` shared/src/commonMain/kotlin/net/newpipe/app/composable/TopAppBar.kt:37 |
| `TabLayout` | 组件库 | tabs | 3 | `<com.google.android.material.tabs.TabLayout` app/src/main/res/layout/fragment_channel.xml:145 |
| `FloatingActionButton` | 组件库 | button | 3 | `<com.google.android.material.floatingactionbutton.FloatingActionButton` app/src/main/res/layout/fragment_choose_tabs.xml:35 |
| `CheckBox` | 组件库 | switch | 3 | `<CheckBox` app/src/main/res/layout/dialog_playback_parameter.xml:503 |
| `AlertDialog` | 组件库 | dialog | 2 | `AlertDialog(` shared/src/commonMain/kotlin/net/newpipe/app/composable/ListPreference.kt:77 |
| `TextInputLayout` | 组件库 |  | 1 | `<com.google.android.material.textfield.TextInputLayout` app/src/main/res/layout/dialog_feed_group_create.xml:34 |
| `Switch` | 组件库 | switch | 1 | `Switch(` shared/src/commonMain/kotlin/net/newpipe/app/composable/SwitchPreference.kt:64 |
| `RadioButton` | 组件库 | switch | 1 | `RadioButton(selected = entry.value == selectedValue, onClick = null)` shared/src/commonMain/kotlin/net/newpipe/app/composable/ListPreference.kt:97 |
| `SeekBar` | 组件库 |  | 1 | `<SeekBar` app/src/main/res/layout/download_dialog.xml:137 |
| `Chip` | 组件库 | tag | 1 | `<com.google.android.material.chip.Chip xmlns:android="http://schemas.android.com/apk/res/android"` app/src/main/res/layout/chip.xml:4 |
| `Tab` | 组件库 | tabs | 1 | `Tab(` shared/src/commonMain/kotlin/net/newpipe/app/screen/about/AboutScreen.kt:87 |
| `ChipGroup` | 组件库 | tag | 1 | `<com.google.android.material.chip.ChipGroup` app/src/main/res/layout/item_metadata_tags.xml:22 |
| `NavigationView` | 组件库 |  | 1 | `<com.google.android.material.navigation.NavigationView xmlns:android="http://schemas.android.com/apk/res/android"` app/src/main/res/layout/drawer_layout.xml:2 |
| `HorizontalPager` | 组件库 | carousel | 1 | `HorizontalPager(modifier = Modifier.fillMaxSize(), state = pagerState) { page ->` shared/src/commonMain/kotlin/net/newpipe/app/screen/about/AboutScreen.kt:104 |
| `NewPipeTextView` | 自定义 View |  | 276 | `<org.schabi.newpipe.views.NewPipeTextView` app/src/main/res/layout/activity_error.xml:29 |
| `NewPipeEditText` | 自定义 View | input | 6 | `<org.schabi.newpipe.views.NewPipeEditText` app/src/main/res/layout/activity_error.xml:112 |
| `CustomCollapsingToolbarLayout` | 自定义 View |  | 3 | `<org.schabi.newpipe.views.CustomCollapsingToolbarLayout` app/src/main/res/layout/fragment_channel.xml:17 |
| `PlayerFastSeekOverlay` | 自定义 View |  | 2 | `<org.schabi.newpipe.views.player.PlayerFastSeekOverlay` app/src/main/res/layout/player.xml:824 |
| `TopAppBar` | Composable | navbar | 10 | `TopAppBar(` shared/src/commonMain/kotlin/net/newpipe/app/composable/TopAppBar.kt:45 |
| `ScrollableTabLayout` | 自定义 View |  | 1 | `<org.schabi.newpipe.views.ScrollableTabLayout` app/src/main/res/layout/fragment_main.xml:8 |
| `SecondsView` | 自定义 View |  | 1 | `<org.schabi.newpipe.views.player.SecondsView` app/src/main/res/layout/player_fast_seek_overlay.xml:15 |
| `NewPipeRecyclerView` | 自定义 View | list | 2 | `<org.schabi.newpipe.views.NewPipeRecyclerView` app/src/main/res/layout/fragment_comments.xml:7 |
| `ExpandableSurfaceView` | 自定义 View |  | 1 | `<org.schabi.newpipe.views.ExpandableSurfaceView` app/src/main/res/layout/player.xml:11 |
| `PreferenceRow` | Composable | cell | 12 | `PreferenceRow(` shared/src/commonMain/kotlin/net/newpipe/app/composable/PreferenceRow.kt:62 |
| `ListPreference` | Composable | list | 20 | `ListPreference(` shared/src/androidMain/kotlin/net/newpipe/app/screen/ExoPlayerSettingsSections.kt:102 |
| `SwitchPreference` | Composable | switch | 25 | `SwitchPreference(` shared/src/androidMain/kotlin/net/newpipe/app/screen/ExoPlayerSettingsSections.kt:112 |
| `PriceTextView` | 自定义 View | price | 2 | `<org.schabi.newpipe.views.PriceTextView` app/src/main/res/layout/fragment_order_detail.xml:171 |
| `FocusAwareSeekBar` | 自定义 View |  | 0 | app/src/main/java/org/schabi/newpipe/views/FocusAwareSeekBar.java |
| `CountDownTextView` | 自定义 View | countdown | 1 | `<org.schabi.newpipe.views.CountDownTextView` app/src/main/res/layout/fragment_order_detail.xml:69 |
| `CircleClipTapView` | 自定义 View |  | 1 | `<org.schabi.newpipe.views.player.CircleClipTapView` app/src/main/res/layout/player_fast_seek_overlay.xml:8 |
| `LibraryListItem` | Composable | cell | 3 | `LibraryListItem(library = library)` shared/src/commonMain/kotlin/net/newpipe/app/composable/about/LibraryListItem.kt:72 |
| `PreferenceCategoryTitle` | Composable |  | 3 | `PreferenceCategoryTitle(title = "Player")` shared/src/commonMain/kotlin/net/newpipe/app/composable/PreferenceCategoryTitle.kt:35 |
| `LinkListItem` | Composable | cell | 3 | `LinkListItem(link = link)` shared/src/commonMain/kotlin/net/newpipe/app/composable/about/LinkListItem.kt:64 |
| `LicenseDialog` | Composable | dialog | 3 | `LicenseDialog(` shared/src/commonMain/kotlin/net/newpipe/app/screen/about/LicenseDialog.kt:92 |
| `NavDisplay` | Composable |  | 2 | `NavDisplay(` shared/src/commonMain/kotlin/net/newpipe/app/App.kt:37 |
| `LicensePage` | Composable |  | 1 | `Page.LICENSE -> LicensePage()` shared/src/commonMain/kotlin/net/newpipe/app/screen/about/AboutScreen.kt:52 |
| `AboutPage` | Composable |  | 1 | `Page.ABOUT -> AboutPage()` shared/src/commonMain/kotlin/net/newpipe/app/screen/about/AboutScreen.kt:51 |
| `ExoPlayerSettingsSections` | Composable | card | 1 | `override fun Render() = ExoPlayerSettingsSections()` shared/src/androidMain/kotlin/net/newpipe/app/platform/AndroidPlayerSettingsSections.kt:17 |
| `Render` | Composable |  | 0 | shared/src/commonMain/kotlin/net/newpipe/app/platform/PlayerSettingsSections.kt |
| `@layout/error_panel` | 业务组件 | empty | 14 | include |
| `@layout/info_cell` | 业务组件 | cell | 7 | include |
| `@layout/list_empty_view` | 业务组件 | empty | 6 | include |
| `@layout/toolbar_layout` | 业务组件 | navbar | 6 | include |
| `@layout/settings_notification_action` | 业务组件 |  | 5 | include |
| … | | | | 共 65 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| navbar | TopAppBar（组件库，用 3 次）<br>TopAppBar（Composable，用 10 次）<br>@layout/toolbar_layout（业务组件，用 6 次）<br>Toolbar（组件库，用 5 次）<br>@layout/settings_category_header_title（业务组件，用 1 次） | `TopAppBar` |
| list | LazyColumn（组件库，用 4 次）<br>RecyclerView（组件库，用 26 次）<br>NewPipeRecyclerView（自定义 View，用 2 次） | `LazyColumn` |
| cell | ListItem（组件库，用 4 次）<br>PreferenceRow（Composable，用 12 次）<br>@layout/info_cell（业务组件，用 7 次） | `ListItem` |
| empty | @layout/error_panel（业务组件，用 14 次）<br>@layout/list_empty_view（业务组件，用 6 次） | `@layout/error_panel` |
| card | Card（Compose Material3，项目未使用）<br>ExoPlayerSettingsSections（Composable，用 1 次） | `Card` |

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
| #FFF8F7 | 3 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #534342 | 3 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #231918 | 3 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #1A1110 | 3 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #00000000 | 2 | `@color/transparent_background_color` | app/src/main/res/layout/activity_player_queue_control.xml、app/src/main/res/layout-land/activity_player_queue_control.xml |
| #60000000 | 2 | — | app/src/main/res/layout/player.xml |
| #C0000000 | 2 | — | app/src/main/res/layout/activity_player_queue_control.xml、app/src/main/res/layout-land/activity_player_queue_control.xml |
| #392E2D | 2 | — | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
| #000000 | 2 | `@color/black_background_color` | shared/src/commonMain/kotlin/net/newpipe/app/theme/Color.kt |
