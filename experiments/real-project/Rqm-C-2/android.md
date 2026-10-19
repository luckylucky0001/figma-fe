
## 页面复用检查：qmuidemo/src/main/java/com/qmuiteam/qmuidemo/manager/QDDataManager.java,qmuidemo/src/main/res/values/attr.xml,qmuidemo/src/main/res/values/colors.xml,qmuidemo/src/main/res/values/dimens.xml,qmuidemo/src/main/res/values/strings.xml,qmuidemo/src/main/res/values/styles.xml,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDOrderDetailFragment.kt,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/view/QDActionBarLayout.kt,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/view/QDCountDownView.kt,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/view/QDNoticeBarView.kt,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/view/QDPriceView.kt,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/view/QDProductCardView.kt,qmuidemo/src/main/java/com/qmuiteam/qmuidemo/view/QDTagView.kt,qmuidemo/src/main/res/layout/fragment_order_detail.xml,qmuidemo/src/main/res/layout/qd_product_card_view.xml  ✅ PASS

- 复用的现成组件（8）：QMUIButton（自定义 View）、QMUICommonListItemView（自定义 View）、QMUIGroupListView（自定义 View）、QMUILinearLayout（自定义 View）、QMUIRadiusImageView（自定义 View）、QMUITopBarLayout（自定义 View）、QMUIWindowInsetLayout（自定义 View）、TextView（组件库）
- 自造轮子（0）
- 有组件却没用（0）
- 疑似手写已有组件（0）
- 页面私有调色板：页面里自己定义了 8 个颜色变量（应使用项目 / 组件库已有的 token）
- 写死颜色（0）

---

# 存量复用体检：/home/user/exp/Rqm-C-2（android）

技术栈：XML View / Jetpack Compose / Material Components / QMUI Android · 扫描文件 573

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 98%（953 / 972） |
| 重复概念（同一 UI 概念有多个实现） | 12 |
| 重复结构（可抽成组件的复制粘贴结构） | 1 组 |
| 硬编码颜色 | 11 处（其中 1 处已有同值 token） |
| 硬编码尺寸 | 361 处 |

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
| `ViewPager` | 组件库 | carousel | 1 | `<androidx.viewpager.widget.ViewPager` qmuidemo/src/main/res/layout/fragment_tab_viewpager_layout.xml:39 |
| `QMUISeekBar` | 组件库 |  | 1 | `<com.qmuiteam.qmui.widget.QMUISeekBar` qmuidemo/src/main/res/layout/fragment_slider.xml:44 |
| `ViewPager2` | 组件库 |  | 1 | `<androidx.viewpager2.widget.ViewPager2` qmuidemo/src/main/res/layout/fragment_tab_viewpager2_layout.xml:39 |
| `QMUIItem` | Composable | cell | 10 | `QMUIItem(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:152 |
| `QMUIPhotoThumbnailWithViewer` | Composable |  | 9 | `QMUIPhotoThumbnailWithViewer(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDPhotoFragment.kt:102 |
| `QMUIDialogAction` | Composable | dialog | 8 | `QMUIDialogAction(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:90 |
| `QMUITopBar` | Composable | navbar | 7 | `QMUITopBar(` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:193 |
| `CommonImageButton` | Composable | button | 6 | `CommonImageButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Edit.kt:338 |
| `QMUIDialog` | Composable | dialog | 5 | `QMUIDialog(modal, horEdge, verEdge, widthLimit, radius, background, content)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:351 |
| `QMUIGesturePhoto` | Composable |  | 4 | `QMUIGesturePhoto(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoViewerActivity.kt:185 |
| `TextButton` | Composable | button | 4 | `TextButton("加粗"){` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEditorFragment.kt:57 |
| `PressWithAlphaBox` | Composable |  | 4 | `PressWithAlphaBox(` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:111 |
| `OnePx` | Composable |  | 4 | `separatorHeight: Dp = OnePx(),` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:179 |
| `QMUITopBarWithLazyScrollState` | Composable |  | 4 | `QMUITopBarWithLazyScrollState(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:358 |
| `BlankBox` | Composable |  | 3 | `BlankBox()` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:80 |
| `QMUIDialogList` | Composable | dialog | 3 | `QMUIDialogList(modal, maxHeight, state, contentPadding) {` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:150 |
| `QMUIPhotoLoading` | Composable | loading | 3 | `QMUIPhotoLoading(lineColor = QMUILocalPickerConfig.current.loadingColor)` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:577 |
| `QMUIPhotoClipper` | Composable |  | 2 | `QMUIPhotoClipper(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoClipActivity.kt:92 |
| `LongImageContent` | Composable | image | 2 | `LongImageContent(request, onSuccess, onError, openBlankColor)` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:107 |
| `OriginOpenButton` | Composable | button | 2 | `OriginOpenButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:207 |
| `PhotoTransformContent` | Composable |  | 2 | `PhotoTransformContent(` photo/src/main/java/com/qmuiteam/photo/compose/GesturePhoto.kt:376 |
| `LongImage` | Composable | image | 2 | `LongImage(onSuccess, onError, openBlankColor)` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:50 |
| `CommonButton` | Composable | button | 2 | `CommonButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Edit.kt:815 |
| `QMUIToast` | Composable | toast | 2 | `QMUIToast(modal, radius, background, content)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIToast.kt:116 |
| `CommonTextButton` | Composable | button | 2 | `CommonTextButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:199 |
| `RowImages` | Composable | image | 2 | `RowImages(images, renderInfo, config, maxWidth, 2, 0, onClick)` photo/src/main/java/com/qmuiteam/photo/compose/Thumbnail.kt:255 |
| `QMUIBitmapRegionItem` | Composable | cell | 2 | `QMUIBitmapRegionItem(image, maxWidth, heightDp)` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:236 |
| `QMUIDialogActions` | Composable | dialog | 2 | `QMUIDialogActions(modal, actions)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:111 |
| `QMUIPhotoPickerPreviewPickedItems` | Composable | cell | 1 | `QMUIPhotoPickerPreviewPickedItems(data, pickedItems, data[pagerState.currentPage].model.id) {` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:528 |
| `QMUIContinuousNestedScrollLayout` | 自定义 View |  | 5 | `<com.qmuiteam.qmui.nestedScroll.QMUIContinuousNestedScrollLayout` qmuidemo/src/main/res/layout/fragment_continuous_nested_scroll.xml:31 |
| `QMUIDialogMarkList` | Composable | dialog | 1 | `QMUIDialogMarkList(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/QDDialogFragment.kt:134 |
| `QMUISpanTouchFixTextView` | 自定义 View |  | 15 | `<com.qmuiteam.qmui.widget.textview.QMUISpanTouchFixTextView` qmui/src/main/res/layout/qmui_empty_view.xml:36 |
| `QMUIPullLayout` | 自定义 View |  | 49 | `<com.qmuiteam.qmui.widget.pullLayout.QMUIPullLayout` qmuidemo/src/main/res/layout/fragment_pull_horizontal_test_layout.xml:32 |
| `QMUIPhotoPickerGridCellMask` | Composable |  | 1 | `QMUIPhotoPickerGridCellMask(pickedIndex)` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:151 |
| `QMUITipDialog` | 自定义 View | dialog | 14 | `@Widget(widgetClass = QMUITipDialog.class, iconRes = R.mipmap.icon_grid_tip_dialog)` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/components/QDTipDialogFragment.java:44 |
| `QMUIChevronIcon` | Composable | icon | 0 | compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUIIcon.kt |
| `TestImageCompress` | Composable | image | 1 | `//                    TestImageCompress()` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDPhotoFragment.kt:93 |
| `QMUIPhotoPickerEditPhotoContent` | Composable |  | 1 | `QMUIPhotoPickerEditPhotoContent(data) {` photo/src/main/java/com/qmuiteam/photo/compose/picker/Edit.kt:124 |
| `QMUIPullRefreshLayout` | 自定义 View |  | 11 | `<com.qmuiteam.qmui.widget.pullRefreshLayout.QMUIPullRefreshLayout` qmuidemo/src/main/res/layout/fragment_continuous_nested_scroll.xml:27 |
| `QMUIDialogMenuItemView` | 自定义 View | dialog | 5 | `itemView.setListener(new QMUIDialogMenuItemView.MenuItemViewListener() {` qmui/src/main/java/com/qmuiteam/qmui/widget/dialog/QMUIDialog.java:494 |
| `QMUIButton` | 自定义 View | button | 5 | `<com.qmuiteam.qmui.layout.QMUIButton` qmuidemo/src/main/res/layout/fragment_order_detail.xml:149 |
| `QMUICommonListItemView` | 自定义 View | cell | 53 | `QMUICommonListItemView itemView = new QMUICommonListItemView(getContext());` qmui/src/main/java/com/qmuiteam/qmui/widget/grouplist/QMUIGroupListView.java:83 |
| `QMUITopBarLayout` | 自定义 View | navbar | 46 | `<com.qmuiteam.qmui.widget.QMUITopBarLayout` qmuidemo/src/main/res/layout/activity_arch_test.xml:23 |
| `QMUIPopup` | 自定义 View | popup | 14 | `return new QMUIPopup(context,` qmui/src/main/java/com/qmuiteam/qmui/widget/popup/QMUIPopups.java:30 |
| `QMUIPhotoThumbnail` | Composable |  | 1 | `QMUIPhotoThumbnail(images, config) { list, index ->` photo/src/main/java/com/qmuiteam/photo/compose/Thumbnail.kt:93 |
| `QMUICheckBox` | Composable | switch | 1 | `QMUICheckBox(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Common.kt:170 |
| `QMUIPhotoBucketItem` | Composable | cell | 1 | `QMUIPhotoBucketItem(it, it.id == currentId, onBucketClick)` photo/src/main/java/com/qmuiteam/photo/compose/picker/Buckets.kt:87 |
| `QMUIPriorityLinearLayout` | 自定义 View |  | 18 | `<com.qmuiteam.qmui.layout.QMUIPriorityLinearLayout` qmuidemo/src/main/res/layout/fragment_priority_linear_layout.xml:45 |
| `QMUIModal` | Composable | dialog | 0 | compose/src/main/java/com/qmuiteam/compose/modal/QMUIModal.kt |
| `QMUIDefaultPhotoConfigProvider` | Composable |  | 0 | photo/src/main/java/com/qmuiteam/photo/compose/PhotoConfig.kt |
| `QMUIPhotoPickerBucketToggleArrow` | Composable | switch | 1 | `QMUIPhotoPickerBucketToggleArrow(iconBgColor, iconColor, isFocusFlow)` photo/src/main/java/com/qmuiteam/photo/compose/picker/TopBarItem.kt:68 |
| `SwipeBackLayout` | 自定义 View |  | 78 | `import static com.qmuiteam.qmui.arch.SwipeBackLayout.DRAG_DIRECTION_BOTTOM_TO_TOP;` arch/src/main/java/com/qmuiteam/qmui/arch/QMUIActivity.java:37 |
| `QDEditor` | Composable |  | 1 | `QDEditor()` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEditorFragment.kt:100 |
| `QMUIProgressBar` | 自定义 View | loading | 6 | `<com.qmuiteam.qmui.widget.QMUIProgressBar` qmuidemo/src/main/res/layout/fragment_progressbar.xml:38 |
| … | | | | 共 168 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| button | QMUIRoundButton（自定义 View，用 29 次）<br>CommonImageButton（Composable，用 6 次）<br>QMUIButton（自定义 View，用 5 次）<br>TextButton（Composable，用 4 次）<br>OriginOpenButton（Composable，用 2 次）<br>CommonButton（Composable，用 2 次）<br>CommonTextButton（Composable，用 2 次）<br>QMUIAlphaImageButton（自定义 View，用 2 次）<br>PickCurrentCheckButton（Composable，用 1 次）<br>QMUIAlphaButton（自定义 View，用 0 次） | `QMUIRoundButton` |
| cell | QMUICommonListItemView（自定义 View，用 53 次）<br>QMUIItem（Composable，用 10 次）<br>QMUIBitmapRegionItem（Composable，用 2 次）<br>QMUIPhotoPickerPreviewPickedItems（Composable，用 1 次）<br>QMUIPhotoBucketItem（Composable，用 1 次） | `QMUICommonListItemView` |
| dialog | QMUITipDialog（自定义 View，用 14 次）<br>QMUIDialog（自定义 View，用 10 次）<br>QMUIDialog（Composable，用 5 次）<br>QMUIModal（Composable，用 0 次）<br>QMUIBaseDialog（自定义 View，用 0 次） | `QMUITipDialog` |
| image | QMUIRadiusImageView（自定义 View，用 8 次）<br>QMUIRadiusImageView2（自定义 View，用 7 次）<br>LongImage（Composable，用 2 次）<br>RowImages（Composable，用 2 次）<br>TestImageCompress（Composable，用 1 次） | `QMUIRadiusImageView` |
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
- qmuidemo/src/main/res/layout/fragment_layout.xml:114 `<Button` → 原生 Button：项目里已有 QMUIAlphaButton（组件库）
- qmuidemo/src/main/res/layout/fragment_layout.xml:125 `<Button` → 原生 Button：项目里已有 QMUIRoundButton（组件库）
- qmuidemo/src/main/res/layout/fragment_layout.xml:148 `<Button` → 原生 Button：项目里已有 QMUIAlphaButton（组件库）
- qmuidemo/src/main/res/layout/fragment_layout.xml:158 `<Button` → 原生 Button：项目里已有 QMUIRoundButton（组件库）
- qmuidemo/src/main/res/layout/fragment_layout.xml:167 `<Button` → 原生 Button：项目里已有 QMUIAlphaButton（组件库）
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
| #222222 | 1 | — | photo/src/main/java/com/qmuiteam/photo/compose/picker/Config.kt |
| #00A8E1 | 1 | `@color/app_color_blue` | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #CCCCCC | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #777777 | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
| #666666 | 1 | — | compose-core/src/main/java/com/qmuiteam/compose/core/ui/DefaultConfig.kt |
