package com.demo.dsview

import android.content.Context
import android.util.AttributeSet
import android.util.TypedValue
import android.view.Gravity
import android.widget.ImageView
import android.widget.LinearLayout
import android.widget.TextView
import androidx.core.content.ContextCompat

/**
 * Legacy key/value row (Figma: `Cell`). Vertical padding 12dp, min height 44dp, title on the left
 * (14sp text_primary), value on the right (14sp text_secondary), optional 16dp chevron 4dp after the value.
 * The Figma `Label` (subtitle) variant is not supported by this legacy view.
 *
 * ```xml
 * <com.demo.dsview.LegacyCellView
 *     android:layout_width="match_parent"
 *     android:layout_height="wrap_content"
 *     app:cellTitle="订单编号"
 *     app:cellValue="DD20261018093021"
 *     app:cellShowArrow="false" />
 * ```
 */
class LegacyCellView @JvmOverloads constructor(
    context: Context,
    attrs: AttributeSet? = null,
    defStyleAttr: Int = 0,
) : LinearLayout(context, attrs, defStyleAttr) {

    private val titleView = TextView(context)
    private val valueView = TextView(context)
    private val arrowView = ImageView(context)

    var title: CharSequence?
        get() = titleView.text
        set(value) { titleView.text = value }

    var value: CharSequence?
        get() = valueView.text
        set(v) { valueView.text = v }

    var showArrow: Boolean
        get() = arrowView.visibility == VISIBLE
        set(v) { arrowView.visibility = if (v) VISIBLE else GONE }

    init {
        orientation = HORIZONTAL
        gravity = Gravity.CENTER_VERTICAL
        minimumHeight = dp(44f)
        val padV = resources.getDimensionPixelSize(R.dimen.ds_space_md)
        setPadding(0, padV, 0, padV)

        titleView.setTextSize(TypedValue.COMPLEX_UNIT_SP, 14f)
        titleView.setTextColor(ContextCompat.getColor(context, R.color.ds_color_text_primary))
        addView(titleView, LayoutParams(0, LayoutParams.WRAP_CONTENT, 1f))

        valueView.setTextSize(TypedValue.COMPLEX_UNIT_SP, 14f)
        valueView.setTextColor(ContextCompat.getColor(context, R.color.ds_color_text_secondary))
        valueView.gravity = Gravity.END
        addView(valueView, LayoutParams(LayoutParams.WRAP_CONTENT, LayoutParams.WRAP_CONTENT))

        arrowView.setImageResource(android.R.drawable.arrow_down_float) // placeholder chevron
        arrowView.rotation = -90f
        addView(arrowView, LayoutParams(dp(16f), dp(16f)).apply { marginStart = dp(4f) })

        val a = context.obtainStyledAttributes(attrs, R.styleable.LegacyCellView, defStyleAttr, 0)
        try {
            title = a.getString(R.styleable.LegacyCellView_cellTitle)
            value = a.getString(R.styleable.LegacyCellView_cellValue)
            showArrow = a.getBoolean(R.styleable.LegacyCellView_cellShowArrow, false)
        } finally {
            a.recycle()
        }
    }

    private fun dp(v: Float): Int =
        TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, v, resources.displayMetrics).toInt()
}
