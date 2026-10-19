import type { CSSProperties } from 'react';
import { ActionBar, Button, Card, Cell, Divider, NavBar, NoticeBar, Price, ProductCard, Tag } from '@ds';

/**
 * 订单详情 — implements design/figma-export/order-detail (375×900).
 * Mock data mirrors the Figma frame; amounts are derived so the summary matches the design.
 */

const yuan = (amount: number) => `¥${amount.toFixed(2)}`;

const order = {
  notice: '商家已接单，预计 10月22日 前发货',
  status: {
    tag: '待发货',
    title: '等待商家发货',
    countdown: '23:59:59',
    desc: '超时未发货将自动赔付 5 元无门槛券',
  },
  address: {
    name: '张三 138****8888',
    detail: '上海市浦东新区世纪大道 100 号 8 楼',
  },
  item: {
    title: '轻薄羽绒服 女款 短款连帽 冬季保暖外套',
    spec: '米白色 / M',
    price: 299,
    count: 1,
  },
  freight: 0,
  discount: 30,
  info: {
    id: '2026102012345678',
    createdAt: '2026-10-20 10:24',
    payMethod: '微信支付',
  },
};

const subtotal = order.item.price * order.item.count;
const payable = subtotal + order.freight - order.discount;

const bodyText: CSSProperties = {
  fontSize: 'var(--ds-font-body-size)',
  lineHeight: 'var(--ds-font-body-line-height)',
  color: 'var(--ds-color-text-primary)',
};

/** Countdown label — not yet in @ds (see design/SPEC.md "CountDown"). Shows the design's static value. */
function CountDown({ value }: { value: string }) {
  return (
    <span
      style={{
        marginLeft: 'auto',
        whiteSpace: 'nowrap',
        fontSize: 'var(--ds-font-caption-size)',
        lineHeight: 'var(--ds-font-caption-line-height)',
        color: 'var(--ds-color-danger)',
      }}
    >
      剩余 {value}
    </span>
  );
}

/** 订单详情 — existing design built from @ds components plus page-local layout. */
export default function OrderDetail() {
  const { status, address, item, info } = order;

  return (
    <div style={{ minHeight: '100vh', background: 'var(--ds-color-bg-page)', paddingBottom: 56 }}>
      <NavBar title="订单详情" onBack={() => window.history.back()} />

      <div style={{ display: 'flex', flexDirection: 'column', gap: 12, padding: 12 }}>
        <NoticeBar tone="warning">{order.notice}</NoticeBar>

        <Card>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <Tag type="warning" style={{ alignSelf: 'center' }}>
              {status.tag}
            </Tag>
            <span
              style={{
                fontSize: 'var(--ds-font-title-size)',
                lineHeight: 'var(--ds-font-title-line-height)',
                fontWeight: 'var(--ds-font-title-weight)',
                color: 'var(--ds-color-text-primary)',
              }}
            >
              {status.title}
            </span>
            <CountDown value={status.countdown} />
          </div>
          <div
            style={{
              marginTop: 4,
              fontSize: 'var(--ds-font-caption-size)',
              lineHeight: 'var(--ds-font-caption-line-height)',
              color: 'var(--ds-color-text-tertiary)',
            }}
          >
            {status.desc}
          </div>
          <Divider style={{ marginTop: 12 }} />
          <Cell title={address.name} label={address.detail} arrow />
        </Card>

        <Card title="商品信息">
          <div style={{ paddingBottom: 12 }}>
            <ProductCard title={item.title} spec={item.spec} price={item.price} count={item.count} />
          </div>
          <Divider />
          <Cell title="商品总额" value={yuan(subtotal)} />
          <Cell title="运费" value={order.freight > 0 ? yuan(order.freight) : '免运费'} />
          <Cell title="店铺优惠" value={`-${yuan(order.discount)}`} />
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'flex-end',
              gap: 8,
              marginTop: 4,
              height: 32,
            }}
          >
            <span style={bodyText}>实付款</span>
            <Price value={payable} size="large" />
          </div>
        </Card>

        <Card title="订单信息">
          <Cell title="订单编号" value={info.id} />
          <Cell title="下单时间" value={info.createdAt} />
          <Cell title="支付方式" value={info.payMethod} />
        </Card>
      </div>

      <ActionBar>
        <Button>联系客服</Button>
        <Button>修改地址</Button>
        <Button type="primary">提醒发货</Button>
      </ActionBar>
    </div>
  );
}
