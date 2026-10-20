# figma-fe：Figma → H5 / Android 组件复用（痛点 4 验证）

**看到新的 Figma 设计稿，先知道项目里哪些组件可以直接复用，而不是重新造轮子。**

做法：先给存量项目做一次体检，盘点出能复用的组件资产；再把新设计稿逐个元素对上去，得到复用清单；开发完检查用了哪些组件、和设计稿截图比对；把偏差反馈给 Agent 修正，最后沉淀成规则和 CI 门禁。

| 交付物 | 位置 |
|---|---|
| **CLI `figmafe`**（Go 1.25） | [`cli/`](cli)：存量项目用 `reuse / match / compare`；新页面闭环用 `scan / map / catalog / plan / audit / capture / diff / verify` |
| **Skill `figma-to-code`** | [`.claude/skills/figma-to-code/`](.claude/skills/figma-to-code/SKILL.md)，存量项目流程见 [`references/brownfield.md`](.claude/skills/figma-to-code/references/brownfield.md) |
| **Figma 导出插件**（免费 Plugin API） | [`figma-plugin/`](figma-plugin) |
| **报告**（PDF 共 10 页） | [`docs/report.pdf`](docs/report.pdf) |
| **分步教程**（每一步都有截图） | [`docs/tutorial/README.md`](docs/tutorial/README.md)：第一部分是新页面闭环，第二部分是存量项目 |
| 复用推荐准确率（26 个开源项目 × 3 份设计稿 + 精修演示稿） | [`experiments/real-match-eval.md`](experiments/real-match-eval.md)、[`experiments/match-judgement.md`](experiments/match-judgement.md) |
| Agent 对照实验（真实项目 / demo 仓库） | [`experiments/real-project/summary.md`](experiments/real-project/summary.md)、[`experiments/summary.md`](experiments/summary.md) |
| Demo：H5 / Android 组件库 + 订单详情页 | [`h5/`](h5)、[`android/`](android)、[`design/`](design) |
| CI 示例 | [`docs/ci/`](docs/ci) |
| 调研 / 方案 / 约定 | [`docs/RESEARCH.md`](docs/RESEARCH.md)、[`docs/PLAN.md`](docs/PLAN.md)、[`docs/CONTRACTS.md`](docs/CONTRACTS.md) |

## 存量项目：快速开始

```bash
make build
bin/figmafe reuse  ~/code/my-h5-app                       # 体检：组件资产、项目惯例、重复实现、不能用的库组件
bin/figmafe match  ~/code/my-h5-app --figma design/figma-export/order-detail
                                                          # 新设计稿 → 复用清单（.md / .json / 叠图 .html）
bin/figmafe reuse  ~/code/my-h5-app --changed origin/main # 写完：只检查新增代码（0=通过 3=未通过）
bin/figmafe reuse  ~/code/my-h5-app --baseline main.json --no-regress   # CI：指标只许变好
bin/figmafe compare --url http://localhost:5173/order --figma design/figma-export/order-detail -o out
```

`--figma` 可以给一个包含多个页面导出的目录，得到整个需求的汇总。Android 项目的用法相同，XML 和 Compose 都支持；新页面用哪种写法会按项目主流自动推断，也可以用 `--target xml|compose` 指定。

在 Claude Code 里直接说「按这份 Figma 在这个项目里实现订单详情」，Skill 会按「体检 → 复用清单 → 确认 → 实现 → 检查」推进。

## 新页面闭环（demo 仓库）

```bash
make fonts && make build && make test
bin/figmafe scan && bin/figmafe map suggest && bin/figmafe map check
bin/figmafe plan   --page OrderDetail --platform h5     # 复用方案（闸门）
# 页面由 Agent 按 Skill 实现到 h5/src/pages/OrderDetail；只想看 verify 效果，可以先拷一份实验里的实现：
#   mkdir -p h5/src/pages/OrderDetail && cp experiments/runs/C-h5-1/src/* h5/src/pages/OrderDetail/
bin/figmafe verify --page OrderDetail --platform h5     # 审计 + 截图对比（0=PASS 3=未达标）
```

## 说明

- **只用免费能力**：Figma Code Connect、Dev Mode MCP（`get_code_connect_map` 等）和 Variables REST API 都要付费，这里一律不用。设计数据通过自己写的插件导出（Plugin API 免费，不需要额外授权）。
- **Android 版本**：demo 代码兼容 compileSdk 32 / minSdk 23（Compose 只用 1.2 已有的 API）。存量检查会读取项目实际的组件库版本和主题，不会推荐版本不够或主题不兼容的组件。
- **环境限制**：本环境装不了 Android SDK，demo 用 Compose Desktop 在 JVM 上渲染；也没有真实 Figma 文件，设计稿是 fixture（生成方式见 `docs/CONTRACTS.md` §6）。详见报告「局限」一节。
