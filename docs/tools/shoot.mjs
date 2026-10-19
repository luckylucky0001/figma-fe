// 文档截图 / PDF 工具（教程和报告用）。依赖 h5/node_modules 里的 playwright-core + 预装 Chromium。
//   node docs/tools/shoot.mjs term <title> <in.txt> <out.png>      终端输出 → 截图
//   node docs/tools/shoot.mjs md   <title> <in.md>  <out.png>      Markdown（pandoc 转 HTML）→ 截图
//   node docs/tools/shoot.mjs html <in.html> <out.png|out.pdf>     HTML → 截图或 PDF（A4）
import { createRequire } from 'node:module';
import { readFileSync, writeFileSync, mkdtempSync, existsSync, readdirSync } from 'node:fs';
import { execFileSync } from 'node:child_process';
import { join, resolve, dirname } from 'node:path';
import { tmpdir } from 'node:os';
import { fileURLToPath, pathToFileURL } from 'node:url';

const here = dirname(fileURLToPath(import.meta.url));
const require = createRequire(join(here, '../../h5/package.json'));
const { chromium } = require('playwright-core');

function chromePath() {
  const root = process.env.PLAYWRIGHT_BROWSERS_PATH || '/opt/pw-browsers';
  for (const d of readdirSync(root)) {
    const p = join(root, d, 'chrome-linux', 'chrome');
    if (d.startsWith('chromium-') && existsSync(p)) return p;
  }
  return undefined;
}

const esc = (s) => s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');

const base = `
  :root { --bg:#0f1419; --fg:#e6e1cf; --mut:#8a93a3; --ok:#7fd962; --err:#ff6a6a; --warn:#ffb454; }
  * { box-sizing: border-box; }
  body { margin: 0; padding: 24px; background: #eef0f3; font-family: 'Noto Sans SC', sans-serif; }`;

function termHTML(title, text) {
  const body = esc(text)
    .replace(/^(\$ .*)$/gm, '<span class="cmd">$1</span>')
    .replace(/(✓[^\n]*|PASS[^\n]*)/g, '<span class="ok">$1</span>')
    .replace(/(❌[^\n]*|✗[^\n]*|FAIL[^\n]*)/g, '<span class="err">$1</span>')
    .replace(/(⚠[^\n]*|\? [A-Z][^\n]*)/g, '<span class="warn">$1</span>');
  return `<!doctype html><meta charset="utf-8"><style>${base}
  .win { width: 1040px; background: var(--bg); border-radius: 10px; box-shadow: 0 8px 28px rgba(0,0,0,.25); overflow: hidden; }
  .bar { height: 34px; background: #1d2430; display: flex; align-items: center; gap: 8px; padding: 0 14px; color: var(--mut); font-size: 13px; }
  .dot { width: 12px; height: 12px; border-radius: 50%; } .r{background:#ff5f57}.y{background:#febc2e}.g{background:#28c840}
  pre { margin: 0; padding: 16px 18px; color: var(--fg); font: 13px/1.55 'DejaVu Sans Mono', 'Noto Sans SC', monospace; white-space: pre-wrap; word-break: break-all; }
  .cmd { color: #73d0ff; } .ok { color: var(--ok); } .err { color: var(--err); } .warn { color: var(--warn); }
  </style><div class="win"><div class="bar"><span class="dot r"></span><span class="dot y"></span><span class="dot g"></span>&nbsp;${esc(title)}</div><pre>${body}</pre></div>`;
}

function mdHTML(title, md) {
  const html = execFileSync('pandoc', ['-f', 'gfm', '-t', 'html'], { input: md }).toString();
  return `<!doctype html><meta charset="utf-8"><style>${base}
  .doc { width: 1040px; background: #fff; border-radius: 10px; box-shadow: 0 8px 28px rgba(0,0,0,.15); overflow: hidden; }
  .bar { height: 34px; background: #f6f8fa; border-bottom: 1px solid #e5e7eb; display: flex; align-items: center; padding: 0 14px; color: #57606a; font-size: 13px; }
  .md { padding: 8px 28px 20px; color: #1f2329; font-size: 14px; line-height: 1.6; }
  .md h1 { font-size: 22px; border-bottom: 1px solid #eee; padding-bottom: 6px; } .md h2 { font-size: 17px; margin-top: 18px; } .md h3 { font-size: 15px; }
  .md code { background: #f3f4f6; padding: 1px 5px; border-radius: 4px; font: 12.5px 'DejaVu Sans Mono', 'Noto Sans SC', monospace; }
  .md pre { background: #f6f8fa; padding: 10px 12px; border-radius: 6px; overflow: hidden; white-space: pre-wrap; } .md pre code { background: none; padding: 0; }
  .md table { border-collapse: collapse; margin: 8px 0; font-size: 13px; } .md td, .md th { border: 1px solid #e5e7eb; padding: 4px 10px; } .md th { background: #f6f8fa; }
  .md blockquote { margin: 8px 0; padding: 4px 12px; border-left: 3px solid #d0d7de; color: #57606a; }
  .md li { margin: 2px 0; }
  </style><div class="doc"><div class="bar">📄 ${esc(title)}</div><div class="md">${html}</div></div>`;
}

async function render(html, out, { pdf = false } = {}) {
  const browser = await chromium.launch({ executablePath: chromePath() });
  const page = await browser.newPage({ deviceScaleFactor: pdf ? 1 : 2, viewport: { width: 1100, height: 800 } });
  const dir = mkdtempSync(join(tmpdir(), 'shoot-'));
  const file = join(dir, 'page.html');
  writeFileSync(file, html);
  await page.goto(pathToFileURL(file).href);
  await page.evaluate(() => document.fonts.ready);
  if (pdf) {
    await page.pdf({ path: out, format: 'A4', printBackground: true, margin: { top: '0', bottom: '0', left: '0', right: '0' } });
  } else {
    const el = await page.$('.win, .doc, body > *');
    await (el || page).screenshot({ path: out });
  }
  await browser.close();
}

const [mode, ...args] = process.argv.slice(2);
if (mode === 'term') await render(termHTML(args[0], readFileSync(args[1], 'utf8')), resolve(args[2]));
else if (mode === 'md') await render(mdHTML(args[0], readFileSync(args[1], 'utf8')), resolve(args[2]));
else if (mode === 'html') {
  const html = readFileSync(args[0], 'utf8').replace('<head>', `<head><base href="${pathToFileURL(resolve(dirname(args[0]))).href}/">`);
  await render(html, resolve(args[1]), { pdf: args[1].endsWith('.pdf') });
} else {
  console.error('usage: shoot.mjs term|md|html …');
  process.exit(2);
}
