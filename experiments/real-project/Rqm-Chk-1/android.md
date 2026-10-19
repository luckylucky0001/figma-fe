
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
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（960 / 960） |
| 重复概念（同一 UI 概念有多个实现） | 13 |
| 重复结构（可抽成组件的复制粘贴结构） | 0 组 |
| 硬编码颜色 | 11 处（其中 1 处已有同值 token） |
| 硬编码尺寸 | 381 处 |

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Android SDK widgets | 142 | TextView×105、ImageView×11、RecyclerView×6、RadioButton×5、Button×5、SeekBar×3、EditText×2、ListView×2 |
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
- `Slider`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `TextInputLayout`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `TextView` | 组件库 |  | 105 | `<TextView` qmui/src/main/res/layout/qmui_common_list_item.xml:45 |
| `ImageView` | 组件库 | image | 11 | `<ImageView` qmuidemo/src/main/res/layout/drawablehelper_createfromview.xml:25 |
| `LazyColumn` | 组件库 | list | 9 | `LazyColumn(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIBottomSheet.kt:40 |
| `RecyclerView` | 组件库 | list | 6 | `<androidx.recyclerview.widget.RecyclerView` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:65 |
| `Button` | 组件库 | button | 5 | `<Button` qmuidemo/src/main/res/layout/fragment_layout.xml:114 |
| `RadioButton` | 组件库 | switch | 5 | `<RadioButton` qmuidemo/src/main/res/layout/fragment_layout.xml:189 |
| `QMUITabSegment` | 组件库 | tabs | 4 | `<com.qmuiteam.qmui.widget.tab.QMUITabSegment` qmuidemo/src/main/res/layout/fragment_fsw_viewpager.xml:30 |
| `SeekBar` | 组件库 |  | 3 | `<SeekBar` qmuidemo/src/main/res/layout/fragment_colorhelper.xml:124 |
| `ListView` | 组件库 | list | 2 | `<ListView` qmuidemo/src/main/res/layout/fragment_listview.xml:25 |
| `EditText` | 组件库 | input | 2 | `<EditText` qmuidemo/src/main/res/layout/fragment_scheme.xml:47 |
| `QMUISeekBar` | 组件库 |  | 1 | `<com.qmuiteam.qmui.widget.QMUISeekBar` qmuidemo/src/main/res/layout/fragment_slider.xml:44 |
| `ProgressBar` | 组件库 | loading | 1 | `<ProgressBar` qmuidemo/src/main/res/layout/fragment_webview_explorer.xml:37 |
| `ViewPager2` | 组件库 |  | 1 | `<androidx.viewpager2.widget.ViewPager2` qmuidemo/src/main/res/layout/fragment_tab_viewpager2_layout.xml:39 |
| `LazyRow` | 组件库 | list | 1 | `LazyRow(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:116 |
| `ViewPager` | 组件库 | carousel | 1 | `<androidx.viewpager.widget.ViewPager` qmuidemo/src/main/res/layout/fragment_tab_viewpager_layout.xml:39 |
| `HorizontalPager` | 组件库 | carousel | 1 | `HorizontalPager(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:39 |
| `SwipeBackLayout` | 自定义 View |  | 78 | `import static com.qmuiteam.qmui.arch.SwipeBackLayout.DRAG_DIRECTION_BOTTOM_TO_TOP;` arch/src/main/java/com/qmuiteam/qmui/arch/QMUIActivity.java:37 |
| `QMUICommonListItemView` | 自定义 View | cell | 56 | `QMUICommonListItemView itemView = new QMUICommonListItemView(getContext());` qmui/src/main/java/com/qmuiteam/qmui/widget/grouplist/QMUIGroupListView.java:83 |
| `QMUIPullLayout` | 自定义 View |  | 49 | `<com.qmuiteam.qmui.widget.pullLayout.QMUIPullLayout` qmuidemo/src/main/res/layout/fragment_pull_horizontal_test_layout.xml:32 |
| `QMUITopBarLayout` | 自定义 View | navbar | 46 | `<com.qmuiteam.qmui.widget.QMUITopBarLayout` qmuidemo/src/main/res/layout/activity_arch_test.xml:23 |
| `QMUIBottomSheet` | 自定义 View | popup | 38 | `mDialog = new QMUIBottomSheet(mContext, style);` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBottomSheetBaseBuilder.java:116 |
| `QMUIRoundButton` | 自定义 View | button | 32 | `<com.qmuiteam.qmui.widget.roundwidget.QMUIRoundButton` qmui/src/main/res/layout/qmui_empty_view.xml:75 |
| `QMUIQuickAction` | 自定义 View |  | 31 | `return new QMUIQuickAction(context, ViewGroup.LayoutParams.WRAP_CONTENT, actionHeight)` qmui/src/main/java/com/qmuiteam/qmui/widget/popup/QMUIPopups.java:71 |
| `QMUIGroupListView` | 自定义 View | card | 29 | `<com.qmuiteam.qmui.widget.grouplist.QMUIGroupListView` qmuidemo/src/main/res/layout/fragment_about.xml:54 |
| `QMUIWindowInsetLayout` | 自定义 View |  | 26 | `<com.qmuiteam.qmui.widget.QMUIWindowInsetLayout xmlns:android="http://schemas.android.com/apk/res/android"` qmuidemo/src/main/res/layout/activity_arch_test.xml:18 |
| `QMUIQQFaceView` | 自定义 View |  | 25 | `<com.qmuiteam.qmui.qqface.QMUIQQFaceView` qmuidemo/src/main/res/layout/fragment_qqface_layout.xml:97 |
| `QMUIPriorityLinearLayout` | 自定义 View |  | 18 | `<com.qmuiteam.qmui.layout.QMUIPriorityLinearLayout` qmuidemo/src/main/res/layout/fragment_priority_linear_layout.xml:45 |
| `QMUIWindowInsetLayout2` | 自定义 View |  | 16 | `<com.qmuiteam.qmui.widget.QMUIWindowInsetLayout2 xmlns:android="http://schemas.android.com/apk/res/android"` qmuidemo/src/main/res/layout/fragment_animation_listview.xml:18 |
| `QMUISpanTouchFixTextView` | 自定义 View |  | 15 | `<com.qmuiteam.qmui.widget.textview.QMUISpanTouchFixTextView` qmui/src/main/res/layout/qmui_empty_view.xml:36 |
| `QMUITipDialog` | 自定义 View | dialog | 14 | `@Widget(widgetClass = QMUITipDialog.class, iconRes = R.mipmap.icon_grid_tip_dialog)` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDTipDialogFragment.java:44 |
| `QMUIPopup` | 自定义 View | popup | 14 | `return new QMUIPopup(context,` qmui/src/main/java/com/qmuiteam/qmui/widget/popup/QMUIPopups.java:30 |
| `QMUIPullRefreshLayout` | 自定义 View |  | 11 | `<com.qmuiteam.qmui.widget.pullRefreshLayout.QMUIPullRefreshLayout` qmuidemo/src/main/res/layout/fragment_continuous_nested_scroll.xml:27 |
| `QMUIFrameLayout` | 自定义 View |  | 11 | `<com.qmuiteam.qmui.layout.QMUIFrameLayout` qmuidemo/src/main/res/layout/fragment_button.xml:39 |
| `QMUIDialog` | 自定义 View | dialog | 10 | `mDialog = new QMUIDialog(mContext, style);` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialogBuilder.java:343 |
| `QMUIRadiusImageView2` | 自定义 View | image | 7 | `<com.qmuiteam.qmui.widget.QMUIRadiusImageView2` qmuidemo/src/main/res/layout/fragment_radius_imageview2.xml:42 |
| `QMUIRadiusImageView` | 自定义 View | image | 7 | `<com.qmuiteam.qmui.widget.QMUIRadiusImageView` qmuidemo/src/main/res/layout/fragment_radius_imageview.xml:42 |
| `QMUILinearLayout` | 自定义 View |  | 7 | `<com.qmuiteam.qmui.layout.QMUILinearLayout style="@style/button_wrapper_style"` qmuidemo/src/main/res/layout/fragment_button.xml:112 |
| `QMUIWebView` | 自定义 View |  | 7 | `private QMUIWebView.OnScrollChangeListener mOnScrollChangeListener;` qmui/src/main/java/com/qmuiteam/qmui/widget/webview/QMUIWebViewContainer.java:34 |
| `QMUIProgressBar` | 自定义 View | loading | 6 | `<com.qmuiteam.qmui.widget.QMUIProgressBar` qmuidemo/src/main/res/layout/fragment_progressbar.xml:38 |
| `QMUITopBar` | 自定义 View | navbar | 6 | `<com.qmuiteam.qmui.widget.QMUITopBar` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:54 |
| `QMUIContinuousNestedBottomRecyclerView` | 自定义 View | list | 6 | `QMUIContinuousNestedBottomRecyclerView recyclerView = new QMUIContinuousNestedBottomRecyclerView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousBottomView.java:96 |
| `QMUIContinuousNestedScrollLayout` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.nestedScroll.QMUIContinuousNestedScrollLayout` qmuidemo/src/main/res/layout/fragment_continuous_nested_scroll.xml:31 |
| `PhotoContentWithRectTransition` | Composable |  | 1 | `PhotoContentWithRectTransition(` photo/src/main/java/com/qmuiteam/photo/compose/GesturePhoto.kt:396 |
| `QMUIDialogMsg` | Composable | dialog | 1 | `QMUIDialogMsg(modal,` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/QDDialogFragment.kt:78 |
| `QMUIDialogMarkList` | Composable | dialog | 1 | `QMUIDialogMarkList(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/QDDialogFragment.kt:134 |
| `QMUIPhotoBucketItem` | Composable | cell | 1 | `QMUIPhotoBucketItem(it, it.id == currentId, onBucketClick)` photo/src/main/java/com/qmuiteam/photo/compose/picker/Buckets.kt:87 |
| `QMUICollapsingTopBarLayout` | 自定义 View |  | 3 | `<com.qmuiteam.qmui.widget.QMUICollapsingTopBarLayout` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:31 |
| `QMUIObservableScrollView` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/widget/QMUIObservableScrollView.java |
| `QMUIWrapContentListView` | 自定义 View | list | 1 | `ListView listView = new QMUIWrapContentListView(context, maxHeight);` qmui/src/main/java/com/qmuiteam/qmui/widget/popup/QMUIPopups.java:58 |
| `PressWithAlphaBox` | Composable |  | 4 | `PressWithAlphaBox(` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:111 |
| `QMUIAlphaRelativeLayout` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/alpha/QMUIAlphaRelativeLayout.java |
| `QMUIButton` | 自定义 View | button | 2 | `QMUIButton button = new QMUIButton(context);` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBottomSheetBaseBuilder.java:198 |
| `PhotoBackgroundWithTransition` | Composable |  | 1 | `PhotoBackgroundWithTransition(backgroundTargetAlpha, transition, transitionDurationMs) {` photo/src/main/java/com/qmuiteam/photo/compose/GesturePhoto.kt:186 |
| `QMUIPhotoPickRadio` | Composable | switch | 1 | `QMUIPhotoPickRadio(isOriginOpen)` photo/src/main/java/com/qmuiteam/photo/compose/picker/Common.kt:142 |
| `QMUIContinuousNestedBottomDelegateLayout` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/nestedScroll/QMUIContinuousNestedBottomDelegateLayout.java |
| `QMUITabView` | 自定义 View | tabs | 2 | `public class QMUITabAdapter extends QMUIItemViewsAdapter<QMUITab, QMUITabView> implements QMUITabView.Callback {` qmui/src/main/java/com/qmuiteam/qmui/widget/tab/QMUITabAdapter.java:23 |
| `QMUIDialogTitle` | Composable | dialog | 1 | `QMUIDialogTitle(title)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:109 |
| `PhotoContentWithAlphaTransition` | Composable |  | 1 | `PhotoContentWithAlphaTransition(` photo/src/main/java/com/qmuiteam/photo/compose/GesturePhoto.kt:368 |
| `QMUIBottomSheetGridItemView` | 自定义 View | popup | 1 | `QMUIBottomSheetGridItemView itemView = new QMUIBottomSheetGridItemView(bottomSheet.getContext());` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBottomSheet.java:511 |
| `QMUIDialogMutiCheckList` | Composable | dialog | 1 | `QMUIDialogMutiCheckList(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/QDDialogFragment.kt:169 |
| … | | | | 共 161 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| button | Button（组件库，用 5 次）<br>QMUIRoundButton（自定义 View，用 32 次）<br>CommonImageButton（Composable，用 6 次）<br>QMUIButton（自定义 View，用 2 次）<br>CommonTextButton（Composable，用 2 次）<br>OriginOpenButton（Composable，用 2 次）<br>CommonButton（Composable，用 2 次）<br>QMUIAlphaImageButton（自定义 View，用 2 次）<br>PickCurrentCheckButton（Composable，用 1 次）<br>QMUIAlphaButton（自定义 View，用 0 次） | `Button` |
| list | LazyColumn（组件库，用 9 次）<br>QMUIContinuousNestedBottomRecyclerView（自定义 View，用 6 次）<br>QMUIContinuousNestedTopRecyclerView（自定义 View，用 4 次）<br>QMUIAnimationListView（自定义 View，用 4 次）<br>RecyclerView（组件库，用 6 次）<br>LazyRow（组件库，用 1 次）<br>ListView（组件库，用 2 次）<br>QMUIWrapContentListView（自定义 View，用 1 次） | `LazyColumn` |
| loading | QMUILoadingView（自定义 View，用 5 次）<br>QMUIProgressBar（自定义 View，用 6 次）<br>QMUIPhotoLoading（Composable，用 3 次）<br>QDLoadingItemView（自定义 View，用 1 次）<br>ProgressBar（组件库，用 1 次） | `QMUILoadingView` |
| dialog | QMUIDialog（自定义 View，用 10 次）<br>QMUITipDialog（自定义 View，用 14 次）<br>QMUIBaseDialog（自定义 View，用 0 次）<br>QMUIModal（Composable，用 0 次） | `QMUIDialog` |
| cell | QMUICommonListItemView（自定义 View，用 56 次）<br>QMUIItem（Composable，用 10 次）<br>QMUIBottomSheetListItemView（自定义 View，用 1 次） | `QMUICommonListItemView` |
| popup | QMUIPopup（自定义 View，用 14 次）<br>QMUIBottomSheet（自定义 View，用 38 次）<br>QMUIFullScreenPopup（自定义 View，用 3 次） | `QMUIPopup` |
| image | ImageView（组件库，用 11 次）<br>QMUIRadiusImageView2（自定义 View，用 7 次）<br>QMUIRadiusImageView（自定义 View，用 7 次） | `ImageView` |
| tabs | QMUITabSegment（组件库，用 4 次）<br>QMUIBasicTabSegment（自定义 View，用 4 次）<br>QMUITabView（自定义 View，用 2 次） | `QMUITabSegment` |
| navbar | QMUITopBar（自定义 View，用 6 次）<br>QMUIGroupListSectionHeaderFooterView（自定义 View，用 2 次）<br>QDSectionHeaderView（自定义 View，用 1 次） | `QMUITopBar` |
| switch | RadioButton（组件库，用 5 次）<br>QMUICheckBox（Composable，用 1 次）<br>QMUIPhotoPickCheckBox（Composable，用 1 次） | `RadioButton` |
| carousel | HorizontalPager（组件库，用 1 次）<br>SwipeBackgroundView（自定义 View，用 2 次）<br>ViewPager（组件库，用 1 次） | `HorizontalPager` |
| input | EmojiEditText（自定义 View，用 3 次）<br>EditText（组件库，用 2 次） | `EmojiEditText` |
| icon | QMUIMarkIcon（Composable，用 1 次）<br>QMUIChevronIcon（Composable，用 0 次） | `QMUIMarkIcon` |

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
