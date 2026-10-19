pluginManagement {
    // Only Maven Central + Gradle Plugin Portal are reachable here (dl.google.com / maven.google.com are blocked).
    repositories { mavenCentral(); gradlePluginPortal() }
}
dependencyResolutionManagement {
    repositories { mavenCentral() }
}
rootProject.name = "demo-android"
include(":ds", ":app", ":harness")
