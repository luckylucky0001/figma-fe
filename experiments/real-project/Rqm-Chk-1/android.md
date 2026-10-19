
## 页面复用检查：qmuidemo/src/main/res/values/colors.xml,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDOrderDetailFragment.kt,qmuidemo/src/main/res/layout/fragment_order_detail.xml  ✅ PASS

- 复用的现成组件（8）：QMUICommonListItemView（自定义 View）、QMUIFrameLayout（自定义 View）、QMUIGroupListView（自定义 View）、QMUILinearLayout（自定义 View）、QMUIRoundButton（自定义 View）、QMUITopBarLayout（自定义 View）、QMUIWindowInsetLayout（自定义 View）、TextView（组件库）
- 自造轮子（0）
- 有组件却没用（0）
- 疑似手写已有组件（0）
- 页面私有调色板：页面里自己定义了 10 个颜色变量（应使用项目 / 组件库已有的 token）
- 写死颜色（0）

---

# 存量复用体检：/home/user/exp/Rqm-Chk-1（android）

技术栈：XML View / Jetpack Compose / Material Components / QMUI Android / Compose Foundation · 扫描文件 554

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（912 / 912） |
| 重复概念（同一 UI 概念有多个实现） | 11 |
| 重复结构（可抽成组件的复制粘贴结构） | 0 组 |
| 硬编码颜色 | 11 处（其中 1 处已有同值 token） |
| 硬编码尺寸 | 381 处 |

## 0. 优先处理（按影响面排序）

1. 「cell（View）」有 2 种实现（共用 57 次）→ 统一到 `QMUICommonListItemView`
2. 「button（Compose）」有 5 种实现（共用 13 次）→ 统一到 `CommonImageButton`
3. 「loading（View）」有 4 种实现（共用 13 次）→ 统一到 `QMUIProgressBar`
4. 「navbar（View）」有 3 种实现（共用 9 次）→ 统一到 `QMUITopBar`
5. 「dialog（View）」有 3 种实现（共用 8 次）→ 统一到 `QMUIDialogView`

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Android SDK widgets | 142 | TextView×105、ImageView×11、RecyclerView×6、RadioButton×5、Button×5、SeekBar×3、EditText×2、ListView×2 |
| Compose Foundation | 11 | LazyColumn×9、HorizontalPager×1、LazyRow×1 |
| QMUI Android | 5 | QMUITabSegment×4、QMUISeekBar×1 |
| Material Components | 0 |  |

**组件库里有、但本项目不能直接用的组件**（match 不会推荐它们）：

- `BottomNavigationView`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `ChipGroup`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `Chip`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `CircularProgressIndicator`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `ExtendedFloatingActionButton`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `LinearProgressIndicator`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `MaterialButtonToggleGroup`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `MaterialButton`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `MaterialCardView`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `MaterialDivider`：需要 Material Components ≥ 1.5.0，项目是 1.4.0，用了会编译失败
- `MaterialSwitch`：需要 Material Components ≥ 1.7.0，项目是 1.4.0，用了会编译失败
- `SearchBar`：需要 Material Components ≥ 1.8.0，项目是 1.4.0，用了会编译失败
- `SearchView`：需要 Material Components ≥ 1.8.0，项目是 1.4.0，用了会编译失败
- `Slider`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `TextInputLayout`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `TextView` | 组件库 |  | 105 | `<TextView` qmui/src/main/res/layout/qmui_common_list_item.xml:45 |
| `ImageView` | 组件库 | image | 11 | `<ImageView` qmuidemo/src/main/res/layout/drawablehelper_createfromview.xml:25 |
| `LazyColumn` | 组件库 | list | 9 | `LazyColumn(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIBottomSheet.kt:40 |
| `RecyclerView` | 组件库 | list | 6 | `<androidx.recyclerview.widget.RecyclerView` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:65 |
| `RadioButton` | 组件库 | radio | 5 | `<RadioButton` qmuidemo/src/main/res/layout/fragment_layout.xml:189 |
| `Button` | 组件库 | button | 5 | `<Button` qmuidemo/src/main/res/layout/fragment_layout.xml:114 |
| `QMUITabSegment` | 组件库 | tabs | 4 | `<com.qmuiteam.qmui.widget.tab.QMUITabSegment` qmuidemo/src/main/res/layout/fragment_fsw_viewpager.xml:30 |
| `SeekBar` | 组件库 |  | 3 | `<SeekBar` qmuidemo/src/main/res/layout/fragment_colorhelper.xml:124 |
| `ListView` | 组件库 | list | 2 | `<ListView` qmuidemo/src/main/res/layout/fragment_listview.xml:25 |
| `EditText` | 组件库 | input | 2 | `<EditText` qmuidemo/src/main/res/layout/fragment_scheme.xml:47 |
| `LazyRow` | 组件库 | list | 1 | `LazyRow(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:116 |
| `QMUISeekBar` | 组件库 |  | 1 | `<com.qmuiteam.qmui.widget.QMUISeekBar` qmuidemo/src/main/res/layout/fragment_slider.xml:44 |
| `HorizontalPager` | 组件库 | carousel | 1 | `HorizontalPager(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:39 |
| `ViewPager2` | 组件库 |  | 1 | `<androidx.viewpager2.widget.ViewPager2` qmuidemo/src/main/res/layout/fragment_tab_viewpager2_layout.xml:39 |
| `ViewPager` | 组件库 | carousel | 1 | `<androidx.viewpager.widget.ViewPager` qmuidemo/src/main/res/layout/fragment_tab_viewpager_layout.xml:39 |
| `ProgressBar` | 组件库 | loading | 1 | `<ProgressBar` qmuidemo/src/main/res/layout/fragment_webview_explorer.xml:37 |
| `QMUIItem` | Composable | cell | 10 | `QMUIItem(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:152 |
| `QMUIPhotoThumbnailWithViewer` | Composable |  | 9 | `QMUIPhotoThumbnailWithViewer(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDPhotoFragment.kt:102 |
| `QMUIDialogAction` | Composable | dialog | 8 | `QMUIDialogAction(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:90 |
| `QMUITopBar` | Composable | navbar | 7 | `QMUITopBar(` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:193 |
| `CommonImageButton` | Composable | button | 6 | `CommonImageButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Edit.kt:338 |
| `QMUIDialog` | Composable | dialog | 5 | `QMUIDialog(modal, horEdge, verEdge, widthLimit, radius, background, content)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:351 |
| `OnePx` | Composable |  | 4 | `separatorHeight: Dp = OnePx(),` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:179 |
| `TextButton` | Composable | button | 4 | `TextButton("加粗"){` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEditorFragment.kt:57 |
| `QMUIGesturePhoto` | Composable |  | 4 | `QMUIGesturePhoto(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoViewerActivity.kt:185 |
| `PressWithAlphaBox` | Composable |  | 4 | `PressWithAlphaBox(` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:111 |
| `QMUITopBarWithLazyScrollState` | Composable |  | 4 | `QMUITopBarWithLazyScrollState(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:358 |
| `QMUIPhotoLoading` | Composable | loading | 3 | `QMUIPhotoLoading(lineColor = QMUILocalPickerConfig.current.loadingColor)` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:577 |
| `BlankBox` | Composable |  | 3 | `BlankBox()` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:80 |
| `QMUIDialogList` | Composable | dialog | 3 | `QMUIDialogList(modal, maxHeight, state, contentPadding) {` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:150 |
| `LongImage` | Composable | image | 2 | `LongImage(onSuccess, onError, openBlankColor)` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:50 |
| `QMUICollapsingTopBarLayout` | 自定义 View |  | 3 | `<com.qmuiteam.qmui.widget.QMUICollapsingTopBarLayout` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:31 |
| `QMUISpanTouchFixTextView` | 自定义 View |  | 15 | `<com.qmuiteam.qmui.widget.textview.QMUISpanTouchFixTextView` qmui/src/main/res/layout/qmui_empty_view.xml:36 |
| `QMUIToast` | Composable | toast | 2 | `QMUIToast(modal, radius, background, content)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIToast.kt:116 |
| `QMUIDialogMutiCheckList` | Composable | checkbox | 1 | `QMUIDialogMutiCheckList(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/QDDialogFragment.kt:169 |
| `QMUIPhotoPickCheckBox` | Composable | checkbox | 1 | `QMUIPhotoPickCheckBox(pickedIndex)` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:163 |
| `QMUIPhotoPickerBucketToggleArrow` | Composable | switch | 1 | `QMUIPhotoPickerBucketToggleArrow(iconBgColor, iconColor, isFocusFlow)` photo/src/main/java/com/qmuiteam/photo/compose/picker/TopBarItem.kt:68 |
| `QMUIBottomSheetGridLineLayout` | 自定义 View |  | 3 | `private QMUIBottomSheetGridLineLayout.ItemWidthCalculator mItemWidthCalculator = null;` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBottomSheet.java:521 |
| `QMUIDialogTitle` | Composable | dialog | 1 | `QMUIDialogTitle(title)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:109 |
| `QMUIQQFaceView` | 自定义 View |  | 25 | `<com.qmuiteam.qmui.qqface.QMUIQQFaceView` qmuidemo/src/main/res/layout/fragment_qqface_layout.xml:97 |
| `RowImages` | Composable | image | 2 | `RowImages(images, renderInfo, config, maxWidth, 2, 0, onClick)` photo/src/main/java/com/qmuiteam/photo/compose/Thumbnail.kt:255 |
| `QMUIEditor` | Composable |  | 1 | `QMUIEditor(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEditorFragment.kt:42 |
| `PickCurrentCheckButton` | Composable | button | 1 | `PickCurrentCheckButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:210 |
| `QMUIContinuousNestedTopDelegateLayout` | 自定义 View |  | 4 | `mTopDelegateLayout = new QMUIContinuousNestedTopDelegateLayout(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll4Fragment.java:56 |
| `QMUIPullRefreshLayout` | 自定义 View |  | 11 | `<com.qmuiteam.qmui.widget.pullRefreshLayout.QMUIPullRefreshLayout` qmuidemo/src/main/res/layout/fragment_continuous_nested_scroll.xml:27 |
| `QMUIWindowInsetLayout` | 自定义 View |  | 26 | `<com.qmuiteam.qmui.widget.QMUIWindowInsetLayout xmlns:android="http://schemas.android.com/apk/res/android"` qmuidemo/src/main/res/layout/activity_arch_test.xml:18 |
| `QMUIPhotoThumbnail` | Composable |  | 1 | `QMUIPhotoThumbnail(images, config) { list, index ->` photo/src/main/java/com/qmuiteam/photo/compose/Thumbnail.kt:93 |
| `PhotoContentWithRectTransition` | Composable |  | 1 | `PhotoContentWithRectTransition(` photo/src/main/java/com/qmuiteam/photo/compose/GesturePhoto.kt:396 |
| `QMUIMarkIcon` | Composable | icon | 1 | `QMUIMarkIcon(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Buckets.kt:161 |
| `QMUIDialogView` | 自定义 View | dialog | 3 | `private QMUIDialogView.OnDecorationListener mOnDecorationListener;` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialogBuilder.java:94 |
| `QMUILinkTextView` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.widget.textview.QMUILinkTextView` qmuidemo/src/main/res/layout/fragment_link_texview_layout.xml:49 |
| `QMUIDialogMenuItemView` | 自定义 View | dialog | 5 | `itemView.setListener(new QMUIDialogMenuItemView.MenuItemViewListener() {` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialog.java:494 |
| `QMUIPhotoPickRadio` | Composable | radio | 1 | `QMUIPhotoPickRadio(isOriginOpen)` photo/src/main/java/com/qmuiteam/photo/compose/picker/Common.kt:142 |
| `PhotoContentWithAlphaTransition` | Composable |  | 1 | `PhotoContentWithAlphaTransition(` photo/src/main/java/com/qmuiteam/photo/compose/GesturePhoto.kt:368 |
| `QMUITopBarLayout` | 自定义 View | navbar | 46 | `<com.qmuiteam.qmui.widget.QMUITopBarLayout` qmuidemo/src/main/res/layout/activity_arch_test.xml:23 |
| `QMUIBasicTabSegment` | 自定义 View | tabs | 4 | `* Please use QMUIBasicTabSegment.OnTabClickListener for a replacement` qmui/src/main/java/com/qmuiteam/qmui/widget/tab/QMUITabSegment.java:307 |
| `QMUIPhotoPickerGridToolBar` | Composable | navbar | 1 | `QMUIPhotoPickerGridToolBar(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:391 |
| `QMUISlider` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.widget.QMUISlider` qmuidemo/src/main/res/layout/fragment_slider.xml:35 |
| `QMUIPriorityLinearLayout` | 自定义 View |  | 18 | `<com.qmuiteam.qmui.layout.QMUIPriorityLinearLayout` qmuidemo/src/main/res/layout/fragment_priority_linear_layout.xml:45 |
| `QMUIAnimationListView` | 自定义 View | list | 4 | `<com.qmuiteam.qmui.widget.QMUIAnimationListView` qmuidemo/src/main/res/layout/fragment_animation_listview.xml:32 |
| … | | | | 共 158 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| button（Compose） | CommonImageButton（Composable，用 6 次）<br>CommonButton（Composable，用 2 次）<br>CommonTextButton（Composable，用 2 次）<br>OriginOpenButton（Composable，用 2 次）<br>PickCurrentCheckButton（Composable，用 1 次） | `CommonImageButton` |
| loading（View） | QMUIProgressBar（自定义 View，用 6 次）<br>QMUILoadingView（自定义 View，用 5 次）<br>QDLoadingItemView（自定义 View，用 1 次）<br>ProgressBar（组件库，用 1 次） | `QMUIProgressBar` |
| navbar（View） | QMUITopBar（自定义 View，用 6 次）<br>QDSectionHeaderView（自定义 View，用 1 次）<br>QMUIGroupListSectionHeaderFooterView（自定义 View，用 2 次） | `QMUITopBar` |
| dialog（View） | QMUIDialogView（自定义 View，用 3 次）<br>QMUIDialogMenuItemView（自定义 View，用 5 次）<br>QMUIBaseDialog（自定义 View，用 0 次） | `QMUIDialogView` |
| cell（View） | QMUICommonListItemView（自定义 View，用 56 次）<br>QMUIBottomSheetListItemView（自定义 View，用 1 次） | `QMUICommonListItemView` |
| dialog（Compose） | QMUIDialogAction（Composable，用 8 次）<br>QMUIModal（Composable，用 0 次） | `QMUIDialogAction` |
| button（View） | Button（组件库，用 5 次）<br>QMUIAlphaImageButton（自定义 View，用 2 次） | `Button` |
| tabs（View） | QMUIBasicTabSegment（自定义 View，用 4 次）<br>QMUITabView（自定义 View，用 2 次） | `QMUIBasicTabSegment` |
| carousel（View） | SwipeBackgroundView（自定义 View，用 2 次）<br>ViewPager（组件库，用 1 次） | `SwipeBackgroundView` |
| checkbox（Compose） | QMUIPhotoPickCheckBox（Composable，用 1 次）<br>QMUICheckBox（Composable，用 1 次） | `QMUIPhotoPickCheckBox` |
| icon（Compose） | QMUIMarkIcon（Composable，用 1 次）<br>QMUIChevronIcon（Composable，用 0 次） | `QMUIMarkIcon` |

## 6. 硬编码颜色热点

| 颜色 | 次数 | 已有 token | 出现在 |
|---|---|---|---|
| #333333 | 2 | — | photo/src/main/java/com/qmuiteam/photo/compose/picker/Config.kt |
| #FF0000 | 2 | — | qmuidemo/src/main/res/layout/fragment_layout.xml、qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDLayoutFragment.java |
| #0000FF | 2 | — | qmuidemo/src/main/res/layout/fragment_layout.xml、qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDLayoutFragment.java |
| #00A8E1 | 1 | `@color/app_color_blue` | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #CCCCCC | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #777777 | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #666666 | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #222222 | 1 | — | photo/src/main/java/com/qmuiteam/photo/compose/picker/Config.kt |
