import { Button, Card, Divider, NavBar, NoticeBar, Price, ProductCard, Tag } from '@ds';
import type { TagProps } from '@ds';
import { orders, orderStatusText, type Order, type OrderStatus } from '@/mock/orders';

const statusTagType: Record<OrderStatus, NonNullable<TagProps['type']>> = {
  pending_payment: 'primary',
  pending_shipment: 'warning',
  shipped: 'success',
  completed: 'success',
  cancelled: 'danger',
};

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
          <Button>修改地址</Button>
          <Button>提醒发货</Button>
        </>
      );
    case 'shipped':
      return (
        <>
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

function OrderCard({ order }: { order: Order }) {
  const count = order.items.reduce((sum, item) => sum + item.count, 0);
  return (
    <Card data-order-id={order.id}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
        <span
          style={{
            fontSize: 'var(--ds-font-subtitle-size)',
            lineHeight: 'var(--ds-font-subtitle-line-height)',
            fontWeight: 600,
            color: 'var(--ds-color-text-primary)',
          }}
        >
          {order.shop}
        </span>
        <Tag type={statusTagType[order.status]}>{orderStatusText[order.status]}</Tag>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: 12, paddingBottom: 12 }}>
        {order.items.map((item) => (
          <ProductCard key={item.id} title={item.title} spec={item.spec} price={item.price} count={item.count} />
        ))}
      </div>

      <Divider />

      <div
        style={{
          display: 'flex',
          alignItems: 'baseline',
          justifyContent: 'flex-end',
          gap: 4,
          padding: '12px 0',
          fontSize: 'var(--ds-font-caption-size)',
          lineHeight: 'var(--ds-font-caption-line-height)',
          color: 'var(--ds-color-text-secondary)',
        }}
      >
        <span>共 {count} 件</span>
        <span>{order.freight > 0 ? `含运费 ¥${order.freight.toFixed(2)}` : '免运费'}</span>
        <span style={{ color: 'var(--ds-color-text-primary)' }}>实付</span>
        <Price value={order.total} size="medium" />
      </div>

      <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
        <OrderActions status={order.status} />
      </div>
    </Card>
  );
}

/** 订单列表 — existing business page built only from @ds components. */
export default function OrderList() {
  return (
    <div style={{ minHeight: '100vh', background: 'var(--ds-color-bg-page)' }}>
      <NavBar title="我的订单" onBack={() => window.history.back()} />
      <div style={{ display: 'flex', flexDirection: 'column', gap: 12, padding: 12 }}>
        <NoticeBar tone="warning">受天气影响，部分地区配送可能延迟 1-2 天，请耐心等待</NoticeBar>
        {orders.map((order) => (
          <OrderCard key={order.id} order={order} />
        ))}
      </div>
    </div>
  );
}
