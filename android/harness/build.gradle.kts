plugins {
    kotlin("jvm")
    id("org.jetbrains.compose")
    application
}
dependencies {
    implementation(project(":app"))
    implementation(project(":ds"))
    implementation(compose.desktop.currentOs)
}
kotlin { jvmToolchain(21) }
application {
    mainClass.set("com.demo.harness.MainKt")
    applicationName = "harness"
    applicationDefaultJvmArgs = listOf("-Djava.awt.headless=true")
}
tasks.named<JavaExec>("run") {
    // Relative outDir args resolve against the directory gradle was invoked from.
    workingDir = rootProject.projectDir.parentFile
}
