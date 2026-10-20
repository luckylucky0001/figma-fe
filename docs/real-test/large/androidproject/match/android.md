# 新设计稿复用清单：TradeOrder 买入下单 → androidproject（android）

Figma 组件 10 个：**可直接复用 7** · 需确认 1 · 需新建 2

| Figma 组件 | 实例 | 结论 | 推荐复用 | 依据 | 项目里的真实用法 |
|---|---|---|---|---|---|
| NavBar（未组件化） | 1 | ✅ 直接复用 | `TitleBar`（组件库 com.hjq.bar） | 0.95：同为 navbar，标准组件，项目里已用 13 次 | `<com.hjq.bar.TitleBar` app/src/main/res/layout/about_activity.xml:11 |
| SearchBar（未组件化） | 1 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| Tabs（未组件化） | 1 | ❓ 需确认 | `TabLayout`（组件库 Material Components） | 0.61：同为 tabs，项目几乎不用 Material Components，标准组件 | `<com.google.android.material.…TabLayout … />` （Material Components 标准用法，项目未用过） |
| Cell（未组件化） | 2 | ✅ 直接复用 | `SettingBar`（自定义 View library/customWidget/src/main/java/com/hjq/custom/widget/layout/SettingBar.java） | 0.92：同为 cell，项目里已用 14 次 | `<com.hjq.custom.widget.layout.SettingBar` app/src/main/res/layout/home_find_fragment.xml:134 |
| Stepper（未组件化） | 2 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| Input（未组件化） | 1 | ✅ 直接复用 | `RegexEditText`（自定义 View library/customWidget/src/main/java/com/hjq/custom/widget/view/RegexEditText.java） | 0.92：同为 input，项目里已用 14 次 | `<com.hjq.custom.widget.view.RegexEditText` app/src/main/res/layout/home_find_fragment.xml:40 |
| Switch（未组件化） | 1 | ✅ 直接复用 | `SwitchButton`（自定义 View library/customWidget/src/main/java/com/hjq/custom/widget/view/SwitchButton.java） | 0.88：同为 switch，项目里已用 8 次 | `<com.hjq.custom.widget.view.SwitchButton` app/src/main/res/layout/home_find_fragment.xml:241 |
| Checkbox（未组件化） | 1 | ✅ 直接复用 | `CheckBox`（组件库 Android SDK widgets） | 0.80：同名，系统原生控件，项目里已用 3 次 | `<androidx.appcompat.widget.AppCompatCheckBox` app/src/main/res/layout/home_find_fragment.xml:181 |
| Button（未组件化） | 1 | ✅ 直接复用 | `Button`（组件库 Android SDK widgets） | 0.95：同名，系统原生控件，项目里已用 54 次 | `<androidx.appcompat.widget.AppCompatButton` app/src/main/res/layout/custom_dialog.xml:34 |
| Empty（未组件化） | 1 | ✅ 直接复用 | `StatusLayout`（自定义 View app/src/main/java/com/hjq/demo/widget/StatusLayout.java） | 0.84：同为 empty，项目里已用 4 次 | `<com.hjq.demo.widget.StatusLayout` app/src/main/res/layout/status_activity.xml:15 |

## 变体 → 属性（推荐组件怎么配置）

- **Input（未组件化）** → `RegexEditText`：
  - 项目惯用 style="@style/EditTextStyle"（14 处）

## 候选明细

- **NavBar（未组件化）**：`TitleBar`（组件库，0.95，用 13 次）、`MaterialToolbar`（组件库，0.61，用 0 次）
- **Input（未组件化）**：`RegexEditText`（自定义 View，0.92，用 14 次）、`ClearEditText`（自定义 View，0.80，用 2 次）、`PasswordEditText`（自定义 View，0.76，用 6 次）、`TextInputEditText`（组件库，0.53，用 0 次）
- **Checkbox（未组件化）**：`CheckBox`（组件库，0.80，用 3 次）、`MaterialCheckBox`（组件库，0.70，用 0 次）
- **Button（未组件化）**：`Button`（组件库，0.95，用 54 次）、`SubmitButton`（自定义 View，0.80，用 2 次）、`PlayButton`（自定义 View，0.62，用 0 次）、`FloatActionButton`（自定义 View，0.50，用 0 次）
