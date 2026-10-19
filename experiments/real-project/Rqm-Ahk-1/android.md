
## 页面复用检查：qmuidemo/src/main/res/values/colors.xml,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDOrderDetailFragment.kt,qmuidemo/src/main/res/layout/fragment_order_detail.xml  ❌ FAIL

- 复用的现成组件（6）：ImageView（组件库）、QMUIRoundButton（自定义 View）、QMUIRoundLinearLayout（自定义 View）、QMUITopBarLayout（自定义 View）、QMUIWindowInsetLayout（自定义 View）、TextView（组件库）
- 自造轮子（0）
- 有组件却没用（0）
- 疑似手写已有组件（1）
  - ✗ qmuidemo/src/main/res/layout/fragment_order_detail.xml:293 页面里有 7 行「标题 + 值」横排是手写的（第 293、317、341、380、… 行），已有 QMUICommonListItemView（自定义 View qmui/src/main/java/com/qmuiteam/qmui/widget/grouplist/QMUICommonListItemView.java）（项目里已用 48 次）
- 页面私有调色板：页面里自己定义了 11 个颜色变量（应使用项目 / 组件库已有的 token）
- 写死颜色（1）
  - · qmuidemo/src/main/res/values/colors.xml:62 #FFFFFF：新定义的颜色 qd_order_detail_white 与已有 token @color/qmui_config_color_white 同值

---

# 存量复用体检：/home/user/exp/Rqm-Ahk-1（android）

技术栈：XML View / Jetpack Compose / Material Components / QMUI Android / Compose Foundation · 扫描文件 554

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（935 / 935） |
| 重复概念（同一 UI 概念有多个实现） | 11 |
| 重复结构（可抽成组件的复制粘贴结构） | 0 组 |
| 硬编码颜色 | 11 处（其中 1 处已有同值 token） |
| 硬编码尺寸 | 477 处 |

## 0. 优先处理（按影响面排序）

1. 「cell（View）」有 2 种实现（共用 49 次）→ 统一到 `QMUICommonListItemView`
2. 「button（Compose）」有 5 种实现（共用 13 次）→ 统一到 `CommonImageButton`
3. 「loading（View）」有 4 种实现（共用 13 次）→ 统一到 `QMUIProgressBar`
4. 「navbar（View）」有 3 种实现（共用 9 次）→ 统一到 `QMUITopBar`
5. 「dialog（View）」有 3 种实现（共用 8 次）→ 统一到 `QMUIDialogView`

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Android SDK widgets | 173 | TextView×135、ImageView×12、RecyclerView×6、Button×5、RadioButton×5、SeekBar×3、ListView×2、EditText×2 |
| Compose Foundation | 11 | LazyColumn×9、LazyRow×1、HorizontalPager×1 |
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

**项目惯例写法**（没有组件、但项目里固定这么写；新设计稿里没有对应组件时沿用它，不要另造）：

- divider：1px View 细线，4 处，例 `<View android:layout_width="match_parent" android:layout_height="1dp" android:layout_marginTop="12dp" android:background="@color/qd_order_de…` qmuidemo/src/main/res/layout/fragment_order_detail.xml:131

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `TextView` | 组件库 |  | 135 | `<TextView` qmui/src/main/res/layout/qmui_common_list_item.xml:45 |
| `ImageView` | 组件库 | image | 12 | `<ImageView` qmuidemo/src/main/res/layout/drawablehelper_createfromview.xml:25 |
| `LazyColumn` | 组件库 | list | 9 | `LazyColumn(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIBottomSheet.kt:40 |
| `RecyclerView` | 组件库 | list | 6 | `<androidx.recyclerview.widget.RecyclerView` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:65 |
| `RadioButton` | 组件库 | radio | 5 | `<RadioButton` qmuidemo/src/main/res/layout/fragment_layout.xml:189 |
| `Button` | 组件库 | button | 5 | `<Button` qmuidemo/src/main/res/layout/fragment_layout.xml:114 |
| `QMUITabSegment` | 组件库 | tabs | 4 | `<com.qmuiteam.qmui.widget.tab.QMUITabSegment` qmuidemo/src/main/res/layout/fragment_fsw_viewpager.xml:30 |
| `SeekBar` | 组件库 |  | 3 | `<SeekBar` qmuidemo/src/main/res/layout/fragment_colorhelper.xml:124 |
| `ListView` | 组件库 | list | 2 | `<ListView` qmuidemo/src/main/res/layout/fragment_listview.xml:25 |
| `EditText` | 组件库 | input | 2 | `<EditText` qmuidemo/src/main/res/layout/fragment_scheme.xml:47 |
| `ViewPager` | 组件库 | carousel | 1 | `<androidx.viewpager.widget.ViewPager` qmuidemo/src/main/res/layout/fragment_tab_viewpager_layout.xml:39 |
| `HorizontalPager` | 组件库 | carousel | 1 | `HorizontalPager(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:39 |
| `LazyRow` | 组件库 | list | 1 | `LazyRow(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:116 |
| `ViewPager2` | 组件库 |  | 1 | `<androidx.viewpager2.widget.ViewPager2` qmuidemo/src/main/res/layout/fragment_tab_viewpager2_layout.xml:39 |
| `QMUISeekBar` | 组件库 |  | 1 | `<com.qmuiteam.qmui.widget.QMUISeekBar` qmuidemo/src/main/res/layout/fragment_slider.xml:44 |
| `ProgressBar` | 组件库 | loading | 1 | `<ProgressBar` qmuidemo/src/main/res/layout/fragment_webview_explorer.xml:37 |
| `QMUIItem` | Composable | cell | 10 | `QMUIItem(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:152 |
| `QMUIPhotoThumbnailWithViewer` | Composable |  | 9 | `QMUIPhotoThumbnailWithViewer(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDPhotoFragment.kt:102 |
| `QMUIDialogAction` | Composable | dialog | 8 | `QMUIDialogAction(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:90 |
| `QMUITopBar` | Composable | navbar | 7 | `QMUITopBar(` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:193 |
| `CommonImageButton` | Composable | button | 6 | `CommonImageButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Edit.kt:338 |
| `QMUIDialog` | Composable | dialog | 5 | `QMUIDialog(modal, horEdge, verEdge, widthLimit, radius, background, content)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:351 |
| `TextButton` | Composable | button | 4 | `TextButton("加粗"){` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEditorFragment.kt:57 |
| `OnePx` | Composable |  | 4 | `separatorHeight: Dp = OnePx(),` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:179 |
| `QMUITopBarWithLazyScrollState` | Composable |  | 4 | `QMUITopBarWithLazyScrollState(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:358 |
| `QMUIGesturePhoto` | Composable |  | 4 | `QMUIGesturePhoto(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoViewerActivity.kt:185 |
| `PressWithAlphaBox` | Composable |  | 4 | `PressWithAlphaBox(` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:111 |
| `BlankBox` | Composable |  | 3 | `BlankBox()` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:80 |
| `QMUIDialogList` | Composable | dialog | 3 | `QMUIDialogList(modal, maxHeight, state, contentPadding) {` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:150 |
| `QMUIPhotoLoading` | Composable | loading | 3 | `QMUIPhotoLoading(lineColor = QMUILocalPickerConfig.current.loadingColor)` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:577 |
| `LongImage` | Composable | image | 2 | `LongImage(onSuccess, onError, openBlankColor)` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:50 |
| `QMUIBitmapRegionItem` | Composable | cell | 2 | `QMUIBitmapRegionItem(image, maxWidth, heightDp)` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:236 |
| `CommonButton` | Composable | button | 2 | `CommonButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Edit.kt:815 |
| `OriginOpenButton` | Composable | button | 2 | `OriginOpenButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:207 |
| `QMUIPhotoClipper` | Composable |  | 2 | `QMUIPhotoClipper(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoClipActivity.kt:92 |
| `QMUIDialogActions` | Composable | dialog | 2 | `QMUIDialogActions(modal, actions)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:111 |
| `RowImages` | Composable | image | 2 | `RowImages(images, renderInfo, config, maxWidth, 2, 0, onClick)` photo/src/main/java/com/qmuiteam/photo/compose/Thumbnail.kt:255 |
| `QMUIToast` | Composable | toast | 2 | `QMUIToast(modal, radius, background, content)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIToast.kt:116 |
| `CommonTextButton` | Composable | button | 2 | `CommonTextButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:199 |
| `QMUIPhotoPickerGridToolBar` | Composable | navbar | 1 | `QMUIPhotoPickerGridToolBar(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:391 |
| `QMUIDraggableScrollBar` | 自定义 View |  | 2 | `QMUIContinuousNestedTopAreaBehavior.Callback, QMUIDraggableScrollBar.Callback {` qmui/src/main/java/com/qmuiteam/qmui/nestedScroll/QMUIContinuousNestedScrollLayout.java:38 |
| `QMUITipDialog` | 自定义 View | toast | 14 | `@Widget(widgetClass = QMUITipDialog.class, iconRes = R.mipmap.icon_grid_tip_dialog)` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDTipDialogFragment.java:44 |
| `QMUIPhotoPickRadio` | Composable | radio | 1 | `QMUIPhotoPickRadio(isOriginOpen)` photo/src/main/java/com/qmuiteam/photo/compose/picker/Common.kt:142 |
| `QMUIDialogView` | 自定义 View | dialog | 3 | `private QMUIDialogView.OnDecorationListener mOnDecorationListener;` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialogBuilder.java:94 |
| `QMUIRoundLinearLayout` | 自定义 View |  | 6 | `<com.qmuiteam.qmui.widget.roundwidget.QMUIRoundLinearLayout` qmuidemo/src/main/res/layout/fragment_order_detail.xml:45 |
| `QMUIBottomSheet` | 自定义 View | sheet | 38 | `mDialog = new QMUIBottomSheet(mContext, style);` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBottomSheetBaseBuilder.java:116 |
| `QMUISlider` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.widget.QMUISlider` qmuidemo/src/main/res/layout/fragment_slider.xml:35 |
| `QMUICollapsingTopBarLayout` | 自定义 View |  | 3 | `<com.qmuiteam.qmui.widget.QMUICollapsingTopBarLayout` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:31 |
| `QMUIGroupListSectionHeaderFooterView` | 自定义 View | navbar | 2 | `return new QMUIGroupListSectionHeaderFooterView(mContext, titleText);` qmui/src/main/java/com/qmuiteam/qmui/widget/grouplist/QMUIGroupListView.java:359 |
| `EmojiconTextView` | 自定义 View |  | 2 | `EmojiconTextView(1);` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/qqface/QDQQFacePerformanceTestFragment.java:153 |
| `HomeController` | 自定义 View |  | 2 | `HomeController.HomeControlListener listener = new HomeController.HomeControlListener() {` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/home/HomeFragment.java:181 |
| `QMUIBottomSheetList` | Composable | sheet | 1 | `QMUIBottomSheetList(it) {` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/QDDialogFragment.kt:217 |
| `QMUIDialogTitle` | Composable | dialog | 1 | `QMUIDialogTitle(title)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:109 |
| `QMUIMarkIcon` | Composable | icon | 1 | `QMUIMarkIcon(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Buckets.kt:161 |
| `QMUIWindowInsetsProvider` | Composable |  | 0 | compose-core/src/main/java/com/qmuiteam/compose/core/provider/WindowInsets.kt |
| `QMUIContinuousNestedTopWebView` | 自定义 View |  | 4 | `mNestedWebView = new QMUIContinuousNestedTopWebView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll1Fragment.java:65 |
| `QMUILinkTextView` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.widget.textview.QMUILinkTextView` qmuidemo/src/main/res/layout/fragment_link_texview_layout.xml:49 |
| `QMUIFloatLayout` | 自定义 View | popup | 4 | `<com.qmuiteam.qmui.widget.QMUIFloatLayout` qmuidemo/src/main/res/layout/fragment_floatlayout.xml:37 |
| `QMUITabView` | 自定义 View | tabs | 2 | `public class QMUITabAdapter extends QMUIItemViewsAdapter<QMUITab, QMUITabView> implements QMUITabView.Callback {` qmui/src/main/java/com/qmuiteam/qmui/widget/tab/QMUITabAdapter.java:23 |
| `QMUICheckBox` | Composable | checkbox | 1 | `QMUICheckBox(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Common.kt:170 |
| … | | | | 共 158 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| button（Compose） | CommonImageButton（Composable，用 6 次）<br>CommonButton（Composable，用 2 次）<br>CommonTextButton（Composable，用 2 次）<br>OriginOpenButton（Composable，用 2 次）<br>PickCurrentCheckButton（Composable，用 1 次） | `CommonImageButton` |
| loading（View） | QMUIProgressBar（自定义 View，用 6 次）<br>QMUILoadingView（自定义 View，用 5 次）<br>QDLoadingItemView（自定义 View，用 1 次）<br>ProgressBar（组件库，用 1 次） | `QMUIProgressBar` |
| navbar（View） | QMUITopBar（自定义 View，用 6 次）<br>QDSectionHeaderView（自定义 View，用 1 次）<br>QMUIGroupListSectionHeaderFooterView（自定义 View，用 2 次） | `QMUITopBar` |
| dialog（View） | QMUIDialogView（自定义 View，用 3 次）<br>QMUIDialogMenuItemView（自定义 View，用 5 次）<br>QMUIBaseDialog（自定义 View，用 0 次） | `QMUIDialogView` |
| cell（View） | QMUICommonListItemView（自定义 View，用 48 次）<br>QMUIBottomSheetListItemView（自定义 View，用 1 次） | `QMUICommonListItemView` |
| dialog（Compose） | QMUIDialogAction（Composable，用 8 次）<br>QMUIModal（Composable，用 0 次） | `QMUIDialogAction` |
| button（View） | Button（组件库，用 5 次）<br>QMUIAlphaImageButton（自定义 View，用 2 次） | `Button` |
| tabs（View） | QMUIBasicTabSegment（自定义 View，用 4 次）<br>QMUITabView（自定义 View，用 2 次） | `QMUIBasicTabSegment` |
| carousel（View） | SwipeBackgroundView（自定义 View，用 2 次）<br>ViewPager（组件库，用 1 次） | `SwipeBackgroundView` |
| checkbox（Compose） | QMUICheckBox（Composable，用 1 次）<br>QMUIPhotoPickCheckBox（Composable，用 1 次） | `QMUICheckBox` |
| icon（Compose） | QMUIMarkIcon（Composable，用 1 次）<br>QMUIChevronIcon（Composable，用 0 次） | `QMUIMarkIcon` |

## 6. 硬编码颜色热点

| 颜色 | 次数 | 已有 token | 出现在 |
|---|---|---|---|
| #333333 | 2 | — | photo/src/main/java/com/qmuiteam/photo/compose/picker/Config.kt |
| #FF0000 | 2 | — | qmuidemo/src/main/res/layout/fragment_layout.xml、qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDLayoutFragment.java |
| #0000FF | 2 | — | qmuidemo/src/main/res/layout/fragment_layout.xml、qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDLayoutFragment.java |
| #666666 | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #222222 | 1 | — | photo/src/main/java/com/qmuiteam/photo/compose/picker/Config.kt |
| #00A8E1 | 1 | `@color/app_color_blue` | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #CCCCCC | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #777777 | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
