package com.demo.ds

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp

/** Tone axis of [DsBanner] (Figma `NoticeBar` → `Tone`). */
enum class DsBannerTone(internal val fg: Color, internal val bg: Color) {
    Info(DsColors.Primary, DsColors.PrimaryLight),
    Warning(DsColors.Warning, DsColors.WarningLight),
}

/**
 * Notice banner (Figma: `NoticeBar`). Height 36, horizontal padding 12, radius.md (8),
 * font.caption, single line, vertically centered. Info = primary on primary-light,
 * Warning = warning on warning-light.
 *
 * @sample
 * ```
 * DsBanner(text = "商家已发货，预计 3 天内送达", tone = DsBannerTone.Info)
 * ```
 */
public @Composable fun DsBanner(
    text: String,
    modifier: Modifier = Modifier,
    tone: DsBannerTone = DsBannerTone.Info,
) {
    Box(
        modifier = modifier
            .dsComponent("DsBanner")
            .fillMaxWidth()
            .height(36.dp)
            .background(tone.bg, RoundedCornerShape(DsRadius.Md))
            .padding(horizontal = DsSpacing.Md),
        contentAlignment = Alignment.CenterStart,
    ) {
        DsTextInternal(text = text, style = DsTypography.Caption, color = tone.fg, maxLines = 1)
    }
}

public @Composable fun DsBannerPreview() {
    DsBanner(text = "商家已发货，预计 3 天内送达")
}
