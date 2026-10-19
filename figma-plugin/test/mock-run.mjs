// Runs the real plugin code (../code.js) against a mocked Figma Plugin API, so the exporter is
// tested without Figma: the mock document is rebuilt from a figmafe export (offset to page
// coordinates, as a real frame would be), the plugin exports it again, and the result must be
// something `figmafe match` reads the same way as the original.
//   node figma-plugin/test/mock-run.mjs <export dir> <out dir>
import { readFileSync, writeFileSync, mkdirSync } from 'node:fs';
import { join, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';
import vm from 'node:vm';

const [src, out] = process.argv.slice(2);
const exp = JSON.parse(readFileSync(join(src, 'export.json'), 'utf8'));
const png = readFileSync(join(src, exp.images.frame));
const OX = 1200, OY = -340; // the frame sits somewhere on the canvas
const mixed = Symbol('mixed');

// component sets / components as Plugin API objects
const sets = {}, comps = {};
for (const [id, s] of Object.entries(exp.componentSets || {})) {
  const defs = {};
  for (const [k, opts] of Object.entries(s.variantProperties || {})) defs[k] = { type: 'VARIANT', variantOptions: opts };
  sets[id] = { id, type: 'COMPONENT_SET', key: s.key || id, name: s.name, description: s.description || '', componentPropertyDefinitions: defs };
}
for (const [id, c] of Object.entries(exp.components || {})) {
  const parent = sets[c.componentSetId] && c.componentSetId !== id ? sets[c.componentSetId] : null;
  comps[id] = { id, type: 'COMPONENT', key: c.key || id, name: c.name || '', description: '', parent };
}

function toPlugin(n) {
  const b = n.absoluteBoundingBox || { x: 0, y: 0, width: 0, height: 0 };
  const node = {
    id: n.id, name: n.name, type: n.type, visible: true,
    absoluteBoundingBox: { x: b.x + OX, y: b.y + OY, width: b.width, height: b.height },
    width: b.width, height: b.height,
    fills: (n.fills || []).map((p) => ({ type: p.type, visible: true, color: { r: p.color.r, g: p.color.g, b: p.color.b }, opacity: p.color.a ?? 1 })),
    strokes: (n.strokes || []).map((p) => ({ type: p.type, visible: true, color: { r: p.color.r, g: p.color.g, b: p.color.b }, opacity: p.color.a ?? 1 })),
    strokeWeight: n.strokeWeight || 0,
    cornerRadius: n.cornerRadius || 0,
    layoutMode: n.layoutMode || 'NONE', itemSpacing: n.itemSpacing || 0,
    paddingLeft: n.paddingLeft || 0, paddingRight: n.paddingRight || 0, paddingTop: n.paddingTop || 0, paddingBottom: n.paddingBottom || 0,
  };
  if (n.type === 'TEXT') {
    Object.assign(node, {
      characters: n.characters,
      fontName: { family: n.style?.fontFamily || 'Inter', style: 'Regular' },
      fontSize: n.style?.fontSize || 14,
      fontWeight: n.style?.fontWeight || 400,
      lineHeight: n.style?.lineHeightPx ? { unit: 'PIXELS', value: n.style.lineHeightPx } : { unit: 'AUTO' },
    });
  }
  if (n.type === 'INSTANCE') {
    const props = {};
    for (const [k, v] of Object.entries(n.componentProperties || {})) props[v.type === 'VARIANT' ? k : `${k}#1:0`] = v;
    node.componentProperties = props;
    node.getMainComponentAsync = async () => comps[n.componentId] || null;
  }
  if (n.children) node.children = n.children.map(toPlugin);
  return node;
}

const frame = toPlugin(exp.frame);
frame.exportAsync = async () => new Uint8Array(png);
let posted;
const figma = {
  mixed,
  showUI() {}, notify(m) { throw new Error('notify: ' + m); },
  fileKey: 'mock', root: { name: 'mock.fig' },
  currentPage: { selection: [frame] },
  ui: { postMessage(m) { posted = m; } },
};
const code = readFileSync(join(dirname(fileURLToPath(import.meta.url)), '..', 'code.js'), 'utf8');
const ctx = vm.createContext({ figma, __html__: '', console, Object, Array, JSON, Math, Date, Symbol });
vm.runInContext(code, ctx);
await figma.ui.onmessage({ type: 'page' });
if (!posted) throw new Error('plugin posted nothing');
mkdirSync(out, { recursive: true });
for (const f of posted.files) writeFileSync(join(out, f.name), f.text ?? Buffer.from(f.bytes));
console.log('plugin exported', posted.files.map((f) => f.name).join(', '), '→', out);
