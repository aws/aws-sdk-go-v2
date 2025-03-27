/*
 * Copyright 2023 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package software.amazon.smithy.aws.go.codegen.customization.service.sqs;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.logging.Logger;
import software.amazon.smithy.go.codegen.CodegenUtils;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.node.Node;
import software.amazon.smithy.model.node.NumberNode;
import software.amazon.smithy.model.shapes.IntegerShape;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.DefaultTrait;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.SmithyUnstableApi;

/**
 * AWS SDK for Go V2 Integrations for the Amazon SQS service
 */
@SmithyUnstableApi
public class SQSCustomizations implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(SQSCustomizations.class.getName());

    private static final ShapeId SQS_SERVICE_ID = ShapeId.from("com.amazonaws.sqs#AmazonSQS");
    private static final ShapeId NON_NULLABLE_INTEGER_ID = ShapeId.from("com.amazonaws.sqs#Integer");

    private static final DefaultTrait DEFAULT_ZERO_TRAIT = new DefaultTrait(NumberNode.from(0L));

    /**
     * Default traits that need to be backfilled
     */
    private static final List<ShapeId> DEFAULT_TRAIT_BACKFILL = ListUtils.of(
            // Structure Shapes
            ShapeId.from("com.amazonaws.sqs#ChangeMessageVisibilityBatchRequestEntry$VisibilityTimeout"),
            ShapeId.from("com.amazonaws.sqs#SendMessageBatchRequestEntry$DelaySeconds"),
            // Top-level Input Shape Members
            ShapeId.from("com.amazonaws.sqs#SendMessageRequest$DelaySeconds"),
            ShapeId.from("com.amazonaws.sqs#ChangeMessageVisibilityRequest$VisibilityTimeout"),
            ShapeId.from("com.amazonaws.sqs#ReceiveMessageRequest$WaitTimeSeconds"),
            ShapeId.from("com.amazonaws.sqs#ReceiveMessageRequest$VisibilityTimeout"),
            ShapeId.from("com.amazonaws.sqs#ReceiveMessageRequest$MaxNumberOfMessages"),
            // Synthetic-equivalent Top-level Input Shape Members
            // Note that "Request" is translated to "Input" in synthetic members
            ShapeId.from("smithy.go.synthetic#SendMessageInput$DelaySeconds"),
            ShapeId.from("smithy.go.synthetic#ChangeMessageVisibilityInput$VisibilityTimeout"),
            ShapeId.from("smithy.go.synthetic#ReceiveMessageInput$WaitTimeSeconds"),
            ShapeId.from("smithy.go.synthetic#ReceiveMessageInput$VisibilityTimeout"),
            ShapeId.from("smithy.go.synthetic#ReceiveMessageInput$MaxNumberOfMessages"));

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        ShapeId serviceId = settings.getService();
        if (!serviceId.equals(SQS_SERVICE_ID)) {
            return model;
        }

        // Add non-nullable integer shape
        model = model.toBuilder()
                .addShape(IntegerShape.builder()
                        .id(NON_NULLABLE_INTEGER_ID)
                        .addTrait(DEFAULT_ZERO_TRAIT)
                        .build())
                .build();

        // Process Default traits that need to be backfilled
        List<Shape> updates = new ArrayList<>();
        for (ShapeId memberShapeId : DEFAULT_TRAIT_BACKFILL) {
            Optional<MemberShape> memberShapeOptional = model
                    .getShape(memberShapeId)
                    .flatMap(s -> s.asMemberShape());
            String memberShapeNamespace = memberShapeId.getNamespace();
            // Synthetic shapes could be missing if the upstream model changes are deduped
            // and the synthetic shapes are no longer needed.
            if (!memberShapeOptional.isPresent() && memberShapeNamespace.equals(CodegenUtils.getSyntheticTypeNamespace())) {
                LOGGER.warning(String.format("SQS service synthetic member shape `" + memberShapeId
                        + "` is not present in the model, so could not be backfilled with a default trait."));
                continue;
            }
            MemberShape memberShape = memberShapeOptional.get();
            // Overwrite default trait to maintain backward compatibility
            if (memberShape.hasTrait(DefaultTrait.class)) {
                DefaultTrait defaultTrait = memberShape.expectTrait(DefaultTrait.class);
                LOGGER.warning("Overwriting default trait for SQS service member shape `" + memberShapeId
                        + "` with value: `" + Node.prettyPrintJson(defaultTrait.toNode()) + "`");
            }
            // Patch member with default trait and change target to non-nullable integer shape
            updates.add(memberShape.toBuilder()
                    .addTrait(DEFAULT_ZERO_TRAIT)
                    .target(NON_NULLABLE_INTEGER_ID)
                    .build());
        }
        return ModelTransformer.create().replaceShapes(model, updates);
    }
}
