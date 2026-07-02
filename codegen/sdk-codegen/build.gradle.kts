/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

import software.amazon.smithy.model.Model
import software.amazon.smithy.model.node.Node
import software.amazon.smithy.model.shapes.ServiceShape
import software.amazon.smithy.gradle.tasks.SmithyBuild
import software.amazon.smithy.aws.traits.ServiceTrait

val smithyVersion: String by project

buildscript {
    val smithyVersion: String by project
    repositories {
        mavenLocal()
        mavenCentral()
    }
    dependencies {
        "classpath"("software.amazon.smithy:smithy-cli:$smithyVersion")
        "classpath"("software.amazon.smithy:smithy-aws-traits:$smithyVersion")
    }
}

plugins {
    val smithyGradleVersion: String by project
    id("software.amazon.smithy") version smithyGradleVersion
}

dependencies {
    val smithyVersion: String by project
    implementation(project(":smithy-aws-go-codegen"))
    implementation("software.amazon.smithy:smithy-smoke-test-traits:$smithyVersion")
    implementation("software.amazon.smithy:smithy-aws-smoke-test-model:$smithyVersion")
}

// This project doesn't produce a JAR.
tasks["jar"].enabled = false

// Run the SmithyBuild task manually since this project needs the built JAR
// from smithy-aws-typescript-codegen.
tasks["smithyBuildJar"].enabled = false

tasks.create<SmithyBuild>("buildSdk") {
    addRuntimeClasspath = true
}

// Generates a smithy-build.json file by creating a new projection for every
// JSON file found in aws-models/. The generated smithy-build.json file is
// not committed to git since it's rebuilt each time codegen is performed.
tasks.register("generate-smithy-build") {
    doLast {
        val projectionsBuilder = Node.objectNodeBuilder()
        val modelsDirProp: String by project
        val models = project.file(modelsDirProp);

        val schemaSerdeEnv = System.getenv("SMITHY_GO_SCHEMA_SERDE") ?: ""
        val forceSchemaSerdeServices = if (schemaSerdeEnv.isNotEmpty()) {
            schemaSerdeEnv.split(",")
        } else {
            emptyList()
        }

        // schema-serde rollout waves
        val w01_rpcv2Cbor = listOf( // 4 svcs
            "com.amazonaws.computeoptimizer#",
            "com.amazonaws.comprehendmedical#",
            "com.amazonaws.marketplaceentitlementservice#",
            "com.amazonaws.snowball#",
        )
        val w02_rpcv2Cbor = listOf( // 5 svcs
            "com.amazonaws.cloudwatch#",
            "com.amazonaws.arcregionswitch#",
            "com.amazonaws.computeoptimizerautomation#",
            "com.amazonaws.gamelift#",
            "com.amazonaws.interconnect#",
        )
        val w03_awsJson = listOf( // 25 svcs
            "com.amazonaws.machinelearning#",
            "com.amazonaws.pinpointsmsvoicev2#",
            "com.amazonaws.migrationhub#",
            "com.amazonaws.migrationhubconfig#",
            "com.amazonaws.kendraranking#",
            "com.amazonaws.kendra#",
            "com.amazonaws.applicationdiscoveryservice#",
            "com.amazonaws.marketplacecommerceanalytics#",
            "com.amazonaws.textract#",
            "com.amazonaws.codestarconnections#",
            "com.amazonaws.fms#",
            "com.amazonaws.marketplacemetering#",
            "com.amazonaws.cloud9#",
            "com.amazonaws.forecast#",
            "com.amazonaws.forecastquery#",
            "com.amazonaws.dax#",
            "com.amazonaws.kinesisanalyticsv2#",
            "com.amazonaws.kinesisanalytics#",
            "com.amazonaws.workmail#",
            "com.amazonaws.proton#",
            "com.amazonaws.memorydb#",
            "com.amazonaws.resourcegroupstaggingapi#",
            "com.amazonaws.codecommit#",
            "com.amazonaws.iotfleetwise#",
            "com.amazonaws.route53domains#",
        )
        val w04_awsJson = listOf( // 25 svcs
            "com.amazonaws.voiceid#",
            "com.amazonaws.globalaccelerator#",
            "com.amazonaws.mediastore#",
            "com.amazonaws.cloudwatchevents#",
            "com.amazonaws.eventbridge#",
            "com.amazonaws.frauddetector#",
            "com.amazonaws.inspector#",
            "com.amazonaws.cloudcontrol#",
            "com.amazonaws.servicediscovery#",
            "com.amazonaws.paymentcryptography#",
            "com.amazonaws.directoryservice#",
            "com.amazonaws.iotsecuretunneling#",
            "com.amazonaws.cloudhsmv2#",
            "com.amazonaws.shield#",
            "com.amazonaws.translate#",
            "com.amazonaws.licensemanager#",
            "com.amazonaws.acmpca#",
            "com.amazonaws.pi#",
            "com.amazonaws.transfer#",
            "com.amazonaws.codebuild#",
            "com.amazonaws.swf#",
            "com.amazonaws.fsx#",
            "com.amazonaws.storagegateway#",
            "com.amazonaws.codepipeline#",
            "com.amazonaws.wafv2#",
        )
        val w05_awsJson = listOf( // 25 svcs
            "com.amazonaws.ecrpublic#",
            "com.amazonaws.cognitoidentity#",
            "com.amazonaws.datasync#",
            "com.amazonaws.organizations#",
            "com.amazonaws.workspaces#",
            "com.amazonaws.comprehend#",
            "com.amazonaws.appstream#",
            "com.amazonaws.personalize#",
            "com.amazonaws.opensearchserverless#",
            "com.amazonaws.apprunner#",
            "com.amazonaws.transcribe#",
            "com.amazonaws.acm#",
            "com.amazonaws.servicequotas#",
            "com.amazonaws.lightsail#",
            "com.amazonaws.emr#",
            "com.amazonaws.route53resolver#",
            "com.amazonaws.cognitoidentityprovider#",
            "com.amazonaws.athena#",
            "com.amazonaws.codedeploy#",
            "com.amazonaws.configservice#",
            "com.amazonaws.partnercentralselling#",
            "com.amazonaws.sqs#",
            "com.amazonaws.mailmanager#",
            "com.amazonaws.sagemaker#",
            "com.amazonaws.cloudtrail#",
        )
        val w06_awsJson = listOf( // 25 svcs
            "com.amazonaws.kinesis#",
            "com.amazonaws.sfn#",
            "com.amazonaws.ecr#",
            "com.amazonaws.secretsmanager#",
            "com.amazonaws.ecs#",
            "com.amazonaws.firehose#",
            "com.amazonaws.kms#",
            "com.amazonaws.cloudwatchlogs#",
            "com.amazonaws.dynamodb#",
            "com.amazonaws.applicationautoscaling#",
            "com.amazonaws.applicationinsights#",
            "com.amazonaws.autoscalingplans#",
            "com.amazonaws.b2bi#",
            "com.amazonaws.backupgateway#",
            "com.amazonaws.bcmdashboards#",
            "com.amazonaws.bcmdataexports#",
            "com.amazonaws.bcmpricingcalculator#",
            "com.amazonaws.bcmrecommendedactions#",
            "com.amazonaws.bedrockdataautomationruntime#",
            "com.amazonaws.billing#",
            "com.amazonaws.budgets#",
            "com.amazonaws.cloudhsm#",
            "com.amazonaws.codeconnections#",
            "com.amazonaws.costandusagereportservice#",
            "com.amazonaws.costexplorer#",
        )
        val w07_awsJson = listOf( // 25 svcs
            "com.amazonaws.costoptimizationhub#",
            "com.amazonaws.databasemigrationservice#",
            "com.amazonaws.datapipeline#",
            "com.amazonaws.devicefarm#",
            "com.amazonaws.directconnect#",
            "com.amazonaws.dynamodbstreams#",
            "com.amazonaws.ec2instanceconnect#",
            "com.amazonaws.evs#",
            "com.amazonaws.freetier#",
            "com.amazonaws.glue#",
            "com.amazonaws.health#",
            "com.amazonaws.healthlake#",
            "com.amazonaws.identitystore#",
            "com.amazonaws.invoicing#",
            "com.amazonaws.iotthingsgraph#",
            "com.amazonaws.keyspaces#",
            "com.amazonaws.keyspacesstreams#",
            "com.amazonaws.lookoutequipment#",
            "com.amazonaws.marketplaceagreement#",
            "com.amazonaws.mturk#",
            "com.amazonaws.mwaaserverless#",
            "com.amazonaws.networkfirewall#",
            "com.amazonaws.odb#",
            "com.amazonaws.partnercentralaccount#",
            "com.amazonaws.partnercentralbenefits#",
        )
        val w08_awsJson = listOf( // 19 svcs
            "com.amazonaws.partnercentralchannel#",
            "com.amazonaws.pcs#",
            "com.amazonaws.pricing#",
            "com.amazonaws.redshiftdata#",
            "com.amazonaws.redshiftserverless#",
            "com.amazonaws.rekognition#",
            "com.amazonaws.route53recoverycluster#",
            "com.amazonaws.servicecatalog#",
            "com.amazonaws.ssm#",
            "com.amazonaws.ssmcontacts#",
            "com.amazonaws.ssoadmin#",
            "com.amazonaws.support#",
            "com.amazonaws.timestreaminfluxdb#",
            "com.amazonaws.timestreamquery#",
            "com.amazonaws.timestreamwrite#",
            "com.amazonaws.verifiedpermissions#",
            "com.amazonaws.waf#",
            "com.amazonaws.wafregional#",
            "com.amazonaws.workspacesinstances#",
        )
        val w09_restJson1 = listOf( // 25 svcs
            "com.amazonaws.mediastoredata#",
            "com.amazonaws.pinpointsmsvoice#",
            "com.amazonaws.taxsettings#",
            "com.amazonaws.chimesdkidentity#",
            "com.amazonaws.chimesdkmediapipelines#",
            "com.amazonaws.chimesdkmeetings#",
            "com.amazonaws.chimesdkmessaging#",
            "com.amazonaws.chimesdkvoice#",
            "com.amazonaws.chime#",
            "com.amazonaws.codeguruprofiler#",
            "com.amazonaws.notifications#",
            "com.amazonaws.ssmsap#",
            "com.amazonaws.notificationscontacts#",
            "com.amazonaws.groundstation#",
            "com.amazonaws.dlm#",
            "com.amazonaws.resourcegroups#",
            "com.amazonaws.securitylake#",
            "com.amazonaws.cognitosync#",
            "com.amazonaws.connectcases#",
            "com.amazonaws.controltower#",
            "com.amazonaws.mgn#",
            "com.amazonaws.lexruntimev2#",
            "com.amazonaws.ssmquicksetup#",
            "com.amazonaws.bedrockagentruntime#",
            "com.amazonaws.appmesh#",
        )
        val w10_restJson1 = listOf( // 25 svcs
            "com.amazonaws.detective#",
            "com.amazonaws.billingconductor#",
            "com.amazonaws.qconnect#",
            "com.amazonaws.wisdom#",
            "com.amazonaws.codegurureviewer#",
            "com.amazonaws.signer#",
            "com.amazonaws.networkmanager#",
            "com.amazonaws.transcribestreaming#",
            "com.amazonaws.outposts#",
            "com.amazonaws.deadline#",
            "com.amazonaws.emrcontainers#",
            "com.amazonaws.servicecatalogappregistry#",
            "com.amazonaws.mediapackagevod#",
            "com.amazonaws.ivsrealtime#",
            "com.amazonaws.ivs#",
            "com.amazonaws.marketplacecatalog#",
            "com.amazonaws.account#",
            "com.amazonaws.mediapackage#",
            "com.amazonaws.drs#",
            "com.amazonaws.resourceexplorer2#",
            "com.amazonaws.workspacesthinclient#",
            "com.amazonaws.xray#",
            "com.amazonaws.iottwinmaker#",
            "com.amazonaws.vpclattice#",
            "com.amazonaws.dataexchange#",
        )
        val w11_restJson1 = listOf( // 25 svcs
            "com.amazonaws.grafana#",
            "com.amazonaws.managedblockchain#",
            "com.amazonaws.dsql#",
            "com.amazonaws.paymentcryptographydata#",
            "com.amazonaws.iot#",
            "com.amazonaws.internetmonitor#",
            "com.amazonaws.neptunegraph#",
            "com.amazonaws.serverlessapplicationrepository#",
            "com.amazonaws.pipes#",
            "com.amazonaws.cloudsearchdomain#",
            "com.amazonaws.amplify#",
            "com.amazonaws.pinpoint#",
            "com.amazonaws.glacier#",
            "com.amazonaws.wellarchitected#",
            "com.amazonaws.connectcontactlens#",
            "com.amazonaws.connect#",
            "com.amazonaws.connectparticipant#",
            "com.amazonaws.lexmodelsv2#",
            "com.amazonaws.amp#",
            "com.amazonaws.schemas#",
            "com.amazonaws.cleanrooms#",
            "com.amazonaws.m2#",
            "com.amazonaws.greengrass#",
            "com.amazonaws.greengrassv2#",
            "com.amazonaws.appflow#",
        )
        val w12_restJson1 = listOf( // 25 svcs
            "com.amazonaws.fis#",
            "com.amazonaws.resiliencehub#",
            "com.amazonaws.resiliencehubv2#",
            "com.amazonaws.inspector2#",
            "com.amazonaws.appsync#",
            "com.amazonaws.accessanalyzer#",
            "com.amazonaws.mediatailor#",
            "com.amazonaws.auditmanager#",
            "com.amazonaws.mediaconnect#",
            "com.amazonaws.mq#",
            "com.amazonaws.bedrockagent#",
            "com.amazonaws.medialive#",
            "com.amazonaws.ram#",
            "com.amazonaws.mediaconvert#",
            "com.amazonaws.bedrock#",
            "com.amazonaws.kafka#",
            "com.amazonaws.batch#",
            "com.amazonaws.appconfig#",
            "com.amazonaws.appconfigdata#",
            "com.amazonaws.securityhub#",
            "com.amazonaws.bedrockruntime#",
            "com.amazonaws.personalizeevents#",
            "com.amazonaws.personalizeruntime#",
            "com.amazonaws.imagebuilder#",
            "com.amazonaws.elasticsearchservice#",
        )
        val w13_restJson1 = listOf( // 25 svcs
            "com.amazonaws.opensearch#",
            "com.amazonaws.backup#",
            "com.amazonaws.quicksight#",
            "com.amazonaws.guardduty#",
            "com.amazonaws.pinpointemail#",
            "com.amazonaws.sesv2#",
            "com.amazonaws.apigateway#",
            "com.amazonaws.apigatewaymanagementapi#",
            "com.amazonaws.apigatewayv2#",
            "com.amazonaws.eks#",
            "com.amazonaws.sagemakeredge#",
            "com.amazonaws.sagemakerfeaturestoreruntime#",
            "com.amazonaws.sagemakermetrics#",
            "com.amazonaws.lambdacore#",
            "com.amazonaws.lambdamicrovms#",
            "com.amazonaws.lambda#",
            "com.amazonaws.aiops#",
            "com.amazonaws.amplifybackend#",
            "com.amazonaws.amplifyuibuilder#",
            "com.amazonaws.appfabric#",
            "com.amazonaws.appintegrations#",
            "com.amazonaws.applicationcostprofiler#",
            "com.amazonaws.applicationsignals#",
            "com.amazonaws.arczonalshift#",
            "com.amazonaws.artifact#",
        )
        val w14_restJson1 = listOf( // 25 svcs
            "com.amazonaws.backupsearch#",
            "com.amazonaws.bedrockagentcore#",
            "com.amazonaws.bedrockagentcorecontrol#",
            "com.amazonaws.bedrockdataautomation#",
            "com.amazonaws.braket#",
            "com.amazonaws.chatbot#",
            "com.amazonaws.cleanroomsml#",
            "com.amazonaws.clouddirectory#",
            "com.amazonaws.cloudfrontkeyvaluestore#",
            "com.amazonaws.cloudtraildata#",
            "com.amazonaws.codeartifact#",
            "com.amazonaws.codecatalyst#",
            "com.amazonaws.codegurusecurity#",
            "com.amazonaws.codestarnotifications#",
            "com.amazonaws.connectcampaigns#",
            "com.amazonaws.connectcampaignsv2#",
            "com.amazonaws.connecthealth#",
            "com.amazonaws.controlcatalog#",
            "com.amazonaws.customerprofiles#",
            "com.amazonaws.databrew#",
            "com.amazonaws.datazone#",
            "com.amazonaws.devopsagent#",
            "com.amazonaws.devopsguru#",
            "com.amazonaws.directoryservicedata#",
            "com.amazonaws.docdbelastic#",
        )
        val w15_restJson1 = listOf( // 25 svcs
            "com.amazonaws.ebs#",
            "com.amazonaws.efs#",
            "com.amazonaws.eksauth#",
            "com.amazonaws.elementalinference#",
            "com.amazonaws.emrserverless#",
            "com.amazonaws.entityresolution#",
            "com.amazonaws.finspace#",
            "com.amazonaws.finspacedata#",
            "com.amazonaws.gameliftstreams#",
            "com.amazonaws.geomaps#",
            "com.amazonaws.geoplaces#",
            "com.amazonaws.georoutes#",
            "com.amazonaws.inspectorscan#",
            "com.amazonaws.iotdataplane#",
            "com.amazonaws.iotdeviceadvisor#",
            "com.amazonaws.iotjobsdataplane#",
            "com.amazonaws.iotmanagedintegrations#",
            "com.amazonaws.iotsitewise#",
            "com.amazonaws.iotwireless#",
            "com.amazonaws.ivschat#",
            "com.amazonaws.kafkaconnect#",
            "com.amazonaws.kinesisvideo#",
            "com.amazonaws.kinesisvideoarchivedmedia#",
            "com.amazonaws.kinesisvideomedia#",
            "com.amazonaws.kinesisvideosignaling#",
        )
        val w16_restJson1 = listOf( // 25 svcs
            "com.amazonaws.kinesisvideowebrtcstorage#",
            "com.amazonaws.lakeformation#",
            "com.amazonaws.launchwizard#",
            "com.amazonaws.lexmodelbuildingservice#",
            "com.amazonaws.lexruntimeservice#",
            "com.amazonaws.licensemanagerlinuxsubscriptions#",
            "com.amazonaws.licensemanagerusersubscriptions#",
            "com.amazonaws.location#",
            "com.amazonaws.macie2#",
            "com.amazonaws.managedblockchainquery#",
            "com.amazonaws.marketplacedeployment#",
            "com.amazonaws.marketplacediscovery#",
            "com.amazonaws.marketplacereporting#",
            "com.amazonaws.mediapackagev2#",
            "com.amazonaws.medicalimaging#",
            "com.amazonaws.migrationhuborchestrator#",
            "com.amazonaws.migrationhubrefactorspaces#",
            "com.amazonaws.migrationhubstrategy#",
            "com.amazonaws.mpa#",
            "com.amazonaws.mwaa#",
            "com.amazonaws.neptunedata#",
            "com.amazonaws.networkflowmonitor#",
            "com.amazonaws.networkmonitor#",
            "com.amazonaws.novaact#",
            "com.amazonaws.oam#",
        )
        val w17_restJson1 = listOf( // 25 svcs
            "com.amazonaws.observabilityadmin#",
            "com.amazonaws.omics#",
            "com.amazonaws.osis#",
            "com.amazonaws.pcaconnectorad#",
            "com.amazonaws.pcaconnectorscep#",
            "com.amazonaws.polly#",
            "com.amazonaws.qapps#",
            "com.amazonaws.qbusiness#",
            "com.amazonaws.rbin#",
            "com.amazonaws.rdsdata#",
            "com.amazonaws.repostspace#",
            "com.amazonaws.rolesanywhere#",
            "com.amazonaws.route53globalresolver#",
            "com.amazonaws.route53profiles#",
            "com.amazonaws.route53recoverycontrolconfig#",
            "com.amazonaws.route53recoveryreadiness#",
            "com.amazonaws.rtbfabric#",
            "com.amazonaws.rum#",
            "com.amazonaws.s3files#",
            "com.amazonaws.s3outposts#",
            "com.amazonaws.s3tables#",
            "com.amazonaws.s3vectors#",
            "com.amazonaws.sagemakera2iruntime#",
            "com.amazonaws.sagemakergeospatial#",
            "com.amazonaws.sagemakerjobruntime#",
        )
        val w18_restJson1 = listOf( // 25 svcs
            "com.amazonaws.sagemakerruntime#",
            "com.amazonaws.sagemakerruntimehttp2#",
            "com.amazonaws.savingsplans#",
            "com.amazonaws.scheduler#",
            "com.amazonaws.securityagent#",
            "com.amazonaws.securityir#",
            "com.amazonaws.signerdata#",
            "com.amazonaws.signin#",
            "com.amazonaws.simpledbv2#",
            "com.amazonaws.snowdevicemanagement#",
            "com.amazonaws.socialmessaging#",
            "com.amazonaws.ssmguiconnect#",
            "com.amazonaws.ssmincidents#",
            "com.amazonaws.sso#",
            "com.amazonaws.ssooidc#",
            "com.amazonaws.supplychain#",
            "com.amazonaws.supportapp#",
            "com.amazonaws.supportauthz#",
            "com.amazonaws.sustainability#",
            "com.amazonaws.synthetics#",
            "com.amazonaws.tnb#",
            "com.amazonaws.trustedadvisor#",
            "com.amazonaws.uxc#",
            "com.amazonaws.wickr#",
            "com.amazonaws.workdocs#",
        )
        val w19_restJson1 = listOf( // 2 svcs
            "com.amazonaws.workmailmessageflow#",
            "com.amazonaws.workspacesweb#",
        )
        val w20_awsQuery = listOf( // 15 svcs
            "com.amazonaws.cloudsearch#",
            "com.amazonaws.elasticbeanstalk#",
            "com.amazonaws.redshift#",
            "com.amazonaws.ses#",
            "com.amazonaws.autoscaling#",
            "com.amazonaws.cloudformation#",
            "com.amazonaws.iam#",
            "com.amazonaws.elasticloadbalancingv2#",
            "com.amazonaws.elasticloadbalancing#",
            "com.amazonaws.docdb#",
            "com.amazonaws.neptune#",
            "com.amazonaws.rds#",
            "com.amazonaws.sts#",
            "com.amazonaws.sns#",
            "com.amazonaws.elasticache#",
        )
        val w21_ec2Query = listOf( // 1 svcs
            "com.amazonaws.ec2#",
        )
        val w22_restXml = listOf( // 3 svcs
            "com.amazonaws.cloudfront#",
            "com.amazonaws.s3control#",
            "com.amazonaws.route53#",
        )
        val w23_s3 = listOf( // 1 svcs
            "com.amazonaws.s3#",
        )

        @OptIn(ExperimentalStdlibApi::class)
        val useLegacySerdeServices = buildList {
            addAll(w01_rpcv2Cbor)
            addAll(w02_rpcv2Cbor)
            addAll(w03_awsJson)
            addAll(w04_awsJson)
            addAll(w05_awsJson)
            addAll(w06_awsJson)
            addAll(w07_awsJson)
            addAll(w08_awsJson)
            addAll(w09_restJson1)
            addAll(w10_restJson1)
            addAll(w11_restJson1)
            addAll(w12_restJson1)
            addAll(w13_restJson1)
            addAll(w14_restJson1)
            addAll(w15_restJson1)
            addAll(w16_restJson1)
            addAll(w17_restJson1)
            addAll(w18_restJson1)
            addAll(w19_restJson1)
            addAll(w20_awsQuery)
            addAll(w21_ec2Query)
            addAll(w22_restXml)
            addAll(w23_s3)
        }

        fileTree(models).filter { it.isFile }.files.forEach eachFile@{ file ->
            val model = Model.assembler()
                    .addImport(file.absolutePath)
                    // Grab the result directly rather than worrying about checking for errors via unwrap.
                    // All we care about here is the service shape, any unchecked errors will be exposed
                    // as part of the actual build task done by the smithy gradle plugin.
                    .assemble().result.get();
            val services = model.shapes(ServiceShape::class.javaObjectType).sorted().toList();
            if (services.size != 1) {
                throw Exception("There must be exactly one service in each aws model file, but found " +
                        "${services.size} in ${file.name}: ${services.map { it.id }}");
            }
            val service = services[0]

            var filteredServices: String = System.getenv("SMITHY_GO_BUILD_API")?: ""
            if (filteredServices.isNotEmpty()) {
                for (filteredService in filteredServices.split(",")) {
                    if (!service.id.toString().startsWith(filteredService)) {
                        return@eachFile
                    }
                }
            }

            val useLegacySerde = !(forceSchemaSerdeServices.any {
                service.id.toString().startsWith(it)
            }) && useLegacySerdeServices.any {
                service.id.toString().startsWith(it)
            }

            val serviceTrait = service.getTrait(ServiceTrait::class.javaObjectType).get();

            val sdkId = serviceTrait.sdkId
                    .replace("-", "")
                    .replace(" ", "")
                    .toLowerCase();
            val projectionContents = Node.objectNodeBuilder()
                    .withMember("imports", Node.fromStrings("${models.absolutePath}${File.separator}${file.name}"))
                    .withMember("plugins", Node.objectNode()
                            .withMember("go-codegen", Node.objectNodeBuilder()
                                    .withMember("service", service.id.toString())
                                    .withMember("module", "github.com/aws/aws-sdk-go-v2/service/$sdkId")
                                    .withMember("useLegacySerde", useLegacySerde)
                                    .build()))
                    .build()
            projectionsBuilder.withMember(sdkId + "." + service.version.toLowerCase(), projectionContents)
        }

        file("smithy-build.json").writeText(Node.prettyPrintJson(Node.objectNodeBuilder()
                .withMember("version", "1.0")
                .withMember("projections", projectionsBuilder.build())
                .build()))
    }
}

// Run the `buildSdk` automatically.
tasks["build"]
        .dependsOn(tasks["generate-smithy-build"])
        .finalizedBy(tasks["buildSdk"])

// ensure built artifacts are put into the SDK's folders
tasks.create<Exec>("copyGoCodegen") {
    dependsOn ("buildSdk")
    commandLine ("$rootDir/copy_go_codegen.sh", "$rootDir/..", (tasks["buildSdk"] as SmithyBuild).outputDirectory.absolutePath)
}
tasks["buildSdk"].finalizedBy(tasks["copyGoCodegen"])
