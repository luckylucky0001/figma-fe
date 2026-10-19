# figmafe match 回归评测（真实开源项目，固定 commit）

由 `python3 experiments/real-match-eval.py` 生成。✓ = 推荐正确，△ = 可用但不是最优，✗ = 推荐错误。
「留出」= 调规则时没用过、第一次跑完才修问题的项目（第一次的结果见 match-judgement.md 第 3 节）。

| 项目 | 技术栈 | 订单详情（组件化设计稿） | 我的持仓（未组件化、中文图层） |
|---|---|---|---|
| newbee | Vue 3 + Vant | ✓11 △0 ✗0 | ✓12 △0 ✗0 |
| qmui | XML + QMUI | ✓11 △0 ✗0 | ✓12 △0 ✗0 |
| nia | Compose + 自有设计系统 | ✓11 △0 ✗0 | ✓12 △0 ✗0 |
| antdm-template | React + antd-mobile | ✓11 △0 ✗0 | ✓12 △0 ✗0 |
| antdm-mobile | React + antd-mobile | ✓11 △0 ✗0 | ✓12 △0 ✗0 |
| jetsnack（留出） | Compose + 自有设计系统 | ✓10 △1 ✗0 | ✓11 △1 ✗0 |
| litemall（留出） | Vue 2 + Vant 2 | ✓11 △0 ✗0 | ✓12 △0 ✗0 |
| hello-uniapp（留出） | uni-app + uni-ui | ✓10 △1 ✗0 | ✓12 △0 ✗0 |
| taro-ui-demo（留出） | Taro + taro-ui | ✓11 △0 ✗0 | ✓12 △0 ✗0 |
| nutui-vite（留出） | Vue 3 + NutUI | ✓11 △0 ✗0 | ✓12 △0 ✗0 |
| newpipe（留出） | XML（AppCompat）+ 少量 Compose | ✓10 △1 ✗0 | ✓11 △1 ✗0 |
| vant-demo-vue3（留出） | Vue 3 + Vant | ✓11 △0 ✗0 | ✓12 △0 ✗0 |
| vant-demo-vue2（留出） | Vue 2 + Vant 2 | ✓11 △0 ✗0 | ✓12 △0 ✗0 |

**合计 ✓294 / △5 / ✗0**
