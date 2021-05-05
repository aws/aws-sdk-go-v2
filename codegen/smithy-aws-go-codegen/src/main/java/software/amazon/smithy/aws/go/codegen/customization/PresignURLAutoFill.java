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
 *
 *
 */

package software.amazon.smithy.aws.go.codegen.customization;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeMap;
import java.util.logging.Logger;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.CodegenUtils;
import software.amazon.smithy.go.codegen.GoCodegenPlugin;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.trait.NoSerializeTrait;
import software.amazon.smithy.go.codegen.trait.UnexportedMemberTrait;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StringShape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.DocumentationTrait;
import software.amazon.smithy.utils.IoUtils;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.utils.SetUtils;

/**
 * PresignURLAutoFill represents a runtime plugin integration class
 * used generate customization to autofill a presign url as
 * an unexported serialized input member
 */
public class PresignURLAutoFill implements GoIntegration {
    /**
     * Map of service shape to Set of operation shapes that need to have this
     * presigned url auto fill customization.
     */
    public static final Map<ShapeId, Set<ShapeId>> SERVICE_TO_OPERATION_MAP = MapUtils.of(
            ShapeId.from("com.amazonaws.rds#AmazonRDSv19"), SetUtils.of(
                    ShapeId.from("com.amazonaws.rds#CopyDBSnapshot"),
                    ShapeId.from("com.amazonaws.rds#CreateDBInstanceReadReplica"),
                    ShapeId.from("com.amazonaws.rds#CopyDBClusterSnapshot"),
                    ShapeId.from("com.amazonaws.rds#CreateDBCluster")),

            ShapeId.from("com.amazonaws.neptune#AmazonRDSv19"), SetUtils.of(
                    ShapeId.from("com.amazonaws.neptune#CopyDBClusterSnapshot"),
                    ShapeId.from("com.amazonaws.neptune#CreateDBCluster")),

            ShapeId.from("com.amazonaws.docdb#AmazonRDSv19"), SetUtils.of(
                    ShapeId.from("com.amazonaws.docdb#CopyDBClusterSnapshot"),
                    ShapeId.from("com.amazonaws.docdb#CreateDBCluster")),

            ShapeId.from("com.amazonaws.ec2#AmazonEC2"), SetUtils.of(
                    ShapeId.from("com.amazonaws.ec2#CopySnapshot"))

            // TODO other services
    );
    private static final Logger LOGGER = Logger.getLogger(PresignURLAutoFill.class.getName());
    private final List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();

    private static void writeMemberSetter(
            GoWriter writer,
            SymbolProvider symbolprovider,
            OperationShape operation,
            StructureShape input,
            MemberShape member
    ) {
        Symbol operationSymbol = symbolprovider.toSymbol(operation);
        Symbol inputSymbol = symbolprovider.toSymbol(input);
        String memberName = symbolprovider.toMemberName(member);

        writer.openBlock("func $L(params interface{}, value string) error {", "}",
                setterFuncName(operationSymbol.getName(), memberName),
                () -> {
                    writer.addUseImports(SmithyGoDependency.FMT);
                    writer.write("input, ok := params.($P)", inputSymbol);
                    writer.openBlock("if !ok {", "}", () -> {
                        writer.write("return fmt.Errorf(\"expect $P type, got %T\", params)", inputSymbol);
                    });
                    writer.write("input.$L = &value", memberName);
                    writer.write("return nil");
                });
    }

    private static void writeMemberGetter(
            GoWriter writer,
            SymbolProvider symbolprovider,
            OperationShape operation,
            StructureShape input,
            MemberShape member
    ) {
        Symbol operationSymbol = symbolprovider.toSymbol(operation);
        Symbol inputSymbol = symbolprovider.toSymbol(input);
        String memberName = symbolprovider.toMemberName(member);

        writer.openBlock("func $L(params interface{}) (string, bool, error) {", "}",
                getterFuncName(operationSymbol.getName(), memberName),
                () -> {
                    writer.addUseImports(SmithyGoDependency.FMT);
                    writer.write("input, ok := params.($P)", inputSymbol);
                    writer.openBlock("if !ok {", "}", () -> {
                        writer.write("return ``, false, fmt.Errorf(\"expect $P type, got %T\", params)", inputSymbol);
                    });
                    writer.openBlock("if input.$L == nil || len(*input.$L) == 0 {", "}", memberName, memberName,
                            () -> writer.write("return ``, false, nil")
                    );
                    writer.write("return *input.$L, true, nil", memberName);
                });
    }

    private static String addPresignMiddlewareFuncName(String operationName) {
        return String.format("add%sPresignURLMiddleware", operationName);
    }

    private static String getterFuncName(String operationName, String memberName) {
        return String.format("get%s%s", operationName, memberName);
    }

    private static String setterFuncName(String operationName, String memberName) {
        return String.format("set%s%s", operationName, memberName);
    }

    private static String copyInputFuncName(String inputName) {
        return String.format("copy%sForPresign", inputName);
    }

    private static String presignFuncName(String operationName) {
        return String.format("Presign%s", operationName);
    }

    private static String autofillPresignClient(String operationName) {
        return String.format("presignAutoFill%sClient", operationName);
    }

    /**
     * Updates the API model to add additional members to the operation input shape that are needed for presign url
     * customization.
     *
     * @param model    API model
     * @param settings Go codegen settings
     * @return updated API model
     */
    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        ShapeId serviceId = settings.getService();
        if (!SERVICE_TO_OPERATION_MAP.containsKey(serviceId)) {
            return model;
        }
        Model.Builder builder = model.toBuilder();

        Shape customString = StringShape.builder().id("sdk.customizations.presignURL#String").build();
        builder.addShape(customString);

        Set<ShapeId> operationIds = SERVICE_TO_OPERATION_MAP.get(serviceId);
        for (ShapeId operationId : operationIds) {
            OperationShape operation = model.expectShape(operationId, OperationShape.class);
            StructureShape input = ProtocolUtils.expectInput(model, operation);

            StructureShape.Builder inputBuilder = input.toBuilder();
            if (input.getAllMembers().containsKey("SourceRegion")) {
                // In the case of EC2 the SourceRegion member is expected to be serialized.
                LOGGER.warning("SourceRegion member is present in model and does not require backfill");
            } else {
                inputBuilder.addMember("SourceRegion", customString.getId(), (memberBuilder) -> {
                    memberBuilder
                            .addTrait(new DocumentationTrait(
                                    "The AWS region the resource is in. The presigned URL will be created with this region, " +
                                            "if the PresignURL member is empty set."))
                            .addTrait(new NoSerializeTrait());
                });
            }

            // Even if the input already contains DestinationRegion replace it with a unexported member.
            inputBuilder.addMember("DestinationRegion", customString.getId(), (memberBuilder) -> {
                memberBuilder
                        .addTrait(new DocumentationTrait(
                                "Used by the SDK's PresignURL autofill customization to specify the region " +
                                        "the of the client's request."))
                        .addTrait(new UnexportedMemberTrait());
            });

            // Add unmodeled parameters needed by presign customization
            builder.addShape(inputBuilder.build());
        }

        return builder.build();
    }

    /**
     * Generates additional types and logic for the presign customization middleware.
     * * Generates member getter and setter functions
     * * Generates API presign client
     * * Generates the connection between middleware, members, and api presign client.
     *
     * @param settings       codegen settings
     * @param model          api model
     * @param symbolProvider codegen symbol provider
     * @param goDelegator    writer provider
     */
    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ShapeId serviceId = settings.getService();
        if (!SERVICE_TO_OPERATION_MAP.containsKey(serviceId)) {
            return;
        }
        ServiceShape service = model.expectShape(serviceId, ServiceShape.class);
        Set<ShapeId> operationIds = SERVICE_TO_OPERATION_MAP.get(serviceId);

        for (ShapeId operationId : operationIds) {
            OperationShape operation = model.expectShape(operationId, OperationShape.class);
            StructureShape input = ProtocolUtils.expectInput(model, operation);

            goDelegator.useShapeWriter(operation, (writer) -> {
                // Need to copy input parameters for presign url.
                writeInputCopy(writer, symbolProvider, input);

                // Members used by the customization need abstract getter and setters
                writeMemberAccessor(writer, symbolProvider, operation, input);

                // Generate the presign client
                writePresignClientCustomization(writer, settings, model, symbolProvider, operation, input);

                // Generate the autofill presign client and its PresignURL method
                writeAutofillPresignClient(writer, symbolProvider, operation, input);
            });

            goDelegator.useShapeTestWriter(operation, (writer) -> {
                writePresignClientCustomizationTest(writer, symbolProvider, operation, input);
            });
        }
    }

    /**
     * Builds the set of runtime plugs used by the presign url customization.
     *
     * @param settings codegen settings
     * @param model    api model
     */
    @Override
    public void processFinalizedModel(GoSettings settings, Model model) {
        ShapeId serviceId = settings.getService();
        if (!SERVICE_TO_OPERATION_MAP.containsKey(serviceId)) {
            return;
        }
        ServiceShape service = settings.getService(model);
        for (ShapeId operationId : SERVICE_TO_OPERATION_MAP.get(serviceId)) {
            final OperationShape operation = model.expectShape(operationId, OperationShape.class);

            // Create a symbol provider because one is not available in this call.
            SymbolProvider symbolProvider = GoCodegenPlugin.createSymbolProvider(model, settings);
            String helperFuncName = addPresignMiddlewareFuncName(symbolProvider.toSymbol(operation).getName());

            runtimeClientPlugins.add(RuntimeClientPlugin.builder()
                    .servicePredicate((m, s) -> s.equals(service))
                    .operationPredicate((m, s, o) -> o.equals(operation))
                    .registerMiddleware(MiddlewareRegistrar.builder()
                            .resolvedFunction(SymbolUtils.createValueSymbolBuilder(helperFuncName)
                                    .build())
                            .useClientOptions()
                            .build())
                    .build());
        }
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return runtimeClientPlugins;
    }

    private void writeInputCopy(
            GoWriter writer,
            SymbolProvider symbolProvider,
            StructureShape input
    ) {
        Symbol inputSymbol = symbolProvider.toSymbol(input);
        writer.openBlock("func $L(params interface{}) (interface{}, error) {", "}",
                copyInputFuncName(inputSymbol.getName()),
                () -> {
                    writer.addUseImports(SmithyGoDependency.FMT);
                    writer.write("input, ok := params.($P)", inputSymbol);
                    writer.openBlock("if !ok {", "}", () -> {
                        writer.write("return nil, fmt.Errorf(\"expect $P type, got %T\", params)", inputSymbol);
                    });
                    writer.write("cpy := *input");
                    writer.write("return &cpy, nil");
                });
    }

    private void writeMemberAccessor(
            GoWriter writer,
            SymbolProvider symbolProvider,
            OperationShape operation,
            StructureShape input
    ) {
        // PresignedUrl member name has inconsistent casing across the services
        MemberShape presignURLMember = CodegenUtils.expectMember(input, "PresignedUrl"::equalsIgnoreCase);
        MemberShape dstRegionMember = CodegenUtils.expectMember(input, "DestinationRegion");
        MemberShape srcRegionMember = CodegenUtils.expectMember(input, "SourceRegion");

        MemberShape[] getMembers = {presignURLMember, srcRegionMember};
        for (MemberShape member : getMembers) {
            writeMemberGetter(writer, symbolProvider, operation, input, member);
        }

        MemberShape[] setMembers = {presignURLMember, dstRegionMember};
        for (MemberShape member : setMembers) {
            writeMemberSetter(writer, symbolProvider, operation, input, member);
        }
    }

    private void writePresignClientCustomization(
            GoWriter writer,
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            OperationShape operation,
            StructureShape input
    ) {
        Symbol smithyStack = SymbolUtils.createPointableSymbolBuilder("Stack",
                SmithyGoDependency.SMITHY_MIDDLEWARE).build();
        Symbol operationSymbol = symbolProvider.toSymbol(operation);
        Symbol inputSymbol = symbolProvider.toSymbol(input);

        String presignURLMember = symbolProvider.toMemberName(
                CodegenUtils.expectMember(input, "PresignedUrl"::equalsIgnoreCase));
        String dstRegionMember = symbolProvider.toMemberName(CodegenUtils.expectMember(input, "DestinationRegion"));
        String srcRegionMember = symbolProvider.toMemberName(CodegenUtils.expectMember(input, "SourceRegion"));

        Symbol parameterAccessor = SymbolUtils.createValueSymbolBuilder("ParameterAccessor",
                AwsCustomGoDependency.PRESIGNEDURL_CUSTOMIZATION).build();
        Symbol addMiddlewareOptions = SymbolUtils.createValueSymbolBuilder("Options",
                AwsCustomGoDependency.PRESIGNEDURL_CUSTOMIZATION).build();
        Symbol addMiddleware = SymbolUtils.createValueSymbolBuilder("AddMiddleware",
                AwsCustomGoDependency.PRESIGNEDURL_CUSTOMIZATION).build();


        // generate middleware mutator to wire up presign client with accessors and custom middleware.
        writer.openBlock("func $L(stack $P, options Options) error {", "}",
                addPresignMiddlewareFuncName(operationSymbol.getName()),
                smithyStack,
                () -> {
                    writer.openBlock("return $T(stack, $T{", "})", addMiddleware, addMiddlewareOptions, () -> {
                        writer.openBlock("Accessor: $T{", "},", parameterAccessor, () -> {
                            writer.write("GetPresignedURL: $L, \n",
                                    getterFuncName(operationSymbol.getName(), presignURLMember));
                            writer.write("GetSourceRegion: $L, \n",
                                    getterFuncName(operationSymbol.getName(), srcRegionMember));
                            writer.write("CopyInput: $L, \n", copyInputFuncName(inputSymbol.getName()));
                            writer.write("SetDestinationRegion: $L,\n",
                                    setterFuncName(operationSymbol.getName(), dstRegionMember));
                            writer.write("SetPresignedURL: $L, \n",
                                    setterFuncName(operationSymbol.getName(), presignURLMember));
                        });
                        writer.write("Presigner: &$L{ client: NewPresignClient(New(options))}, \n",
                                autofillPresignClient(operationSymbol.getName()));
                    });
                });
    }


    private void writeAutofillPresignClient(
            GoWriter writer,
            SymbolProvider symbolprovider,
            OperationShape operation,
            StructureShape input
    ) {
        Symbol operationSymbol = symbolprovider.toSymbol(operation);
        Symbol inputSymbol = symbolprovider.toSymbol(input);
        Symbol removeMiddleware = SymbolUtils.createValueSymbolBuilder("RemoveMiddleware",
                AwsCustomGoDependency.PRESIGNEDURL_CUSTOMIZATION).build();

        // generate autofill presign client
        writer.openBlock("type $L struct {", "}",
                autofillPresignClient(operationSymbol.getName()), () -> {
                    writer.write("client *PresignClient");
                });

        writer.write("");

        // generate PresignURL method that satisfies URLPresigner interface of middleware
        writer.writeDocs("PresignURL is a middleware accessor that satisfies URLPresigner interface.");
        writer.openBlock("func (c *$L) PresignURL(ctx context.Context, srcRegion string, params interface{}) "
                        + "(*v4.PresignedHTTPRequest, error) {", "}",

                autofillPresignClient(operationSymbol.getName()), () -> {
                    writer.addUseImports(SmithyGoDependency.FMT);
                    // check input
                    writer.write("input, ok := params.($P)", inputSymbol);
                    writer.openBlock("if !ok {", "}", () -> {
                        writer.write("return nil, fmt.Errorf(\"expect $P type, got %T\", params)", inputSymbol);
                    });

                    // generate client options
                    writer.openBlock("optFn := func(o *Options) {", "}", () -> {
                        writer.write("o.Region = srcRegion");
                        writer.write("o.APIOptions = append(o.APIOptions, $T)", removeMiddleware);
                    });

                    // getPresignAPIOptions
                    writer.write("presignOptFn := WithPresignClientFromClientOptions(optFn)");

                    // call the exported function
                    writer.write("return c.client.$L(ctx, input, presignOptFn)",
                            presignFuncName(operationSymbol.getName()));
                });
        writer.write("");
    }

    private void writePresignClientCustomizationTest(
            GoWriter writer,
            SymbolProvider symbolProvider,
            OperationShape operation,
            StructureShape input
    ) {
        String presignURLMember = symbolProvider.toMemberName(
                CodegenUtils.expectMember(input, "PresignedUrl"::equalsIgnoreCase));
        String dstRegionMember = symbolProvider.toMemberName(CodegenUtils.expectMember(input, "DestinationRegion"));
        String srcRegionMember = symbolProvider.toMemberName(CodegenUtils.expectMember(input, "SourceRegion"));

        Map<String, String> replaces = new TreeMap<>();
        replaces.put("operation", symbolProvider.toSymbol(operation).getName());
        replaces.put("input", symbolProvider.toSymbol(input).getName());
        replaces.put("presignURL_member", presignURLMember);
        replaces.put("srcRegion_member", srcRegionMember);
        replaces.put("dstRegion_member", dstRegionMember);

        String template = IoUtils.readUtf8Resource(getClass(), "presign_url.go.template");
        for (Map.Entry<String, String> entry : replaces.entrySet()) {
            template = template.replaceAll("\\$\\{" + entry.getKey() + "}", entry.getValue());
        }

        writer.addUseImports(SmithyGoDependency.NET_HTTP);
        writer.addUseImports(SmithyGoDependency.CONTEXT);
        writer.addUseImports(SmithyGoDependency.STRINGS);
        writer.addUseImports(SmithyGoDependency.TESTING);

        writer.addUseImports(SmithyGoDependency.SMITHY_HTTP_TRANSPORT);
        writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);

        writer.addUseImports(AwsGoDependency.AWS_CORE);
        writer.addUseImports(AwsGoDependency.AWSTESTING_UNIT);
        writer.addUseImports(AwsCustomGoDependency.PRESIGNEDURL_CUSTOMIZATION);

        writer.write(template);
    }

}
