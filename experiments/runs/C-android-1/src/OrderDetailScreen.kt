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

/** Height of [DsBottomBar]; the scrolling content reserves this much space at the bottom. */
private val BottomBarHeight = 56.dp

/** 订单详情（Figma frame 1:1）— built only from com.demo.ds components and tokens. */
@Composable
fun OrderDetailScreen() {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(DsColors.BgPage),
    ) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(bottom = BottomBarHeight),
        ) {
            DsTopBar(title = "订单详情", onBack = {}, modifier = Modifier.figmaId("1:2"))
            Column(
                modifier = Modifier.fillMaxWidth().padding(DsSpacing.Md),
                verticalArrangement = Arrangement.spacedBy(DsSpacing.Md),
            ) {
                DsBanner(
                    text = "商家已接单，预计 10月22日 前发货",
                    tone = DsBannerTone.Warning,
                    modifier = Modifier.figmaId("1:4"),
                )
                StatusSection()
                GoodsSection()
                OrderInfoSection()
            }
        }

        DsBottomBar(modifier = Modifier.align(Alignment.BottomCenter).figmaId("1:56")) {
            DsButton(
                text = "联系客服",
                onClick = {},
                size = DsButtonSize.Medium,
                style = DsButtonStyle.Default,
                modifier = Modifier.figmaId("1:57"),
            )
            DsButton(
                text = "修改地址",
                onClick = {},
                size = DsButtonSize.Medium,
                style = DsButtonStyle.Default,
                modifier = Modifier.figmaId("1:59"),
            )
            DsButton(
                text = "提醒发货",
                onClick = {},
                size = DsButtonSize.Medium,
                style = DsButtonStyle.Primary,
                modifier = Modifier.figmaId("1:61"),
            )
        }
    }
}

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
                style = DsTypography.Title,
                color = DsColors.TextPrimary,
                maxLines = 1,
                modifier = Modifier.figmaId("1:10"),
            )
            Spacer(Modifier.weight(1f))
            // CountDown: no component on either platform yet (component-map gap) → static caption text.
            DsText(
                text = "剩余 23:59:59",
                style = DsTypography.Caption,
                color = DsColors.Danger,
                maxLines = 1,
                modifier = Modifier.figmaId("1:11"),
            )
        }
        Spacer(Modifier.height(DsSpacing.Xs))
        DsText(
            text = "超时未发货将自动赔付 5 元无门槛券",
            style = DsTypography.Caption,
            color = DsColors.TextTertiary,
            modifier = Modifier.fillMaxWidth().figmaId("1:13"),
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

@Composable
private fun GoodsSection() {
    DsCard(title = "商品信息", modifier = Modifier.fillMaxWidth().figmaId("1:18")) {
        DsGoodsCard(
            title = "轻薄羽绒服 女款 短款连帽 冬季保暖外套",
            price = 299.00,
            count = 1,
            spec = "米白色 / M",
            modifier = Modifier.figmaId("1:20"),
        )
        Spacer(Modifier.height(DsSpacing.Md))
        DsDivider(modifier = Modifier.figmaId("1:29"))
        DsListItem(title = "商品总额", trailingText = "¥299.00", modifier = Modifier.figmaId("1:30"))
        DsListItem(title = "运费", trailingText = "免运费", modifier = Modifier.figmaId("1:33"))
        DsListItem(title = "店铺优惠", trailingText = "-¥30.00", modifier = Modifier.figmaId("1:36"))
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(top = DsSpacing.Xs)
                .height(32.dp)
                .figmaId("1:39"),
            horizontalArrangement = Arrangement.spacedBy(DsSpacing.Sm, Alignment.End),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            DsText(
                text = "实付款",
                style = DsTypography.Body,
                color = DsColors.TextPrimary,
                maxLines = 1,
                modifier = Modifier.figmaId("1:40"),
            )
            DsPriceText(
                amount = 269.00,
                size = DsPriceSize.Large,
                tone = DsPriceTone.Primary,
                modifier = Modifier.figmaId("1:41"),
            )
        }
    }
}

@Composable
private fun OrderInfoSection() {
    DsCard(title = "订单信息", modifier = Modifier.fillMaxWidth().figmaId("1:45")) {
        DsListItem(title = "订单编号", trailingText = "2026102012345678", modifier = Modifier.figmaId("1:47"))
        DsListItem(title = "下单时间", trailingText = "2026-10-20 10:24", modifier = Modifier.figmaId("1:50"))
        DsListItem(title = "支付方式", trailingText = "微信支付", modifier = Modifier.figmaId("1:53"))
    }
}
