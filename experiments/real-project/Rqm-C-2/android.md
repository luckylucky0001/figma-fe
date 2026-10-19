
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
| 重复概念（同一 UI 概念有多个实现） | 13 |
| 重复结构（可抽成组件的复制粘贴结构） | 0 组 |
| 硬编码颜色 | 11 处（其中 1 处已有同值 token） |
| 硬编码尺寸 | 361 处 |

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Android SDK widgets | 149 | TextView×112、ImageView×11、RecyclerView×6、RadioButton×5、Button×5、SeekBar×3、ListView×2、EditText×2 |
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
| `TextView` | 组件库 |  | 112 | `<TextView` qmui/src/main/res/layout/qmui_common_list_item.xml:45 |
| `ImageView` | 组件库 | image | 11 | `<ImageView` qmuidemo/src/main/res/layout/drawablehelper_createfromview.xml:25 |
| `LazyColumn` | 组件库 | list | 9 | `LazyColumn(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIBottomSheet.kt:40 |
| `RecyclerView` | 组件库 | list | 6 | `<androidx.recyclerview.widget.RecyclerView` qmuidemo/src/main/res/layout/fragment_collapsing_topbar_layout.xml:65 |
| `RadioButton` | 组件库 | switch | 5 | `<RadioButton` qmuidemo/src/main/res/layout/fragment_layout.xml:189 |
| `Button` | 组件库 | button | 5 | `<Button` qmuidemo/src/main/res/layout/fragment_layout.xml:114 |
| `QMUITabSegment` | 组件库 | tabs | 4 | `<com.qmuiteam.qmui.widget.tab.QMUITabSegment` qmuidemo/src/main/res/layout/fragment_fsw_viewpager.xml:30 |
| `SeekBar` | 组件库 |  | 3 | `<SeekBar` qmuidemo/src/main/res/layout/fragment_colorhelper.xml:124 |
| `EditText` | 组件库 | input | 2 | `<EditText` qmuidemo/src/main/res/layout/fragment_scheme.xml:47 |
| `ListView` | 组件库 | list | 2 | `<ListView` qmuidemo/src/main/res/layout/fragment_listview.xml:25 |
| `ViewPager` | 组件库 | carousel | 1 | `<androidx.viewpager.widget.ViewPager` qmuidemo/src/main/res/layout/fragment_tab_viewpager_layout.xml:39 |
| `HorizontalPager` | 组件库 | carousel | 1 | `HorizontalPager(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:39 |
| `ProgressBar` | 组件库 | loading | 1 | `<ProgressBar` qmuidemo/src/main/res/layout/fragment_webview_explorer.xml:37 |
| `LazyRow` | 组件库 | list | 1 | `LazyRow(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Preview.kt:116 |
| `ViewPager2` | 组件库 |  | 1 | `<androidx.viewpager2.widget.ViewPager2` qmuidemo/src/main/res/layout/fragment_tab_viewpager2_layout.xml:39 |
| `QMUISeekBar` | 组件库 |  | 1 | `<com.qmuiteam.qmui.widget.QMUISeekBar` qmuidemo/src/main/res/layout/fragment_slider.xml:44 |
| `QMUIItem` | Composable | cell | 10 | `QMUIItem(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:152 |
| `QMUIPhotoThumbnailWithViewer` | Composable |  | 9 | `QMUIPhotoThumbnailWithViewer(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDPhotoFragment.kt:102 |
| `QMUIDialogAction` | Composable | dialog | 8 | `QMUIDialogAction(` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:90 |
| `QMUITopBar` | Composable | navbar | 7 | `QMUITopBar(` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:193 |
| `CommonImageButton` | Composable | button | 6 | `CommonImageButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Edit.kt:338 |
| `QMUIDialog` | Composable | dialog | 5 | `QMUIDialog(modal, horEdge, verEdge, widthLimit, radius, background, content)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:351 |
| `QMUIGesturePhoto` | Composable |  | 4 | `QMUIGesturePhoto(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoViewerActivity.kt:185 |
| `PressWithAlphaBox` | Composable |  | 4 | `PressWithAlphaBox(` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:111 |
| `QMUITopBarWithLazyScrollState` | Composable |  | 4 | `QMUITopBarWithLazyScrollState(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:358 |
| `OnePx` | Composable |  | 4 | `separatorHeight: Dp = OnePx(),` compose-core/src/main/java/com/qmuiteam/compose/core/ui/QMUITopBar.kt:179 |
| `TextButton` | Composable | button | 4 | `TextButton("加粗"){` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEditorFragment.kt:57 |
| `BlankBox` | Composable |  | 3 | `BlankBox()` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:80 |
| `QMUIPhotoLoading` | Composable | loading | 3 | `QMUIPhotoLoading(lineColor = QMUILocalPickerConfig.current.loadingColor)` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:577 |
| `QMUIDialogList` | Composable | dialog | 3 | `QMUIDialogList(modal, maxHeight, state, contentPadding) {` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:150 |
| `QMUIPhotoClipper` | Composable |  | 2 | `QMUIPhotoClipper(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoClipActivity.kt:92 |
| `CommonTextButton` | Composable | button | 2 | `CommonTextButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:199 |
| `OriginOpenButton` | Composable | button | 2 | `OriginOpenButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:207 |
| `LongImage` | Composable | image | 2 | `LongImage(onSuccess, onError, openBlankColor)` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:50 |
| `CommonButton` | Composable | button | 2 | `CommonButton(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Edit.kt:815 |
| `QMUIDialogActions` | Composable | dialog | 2 | `QMUIDialogActions(modal, actions)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:111 |
| `QMUIToast` | Composable | toast | 2 | `QMUIToast(modal, radius, background, content)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIToast.kt:116 |
| `QMUIBitmapRegionItem` | Composable | cell | 2 | `QMUIBitmapRegionItem(image, maxWidth, heightDp)` photo-coil/src/main/java/com/qmuiteam/photo/coil/QMUICoilPhoto.kt:236 |
| `RowImages` | Composable | image | 2 | `RowImages(images, renderInfo, config, maxWidth, 2, 0, onClick)` photo/src/main/java/com/qmuiteam/photo/compose/Thumbnail.kt:255 |
| `TestImageCompress` | Composable | image | 1 | `//                    TestImageCompress()` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDPhotoFragment.kt:93 |
| `PickerResult` | Composable |  | 1 | `PickerResult()` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDPhotoFragment.kt:89 |
| `QMUICheckBox` | Composable | switch | 1 | `QMUICheckBox(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Common.kt:170 |
| `QMUIPhotoThumbnail` | Composable |  | 1 | `QMUIPhotoThumbnail(images, config) { list, index ->` photo/src/main/java/com/qmuiteam/photo/compose/Thumbnail.kt:93 |
| `QDEditor` | Composable |  | 1 | `QDEditor()` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEditorFragment.kt:100 |
| `QMUIDialogMsg` | Composable | dialog | 1 | `QMUIDialogMsg(modal,` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/QDDialogFragment.kt:78 |
| `PhotoBackground` | Composable |  | 1 | `PhotoBackground(alpha = it)` photo/src/main/java/com/qmuiteam/photo/compose/GesturePhoto.kt:187 |
| `PhotoContentWithAlphaTransition` | Composable |  | 1 | `PhotoContentWithAlphaTransition(` photo/src/main/java/com/qmuiteam/photo/compose/GesturePhoto.kt:368 |
| `QMUIPhotoPickerGridCellMask` | Composable |  | 1 | `QMUIPhotoPickerGridCellMask(pickedIndex)` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:151 |
| `QMUIPhotoPickerGrid` | Composable |  | 1 | `QMUIPhotoPickerGrid(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:372 |
| `QMUIPhotoPickerEdit` | Composable |  | 1 | `QMUIPhotoPickerEdit(onBackPressedDispatcher, scene.current) {` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:569 |
| `PhotoContentWithRectTransition` | Composable |  | 1 | `PhotoContentWithRectTransition(` photo/src/main/java/com/qmuiteam/photo/compose/GesturePhoto.kt:396 |
| `QMUIPhotoPickCheckBox` | Composable | switch | 1 | `QMUIPhotoPickCheckBox(pickedIndex)` photo/src/main/java/com/qmuiteam/photo/compose/picker/Grid.kt:163 |
| `QMUIPhotoPickerGridToolBar` | Composable | navbar | 1 | `QMUIPhotoPickerGridToolBar(` photo/src/main/java/com/qmuiteam/photo/activity/QMUIPhotoPickerActivity.kt:391 |
| `QMUIMarkIcon` | Composable | icon | 1 | `QMUIMarkIcon(` photo/src/main/java/com/qmuiteam/photo/compose/picker/Buckets.kt:161 |
| `QMUIDialogMarkList` | Composable | dialog | 1 | `QMUIDialogMarkList(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/QDDialogFragment.kt:134 |
| `QMUIDialogTitle` | Composable | dialog | 1 | `QMUIDialogTitle(title)` compose/src/main/java/com/qmuiteam/compose/modal/QMUIDialog.kt:109 |
| `QMUIPhotoPickRadio` | Composable | switch | 1 | `QMUIPhotoPickRadio(isOriginOpen)` photo/src/main/java/com/qmuiteam/photo/compose/picker/Common.kt:142 |
| `QMUIDialogMutiCheckList` | Composable | dialog | 1 | `QMUIDialogMutiCheckList(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/QDDialogFragment.kt:169 |
| `QMUIBottomSheetList` | Composable | popup | 1 | `QMUIBottomSheetList(it) {` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/QDDialogFragment.kt:217 |
| `QMUIEditor` | Composable |  | 1 | `QMUIEditor(` qmuidemo/src/main/java/com/qmuiteam/qmuidemo/fragment/lab/QDEditorFragment.kt:42 |
| … | | | | 共 167 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| button | QMUIButton（自定义 View，用 5 次）<br>Button（组件库，用 5 次）<br>QMUIRoundButton（自定义 View，用 29 次）<br>CommonImageButton（Composable，用 6 次）<br>CommonTextButton（Composable，用 2 次）<br>OriginOpenButton（Composable，用 2 次）<br>CommonButton（Composable，用 2 次）<br>QMUIAlphaImageButton（自定义 View，用 2 次）<br>PickCurrentCheckButton（Composable，用 1 次）<br>QMUIAlphaButton（自定义 View，用 0 次） | `QMUIButton` |
| list | LazyColumn（组件库，用 9 次）<br>QMUIContinuousNestedBottomRecyclerView（自定义 View，用 6 次）<br>QMUIContinuousNestedTopRecyclerView（自定义 View，用 4 次）<br>QMUIAnimationListView（自定义 View，用 4 次）<br>RecyclerView（组件库，用 6 次）<br>LazyRow（组件库，用 1 次）<br>ListView（组件库，用 2 次）<br>QMUIWrapContentListView（自定义 View，用 1 次） | `LazyColumn` |
| loading | QMUILoadingView（自定义 View，用 5 次）<br>QMUIProgressBar（自定义 View，用 6 次）<br>QMUIPhotoLoading（Composable，用 3 次）<br>QDLoadingItemView（自定义 View，用 1 次）<br>ProgressBar（组件库，用 1 次） | `QMUILoadingView` |
| dialog | QMUIDialog（自定义 View，用 10 次）<br>QMUITipDialog（自定义 View，用 14 次）<br>QMUIModal（Composable，用 0 次）<br>QMUIBaseDialog（自定义 View，用 0 次） | `QMUIDialog` |
| cell | QMUICommonListItemView（自定义 View，用 53 次）<br>QMUIItem（Composable，用 10 次）<br>QMUIBottomSheetListItemView（自定义 View，用 1 次） | `QMUICommonListItemView` |
| popup | QMUIPopup（自定义 View，用 14 次）<br>QMUIBottomSheet（自定义 View，用 38 次）<br>QMUIFullScreenPopup（自定义 View，用 3 次） | `QMUIPopup` |
| image | ImageView（组件库，用 11 次）<br>QMUIRadiusImageView（自定义 View，用 8 次）<br>QMUIRadiusImageView2（自定义 View，用 7 次） | `ImageView` |
| tabs | QMUITabSegment（组件库，用 4 次）<br>QMUIBasicTabSegment（自定义 View，用 4 次）<br>QMUITabView（自定义 View，用 2 次） | `QMUITabSegment` |
| navbar | QMUITopBar（自定义 View，用 6 次）<br>QMUIGroupListSectionHeaderFooterView（自定义 View，用 2 次）<br>QDSectionHeaderView（自定义 View，用 1 次） | `QMUITopBar` |
| switch | RadioButton（组件库，用 5 次）<br>QMUICheckBox（Composable，用 1 次）<br>QMUIPhotoPickCheckBox（Composable，用 1 次） | `RadioButton` |
| carousel | HorizontalPager（组件库，用 1 次）<br>SwipeBackgroundView（自定义 View，用 2 次）<br>ViewPager（组件库，用 1 次） | `HorizontalPager` |
| input | EmojiEditText（自定义 View，用 3 次）<br>EditText（组件库，用 2 次） | `EmojiEditText` |
| icon | QMUIMarkIcon（Composable，用 1 次）<br>QMUIChevronIcon（Composable，用 0 次） | `QMUIMarkIcon` |

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
