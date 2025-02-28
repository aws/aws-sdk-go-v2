import software.amazon.smithy.gradle.tasks.SmithyBuild

val smithyVersion: String by project

buildscript {
    val smithyVersion: String by project
    repositories {
        mavenLocal()
        mavenCentral()
    }
    dependencies {
        "classpath"("software.amazon.smithy:smithy-cli:$smithyVersion")
    }
}

plugins {
    val smithyGradleVersion: String by project
    id("software.amazon.smithy") version smithyGradleVersion
}

dependencies {
    implementation("software.amazon.smithy:smithy-cli:$smithyVersion")
    implementation(project(":smithy-aws-go-codegen"))
}

// This project doesn't produce a JAR.
tasks["jar"].enabled = false

// Run the SmithyBuild task manually since this project needs the built JAR
// from smithy-aws-typescript-codegen.
tasks["smithyBuildJar"].enabled = false

tasks.create<SmithyBuild>("buildSdk") {
    addRuntimeClasspath = true
}

// Run the `buildSdk` automatically.
tasks["build"].finalizedBy(tasks["buildSdk"])

// ensure built artifacts are put into the SDK's folders
tasks.create<Exec>("copyGoCodegen") {
    dependsOn("buildSdk")
    commandLine("$rootDir/copy_go_codegen.sh", "$rootDir/..", (tasks["buildSdk"] as SmithyBuild).outputDirectory.absolutePath)
}
tasks["buildSdk"].finalizedBy(tasks["copyGoCodegen"])

java.sourceSets["main"].java {
    srcDirs("model", "src/main/smithy")
}
