# H5 端细节

- 工程：`h5/`（React 18 + TS + Vite）。组件库 `h5/src/ds/`，从 `@ds` 导入：`import { Button, Cell } from '@ds';`
- 页面：`h5/src/pages/<Page>/index.tsx`，默认导出页面组件；访问 `?page=<Page>`。样式可以放同目录 `index.css`。
- Figma id：组件上写 `data-figma-id="1:23"`（组件会透传到根节点）；页面级文本元素也可以挂。
- 页面级文本 / 容器：只用 `h5/src/ds/tokens.css` 的 CSS 变量，例如
  `font-size: var(--ds-font-title-size); line-height: var(--ds-font-title-line-height); font-weight: var(--ds-font-title-weight); color: var(--ds-color-text-primary);`
- 页面背景 `var(--ds-color-bg-page)`；内容区横向内边距和卡片间距看 plan.md 里容器的「间距 / 内边距」。
- `ActionBar` 是 `position: fixed` 吸底，内容区底部要留出 56px，避免被遮挡。
- verify 会自动 `npm run build`，构建失败时报告里会带出错误。
