package software.amazon.smithy.aws.go.codegen;

import java.util.Set;
import software.amazon.smithy.aws.traits.protocols.Ec2QueryTrait;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;

/**
 * Handles generating the ec2 query protocol for services.
 *
 * @inheritDoc
 *
 * @see AwsQuery
 */
final class Ec2Query extends AwsQuery {
    @Override
    public ShapeId getProtocol() {
        return Ec2QueryTrait.ID;
    }

    @Override
    protected void generateDocumentBodyShapeSerializers(GenerationContext context, Set<Shape> shapes) {
        Ec2QueryShapeSerVisitor visitor = new Ec2QueryShapeSerVisitor(context);
        shapes.forEach(shape -> shape.accept(visitor));
    }

    @Override
    protected void generateDocumentBodyShapeDeserializers(GenerationContext context, Set<Shape> shapes) {
        XmlShapeDeserVisitor visitor = new XmlShapeDeserVisitor(context);
        shapes.forEach(shape -> shape.accept(visitor));
    }

    @Override
    protected void writeErrorMessageCodeDeserializer(GenerationContext context) {
        GoWriter writer = context.getWriter();
        writer.addUseImports(AwsGoDependency.AWS_EC2QUERY_PROTOCOL);
        writer.write("errorComponents, err := ec2query.GetErrorResponseComponents(errorBody)");
        writer.write("if err != nil { return err }");
        writer.insertTrailingNewline();

        writer.addUseImports(AwsGoDependency.AWS_MIDDLEWARE);
        writer.write("awsmiddleware.SetRequestIDMetadata(metadata, errorComponents.RequestID)");
        writer.insertTrailingNewline();

        writer.write("if len(errorComponents.Code) != 0 { errorCode = errorComponents.Code}");
        writer.write("if len(errorComponents.Message) != 0 { errorMessage = errorComponents.Message}");
        writer.insertTrailingNewline();

        writer.write("errorBody.Seek(0, io.SeekStart)");
        writer.insertTrailingNewline();
    }
}
