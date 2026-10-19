# figmafe exporter（Figma 插件）

免费获取 Figma 数据的主通道：Plugin API 不消耗 REST / MCP 的调用配额（免费档 REST 每月只有 6 次，Variables REST 只对 Enterprise 开放）。

## 安装（开发插件，免费账号可用）
Figma 桌面端 → Plugins → Development → Import plugin from manifest… → 选择本目录的 `manifest.json`。

## 使用
1. **组件库文件**里运行 → 「导出组件库」→ 得到 `library.json`、`variables.tokens.json`，放到 `design/figma-export/`。
2. **页面文件**里选中页面 Frame → 运行 → 「导出选中页面」→ 得到 `export.json`、`frame@2x.png`，放到 `design/figma-export/<page>/`，并在 `figmafe.yaml` 的 `figma.pages` 登记。

导出格式见 `docs/CONTRACTS.md` §5（字段与 Figma REST API 一致）。

> 说明：本仓库的 demo 环境无法访问 Figma（`api.figma.com` 被网络策略拒绝，也没有真实 Figma 文件），插件代码未在真实 Figma 中运行过；demo 用的是 `figmafe fixture` 生成的同格式数据。拿到真实文件后请先导出一次核对字段。
