# 新设计稿复用清单：ProductDetail 商品详情 → gkd（android）

Figma 组件 13 个：**可直接复用 5** · 需确认 2 · 需新建 6

| Figma 组件 | 实例 | 结论 | 推荐复用 | 依据 | 项目里的真实用法 |
|---|---|---|---|---|---|
| Swiper | 1 | ❓ 需确认 | `HorizontalPager`（组件库 Compose Foundation） | 0.65：同为 carousel，专用组件（pager），项目里已用 1 次 | `HorizontalPager(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/page/GkImagePreviewContent.kt:156 |
| NavBar | 1 | ✅ 直接复用 | `GkTopAppBar`（Composable gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkTopAppBar.kt） | 1.00：同为 navbar，项目里已用 28 次，设计系统组件 | `GkTopAppBar(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkEditorLayout.kt:23 |
| Price | 1 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| CountDown | 1 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| Card | 3 | ✅ 直接复用 | `Card`（组件库 Compose Material3） | 1.00：同名，项目里已用 25 次 | `Card(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkAppRuleRestrictionCard.kt:28 |
| Tag | 4 | ❓ 需确认 | `InputChip`（组件库 Compose Material3） | 0.69：同为 tag，专用组件（input），标准组件，项目里已用 1 次 | `InputChip(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/subscription/ActionLogPage.kt:204 |
| NoticeBar | 1 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| Cell | 3 | ✅ 直接复用 | `GkSettingItem`（Composable gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/GkSettingItem.kt） | 1.00：同为 cell，变体轴匹配 1/2，项目里已用 23 次，设计系统组件 | `GkSettingItem(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/home/SettingsPage.kt:97 |
| Divider | 2 | ✅ 直接复用 | `HorizontalDivider`（组件库 Compose Material3） | 0.87：同为 divider，项目里已用 7 次 | `HorizontalDivider()` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/crash/CrashReportPage.kt:266 |
| Avatar | 1 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| Rate | 1 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| ActionBar | 1 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| Button | 2 | ✅ 直接复用 | `TextButton`（组件库 Compose Material3） | 0.97：同为 button，项目里已用 72 次 | `TextButton(` gkd-app/src/commonMain/kotlin/li/gkd/app/ui/component/DialogRequests.kt:122 |

## 变体 → 属性（推荐组件怎么配置）

- **Card** → `Card`：
  - Title：组件没有对应属性，需确认（用插槽 / 自行组合）
- **Tag** → `InputChip`：
  - Tone：Primary / Success 没有对应的枚举属性，需确认（可能要用 color / class 等实现）
  - Variant：Light / Solid 没有对应的枚举属性，需确认（可能要用 color / class 等实现）
- **Cell** → `GkSettingItem`：
  - Arrow：组件没有对应属性，需确认（用插槽 / 自行组合）
  - Label=On → subtitle = true / 传入内容；Label=Off → 不传
- **Button** → `TextButton`：
  - Size：Large 没有对应的枚举属性，需确认（可能要用 color / class 等实现）
  - Type：Default / Primary 没有对应的枚举属性，需确认（可能要用 color / class 等实现）

## 候选明细

- **Swiper**：`HorizontalPager`（组件库，0.65，用 1 次）、`VerticalPager`（组件库，0.64，用 0 次）
- **NavBar**：`GkTopAppBar`（Composable，1.00，用 28 次）、`GkMultiSelectionTopAppBar`（Composable，0.85，用 9 次）、`GkPrivilegeTopBar`（Composable，0.83，用 2 次）、`GkRuleListHeader`（Composable，0.76，用 2 次）
- **Card**：`Card`（组件库，1.00，用 25 次）、`GkCopyTextCard`（Composable，0.90，用 3 次）、`GkExpandableSection`（Composable，0.85，用 1 次）、`OutlinedCard`（组件库，0.77，用 1 次）
- **Tag**：`InputChip`（组件库，0.69，用 1 次）、`AssistChip`（组件库，0.65，用 1 次）、`FilterChip`（组件库，0.64，用 3 次）、`SuggestionChip`（组件库，0.60，用 0 次）
- **Cell**：`GkSettingItem`（Composable，1.00，用 23 次）、`GkRuleListItem`（Composable，0.87，用 5 次）、`GkBatchActionMenuItem`（Composable，0.76，用 22 次）、`GkTemplateVariableRow`（Composable，0.68，用 7 次）
- **Divider**：`HorizontalDivider`（组件库，0.87，用 7 次）、`GkTwoLineText`（Composable，0.82，用 10 次）、`Divider`（组件库，0.55，用 0 次）
- **Button**：`TextButton`（组件库，0.97，用 72 次）、`GkIconButton`（Composable，0.95，用 71 次）、`Button`（组件库，0.95，用 3 次）、`OutlinedButton`（组件库，0.92，用 16 次）
