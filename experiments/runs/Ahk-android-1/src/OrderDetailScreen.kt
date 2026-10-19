package com.demo.app.screens

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.background
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.demo.ds.DsBadge
import com.demo.ds.DsBadgeTone
import com.demo.ds.DsBanner
import com.demo.ds.DsBannerTone
import com.demo.ds.DsBottomBar
import com.demo.ds.DsButton
import com.demo.ds.DsButtonStyle
import com.demo.ds.DsCard
import com.demo.ds.DsColors
import com.demo.ds.DsDivider
import com.demo.ds.DsGoodsCard
import com.demo.ds.DsListItem
import com.demo.ds.DsPriceSize
import com.demo.ds.DsPriceText
import com.demo.ds.DsPriceTone
import com.demo.ds.DsSpacing
import com.demo.ds.DsText
import com.demo.ds.DsTopBar
import com.demo.ds.DsTypography
import com.demo.ds.figmaId
import java.util.Locale

/** Data for the 订单详情 page. Amounts are in yuan. */
data class OrderDetailItem(val title: String, val spec: String, val price: Double, val count: Int)

data class OrderDetail(
    val notice: String,
    val statusText: String,
    val statusTone: DsBadgeTone,
    val statusTitle: String,
    val countdown: String,
    val statusDesc: String,
    val receiver: String,
    val address: String,
    val item: OrderDetailItem,
    val goodsTotal: Double,
    val freight: Double,
    val discount: Double,
    val paid: Double,
    val orderNo: String,
    val createdAt: String,
    val payMethod: String,
)

/** Mock data matching the Figma frame (design/figma-export/order-detail). */
val sampleOrderDetail = OrderDetail(
    notice = "商家已接单，预计 10月22日 前发货",
    statusText = "待发货",
    statusTone = DsBadgeTone.Warning,
    statusTitle = "等待商家发货",
    countdown = "剩余 23:59:59",
    statusDesc = "超时未发货将自动赔付 5 元无门槛券",
    receiver = "张三 138****8888",
    address = "上海市浦东新区世纪大道 100 号 8 楼",
    item = OrderDetailItem(
        title = "轻薄羽绒服 女款 短款连帽 冬季保暖外套",
        spec = "米白色 / M",
        price = 299.0,
        count = 1,
    ),
    goodsTotal = 299.0,
    freight = 0.0,
    discount = 30.0,
    paid = 269.0,
    orderNo = "2026102012345678",
    createdAt = "2026-10-20 10:24",
    payMethod = "微信支付",
)

private fun yuan(value: Double): String = "¥" + String.format(Locale.US, "%.2f", value)

/** 订单状态卡片: status tag + title + countdown, description, receiver address. */
@Composable
private fun OrderStatusCard(detail: OrderDetail) {
    DsCard(modifier = Modifier.fillMaxWidth().figmaId("1:6")) {
        Row(
            modifier = Modifier.fillMaxWidth().figmaId("1:7"),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                DsBadge(text = detail.statusText, tone = detail.statusTone, modifier = Modifier.figmaId("1:8"))
                Spacer(Modifier.width(DsSpacing.Sm))
                DsText(text = detail.statusTitle, style = DsTypography.Title, modifier = Modifier.figmaId("1:10"))
            }
            DsText(
                text = detail.countdown,
                style = DsTypography.Caption,
                color = DsColors.Danger,
                modifier = Modifier.figmaId("1:11"),
            )
        }
        Spacer(Modifier.height(DsSpacing.Xs))
        DsText(
            text = detail.statusDesc,
            style = DsTypography.Caption,
            color = DsColors.TextTertiary,
            modifier = Modifier.fillMaxWidth().figmaId("1:13"),
        )
        Spacer(Modifier.height(DsSpacing.Md))
        DsDivider(modifier = Modifier.figmaId("1:14"))
        DsListItem(
            title = detail.receiver,
            subtitle = detail.address,
            showChevron = true,
            modifier = Modifier.figmaId("1:15"),
        )
    }
}

/** 商品信息卡片: goods row, amount breakdown, and 实付款. */
@Composable
private fun GoodsCard(detail: OrderDetail) {
    DsCard(modifier = Modifier.fillMaxWidth().figmaId("1:18"), title = "商品信息") {
        DsGoodsCard(
            title = detail.item.title,
            price = detail.item.price,
            spec = detail.item.spec,
            count = detail.item.count,
            modifier = Modifier.figmaId("1:20"),
        )
        Spacer(Modifier.height(DsSpacing.Md))
        DsDivider(modifier = Modifier.figmaId("1:29"))
        DsListItem(title = "商品总额", trailingText = yuan(detail.goodsTotal), modifier = Modifier.figmaId("1:30"))
        DsListItem(
            title = "运费",
            trailingText = if (detail.freight > 0) yuan(detail.freight) else "免运费",
            modifier = Modifier.figmaId("1:33"),
        )
        DsListItem(title = "店铺优惠", trailingText = "-" + yuan(detail.discount), modifier = Modifier.figmaId("1:36"))
        Spacer(Modifier.height(DsSpacing.Xs))
        Row(
            modifier = Modifier.fillMaxWidth().height(32.dp).figmaId("1:39"),
            horizontalArrangement = Arrangement.End,
            verticalAlignment = Alignment.CenterVertically,
        ) {
            DsText(text = "实付款", style = DsTypography.Body, modifier = Modifier.figmaId("1:40"))
            Spacer(Modifier.width(DsSpacing.Sm))
            DsPriceText(
                amount = detail.paid,
                size = DsPriceSize.Large,
                tone = DsPriceTone.Primary,
                modifier = Modifier.figmaId("1:41"),
            )
        }
    }
}

/** 订单信息卡片: order number, order time, payment method. */
@Composable
private fun OrderInfoCard(detail: OrderDetail) {
    DsCard(modifier = Modifier.fillMaxWidth().figmaId("1:45"), title = "订单信息") {
        DsListItem(title = "订单编号", trailingText = detail.orderNo, modifier = Modifier.figmaId("1:47"))
        DsListItem(title = "下单时间", trailingText = detail.createdAt, modifier = Modifier.figmaId("1:50"))
        DsListItem(title = "支付方式", trailingText = detail.payMethod, modifier = Modifier.figmaId("1:53"))
    }
}

/** 订单详情 — Figma: order-detail (375×900). Built only from com.demo.ds components. */
@Composable
fun OrderDetailScreen(detail: OrderDetail = sampleOrderDetail) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(DsColors.BgPage),
    ) {
        DsTopBar(title = "订单详情", onBack = {}, modifier = Modifier.figmaId("1:2"))
        Column(
            modifier = Modifier
                .weight(1f)
                .fillMaxWidth()
                .verticalScroll(rememberScrollState())
                // 36dp bottom keeps the last card 36dp above the action bar, as in the frame.
                .padding(start = DsSpacing.Md, top = DsSpacing.Md, end = DsSpacing.Md, bottom = DsSpacing.Xl + DsSpacing.Md),
            verticalArrangement = Arrangement.spacedBy(DsSpacing.Md),
        ) {
            DsBanner(text = detail.notice, tone = DsBannerTone.Warning, modifier = Modifier.figmaId("1:4"))
            OrderStatusCard(detail)
            GoodsCard(detail)
            OrderInfoCard(detail)
        }
        DsBottomBar(modifier = Modifier.figmaId("1:56")) {
            DsButton(text = "联系客服", onClick = {}, style = DsButtonStyle.Default, modifier = Modifier.figmaId("1:57"))
            DsButton(text = "修改地址", onClick = {}, style = DsButtonStyle.Default, modifier = Modifier.figmaId("1:59"))
            DsButton(text = "提醒发货", onClick = {}, style = DsButtonStyle.Primary, modifier = Modifier.figmaId("1:61"))
        }
    }
}
