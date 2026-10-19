package com.demo.dsview

import android.content.Context
import android.graphics.drawable.GradientDrawable
import android.util.AttributeSet
import android.util.TypedValue
import android.view.Gravity
import androidx.annotation.ColorRes
import androidx.appcompat.widget.AppCompatTextView
import androidx.core.content.ContextCompat

/**
 * Legacy status tag (Figma: `Tag`, Light variant only). Height 20dp, horizontal padding 6dp,
 * radius 4dp, 12sp medium text. Background `{tone}_light`, text `{tone}`.
 *
 * ```xml
 * <com.demo.dsview.LegacyTagView
 *     android:layout_width="wrap_content"
 *     android:layout_height="20dp"
 *     android:text="待付款"
 *     app:tagTone="primary" />
 * ```
 */
class LegacyTagView @JvmOverloads constructor(
    context: Context,
    attrs: AttributeSet? = null,
    defStyleAttr: Int = android.R.attr.textViewStyle,
) : AppCompatTextView(context, attrs, defStyleAttr) {

    enum class Tone(@ColorRes val fg: Int, @ColorRes val bg: Int) {
        PRIMARY(R.color.ds_color_primary, R.color.ds_color_primary_light),
        SUCCESS(R.color.ds_color_success, R.color.ds_color_success_light),
        WARNING(R.color.ds_color_warning, R.color.ds_color_warning_light),
        DANGER(R.color.ds_color_danger, R.color.ds_color_danger_light),
    }

    var tone: Tone = Tone.PRIMARY
        set(value) {
            field = value
            applyTone()
        }

    init {
        val a = context.obtainStyledAttributes(attrs, R.styleable.LegacyTagView, defStyleAttr, 0)
        try {
            tone = Tone.values()[a.getInt(R.styleable.LegacyTagView_tagTone, 0)]
        } finally {
            a.recycle()
        }
        gravity = Gravity.CENTER
        includeFontPadding = false
        maxLines = 1
        setTextSize(TypedValue.COMPLEX_UNIT_SP, 12f)
        val padH = dp(6f)
        setPadding(padH, 0, padH, 0)
        minHeight = dp(20f)
        applyTone()
    }

    private fun applyTone() {
        setTextColor(ContextCompat.getColor(context, tone.fg))
        background = GradientDrawable().apply {
            cornerRadius = resources.getDimension(R.dimen.ds_radius_sm)
            setColor(ContextCompat.getColor(context, tone.bg))
        }
    }

    private fun dp(v: Float): Int =
        TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, v, resources.displayMetrics).toInt()
}
