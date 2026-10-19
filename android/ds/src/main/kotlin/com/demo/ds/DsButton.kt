package com.demo.ds

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.semantics.Role
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.Dp
import androidx.compose.ui.unit.TextUnit
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

/** Style axis of [DsButton] (Figma `Button` → `Type`). */
enum class DsButtonStyle { Primary, Default, Text }

/** Size axis of [DsButton] (Figma `Button` → `Size`). */
enum class DsButtonSize(
    internal val height: Dp,
    internal val paddingH: Dp,
    internal val fontSize: TextUnit,
    internal val weight: FontWeight,
    internal val radius: Dp,
) {
    Large(44.dp, 24.dp, 16.sp, FontWeight.SemiBold, 22.dp),
    Medium(32.dp, 16.dp, 14.sp, FontWeight.Medium, 16.dp),
    Small(24.dp, 12.dp, 12.sp, FontWeight.Medium, 12.dp),
}

/**
 * Button (Figma: `Button`, variants `Type=Primary|Default|Text`, `Size=Large|Medium|Small`).
 * - Primary: bg primary, white text
 * - Default: bg-card, 1dp border, text-primary
 * - Text: no background / border, primary text
 * Large 44 / pad 24 / 16·600 / r22; Medium 32 / pad 16 / 14·500 / r16; Small 24 / pad 12 / 12·500 / r12.
 *
 * @sample
 * ```
 * DsButton(text = "提醒发货", onClick = { viewModel.remind() }, style = DsButtonStyle.Primary, size = DsButtonSize.Medium)
 * ```
 */
public @Composable fun DsButton(
    text: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    style: DsButtonStyle = DsButtonStyle.Primary,
    size: DsButtonSize = DsButtonSize.Medium,
) {
    val shape = RoundedCornerShape(size.radius)
    val textColor: Color = when (style) {
        DsButtonStyle.Primary -> DsColors.White
        DsButtonStyle.Default -> DsColors.TextPrimary
        DsButtonStyle.Text -> DsColors.Primary
    }
    val styled = when (style) {
        DsButtonStyle.Primary -> Modifier.background(DsColors.Primary, shape)
        DsButtonStyle.Default -> Modifier.background(DsColors.BgCard, shape).border(1.dp, DsColors.Border, shape)
        DsButtonStyle.Text -> Modifier
    }
    Box(
        modifier = modifier
            .dsComponent("DsButton")
            .height(size.height)
            .clip(shape)
            .then(styled)
            .clickable(role = Role.Button, onClick = onClick)
            .padding(horizontal = size.paddingH),
        contentAlignment = Alignment.Center,
    ) {
        DsTextInternal(
            text = text,
            style = TextStyle(fontFamily = DsFont, fontSize = size.fontSize, fontWeight = size.weight, lineHeight = size.fontSize * 1.4f),
            color = textColor,
            maxLines = 1,
        )
    }
}

public @Composable fun DsButtonPreview() {
    Row {
        DsButton(text = "取消订单", onClick = {}, style = DsButtonStyle.Default, modifier = Modifier.padding(end = 8.dp))
        DsButton(text = "立即支付", onClick = {})
    }
}
