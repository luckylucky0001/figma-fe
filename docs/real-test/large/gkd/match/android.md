# 新设计稿复用清单：TradeOrder 买入下单 → gkd（android）

Figma 组件 10 个：**可直接复用 8** · 需确认 1 · 需新建 1

| Figma 组件 | 实例 | 结论 | 推荐复用 | 依据 | 项目里的真实用法 |
|---|---|---|---|---|---|
| NavBar（未组件化） | 1 | ✅ 直接复用 | `GkTopAppBar`（Composable gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkTopAppBar.kt） | 1.00：同为 navbar，项目里已用 28 次，设计系统组件 | `GkTopAppBar(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkEditorLayout.kt:23 |
| SearchBar（未组件化） | 1 | ✅ 直接复用 | `GkAppBarTextField`（Composable gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkAppBarTextField.kt） | 0.94：同为 search，6 处调用里有 6 处传的是搜索词，是项目的搜索框，项目里已用 6 次，设计系统组件 | `GkAppBarTextField(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/AppListPage.kt:171 |
| Tabs（未组件化） | 1 | ❓ 需确认 | `TabRow`（组件库 Compose Material3） | 0.72：同为 tabs | `TabRow(…)  // import androidx.compose.material3.TabRow` （Compose Material3 标准用法，项目未用过） |
| Cell（未组件化） | 2 | ✅ 直接复用 | `GkSettingItem`（Composable gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkSettingItem.kt） | 1.00：同为 cell，项目里已用 23 次，设计系统组件 | `GkSettingItem(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/SettingsPage.kt:97 |
| Stepper（未组件化） | 2 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| Input（未组件化） | 1 | ✅ 直接复用 | `GkMultiTextField`（Composable gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkMultiTextField.kt） | 0.94：同为 input，项目里已用 6 次，设计系统组件 | `GkMultiTextField(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkMultiTextField.kt:83 |
| Switch（未组件化） | 1 | ✅ 直接复用 | `GkSwitch`（Composable gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkSwitch.kt） | 1.00：同名，可选中组件，设计稿里是静态展示，项目里已用 7 次，设计系统组件 | `GkSwitch(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkSubsItemCard.kt:260 |
| Checkbox（未组件化） | 1 | ✅ 直接复用 | `GkCheckbox`（Composable gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkCheckbox.kt） | 1.00：同名，可选中组件，设计稿里是静态展示，项目里已用 6 次，设计系统组件 | `GkCheckbox(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkAppCheckboxCard.kt:67 |
| Button（未组件化） | 1 | ✅ 直接复用 | `TextButton`（组件库 Compose Material3） | 0.97：同为 button，项目里已用 72 次 | `TextButton(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/DialogRequests.kt:122 |
| Empty（未组件化） | 1 | ✅ 直接复用 | `GkEmptyState`（Composable gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkEmptyState.kt） | 0.88：名称包含，项目里已用 15 次，设计系统组件 | `GkEmptyState(text = stringResource(Res.string.data_empty))` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkLogTimeline.kt:116 |

## 候选明细

- **NavBar（未组件化）**：`GkTopAppBar`（Composable，1.00，用 28 次）、`GkMultiSelectionTopAppBar`（Composable，0.85，用 9 次）、`GkPrivilegeTopBar`（Composable，0.83，用 2 次）、`GkRuleListHeader`（Composable，0.76，用 2 次）
- **SearchBar（未组件化）**：`GkAppBarTextField`（Composable，0.94，用 6 次）、`SearchBar`（组件库，0.77，用 0 次）
- **Tabs（未组件化）**：`TabRow`（组件库，0.72，用 0 次）、`GkRuleCategoryTabs`（Composable，0.70，用 2 次）、`Tab`（组件库，0.69，用 2 次）、`PrimaryScrollableTabRow`（组件库，0.65，用 1 次）
- **Cell（未组件化）**：`GkSettingItem`（Composable，1.00，用 23 次）、`GkRuleListItem`（Composable，0.87，用 5 次）、`GkBatchActionMenuItem`（Composable，0.76，用 22 次）、`GkTemplateVariableRow`（Composable，0.68，用 7 次）
- **Input（未组件化）**：`GkMultiTextField`（Composable，0.94，用 6 次）、`OutlinedTextField`（组件库，0.89，用 10 次）、`GkOutlinedTextField`（Composable，0.88，用 2 次）、`TextField`（组件库，0.69，用 2 次）
- **Switch（未组件化）**：`GkSwitch`（Composable，1.00，用 7 次）、`GkTextSwitch`（Composable，0.99，用 38 次）、`Switch`（组件库，0.78，用 2 次）、`GkTriStateSwitch`（Composable，0.70，用 2 次）
- **Checkbox（未组件化）**：`GkCheckbox`（Composable，1.00，用 6 次）、`GkMenuItemCheckbox`（Composable，0.77，用 7 次）、`Checkbox`（组件库，0.75，用 1 次）
- **Button（未组件化）**：`TextButton`（组件库，0.97，用 72 次）、`GkIconButton`（Composable，0.95，用 71 次）、`Button`（组件库，0.95，用 3 次）、`OutlinedButton`（组件库，0.92，用 16 次）
