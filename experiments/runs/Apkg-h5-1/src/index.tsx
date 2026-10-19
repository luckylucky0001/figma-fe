import { useEffect, useState, type CSSProperties } from 'react';
import { ActionBar, Button, Card, Cell, Divider, NavBar, NoticeBar, Price, ProductCard, Tag } from '@company/mobile-ui';
import type { TagProps } from '@company/mobile-ui';
import { orderDetail } from '@/mock/orderDetail';
import { orderStatusText, type OrderStatus } from '@/mock/orders';

const statusTagType: Record<OrderStatus, NonNullable<TagProps['type']>> = {
  pending_payment: 'primary',
  pending_shipment: 'warning',
  shipped: 'success',
  completed: 'success',
  cancelled: 'danger',
};

const yuan = (n: number) => `¥${n.toFixed(2)}`;
const pad = (n: number) => String(n).padStart(2, '0');

/**
 * CountDown — Figma component `CountDown` (10:11) has no implementation in @company/mobile-ui yet,
 * so it lives here: 「剩余 HH:MM:SS」, font.caption, color.danger.
 */
function CountDown({ seconds }: { seconds: number }) {
  const [left, setLeft] = useState(seconds);

  useEffect(() => {
    setLeft(seconds);
    const timer = window.setInterval(() => setLeft((s) => Math.max(0, s - 1)), 1000);
    return () => window.clearInterval(timer);
  }, [seconds]);

  const h = Math.floor(left / 3600);
  const m = Math.floor((left % 3600) / 60);
  const s = left % 60;
  return (
    <span
      style={{
        flexShrink: 0,
        marginLeft: 'auto',
        fontSize: 'var(--ds-font-caption-size)',
        lineHeight: 'var(--ds-font-caption-line-height)',
        fontWeight: 'var(--ds-font-caption-weight)',
        color: 'var(--ds-color-danger)',
        whiteSpace: 'nowrap',
      }}
    >
      剩余 {pad(h)}:{pad(m)}:{pad(s)}
    </span>
  );
}

function OrderActions({ status }: { status: OrderStatus }) {
  switch (status) {
    case 'pending_payment':
      return (
        <>
          <Button>取消订单</Button>
          <Button type="primary">去支付</Button>
        </>
      );
    case 'pending_shipment':
      return (
        <>
          <Button>联系客服</Button>
          <Button>修改地址</Button>
          <Button type="primary">提醒发货</Button>
        </>
      );
    case 'shipped':
      return (
        <>
          <Button>联系客服</Button>
          <Button>查看物流</Button>
          <Button type="primary">确认收货</Button>
        </>
      );
    case 'completed':
      return (
        <>
          <Button>申请售后</Button>
          <Button type="primary">再次购买</Button>
        </>
      );
    case 'cancelled':
      return (
        <>
          <Button type="text">删除订单</Button>
          <Button>再次购买</Button>
        </>
      );
  }
}

const statusTitleStyle: CSSProperties = {
  minWidth: 0,
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
  fontSize: 'var(--ds-font-title-size)',
  lineHeight: 'var(--ds-font-title-line-height)',
  fontWeight: 'var(--ds-font-title-weight)',
  color: 'var(--ds-color-text-primary)',
};

const captionTertiaryStyle: CSSProperties = {
  fontSize: 'var(--ds-font-caption-size)',
  lineHeight: 'var(--ds-font-caption-line-height)',
  fontWeight: 'var(--ds-font-caption-weight)',
  color: 'var(--ds-color-text-tertiary)',
};

/** 订单详情 — Figma frame `OrderDetail 订单详情` (1:1), built from @company/mobile-ui. */
export default function OrderDetail() {
  const order = orderDetail;

  return (
    <div style={{ minHeight: '100vh', background: 'var(--ds-color-bg-page)', paddingBottom: 56 + 12 }}>
      <NavBar title="订单详情" onBack={() => window.history.back()} />

      <div style={{ display: 'flex', flexDirection: 'column', gap: 12, padding: 12 }}>
        <NoticeBar tone="warning">{order.notice}</NoticeBar>

        {/* 状态 + 收货地址 */}
        <Card>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8, height: 24 }}>
            <Tag type={statusTagType[order.status]} style={{ alignSelf: 'center' }}>
              {orderStatusText[order.status]}
            </Tag>
            <span style={statusTitleStyle}>{order.statusTitle}</span>
            <CountDown seconds={order.remainingSeconds} />
          </div>
          <div style={{ ...captionTertiaryStyle, marginTop: 4, marginBottom: 12 }}>{order.statusDesc}</div>
          <Divider />
          <Cell
            title={`${order.receiver.name} ${order.receiver.phone}`}
            label={order.receiver.address}
            arrow
            style={{ cursor: 'pointer' }}
          />
        </Card>

        {/* 商品信息 */}
        <Card title="商品信息">
          {order.items.map((item) => (
            <ProductCard
              key={item.id}
              title={item.title}
              spec={item.spec}
              price={item.price}
              count={item.count}
              style={{ marginBottom: 12 }}
            />
          ))}
          <Divider />
          <Cell title="商品总额" value={yuan(order.goodsAmount)} />
          <Cell title="运费" value={order.freight > 0 ? yuan(order.freight) : '免运费'} />
          <Cell title="店铺优惠" value={`-${yuan(order.discount)}`} />
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

        {/* 订单信息 */}
        <Card title="订单信息">
          <Cell title="订单编号" value={order.id} />
          <Cell title="下单时间" value={order.createdAt} />
          <Cell title="支付方式" value={order.payMethod} />
        </Card>
      </div>

      <ActionBar>
        <OrderActions status={order.status} />
      </ActionBar>
    </div>
  );
}
