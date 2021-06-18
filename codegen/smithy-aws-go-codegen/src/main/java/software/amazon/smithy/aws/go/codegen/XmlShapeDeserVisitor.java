package software.amazon.smithy.aws.go.codegen;

import java.util.Collections;
import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.TreeSet;
import java.util.function.Predicate;
import java.util.logging.Logger;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.UnionGenerator;
import software.amazon.smithy.go.codegen.integration.DocumentShapeDeserVisitor;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.CollectionShape;
import software.amazon.smithy.model.shapes.DocumentShape;
import software.amazon.smithy.model.shapes.ListShape;
import software.amazon.smithy.model.shapes.MapShape;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.SetShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.SimpleShape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.shapes.UnionShape;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.model.traits.XmlAttributeTrait;
import software.amazon.smithy.model.traits.XmlNameTrait;
import software.amazon.smithy.utils.FunctionalUtils;

/**
 * Visitor to generate deserialization functions for shapes in XML protocol
 * document bodies.
 * <p>
 * This class handles function body generation for all types expected by the
 * {@code DocumentShapeDeserVisitor}. No other shape type serialization is overwritten.
 * <p>
 * Timestamps are serialized to {@link TimestampFormatTrait.Format}.DATE_TIME by default.
 */
public class XmlShapeDeserVisitor extends DocumentShapeDeserVisitor {
    private static final TimestampFormatTrait.Format DEFAULT_TIMESTAMP_FORMAT = TimestampFormatTrait.Format.DATE_TIME;
    private static final Logger LOGGER = Logger.getLogger(XmlShapeDeserVisitor.class.getName());

    private final Predicate<MemberShape> memberFilter;

    /**
     * @param context The generation context.
     */
    public XmlShapeDeserVisitor(GenerationContext context) {
        this(context, FunctionalUtils.alwaysTrue());
    }

    /**
     * @param context      The generation context.
     * @param memberFilter A filter that is applied to structure members. This is useful for
     *                     members that won't be in the body.
     */
    public XmlShapeDeserVisitor(GenerationContext context, Predicate<MemberShape> memberFilter) {
        super(context);
        this.memberFilter = memberFilter;
    }

    @Override
    protected Map<String, String> getAdditionalArguments() {
        return Collections.singletonMap("decoder", "smithyxml.NodeDecoder");
    }

    private XmlMemberDeserVisitor getMemberDeserVisitor(
            MemberShape member,
            String dataDest,
            boolean isXmlAttributeMember
    ) {
        // Get the timestamp format to be used, defaulting to rfc 3339 date-time format.
        TimestampFormatTrait.Format format = member.getMemberTrait(getContext().getModel(), TimestampFormatTrait.class)
                .map(TimestampFormatTrait::getFormat).orElse(DEFAULT_TIMESTAMP_FORMAT);
        return new XmlMemberDeserVisitor(getContext(), member, dataDest, format, isXmlAttributeMember);
    }

    // generates code to define and initialize output variable for an aggregate shape
    private void generatesInitializerForOutputVariable(GenerationContext context, Shape shape) {
        GoWriter writer = context.getWriter().get();
        Symbol shapeSymbol = context.getSymbolProvider().toSymbol(shape);

        writer.write("var sv $P", shapeSymbol);
        writer.openBlock("if *v == nil {", "", () -> {
            if (shape.isStructureShape()) {
                writer.write("sv = &$T{}", shapeSymbol);
            } else {
                writer.write("sv = make($P, 0)", shapeSymbol);
            }
            writer.openBlock("} else {", "}", () -> {
                writer.write("sv = *v");
            });
        });
    }

    @Override
    public Void mapShape(MapShape shape) {
        super.mapShape(shape);
        generateFlattenedMapDeserializer(getContext(), shape);
        return null;
    }

    @Override
    public Void setShape(SetShape shape) {
        super.setShape(shape);
        generateFlattenedCollectionDeserializer(getContext(), shape);
        return null;
    }

    @Override
    public Void listShape(ListShape shape) {
        super.listShape(shape);
        generateFlattenedCollectionDeserializer(getContext(), shape);
        return null;
    }

    /**
     * Deserializes the collection shapes.
     * <p>
     * In case of nested collections we will have nested `Member` element tags.
     * for eg: <ParentList><Member><ChildList><Member>abc</Member></ChildList></Member></ParentList>
     * <p>
     * The XMLNodeDecoder decodes per xml element node level and exits when it encounters an end element
     * with xml name that matches the xml name of start element.
     * <p>
     * For simple type members their is no function scoping, instead we use a loop to provide appropriate scoping.
     * This helps ensure we do not exit early when we have nested tags with same element name.
     *
     * @param context the generation context.
     * @param shape   the Collection shape to be deserialized.
     */
    @Override
    protected void deserializeCollection(GenerationContext context, CollectionShape shape) {
        GoWriter writer = context.getWriter().get();

        // initialize the output member variable
        generatesInitializerForOutputVariable(context, shape);

        writer.write("originalDecoder := decoder");
        // Iterate through the decoder. The member visitor will handle popping xml tokens
        // enclosed within a xml start and end element.
        writer.openBlock("for {", "}", () -> {
            writer.write("t, done, err := decoder.Token()");
            writer.write("if err != nil { return err }");
            writer.write("if done { break }");

            MemberShape member = shape.getMember();
            Shape target = context.getModel().expectShape(member.getTarget());
            String serializedMemberName = getSerializedMemberName(member);

            // If target is a simple shape, we must get an explicit member decoder to handle `member` element tag for
            // each member element of the list. This is not needed for the aggregate shapes as visitor handles it directly.
            if (target instanceof SimpleShape) {
                writer.write("memberDecoder := smithyxml.WrapNodeDecoder(decoder.Decoder, t)");
                writer.write("decoder = memberDecoder");
            }

            writer.addUseImports(SmithyGoDependency.STRINGS);
            writer.openBlock("switch {", "}", () -> {
                writer.openBlock("case strings.EqualFold($S, t.Name.Local):", "", serializedMemberName, () -> {
                    writer.write("var col $P", context.getSymbolProvider().toSymbol(member));
                    target.accept(getMemberDeserVisitor(member, "col", false));
                    writer.write("sv = append(sv, col)");
                });
                writer.openBlock("default:", "", () -> {
                    writer.write("err = decoder.Decoder.Skip()");
                    writer.write("if err != nil { return err }");
                });
            });
            writer.write("decoder = originalDecoder");
        });
        writer.write("*v = sv");
        writer.write("return nil");
    }

    // Generates deserializer function for collection shapes with xml flattened trait.
    public void generateFlattenedCollectionDeserializer(GenerationContext context, CollectionShape shape) {
        GoWriter writer = context.getWriter().get();
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol symbol = symbolProvider.toSymbol(shape);

        MemberShape member = shape.getMember();
        Symbol memberSymbol = symbolProvider.toSymbol(member);
        Shape target = context.getModel().expectShape(member.getTarget());

        writer.openBlock("func $L(v *$P, decoder smithyxml.NodeDecoder) error {", "}",
                getUnwrappedMapDelegateFunctionName(context, shape), symbol, () -> {
                    // initialize the output member variable
                    generatesInitializerForOutputVariable(context, shape);
                    writer.openBlock(" switch { default: ", "}", () -> {
                        writer.write("var mv $P", memberSymbol);
                        writer.write("t := decoder.StartEl");
                        writer.write("_ = t");
                        target.accept(getMemberDeserVisitor(member, "mv", false));
                        writer.write("sv = append(sv, mv)");
                    });

                    writer.write("*v = sv");
                    writer.write("return nil");
                });
    }

    @Override
    protected void deserializeMap(GenerationContext context, MapShape shape) {
        GoWriter writer = context.getWriter().get();

        // initialize the output member variable
        generatesInitializerForOutputVariable(context, shape);

        // Iterate through the decoder. The member visitor will handle popping xml tokens
        // enclosed within a xml start and end element.
        writer.openBlock("for {", "}", () -> {
            writer.write("t, done, err := decoder.Token()");
            writer.write("if err != nil { return err }");
            writer.write("if done { break }");

            // non-flattened maps
            writer.addUseImports(SmithyGoDependency.STRINGS);
            writer.openBlock("switch {", "}", () -> {
                writer.openBlock("case strings.EqualFold(\"entry\", t.Name.Local):", "", () -> {
                    writer.write("entryDecoder := smithyxml.WrapNodeDecoder(decoder.Decoder, t)");
                    // delegate to unwrapped map deserializer function
                    writer.openBlock("if err := $L(&sv, entryDecoder); err != nil {", "}",
                            getUnwrappedMapDelegateFunctionName(context, shape), () -> {
                                writer.write("return err");
                            });
                });
                writer.openBlock("default:", "", () -> {
                    writer.write("err = decoder.Decoder.Skip()");
                    writer.write("if err != nil { return err }");
                });
            });
        });

        writer.write("*v = sv");
        writer.write("return nil");
    }

    // Generates deserializer function for flattened maps.
    protected void generateFlattenedMapDeserializer(GenerationContext context, MapShape shape) {
        GoWriter writer = context.getWriter().get();
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol symbol = symbolProvider.toSymbol(shape);

        writer.addUseImports(SmithyGoDependency.SMITHY_XML);
        writer.openBlock("func $L(v *$P, decoder smithyxml.NodeDecoder) error {", "}",
                getUnwrappedMapDelegateFunctionName(context, shape), symbol, () -> {
                    // initialize the output member variable
                    generatesInitializerForOutputVariable(context, shape);
                    MemberShape valueShape = shape.getValue();
                    MemberShape keyShape = shape.getKey();

                    Symbol keySymbol = context.getSymbolProvider().toSymbol(keyShape);
                    Symbol valueSymbol = context.getSymbolProvider().toSymbol(valueShape);

                    Shape targetKey = context.getModel().expectShape(keyShape.getTarget());

                    writer.write("var ek $P", keySymbol);
                    writer.write("var ev $P", valueSymbol);
                    writer.insertTrailingNewline();

                    // Iterate through the decoder. The member visitor will handle popping xml tokens
                    // enclosed within a xml start and end element.
                    writer.openBlock("for {", "}", () -> {
                        writer.write("t, done, err := decoder.Token()");
                        writer.write("if err != nil { return err }");
                        writer.openBlock("if done {", "}", () -> {
                            // set the key value pair in map
                            if (keyShape.hasTrait(EnumTrait.class) || targetKey.hasTrait(EnumTrait.class)) {
                                writer.write("sv[string(ek)] = ev");
                            } else {
                                writer.write("sv[ek] = ev");
                            }
                            writer.write("break");
                        });

                        writer.write("originalDecoder := decoder");
                        writer.write("decoder = smithyxml.WrapNodeDecoder(originalDecoder.Decoder, t)");
                        writer.insertTrailingNewline();

                        writer.openBlock("switch {", "}", () -> {
                            writer.addUseImports(SmithyGoDependency.STRINGS);
                            writer.openBlock("case strings.EqualFold($S, t.Name.Local):", "", getSerializedMemberName(keyShape), () -> {
                                String dest = "ek";
                                context.getModel().expectShape(keyShape.getTarget()).accept(
                                        getMemberDeserVisitor(keyShape, dest, false));
                            });

                            writer.openBlock("case strings.EqualFold($S, t.Name.Local):", "", getSerializedMemberName(valueShape), () -> {
                                String dest = "ev";
                                context.getModel().expectShape(valueShape.getTarget()).accept(
                                        getMemberDeserVisitor(valueShape, dest, false));
                            });

                            writer.openBlock("default:", "", () -> {
                                writer.write("err = decoder.Decoder.Skip()");
                                writer.write("if err != nil { return err }");
                            });
                        });
                        writer.write("decoder = originalDecoder");
                    });
                    writer.write("*v = sv");
                    writer.write("return nil");
                });
    }

    private String getUnwrappedMapDelegateFunctionName(GenerationContext context, Shape shape) {
        return ProtocolGenerator.getDocumentDeserializerFunctionName(shape, context.getService(), context.getProtocolName()) + "Unwrapped";
    }

    @Override
    protected void deserializeStructure(GenerationContext context, StructureShape shape) {
        GoWriter writer = context.getWriter().get();
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Model model = context.getModel();

        // initialize the output member variable
        generatesInitializerForOutputVariable(context, shape);
        // Deserialize member shapes modeled with xml attribute trait
        if (hasXmlAttributeTraitMember(shape)) {
            writer.openBlock("for _, attr := range decoder.StartEl.Attr {", "}", () -> {
                writer.write("name := attr.Name.Local");
                writer.openBlock("if len(attr.Name.Space) != 0 {", "}", () -> {
                    writer.addUseImports(SmithyGoDependency.STRINGS);
                    writer.write("name = attr.Name.Space + `:` + attr.Name.Local");
                });
                writer.openBlock("switch {", "}", () -> {
                    Set<MemberShape> members = new TreeSet<>(shape.members());
                    for (MemberShape member : members) {
                        // check if member does not conform with the member filter or does not have a xmlAttribute trait
                        if (!memberFilter.test(member) || !member.hasTrait(XmlAttributeTrait.ID)) {
                            continue;
                        }

                        String memberName = symbolProvider.toMemberName(member);
                        String serializedMemberName = getSerializedMemberName(member);
                        writer.addUseImports(SmithyGoDependency.STRINGS);
                        writer.openBlock("case strings.EqualFold($S, name):", "", serializedMemberName, () -> {
                            String dest = "sv." + memberName;
                            context.getModel().expectShape(member.getTarget()).accept(
                                    getMemberDeserVisitor(member, dest, true));
                        });
                    }
                });
            });
        }

        // Iterate through the decoder. The member visitor will handle popping xml tokens
        // enclosed within a xml start and end element.
        writer.openBlock("for {", "}", () -> {
            writer.write("t, done, err := decoder.Token()");
            writer.write("if err != nil { return err }");
            writer.write("if done { break }");

            // Create a new decoder for each member
            writer.write("originalDecoder := decoder");
            writer.write("decoder = smithyxml.WrapNodeDecoder(originalDecoder.Decoder, t)");
            writer.insertTrailingNewline();

            writer.openBlock("switch {", "}", () -> {
                Set<MemberShape> members = new TreeSet<>(shape.members());
                for (MemberShape member : members) {
                    // check if member is not a document binding or has a xmlAttribute trait
                    if (!memberFilter.test(member) || member.hasTrait(XmlAttributeTrait.ID)) {
                        continue;
                    }
                    String memberName = symbolProvider.toMemberName(member);
                    String serializedMemberName = getSerializedMemberName(member);
                    writer.addUseImports(SmithyGoDependency.STRINGS);
                    writer.openBlock("case strings.EqualFold($S, t.Name.Local):", "", serializedMemberName, () -> {
                        String dest = "sv." + memberName;
                        model.expectShape(member.getTarget()).accept(
                                getMemberDeserVisitor(member, dest, false));
                    });
                }

                writer.openBlock("default:", "", () -> {
                    writer.writeDocs("Do nothing and ignore the unexpected tag element");
                    writer.write("err = decoder.Decoder.Skip()");
                    writer.write("if err != nil { return err }");
                });
            });
            // re-assign the  original decoder
            writer.write("decoder = originalDecoder");
        });

        writer.write("*v = sv");
        writer.write("return nil");
    }

    // return true if any member of the shape is decorated with XmlAttributeTrait
    private boolean hasXmlAttributeTraitMember(Shape shape) {
        for (MemberShape member : shape.members()) {
            if (member.hasTrait(XmlAttributeTrait.ID)) {
                return true;
            }
        }
        return false;
    }

    private String getSerializedMemberName(MemberShape memberShape) {
        Optional<XmlNameTrait> xmlNameTrait = memberShape.getTrait(XmlNameTrait.class);
        return xmlNameTrait.isPresent() ? xmlNameTrait.get().getValue() : memberShape.getMemberName();
    }

    @Override
    protected void deserializeDocument(GenerationContext context, DocumentShape shape) {
        GoWriter writer = context.getWriter().get();
        LOGGER.warning("The document type is unsupported for XML protocols.");
        writer.addUseImports(SmithyGoDependency.SMITHY);
        writer.write("return &smithy.DeserializationError{Err: fmt.Errorf("
                + "\"Document type is unsupported for XML protocols.\")}");
    }

    @Override
    protected void deserializeUnion(GenerationContext context, UnionShape shape) {
        GoWriter writer = context.getWriter().get();
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol symbol = symbolProvider.toSymbol(shape);
        Model model = context.getModel();

        writer.write("var uv $P", symbol);
        writer.write("var memberFound bool");

        // Iterate through the decoder. The member visitor will handle popping
        // xml tokens enclosed within a xml start and end element.
        writer.openBlock("for {", "}", () -> {
            writer.write("t, done, err := decoder.Token()");
            writer.write("if err != nil { return err }");
            writer.write("if done { break }");
            writer.openBlock("if memberFound {", "}", () -> {
                writer.write("if err = decoder.Decoder.Skip(); err != nil { return err }");
            });

            // Create a new decoder for each member
            writer.write("originalDecoder := decoder");
            writer.write("decoder = smithyxml.WrapNodeDecoder(originalDecoder.Decoder, t)");
            writer.insertTrailingNewline();

            writer.openBlock("switch {", "}", () -> {
                Set<MemberShape> members = new TreeSet<>(shape.members());
                for (MemberShape member : members) {
                    // check if member is not a document binding or has a xmlAttribute trait
                    if (!memberFilter.test(member) || member.hasTrait(XmlAttributeTrait.ID)) {
                        continue;
                    }
                    Symbol targetSymbol = symbolProvider.toSymbol(member);
                    Symbol memberSymbol = SymbolUtils.createValueSymbolBuilder(
                            symbolProvider.toMemberName(member),
                            symbol.getNamespace()
                    ).build();

                    writer.addUseImports(SmithyGoDependency.STRINGS);

                    String serializedMemberName = getSerializedMemberName(member);
                    writer.openBlock("case strings.EqualFold($S, t.Name.Local):", "", serializedMemberName, () -> {
                        writer.write("var mv $P", targetSymbol);
                        model.expectShape(member.getTarget()).accept(getMemberDeserVisitor(member, "mv", false));
                        writer.write("uv = &$T{Value: mv}", memberSymbol);
                        writer.write("memberFound = true");
                    });
                }

                writer.openBlock("default:", "", () -> {
                    // This is the function to take a value and convert it to the union type.
                    Symbol unknownMemberSymbol = SymbolUtils.createValueSymbolBuilder(
                            UnionGenerator.UNKNOWN_MEMBER_NAME,
                            symbol.getNamespace()
                    ).build();
                    writer.write("uv = &$T{Tag: t.Name.Local}", unknownMemberSymbol);
                    writer.write("memberFound = true");
                });
            });
            // re-assign the  original decoder
            writer.write("decoder = originalDecoder");
        });

        writer.write("*v = uv");
        writer.write("return nil");
    }
}
