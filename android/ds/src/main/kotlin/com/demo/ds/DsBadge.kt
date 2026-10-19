package com.demo.ds

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

/** Tone axis of [DsBadge] (Figma `Tag` → `Tone`). */
enum class DsBadgeTone(internal val fg: Color, internal val bg: Color) {
    Primary(DsColors.Primary, DsColors.PrimaryLight),
    Success(DsColors.Success, DsColors.SuccessLight),
    Warning(DsColors.Warning, DsColors.WarningLight),
    Danger(DsColors.Danger, DsColors.DangerLight),
}

/**
 * Status badge (Figma: `Tag`). Height 20, horizontal padding 6, radius.sm (4), font 12 / 20 / 500.
 * Only the Light look is implemented: background `{tone}-light`, text `{tone}`.
 * (Figma's `Variant=Solid` has no Android counterpart yet.)
 *
 * @sample
 * ```
 * DsBadge(text = "待发货", tone = DsBadgeTone.Warning)
 * ```
 */
public @Composable fun DsBadge(
    text: String,
    modifier: Modifier = Modifier,
    tone: DsBadgeTone = DsBadgeTone.Primary,
) {
    Box(
        modifier = modifier
            .dsComponent("DsBadge")
            .height(20.dp)
            .background(tone.bg, RoundedCornerShape(DsRadius.Sm))
            .padding(horizontal = 6.dp),
        contentAlignment = Alignment.Center,
    ) {
        DsTextInternal(text = text, style = BadgeTextStyle, color = tone.fg, maxLines = 1)
    }
}

private val BadgeTextStyle: TextStyle
    get() = TextStyle(fontFamily = DsFont, fontSize = 12.sp, lineHeight = 20.sp, fontWeight = FontWeight.Medium)

public @Composable fun DsBadgePreview() {
    Row {
        DsBadgeTone.values().forEach { DsBadge(text = it.name, tone = it, modifier = Modifier.padding(end = 8.dp)) }
    }
}
