// figmafe exporter — 免费的 Figma 插件（Plugin API，没有 REST 调用配额限制）。
// 用法：在 Figma 桌面端 Plugins → Development → Import plugin from manifest… 选择本目录的 manifest.json；
//   1) 在页面文件里选中一个 Frame → 运行插件 → 「导出页面」：下载 <page>.zip（export.json + frame@2x.png）
//   2) 在组件库文件里运行 → 「导出组件库」：下载 library.json + variables.tokens.json（DTCG）
// 导出格式与 docs/CONTRACTS.md §5 一致，字段沿用 Figma REST API，figmafe 直接读取。

figma.showUI(__html__, { width: 320, height: 220 });

const round = (v) => Math.round(v * 100) / 100;

function paints(list) {
  if (!Array.isArray(list)) return undefined;
  const out = list
    .filter((p) => p.visible !== false && p.type === 'SOLID')
    .map((p) => ({ type: 'SOLID', color: { r: p.color.r, g: p.color.g, b: p.color.b, a: p.opacity == null ? 1 : p.opacity } }));
  return out.length ? out : undefined;
}

async function serialize(node, frameBox, components, componentSets) {
  const b = node.absoluteBoundingBox || { x: 0, y: 0, width: node.width, height: node.height };
  const n = {
    id: node.id,
    name: node.name,
    type: node.type,
    absoluteBoundingBox: { x: round(b.x - frameBox.x), y: round(b.y - frameBox.y), width: round(b.width), height: round(b.height) },
  };
  if ('fills' in node && node.fills !== figma.mixed) n.fills = paints(node.fills);
  if ('strokes' in node) { n.strokes = paints(node.strokes); if (n.strokes) n.strokeWeight = node.strokeWeight; }
  if ('cornerRadius' in node && typeof node.cornerRadius === 'number' && node.cornerRadius) n.cornerRadius = node.cornerRadius;
  if ('layoutMode' in node && node.layoutMode !== 'NONE') {
    n.layoutMode = node.layoutMode;
    n.itemSpacing = node.itemSpacing;
    ['paddingLeft', 'paddingRight', 'paddingTop', 'paddingBottom'].forEach((k) => { if (node[k]) n[k] = node[k]; });
  }
  if (node.type === 'TEXT') {
    n.characters = node.characters;
    const fontName = node.fontName === figma.mixed ? node.getRangeFontName(0, 1) : node.fontName;
    const size = node.fontSize === figma.mixed ? node.getRangeFontSize(0, 1) : node.fontSize;
    const lh = node.lineHeight === figma.mixed ? node.getRangeLineHeight(0, 1) : node.lineHeight;
    const weight = node.fontWeight === figma.mixed ? node.getRangeFontWeight(0, 1) : node.fontWeight;
    n.style = {
      fontFamily: fontName.family,
      fontWeight: weight,
      fontSize: size,
      lineHeightPx: lh.unit === 'PIXELS' ? lh.value : lh.unit === 'PERCENT' ? (size * lh.value) / 100 : size * 1.2,
    };
  }
  if (node.type === 'INSTANCE') {
    const main = await node.getMainComponentAsync();
    if (main) {
      n.componentId = main.id;
      const set = main.parent && main.parent.type === 'COMPONENT_SET' ? main.parent : null;
      components[main.id] = { key: main.key, name: main.name, componentSetId: set ? set.id : main.id };
      if (set) componentSets[set.id] = setInfo(set);
      else componentSets[main.id] = { key: main.key, name: main.name, description: main.description, variantProperties: {} };
    }
    const props = {};
    for (const [k, v] of Object.entries(node.componentProperties || {})) {
      props[k.split('#')[0]] = { type: v.type, value: v.value };
    }
    if (Object.keys(props).length) n.componentProperties = props;
  }
  if ('children' in node && node.visible !== false) {
    const kids = [];
    for (const c of node.children) if (c.visible !== false) kids.push(await serialize(c, frameBox, components, componentSets));
    if (kids.length) n.children = kids;
  }
  return n;
}

function setInfo(set) {
  const vp = {};
  for (const [k, def] of Object.entries(set.componentPropertyDefinitions || {})) {
    if (def.type === 'VARIANT') vp[k] = def.variantOptions;
  }
  return { key: set.key, name: set.name, description: set.description, variantProperties: vp };
}

function hex(c) {
  const h = (v) => Math.round(v * 255).toString(16).padStart(2, '0').toUpperCase();
  return '#' + h(c.r) + h(c.g) + h(c.b);
}

async function exportVariables() {
  // Figma Variables → DTCG 2025.10（Plugin API 免费可读；REST 的 Variables API 需要 Enterprise）
  const out = {};
  const collections = await figma.variables.getLocalVariableCollectionsAsync();
  for (const col of collections) {
    const mode = col.defaultModeId;
    for (const id of col.variableIds) {
      const v = await figma.variables.getVariableByIdAsync(id);
      const parts = v.name.split('/').map((s) => s.trim().toLowerCase().replace(/\s+/g, '-'));
      const val = v.valuesByMode[mode];
      let token;
      if (v.resolvedType === 'COLOR' && val && val.r != null) token = { $type: 'color', $value: hex(val) };
      else if (v.resolvedType === 'FLOAT') token = { $type: 'dimension', $value: { value: val, unit: 'px' } };
      else continue;
      let cur = out;
      for (const p of parts.slice(0, -1)) cur = cur[p] = cur[p] || {};
      cur[parts[parts.length - 1]] = token;
    }
  }
  return out;
}

figma.ui.onmessage = async (msg) => {
  try {
    if (msg.type === 'page') {
      const frame = figma.currentPage.selection[0];
      if (!frame || (frame.type !== 'FRAME' && frame.type !== 'COMPONENT')) {
        figma.notify('请先选中一个页面 Frame');
        return;
      }
      const box = frame.absoluteBoundingBox;
      const components = {}, componentSets = {};
      const tree = await serialize(frame, box, components, componentSets);
      tree.absoluteBoundingBox = { x: 0, y: 0, width: round(box.width), height: round(box.height) };
      const exp = {
        schema: 'figmafe-export/1',
        source: { fileKey: figma.fileKey || '', fileName: figma.root.name, exportedAt: new Date().toISOString(), exporter: 'figmafe-plugin@0.1' },
        frame: tree, components, componentSets,
        images: { frame: 'frame@2x.png', scale: 2 },
      };
      const png = await frame.exportAsync({ format: 'PNG', constraint: { type: 'SCALE', value: 2 } });
      figma.ui.postMessage({ type: 'download', files: [
        { name: 'export.json', text: JSON.stringify(exp, null, 2) },
        { name: 'frame@2x.png', bytes: Array.from(png) },
      ], zip: frame.name.replace(/[^\w-]+/g, '-') });
    } else if (msg.type === 'library') {
      await figma.loadAllPagesAsync();
      const sets = {};
      for (const s of figma.root.findAllWithCriteria({ types: ['COMPONENT_SET'] })) sets[s.id] = setInfo(s);
      for (const c of figma.root.findAllWithCriteria({ types: ['COMPONENT'] })) {
        if (c.parent && c.parent.type === 'COMPONENT_SET') continue;
        sets[c.id] = { key: c.key, name: c.name, description: c.description, variantProperties: {} };
      }
      const lib = { schema: 'figmafe-library/1', fileName: figma.root.name, componentSets: sets };
      figma.ui.postMessage({ type: 'download', files: [
        { name: 'library.json', text: JSON.stringify(lib, null, 2) },
        { name: 'variables.tokens.json', text: JSON.stringify(await exportVariables(), null, 2) },
      ] });
    }
  } catch (e) {
    figma.notify('导出失败：' + e.message, { error: true });
  }
};
