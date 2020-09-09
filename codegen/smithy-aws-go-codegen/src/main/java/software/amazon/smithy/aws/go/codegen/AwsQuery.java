package software.amazon.smithy.aws.go.codegen;

import static software.amazon.smithy.go.codegen.integration.HttpProtocolGeneratorUtils.isShapeWithResponseBindings;
import static software.amazon.smithy.aws.go.codegen.AwsProtocolUtils.handleDecodeError;
import static software.amazon.smithy.aws.go.codegen.XmlProtocolUtils.initializeXmlDecoder;

import java.util.Set;
import software.amazon.smithy.aws.traits.protocols.AwsQueryTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.HttpRpcProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.model.knowledge.HttpBinding;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;

/**
 * Handles generating the aws query protocol for services.
 *
 * @inheritDoc
 *
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
        GoWriter writer = context.getWriter();
        StructureShape input = ProtocolUtils.expectInput(context.getModel(), operation);
        String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(input, getProtocolName());
        writer.addUseImports(AwsGoDependency.AWS_QUERY_PROTOCOL);

        writer.addUseImports(SmithyGoDependency.BYTES);
        writer.write("bodyWriter := bytes.NewBuffer(nil)");
        writer.write("bodyEncoder := query.NewEncoder(bodyWriter)");
        writer.write("body := bodyEncoder.Object()");
        writer.write("body.Key(\"Action\").String($S)", operation.getId().getName());
        writer.write("body.Key(\"Version\").String($S)", context.getService().getVersion());
        writer.write("");

        if (!input.members().isEmpty()) {
            writer.openBlock("if err := $L(input, bodyEncoder.Value); err != nil {", "}", functionName, () -> {
                writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
            }).write("");
        } else {
            writer.write("_ = input");
            writer.write("");
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
        GoWriter writer = context.getWriter();
        StructureShape output = ProtocolUtils.expectOutput(context.getModel(), operation);
        String functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(output, getProtocolName());
        initializeXmlDecoder(writer, "response.Body", "out, metadata, ","nil");
        writer.write("err = $L(&output, decoder)", functionName);
        handleDecodeError(writer, "out, metadata, ");
    }

    @Override
    protected void deserializeError(GenerationContext context, StructureShape shape) {
        GoWriter writer = context.getWriter();
        Symbol symbol = context.getSymbolProvider().toSymbol(shape);

        writer.write("output := &$T{}", symbol);
        writer.insertTrailingNewline();
        if (isShapeWithResponseBindings(context.getModel(), shape, HttpBinding.Location.DOCUMENT)) {
            String documentDeserFunctionName = ProtocolGenerator.getDocumentDeserializerFunctionName(
                    shape, getProtocolName());
            writer.addUseImports(SmithyGoDependency.IO);
            initializeXmlDecoder(writer, "errorBody", "output");
            writer.write("err = $L(&output, decoder)", documentDeserFunctionName);
            XmlProtocolUtils.handleDecodeError(writer, "");
            writer.insertTrailingNewline();
        }
        writer.write("return output");
    }

    @Override
    protected void writeErrorMessageCodeDeserializer(GenerationContext context) {
        XmlProtocolUtils.writeXmlErrorMessageCodeDeserializer(context);
    }

    @Override
    public void generateProtocolTests(GenerationContext context) {
        AwsProtocolUtils.generateHttpProtocolTests(context);
    }
}
