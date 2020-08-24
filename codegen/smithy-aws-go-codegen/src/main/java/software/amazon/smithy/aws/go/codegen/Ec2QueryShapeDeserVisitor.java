package software.amazon.smithy.aws.go.codegen;

import java.util.Optional;
import java.util.function.Predicate;
import software.amazon.smithy.aws.traits.protocols.Ec2QueryNameTrait;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;
import software.amazon.smithy.model.shapes.MemberShape;

final class Ec2QueryShapeDeserVisitor extends XmlShapeDeserVisitor {

    /**
     * @param context The generation context.
     */
    public Ec2QueryShapeDeserVisitor(GenerationContext context) {
        super(context);
    }

    public Ec2QueryShapeDeserVisitor(GenerationContext context, Predicate<MemberShape> memberFilter) {
        super(context, memberFilter);
    }

}
