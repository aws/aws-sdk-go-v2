/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *  http://aws.amazon.com/apache2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

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
import software.amazon.smithy.go.codegen.GoValueAccessUtils;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.ProtocolDocumentGenerator;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.DocumentShapeSerVisitor;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;
import software.amazon.smithy.go.codegen.knowledge.GoPointableIndex;
import software.amazon.smithy.go.codegen.trait.NoSerializeTrait;
import software.amazon.smithy.model.knowledge.NullableIndex;
import software.amazon.smithy.model.shapes.CollectionShape;
import software.amazon.smithy.model.shapes.DocumentShape;
import software.amazon.smithy.model.shapes.MapShape;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.shapes.UnionShape;
import software.amazon.smithy.model.traits.JsonNameTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait.Format;
import software.amazon.smithy.utils.FunctionalUtils;
import software.amazon.smithy.utils.SmithyBuilder;

/**
 * Visitor to generate serialization functions for shapes in AWS JSON protocol
 * document bodies.
 * <p>
 * This class handles function body generation for all types expected by the
 * {@code DocumentShapeSerVisitor}. No other shape type serialization is overwritten.
 * <p>
 * Timestamps are serialized to {@link Format}.EPOCH_SECONDS by default.
 */
final class JsonShapeSerVisitor extends DocumentShapeSerVisitor {
    private static final Format DEFAULT_TIMESTAMP_FORMAT = Format.EPOCH_SECONDS;
    private static final Logger LOGGER = Logger.getLogger(JsonShapeSerVisitor.class.getName());

    private final Predicate<MemberShape> memberFilter;
    private final GoPointableIndex pointableIndex;
    private final NullableIndex nullableIndex;
    private final boolean supportJsonName;

    /**
     * Returns a new builder for building the JsonShapeSerVisitor.
     * @return Builder
     */
    public static Builder builder() {
        return new Builder();
    }

    protected JsonShapeSerVisitor(Builder builder) {
        super(SmithyBuilder.requiredState("context", builder.context), builder.serNameProvider);
        if (builder.memberFilter != null) {
            this.memberFilter = builder.memberFilter;
        } else {
            this.memberFilter = NoSerializeTrait.excludeNoSerializeMembers().and(
                    FunctionalUtils.alwaysTrue());
        }
        this.supportJsonName = builder.supportJsonName;
        this.pointableIndex = GoPointableIndex.of(this.getContext().getModel());
        this.nullableIndex = NullableIndex.of(this.getContext().getModel());
    }

    private DocumentMemberSerVisitor getMemberSerVisitor(MemberShape member, String source, String dest) {
        // Get the timestamp format to be used, defaulting to epoch seconds.
        Format format = member.getMemberTrait(getContext().getModel(), TimestampFormatTrait.class)
                .map(TimestampFormatTrait::getFormat).orElse(DEFAULT_TIMESTAMP_FORMAT);
        return new DocumentMemberSerVisitor(getContext(), member, source, dest, format);
    }

    @Override
    protected Map<String, String> getAdditionalSerArguments() {
        return Collections.singletonMap("value", "smithyjson.Value");
    }

    @Override
    protected void serializeCollection(GenerationContext context, CollectionShape shape) {
        GoWriter writer = context.getWriter().get();
        Shape target = context.getModel().expectShape(shape.getMember().getTarget());
        writer.addUseImports(SmithyGoDependency.SMITHY_JSON);
        writer.write("array := value.Array()");
        writer.write("defer array.Close()");
        writer.write("");

        writer.openBlock("for i := range v {", "}", () -> {
            writer.write("av := array.Value()");

            // Null values in lists should be serialized as such. Enums can't be null, so we don't bother
            // putting this in for their case.
            if (pointableIndex.isNillable(shape.getMember())) {
                writer.openBlock("if vv := v[i]; vv == nil {", "}", () -> {
                    if (nullableIndex.isNullable(shape.getMember())) {
                        writer.write("av.Null()");
                    }
                    writer.write("continue");
                });
            }

            target.accept(getMemberSerVisitor(shape.getMember(), "v[i]", "av"));
        });

        writer.write("return nil");
    }

    @Override
    protected void serializeDocument(GenerationContext context, DocumentShape shape) {
        GoWriter writer = context.getWriter().get();

        Symbol isInterface = ProtocolDocumentGenerator.Utilities.getInternalDocumentSymbolBuilder(context.getSettings(),
                ProtocolDocumentGenerator.INTERNAL_IS_DOCUMENT_INTERFACE).build();

        writer.openBlock("if v == nil {", "}", () -> writer.write("return nil"));

        writer.openBlock("if !$T(v) {", "}", isInterface, () -> {
            writer.addUseImports(SmithyGoDependency.FMT);
            writer.write("return fmt.Errorf(\"%T is not a compatible document type\", v)");
        });

        writer.write("db, err := v.MarshalSmithyDocument()");
        writer.openBlock("if err != nil {", "}", () -> writer.write("return err"));
        writer.write("value.Write(db)");
        writer.write("return nil");
    }

    @Override
    protected void serializeMap(GenerationContext context, MapShape shape) {
        GoWriter writer = context.getWriter().get();
        Shape target = context.getModel().expectShape(shape.getValue().getTarget());
        writer.addUseImports(SmithyGoDependency.SMITHY_JSON);

        writer.write("object := value.Object()");
        writer.write("defer object.Close()");
        writer.write("");

        writer.openBlock("for key := range v {", "}", () -> {
            writer.write("om := object.Key(key)");

            // Null values in maps should be serialized as such. Enums can't be null, so we don't bother
            // putting this in for their case.
            if (pointableIndex.isNillable(shape.getValue())) {
                writer.openBlock("if vv := v[key]; vv == nil {", "}", () -> {
                    if (nullableIndex.isNullable(shape.getValue())) {
                        writer.write("om.Null()");
                    }
                    writer.write("continue");
                });
            }

            target.accept(getMemberSerVisitor(shape.getValue(), "v[key]", "om"));
        });

        writer.write("return nil");
    }

    @Override
    protected void serializeStructure(GenerationContext context, StructureShape shape) {
        GoWriter writer = context.getWriter().get();
        writer.addUseImports(SmithyGoDependency.SMITHY_JSON);

        writer.write("object := value.Object()");
        writer.write("defer object.Close()");
        writer.write("");

        // Use a TreeSet to sort the members.
        Set<MemberShape> members = new TreeSet<>(shape.getAllMembers().values());
        for (MemberShape member : members) {
            if (!memberFilter.test(member)) {
                continue;
            }
            Shape target = context.getModel().expectShape(member.getTarget());
            String serializedMemberName = getSerializedMemberName(member);

            GoValueAccessUtils.writeIfNonZeroValueMember(context.getModel(), context.getSymbolProvider(), writer,
                    member, "v", true, member.isRequired(), (operand) -> {
                        writer.write("ok := object.Key($S)", serializedMemberName);
                        target.accept(getMemberSerVisitor(member, operand, "ok"));
                    });
            writer.write("");
        }

        writer.write("return nil");
    }

    private String getSerializedMemberName(MemberShape memberShape) {
        if (this.supportJsonName) {
            Optional<JsonNameTrait> jsonNameTrait = memberShape.getTrait(JsonNameTrait.class);
            return jsonNameTrait.isPresent() ? jsonNameTrait.get().getValue() : memberShape.getMemberName();
        }
        return memberShape.getMemberName();
    }

    @Override
    protected void serializeUnion(GenerationContext context, UnionShape shape) {
        GoWriter writer = context.getWriter().get();
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol symbol = symbolProvider.toSymbol(shape);
        writer.addUseImports(SmithyGoDependency.FMT);
        writer.addUseImports(SmithyGoDependency.SMITHY_JSON);

        writer.write("object := value.Object()");
        writer.write("defer object.Close()");
        writer.write("");

        writer.openBlock("switch uv := v.(type) {", "}", () -> {
            // Use a TreeSet to sort the members.
            Set<MemberShape> members = new TreeSet<>(shape.getAllMembers().values());
            for (MemberShape member : members) {
                Shape target = context.getModel().expectShape(member.getTarget());
                Symbol memberSymbol = SymbolUtils.createValueSymbolBuilder(
                        symbolProvider.toMemberName(member),
                        symbol.getNamespace()
                ).build();
                String serializedMemberName = getSerializedMemberName(member);

                writer.openBlock("case *$T:", "", memberSymbol, () -> {
                    writer.write("av := object.Key($S)", serializedMemberName);
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

    public static final class Builder implements SmithyBuilder<JsonShapeSerVisitor> {
        private GenerationContext context;
        private Predicate<MemberShape> memberFilter;
        private SerializerNameProvider serNameProvider;
        private boolean supportJsonName;

        private Builder() {}

        public Builder context(GenerationContext context) {
            this.context = context;
            return this;
        }

        public Builder memberFilter(Predicate<MemberShape> filter) {
            this.memberFilter = filter;
            return this;
        }

        public Builder serializerNameProvider(SerializerNameProvider nameProvider) {
            this.serNameProvider = nameProvider;
            return this;
        }

        public Builder supportJsonName(boolean v) {
            this.supportJsonName = v;
            return this;
        }

        @Override
        public JsonShapeSerVisitor build() {
            return new JsonShapeSerVisitor(this);
        }
    }
}
