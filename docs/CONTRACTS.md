# 工程约定（CLI / 组件库 / 页面之间的接口）

## 1. 目录

```
design/SPEC.md                         设计规范（Figma 组件库的文字版）
design/source/order-detail.html        生成 Figma fixture 用的设计源（不依赖任何组件库，见第 6 节）
design/figma-export/<page>/            Figma 导出产物（插件导出，或由 fixture 生成）
  export.json                          节点树 + 组件 / 组件集（格式见第 5 节）
  frame@2x.png                         设计图
  variables.tokens.json                Figma Variables（DTCG）
design/component-map.yaml              Figma ↔ H5 / Android 映射（人工确认后的结果）
h5/                                    H5：React 18 + TS + Vite
  src/ds/                              ★ H5 组件库（资产）
  src/pages/<Page>/index.tsx           页面
  src/mock/                            mock 数据
android/                               Android：Kotlin + Compose（用 Compose Desktop 在 JVM 上渲染）
  ds/src/main/kotlin/com/demo/ds/      ★ Android 组件库（资产）
  app/src/main/kotlin/com/demo/app/    页面（screens/）+ Screens.kt 注册表
  harness/                             渲染器：把 Screen 渲染成 PNG + semantics 节点 JSON
cli/                                   Go：figmafe
figma-plugin/                          Figma 导出插件（免费 Plugin API）
.claude/skills/figma-to-code/          Skill
```

## 2. H5 组件库约定（`h5/src/ds/`）

- 每个组件一个目录：`src/ds/<Name>/index.tsx` + `<Name>.module.css`（或 index.css）
- 统一出口：`src/ds/index.ts`，`export { Button } from './Button'` 和 `export type { ButtonProps } from './Button'`
- Props 定义为 `export interface <Name>Props`；变体轴用字符串字面量联合类型，例如 `type?: 'primary' | 'default' | 'text'`
- 每个组件文件顶部有 JSDoc，其中包含 `@example` 示例
- **根元素必须**：
  - 带 `data-ds="<Name>"`（用于运行时识别「这个节点是不是组件库渲染的」）
  - 透传剩余的 `data-*`、`className`、`style` 属性（页面会给它挂 `data-figma-id`）
- 只能用 `src/ds/tokens.css` 里的 CSS 变量（`--ds-color-primary` 等），不能写死颜色
- 全局：`body { margin:0; font-family: 'Noto Sans SC', sans-serif; -webkit-font-smoothing: antialiased; }`

## 3. Android 组件库约定（`android/ds/`）

- 包名 `com.demo.ds`，每个组件一个文件：`DsButton.kt` 等，`public @Composable fun DsXxx(...)`
- 变体轴用 `enum class`，例如 `enum class DsButtonType { Primary, Default, Text }`
- 每个组件函数上方有 KDoc，其中包含 `@sample` 段落（示例代码）
- 每个组件都有 `@Preview` 风格的示例函数：`@Composable fun DsButtonPreview()`（Compose Desktop 里没有 `@Preview` 注解也没关系，命名约定即可）
- **根节点必须**：
  - 接收 `modifier: Modifier = Modifier` 并应用到根节点
  - 打上组件标记：`Modifier.dsComponent("DsButton")`（自定义 SemanticsPropertyKey `DsComponentKey`，定义在 `com.demo.ds.Semantics.kt`）
- 页面给节点挂 Figma id：`Modifier.figmaId("12:34")`（= `testTag`，定义在 `Semantics.kt`）
- Tokens：`DsColors`、`DsSpacing`、`DsRadius`、`DsTypography` 对象（`Tokens.kt`）；字体 `DsFont`（从 `~/.local/share/fonts` 加载 Noto Sans SC）
- 页面里只能用 `com.demo.ds.*`；不允许直接调用 `androidx.compose.material3.*` 的组件

## 4. 页面挂载与渲染

- H5：`h5/src/main.tsx` 读取 `?page=<Page>`，渲染 `src/pages/<Page>/index.tsx` 的默认导出。视口宽度等于设计帧宽度，DPR 2
- Android：`android/app/.../Screens.kt` 里有 `val screens: Map<String, @Composable () -> Unit>`；harness 用法：
  `gradle -q -p android :harness:run --args="<Page> <outDir> <widthDp> <heightDp>"`
  输出 `<outDir>/screen.png`（@2x）和 `<outDir>/nodes.json`

`nodes.json`（两端统一格式，H5 端由 CLI 通过 DOM 生成）：

```json
{ "platform": "android", "scale": 2, "width": 375, "height": 900,
  "nodes": [ { "figmaId": "12:34", "ds": "DsButton", "text": "提醒发货",
               "x": 0, "y": 0, "w": 0, "h": 0 } ] }
```

坐标统一用 1x（dp/px），只导出有 figmaId、有 ds 标记或有文本的节点。

## 5. Figma 导出格式（`export.json`，schema `figmafe-export/1`）

节点字段与 Figma REST API `GET /v1/files/:key/nodes` 保持一致，插件导出和 REST 导出可以互换：

```json
{
  "schema": "figmafe-export/1",
  "source": { "fileKey": "…", "fileName": "…", "exportedAt": "…", "exporter": "figmafe-plugin@0.1" },
  "frame": { "id": "1:2", "name": "OrderDetail", "type": "FRAME",
             "absoluteBoundingBox": {"x":0,"y":0,"width":375,"height":900},
             "fills": [{"type":"SOLID","color":{"r":0.96,"g":0.96,"b":0.97,"a":1}}],
             "children": [ … ] },
  "components":    { "<componentId>": { "key": "…", "name": "Type=Primary, Size=Medium", "componentSetId": "<setId>" } },
  "componentSets": { "<setId>": { "key": "…", "name": "Button",
                                  "variantProperties": { "Type": ["Primary","Default","Text"], "Size": ["Large","Medium","Small"] } } },
  "images": { "frame": "frame@2x.png", "scale": 2 }
}
```

节点会出现的字段：`id name type children absoluteBoundingBox fills strokes cornerRadius paddingLeft paddingRight paddingTop paddingBottom itemSpacing layoutMode characters style{fontFamily,fontWeight,fontSize,lineHeightPx} componentId componentProperties{Name:{type:"VARIANT"|"BOOLEAN"|"TEXT",value}}`。

## 6. 关于 fixture（诚实说明）

当前没有真实 Figma 文件，也没有 Figma 网络访问（`api.figma.com` 被网络策略拒绝）。所以 `design/figma-export/order-detail/` 由 `design/source/order-detail.html` 生成：

- 这个设计源是**独立手写的绝对尺寸 HTML**，不 import 任何一端的组件库
- 用 `figmafe fixture` 渲染后，按第 5 节格式导出节点树和 @2x 图
- 拿到真实 Figma 后，用 `figma-plugin/` 导出同格式产物替换即可，后续流程不变

## 7. 存量项目：`reuse` / `match` / `compare` 的输入输出

这三条命令不依赖 demo 仓库，可以直接对任何 H5 / Android 项目运行（只读，不改项目代码）。

| 命令 | 输出目录（默认） | 文件 |
|---|---|---|
| `figmafe reuse <项目>` | `<项目>/.figmafe/reuse/` | `<h5\|android>.md`（体检报告）、`<h5\|android>.json`（`Report`） |
| `… --focus <文件/目录>` 或 `--changed <git ref>` | 同上 | 另加 `<平台>.page.json`（`PageCheck`），报告顶部是页面检查；不通过时退出码 3 |
| `… --baseline <旧 json> [--no-regress]` | 同上 | 报告顶部是指标对比；有指标变差且带 `--no-regress` 时退出码 3 |
| `figmafe match <项目> --figma <页面导出目录>` | `<项目>/.figmafe/match/` | `<平台>.md`（复用清单）、`<平台>.json`（`MatchResult`）、`<平台>.html`（叠在设计图上的可视化）、`component-map.draft.<平台>.yaml` |
| `… --figma <多个页面的父目录>` | 每个页面一个子目录 | 同上，另加 `<平台>.summary.md`（整个需求的汇总：哪些组件被多个页面用到、要新建什么） |
| `figmafe compare --url <页面> --figma <页面导出目录>` | `-o` 指定 | 截图、差异图、像素差异 / 文本位置偏差报告 |

主要 JSON 字段（完整定义见 `cli/internal/reuse/types.go`、`match.go`、`focus.go`）：

- `Report`：`platform`、`stack`、`libraries`（组件库及常用组件）、`assets`（每个组件的 `kind` / `concept` / `usage` / 真实用法 `example` + `exampleAt`；兼容性 `caveat`、`unavailable`；设计系统封装 `wrappedBy`；可配置项 `enums`、`required`、`styles`、`attrs`）、`conventions`（项目惯例写法）、`clusters`（同一概念的多个实现）、`findings`（有组件却没用 / 绕过封装）、`clones`（复制粘贴结构）、`colors`、`metrics`（组件化率等，`--baseline` 比较的就是它）。
- `MatchResult`：`items[]` 每项是一个 Figma 元素：`figma`、`nodeIds`、`variants`、`verdict`（`reuse` ≥0.8 / `confirm` / `new`）、`decision`（component-map 里登记过的决定）、`candidates[]`（`name`、`score`、`why`、`example`、`props`：Figma 变体 → 组件属性的写法）。
- `PageCheck`：`used`、`reinvented`、`missedReuse`、`handBuilt`、`inconsistent`、`colors`、`localTokens`、`pass`。

项目里 `design/component-map.yaml`、`.figmafe/component-map.yaml` 或根目录 `component-map.yaml`（按这个顺序找）里登记的映射和缺口决定会被 `match` 优先采用，也会让页面检查不再把「已决定新建」的组件判成自造轮子。格式错误时命令直接报错，不会静默忽略。
