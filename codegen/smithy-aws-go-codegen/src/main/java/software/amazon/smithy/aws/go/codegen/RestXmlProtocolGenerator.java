package software.amazon.smithy.aws.go.codegen;

import java.util.Set;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.HttpBindingProtocolGenerator;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.TimestampFormatTrait;

abstract class RestXmlProtocolGenerator extends HttpBindingProtocolGenerator {
    /**
     * Creates a AWS REST XML protocol generator.
     */
    RestXmlProtocolGenerator() {
        super(true);
    }

    @Override
    protected TimestampFormatTrait.Format getDocumentTimestampFormat() {
        return null;
    }

    @Override
    protected String getDocumentContentType() {
        return null;
    }

    @Override
    protected void generateOperationDocumentSerializer(
            GenerationContext context, OperationShape operation
    ) {

    }

    @Override
    protected void writeMiddlewareDocumentSerializerDelegator(
            Model model, SymbolProvider symbolProvider, OperationShape operation,
            GoStackStepMiddlewareGenerator generator, GoWriter writer
    ) {

    }

    @Override
    protected void writeMiddlewareDocumentDeserializerDelegator(
            GoWriter writer, Model model, SymbolProvider symbolProvider, OperationShape operation,
            GoStackStepMiddlewareGenerator generator
    ) {

    }

    @Override
    protected void writeMiddlewarePayloadSerializerDelegator(
            Model model, SymbolProvider symbolProvider, OperationShape operation, MemberShape memberShape,
            GoStackStepMiddlewareGenerator generator, GoWriter writer
    ) {

    }

    @Override
    protected void generateDocumentBodyShapeSerializers(
            GenerationContext context, Set<Shape> shapes
    ) {

    }

    @Override
    protected void generateOperationDocumentDeserializer(
            GenerationContext context, OperationShape operation
    ) {

    }

    @Override
    protected void generateDocumentBodyShapeDeserializers(
            GenerationContext context, Set<Shape> shapes
    ) {

    }

    @Override
    protected void generateErrorDocumentBindingDeserializer(
            GenerationContext context, ShapeId shapeId
    ) {

    }

    @Override
    public ShapeId getProtocol() {
        return null;
    }
}
