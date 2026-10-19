
## 页面复用检查：qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDOrderDetailFragment.kt,qmuidemo/src/main/res/layout/fragment_order_detail.xml  ❌ FAIL

- 复用的现成组件（8）：ImageView（组件库）、QMUIRadiusImageView2（自定义 View）、QMUIRoundButton（自定义 View）、QMUIRoundFrameLayout（自定义 View）、QMUIRoundLinearLayout（自定义 View）、QMUITopBarLayout（自定义 View）、QMUIWindowInsetLayout（自定义 View）、TextView（组件库）
- 自造轮子（0）
- 有组件却没用（0）
- 疑似手写已有组件（1）
  - ✗ qmuidemo/src/main/res/layout/fragment_order_detail.xml:148 页面手写了「cell」（.order_address_cell），已有 QMUICommonListItemView（自定义 View qmui/src/main/java/com/qmuiteam/qmui/widget/grouplist/QMUICommonListItemView.java）（项目里已用 48 次）
- 写死颜色（57）
  - · qmuidemo/src/main/res/layout/fragment_order_detail.xml:21 #F5F6F7：写死颜色
  - · qmuidemo/src/main/res/layout/fragment_order_detail.xml:21 #FFFFFF：写死颜色，已有同值 token @color/qmui_config_color_white
  - · qmuidemo/src/main/res/layout/fragment_order_detail.xml:21 #E5E6EB：写死颜色
  - · qmuidemo/src/main/res/layout/fragment_order_detail.xml:22 #1F2329：写死颜色
  - · qmuidemo/src/main/res/layout/fragment_order_detail.xml:22 #646A73：写死颜色
  - · qmuidemo/src/main/res/layout/fragment_order_detail.xml:22 #8F959E：写死颜色
  - · qmuidemo/src/main/res/layout/fragment_order_detail.xml:23 #FF4D1A：写死颜色
  - · qmuidemo/src/main/res/layout/fragment_order_detail.xml:23 #FF8F1F：写死颜色
  - · qmuidemo/src/main/res/layout/fragment_order_detail.xml:23 #FFF3E8：写死颜色
  - · qmuidemo/src/main/res/layout/fragment_order_detail.xml:23 #F53F3F：写死颜色
  - … 共 57 处

---

# 存量复用体检：/home/user/exp/Rqm-A-1（android）

技术栈：XML View / Jetpack Compose / Material Components / QMUI Android / Compose Foundation · 扫描文件 554

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（930 / 930） |
| 重复概念（同一 UI 概念有多个实现） | 11 |
| 重复结构（可抽成组件的复制粘贴结构） | 0 组 |
| 硬编码颜色 | 47 处（其中 3 处已有同值 token） |
| 硬编码尺寸 | 462 处 |

## 0. 优先处理（按影响面排序）

1. 「cell（View）」有 2 种实现（共用 49 次）→ 统一到 `QMUICommonListItemView`
2. 「button（Compose）」有 5 种实现（共用 13 次）→ 统一到 `CommonImageButton`
3. 「loading（View）」有 4 种实现（共用 13 次）→ 统一到 `QMUIProgressBar`
4. 「navbar（View）」有 3 种实现（共用 9 次）→ 统一到 `QMUITopBar`
5. 「dialog（View）」有 3 种实现（共用 8 次）→ 统一到 `QMUIDialogView`

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Android SDK widgets | 168 | TextView×130、ImageView×12、RecyclerView×6、Button×5、RadioButton×5、SeekBar×3、ListView×2、EditText×2 |
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

- divider：1px View 细线，4 处，例 `<View android:layout_width="match_parent" android:layout_height="1px" android:background="?attr/qmui_skin_support_color_separator_darken"/>` qmuidemo/src/main/res/layout/fragment_touch_span_fix_layout.xml:81

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `TextView` | 组件库 |  | 130 | `<TextView` qmui/src/main/res/layout/qmui_common_list_item.xml:45 |
| `ImageView` | 组件库 | image | 12 | `<ImageView` qmuidemo/src/main/res/layout/drawablehelper_createfromview.xml:25 |
| `LazyColumn` | 组件库 | list | 9 | `LazyColumn(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIBottomSheet.kt:40 |
| `RecyclerView` | 组件库 | list | 6 | `<androidx.recyclerview.widget.RecyclerView` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:65 |
| `Button` | 组件库 | button | 5 | `<Button` qmuidemo/src/main/res/layout/fragment_layout.xml:114 |
| `RadioButton` | 组件库 | radio | 5 | `<RadioButton` qmuidemo/src/main/res/layout/fragment_layout.xml:189 |
| `QMUITabSegment` | 组件库 | tabs | 4 | `<com.qmuiteam.qmui.widget.tab.QMUITabSegment` qmuidemo/src/main/res/layout/fragment_fsw_viewpager.xml:30 |
| `SeekBar` | 组件库 |  | 3 | `<SeekBar` qmuidemo/src/main/res/layout/fragment_colorhelper.xml:124 |
| `ListView` | 组件库 | list | 2 | `<ListView` qmuidemo/src/main/res/layout/fragment_listview.xml:25 |
| `EditText` | 组件库 | input | 2 | `<EditText` qmuidemo/src/main/res/layout/fragment_scheme.xml:47 |
| `HorizontalPager` | 组件库 | carousel | 1 | `HorizontalPager(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:39 |
| `ProgressBar` | 组件库 | loading | 1 | `<ProgressBar` qmuidemo/src/main/res/layout/fragment_webview_explorer.xml:37 |
| `QMUISeekBar` | 组件库 |  | 1 | `<com.qmuiteam.qmui.widget.QMUISeekBar` qmuidemo/src/main/res/layout/fragment_slider.xml:44 |
| `ViewPager` | 组件库 | carousel | 1 | `<androidx.viewpager.widget.ViewPager` qmuidemo/src/main/res/layout/fragment_tab_viewpager_layout.xml:39 |
| `LazyRow` | 组件库 | list | 1 | `LazyRow(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:116 |
| `ViewPager2` | 组件库 |  | 1 | `<androidx.viewpager2.widget.ViewPager2` qmuidemo/src/main/res/layout/fragment_tab_viewpager2_layout.xml:39 |
| `SwipeBackLayout` | 自定义 View |  | 78 | `import static com.qmuiteam.qmui.arch.SwipeBackLayout.DRAG_DIRECTION_BOTTOM_TO_TOP;` arch/src/main/java/com/qmuiteam/qmui/arch/QMUIActivity.java:37 |
| `QMUIPullLayout` | 自定义 View |  | 49 | `<com.qmuiteam.qmui.widget.pullLayout.QMUIPullLayout` qmuidemo/src/main/res/layout/fragment_pull_horizontal_test_layout.xml:32 |
| `QMUICommonListItemView` | 自定义 View | cell | 48 | `QMUICommonListItemView itemView = new QMUICommonListItemView(getContext());` qmui/src/main/java/com/qmuiteam/qmui/widget/grouplist/QMUIGroupListView.java:83 |
| `QMUITopBarLayout` | 自定义 View | navbar | 46 | `<com.qmuiteam.qmui.widget.QMUITopBarLayout` qmuidemo/src/main/res/layout/activity_arch_test.xml:23 |
| `QMUIBottomSheet` | 自定义 View | sheet | 38 | `mDialog = new QMUIBottomSheet(mContext, style);` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBottomSheetBaseBuilder.java:116 |
| `QMUIRoundButton` | 自定义 View | button | 33 | `<com.qmuiteam.qmui.widget.roundwidget.QMUIRoundButton` qmui/src/main/res/layout/qmui_empty_view.xml:75 |
| `QMUIWindowInsetLayout` | 自定义 View |  | 26 | `<com.qmuiteam.qmui.widget.QMUIWindowInsetLayout xmlns:android="http://schemas.android.com/apk/res/android"` qmuidemo/src/main/res/layout/activity_arch_test.xml:18 |
| `QMUIGroupListView` | 自定义 View | card | 25 | `<com.qmuiteam.qmui.widget.grouplist.QMUIGroupListView` qmuidemo/src/main/res/layout/fragment_about.xml:54 |
| `QMUIQQFaceView` | 自定义 View |  | 25 | `<com.qmuiteam.qmui.qqface.QMUIQQFaceView` qmuidemo/src/main/res/layout/fragment_qqface_layout.xml:97 |
| `QMUIPriorityLinearLayout` | 自定义 View |  | 18 | `<com.qmuiteam.qmui.layout.QMUIPriorityLinearLayout` qmuidemo/src/main/res/layout/fragment_priority_linear_layout.xml:45 |
| `QMUIWindowInsetLayout2` | 自定义 View |  | 16 | `<com.qmuiteam.qmui.widget.QMUIWindowInsetLayout2 xmlns:android="http://schemas.android.com/apk/res/android"` qmuidemo/src/main/res/layout/fragment_animation_listview.xml:18 |
| `QMUISpanTouchFixTextView` | 自定义 View |  | 15 | `<com.qmuiteam.qmui.widget.textview.QMUISpanTouchFixTextView` qmui/src/main/res/layout/qmui_empty_view.xml:36 |
| `QMUITipDialog` | 自定义 View | toast | 14 | `@Widget(widgetClass = QMUITipDialog.class, iconRes = R.mipmap.icon_grid_tip_dialog)` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDTipDialogFragment.java:44 |
| `QMUIPullRefreshLayout` | 自定义 View |  | 11 | `<com.qmuiteam.qmui.widget.pullRefreshLayout.QMUIPullRefreshLayout` qmuidemo/src/main/res/layout/fragment_continuous_nested_scroll.xml:27 |
| `QMUIFrameLayout` | 自定义 View |  | 10 | `<com.qmuiteam.qmui.layout.QMUIFrameLayout` qmuidemo/src/main/res/layout/fragment_button.xml:39 |
| `QMUIDialog` | 自定义 View | dialog | 10 | `mDialog = new QMUIDialog(mContext, style);` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialogBuilder.java:343 |
| `QMUIRadiusImageView2` | 自定义 View | image | 8 | `<com.qmuiteam.qmui.widget.QMUIRadiusImageView2` qmuidemo/src/main/res/layout/fragment_order_detail.xml:218 |
| `QMUIWebView` | 自定义 View |  | 7 | `private QMUIWebView.OnScrollChangeListener mOnScrollChangeListener;` qmui/src/main/java/com/qmuiteam/qmui/widget/webview/QMUIWebViewContainer.java:34 |
| `QMUIRadiusImageView` | 自定义 View | image | 7 | `<com.qmuiteam.qmui.widget.QMUIRadiusImageView` qmuidemo/src/main/res/layout/fragment_radius_imageview.xml:42 |
| `QMUIProgressBar` | 自定义 View | loading | 6 | `<com.qmuiteam.qmui.widget.QMUIProgressBar` qmuidemo/src/main/res/layout/fragment_progressbar.xml:38 |
| `QMUIContinuousNestedBottomRecyclerView` | 自定义 View | list | 6 | `QMUIContinuousNestedBottomRecyclerView recyclerView = new QMUIContinuousNestedBottomRecyclerView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousBottomView.java:96 |
| `QMUILinearLayout` | 自定义 View |  | 6 | `<com.qmuiteam.qmui.layout.QMUILinearLayout style="@style/button_wrapper_style"` qmuidemo/src/main/res/layout/fragment_button.xml:112 |
| `QMUITopBar` | 自定义 View | navbar | 6 | `<com.qmuiteam.qmui.widget.QMUITopBar` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:54 |
| `QMUIDialogMenuItemView` | 自定义 View | dialog | 5 | `itemView.setListener(new QMUIDialogMenuItemView.MenuItemViewListener() {` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialog.java:494 |
| `QMUISlider` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.widget.QMUISlider` qmuidemo/src/main/res/layout/fragment_slider.xml:35 |
| `QMUILoadingView` | 自定义 View | loading | 5 | `<com.qmuiteam.qmui.widget.QMUILoadingView` qmui/src/main/res/layout/qmui_empty_view.xml:22 |
| `QMUIContinuousNestedScrollLayout` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.nestedScroll.QMUIContinuousNestedScrollLayout` qmuidemo/src/main/res/layout/fragment_continuous_nested_scroll.xml:31 |
| `QMUIVerticalTextView` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.widget.QMUIVerticalTextView` qmuidemo/src/main/res/layout/fragment_pull_horizontal_test_layout.xml:48 |
| `QMUILinkTextView` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.widget.textview.QMUILinkTextView` qmuidemo/src/main/res/layout/fragment_link_texview_layout.xml:49 |
| `QMUIContinuousNestedTopRecyclerView` | 自定义 View | list | 4 | `mTopRecyclerView = new QMUIContinuousNestedTopRecyclerView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll3Fragment.java:51 |
| `QMUIBasicTabSegment` | 自定义 View | tabs | 4 | `* Please use QMUIBasicTabSegment.OnTabClickListener for a replacement` qmui/src/main/java/com/qmuiteam/qmui/widget/tab/QMUITabSegment.java:307 |
| `QMUIFloatLayout` | 自定义 View | popup | 4 | `<com.qmuiteam.qmui.widget.QMUIFloatLayout` qmuidemo/src/main/res/layout/fragment_floatlayout.xml:37 |
| `QMUIAnimationListView` | 自定义 View | list | 4 | `<com.qmuiteam.qmui.widget.QMUIAnimationListView` qmuidemo/src/main/res/layout/fragment_animation_listview.xml:32 |
| `QMUIContinuousNestedTopWebView` | 自定义 View |  | 4 | `mNestedWebView = new QMUIContinuousNestedTopWebView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll1Fragment.java:65 |
| `QMUIContinuousNestedTopDelegateLayout` | 自定义 View |  | 4 | `mTopDelegateLayout = new QMUIContinuousNestedTopDelegateLayout(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll4Fragment.java:56 |
| `QMUIStickySectionLayout` | 自定义 View |  | 3 | `<com.qmuiteam.qmui.widget.section.QMUIStickySectionLayout` qmuidemo/src/main/res/layout/fragment_section_layout.xml:34 |
| `QDContinuousBottomView` | 自定义 View |  | 3 | `mBottomView = new QDContinuousBottomView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll2Fragment.java:48 |
| `QMUICollapsingTopBarLayout` | 自定义 View |  | 3 | `<com.qmuiteam.qmui.widget.QMUICollapsingTopBarLayout` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:31 |
| `QMUIRoundLinearLayout` | 自定义 View |  | 3 | `<com.qmuiteam.qmui.widget.roundwidget.QMUIRoundLinearLayout` qmuidemo/src/main/res/layout/fragment_order_detail.xml:72 |
| `QMUIDialogView` | 自定义 View | dialog | 3 | `private QMUIDialogView.OnDecorationListener mOnDecorationListener;` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialogBuilder.java:94 |
| `EmojiEditText` | 自定义 View | input | 3 | `val editText = EmojiEditText(context).apply {` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEmojiInputFragment.kt:38 |
| `QMUIBottomSheetGridLineLayout` | 自定义 View |  | 3 | `private QMUIBottomSheetGridLineLayout.ItemWidthCalculator mItemWidthCalculator = null;` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBottomSheet.java:521 |
| `QMUIDraggableScrollBar` | 自定义 View |  | 2 | `QMUIContinuousNestedTopAreaBehavior.Callback, QMUIDraggableScrollBar.Callback {` qmui/src/main/java/com/qmuiteam/qmui/nestedScroll/QMUIContinuousNestedScrollLayout.java:38 |
| `BlankBox` | Composable |  | 3 | `BlankBox()` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:80 |
| … | | | | 共 158 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| button（Compose） | CommonImageButton（Composable，用 6 次）<br>CommonTextButton（Composable，用 2 次）<br>CommonButton（Composable，用 2 次）<br>OriginOpenButton（Composable，用 2 次）<br>PickCurrentCheckButton（Composable，用 1 次） | `CommonImageButton` |
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
| #1F2329 | 15 | — | qmuidemo/src/main/res/layout/fragment_order_detail.xml |
| #646A73 | 6 | — | qmuidemo/src/main/res/layout/fragment_order_detail.xml |
| #8F959E | 4 | — | qmuidemo/src/main/res/layout/fragment_order_detail.xml |
| #E5E6EB | 3 | — | qmuidemo/src/main/res/layout/fragment_order_detail.xml |
| #333333 | 2 | — | photo/src/main/java/com/qmuiteam/photo/compose/picker/Config.kt |
| #FF0000 | 2 | — | qmuidemo/src/main/res/layout/fragment_layout.xml、qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDLayoutFragment.java |
| #FF8F1F | 2 | — | qmuidemo/src/main/res/layout/fragment_order_detail.xml |
| #FFFFFF | 2 | `@color/qmui_config_color_white` | qmuidemo/src/main/res/layout/fragment_order_detail.xml |
| #0000FF | 2 | — | qmuidemo/src/main/res/layout/fragment_layout.xml、qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDLayoutFragment.java |
| #666666 | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #F5F6F7 | 1 | — | qmuidemo/src/main/res/layout/fragment_order_detail.xml |
| #F53F3F | 1 | — | qmuidemo/src/main/res/layout/fragment_order_detail.xml |
