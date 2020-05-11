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
import java.util.stream.Collectors;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDependency;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.HttpBindingProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.HttpBinding;
import software.amazon.smithy.model.knowledge.HttpBindingIndex;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeType;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.model.traits.MediaTypeTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait;

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
            GenerationContext context, OperationShape operation
    ) {
        Model model = context.getModel();

        HttpBindingIndex bindingIndex = model.getKnowledge(HttpBindingIndex.class);

        Set<MemberShape> bindingMap = bindingIndex.getRequestBindings(operation).values().stream()
                .filter(binding -> binding.getLocation() == HttpBinding.Location.DOCUMENT)
                .map(HttpBinding::getMember)
                .collect(Collectors.toSet());

        if (bindingMap.size() == 0) {
            return;
        }

        Shape inputShape = model.expectShape(operation.getInput()
                .orElseThrow(() -> new CodegenException("input shapre missing for operation " + operation.getId())));

        GoWriter writer = context.getWriter();

        writeJsonShapeSerializerFunction(writer, model, context.getSymbolProvider(), inputShape,
                bindingMap::contains);
        writer.write("");
    }

    private void writeJsonShapeSerializerFunction(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            Shape inputShape,
            Function<MemberShape, Boolean> filterMemberShapes
    ) {
        Symbol jsonEncoder = SymbolUtils.createValueSymbolBuilder("Value", GoDependency.AWS_JSON_PROTOCOL)
                .build();

        Symbol inputSymbol = symbolProvider.toSymbol(inputShape);

        String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(inputShape, getProtocolName());

        writer.addUseImports(SymbolUtils.createValueSymbolBuilder(null, GoDependency.FMT).build());
        writer.addUseImports(inputSymbol);
        writer.addUseImports(jsonEncoder);

        writer.openBlock("func $L(v $P, value $P) error {", "}", functionName, inputSymbol,
                jsonEncoder, () -> {
                    writer.openBlock("if v == nil {", "}", () -> {
                        writer.write("return fmt.Errorf(\"unsupported serialization of nil %T\", v)");
                    });
                    writer.write("");

                    switch (inputShape.getType()) {
                        case MAP:
                        case STRUCTURE:
                            writeShapeToJsonObject(model, symbolProvider, writer, inputShape,
                                    filterMemberShapes);
                            break;
                        case LIST:
                        case SET:
                            writeShapeToJsonArray(model, symbolProvider, writer, inputShape);
                            break;
                        case UNION:
                            writer.write("// TODO: Support " + inputShape.getType().name() + " Serialization");
                        case DOCUMENT:
                        default:
                            throw new CodegenException("unexpected type");
                    }

                    writer.write("return nil");
                });
    }

    private void writeShapeToJsonArray(
            Model model,
            SymbolProvider symbolProvider,
            GoWriter writer,
            Shape shape
    ) {
        if (shape.members().size() > 1) {
            throw new CodegenException("not possible to serialize shape with only multiple member shapes"
                    + " to an array");
        }

        MemberShape memberShape = shape.members().iterator().next();
        Shape targetShape = model.expectShape(memberShape.getTarget());

        writer.write("array := value.Array()");
        writer.write("defer array.Close()");
        writer.write("");

        writer.openBlock("for i := range v {", "}", () -> {
            writer.write("av := array.Value()");
            if (isShapeTypeDocumentSerializerRequired(targetShape.getType())) {
                String serFunctionName = ProtocolGenerator.getDocumentSerializerFunctionName(targetShape,
                        getProtocolName());
                writer.openBlock("if err := $L(v[i], av); err != nil {", "}", serFunctionName, () -> {
                    writer.write("return err");
                });
            } else {
                writer.write("av" + writeSimpleShapeToJsonValue(targetShape,
                        symbolProvider.toSymbol(targetShape), "v[i]"));
            }
        });
    }

    private void writeShapeToJsonObject(
            Model model,
            SymbolProvider symbolProvider,
            GoWriter writer,
            Shape shape,
            Function<MemberShape, Boolean> filterMemberShapes
    ) {
        writer.write("object := value.Object()");
        writer.write("defer object.Close()");
        writer.write("");

        shape.members().forEach(memberShape -> {
            if (!filterMemberShapes.apply(memberShape)) {
                return;
            }

            Shape targetShape = model.expectShape(memberShape.getTarget());
            String fieldName = symbolProvider.toMemberName(memberShape);

            writeSafeOperandAccessor(model, symbolProvider, memberShape, "v", writer,
                    (bodyWriter, operand) -> {
                        if (isShapeTypeDocumentSerializerRequired(targetShape.getType())) {
                            String serFunctionName = ProtocolGenerator
                                    .getDocumentSerializerFunctionName(targetShape,
                                            getProtocolName());
                            writer.openBlock("if err := $L(v.$L, object.Key($S)); err != nil {", "}", serFunctionName,
                                    fieldName, memberShape.getMemberName(), () -> {
                                        writer.write("return err");
                                    });
                        } else {
                            writer.write("object.Key($S)" + writeSimpleShapeToJsonValue(targetShape,
                                    symbolProvider.toSymbol(targetShape), operand), memberShape.getMemberName());
                        }
                    });
            writer.write("");
        });
    }

    private String writeSimpleShapeToJsonValue(Shape targetShape, Symbol targetSymbol, String operand) {
        operand = isDereferenceRequired(targetShape, targetSymbol) ? "*" + operand : operand;

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
            default:
                throw new CodegenException("unsupported shape type " + targetShape.getType());
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
        Boolean hasDocumentBindings = model.getKnowledge(HttpBindingIndex.class).getRequestBindings(operation)
                .values().stream()
                .map(binding -> binding.getLocation() == HttpBinding.Location.DOCUMENT)
                .filter(aBoolean -> aBoolean)
                .findFirst().orElse(false);

        Optional<HttpBinding> payloadBinding = model.getKnowledge(HttpBindingIndex.class).getRequestBindings(operation)
                .values().stream()
                .filter(binding -> binding.getLocation() == HttpBinding.Location.PAYLOAD)
                .findFirst();

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
                    .orElseThrow(() -> new CodegenException("input shape is missing on " + operation.getId())));
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
            writeJsonShapeSerializerFunction(writer, model, symbolProvider, shape, memberShape -> true);
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
