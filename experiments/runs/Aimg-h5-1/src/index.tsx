import { useEffect, useState, type CSSProperties } from 'react';
import { ActionBar, Button, Card, Cell, Divider, NavBar, NoticeBar, Price, ProductCard, Tag } from '@ds';
import { orderDetail as order } from '@/mock/orderDetail';
import { orderStatusText } from '@/mock/orders';

const yuan = (n: number) => `¥${n.toFixed(2)}`;
const pad = (n: number) => String(n).padStart(2, '0');

/**
 * CountDown — "剩余 HH:MM:SS", font.caption, color.danger.
 * Exists in the Figma library but not yet in @ds, so it lives here for now.
 */
function CountDown({ seconds }: { seconds: number }) {
  const [left, setLeft] = useState(seconds);
  useEffect(() => {
    if (left <= 0) return;
    const timer = window.setTimeout(() => setLeft((s) => s - 1), 1000);
    return () => window.clearTimeout(timer);
  }, [left]);
  const h = Math.floor(left / 3600);
  const m = Math.floor((left % 3600) / 60);
  const s = left % 60;
  return (
    <span
      style={{
        flexShrink: 0,
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

/** 订单详情 — built from @ds components; CountDown is a local stand-in until @ds ships it. */
export default function OrderDetail() {
  const { receiver } = order;
  return (
    <div style={{ minHeight: '100vh', background: 'var(--ds-color-bg-page)', paddingBottom: 56 }}>
      <NavBar title="订单详情" onBack={() => window.history.back()} />

      <div style={{ display: 'flex', flexDirection: 'column', gap: 12, padding: 12 }}>
        <NoticeBar tone="warning">{order.notice}</NoticeBar>

        {/* 订单状态 + 收货地址 */}
        <Card>
          <div style={{ paddingBottom: 12 }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
              <Tag type="warning" style={{ alignSelf: 'center' }}>
                {orderStatusText[order.status]}
              </Tag>
              <span style={statusTitleStyle}>{order.statusTitle}</span>
              <CountDown seconds={order.remainSeconds} />
            </div>
            <div style={{ ...captionTertiary, marginTop: 4 }}>{order.statusDesc}</div>
          </div>
          <Divider />
          <Cell title={`${receiver.name} ${receiver.phone}`} label={receiver.address} arrow />
        </Card>

        {/* 商品信息 + 金额 */}
        <Card title="商品信息">
          <div style={{ display: 'flex', flexDirection: 'column', gap: 12, paddingBottom: 12 }}>
            {order.items.map((item) => (
              <ProductCard key={item.id} title={item.title} spec={item.spec} price={item.price} count={item.count} />
            ))}
          </div>
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
              height: 36,
              paddingTop: 4,
              fontSize: 'var(--ds-font-body-size)',
              lineHeight: 'var(--ds-font-body-line-height)',
              color: 'var(--ds-color-text-primary)',
            }}
          >
            <span>实付款</span>
            <Price value={order.total} size="large" />
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
        <Button>联系客服</Button>
        <Button>修改地址</Button>
        <Button type="primary">提醒发货</Button>
      </ActionBar>
    </div>
  );
}
