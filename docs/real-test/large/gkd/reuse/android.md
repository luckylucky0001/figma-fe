# 存量复用体检：gkd（android）

技术栈：Jetpack Compose / Compose Foundation / Compose Material3 · 扫描文件 839

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（910 / 910） |
| 重复概念（同一 UI 概念有多个实现） | 4 |
| 重复结构（可抽成组件的复制粘贴结构） | 0 组 |
| 硬编码颜色 | 31 处（其中 0 处已有同值 token） |

## 0. 优先处理（按影响面排序）

1. 「dialog（Compose）」有 3 种实现（共用 36 次）→ 统一到 `GkDialog`
2. 「card（Compose）」有 4 种实现（共用 30 次）→ 统一到 `Card`
3. 「input（Compose）」有 3 种实现（共用 18 次）→ 统一到 `GkMultiTextField`
4. 1 处TextField（Compose Material3），项目设计系统已封装为 GkMultiTextField，应使用封装
5. 1 处FloatingActionButton（Compose Material3），项目设计系统已封装为 GkAnimatedFloatingActionButton，应使用封装

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Compose Material3 | 174 | TextButton×72、Card×25、OutlinedButton×16、OutlinedTextField×10、CircularProgressIndicator×7、HorizontalDivider×7、RadioButton×5、LinearProgressIndicator×4 |
| Compose Foundation | 20 | LazyColumn×16、LazyVerticalGrid×2、HorizontalPager×1、LazyRow×1 |

**组件库里有、但本项目不能直接用的组件**（match 不会推荐它们）：

- `Divider`：已废弃，改用 HorizontalDivider

**设计系统封装**（组件库组件 → 项目的封装；新页面一律用右边的封装）：`Checkbox` → `GkCheckbox`、`FloatingActionButton` → `GkAnimatedFloatingActionButton`、`ModalBottomSheet` → `GkModalBottomSheet`、`Switch` → `GkSwitch`、`TextField` → `GkMultiTextField`、`TopAppBar` → `GkTopAppBar`

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `TextButton` | 组件库 | button | 72 | `TextButton(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/DialogRequests.kt:122 |
| `Card` | 组件库 | card | 25 | `Card(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkAppRuleRestrictionCard.kt:28 |
| `LazyColumn` | 组件库 | list | 16 | `LazyColumn(` gkd-app/src/androidMain/kotlin/li/gkd/app/service/EventService.kt:126 |
| `OutlinedButton` | 组件库 | button | 16 | `OutlinedButton(onClick = flushSettings) { Text("重试保存") }` gkd-app/src/jvmMain/kotlin/li/gkd/app/DesktopControls.kt:73 |
| `OutlinedTextField` | 组件库 | input | 10 | `OutlinedTextField(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkGithubUploadDialogs.kt:68 |
| `CircularProgressIndicator` | 组件库 | loading | 7 | `CircularProgressIndicator()` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkSubscriptionPageContent.kt:75 |
| `HorizontalDivider` | 组件库 | divider | 7 | `HorizontalDivider()` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/crash/CrashReportPage.kt:266 |
| `RadioButton` | 组件库 | radio | 5 | `RadioButton(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkMenu.kt:72 |
| `LinearProgressIndicator` | 组件库 | loading | 4 | `LinearProgressIndicator(progress = { progress }, modifier = Modifier.fillMaxWidth())` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkGithubUploadDialogs.kt:181 |
| `Button` | 组件库 | button | 3 | `Button(onClick = onRetry) {` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkWebViewErrorContent.kt:72 |
| `FilterChip` | 组件库 | tag | 3 | `FilterChip(selected = selectedService == service, onClick = { selectedService = service },` gkd-app/src/jvmMain/kotlin/li/gkd/app/DesktopControls.kt:369 |
| `AlertDialog` | 组件库 | dialog | 2 | `AlertDialog(` gkd-app/src/jvmMain/kotlin/li/gkd/app/DesktopOverlays.kt:30 |
| `BottomAppBar` | 组件库 | bottomappbar | 2 | `BottomAppBar {` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/crash/CrashReportPage.kt:133 |
| `FloatingActionButton` | 组件库 | button | 2 | `FloatingActionButton(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkAnimatedFloatingActionButton.kt:60 |
| `LazyVerticalGrid` | 组件库 |  | 2 | `LazyVerticalGrid(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/snapshot/SnapshotPage.kt:404 |
| `Slider` | 组件库 |  | 2 | `Slider(` gkd-app/src/jvmMain/kotlin/li/gkd/app/DesktopControls.kt:312 |
| `Switch` | 组件库 | switch | 2 | `Switch(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkRuleEnableControl.kt:105 |
| `Tab` | 组件库 | tabs | 2 | `if (firstTabTitle != null) Tab(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkRuleCategoryTabs.kt:32 |
| `TextField` | 组件库 | input | 2 | `TextField(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkMultiTextField.kt:47 |
| `AssistChip` | 组件库 | tag | 1 | `AssistChip(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/subscription/GkRuleGroupDialog.kt:166 |
| `Checkbox` | 组件库 | checkbox | 1 | `Checkbox(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkCheckbox.kt:65 |
| `HorizontalPager` | 组件库 | carousel | 1 | `HorizontalPager(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/page/GkImagePreviewContent.kt:156 |
| `InputChip` | 组件库 | tag | 1 | `InputChip(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/subscription/ActionLogPage.kt:204 |
| `LazyRow` | 组件库 | list | 1 | `LazyRow(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/subscription/ActionLogPage.kt:194 |
| `ModalBottomSheet` | 组件库 | sheet | 1 | `ModalBottomSheet(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkDialog.kt:99 |
| `NavigationBar` | 组件库 | tabbar | 1 | `NavigationBar(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/GkHomeScaffold.kt:26 |
| `OutlinedCard` | 组件库 | card | 1 | `OutlinedCard(modifier = Modifier.padding(horizontal = 16.dp).fillMaxWidth()) {` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/settings/ExternalCallsPage.kt:74 |
| `PrimaryScrollableTabRow` | 组件库 | tabs | 1 | `PrimaryScrollableTabRow(selectedTabIndex = selectedIndex, edgePadding = 8.dp) {` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkRuleCategoryTabs.kt:31 |
| `TopAppBar` | 组件库 | navbar | 1 | `TopAppBar(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkTopAppBar.kt:44 |
| `GkIconButton` | Composable | button | 71 | `GkIconButton(` gkd-app/src/androidMain/kotlin/li/gkd/app/service/EventService.kt:151 |
| `GkIcon` | Composable | icon | 59 | `GkIcon(` gkd-app/src/androidMain/kotlin/li/gkd/app/service/ActivityService.kt:151 |
| `GkPageBottomSpace` | Composable |  | 38 | `GkPageBottomSpace()` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkLogTimeline.kt:118 |
| `GkTextSwitch` | Composable | switch | 38 | `GkTextSwitch(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/SettingsPage.kt:144 |
| `GkAlertDialog` | Composable | dialog | 29 | `GkAlertDialog(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/DialogRequests.kt:113 |
| `GkTopAppBar` | Composable | navbar | 28 | `GkTopAppBar(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkEditorLayout.kt:23 |
| `GkScaffold` | Composable |  | 27 | `GkScaffold(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkEditorLayout.kt:19 |
| `GkSettingItem` | Composable | cell | 23 | `GkSettingItem(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/SettingsPage.kt:97 |
| `GkAppNameText` | Composable |  | 22 | `GkAppNameText(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkAppNameText.kt:130 |
| `GkBatchActionMenuItem` | Composable | cell | 22 | `GkBatchActionMenuItem(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkRuleBatchMenuItems.kt:14 |
| `GkEmptyState` | Composable |  | 15 | `GkEmptyState(text = stringResource(Res.string.data_empty))` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkLogTimeline.kt:116 |
| `GkAppIcon` | Composable | icon | 10 | `icon = { GkAppIcon(appInfo.id, 32.dp) },` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/AppListPage.kt:338 |
| `GkTwoLineText` | Composable | divider | 10 | `GkTwoLineText(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/subscription/ActionLogPage.kt:151 |
| `GkAnimatedFloatingActionButton` | Composable | button | 9 | `GkAnimatedFloatingActionButton(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/AppListPage.kt:297 |
| `GkFilterIconButton` | Composable | button | 9 | `GkFilterIconButton(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkRuleSortButton.kt:19 |
| `GkMultiSelectionActions` | Composable |  | 9 | `GkMultiSelectionActions(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/SubsManagePage.kt:181 |
| `GkMultiSelectionTopAppBar` | Composable | navbar | 9 | `GkMultiSelectionTopAppBar(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/SubsManagePage.kt:171 |
| `GkTooltipIconButtonBox` | Composable | button | 9 | `GkTooltipIconButtonBox(contentDescription) {` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkAnimatedFloatingActionButton.kt:53 |
| `GkMenuItemCheckbox` | Composable | checkbox | 7 | `GkMenuItemCheckbox(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkAppFilterContent.kt:34 |
| `GkSwitch` | Composable | switch | 7 | `GkSwitch(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkSubsItemCard.kt:260 |
| `GkTemplateVariableRow` | Composable | cell | 7 | `GkTemplateVariableRow(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/ActionToastPage.kt:180 |
| `GkTriggerCount` | Composable |  | 7 | `GkTriggerCount(count, lastTriggerTime = triggerSummary?.lastTime)` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkRuleGroupCard.kt:57 |
| `GkAppBarTextField` | Composable | search | 6 | `GkAppBarTextField(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/AppListPage.kt:171 |
| `GkCheckbox` | Composable | checkbox | 6 | `GkCheckbox(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkAppCheckboxCard.kt:67 |
| `GkCopyableText` | Composable |  | 6 | `GkCopyableText(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkCopyText.kt:135 |
| `GkEditor` | Composable |  | 6 | `GkEditor(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/ActionToastPage.kt:64 |
| `GkMenuItemRadioButton` | Composable | radio | 6 | `GkMenuItemRadioButton(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkAppFilterContent.kt:23 |
| `GkMultiTextField` | Composable | input | 6 | `GkMultiTextField(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkMultiTextField.kt:83 |
| `GkSearchCloseIconButton` | Composable | button | 6 | `GkSearchCloseIconButton(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/AppListPage.kt:229 |
| `GkDialog` | Composable | dialog | 5 | `) = GkDialog(` gkd-app/src/androidMain/kotlin/li/gkd/app/ui/component/GkFullscreenDialog.kt:24 |
| `GkFixedTimeText` | Composable |  | 5 | `GkFixedTimeText(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkEventLogOverlayCard.kt:51 |
| … | | | | 共 175 个，见 reuse.json |

## 3. 重复资产（同一概念多个实现 → 收敛到推荐组件）

| 概念 | 实现 | 推荐统一为 |
|---|---|---|
| card（Compose） | Card（组件库，用 25 次）<br>GkCopyTextCard（Composable，用 3 次）<br>GkExpandableSection（Composable，用 1 次）<br>OutlinedCard（组件库，用 1 次） | `Card` |
| dialog（Compose） | GkDialog（Composable，用 5 次）<br>GkAlertDialog（Composable，用 29 次）<br>AlertDialog（组件库，用 2 次） | `GkDialog` |
| input（Compose） | GkMultiTextField（Composable，用 6 次）<br>OutlinedTextField（组件库，用 10 次）<br>GkOutlinedTextField（Composable，用 2 次） | `GkMultiTextField` |
| carousel（Compose） | GkMatchingBanner（Composable，用 1 次）<br>HorizontalPager（组件库，用 1 次） | `GkMatchingBanner` |

## 4. 有组件却没用（2 处，含绕过设计系统封装直接用组件库）

- gkd-app/src/commonMain/kotlin/li/gkd/app/ui/subscription/SubsAppListPage.kt:226 `FloatingActionButton(onClick = {` → 直接用了 FloatingActionButton（Compose Material3），项目设计系统已封装为 GkAnimatedFloatingActionButton，应使用封装
- gkd-app/src/commonMain/kotlin/li/gkd/app/ui/subscription/UpsertRuleGroupPage.kt:78 `TextField(` → 直接用了 TextField（Compose Material3），项目设计系统已封装为 GkMultiTextField，应使用封装

## 6. 硬编码颜色热点

| 颜色 | 次数 | 已有 token | 出现在 |
|---|---|---|---|
| #5F6368 | 4 | — | gkd-app/src/commonMain/kotlin/li/gkd/app/ui/icon/DragPan.kt、gkd-app/src/commonMain/kotlin/li/gkd/app/ui/icon/LockOpenRight.kt、gkd-app/src/commonMain/kotlin/li/gkd/app/ui/icon/ResetSettings.kt |
| #000000 | 3 | — | gkd-app/src/commonMain/kotlin/li/gkd/app/ui/icon/SportsBasketball.kt、gkd-app/src/commonMain/kotlin/li/gkd/app/ui/style/Color.kt |
| #AF00DB | 2 | — | gkd-app/src/commonMain/kotlin/li/gkd/app/ui/style/Color.kt |
| #E1E4E8 | 2 | — | gkd-app/src/commonMain/kotlin/li/gkd/app/ui/style/Color.kt |
| #FFA07A | 2 | — | gkd-app/src/commonMain/kotlin/li/gkd/app/ui/style/Color.kt |
| #0000FF | 1 | — | gkd-app/src/commonMain/kotlin/li/gkd/app/ui/style/Color.kt |
| #001080 | 1 | — | gkd-app/src/commonMain/kotlin/li/gkd/app/ui/style/Color.kt |
| #098658 | 1 | — | gkd-app/src/commonMain/kotlin/li/gkd/app/ui/style/Color.kt |
| #2E8B57 | 1 | — | gkd-app/src/commonMain/kotlin/li/gkd/app/ui/style/Color.kt |
| #34C759 | 1 | — | gkd-app/src/jvmMain/kotlin/li/gkd/app/GkAndroidWindow.kt |
| #3C3C3C | 1 | — | gkd-app/src/jvmMain/kotlin/li/gkd/app/GkAndroidWindow.kt |
| #669900 | 1 | — | gkd-app/src/commonMain/kotlin/li/gkd/app/ui/style/Color.kt |
