
## 页面复用检查：qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDOrderDetailFragment.kt,qmuidemo/src/main/res/layout/fragment_order_detail.xml,qmuidemo/src/main/res/values/order_detail.xml  ❌ FAIL

- 复用的现成组件（8）：QMUILinearLayout（自定义 View）、QMUIRadiusImageView2（自定义 View）、QMUIRoundButton（自定义 View）、QMUIRoundFrameLayout（自定义 View）、QMUIRoundLinearLayout（自定义 View）、QMUITopBarLayout（自定义 View）、QMUIWindowInsetLayout（自定义 View）、TextView（组件库）
- 自造轮子（0）
- 有组件却没用（1）
  - · qmuidemo/src/main/res/layout/fragment_order_detail.xml:165 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）
- 疑似手写已有组件（1）
  - ✗ qmuidemo/src/main/res/layout/fragment_order_detail.xml:132 页面手写了「cell」（.order_detail_address_cell），已有 QMUICommonListItemView（自定义 View qmui/src/main/java/com/qmuiteam/qmui/widget/grouplist/QMUICommonListItemView.java）（项目里已用 48 次）
- 页面私有调色板：页面里自己定义了 11 个颜色变量（应使用项目 / 组件库已有的 token）
- 写死颜色（4）
  - · qmuidemo/src/main/res/values/order_detail.xml:5 #FFFFFF：新定义的颜色 order_detail_card_bg 与已有 token @color/qmui_config_color_white 同值
  - · qmuidemo/src/main/res/values/order_detail.xml:43 #E5E6EB：写死颜色
  - · qmuidemo/src/main/res/values/order_detail.xml:99 #E5E6EB：写死颜色
  - · qmuidemo/src/main/res/values/order_detail.xml:107 #FF4D1A：写死颜色

---

# 存量复用体检：/home/user/exp/Rqm-A-2（android）

技术栈：XML View / Jetpack Compose / Material Components / QMUI Android · 扫描文件 566

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 98%（957 / 977） |
| 重复概念（同一 UI 概念有多个实现） | 12 |
| 重复结构（可抽成组件的复制粘贴结构） | 1 组 |
| 硬编码颜色 | 11 处（其中 1 处已有同值 token） |
| 硬编码尺寸 | 403 处 |

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Android SDK widgets | 148 | TextView×130、RecyclerView×6、RadioButton×5、SeekBar×3、ListView×2、ViewPager2×1、ViewPager×1 |
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
| `TextView` | 组件库 |  | 130 | `<TextView` qmui/src/main/res/layout/qmui_common_list_item.xml:45 |
| `RecyclerView` | 组件库 | list | 6 | `<androidx.recyclerview.widget.RecyclerView` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:65 |
| `RadioButton` | 组件库 | switch | 5 | `<RadioButton` qmuidemo/src/main/res/layout/fragment_layout.xml:189 |
| `QMUITabSegment` | 组件库 | tabs | 4 | `<com.qmuiteam.qmui.widget.tab.QMUITabSegment` qmuidemo/src/main/res/layout/fragment_fsw_viewpager.xml:30 |
| `SeekBar` | 组件库 |  | 3 | `<SeekBar` qmuidemo/src/main/res/layout/fragment_colorhelper.xml:124 |
| `ListView` | 组件库 | list | 2 | `<ListView` qmuidemo/src/main/res/layout/fragment_listview.xml:25 |
| `QMUISeekBar` | 组件库 |  | 1 | `<com.qmuiteam.qmui.widget.QMUISeekBar` qmuidemo/src/main/res/layout/fragment_slider.xml:44 |
| `ViewPager2` | 组件库 |  | 1 | `<androidx.viewpager2.widget.ViewPager2` qmuidemo/src/main/res/layout/fragment_tab_viewpager2_layout.xml:39 |
| `ViewPager` | 组件库 | carousel | 1 | `<androidx.viewpager.widget.ViewPager` qmuidemo/src/main/res/layout/fragment_tab_viewpager_layout.xml:39 |
| `QMUIItem` | Composable | cell | 10 | `QMUIItem(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:152 |
| `QMUIPhotoThumbnailWithViewer` | Composable |  | 9 | `QMUIPhotoThumbnailWithViewer(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDPhotoFragment.kt:102 |
| `QMUIDialogAction` | Composable | dialog | 8 | `QMUIDialogAction(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:90 |
| `QMUITopBar` | Composable | navbar | 7 | `QMUITopBar(` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:193 |
| `CommonImageButton` | Composable | button | 6 | `CommonImageButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Edit.kt:338 |
| `QMUIDialog` | Composable | dialog | 5 | `QMUIDialog(modal, horEdge, verEdge, widthLimit, radius, background, content)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:351 |
| `TextButton` | Composable | button | 4 | `TextButton("加粗"){` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEditorFragment.kt:57 |
| `QMUITopBarWithLazyScrollState` | Composable |  | 4 | `QMUITopBarWithLazyScrollState(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:358 |
| `PressWithAlphaBox` | Composable |  | 4 | `PressWithAlphaBox(` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:111 |
| `OnePx` | Composable |  | 4 | `separatorHeight: Dp = OnePx(),` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:179 |
| `QMUIGesturePhoto` | Composable |  | 4 | `QMUIGesturePhoto(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoViewerActivity.kt:185 |
| `QMUIPhotoLoading` | Composable | loading | 3 | `QMUIPhotoLoading(lineColor = QMUILocalPickerConfig.current.loadingColor)` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:577 |
| `BlankBox` | Composable |  | 3 | `BlankBox()` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:80 |
| `QMUIDialogList` | Composable | dialog | 3 | `QMUIDialogList(modal, maxHeight, state, contentPadding) {` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:150 |
| `CommonTextButton` | Composable | button | 2 | `CommonTextButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:199 |
| `QMUIPhotoClipper` | Composable |  | 2 | `QMUIPhotoClipper(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoClipActivity.kt:92 |
| `QMUIDialogActions` | Composable | dialog | 2 | `QMUIDialogActions(modal, actions)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:111 |
| `LongImage` | Composable | image | 2 | `LongImage(onSuccess, onError, openBlankColor)` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:50 |
| `QMUIBitmapRegionItem` | Composable | cell | 2 | `QMUIBitmapRegionItem(image, maxWidth, heightDp)` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:236 |
| `LongImageContent` | Composable | image | 2 | `LongImageContent(request, onSuccess, onError, openBlankColor)` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:107 |
| `RowImages` | Composable | image | 2 | `RowImages(images, renderInfo, config, maxWidth, 2, 0, onClick)` photo/src/main/java/com/qmuiteam/photo/compose/Thumbnail.kt:255 |
| `QMUIToast` | Composable | toast | 2 | `QMUIToast(modal, radius, background, content)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIToast.kt:116 |
| `OriginOpenButton` | Composable | button | 2 | `OriginOpenButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:207 |
| `CommonButton` | Composable | button | 2 | `CommonButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Edit.kt:815 |
| `PhotoTransformContent` | Composable |  | 2 | `PhotoTransformContent(` photo/src/main/java/com/qmuiteam/photo/compose/GesturePhoto.kt:376 |
| `PickCurrentCheckButton` | Composable | button | 1 | `PickCurrentCheckButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:210 |
| `QMUIDialogMsg` | Composable | dialog | 1 | `QMUIDialogMsg(modal,` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/QDDialogFragment.kt:78 |
| `QMUIWebView` | 自定义 View |  | 7 | `private QMUIWebView.OnScrollChangeListener mOnScrollChangeListener;` qmui/src/main/java/com/qmuiteam/qmui/widget/webview/QMUIWebViewContainer.java:34 |
| `QMUIFloatLayout` | 自定义 View |  | 4 | `<com.qmuiteam.qmui.widget.QMUIFloatLayout` qmuidemo/src/main/res/layout/fragment_floatlayout.xml:37 |
| `QMUIGroupListView` | 自定义 View | card | 25 | `<com.qmuiteam.qmui.widget.grouplist.QMUIGroupListView` qmuidemo/src/main/res/layout/fragment_about.xml:54 |
| `QMUIAlphaTextView` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/alpha/QMUIAlphaTextView.java |
| `QMUIContinuousNestedBottomRecyclerView` | 自定义 View |  | 6 | `QMUIContinuousNestedBottomRecyclerView recyclerView = new QMUIContinuousNestedBottomRecyclerView(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousBottomView.java:96 |
| `QMUIProgressBar` | 自定义 View | loading | 6 | `<com.qmuiteam.qmui.widget.QMUIProgressBar` qmuidemo/src/main/res/layout/fragment_progressbar.xml:38 |
| `QMUIPullLayout` | 自定义 View |  | 49 | `<com.qmuiteam.qmui.widget.pullLayout.QMUIPullLayout` qmuidemo/src/main/res/layout/fragment_pull_horizontal_test_layout.xml:32 |
| `QMUIContinuousNestedTopLinearLayout` | 自定义 View |  | 1 | `mTopLinearLayout = new QMUIContinuousNestedTopLinearLayout(getContext());` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDContinuousNestedScroll6Fragment.java:55 |
| `QDQQFaceBasePagerView` | 自定义 View |  | 0 | qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/qqface/pageView/QDQQFaceBasePagerView.java |
| `QMUITopBarContent` | Composable |  | 1 | `QMUITopBarContent(` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:244 |
| `QMUIPhotoPickerEditPhotoContent` | Composable |  | 1 | `QMUIPhotoPickerEditPhotoContent(data) {` photo/src/main/java/com/qmuiteam/photo/compose/picker/Edit.kt:124 |
| `EmojiEditText` | 自定义 View |  | 3 | `val editText = EmojiEditText(context).apply {` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEmojiInputFragment.kt:38 |
| `QMUILinearLayout` | 自定义 View |  | 7 | `<com.qmuiteam.qmui.layout.QMUILinearLayout style="@style/button_wrapper_style"` qmuidemo/src/main/res/layout/fragment_button.xml:112 |
| `EmojiLayout` | 自定义 View |  | 2 | `return EmojiLayout(requireContext())` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEmojiInputFragment.kt:29 |
| `QMUIAlphaConstraintLayout` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/alpha/QMUIAlphaConstraintLayout.java |
| `QMUIPhotoPickCheckBox` | Composable |  | 1 | `QMUIPhotoPickCheckBox(pickedIndex)` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:163 |
| `QMUIObservableScrollView` | 自定义 View |  | 0 | qmui/src/main/java/com/qmuiteam/qmui/widget/QMUIObservableScrollView.java |
| `QMUICheckBox` | Composable | switch | 1 | `QMUICheckBox(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Common.kt:170 |
| `QMUIBaseDialog` | 自定义 View | dialog | 0 | qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIBaseDialog.java |
| `QMUIEditor` | Composable |  | 1 | `QMUIEditor(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEditorFragment.kt:42 |
| `PhotoContentWithRectTransition` | Composable |  | 1 | `PhotoContentWithRectTransition(` photo/src/main/java/com/qmuiteam/photo/compose/GesturePhoto.kt:396 |
| `QMUIPhotoPickRadio` | Composable | switch | 1 | `QMUIPhotoPickRadio(isOriginOpen)` photo/src/main/java/com/qmuiteam/photo/compose/picker/Common.kt:142 |
| `QMUIBottomSheetRootLayout` | 自定义 View |  | 1 | `<com.qmuiteam.qmui.widget.dialog.QMUIBottomSheetRootLayout` qmui/src/main/res/layout/qmui_bottom_sheet_dialog.xml:38 |
| `QMUIDialog` | 自定义 View | dialog | 10 | `mDialog = new QMUIDialog(mContext, style);` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialogBuilder.java:343 |
| … | | | | 共 162 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| button | QMUIRoundButton（自定义 View，用 33 次）<br>CommonImageButton（Composable，用 6 次）<br>TextButton（Composable，用 4 次）<br>CommonTextButton（Composable，用 2 次）<br>OriginOpenButton（Composable，用 2 次）<br>CommonButton（Composable，用 2 次）<br>QMUIAlphaImageButton（自定义 View，用 2 次）<br>QMUIButton（自定义 View，用 2 次）<br>PickCurrentCheckButton（Composable，用 1 次）<br>QMUIAlphaButton（自定义 View，用 0 次） | `QMUIRoundButton` |
| cell | QMUICommonListItemView（自定义 View，用 48 次）<br>QMUIItem（Composable，用 10 次）<br>QMUIBitmapRegionItem（Composable，用 2 次）<br>QMUIPhotoPickerPreviewPickedItems（Composable，用 1 次）<br>QMUIPhotoBucketItem（Composable，用 1 次） | `QMUICommonListItemView` |
| dialog | QMUITipDialog（自定义 View，用 14 次）<br>QMUIDialog（自定义 View，用 10 次）<br>QMUIDialog（Composable，用 5 次）<br>QMUIBaseDialog（自定义 View，用 0 次）<br>QMUIModal（Composable，用 0 次） | `QMUITipDialog` |
| image | QMUIRadiusImageView2（自定义 View，用 8 次）<br>QMUIRadiusImageView（自定义 View，用 7 次）<br>LongImage（Composable，用 2 次）<br>RowImages（Composable，用 2 次）<br>TestImageCompress（Composable，用 1 次） | `QMUIRadiusImageView2` |
| navbar | QMUITopBar（Composable，用 7 次）<br>QMUITopBar（自定义 View，用 6 次）<br>QMUIGroupListSectionHeaderFooterView（自定义 View，用 2 次）<br>QDSectionHeaderView（自定义 View，用 1 次） | `QMUITopBar` |
| loading | QMUIProgressBar（自定义 View，用 6 次）<br>QMUILoadingView（自定义 View，用 5 次）<br>QMUIPhotoLoading（Composable，用 3 次）<br>QDLoadingItemView（自定义 View，用 1 次） | `QMUIProgressBar` |
| list | QMUIAnimationListView（自定义 View，用 4 次）<br>QMUIWrapContentListView（自定义 View，用 1 次）<br>RecyclerView（组件库，用 6 次）<br>ListView（组件库，用 2 次） | `QMUIAnimationListView` |
| switch | QMUICheckBox（Composable，用 1 次）<br>QMUIPhotoPickRadio（Composable，用 1 次）<br>QMUIPhotoPickerBucketToggleArrow（Composable，用 1 次）<br>RadioButton（组件库，用 5 次） | `QMUICheckBox` |
| popup | QMUIBottomSheet（自定义 View，用 38 次）<br>QMUIPopup（自定义 View，用 14 次）<br>QMUIFullScreenPopup（自定义 View，用 3 次） | `QMUIBottomSheet` |
| tabs | QMUITabSegment（组件库，用 4 次）<br>QMUIBasicTabSegment（自定义 View，用 4 次）<br>QMUITabView（自定义 View，用 2 次） | `QMUITabSegment` |
| carousel | SwipeBackgroundView（自定义 View，用 2 次）<br>ViewPager（组件库，用 1 次） | `SwipeBackgroundView` |
| icon | QMUIMarkIcon（Composable，用 1 次）<br>QMUIChevronIcon（Composable，用 0 次） | `QMUIMarkIcon` |

## 4. 有组件却没用（20 处）

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
- qmuidemo/src/main/res/layout/fragment_layout.xml:148 `<Button` → 原生 Button：项目里已有 QMUIRoundButton（组件库）
- qmuidemo/src/main/res/layout/fragment_layout.xml:158 `<Button` → 原生 Button：项目里已有 QMUIAlphaButton（组件库）
- qmuidemo/src/main/res/layout/fragment_layout.xml:167 `<Button` → 原生 Button：项目里已有 QMUIRoundButton（组件库）
- qmuidemo/src/main/res/layout/fragment_order_detail.xml:165 `<ImageView` → 原生 ImageView：项目里已有 ShapeableImageView（组件库）
- qmuidemo/src/main/res/layout/fragment_scheme.xml:47 `<EditText` → 原生 EditText：项目里已有 TextInputEditText（组件库）
- qmuidemo/src/main/res/layout/fragment_verticaltextview.xml:62 `<EditText` → 原生 EditText：项目里已有 TextInputEditText（组件库）
- qmuidemo/src/main/res/layout/fragment_webview_explorer.xml:37 `<ProgressBar` → 原生 ProgressBar：项目里已有 QMUILoadingView（组件库）
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
