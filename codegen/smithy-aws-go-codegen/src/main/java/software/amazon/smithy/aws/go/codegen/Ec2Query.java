package software.amazon.smithy.aws.go.codegen;

import java.util.Set;
import software.amazon.smithy.aws.traits.protocols.Ec2QueryTrait;
import software.amazon.smithy.go.codegen.integration.HttpRpcProtocolGenerator;
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
}
