package com.demo.ds

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Path
import androidx.compose.ui.graphics.StrokeCap
import androidx.compose.ui.graphics.StrokeJoin
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.unit.dp

/**
 * Top navigation bar (Figma: `NavBar`). Height 44, full width, bg-card,
 * optional 24×24 back chevron 12dp from the left, centered title (font.title).
 *
 * @sample
 * ```
 * DsTopBar(title = "订单详情", onBack = { navController.popBackStack() })
 * ```
 */
public @Composable fun DsTopBar(
    title: String,
    modifier: Modifier = Modifier,
    onBack: (() -> Unit)? = null,
) {
    Box(
        modifier = modifier
            .dsComponent("DsTopBar")
            .fillMaxWidth()
            .height(44.dp)
            .background(DsColors.BgCard),
    ) {
        if (onBack != null) {
            DsChevron(
                left = true,
                strokeWidthDp = 2f,
                modifier = Modifier
                    .align(Alignment.CenterStart)
                    .padding(start = DsSpacing.Md)
                    .size(24.dp)
                    .clickable(onClick = onBack),
            )
        }
        DsTextInternal(
            text = title,
            style = DsTypography.Title,
            color = DsColors.TextPrimary,
            maxLines = 1,
            modifier = Modifier.align(Alignment.Center).padding(horizontal = 48.dp),
        )
    }
}

/** Chevron icon drawn with a Path on a square canvas (viewBox 24). */
@Composable
internal fun DsChevron(left: Boolean, strokeWidthDp: Float, modifier: Modifier = Modifier, color: androidx.compose.ui.graphics.Color = DsColors.TextPrimary) {
    Canvas(modifier) {
        val s = size.minDimension / 24f
        val path = Path().apply {
            if (left) {
                moveTo(15f * s, 6f * s); lineTo(9f * s, 12f * s); lineTo(15f * s, 18f * s)
            } else {
                moveTo(9f * s, 6f * s); lineTo(15f * s, 12f * s); lineTo(9f * s, 18f * s)
            }
        }
        drawPath(path, color, style = Stroke(width = strokeWidthDp.dp.toPx(), cap = StrokeCap.Round, join = StrokeJoin.Round))
    }
}

public @Composable fun DsTopBarPreview() {
    DsTopBar(title = "订单详情", onBack = {})
}
