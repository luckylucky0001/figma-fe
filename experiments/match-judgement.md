# figmafe match 准确率

## 1. demo 仓库（有人工确认过的标准答案，自动评估）

运行 `python3 experiments/match-accuracy.py`：H5 **11/11**，Android（`--target compose`）**11/11**。

> 注意：demo 仓库同时被用来校准打分规则，结果偏乐观。下面两个真实开源项目的结果更有参考价值。

## 2. 真实项目（人工逐条判断）

Figma：订单详情（11 个组件）。判断标准：推荐的组件是否就是这个项目里实现该元素时应该用的那个。

### newbee-mall-vue3-app（Vue 3 + Vant）

| Figma | 推荐 | 结论 | 人工判断 |
|---|---|---|---|
| NavBar | SimpleHeader（业务组件，项目里已用 11 次） | 直接复用 | ✓ 项目里的顶部导航就是它。名字叫 NavBar 的那个组件其实是底部标签栏，工具已经通过样式判断出来，没有选它 |
| NoticeBar | Vant NoticeBar | 直接复用 | ✓ |
| Card | Vant CellGroup（inset） | 需确认 | ✓ Vant 里没有通用卡片，用 inset 圆角分组是常见做法。标为需确认是合理的 |
| Tag | Vant Tag | 直接复用 | ✓ |
| CountDown | Vant CountDown | 直接复用 | ✓ 正是「不造轮子」的典型例子 |
| Divider | Vant Divider | 直接复用 | ✓ |
| Cell | Vant Cell | 直接复用 | ✓ |
| ProductCard | Vant Card（项目里已用 2 次） | 直接复用 | ✓ van-card 就是商品卡 |
| Price | — | 需新建 | ✓ Vant 和项目里都没有独立的价格组件 |
| ActionBar | Vant ActionBar（项目里已用 1 次） | 直接复用 | ✓ |
| Button | Vant Button（项目里已用 14 次） | 直接复用 | ✓ |

**11/11**

### QMUI_Android（XML View，`--target xml`）

| Figma | 推荐 | 结论 | 人工判断 |
|---|---|---|---|
| NavBar | QMUITopBar（View，项目里已用 6 次） | 直接复用 | ✓ |
| NoticeBar | — | 需新建 | ✓ QMUI 没有通知条组件 |
| Card | QMUIGroupListView（项目里已用 25 次） | 直接复用 | ✓ QMUI 惯用分组列表来做白底分区 |
| Tag | Material Chip | 需确认 | △ 可以用，但 QMUI demo 里实际常用 QMUIRoundButton 做标签，工具没有识别出这种「借用」写法 |
| CountDown | — | 需新建 | ✓ |
| Divider | MaterialDivider | 直接复用 | △ 项目里更常见的写法是分组列表自带的分隔线，或者一个 1px 的 View |
| Cell | QMUICommonListItemView（项目里已用 48 次） | 直接复用 | ✓ |
| ProductCard | MaterialCardView | 需确认 | △ 只能当容器用，商品卡的内容仍要自己写 |
| Price | — | 需新建 | ✓ |
| ActionBar | — | 需新建 | ✓ |
| Button | QMUIRoundButton（项目里已用 29 次） | 直接复用 | ✓ |

**8/11 完全正确，3/11 可用但不是最优**（都标成了「需确认」或者给了说明，交给人做最终决定）。

### 小结

- 「直接复用」的结论在两个真实项目里都没有出现明显错误。拿不准的情况，工具给出的都是「需确认」。
- 还没覆盖的：「借用」写法（例如用按钮组件来充当标签）、项目里约定俗成的写法（例如用 1px 的 View 做分割线）。这些可以通过 `component-map.yaml` 人工补充，之后就会固定下来。
