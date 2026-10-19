package com.demo.harness

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.ui.ExperimentalComposeUiApi
import androidx.compose.ui.ImageComposeScene
import androidx.compose.ui.Modifier
import androidx.compose.ui.semantics.SemanticsNode
import androidx.compose.ui.semantics.SemanticsProperties
import androidx.compose.ui.semantics.getOrNull
import androidx.compose.ui.unit.Density
import com.demo.app.screens
import com.demo.ds.DsComponentKey
import org.jetbrains.skia.EncodedImageFormat
import java.io.File
import java.util.Locale
import kotlin.system.exitProcess

private const val SCALE = 2f

private data class Node(
    val figmaId: String?,
    val ds: String?,
    val text: String?,
    val x: Float,
    val y: Float,
    val w: Float,
    val h: Float,
)

private fun usage(): Nothing {
    System.err.println("usage: harness <ScreenName> <outDir> [widthDp=375] [heightDp=812]")
    System.err.println("screens: ${screens.keys.joinToString(", ")}")
    exitProcess(2)
}

@OptIn(ExperimentalComposeUiApi::class)
fun main(args: Array<String>) {
    if (args.size < 2) usage()
    val name = args[0]
    val outDir = File(args[1])
    val widthDp = args.getOrNull(2)?.toIntOrNull() ?: 375
    val heightDp = args.getOrNull(3)?.toIntOrNull() ?: 812
    val screen = screens[name]
    if (screen == null) {
        System.err.println("error: screen '$name' is not registered in com.demo.app.screens")
        System.err.println("screens: ${screens.keys.joinToString(", ")}")
        exitProcess(1)
    }
    outDir.mkdirs()

    val wPx = (widthDp * SCALE).toInt()
    val hPx = (heightDp * SCALE).toInt()
    val scene = ImageComposeScene(wPx, hPx, Density(SCALE)) {
        Box(Modifier.fillMaxSize()) { screen() }
    }
    try {
        // Render twice so lazily-loaded fonts / first-frame layout settle.
        scene.render()
        val image = scene.render()
        val png = image.encodeToData(EncodedImageFormat.PNG) ?: error("PNG encoding failed")
        File(outDir, "screen.png").writeBytes(png.bytes)

        val nodes = mutableListOf<Node>()
        scene.roots.forEach { owner -> collect(owner.semanticsOwner.unmergedRootSemanticsNode, nodes) }
        File(outDir, "nodes.json").writeText(toJson(nodes, widthDp, heightDp))
        println("wrote ${File(outDir, "screen.png").path} (${wPx}x$hPx) and nodes.json (${nodes.size} nodes)")
    } finally {
        scene.close()
    }
    exitProcess(0)
}

private fun SemanticsNode.ownText(): String? =
    config.getOrNull(SemanticsProperties.Text)?.joinToString("") { it.text }?.takeIf { it.isNotEmpty() }

private fun SemanticsNode.isDs(): Boolean = config.getOrNull(DsComponentKey) != null

/** Texts of descendants, in tree order. Returns null if a nested ds component is found (container components). */
private fun SemanticsNode.leafTexts(): List<String>? {
    val out = mutableListOf<String>()
    fun walk(n: SemanticsNode): Boolean {
        for (c in n.children) {
            if (c.isDs()) return false
            c.ownText()?.let { out += it }
            if (!walk(c)) return false
        }
        return true
    }
    return if (walk(this)) out else null
}

private fun collect(node: SemanticsNode, out: MutableList<Node>) {
    val tag = node.config.getOrNull(SemanticsProperties.TestTag)
    val ds = node.config.getOrNull(DsComponentKey)
    var text = node.ownText()
    if (ds != null && text == null) {
        // Leaf components (DsButton, DsBadge, DsListItem...) report the text they render, like the H5 DOM side.
        text = node.leafTexts()?.takeIf { it.isNotEmpty() }?.joinToString(" ")
    }
    if (tag != null || ds != null || text != null) {
        // Unclipped geometry (like DOM getBoundingClientRect): nodes below the viewport keep their real size.
        val pos = node.positionInRoot
        out += Node(tag, ds, text, pos.x / SCALE, pos.y / SCALE, node.size.width / SCALE, node.size.height / SCALE)
    }
    node.children.forEach { collect(it, out) }
}

private fun r1(v: Float): String {
    val d = Math.round(v * 10.0) / 10.0
    return if (d == Math.floor(d)) d.toLong().toString() else String.format(Locale.US, "%.1f", d)
}

private fun str(s: String?): String {
    if (s == null) return "null"
    val sb = StringBuilder("\"")
    for (ch in s) {
        when (ch) {
            '"' -> sb.append("\\\"")
            '\\' -> sb.append("\\\\")
            '\n' -> sb.append("\\n")
            '\r' -> sb.append("\\r")
            '\t' -> sb.append("\\t")
            else -> if (ch < ' ') sb.append(String.format("\\u%04x", ch.code)) else sb.append(ch)
        }
    }
    return sb.append('"').toString()
}

private fun toJson(nodes: List<Node>, w: Int, h: Int): String = buildString {
    append("{ \"platform\": \"android\", \"scale\": 2, \"width\": $w, \"height\": $h,\n  \"nodes\": [\n")
    nodes.forEachIndexed { i, n ->
        append("    { \"figmaId\": ${str(n.figmaId)}, \"ds\": ${str(n.ds)}, \"text\": ${str(n.text)}, ")
        append("\"x\": ${r1(n.x)}, \"y\": ${r1(n.y)}, \"w\": ${r1(n.w)}, \"h\": ${r1(n.h)} }")
        append(if (i < nodes.size - 1) ",\n" else "\n")
    }
    append("  ]\n}\n")
}
