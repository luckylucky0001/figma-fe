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

/** 订单详情 — NavBar 固定在顶部，内容区可滚动，ActionBar 吸底。 */
@Composable
fun OrderDetailScreen() {
    Box(Modifier.fillMaxSize().background(DsColors.BgPage)) {
        Column(Modifier.fillMaxSize()) {
            DsTopBar(title = "订单详情", modifier = Modifier.figmaId("1:2"), onBack = {})
            Column(
                modifier = Modifier
                    .weight(1f)
                    .fillMaxWidth()
                    .verticalScroll(rememberScrollState())
                    .padding(start = DsSpacing.Md, top = DsSpacing.Md, end = DsSpacing.Md, bottom = 56.dp),
                verticalArrangement = Arrangement.spacedBy(DsSpacing.Md),
            ) {
                DsBanner(
                    text = "商家已接单，预计 10月22日 前发货",
                    tone = DsBannerTone.Warning,
                    modifier = Modifier.figmaId("1:4"),
                )
                StatusSection()
                GoodsSection()
                InfoSection()
            }
        }
        DsBottomBar(modifier = Modifier.align(Alignment.BottomCenter).figmaId("1:56")) {
            DsButton(text = "联系客服", onClick = {}, style = DsButtonStyle.Default, modifier = Modifier.figmaId("1:57"))
            DsButton(text = "修改地址", onClick = {}, style = DsButtonStyle.Default, modifier = Modifier.figmaId("1:59"))
            DsButton(text = "提醒发货", onClick = {}, style = DsButtonStyle.Primary, modifier = Modifier.figmaId("1:61"))
        }
    }
}

/** 状态与收货人信息。 */
@Composable
private fun StatusSection() {
    DsCard(modifier = Modifier.fillMaxWidth().figmaId("1:6")) {
        Row(
            modifier = Modifier.fillMaxWidth().figmaId("1:7"),
            horizontalArrangement = Arrangement.spacedBy(DsSpacing.Sm),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            DsBadge(text = "待发货", tone = DsBadgeTone.Warning, modifier = Modifier.figmaId("1:8"))
            DsText(
                text = "等待商家发货",
                modifier = Modifier.figmaId("1:10"),
                style = DsTypography.Title,
                color = DsColors.TextPrimary,
            )
            Spacer(Modifier.weight(1f))
            DsText(
                text = "剩余 23:59:59",
                modifier = Modifier.figmaId("1:11"),
                style = DsTypography.Caption,
                color = DsColors.Danger,
            )
        }
        Spacer(Modifier.height(DsSpacing.Xs))
        DsText(
            text = "超时未发货将自动赔付 5 元无门槛券",
            modifier = Modifier.fillMaxWidth().figmaId("1:13"),
            style = DsTypography.Caption,
            color = DsColors.TextTertiary,
        )
        Spacer(Modifier.height(DsSpacing.Md))
        DsDivider(modifier = Modifier.figmaId("1:14"))
        DsListItem(
            title = "张三 138****8888",
            subtitle = "上海市浦东新区世纪大道 100 号 8 楼",
            showChevron = true,
            modifier = Modifier.figmaId("1:15"),
        )
    }
}

/** 商品与金额明细。 */
@Composable
private fun GoodsSection() {
    DsCard(modifier = Modifier.fillMaxWidth().figmaId("1:18"), title = "商品信息") {
        DsGoodsCard(
            title = "轻薄羽绒服 女款 短款连帽 冬季保暖外套",
            price = 299.0,
            spec = "米白色 / M",
            count = 1,
            modifier = Modifier.figmaId("1:20"),
        )
        Spacer(Modifier.height(DsSpacing.Md))
        DsDivider(modifier = Modifier.figmaId("1:29"))
        DsListItem(title = "商品总额", trailingText = "¥299.00", modifier = Modifier.figmaId("1:30"))
        DsListItem(title = "运费", trailingText = "免运费", modifier = Modifier.figmaId("1:33"))
        DsListItem(title = "店铺优惠", trailingText = "-¥30.00", modifier = Modifier.figmaId("1:36"))
        Spacer(Modifier.height(DsSpacing.Xs))
        Row(
            modifier = Modifier.fillMaxWidth().height(32.dp).figmaId("1:39"),
            horizontalArrangement = Arrangement.spacedBy(DsSpacing.Sm, Alignment.End),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            DsText(
                text = "实付款",
                modifier = Modifier.figmaId("1:40"),
                style = DsTypography.Body,
                color = DsColors.TextPrimary,
            )
            DsPriceText(
                amount = 269.0,
                modifier = Modifier.figmaId("1:41"),
                size = DsPriceSize.Large,
                tone = DsPriceTone.Primary,
            )
        }
    }
}

/** 订单编号、下单时间、支付方式。 */
@Composable
private fun InfoSection() {
    DsCard(modifier = Modifier.fillMaxWidth().figmaId("1:45"), title = "订单信息") {
        DsListItem(title = "订单编号", trailingText = "2026102012345678", modifier = Modifier.figmaId("1:47"))
        DsListItem(title = "下单时间", trailingText = "2026-10-20 10:24", modifier = Modifier.figmaId("1:50"))
        DsListItem(title = "支付方式", trailingText = "微信支付", modifier = Modifier.figmaId("1:53"))
    }
}
