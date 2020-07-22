package software.amazon.smithy.aws.go.codegen;

import static software.amazon.smithy.go.codegen.integration.ProtocolUtils.writeSafeMemberAccessor;

import java.util.Collections;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;
import java.util.function.Predicate;
import java.util.logging.Logger;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDependency;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.DocumentShapeSerVisitor;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;
import software.amazon.smithy.model.shapes.CollectionShape;
import software.amazon.smithy.model.shapes.DocumentShape;
import software.amazon.smithy.model.shapes.MapShape;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.shapes.UnionShape;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait.Format;
import software.amazon.smithy.model.traits.XmlFlattenedTrait;
import software.amazon.smithy.model.traits.XmlNameTrait;
import software.amazon.smithy.utils.FunctionalUtils;

/**
 * Visitor to generate serialization functions for shapes in AWS Query protocol
 * document bodies.
 *
 * This class handles function body generation for all types expected by the
 * {@code DocumentShapeSerVisitor}. No other shape type serialization is overwritten.
 *
 * Timestamps are serialized to {@link Format}.DATE_TIME by default.
 */
final class QueryShapeSerVisitor extends DocumentShapeSerVisitor {
    private static final Format DEFAULT_TIMESTAMP_FORMAT = Format.DATE_TIME;
    private static final Logger LOGGER = Logger.getLogger(JsonShapeSerVisitor.class.getName());

    private final Predicate<MemberShape> memberFilter;

    public QueryShapeSerVisitor(GenerationContext context) {
        this(context, FunctionalUtils.alwaysTrue());
    }

    public QueryShapeSerVisitor(GenerationContext context, Predicate<MemberShape> memberFilter) {
        super(context);
        this.memberFilter = memberFilter;
    }

    private DocumentMemberSerVisitor getMemberSerVisitor(MemberShape member, String source, String dest) {
        // Get the timestamp format to be used, defaulting to epoch seconds.
        Format format = member.getMemberTrait(getContext().getModel(), TimestampFormatTrait.class)
                .map(TimestampFormatTrait::getFormat).orElse(DEFAULT_TIMESTAMP_FORMAT);
        return new DocumentMemberSerVisitor(getContext(), source, dest, format);
    }

    @Override
    protected Map<String, String> getAdditionalSerArguments() {
        return Collections.singletonMap("value", "query.Value");
    }

    @Override
    protected void serializeCollection(GenerationContext context, CollectionShape shape) {
        GoWriter writer = context.getWriter();
        MemberShape member = shape.getMember();
        Shape target = context.getModel().expectShape(member.getTarget());

        writer.write("array := value.Array($S)", getSerializedLocationName(member, "member"));
        writer.write("");

        writer.openBlock("for i := range v {", "}", () -> {
            writer.write("av := array.Value()");

            // Null values should be omitted for query.
            if (!target.hasTrait(EnumTrait.class)) {
                writer.openBlock("if vv := v[i]; vv == nil {", "}", () -> {
                    writer.write("continue");
                });
            }

            target.accept(getMemberSerVisitor(shape.getMember(), "v[i]", "av"));
        });
        writer.write("return nil");
    }

    @Override
    protected void serializeDocument(GenerationContext context, DocumentShape shape) {
        LOGGER.warning("Document type is unsupported for Query serialization.");
        context.getWriter().write("return &smithy.SerializationError{Err: fmt.Errorf("
                + "\"Document type is unsupported for the query protocol.\")}");
    }

    @Override
    protected void serializeMap(GenerationContext context, MapShape shape) {
        GoWriter writer = context.getWriter();
        Shape target = context.getModel().expectShape(shape.getValue().getTarget());

        String keyLocationName = getSerializedLocationName(shape.getKey(), "key");
        String valueLocationName = getSerializedLocationName(shape.getValue(), "value");
        writer.write("object := value.Map($S, $S)", keyLocationName, valueLocationName);
        writer.write("");

        // Create a sorted list of the map's keys so we can have a stable body.
        // Ideally this would be a function we dispatch to, but the lack of generics make
        // that impractical since you can't make a function for a map[string]any
        writer.write("keys := make([]string, 0, len(v))");
        writer.write("for key := range v { keys = append(keys, key) }");
        writer.addUseImports(GoDependency.standardLibraryDependency("sort", "1.14"));
        writer.write("sort.Strings(keys)");
        writer.write("");

        writer.addUseImports(SmithyGoDependency.FMT);
        writer.openBlock("for _, key := range keys {", "}", () -> {
            writer.write("om := object.Key(key)");

            // Null values should be omitted for query.
            if (!target.hasTrait(EnumTrait.class)) {
                writer.openBlock("if vv := v[key]; vv == nil {", "}", () -> {
                    writer.write("continue");
                });
            }

            target.accept(getMemberSerVisitor(shape.getValue(), "v[key]", "om"));
        });

        writer.write("return nil");
    }

    @Override
    protected void serializeStructure(GenerationContext context, StructureShape shape) {
        GoWriter writer = context.getWriter();
        writer.write("object := value.Object()");
        writer.write("_ = object");
        writer.write("");

        // Use a TreeSet to sort the members.
        Set<MemberShape> members = new TreeSet<>(shape.getAllMembers().values());
        for (MemberShape member : members) {
            if (!memberFilter.test(member)) {
                continue;
            }
            Shape target = context.getModel().expectShape(member.getTarget());
            writeSafeMemberAccessor(context, member, "v", (operand) -> {
                String locationName = getSerializedLocationName(member, member.getMemberName());
                if (isFlattened(context, member)) {
                    writer.write("ok := object.FlatKey($S)", locationName);
                } else {
                    writer.write("ok := object.Key($S)", locationName);
                }
                target.accept(getMemberSerVisitor(member, operand, "ok"));
            });
            writer.write("");
        }

        writer.write("return nil");
    }

    /**
     * Retrieves the correct serialization location based on the member's
     * xmlName trait or uses the default value.
     *
     * @param memberShape The member being serialized.
     * @param defaultValue A default value for the location.
     * @return The location where the member will be serialized.
     */
    protected String getSerializedLocationName(MemberShape memberShape, String defaultValue) {
        return memberShape.getTrait(XmlNameTrait.class)
                .map(XmlNameTrait::getValue)
                .orElse(defaultValue);
    }

    /**
     * Tells whether the contents of the member should be flattened
     * when serialized.
     *
     * @param context The generation context.
     * @param memberShape The member being serialized.
     * @return If the member's contents should be flattened when serialized.
     */
    protected boolean isFlattened(GenerationContext context, MemberShape memberShape) {
        return memberShape.hasTrait(XmlFlattenedTrait.class);
    }

    @Override
    protected void serializeUnion(GenerationContext context, UnionShape shape) {
        GoWriter writer = context.getWriter();
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol symbol = symbolProvider.toSymbol(shape);
        writer.addUseImports(SmithyGoDependency.FMT);

        writer.write("object := value.Object()");
        writer.write("");

        writer.openBlock("switch uv := v.(type) {", "}", () -> {
            // Use a TreeSet to sort the members.
            Set<MemberShape> members = new TreeSet<>(shape.getAllMembers().values());
            for (MemberShape member : members) {
                Shape target = context.getModel().expectShape(member.getTarget());
                String exportedMemberName = symbol.getName() + symbolProvider.toMemberName(member);

                writer.openBlock("case *$L:", "", exportedMemberName, () -> {
                    String locationName = getSerializedLocationName(member, member.getMemberName());
                    if (isFlattened(context, member)) {
                        writer.write("ok := object.FlatKey($S)", locationName);
                    } else {
                        writer.write("ok := object.Key($S)", locationName);
                    }
                    target.accept(getMemberSerVisitor(member, "uv.Value()", "av"));
                });
            }

            // Handle unknown union values
            writer.openBlock("case *$LUnknown:", "", symbol.getName(), () -> writer.write("fallthrough"));
            writer.openBlock("default:", "", () -> {
                writer.write("return fmt.Errorf(\"attempted to serialize unknown member type %T"
                        + " for union %T\", uv, v)");
            });
        });

        writer.write("return nil");
    }
}
