#!/usr/bin/env bash
# 重新生成教程截图：在一个干净的工作区里真实执行每一步，把输出截成图放到 docs/tutorial/img/。
# 用法：docs/tutorial/capture.sh   （需要：Go、Node、Chromium、Noto Sans SC 字体、Gradle）
set -uo pipefail
repo=$(cd "$(dirname "$0")/../.." && pwd)
img=$repo/docs/tutorial/img
ws=${TUTORIAL_WS:-/tmp/figmafe-tutorial}
shoot() { node "$repo/docs/tools/shoot.mjs" "$@"; }
mkdir -p "$img"
rm -rf "$ws"
git -C "$repo" worktree remove --force "$ws" >/dev/null 2>&1 || true
git -C "$repo" worktree add -q --detach "$ws" HEAD
ln -s "$repo/h5/node_modules" "$ws/h5/node_modules"
cd "$ws"
t=$(mktemp -d)

# term <img-name> <title> <command…>：执行命令，连同提示符一起截图
term() {
  local name=$1 title=$2; shift 2
  { echo "\$ $*"; bash -c "$*" 2>&1; } > "$t/$name.txt"
  shoot term "$title" "$t/$name.txt" "$img/$name.png"
  echo "  📸 $name"
}
# md <img-name> <title> <file> [head-lines]
md() {
  local f=$3
  if [ -n "${4:-}" ]; then head -n "$4" "$3" > "$t/$1.md"; f=$t/$1.md; fi
  shoot md "$2" "$f" "$img/$1.png"
  echo "  📸 $1"
}

echo "Step 1 环境准备"
term 01-build "构建 CLI + 单元测试" "make build >/dev/null && make test 2>&1 | grep -v 'no test files'"

echo "Step 2 Figma 导出"
cp design/figma-export/order-detail/frame@2x.png "$img/02-design.png"
term 02-export "Figma 导出产物（插件导出 / fixture）" "ls -l design/figma-export design/figma-export/order-detail && python3 -c \"import json;e=json.load(open('design/figma-export/order-detail/export.json'));n=e['frame']['children'][4]['children'][1];print(json.dumps({k:n[k] for k in ('id','name','type','componentId','componentProperties','absoluteBoundingBox')},ensure_ascii=False,indent=2))\""

echo "Step 3 资产扫描"
term 03-scan "Phase 0：扫描两端组件资产" "bin/figmafe scan"

echo "Step 4 映射"
term 04-map-suggest "Phase 1：映射建议" "bin/figmafe map suggest"
bin/figmafe map check >/dev/null
md 05-map-check "Phase 1：映射检查（.figmafe/map-check.md）" .figmafe/map-check.md

echo "Step 5 组件手册"
term 06-catalog "生成 Skill 的组件手册" "bin/figmafe catalog"
md 06-catalog-h5 "组件手册（references/catalog-h5.md，节选）" .claude/skills/figma-to-code/references/catalog-h5.md 70

echo "Step 6 复用方案"
bin/figmafe plan --page OrderDetail --platform h5 -q >/dev/null
md 07-plan-h5 "Phase 2：复用方案（.figmafe/h5/OrderDetail/plan.md，节选）" .figmafe/h5/OrderDetail/plan.md 48
bin/figmafe plan --page OrderDetail --platform android -q >/dev/null
sed -n '/## 3/,/## 4/p' .figmafe/android/OrderDetail/plan.md | head -24 > "$t/plan-android.md"
md 07-plan-android "Phase 2：Android 复用方案（节选）" "$t/plan-android.md"

echo "Step 7 Skill"
md 08-skill "Skill：.claude/skills/figma-to-code/SKILL.md" .claude/skills/figma-to-code/SKILL.md

echo "Step 8 第一版实现（故意带 3 个常见问题）→ verify"
mkdir -p h5/src/pages/OrderDetail
cp "$repo/experiments/runs/C-h5-1/src/"* h5/src/pages/OrderDetail/
python3 - <<'EOF'
import re
p='h5/src/pages/OrderDetail/index.tsx'
s=open(p).read()
# 问题 1：自己拼状态标签（不用 Tag，写死颜色）
s=re.sub(r'<Tag\b[^>]*>([^<]*)</Tag>', r'<span style={{ background: "#FFF3E8", color: "#FF8F1F", padding: "0 6px", borderRadius: 4, fontSize: 12 }}>\1</span>', s, count=1)
s=s.replace(' Tag }', ' }').replace('Tag, ', '')  # 不再使用 Tag，去掉 import
# 问题 2：用原生 button 代替 Button
s=re.sub(r'<Button\b[^>]*>\s*(联系客服)\s*</Button>', r'<button className="btn">\1</button>', s, count=1, flags=re.S)
# 问题 3：一行左右结构的 Cell 自己用 div 拼
s=re.sub(r'<Cell\b[^>]*title="支付方式"[^>]*/>', '<div style={{ display: "flex", justifyContent: "space-between", padding: "10px 0" }}><span>支付方式</span><span>微信支付</span></div>', s, count=1)
open(p,'w').write(s)
EOF
grep -n 'span style\|<button\|支付方式' h5/src/pages/OrderDetail/index.tsx | head -5 > "$t/09-diff.txt"
shoot term "第一版实现中的 3 处问题（示例）" "$t/09-diff.txt" "$img/09-bad-code.png"
term 09-verify-fail "Phase 4：verify 第 1 轮（FAIL）" "bin/figmafe verify --page OrderDetail --platform h5 -q; echo exit=\$?"
md 09-report-fail "verify 报告（.figmafe/h5/OrderDetail/report.md）" .figmafe/h5/OrderDetail/report.md 46
cp .figmafe/h5/OrderDetail/diff.png "$img/09-diff-fail.png"
cp .figmafe/h5/OrderDetail/screen.png "$img/09-screen-fail.png"

echo "Step 9 按报告修正 → verify PASS"
cp "$repo/experiments/runs/C-h5-1/src/"* h5/src/pages/OrderDetail/
term 10-verify-pass "Phase 4：verify 第 2 轮（PASS）" "bin/figmafe verify --page OrderDetail --platform h5 -q; echo exit=\$?"
md 10-report-pass "verify 报告（PASS）" .figmafe/h5/OrderDetail/report.md 40
cp .figmafe/h5/OrderDetail/screen.png "$img/10-screen-pass.png"
cp .figmafe/h5/OrderDetail/diff.png "$img/10-diff-pass.png"

echo "Step 10 Android"
mkdir -p android/app/src/main/kotlin/com/demo/app/screens/OrderDetail
cp "$repo/experiments/runs/C-android-1/src/"* android/app/src/main/kotlin/com/demo/app/screens/OrderDetail/
python3 - <<'EOF'
p='android/app/src/main/kotlin/com/demo/app/Screens.kt'
s=open(p).read()
if 'OrderDetail' not in s:
    s=s.replace('import com.demo.app.screens.OrderListScreen','import com.demo.app.screens.OrderListScreen\nimport com.demo.app.screens.OrderDetail.OrderDetailScreen')
    s=s.replace('"Gallery" to { GalleryScreen() },','"Gallery" to { GalleryScreen() },\n    "OrderDetail" to { OrderDetailScreen() },')
open(p,'w').write(s)
EOF
term 11-verify-android "Android：verify" "bin/figmafe verify --page OrderDetail --platform android -q; echo exit=\$?"
md 11-report-android "Android verify 报告" .figmafe/android/OrderDetail/report.md 44
cp .figmafe/android/OrderDetail/screen.png "$img/11-screen-android.png"
cp .figmafe/android/OrderDetail/diff.png "$img/11-diff-android.png"

echo "Step 10b 三联对比图"
# compare <img-name> <title> <img1> <label1> <img2> <label2> <img3> <label3>
compare() {
  local name=$1 title=$2; shift 2
  local cells=""
  while [ $# -gt 0 ]; do cells+="<figure><img src=\"$img/$1\"><figcaption>$2</figcaption></figure>"; shift 2; done
  cat > "$t/$name.html" <<HTML
<!doctype html><html><head><meta charset="utf-8"><style>
body{margin:0;padding:20px;background:#eef0f3;font-family:'Noto Sans SC',sans-serif}
.doc{display:inline-block;background:#fff;border-radius:10px;padding:16px 18px;box-shadow:0 8px 28px rgba(0,0,0,.15)}
h3{margin:0 0 12px;font-size:16px;color:#1f2329} .row{display:flex;gap:16px}
figure{margin:0;text-align:center} img{width:300px;border:1px solid #e5e6eb;border-radius:6px;display:block}
figcaption{margin-top:6px;font-size:13px;color:#646a73}
</style></head><body><div class="doc"><h3>$title</h3><div class="row">$cells</div></div></body></html>
HTML
  shoot html "$t/$name.html" "$img/$name.png"
  echo "  📸 $name"
}
compare 09-compare "第 1 轮：设计稿 / 实现 / 差异（红色为不一致像素）" 02-design.png "Figma 设计稿" 09-screen-fail.png "第一版实现" 09-diff-fail.png "差异图"
compare 10-compare "第 2 轮：修正后" 02-design.png "Figma 设计稿" 10-screen-pass.png "修正后（H5）" 10-diff-pass.png "差异图"
compare 11-compare "Android（Compose）" 02-design.png "Figma 设计稿" 11-screen-android.png "Android 实现" 11-diff-android.png "差异图"

echo "Step 11 沉淀规则"
md 12-rules "沉淀规则（references/rules.md）" .claude/skills/figma-to-code/references/rules.md

echo "Part 2 存量项目（真实开源项目，固定 commit，由 experiments/real-match-eval.py 拉取到 experiments/.cache/）"
python3 "$repo/experiments/real-match-eval.py" newbee qmui nia >/dev/null 2>&1 || true
cache=$repo/experiments/.cache
# html <img-name> <file.html>
html() { shoot html "$2" "$img/$1.png"; echo "  📸 $1"; }
# section <img-name> <title> <md file> <from-heading-regex> <to-heading-regex>
section() { sed -n "/$4/,/$5/p" "$3" | sed '$d' > "$t/$1.md"; shoot md "$2" "$t/$1.md" "$img/$1.png"; echo "  📸 $1"; }

term 20-reuse-newbee "存量体检：newbee-mall（Vue 3 + Vant）" "bin/figmafe reuse $cache/newbee -o $t/nb -q"
head -n 40 "$t/nb/h5.md" > "$t/20-report.md"; shoot md "体检报告节选（.figmafe/reuse/h5.md）" "$t/20-report.md" "$img/20-reuse-report.png"; echo "  📸 20-reuse-report"
term 21-reuse-qmui "存量体检：QMUI_Android（XML + QMUI）" "bin/figmafe reuse $cache/qmui -o $t/qm -q"
section 21-qmui-compat "QMUI：组件库里有、但本项目不能直接用的组件" "$t/qm/android.md" "组件库里有" "^## 2"
term 22-reuse-nia "存量体检：Now in Android（Compose + 自有设计系统）" "bin/figmafe reuse $cache/nia -o $t/nia -q"
section 22-nia-wrappers "Now in Android：设计系统封装" "$t/nia/android.md" "设计系统封装" "^## 2"

term 23-match-newbee "新设计稿 → newbee 的复用清单" "bin/figmafe match --figma design/figma-export/order-detail $cache/newbee -o $t/m-nb"
html 23-match-newbee-overlay "$t/m-nb/h5.html"
section 23-match-props "变体 → 属性（推荐组件怎么配置）" "$t/m-nb/h5.md" "^## 变体" "^## 候选"
term 24-match-qmui "同一份设计稿 → QMUI（XML）" "bin/figmafe match --figma design/figma-export/order-detail $cache/qmui -o $t/m-qm -q"
html 24-match-qmui-overlay "$t/m-qm/android.html"
term 25-match-holdings "没组件化、中文图层名的设计稿 → newbee" "bin/figmafe match --figma design/figma-export/holdings $cache/newbee -o $t/m-h -q"
html 25-match-holdings-overlay "$t/m-h/h5.html"

term 26-check-a "页面复用检查：不用 Skill 的实现（newbee，A 组）" "cat $repo/experiments/real-project/Rnb-A-1.txt | sed -n '1,22p'"
term 26-check-c "页面复用检查：用 Skill 的实现（newbee，C 组）" "cat $repo/experiments/real-project/Rnb-C-1.txt | sed -n '1,12p'"
term 27-gate "PR 闸门：只检查改动的文件（Now in Android，C 组，决定已写进 component-map.yaml）" "cat $repo/experiments/real-project/Rnia-C-2.txt | sed -n '1,14p'"
md 28-eval "match 回归评测（14 个真实开源项目）" "$repo/experiments/real-match-eval.md"
md 29-summary "真实项目对照实验汇总" "$repo/experiments/real-project/summary.md" 60

cd "$repo" && git worktree remove --force "$ws" >/dev/null 2>&1
echo "✓ screenshots → docs/tutorial/img/"
