package com.demo.ds

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ColumnScope
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip

/**
 * Card container (Figma: `Card`, variant `Title=On|Off`). bg-card, radius.lg (12), padding 16.
 * Optional title (font.subtitle) followed by a 12dp gap; children stacked vertically with no spacing.
 *
 * @sample
 * ```
 * DsCard(modifier = Modifier.fillMaxWidth(), title = "订单信息") {
 *     DsListItem(title = "订单编号", trailingText = "2024052012345")
 * }
 * ```
 */
public @Composable fun DsCard(
    modifier: Modifier = Modifier,
    title: String? = null,
    content: @Composable ColumnScope.() -> Unit,
) {
    Column(
        modifier = modifier
            .dsComponent("DsCard")
            .clip(RoundedCornerShape(DsRadius.Lg))
            .background(DsColors.BgCard)
            .padding(DsSpacing.Lg),
    ) {
        if (title != null) {
            DsTextInternal(text = title, style = DsTypography.Subtitle, color = DsColors.TextPrimary, maxLines = 1)
            Spacer(Modifier.height(DsSpacing.Md))
        }
        content()
    }
}

public @Composable fun DsCardPreview() {
    DsCard(title = "订单信息") {
        DsListItem(title = "订单编号", trailingText = "2024052012345")
        DsListItem(title = "下单时间", trailingText = "2024-05-20 12:30")
    }
}
