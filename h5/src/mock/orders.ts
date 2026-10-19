export type OrderStatus = 'pending_payment' | 'pending_shipment' | 'shipped' | 'completed' | 'cancelled';

export interface OrderItem {
  id: string;
  title: string;
  spec: string;
  price: number;
  count: number;
}

export interface Order {
  id: string;
  shop: string;
  status: OrderStatus;
  createdAt: string;
  items: OrderItem[];
  freight: number;
  total: number;
}

export const orderStatusText: Record<OrderStatus, string> = {
  pending_payment: '待付款',
  pending_shipment: '待发货',
  shipped: '已发货',
  completed: '已完成',
  cancelled: '已取消',
};

export const orders: Order[] = [
  {
    id: 'DD20261018093021',
    shop: '数码好物旗舰店',
    status: 'pending_payment',
    createdAt: '2026-10-18 09:30',
    items: [
      { id: 'p1', title: '无线降噪蓝牙耳机 Pro 主动降噪 长续航 入耳式运动耳机', spec: '曜石黑 / 标准版', price: 899, count: 1 },
    ],
    freight: 0,
    total: 899,
  },
  {
    id: 'DD20261016184512',
    shop: '生活家居馆',
    status: 'pending_shipment',
    createdAt: '2026-10-16 18:45',
    items: [
      { id: 'p2', title: '北欧风陶瓷马克杯 大容量 带盖带勺', spec: '奶白色 / 450ml', price: 39.9, count: 2 },
      { id: 'p3', title: '竹纤维抹布 吸水不掉毛 厨房专用', spec: '10 条装', price: 19.5, count: 1 },
    ],
    freight: 6,
    total: 105.3,
  },
  {
    id: 'DD20261012110233',
    shop: '鲜果优选',
    status: 'shipped',
    createdAt: '2026-10-12 11:02',
    items: [
      { id: 'p4', title: '新疆阿克苏冰糖心苹果 当季现摘', spec: '5 斤装 / 果径 80mm+', price: 49.8, count: 1 },
    ],
    freight: 0,
    total: 49.8,
  },
  {
    id: 'DD20261001201809',
    shop: '运动户外专营店',
    status: 'completed',
    createdAt: '2026-10-01 20:18',
    items: [
      { id: 'p5', title: '轻量透气跑步鞋 男女同款 缓震回弹', spec: '雾灰 / 42', price: 329, count: 1 },
    ],
    freight: 0,
    total: 329,
  },
  {
    id: 'DD20260925081544',
    shop: '图书音像城',
    status: 'cancelled',
    createdAt: '2026-09-25 08:15',
    items: [
      { id: 'p6', title: '设计系统实战：从组件库到多端落地', spec: '平装', price: 68, count: 1 },
    ],
    freight: 0,
    total: 68,
  },
];
