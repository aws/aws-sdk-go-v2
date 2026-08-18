/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package software.amazon.smithy.aws.go.codegen.customization.service.s3;

import java.util.HashSet;
import java.util.List;
import java.util.Optional;
import java.util.Set;

import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.HttpPayloadTrait;
import software.amazon.smithy.model.traits.StreamingTrait;

/**
 * Adds middleware to handle S3 response errors with 200 ok status code.
 *
 * Per internal specification, this customization MUST be applied to all S3
 * operations with a structured XML response. The response MUST NOT have a
 * streaming binary payload or event-stream. In Smithy terms, this means all
 * S3 operations whose output does not contain a member with @httpPayload
 * targeting a @streaming blob, @streaming union, or string.
 *
 * Coverage is being rolled out in waves gated by call volume. Once all waves
 * have soaked, remove the wave sets and the isInEnabledWave check — the
 * model-driven logic in supports200Error is the correct final behavior.
 */
public class S3ErrorWith200Status implements GoIntegration {
    private static final String ADD_ERROR_HANDLER_INTERNAL = "HandleResponseErrorWith200Status";

    // Already covered on main — always enabled.
    private static final Set<String> ORIGINAL_OPERATIONS = Set.of(
            "CopyObject", "UploadPartCopy", "CompleteMultipartUpload"
    );

    // Wave 1: < 1B req/week. Very low risk.
    private static final Set<String> WAVE_1_OPERATIONS = Set.of(
            "AbortMultipartUpload",
            "CreateBucketMetadataConfiguration",
            "CreateBucketMetadataTableConfiguration",
            "CreateSession",
            "DeleteBucket",
            "DeleteBucketAnalyticsConfiguration",
            "DeleteBucketCors",
            "DeleteBucketEncryption",
            "DeleteBucketIntelligentTieringConfiguration",
            "DeleteBucketInventoryConfiguration",
            "DeleteBucketLifecycle",
            "DeleteBucketMetadataConfiguration",
            "DeleteBucketMetadataTableConfiguration",
            "DeleteBucketMetricsConfiguration",
            "DeleteBucketOwnershipControls",
            "DeleteBucketPolicy",
            "DeleteBucketReplication",
            "DeleteBucketTagging",
            "DeleteBucketWebsite",
            "DeleteObjectAnnotation",
            "DeleteObjectTagging",
            "DeletePublicAccessBlock",
            "GetBucketAbac",
            "GetBucketAnalyticsConfiguration",
            "GetBucketIntelligentTieringConfiguration",
            "GetBucketInventoryConfiguration",
            "GetBucketMetadataConfiguration",
            "GetBucketMetadataTableConfiguration",
            "GetBucketMetricsConfiguration",
            "GetBucketTagging",
            "GetObjectAcl",
            "GetObjectAttributes",
            "GetObjectLegalHold",
            "GetObjectRetention",
            "GetPublicAccessBlock",
            "ListBucketAnalyticsConfigurations",
            "ListBucketIntelligentTieringConfigurations",
            "ListBucketInventoryConfigurations",
            "ListBucketMetricsConfigurations",
            "ListMultipartUploads",
            "ListObjectAnnotations",
            "ListParts",
            "PutBucketAbac",
            "PutBucketAccelerateConfiguration",
            "PutBucketCors",
            "PutBucketIntelligentTieringConfiguration",
            "PutBucketInventoryConfiguration",
            "PutBucketLifecycleConfiguration",
            "PutBucketLogging",
            "PutBucketMetricsConfiguration",
            "PutBucketNotificationConfiguration",
            "PutBucketOwnershipControls",
            "PutBucketReplication",
            "PutBucketRequestPayment",
            "PutBucketTagging",
            "PutBucketVersioning",
            "PutBucketWebsite",
            "PutObjectAnnotation",
            "PutObjectLegalHold",
            "PutObjectLockConfiguration",
            "PutPublicAccessBlock",
            "RenameObject",
            "RestoreObject",
            "UpdateBucketMetadataAnnotationTableConfiguration",
            "UpdateBucketMetadataInventoryTableConfiguration",
            "UpdateBucketMetadataJournalTableConfiguration",
            "UpdateObjectEncryption",
            "WriteGetObjectResponse"
    );

    // Wave 2: 600M–5B req/week. Low-medium risk.
    private static final Set<String> WAVE_2_OPERATIONS = Set.of(
            "GetBucketAccelerateConfiguration",
            "GetBucketCors",
            "GetBucketLifecycleConfiguration",
            "GetBucketLogging",
            "GetBucketNotificationConfiguration",
            "GetBucketOwnershipControls",
            "GetBucketPolicyStatus",
            "GetBucketReplication",
            "GetBucketRequestPayment",
            "GetBucketVersioning",
            "GetBucketWebsite",
            "GetObjectLockConfiguration",
            "GetPublicAccessBlock",
            "ListBuckets",
            "ListDirectoryBuckets",
            "PutBucketAcl",
            "PutBucketEncryption",
            "PutBucketPolicy"
    );

    // Wave 3: 5B–20B req/week. Medium risk.
    private static final Set<String> WAVE_3_OPERATIONS = Set.of(
            "CreateMultipartUpload",
            "DeleteObjects",
            "GetBucketAcl",
            "GetBucketEncryption",
            "GetBucketLocation",
            "GetObjectTagging",
            "HeadBucket",
            "ListObjectVersions",
            "PutObjectRetention",
            "PutObjectTagging"
    );

    // Wave 4: 50B+ req/week. Medium-high risk (volume).
    private static final Set<String> WAVE_4_OPERATIONS = Set.of(
            "CreateBucket",
            "DeleteObject",
            "HeadObject",
            "ListObjects",
            "ListObjectsV2",
            "PutObject",
            "UploadPart"
    );

    /**
     * Combined set of all currently enabled operations. To enable a wave,
     * add it to this set. To finish rollout, delete all wave sets and the
     * isInEnabledWave check entirely.
     */
    private static final Set<String> ENABLED_OPERATIONS = buildEnabledOperations();

    private static Set<String> buildEnabledOperations() {
        Set<String> enabled = new HashSet<>(ORIGINAL_OPERATIONS);
        enabled.addAll(WAVE_1_OPERATIONS);
        enabled.addAll(WAVE_2_OPERATIONS);
        // enabled.addAll(WAVE_3_OPERATIONS);
        // enabled.addAll(WAVE_4_OPERATIONS);
        return Set.copyOf(enabled);
    }

    @Override
    public byte getOrder() {
        // The associated customization ordering is relative to operation deserializers
        // and thus the integration should be added at the end.
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return List.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate(S3ErrorWith200Status::supports200Error)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ADD_ERROR_HANDLER_INTERNAL,
                                        AwsCustomGoDependency.S3_CUSTOMIZATION).build())
                                .build())
                        .build()
        );
    }

    /**
     * Returns true if the operation supports error response with 200 ok status code.
     *
     * Uses a two-layer check:
     * 1. Model-driven: excludes operations with @httpPayload targeting @streaming or string.
     * 2. Wave gate: only enables operations that are in a currently-active wave.
     *
     * Once all waves are shipped and baked, remove the wave gate (step 2) and
     * the model-driven logic alone becomes the final behavior.
     */
    private static boolean supports200Error(Model model, ServiceShape service, OperationShape operation) {
        if (!isS3Service(model, service)) {
            return false;
        }

        // Model-driven exclusion: this is the correct final logic.
        Optional<ShapeId> output = operation.getOutput();
        if (output.isPresent()) {
            StructureShape outputShape = model.expectShape(output.get(), StructureShape.class);
            boolean excluded = outputShape.getAllMembers().values().stream()
                    .filter(member -> member.hasTrait(HttpPayloadTrait.class))
                    .map(member -> model.expectShape(member.getTarget()))
                    .anyMatch(target -> target.hasTrait(StreamingTrait.class) || target.isStringShape());
            if (excluded) {
                return false;
            }
        }

        // Wave gate: remove this check once all waves are shipped and baked.
        return ENABLED_OPERATIONS.contains(operation.getId().getName(service));
    }

    // returns true if service is s3
    private static boolean isS3Service(Model model, ServiceShape service) {
        return S3ModelUtils.isServiceS3(model, service);
    }
}
