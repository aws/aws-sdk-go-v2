package software.amazon.smithy.aws.go.codegen;

import java.util.Collections;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;
import java.util.function.Predicate;
import java.util.logging.Logger;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoValueAccessUtils;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.DocumentShapeSerVisitor;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;
import software.amazon.smithy.go.codegen.trait.NoSerializeTrait;
import software.amazon.smithy.model.shapes.CollectionShape;
import software.amazon.smithy.model.shapes.DocumentShape;
import software.amazon.smithy.model.shapes.MapShape;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.shapes.UnionShape;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.model.traits.XmlFlattenedTrait;
import software.amazon.smithy.model.traits.XmlNameTrait;
import software.amazon.smithy.model.traits.XmlNamespaceTrait;
import software.amazon.smithy.utils.FunctionalUtils;

final class XmlShapeSerVisitor extends DocumentShapeSerVisitor {
    private static final TimestampFormatTrait.Format DEFAULT_TIMESTAMP_FORMAT = TimestampFormatTrait.Format.DATE_TIME;
    private static final Logger LOGGER = Logger.getLogger(XmlShapeSerVisitor.class.getName());

    private final Predicate<MemberShape> memberFilter;

    public XmlShapeSerVisitor(GenerationContext context) {
        this(context, NoSerializeTrait.excludeNoSerializeMembers().and(FunctionalUtils.alwaysTrue()));
    }

    public XmlShapeSerVisitor(GenerationContext context, Predicate<MemberShape> memberFilter) {
        super(context);
        this.memberFilter = NoSerializeTrait.excludeNoSerializeMembers().and(memberFilter);
    }

    private DocumentMemberSerVisitor getMemberSerVisitor(MemberShape member, String source, String dest) {
        // Get the timestamp format to be used, defaulting to date-time format.
        TimestampFormatTrait.Format format = member.getMemberTrait(getContext().getModel(), TimestampFormatTrait.class)
                .map(TimestampFormatTrait::getFormat).orElse(DEFAULT_TIMESTAMP_FORMAT);
        return new DocumentMemberSerVisitor(getContext(), member, source, dest, format);
    }

    @Override
    protected Map<String, String> getAdditionalSerArguments() {
        return Collections.singletonMap("value", "smithyxml.Value");
    }

    @Override
    protected void serializeCollection(GenerationContext context, CollectionShape shape) {
        GoWriter writer = context.getWriter().get();
        Shape target = context.getModel().expectShape(shape.getMember().getTarget());
        MemberShape member = shape.getMember();
        writer.write("var array *smithyxml.Array");
        writer.openBlock("if !value.IsFlattened() {", "}", () -> {
            writer.write("defer value.Close()");
        });

        if (member.hasTrait(XmlNameTrait.class) || member.hasTrait(XmlNamespaceTrait.class)) {
            XmlProtocolUtils.generateXMLStartElement(context, member, "customMemberName", "v");
            writer.write("array = value.ArrayWithCustomName(customMemberName)");
        } else {
            writer.write("array = value.Array()");
        }
        writer.insertTrailingNewline();

        writer.openBlock("for i := range v {", "}", () -> {
            // Serialize zero members as empty values.
            GoValueAccessUtils.writeIfZeroValue(context.getModel(), writer, member, "v[i]", () -> {
                writer.write("am := array.Member()");
                writer.write("am.Close()");
                writer.write("continue");
            });

            writer.write("am := array.Member()");
            target.accept(getMemberSerVisitor(shape.getMember(), "v[i]", "am"));
        });

        writer.write("return nil");
    }

    @Override
    protected void serializeMap(GenerationContext context, MapShape shape) {
        GoWriter writer = context.getWriter().get();
        Shape targetKey = context.getModel().expectShape(shape.getKey().getTarget());
        Shape targetValue = context.getModel().expectShape(shape.getValue().getTarget());

        writer.openBlock("if !value.IsFlattened() {", "}", () -> {
            writer.write("defer value.Close()");
        });

        writer.write("m := value.Map()");
        writer.insertTrailingNewline();

        writer.openBlock("for key := range v {", "}", () -> {
            writer.write("entry := m.Entry()");
            writer.insertTrailingNewline();

            // Serialize zero values as empty values.
            GoValueAccessUtils.writeIfZeroValue(context.getModel(), writer, shape.getValue(), "v[key]", () -> {
                writer.write("entry.Close()");
                writer.write("continue");
            });

            // map entry key
            XmlProtocolUtils.generateXMLStartElement(context, shape.getKey(), "keyElement", "v");
            targetKey.accept(getMemberSerVisitor(shape.getKey(), "key", "entry.MemberElement(keyElement)"));
            writer.insertTrailingNewline();

            // map entry value
            XmlProtocolUtils.generateXMLStartElement(context, shape.getValue(), "valueElement", "v");
            String dest = "entry.MemberElement(valueElement)";
            if (shape.getValue().hasTrait(XmlFlattenedTrait.class)) {
                dest = "entry.FlattenedElement(valueElement)";
            }
            targetValue.accept(getMemberSerVisitor(shape.getValue(), "v[key]", dest));

            // close the map entry
            writer.write("entry.Close()");
            writer.insertTrailingNewline();
        });
        writer.write("return nil");
    }

    @Override
    protected void serializeStructure(GenerationContext context, StructureShape shape) {
        GoWriter writer = context.getWriter().get();

        // defer close xml.value
        writer.write("defer value.Close()");
        writer.insertTrailingNewline();

        // Use a tree sort to sort the members.
        Set<MemberShape> members = new TreeSet<>(shape.getAllMembers().values());
        for (MemberShape member : members) {
            if (!memberFilter.test(member)) {
                continue;
            }

            Shape target = context.getModel().expectShape(member.getTarget());

            writer.addUseImports(SmithyGoDependency.SMITHY_XML);

            GoValueAccessUtils.writeIfNonZeroValueMember(context.getModel(), context.getSymbolProvider(), writer,
                    member, "v", true, member.isRequired(), (operand) -> {
                        XmlProtocolUtils.generateXMLStartElement(context, member, "root", "v");

                        // check if member shape has flattened trait
                        if (member.hasTrait(XmlFlattenedTrait.class)) {
                            writer.write("el := value.FlattenedElement($L)", "root");
                        } else {
                            writer.write("el := value.MemberElement($L)", "root");
                        }
                        target.accept(getMemberSerVisitor(member, operand, "el"));
                    });

            writer.insertTrailingNewline();
        }

        writer.write("return nil");
    }

    @Override
    protected void serializeUnion(GenerationContext context, UnionShape shape) {
        GoWriter writer = context.getWriter().get();
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol symbol = symbolProvider.toSymbol(shape);
        writer.addUseImports(SmithyGoDependency.FMT);

        writer.write("defer value.Close()");
        writer.insertTrailingNewline();

        writer.openBlock("switch uv := v.(type) {", "}", () -> {
            // Use a TreeSet to sort the members.
            Set<MemberShape> members = new TreeSet<>(shape.getAllMembers().values());
            for (MemberShape member : members) {
                Shape target = context.getModel().expectShape(member.getTarget());
                Symbol memberSymbol = SymbolUtils.createValueSymbolBuilder(
                        symbolProvider.toMemberName(member),
                        symbol.getNamespace()
                ).build();

                writer.openBlock("case *$T:", "", memberSymbol, () -> {
                    XmlProtocolUtils.generateXMLStartElement(context, member, "customMemberName", "v");
                    writer.write("av := value.MemberElement(customMemberName)");
                    target.accept(getMemberSerVisitor(member, "uv.Value", "av"));
                });
            }

            // Handle unknown union values
            writer.openBlock("default:", "", () -> {
                writer.write("return fmt.Errorf(\"attempted to serialize unknown member type %T"
                        + " for union %T\", uv, v)");
            });
        });

        writer.write("return nil");
    }

    @Override
    protected void serializeDocument(GenerationContext context, DocumentShape shape) {
        // TODO: implement document serialization
        LOGGER.warning("Document type is currently unsupported for XML serialization.");
        context.getWriter().get().writeDocs("TODO: implement document serialization.");
        context.getWriter().get().write("return nil");
    }
}
