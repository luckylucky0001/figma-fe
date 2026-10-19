package com.demo.app.screens

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ColumnScope
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
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

@Composable
private fun Section(name: String, content: @Composable ColumnScope.() -> Unit) {
    Column(Modifier.fillMaxWidth().padding(start = DsSpacing.Md, end = DsSpacing.Md, top = DsSpacing.Lg)) {
        DsText(text = name, style = DsTypography.Subtitle, modifier = Modifier.padding(bottom = DsSpacing.Sm))
        Column(verticalArrangement = Arrangement.spacedBy(DsSpacing.Md), content = content)
    }
}

@Composable
private fun Variant(label: String, content: @Composable () -> Unit) {
    Column {
        DsText(text = label, style = DsTypography.Caption, color = DsColors.TextTertiary, modifier = Modifier.padding(bottom = 6.dp))
        content()
    }
}

private val badgeSamples = mapOf(
    DsBadgeTone.Primary to "待付款",
    DsBadgeTone.Success to "已发货",
    DsBadgeTone.Warning to "待发货",
    DsBadgeTone.Danger to "已取消",
)

/** Gallery — every ds component in every implemented variant, for visual QA. */
@Composable
fun GalleryScreen() {
    Box(Modifier.fillMaxSize().background(DsColors.BgPage)) {
        Column(
            Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(bottom = 56.dp + DsSpacing.Lg),
        ) {
            Section("DsTopBar") {
                Variant("default") { DsTopBar(title = "订单详情", onBack = {}) }
            }

            Section("DsCard") {
                Variant("title") {
                    DsCard(modifier = Modifier.fillMaxWidth(), title = "收货信息") { DsText(text = "卡片内容区域") }
                }
                Variant("no title") {
                    DsCard(modifier = Modifier.fillMaxWidth()) { DsText(text = "卡片内容区域") }
                }
            }

            Section("DsListItem") {
                DsCard(modifier = Modifier.fillMaxWidth()) {
                    DsListItem(title = "订单编号", trailingText = "DD20261018093021")
                    DsDivider()
                    DsListItem(title = "配送方式", trailingText = "快递", showChevron = true)
                    DsDivider()
                    DsListItem(title = "发票", subtitle = "电子普通发票 · 个人", trailingText = "已开具")
                    DsDivider()
                    DsListItem(title = "优惠券", subtitle = "满 300 减 30", trailingText = "-¥30.00", showChevron = true)
                    DsDivider()
                    DsListItem(title = "仅标题")
                }
            }

            Section("DsBadge") {
                Variant("Light (only variant on Android)") {
                    Row(horizontalArrangement = Arrangement.spacedBy(DsSpacing.Sm)) {
                        badgeSamples.forEach { (tone, text) -> DsBadge(text = text, tone = tone) }
                    }
                }
            }

            Section("DsPriceText") {
                DsPriceTone.values().forEach { tone ->
                    Variant("tone=${tone.name} · Large / Medium / Small") {
                        Row(horizontalArrangement = Arrangement.spacedBy(DsSpacing.Lg)) {
                            DsPriceText(amount = 1299.0, size = DsPriceSize.Large, tone = tone, modifier = Modifier.alignByBaseline())
                            DsPriceText(amount = 59.9, size = DsPriceSize.Medium, tone = tone, modifier = Modifier.alignByBaseline())
                            DsPriceText(amount = 8.05, size = DsPriceSize.Small, tone = tone, modifier = Modifier.alignByBaseline())
                        }
                    }
                }
            }

            Section("DsButton") {
                DsButtonStyle.values().forEach { style ->
                    Variant("style=${style.name} · Large / Medium / Small") {
                        Row(
                            horizontalArrangement = Arrangement.spacedBy(DsSpacing.Sm),
                            verticalAlignment = Alignment.CenterVertically,
                        ) {
                            DsButtonSize.values().forEach { size ->
                                DsButton(
                                    text = if (style == DsButtonStyle.Text) "查看详情" else "提醒发货",
                                    onClick = {},
                                    style = style,
                                    size = size,
                                )
                            }
                        }
                    }
                }
                Variant("full width") {
                    DsButton(text = "立即支付", onClick = {}, size = DsButtonSize.Large, modifier = Modifier.fillMaxWidth())
                }
            }

            Section("DsDivider") {
                DsCard(modifier = Modifier.fillMaxWidth()) {
                    DsText(text = "上方内容", modifier = Modifier.padding(bottom = DsSpacing.Md))
                    DsDivider()
                    DsText(text = "下方内容", modifier = Modifier.padding(top = DsSpacing.Md))
                }
            }

            Section("DsGoodsCard") {
                DsCard(modifier = Modifier.fillMaxWidth()) {
                    DsGoodsCard(
                        title = "无线降噪蓝牙耳机 Pro 主动降噪 长续航 入耳式运动耳机 超长标题测试两行省略",
                        price = 899.0,
                        spec = "曜石黑 / 标准版",
                        count = 1,
                    )
                }
                DsCard(modifier = Modifier.fillMaxWidth()) {
                    DsGoodsCard(title = "北欧风陶瓷马克杯", price = 39.9, count = 2)
                }
            }

            Section("DsBanner") {
                Variant("tone=Info") {
                    DsBanner(text = "订单已提交，请在 30 分钟内完成支付", tone = DsBannerTone.Info)
                }
                Variant("tone=Warning (overflow)") {
                    DsBanner(text = "受天气影响，部分地区配送可能延迟 1-2 天，请耐心等待，给您带来不便敬请谅解", tone = DsBannerTone.Warning)
                }
            }

            Section("DsBottomBar") {
                DsText(text = "pinned to the bottom of the screen ↓", style = DsTypography.Caption, color = DsColors.TextTertiary)
            }
        }

        DsBottomBar(Modifier.align(Alignment.BottomCenter)) {
            DsButton(text = "取消订单", onClick = {}, style = DsButtonStyle.Default)
            DsButton(text = "去支付", onClick = {}, style = DsButtonStyle.Primary)
        }
    }
}
