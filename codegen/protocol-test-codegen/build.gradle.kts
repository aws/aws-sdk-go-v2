/*
 *     Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 *     Licensed under the Apache License, Version 2.0 (the "License").
 *     You may not use this file except in compliance with the License.
 *     A copy of the License is located at
 *
 *      http://aws.amazon.com/apache2.0
 *
 *     or in the "license" file accompanying this file. This file is distributed
 *     on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 *     express or implied. See the License for the specific language governing
 *     permissions and limitations under the License.
 */

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
    implementation("software.amazon.smithy:smithy-aws-protocol-tests:1.42.0") // TODO(rpc2cbor) match smithy version when merged
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
