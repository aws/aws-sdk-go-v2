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

import java.net.URLClassLoader
import software.amazon.smithy.model.Model
import software.amazon.smithy.model.node.Node
import software.amazon.smithy.model.shapes.ServiceShape
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
    implementation("software.amazon.smithy:smithy-protocol-tests:$smithyVersion")
    implementation("software.amazon.smithy:smithy-aws-protocol-tests:$smithyVersion")
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

// Services to exclude from codegen (real-service-specific tests, validation, etc.)
val excludedServices = setOf(
    "com.amazonaws.glacier#Glacier",
    "com.amazonaws.apigateway#BackplaneControlService",
    "com.amazonaws.machinelearning#AmazonML_20141212",
    "com.amazonaws.s3#AmazonS3",
    "aws.protocoltests.restjson.validation#RestJsonValidation",
)

// Override (projectionName, moduleSuffix) for services where auto-derivation
// doesn't match the existing convention.
// Pair(projection-name, folder-name)
val overrides = mapOf(
    "aws.protocoltests.ec2#AwsEc2" to Pair("aws-ec2", "ec2query"),
    "aws.protocoltests.json#JsonProtocol" to Pair("aws-json", "jsonrpc"),
    "aws.protocoltests.json10#JsonRpc10" to Pair("aws-json-10", "jsonrpc10"),
    "aws.protocoltests.query#AwsQuery" to Pair("aws-query", "query"),
    "aws.protocoltests.restjson#RestJson" to Pair("aws-restjson", "awsrestjson"),
    "aws.protocoltests.restxml#RestXml" to Pair("aws-restxml", "restxml"),
    "aws.protocoltests.restxml.xmlns#RestXmlWithNamespace" to Pair("aws-restxml-with-namespace", "restxmlwithnamespace"),
    "smithy.protocoltests.rpcv2Cbor#RpcV2Protocol" to Pair("smithy-rpcv2-cbor", "smithyrpcv2cbor"),
    "aws.protocoltests.json#Json10QueryCompatible" to Pair("aws-json-10-querycompatible", "jsonrpc10querycompatible"),
    "aws.protocoltests.json10#QueryCompatibleJsonRpc10" to Pair("querycompatiblejsonrpc10", "querycompatiblejsonrpc10"),
    "aws.protocoltests.rpcv2cbor#NonQueryCompatibleRpcV2Protocol" to Pair("nonquerycompatiblerpcv2protocol", "nonquerycompatiblerpcv2protocol"),
    "aws.protocoltests.rpcv2cbor#QueryCompatibleRpcV2Protocol" to Pair("querycompatiblerpcv2protocol", "querycompatiblerpcv2protocol"),
    "smithy.protocoltests.rpcv2Cbor#RpcV2CborQueryCompatible" to Pair("smithy-rpcv2-cbor-querycompatible", "smithyrpcv2cborquerycompatible"),
)

fun deriveNames(shapeId: String): Pair<String, String> {
    val svcName = shapeId.substringAfter("#")
    val hyphenated = svcName.replace(Regex("([a-z0-9])([A-Z])"), "$1-$2").toLowerCase()
    return Pair(hyphenated, hyphenated.replace("-", ""))
}

// Generates smithy-build.json by discovering all service shapes from the model
// dependencies and local models/ directory.
tasks.register("generate-smithy-build") {
    doLast {
        val urls = project.configurations.getByName("runtimeClasspath")
            .map { it.toURI().toURL() }.toTypedArray()
        val cl = URLClassLoader(urls, Model::class.java.classLoader)
        val model = Model.assembler(cl)
            .discoverModels(cl)
            .addImport("${project.projectDir}/models/")
            .assemble()
            .result
            .get()

        val modulePrefix = "github.com/aws/aws-sdk-go-v2/internal/protocoltest"
        val projectionsBuilder = Node.objectNodeBuilder()

        model.shapes(ServiceShape::class.javaObjectType).sorted().forEach { service ->
            val shapeId = service.id.toString()
            if (shapeId in excludedServices) return@forEach

            val (projName, modSuffix) = overrides[shapeId] ?: deriveNames(shapeId)

            projectionsBuilder.withMember(projName, Node.objectNodeBuilder()
                .withMember("transforms", Node.fromNodes(
                    Node.objectNodeBuilder()
                        .withMember("name", "includeServices")
                        .withMember("args", Node.objectNode()
                            .withMember("services", Node.fromStrings(shapeId)))
                        .build(),
                    Node.objectNodeBuilder()
                        .withMember("name", "removeUnusedShapes")
                        .build(),
                ))
                .withMember("plugins", Node.objectNode()
                    .withMember("go-codegen", Node.objectNodeBuilder()
                        .withMember("service", shapeId)
                        .withMember("module", "$modulePrefix/$modSuffix")
                        .build()))
                .build())
        }

        file("smithy-build.json").writeText(Node.prettyPrintJson(Node.objectNodeBuilder()
            .withMember("version", "1.0")
            .withMember("sources", Node.fromStrings("models"))
            .withMember("projections", projectionsBuilder.build())
            .build()))
    }
}

// Run the `buildSdk` automatically.
tasks["build"]
    .dependsOn("generate-smithy-build")
    .finalizedBy(tasks["buildSdk"])

val protocolTestDir = file("$rootDir/../internal/protocoltest")
val manualDir = file("${project.projectDir}/manual")

tasks.create<Delete>("cleanProtocolTests") {
    dependsOn("buildSdk")
    delete(protocolTestDir)
}

// ensure built artifacts are put into the SDK's folders
tasks.create<Exec>("copyGoCodegen") {
    dependsOn("cleanProtocolTests")
    commandLine("$rootDir/copy_go_codegen.sh", "$rootDir/..", (tasks["buildSdk"] as SmithyBuild).outputDirectory.absolutePath)
}

tasks.create<Copy>("copyManualFiles") {
    dependsOn("copyGoCodegen")
    from(manualDir)
    into(protocolTestDir)
}

tasks["buildSdk"].finalizedBy(tasks["copyGoCodegen"])
tasks["copyGoCodegen"].finalizedBy(tasks["copyManualFiles"])

java.sourceSets["main"].java {
    srcDirs("models")
}
