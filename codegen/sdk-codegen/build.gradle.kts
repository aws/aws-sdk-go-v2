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

        // Schema-serde rollout waves. Services use legacy serde until their
        // wave is removed from useLegacySerdeServices below. Waves are ordered
        // by blast radius (based on SDK DW request volume, May 2026).
        val wave1 = listOf( // no metrics (219 services)
            "com.amazonaws.aiops#",
            "com.amazonaws.amplifybackend#",
            "com.amazonaws.amplifyuibuilder#",
            "com.amazonaws.apigatewaymanagementapi#",
            "com.amazonaws.appconfigdata#",
            "com.amazonaws.appfabric#",
            "com.amazonaws.appintegrations#",

            "com.amazonaws.applicationcostprofiler#",
            "com.amazonaws.applicationinsights#",
            "com.amazonaws.applicationsignals#",
            "com.amazonaws.arcregionswitch#",
            "com.amazonaws.arczonalshift#",
            "com.amazonaws.artifact#",
            "com.amazonaws.autoscalingplans#",
            "com.amazonaws.b2bi#",
            "com.amazonaws.backupgateway#",
            "com.amazonaws.backupsearch#",
            "com.amazonaws.bcmdashboards#",
            "com.amazonaws.bcmdataexports#",
            "com.amazonaws.bcmpricingcalculator#",
            "com.amazonaws.bcmrecommendedactions#",
            "com.amazonaws.bedrockagentcore#",
            "com.amazonaws.bedrockagentcorecontrol#",
            "com.amazonaws.bedrockdataautomation#",
            "com.amazonaws.bedrockdataautomationruntime#",
            "com.amazonaws.billing#",
            "com.amazonaws.braket#",
            "com.amazonaws.budgets#",
            "com.amazonaws.chatbot#",
            "com.amazonaws.chimesdkidentity#",
            "com.amazonaws.chimesdkmediapipelines#",
            "com.amazonaws.chimesdkmeetings#",
            "com.amazonaws.chimesdkmessaging#",
            "com.amazonaws.chimesdkvoice#",
            "com.amazonaws.cleanroomsml#",
            "com.amazonaws.clouddirectory#",
            "com.amazonaws.cloudfrontkeyvaluestore#",
            "com.amazonaws.cloudhsm#",
            "com.amazonaws.cloudsearchdomain#",
            "com.amazonaws.cloudtraildata#",
            "com.amazonaws.cloudwatchevents#",
            "com.amazonaws.codeartifact#",
            "com.amazonaws.codecatalyst#",
            "com.amazonaws.codeconnections#",
            "com.amazonaws.codegurusecurity#",

            "com.amazonaws.computeoptimizerautomation#",
            "com.amazonaws.connectcampaigns#",
            "com.amazonaws.connectcampaignsv2#",
            "com.amazonaws.connectcontactlens#",
            "com.amazonaws.connecthealth#",
            "com.amazonaws.connectparticipant#",
            "com.amazonaws.controlcatalog#",
            "com.amazonaws.costandusagereportservice#",
            "com.amazonaws.costoptimizationhub#",


            "com.amazonaws.databrew#",
            "com.amazonaws.datapipeline#",
            "com.amazonaws.datazone#",
            "com.amazonaws.devicefarm#",
            "com.amazonaws.devopsagent#",
            "com.amazonaws.devopsguru#",
            "com.amazonaws.directconnect#",
            "com.amazonaws.directoryservicedata#",
            "com.amazonaws.docdb#",
            "com.amazonaws.docdbelastic#",

            "com.amazonaws.ebs#",
            "com.amazonaws.ec2instanceconnect#",
            "com.amazonaws.eksauth#",
            "com.amazonaws.elasticloadbalancing#",
            "com.amazonaws.elementalinference#",
            "com.amazonaws.emrserverless#",
            "com.amazonaws.entityresolution#",
            "com.amazonaws.evs#",
            "com.amazonaws.finspace#",
            "com.amazonaws.finspacedata#",
            "com.amazonaws.forecastquery#",
            "com.amazonaws.freetier#",

            "com.amazonaws.gameliftstreams#",
            "com.amazonaws.geomaps#",
            "com.amazonaws.geoplaces#",
            "com.amazonaws.georoutes#",
            "com.amazonaws.greengrassv2#",
            "com.amazonaws.health#",
            "com.amazonaws.healthlake#",
            "com.amazonaws.identitystore#",
            "com.amazonaws.inspectorscan#",
            "com.amazonaws.interconnect#",
            "com.amazonaws.invoicing#",
            "com.amazonaws.iotdeviceadvisor#",
            "com.amazonaws.ioteventsdata#",
            "com.amazonaws.iotmanagedintegrations#",
            "com.amazonaws.iotsitewise#",
            "com.amazonaws.iotthingsgraph#",
            "com.amazonaws.iotwireless#",
            "com.amazonaws.ivschat#",
            "com.amazonaws.ivsrealtime#",
            "com.amazonaws.kafkaconnect#",
            "com.amazonaws.kendra#",
            "com.amazonaws.keyspacesstreams#",
            "com.amazonaws.kinesisanalyticsv2#",
            "com.amazonaws.kinesisvideo#",
            "com.amazonaws.kinesisvideoarchivedmedia#",
            "com.amazonaws.kinesisvideomedia#",
            "com.amazonaws.kinesisvideosignaling#",
            "com.amazonaws.kinesisvideowebrtcstorage#",
            "com.amazonaws.lakeformation#",
            "com.amazonaws.launchwizard#",
            "com.amazonaws.lexmodelbuildingservice#",
            "com.amazonaws.lexruntimeservice#",
            "com.amazonaws.licensemanagerlinuxsubscriptions#",
            "com.amazonaws.licensemanagerusersubscriptions#",
            "com.amazonaws.location#",
            "com.amazonaws.lookoutequipment#",
            "com.amazonaws.macie2#",
            "com.amazonaws.mailmanager#",
            "com.amazonaws.managedblockchainquery#",
            "com.amazonaws.marketplaceagreement#",
            "com.amazonaws.marketplacedeployment#",
            "com.amazonaws.marketplacediscovery#",
            "com.amazonaws.marketplacereporting#",
            "com.amazonaws.mediapackagev2#",
            "com.amazonaws.mediastoredata#",
            "com.amazonaws.medicalimaging#",
            "com.amazonaws.migrationhubconfig#",
            "com.amazonaws.migrationhubrefactorspaces#",
            "com.amazonaws.migrationhubstrategy#",
            "com.amazonaws.mpa#",
            "com.amazonaws.mturk#",
            "com.amazonaws.mwaa#",
            "com.amazonaws.mwaaserverless#",
            "com.amazonaws.neptune#",
            "com.amazonaws.neptunedata#",
            "com.amazonaws.networkfirewall#",
            "com.amazonaws.networkflowmonitor#",
            "com.amazonaws.networkmonitor#",
            "com.amazonaws.novaact#",
            "com.amazonaws.oam#",
            "com.amazonaws.observabilityadmin#",
            "com.amazonaws.odb#",
            "com.amazonaws.opensearch#",
            "com.amazonaws.osis#",
            "com.amazonaws.panorama#",
            "com.amazonaws.partnercentralaccount#",
            "com.amazonaws.partnercentralbenefits#",
            "com.amazonaws.partnercentralchannel#",
            "com.amazonaws.paymentcryptographydata#",
            "com.amazonaws.pcaconnectorad#",
            "com.amazonaws.pcaconnectorscep#",
            "com.amazonaws.pcs#",
            "com.amazonaws.personalizeevents#",
            "com.amazonaws.personalizeruntime#",
            "com.amazonaws.pinpointemail#",
            "com.amazonaws.pinpointsmsvoicev2#",

            "com.amazonaws.pricing#",
            "com.amazonaws.qapps#",
            "com.amazonaws.qbusiness#",
            "com.amazonaws.qconnect#",
            "com.amazonaws.rbin#",
            "com.amazonaws.redshiftdata#",
            "com.amazonaws.redshiftserverless#",

            "com.amazonaws.repostspace#",
            "com.amazonaws.rolesanywhere#",
            "com.amazonaws.route53globalresolver#",
            "com.amazonaws.route53profiles#",
            "com.amazonaws.route53recoverycluster#",
            "com.amazonaws.route53recoverycontrolconfig#",
            "com.amazonaws.route53recoveryreadiness#",
            "com.amazonaws.rtbfabric#",
            "com.amazonaws.rum#",
            "com.amazonaws.sagemakera2iruntime#",
            "com.amazonaws.sagemakeredge#",
            "com.amazonaws.sagemakerfeaturestoreruntime#",
            "com.amazonaws.sagemakergeospatial#",
            "com.amazonaws.sagemakermetrics#",
            "com.amazonaws.sagemakerruntime#",
            "com.amazonaws.sagemakerruntimehttp2#",
            "com.amazonaws.savingsplans#",
            "com.amazonaws.scheduler#",
            "com.amazonaws.securityagent#",
            "com.amazonaws.securityir#",
            "com.amazonaws.servicecatalog#",
            "com.amazonaws.servicequotas#",
            "com.amazonaws.signerdata#",
            "com.amazonaws.simpledbv2#",
            "com.amazonaws.simspaceweaver#",
            "com.amazonaws.socialmessaging#",
            "com.amazonaws.ssmcontacts#",
            "com.amazonaws.ssmguiconnect#",
            "com.amazonaws.ssmincidents#",
            "com.amazonaws.ssmsap#",
            "com.amazonaws.ssoadmin#",
            "com.amazonaws.supplychain#",
            "com.amazonaws.support#",
            "com.amazonaws.supportapp#",
            "com.amazonaws.sustainability#",
            "com.amazonaws.synthetics#",
            "com.amazonaws.timestreaminfluxdb#",
            "com.amazonaws.timestreamquery#",
            "com.amazonaws.timestreamwrite#",
            "com.amazonaws.tnb#",
            "com.amazonaws.trustedadvisor#",
            "com.amazonaws.uxc#",
            "com.amazonaws.verifiedpermissions#",
            "com.amazonaws.waf#",
            "com.amazonaws.wafregional#",
            "com.amazonaws.wickr#",
            "com.amazonaws.workdocs#",
            "com.amazonaws.workmailmessageflow#",
            "com.amazonaws.workspacesinstances#",
            "com.amazonaws.workspacesthinclient#",
            "com.amazonaws.workspacesweb#",
        )
        val wave2 = listOf( // bottom half + missing-but-used (97 services)
            "com.amazonaws.account#",
            "com.amazonaws.amplify#",
            "com.amazonaws.applicationdiscoveryservice#",
            "com.amazonaws.appmesh#",
            "com.amazonaws.bedrockagentruntime#",
            "com.amazonaws.billingconductor#",
            "com.amazonaws.chime#",
            "com.amazonaws.cloud9#",
            "com.amazonaws.cloudcontrol#",
            "com.amazonaws.cloudhsmv2#",
            "com.amazonaws.cloudsearch#",
            "com.amazonaws.codecommit#",
            "com.amazonaws.codeguruprofiler#",
            "com.amazonaws.codegurureviewer#",
            "com.amazonaws.codestarconnections#",
            "com.amazonaws.codestarnotifications#",
            "com.amazonaws.cognitosync#",
            "com.amazonaws.comprehendmedical#",
            "com.amazonaws.connectcases#",
            "com.amazonaws.controltower#",
            "com.amazonaws.costexplorer#",
            "com.amazonaws.dataexchange#",
            "com.amazonaws.dax#",
            "com.amazonaws.deadline#",
            "com.amazonaws.detective#",
            "com.amazonaws.directoryservice#",
            "com.amazonaws.dlm#",
            "com.amazonaws.drs#",
            "com.amazonaws.dsql#",
            "com.amazonaws.efs#",
            "com.amazonaws.elasticache#",
            "com.amazonaws.emrcontainers#",
            "com.amazonaws.eventbridge#",
            "com.amazonaws.fms#",
            "com.amazonaws.forecast#",
            "com.amazonaws.frauddetector#",
            "com.amazonaws.globalaccelerator#",
            "com.amazonaws.glue#",
            "com.amazonaws.grafana#",
            "com.amazonaws.groundstation#",
            "com.amazonaws.inspector#",
            "com.amazonaws.internetmonitor#",
            "com.amazonaws.iotevents#",
            "com.amazonaws.iotfleetwise#",
            "com.amazonaws.iotjobsdataplane#",
            "com.amazonaws.iotsecuretunneling#",
            "com.amazonaws.iottwinmaker#",
            "com.amazonaws.ivs#",
            "com.amazonaws.kendraranking#",
            "com.amazonaws.keyspaces#",
            "com.amazonaws.kinesisanalytics#",
            "com.amazonaws.lexruntimev2#",
            "com.amazonaws.machinelearning#",
            "com.amazonaws.managedblockchain#",
            "com.amazonaws.marketplacecatalog#",
            "com.amazonaws.marketplacecommerceanalytics#",
            "com.amazonaws.marketplaceentitlementservice#",
            "com.amazonaws.marketplacemetering#",
            "com.amazonaws.mediapackage#",
            "com.amazonaws.mediapackagevod#",
            "com.amazonaws.mediastore#",
            "com.amazonaws.memorydb#",
            "com.amazonaws.mgn#",
            "com.amazonaws.migrationhub#",
            "com.amazonaws.migrationhuborchestrator#",
            "com.amazonaws.neptunegraph#",
            "com.amazonaws.networkmanager#",
            "com.amazonaws.notifications#",
            "com.amazonaws.notificationscontacts#",
            "com.amazonaws.outposts#",
            "com.amazonaws.paymentcryptography#",
            "com.amazonaws.pinpoint#",
            "com.amazonaws.pinpointsmsvoice#",
            "com.amazonaws.pipes#",
            "com.amazonaws.proton#",
            "com.amazonaws.rdsdata#",
            "com.amazonaws.resourceexplorer2#",
            "com.amazonaws.resourcegroups#",
            "com.amazonaws.resourcegroupstaggingapi#",
            "com.amazonaws.route53domains#",
            "com.amazonaws.securitylake#",
            "com.amazonaws.serverlessapplicationrepository#",
            "com.amazonaws.servicecatalogappregistry#",
            "com.amazonaws.servicediscovery#",
            "com.amazonaws.sesv2#",
            "com.amazonaws.signer#",
            "com.amazonaws.snowball#",
            "com.amazonaws.snowdevicemanagement#",
            "com.amazonaws.ssmquicksetup#",
            "com.amazonaws.taxsettings#",
            "com.amazonaws.textract#",
            "com.amazonaws.transcribestreaming#",
            "com.amazonaws.voiceid#",
            "com.amazonaws.vpclattice#",
            "com.amazonaws.wisdom#",
            "com.amazonaws.workmail#",
            "com.amazonaws.xray#",
        )
        val wave3 = listOf( // top half (92 services)
            "com.amazonaws.accessanalyzer#",
            "com.amazonaws.acm#",
            "com.amazonaws.acmpca#",
            "com.amazonaws.amp#",
            "com.amazonaws.apigateway#",
            "com.amazonaws.apigatewayv2#",
            "com.amazonaws.appconfig#",
            "com.amazonaws.appflow#",
            "com.amazonaws.apprunner#",
            "com.amazonaws.appstream#",
            "com.amazonaws.appsync#",
            "com.amazonaws.athena#",
            "com.amazonaws.auditmanager#",
            "com.amazonaws.autoscaling#",
            "com.amazonaws.backup#",
            "com.amazonaws.batch#",
            "com.amazonaws.bedrock#",
            "com.amazonaws.bedrockagent#",
            "com.amazonaws.bedrockruntime#",
            "com.amazonaws.cleanrooms#",
            "com.amazonaws.cloudformation#",
            "com.amazonaws.cloudfront#",
            "com.amazonaws.cloudtrail#",
            "com.amazonaws.codebuild#",
            "com.amazonaws.codedeploy#",
            "com.amazonaws.codepipeline#",
            "com.amazonaws.cognitoidentity#",
            "com.amazonaws.cognitoidentityprovider#",
            "com.amazonaws.comprehend#",
            "com.amazonaws.configservice#",
            "com.amazonaws.connect#",
            "com.amazonaws.datasync#",
            "com.amazonaws.ecr#",
            "com.amazonaws.ecrpublic#",
            "com.amazonaws.ecs#",
            "com.amazonaws.eks#",
            "com.amazonaws.elasticbeanstalk#",
            "com.amazonaws.elasticloadbalancingv2#",
            "com.amazonaws.elasticsearchservice#",
            "com.amazonaws.emr#",
            "com.amazonaws.firehose#",
            "com.amazonaws.fis#",
            "com.amazonaws.fsx#",
            "com.amazonaws.glacier#",
            "com.amazonaws.greengrass#",
            "com.amazonaws.guardduty#",
            "com.amazonaws.imagebuilder#",
            "com.amazonaws.inspector2#",
            "com.amazonaws.iot#",
            "com.amazonaws.iotdataplane#",
            "com.amazonaws.kafka#",
            "com.amazonaws.kinesis#",
            "com.amazonaws.kms#",
            "com.amazonaws.lambda#",
            "com.amazonaws.lexmodelsv2#",
            "com.amazonaws.licensemanager#",
            "com.amazonaws.lightsail#",
            "com.amazonaws.m2#",
            "com.amazonaws.mediaconnect#",
            "com.amazonaws.mediaconvert#",
            "com.amazonaws.medialive#",
            "com.amazonaws.mediatailor#",
            "com.amazonaws.mq#",
            "com.amazonaws.omics#",
            "com.amazonaws.opensearchserverless#",
            "com.amazonaws.organizations#",
            "com.amazonaws.partnercentralselling#",
            "com.amazonaws.personalize#",
            "com.amazonaws.pi#",
            "com.amazonaws.quicksight#",
            "com.amazonaws.ram#",
            "com.amazonaws.rds#",
            "com.amazonaws.redshift#",
            "com.amazonaws.resiliencehub#",
            "com.amazonaws.route53resolver#",
            "com.amazonaws.sagemaker#",
            "com.amazonaws.schemas#",
            "com.amazonaws.secretsmanager#",
            "com.amazonaws.securityhub#",
            "com.amazonaws.ses#",
            "com.amazonaws.sfn#",
            "com.amazonaws.shield#",
            "com.amazonaws.sqs#",
            "com.amazonaws.ssm#",
            "com.amazonaws.storagegateway#",
            "com.amazonaws.swf#",
            "com.amazonaws.transcribe#",
            "com.amazonaws.transfer#",
            "com.amazonaws.translate#",
            "com.amazonaws.wafv2#",
            "com.amazonaws.wellarchitected#",
            "com.amazonaws.workspaces#",
        )
        val wave4 = listOf( // credential providers
            "com.amazonaws.iam#",
            "com.amazonaws.signin#",
            "com.amazonaws.sso#",
            "com.amazonaws.ssooidc#",
            "com.amazonaws.sts#",
        )
        val wave5 = listOf( // tier zero
            "com.amazonaws.applicationautoscaling#",
            "com.amazonaws.cloudwatch#",
            "com.amazonaws.cloudwatchlogs#",
            "com.amazonaws.computeoptimizer#",
            "com.amazonaws.customerprofiles#",
            "com.amazonaws.databasemigrationservice#",
            "com.amazonaws.dynamodb#",
            "com.amazonaws.dynamodbstreams#",
            "com.amazonaws.ec2#",
            "com.amazonaws.gamelift#",
            "com.amazonaws.polly#",
            "com.amazonaws.rekognition#",
            "com.amazonaws.sns#",
        )
        val wave6 = listOf( // S3, route53
            "com.amazonaws.s3#",
            "com.amazonaws.s3control#",
            "com.amazonaws.s3files#",
            "com.amazonaws.s3outposts#",
            "com.amazonaws.s3tables#",
            "com.amazonaws.s3vectors#",
            "com.amazonaws.route53#",
        )

        @OptIn(ExperimentalStdlibApi::class)
        val useLegacySerdeServices = buildList {
            // addAll(wave1) (initial release)
            addAll(wave2)
            addAll(wave3)
            addAll(wave4)
            addAll(wave5)
            addAll(wave6)
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
