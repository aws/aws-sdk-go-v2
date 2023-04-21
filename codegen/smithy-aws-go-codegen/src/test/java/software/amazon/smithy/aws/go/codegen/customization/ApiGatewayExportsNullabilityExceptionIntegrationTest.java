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
import static software.amazon.smithy.aws.go.codegen.customization.ApiGatewayExportsNullabilityExceptionIntegration.NULLABILITY_EXCEPTIONS_FILE;

import java.net.URL;
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
import software.amazon.smithy.model.node.Node;
import software.amazon.smithy.model.node.StringNode;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.DefaultTrait;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.utils.IoUtils;

public class ApiGatewayExportsNullabilityExceptionIntegrationTest {
    private static final Logger LOGGER = Logger
            .getLogger(ApiGatewayExportsNullabilityExceptionIntegrationTest.class.getName());
    private static final String PATH_PREFIX = "../sdk-codegen/aws-models/";

    @Test
    public void test_default_root_shapes() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = new ApiGatewayExportsNullabilityExceptionIntegration();
        GoSettings settings = new GoSettings();
        URL modelUrl = getClass().getClassLoader().getResource("APIGW_exports_cases/default.smithy");
        Model model = new ModelAssembler()
                .addImport(modelUrl)
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        model = AddOperationShapes.execute(model, settings.getService());
        model = integration.preprocessModel(model, settings);
        assertEquals(false, model.expectShape(ShapeId.from("com.amazonaws.greengrass#__boolean"))
                .expectTrait(DefaultTrait.class)
                .toNode()
                .expectBooleanNode()
                .getValue());
        assertEquals(0L, model.expectShape(ShapeId.from("com.amazonaws.greengrass#__integer"))
                .expectTrait(DefaultTrait.class)
                .toNode()
                .expectNumberNode()
                .getValue());
    }

    @Test
    public void test_nondefault_root_shapes() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = new ApiGatewayExportsNullabilityExceptionIntegration();
        GoSettings settings = new GoSettings();
        URL modelUrl = getClass().getClassLoader().getResource("APIGW_exports_cases/nondefault.smithy");
        Model model = new ModelAssembler()
                .addImport(modelUrl)
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        model = AddOperationShapes.execute(model, settings.getService());
        model = integration.preprocessModel(model, settings);
        assertTrue(model.expectShape(ShapeId.from("com.amazonaws.greengrass#NonTargetBoolean"))
                .getTrait(DefaultTrait.class)
                .isEmpty());
        assertTrue(model.expectShape(ShapeId.from("com.amazonaws.greengrass#NonTargetInteger"))
                .getTrait(DefaultTrait.class)
                .isEmpty());
    }

    @Test
    public void test_existing_members_target_default_root_shapes() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = new ApiGatewayExportsNullabilityExceptionIntegration();
        GoSettings settings = new GoSettings();
        URL modelUrl = getClass().getClassLoader()
                .getResource("APIGW_exports_cases/existing_target_default.smithy");
        Model model = new ModelAssembler()
                .addImport(modelUrl)
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        model = AddOperationShapes.execute(model, settings.getService());
        model = integration.preprocessModel(model, settings);
        assertEquals(0L, model.expectShape(ShapeId.from("com.amazonaws.greengrass#TestStructure$integerMember"))
                .expectTrait(DefaultTrait.class)
                .toNode()
                .expectNumberNode()
                .getValue());
        assertEquals(false,
                model.expectShape(ShapeId.from("com.amazonaws.greengrass#TestStructure$booleanMember"))
                        .expectTrait(DefaultTrait.class)
                        .toNode()
                        .expectBooleanNode()
                        .getValue());
    }

    @Test
    public void test_existing_members_target_nondefault_root_shapes() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = new ApiGatewayExportsNullabilityExceptionIntegration();
        GoSettings settings = new GoSettings();
        URL modelUrl = getClass().getClassLoader()
                .getResource("APIGW_exports_cases/existing_target_nondefault.smithy");
        Model model = new ModelAssembler()
                .addImport(modelUrl)
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        model = AddOperationShapes.execute(model, settings.getService());
        model = integration.preprocessModel(model, settings);
        assertTrue(model.expectShape(ShapeId.from("com.amazonaws.greengrass#TestStructure$integerMember"))
                .getTrait(DefaultTrait.class)
                .isEmpty());
        assertTrue(model.expectShape(ShapeId.from("com.amazonaws.greengrass#TestStructure$booleanMember"))
                .getTrait(DefaultTrait.class)
                .isEmpty());
    }

    private static ApiGatewayExportsNullabilityExceptionIntegration overridePreviousModel(String modelPath) {
        class MockPreviousModelApiGatewayExportsNullabilityExceptionIntegration
                extends ApiGatewayExportsNullabilityExceptionIntegration {
            @Override
            protected Model getPreviousModel(ShapeId service, Model model) {
                URL beforeModelUrl = ApiGatewayExportsNullabilityExceptionIntegrationTest.class
                        .getClassLoader()
                        .getResource(modelPath);
                Model previousModel = new ModelAssembler()
                        .addImport(beforeModelUrl)
                        .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                        .assemble()
                        .unwrap();
                return AddOperationShapes.execute(previousModel, service);
            }
        }
        return new MockPreviousModelApiGatewayExportsNullabilityExceptionIntegration();
    }

    /**
     * TODO(APIGW): change to successful processing after model fix
     */
    @Test
    public void test_adding_new_default_root_shapes() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = overridePreviousModel(
                "APIGW_exports_cases/diffs/added_default.before.smithy");
        GoSettings settings = new GoSettings();
        URL modelUrl = getClass().getClassLoader()
                .getResource("APIGW_exports_cases/diffs/added_default.after.smithy");
        Model model = new ModelAssembler()
                .addImport(modelUrl)
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        final Model thrownModel = AddOperationShapes.execute(model, settings.getService());
        assertThrows(CodegenException.class, () -> {
            try {
                integration.preprocessModel(thrownModel, settings);
            } catch (Exception e) {
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#NewDefaultedInteger: Shape must be added to the APIGW exports nullability exceptions | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#NewDefaultedBoolean: Shape must be added to the APIGW exports nullability exceptions | AwsSdkGoV2ChangedNullability"));
                throw e;
            }
        });
    }

    @Test
    public void test_adding_new_nondefault_root_shapes() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = overridePreviousModel(
                "APIGW_exports_cases/diffs/added_nondefault.before.smithy");
        GoSettings settings = new GoSettings();
        URL modelUrl = getClass().getClassLoader()
                .getResource("APIGW_exports_cases/diffs/added_nondefault.after.smithy");
        Model model = new ModelAssembler()
                .addImport(modelUrl)
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        model = AddOperationShapes.execute(model, settings.getService());
        model = integration.preprocessModel(model, settings);
        assertTrue(model.expectShape(ShapeId.from("com.amazonaws.greengrass#NewNonDefaultedInteger"))
                .getTrait(DefaultTrait.class)
                .isEmpty());
        assertTrue(model.expectShape(ShapeId.from("com.amazonaws.greengrass#NewNonDefaultedBoolean"))
                .getTrait(DefaultTrait.class)
                .isEmpty());
    }

    // TODO(APIGW): provide a green path forward?
    @Test
    public void test_changing_existing_member_target_from_default_to_nondefault() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = overridePreviousModel(
                "APIGW_exports_cases/diffs/existing_target_default_to_nondefault.before.smithy");
        GoSettings settings = new GoSettings();
        URL modelUrl = getClass().getClassLoader()
                .getResource("APIGW_exports_cases/diffs/existing_target_default_to_nondefault.after.smithy");
        Model model = new ModelAssembler()
                .addImport(modelUrl)
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        final Model thrownModel = AddOperationShapes.execute(model, settings.getService());
        assertThrows(CodegenException.class, () -> {
            try {
                integration.preprocessModel(thrownModel, settings);
            } catch (Exception e) {
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$integerMember: Shape changed GoPointableIndex::isDereferencable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$integerMember: Shape changed GoPointableIndex::isNillable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$integerMember: Shape changed GoPointableIndex::isPointable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$integerMember: Shape changed NullableIndex::isNullable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$integerMember: Shape changed NullableIndex::isMemberNullable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$booleanMember: Shape changed GoPointableIndex::isDereferencable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$booleanMember: Shape changed GoPointableIndex::isNillable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$booleanMember: Shape changed GoPointableIndex::isPointable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$booleanMember: Shape changed NullableIndex::isNullable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$booleanMember: Shape changed NullableIndex::isMemberNullable() | AwsSdkGoV2ChangedNullabilit"));
                throw e;
            }
        });
    }

    // TODO(APIGW): provide a green path forward?
    @Test
    public void test_changing_existing_member_target_from_nondefault_to_default() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = overridePreviousModel(
                "APIGW_exports_cases/diffs/existing_target_nondefault_to_default.before.smithy");
        GoSettings settings = new GoSettings();
        URL modelUrl = getClass().getClassLoader()
                .getResource("APIGW_exports_cases/diffs/existing_target_nondefault_to_default.after.smithy");
        Model model = new ModelAssembler()
                .addImport(modelUrl)
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        final Model thrownModel = AddOperationShapes.execute(model, settings.getService());
        assertThrows(CodegenException.class, () -> {
            try {
                integration.preprocessModel(thrownModel, settings);
            } catch (Exception e) {
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$integerMember: Shape changed GoPointableIndex::isDereferencable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$integerMember: Shape changed GoPointableIndex::isNillable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$integerMember: Shape changed GoPointableIndex::isPointable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$integerMember: Shape changed NullableIndex::isNullable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$integerMember: Shape changed NullableIndex::isMemberNullable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$booleanMember: Shape changed GoPointableIndex::isDereferencable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$booleanMember: Shape changed GoPointableIndex::isNillable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$booleanMember: Shape changed GoPointableIndex::isPointable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$booleanMember: Shape changed NullableIndex::isNullable() | AwsSdkGoV2ChangedNullability"));
                assertTrue(e.getMessage().contains(
                        "[ERROR] com.amazonaws.greengrass#TestStructure$booleanMember: Shape changed NullableIndex::isMemberNullable() | AwsSdkGoV2ChangedNullabilit"));
                throw e;
            }
        });
    }

    @Test
    public void test_added_members_target_default_root_shapes() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = overridePreviousModel(
                "APIGW_exports_cases/diffs/added_target_default.before.smithy");
        GoSettings settings = new GoSettings();
        URL modelUrl = getClass().getClassLoader()
                .getResource("APIGW_exports_cases/diffs/added_target_default.after.smithy");
        Model model = new ModelAssembler()
                .addImport(modelUrl)
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        model = AddOperationShapes.execute(model, settings.getService());
        model = integration.preprocessModel(model, settings);
        assertEquals(0L, model.expectShape(ShapeId.from("com.amazonaws.greengrass#TestStructure$integerMember"))
                .expectTrait(DefaultTrait.class)
                .toNode()
                .expectNumberNode()
                .getValue());
        assertEquals(false,
                model.expectShape(ShapeId.from("com.amazonaws.greengrass#TestStructure$booleanMember"))
                        .expectTrait(DefaultTrait.class)
                        .toNode()
                        .expectBooleanNode()
                        .getValue());
    }

    @Test
    public void test_added_members_target_nondefault_root_shapes() {
        ApiGatewayExportsNullabilityExceptionIntegration integration = overridePreviousModel(
                "APIGW_exports_cases/diffs/added_target_nondefault.before.smithy");
        GoSettings settings = new GoSettings();
        URL modelUrl = getClass().getClassLoader()
                .getResource("APIGW_exports_cases/diffs/added_target_nondefault.after.smithy");
        Model model = new ModelAssembler()
                .addImport(modelUrl)
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        model = AddOperationShapes.execute(model, settings.getService());
        model = integration.preprocessModel(model, settings);
        assertTrue(model.expectShape(ShapeId.from("com.amazonaws.greengrass#TestStructure$integerMember"))
                .getTrait(DefaultTrait.class)
                .isEmpty());
        assertTrue(model.expectShape(ShapeId.from("com.amazonaws.greengrass#TestStructure$booleanMember"))
                .getTrait(DefaultTrait.class)
                .isEmpty());
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
            return Node.parseJsonWithComments(IoUtils.readUtf8File(Path.of(PATH_PREFIX, modelFile)));
        } catch (Exception e) {
            throw new CodegenException(e);
        }
    }

    // TODO(APIGW): Should be deleted after APIGW exports models are fixed.
    private Model stripDefaultsFromModel(Model model, ShapeId service) {
        Set<ShapeId> shapeIdsToReplace = Node
                .parseJsonWithComments(
                        IoUtils.readUtf8Resource(getClass(), NULLABILITY_EXCEPTIONS_FILE))
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
