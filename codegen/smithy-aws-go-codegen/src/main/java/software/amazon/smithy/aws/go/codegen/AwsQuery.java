package software.amazon.smithy.aws.go.codegen;

import static software.amazon.smithy.go.codegen.integration.HttpProtocolGeneratorUtils.isShapeWithResponseBindings;
import static software.amazon.smithy.aws.go.codegen.AwsProtocolUtils.handleDecodeError;
import static software.amazon.smithy.aws.go.codegen.XmlProtocolUtils.initializeXmlDecoder;

import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.TreeMap;
import software.amazon.smithy.aws.traits.protocols.AwsQueryErrorTrait;
import software.amazon.smithy.aws.traits.protocols.AwsQueryTrait;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.SyntheticClone;
import software.amazon.smithy.go.codegen.integration.HttpRpcProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.model.knowledge.EventStreamInfo;
import software.amazon.smithy.model.knowledge.HttpBinding;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.shapes.UnionShape;
import software.amazon.smithy.model.traits.ErrorTrait;

/**
 * Handles generating the aws query protocol for services.
 *
 * @inheritDoc
 * @see HttpRpcProtocolGenerator
 */
class AwsQuery extends HttpRpcProtocolGenerator {

    @Override
    public ShapeId getProtocol() {
        return AwsQueryTrait.ID;
    }

    @Override
    protected String getOperationPath(GenerationContext context, OperationShape operation) {
        return "/";
    }

    @Override
    protected String getDocumentContentType() {
        return "application/x-www-form-urlencoded";
    }

    @Override
    protected void generateDocumentBodyShapeSerializers(GenerationContext context, Set<Shape> shapes) {
        QueryShapeSerVisitor visitor = new QueryShapeSerVisitor(context);
        shapes.forEach(shape -> shape.accept(visitor));
    }

    @Override
    protected void serializeInputDocument(GenerationContext context, OperationShape operation) {
        GoWriter writer = context.getWriter().get();
        ServiceShape service = context.getService();
        StructureShape input = ProtocolUtils.expectInput(context.getModel(), operation);
        String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(
                input, context.getService(), getProtocolName());
        writer.addUseImports(AwsGoDependency.AWS_QUERY_PROTOCOL);

        writer.addUseImports(SmithyGoDependency.BYTES);
        writer.write("bodyWriter := bytes.NewBuffer(nil)");
        writer.write("bodyEncoder := query.NewEncoder(bodyWriter)");
        writer.write("body := bodyEncoder.Object()");
        writer.write("body.Key(\"Action\").String($S)", operation.getId().getName(service));
        writer.write("body.Key(\"Version\").String($S)", service.getVersion());
        writer.write("");

        if (!input.members().isEmpty()) {
            writer.openBlock("if err := $L(input, bodyEncoder.Value); err != nil {", "}", functionName, () -> {
                writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
            }).write("");
        }

        writer.write("err = bodyEncoder.Encode()");
        writer.openBlock("if err != nil {", "}", () -> {
            writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
        }).write("");

        writer.openBlock("if request, err = request.SetStream(bytes.NewReader(bodyWriter.Bytes())); err != nil {",
                "}", () -> {
                    writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
                });
    }

    @Override
    protected void generateDocumentBodyShapeDeserializers(GenerationContext context, Set<Shape> shapes) {
        XmlShapeDeserVisitor visitor = new XmlShapeDeserVisitor(context);
        shapes.forEach(shape -> shape.accept(visitor));
    }

    @Override
    protected void deserializeOutputDocument(GenerationContext context, OperationShape operation) {
        GoWriter writer = context.getWriter().get();
        StructureShape output = ProtocolUtils.expectOutput(context.getModel(), operation);
        String functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(
                output, context.getService(), getProtocolName());
        initializeXmlDecoder(writer, "response.Body", "out, metadata, ", "nil");
        unwrapOutputDocument(context, operation);
        writer.write("err = $L(&output, decoder)", functionName);
        handleDecodeError(writer, "out, metadata, ");
    }

    @Override
    protected void deserializeError(GenerationContext context, StructureShape shape) {
        GoWriter writer = context.getWriter().get();
        Symbol symbol = context.getSymbolProvider().toSymbol(shape);

        writer.write("output := &$T{}", symbol);
        writer.insertTrailingNewline();
        if (isShapeWithResponseBindings(context.getModel(), shape, HttpBinding.Location.DOCUMENT)) {
            String functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(
                    shape, context.getService(), context.getProtocolName());

            writer.addUseImports(SmithyGoDependency.IO);
            initializeXmlDecoder(writer, "errorBody", "output");
            unwrapErrorElement(context);
            writer.write("err = $L(&output, decoder)", functionName);
            XmlProtocolUtils.handleDecodeError(writer, "");
            writer.insertTrailingNewline();
        }
        writer.write("return output");
    }

    protected void unwrapOutputDocument(GenerationContext context, OperationShape shape) {
        GoWriter writer = context.getWriter().get();
        ServiceShape service = context.getService();
        writer.write("t, err = decoder.GetElement(\"$LResult\")", shape.getId().getName(service));
        handleDecodeError(writer, "out, metadata, ");
        Symbol wrapNodeDecoder = SymbolUtils.createValueSymbolBuilder("WrapNodeDecoder",
                SmithyGoDependency.SMITHY_XML).build();
        writer.write("decoder = $T(decoder.Decoder, t)", wrapNodeDecoder);
    }

    protected void unwrapErrorElement(GenerationContext context) {
        GoWriter writer = context.getWriter().get();
        writer.write("t, err = decoder.GetElement(\"Error\")");
        XmlProtocolUtils.handleDecodeError(writer, "");
        Symbol wrapNodeDecoder = SymbolUtils.createValueSymbolBuilder("WrapNodeDecoder",
                SmithyGoDependency.SMITHY_XML).build();
        writer.write("decoder = $T(decoder.Decoder, t)", wrapNodeDecoder);
    }

    @Override
    protected void writeErrorMessageCodeDeserializer(GenerationContext context) {
        XmlProtocolUtils.writeXmlErrorMessageCodeDeserializer(context);
    }

    @Override
    public void generateProtocolTests(GenerationContext context) {
        AwsProtocolUtils.generateHttpProtocolTests(context);
    }

    @Override
    public Map<String, ShapeId> getOperationErrors(GenerationContext context, OperationShape operation) {
        Map<String, ShapeId> errors = new TreeMap<>();

        operation.getErrors().forEach(shapeId -> {
            Shape errorShape = context.getModel().expectShape(shapeId);
            String errorName = shapeId.getName(context.getService());

            Optional<AwsQueryErrorTrait> errorShapeTrait = errorShape.getTrait(AwsQueryErrorTrait.class);
            if (errorShapeTrait.isPresent()) {
                errors.put(errorShapeTrait.get().getCode(), shapeId);
            } else {
                errors.put(errorName, shapeId);
            }
        });

        return errors;
    }

    @Override
    public String getErrorCode(ServiceShape service, StructureShape errorShape) {
        Optional<AwsQueryErrorTrait> trait = errorShape.getTrait(AwsQueryErrorTrait.class);
        if (trait.isPresent()) {
            return trait.get().getCode();
        }

        return super.getErrorCode(service, errorShape);
    }

    @Override
    protected void generateEventStreamSerializers(
            GenerationContext context,
            UnionShape eventUnion,
            Set<EventStreamInfo> eventStreamInfos
    ) {
        throw new CodegenException("event streams not supported with AWS QUERY protocol.");
    }

    @Override
    protected void generateEventStreamDeserializers(
            GenerationContext context,
            UnionShape eventUnion,
            Set<EventStreamInfo> eventStreamInfos
    ) {
        throw new CodegenException("event streams not supported with AWS QUERY protocol.");
    }

    @Override
    public void generateEventStreamComponents(GenerationContext context) {
        throw new CodegenException("event streams not supported with AWS QUERY protocol.");
    }

    @Override
    protected void writeOperationSerializerMiddlewareEventStreamSetup(GenerationContext context, EventStreamInfo info) {
        throw new CodegenException("event streams not supported with AWS QUERY protocol.");
    }
}
