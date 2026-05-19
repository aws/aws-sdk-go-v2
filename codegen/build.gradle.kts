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

allprojects {
    repositories {
        mavenLocal()
        mavenCentral()
    }
}

// Exclude protocol-test-codegen from the main build loop. Protocol tests are
// generated separately via :protocol-test-codegen:build (see Makefile target
// smithy-generate-protocol-tests) to avoid regenerating them on every SDK
// codegen run. When protocol-test-codegen is not explicitly targeted, all its
// tasks are disabled so finalizers like buildSdk don't execute.
project(":protocol-test-codegen") {
    if (!gradle.startParameter.taskNames.any { it.contains("protocol-test-codegen") }) {
        tasks.configureEach {
            enabled = false
        }
    }
}
