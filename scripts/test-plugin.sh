#!/usr/bin/env bash
# Runs the Figma plugin (figma-plugin/code.js) against a mocked Plugin API and checks that its export
# reads exactly like the fixture: same tree / geometry / text / fills, same `figmafe match` result.
set -euo pipefail
cd "$(dirname "$0")/.."
tmp=$(mktemp -d); trap 'rm -rf "$tmp"' EXIT
for page in order-detail holdings; do
  node figma-plugin/test/mock-run.mjs "design/figma-export/$page" "$tmp/$page" >/dev/null
  python3 - "design/figma-export/$page/export.json" "$tmp/$page/export.json" <<'PY'
import json, sys
a, b = (json.load(open(p))["frame"] for p in sys.argv[1:3])
def key(n):
    bb = n.get("absoluteBoundingBox") or {}
    fills = [tuple(round(p["color"][c], 3) for c in "rgb") for p in n.get("fills") or []]
    return (n["id"], n["name"], n["type"], n.get("characters"), n.get("componentId"), (n.get("style") or {}).get("fontSize"),
            tuple(round(bb.get(k, 0), 1) for k in ("x", "y", "width", "height")), tuple(fills))
def flat(n, out):
    out.append(key(n)); [flat(c, out) for c in n.get("children") or []]; return out
x, y = flat(a, []), flat(b, [])
bad = [(p, q) for p, q in zip(x, y) if p != q]
if len(x) != len(y) or bad:
    sys.exit(f"{sys.argv[1]}: {len(x)} vs {len(y)} nodes, first difference {bad[:1]}")
PY
  bin/figmafe match --figma "design/figma-export/$page" h5 --platform h5 -o "$tmp/m-fixture" -q >/dev/null
  bin/figmafe match --figma "$tmp/$page" h5 --platform h5 -o "$tmp/m-plugin" -q >/dev/null
  summary='import json,sys;print([(i["figma"],i["verdict"],i["instances"],i.get("variants"),[c["name"] for c in i["candidates"] or []]) for i in json.load(open(sys.argv[1]))["items"]])'
  python3 -c "$summary" "$tmp/m-fixture/h5.json" > "$tmp/a.txt"
  python3 -c "$summary" "$tmp/m-plugin/h5.json" > "$tmp/b.txt"
  if ! cmp -s "$tmp/a.txt" "$tmp/b.txt"; then
    echo "✗ $page: figmafe match reads the plugin export differently"; exit 1
  fi
  echo "✓ plugin export of $page = fixture (tree, geometry, text, fills, instances; match result)"
done
