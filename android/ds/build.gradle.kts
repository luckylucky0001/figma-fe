plugins {
    kotlin("jvm")
    id("org.jetbrains.compose")
}
dependencies {
    // Compose Desktop = same androidx.compose.* APIs as Android, rendered on JVM via Skia.
    api(compose.desktop.currentOs)
    api(compose.foundation)
}
kotlin { jvmToolchain(21) }
