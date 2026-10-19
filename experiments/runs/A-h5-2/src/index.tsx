import type { CSSProperties } from 'react';
import { ActionBar, Button, Card, Cell, Divider, NavBar, NoticeBar, Price, ProductCard, Tag } from '@ds';
import { orderStatusText } from '@/mock/orders';
import { orderDetail as order } from '@/mock/orderDetail';
import { CountDown } from './CountDown';

const yuan = (n: number) => `${n < 0 ? '-' : ''}¥${Math.abs(n).toFixed(2)}`;

const statusTitleStyle: CSSProperties = {
  flex: 1,
  minWidth: 0,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
  fontSize: 'var(--ds-font-title-size)',
  lineHeight: 'var(--ds-font-title-line-height)',
  fontWeight: 'var(--ds-font-title-weight)',
  color: 'var(--ds-color-text-primary)',
};

const captionTertiary: CSSProperties = {
  fontSize: 'var(--ds-font-caption-size)',
  lineHeight: 'var(--ds-font-caption-line-height)',
  fontWeight: 'var(--ds-font-caption-weight)',
  color: 'var(--ds-color-text-tertiary)',
};

/** 订单详情 — Figma frame "OrderDetail 订单详情" (design/figma-export/order-detail). */
export default function OrderDetail() {
  return (
    <div style={{ minHeight: '100vh', background: 'var(--ds-color-bg-page)', paddingBottom: 56 }}>
      <NavBar title="订单详情" onBack={() => window.history.back()} />

      <div style={{ display: 'flex', flexDirection: 'column', gap: 12, padding: 12 }}>
        <NoticeBar tone="warning">{order.notice}</NoticeBar>

        {/* Card / 状态 */}
        <Card>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8, height: 24 }}>
            <Tag type="warning" style={{ alignSelf: 'center' }}>
              {orderStatusText[order.status]}
            </Tag>
            <span style={statusTitleStyle}>{order.statusTitle}</span>
            <CountDown seconds={order.remainingSeconds} />
          </div>
          <div style={{ ...captionTertiary, marginTop: 4 }}>{order.statusDesc}</div>
          <Divider style={{ marginTop: 12 }} />
          <Cell
            title={`${order.receiver.name} ${order.receiver.phone}`}
            label={order.receiver.address}
            arrow
          />
        </Card>

        {/* Card / 商品 */}
        <Card title="商品信息">
          <div style={{ display: 'flex', flexDirection: 'column', gap: 12, paddingBottom: 12 }}>
            {order.items.map((item) => (
              <ProductCard key={item.id} title={item.title} spec={item.spec} price={item.price} count={item.count} />
            ))}
          </div>
          <Divider />
          <Cell title="商品总额" value={yuan(order.goodsAmount)} />
          <Cell title="运费" value={order.freight > 0 ? yuan(order.freight) : '免运费'} />
          <Cell title="店铺优惠" value={yuan(-order.discount)} />
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'flex-end',
              gap: 8,
              height: 32,
              marginTop: 4,
            }}
          >
            <span
              style={{
                fontSize: 'var(--ds-font-body-size)',
                lineHeight: 'var(--ds-font-body-line-height)',
                color: 'var(--ds-color-text-primary)',
              }}
            >
              实付款
            </span>
            <Price value={order.total} size="large" tone="primary" />
          </div>
        </Card>

        {/* Card / 订单信息 */}
        <Card title="订单信息">
          <Cell title="订单编号" value={order.id} />
          <Cell title="下单时间" value={order.createdAt} />
          <Cell title="支付方式" value={order.payMethod} />
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
