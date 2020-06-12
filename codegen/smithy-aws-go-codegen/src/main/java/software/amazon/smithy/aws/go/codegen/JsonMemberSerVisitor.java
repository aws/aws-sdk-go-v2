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
 * Visitor to generate member values for aggregate types serialized in documents.
 */
public class JsonMemberSerVisitor implements ShapeVisitor<Void> {
    private final GenerationContext context;
    private final String dataSource;
    private final String dataDest;
    private final Format timestampFormat;

    public JsonMemberSerVisitor(
            GenerationContext context,
            String dataSource,
            String dataDest,
            Format timestampFormat
    ) {
        this.context = context;
        this.dataSource = dataSource;
        this.dataDest = dataDest;
        this.timestampFormat = timestampFormat;
    }

    @Override
    public Void blobShape(BlobShape shape) {
        String source = conditionallyDereference(shape, dataSource);
        context.getWriter().write("$L.Base64EncodeBytes($L)", dataDest, source);
        return null;
    }

    @Override
    public Void booleanShape(BooleanShape shape) {
        String source = conditionallyDereference(shape, dataSource);
        context.getWriter().write("$L.Boolean($L)", dataDest, source);
        return null;
    }

    @Override
    public Void byteShape(ByteShape shape) {
        String source = conditionallyDereference(shape, dataSource);
        context.getWriter().write("$L.Byte($L)", dataDest, source);
        return null;
    }

    @Override
    public Void shortShape(ShortShape shape) {
        String source = conditionallyDereference(shape, dataSource);
        context.getWriter().write("$L.Short($L)", dataDest, source);
        return null;
    }

    @Override
    public Void integerShape(IntegerShape shape) {
        String source = conditionallyDereference(shape, dataSource);
        context.getWriter().write("$L.Integer($L)", dataDest, source);
        return null;
    }

    @Override
    public Void longShape(LongShape shape) {
        String source = conditionallyDereference(shape, dataSource);
        context.getWriter().write("$L.Long($L)", dataDest, source);
        return null;
    }

    @Override
    public Void floatShape(FloatShape shape) {
        String source = conditionallyDereference(shape, dataSource);
        context.getWriter().write("$L.Float($L)", dataDest, source);
        return null;
    }

    @Override
    public Void doubleShape(DoubleShape shape) {
        String source = conditionallyDereference(shape, dataSource);
        context.getWriter().write("$L.Double($L)", dataDest, source);
        return null;
    }

    @Override
    public Void timestampShape(TimestampShape shape) {
        String source = conditionallyDereference(shape, dataSource);
        GoWriter writer = context.getWriter();
        writer.addUseImports(SmithyGoDependency.SMITHY_TIME);

        switch (timestampFormat) {
            case DATE_TIME:
                writer.write("$L.String(smithytime.FormatDateTime($L))", dataDest, source);
                break;
            case HTTP_DATE:
                writer.write("$L.String(smithytime.FormatHTTPDate($L))", dataDest, source);
                break;
            case EPOCH_SECONDS:
                writer.write("$L.Double(smithytime.FormatEpochSeconds($L))", dataDest, source);
                break;
            case UNKNOWN:
                throw new CodegenException("Unknown timestamp format");
        }
        return null;
    }

    @Override
    public Void stringShape(StringShape shape) {
        String source = conditionallyDereference(shape, dataSource);
        if (shape.hasTrait(EnumTrait.class)) {
            source = String.format("string(%s)", source);
        }
        context.getWriter().write("$L.String($L)", dataDest, source);
        return null;
    }

    private String conditionallyDereference(Shape shape, String dataSource) {
        return CodegenUtils.isShapePassByReference(shape) ? "*" + dataSource : dataSource;
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
        throw new CodegenException(String.format("Cannot serialize shape type %s on protocol, shape: %s.",
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
        writeDelegateFunction(shape);
        return null;
    }

    @Override
    public Void setShape(SetShape shape) {
        writeDelegateFunction(shape);
        return null;
    }

    @Override
    public Void mapShape(MapShape shape) {
        writeDelegateFunction(shape);
        return null;
    }

    private void writeDelegateFunction(Shape shape) {
        String serFunctionName = ProtocolGenerator.getDocumentSerializerFunctionName(shape, context.getProtocolName());
        GoWriter writer = context.getWriter();
        writer.openBlock("if err := $L($L, $L); err != nil {", "}", serFunctionName, dataSource, dataDest, () -> {
            writer.write("return err");
        });
    }
}
