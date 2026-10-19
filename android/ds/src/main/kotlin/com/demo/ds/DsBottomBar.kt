package com.demo.ds

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.RowScope
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.drawBehind
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.unit.dp

/**
 * Bottom action bar (Figma: `ActionBar`). Height 56, bg-card, 1dp top border (color.border),
 * horizontal padding 16, children End-aligned with 8dp spacing, vertically centered.
 * It does not position itself — the screen places it at the bottom.
 *
 * @sample
 * ```
 * Box(Modifier.fillMaxSize()) {
 *     DsBottomBar(Modifier.align(Alignment.BottomCenter)) {
 *         DsButton(text = "取消订单", onClick = {}, style = DsButtonStyle.Default)
 *         DsButton(text = "立即支付", onClick = {})
 *     }
 * }
 * ```
 */
public @Composable fun DsBottomBar(
    modifier: Modifier = Modifier,
    content: @Composable RowScope.() -> Unit,
) {
    Row(
        modifier = modifier
            .dsComponent("DsBottomBar")
            .fillMaxWidth()
            .height(56.dp)
            .background(DsColors.BgCard)
            .drawBehind {
                drawRect(DsColors.Border, topLeft = Offset.Zero, size = Size(size.width, 1.dp.toPx()))
            }
            .padding(horizontal = DsSpacing.Lg),
        horizontalArrangement = Arrangement.spacedBy(DsSpacing.Sm, Alignment.End),
        verticalAlignment = Alignment.CenterVertically,
        content = content,
    )
}

public @Composable fun DsBottomBarPreview() {
    DsBottomBar {
        DsButton(text = "取消订单", onClick = {}, style = DsButtonStyle.Default)
        DsButton(text = "立即支付", onClick = {})
    }
}
