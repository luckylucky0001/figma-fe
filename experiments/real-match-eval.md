# figmafe match 回归评测（真实开源项目，固定 commit）

由 `python3 experiments/real-match-eval.py` 生成。✓ = 推荐正确，△ = 可用但不是最优，✗ = 推荐错误。
「留出」= 调规则时没用过、第一次跑完才修问题的项目（第一次的结果见 match-judgement.md 第 3 节）。

| 项目 | 技术栈 | 订单详情（组件化设计稿） | 我的持仓（未组件化、中文图层） | 买入下单（未组件化、英文图层、新概念） |
|---|---|---|---|---|
| newbee | Vue 3 + Vant | ✓11 △0 ✗0 | ✓12 △0 ✗0 | ✓10 △0 ✗0 |
| qmui | XML + QMUI | ✓11 △0 ✗0 | ✓12 △0 ✗0 | ✓8 △2 ✗0 |
| nia | Compose + 自有设计系统 | ✓11 △0 ✗0 | ✓12 △0 ✗0 | ✓9 △1 ✗0 |
| antdm-template | React + antd-mobile | ✓11 △0 ✗0 | ✓12 △0 ✗0 | ✓10 △0 ✗0 |
| antdm-mobile | React + antd-mobile | ✓11 △0 ✗0 | ✓12 △0 ✗0 | ✓10 △0 ✗0 |
| jetsnack（留出） | Compose + 自有设计系统 | ✓10 △1 ✗0 | ✓11 △1 ✗0 | ✓10 △0 ✗0 |
| litemall（留出） | Vue 2 + Vant 2 | ✓11 △0 ✗0 | ✓12 △0 ✗0 | ✓10 △0 ✗0 |
| hello-uniapp（留出） | uni-app + uni-ui | ✓10 △1 ✗0 | ✓12 △0 ✗0 | ✓9 △1 ✗0 |
| taro-ui-demo（留出） | Taro + taro-ui | ✓11 △0 ✗0 | ✓12 △0 ✗0 | ✓10 △0 ✗0 |
| nutui-vite（留出） | Vue 3 + NutUI | ✓11 △0 ✗0 | ✓12 △0 ✗0 | ✓10 △0 ✗0 |
| sunflower（留出） | Compose + Material3（无自有设计系统） | ✓10 △1 ✗0 | ✓11 △1 ✗0 | ✓10 △0 ✗0 |
| pokedex（留出） | XML + Material（MaterialComponents 主题） | ✓11 △0 ✗0 | ✓12 △0 ✗0 | ✓10 △0 ✗0 |
| newpipe（留出） | XML（AppCompat）+ 少量 Compose | ✓10 △1 ✗0 | ✓11 △1 ✗0 | ✓8 △2 ✗0 |
| vant-demo-vue3（留出） | Vue 3 + Vant | ✓11 △0 ✗0 | ✓12 △0 ✗0 | ✓10 △0 ✗0 |
| vant-demo-vue2（留出） | Vue 2 + Vant 2 | ✓11 △0 ✗0 | ✓12 △0 ✗0 | ✓10 △0 ✗0 |
| vue3-vant-mobile（留出） | Vue 3 + Vant 4 | ✓11 △0 ✗0 | ✓12 △0 ✗0 | ✓10 △0 ✗0 |
| mall-app-web（留出） | uni-app + 自有组件 | ✓10 △1 ✗0 | ✓12 △0 ✗0 | ✓10 △0 ✗0 |
| antennapod（留出） | XML + Material3 主题 | ✓10 △1 ✗0 | ✓11 △1 ✗0 | ✓9 △1 ✗0 |
| wanandroid（留出） | XML（AppCompat）+ 三方 ActionBarEx | ✓10 △1 ✗0 | ✓11 △1 ✗0 | ✓9 △1 ✗0 |
| seal（留出） | Compose + Material3 + 自有组件 | ✓10 △1 ✗0 | ✓11 △1 ✗0 | ✓9 △1 ✗0 |
| kunminx（留出） | XML + DataBinding（AppCompat 主题） | ✓10 △1 ✗0 | ✓11 △1 ✗0 | ✓7 △3 ✗0 |

**合计 ✓665 / △28 / ✗0**
