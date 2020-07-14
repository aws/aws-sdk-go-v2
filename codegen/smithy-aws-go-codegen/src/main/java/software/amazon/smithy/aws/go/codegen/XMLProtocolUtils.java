package software.amazon.smithy.aws.go.codegen;

import static software.amazon.smithy.go.codegen.integration.ProtocolUtils.writeSafeMemberAccessor;

import java.util.Optional;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SyntheticClone;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.model.traits.XmlAttributeTrait;
import software.amazon.smithy.model.traits.XmlNameTrait;
import software.amazon.smithy.model.traits.XmlNamespaceTrait;

/**
 * Utility methods for generating XML protocols.
 */
final class XMLProtocolUtils {
    private XMLProtocolUtils() {
    }

    /**
     * generateXMLStartElement generates the XML start element for a shape. It is used to generate smithy xml's startElement.
     *
     * @param context  is the generation context.
     * @param shape    is the Shape for which xml start element is to be generated.
     * @param dst      is the operand name which holds the generated start element.
     * @param inputSrc is the input variable for the shape with values to be serialized.
     */
    public static void generateXMLStartElement(
            ProtocolGenerator.GenerationContext context, Shape shape, String dst, String inputSrc
    ) {
        GoWriter writer = context.getWriter();
        String attrName = dst + "Attr";
        writer.write("$L := []smithyxml.Attr{}", attrName);

        Optional<XmlNamespaceTrait> xmlNamespaceTrait = shape.getTrait(XmlNamespaceTrait.class);
        if (xmlNamespaceTrait.isPresent()) {
            XmlNamespaceTrait namespace = xmlNamespaceTrait.get();
            writer.write("$L = append($L, smithyxml.NewNamespaceAttribute($S, $S))",
                    attrName, attrName,
                    namespace.getPrefix().isPresent() ? namespace.getPrefix().get() : "", namespace.getUri()
            );
        }

        // Traverse member shapes to get attributes
        shape.members().stream().forEach(memberShape -> {
            if (memberShape.hasTrait(XmlAttributeTrait.class)) {
                writeSafeMemberAccessor(context, memberShape, inputSrc, (operand) -> {
                    // xml attributes should always be string
                    String dest = "av";
                    formatXmlAttributeValueAsString(context, memberShape, operand, dest);
                    writer.write("$L = append($L, smithyxml.NewAttribute($S, $L))",
                            attrName, attrName, getSerializedXMLMemberName(memberShape), dest);
                });
            }
        });

        writer.openBlock("$L := smithyxml.StartElement{ ", "}", dst, () -> {
            writer.openBlock("Name:smithyxml.Name{", "},", () -> {
                writer.write("Local: $S,", getSerializedXMLShapeName(context, shape));
            });
            writer.write("Attr : $L,", attrName);
        });
    }

    // generates code to format xml attributes. If a shape type is timestamp, number, or boolean
    // it will be formatted into a string.
    private static void formatXmlAttributeValueAsString(
            ProtocolGenerator.GenerationContext context,
            MemberShape member, String src, String dest
    ) {
        GoWriter writer = context.getWriter();
        Shape target = context.getModel().expectShape(member.getTarget());

        // declare destination variable
        writer.write("var $L string", dest);

        if (target.isStringShape()) {
            writer.write("$L = *$L", dest, src);
            return;
        }

        if (target.isTimestampShape() || target.hasTrait(TimestampFormatTrait.class)) {
            TimestampFormatTrait.Format format = member.getMemberTrait(context.getModel(), TimestampFormatTrait.class)
                    .map(TimestampFormatTrait::getFormat).orElse(TimestampFormatTrait.Format.DATE_TIME);
            writer.addUseImports(SmithyGoDependency.SMITHY_TIME);
            switch (format) {
                case DATE_TIME:
                    writer.write("$L = smithytime.FormatDateTime(*$L)", dest, src);
                    break;
                case HTTP_DATE:
                    writer.write("$L = smithytime.FormatHTTPDate(*$L)", dest, src);
                    break;
                case EPOCH_SECONDS:
                    writer.addUseImports(SmithyGoDependency.STRCONV);
                    writer.write("$L = strconv.FormatFloat(smithytime.FormatEpochSeconds(*$L), 'f', -1, 64)", dest,
                            src);
                    break;
                case UNKNOWN:
                    throw new CodegenException("Unknown timestamp format");
            }
            return;
        }

        if (target.isBooleanShape()) {
            writer.write(SmithyGoDependency.STRCONV);
            writer.write("$L = strconv.FormatBool(*$L)", dest, src);
            return;
        }

        if (target.isByteShape() || target.isShortShape() || target.isIntegerShape() || target.isLongShape()) {
            writer.write(SmithyGoDependency.STRCONV);
            writer.write("$L = strconv.FormatInt(int64(*$L), 10)", dest, src);
            return;
        }

        if (target.isFloatShape()) {
            writer.write(SmithyGoDependency.STRCONV);
            writer.write("$L = strconv.FormatFloat(float64(*$L),'f', -1, 32)", dest, src);
            return;
        }

        if (target.isDoubleShape()) {
            writer.write(SmithyGoDependency.STRCONV);
            writer.write("$L = strconv.FormatFloat(*$L,'f', -1, 64)", dest, src);
            return;
        }

        if (target.isBigIntegerShape() || target.isBigDecimalShape()) {
            throw new CodegenException(String.format("Cannot serialize shape type %s on protocol, shape: %s.",
                    target.getType(), target.getId()));
        }

        throw new CodegenException(
                "Members serialized as XML attributes can only be of string, number, boolean or timestamp format");
    }

    /**
     * getSerializedXMLMemberName returns a xml member name used for serializing. If a member shape has
     * XML name trait, xml name would be given precedence over member name.
     *
     * @param memberShape is the member shape for which serializer name is queried.
     * @return name of a xml member shape used by serializers
     */
    private static String getSerializedXMLMemberName(MemberShape memberShape) {
        Optional<XmlNameTrait> xmlNameTrait = memberShape.getTrait(XmlNameTrait.class);
        return xmlNameTrait.isPresent() ? xmlNameTrait.get().getValue() : memberShape.getMemberName();
    }

    /**
     * getSerializedXMLShapeName returns a xml shape name used for serializing. If a member shape
     * has xml name trait, xml name would be given precedence over member name.
     * This correctly handles renamed shapes, and returns the original shape name.
     *
     * @param context is the generation context for which
     * @param shape   is the Shape for which serializer name is queried.
     * @return name of a xml member shape used by serializers.
     */
    private static String getSerializedXMLShapeName(ProtocolGenerator.GenerationContext context, Shape shape) {
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol shapeSymbol = symbolProvider.toSymbol(shape);
        String shapeName = shapeSymbol.getName();

        // check if synthetic cloned shape
        Optional<SyntheticClone> clone = shape.getTrait(SyntheticClone.class);
        if (clone.isPresent()) {
            SyntheticClone cl = clone.get();
            shapeName = cl.getArchetype().getName();
        }

        // check if shape is member shape
        Optional<MemberShape> member = shape.asMemberShape();
        if (member.isPresent()) {
            return getSerializedXMLMemberName(member.get());
        }

        return shape.getTrait(XmlNameTrait.class).map(XmlNameTrait::getValue).orElse(shapeName);
    }
}
