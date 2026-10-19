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
import androidx.compose.foundation.layout.width
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
import com.demo.app.mock.orderDetail
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
import kotlinx.coroutines.delay
import java.util.Locale

private val BottomBarHeight = 56.dp

private fun money(v: Double): String = "¥" + String.format(Locale.US, "%.2f", v)

/**
 * Figma `CountDown` has no ds implementation on either platform yet (see design/SPEC.md §2),
 * so it is built locally from DsText: 「剩余 HH:MM:SS」, font.caption, color.danger.
 */
@Composable
private fun CountDown(seconds: Long, modifier: Modifier = Modifier) {
    var left by remember(seconds) { mutableStateOf(seconds) }
    LaunchedEffect(seconds) {
        while (left > 0) {
            delay(1000)
            left -= 1
        }
    }
    val text = String.format(Locale.US, "剩余 %02d:%02d:%02d", left / 3600, left % 3600 / 60, left % 60)
    DsText(text = text, style = DsTypography.Caption, color = DsColors.Danger, maxLines = 1, modifier = modifier)
}

@Composable
private fun StatusCard(order: OrderDetail) {
    DsCard(modifier = Modifier.fillMaxWidth()) {
        Row(modifier = Modifier.fillMaxWidth(), verticalAlignment = Alignment.CenterVertically) {
            DsBadge(text = order.status.text, tone = DsBadgeTone.Warning)
            Spacer(Modifier.width(DsSpacing.Sm))
            DsText(text = order.statusTitle, style = DsTypography.Title, maxLines = 1, modifier = Modifier.weight(1f))
            Spacer(Modifier.width(DsSpacing.Sm))
            CountDown(seconds = order.remainingSeconds)
        }
        Spacer(Modifier.height(DsSpacing.Xs))
        DsText(text = order.statusDesc, style = DsTypography.Caption, color = DsColors.TextTertiary)
        Spacer(Modifier.height(DsSpacing.Md))
        DsDivider()
        DsListItem(title = order.receiver, subtitle = order.address, showChevron = true)
    }
}

@Composable
private fun GoodsCard(order: OrderDetail) {
    DsCard(modifier = Modifier.fillMaxWidth(), title = "商品信息") {
        DsGoodsCard(title = order.item.title, price = order.item.price, spec = order.item.spec, count = order.item.count)
        Spacer(Modifier.height(DsSpacing.Md))
        DsDivider()
        DsListItem(title = "商品总额", trailingText = money(order.goodsTotal))
        DsListItem(title = "运费", trailingText = if (order.freight > 0) money(order.freight) else "免运费")
        DsListItem(title = "店铺优惠", trailingText = "-" + money(order.shopDiscount))
        Row(
            modifier = Modifier.fillMaxWidth().padding(top = DsSpacing.Xs),
            horizontalArrangement = Arrangement.spacedBy(DsSpacing.Sm, Alignment.End),
            // The design centers the label on the price box (not baseline-aligned).
            verticalAlignment = Alignment.CenterVertically,
        ) {
            DsText(text = "实付款")
            DsPriceText(amount = order.total, size = DsPriceSize.Large, tone = DsPriceTone.Primary)
        }
    }
}

@Composable
private fun InfoCard(order: OrderDetail) {
    DsCard(modifier = Modifier.fillMaxWidth(), title = "订单信息") {
        DsListItem(title = "订单编号", trailingText = order.id)
        DsListItem(title = "下单时间", trailingText = order.createdAt)
        DsListItem(title = "支付方式", trailingText = order.payMethod)
    }
}

/** 订单详情 — design/figma-export/order-detail, built only from com.demo.ds components. */
@Composable
fun OrderDetailScreen(order: OrderDetail = orderDetail) {
    Box(Modifier.fillMaxSize().background(DsColors.BgPage)) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(bottom = BottomBarHeight),
        ) {
            DsTopBar(title = "订单详情", onBack = {})
            Column(
                modifier = Modifier.fillMaxWidth().padding(DsSpacing.Md),
                verticalArrangement = Arrangement.spacedBy(DsSpacing.Md),
            ) {
                DsBanner(text = order.notice, tone = DsBannerTone.Warning)
                StatusCard(order)
                GoodsCard(order)
                InfoCard(order)
            }
        }

        DsBottomBar(Modifier.align(Alignment.BottomCenter)) {
            DsButton(text = "联系客服", onClick = {}, style = DsButtonStyle.Default)
            DsButton(text = "修改地址", onClick = {}, style = DsButtonStyle.Default)
            DsButton(text = "提醒发货", onClick = {}, style = DsButtonStyle.Primary)
        }
    }
}
