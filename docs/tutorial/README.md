# 教程：用 figma-to-code Skill + figmafe，把 Figma 页面复用组件还原到 H5 / Android

> 本教程的每一张截图都是真实运行产生的，可以用 `docs/tutorial/capture.sh` 重新生成。
>
> - **第一部分（Step 1-11）**：已经接入的项目怎么用：组件映射、复用方案、截图逐节点比对。示例是本仓库的 demo（H5 `h5/src/ds`，Android `android/ds` / `android/ds-view`）。
> - **第二部分（Step 12-19）**：**任何存量项目**怎么用，不需要先接入：先体检，拿到新设计稿后自动给出「哪些组件可以直接复用、怎么配置」，实现后检查有没有自造轮子，接进 PR 闸门。示例是真实开源项目（newbee-mall、QMUI_Android、Now in Android）。只关心存量项目可以直接跳到 [第二部分](#第二部分存量项目不需要先接入)。

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

需要 Go 1.25+、Node 18+、Chromium（Playwright 自带的即可）、Gradle 8+（JDK 17+），以及 Noto Sans SC 字体（`make fonts`）。字体必须装上：设计图、H5 截图、Android 渲染都用这一套字体，像素对比才有意义。

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

# 第二部分：存量项目（不需要先接入）

存量项目的问题是：组件散在各处，有自己写的、有组件库的、有被封装过的，新人和 Agent 都不知道「已经有什么」，于是看着设计稿重新造。第二部分的工具不需要 `figmafe.yaml`、不需要 component-map，直接对着项目目录跑。支持：

- **H5**：Vue 2 / 3、React（含写在 .js 里的 JSX）、uni-app、Taro；组件库 Vant（2/3/4）、React Vant、NutUI、antd-mobile（v2/v5）、TDesign Mobile、uView、uni-ui、taro-ui、Element Plus，以及单组件的 npm 包（如 vue-countdown）
- **Android**：XML View（layout、自定义 View、`attrs.xml`、`<include>`、style）、Jetpack Compose；组件库 Material Components、Compose Material3 / Foundation、QMUI，以及项目自己的设计系统模块

## Step 12 存量体检：项目里已经有什么

```bash
bin/figmafe reuse <项目目录>          # 报告在 <项目>/.figmafe/reuse/<platform>.md（.figmafe 会自动被 git 忽略）
```

![newbee 体检](img/20-reuse-newbee.png)

报告第 0 节是「优先处理」：按影响面排好序的几件事，团队可以从这里开始治理：

![体检报告](img/20-reuse-report.png)

同一个工具对 Android 也一样。QMUI_Android 里，工具从依赖和主题推断出：项目用的是 Material 1.4.0、App 主题继承 AppCompat，所以 `MaterialDivider`（1.5 才有）、`Chip`（需要 MaterialComponents 主题，否则运行时崩溃）这些组件虽然在组件库里，**但这个项目不能直接用**，后面的推荐会自动排除它们：

![QMUI 体检](img/21-reuse-qmui.png)

![QMUI 不能直接用的组件](img/21-qmui-compat.png)

Now in Android 有自己的设计系统（`core/designsystem`），工具找出了它封装了哪些 Material3 组件。之后新代码直接用 `Button` 而不是 `NiaButton`，会被判为「绕过封装」：

![Now in Android 体检](img/22-reuse-nia.png)

![设计系统封装](img/22-nia-wrappers.png)

## Step 13 新设计稿 → 复用清单（这一步回答「哪些可以直接复用」）

```bash
bin/figmafe match --figma <Figma 导出目录> <项目目录>
```

![match](img/23-match-newbee.png)

同时生成一张叠在设计图上的可视版（`.html`），绿色 = 直接复用，橙色 = 需确认，红色 = 需新建。方案闸门时直接把这一页发给人看：

![复用清单叠图](img/23-match-newbee-overlay.png)

清单不只说「用哪个组件」，还说「怎么配置」：设计稿的变体映射到组件的属性，以及项目自己的惯用写法（比如 newbee 的按钮都是 `block color="#1baeae"`）：

![变体到属性](img/23-match-props.png)

一个需求往往有好几个页面。`--figma` 指向放着多个页面导出的文件夹，会多出一份整个需求的复用清单：

![需求汇总](img/23b-feature-summary.png)

## Step 14 同一份设计稿，不同的项目给出不同的答案

推荐来自**这个项目**里的东西，而不是一张通用对照表。同一份设计稿放到 QMUI（XML）上：导航栏是 `QMUITopBarLayout`、卡片是 `QMUIGroupListView`、单元格是 `QMUICommonListItemView`、按钮是项目里用了 29 次的 `QMUIRoundButton`（附带项目惯用的 `style="@style/QDRoundButtonStyle"`）；而 Chip、MaterialDivider 因为版本 / 主题不兼容，不会被推荐。没有现成组件的分割线，推荐的是 QMUI 容器自带的分割线属性：

![QMUI](img/24-match-qmui.png)

![QMUI 叠图](img/24-match-qmui-overlay.png)

## Step 15 没组件化、图层名不规范的设计稿也能用

真实设计稿经常没有组件化，图层叫「导航栏」「涨跌标签」「收盘倒计时」。工具按图层名推断 UI 概念，结果标「（未组件化）」：

![持仓页](img/25-match-holdings.png)

![持仓页叠图](img/25-match-holdings-overlay.png)

英文图层名、前两份设计稿没有的组件（搜索、步进器、开关、复选框、空状态），放到 Compose 项目上：

![买入下单](img/25b-match-trade.png)

![买入下单叠图](img/25b-match-trade-overlay.png)

## Step 16 让 Agent 按清单实现，再检查有没有自造轮子

Skill 的存量流程（`references/brownfield.md`）要求 Agent：先体检 → 拿复用清单 → 「直接复用」的必须用、「需确认」的写理由、只有「需新建」的才能新建（放进组件目录）→ 实现完跑页面复用检查：

```bash
bin/figmafe reuse <项目目录> --focus <新页面文件>      # 退出码 3 = 不通过
```

同一个任务，在 newbee 上的检查结果。上面是不用 Skill 的 Agent（A 组），下面是用 Skill 的 Agent（C 组）：

![A 组](img/26-check-a.png)

![C 组](img/26-check-c.png)

不用 Skill 时，Agent 手写了卡片、分割线、商品卡、底部操作栏，还自己定义了 10 个颜色变量；用 Skill 时这些都换成了项目已有的组件，缺的价格组件放进了 `src/components`。

## Step 17 PR 闸门：只检查这次改动

```bash
bin/figmafe reuse . --changed origin/main      # 新文件整份检查，修改过的文件只看新增的行
```

有意不用推荐组件、或者新建了和已有组件同类的组件时，把决定和理由写进 `component-map.yaml` 的 `gaps`，闸门会认可这个决定，不会每次都报：

![PR 闸门](img/27-gate.png)

GitHub Actions / GitLab CI 的配置示例在 `docs/ci/`。

存量的问题一次改不完，可以先保证「不变差」：

![baseline](img/27b-baseline.png)

## Step 18 页面能跑起来时：截图比对

存量项目没有挂 Figma 节点 id，做不了第一部分那样的逐节点比对，但可以对任意一个能打开的页面做整页比对和文字覆盖检查：

```bash
bin/figmafe compare --url http://localhost:5173/#/order/detail --figma <Figma 导出目录>
```

![compare](img/27c-compare.png)

## Step 19 效果：准确率和对照实验

**复用清单的准确率**：18 个真实开源项目（其中 13 个在调规则时没用过，第一次跑出的错误都如实记录在 `experiments/match-judgement.md`）× 3 份设计稿，每一项人工判断推荐得对不对：

![准确率](img/28-eval.png)

**对照实验**：在 5 个真实项目上让 Agent 实现同一个页面，A 组只给任务，C 组用 Skill，统一用页面复用检查打分：

![汇总](img/29-summary.png)

---

## 接入你们自己的项目（清单）

**存量项目，最快的用法**（不需要任何配置）：`figmafe reuse` 体检 → 把 `.claude/skills/figma-to-code/` 放进仓库 → 新需求先跑 `figmafe match` 出复用清单 → 把 `figmafe reuse --changed origin/main` 接进 PR CI（`docs/ci/`）。

**想要逐节点的截图比对**（第一部分的完整流程）：

1. 用插件导出组件库，得到 `library.json` 和 `variables.tokens.json`；导出第一个典型页面，得到 `export.json` 和 `frame@2x.png`
2. 修改 `figmafe.yaml`：H5 组件库路径（源码目录，或 npm 包的 `dist/types`）、Android Compose 组件库和 XML View 组件库的路径、页面目录、构建和渲染命令
3. `figmafe scan` → `map suggest` → 人工确认后写 `component-map.yaml` → `map check`，把缺口和漂移交给组件库 owner
4. `figmafe catalog`，然后把 `.claude/skills/figma-to-code/` 提交进仓库
5. 先选一个典型页面跑通 plan → 实现 → verify，校准阈值（`thresholds`）
6. 把 `figmafe verify` 接入 PR CI（退出码不为 0 时阻塞合并），让 `rules.md` 随着页面增多持续积累
