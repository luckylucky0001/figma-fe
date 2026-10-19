# 教程：用 figma-to-code Skill + figmafe，把 Figma 页面复用组件还原到 H5 / Android

> 本教程的每一张截图都是真实运行产生的，可以用 `docs/tutorial/capture.sh` 重新生成。
> 示例页面：订单详情（`OrderDetail`，375×900）。两端组件库：H5 `h5/src/ds`（React），Android `android/ds`（Compose）和 `android/ds-view`（XML View 存量组件）。

**整体流程**

```
Figma ──导出插件──▶ design/figma-export        (Step 2)
组件库 ──scan─────▶ 资产清单 inventory          (Step 3  Phase 0)
       ──map──────▶ component-map + 缺口/漂移报告 (Step 4  Phase 1)
       ──catalog──▶ Skill 组件手册               (Step 5)
页面   ──plan─────▶ 复用方案（闸门：人确认）      (Step 6  Phase 2)
       ──Agent 按 Skill 实现──                     (Step 7  Phase 3)
       ──verify───▶ 审计 + 截图对比 → 修正 → PASS (Step 8-10 Phase 4)
       ──rules.md─▶ 沉淀规则                     (Step 11 Phase 5)
```

---

## Step 1 环境准备

需要 Go 1.22+、Node 18+、Chromium（Playwright 自带的即可）、Gradle 8+（JDK 17+），以及 Noto Sans SC 字体（`make fonts`）。字体必须装上：设计图、H5 截图、Android 渲染都用这一套字体，像素对比才有意义。

```bash
make fonts      # 安装 Noto Sans SC
make build      # 编译 CLI → bin/figmafe
make test       # 单元测试
```

![构建与测试](img/01-build.png)

所有路径和阈值都在仓库根目录的 `figmafe.yaml` 里配置：组件库位置、页面目录、构建命令、像素阈值等。

---

## Step 2 从 Figma 导出设计数据（免费）

**不用 Figma MCP 和 Code Connect（都要付费），也不调用 REST API（免费档每月 6 次）。** 用仓库里的免费插件 `figma-plugin/`：

1. Figma 桌面端 → Plugins → Development → Import plugin from manifest… → 选 `figma-plugin/manifest.json`
2. 在**组件库文件**里运行 →「导出组件库」→ 得到 `library.json` 和 `variables.tokens.json`（DTCG 格式），放进 `design/figma-export/`
3. 在**页面文件**里选中页面 Frame → 运行 →「导出选中页面」→ 得到 `export.json` 和 `frame@2x.png`，放进 `design/figma-export/order-detail/`
4. 在 `figmafe.yaml` 的 `figma.pages` 里登记页面

> 权限：插件不需要 token 或付费席位，但必须在**有编辑权限**的文件里运行。只有查看权限的话，先「Duplicate to your drafts」复制一份，在副本里导出。

设计稿（`frame@2x.png`）：

<img src="img/02-design.png" width="300">

`export.json` 的字段和 Figma REST API 一致。组件实例带有 `componentId` 和变体属性，这是后面识别复用的依据：

![导出产物](img/02-export.png)

> 本 demo 环境访问不了 Figma，示例数据由 `figmafe fixture` 根据 `design/source/order-detail.html` 生成，格式和插件导出完全一致。

---

## Step 3 Phase 0：盘点组件资产（先盘点，再谈复用）

```bash
bin/figmafe scan
```

扫描两端组件库，输出资产清单（`.figmafe/inventory-*.json`），包括组件、props / 变体轴、示例、业务代码里的使用次数、token，以及存量代码里的坏实践（基线）：
- **H5**：解析 `src/ds/index.ts` 的导出，以及每个组件的 `Props` 接口和 JSDoc `@example`。组件库是 npm 包时，改为扫描它的 `.d.ts` 类型声明。
- **Android Compose**：解析 `@Composable fun Ds*` 的参数、`enum class` 变体和 KDoc `@sample`。
- **Android XML View**：解析 View 子类、`attrs.xml` 的 `declare-styleable`、layout 里的使用次数，以及原生控件、写死颜色这类存量问题。

![资产扫描](img/03-scan.png)

---

## Step 4 Phase 1：建立 Figma ↔ 代码映射

```bash
bin/figmafe map suggest     # 按名称（含同义词）+ 变体结构打分，给出候选
```

![映射建议](img/04-map-suggest.png)

看图上的几个例子：
- ✓ 高置信，直接采用：`NavBar → DsTopBar`、`Cell → DsListItem` 这类名字不同的，靠同义词表匹配上。
- 一个 Figma 组件同时有两套资产（`DsButton / LegacyButton`）时，由人来决定。本例的决定是：新页面用 Compose，Legacy* 只给存量 XML 页面用。
- ✗ 没有候选：`CountDown`，属于缺口。

人工确认后，写进 `design/component-map.yaml`。每个 Figma 组件写清楚两端各自的实现：组件、import、变体值对应的 props、文本图层对应的属性、使用示例和约束。然后检查映射：

```bash
bin/figmafe map check
```

![映射检查](img/05-map-check.png)

这份报告本身就是很有价值的产出，图中能看到：
- **CountDown 两端都没有**：已登记处理方案。
- **Android 的 DsBadge 不支持 Solid 变体**。
- **token 漂移**：Android `color.warning` 是 #FF8800，Figma 是 #FF8F1F。

这些问题全部记录在 `docs/component-issues.md`，交给组件库 owner。

---

## Step 5 生成 Skill 的组件手册

```bash
bin/figmafe catalog
```

![catalog](img/06-catalog.png)

`references/catalog-h5.md` 和 `catalog-android.md` 是给 Agent 看的组件手册，借鉴的是 Mística 在 npm 包里附带 `llms.md` 的做法。Agent 读这份手册，就不需要自己去翻源码：

![组件手册](img/06-catalog-h5.png)

---

## Step 6 Phase 2：生成复用方案（闸门）

```bash
bin/figmafe plan --page OrderDetail --platform h5
bin/figmafe plan --page OrderDetail --platform android
```

方案逐个 Figma 节点给出：用哪个组件、可以直接粘贴的代码（变体 → props、文本 → 属性、`data-figma-id`）、页面级文本对应的 token、待决项和约束。**必须经人确认后才能写代码。**

![H5 复用方案](img/07-plan-h5.png)

同一份 Figma，Android 端会自动换成 Android 的组件和 token：

![Android 复用方案](img/07-plan-android.png)

---

## Step 7 Phase 3：让 Agent 按 Skill 实现

Skill 放在 `.claude/skills/figma-to-code/`。在 Claude Code 里这样发起任务，Skill 会自动触发：

```
按 Figma 实现订单详情页（design/figma-export/order-detail），H5 端，页面放 h5/src/pages/OrderDetail。
```

也可以直接输入 `/figma-to-code`。Skill 规定的流程和硬性规则如下：

![SKILL.md](img/08-skill.png)

---

## Step 8 Phase 4：verify 第 1 轮，抓出偏差

为了演示，第一版实现里故意放了 3 个常见问题：自己拼状态标签、用原生 `<button>`、自己用 div 拼左右结构的行。

![问题代码](img/09-bad-code.png)

```bash
bin/figmafe verify --page OrderDetail --platform h5
```

verify 一次跑完四件事：静态审计、构建、截图（Chromium，375×900，2 倍图）、像素和节点级对比。没达标时退出码为 3：

![verify FAIL](img/09-verify-fail.png)

报告列出了每一处问题，静态问题带文件和行号，节点偏差带 Figma 节点 id，Agent 可以照着逐条修改：

![报告](img/09-report-fail.png)

差异图（红色为不一致的像素）：

![第 1 轮对比](img/09-compare.png)

---

## Step 9 按报告修正 → PASS

把 3 处问题改回组件库组件后，再跑一次 verify：

![verify PASS](img/10-verify-pass.png)

![PASS 报告](img/10-report-pass.png)

![第 2 轮对比](img/10-compare.png)

剩下的轻微偏差（ProductCard 内部下移约 4px、按钮宽 2px）出在**组件库内部**。按规则，页面不做补偿，在交付回复里列出，并登记到 `docs/component-issues.md`。

---

## Step 10 Android 同样的流程

页面放在 `android/app/.../screens/OrderDetail/`，并在 `Screens.kt` 里注册。verify 会自动执行 `gradle :harness:installDist`，在 JVM 上渲染 Compose 页面（Compose Desktop，API 和 Android 一致），再导出 semantics 树和 2 倍图：

![Android verify](img/11-verify-android.png)

![Android 报告](img/11-report-android.png)

![Android 对比](img/11-compare.png)

Android 端剩下的 2 处严重偏差（DsPriceText 高度、DsGoodsCard 规格行）同样是组件库的问题，已经登记到组件问题清单。这正是 verify 的价值之一：**它能把「页面写错了」和「组件本身不符合设计」区分开。**

> 真实项目接入时，把渲染器换成 Roborazzi 或 Paparazzi（两者都支持 Compose 和 View），在 CI 里输出同格式的 `screen.png` 和 `nodes.json` 即可，CLI 的对比逻辑不用改。

---

## Step 11 Phase 5：沉淀规则

每做完一个页面，Agent 把「下次还会遇到的问题」追加到 `references/rules.md`，再由人整理合并。下一个页面的 Agent 会先读到这些规则：

![rules](img/12-rules.png)

---

## 接入你们自己的项目（清单）

1. 用插件导出组件库，得到 `library.json` 和 `variables.tokens.json`；导出第一个典型页面，得到 `export.json` 和 `frame@2x.png`
2. 修改 `figmafe.yaml`：H5 组件库路径（源码目录，或 npm 包的 `dist/types`）、Android Compose 组件库和 XML View 组件库的路径、页面目录、构建和渲染命令
3. `figmafe scan` → `map suggest` → 人工确认后写 `component-map.yaml` → `map check`，把缺口和漂移交给组件库 owner
4. `figmafe catalog`，然后把 `.claude/skills/figma-to-code/` 提交进仓库
5. 先选一个典型页面跑通 plan → 实现 → verify，校准阈值（`thresholds`）
6. 把 `figmafe verify` 接入 PR CI（退出码不为 0 时阻塞合并），让 `rules.md` 随着页面增多持续积累
