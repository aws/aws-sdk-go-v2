/*
 * Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

import java.util.Optional;
import java.util.Set;
import java.util.function.Function;
import java.util.function.Predicate;
import java.util.stream.Collectors;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.CodegenUtils;
import software.amazon.smithy.go.codegen.GoDependency;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.HttpBindingProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.HttpBinding;
import software.amazon.smithy.model.knowledge.HttpBindingIndex;
import software.amazon.smithy.model.shapes.CollectionShape;
import software.amazon.smithy.model.shapes.MapShape;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeType;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.model.traits.MediaTypeTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.utils.FunctionalUtils;

/**
 * Handles general components across the AWS JSON protocols that have HTTP bindings.
 * It handles reading and writing from document bodies, including generating any
 * functions needed for performing serde.
 *
 * @see <a href="https://awslabs.github.io/smithy/spec/http.html">Smithy HTTP protocol bindings.</a>
 */
abstract class RestJsonProtocolGenerator extends HttpBindingProtocolGenerator {
    /**
     * Creates a AWS JSON RPC protocol generator.
     */
    RestJsonProtocolGenerator() {
        super(true);
    }

    @Override
    protected TimestampFormatTrait.Format getDocumentTimestampFormat() {
        return TimestampFormatTrait.Format.EPOCH_SECONDS;
    }

    @Override
    protected void generateOperationDocumentSerializer(
            GenerationContext context,
            OperationShape operation
    ) {
        Model model = context.getModel();
        HttpBindingIndex bindingIndex = model.getKnowledge(HttpBindingIndex.class);
        Set<MemberShape> documentBindings = bindingIndex.getRequestBindings(operation, HttpBinding.Location.DOCUMENT)
                .stream()
                .map(HttpBinding::getMember)
                .collect(Collectors.toSet());

        if (documentBindings.size() == 0) {
            return;
        }

        Shape inputShape = model.expectShape(operation.getInput()
                .orElseThrow(() -> new CodegenException("Input shape missing for operation " + operation.getId())));
        GoWriter writer = context.getWriter();

        writeJsonShapeSerializerFunction(writer, model, context.getSymbolProvider(), inputShape,
                documentBindings::contains);
        writer.write("");
    }

    private void writeJsonShapeSerializerFunction(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            Shape inputShape,
            Predicate<MemberShape> filterMemberShapes
    ) {
        Symbol jsonEncoder = SymbolUtils.createPointableSymbolBuilder("Value", GoDependency.AWS_JSON_PROTOCOL)
                .build();
        Symbol inputSymbol = symbolProvider.toSymbol(inputShape);
        String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(inputShape, getProtocolName());

        writer.addUseImports(GoDependency.FMT);
        writer.openBlock("func $L(v $P, value $T) error {", "}", functionName, inputSymbol,
                jsonEncoder, () -> {
                    writer.openBlock("if v == nil {", "}", () -> {
                        writer.write("return fmt.Errorf(\"unsupported serialization of nil %T\", v)");
                    });
                    writer.write("");

                    switch (inputShape.getType()) {
                        case MAP:
                        case STRUCTURE:
                            writeShapeToJsonObject(model, symbolProvider, writer, inputShape, filterMemberShapes);
                            break;
                        case LIST:
                        case SET:
                            writeShapeToJsonArray(model, writer, (CollectionShape) inputShape);
                            break;
                        case UNION:
                        case DOCUMENT:
                            writer.write("// TODO: Support " + inputShape.getType().name() + " Serialization");
                            break;
                        default:
                            throw new CodegenException("Unexpected shape serialization to JSON");
                    }

                    writer.write("return nil");
                });
    }

    private void writeShapeToJsonArray(
            Model model,
            GoWriter writer,
            CollectionShape shape
    ) {
        MemberShape memberShape = shape.members().iterator().next();
        Shape targetShape = model.expectShape(memberShape.getTarget());

        writer.write("array := value.Array()");
        writer.write("defer array.Close()");
        writer.write("");

        writer.openBlock("for i := range v {", "}", () -> {
            writer.write("av := array.Value()");
            String operand = "v[i]";
            if (isShapeTypeDocumentSerializerRequired(targetShape.getType())) {
                String serFunctionName = ProtocolGenerator.getDocumentSerializerFunctionName(targetShape,
                        getProtocolName());
                operand = CodegenUtils.isShapePassByReference(targetShape) ? "&" + operand : operand;
                writer.openBlock("if err := $L($L, av); err != nil {", "}", serFunctionName, operand, () -> {
                    writer.write("return err");
                });
            } else {
                writer.write("av" + writeSimpleShapeToJsonValue(model, memberShape, operand));
            }
        });
    }

    private void writeShapeToJsonObject(
            Model model,
            SymbolProvider symbolProvider,
            GoWriter writer,
            Shape shape,
            Predicate<MemberShape> filterMemberShapes
    ) {
        writer.write("object := value.Object()");
        writer.write("defer object.Close()");
        writer.write("");

        switch (shape.getType()) {
            case MAP:
                writeMapShapeToJsonObject(model, writer, (MapShape) shape);
                break;
            case STRUCTURE:
                writeStructuredShapeToJsonObject(model, symbolProvider, writer, (StructureShape) shape,
                        filterMemberShapes);
                break;
            default:
                throw new CodegenException("Unexpected shape serialization to JSON Object");
        }
    }

    private void writeStructuredShapeToJsonObject(
            Model model,
            SymbolProvider symbolProvider,
            GoWriter writer,
            StructureShape shape,
            Predicate<MemberShape> filterMemberShapes
    ) {
        shape.members().forEach(memberShape -> {
            if (!filterMemberShapes.test(memberShape)) {
                return;
            }

            Shape targetShape = model.expectShape(memberShape.getTarget());

            writeSafeOperandAccessor(model, symbolProvider, memberShape, "v", writer,
                    (bodyWriter, operand) -> {
                        if (isShapeTypeDocumentSerializerRequired(targetShape.getType())) {
                            String serFunctionName = ProtocolGenerator.getDocumentSerializerFunctionName(targetShape,
                                    getProtocolName());
                            operand = CodegenUtils.isShapePassByReference(targetShape) ? "&" + operand : operand;
                            writer.openBlock("if err := $L($L, object.Key($S)); err != nil {", "}", serFunctionName,
                                    operand, memberShape.getMemberName(), () -> {
                                        writer.write("return err");
                                    });
                        } else {
                            writer.write("object.Key($S)" + writeSimpleShapeToJsonValue(model, memberShape, operand),
                                    memberShape.getMemberName());
                        }
                    });
            writer.write("");
        });
    }

    private void writeMapShapeToJsonObject(Model model, GoWriter writer, MapShape shape) {
        MemberShape memberShape = shape.getValue();
        Shape targetShape = model.expectShape(memberShape.getTarget());

        writer.openBlock("for key := range v {", "}", () -> {
            if (isShapeTypeDocumentSerializerRequired(targetShape.getType())) {
                String serFunctionName = ProtocolGenerator
                        .getDocumentSerializerFunctionName(targetShape,
                                getProtocolName());
                String operand = "v[key]";
                operand = CodegenUtils.isShapePassByReference(targetShape) ? "&" + operand : operand;
                writer.openBlock("if err := $L($L, object.Key(key)); err != nil {", "}", serFunctionName, operand,
                        () -> {
                            writer.write("return err");
                        });
            } else {
                writer.write("object.Key(key)" + writeSimpleShapeToJsonValue(model, memberShape, "v[key]"));
            }
        });
    }

    private String writeSimpleShapeToJsonValue(Model model, MemberShape memberShape, String operand) {
        Shape targetShape = model.expectShape(memberShape.getTarget());

        // JSON encoder helper methods take a value not a reference so we need to determine if we need to dereference.
        operand = CodegenUtils.isShapePassByReference(targetShape)
                ? "*" + operand : operand;

        switch (targetShape.getType()) {
            case BOOLEAN:
                return ".Boolean(" + operand + ")";
            case STRING:
                operand = targetShape.hasTrait(EnumTrait.class) ? "string(" + operand + ")" : operand;
                return ".String(" + operand + ")";
            case TIMESTAMP:
                // TODO: This needs to handle formats
                return ".UnixTime(" + operand + ")";
            case BYTE:
                return ".Byte(" + operand + ")";
            case SHORT:
                return ".Short(" + operand + ")";
            case INTEGER:
                return ".Integer(" + operand + ")";
            case LONG:
                return ".Long(" + operand + ")";
            case FLOAT:
                return ".Float(" + operand + ")";
            case DOUBLE:
                return ".Double(" + operand + ")";
            case BLOB:
                return ".Blob(" + operand + ")";
            default:
                throw new CodegenException("Unsupported shape type " + targetShape.getType());
        }
    }

    @Override
    protected void writeMiddlewareDocumentSerializerDelegator(
            Model model,
            SymbolProvider symbolProvider,
            OperationShape operation,
            GoStackStepMiddlewareGenerator generator,
            GoWriter writer
    ) {
        HttpBindingIndex httpBindingIndex = model.getKnowledge(HttpBindingIndex.class);
        boolean hasDocumentBindings = httpBindingIndex.getRequestBindings(operation, HttpBinding.Location.DOCUMENT)
                .size() > 0;
        Optional<HttpBinding> payloadBinding = httpBindingIndex.getRequestBindings(operation,
                HttpBinding.Location.PAYLOAD).stream().findFirst();

        if (!(payloadBinding.isPresent() || hasDocumentBindings)) {
            return;
        }

        writer.addUseImports(SymbolUtils.createValueSymbolBuilder(null, GoDependency.AWS_JSON_PROTOCOL).build());

        writer.write("var documentPayload []byte");
        writer.write("");

        if (payloadBinding.isPresent()) {
            MemberShape memberShape = payloadBinding.get().getMember();
            Shape payloadShape = model.expectShape(memberShape.getTarget());
            ShapeType shapeType = payloadShape.getType();
            String memberName = symbolProvider.toMemberName(memberShape);

            Optional<MediaTypeTrait> mediaTypeTrait = payloadShape.getTrait(MediaTypeTrait.class);
            mediaTypeTrait.ifPresent(typeTrait -> writer.write("restEncoder.SetHeader(\"Content-Type\").String($S)",
                    typeTrait.getValue()));

            if (shapeType == ShapeType.BLOB) {
                writer.write("documentPayload = input.$L", memberName);
            } else if (shapeType == ShapeType.STRING) {
                writer.write("documentPayload = []byte(input.$L)", memberName);
            } else {
                String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(payloadShape,
                        getProtocolName());
                writer.write("jsonEncoder := json.NewEncoder()");
                writer.openBlock("if err := $L(input.$L, jsonEncoder.Value); err != nil {", "}", functionName,
                        memberName, () -> {
                            writer.write("return err");
                        });
                writer.write("documentPayload = jsonEncoder.Bytes()");
            }
        } else {
            writer.write("restEncoder.SetHeader(\"Content-Type\").String($S)", getDocumentContentType());
            writer.write("");

            Shape inputShape = model.expectShape(operation.getInput()
                    .orElseThrow(() -> new CodegenException("Input shape is missing on " + operation.getId())));
            String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(inputShape, getProtocolName());
            writer.write("jsonEncoder := json.NewEncoder()");
            writer.openBlock("if err := $L(input, jsonEncoder.Value); err != nil {", "}", functionName, () -> {
                writer.write("return err");
            });
            writer.write("documentPayload = jsonEncoder.Bytes()");
        }
        writer.write("");

        writer.openBlock("if request, err = request.SetStream(bytes.NewReader(documentPayload)); err != nil {", "}",
                () -> {
                    writer.write("return out, metadata, &aws.SerializationError{Err: err}");
                });
    }

    @Override
    protected void generateDocumentBodyShapeSerializers(GenerationContext context, Set<Shape> shapes) {
        GoWriter writer = context.getWriter();
        Model model = context.getModel();
        SymbolProvider symbolProvider = context.getSymbolProvider();

        shapes.forEach(shape -> {
            writeJsonShapeSerializerFunction(writer, model, symbolProvider, shape, FunctionalUtils.alwaysTrue());
            writer.write("");
        });
    }

    @Override
    protected void generateDocumentBodyShapeDeserializers(GenerationContext context, Set<Shape> shapes) {
    }

    @Override
    public void generateSharedSerializerComponents(GenerationContext context) {
        super.generateSharedSerializerComponents(context);
        // pass
    }
}
