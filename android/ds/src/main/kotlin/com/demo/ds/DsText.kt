package com.demo.ds

import androidx.compose.foundation.text.BasicText
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.layout.layout
import androidx.compose.ui.text.style.TextOverflow
import kotlin.math.abs
import kotlin.math.roundToInt

/** Internal text primitive used by ds components (single source for font + color handling). */
@Composable
internal fun DsTextInternal(
    text: String,
    style: TextStyle,
    color: Color,
    modifier: Modifier = Modifier,
    maxLines: Int = Int.MAX_VALUE,
) {
    BasicText(
        text = text,
        modifier = modifier.snapToLineHeight(style),
        style = style.copy(color = color),
        overflow = TextOverflow.Ellipsis,
        softWrap = true,
        maxLines = maxLines,
    )
}

/**
 * Public text helper for screens (pages may only use `com.demo.ds.*`). Defaults to font.body / text-primary.
 * Intentionally NOT tagged with [dsComponent]: it is a typography primitive, not a Figma component.
 *
 * @sample
 * ```
 * DsText(text = "数码好物旗舰店", style = DsTypography.Subtitle)
 * ```
 */
public @Composable fun DsText(
    text: String,
    modifier: Modifier = Modifier,
    style: TextStyle = DsTypography.Body,
    color: Color = DsColors.TextPrimary,
    maxLines: Int = Int.MAX_VALUE,
) {
    DsTextInternal(text = text, style = style, color = color, modifier = modifier, maxLines = maxLines)
}

/**
 * Keeps a text node's height an exact multiple of the token line height (e.g. 20dp for font.body).
 * Skia's paragraph layout can come out 1px taller because of float rounding (14sp × 20/14); this
 * snaps such ±1px noise away so containers match the spec (e.g. DsListItem = 44). Plain
 * `Modifier.layout`, available since Compose 1.0.
 */
private fun Modifier.snapToLineHeight(style: TextStyle): Modifier =
    if (!style.lineHeight.isSp) this else layout { measurable, constraints ->
        val p = measurable.measure(constraints)
        val line = style.lineHeight.toPx()
        val lines = (p.height / line).roundToInt().coerceAtLeast(1)
        val snapped = (lines * line).roundToInt()
        val h = if (abs(p.height - snapped) <= 1) snapped.coerceIn(constraints.minHeight, constraints.maxHeight) else p.height
        layout(p.width, h) { p.placeRelative(0, 0) }
    }
