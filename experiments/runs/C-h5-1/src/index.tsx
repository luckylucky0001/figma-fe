import { ActionBar, Button, Card, Cell, Divider, NavBar, NoticeBar, Price, ProductCard, Tag } from '@ds';
import './index.css';

/** 订单详情 — implemented from Figma frame 1:1 (OrderDetail), built only from @ds components. */
export default function OrderDetail() {
  return (
    <div className="order-detail">
      <NavBar title="订单详情" onBack={() => window.history.back()} data-figma-id="1:2" />

      <div className="order-detail__body">
        <NoticeBar tone="warning" data-figma-id="1:4">
          商家已接单，预计 10月22日 前发货
        </NoticeBar>

        <Card data-figma-id="1:6">
          <div className="order-detail__status-row" data-figma-id="1:7">
            <Tag type="warning" className="order-detail__v-center" data-figma-id="1:8">
              待发货
            </Tag>
            <span className="order-detail__status-title" data-figma-id="1:10">
              等待商家发货
            </span>
            {/* CountDown has no component yet (component-map gap): static text, font.caption + color.danger */}
            <span className="order-detail__countdown" data-figma-id="1:11">
              剩余 23:59:59
            </span>
          </div>
          <div className="order-detail__status-desc" data-figma-id="1:13">
            超时未发货将自动赔付 5 元无门槛券
          </div>
          <Divider data-figma-id="1:14" />
          <Cell title="张三 138****8888" label="上海市浦东新区世纪大道 100 号 8 楼" arrow data-figma-id="1:15" />
        </Card>

        <Card title="商品信息" data-figma-id="1:18">
          <ProductCard
            className="order-detail__product"
            title="轻薄羽绒服 女款 短款连帽 冬季保暖外套"
            spec="米白色 / M"
            price={299.0}
            count={1}
            data-figma-id="1:20"
          />
          <Divider data-figma-id="1:29" />
          <Cell title="商品总额" value="¥299.00" data-figma-id="1:30" />
          <Cell title="运费" value="免运费" data-figma-id="1:33" />
          <Cell title="店铺优惠" value="-¥30.00" data-figma-id="1:36" />
          <div className="order-detail__total-row" data-figma-id="1:39">
            <span className="order-detail__total-label" data-figma-id="1:40">
              实付款
            </span>
            <Price value={269.0} size="large" tone="primary" data-figma-id="1:41" />
          </div>
        </Card>

        <Card title="订单信息" data-figma-id="1:45">
          <Cell title="订单编号" value="2026102012345678" data-figma-id="1:47" />
          <Cell title="下单时间" value="2026-10-20 10:24" data-figma-id="1:50" />
          <Cell title="支付方式" value="微信支付" data-figma-id="1:53" />
        </Card>
      </div>

      <ActionBar data-figma-id="1:56">
        <Button size="medium" type="default" data-figma-id="1:57">
          联系客服
        </Button>
        <Button size="medium" type="default" data-figma-id="1:59">
          修改地址
        </Button>
        <Button size="medium" type="primary" data-figma-id="1:61">
          提醒发货
        </Button>
      </ActionBar>
    </div>
  );
}
