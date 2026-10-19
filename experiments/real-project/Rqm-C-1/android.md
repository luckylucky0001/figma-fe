
## 页面复用检查：qmuidemo/src/main/java/com/qmuiteam/qmuidemo/manager/QDDataManager.java,qmuidemo/src/main/res/values/attr.xml,qmuidemo/src/main/res/values/colors.xml,qmuidemo/src/main/res/values/styles.xml,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDOrderDetailFragment.kt,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/view/QDActionBar.kt,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/view/QDCountDownView.kt,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/view/QDNoticeBarView.kt,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/view/QDPriceView.kt,qmuidemo/src/main/res/layout/fragment_order_detail.xml  ✅ PASS

- 复用的现成组件（9）：QMUIButton（自定义 View）、QMUICommonListItemView（自定义 View）、QMUIGroupListView（自定义 View）、QMUILinearLayout（自定义 View）、QMUIRadiusImageView2（自定义 View）、QMUIRoundButton（自定义 View）、QMUITopBarLayout（自定义 View）、QMUIWindowInsetLayout（自定义 View）、TextView（组件库）
- 自造轮子（0）
- 有组件却没用（0）
- 疑似手写已有组件（0）
- 页面私有调色板：页面里自己定义了 10 个颜色变量（应使用项目 / 组件库已有的 token）
- 写死颜色（0）

---

# 存量复用体检：/home/user/exp/Rqm-C-1（android）

技术栈：XML View / Jetpack Compose / Material Components / QMUI Android · 扫描文件 570

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 98%（953 / 972） |
| 重复概念（同一 UI 概念有多个实现） | 12 |
| 重复结构（可抽成组件的复制粘贴结构） | 1 组 |
| 硬编码颜色 | 11 处（其中 1 处已有同值 token） |
| 硬编码尺寸 | 394 处 |

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Android SDK widgets | 130 | TextView×112、RecyclerView×6、RadioButton×5、SeekBar×3、ListView×2、ViewPager×1、ViewPager2×1 |
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
- `Slider`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `TextInputLayout`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `TextView` | 组件库 |  | 112 | `<TextView` qmui/src/main/res/layout/qmui_common_list_item.xml:45 |
| `RecyclerView` | 组件库 | list | 6 | `<androidx.recyclerview.widget.RecyclerView` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:65 |
| `RadioButton` | 组件库 | switch | 5 | `<RadioButton` qmuidemo/src/main/res/layout/fragment_layout.xml:189 |
| `QMUITabSegment` | 组件库 | tabs | 4 | `<com.qmuiteam.qmui.widget.tab.QMUITabSegment` qmuidemo/src/main/res/layout/fragment_fsw_viewpager.xml:30 |
| `SeekBar` | 组件库 |  | 3 | `<SeekBar` qmuidemo/src/main/res/layout/fragment_colorhelper.xml:124 |
| `ListView` | 组件库 | list | 2 | `<ListView` qmuidemo/src/main/res/layout/fragment_listview.xml:25 |
| `QMUISeekBar` | 组件库 |  | 1 | `<com.qmuiteam.qmui.widget.QMUISeekBar` qmuidemo/src/main/res/layout/fragment_slider.xml:44 |
| `ViewPager2` | 组件库 |  | 1 | `<androidx.viewpager2.widget.ViewPager2` qmuidemo/src/main/res/layout/fragment_tab_viewpager2_layout.xml:39 |
| `ViewPager` | 组件库 | carousel | 1 | `<androidx.viewpager.widget.ViewPager` qmuidemo/src/main/res/layout/fragment_tab_viewpager_layout.xml:39 |
| `SwipeBackLayout` | 自定义 View |  | 78 | `import static com.qmuiteam.qmui.arch.SwipeBackLayout.DRAG_DIRECTION_BOTTOM_TO_TOP;` arch/src/main/java/com/qmuiteam/qmui/arch/QMUIActivity.java:37 |
| `QMUICommonListItemView` | 自定义 View | cell | 52 | `QMUICommonListItemView itemView = new QMUICommonListItemView(getContext());` qmui/src/main/java/com/qmuiteam/qmui/widget/grouplist/QMUIGroupListView.java:83 |
| `QMUIPullLayout` | 自定义 View |  | 49 | `<com.qmuiteam.qmui.widget.pullLayout.QMUIPullLayout` qmuidemo/src/main/res/layout/fragment_pull_horizontal_test_layout.xml:32 |
| `QMUITopBarLayout` | 自定义 View | navbar | 46 | `<com.qmuiteam.qmui.widget.QMUITopBarLayout` qmuidemo/src/main/res/layout/activity_arch_test.xml:23 |
| `QMUIBottomSheet` | 自定义 View | popup | 38 | `mDialog = new QMUIBottomSheet(mContext, style);` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBottomSheetBaseBuilder.java:116 |
| `QMUIQuickAction` | 自定义 View |  | 31 | `return new QMUIQuickAction(context, ViewGroup.LayoutParams.WRAP_CONTENT, actionHeight)` qmui/src/main/java/com/qmuiteam/qmui/widget/popup/QMUIPopups.java:71 |
| `QMUIGroupListView` | 自定义 View | card | 31 | `<com.qmuiteam.qmui.widget.grouplist.QMUIGroupListView` qmuidemo/src/main/res/layout/fragment_about.xml:54 |
| `QMUIRoundButton` | 自定义 View | button | 30 | `<com.qmuiteam.qmui.widget.roundwidget.QMUIRoundButton` qmui/src/main/res/layout/qmui_empty_view.xml:75 |
| `QMUIWindowInsetLayout` | 自定义 View |  | 26 | `<com.qmuiteam.qmui.widget.QMUIWindowInsetLayout xmlns:android="http://schemas.android.com/apk/res/android"` qmuidemo/src/main/res/layout/activity_arch_test.xml:18 |
| `QMUIQQFaceView` | 自定义 View |  | 25 | `<com.qmuiteam.qmui.qqface.QMUIQQFaceView` qmuidemo/src/main/res/layout/fragment_qqface_layout.xml:97 |
| `QMUIPriorityLinearLayout` | 自定义 View |  | 18 | `<com.qmuiteam.qmui.layout.QMUIPriorityLinearLayout` qmuidemo/src/main/res/layout/fragment_priority_linear_layout.xml:45 |
| `QMUIWindowInsetLayout2` | 自定义 View |  | 16 | `<com.qmuiteam.qmui.widget.QMUIWindowInsetLayout2 xmlns:android="http://schemas.android.com/apk/res/android"` qmuidemo/src/main/res/layout/fragment_animation_listview.xml:18 |
| `QMUISpanTouchFixTextView` | 自定义 View |  | 15 | `<com.qmuiteam.qmui.widget.textview.QMUISpanTouchFixTextView` qmui/src/main/res/layout/qmui_empty_view.xml:36 |
| `QMUITipDialog` | 自定义 View | dialog | 14 | `@Widget(widgetClass = QMUITipDialog.class, iconRes = R.mipmap.icon_grid_tip_dialog)` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDTipDialogFragment.java:44 |
| `QMUIPopup` | 自定义 View | popup | 14 | `return new QMUIPopup(context,` qmui/src/main/java/com/qmuiteam/qmui/widget/popup/QMUIPopups.java:30 |
| `QMUIPullRefreshLayout` | 自定义 View |  | 11 | `<com.qmuiteam.qmui.widget.pullRefreshLayout.QMUIPullRefreshLayout` qmuidemo/src/main/res/layout/fragment_continuous_nested_scroll.xml:27 |
| `QMUILinearLayout` | 自定义 View |  | 10 | `<com.qmuiteam.qmui.layout.QMUILinearLayout style="@style/button_wrapper_style"` qmuidemo/src/main/res/layout/fragment_button.xml:112 |
| `QMUIDialog` | 自定义 View | dialog | 10 | `mDialog = new QMUIDialog(mContext, style);` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialogBuilder.java:343 |
| `QMUIFrameLayout` | 自定义 View |  | 10 | `<com.qmuiteam.qmui.layout.QMUIFrameLayout` qmuidemo/src/main/res/layout/fragment_button.xml:39 |
| `QMUIRadiusImageView2` | 自定义 View | image | 8 | `<com.qmuiteam.qmui.widget.QMUIRadiusImageView2` qmuidemo/src/main/res/layout/fragment_order_detail.xml:134 |
| `QMUIRadiusImageView` | 自定义 View | image | 7 | `<com.qmuiteam.qmui.widget.QMUIRadiusImageView` qmuidemo/src/main/res/layout/fragment_radius_imageview.xml:42 |
| `QMUIWebView` | 自定义 View |  | 7 | `private QMUIWebView.OnScrollChangeListener mOnScrollChangeListener;` qmui/src/main/java/com/qmuiteam/qmui/widget/webview/QMUIWebViewContainer.java:34 |
| `QMUIContinuousNestedBottomRecyclerView` | 自定义 View |  | 6 | `QMUIContinuousNestedBottomRecyclerView recyclerView = new QMUIContinuousNestedBottomRecyclerView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousBottomView.java:96 |
| `QMUITopBar` | 自定义 View | navbar | 6 | `<com.qmuiteam.qmui.widget.QMUITopBar` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:54 |
| `QMUIProgressBar` | 自定义 View | loading | 6 | `<com.qmuiteam.qmui.widget.QMUIProgressBar` qmuidemo/src/main/res/layout/fragment_progressbar.xml:38 |
| `QMUIContinuousNestedScrollLayout` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.nestedScroll.QMUIContinuousNestedScrollLayout` qmuidemo/src/main/res/layout/fragment_continuous_nested_scroll.xml:31 |
| `QMUISlider` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.widget.QMUISlider` qmuidemo/src/main/res/layout/fragment_slider.xml:35 |
| `QMUILoadingView` | 自定义 View | loading | 5 | `<com.qmuiteam.qmui.widget.QMUILoadingView` qmui/src/main/res/layout/qmui_empty_view.xml:22 |
| `QMUIButton` | 自定义 View | button | 5 | `<com.qmuiteam.qmui.layout.QMUIButton` qmuidemo/src/main/res/layout/fragment_order_detail.xml:273 |
| `QMUIVerticalTextView` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.widget.QMUIVerticalTextView` qmuidemo/src/main/res/layout/fragment_pull_horizontal_test_layout.xml:48 |
| `QMUILinkTextView` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.widget.textview.QMUILinkTextView` qmuidemo/src/main/res/layout/fragment_link_texview_layout.xml:49 |
| `QMUIDialogMenuItemView` | 自定义 View | dialog | 5 | `itemView.setListener(new QMUIDialogMenuItemView.MenuItemViewListener() {` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialog.java:494 |
| `QMUIContinuousNestedTopWebView` | 自定义 View |  | 4 | `mNestedWebView = new QMUIContinuousNestedTopWebView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll1Fragment.java:65 |
| `QMUIBasicTabSegment` | 自定义 View | tabs | 4 | `* Please use QMUIBasicTabSegment.OnTabClickListener for a replacement` qmui/src/main/java/com/qmuiteam/qmui/widget/tab/QMUITabSegment.java:307 |
| `QMUIAnimationListView` | 自定义 View | list | 4 | `<com.qmuiteam.qmui.widget.QMUIAnimationListView` qmuidemo/src/main/res/layout/fragment_animation_listview.xml:32 |
| `QMUIContinuousNestedTopRecyclerView` | 自定义 View |  | 4 | `mTopRecyclerView = new QMUIContinuousNestedTopRecyclerView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll3Fragment.java:51 |
| `QMUIContinuousNestedTopDelegateLayout` | 自定义 View |  | 4 | `mTopDelegateLayout = new QMUIContinuousNestedTopDelegateLayout(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll4Fragment.java:56 |
| `QMUIFloatLayout` | 自定义 View |  | 4 | `<com.qmuiteam.qmui.widget.QMUIFloatLayout` qmuidemo/src/main/res/layout/fragment_floatlayout.xml:37 |
| `QDContinuousBottomView` | 自定义 View |  | 3 | `mBottomView = new QDContinuousBottomView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll2Fragment.java:48 |
| `EmojiEditText` | 自定义 View |  | 3 | `val editText = EmojiEditText(context).apply {` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEmojiInputFragment.kt:38 |
| `QMUIStickySectionLayout` | 自定义 View |  | 3 | `<com.qmuiteam.qmui.widget.section.QMUIStickySectionLayout` qmuidemo/src/main/res/layout/fragment_section_layout.xml:34 |
| `QMUIDialogView` | 自定义 View | dialog | 3 | `private QMUIDialogView.OnDecorationListener mOnDecorationListener;` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialogBuilder.java:94 |
| `QMUIBottomSheetGridLineLayout` | 自定义 View |  | 3 | `private QMUIBottomSheetGridLineLayout.ItemWidthCalculator mItemWidthCalculator = null;` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBottomSheet.java:521 |
| `QMUIFullScreenPopup` | 自定义 View | popup | 3 | `return new QMUIFullScreenPopup(context);` qmui/src/main/java/com/qmuiteam/qmui/widget/popup/QMUIPopups.java:67 |
| `QMUICollapsingTopBarLayout` | 自定义 View |  | 3 | `<com.qmuiteam.qmui.widget.QMUICollapsingTopBarLayout` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:31 |
| `Test` | 自定义 View |  | 2 | `<com.qmuiteam.qmuidemo.fragment.components.qqface.Test` qmuidemo/src/main/res/layout/fragment_qqface_layout.xml:64 |
| `QMUIDialogTitle` | Composable | dialog | 1 | `QMUIDialogTitle(title)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:109 |
| `QMUIPhotoPickerGrid` | Composable |  | 1 | `QMUIPhotoPickerGrid(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:372 |
| `EmojiLayout` | 自定义 View |  | 2 | `return EmojiLayout(requireContext())` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEmojiInputFragment.kt:29 |
| `QMUIRoundFrameLayout` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/widget/roundwidget/QMUIRoundFrameLayout.java |
| `QMUIWebViewContainer` | 自定义 View |  | 1 | `<com.qmuiteam.qmui.widget.webview.QMUIWebViewContainer` qmuidemo/src/main/res/layout/fragment_webview_explorer.xml:24 |
| … | | | | 共 166 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| button | QMUIRoundButton（自定义 View，用 30 次）<br>CommonImageButton（Composable，用 6 次）<br>QMUIButton（自定义 View，用 5 次）<br>TextButton（Composable，用 4 次）<br>CommonButton（Composable，用 2 次）<br>QMUIAlphaImageButton（自定义 View，用 2 次）<br>CommonTextButton（Composable，用 2 次）<br>OriginOpenButton（Composable，用 2 次）<br>PickCurrentCheckButton（Composable，用 1 次）<br>QMUIAlphaButton（自定义 View，用 0 次） | `QMUIRoundButton` |
| cell | QMUICommonListItemView（自定义 View，用 52 次）<br>QMUIItem（Composable，用 10 次）<br>QMUIBitmapRegionItem（Composable，用 2 次）<br>QMUIPhotoPickerPreviewPickedItems（Composable，用 1 次）<br>QMUIPhotoBucketItem（Composable，用 1 次） | `QMUICommonListItemView` |
| dialog | QMUITipDialog（自定义 View，用 14 次）<br>QMUIDialog（自定义 View，用 10 次）<br>QMUIDialog（Composable，用 5 次）<br>QMUIBaseDialog（自定义 View，用 0 次）<br>QMUIModal（Composable，用 0 次） | `QMUITipDialog` |
| image | QMUIRadiusImageView2（自定义 View，用 8 次）<br>QMUIRadiusImageView（自定义 View，用 7 次）<br>RowImages（Composable，用 2 次）<br>LongImage（Composable，用 2 次）<br>TestImageCompress（Composable，用 1 次） | `QMUIRadiusImageView2` |
| navbar | QMUITopBar（Composable，用 7 次）<br>QMUITopBar（自定义 View，用 6 次）<br>QMUIGroupListSectionHeaderFooterView（自定义 View，用 2 次）<br>QDSectionHeaderView（自定义 View，用 1 次） | `QMUITopBar` |
| loading | QMUIProgressBar（自定义 View，用 6 次）<br>QMUILoadingView（自定义 View，用 5 次）<br>QMUIPhotoLoading（Composable，用 3 次）<br>QDLoadingItemView（自定义 View，用 1 次） | `QMUIProgressBar` |
| list | QMUIAnimationListView（自定义 View，用 4 次）<br>QMUIWrapContentListView（自定义 View，用 1 次）<br>RecyclerView（组件库，用 6 次）<br>ListView（组件库，用 2 次） | `QMUIAnimationListView` |
| switch | QMUICheckBox（Composable，用 1 次）<br>QMUIPhotoPickerBucketToggleArrow（Composable，用 1 次）<br>QMUIPhotoPickRadio（Composable，用 1 次）<br>RadioButton（组件库，用 5 次） | `QMUICheckBox` |
| popup | QMUIBottomSheet（自定义 View，用 38 次）<br>QMUIPopup（自定义 View，用 14 次）<br>QMUIFullScreenPopup（自定义 View，用 3 次） | `QMUIBottomSheet` |
| tabs | QMUITabSegment（组件库，用 4 次）<br>QMUIBasicTabSegment（自定义 View，用 4 次）<br>QMUITabView（自定义 View，用 2 次） | `QMUITabSegment` |
| carousel | SwipeBackgroundView（自定义 View，用 2 次）<br>ViewPager（组件库，用 1 次） | `SwipeBackgroundView` |
| icon | QMUIMarkIcon（Composable，用 1 次）<br>QMUIChevronIcon（Composable，用 0 次） | `QMUIMarkIcon` |

## 4. 有组件却没用（19 处）

- qmuidemo/src/main/res/layout/drawablehelper_createfromview.xml:25 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）
- qmuidemo/src/main/res/layout/fragment_about.xml:40 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）
- qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:44 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）
- qmuidemo/src/main/res/layout/fragment_colorhelper.xml:79 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）
- qmuidemo/src/main/res/layout/fragment_drawablehelper.xml:54 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）
- qmuidemo/src/main/res/layout/fragment_drawablehelper.xml:66 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）
- qmuidemo/src/main/res/layout/fragment_drawablehelper.xml:85 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）
- qmuidemo/src/main/res/layout/fragment_drawablehelper.xml:91 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）
- qmuidemo/src/main/res/layout/fragment_drawablehelper.xml:98 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）
- qmuidemo/src/main/res/layout/fragment_layout.xml:114 `<Button` → 原生 Button：项目里已有 QMUIRoundButton（组件库）
- qmuidemo/src/main/res/layout/fragment_layout.xml:125 `<Button` → 原生 Button：项目里已有 QMUIRoundButton（组件库）
- qmuidemo/src/main/res/layout/fragment_layout.xml:148 `<Button` → 原生 Button：项目里已有 QMUIAlphaButton（组件库）
- qmuidemo/src/main/res/layout/fragment_layout.xml:158 `<Button` → 原生 Button：项目里已有 QMUIRoundButton（组件库）
- qmuidemo/src/main/res/layout/fragment_layout.xml:167 `<Button` → 原生 Button：项目里已有 QMUIRoundButton（组件库）
- qmuidemo/src/main/res/layout/fragment_scheme.xml:47 `<EditText` → 原生 EditText：项目里已有 TextInputEditText（组件库）
- qmuidemo/src/main/res/layout/fragment_verticaltextview.xml:62 `<EditText` → 原生 EditText：项目里已有 TextInputEditText（组件库）
- qmuidemo/src/main/res/layout/fragment_webview_explorer.xml:37 `<ProgressBar` → 原生 ProgressBar：项目里已有 QMUIProgressBar（组件库）
- qmuidemo/src/main/res/layout/home_item_layout.xml:28 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）
- qmuidemo/src/main/res/layout/tipdialog_custom.xml:26 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）

## 5. 重复结构（复制粘贴的 UI 片段，可抽成组件）

- **5 个元素 × 6 处**：qmuidemo/src/main/res/layout/fragment_priority_linear_layout.xml:45、qmuidemo/src/main/res/layout/fragment_priority_linear_layout.xml:92、qmuidemo/src/main/res/layout/fragment_priority_linear_layout.xml:136、qmuidemo/src/main/res/layout/fragment_priority_linear_layout.xml:180、qmuidemo/src/main/res/layout/fragment_priority_linear_layout.xml:224、qmuidemo/src/main/res/layout/fragment_priority_linear_layout.xml:270

```
<com.qmuiteam.qmui.layout.QMUIPriorityLinearLayout>
  <TextView>
  <TextView>
  <TextView>
  <TextView>
```

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
