package com.demo.ds

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

/**
 * Goods row (Figma: `ProductCard`). 80×80 placeholder thumbnail (radius.md, bg-placeholder),
 * 12dp gap, then an 80dp-tall column: title (font.body, ≤2 lines, ellipsis), optional spec
 * (font.caption, text-tertiary, 4dp below), and a bottom row pinned to the bottom with
 * `DsPriceText(Medium, Default)` on the left and "x{count}" (font.caption, text-tertiary) on the right.
 *
 * @sample
 * ```
 * DsGoodsCard(title = "无线降噪耳机 Pro", price = 299.0, spec = "星空黑 / 标准版", count = 1)
 * ```
 */
public @Composable fun DsGoodsCard(
    title: String,
    price: Double,
    modifier: Modifier = Modifier,
    spec: String? = null,
    count: Int = 1,
) {
    Row(modifier = modifier.dsComponent("DsGoodsCard").fillMaxWidth()) {
        Box(
            Modifier
                .size(80.dp)
                .background(DsColors.BgPlaceholder, RoundedCornerShape(DsRadius.Md)),
        )
        Spacer(Modifier.width(DsSpacing.Md))
        // Box (not a Column with a weighted spacer) so the bottom row stays pinned to the bottom
        // even when a 2-line title + spec fill the 80dp height (avoids squeezing / clipping the price).
        Box(Modifier.weight(1f).height(80.dp)) {
            Column(Modifier.fillMaxWidth()) {
                DsTextInternal(text = title, style = DsTypography.Body, color = DsColors.TextPrimary, maxLines = 2)
                if (spec != null) {
                    Spacer(Modifier.height(DsSpacing.Xs))
                    DsTextInternal(text = spec, style = DsTypography.Caption, color = DsColors.TextTertiary, maxLines = 1)
                }
            }
            Row(Modifier.fillMaxWidth().align(Alignment.BottomStart)) {
                DsPriceText(amount = price, size = DsPriceSize.Medium, tone = DsPriceTone.Default, modifier = Modifier.alignByBaseline())
                Spacer(Modifier.weight(1f))
                DsTextInternal(
                    text = "x$count",
                    style = DsTypography.Caption,
                    color = DsColors.TextTertiary,
                    maxLines = 1,
                    modifier = Modifier.alignByBaseline(),
                )
            }
        }
    }
}

public @Composable fun DsGoodsCardPreview() {
    DsGoodsCard(title = "无线降噪耳机 Pro 主动降噪 蓝牙5.3 超长续航", price = 299.0, spec = "星空黑 / 标准版", count = 1)
}
