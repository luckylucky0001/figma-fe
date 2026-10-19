package com.demo.app

import androidx.compose.runtime.Composable
import com.demo.app.screens.GalleryScreen
import com.demo.app.screens.OrderListScreen

/** Screen registry used by the harness: `harness <ScreenName> <outDir> [widthDp] [heightDp]`. */
val screens: Map<String, @Composable () -> Unit> = linkedMapOf(
    "OrderList" to { OrderListScreen() },
    "Gallery" to { GalleryScreen() },
)
