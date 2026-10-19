package com.demo.app.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.RowScope
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import com.demo.app.mock.Order
import com.demo.app.mock.OrderStatus
import com.demo.app.mock.orders
import com.demo.ds.DsBadge
import com.demo.ds.DsBadgeTone
import com.demo.ds.DsBanner
import com.demo.ds.DsBannerTone
import com.demo.ds.DsButton
import com.demo.ds.DsButtonStyle
import com.demo.ds.DsCard
import com.demo.ds.DsColors
import com.demo.ds.DsDivider
import com.demo.ds.DsGoodsCard
import com.demo.ds.DsPriceText
import com.demo.ds.DsPriceTone
import com.demo.ds.DsSpacing
import com.demo.ds.DsText
import com.demo.ds.DsTopBar
import com.demo.ds.DsTypography
import java.util.Locale

private fun OrderStatus.badgeTone(): DsBadgeTone = when (this) {
    OrderStatus.PendingPayment -> DsBadgeTone.Primary
    OrderStatus.PendingShipment -> DsBadgeTone.Warning
    OrderStatus.Shipped, OrderStatus.Completed -> DsBadgeTone.Success
    OrderStatus.Cancelled -> DsBadgeTone.Danger
}

@Composable
private fun RowScope.OrderActions(status: OrderStatus) {
    when (status) {
        OrderStatus.PendingPayment -> {
            DsButton(text = "取消订单", onClick = {}, style = DsButtonStyle.Default)
            DsButton(text = "去支付", onClick = {}, style = DsButtonStyle.Primary)
        }
        OrderStatus.PendingShipment -> {
            DsButton(text = "修改地址", onClick = {}, style = DsButtonStyle.Default)
            DsButton(text = "提醒发货", onClick = {}, style = DsButtonStyle.Default)
        }
        OrderStatus.Shipped -> {
            DsButton(text = "查看物流", onClick = {}, style = DsButtonStyle.Default)
            DsButton(text = "确认收货", onClick = {}, style = DsButtonStyle.Primary)
        }
        OrderStatus.Completed -> {
            DsButton(text = "申请售后", onClick = {}, style = DsButtonStyle.Default)
            DsButton(text = "再次购买", onClick = {}, style = DsButtonStyle.Primary)
        }
        OrderStatus.Cancelled -> {
            DsButton(text = "删除订单", onClick = {}, style = DsButtonStyle.Text)
            DsButton(text = "再次购买", onClick = {}, style = DsButtonStyle.Default)
        }
    }
}

@Composable
private fun OrderCard(order: Order) {
    val count = order.items.sumOf { it.count }
    DsCard(modifier = Modifier.fillMaxWidth()) {
        Row(
            modifier = Modifier.fillMaxWidth().padding(bottom = DsSpacing.Md),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically,
        ) {
            DsText(text = order.shop, style = DsTypography.Subtitle, maxLines = 1)
            DsBadge(text = order.status.text, tone = order.status.badgeTone())
        }

        Column(
            modifier = Modifier.fillMaxWidth().padding(bottom = DsSpacing.Md),
            verticalArrangement = Arrangement.spacedBy(DsSpacing.Md),
        ) {
            order.items.forEach { item ->
                DsGoodsCard(title = item.title, price = item.price, spec = item.spec, count = item.count)
            }
        }

        DsDivider()

        Row(
            modifier = Modifier.fillMaxWidth().padding(vertical = DsSpacing.Md),
            horizontalArrangement = Arrangement.spacedBy(DsSpacing.Xs, Alignment.End),
        ) {
            val freight = if (order.freight > 0) "含运费 ¥" + String.format(Locale.US, "%.2f", order.freight) else "免运费"
            DsText(text = "共 $count 件", style = DsTypography.Caption, color = DsColors.TextSecondary, modifier = Modifier.alignByBaseline())
            DsText(text = freight, style = DsTypography.Caption, color = DsColors.TextSecondary, modifier = Modifier.alignByBaseline())
            DsText(text = "实付", style = DsTypography.Caption, color = DsColors.TextPrimary, modifier = Modifier.alignByBaseline())
            DsPriceText(amount = order.total, tone = DsPriceTone.Primary, modifier = Modifier.alignByBaseline())
        }

        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.spacedBy(DsSpacing.Sm, Alignment.End),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            OrderActions(order.status)
        }
    }
}

/** 订单列表 — existing business screen built only from com.demo.ds components. */
@Composable
fun OrderListScreen() {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(DsColors.BgPage)
            .verticalScroll(rememberScrollState()),
    ) {
        DsTopBar(title = "我的订单", onBack = {})
        Column(
            modifier = Modifier.fillMaxWidth().padding(DsSpacing.Md),
            verticalArrangement = Arrangement.spacedBy(DsSpacing.Md),
        ) {
            DsBanner(text = "受天气影响，部分地区配送可能延迟 1-2 天，请耐心等待", tone = DsBannerTone.Warning)
            orders.forEach { OrderCard(it) }
        }
    }
}
