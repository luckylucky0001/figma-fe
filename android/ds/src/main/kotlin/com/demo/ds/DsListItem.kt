package com.demo.ds

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp

/**
 * Key/value list row (Figma: `Cell`, variants `Arrow=On|Off`, `Label=On|Off`).
 * Vertical padding 12, no horizontal padding (provided by the Card), min height 44.
 * Left: title (font.body, text-primary) + optional subtitle (font.caption, text-tertiary, 2dp gap).
 * Right: trailing value (font.body, text-secondary) + optional 16×16 chevron (4dp gap).
 *
 * @sample
 * ```
 * DsListItem(title = "收货地址", subtitle = "张三 138****0000", trailingText = "修改", showChevron = true)
 * ```
 */
public @Composable fun DsListItem(
    title: String,
    modifier: Modifier = Modifier,
    subtitle: String? = null,
    trailingText: String? = null,
    showChevron: Boolean = false,
) {
    Row(
        modifier = modifier
            .dsComponent("DsListItem")
            .fillMaxWidth()
            .heightIn(min = 44.dp)
            .padding(vertical = DsSpacing.Md),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        Column(Modifier.weight(1f)) {
            DsTextInternal(text = title, style = DsTypography.Body, color = DsColors.TextPrimary)
            if (subtitle != null) {
                Spacer(Modifier.height(2.dp))
                DsTextInternal(text = subtitle, style = DsTypography.Caption, color = DsColors.TextTertiary)
            }
        }
        if (trailingText != null) {
            Spacer(Modifier.width(DsSpacing.Sm))
            DsTextInternal(
                text = trailingText,
                style = DsTypography.Body.copy(textAlign = TextAlign.End),
                color = DsColors.TextSecondary,
            )
        }
        if (showChevron) {
            Spacer(Modifier.width(DsSpacing.Xs))
            DsChevron(left = false, strokeWidthDp = 1.5f, color = DsColors.TextTertiary, modifier = Modifier.size(16.dp))
        }
    }
}

public @Composable fun DsListItemPreview() {
    Column {
        DsListItem(title = "订单编号", trailingText = "2024052012345")
        DsListItem(title = "收货地址", subtitle = "张三 138****0000", trailingText = "修改", showChevron = true)
    }
}
