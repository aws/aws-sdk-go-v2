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

import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.CodegenUtils;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;
import software.amazon.smithy.model.shapes.BigDecimalShape;
import software.amazon.smithy.model.shapes.BigIntegerShape;
import software.amazon.smithy.model.shapes.BlobShape;
import software.amazon.smithy.model.shapes.BooleanShape;
import software.amazon.smithy.model.shapes.ByteShape;
import software.amazon.smithy.model.shapes.CollectionShape;
import software.amazon.smithy.model.shapes.DocumentShape;
import software.amazon.smithy.model.shapes.DoubleShape;
import software.amazon.smithy.model.shapes.FloatShape;
import software.amazon.smithy.model.shapes.IntegerShape;
import software.amazon.smithy.model.shapes.ListShape;
import software.amazon.smithy.model.shapes.LongShape;
import software.amazon.smithy.model.shapes.MapShape;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ResourceShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.SetShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeVisitor;
import software.amazon.smithy.model.shapes.ShortShape;
import software.amazon.smithy.model.shapes.StringShape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.shapes.TimestampShape;
import software.amazon.smithy.model.shapes.UnionShape;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait.Format;

/**
 * Visitor to generate member values for aggregate types deserialized from documents.
 */
public class JsonMemberDeserVisitor implements ShapeVisitor<Void> {
    private final GenerationContext context;
    private final String dataDest;
    private final Format timestampFormat;

    public JsonMemberDeserVisitor(GenerationContext context, String dataDest, Format timestampFormat) {
        this.context = context;
        this.dataDest = dataDest;
        this.timestampFormat = timestampFormat;
    }

    @Override
    public Void blobShape(BlobShape shape) {
        GoWriter writer = context.getWriter();
        writer.write("err := decoder.Decode(&$L)", dataDest);
        writer.write("if err != nil { return err }");
        return null;
    }

    @Override
    public Void booleanShape(BooleanShape shape) {
        GoWriter writer = context.getWriter();
        writer.addUseImports(SmithyGoDependency.FMT);
        consumeToken();
        writer.openBlock("if val != nil {", "}", () -> {
            writer.write("jtv, ok := val.(bool)");
            writer.openBlock("if !ok {", "}", () -> {
                writer.write("return fmt.Errorf(\"expected $L to be of type *bool, got %T instead\", val)",
                        shape.getId().getName());
            });
            writer.write("$L = &jtv", dataDest);
        });
        return null;
    }

    /**
     * Consumes a single token into the variable "val", returning on any error.
     */
    private void consumeToken() {
        GoWriter writer = context.getWriter();
        writer.write("val, err := decoder.Token()");
        writer.write("if err != nil { return err }");
    }

    @Override
    public Void byteShape(ByteShape shape) {
        GoWriter writer = context.getWriter();
        handleInteger(shape, CodegenUtils.generatePointerValueIfPointable(writer, shape, "int8(i64)"));
        return null;
    }

    @Override
    public Void shortShape(ShortShape shape) {
        GoWriter writer = context.getWriter();
        handleInteger(shape, CodegenUtils.generatePointerValueIfPointable(writer, shape, "int16(i64)"));
        return null;
    }

    @Override
    public Void integerShape(IntegerShape shape) {
        GoWriter writer = context.getWriter();
        handleInteger(shape, CodegenUtils.generatePointerValueIfPointable(writer, shape, "int32(i64)"));
        return null;
    }

    @Override
    public Void longShape(LongShape shape) {
        handleInteger(shape, "&i64");
        return null;
    }

    /**
     * Deserializes a number without a fractional value.
     *
     * The 64-bit integer representation of the number is stored in the variable {@code i64}.
     *
     * @param shape The shape being deserialized.
     * @param cast A wrapping of {@code i64} to cast it to the proper type.
     */
    private void handleInteger(Shape shape, String cast) {
        GoWriter writer = context.getWriter();
        handleNumber(shape, () -> {
            writer.write("i64, err := jtv.Int64()");
            writer.write("if err != nil { return err }");
            writer.write("$L = $L", dataDest, cast);
        });
    }

    /**
     * Deserializes a json number into a json token.
     *
     * The number token is stored under the variable {@code jtv}.
     *
     * @param shape The shape being deserialized.
     * @param r A runnable that runs after the value has been parsed, before the scope closes.
     */
    private void handleNumber(Shape shape, Runnable r) {
        GoWriter writer = context.getWriter();
        writer.addUseImports(SmithyGoDependency.FMT);
        consumeToken();

        writer.openBlock("if val != nil {", "}", () -> {
            writer.write("jtv, ok := val.(json.Number)");
            writer.openBlock("if !ok {", "}", () -> {
                writer.write("return fmt.Errorf(\"expected $L to be json.Number, got %T instead\", val)",
                        shape.getId().getName());
            });
            r.run();
        });
    }

    @Override
    public Void floatShape(FloatShape shape) {
        GoWriter writer = context.getWriter();
        handleFloat(shape, CodegenUtils.generatePointerValueIfPointable(writer, shape, "float32(f64)"));
        return null;
    }

    @Override
    public Void doubleShape(DoubleShape shape) {
        handleFloat(shape, "&f64");
        return null;
    }

    /**
     * Deserializes a number with a fractional value.
     *
     * The 64-bit float representation of the number is stored in the variable {@code f64}.
     *
     * @param shape The shape being deserialized.
     * @param cast A wrapping of {@code f64} to cast it to the proper type.
     */
    private void handleFloat(Shape shape, String cast) {
        GoWriter writer = context.getWriter();
        handleNumber(shape, () -> {
            writer.write("f64, err := jtv.Float64()");
            writer.write("if err != nil { return err }");
            writer.write("$L = $L", dataDest, cast);
        });
    }

    @Override
    public Void stringShape(StringShape shape) {
        GoWriter writer = context.getWriter();
        Symbol symbol = context.getSymbolProvider().toSymbol(shape);

        if (shape.hasTrait(EnumTrait.class)) {
            handleString(shape, () -> writer.write("$L = $P(jtv)", dataDest, symbol));
        } else {
            handleString(shape, () -> writer.write("$L = &jtv", dataDest));
        }

        return null;
    }

    /**
     * Deserializes a json string into a json token.
     *
     * The number token is stored under the variable {@code jtv}.
     *
     * @param shape The shape being deserialized.
     * @param r A runnable that runs after the value has been parsed, before the scope closes.
     */
    private void handleString(Shape shape, Runnable r) {
        GoWriter writer = context.getWriter();
        writer.addUseImports(SmithyGoDependency.FMT);
        consumeToken();

        writer.openBlock("if val != nil {", "}", () -> {
            writer.write("jtv, ok := val.(string)");
            writer.openBlock("if !ok {", "}", () -> {
                writer.write("return fmt.Errorf(\"expected $L to be of type string, got %T instead\", val)",
                        shape.getId().getName());
            });
            r.run();
        });
    }

    @Override
    public Void timestampShape(TimestampShape shape) {
        GoWriter writer = context.getWriter();
        writer.addUseImports(SmithyGoDependency.SMITHY_TIME);

        switch (timestampFormat) {
            case DATE_TIME:
                handleString(shape, () -> {
                    writer.write("t, err := smithytime.ParseDateTime(jtv)");
                    writer.write("if err != nil { return err }");
                    writer.write("$L = &t", dataDest);
                });
                break;
            case HTTP_DATE:
                handleString(shape, () -> {
                    writer.write("t, err := smithytime.ParseHTTPDate(jtv)");
                    writer.write("if err != nil { return err }");
                    writer.write("$L = &t", dataDest);
                });
                break;
            case EPOCH_SECONDS:
                writer.addUseImports(SmithyGoDependency.SMITHY_PTR);
                handleFloat(shape, "ptr.Time(smithytime.ParseEpochSeconds(f64))");
                break;
            default:
                throw new CodegenException(String.format("Unknown timestamp format %s", timestampFormat));
        }
        return null;
    }

    @Override
    public Void bigIntegerShape(BigIntegerShape shape) {
        // Fail instead of losing precision through Number.
        unsupportedShape(shape);
        return null;
    }

    @Override
    public Void bigDecimalShape(BigDecimalShape shape) {
        // Fail instead of losing precision through Number.
        unsupportedShape(shape);
        return null;
    }

    private String unsupportedShape(Shape shape) {
        throw new CodegenException(String.format("Cannot deserialize shape type %s on protocol, shape: %s.",
                shape.getType(), shape.getId()));
    }

    @Override
    public Void operationShape(OperationShape shape) {
        throw new CodegenException("Operation shapes cannot be bound to documents.");
    }

    @Override
    public Void resourceShape(ResourceShape shape) {
        throw new CodegenException("Resource shapes cannot be bound to documents.");
    }

    @Override
    public Void serviceShape(ServiceShape shape) {
        throw new CodegenException("Service shapes cannot be bound to documents.");
    }

    @Override
    public Void memberShape(MemberShape shape) {
        throw new CodegenException("Member shapes cannot be bound to documents.");
    }

    @Override
    public Void documentShape(DocumentShape shape) {
        writeDelegateFunction(shape);
        return null;
    }

    @Override
    public Void structureShape(StructureShape shape) {
        writeDelegateFunction(shape);
        return null;
    }

    @Override
    public Void unionShape(UnionShape shape) {
        writeDelegateFunction(shape);
        return null;
    }

    @Override
    public Void listShape(ListShape shape) {
        return collectionShape(shape);
    }

    @Override
    public Void setShape(SetShape shape) {
        return collectionShape(shape);
    }

    private Void collectionShape(CollectionShape shape) {
        writeDelegateFunction(shape);
        return null;
    }

    @Override
    public Void mapShape(MapShape shape) {
        writeDelegateFunction(shape);
        return null;
    }

    private void writeDelegateFunction(Shape shape) {
        String functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(shape, context.getProtocolName());
        GoWriter writer = context.getWriter();
        writer.openBlock("if err := $L(&$L, decoder); err != nil {", "}", functionName, dataDest, () -> {
            writer.write("return err");
        });
    }
}
