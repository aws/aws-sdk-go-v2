package software.amazon.smithy.aws.go.codegen;

import java.util.Set;
import software.amazon.smithy.aws.traits.protocols.AwsQueryTrait;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.HttpRpcProtocolGenerator;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;

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

    }

    @Override
    protected void serializeInputDocument(GenerationContext context, OperationShape operation) {
        GoWriter writer = context.getWriter();
        writer.write("_ = input");
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
