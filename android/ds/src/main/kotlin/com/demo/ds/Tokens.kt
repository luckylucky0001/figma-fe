// Compose 1.2 compatible (compileSdk 32 / minSdk 23)
// Only APIs present in Jetpack Compose 1.2.x are used in :ds and :app. The one desktop-only bit is
// font loading below (androidx.compose.ui.text.platform.Font(File)); on Android swap it for
// Font(R.font.noto_sans_sc_regular, FontWeight.Normal) etc.
package com.demo.ds

import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.platform.Font
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import java.io.File

/** Color tokens (design/SPEC.md §1). */
object DsColors {
    val Primary = Color(0xFFFF4D1A)
    val PrimaryLight = Color(0xFFFFF1EB)
    val Success = Color(0xFF00B578)
    val SuccessLight = Color(0xFFE6F8F1)
    // NOTE: intentional drift from Figma (color.warning = #FF8F1F) — kept for the drift-report demo.
    val Warning = Color(0xFFFF8800)
    val WarningLight = Color(0xFFFFF3E8)
    val Danger = Color(0xFFF53F3F)
    val DangerLight = Color(0xFFFFECE8)
    val TextPrimary = Color(0xFF1F2329)
    val TextSecondary = Color(0xFF646A73)
    val TextTertiary = Color(0xFF8F959E)
    val BgPage = Color(0xFFF5F6F7)
    val BgCard = Color(0xFFFFFFFF)
    val BgPlaceholder = Color(0xFFF2F3F5)
    val Border = Color(0xFFE5E6EB)
    val White = Color(0xFFFFFFFF)
}

/** Spacing tokens: xs/sm/md/lg/xl = 4/8/12/16/24. */
object DsSpacing {
    val Xs = 4.dp
    val Sm = 8.dp
    val Md = 12.dp
    val Lg = 16.dp
    val Xl = 24.dp
}

/** Radius tokens: sm/md/lg/full = 4/8/12/999. */
object DsRadius {
    val Sm = 4.dp
    val Md = 8.dp
    val Lg = 12.dp
    val Full = 999.dp
}

private fun loadFonts(): FontFamily {
    val dir = File(System.getProperty("ds.fontDir") ?: (System.getProperty("user.home") + "/.local/share/fonts"))
    val files = listOf(
        "NotoSansSC_400Regular.ttf" to FontWeight.Normal,
        "NotoSansSC_500Medium.ttf" to FontWeight.Medium,
        "NotoSansSC_600SemiBold.ttf" to FontWeight.SemiBold,
        "NotoSansSC_700Bold.ttf" to FontWeight.Bold,
    ).map { (name, w) -> File(dir, name) to w }.filter { it.first.isFile }
    if (files.isEmpty()) {
        System.err.println("[ds] Noto Sans SC not found in $dir, falling back to default font family")
        return FontFamily.Default
    }
    return FontFamily(files.map { (f, w) -> Font(f, w) })
}

/** Noto Sans SC loaded from `~/.local/share/fonts` (override with `-Dds.fontDir=`). */
val DsFont: FontFamily by lazy { loadFonts() }

/** Typography tokens: size / lineHeight / weight. */
object DsTypography {
    val Title = TextStyle(fontFamily = DsFont, fontSize = 17.sp, lineHeight = 24.sp, fontWeight = FontWeight.SemiBold)
    val Subtitle = TextStyle(fontFamily = DsFont, fontSize = 15.sp, lineHeight = 22.sp, fontWeight = FontWeight.SemiBold)
    val Body = TextStyle(fontFamily = DsFont, fontSize = 14.sp, lineHeight = 20.sp, fontWeight = FontWeight.Normal)
    val Caption = TextStyle(fontFamily = DsFont, fontSize = 12.sp, lineHeight = 18.sp, fontWeight = FontWeight.Normal)
}
