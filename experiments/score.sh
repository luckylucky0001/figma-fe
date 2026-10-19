#!/usr/bin/env bash
# 统一打分：用主仓库的 figmafe + 映射，对某次实验的工作目录跑 verify，结果存到 experiments/runs/<run>/。
# 用法：experiments/score.sh <run-name> <h5|android> <worktree-dir> [pkg]
#   pkg：组件库以 npm 包形式提供的场景，使用 experiments/figmafe.pkg.yaml 和 component-map.pkg.yaml
set -uo pipefail
run=$1; pf=$2; wt=$3; variant=${4:-}
repo=$(cd "$(dirname "$0")/.." && pwd)
out=$repo/experiments/runs/$run
mkdir -p "$out"
# A 组的工作目录里没有 figmafe 配置和映射：临时放进去，打完分再删掉
added=()
for f in figmafe.yaml design/component-map.yaml design/figma-export/library.json design/figma-export/variables.tokens.json design/figma-export/order-detail/export.json; do
  if [ ! -e "$wt/$f" ]; then
    src=$repo/$f
    if [ "$variant" = pkg ] && [[ $f == *.yaml ]]; then src=$repo/experiments/$(basename "${f%.yaml}").pkg.yaml; fi
    cp "$src" "$wt/$f"; added+=("$f")
  fi
done
"$repo/bin/figmafe" verify --page OrderDetail --platform "$pf" --root "$wt" -q > "$out/verify.log" 2>&1
echo "exit=$?" >> "$out/verify.log"
cp "$wt/.figmafe/$pf/OrderDetail/"{report.md,report.json,plan.md,screen.png,diff.png} "$out/" 2>/dev/null
# 记录这次运行改动了哪些文件（组件库被改动要单独标出）
(cd "$wt" && git status --porcelain --untracked-files=all | grep -v '^?? .figmafe' | grep -v 'node_modules' > "$out/changed-files.txt")
if [ "$pf" = h5 ]; then src=h5/src/pages/OrderDetail; else src=android/app/src/main/kotlin/com/demo/app/screens/OrderDetail; fi
mkdir -p "$out/src" && cp -r "$wt/$src/." "$out/src/" 2>/dev/null
for f in "${added[@]}"; do rm -f "$wt/$f"; done
tail -1 "$out/verify.log"; sed -n '3p' "$out/report.md"
