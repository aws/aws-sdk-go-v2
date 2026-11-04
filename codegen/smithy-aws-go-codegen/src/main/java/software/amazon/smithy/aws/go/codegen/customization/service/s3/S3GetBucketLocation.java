package software.amazon.smithy.aws.go.codegen.customization.service.s3;

import java.util.List;
import java.util.Map;
import software.amazon.smithy.aws.go.codegen.XmlProtocolUtils;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.Writable;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.middleware.DeserializeStepMiddleware;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.utils.ListUtils;

/**
 * This integration generates a custom deserializer for GetBucketLocation response.
 * Amazon S3 service does not wrap the GetBucketLocation response with Operation
 * name xml tags, and thus custom deserialization is required.
 * <p>
 * Related to aws/aws-sdk-go-v2#908
 */
public class S3GetBucketLocation implements GoIntegration {

    private final String protocolName = "awsRestxml";
    private final String swapDeserializerFuncName = "swapDeserializerHelper";
    private final String getBucketLocationOpID = "GetBucketLocation";
    private static final String UNWRAPPED_DESERIALIZER_NAME = "getBucketLocationDeserializer";

    /**
     * Return true if service is Amazon S3.
     *
     * @param model   is the generation model.
     * @param service is the service shape being audited.
     */
    private static boolean isS3Service(Model model, ServiceShape service) {
        return S3ModelUtils.isServiceS3(model, service);
    }

    /**
     * returns name of the deserializer middleware written wrt this customization.
     *
     * @param service the service closure for the operation.
     * @param operation the operation for which custom deserializer is generated.
     */
    private String getDeserializeMiddlewareName(ServiceShape service, OperationShape operation) {
        return ProtocolGenerator.getDeserializeMiddlewareName(operation.getId(), service, protocolName) + "_custom";
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            return isS3Service(model, service) && operation.getId().getName()
                                    .equals(getBucketLocationOpID);
                        })
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(
                                        SymbolUtils.createValueSymbolBuilder(swapDeserializerFuncName).build())
                                .build())
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ShapeId serviceId = settings.getService();
        ServiceShape service = model.expectShape(serviceId, ServiceShape.class);
        if (!isS3Service(model, service)) {
            return;
        }

        for (OperationShape operation : TopDownIndex.of(model).getContainedOperations(service)) {
            if (!(operation.getId().getName(service).equals(getBucketLocationOpID))) {
                continue;
            }

            if (!settings.useLegacySerde()) {
                goDelegator.useShapeWriter(operation, this::writeUnwrappedXmlOutputDeserializer);
                continue;
            }

            goDelegator.useShapeWriter(operation, writer -> {
                writeCustomDeserializer(writer, model, symbolProvider, service, operation);
                writeDeserializerSwapFunction(writer, service, operation);
            });
        }

    }

    /**
     * Schema serde has no per-operation document deserializer to delegate to, so instead of replacing the
     * OperationDeserializer this wraps it. S3 returns the bare member element (<LocationConstraint>us-west-2</...>)
     * with no element for the output structure, so the wrapper puts the body of a successful response inside a
     * synthetic root element before the standard deserializer reads it.
     */
    private void writeUnwrappedXmlOutputDeserializer(GoWriter writer) {
        writer.addUseImports(SmithyGoDependency.FMT);
        writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);

        new UnwrappedXmlOutputDeserializer().accept(writer);
        writer.write("");

        var id = ProtocolUtils.OPERATION_DESERIALIZER_MIDDLEWARE_ID.getString().orElseThrow();
        writer.writeDocs("Helper to wrap the standard deserializer");
        writer.openBlock("func $L(stack *middleware.Stack) error {", "}", swapDeserializerFuncName, () -> {
            writer.write("m, ok := stack.Deserialize.Get($S)", id);
            writer.openBlock("if !ok {", "}", () -> {
                writer.write("return fmt.Errorf($S)", id + " middleware not found");
            });
            writer.write("");
            writer.write("_, err := stack.Deserialize.Swap($S, &$L{delegate: m})", id, UNWRAPPED_DESERIALIZER_NAME);
            writer.write("return err");
        });
    }

    private static class UnwrappedXmlOutputDeserializer extends DeserializeStepMiddleware {
        @Override
        public String getStructName() {
            return UNWRAPPED_DESERIALIZER_NAME;
        }

        @Override
        public String getId() {
            return ProtocolUtils.OPERATION_DESERIALIZER_MIDDLEWARE_ID.getString().orElseThrow();
        }

        @Override
        public Map<String, Symbol> getFields() {
            return Map.of("delegate", SmithyGoDependency.SMITHY_MIDDLEWARE.interfaceSymbol("DeserializeMiddleware"));
        }

        @Override
        public Writable getFuncBody() {
            return writer -> {
                writer.addUseImports(SmithyGoDependency.IO);
                writer.addUseImports(SmithyGoDependency.STRINGS);
                writer.addUseImports(SmithyGoDependency.SMITHY_HTTP_TRANSPORT);

                writer.openBlock("return m.delegate.HandleDeserialize(ctx, in, middleware.DeserializeHandlerFunc(func(",
                        "}))", () -> {
                    writer.write("ctx context.Context, in middleware.DeserializeInput,");
                    writer.write(") (");
                    writer.write("middleware.DeserializeOutput, middleware.Metadata, error,");
                    writer.write(") {");
                    writer.write("out, md, err := next.HandleDeserialize(ctx, in)");
                    writer.openBlock("if err != nil {", "}", () -> {
                        writer.write("return out, md, err");
                    });
                    writer.write("");
                    writer.write("resp, ok := out.RawResponse.(*smithyhttp.Response)");
                    writer.openBlock("if !ok || resp.StatusCode < 200 || resp.StatusCode >= 300 {", "}", () -> {
                        writer.write("return out, md, nil");
                    });
                    writer.write("");
                    writer.openBlock("resp.Body = struct {", "}{", () -> {
                        writer.write("io.Reader");
                        writer.write("io.Closer");
                    });
                    writer.indent();
                    writer.openBlock("io.MultiReader(", "),", () -> {
                        writer.write("strings.NewReader($S),", "<LocationConstraint>");
                        writer.write("resp.Body,");
                        writer.write("strings.NewReader($S),", "</LocationConstraint>");
                    });
                    writer.write("resp.Body,");
                    writer.dedent();
                    writer.write("}");
                    writer.write("return out, md, nil");
                });
            };
        }
    }

    /**
     * writes helper function to swap deserialization middleware with the generated
     * custom deserializer middleware.
     *
     * @param writer    is the go writer used
     * @param operation is the operation for which swap function is written.
     */
    private void writeDeserializerSwapFunction(
            GoWriter writer,
            ServiceShape service,
            OperationShape operation
    ) {
        writer.writeDocs("Helper to swap in a custom deserializer");
        writer.openBlock("func $L(stack *middleware.Stack) error{", "}",
                swapDeserializerFuncName, () -> {
                    writer.write("_, err := stack.Deserialize.Swap($S, &$L{})",
                            ProtocolUtils.OPERATION_DESERIALIZER_MIDDLEWARE_ID.getString(),
                            getDeserializeMiddlewareName(service, operation)
                    );
                    writer.write("if err != nil { return err }");
                    writer.write("return nil");
                });
    }

    /**
     * writes a custom deserializer middleware for the provided operation.
     *
     * @param goWriter       is the go writer used.
     * @param model          is the generation model.
     * @param symbolProvider is the symbol provider.
     * @param operation      is the operation shape for which custom deserializer is written.
     */
    private void writeCustomDeserializer(
            GoWriter goWriter,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape service,
            OperationShape operation
    ) {
        GoStackStepMiddlewareGenerator middleware = GoStackStepMiddlewareGenerator.createDeserializeStepMiddleware(
                getDeserializeMiddlewareName(service, operation), ProtocolUtils.OPERATION_DESERIALIZER_MIDDLEWARE_ID);

        String errorFunctionName = ProtocolGenerator.getOperationErrorDeserFunctionName(
                operation, service, protocolName);

        middleware.writeMiddleware(goWriter, (generator, writer) -> {
            writer.addUseImports(SmithyGoDependency.FMT);

            writer.write("out, metadata, err = next.$L(ctx, in)", generator.getHandleMethodName());
            writer.write("if err != nil { return out, metadata, err }");
            writer.write("");

            writer.addUseImports(SmithyGoDependency.SMITHY_HTTP_TRANSPORT);
            writer.write("response, ok := out.RawResponse.(*smithyhttp.Response)");
            writer.openBlock("if !ok {", "}", () -> {
                writer.addUseImports(SmithyGoDependency.SMITHY);
                writer.write(String.format("return out, metadata, &smithy.DeserializationError{Err: %s}",
                        "fmt.Errorf(\"unknown transport type %T\", out.RawResponse)"));
            });
            writer.write("");

            writer.openBlock("if response.StatusCode < 200 || response.StatusCode >= 300 {", "}", () -> {
                writer.write("return out, metadata, $L(response, &metadata)", errorFunctionName);
            });

            Shape outputShape = model.expectShape(operation.getOutput()
                    .orElseThrow(() -> new CodegenException("expect output shape for operation: " + operation.getId()))
            );

            Symbol outputSymbol = symbolProvider.toSymbol(outputShape);

            // initialize out.Result as output structure shape
            writer.write("output := &$T{}", outputSymbol);
            writer.write("out.Result = output");
            writer.write("");

            writer.addUseImports(SmithyGoDependency.XML);
            writer.addUseImports(SmithyGoDependency.SMITHY_XML);
            writer.addUseImports(SmithyGoDependency.IO);
            writer.addUseImports(SmithyGoDependency.SMITHY_IO);

            writer.write("var buff [1024]byte");
            writer.write("ringBuffer := smithyio.NewRingBuffer(buff[:])");
            writer.write("body := io.TeeReader(response.Body, ringBuffer)");
            writer.write("rootDecoder := xml.NewDecoder(body)");

            // define a decoder with empty start element since we s3 does not wrap Location Constraint
            // xml tag with operation specific xml tag.
            writer.write("decoder := smithyxml.WrapNodeDecoder(rootDecoder, xml.StartElement{})");

            String deserFuncName = ProtocolGenerator.getDocumentDeserializerFunctionName(outputShape, service, protocolName);
            writer.addUseImports(SmithyGoDependency.IO);

            // delegate to already generated inner body deserializer function.
            writer.write("err = $L(&output, decoder)", deserFuncName);

            // EOF error is valid in this case, as we provide a NOP start element at start.
            // Note that we correctly handle unexpected EOF.
            writer.addUseImports(SmithyGoDependency.IO);
            writer.write("if err == io.EOF { err = nil }");

            XmlProtocolUtils.handleDecodeError(writer, "out, metadata,");

            writer.write("");
            writer.write("return out, metadata, err");
        });
    }
}
