package com.demo.app.mock

/** Mock data for the order list (mirrors h5/src/mock/orders.ts). */
enum class OrderStatus(val text: String) {
    PendingPayment("待付款"),
    PendingShipment("待发货"),
    Shipped("已发货"),
    Completed("已完成"),
    Cancelled("已取消"),
}

data class OrderItem(val id: String, val title: String, val spec: String, val price: Double, val count: Int)

data class Order(
    val id: String,
    val shop: String,
    val status: OrderStatus,
    val createdAt: String,
    val items: List<OrderItem>,
    val freight: Double,
    val total: Double,
)

val orders: List<Order> = listOf(
    Order(
        id = "DD20261018093021", shop = "数码好物旗舰店", status = OrderStatus.PendingPayment, createdAt = "2026-10-18 09:30",
        items = listOf(OrderItem("p1", "无线降噪蓝牙耳机 Pro 主动降噪 长续航 入耳式运动耳机", "曜石黑 / 标准版", 899.0, 1)),
        freight = 0.0, total = 899.0,
    ),
    Order(
        id = "DD20261016184512", shop = "生活家居馆", status = OrderStatus.PendingShipment, createdAt = "2026-10-16 18:45",
        items = listOf(
            OrderItem("p2", "北欧风陶瓷马克杯 大容量 带盖带勺", "奶白色 / 450ml", 39.9, 2),
            OrderItem("p3", "竹纤维抹布 吸水不掉毛 厨房专用", "10 条装", 19.5, 1),
        ),
        freight = 6.0, total = 105.3,
    ),
    Order(
        id = "DD20261012110233", shop = "鲜果优选", status = OrderStatus.Shipped, createdAt = "2026-10-12 11:02",
        items = listOf(OrderItem("p4", "新疆阿克苏冰糖心苹果 当季现摘", "5 斤装 / 果径 80mm+", 49.8, 1)),
        freight = 0.0, total = 49.8,
    ),
    Order(
        id = "DD20261001201809", shop = "运动户外专营店", status = OrderStatus.Completed, createdAt = "2026-10-01 20:18",
        items = listOf(OrderItem("p5", "轻量透气跑步鞋 男女同款 缓震回弹", "雾灰 / 42", 329.0, 1)),
        freight = 0.0, total = 329.0,
    ),
    Order(
        id = "DD20260925081544", shop = "图书音像城", status = OrderStatus.Cancelled, createdAt = "2026-09-25 08:15",
        items = listOf(OrderItem("p6", "设计系统实战：从组件库到多端落地", "平装", 68.0, 1)),
        freight = 0.0, total = 68.0,
    ),
)
