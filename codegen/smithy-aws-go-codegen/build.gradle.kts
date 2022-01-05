/*
 * Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *  http://aws.amazon.com/apache2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

val smithyVersion: String by project

plugins {
    `java-library`
    jacoco
}

group = "software.amazon.smithy.go"
version = "0.1.0"

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(17))
    }
}

tasks.withType<Test> {
    useJUnitPlatform()
}

dependencies {
    api("software.amazon.smithy:smithy-aws-traits:$smithyVersion")
    api("software.amazon.smithy:smithy-aws-iam-traits:$smithyVersion")
    api("software.amazon.smithy:smithy-aws-cloudformation-traits:$smithyVersion")
    api("software.amazon.smithy.go:smithy-go-codegen:0.1.0")
    testImplementation("org.junit.jupiter:junit-jupiter-api:5.4.0")
    testImplementation("org.junit.jupiter:junit-jupiter-engine:5.4.0")
    testCompileOnly("org.junit.jupiter:junit-jupiter-params:5.4.0")
    testImplementation("org.hamcrest:hamcrest:2.1")
}

// == Code coverage ==
// Always run the jacoco test report after testing.
tasks["test"].finalizedBy(tasks["jacocoTestReport"])

// Configure jacoco to generate an HTML report.
tasks.withType<JacocoReport> {
    reports {
        xml.isEnabled = false
        csv.isEnabled = false
        html.destination = file("$buildDir/reports/jacoco")
    }
}
