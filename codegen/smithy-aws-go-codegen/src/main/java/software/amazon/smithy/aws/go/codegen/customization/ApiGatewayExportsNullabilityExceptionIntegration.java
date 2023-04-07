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

package software.amazon.smithy.aws.go.codegen.customization;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.logging.Logger;
import java.util.stream.Collectors;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.node.BooleanNode;
import software.amazon.smithy.model.node.Node;
import software.amazon.smithy.model.node.NumberNode;
import software.amazon.smithy.model.node.StringNode;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.NumberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.DefaultTrait;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.utils.IoUtils;
import software.amazon.smithy.utils.SetUtils;
import software.amazon.smithy.utils.SmithyInternalApi;

/**
 * Due to internal model fixes for API Gateway (APIGW) exports, certain shapes
 * and members that target those shapes that used to have defaults became
 * nullable. In order to not break the existing Go v2 client API nullability,
 * this customization adds back default values for affected shapes:
 *
 * - Root Boolean shapes: false
 * - Root Number shapes: 0
 * - Snapshotted Members shapes: inherit defaults
 *
 * The class of services affected are APIGW services and APIGW services that
 * migrated to Smithy, seen in `APIGW_NULLABILITY_EXCEPTION_SERVICES`.
 *
 * A "snapshot" of root-level and member shapes are captured in
 * `APIGW_exports_nullability_exceptions.json`, a mapping of service shape IDs
 * to a list of affected root-level and member shapes.
 * 
 * Cases:
 * - Snapshotted Root level:
 *   - If a snapshotted root level shape doesn't exist, an exception will be
 *     thrown.
 *   - Otherwise, a default will be patched if a default doesn't already exist
 * - Snapshotted Member level:
 *   - If a snapshotted member level shape doesn't exist, an exception will be
 *     thrown.
 *   - Otherwise, a default will be patched if a default doesn't already exist
 * - Nonsnapshotted Member level:
 *   - All nonsnapshotted member level shape that target a snapshotted root level
 *     will be identified and throw an error.
 *   - This is prevent breaking changes if a nonsnapshotted member changes target
 *     from a snapshotted root shape to a nonsnapshotted root shape, as the SDK
 *     has no prior info on either shape.
 */
@SmithyInternalApi
public class ApiGatewayExportsNullabilityExceptionIntegration implements GoIntegration {
    private static final Logger LOGGER = Logger
            .getLogger(ApiGatewayExportsNullabilityExceptionIntegration.class.getName());
    private static final String NULLABILITY_EXCEPTIONS_FILE = "APIGW_exports_nullability_exceptions.json";
    private static final Set<ShapeId> APIGW_NULLABILITY_EXCEPTION_SERVICES = SetUtils.of(
            // APIGW services
            ShapeId.from("com.amazonaws.greengrass#Greengrass"),
            ShapeId.from("com.amazonaws.amplifybackend#AmplifyBackend"),
            ShapeId.from("com.amazonaws.mediaconnect#MediaConnect"),
            ShapeId.from("com.amazonaws.route53recoverycontrolconfig#Route53RecoveryControlConfig"),
            ShapeId.from("com.amazonaws.pinpoint#Pinpoint"),
            ShapeId.from("com.amazonaws.apigatewayv2#ApiGatewayV2"),
            ShapeId.from("com.amazonaws.mediaconvert#MediaConvert"),
            ShapeId.from("com.amazonaws.medialive#MediaLive"),
            ShapeId.from("com.amazonaws.macie2#Macie2"),
            ShapeId.from("com.amazonaws.mediapackage#MediaPackage"),
            ShapeId.from("com.amazonaws.apigatewaymanagementapi#ApiGatewayManagementApi"),
            ShapeId.from("com.amazonaws.kafka#Kafka"),
            ShapeId.from("com.amazonaws.mediapackagevod#MediaPackageVod"),
            ShapeId.from("com.amazonaws.mq#mq"),
            ShapeId.from("com.amazonaws.iot1clickdevicesservice#IoT1ClickDevicesService"),
            ShapeId.from("com.amazonaws.serverlessapplicationrepository#ServerlessApplicationRepository"),
            ShapeId.from("com.amazonaws.schemas#schemas"),
            ShapeId.from("com.amazonaws.pinpointsmsvoice#PinpointSMSVoice"),
            ShapeId.from("com.amazonaws.route53recoveryreadiness#Route53RecoveryReadiness"),
            // APIGW services migrated to Smithy
            ShapeId.from("com.amazonaws.dataexchange#DataExchange"),
            ShapeId.from("com.amazonaws.kafkaconnect#KafkaConnect"),
            ShapeId.from("com.amazonaws.mediatailor#MediaTailor"));

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        ShapeId service = settings.getService();
        if (!APIGW_NULLABILITY_EXCEPTION_SERVICES.contains(service)) {
            return model;
        }
        return handleApiGateWayExportsNullabilityExceptions(model, service);
    }

    private Model handleApiGateWayExportsNullabilityExceptions(Model model, ShapeId service) {
        LOGGER.info("Handling APIGW exports nullability exceptions for service: " + service.toString());

        // Read nullability exceptions
        Set<ShapeId> nullabilityExceptions = Node
                .parse(IoUtils.readUtf8Resource(getClass(), NULLABILITY_EXCEPTIONS_FILE))
                .expectObjectNode()
                .expectArrayMember(service.toString())
                .getElementsAs((StringNode s) -> s.expectShapeId())
                .stream()
                .collect(Collectors.toSet());

        // Knowledge index
        Set<Shape> shapesToReplace = new HashSet<>();
        Set<NumberShape> numberShapes = model.toSet(NumberShape.class);
        Set<MemberShape> memberShapes = model.getMemberShapes();

        patchDefaultsForRootLevelSnapshottedShapes(
                model, nullabilityExceptions, shapesToReplace, numberShapes, memberShapes);
        patchDefaultsForMemberLevelSnapshottedShapes(
                model, nullabilityExceptions, shapesToReplace, numberShapes, memberShapes);
        identityNonSnapshottedMemberLevelShapes(
                nullabilityExceptions, memberShapes);

        // Replace nullability exception shapes
        return ModelTransformer.create().replaceShapes(model, shapesToReplace);
    }

    private void patchDefaultsForRootLevelSnapshottedShapes(
            Model model,
            Set<ShapeId> nullabilityExceptions,
            Set<Shape> shapesToReplace,
            Set<NumberShape> numberShapes,
            Set<MemberShape> memberShapes) {
        for (ShapeId shapeId : nullabilityExceptions) {
            Optional<Shape> shapeOptional = model.getShape(shapeId);
            if (!shapeOptional.isPresent()) {
                LOGGER.severe("ShapeId `" + shapeId.toString() + "` is not present in the model");
                continue;
            }
            Shape shape = shapeOptional.get();
            if (shape.hasTrait(DefaultTrait.class)) {
                continue;
            }
            if (isMemberLevelShape(shape, memberShapes)) {
                continue;
            }
            Boolean isBooleanShape = shape.isBooleanShape();
            Boolean isNumberShape = numberShapes.contains(shape);
            if (!isBooleanShape && !isNumberShape) {
                throw new CodegenException("Root level shape `" + shape.toShapeId().toString()
                        + "` has an invalid shape type `" + shape.getType() + "`");
            }
            DefaultTrait patchedDefaultTrait = new DefaultTrait(isBooleanShape
                    ? BooleanNode.from(false)
                    : NumberNode.from(0L));
            shapesToReplace.add(Shape.shapeToBuilder(shape)
                    .addTrait(patchedDefaultTrait)
                    .build());
        }
    }

    private void patchDefaultsForMemberLevelSnapshottedShapes(
            Model model,
            Set<ShapeId> nullabilityExceptions,
            Set<Shape> shapesToReplace,
            Set<NumberShape> numberShapes,
            Set<MemberShape> memberShapes) {
        for (ShapeId shapeId : nullabilityExceptions) {
            Optional<Shape> shapeOptional = model.getShape(shapeId);
            if (!shapeOptional.isPresent()) {
                LOGGER.severe("ShapeId `" + shapeId.toString() + "` is not present in the model");
                continue;
            }
            Shape shape = shapeOptional.get();
            if (shape.hasTrait(DefaultTrait.class)) {
                continue;
            }
            if (!isMemberLevelShape(shape, memberShapes)) {
                continue;
            }
            patchDefaultForMemberShape(
                    shape.asMemberShape().get(), model, nullabilityExceptions, shapesToReplace, numberShapes);
        }
    }

    private Boolean isMemberLevelShape(Shape shape, Set<MemberShape> memberShapes) {
        return memberShapes.contains(shape);
    }

    private void identityNonSnapshottedMemberLevelShapes(
            Set<ShapeId> nullabilityExceptions,
            Set<MemberShape> memberShapes) {
        List<ShapeId> nonSnapshottedMemberShapes = new ArrayList<>();
        for (MemberShape shape : memberShapes) {
            if (shape.hasTrait(DefaultTrait.class)) {
                continue;
            }
            if (nullabilityExceptions.contains(shape.toShapeId())) {
                continue;
            }
            // Only replace member shapes that target root shape nullability exceptions
            ShapeId targetShapeId = shape.getTarget();
            if (!nullabilityExceptions.contains(targetShapeId)) {
                continue;
            }
            nonSnapshottedMemberShapes.add(shape.toShapeId());
        }
        if (!nonSnapshottedMemberShapes.isEmpty()) {
            throw new CodegenException("Member level shapes that target nullability exception root "
                    + "level shapes are missing from the snapshot. These members MUST be added to the "
                    + "snapshotted shapes in `" + NULLABILITY_EXCEPTIONS_FILE + "` to avoid breaking "
                    + "changes in aws-sdk-go-v2: "
                    + nonSnapshottedMemberShapes.stream()
                            .map(ShapeId::toString)
                            .collect(Collectors.joining(", ")));
        }
    }

    /**
     * Patch default for a member shape
     */
    private void patchDefaultForMemberShape(
            MemberShape shape,
            Model model,
            Set<ShapeId> nullabilityExceptions,
            Set<Shape> shapesToReplace,
            Set<NumberShape> numberShapes) {
        ShapeId targetShapeId = shape.getTarget();
        Shape targetShape = model.expectShape(targetShapeId);
        Boolean isBooleanShape = targetShape.isBooleanShape();
        Boolean isNumberShape = numberShapes.contains(targetShape);
        if (!isBooleanShape && !isNumberShape) {
            throw new CodegenException("Member level shape target `" + targetShapeId.toString()
                    + "` has an invalid target shape type `" + targetShape.getType() + "`");
        }
        DefaultTrait patchedDefaultTrait = new DefaultTrait(isBooleanShape
                ? BooleanNode.from(false)
                : NumberNode.from(0L));
        shapesToReplace.add(Shape.shapeToBuilder(shape)
                .addTrait(patchedDefaultTrait)
                .build());
    }
}
