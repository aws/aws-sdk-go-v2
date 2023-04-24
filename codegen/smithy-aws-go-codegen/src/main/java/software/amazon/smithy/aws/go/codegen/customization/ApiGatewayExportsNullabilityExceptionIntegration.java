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

import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.logging.Logger;
import java.util.stream.Collectors;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.diff.ModelDiff;
import software.amazon.smithy.go.codegen.AddOperationShapes;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.Synthetic;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.loader.ModelAssembler;
import software.amazon.smithy.model.node.ArrayNode;
import software.amazon.smithy.model.node.BooleanNode;
import software.amazon.smithy.model.node.Node;
import software.amazon.smithy.model.node.NumberNode;
import software.amazon.smithy.model.node.ObjectNode;
import software.amazon.smithy.model.node.StringNode;
import software.amazon.smithy.model.shapes.BooleanShape;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.NumberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.ToShapeId;
import software.amazon.smithy.model.traits.BoxTrait;
import software.amazon.smithy.model.traits.DefaultTrait;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.model.validation.ValidationEvent;
import software.amazon.smithy.utils.IoUtils;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.SetUtils;
import software.amazon.smithy.utils.SmithyInternalApi;

/**
 * Due to internal model fixes for API Gateway (APIGW) exports, certain shapes
 * and members that target those shapes that used to have defaults became
 * nullable. In order to not break the existing Go v2 client API nullability,
 * this customization adds back default values for affected shapes.
 *
 * Nullability exceptions are root shapes captured in
 * `APIGW_exports_nullability_exceptions.json`.
 *
 * Definitions:
 * - “defaulted” shapes are shapes that have the default trait applied.
 * - “non-defaulted” shapes are shapes that have the default trait applied.
 * - “Existing” shapes are shapes that are present in a previous version of a
 *   model
 * - “Added” shapes are shapes that are new in a current version of a model, not
 *   present in a previous version of a model
 *
 * Shape scenarios we need to be aware of when defining the customization:
 *
 * - Existing “defaulted” root boolean or number shapes MUST be backfilled with
 *   a default trait.
 * - TODO(APIGW): remove after model fixes
 *   - Added root boolean or number shapes with default traits MUST be backfilled
 *     with a default trait.
 *     - This is due to not having context of the previous models: there is no way
 *       to tell if the added root shape will be affected by the C2J model fixes
 * - TODO(APIGW): will be true after model fixes
 *   - Added root boolean or number shapes with default traits MUST NOT be
 *     backfilled with a default trait.
 * - Added root boolean or number shapes without default traits MUST NOT be
 *   backfilled with a default trait.
 * - Existing member shapes that target a “defaulted” root shape MUST be
 *   backfilled with a default trait.
 * - Existing member shapes that target a “non-defaulted” root shape MUST NOT be
 *   backfilled with a default trait.
 * - Existing member shapes that change targets from a “defaulted” root shape to
 *   a “non-defaulted” root shape will throw an error.
 *   - TODO(APIGW): provide a green path forward?
 * - Existing member shapes that change targets from a “non-defaulted” root
 *   shape to a “defaulted” root shape will throw an error.
 *   - TODO(APIGW): provide a green path forward?
 * - Added member shapes that target a “defaulted” root shape MUST be backfilled
 *   with a default trait.
 * - Added member shapes that target a “non-defaulted” root shape MUST NOT be
 *   backfilled with a default trait.
 */
@SmithyInternalApi
public class ApiGatewayExportsNullabilityExceptionIntegration implements GoIntegration {
    static final String NULLABILITY_EXCEPTIONS_FILE = "APIGW_exports_nullability_exceptions.json";
    private static final Logger LOGGER = Logger
            .getLogger(ApiGatewayExportsNullabilityExceptionIntegration.class.getName());
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
        Set<ShapeId> nullabilityExceptions = getNullabilityExceptions(service);
        Model previousModel = handleApiGateWayExportsNullabilityExceptions(
                getPreviousModel(service, model), service, nullabilityExceptions, true);
        model = handleApiGateWayExportsNullabilityExceptions(
                model, service, nullabilityExceptions, true);
        List<ValidationEvent> awsSdkGoV2ChangedNullabilityEvents = getAwsSdkGoV2ChangedNullabilityEvents(
                previousModel,
                model);
        if (!awsSdkGoV2ChangedNullabilityEvents.isEmpty()) {
            StringBuilder sb = new StringBuilder().append("AwsSdkGoV2ChangedNullability Validation events found:\n");
            for (ValidationEvent e : awsSdkGoV2ChangedNullabilityEvents) {
                sb.append(" " + e.toString() + "\n");
            }
            throw new CodegenException(sb.toString());
        }
        validateNullabilityExceptions(nullabilityExceptions, model, service);
        return model;
    }

    static Set<ShapeId> getNullabilityExceptions(ToShapeId service) {
        String nullabilityExceptionsString = IoUtils.readUtf8Resource(
                ApiGatewayExportsNullabilityExceptionIntegration.class,
                NULLABILITY_EXCEPTIONS_FILE);
        Set<ShapeId> nullabilityExceptions = Node
                .parseJsonWithComments(nullabilityExceptionsString)
                .expectObjectNode()
                .expectArrayMember(service.toShapeId().toString())
                .getElementsAs((StringNode s) -> s.expectShapeId())
                .stream()
                .collect(Collectors.toSet());
        return nullabilityExceptions;
    }

    protected Model getPreviousModel(ShapeId service, Model model) {
        LOGGER.info("Getting Previous Model for: " + service.toString());
        String DIFF_WORKTREE_BRANCH = "__nullability-worktree-" + service;
        Path root = getRootPath();

        String sdkId = model.getServiceShapesWithTrait(ServiceTrait.class).iterator().next()
                .expectTrait(ServiceTrait.class).getSdkId()
                .replace(" ", "-")
                .toLowerCase();

        String sha = getPreviousModelSha(root, sdkId, service);
        LOGGER.info("Previous model SHA for " + service.toString() + ": " + sha);

        Path worktreePath = Paths.get("/tmp").resolve(DIFF_WORKTREE_BRANCH);
        LOGGER.info("Git Worktree Path for " + service.toString() + ": " + worktreePath);

        String modelPath = worktreePath + "/codegen/sdk-codegen/aws-models/" + sdkId + ".json";
        LOGGER.info("Git Worktree Model Path for " + service.toString() + ": " + modelPath);
        if (!Files.isDirectory(worktreePath)) {
            // First, prune old work trees
            exec(ListUtils.of("git", "worktree", "prune"), root, "Error pruning worktrees");
            // Now create the worktree using a dedicated branch. The branch allows other
            // worktrees to checkout the same branch or SHA without conflicting.
            exec(ListUtils.of("git", "worktree", "add", "--quiet", "--force", "-B", DIFF_WORKTREE_BRANCH,
                    worktreePath.toString(), sha),
                    root, "Unable to create git worktree");
        } else {
            // Checkout the right model version based on commit in the worktree.
            exec(ListUtils.of("git", "checkout", "--quiet", sha, "--", modelPath),
                    worktreePath, "Unable to checkout " + modelPath + "at" + sha + " in git worktree");
        }
        Model previousModel = new ModelAssembler()
                .addImport(modelPath)
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        return AddOperationShapes.execute(previousModel, service);
    }

    /**
     * TODO(APIGW): there must be a better way to resolve the correct path...
     */
    private Path getRootPath() {
        Path root = Paths.get(System.getProperty("user.dir"));
        if (root.getFileName().endsWith("smithy-aws-go-codegen")) {
            root = root.getParent();
        }
        if (root.getFileName().equals("codegen")) {
            throw new CodegenException("Expected codegen/ for file operations");
        }
        return root;
    }

    private static String getPreviousModelSha(Path root, String sdkId, ShapeId service) {
        // Determine the SHA of the previous model
        String modelPath = root.resolve("sdk-codegen/aws-models/" + sdkId + ".json").toString();
        LOGGER.info("Repository Model Path for " + service.toString() + ": " + modelPath);
        List<String> args = ListUtils.of("git", "log", "-1", "--skip", "1",
                "--pretty=format:%h",
                modelPath);
        return exec(args, root, "Invalid git revision").trim();
    }

    private static String exec(List<String> list, Path root, String errorPrefix) {
        StringBuilder output = new StringBuilder();
        int code = IoUtils.runCommand(list, root, output, Collections.emptyMap());
        if (code != 0) {
            throw new CodegenException(errorPrefix + ": " + output);
        }
        return output.toString();
    }

    private static List<ValidationEvent> getAwsSdkGoV2ChangedNullabilityEvents(
            Model previousModel,
            Model currentModel) {
        return ModelDiff.compare(previousModel, currentModel)
                .stream()
                .filter(e -> e.getId().equals(AwsSdkGoV2ChangedNullability.class.getSimpleName()))
                .collect(Collectors.toList());
    }

    private static Model handleApiGateWayExportsNullabilityExceptions(
            Model model,
            ShapeId service,
            Set<ShapeId> nullabilityExceptions,
            boolean relaxed) {
        LOGGER.info("Handling APIGW exports nullability exceptions for service: " + service.toString());

        // Knowledge index
        Set<Shape> shapesToReplace = new HashSet<>();
        Set<NumberShape> numberShapes = model.toSet(NumberShape.class);

        // Patch default traits to nullability exceptions
        for (ShapeId shapeId : nullabilityExceptions) {
            if (relaxed && !model.getShape(shapeId).isPresent()) {
                LOGGER.warning("Shape `" + shapeId + "` nullability exception is not present in the model");
                continue;
            }
            Shape shape = model.expectShape(shapeId);
            if (shape.isBooleanShape()) {
                DefaultTrait patchedDefaultTrait = new DefaultTrait(BooleanNode.from(false));
                shapesToReplace.add(Shape.shapeToBuilder(shape)
                        .removeTrait(BoxTrait.ID)
                        .addTrait(patchedDefaultTrait)
                        .build());
            } else if (numberShapes.contains(shape)) {
                DefaultTrait patchedDefaultTrait = new DefaultTrait(NumberNode.from(0L));
                shapesToReplace.add(Shape.shapeToBuilder(shape)
                        .removeTrait(BoxTrait.ID)
                        .addTrait(patchedDefaultTrait)
                        .build());
            } else {
                throw new CodegenException(
                        "Defaulted root shapes can only be boolean or number shapes, but `"
                                + shapeId + "` of type: " + shape.getType());
            }
        }

        // Patch default traits to members that target nullability exceptions
        for (MemberShape shape : model.toSet(MemberShape.class)) {
            if (!nullabilityExceptions.contains(shape.getTarget())) {
                continue;
            }
            Shape targetShape = model.expectShape(shape.getTarget());
            if (targetShape.isBooleanShape()) {
                DefaultTrait patchedDefaultTrait = new DefaultTrait(BooleanNode.from(false));
                shapesToReplace.add(Shape.shapeToBuilder(shape)
                        .removeTrait(BoxTrait.ID)
                        .addTrait(patchedDefaultTrait)
                        .build());
            } else if (numberShapes.contains(targetShape)) {
                DefaultTrait patchedDefaultTrait = new DefaultTrait(NumberNode.from(0L));
                shapesToReplace.add(Shape.shapeToBuilder(shape)
                        .removeTrait(BoxTrait.ID)
                        .addTrait(patchedDefaultTrait)
                        .build());
            } else {
                throw new CodegenException(
                        "Member shapes can only target boolean or number shapes in nullabity exceptions, but `"
                                + targetShape.toShapeId() + "` is of type: " + targetShape.getType());
            }
        }

        return ModelTransformer.create().replaceShapes(model, shapesToReplace);
    }

    private static void validateNullabilityExceptions(Set<ShapeId> nullabilityExceptions, Model model, ShapeId service) {
        Map<ShapeId, Shape> nullabilityExceptionMap = new HashMap<>();
        for (ShapeId shapeId : nullabilityExceptions) {
            if (model.getShape(shapeId).isPresent()) {
                nullabilityExceptionMap.put(shapeId, model.expectShape(shapeId));
            } else {
                LOGGER.warning("Shape `" + shapeId + "` nullability exception is not present in the model");
            }
        }

        for (BooleanShape shape : model.getBooleanShapesWithTrait(DefaultTrait.class)) {
            ShapeId shapeId = shape.toShapeId();
            String namespace = shapeId.getNamespace();
            if (!namespace.equals(service.getNamespace()) && !namespace.equals(Synthetic.ID.getNamespace())) {
                continue;
            }
            if (!nullabilityExceptions.contains(shapeId)) {
                throw new CodegenException("Shape `" + shapeId + "` should be in nullability exceptions");
            }
        }

        for (NumberShape shape : model.toSet(NumberShape.class).stream()
                .filter(s -> s.hasTrait(DefaultTrait.class))
                .collect(Collectors.toList())) {
            ShapeId shapeId = shape.toShapeId();
            String namespace = shapeId.getNamespace();
            if (!namespace.equals(service.getNamespace()) && !namespace.equals(Synthetic.ID.getNamespace())) {
                continue;
            }
            if (!nullabilityExceptions.contains(shapeId)) {
                throw new CodegenException("Shape `" + shapeId + "` should be in nullability exceptions");
            }
        }

        // Existing “defaulted” root boolean or number shapes MUST be backfilled with a
        // default trait.
        for (Map.Entry<ShapeId, Shape> entry : nullabilityExceptionMap.entrySet()) {
            ShapeId shapeId = entry.getKey();
            Shape shape = entry.getValue();
            DefaultTrait trait = shape.expectTrait(DefaultTrait.class);
            if (shape.isBooleanShape()) {
                if (trait.toNode().expectBooleanNode().getValue() != false) {
                    throw new CodegenException(
                            "Expected nullability exception `" + shapeId + "` to have a default value of false");
                }
            } else { // NumberShape
                if (!trait.toNode().expectNumberNode().getValue().equals(0L)) {
                    throw new CodegenException(
                            "Expected nullability exception `" + shapeId + "` to have a default value of 0");
                }
            }
        }
        // Existing member shapes that target a “defaulted” root shape MUST be
        // backfilled with a default trait.
        for (MemberShape shape : model.getMemberShapes()) {
            ShapeId targetShapeId = shape.getTarget();
            if (!nullabilityExceptions.contains(targetShapeId)) {
                continue;
            }
            if (shape.toShapeId().getName().startsWith("__listOf")
                    || shape.toShapeId().getName().startsWith("MapOf")) {
                continue;
            }
            DefaultTrait trait = shape.expectTrait(DefaultTrait.class);
            DefaultTrait targetTrait = nullabilityExceptionMap.get(targetShapeId).expectTrait(DefaultTrait.class);
            if (!trait.equals(targetTrait)) {
                throw new CodegenException(
                        "Expected member shape `" + shape.toShapeId()
                                + "` to have the same default value as the nullability exception `" + shape.getTarget()
                                + "`");
            }
        }
    }

    /**
     * TODO(APIGW): remove after models are fixed
     */
    private static void writeNullabilityExceptions(ToShapeId service, Model model) {
        String ROOT_NULLABILITY_EXCEPTIONS_FILE = "root_" + service.toShapeId() + "_" + NULLABILITY_EXCEPTIONS_FILE;
        ObjectNode nullabilityExceptions = Node.objectNode();
        List<Shape> shapesToWrite = new ArrayList<>();
        Set<BooleanShape> booleanShapes = model.getBooleanShapes();
        for (BooleanShape shape : booleanShapes) {
            if (!shape.toShapeId().getNamespace().equals("smithy.go.synthetic")
                    && !shape.toShapeId().getNamespace().equals(service.toShapeId().getNamespace())) {
                continue;
            }
            if (shape.hasTrait(DefaultTrait.class)) {
                shapesToWrite.add(shape);
            }
        }
        Set<NumberShape> numberShapes = model.toSet(NumberShape.class);
        for (NumberShape shape : numberShapes) {
            if (!shape.toShapeId().getNamespace().equals("smithy.go.synthetic")
                    && !shape.toShapeId().getNamespace().equals(service.toShapeId().getNamespace())) {
                continue;
            }
            if (shape.hasTrait(DefaultTrait.class)) {
                shapesToWrite.add(shape);
            }
        }
        shapesToWrite.sort((s1, s2) -> s1.toShapeId().compareTo(s2.toShapeId()));
        ArrayNode arrayNode = Node.arrayNode();
        for (Shape shape : shapesToWrite) {
            arrayNode = arrayNode.withValue(StringNode.from(shape.toShapeId().toString()));
        }
        nullabilityExceptions = nullabilityExceptions
                .withMember(service.toShapeId().toString(), arrayNode);
        Path writePath = Paths.get("src/main/resources/software/amazon/smithy/aws/go/codegen/customization",
                ROOT_NULLABILITY_EXCEPTIONS_FILE);
        try {
            LOGGER.info("Writing nullability exceptions for " + service.toShapeId().toString() + ": "
                    + writePath.toAbsolutePath());
            Files.writeString(writePath, Node.prettyPrintJson(nullabilityExceptions));
        } catch (Exception e) {
            throw new CodegenException(e);
        }
    }
}
