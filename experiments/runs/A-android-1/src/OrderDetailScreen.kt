package com.demo.app.screens.OrderDetail

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.demo.app.mock.OrderDetail
import com.demo.app.mock.OrderStatus
import com.demo.app.mock.orderDetail
import com.demo.ds.DsBadge
import com.demo.ds.DsBadgeTone
import com.demo.ds.DsBanner
import com.demo.ds.DsBannerTone
import com.demo.ds.DsBottomBar
import com.demo.ds.DsButton
import com.demo.ds.DsButtonSize
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
import kotlinx.coroutines.delay
import java.math.BigDecimal
import java.math.RoundingMode

private val ActionBarHeight = 56.dp

private fun money(amount: Double): String =
    "¥" + BigDecimal.valueOf(amount).setScale(2, RoundingMode.HALF_UP).toPlainString()

private fun OrderStatus.badgeTone(): DsBadgeTone = when (this) {
    OrderStatus.PendingPayment -> DsBadgeTone.Primary
    OrderStatus.PendingShipment -> DsBadgeTone.Warning
    OrderStatus.Shipped, OrderStatus.Completed -> DsBadgeTone.Success
    OrderStatus.Cancelled -> DsBadgeTone.Danger
}

private fun formatHms(totalSeconds: Long): String {
    val s = totalSeconds.coerceAtLeast(0)
    return "%02d:%02d:%02d".format(s / 3600, s % 3600 / 60, s % 60)
}

/**
 * Figma `CountDown` (10:11) — no counterpart in com.demo.ds yet, so it is built locally from
 * DsText: 「剩余 HH:MM:SS」, font.caption, color.danger. Ticks once per second.
 */
@Composable
private fun CountDown(remainSeconds: Long, modifier: Modifier = Modifier) {
    var left by remember(remainSeconds) { mutableStateOf(remainSeconds) }
    LaunchedEffect(remainSeconds) {
        while (left > 0) {
            delay(1_000)
            left -= 1
        }
    }
    DsText(
        text = "剩余 ${formatHms(left)}",
        style = DsTypography.Caption,
        color = DsColors.Danger,
        maxLines = 1,
        modifier = modifier,
    )
}

/** Card / 状态 (1:6): status row + description, divider, receiver address cell. */
@Composable
private fun StatusCard(detail: OrderDetail) {
    DsCard(modifier = Modifier.fillMaxWidth().figmaId("1:6")) {
        Row(
            modifier = Modifier.fillMaxWidth().height(24.dp).figmaId("1:7"),
            horizontalArrangement = Arrangement.spacedBy(DsSpacing.Sm),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            DsBadge(text = detail.status.text, tone = detail.status.badgeTone(), modifier = Modifier.figmaId("1:8"))
            DsText(
                text = detail.statusTitle,
                style = DsTypography.Title,
                maxLines = 1,
                modifier = Modifier.weight(1f).figmaId("1:10"),
            )
            CountDown(remainSeconds = detail.remainSeconds, modifier = Modifier.figmaId("1:11"))
        }
        DsText(
            text = detail.statusDesc,
            style = DsTypography.Caption,
            color = DsColors.TextTertiary,
            modifier = Modifier.padding(top = DsSpacing.Xs).figmaId("1:13"),
        )
        Spacer(Modifier.height(DsSpacing.Md))
        DsDivider(modifier = Modifier.figmaId("1:14"))
        DsListItem(
            title = "${detail.receiver} ${detail.phone}",
            subtitle = detail.address,
            showChevron = true,
            modifier = Modifier.figmaId("1:15"),
        )
    }
}

/** Card / 商品 (1:18): goods, amount breakdown and the paid total. */
@Composable
private fun GoodsCard(detail: OrderDetail) {
    DsCard(modifier = Modifier.fillMaxWidth().figmaId("1:18"), title = "商品信息") {
        Column(verticalArrangement = Arrangement.spacedBy(DsSpacing.Md)) {
            detail.items.forEach { item ->
                DsGoodsCard(
                    title = item.title,
                    price = item.price,
                    spec = item.spec,
                    count = item.count,
                    modifier = Modifier.figmaId("1:20"),
                )
            }
        }
        Spacer(Modifier.height(DsSpacing.Md))
        DsDivider(modifier = Modifier.figmaId("1:29"))
        DsListItem(title = "商品总额", trailingText = money(detail.goodsAmount), modifier = Modifier.figmaId("1:30"))
        DsListItem(
            title = "运费",
            trailingText = if (detail.freight > 0) money(detail.freight) else "免运费",
            modifier = Modifier.figmaId("1:33"),
        )
        DsListItem(title = "店铺优惠", trailingText = "-" + money(detail.discount), modifier = Modifier.figmaId("1:36"))
        Row(
            modifier = Modifier.fillMaxWidth().padding(top = DsSpacing.Xs).height(32.dp).figmaId("1:39"),
            horizontalArrangement = Arrangement.spacedBy(DsSpacing.Sm, Alignment.End),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            DsText(text = "实付款", maxLines = 1, modifier = Modifier.figmaId("1:40"))
            DsPriceText(
                amount = detail.total,
                size = DsPriceSize.Large,
                tone = DsPriceTone.Primary,
                modifier = Modifier.figmaId("1:41"),
            )
        }
    }
}

/** Card / 订单信息 (1:45). */
@Composable
private fun OrderInfoCard(detail: OrderDetail) {
    DsCard(modifier = Modifier.fillMaxWidth().figmaId("1:45"), title = "订单信息") {
        DsListItem(title = "订单编号", trailingText = detail.id, modifier = Modifier.figmaId("1:47"))
        DsListItem(title = "下单时间", trailingText = detail.createdAt, modifier = Modifier.figmaId("1:50"))
        DsListItem(title = "支付方式", trailingText = detail.payMethod, modifier = Modifier.figmaId("1:53"))
    }
}

/** 订单详情 — Figma frame `OrderDetail 订单详情` (1:1), built from com.demo.ds components. */
@Composable
fun OrderDetailScreen(
    detail: OrderDetail = orderDetail,
    onBack: () -> Unit = {},
    onContactService: () -> Unit = {},
    onChangeAddress: () -> Unit = {},
    onRemindShipment: () -> Unit = {},
) {
    Box(Modifier.fillMaxSize().background(DsColors.BgPage).figmaId("1:1")) {
        Column(Modifier.fillMaxSize()) {
            DsTopBar(title = "订单详情", onBack = onBack, modifier = Modifier.figmaId("1:2"))
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .weight(1f)
                    .verticalScroll(rememberScrollState())
                    .padding(start = DsSpacing.Md, end = DsSpacing.Md, top = DsSpacing.Md, bottom = DsSpacing.Md + ActionBarHeight),
                verticalArrangement = Arrangement.spacedBy(DsSpacing.Md),
            ) {
                DsBanner(text = detail.notice, tone = DsBannerTone.Warning, modifier = Modifier.figmaId("1:4"))
                StatusCard(detail)
                GoodsCard(detail)
                OrderInfoCard(detail)
            }
        }
        DsBottomBar(modifier = Modifier.align(Alignment.BottomCenter).figmaId("1:56")) {
            DsButton(
                text = "联系客服",
                onClick = onContactService,
                style = DsButtonStyle.Default,
                size = DsButtonSize.Medium,
                modifier = Modifier.figmaId("1:57"),
            )
            DsButton(
                text = "修改地址",
                onClick = onChangeAddress,
                style = DsButtonStyle.Default,
                size = DsButtonSize.Medium,
                modifier = Modifier.figmaId("1:59"),
            )
            DsButton(
                text = "提醒发货",
                onClick = onRemindShipment,
                style = DsButtonStyle.Primary,
                size = DsButtonSize.Medium,
                modifier = Modifier.figmaId("1:61"),
            )
        }
    }
}
