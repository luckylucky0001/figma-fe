# 复用方案：OrderDetail（h5）

> 由 `figmafe plan` 生成。**闸门：请确认本方案后再写代码。**
> 设计帧 375×900；节点 id 即 Figma 节点 id，实现时挂到对应组件上（H5 `data-figma-id`，Android `Modifier.figmaId`）。

## 1. 组件清单

| 组件 | 数量 |
|---|---|
| `ActionBar` | 1 |
| `Button` | 3 |
| `Card` | 3 |
| `Cell` | 7 |
| `Divider` | 2 |
| `NavBar` | 1 |
| `NoticeBar` | 1 |
| `Price` | 1 |
| `ProductCard` | 1 |
| `Tag` | 1 |
| **合计** | **21** |

## 2. 待决项（Figma 有、代码没有）

- **CountDown**：两端暂无组件。本期按普通文本实现（font.caption + color.danger，静态文案），已登记为组件需求；不要在页面里自造倒计时组件

## 3. 页面结构（节点 → 组件 → 代码）

  - [1:2] **NavBar** → `<NavBar title="订单详情" data-figma-id="1:2" />` (0,0 375×44)
  - [1:4] **NoticeBar** → `<NoticeBar tone="warning" data-figma-id="1:4">商家已接单，预计 10月22日 前发货</NoticeBar>` (12,56 351×36)
  - [1:6] **Card** → `<Card data-figma-id="1:6" />` (12,104 351×155)
    - [1:7] 容器 StatusRow 横向 间距 8(space.sm) (28,120 319×24)
      - [1:8] **Tag** → `<Tag type="warning" data-figma-id="1:8">待发货</Tag>` (28,122 48.015625×20)
      - [1:10] 文本「等待商家发货」 样式 var(--ds-font-title-*) + var(--ds-color-text-primary) (84.015625,120 102×24)
      - [1:11] **CountDown** (273.671875,123 73.328125×18) ⚠ 未映射：CountDown → 处理方案：两端暂无组件。本期按普通文本实现（font.caption + color.danger，静态文案），已登记为组件需求；不要在页面里自造倒计时组件 — 文本「剩余 23:59:59」
    - [1:13] 文本「超时未发货将自动赔付 5 元无门槛券」 样式 var(--ds-font-caption-*) + var(--ds-color-text-tertiary) (28,148 319×18)
    - [1:14] **Divider** → `<Divider data-figma-id="1:14" />` (28,178 319×1)
    - [1:15] **Cell** → `<Cell title="张三 138****8888" label="上海市浦东新区世纪大道 100 号 8 楼" arrow data-figma-id="1:15" />` (28,179 319×64)
  - [1:18] **Card** → `<Card title="商品信息" data-figma-id="1:18" />` (12,271 351×327)
    - [1:20] **ProductCard** → `<ProductCard title="轻薄羽绒服 女款 短款连帽 冬季保暖外套" price={299.00} count={1} spec="米白色 / M" data-figma-id="1:20" />` (28,321 319×80)
    - [1:29] **Divider** → `<Divider data-figma-id="1:29" />` (28,413 319×1)
    - [1:30] **Cell** → `<Cell title="商品总额" value="¥299.00" data-figma-id="1:30" />` (28,414 319×44)
    - [1:33] **Cell** → `<Cell title="运费" value="免运费" data-figma-id="1:33" />` (28,458 319×44)
    - [1:36] **Cell** → `<Cell title="店铺优惠" value="-¥30.00" data-figma-id="1:36" />` (28,502 319×44)
    - [1:39] 容器 TotalRow 横向 间距 8(space.sm) (28,550 319×32)
      - [1:40] 文本「实付款」 样式 var(--ds-font-body-*) + var(--ds-color-text-primary) (229.3125,556 42×20)
      - [1:41] **Price** → `<Price value={269.00} size="large" tone="primary" data-figma-id="1:41" />` (279.3125,552.796875 67.6875×26.390625)
  - [1:45] **Card** → `<Card title="订单信息" data-figma-id="1:45" />` (12,610 351×198)
    - [1:47] **Cell** → `<Cell title="订单编号" value="2026102012345678" data-figma-id="1:47" />` (28,660 319×44)
    - [1:50] **Cell** → `<Cell title="下单时间" value="2026-10-20 10:24" data-figma-id="1:50" />` (28,704 319×44)
    - [1:53] **Cell** → `<Cell title="支付方式" value="微信支付" data-figma-id="1:53" />` (28,748 319×44)
  - [1:56] **ActionBar** → `<ActionBar data-figma-id="1:56" />` (0,844 375×56)
    - [1:57] **Button** → `<Button size="medium" type="default" data-figma-id="1:57">联系客服</Button>` (75,856.5 90×32)
    - [1:59] **Button** → `<Button size="medium" type="default" data-figma-id="1:59">修改地址</Button>` (173,856.5 90×32)
    - [1:61] **Button** → `<Button size="medium" type="primary" data-figma-id="1:61">提醒发货</Button>` (271,856.5 88×32)

## 4. 约束（来自 component-map）

- **NavBar**：页面顶部只放一个 NavBar；不要自己用 div/Row 拼标题栏
- **Card**：Card 已内置 16 内边距和 12 圆角，页面不要再给 Card 加 padding / 背景 / 圆角
- **Card**：卡片之间的间距用 space.md (12)，由页面容器的 gap 提供
- **Tag**：状态色只能通过 type/tone 表达，禁止自己写背景色和文字色
- **Divider**：不要用 border-bottom / material3.Divider 代替
- **Cell**：左标题 / 右值的行一律用 Cell，不要自己用 flex/Row + 两个文本拼
- **ProductCard**：商品行一律用 ProductCard，内部的缩略图 / 标题 / 规格 / 价格 / 数量不要拆开自己拼
- **Price**：金额一律用 Price，不要拼「¥」+ 数字文本（字号层级和基线对齐会错）
- **Price**：例外：Cell 右侧的值（如「¥299.00」「-¥30.00」）在设计稿里是普通文本，用 Cell 的 value 传字符串
- **ActionBar**：吸底操作区一律用 ActionBar；页面内容区底部要留出 56 的空间
- **Button**：一屏最多一个 Primary 按钮
- **Button**：H5 禁止原生 <button>；Android 禁止 material3.Button / OutlinedButton
