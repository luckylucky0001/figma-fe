# 存量复用体检：androidproject（android）

技术栈：XML View / Material Components · 扫描文件 198

| 指标 | 值 |
|---|---|
| 组件化率（组件使用 / 组件使用 + 可替换的原生元素） | 100%（235 / 236） |
| 重复概念（同一 UI 概念有多个实现） | 0 |
| 重复结构（可抽成组件的复制粘贴结构） | 2 组 |
| 硬编码颜色 | 44 处（其中 5 处已有同值 token） |
| 硬编码尺寸 | 3 处 |

## 0. 优先处理（按影响面排序）

1. 5 处写死的颜色已有同值 token，可以直接替换（见第 6 节）
2. 23 个元素的结构被复制了 2 次（app/src/main/res/layout/dialog_activity.xml:2 等）→ 抽成组件
3. 1 处原生 Toolbar：项目里已有 TitleBar（组件库）

## 1. 三方组件库使用

| 组件库 | 使用次数 | 最常用 |
|---|---|---|
| Android SDK widgets | 103 | Button×54、TextView×23、ImageView×13、CheckBox×3、RadioButton×3、RecyclerView×3、CardView×2、ProgressBar×1 |
| com.hjq.bar | 13 | TitleBar×13 |
| com.hjq.shape | 5 | ShapeTextView×4、ShapeLinearLayout×1 |
| Material Components | 1 | AppBarLayout×1 |
| com.airbnb.lottie | 1 | LottieAnimationView×1 |
| com.hjq.nested | 1 | NestedScrollViewPager×1 |
| me.relex.circleindicator | 1 | CircleIndicator3×1 |

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
- `MaterialDivider`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `MaterialSwitch`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `SearchBar`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `SearchView`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `Slider`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃
- `TextInputLayout`：需要 Theme.MaterialComponents 主题，项目主题 AppTheme 继承 AppCompat，直接用会崩溃

## 2. 组件资产（可复用的现成组件）

| 组件 | 类型 | 概念 | 使用次数 | 真实用法示例 |
|---|---|---|---|---|
| `Button` | 组件库 | button | 54 | `<androidx.appcompat.widget.AppCompatButton` app/src/main/res/layout/custom_dialog.xml:34 |
| `TextView` | 组件库 |  | 23 | `<androidx.appcompat.widget.AppCompatTextView` app/src/main/res/layout/about_activity.xml:34 |
| `ImageView` | 组件库 | image | 13 | `<androidx.appcompat.widget.AppCompatImageView` app/src/main/res/layout/custom_dialog.xml:20 |
| `TitleBar` | 组件库 | navbar | 13 | `<com.hjq.bar.TitleBar` app/src/main/res/layout/about_activity.xml:11 |
| `ShapeTextView` | 组件库 |  | 4 | `<com.hjq.shape.view.ShapeTextView` app/src/main/res/layout/home_find_fragment.xml:312 |
| `CheckBox` | 组件库 | checkbox | 3 | `<androidx.appcompat.widget.AppCompatCheckBox` app/src/main/res/layout/home_find_fragment.xml:181 |
| `RadioButton` | 组件库 | radio | 3 | `<androidx.appcompat.widget.AppCompatRadioButton` app/src/main/res/layout/home_find_fragment.xml:211 |
| `RecyclerView` | 组件库 | list | 3 | `<androidx.recyclerview.widget.RecyclerView` app/src/main/res/layout/home_activity.xml:16 |
| `CardView` | 组件库 | card | 2 | `<androidx.cardview.widget.CardView xmlns:android="http://schemas.android.com/apk/res/android"` app/src/main/res/layout/custom_dialog.xml:2 |
| `AppBarLayout` | 组件库 |  | 1 | `<com.google.android.material.appbar.AppBarLayout` app/src/main/res/layout/home_main_fragment.xml:9 |
| `CircleIndicator3` | 组件库 |  | 1 | `<me.relex.circleindicator.CircleIndicator3` app/src/main/res/layout/guide_activity.xml:22 |
| `LottieAnimationView` | 组件库 |  | 1 | `<com.airbnb.lottie.LottieAnimationView` app/src/main/res/layout/splash_activity.xml:11 |
| `NestedScrollViewPager` | 组件库 | carousel | 1 | `<com.hjq.nested.scroll.layout.NestedScrollViewPager` app/src/main/res/layout/home_main_fragment.xml:92 |
| `ProgressBar` | 组件库 | loading | 1 | `<ProgressBar` app/src/main/res/layout/update_dialog.xml:87 |
| `ShapeLinearLayout` | 组件库 |  | 1 | `<com.hjq.shape.layout.ShapeLinearLayout xmlns:android="http://schemas.android.com/apk/res/android"` app/src/main/res/layout/pay_password_dialog.xml:2 |
| `ViewPager2` | 组件库 |  | 1 | `<androidx.viewpager2.widget.ViewPager2` app/src/main/res/layout/guide_activity.xml:10 |
| `BaseDialog` | 自定义 View | dialog | 29 | `public final class DialogManager implements LifecycleEventObserver, BaseDialog.OnDismissListener {` app/src/main/java/com/hjq/demo/manager/DialogManager.java:21 |
| `RegexEditText` | 自定义 View | input | 14 | `<com.hjq.custom.widget.view.RegexEditText` app/src/main/res/layout/home_find_fragment.xml:40 |
| `SettingBar` | 自定义 View | cell | 14 | `<com.hjq.custom.widget.layout.SettingBar` app/src/main/res/layout/home_find_fragment.xml:134 |
| `SwitchButton` | 自定义 View | switch | 8 | `<com.hjq.custom.widget.view.SwitchButton` app/src/main/res/layout/home_find_fragment.xml:241 |
| `CountdownView` | 自定义 View | countdown | 6 | `<com.hjq.custom.widget.view.CountdownView` app/src/main/res/layout/home_find_fragment.xml:270 |
| `PasswordEditText` | 自定义 View | input | 6 | `<com.hjq.custom.widget.view.PasswordEditText` app/src/main/res/layout/home_find_fragment.xml:32 |
| `SimpleRatingBar` | 自定义 View |  | 4 | `<com.hjq.custom.widget.view.SimpleRatingBar` app/src/main/res/layout/home_find_fragment.xml:282 |
| `StatusLayout` | 自定义 View | empty | 4 | `<com.hjq.demo.widget.StatusLayout` app/src/main/res/layout/status_activity.xml:15 |
| `PasswordView` | 自定义 View |  | 3 | `<com.hjq.demo.widget.PasswordView` app/src/main/res/layout/pay_password_dialog.xml:66 |
| `PlayerView` | 自定义 View |  | 3 | `implements PlayerView.OnPlayListener {` app/src/main/java/com/hjq/demo/ui/activity/common/VideoPlayActivity.java:25 |
| `ScaleImageView` | 自定义 View | image | 3 | `<com.hjq.custom.widget.view.ScaleImageView` app/src/main/res/layout/about_activity.xml:16 |
| `SmartTextView` | 自定义 View |  | 3 | `<com.hjq.custom.widget.view.SmartTextView` app/src/main/res/layout/pay_password_dialog.xml:20 |
| `XCollapsingToolbarLayout` | 自定义 View |  | 3 | `<com.hjq.demo.widget.XCollapsingToolbarLayout` app/src/main/res/layout/home_main_fragment.xml:16 |
| `ClearEditText` | 自定义 View | input | 2 | `<com.hjq.custom.widget.view.ClearEditText` app/src/main/res/layout/home_find_fragment.xml:25 |
| `PressAlphaTextView` | 自定义 View |  | 2 | `return new PressAlphaTextView(context);` app/src/main/java/com/hjq/demo/other/TitleBarStyle.java:31 |
| `SubmitButton` | 自定义 View | button | 2 | `<com.hjq.custom.widget.view.SubmitButton` app/src/main/res/layout/login_activity.xml:86 |
| `BottomSheetDialog` | 自定义 View | sheet | 1 | `BottomSheetDialog dialog = new BottomSheetDialog(context, themeId);` app/src/main/java/com/hjq/demo/ui/dialog/common/AlbumDialog.java:97 |
| `DrawableTextView` | 自定义 View |  | 1 | `<com.hjq.custom.widget.view.DrawableTextView` app/src/main/res/layout/home_find_fragment.xml:300 |
| `SlantedTextView` | 自定义 View |  | 1 | `<com.hjq.custom.widget.view.SlantedTextView` app/src/main/res/layout/splash_activity.xml:21 |
| `WrapRecyclerView` | 自定义 View | list | 1 | `<com.hjq.custom.widget.layout.WrapRecyclerView` app/src/main/res/layout/status_fragment.xml:15 |
| `BrowserView` | 自定义 View |  | 0 | app/src/main/java/com/hjq/demo/widget/webview/BrowserView.java |
| `CustomViewStub` | 自定义 View |  | 0 | library/customWidget/src/main/java/com/hjq/custom/widget/layout/CustomViewStub.java |
| `FloatActionButton` | 自定义 View | button | 0 | library/customWidget/src/main/java/com/hjq/custom/widget/view/FloatActionButton.java |
| `PlayButton` | 自定义 View | button | 0 | library/customWidget/src/main/java/com/hjq/custom/widget/view/PlayButton.java |
| `SimpleLayout` | 自定义 View |  | 0 | library/customWidget/src/main/java/com/hjq/custom/widget/layout/SimpleLayout.java |

## 4. 有组件却没用（1 处，含绕过设计系统封装直接用组件库）

- app/src/main/res/layout/home_main_fragment.xml:31 `<androidx.appcompat.widget.Toolbar` → 原生 Toolbar：项目里已有 TitleBar（组件库）

## 5. 重复结构（复制粘贴的 UI 片段，可抽成组件）

- **23 个元素 × 2 处**：app/src/main/res/layout/dialog_activity.xml:2、app/src/main/res/layout/home_mine_fragment.xml:2

```
<LinearLayout>
  <com.hjq.bar.TitleBar>
  <androidx.core.widget.NestedScrollView>
    <LinearLayout>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
      <androidx.appcompat.widget.AppCompatButton>
```
- **5 个元素 × 2 处**：app/src/main/res/layout/register_activity.xml:38、app/src/main/res/layout/safe_dialog.xml:31

```
<LinearLayout>
  <FrameLayout>
    <com.hjq.custom.widget.view.RegexEditText>
    <View>
  <com.hjq.custom.widget.view.CountdownView>
```

## 6. 硬编码颜色热点

| 颜色 | 次数 | 已有 token | 出现在 |
|---|---|---|---|
| #333333 | 4 | — | app/src/main/res/layout/pay_password_dialog.xml、library/customWidget/src/main/java/com/hjq/custom/widget/view/SwitchButton.java |
| #222222 | 3 | — | app/src/main/res/layout/about_activity.xml、app/src/main/java/com/hjq/demo/ui/dialog/common/AddressDialog.java |
| #ECECEC | 3 | — | app/src/main/res/layout/pay_password_dialog.xml、app/src/main/java/com/hjq/demo/widget/PasswordView.java、library/customWidget/src/main/java/com/hjq/custom/widget/layout/SettingBar.java |
| #000000 | 2 | `@color/black` | app/src/main/java/com/hjq/demo/other/MaterialHeader.java |
| #11BBFF | 2 | — | app/src/main/java/com/hjq/demo/other/MaterialHeader.java |
| #A4A4A4 | 2 | — | app/src/main/res/layout/home_navigation_item.xml、app/src/main/res/layout/phone_reset_activity.xml |
| #FFFFFF | 2 | `@color/white` | library/customWidget/src/main/java/com/hjq/custom/widget/view/SwitchButton.java |
| #0099CC | 1 | — | app/src/main/java/com/hjq/demo/other/MaterialHeader.java |
| #142DCC | 1 | — | app/src/main/java/com/hjq/demo/other/SmartBallPulseFooter.java |
| #19CC95 | 1 | — | library/customWidget/src/main/java/com/hjq/custom/widget/view/SubmitButton.java |
| #287BDE | 1 | — | app/src/main/java/com/hjq/demo/ui/activity/common/CrashActivity.java |
| #30B399 | 1 | — | app/src/main/java/com/hjq/demo/other/SmartBallPulseFooter.java |
