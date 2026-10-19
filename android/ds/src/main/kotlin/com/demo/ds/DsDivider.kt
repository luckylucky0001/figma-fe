package com.demo.ds

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

/**
 * 1dp horizontal divider in color.border, no margins (Figma: `Divider`).
 *
 * @sample
 * ```
 * DsCard { DsListItem(title = "A"); DsDivider(); DsListItem(title = "B") }
 * ```
 */
public @Composable fun DsDivider(modifier: Modifier = Modifier) {
    Box(
        modifier = modifier
            .dsComponent("DsDivider")
            .fillMaxWidth()
            .height(1.dp)
            .background(DsColors.Border),
    )
}

public @Composable fun DsDividerPreview() {
    DsDivider()
}
