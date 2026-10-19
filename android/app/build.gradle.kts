plugins {
    kotlin("jvm")
    id("org.jetbrains.compose")
}
dependencies {
    implementation(project(":ds"))
    implementation(compose.foundation)
}
kotlin { jvmToolchain(21) }
