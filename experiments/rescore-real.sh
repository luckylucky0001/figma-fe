#!/usr/bin/env bash
# Re-score every real-project run with the current figmafe (the runs themselves live outside the repo,
# under $EXP, default /home/user/exp). newbee runs are judged on the new page, the others on the change
# against the untouched base commit (--changed HEAD: new files whole, modified files by added lines).
set -uo pipefail
cd "$(dirname "$0")/.."
EXP=${EXP:-/home/user/exp}
for d in "$EXP"/R*-*; do
  r=$(basename "$d")
  case $r in
    Rnb-*) args=(--platform h5 --focus src/views/OrderDetailV2.vue) ;;
    Rlm-*) args=(--platform h5 --changed HEAD) ;;
    *)     args=(--platform android --changed HEAD) ;;
  esac
  bin/figmafe reuse "$d" "${args[@]}" -q -o "experiments/real-project/$r" > "experiments/real-project/$r.txt" 2>&1
  printf '%-10s %s\n' "$r" "$(grep -o '✅ PASS\|❌ FAIL' "experiments/real-project/$r.txt" | head -1)"
done
python3 experiments/real-summary.py > /dev/null && echo "✓ experiments/real-project/summary.md"
