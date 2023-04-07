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

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.logging.Logger;
import java.util.stream.Collectors;
import java.util.stream.Stream;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.go.codegen.AddOperationShapes;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.knowledge.GoPointableIndex;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.loader.ModelAssembler;
import software.amazon.smithy.model.node.ExpectationNotMetException;
import software.amazon.smithy.model.node.Node;
import software.amazon.smithy.model.node.StringNode;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.DefaultTrait;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.utils.IoUtils;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

public class ApiGatewayExportsNullabilityExceptionIntegrationTest {
    private static final Logger LOGGER = Logger
            .getLogger(ApiGatewayExportsNullabilityExceptionIntegrationTest.class.getName());
    private static final String PATH_PREFIX = "../sdk-codegen/aws-models/";
    private static final String NULLABILITY_EXCEPTIONS_FILE = "APIGW_exports_nullability_exceptions.json";

    /**
     * Sanity test for service models
     */
    @ParameterizedTest
    @MethodSource("apigwNullabilityExceptionServices")
    public void test_APIGW_exports_nullability_exception_services(String modelFile) {
        try {
            loadPreprocessedModel(modelFile);
        } catch (Exception e) {
            LOGGER.severe(e.getMessage());
            throw e;
        }
    }

    private static Stream<Arguments> apigwNullabilityExceptionServices() {
        return Stream.of(
                Arguments.of("greengrass.json"),
                Arguments.of("amplifybackend.json"),
                Arguments.of("mediaconnect.json"),
                Arguments.of("route53-recovery-control-config.json"),
                Arguments.of("pinpoint.json"),
                Arguments.of("apigatewayv2.json"),
                Arguments.of("mediaconvert.json"),
                Arguments.of("medialive.json"),
                Arguments.of("macie2.json"),
                Arguments.of("mediapackage.json"),
                Arguments.of("apigatewaymanagementapi.json"),
                Arguments.of("kafka.json"),
                Arguments.of("mediapackage-vod.json"),
                Arguments.of("mq.json"),
                Arguments.of("iot-1click-devices-service.json"),
                Arguments.of("serverlessapplicationrepository.json"),
                Arguments.of("schemas.json"),
                Arguments.of("pinpoint-sms-voice.json"),
                Arguments.of("route53-recovery-readiness.json"),
                Arguments.of("dataexchange.json"),
                Arguments.of("kafkaconnect.json"),
                Arguments.of("mediatailor.json"));
    }

    /**
     * MediaTailor: APIGW exports -> Smithy migrated service
     *
     * Integration should NOT affected shapes modeled in Smithy after migration.
     *
     * See `mediatailor.json` change in this commit:
     *
     * https://github.com/aws/aws-sdk-go-v2/commit/18e7f160c88a16ad2010ac3208e2939505760e97.
     */
    @Test
    public void test_unaffected_MediaTailor_APIGW_to_Smithy_migrated() {
        Model preprocessedModel = loadPreprocessedModel("mediatailor.json");
        StructureShape segmentationDescriptorShape = preprocessedModel
                .expectShape(ShapeId.from("com.amazonaws.mediatailor#SegmentationDescriptor"))
                .asStructureShape()
                .get();
        assertFalse(segmentationDescriptorShape.getMember("SegmentationEventId").get().hasTrait(DefaultTrait.class));
        GoPointableIndex goPointableIndex = new GoPointableIndex(preprocessedModel);
        List<MemberShape> nullableMemberSnapshot = new ArrayList<>() {
            {
                add(segmentationDescriptorShape.getMember("SegmentationUpidType").get());
                add(segmentationDescriptorShape.getMember("SegmentationUpid").get());
                add(segmentationDescriptorShape.getMember("SegmentationTypeId").get());
                add(segmentationDescriptorShape.getMember("SegmentNum").get());
                add(segmentationDescriptorShape.getMember("SegmentsExpected").get());
                add(segmentationDescriptorShape.getMember("SubSegmentNum").get());
                add(segmentationDescriptorShape.getMember("SubSegmentsExpected").get());
            }
        };
        for (MemberShape memberShape : nullableMemberSnapshot) {
            assertFalse(memberShape.hasTrait(DefaultTrait.class));
            assertTrue(goPointableIndex.isNillable(memberShape.toShapeId()));
            assertTrue(goPointableIndex.isPointable(memberShape.toShapeId()));
            assertTrue(goPointableIndex.isDereferencable(memberShape.toShapeId()));
        }
    }

    /**
     * MediaLive: APIGW exports -> Other default values (e.g. strings) should not
     * be affected.
     *
     * Member com.amazonaws.medialive#DescribeInputDeviceThumbnailResponse$Body
     * should still have a default value of ""
     */
    @Test
    public void test_unaffected_MediaLive_APIGW_default_strings() {
        Model preprocessedModel = loadPreprocessedModel("medialive.json");
        MemberShape bodyShape = preprocessedModel
                .expectShape(ShapeId.from("com.amazonaws.medialive#DescribeInputDeviceThumbnailResponse$Body"))
                .asMemberShape()
                .get();
        assertTrue(bodyShape.hasTrait(DefaultTrait.class));
        assertEquals("",
                bodyShape.expectTrait(DefaultTrait.class).toNode().asStringNode().get().getValue());
    }

    /*
    @Test
    public void test_missing_snapshotted_root_level_shape() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = new ApiGatewayExportsNullabilityExceptionIntegration();
        GoSettings settings = new GoSettings();
        Model model = new ModelAssembler()
                .addDocumentNode(getModel("pinpoint-sms-voice.json"))
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        // TODO(APIGW): Should be deleted after APIGW exports models are fixed.
        model = stripDefaultsFromModel(model, service);
        final Model renamedShapesModel = ModelTransformer.create().renameShapes(model, MapUtils.of(
                ShapeId.from("com.amazonaws.pinpointsmsvoice#Boolean"),
                ShapeId.from("com.amazonaws.pinpointsmsvoice#RenamedBoolean")));
        assertThrows(ExpectationNotMetException.class, () -> integration.preprocessModel(renamedShapesModel, settings));
    }
    */

    /*
    @Test
    public void test_missing_snapshotted_member_level_shape() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = new ApiGatewayExportsNullabilityExceptionIntegration();
        GoSettings settings = new GoSettings();
        Model model = new ModelAssembler()
                .addDocumentNode(getModel("pinpoint-sms-voice.json"))
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        // TODO(APIGW): Should be deleted after APIGW exports models are fixed.
        model = stripDefaultsFromModel(model, service);
        final Model removedShapesModel = ModelTransformer.create().removeShapes(model, ListUtils.of(
                model.expectShape(ShapeId.from("com.amazonaws.pinpointsmsvoice#EventDestination$Enabled"))));
        assertThrows(ExpectationNotMetException.class, () -> integration.preprocessModel(removedShapesModel, settings));
    }
    */

    @Test
    public void test_identify_nonsnapshotted_member_level_shape() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = new ApiGatewayExportsNullabilityExceptionIntegration();
        GoSettings settings = new GoSettings();
        Model model = new ModelAssembler()
                .addDocumentNode(getModel("pinpoint-sms-voice.json"))
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        // TODO(APIGW): Should be deleted after APIGW exports models are fixed.
        model = stripDefaultsFromModel(model, service);
        StructureShape shapeToReplace = model.expectShape(
                ShapeId.from("com.amazonaws.pinpointsmsvoice#EventDestination"),
                StructureShape.class);
        StructureShape replacingShape = shapeToReplace.toBuilder()
                .addMember("NewMember", ShapeId.from("com.amazonaws.pinpointsmsvoice#Boolean"))
                .build();
        final Model replacedShapesModel = ModelTransformer.create().replaceShapes(model, ListUtils.of(replacingShape));
        assertThrows(CodegenException.class, () -> integration.preprocessModel(replacedShapesModel, settings));
    }

    private Model loadPreprocessedModel(String modelFile) {
        ApiGatewayExportsNullabilityExceptionIntegration integration = new ApiGatewayExportsNullabilityExceptionIntegration();
        GoSettings settings = new GoSettings();
        Model model = new ModelAssembler()
                .addDocumentNode(getModel(modelFile))
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        model = AddOperationShapes.execute(model, settings.getService());
        // TODO(APIGW): Should be deleted after APIGW exports models are fixed.
        model = stripDefaultsFromModel(model, service);
        return integration.preprocessModel(model, settings);
    }

    private Node getModel(String modelFile) {
        try {
            return Node.parse(IoUtils.readUtf8File(Path.of(PATH_PREFIX, modelFile)));
        } catch (Exception e) {
            throw new CodegenException(e);
        }
    }

    // TODO(APIGW): Should be deleted after APIGW exports models are fixed.
    private Model stripDefaultsFromModel(Model model, ShapeId service) {
        Set<ShapeId> shapeIdsToReplace = Node
                .parse(IoUtils.readUtf8Resource(getClass(), NULLABILITY_EXCEPTIONS_FILE))
                .expectObjectNode()
                .expectArrayMember(service.toString())
                .getElementsAs(StringNode.class)
                .stream()
                .map(StringNode::getValue)
                .map(ShapeId::from)
                .collect(Collectors.toSet());
        List<Shape> shapesToReplace = new ArrayList<>();
        // Strip root shapes
        for (ShapeId shapeId : shapeIdsToReplace) {
            // TODO: clean this up later
            Optional<Shape> shape = model.getShape(shapeId);
            if (shape.isPresent()) {
                if (shape.get().hasTrait(DefaultTrait.class)) {
                    shapesToReplace.add(Shape.shapeToBuilder(shape.get())
                            .removeTrait(DefaultTrait.ID)
                            .build());
                }
            } else {
                LOGGER.severe("ShapeId `" + shapeId.toString() + "` is not present in the model");
            }
        }
        // Strip member shapes that target affected root shapes
        Set<MemberShape> memberShapes = model.getMemberShapes();
        for (MemberShape shape : memberShapes) {
            ShapeId targetShapeId = shape.getTarget();
            if (!shapeIdsToReplace.contains(targetShapeId)) {
                continue;
            }
            if (shape.hasTrait(DefaultTrait.class)) {
                shapesToReplace.add(Shape.shapeToBuilder(shape)
                        .removeTrait(DefaultTrait.ID)
                        .build());
            }
        }
        // Replace shapes
        Model strippedModel = ModelTransformer.create().replaceShapes(model, shapesToReplace);
        // Assert root shape defaults are removed
        for (ShapeId shapeId : shapeIdsToReplace) {
            if (!strippedModel.getShape(shapeId).isPresent()) {
                continue;
            }
            assertFalse(strippedModel.expectShape(shapeId).hasTrait(DefaultTrait.class));
        }
        // Assert member shape defaults are removed
        for (MemberShape shape : strippedModel.getMemberShapesWithTrait(DefaultTrait.class)) {
            assertFalse(shapeIdsToReplace.contains(shape.getTarget()));
        }
        return strippedModel;
    }
}
