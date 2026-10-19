package com.demo.ds

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.text.BasicText
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.SpanStyle
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.buildAnnotatedString
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.withStyle
import androidx.compose.ui.unit.TextUnit
import androidx.compose.ui.unit.sp
import java.math.BigDecimal
import java.math.RoundingMode

/** Size axis of [DsPriceText]: font sizes of (¥, integer, decimals). */
enum class DsPriceSize(internal val symbol: TextUnit, internal val integer: TextUnit, internal val decimal: TextUnit) {
    Large(14.sp, 22.sp, 14.sp),
    Medium(12.sp, 18.sp, 12.sp),
    Small(12.sp, 14.sp, 12.sp),
}

/** Tone axis of [DsPriceText]: Primary = color.primary, Default = text-primary. */
enum class DsPriceTone(internal val color: Color) {
    Primary(DsColors.Primary),
    Default(DsColors.TextPrimary),
}

/**
 * Price text (Figma: `Price`): "¥" + integer + "." + two decimals in a single baseline-aligned
 * AnnotatedString. Integer weight 700, "¥" and decimals 600. Line height = 1.2 × integer size.
 *
 * @sample
 * ```
 * DsPriceText(amount = 269.0, size = DsPriceSize.Large, tone = DsPriceTone.Primary)
 * ```
 */
public @Composable fun DsPriceText(
    amount: Double,
    modifier: Modifier = Modifier,
    size: DsPriceSize = DsPriceSize.Medium,
    tone: DsPriceTone = DsPriceTone.Primary,
) {
    val formatted = BigDecimal.valueOf(amount).setScale(2, RoundingMode.HALF_UP).toPlainString()
    val intPart = formatted.substringBefore('.')
    val decPart = formatted.substringAfter('.')
    val text = buildAnnotatedString {
        withStyle(SpanStyle(fontSize = size.symbol, fontWeight = FontWeight.SemiBold)) { append("¥") }
        withStyle(SpanStyle(fontSize = size.integer, fontWeight = FontWeight.Bold)) { append(intPart) }
        withStyle(SpanStyle(fontSize = size.decimal, fontWeight = FontWeight.SemiBold)) { append(".$decPart") }
    }
    BasicText(
        text = text,
        modifier = modifier.dsComponent("DsPriceText"),
        style = TextStyle(
            fontFamily = DsFont,
            color = tone.color,
            fontSize = size.integer,
            lineHeight = (size.integer.value * 1.2f).sp,
        ),
        softWrap = false,
        maxLines = 1,
    )
}

public @Composable fun DsPriceTextPreview() {
    Column {
        DsPriceTextSamples()
    }
}

@Composable
private fun DsPriceTextSamples() {
    DsPriceSize.values().forEach { s ->
        DsPriceTone.values().forEach { t -> DsPriceText(amount = 1299.5, size = s, tone = t) }
    }
}
