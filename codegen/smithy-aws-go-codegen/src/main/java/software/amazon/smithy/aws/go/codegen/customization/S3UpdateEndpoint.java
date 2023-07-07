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
import java.util.Set;
import java.util.TreeSet;
import java.util.stream.Collectors;
import software.amazon.smithy.aws.go.codegen.AwsSignatureVersion4aUtils;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoCodegenPlugin;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.SetUtils;

/**
 * S3UpdateEndpoint integration serves to apply customizations for S3 service,
 * and modifies the resolved endpoint based on S3 client config or input shape values.
 */
public class S3UpdateEndpoint implements GoIntegration {
    // Middleware name
    private static final String UPDATE_ENDPOINT_INTERNAL_ADDER = "UpdateEndpoint";
    // Middleware options
    private static final String UPDATE_ENDPOINT_INTERNAL_OPTIONS =
            UPDATE_ENDPOINT_INTERNAL_ADDER + "Options";
    private static final String UPDATE_ENDPOINT_INTERNAL_PARAMETER_ACCESSOR =
            UPDATE_ENDPOINT_INTERNAL_ADDER + "ParameterAccessor";

    // s3 shared option constants
    private static final String USE_DUALSTACK_OPTION = "UseDualstack";
    private static final String USE_ARNREGION_OPTION = "UseARNRegion";

    // list of runtime-client plugins
    private final List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();

    private static boolean isS3SharedService(Model model, ServiceShape service) {
        return S3ModelUtils.isServiceS3(model, service) || S3ModelUtils.isServiceS3Control(model, service);
    }

    private static String copyInputFuncName(String inputName) {
        return String.format("copy%sForUpdateEndpoint", inputName);
    }

    private static String getterFuncName(String operationName, String memberName) {
        return String.format("get%s%s", operationName, memberName);
    }

    private static String setterFuncName(String operationName, String memberName) {
        return String.format("set%s%s", operationName, memberName);
    }

    private static String addMiddlewareFuncName(String operationname, String middlewareName) {
        return String.format("add%s%s", operationname, middlewareName);
    }

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first.
     *
     * @return Returns the sort order, defaults to -40.
     */
    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceShape service = settings.getService(model);

        // if service is s3control
        if (S3ModelUtils.isServiceS3Control(model, service)) {
            S3control s3control = new S3control(service);
            s3control.writeAdditionalFiles(settings, model, symbolProvider, goDelegator);
        }

        // check if service is s3
        if (S3ModelUtils.isServiceS3(model, service)) {
            S3 s3 = new S3(service);
            s3.writeAdditionalFiles(settings, model, symbolProvider, goDelegator);
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
        ServiceShape service = settings.getService(model);
        for (final OperationShape operation : TopDownIndex.of(model).getContainedOperations(service)) {
            // Create a symbol provider because one is not available in this call.
            SymbolProvider symbolProvider = GoCodegenPlugin.createSymbolProvider(model, settings);
            String helperFuncName = addMiddlewareFuncName(
                    symbolProvider.toSymbol(operation).getName(),
                    UPDATE_ENDPOINT_INTERNAL_ADDER
            );
        }
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        runtimeClientPlugins.addAll(ListUtils.of(
                // Add S3 shared config's dualstack option
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3UpdateEndpoint::isS3SharedService)
                        .configFields(ListUtils.of(
                                ConfigField.builder()
                                        .name(USE_DUALSTACK_OPTION)
                                        .type(SymbolUtils.createValueSymbolBuilder("bool")
                                                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                                .build())
                                        .documentation("Allows you to enable dual-stack endpoint support for the "
                                                       + "service.")
                                        .deprecated("""
                                                    Set dual-stack by setting UseDualStackEndpoint on
                                                    EndpointResolverOptions. When EndpointResolverOptions'
                                                    UseDualStackEndpoint field is set it overrides this field value.""")
                                        .build(),
                                ConfigField.builder()
                                        .name(USE_ARNREGION_OPTION)
                                        .type(SymbolUtils.createValueSymbolBuilder("bool")
                                                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                                .build())
                                        .documentation("Allows you to enable arn region support for the service.")
                                        .build()
                        ))
                        .build()
        ));
        runtimeClientPlugins.addAll(S3.getClientPlugins());
        return runtimeClientPlugins;
    }

    /*
     * s3 class is the private class handling s3 goIntegration for endpoint mutations
     */
    private static class S3 {
        // options to be generated on Client's options type
        private static final String USE_PATH_STYLE_OPTION = "UsePathStyle";
        private static final String USE_ACCELERATE_OPTION = "UseAccelerate";
        private static final String DISABLE_MRAP_OPTION = "DisableMultiRegionAccessPoints";

        // private function getter constant
        private static final String NOP_BUCKET_ACCESSOR = "nopGetBucketAccessor";
        // service shape representing s3
        private final ServiceShape service;
        // list of operations that take in bucket as input
        private final Set<String> SUPPORT_BUCKET_AS_INPUT = new TreeSet<>();
        // list of operations that do not support accelerate
        private final Set<String> NOT_SUPPORT_ACCELERATE = SetUtils.of(
                "ListBuckets", "CreateBucket", "DeleteBucket"
        );

        private final Set<String> TARGET_OBJECT_LAMBDAS = SetUtils.of("WriteGetObjectResponse");

        private S3(ServiceShape service) {
            this.service = service;
        }

        // getClientPlugins returns a list of client plugins for s3 service
        private static List<RuntimeClientPlugin> getClientPlugins() {
            List<RuntimeClientPlugin> list = ListUtils.of(
                    // Add S3 config to use path style host addressing.
                    RuntimeClientPlugin.builder()
                            .servicePredicate((model, service1) -> S3ModelUtils.isServiceS3(model, service1))
                            .configFields(ListUtils.of(
                                    ConfigField.builder()
                                            .name(USE_PATH_STYLE_OPTION)
                                            .type(SymbolUtils.createValueSymbolBuilder("bool")
                                                    .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                                    .build())
                                            .documentation(
                                                    "Allows you to enable the client to use path-style addressing, "
                                                    + "i.e., `https://s3.amazonaws.com/BUCKET/KEY`. By default, the S3 client "
                                                    + "will use virtual hosted bucket addressing when possible"
                                                    + "(`https://BUCKET.s3.amazonaws.com/KEY`).")
                                            .build(),
                                    ConfigField.builder()
                                            .name(USE_ACCELERATE_OPTION)
                                            .type(SymbolUtils.createValueSymbolBuilder("bool")
                                                    .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                                    .build())
                                            .documentation("Allows you to enable S3 Accelerate feature. All operations "
                                                           + "compatible with S3 Accelerate will use the accelerate endpoint for "
                                                           + "requests. Requests not compatible will fall back to normal S3 requests. "
                                                           + "The bucket must be enabled for accelerate to be used with S3 client with "
                                                           + "accelerate enabled. If the bucket is not enabled for accelerate an error "
                                                           + "will be returned. The bucket name must be DNS compatible to work "
                                                           + "with accelerate.")
                                            .build(),
                                    ConfigField.builder()
                                            .name(DISABLE_MRAP_OPTION)
                                            .type(SymbolUtils.createValueSymbolBuilder("bool")
                                                    .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                                    .build())
                                            .documentation("Allows you to disable S3 Multi-Region access points feature.")
                                            .build(),
                                    ConfigField.builder()
                                            .name(AwsSignatureVersion4aUtils.V4A_SIGNER_INTERFACE_NAME)
                                            .type(SymbolUtils.createValueSymbolBuilder(
                                                            AwsSignatureVersion4aUtils.V4A_SIGNER_INTERFACE_NAME)
                                                    .build())
                                            .documentation("Signature Version 4a (SigV4a) Signer")
                                            .build()
                            ))
                            .build()
            );
            return list;
        }

        // retrieves function name for get bucket accessor function
        private String getBucketAccessorFuncName(String operationName) {
            return getterFuncName(operationName, "BucketMember");
        }

        private void writeAdditionalFiles(
                GoSettings settings,
                Model model,
                SymbolProvider symbolProvider,
                GoDelegator goDelegator
        ) {

            for (OperationShape operation : TopDownIndex.of(model).getContainedOperations(service)) {
                goDelegator.useShapeWriter(operation, writer -> {
                    // generate get bucket member helper function
                    writeGetBucketMemberHelper(writer, model, symbolProvider, operation);
                    // generate update endpoint middleware helper function
                    writeMiddlewareHelper(writer, model, symbolProvider, operation);
                });
            }

            goDelegator.useShapeWriter(service, writer -> {
                // generate NOP bucket accessor helper
                writeNOPBucketAccessorHelper(writer);
            });
        }

        private void writeMiddlewareHelper(
                GoWriter writer,
                Model model,
                SymbolProvider symbolProvider,
                OperationShape operationShape
        ) {
            // imports
            writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);

            // operation name
            String operationName = symbolProvider.toSymbol(operationShape).getName();

            writer.openBlock("func $L(stack *middleware.Stack, options Options) error {", "}",
                    addMiddlewareFuncName(symbolProvider.toSymbol(operationShape).getName(),
                            UPDATE_ENDPOINT_INTERNAL_ADDER), () -> {
                        writer.write("return $T(stack, $T{ \n"
                                     + "Accessor : $T{\n"
                                     + "GetBucketFromInput: $L,\n},\n"
                                     + "UsePathStyle: options.$L,\n"
                                     + "UseAccelerate: options.$L,\n"
                                     + "SupportsAccelerate: $L,\n"
                                     + "TargetS3ObjectLambda: $L,\n"
                                     + "EndpointResolver: options.EndpointResolver,\n"
                                     + "EndpointResolverOptions: options.EndpointOptions,\n"
                                     + "UseARNRegion: options.$L,\n"
                                     + "DisableMultiRegionAccessPoints: options.$L,\n"
                                     + "})",
                                SymbolUtils.createValueSymbolBuilder(UPDATE_ENDPOINT_INTERNAL_ADDER,
                                        AwsCustomGoDependency.S3_CUSTOMIZATION).build(),
                                SymbolUtils.createValueSymbolBuilder(UPDATE_ENDPOINT_INTERNAL_OPTIONS,
                                        AwsCustomGoDependency.S3_CUSTOMIZATION).build(),
                                SymbolUtils.createValueSymbolBuilder(UPDATE_ENDPOINT_INTERNAL_PARAMETER_ACCESSOR,
                                        AwsCustomGoDependency.S3_CUSTOMIZATION).build(),
                                SUPPORT_BUCKET_AS_INPUT.contains(operationName) ?
                                        getBucketAccessorFuncName(operationName) : NOP_BUCKET_ACCESSOR,
                                USE_PATH_STYLE_OPTION,
                                USE_ACCELERATE_OPTION,
                                !NOT_SUPPORT_ACCELERATE.contains(operationName),
                                TARGET_OBJECT_LAMBDAS.contains(operationName),
                                USE_ARNREGION_OPTION,
                                DISABLE_MRAP_OPTION
                        );
                    });
            writer.insertTrailingNewline();
        }

        private void writeNOPBucketAccessorHelper(
                GoWriter writer
        ) {
            writer.writeDocs(
                    String.format("%s is no-op accessor for operation that don't support bucket member as input",
                            NOP_BUCKET_ACCESSOR)
            );
            writer.openBlock("func $L(input interface{}) (*string, bool) {", "}", NOP_BUCKET_ACCESSOR,
                    () -> {
                        writer.write("return nil, false");
                    });
            writer.insertTrailingNewline();
        }

        private void writeGetBucketMemberHelper(
                GoWriter writer,
                Model model,
                SymbolProvider symbolProvider,
                OperationShape operation
        ) {
            Symbol operationSymbol = symbolProvider.toSymbol(operation);
            String operationName = operationSymbol.getName();
            String funcName = getBucketAccessorFuncName(operationSymbol.getName());

            StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);

            List<MemberShape> targetBucketShape = input.getAllMembers().values().stream()
                    .filter(m -> m.getTarget().getName(service).equals("BucketName"))
                    .collect(Collectors.toList());
            // if model has multiple top level shapes targeting `BucketName`, we throw a codegen exception
            if (targetBucketShape.size() > 1) {
                throw new CodegenException(
                        "BucketName shape should be targeted by only one input member, found " +
                        targetBucketShape.size() + " for Input shape: " + input.getId());
            }

            if (targetBucketShape.isEmpty()) {
                return;
            }

            // add operation to set denoting operation supports bucket input member
            SUPPORT_BUCKET_AS_INPUT.add(operationName);

            writer.writeDocs(
                    String.format("%s returns a pointer to string denoting a provided bucket member value"
                                  + "and a boolean indicating if the input has a modeled bucket name,", funcName)
            );
            writer.openBlock("func $L(input interface{}) (*string, bool) {", "}", funcName,
                    () -> {
                        String targetShapeName = targetBucketShape.get(0).getMemberName();
                        writer.write("in := input.($P)", symbolProvider.toSymbol(input));
                        writer.openBlock("if in.$L == nil {", "}", targetShapeName, () -> {
                            writer.write("return nil, false");
                        });
                        writer.write("return in.$L, true", targetShapeName);
                    });
            writer.insertTrailingNewline();
        }
    }

    /**
     * s3control class is the private class handling s3control goIntegration for endpoint mutations
     */
    private static class S3control {

        // nop accessor constants
        private static final String NOP_GET_ARN_ACCESSOR = "nopGetARNAccessor";
        private static final String NOP_SET_ARN_ACCESSOR = "nopSetARNAccessor";
        private static final String NOP_BACKFILL_ACCOUNT_ID_HELPER = "nopBackfillAccountIDAccessor";
        private static final String NOP_GET_OUTPOST_ID_FROM_INPUT = "nopGetOutpostIDFromInput";
        // Map of service, list of operationName that support ARNs as input
        private static final Set<String> supportsARN = new TreeSet<>();
        // service associated with this class
        private final ServiceShape service;
        // List of operations that use Accesspoint field as ARN input source.
        private final Set<String> LIST_ACCESSPOINT_ARN_INPUT = SetUtils.of(
                "GetAccessPoint", "DeleteAccessPoint", "PutAccessPointPolicy",
                "GetAccessPointPolicy", "DeleteAccessPointPolicy"
        );
        // List of operations that use OutpostID to resolve endpoint
        private final Set<String> LIST_OUTPOST_ID_INPUT = SetUtils.of(
                "CreateBucket", "ListRegionalBuckets"
        );

        private S3control(ServiceShape service) {
            this.service = service;
        }

        // returns a function identifier string for backfillAccountID function
        private static final String backFillAccountIDFuncName(String operation) {
            return String.format("backFill%s%s", operation, "AccountID");
        }

        // returns a function identifier string for arn member setter function
        private static final String setARNMemberFuncName(String operation) {
            return setterFuncName(operation, "ARNMember");
        }

        // returns a function identifier string for arn member getter function
        private static final String getARNMemberFuncName(String operation) {
            return getterFuncName(operation, "ARNMember");
        }

        // returns a function identifier string for outpost id member getter function
        private static final String getOutpostIDMemberFuncName(String operation) {
            return getterFuncName(operation, "OutpostIDMember");
        }

        void writeAdditionalFiles(
                GoSettings settings,
                Model model,
                SymbolProvider symbolProvider,
                GoDelegator goDelegator
        ) {
            for (OperationShape operation : TopDownIndex.of(model).getContainedOperations(service)) {
                goDelegator.useShapeWriter(operation, writer -> {
                    // get input shape from operation
                    StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);
                    // generate input copy function
                    writeInputCopy(writer, symbolProvider, input);
                    // generate outpost id accessor function
                    writeOutpostIDHelper(writer, model, symbolProvider, operation);
                    // generate arn helper function
                    writeARNHelper(writer, model, symbolProvider, operation);
                    // generate backfill account id helper function
                    writeBackfillAccountIDHelper(writer, model, symbolProvider, operation);
                    // generate update endpoint middleware helper function
                    writeMiddlewareHelper(writer, model, symbolProvider, operation);
                });
            }

            // services must be processed later than operations to help generation of Nop helpers
            goDelegator.useShapeWriter(service, writer -> {
                // generate outpost id helper function
                writeNopOutpostIDHelper(writer);
                // generate Nop methods for ARNHelper
                writeNopARNHelper(writer);
                // generate Nop methods for BackfillAccountIDHelper
                writeNopBackfillAccountIDHelper(writer);
            });
        }

        private void writeNopARNHelper(
                GoWriter writer
        ) {
            // generate get arn member accessor getter function
            writer.writeDocs("nopGetARNAccessor provides a nop get accessor function to be used "
                             + "when a certain operation does not support ARNs");
            writer.openBlock("func $L (input interface{}) (*string, bool) { ", "}",
                    NOP_GET_ARN_ACCESSOR, () -> {
                        writer.write("return nil, false");
                    });
            writer.insertTrailingNewline();

            // generate set arn member accessor setter function
            writer.writeDocs("nopSetARNAccessor provides a nop set accessor function to be used "
                             + "when a certain operation does not support ARNs");
            writer.openBlock("func $L (input interface{}, v string) error {", "}",
                    NOP_SET_ARN_ACCESSOR, () -> {
                        writer.write("return nil");
                    });
            writer.insertTrailingNewline();
        }

        private void writeNopBackfillAccountIDHelper(
                GoWriter writer
        ) {
            // generate arn member accessor getter function
            writer.writeDocs("nopBackfillAccountIDAccessor provides a nop accessor function to be used "
                             + "when a certain operation does not need to validate and backfill account id");
            writer.openBlock("func $L (input interface{}, v string) error {", "}",
                    NOP_BACKFILL_ACCOUNT_ID_HELPER, () -> {
                        writer.write("return nil");
                    });
            writer.insertTrailingNewline();
        }

        private void writeNopOutpostIDHelper(
                GoWriter writer
        ) {
            writer.writeDocs("nopGetOutpostIDFromInput provides a nop accessor function to be used "
                             + "when endpoint customization behavior is not based on presence of outpost id member if any");
            writer.openBlock("func $L (input interface{}) (*string, bool) {", "}",
                    NOP_GET_OUTPOST_ID_FROM_INPUT, () -> {
                        writer.write("return nil, false");
                    });
            writer.insertTrailingNewline();
        }

        private void writeMiddlewareHelper(
                GoWriter writer, Model model, SymbolProvider symbolProvider, OperationShape operationShape
        ) {
            // imports
            writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);

            // input shape
            StructureShape inputShape = model.expectShape(operationShape.getInput().get(), StructureShape.class);
            String operationName = symbolProvider.toSymbol(operationShape).getName();

            writer.openBlock("func $L(stack *middleware.Stack, options Options) error {", "}",
                    addMiddlewareFuncName(symbolProvider.toSymbol(operationShape).getName(),
                            UPDATE_ENDPOINT_INTERNAL_ADDER), () -> {
                        writer.write("return $T(stack, $T{ \n"
                                     + "Accessor : $T{GetARNInput: $L,\n BackfillAccountID: $L,\n"
                                     + "GetOutpostIDInput: $L, \n UpdateARNField: $L,\n CopyInput: $L,\n }, \n"
                                     + "EndpointResolver: options.EndpointResolver,\n "
                                     + "EndpointResolverOptions: options.EndpointOptions,\n"
                                     + "UseARNRegion: options.$L, \n })",
                                SymbolUtils.createValueSymbolBuilder(UPDATE_ENDPOINT_INTERNAL_ADDER,
                                        AwsCustomGoDependency.S3CONTROL_CUSTOMIZATION).build(),
                                SymbolUtils.createValueSymbolBuilder(UPDATE_ENDPOINT_INTERNAL_OPTIONS,
                                        AwsCustomGoDependency.S3CONTROL_CUSTOMIZATION).build(),
                                SymbolUtils.createValueSymbolBuilder(UPDATE_ENDPOINT_INTERNAL_PARAMETER_ACCESSOR,
                                        AwsCustomGoDependency.S3CONTROL_CUSTOMIZATION).build(),
                                supportsARN.contains(operationName) ? getARNMemberFuncName(
                                        operationName) : NOP_GET_ARN_ACCESSOR,
                                supportsARN.contains(operationName) ? backFillAccountIDFuncName(
                                        operationName) : NOP_BACKFILL_ACCOUNT_ID_HELPER,
                                LIST_OUTPOST_ID_INPUT.contains(operationName) ? getOutpostIDMemberFuncName(
                                        operationName) : NOP_GET_OUTPOST_ID_FROM_INPUT,
                                supportsARN.contains(operationName) ? setARNMemberFuncName(
                                        operationName) : NOP_SET_ARN_ACCESSOR,
                                copyInputFuncName(symbolProvider.toSymbol(inputShape).getName()),
                                USE_ARNREGION_OPTION
                        );
                    });
            writer.insertTrailingNewline();
        }

        /**
         * Writes a accessor function that returns an address to copy of passed in input
         *
         * @param writer
         * @param symbolProvider
         * @param input
         */
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

        /**
         * writes BackfillAccountID Helper function for s3 api operation
         * <p>
         * Generates code:
         * === api_operation.go===
         * func backfillAccountID(input interface{}, v string) error {
         * in := input.(*OpInputType)
         * if in.AccountId!=nil {
         * iv := *in.AccountId
         * if !strings.EqualFold(iv, v) {
         * return fmt.Errorf("error backfilling account id")
         * }
         * return nil
         * }
         * <p>
         * in.AccountId = &v
         * return nil
         * }
         */
        private void writeBackfillAccountIDHelper(
                GoWriter writer, Model model, SymbolProvider symbolProvider, OperationShape operation
        ) {
            StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);
            List<MemberShape> targetAccountIDShape = input.getAllMembers().values().stream()
                    .filter(m -> m.getMemberName().equals("AccountId"))
                    .collect(Collectors.toList());
            // if model has multiple top level shapes targeting `AccountId`, we throw a codegen exception
            if (targetAccountIDShape.size() > 1) {
                throw new CodegenException("AccountId shape should be targeted by only one input member, found " +
                                           targetAccountIDShape.size() + " for Input shape: " + input.getId());
            }

            if (targetAccountIDShape.isEmpty()) {
                return;
            }

            Symbol inputSymbol = symbolProvider.toSymbol(input);
            writer.write("func $L (input interface{}, v string) error { ",
                    backFillAccountIDFuncName(symbolProvider.toSymbol(operation).getName()));
            String memberName = targetAccountIDShape.get(0).getMemberName();
            writer.write("in := input.($P)", inputSymbol);
            writer.write("if in.$L != nil {", memberName);

            writer.addUseImports(SmithyGoDependency.STRINGS);
            writer.write("if !strings.EqualFold(*in.$L, v) {", memberName);

            writer.addUseImports(SmithyGoDependency.FMT);
            writer.write("return fmt.Errorf(\"error backfilling account id\") }");
            writer.write("return nil }");
            writer.write("in.$L = &v", memberName);
            writer.write("return nil }");

            writer.insertTrailingNewline();
        }

        /**
         * writes getARNMemberValue and updateARNMemberValue update function for all api input operations
         */
        private void writeARNHelper(
                GoWriter writer, Model model, SymbolProvider symbolProvider, OperationShape operation
        ) {
            // list of outpost id input require special behavior
            if (LIST_OUTPOST_ID_INPUT.contains(operation.getId().getName(service))) {
                return;
            }

            // arn target shape
            String arnType = LIST_ACCESSPOINT_ARN_INPUT.contains(
                    operation.getId().getName(service)
            ) ? "AccessPointName" : "BucketName";

            StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);
            List<MemberShape> listOfARNMembers = input.getAllMembers().values().stream()
                    .filter(m -> m.getTarget().getName(service).equals(arnType))
                    .collect(Collectors.toList());
            // if model has multiple top level shapes targeting arnable field, we throw a codegen exception
            if (listOfARNMembers.size() > 1) {
                throw new CodegenException(arnType + " shape should be targeted by only one input member, found " +
                                           listOfARNMembers.size() + " for Input shape: " + input.getId());
            }

            if (listOfARNMembers.isEmpty()) {
                return;
            }

            String operationName = symbolProvider.toSymbol(operation).getName();
            // this operation supports taking arn as input
            supportsARN.add(operationName);

            Symbol inputSymbol = symbolProvider.toSymbol(input);
            String memberName = listOfARNMembers.get(0).getMemberName();

            // generate arn member accessor getter function
            writer.write("func $L (input interface{}) (*string, bool) {",
                    getARNMemberFuncName(symbolProvider.toSymbol(operation).getName()));
            writer.write("in := input.($P)", inputSymbol);
            writer.write("if in.$L == nil {return nil, false }", memberName);
            writer.write("return in.$L, true }", memberName);

            writer.insertTrailingNewline();

            // generate arn member accessor setter function
            writer.write("func $L (input interface{}, v string) error {",
                    setARNMemberFuncName(symbolProvider.toSymbol(operation).getName()));
            writer.write("in := input.($P)", inputSymbol);
            writer.write("in.$L = &v", memberName);
            writer.write("return nil }");

            writer.insertTrailingNewline();
        }

        /**
         * writes OutpostID Helper function for operations CreateBucket and ListRegionalBuckets
         * <p>
         * Generates code:
         * func get<OpName>OutpostIDHelper</> (in interface{}) (*string, bool) {
         * i, ok := input.(*OpName)
         * if !ok {
         * return nil, fmt.Errorf("Expected input of type *OpName, got %T", input)
         * }
         * return i.<MemberName>, nil
         * }
         */
        private void writeOutpostIDHelper(
                GoWriter writer,
                Model model,
                SymbolProvider symbolProvider,
                OperationShape operation
        ) {
            String operationName = symbolProvider.toSymbol(operation).getName();
            if (!LIST_OUTPOST_ID_INPUT.contains(operationName)) {
                return;
            }

            String funcName = getOutpostIDMemberFuncName(operationName);

            writer.writeDocs(
                    String.format("%s returns a pointer to string denoting a provided outpost-id member value"
                                  + " and a boolean indicating if the input has a modeled outpost-id,", funcName));
            writer.openBlock("func $L (input interface{}) (*string, bool) {", "}",
                    funcName, () -> {
                        StructureShape input = model.expectShape(operation.getInput().get(),
                                StructureShape.class);
                        List<MemberShape> outpostIDMemberShapes = input.getAllMembers().values().stream()
                                .filter(m -> m.getMemberName().equalsIgnoreCase("OutpostId"))
                                .collect(Collectors.toList());
                        // if model has multiple top level shapes targeting `OutpostId`, we throw a codegen exception
                        if (outpostIDMemberShapes.size() > 1) {
                            throw new CodegenException(
                                    "OutpostID shape should be targeted by only one input member, found " +
                                    outpostIDMemberShapes.size() + " for Input shape: " + input.getId());
                        }

                        if (outpostIDMemberShapes.isEmpty()) {
                            LIST_OUTPOST_ID_INPUT.remove(operationName);
                        }

                        Symbol inputSymbol = symbolProvider.toSymbol(input);
                        String memberName = outpostIDMemberShapes.get(0).getMemberName();
                        writer.write("in := input.($P)", inputSymbol);
                        writer.openBlock("if in.$L == nil  {", "}", memberName, () -> {
                            writer.write("return nil, false");
                        });
                        writer.write("return in.$L, true", memberName);
                    });
        }
    }
}
