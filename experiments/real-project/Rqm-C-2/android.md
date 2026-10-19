
## 页面复用检查：15 个文件（qmuidemo/src/main/java/com/qmuiteam/qmuidemo/manager/QDDataManager.java、qmuidemo/src/main/res/values/attr.xml …）  ✅ PASS

- 复用的现成组件（8）：QMUIButton（自定义 View）、QMUICommonListItemView（自定义 View）、QMUIGroupListView（自定义 View）、QMUILinearLayout（自定义 View）、QMUIRadiusImageView（自定义 View）、QMUITopBarLayout（自定义 View）、QMUIWindowInsetLayout（自定义 View）、TextView（组件库）
- 自造轮子（0）
- 有组件却没用（0）
- 疑似手写已有组件（0）
- 页面私有调色板：页面里自己定义了 8 个颜色变量（应使用项目 / 组件库已有的 token）
- 写死颜色（0）

---

# 存量复用体检：/home/user/exp/Rqm-C-2（android）

技术栈：XML View / Jetpack Compose / Material Components / QMUI Android / Compose Foundation · 扫描文件 561

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（974 / 974） |
| 重复概念（同一 UI 概念有多个实现） | 15 |
| 重复结构（可抽成组件的复制粘贴结构） | 0 组 |
| 硬编码颜色 | 11 处（其中 1 处已有同值 token） |
| 硬编码尺寸 | 361 处 |

## 0. 优先处理（按影响面排序）

1. 「cell（View）」有 2 种实现（共用 54 次）→ 统一到 `QMUICommonListItemView`
2. 「button（View）」有 5 种实现（共用 41 次）→ 统一到 `QMUIRoundButton`
3. 「image（View）」有 3 种实现（共用 26 次）→ 统一到 `QMUIRadiusImageView`
4. 「list（View）」有 6 种实现（共用 23 次）→ 统一到 `QMUIContinuousNestedBottomRecyclerView`
5. 「popup（View）」有 3 种实现（共用 21 次）→ 统一到 `QMUIPopup`

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Android SDK widgets | 149 | TextView×112、ImageView×11、RecyclerView×6、RadioButton×5、Button×5、SeekBar×3、ListView×2、EditText×2 |
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
| `TextView` | 组件库 |  | 112 | `<TextView` qmui/src/main/res/layout/qmui_common_list_item.xml:45 |
| `ImageView` | 组件库 | image | 11 | `<ImageView` qmuidemo/src/main/res/layout/drawablehelper_createfromview.xml:25 |
| `LazyColumn` | 组件库 | list | 9 | `LazyColumn(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIBottomSheet.kt:40 |
| `RecyclerView` | 组件库 | list | 6 | `<androidx.recyclerview.widget.RecyclerView` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:65 |
| `RadioButton` | 组件库 | radio | 5 | `<RadioButton` qmuidemo/src/main/res/layout/fragment_layout.xml:189 |
| `Button` | 组件库 | button | 5 | `<Button` qmuidemo/src/main/res/layout/fragment_layout.xml:114 |
| `QMUITabSegment` | 组件库 | tabs | 4 | `<com.qmuiteam.qmui.widget.tab.QMUITabSegment` qmuidemo/src/main/res/layout/fragment_fsw_viewpager.xml:30 |
| `SeekBar` | 组件库 |  | 3 | `<SeekBar` qmuidemo/src/main/res/layout/fragment_colorhelper.xml:124 |
| `EditText` | 组件库 | input | 2 | `<EditText` qmuidemo/src/main/res/layout/fragment_scheme.xml:47 |
| `ListView` | 组件库 | list | 2 | `<ListView` qmuidemo/src/main/res/layout/fragment_listview.xml:25 |
| `HorizontalPager` | 组件库 | carousel | 1 | `HorizontalPager(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:39 |
| `ProgressBar` | 组件库 | loading | 1 | `<ProgressBar` qmuidemo/src/main/res/layout/fragment_webview_explorer.xml:37 |
| `ViewPager` | 组件库 | carousel | 1 | `<androidx.viewpager.widget.ViewPager` qmuidemo/src/main/res/layout/fragment_tab_viewpager_layout.xml:39 |
| `LazyRow` | 组件库 | list | 1 | `LazyRow(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:116 |
| `ViewPager2` | 组件库 |  | 1 | `<androidx.viewpager2.widget.ViewPager2` qmuidemo/src/main/res/layout/fragment_tab_viewpager2_layout.xml:39 |
| `QMUISeekBar` | 组件库 |  | 1 | `<com.qmuiteam.qmui.widget.QMUISeekBar` qmuidemo/src/main/res/layout/fragment_slider.xml:44 |
| `SwipeBackLayout` | 自定义 View |  | 78 | `import static com.qmuiteam.qmui.arch.SwipeBackLayout.DRAG_DIRECTION_BOTTOM_TO_TOP;` arch/src/main/java/com/qmuiteam/qmui/arch/QMUIActivity.java:37 |
| `QMUICommonListItemView` | 自定义 View | cell | 53 | `QMUICommonListItemView itemView = new QMUICommonListItemView(getContext());` qmui/src/main/java/com/qmuiteam/qmui/widget/grouplist/QMUIGroupListView.java:83 |
| `QMUIPullLayout` | 自定义 View |  | 49 | `<com.qmuiteam.qmui.widget.pullLayout.QMUIPullLayout` qmuidemo/src/main/res/layout/fragment_pull_horizontal_test_layout.xml:32 |
| `QMUITopBarLayout` | 自定义 View | navbar | 46 | `<com.qmuiteam.qmui.widget.QMUITopBarLayout` qmuidemo/src/main/res/layout/activity_arch_test.xml:23 |
| `QMUIBottomSheet` | 自定义 View | sheet | 38 | `mDialog = new QMUIBottomSheet(mContext, style);` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBottomSheetBaseBuilder.java:116 |
| `QMUIQuickAction` | 自定义 View |  | 31 | `return new QMUIQuickAction(context, ViewGroup.LayoutParams.WRAP_CONTENT, actionHeight)` qmui/src/main/java/com/qmuiteam/qmui/widget/popup/QMUIPopups.java:71 |
| `QMUIDialogAction` | Composable | dialog | 8 | `QMUIDialogAction(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:90 |
| `QDCountDownView` | 自定义 View | countdown | 1 | `<com.qmuiteam.qmuidemo.view.QDCountDownView` qmuidemo/src/main/res/layout/fragment_order_detail.xml:68 |
| `QMUIWrapContentScrollView` | 自定义 View |  | 1 | `QMUIWrapContentScrollView scrollView = new QMUIWrapContentScrollView(view.getContext());` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialogBuilder.java:498 |
| `QMUIRoundLinearLayout` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/widget/roundwidget/QMUIRoundLinearLayout.java |
| `QMUIAlphaLinearLayout` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/alpha/QMUIAlphaLinearLayout.java |
| `QMUIProgressBar` | 自定义 View | loading | 6 | `<com.qmuiteam.qmui.widget.QMUIProgressBar` qmuidemo/src/main/res/layout/fragment_progressbar.xml:38 |
| `QMUIItem` | Composable | cell | 10 | `QMUIItem(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:152 |
| `QMUIContinuousNestedTopWebView` | 自定义 View |  | 4 | `mNestedWebView = new QMUIContinuousNestedTopWebView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll1Fragment.java:65 |
| `QMUIFloatLayout` | 自定义 View | popup | 4 | `<com.qmuiteam.qmui.widget.QMUIFloatLayout` qmuidemo/src/main/res/layout/fragment_floatlayout.xml:37 |
| `QMUIWebView` | 自定义 View |  | 7 | `private QMUIWebView.OnScrollChangeListener mOnScrollChangeListener;` qmui/src/main/java/com/qmuiteam/qmui/widget/webview/QMUIWebViewContainer.java:34 |
| `QMUIAlphaImageButton` | 自定义 View | button | 2 | `QMUIAlphaImageButton imageButton = new QMUIAlphaImageButton(getContext());` qmui/src/main/java/com/qmuiteam/qmui/widget/QMUITopBar.java:656 |
| `QMUIAlphaConstraintLayout` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/alpha/QMUIAlphaConstraintLayout.java |
| `QMUIWindowInsetLayout` | 自定义 View |  | 26 | `<com.qmuiteam.qmui.widget.QMUIWindowInsetLayout xmlns:android="http://schemas.android.com/apk/res/android"` qmuidemo/src/main/res/layout/activity_arch_test.xml:18 |
| `QDLoadingItemView` | 自定义 View | loading | 1 | `return new ViewHolder(new QDLoadingItemView(viewGroup.getContext()));` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/section/QDGridSectionAdapter.java:72 |
| `QMUIAlphaTextView` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/alpha/QMUIAlphaTextView.java |
| `QMUIAlphaRelativeLayout` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/alpha/QMUIAlphaRelativeLayout.java |
| `QMUIPullLoadMoreView` | 自定义 View |  | 1 | `<com.qmuiteam.qmui.widget.pullLayout.QMUIPullLoadMoreView` qmuidemo/src/main/res/layout/fragment_pull_refresh_and_load_more_test_layout.xml:57 |
| `QMUIFontFitTextView` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/widget/QMUIFontFitTextView.java |
| `QMUIBottomSheetGridLineLayout` | 自定义 View |  | 3 | `private QMUIBottomSheetGridLineLayout.ItemWidthCalculator mItemWidthCalculator = null;` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBottomSheet.java:521 |
| `QDShadowAdjustLayout` | 自定义 View |  | 1 | `<com.qmuiteam.qmuidemo.view.QDShadowAdjustLayout` qmuidemo/src/main/res/layout/fragment_layout.xml:23 |
| `QMUILoadingView` | 自定义 View | loading | 5 | `<com.qmuiteam.qmui.widget.QMUILoadingView` qmui/src/main/res/layout/qmui_empty_view.xml:22 |
| `QMUISlider` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.widget.QMUISlider` qmuidemo/src/main/res/layout/fragment_slider.xml:35 |
| `EmojiLayout` | 自定义 View |  | 2 | `return EmojiLayout(requireContext())` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEmojiInputFragment.kt:29 |
| `QDQQFacePagerView` | 自定义 View |  | 1 | `view = new QDQQFacePagerView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/qqface/QDQQFacePerformanceTestFragment.java:142 |
| `QMUITabView` | 自定义 View | tabs | 2 | `public class QMUITabAdapter extends QMUIItemViewsAdapter<QMUITab, QMUITabView> implements QMUITabView.Callback {` qmui/src/main/java/com/qmuiteam/qmui/widget/tab/QMUITabAdapter.java:23 |
| `QMUIButton` | 自定义 View | button | 5 | `<com.qmuiteam.qmui.layout.QMUIButton` qmuidemo/src/main/res/layout/fragment_order_detail.xml:149 |
| `QMUIConstraintLayout` | 自定义 View |  | 2 | `<com.qmuiteam.qmui.layout.QMUIConstraintLayout xmlns:android="http://schemas.android.com/apk/res/android"` qmuidemo/src/main/res/layout/fragment_verticaltextview.xml:18 |
| `QMUIRadiusImageView2` | 自定义 View | image | 7 | `<com.qmuiteam.qmui.widget.QMUIRadiusImageView2` qmuidemo/src/main/res/layout/fragment_radius_imageview2.xml:42 |
| `QDNoticeBarView` | 自定义 View | notice | 1 | `<com.qmuiteam.qmuidemo.view.QDNoticeBarView` qmuidemo/src/main/res/layout/fragment_order_detail.xml:30 |
| `QMUIAppBarLayout` | 自定义 View |  | 1 | `<com.qmuiteam.qmui.widget.QMUIAppBarLayout` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:26 |
| `QMUIContinuousNestedTopLinearLayout` | 自定义 View |  | 1 | `mTopLinearLayout = new QMUIContinuousNestedTopLinearLayout(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll6Fragment.java:55 |
| `QMUIRoundRelativeLayout` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/widget/roundwidget/QMUIRoundRelativeLayout.java |
| `QMUIBaseDialog` | 自定义 View | dialog | 0 | qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBaseDialog.java |
| `QMUIDraggableScrollBar` | 自定义 View |  | 2 | `QMUIContinuousNestedTopAreaBehavior.Callback, QMUIDraggableScrollBar.Callback {` qmui/src/main/java/com/qmuiteam/qmui/nestedScroll/QMUIContinuousNestedScrollLayout.java:38 |
| `QMUIFrameLayout` | 自定义 View |  | 10 | `<com.qmuiteam.qmui.layout.QMUIFrameLayout` qmuidemo/src/main/res/layout/fragment_button.xml:39 |
| `HomeController` | 自定义 View |  | 2 | `HomeController.HomeControlListener listener = new HomeController.HomeControlListener() {` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/home/HomeFragment.java:181 |
| `QMUIRelativeLayout` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/layout/QMUIRelativeLayout.java |
| `QMUIDialogView` | 自定义 View | dialog | 3 | `private QMUIDialogView.OnDecorationListener mOnDecorationListener;` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialogBuilder.java:94 |
| … | | | | 共 167 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| list（View） | QMUIContinuousNestedBottomRecyclerView（自定义 View，用 6 次）<br>QMUIContinuousNestedTopRecyclerView（自定义 View，用 4 次）<br>QMUIAnimationListView（自定义 View，用 4 次）<br>RecyclerView（组件库，用 6 次）<br>QMUIWrapContentListView（自定义 View，用 1 次）<br>ListView（组件库，用 2 次） | `QMUIContinuousNestedBottomRecyclerView` |
| button（View） | QMUIRoundButton（自定义 View，用 29 次）<br>QMUIButton（自定义 View，用 5 次）<br>Button（组件库，用 5 次）<br>QMUIAlphaImageButton（自定义 View，用 2 次）<br>QMUIAlphaButton（自定义 View，用 0 次） | `QMUIRoundButton` |
| button（Compose） | CommonImageButton（Composable，用 6 次）<br>CommonButton（Composable，用 2 次）<br>CommonTextButton（Composable，用 2 次）<br>OriginOpenButton（Composable，用 2 次）<br>PickCurrentCheckButton（Composable，用 1 次） | `CommonImageButton` |
| loading（View） | QMUIProgressBar（自定义 View，用 6 次）<br>QMUILoadingView（自定义 View，用 5 次）<br>QDLoadingItemView（自定义 View，用 1 次）<br>ProgressBar（组件库，用 1 次） | `QMUIProgressBar` |
| image（View） | QMUIRadiusImageView（自定义 View，用 8 次）<br>QMUIRadiusImageView2（自定义 View，用 7 次）<br>ImageView（组件库，用 11 次） | `QMUIRadiusImageView` |
| popup（View） | QMUIPopup（自定义 View，用 14 次）<br>QMUIFloatLayout（自定义 View，用 4 次）<br>QMUIFullScreenPopup（自定义 View，用 3 次） | `QMUIPopup` |
| navbar（View） | QMUITopBar（自定义 View，用 6 次）<br>QMUIGroupListSectionHeaderFooterView（自定义 View，用 2 次）<br>QDSectionHeaderView（自定义 View，用 1 次） | `QMUITopBar` |
| cell（View） | QMUICommonListItemView（自定义 View，用 53 次）<br>QMUIBottomSheetListItemView（自定义 View，用 1 次） | `QMUICommonListItemView` |
| dialog（View） | QMUIDialog（自定义 View，用 10 次）<br>QMUIBaseDialog（自定义 View，用 0 次） | `QMUIDialog` |
| dialog（Compose） | QMUIDialogAction（Composable，用 8 次）<br>QMUIModal（Composable，用 0 次） | `QMUIDialogAction` |
| tabs（View） | QMUIBasicTabSegment（自定义 View，用 4 次）<br>QMUITabView（自定义 View，用 2 次） | `QMUIBasicTabSegment` |
| input（View） | EmojiEditText（自定义 View，用 3 次）<br>EditText（组件库，用 2 次） | `EmojiEditText` |
| carousel（View） | SwipeBackgroundView（自定义 View，用 2 次）<br>ViewPager（组件库，用 1 次） | `SwipeBackgroundView` |
| checkbox（Compose） | QMUICheckBox（Composable，用 1 次）<br>QMUIPhotoPickCheckBox（Composable，用 1 次） | `QMUICheckBox` |
| icon（Compose） | QMUIMarkIcon（Composable，用 1 次）<br>QMUIChevronIcon（Composable，用 0 次） | `QMUIMarkIcon` |

## 6. 硬编码颜色热点

| 颜色 | 次数 | 已有 token | 出现在 |
|---|---|---|---|
| #0000FF | 2 | — | qmuidemo/src/main/res/layout/fragment_layout.xml、qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDLayoutFragment.java |
| #333333 | 2 | — | photo/src/main/java/com/qmuiteam/photo/compose/picker/Config.kt |
| #FF0000 | 2 | — | qmuidemo/src/main/res/layout/fragment_layout.xml、qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDLayoutFragment.java |
| #00A8E1 | 1 | `@color/app_color_blue` | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #CCCCCC | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #777777 | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #666666 | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #222222 | 1 | — | photo/src/main/java/com/qmuiteam/photo/compose/picker/Config.kt |
