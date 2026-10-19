package com.demo.ds

import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.testTag
import androidx.compose.ui.semantics.SemanticsPropertyKey
import androidx.compose.ui.semantics.SemanticsPropertyReceiver
import androidx.compose.ui.semantics.semantics

/** Marks the root semantics node of a design-system component, e.g. `"DsButton"`. */
val DsComponentKey = SemanticsPropertyKey<String>("DsComponent")

var SemanticsPropertyReceiver.dsComponentName by DsComponentKey

/** Tag the root of a ds component. Every `DsXxx` applies this to its root node. */
fun Modifier.dsComponent(name: String): Modifier = this.semantics { this[DsComponentKey] = name }

/** Attach a Figma node id to a node (stored as `testTag`), used by pages for design ↔ code mapping. */
fun Modifier.figmaId(id: String): Modifier = this.testTag(id)
