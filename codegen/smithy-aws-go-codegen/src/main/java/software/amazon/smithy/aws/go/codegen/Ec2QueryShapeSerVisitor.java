package software.amazon.smithy.aws.go.codegen;

import java.util.Optional;
import java.util.function.Predicate;
import software.amazon.smithy.aws.traits.protocols.Ec2QueryNameTrait;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;
import software.amazon.smithy.go.codegen.knowledge.GoPointableIndex;
import software.amazon.smithy.model.shapes.CollectionShape;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeType;
import software.amazon.smithy.model.traits.TimestampFormatTrait.Format;
import software.amazon.smithy.model.traits.XmlNameTrait;
import software.amazon.smithy.utils.StringUtils;

/**
 * Visitor to generate serialization functions for shapes in EC2 Query protocol
 * document bodies.
 *
 * This class uses the implementations provided by {@code QueryShapeSerVisitor} but with
 * the following protocol specific customizations for ec2 query:
 *
 * <ul>
 *   <li>ec2 query flattens all lists, sets, and maps regardless of the {@code @xmlFlattened} trait.</li>
 *   <li>ec2 query respects the {@code @ec2QueryName} trait, then the {@code xmlName}
 *     trait value with the first letter capitalized.</li>
 * </ul>
 *
 * Timestamps are serialized to {@link Format}.DATE_TIME by default.
 *
 * @see QueryShapeSerVisitor
 */
final class Ec2QueryShapeSerVisitor extends QueryShapeSerVisitor {
    public Ec2QueryShapeSerVisitor(GenerationContext context) {
        super(context);
    }

    public Ec2QueryShapeSerVisitor(GenerationContext context, Predicate<MemberShape> memberFilter) {
        super(context, memberFilter);
    }

    @Override
    protected String getSerializedLocationName(MemberShape memberShape, String defaultValue) {
        // The serialization for aws.ec2 prioritizes the @ec2QueryName trait for serialization.
        Optional<Ec2QueryNameTrait> trait = memberShape.getTrait(Ec2QueryNameTrait.class);
        if (trait.isPresent()) {
            return trait.get().getValue();
        }

        // Fall back to the capitalized @xmlName trait if present on the member,
        // otherwise use the capitalized default value.
        return StringUtils.capitalize(memberShape.getTrait(XmlNameTrait.class)
                .map(XmlNameTrait::getValue)
                .orElse(defaultValue));
    }

    // EC2Query specifically does not serialize non-nil, empty lists
    protected void serializeCollection(GenerationContext context, CollectionShape shape) {
        context.getWriter().get()
                .write("if len(v) == 0 { return nil }");

        super.serializeCollection(context, shape);
    }

    @Override
    protected boolean isFlattened(GenerationContext context, MemberShape memberShape) {
        // All lists, sets, and maps are flattened in aws.ec2.
        ShapeType targetType = context.getModel().expectShape(memberShape.getTarget()).getType();
        return targetType == ShapeType.LIST || targetType == ShapeType.SET || targetType == ShapeType.MAP;
    }
}
