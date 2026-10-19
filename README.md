# figma-fe：Figma → H5 / Android 组件复用与还原度（痛点 4 验证 Demo）

先盘点组件资产，再识别复用；开发前出方案并确认，开发后做审计和截图对比，把偏差反馈给 Agent 修正，最后把经验沉淀成规则。

| 交付物 | 位置 |
|---|---|
| **CLI `figmafe`**（Go） | [`cli/`](cli) —— `scan / map / catalog / plan / audit / capture / diff / verify` |
| **Skill `figma-to-code`** | [`.claude/skills/figma-to-code/`](.claude/skills/figma-to-code/SKILL.md) |
| **Figma 导出插件**（免费 Plugin API） | [`figma-plugin/`](figma-plugin) |
| **Demo**：H5 组件库 + Android 组件库 + 订单详情页 | [`h5/`](h5)、[`android/`](android)、[`design/`](design) |
| **对照实验**（13 次运行） | [`experiments/`](experiments/summary.md) |
| **报告**（PDF 共 3 页） | [`docs/report.pdf`](docs/report.pdf)、[`docs/REPORT.md`](docs/REPORT.md) |
| **分步教程**（每一步都有截图） | [`docs/tutorial/README.md`](docs/tutorial/README.md) |
| 调研 / 方案 / 约定 | [`docs/RESEARCH.md`](docs/RESEARCH.md)、[`docs/PLAN.md`](docs/PLAN.md)、[`docs/CONTRACTS.md`](docs/CONTRACTS.md) |
| 组件库待修清单 | [`docs/component-issues.md`](docs/component-issues.md) |

## 快速开始

```bash
make fonts && make build && make test
bin/figmafe scan                                        # Phase 0 盘点两端组件资产
bin/figmafe map suggest && bin/figmafe map check        # Phase 1 映射 + 缺口 / 漂移
bin/figmafe catalog                                     # 生成 Skill 组件手册
bin/figmafe plan   --page OrderDetail --platform h5     # Phase 2 复用方案（闸门）
bin/figmafe verify --page OrderDetail --platform h5     # Phase 4 审计 + 截图对比（0=PASS 3=未达标）
```

在 Claude Code 里直接说「按 Figma 实现订单详情页」，或者输入 `/figma-to-code`，Skill 就会按上面的流程推进。

## 说明

- **只用免费能力**：Figma 官方 MCP、Code Connect 和 Variables REST API 都需要付费，这里一律不用。设计数据通过自己写的插件导出。
- **Android 渲染方式**：本环境访问不了 `dl.google.com`，装不了 Android SDK，所以用 Compose Desktop 在 JVM 上渲染，API 和 Android 的 Compose 一致；代码只用 Compose 1.2 已有的 API，以兼容 compileSdk 32。在真实项目里，把截图这一步换成 Roborazzi 或 Paparazzi 即可。
- **设计数据是 fixture**：本环境没有真实的 Figma 文件，demo 用的是 fixture，生成方式见 `docs/CONTRACTS.md` §6。
