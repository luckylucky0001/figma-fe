package com.demo.dsview

import android.content.Context
import android.graphics.Color
import android.graphics.Typeface
import android.graphics.drawable.GradientDrawable
import android.util.AttributeSet
import android.util.TypedValue
import android.view.Gravity
import androidx.appcompat.widget.AppCompatButton
import androidx.core.content.ContextCompat

/**
 * Legacy button (Figma: `Button`, `Type=Primary|Default|Text`, `Size=Large|Medium|Small`).
 * - primary: bg ds_color_primary, white text
 * - default: bg ds_color_bg_card, 1dp ds_color_border stroke, ds_color_text_primary text
 * - text: no background, ds_color_primary text
 * - large 44dp / pad 24 / 16sp bold / r22; medium 32dp / pad 16 / 14sp / r16; small 24dp / pad 12 / 12sp / r12
 *
 * ```xml
 * <com.demo.dsview.LegacyButton
 *     android:layout_width="wrap_content"
 *     android:layout_height="wrap_content"
 *     android:text="去支付"
 *     app:btnType="primary"
 *     app:btnSize="medium" />
 * ```
 */
class LegacyButton @JvmOverloads constructor(
    context: Context,
    attrs: AttributeSet? = null,
    defStyleAttr: Int = androidx.appcompat.R.attr.buttonStyle,
) : AppCompatButton(context, attrs, defStyleAttr) {

    enum class Type { PRIMARY, DEFAULT, TEXT }

    enum class Size(val heightDp: Float, val padHDp: Float, val textSp: Float, val radiusDp: Float, val bold: Boolean) {
        LARGE(44f, 24f, 16f, 22f, true),
        MEDIUM(32f, 16f, 14f, 16f, false),
        SMALL(24f, 12f, 12f, 12f, false),
    }

    var type: Type = Type.PRIMARY
        set(value) { field = value; applyStyle() }

    var size: Size = Size.MEDIUM
        set(value) { field = value; applyStyle() }

    init {
        val a = context.obtainStyledAttributes(attrs, R.styleable.LegacyButton, defStyleAttr, 0)
        try {
            type = Type.values()[a.getInt(R.styleable.LegacyButton_btnType, 0)]
            size = Size.values()[a.getInt(R.styleable.LegacyButton_btnSize, 1)]
        } finally {
            a.recycle()
        }
        isAllCaps = false
        gravity = Gravity.CENTER
        includeFontPadding = false
        stateListAnimator = null
        applyStyle()
    }

    private fun applyStyle() {
        minHeight = dp(size.heightDp)
        minimumHeight = dp(size.heightDp)
        minWidth = 0
        minimumWidth = 0
        val padH = dp(size.padHDp)
        setPadding(padH, 0, padH, 0)
        setTextSize(TypedValue.COMPLEX_UNIT_SP, size.textSp)
        typeface = if (size.bold) Typeface.DEFAULT_BOLD else Typeface.DEFAULT
        val bg = GradientDrawable().apply { cornerRadius = dp(size.radiusDp).toFloat() }
        when (type) {
            Type.PRIMARY -> {
                bg.setColor(ContextCompat.getColor(context, R.color.ds_color_primary))
                setTextColor(Color.WHITE)
            }
            Type.DEFAULT -> {
                bg.setColor(ContextCompat.getColor(context, R.color.ds_color_bg_card))
                bg.setStroke(dp(1f), ContextCompat.getColor(context, R.color.ds_color_border))
                setTextColor(ContextCompat.getColor(context, R.color.ds_color_text_primary))
            }
            Type.TEXT -> {
                bg.setColor(Color.TRANSPARENT)
                setTextColor(ContextCompat.getColor(context, R.color.ds_color_primary))
            }
        }
        background = bg
    }

    private fun dp(v: Float): Int =
        TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, v, resources.displayMetrics).toInt()
}
