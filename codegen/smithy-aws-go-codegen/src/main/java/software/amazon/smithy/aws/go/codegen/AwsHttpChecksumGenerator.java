package software.amazon.smithy.aws.go.codegen;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.Optional;
import java.util.Set;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.aws.traits.auth.UnsignedPayloadTrait;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoCodegenPlugin;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.aws.traits.HttpChecksumTrait;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.StreamingTrait;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.utils.SetUtils;

public class AwsHttpChecksumGenerator implements GoIntegration {
    // constant map with service to list of operation for which we should ignore multipart checksum validation.
    private static final Map<ShapeId, Set<ShapeId>> ignoreMultipartChecksumValidationMap = MapUtils.of(
            ShapeId.from("com.amazonaws.s3#AmazonS3"), SetUtils.of(
                    ShapeId.from("com.amazonaws.s3#GetObject")
            )
    );
    // list of runtime-client plugins
    private final List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();

    private static String getRequestAlgorithmAccessorFuncName(String operationName) {
        return String.format("get%s%s", operationName, "RequestAlgorithmMember");
    }

    private static String getRequestValidationModeAccessorFuncName(String operationName) {
        return String.format("get%s%s", operationName, "RequestValidationModeMember");
    }

    private static String getAddInputMiddlewareFuncName(String operationName) {
        return String.format("add%sInputChecksumMiddlewares", operationName);
    }

    private static String getAddOutputMiddlewareFuncName(String operationName) {
        return String.format("add%sOutputChecksumMiddlewares", operationName);
    }

    @Override
    public byte getOrder() {
        return 127;
    }

    /**
     * Builds the set of runtime plugs.
     *
     * @param settings codegen settings
     * @param model    api model
     */
    @Override
    public void processFinalizedModel(GoSettings settings, Model model) {
        ServiceShape service = settings.getService(model);
        for (ShapeId operationId : service.getAllOperations()) {
            final OperationShape operation = model.expectShape(operationId, OperationShape.class);

            // Create a symbol provider because one is not available in this call.
            SymbolProvider symbolProvider = GoCodegenPlugin.createSymbolProvider(model, settings);

            // Input helper
            String inputHelperFuncName = getAddInputMiddlewareFuncName(
                    symbolProvider.toSymbol(operation).getName()
            );
            runtimeClientPlugins.add(RuntimeClientPlugin.builder()
                    .operationPredicate((m, s, o) -> {
                        if (!hasInputChecksumTrait(m, s, o)) {
                            return false;
                        }
                        return o.equals(operation);
                    })
                    .registerMiddleware(MiddlewareRegistrar.builder()
                            .resolvedFunction(SymbolUtils.createValueSymbolBuilder(inputHelperFuncName)
                                    .build())
                            .useClientOptions()
                            .build())
                    .build());

            // Output helper
            String outputHelperFuncName = getAddOutputMiddlewareFuncName(
                    symbolProvider.toSymbol(operation).getName()
            );
            runtimeClientPlugins.add(RuntimeClientPlugin.builder()
                    .operationPredicate((m, s, o) -> {
                        if (!hasOutputChecksumTrait(m, s, o)) {
                            return false;
                        }
                        return o.equals(operation);
                    })
                    .registerMiddleware(MiddlewareRegistrar.builder()
                            .resolvedFunction(SymbolUtils.createValueSymbolBuilder(outputHelperFuncName)
                                    .build())
                            .useClientOptions()
                            .build())
                    .build());
        }
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceShape service = settings.getService(model);
        boolean supportsComputeInputChecksumsWorkflow = false;
        boolean supportsChecksumValidationWorkflow = false;

        for (ShapeId operationID : service.getAllOperations()) {
            OperationShape operation = model.expectShape(operationID, OperationShape.class);
            if (!hasChecksumTrait(model, service, operation)) {
                continue;
            }

            final boolean generateComputeInputChecksums = hasInputChecksumTrait(model, service, operation);
            if (generateComputeInputChecksums) {
                supportsComputeInputChecksumsWorkflow = true;
            }

            final boolean generateOutputChecksumValidation = hasOutputChecksumTrait(model, service, operation);
            if (generateOutputChecksumValidation) {
                supportsChecksumValidationWorkflow = true;
            }

            goDelegator.useShapeWriter(operation, writer -> {
                // generate getter helper function to access input member value
                writeGetInputMemberAccessorHelper(writer, model, symbolProvider, operation);

                // generate middleware helper function
                if (generateComputeInputChecksums) {
                    writeInputMiddlewareHelper(writer, model, symbolProvider, service, operation);
                }

                if (generateOutputChecksumValidation) {
                    writeOutputMiddlewareHelper(writer, model, symbolProvider, service, operation);
                }
            });
        }

        if (supportsComputeInputChecksumsWorkflow) {
            goDelegator.useShapeWriter(service, writer -> {
                generateInputComputedChecksumMetadataHelpers(writer, model, symbolProvider, service);
            });
        }

        if (supportsChecksumValidationWorkflow) {
            goDelegator.useShapeWriter(service, writer -> {
                generateOutputChecksumValidationMetadataHelpers(writer, model, symbolProvider, service);
            });
        }
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return runtimeClientPlugins;
    }

    // return true if operation shape is decorated with `httpChecksum` trait.
    private boolean hasChecksumTrait(Model model, ServiceShape service, OperationShape operation) {
        return operation.hasTrait(HttpChecksumTrait.class);
    }

    private boolean hasInputChecksumTrait(Model model, ServiceShape service, OperationShape operation) {
        if (!hasChecksumTrait(model, service, operation)) {
            return false;
        }
        HttpChecksumTrait trait = operation.expectTrait(HttpChecksumTrait.class);
        return trait.isRequestChecksumRequired() || trait.getRequestAlgorithmMember().isPresent();
    }

    private boolean hasOutputChecksumTrait(Model model, ServiceShape service, OperationShape operation) {
        if (!hasChecksumTrait(model, service, operation)) {
            return false;
        }
        HttpChecksumTrait trait = operation.expectTrait(HttpChecksumTrait.class);
        return trait.getRequestValidationModeMember().isPresent() && !trait.getResponseAlgorithms().isEmpty();
    }

    private boolean isS3ServiceShape(Model model, ServiceShape service) {
        String serviceId = service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("S3");
    }

    private void writeInputMiddlewareHelper(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape service,
            OperationShape operation
    ) {
        Symbol operationSymbol = symbolProvider.toSymbol(operation);
        String operationName = operationSymbol.getName();
        StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);

        HttpChecksumTrait trait = operation.expectTrait(HttpChecksumTrait.class);
        boolean isRequestChecksumRequired = trait.isRequestChecksumRequired();
        boolean hasRequestAlgorithmMember = trait.getRequestAlgorithmMember().isPresent();

        boolean supportsTrailingChecksum = false;
        for (MemberShape memberShape : input.getAllMembers().values()) {
            Shape targetShape = model.expectShape(memberShape.getTarget());
            if (targetShape.hasTrait(StreamingTrait.class) &&
                    !StreamingTrait.isEventStream(model, memberShape)
            ) {
                if (isS3ServiceShape(model, service) || (
                        AwsSignatureVersion4.hasSigV4AuthScheme(model, service, operation)
                                && !operation.hasTrait(UnsignedPayloadTrait.class))) {
                    supportsTrailingChecksum = true;
                }
            }
        }

        boolean supportsRequestTrailingChecksum = supportsTrailingChecksum;
        boolean supportsDecodedContentLengthHeader = isS3ServiceShape(model, service);

        // imports
        writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);

        writer.openBlock("func $L(stack *middleware.Stack, options Options) error {", "}",
                getAddInputMiddlewareFuncName(operationName), () -> {
                    writer.write("""
                                    return $T(stack, $T{
                                        GetAlgorithm: $L,
                                        RequireChecksum: $L,
                                        EnableTrailingChecksum: $L,
                                        EnableComputeSHA256PayloadHash: true,
                                        EnableDecodedContentLengthHeader: $L,
                                    })""",
                            SymbolUtils.createValueSymbolBuilder("AddInputMiddleware",
                                    AwsGoDependency.SERVICE_INTERNAL_CHECKSUM).build(),
                            SymbolUtils.createValueSymbolBuilder("InputMiddlewareOptions",
                                    AwsGoDependency.SERVICE_INTERNAL_CHECKSUM).build(),
                            hasRequestAlgorithmMember ?
                                    getRequestAlgorithmAccessorFuncName(operationName) : "nil",
                            isRequestChecksumRequired,
                            supportsRequestTrailingChecksum,
                            supportsDecodedContentLengthHeader);
                }
        );
        writer.insertTrailingNewline();
    }

    private void writeOutputMiddlewareHelper(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape service,
            OperationShape operation
    ) {
        Symbol operationSymbol = symbolProvider.toSymbol(operation);
        String operationName = operationSymbol.getName();
        StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);

        HttpChecksumTrait trait = operation.expectTrait(HttpChecksumTrait.class);
        List<String> responseAlgorithms = trait.getResponseAlgorithms();

        // imports
        writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);

        writer.openBlock("func $L(stack *middleware.Stack, options Options) error {", "}",
                getAddOutputMiddlewareFuncName(operationName), () -> {
                    writer.write("""
                                    return $T(stack, $T{
                                        GetValidationMode: $L,
                                        ValidationAlgorithms: $L,
                                        IgnoreMultipartValidation: $L,
                                        LogValidationSkipped: true,
                                        LogMultipartValidationSkipped: true,
                                    })""",
                            SymbolUtils.createValueSymbolBuilder("AddOutputMiddleware",
                                    AwsGoDependency.SERVICE_INTERNAL_CHECKSUM).build(),
                            SymbolUtils.createValueSymbolBuilder("OutputMiddlewareOptions",
                                    AwsGoDependency.SERVICE_INTERNAL_CHECKSUM).build(),

                            getRequestValidationModeAccessorFuncName(operationName),
                            convertToGoStringList(responseAlgorithms),
                            ignoreMultipartChecksumValidationMap.getOrDefault(
                                    service.toShapeId(), new HashSet<>()).contains(operation.toShapeId())
                    );
                });
        writer.insertTrailingNewline();
    }

    private String convertToGoStringList(List<String> list) {
        StringBuilder sb = new StringBuilder();
        sb.append("[]string{");
        for (String item : list) {
            sb.append("\"").append(item).append("\"");
            sb.append(",");
        }
        if (!list.isEmpty()) {
            sb.deleteCharAt(sb.length() - 1);
        }
        sb.append("}");
        return sb.toString();
    }

    private void writeGetInputMemberAccessorHelper(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            OperationShape operation
    ) {
        Symbol operationSymbol = symbolProvider.toSymbol(operation);
        StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);

        HttpChecksumTrait trait = operation.expectTrait(HttpChecksumTrait.class);

        // Input parameter for computing request payload's checksum.
        if (trait.getRequestAlgorithmMember().isPresent()) {
            Optional<MemberShape> memberShape = input.getAllMembers().values().stream()
                    .filter(m -> m.getMemberName().toLowerCase(Locale.ENGLISH)
                            .equals(trait.getRequestAlgorithmMember().get().toLowerCase(Locale.ENGLISH)))
                    .findFirst();
            if (!memberShape.isPresent()) {
                throw new CodegenException(String.format(
                        "Found no matching input member named %s modeled with HttpChecksum trait",
                        trait.getRequestAlgorithmMember().get()));
            }

            String memberName = memberShape.get().getMemberName();
            String funcName = getRequestAlgorithmAccessorFuncName(operationSymbol.getName());
            writer.writeDocs(
                    String.format("%s gets the request checksum algorithm value provided as input.", funcName));
            getInputTemplate(writer, symbolProvider, input, funcName, memberName);
            writer.insertTrailingNewline();
        }

        // Output parameter for validating response payload's checksum
        if (trait.getRequestValidationModeMember().isPresent()) {
            Optional<MemberShape> memberShape = input.getAllMembers().values().stream()
                    .filter(m -> m.getMemberName().toLowerCase(Locale.ENGLISH)
                            .equals(trait.getRequestValidationModeMember().get().toLowerCase(Locale.ENGLISH)))
                    .findFirst();
            if (!memberShape.isPresent()) {
                throw new CodegenException(String.format(
                        "Found no matching input member named %s modeled with HttpChecksum trait",
                        trait.getRequestValidationModeMember().get()));
            }

            String memberName = memberShape.get().getMemberName();
            String funcName = getRequestValidationModeAccessorFuncName(operationSymbol.getName());
            writer.writeDocs(
                    String.format("%s gets the request checksum validation mode provided as input.", funcName));
            getInputTemplate(writer, symbolProvider, input, funcName, memberName);
            writer.insertTrailingNewline();
        }
    }

    private void getInputTemplate(
            GoWriter writer,
            SymbolProvider symbolProvider,
            StructureShape input,
            String funcName,
            String memberName
    ) {
        writer.openBlock("func $L(input interface{}) (string, bool) {", "}", funcName,
                () -> {
                    writer.write("in := input.($P)", symbolProvider.toSymbol(input));
                    writer.openBlock("if len(in.$L) == 0 {", "}", memberName, () -> {
                        writer.write("return \"\", false");
                    });
                    writer.write("return string(in.$L), true", memberName);
                });
        writer.write("");
    }

    private void generateInputComputedChecksumMetadataHelpers(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape service
    ) {
        String metadataStructName = "ComputedInputChecksumsMetadata";

        writer.writeDocs(String.format("""
                %s provides information about the algorithms used to compute the checksum(s) of the
                input payload.
                """, metadataStructName));
        writer.openBlock("type $L struct {", "}", metadataStructName, () -> {
            writer.writeDocs("""
                    ComputedChecksums is a map of algorithm name to checksum value of the computed
                    input payload's checksums.
                    """);
            writer.write("ComputedChecksums map[string]string");
        });

        Symbol metadataStructSymbol =
                SymbolUtils.createValueSymbolBuilder(metadataStructName).build();
        String metadataGetterFuncName = "Get" + metadataStructName;
        Symbol getAlgorithmUsed = SymbolUtils.createValueSymbolBuilder(
                "GetComputedInputChecksums", AwsGoDependency.SERVICE_INTERNAL_CHECKSUM).build();

        writer.writeDocs(String.format("""
                %s retrieves from the result metadata the map of algorithms and input payload checksums values.
                """, metadataGetterFuncName));
        writer.openBlock("func $L(m $T) ($T, bool) {", "}",
                metadataGetterFuncName,
                SymbolUtils.createValueSymbolBuilder("Metadata", SmithyGoDependency.SMITHY_MIDDLEWARE).build(),
                metadataStructSymbol,
                () -> {
                    writer.write("values, ok := $T(m)", getAlgorithmUsed);
                    writer.write("if !ok { return $T{}, false }", metadataStructSymbol);
                    writer.write("""
                            return $T{
                                ComputedChecksums: values,
                            }, true
                            """, metadataStructSymbol);
                });
        writer.write("");
    }

    private void generateOutputChecksumValidationMetadataHelpers(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape service
    ) {
        String metadataStructName = "ChecksumValidationMetadata";

        writer.writeDocs(String.format(
                "%s contains metadata such as the checksum algorithm used for data integrity validation.",
                metadataStructName));
        writer.openBlock("type $L struct {", "}", metadataStructName, () -> {
            writer.writeDocs("""
                    AlgorithmsUsed is the set of the checksum algorithms used to validate the response payload.
                    The response payload must be completely read in order for the checksum validation to be
                    performed. An error is returned by the operation output's response io.ReadCloser
                    if the computed checksums are invalid.
                    """);
            writer.write("AlgorithmsUsed []string");
        });

        Symbol metadataStructSymbol =
                SymbolUtils.createValueSymbolBuilder(metadataStructName).build();
        String metadataGetterFuncName = "Get" + metadataStructName;
        Symbol getAlgorithmUsed = SymbolUtils.createValueSymbolBuilder(
                "GetOutputValidationAlgorithmsUsed", AwsGoDependency.SERVICE_INTERNAL_CHECKSUM).build();

        writer.writeDocs(String.format("""
                %s returns the set of algorithms that will be used to validate the response payload with. The
                response payload must be completely read in order for the checksum validation to be performed.
                An error is returned by the operation output's response io.ReadCloser if the computed checksums
                are invalid. Returns false if no checksum algorithm used metadata was found.
                """, metadataGetterFuncName));
        writer.openBlock("func $L(m $T) ($T, bool) {", "}", metadataGetterFuncName,
                SymbolUtils.createValueSymbolBuilder("Metadata", SmithyGoDependency.SMITHY_MIDDLEWARE).build(),
                metadataStructSymbol,
                () -> {
                    writer.write("values, ok := $T(m)", getAlgorithmUsed);
                    writer.write("if !ok { return $T{}, false }", metadataStructSymbol);
                    writer.write("""
                            return $T{
                                AlgorithmsUsed: append(make([]string, 0, len(values)), values...),
                            }, true
                            """, metadataStructSymbol);
                });
        writer.write("");
    }
}
