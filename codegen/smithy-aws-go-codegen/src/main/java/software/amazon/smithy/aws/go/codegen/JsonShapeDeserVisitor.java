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
import java.util.logging.Logger;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.DocumentShapeDeserVisitor;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;
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

/**
 * Visitor to generate deserialization functions for shapes in AWS JSON protocol
 * document bodies.
 *
 * This class handles function body generation for all types expected by the
 * {@code DocumentShapeDeserVisitor}. No other shape type serialization is overwritten.
 *
 * Timestamps are serialized to {@link Format}.EPOCH_SECONDS by default.
 */
public class JsonShapeDeserVisitor extends DocumentShapeDeserVisitor {
    private static final Format DEFAULT_TIMESTAMP_FORMAT = Format.EPOCH_SECONDS;
    private static final Logger LOGGER = Logger.getLogger(JsonShapeDeserVisitor.class.getName());

    public JsonShapeDeserVisitor(GenerationContext context) {
        super(context);
    }

    private JsonMemberDeserVisitor getMemberDeserVisitor(MemberShape member, String dataDest) {
        // Get the timestamp format to be used, defaulting to epoch seconds.
        Format format = member.getMemberTrait(getContext().getModel(), TimestampFormatTrait.class)
                .map(TimestampFormatTrait::getFormat).orElse(DEFAULT_TIMESTAMP_FORMAT);
        return new JsonMemberDeserVisitor(getContext(), dataDest, format);
    }

    @Override
    protected Map<String, String> getAdditionalArguments() {
        return Collections.singletonMap("decoder", "*json.Decoder");
    }

    @Override
    protected void deserializeCollection(GenerationContext context, CollectionShape shape) {
        GoWriter writer = context.getWriter();
        MemberShape member = shape.getMember();
        Shape target = context.getModel().expectShape(member.getTarget());
        writeJsonTokenizerStartStub(writer, shape);
        Symbol symbol = context.getSymbolProvider().toSymbol(shape);

        // Initialize the value now that the start stub has verified that there's something there.
        writer.write("var cv $P", symbol);
        writer.openBlock("if *v == nil {", "", () -> {
            writer.write("cv = $P{}", symbol);
            writer.openBlock("} else {", "}", () -> {
                writer.write("cv = *v");
            });
        });

        // Iterate through the decoder. The member visitor will handle popping tokens.
        writer.openBlock("for decoder.More() {", "}", () -> {
            // We need to write out an intermediate variable to assign the value of the
            // member to so that we can use it in the append function later.
            writer.write("var col $P", context.getSymbolProvider().toSymbol(target));
            target.accept(getMemberDeserVisitor(member, "col"));
            writer.write("cv = append(cv, col)");
            writer.write("");
        });

        writeJsonTokenizerEndStub(writer, shape);
        writer.write("*v = cv");
        writer.write("return nil");
    }

    @Override
    protected void deserializeDocument(GenerationContext context, DocumentShape shape) {
        GoWriter writer = context.getWriter();
        // TODO: implement document deserialization
        LOGGER.warning("Document type is currently unsupported for JSON serialization.");
        context.getWriter().writeDocs("TODO: implement document serialization.");
        writer.write("return nil");
    }

    @Override
    protected void deserializeMap(GenerationContext context, MapShape shape) {
        GoWriter writer = context.getWriter();
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol symbol = symbolProvider.toSymbol(shape);
        MemberShape member = shape.getValue();
        Symbol targetSymbol = symbolProvider.toSymbol(member);
        writeJsonTokenizerStartStub(writer, shape);

        // Initialize the value now that the start stub has verified that there's something there.
        writer.write("var mv $P", symbol);
        writer.openBlock("if *v == nil {", "", () -> {
            writer.write("mv = $P{}", symbol);
            writer.openBlock("} else {", "}", () -> {
                writer.write("mv = *v");
            });
        });

        // Iterate through the decoder. The member visitor will handle popping tokens.
        writer.openBlock("for decoder.More() {", "}", () -> {
            // Deserialize the key
            writer.write("token, err := decoder.Token()");
            writer.write("if err != nil { return err }");
            writer.write("");
            writer.write("key, ok := token.(string)");
            writer.write("if !ok { return fmt.Errorf(\"expected map-key of type string, found type %T\", token)}");
            writer.write("");

            // Deserialize the value. We need to write out an intermediate variable here
            // since we can't just pass in &mv[key]
            writer.write("var parsedVal $P", targetSymbol);
            context.getModel().expectShape(member.getTarget()).accept(getMemberDeserVisitor(member, "parsedVal"));
            writer.write("mv[key] = parsedVal");
            writer.write("");
        });

        writeJsonTokenizerEndStub(writer, shape);
        writer.write("*v = mv");
        writer.write("return nil");
    }

    @Override
    protected void deserializeStructure(GenerationContext context, StructureShape shape) {
        GoWriter writer = context.getWriter();
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol symbol = symbolProvider.toSymbol(shape);

        writeJsonTokenizerStartStub(writer, shape);

        // Initialize the value now that the start stub has verified that there's something there.
        writer.write("var sv $P", symbol);
        writer.openBlock("if *v == nil {", "", () -> {
            writer.write("sv = &$T{}", symbol);
            writer.openBlock("} else {", "}", () -> {
                writer.write("sv = *v");
            });
        });

        // Iterate through the decoder. The member visitor will handle popping tokens.
        writer.openBlock("for decoder.More() {", "}", () -> {
            writer.write("t, err := decoder.Token()");
            writer.write("if err != nil { return err }");
            writer.openBlock("switch t {", "}", () -> {
                Set<MemberShape> members = new TreeSet<>(shape.members());
                for (MemberShape member : members) {
                    String memberName = symbolProvider.toMemberName(member);
                    String serializedMemberName = getSerializedMemberName(member);
                    writer.openBlock("case $S:", "", serializedMemberName, () -> {
                        String dest = "sv." + memberName;
                        context.getModel().expectShape(member.getTarget()).accept(getMemberDeserVisitor(member, dest));
                    });
                }

                writer.openBlock("default:", "", () -> {
                    // TODO: do I need to import from restjson here?
                    writer.addUseImports(AwsGoDependency.AWS_REST_JSON_PROTOCOL);
                    writer.write("err := restjson.DiscardUnknownField(decoder)");
                    writer.write("if err != nil {return err}");
                });
            });
        });

        writeJsonTokenizerEndStub(writer, shape);
        writer.write("*v = sv");
        writer.write("return nil");
    }

    @Override
    protected void deserializeUnion(GenerationContext context, UnionShape shape) {
        GoWriter writer = context.getWriter();
        writer.write("return nil");
    }

    private String getSerializedMemberName(MemberShape memberShape) {
        Optional<JsonNameTrait> jsonNameTrait = memberShape.getTrait(JsonNameTrait.class);
        return jsonNameTrait.isPresent() ? jsonNameTrait.get().getValue() : memberShape.getMemberName();
    }

    /**
     * Writes out a stub to initialize decoding.
     *
     * @param writer The GoWriter to use.
     * @param shape The shape the stub is intended to start parsing.
     */
    private void writeJsonTokenizerStartStub(GoWriter writer, Shape shape) {
        writer.addUseImports(SmithyGoDependency.JSON);
        String startToken = shape instanceof CollectionShape ? "[" : "{";
        writer.write("startToken, err := decoder.Token()");
        writer.write("if err == io.EOF { return nil }");
        writer.write("if err != nil { return err }");
        writer.write("if startToken == nil { return nil }");
        writer.openBlock("if t, ok := startToken.(json.Delim); !ok || t != '$L' {","}", startToken, () -> {
            writer.addUseImports(SmithyGoDependency.FMT);
            writer.write("return fmt.Errorf(\"expect `$L` as start token\")", startToken);
        });
        writer.write("");

    }

    /**
     * Writes out a stub to finalize decoding.
     *
     * @param writer The GoWriter to use.
     * @param shape The shape the stub is intended to finalize parsing for.
     */
    private void writeJsonTokenizerEndStub(GoWriter writer, Shape shape) {
        String endToken = shape instanceof CollectionShape ? "]" : "}";
        writer.write("endToken, err := decoder.Token()");
        writer.write("if err != nil { return err }");
        writer.openBlock("if t, ok := endToken.(json.Delim); !ok || t != '$L' {", "}", endToken, () -> {
            writer.write("return fmt.Errorf(\"expect `$L` as end token\")", endToken);
        });
        writer.write("");
    }
}
