# 新设计稿复用清单：ProductDetail 商品详情 → androidproject（android）

Figma 组件 13 个：**可直接复用 5** · 需确认 2 · 需新建 6

| Figma 组件 | 实例 | 结论 | 推荐复用 | 依据 | 项目里的真实用法 |
|---|---|---|---|---|---|
| Swiper | 1 | ❓ 需确认 | `NestedScrollViewPager`（组件库 com.hjq.nested） | 0.69：同为 carousel，专用组件（nestedscroll），标准组件，项目里已用 1 次 | `<com.hjq.nested.scroll.layout.NestedScrollViewPager` app/src/main/res/layout/home_main_fragment.xml:92 |
| NavBar | 1 | ✅ 直接复用 | `TitleBar`（组件库 com.hjq.bar） | 0.95：同为 navbar，标准组件，项目里已用 13 次 | `<com.hjq.bar.TitleBar` app/src/main/res/layout/about_activity.xml:11 |
| Price | 1 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| CountDown | 1 | ❓ 需确认 | `CountdownView`（自定义 View library/customWidget/src/main/java/com/hjq/custom/widget/view/CountdownView.java） | 0.89：同名，专门用途（验证码），项目里已用 6 次 | `<com.hjq.custom.widget.view.CountdownView` app/src/main/res/layout/home_find_fragment.xml:270 |
| Card | 3 | ✅ 直接复用 | `CardView`（组件库 Android SDK widgets） | 0.93：同名，项目里已用 2 次 | `<androidx.cardview.widget.CardView xmlns:android="http://schemas.android.com/apk/res/android"` app/src/main/res/layout/custom_dialog.xml:2 |
| Tag | 4 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| NoticeBar | 1 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| Cell | 3 | ✅ 直接复用 | `SettingBar`（自定义 View library/customWidget/src/main/java/com/hjq/custom/widget/layout/SettingBar.java） | 0.92：同为 cell，项目里已用 14 次 | `<com.hjq.custom.widget.layout.SettingBar` app/src/main/res/layout/home_find_fragment.xml:134 |
| Divider | 2 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| Avatar | 1 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| Rate | 1 | ✅ 直接复用 | `SimpleRatingBar`（自定义 View library/customWidget/src/main/java/com/hjq/custom/widget/view/SimpleRatingBar.java） | 0.84：同为 rate，项目里已用 4 次 | `<com.hjq.custom.widget.view.SimpleRatingBar` app/src/main/res/layout/home_find_fragment.xml:282 |
| ActionBar | 1 | 🆕 需新建 | — | 项目和组件库里都没有 | — |
| Button | 2 | ✅ 直接复用 | `Button`（组件库 Android SDK widgets） | 0.95：同名，系统原生控件，项目里已用 54 次 | `<androidx.appcompat.widget.AppCompatButton` app/src/main/res/layout/custom_dialog.xml:34 |

## 变体 → 属性（推荐组件怎么配置）

- **CountDown** → `CountdownView`：
  - 项目惯用 style="@style/CountdownViewStyle"（6 处）
- **Card** → `CardView`：
  - Title：组件没有对应属性，需确认（用插槽 / 自行组合）
- **Cell** → `SettingBar`：
  - Arrow：组件没有对应属性，需确认（用插槽 / 自行组合）
  - Label：组件没有对应属性，需确认（用插槽 / 自行组合）
- **Rate** → `SimpleRatingBar`：
  - Size：Small 没有对应的枚举属性，需确认（可能要用 color / class 等实现）
- **Button** → `Button`：
  - Size：Large 没有对应的枚举属性，需确认（可能要用 color / class 等实现）
  - Type：Default / Primary 没有对应的枚举属性，需确认（可能要用 color / class 等实现）

## 候选明细

- **NavBar**：`TitleBar`（组件库，0.95，用 13 次）、`MaterialToolbar`（组件库，0.61，用 0 次）
- **Button**：`Button`（组件库，0.95，用 54 次）、`SubmitButton`（自定义 View，0.80，用 2 次）、`PlayButton`（自定义 View，0.67，用 0 次）、`FloatActionButton`（自定义 View，0.50，用 0 次）
