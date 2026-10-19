# Android 端细节

- 工程：`android/`。Compose 组件库 `android/ds/`（包 `com.demo.ds`，`Ds*`）；`android/ds-view/` 是 XML View 存量组件（`Legacy*`），**新页面不要用**。
- 目标项目 compileSdk 32 / minSdk 23：只能用 Compose 1.2 已有的 API（Box/Row/Column、padding/size/background、Arrangement.spacedBy 等）。
- 本环境没有 Android SDK，用 Compose Desktop 在 JVM 上渲染（API 与 Android Compose 相同）。
- 页面：`android/app/src/main/kotlin/com/demo/app/screens/<Page>/<Page>Screen.kt`，并在 `android/app/src/main/kotlin/com/demo/app/Screens.kt` 注册 `"<Page>" to { <Page>Screen() }`。
- Figma id：`Modifier.figmaId("1:23")`（`import com.demo.ds.figmaId`），传给组件的 `modifier` 参数。
- 页面级文本：`DsText(text, style = DsTypography.Title, color = DsColors.TextPrimary)`；间距 `DsSpacing.*`，颜色 `DsColors.*`。
- `DsBottomBar` 不负责定位：页面根用 `Box(Modifier.fillMaxSize())`，内容 `Column` 加底部留白 56.dp，`DsBottomBar(Modifier.align(Alignment.BottomCenter))`。
- 构建 + 渲染由 `figmafe verify` 自动完成（`gradle -q :harness:installDist` → `harness <Page> <out> 375 900`）。Maven Central 偶尔 429，重跑即可。
