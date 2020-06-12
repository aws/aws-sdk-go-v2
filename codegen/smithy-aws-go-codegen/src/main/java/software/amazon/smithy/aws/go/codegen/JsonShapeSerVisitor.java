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

import static software.amazon.smithy.aws.go.codegen.AwsProtocolUtils.writeSafeMemberAccessor;

import java.util.Collections;
import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.TreeSet;
import java.util.logging.Logger;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
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
import software.amazon.smithy.model.traits.JsonNameTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait.Format;

/**
 * Visitor to generate serialization functions for shapes in AWS JSON protocol
 * document bodies.
 *
 * This class handles function body generation for all types expected by the
 * {@code DocumentShapeSerVisitor}. No other shape type serialization is overwritten.
 *
 * Timestamps are serialized to {@link Format}.EPOCH_SECONDS by default.
 */
final class JsonShapeSerVisitor extends DocumentShapeSerVisitor {
    private static final Format DEFAULT_TIMESTAMP_FORMAT = Format.EPOCH_SECONDS;
    private static final Logger LOGGER = Logger.getLogger(JsonShapeSerVisitor.class.getName());

    public JsonShapeSerVisitor(GenerationContext context) {
        super(context);
    }

    private JsonMemberSerVisitor getMemberSerVisitor(MemberShape member, String source, String dest) {
        // Get the timestamp format to be used, defaulting to epoch seconds.
        Format format = member.getMemberTrait(getContext().getModel(), TimestampFormatTrait.class)
                .map(TimestampFormatTrait::getFormat).orElse(DEFAULT_TIMESTAMP_FORMAT);
        return new JsonMemberSerVisitor(getContext(), source, dest, format);
    }

    @Override
    protected Map<String, String> getAdditionalSerArguments() {
        return Collections.singletonMap("value", "smithyjson.Value");
    }

    @Override
    protected void serializeCollection(GenerationContext context, CollectionShape shape) {
        GoWriter writer = context.getWriter();
        Shape target = context.getModel().expectShape(shape.getMember().getTarget());

        writer.write("array := value.Array()");
        writer.write("defer array.Close()");
        writer.write("");

        writer.openBlock("for i := range v {", "}", () -> {
            writer.write("av := array.Value()");

            // Null values in lists should be serialized as such. Enums can't be null, so we don't bother
            // putting this in for their case.
            if (!target.hasTrait(EnumTrait.class)) {
                writer.openBlock("if vv := v[i]; vv == nil {", "}", () -> {
                    writer.write("av.Null()");
                    writer.write("continue");
                });
            }

            target.accept(getMemberSerVisitor(shape.getMember(), "v[i]", "av"));
        });

        writer.write("return nil");
    }

    @Override
    protected void serializeDocument(GenerationContext context, DocumentShape shape) {
        // TODO: implement document serialization
        LOGGER.warning("Document type is currently unsupported for JSON serialization.");
        context.getWriter().writeDocs("TODO: implement document serialization.");
        context.getWriter().write("return nil");
    }

    @Override
    protected void serializeMap(GenerationContext context, MapShape shape) {
        GoWriter writer = context.getWriter();
        Shape target = context.getModel().expectShape(shape.getValue().getTarget());

        writer.write("object := value.Object()");
        writer.write("defer object.Close()");
        writer.write("");

        writer.openBlock("for key := range v {", "}", () -> {
            writer.write("om := object.Key(key)");

            // Null values in maps should be serialized as such. Enums can't be null, so we don't bother
            // putting this in for their case.
            if (!target.hasTrait(EnumTrait.class)) {
                writer.openBlock("if vv := v[key]; vv == nil {", "}", () -> {
                    writer.write("om.Null()");
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
        writer.write("defer object.Close()");
        writer.write("");

        // Use a TreeSet to sort the members.
        Set<MemberShape> members = new TreeSet<>(shape.getAllMembers().values());
        for (MemberShape member : members) {
            Shape target = context.getModel().expectShape(member.getTarget());
            String serializedMemberName = getSerializedMemberName(member);
            writeSafeMemberAccessor(context, member, "v", (operand) -> {
                writer.write("ok := object.Key($S)", serializedMemberName);
                target.accept(getMemberSerVisitor(member, operand, "ok"));
            });
            writer.write("");
        }

        writer.write("return nil");
    }

    private String getSerializedMemberName(MemberShape memberShape) {
        Optional<JsonNameTrait> jsonNameTrait = memberShape.getTrait(JsonNameTrait.class);
        return jsonNameTrait.isPresent() ? jsonNameTrait.get().getValue() : memberShape.getMemberName();
    }

    @Override
    protected void serializeUnion(GenerationContext context, UnionShape shape) {
        GoWriter writer = context.getWriter();
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol symbol = symbolProvider.toSymbol(shape);
        writer.addUseImports(SmithyGoDependency.FMT);

        writer.write("object := value.Object()");
        writer.write("defer object.Close()");
        writer.write("");

        writer.openBlock("switch uv := v.(type) {", "}", () -> {
            // Use a TreeSet to sort the members.
            Set<MemberShape> members = new TreeSet<>(shape.getAllMembers().values());
            for (MemberShape member : members) {
                Shape target = context.getModel().expectShape(member.getTarget());
                Symbol memberSymbol = symbolProvider.toSymbol(member);
                String exportedMemberName = symbol.getName() + symbolProvider.toMemberName(member);

                writer.openBlock("case *$L:", "", exportedMemberName, () -> {
                    writer.write("av := object.Key($S)", memberSymbol.getName());
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
