package software.amazon.smithy.aws.go.codegen;

import java.util.Set;
import software.amazon.smithy.aws.traits.protocols.AwsQueryTrait;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.HttpRpcProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
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
final class AwsQuery extends HttpRpcProtocolGenerator {

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

        writer.write("bodyEncoder := query.NewEncoder()");
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

        writer.write("encodedBody, err := bodyEncoder.Encode()");
        writer.openBlock("if err != nil {", "}", () -> {
            writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
        }).write("");

        writer.addUseImports(SmithyGoDependency.BYTES);
        writer.openBlock("if request, err = request.SetStream(bytes.NewReader(encodedBody)); err != nil {",
                "}", () -> {
            writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
        });
    }

    @Override
    protected void generateDocumentBodyShapeDeserializers(GenerationContext context, Set<Shape> shapes) {
        // TODO: support query deser
    }

    @Override
    protected void deserializeOutputDocument(GenerationContext context, OperationShape operation) {
        // TODO: support query deser
    }

    @Override
    protected void deserializeError(GenerationContext context, StructureShape shape) {
        // TODO: support query error deser
        GoWriter writer = context.getWriter();
        writer.writeDocs("TODO: support query error deser");
        writer.write("return &smithy.DeserializationError{Err: fmt.Errorf(\"TODO: support query error deser\")}");
    }

    @Override
    protected void writeErrorMessageCodeDeserializer(GenerationContext context) {
        // TODO: support query error message / code deser
        context.getWriter().write("_ = errorBody");
    }

    @Override
    public void generateProtocolTests(GenerationContext context) {
        AwsProtocolUtils.generateHttpProtocolTests(context);
    }
}
