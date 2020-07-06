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

import java.io.Writer;
import java.util.Collection;
import java.util.Optional;
import java.util.Set;
import java.util.TreeSet;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Predicate;
import java.util.stream.Collectors;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.CodegenUtils;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.HttpBindingProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUnitTestGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUnitTestResponseGenerator;
import software.amazon.smithy.go.codegen.integration.IdempotencyTokenMiddlewareGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolTestGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUnitTestRequestGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUnitTestResponseErrorGenerator;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.HttpBinding;
import software.amazon.smithy.model.knowledge.HttpBindingIndex;
import software.amazon.smithy.model.shapes.CollectionShape;
import software.amazon.smithy.model.shapes.MapShape;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.ShapeType;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.shapes.UnionShape;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.model.traits.HttpErrorTrait;
import software.amazon.smithy.model.traits.JsonNameTrait;
import software.amazon.smithy.model.traits.MediaTypeTrait;
import software.amazon.smithy.model.traits.StreamingTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.utils.FunctionalUtils;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.SetUtils;

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
        String functionName = ProtocolGenerator.getOperationDocumentSerFunctionName(inputShape, getProtocolName());

        writeJsonShapeSerializerFunction(writer, model, context.getSymbolProvider(), functionName, inputShape,
                documentBindings::contains);
        writer.write("");
    }

    private void writeJsonShapeSerializerFunction(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            String functionName,
            Shape inputShape,
            Predicate<MemberShape> filterMemberShapes
    ) {
        Symbol jsonEncoder = SymbolUtils.createPointableSymbolBuilder("Value", SmithyGoDependency.SMITHY_JSON).build();
        Symbol inputSymbol = symbolProvider.toSymbol(inputShape);

        writer.addUseImports(SmithyGoDependency.FMT);
        writer.openBlock("func $L(v $P, value $T) error {", "}", functionName, inputSymbol,
                jsonEncoder, () -> {
                    writer.openBlock("if v == nil {", "}", () -> {
                        writer.write("return fmt.Errorf(\"unsupported serialization of nil %T\", v)");
                    });
                    writer.write("");

                    switch (inputShape.getType()) {
                        case UNION:
                        case MAP:
                        case STRUCTURE:
                            writeShapeToJsonObject(model, symbolProvider, writer, inputShape, filterMemberShapes);
                            break;
                        case LIST:
                        case SET:
                            writeShapeToJsonArray(model, writer, (CollectionShape) inputShape);
                            break;
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
            if (!targetShape.hasTrait(EnumTrait.class)) {
                writer.openBlock("if vv := v[i]; vv == nil {", "}", () -> {
                    writer.write("av.Null()");
                    writer.write("continue");
                });
            }
            String operand = "v[i]";
            if (isShapeTypeDocumentSerializerRequired(targetShape.getType())) {
                String serFunctionName = ProtocolGenerator.getDocumentSerializerFunctionName(targetShape,
                        getProtocolName());

                writer.openBlock("if err := $L($L, av); err != nil {", "}", serFunctionName, operand, () -> {
                    writer.write("return err");
                });
            } else {
                generateSimpleShapeToJsonValue(model, writer, memberShape, operand, (w, s) -> w.write("av.$L", s));
            }
        });
    }

    @Override
    public void generateProtocolTests(GenerationContext context) {
        Set<HttpProtocolUnitTestGenerator.ConfigValue> configValues = new TreeSet<>();
        configValues.addAll(SetUtils.of(
                HttpProtocolUnitTestGenerator.ConfigValue.builder()
                        .name(AddAwsConfigFields.REGION_CONFIG_NAME)
                        .value(writer -> {
                            writer.write("$S,", "us-west-2");
                        })
                        .build(),
                HttpProtocolUnitTestGenerator.ConfigValue.builder()
                        .name(AddAwsConfigFields.HTTP_CLIENT_CONFIG_NAME)
                        .value(writer -> {
                            writer.addUseImports(AwsGoDependency.AWS_CORE);
                            writer.write("aws.NewBuildableHTTPClient(),");
                        })
                        .build(),
                HttpProtocolUnitTestGenerator.ConfigValue.builder()
                        .name(AddAwsConfigFields.ENDPOINT_RESOLVER_CONFIG_NAME)
                        .value(writer -> {
                            writer.addUseImports(AwsGoDependency.AWS_CORE);
                            writer.openBlock("aws.EndpointResolverFunc("
                                            + "func(service, region string) (e aws.Endpoint, err error) {",
                                    "}),", () -> {
                                        writer.write("e.URL = server.URL");
                                        writer.write("e.SigningRegion = \"us-west-2\"");
                                        writer.write("return e, err");
                                    });
                        })
                        .build(),
                HttpProtocolUnitTestGenerator.ConfigValue.builder()
                        .name("APIOptions")
                        .value(writer -> {
                            writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);
                            writer.openBlock("[]APIOptionFunc{", "},", () -> {
                                writer.openBlock("func(s *middleware.Stack) error {", "},", () -> {
                                    writer.write("s.Build.Clear()");
                                    writer.write("s.Finalize.Clear()");
                                    writer.write("return nil");
                                });
                            });
                        })
                        .build()
        ));

        // TODO can this check be replaced with a lookup into the runtime plugins?
        if (IdempotencyTokenMiddlewareGenerator.hasOperationsWithIdempotencyToken(context.getModel(), context.getService())) {
            configValues.add(
                    HttpProtocolUnitTestGenerator.ConfigValue.builder()
                            .name(IdempotencyTokenMiddlewareGenerator.IDEMPOTENCY_CONFIG_NAME)
                            .value(writer -> {
                                writer.addUseImports(SmithyGoDependency.SMITHY_RAND);
                                writer.addUseImports(SmithyGoDependency.SMITHY_TESTING);
                                writer.write("smithyrand.NewUUIDIdempotencyToken(&smithytesting.ByteLoop{}),");
                            })
                    .build()
            );
        }

        new HttpProtocolTestGenerator(context,
                (HttpProtocolUnitTestRequestGenerator.Builder) new HttpProtocolUnitTestRequestGenerator
                        .Builder()
                        .addClientConfigValues(configValues),
                (HttpProtocolUnitTestResponseGenerator.Builder) new HttpProtocolUnitTestResponseGenerator
                        .Builder()
                        .addClientConfigValues(configValues),
                (HttpProtocolUnitTestResponseErrorGenerator.Builder) new HttpProtocolUnitTestResponseErrorGenerator
                        .Builder()
                        .addClientConfigValues(configValues)
        ).generateProtocolTests();
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
            case UNION:
                writeUnionShapeToJsonObject(model, symbolProvider, writer, (UnionShape) shape);
                break;
            default:
                throw new CodegenException("Unexpected shape serialization to JSON Object");
        }
    }

    private void writeUnionShapeToJsonObject(
            Model model,
            SymbolProvider symbolProvider,
            GoWriter writer,
            UnionShape shape
    ) {
        Symbol symbol = symbolProvider.toSymbol(shape);

        writer.addUseImports(SmithyGoDependency.FMT);

        writer.openBlock("switch uv := v.(type) {", "}", () -> {
            for (MemberShape memberShape : shape.getAllMembers().values()) {
                Shape targetShape = model.expectShape(memberShape.getTarget());
                String exportedMemberName = symbol.getName() + symbolProvider.toMemberName(memberShape);

                writer.openBlock("case *$L:", "", exportedMemberName, () -> {
                    if (isShapeTypeDocumentSerializerRequired(targetShape.getType())) {
                        String serFunctionName = ProtocolGenerator.getDocumentSerializerFunctionName(targetShape,
                                getProtocolName());
                        writer.write("av := object.key($S)", getSerializedMemberName(memberShape));
                        writer.openBlock("if err := $L(uv.Value(), av); err != nil {", "}", serFunctionName, () -> {
                            writer.write("return err");
                        });
                    } else {
                        generateSimpleShapeToJsonValue(model, writer, memberShape, "uv.Value()", (w, s) -> {
                            writer.write("object.Key($S).$L", getSerializedMemberName(memberShape), s);
                        });
                    }
                });
            }
            writer.openBlock("case *$LUnknown:", "", symbol.getName(), () -> writer.write("fallthrough"));
            writer.openBlock("default:", "", () -> {
                writer.write("return fmt.Errorf(\"attempted to serialize unknown member type %T"
                        + " for union %T\", uv, v)");
            });
        });
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
                            writer.openBlock("if err := $L($L, object.Key($S)); err != nil {", "}", serFunctionName,
                                    operand, getSerializedMemberName(memberShape), () -> {
                                        writer.write("return err");
                                    });
                        } else {
                            generateSimpleShapeToJsonValue(model, writer, memberShape, operand, (w, s) -> {
                                writer.write("object.Key($S).$L", getSerializedMemberName(memberShape), s);
                            });
                        }
                    });
            writer.write("");
        });
    }

    private void writeMapShapeToJsonObject(Model model, GoWriter writer, MapShape shape) {
        MemberShape memberShape = shape.getValue();
        Shape targetShape = model.expectShape(memberShape.getTarget());

        writer.openBlock("for key := range v {", "}", () -> {
            writer.write("om := object.Key(key)");
            if (!targetShape.hasTrait(EnumTrait.class)) {
                writer.openBlock("if vv := v[key]; vv == nil {", "}", () -> {
                    writer.write("om.Null()");
                    writer.write("continue");
                });
            }
            if (isShapeTypeDocumentSerializerRequired(targetShape.getType())) {
                String serFunctionName = ProtocolGenerator
                        .getDocumentSerializerFunctionName(targetShape,
                                getProtocolName());
                writer.openBlock("if err := $L($L, om); err != nil {", "}", serFunctionName, "v[key]",
                        () -> {
                            writer.write("return err");
                        });
            } else {
                generateSimpleShapeToJsonValue(model, writer, memberShape, "v[key]", (w, s) -> {
                    writer.write("om.$L", s);
                });
            }
        });
    }

    private void generateSimpleShapeToJsonValue(
            Model model,
            GoWriter writer,
            MemberShape memberShape,
            String operand,
            BiConsumer<GoWriter, String> locationEncoder
    ) {
        Shape targetShape = model.expectShape(memberShape.getTarget());

        // JSON encoder helper methods take a value not a reference so we need to determine if we need to dereference.
        operand = CodegenUtils.isShapePassByReference(targetShape)
                && targetShape.getType() != ShapeType.BIG_INTEGER
                && targetShape.getType() != ShapeType.BIG_DECIMAL
                ? "*" + operand : operand;

        switch (targetShape.getType()) {
            case BOOLEAN:
                locationEncoder.accept(writer, "Boolean(" + operand + ")");
                break;
            case STRING:
                operand = targetShape.hasTrait(EnumTrait.class) ? "string(" + operand + ")" : operand;
                locationEncoder.accept(writer, "String(" + operand + ")");
                break;
            case TIMESTAMP:
                generateDocumentTimestampSerializer(model, writer, memberShape, operand, locationEncoder);
                break;
            case BYTE:
                locationEncoder.accept(writer, "Byte(" + operand + ")");
                break;
            case SHORT:
                locationEncoder.accept(writer, "Short(" + operand + ")");
                break;
            case INTEGER:
                locationEncoder.accept(writer, "Integer(" + operand + ")");
                break;
            case LONG:
                locationEncoder.accept(writer, "Long(" + operand + ")");
                break;
            case FLOAT:
                locationEncoder.accept(writer, "Float(" + operand + ")");
                break;
            case DOUBLE:
                locationEncoder.accept(writer, "Double(" + operand + ")");
                break;
            case BLOB:
                locationEncoder.accept(writer, "Base64EncodeBytes(" + operand + ")");
                break;
            default:
                throw new CodegenException("Unsupported shape type " + targetShape.getType());
        }
    }

    @Override
    protected void writeMiddlewarePayloadSerializerDelegator(
            Model model,
            SymbolProvider symbolProvider,
            OperationShape operation,
            MemberShape memberShape,
            GoStackStepMiddlewareGenerator generator,
            GoWriter writer
    ) {
        Shape payloadShape = model.expectShape(memberShape.getTarget());

        writeSafeOperandAccessor(model, symbolProvider, memberShape, "input", writer, (w, s) -> {
            writer.openBlock("if !restEncoder.HasHeader(\"Content-Type\") {", "}", () -> {
                writer.write("restEncoder.SetHeader(\"Content-Type\").String($S)", getPayloadShapeMediaType(payloadShape));
            });
            writer.write("");

            if (payloadShape.hasTrait(StreamingTrait.class)) {
                w.write("payload := $L", s);

            } else if (payloadShape.isBlobShape()) {
                w.addUseImports(SmithyGoDependency.BYTES);
                w.write("payload := bytes.NewReader($L)", s);

            } else if (payloadShape.isStringShape()) {
                w.addUseImports(SmithyGoDependency.STRINGS);
                w.write("payload := strings.NewReader(*$L)", s);

            } else {
                String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(payloadShape,
                        getProtocolName());
                w.addUseImports(SmithyGoDependency.SMITHY_JSON);
                w.write("jsonEncoder := smithyjson.NewEncoder()");
                w.openBlock("if err := $L($L, jsonEncoder.Value); err != nil {", "}", functionName,
                        s, () -> {
                            w.write("return out, metadata, &smithy.SerializationError{Err: err}");
                        });
                w.write("payload := bytes.NewReader(jsonEncoder.Bytes())");
            }

            w.openBlock("if request, err = request.SetStream(payload); err != nil {", "}",
                    () -> {
                        w.write("return out, metadata, &smithy.SerializationError{Err: err}");
                    });
        });
    }

    /**
     * Retruns the MediaType for the payload shape derived from the MediaTypeTrait, shape type, or document content type.
     *
     * @param payloadShape shape bound to the payload.
     * @return string for media type.
     */
    private String getPayloadShapeMediaType(Shape payloadShape) {
        Optional<MediaTypeTrait> mediaTypeTrait = payloadShape.getTrait(MediaTypeTrait.class);

        if (mediaTypeTrait.isPresent()) {
            return mediaTypeTrait.get().getValue();
        }

        if (payloadShape.isBlobShape()) {
            return "application/octet-stream";
        }

        if (payloadShape.isStringShape()) {
            return "test/plain";
        }

        return getDocumentContentType();
    }


    @Override
    protected void writeMiddlewareDocumentSerializerDelegator(
            Model model,
            SymbolProvider symbolProvider,
            OperationShape operation,
            GoStackStepMiddlewareGenerator generator,
            GoWriter writer
    ) {

        writer.addUseImports(SmithyGoDependency.SMITHY);
        writer.addUseImports(SmithyGoDependency.SMITHY_JSON);

        writer.write("restEncoder.SetHeader(\"Content-Type\").String($S)", getDocumentContentType());
        writer.write("");

        Shape inputShape = model.expectShape(operation.getInput()
                .orElseThrow(() -> new CodegenException("Input shape is missing on " + operation.getId())));

        String functionName = ProtocolGenerator.getOperationDocumentSerFunctionName(inputShape, getProtocolName());
        writer.addUseImports(SmithyGoDependency.SMITHY_JSON);
        writer.write("jsonEncoder := smithyjson.NewEncoder()");
        writer.openBlock("if err := $L(input, jsonEncoder.Value); err != nil {", "}", functionName, () -> {
            writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
        });
        writer.write("");

        writer.addUseImports(SmithyGoDependency.BYTES);
        writer.openBlock("if request, err = request.SetStream(bytes.NewReader(jsonEncoder.Bytes())); err != nil {", "}",
                () -> {
                    writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
                });
    }

    @Override
    protected void generateDocumentBodyShapeSerializers(GenerationContext context, Set<Shape> shapes) {
        GoWriter writer = context.getWriter();
        Model model = context.getModel();
        SymbolProvider symbolProvider = context.getSymbolProvider();

        shapes.forEach(shape -> {
            String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(shape, getProtocolName());
            writeJsonShapeSerializerFunction(writer, model, symbolProvider, functionName, shape,
                    FunctionalUtils.alwaysTrue());
            writer.write("");
        });
    }

    @Override
    public void generateSharedSerializerComponents(GenerationContext context) {
        super.generateSharedSerializerComponents(context);
        // pass
    }

    /**
     * Get the serialized name to be used for the member shape.
     *
     * @param memberShape the member shape
     * @return the serialized member name
     */
    private String getSerializedMemberName(MemberShape memberShape) {
        Optional<JsonNameTrait> jsonNameTrait = memberShape.getTrait(JsonNameTrait.class);
        return jsonNameTrait.isPresent() ? jsonNameTrait.get().getValue() : memberShape.getMemberName();
    }

    /**
     * Generate the serializer statement for the document timestamp
     *
     * @param model       the model
     * @param writer      the writer
     * @param memberShape the timestamp member shape to serialize
     * @param operand     the go operand
     */
    private void generateDocumentTimestampSerializer(
            Model model,
            GoWriter writer,
            MemberShape memberShape,
            String operand,
            BiConsumer<GoWriter, String> locationEncoder
    ) {
        writer.addUseImports(SmithyGoDependency.SMITHY_TIME);

        TimestampFormatTrait.Format format = memberShape.getMemberTrait(model, TimestampFormatTrait.class)
                .map(TimestampFormatTrait::getFormat)
                .orElse(TimestampFormatTrait.Format.EPOCH_SECONDS);

        switch (format) {
            case DATE_TIME:
                locationEncoder.accept(writer, "String(smithytime.FormatDateTime(" + operand + "))");
                break;
            case HTTP_DATE:
                locationEncoder.accept(writer, "String(smithytime.FormatHTTPDate(" + operand + "))");
                break;
            case EPOCH_SECONDS:
                locationEncoder.accept(writer, "Double(smithytime.FormatEpochSeconds(" + operand + "))");
                break;
            case UNKNOWN:
                throw new CodegenException("Unknown timestamp format");
        }
    }

    @Override
    protected void writeMiddlewareDocumentDeserializerDelegator(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            OperationShape operation,
            GoStackStepMiddlewareGenerator generator
    ) {
        Shape targetShape = model.expectShape(operation.getOutput().get());
        String operand = "output";

        boolean isShapeWithPayloadBinding = isShapeWithResponseBindings(model, operation, HttpBinding.Location.PAYLOAD);
        if (isShapeWithPayloadBinding) {
            // since payload trait can only be applied to a single member in a output shape
            MemberShape memberShape = model.getKnowledge(HttpBindingIndex.class)
                    .getResponseBindings(operation, HttpBinding.Location.PAYLOAD).stream()
                    .findFirst()
                    .orElseThrow(() -> new CodegenException("Expected payload binding member"))
                    .getMember();

            Shape payloadShape = model.expectShape(memberShape.getTarget());

            // if target shape is of type String or type Blob, then delegate deserializers for explicit payload shapes
            if (payloadShape.isStringShape() || payloadShape.isBlobShape()) {
                writeMiddlewarePayloadBindingDeserializerDelegator(writer, targetShape, true);
                return;
            }
            // for other payload target types we should deserialize using the appropriate document deserializer
            targetShape = payloadShape;
            operand += "." + symbolProvider.toMemberName(memberShape);
        }

        writer.addUseImports(SmithyGoDependency.SMITHY_IO);
        writer.write("buff := make([]byte, 1024)");
        writer.write("ringBuffer := smithyio.NewRingBuffer(buff)");
        writer.write("");

        writer.addUseImports(SmithyGoDependency.IO);
        writer.write("body := io.TeeReader(response.Body, ringBuffer)");
        writer.write("defer response.Body.Close()");
        writer.write("");

        writer.addUseImports(SmithyGoDependency.JSON);
        writer.write("decoder := json.NewDecoder(body)");
        writer.write("decoder.UseNumber()");
        writer.write("");

        writeMiddlewareDocumentBindingDeserializerDelegator(writer, targetShape, operand, !isShapeWithPayloadBinding);
    }

    // Writes middleware that delegates to deserializers for shapes that have explicit payload.
    private void writeMiddlewarePayloadBindingDeserializerDelegator(
            GoWriter writer,
            Shape shape,
            boolean isOperationShape
    ) {
        String deserFuncName = isOperationShape ?
                ProtocolGenerator.getDocumentOutputDeserializerFunctionName(shape, getProtocolName()) :
                ProtocolGenerator.getDocumentDeserializerFunctionName(shape, getProtocolName());
        writer.write("err = $L(output, response.Body)", deserFuncName);
        writer.openBlock("if err != nil {", "}", () -> {
            writer.addUseImports(SmithyGoDependency.SMITHY);
            writer.write(String.format("return out, metadata, &smithy.DeserializationError{Err:%s}",
                    "fmt.Errorf(\"failed to deserialize response payload, %w\", err)"));
        });
    }


    // Write middleware that delegates to deserializers for shapes that have implicit payload and deserializer
    private void writeMiddlewareDocumentBindingDeserializerDelegator(
            GoWriter writer,
            Shape shape,
            String operand,
            boolean isOperationShape
    ) {
        String deserFuncName = isOperationShape ?
                ProtocolGenerator.getDocumentOutputDeserializerFunctionName(shape, getProtocolName()) :
                ProtocolGenerator.getDocumentDeserializerFunctionName(shape, getProtocolName());
        writer.write("err = $L(&$L, decoder)", deserFuncName, operand);
        writer.openBlock("if err != nil {", "}", () -> {
            writer.addUseImports(SmithyGoDependency.BYTES);
            writer.addUseImports(SmithyGoDependency.SMITHY);
            writer.write("var snapshot bytes.Buffer");
            writer.write("io.Copy(&snapshot, ringBuffer)");
            writer.openBlock("return out, metadata, &smithy.DeserializationError {", "}", () -> {
                writer.write("Err: fmt.Errorf(\"failed to decode response body with invalid JSON, %w\", err),");
                writer.write("Snapshot: snapshot.Bytes(),");
            });
        });
    }

    @Override
    protected void writeMiddlewareErrorDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            OperationShape operationShape,
            GoStackStepMiddlewareGenerator generator
    ) {
        Collection<ShapeId> ErrorShapeIds = operationShape.getErrors();

        // checks if response has an error and retrieve the error code from the response
        writer.openBlock("if response.StatusCode < 200 || response.StatusCode >= 300 {", "}", () -> {
            // Retrieve error shape name from response. For REST JSON protocol, the error shape name can be found either
            // at Header `X-Amzn-Errortype` or a body field with the name `code`, or a body field named `__type`.
            writer.write("errorType := response.Header.Get($S)", "X-Amzn-Errortype");

            writer.addUseImports(AwsGoDependency.AWS_REST_JSON_PROTOCOL);
            writer.write("errorType = restjson.SanitizeErrorCode(errorType)");
            writer.write("");

            // if no modeled exceptions for the operation shape, return the response body as is
            if (ErrorShapeIds.size() == 0) {
                writer.addUseImports(SmithyGoDependency.JSON);
                writer.write("decoder := json.NewDecoder(response.Body)");
                writer.write("decoder.UseNumber()");
                writer.write("defer response.Body.Close()");
                writer.write("");

                writer.addUseImports(AwsGoDependency.AWS_REST_JSON_PROTOCOL);
                writer.addUseImports(SmithyGoDependency.SMITHY);
                writer.write("genericError, err := restjson.GetSmithyGenericAPIError(decoder, errorType)");
                writer.write("if err != nil { return out, metadata, &smithy.DeserializationError{ Err: err}}");
                writer.write("return out, metadata, genericError");
                return;
            }

            writer.addUseImports(SmithyGoDependency.SMITHY_IO);
            writer.write("buff := make([]byte, 1024)");
            writer.write("ringBuffer := smithyio.NewRingBuffer(buff)");
            writer.write("");

            writer.addUseImports(SmithyGoDependency.BYTES);
            writer.write("var errorBuffer bytes.Buffer");

            writer.addUseImports(SmithyGoDependency.SMITHY);
            writer.addUseImports(SmithyGoDependency.IO);
            writer.write("defer response.Body.Close()");
            writer.write("_, err := io.Copy(&errorBuffer, response.Body)");
            writer.openBlock("if err != nil {", "}", () -> {
                writer.write(String.format("return out, metadata, &smithy.DeserializationError{Err: %s}",
                        "fmt.Errorf(\"failed to copy error response body, %w\", err)"));
            });
            writer.write("");

            writer.write("errorBody := bytes.NewReader(errorBuffer.Bytes())");

            writer.addUseImports(SmithyGoDependency.JSON);
            writer.write("decoder := json.NewDecoder(io.TeeReader(errorBody, ringBuffer))");
            writer.write("decoder.UseNumber()");
            writer.write("");

            writer.write("var errorMessage string");
            writer.addUseImports(AwsGoDependency.AWS_REST_JSON_PROTOCOL);

            // If errorType is empty, look for error type in a body field with the name `code`,
            // or a body field named `__type`.
            writer.openBlock("if len(errorType) == 0 {", "}", () -> {
                writer.write("errorType, errorMessage, err = restjson.GetErrorInfo(decoder)");
                writer.openBlock("if err != nil {", "}", () -> {
                    writer.addUseImports(SmithyGoDependency.SMITHY);
                    writer.write("var snapshot bytes.Buffer");
                    writer.write("io.Copy(&snapshot, ringBuffer)");
                    writer.openBlock("return out, metadata, &smithy.DeserializationError {", "}", () -> {
                        writer.write(
                                "Err: fmt.Errorf(\"failed to decode response error with invalid JSON, %w\", err),");
                        writer.write("Snapshot: snapshot.Bytes(),");
                    });
                });
            });

            writer.write("");

            writer.openBlock("if len(errorType) == 0 {", "}", () -> {
                writer.openBlock("switch response.StatusCode {", "}", () -> {
                    for (ShapeId errorShapeId : ErrorShapeIds) {
                        Shape errorShape = model.expectShape(errorShapeId);
                        if (errorShape.hasTrait(HttpErrorTrait.class)) {
                            int statusCode = errorShape.getTrait(HttpErrorTrait.class).get().getCode();
                            writer.write("case $L: errorType = $S", statusCode, errorShapeId.getName());
                        }
                    }
                });
            });

            writer.write("");
            writer.write("// reset the ring buffer");
            writer.write("ringBuffer.Reset()");

            writer.addUseImports(SmithyGoDependency.IO);
            writer.write("// seek start of error body");
            writer.write("errorBody.Seek(0, io.SeekStart)");
            writer.write("");

            // generate middleware for modeled error shapes
            writeErrorShapeDeserializerDelegator(writer, model, symbolProvider, ErrorShapeIds);
            writer.write("");

            writer.openBlock("if len(errorMessage) != 0 {", "}", () -> {
                writer.addUseImports(SmithyGoDependency.SMITHY);
                writer.openBlock("genericError := &smithy.GenericAPIError{", "}", () -> {
                    writer.write("Code : errorType,");
                    writer.write("Message : errorMessage,");
                });
                writer.write("");
                writer.write("return out, metadata, genericError");
            });

            writer.write("");
            writer.addUseImports(AwsGoDependency.AWS_REST_JSON_PROTOCOL);
            writer.addUseImports(SmithyGoDependency.SMITHY);
            writer.write("genericError, err := restjson.GetSmithyGenericAPIError(decoder, errorType)");
            writer.write("if err != nil { return out, metadata, &smithy.DeserializationError{ Err: err }}");
            writer.write("");
            writer.write("return out, metadata, genericError");
        });
        writer.write("");
    }

    // writeErrorShapeMiddlewareDelegator takes in the list of error shapes, and generates
    // middleware error shape delegators.  It delegates based on whether the error shape has
    // rest bindings, payload bindings, document bindings.
    private void writeErrorShapeDeserializerDelegator(
            GoWriter writer, Model model, SymbolProvider symbolProvider,
            Collection<ShapeId> ErrorShapeIds
    ) {
        writer.write("decoder = json.NewDecoder(io.TeeReader(errorBody, ringBuffer))");
        writer.write("decoder.UseNumber()");
        writer.write("");

        for (ShapeId errorShapeId : ErrorShapeIds) {
            Shape errorShape = model.expectShape(errorShapeId);
            Symbol errorSymbol = symbolProvider.toSymbol(errorShape);

            writer.openBlock("if errorType == $S {", "}", errorShapeId.getName(), () -> {
                writer.write("errResult := &$T{}", errorSymbol);
                writer.write("output := errResult");
                writer.write("_ = output");
                writer.write("");


                if (isShapeWithRestResponseBindings(model, errorShape)) {
                    String deserFuncName = ProtocolGenerator.getOperationHttpBindingsDeserFunctionName(
                            errorShape, getProtocolName());

                    writer.write("err = $L(output, response)", deserFuncName);
                    writer.openBlock("if err != nil {", "}", () -> {
                        writer.addUseImports(SmithyGoDependency.SMITHY);
                        writer.write(String.format("return out, metadata, &smithy.DeserializationError{Err: %s}",
                                "fmt.Errorf(\"failed to decode response error with invalid Http bindings, %w\", err)"));
                    });
                    writer.write("");
                }

                if (isShapeWithResponseBindings(model, errorShape, HttpBinding.Location.DOCUMENT)
                        || isShapeWithResponseBindings(model, errorShape, HttpBinding.Location.PAYLOAD)) {
                    writeMiddlewareDocumentBindingDeserializerDelegator(writer, errorShape, "output", false);
                }

                // TODO: fix variable scoping and shadowing
                writer.write("return out, metadata, errResult");
            });
            writer.write("");
        }
    }

    @Override
    protected void generateOperationDocumentDeserializer(
            GenerationContext context,
            OperationShape operation
    ) {
        Model model = context.getModel();
        HttpBindingIndex bindingIndex = model.getKnowledge(HttpBindingIndex.class);
        Set<MemberShape> documentBindings = bindingIndex.getResponseBindings(operation, HttpBinding.Location.DOCUMENT)
                .stream()
                .map(HttpBinding::getMember)
                .collect(Collectors.toSet());

        Shape outputShape = model.expectShape(operation.getOutput()
                .orElseThrow(() -> new CodegenException("Output shape missing for operation " + operation.getId())));
        GoWriter writer = context.getWriter();

        if (documentBindings.size() != 0) {
            writeDocumentBindingDeserializer(writer, model, context.getSymbolProvider(), outputShape,
                    documentBindings::contains, true);
            writer.write("");
        }

        Set<MemberShape> payloadBindings = bindingIndex.getResponseBindings(operation, HttpBinding.Location.PAYLOAD)
                .stream()
                .map(HttpBinding::getMember)
                .collect(Collectors.toSet());

        if (payloadBindings.size() == 0) {
            return;
        }

        writePayloadBindingDeserializer(writer, model, context.getSymbolProvider(), outputShape,
                payloadBindings::contains, true);
        writer.write("");
    }

    @Override
    protected void generateErrorDocumentBindingDeserializer(GenerationContext context, ShapeId shapeId) {
        Model model = context.getModel();
        Shape shape = model.expectShape(shapeId);
        GoWriter writer = context.getWriter();

        HttpBindingIndex bindingIndex = model.getKnowledge(HttpBindingIndex.class);
        Set<MemberShape> errorDocumentBinding = bindingIndex.getResponseBindings(shapeId, HttpBinding.Location.DOCUMENT)
                .stream()
                .map(HttpBinding::getMember)
                .collect(Collectors.toSet());

        if (errorDocumentBinding.size() != 0) {
            writeDocumentBindingDeserializer(writer, model, context.getSymbolProvider(), shape,
                    errorDocumentBinding::contains, false);
            writer.write("");
        }

        Set<MemberShape> errorPayloadBinding = bindingIndex.getResponseBindings(shapeId, HttpBinding.Location.PAYLOAD)
                .stream()
                .map(HttpBinding::getMember)
                .collect(Collectors.toSet());

        // do not generate if no payload binding deserializer for Error Binding
        if (errorPayloadBinding.size() == 0) {
            return;
        }

        writePayloadBindingDeserializer(writer, model, context.getSymbolProvider(), shape,
                errorPayloadBinding::contains, false);
        writer.write("");
    }


    @Override
    protected void generateDocumentBodyShapeDeserializers(GenerationContext context, Set<Shape> shapes) {
        GoWriter writer = context.getWriter();
        Model model = context.getModel();
        SymbolProvider symbolProvider = context.getSymbolProvider();

        shapes.forEach(shape -> {
            writeDocumentBindingDeserializer(writer, model, symbolProvider, shape, FunctionalUtils.alwaysTrue(), false);
            writer.write("");
        });
    }


    // Generate deserializers for shapes with payload binding
    private void writePayloadBindingDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            Shape shape,
            Predicate<MemberShape> filterMemberShapes,
            Boolean isOutputShape
    ) {
        Symbol shapeSymbol = symbolProvider.toSymbol(shape);
        String funcName = isOutputShape ? ProtocolGenerator.getDocumentOutputDeserializerFunctionName(shape,
                getProtocolName()) : ProtocolGenerator.getDocumentDeserializerFunctionName(shape, getProtocolName());

        for (MemberShape memberShape : shape.members()) {
            if (!filterMemberShapes.test(memberShape)) {
                continue;
            }

            String memberName = symbolProvider.toMemberName(memberShape);
            Shape targetShape = model.expectShape(memberShape.getTarget());
            if (targetShape.isStringShape() || targetShape.isBlobShape()) {
                writer.openBlock("func $L(v $P, body io.ReadCloser) error {", "}",
                        funcName, shapeSymbol, () -> {
                            writer.openBlock("if v == nil {", "}", () -> {
                                writer.write("return fmt.Errorf(\"unsupported deserialization of nil %T\", v)");
                            });
                            writer.write("");

                            if (!targetShape.hasTrait(StreamingTrait.class)) {
                                writer.addUseImports(SmithyGoDependency.IOUTIL);
                                writer.write("bs, err := ioutil.ReadAll(body)");
                                writer.write("if err != nil { return err }");
                                writer.openBlock("if len(bs) > 0 {", "}", () -> {
                                    if (targetShape.isBlobShape()) {
                                        writer.write("v.$L = bs", memberName);
                                    } else { // string
                                        writer.addUseImports(SmithyGoDependency.SMITHY_PTR);
                                        if (targetShape.hasTrait(EnumTrait.class)) {
                                            writer.write("v.$L = string(bs)", memberName);
                                        } else {
                                            writer.write("v.$L = ptr.String(string(bs))", memberName);
                                        }
                                    }
                                });
                            } else {
                                writer.write("v.$L = body", memberName);
                            }

                            writer.write("return nil");
                        });
            } else {
                // delegate to Json Document Binding Deserializer
                writeDocumentBindingDeserializer(writer, model, symbolProvider, shape, filterMemberShapes,
                        isOutputShape);
            }
        }
    }

    // Generate deserializers for shape with document binding
    private void writeDocumentBindingDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            Shape shape,
            Predicate<MemberShape> filterMemberShapes,
            Boolean isOutputShape
    ) {
        Symbol jsonDecoder = SymbolUtils.createPointableSymbolBuilder("Decoder", SmithyGoDependency.JSON).build();
        Symbol shapeSymbol = symbolProvider.toSymbol(shape);
        String functionName = isOutputShape ?
                ProtocolGenerator.getDocumentOutputDeserializerFunctionName(shape, getProtocolName()) :
                ProtocolGenerator.getDocumentDeserializerFunctionName(shape, getProtocolName());

        writer.addUseImports(SmithyGoDependency.FMT);
        switch (shape.getType()) {
            case STRUCTURE:
                writer.openBlock("func $L(v *$P, decoder $P) error {", "}", functionName, shapeSymbol,
                        jsonDecoder, () -> {
                            writer.openBlock("if v == nil {", "}", () -> {
                                writer.write("return fmt.Errorf(\"unexpected nil of type %T\", v)");
                            });
                            writer.write("");
                            generateDocumentBindingStructureShapeDeserializer(writer, model, symbolProvider, shape,
                                    filterMemberShapes);
                            writer.write("");
                            writer.write("return nil");
                        });
                break;
            case SET:
            case LIST:
                writer.openBlock("func $L(vp *$P, decoder $P) error {", "}", functionName, shapeSymbol,
                        jsonDecoder, () -> {
                            writer.openBlock("if vp == nil {", "}", () -> {
                                writer.write("return fmt.Errorf(\"unexpected nil of type %T\", vp)");
                            });
                            writer.write("");
                            generateDocumentBindingCollectionShapeDeserializer(writer, model, symbolProvider, shape,
                                    filterMemberShapes);
                            writer.write("");
                            writer.write("return nil");
                        });
                break;
            case MAP:
                writer.openBlock("func $L(vp *$P, decoder $P) error {", "}", functionName, shapeSymbol,
                        jsonDecoder, () -> {
                            writer.openBlock("if vp == nil {", "}", () -> {
                                writer.write("return fmt.Errorf(\"unexpected nil of type %T\", vp)");
                            });
                            writer.write("");
                            generateDocumentBindingMapShapeDeserializer(writer, model, symbolProvider, shape,
                                    filterMemberShapes);
                            writer.write("");
                            writer.write("return nil");
                        });
                break;
            case DOCUMENT:
                writer.openBlock("func $L(v *$P, decoder $P) error {", "}", functionName, shapeSymbol,
                        jsonDecoder, () -> {
                            // TODO: Requires Document Implementation
                            writer.addUseImports(SmithyGoDependency.FMT);
                            writer.write("return fmt.Errorf(\"document types not implemented\")");
                        });
                break;
            case UNION:
                writer.openBlock("func $L(v *$P, decoder $P) error {", "}", functionName, shapeSymbol,
                        jsonDecoder, () -> {
                            // TODO: Requires Union Implementation
                            writer.addUseImports(SmithyGoDependency.FMT);
                            writer.write("return fmt.Errorf(\"union types not implemented\")");
                        });
                break;
            default:
                break;
        }
    }

    // Generates deserializers for structure Shapes
    private void generateDocumentBindingStructureShapeDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            Shape shape,
            Predicate<MemberShape> filterMemberShapes
    ) {
        Symbol symbol = symbolProvider.toSymbol(shape);
        writeJsonTokenizerStartStub(writer, shape);
        writer.write("var sv $P", symbol);
        writer.openBlock("if *v == nil {", "", () -> {
            writer.write("sv = &$T{}", symbol);
            writer.openBlock("} else {", "}", () -> {
                writer.write("sv = *v");
            });
        });
        writer.openBlock("for decoder.More() {", "}",
                () -> {
                    writer.write("t, err := decoder.Token()");
                    writer.write("if err != nil { return err }");
                    writer.openBlock("switch t {", "}", () -> {
                        for (MemberShape memberShape : shape.members()) {
                            if (!filterMemberShapes.test(memberShape)) {
                                continue;
                            }

                            String memberName = symbolProvider.toMemberName(memberShape);
                            writer.openBlock("case $S:", "", getSerializedMemberName(memberShape), () -> {
                                String operand = generateDocumentBindingMemberShapeDeserializer(writer, model,
                                        symbolProvider, memberShape);
                                writer.write(String.format("sv.%s = %s", memberName, operand));
                            });
                        }

                        // default case to handle unknown fields
                        writer.openBlock("default : ", "", () -> {
                            writer.addUseImports(AwsGoDependency.AWS_REST_JSON_PROTOCOL);
                            writer.write("err := restjson.DiscardUnknownField(decoder)");
                            writer.write("if err != nil {return err}");
                        });
                    });
                });
        writeJsonTokenizerEndStub(writer, shape);
        writer.write("*v = sv");
    }


    // Generates deserializers for collection shapes.
    private void generateDocumentBindingCollectionShapeDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            Shape shape,
            Predicate<MemberShape> filterMemberShapes
    ) {
        Symbol symbol = symbolProvider.toSymbol(shape);
        writeJsonTokenizerStartStub(writer, shape);
        writer.write("var cv $P", symbol);
        writer.openBlock("if *vp == nil {", "", () -> {
            writer.write("cv = $P{}", symbol);
            writer.openBlock("} else {", "}", () -> {
                writer.write("cv = *vp");
            });
        });
        writer.openBlock("for decoder.More() {", "}", () -> {
            MemberShape memberShape = shape.members().iterator().next();
            String operand = generateDocumentBindingMemberShapeDeserializer(writer, model, symbolProvider, memberShape);

            writer.write(String.format("cv = append(cv, %s)", operand));
            writer.write("");
        });
        writeJsonTokenizerEndStub(writer, shape);
        writer.write("*vp = cv");
    }

    // Generates deserializers for map shapes.
    private void generateDocumentBindingMapShapeDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            Shape shape,
            Predicate<MemberShape> filterMemberShapes
    ) {
        writeJsonTokenizerStartStub(writer, shape);
        Symbol symbol = symbolProvider.toSymbol(shape);
        writer.write("var mv $P", symbol);
        writer.openBlock("if *vp == nil {", "", () -> {
            writer.write("mv = $P{}", symbol);
            writer.openBlock("} else {", "}", () -> {
                writer.write("mv = *vp");
            });
        });
        writer.openBlock("for decoder.More() {", "}", () -> {
            MemberShape memberShape = shape.asMapShape().get().getValue();

            writer.write("token, err := decoder.Token()");
            writer.write("if err != nil { return err }");
            writer.write("");
            writer.write("key, ok := token.(string)");
            writer.write("if !ok { return fmt.Errorf(\"expected map-key of type string, found type %T\", token)}");
            writer.write("");

            String operand = generateDocumentBindingMemberShapeDeserializer(writer, model, symbolProvider, memberShape);
            writer.write(String.format("mv[key] = %s", operand));
            writer.write("");
        });

        writeJsonTokenizerEndStub(writer, shape);
        writer.write("*vp = mv");
    }

    // generateDocumentBindingMemberShapeDeserializer delegates to the appropriate
    // deserializer generator function for the member shapes.
    private String generateDocumentBindingMemberShapeDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            MemberShape memberShape
    ) {
        Shape targetShape = model.expectShape(memberShape.getTarget());
        switch (targetShape.getType()) {
            case STRING:
                return generateDocumentBindingStringMemberDeserializer(writer, model, symbolProvider, memberShape);
            case BOOLEAN:
                return generateDocumentBindingBooleanMemberDeserializer(writer, model, symbolProvider, memberShape);
            case TIMESTAMP:
                return generateDocumentBindingTimestampMemberDeserializer(writer, model, symbolProvider, memberShape);
            case BLOB:
                return generateDocumentBindingBlobMemberDeserializer(writer, model, symbolProvider, memberShape);
            case BYTE:
            case SHORT:
            case INTEGER:
            case LONG:
                return generateDocumentBindingIntegerMemberDeserializer(writer, model, symbolProvider, memberShape);
            case BIG_INTEGER:
            case BIG_DECIMAL:
                return generateDocumentBindingBigMemberDeserializer(writer, model, memberShape);
            case FLOAT:
            case DOUBLE:
                return generateDocumentBindingFloatMemberDeserializer(writer, model, memberShape);
            case SET:
            case LIST:
            case MAP:
                return generateDocumentBindingCollectionMemberDeserializer(writer, model, symbolProvider, memberShape);
            case STRUCTURE:
                return generateDocumentBindingStructureMemberDeserializer(writer, model, symbolProvider, memberShape);
            case UNION:
            case DOCUMENT:
                writer.write("// TODO: Support " + targetShape.getType() + " Deserialization");
                return "nil";
            default:
                throw new CodegenException("Unexpected shape deserialization to JSON");
        }
    }


    // Generates deserializer for String member shape.
    private String generateDocumentBindingStringMemberDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            MemberShape memberShape
    ) {
        String memberName = symbolProvider.toMemberName(memberShape);
        Shape targetShape = model.expectShape(memberShape.getTarget());
        Symbol targetSymbol = symbolProvider.toSymbol(targetShape);
        writer.addUseImports(SmithyGoDependency.FMT);
        writer.write("val, err := decoder.Token()");
        writer.write("if err != nil { return err }");

        BiConsumer<GoWriter, Consumer<GoWriter>> handleJsonString = (w, c) -> {
            w.openBlock("if val != nil {", "}", () -> {
                w.write("jtv, ok := val.(string)");
                w.openBlock("if !ok {", "}", () -> {
                    w.write("return fmt.Errorf(\"expected $L to be of type string, got %T instead\", val)",
                            memberName);
                });
                c.accept(w);
            });
        };

        if (targetShape.hasTrait(EnumTrait.class)) {
            writer.write("var ev $P", targetSymbol);
            handleJsonString.accept(writer, w -> {
                w.write("ev = $P(jtv)", targetSymbol);
            });
            return "ev";
        }

        writer.write("var st *string");
        handleJsonString.accept(writer, w -> {
            w.write("st = &jtv");
        });
        return "st";
    }

    // Generates deserializer for Boolean member shape.
    private String generateDocumentBindingBooleanMemberDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            MemberShape memberShape
    ) {
        String shapeName = symbolProvider.toMemberName(memberShape);
        writer.addUseImports(SmithyGoDependency.FMT);
        writer.write("val, err := decoder.Token()");
        writer.write("if err != nil { return err }");
        writer.write("var b *bool");
        writer.openBlock("if val != nil {", "}", () -> {
            writer.write("jtv, ok := val.(bool)");
            writer.openBlock("if !ok {", "}", () -> {
                writer.write("return fmt.Errorf(\"expected $L to be of type *bool, got %T instead\", val)",
                        shapeName);
            });
            writer.write("b = &jtv");
        });
        return "b";
    }

    // Generates deserializer for Byte, Short, Integer, Long member shape.
    private String generateDocumentBindingIntegerMemberDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            MemberShape memberShape
    ) {
        writer.addUseImports(SmithyGoDependency.FMT);
        writer.write("val, err := decoder.Token()");
        writer.write("if err != nil { return err }");

        BiConsumer<GoWriter, Consumer<GoWriter>> handleJsonNumber = (w, c) -> {
            w.openBlock("if val != nil {", "}", () -> {
                w.write("jtv, ok := val.(json.Number)");
                w.openBlock("if !ok {", "}", () -> {
                    w.write("return fmt.Errorf(\"expected $L to be json.Number, got %T instead\", val)",
                            symbolProvider.toMemberName(memberShape));
                });
                w.write("i64, err := jtv.Int64()");
                w.write("if err != nil { return err }");
                c.accept(w);
            });
        };

        Shape targetShape = model.expectShape(memberShape.getTarget());
        switch (targetShape.getType()) {
            case BYTE:
                writer.write("var st *int8");
                handleJsonNumber.accept(writer, w -> {
                    w.write("st = " + CodegenUtils.generatePointerValueIfPointable(w, targetShape, "int8(i64)"));
                });
                break;
            case SHORT:
                writer.write("var st *int16");
                handleJsonNumber.accept(writer, w -> {
                    w.write("st = " + CodegenUtils.generatePointerValueIfPointable(w, targetShape, "int16(i64)"));
                });
                break;
            case INTEGER:
                writer.write("var st *int32");
                handleJsonNumber.accept(writer, w -> {
                    w.write("st = " + CodegenUtils.generatePointerValueIfPointable(w, targetShape, "int32(i64)"));
                });
                break;
            case LONG:
                writer.write("var st *int64");
                handleJsonNumber.accept(writer, w -> {
                    w.write("st = &i64");
                });
                break;
            default:
                throw new CodegenException(
                        "unexpected integer number type " + targetShape.getType() + ", " + memberShape.getId());
        }

        return "st";
    }

    // Generates deserializer for Big Integer, Big Decimal member shape.
    private String generateDocumentBindingBigMemberDeserializer(
            GoWriter writer,
            Model model,
            MemberShape memberShape
    ) {
        writer.addUseImports(SmithyGoDependency.FMT);

        writer.write("val, err := decoder.Token()");
        writer.write("if err != nil { return err }");

        BiConsumer<GoWriter, Consumer<GoWriter>> handleJsonString = (w, c) -> {
            w.openBlock("if val != nil {", "}", () -> {
                w.write("jtv, ok := val.(string)");
                w.openBlock("if !ok {", "}", () -> {
                    w.write("return fmt.Errorf(\"expected $L to be of type string, got %T instead\", val)",
                            memberShape.getMemberName());
                });
                c.accept(w);
            });
        };


        writer.addUseImports(SmithyGoDependency.BIG);
        Shape targetShape = model.expectShape(memberShape.getTarget());
        switch (targetShape.getType()) {
            case BIG_INTEGER:
                writer.write("var bv *big.Int");
                handleJsonString.accept(writer, w -> {
                    w.write("bv, ok = new(big.Int).SetString(jtv, 10)");
                    w.write("if !ok { return fmt.Errorf(\"error deserializing big integer type\")}");
                });
                break;
            case BIG_DECIMAL:
                writer.write("var bv *big.Float");
                handleJsonString.accept(writer, w -> {
                    w.write("bv, _, err = new(big.Float).Parse(jtv, 10)");
                    w.write("if !ok { return fmt.Errorf(\"error deserializing big decimal type\")}");
                });
                break;
            default:
                throw new CodegenException(
                        "unexpected big number type " + targetShape.getType() + ", " + memberShape.getId());
        }

        return "bv";
    }

    // Generates deserializer for Float, Double member shape.
    private String generateDocumentBindingFloatMemberDeserializer(
            GoWriter writer,
            Model model,
            MemberShape memberShape
    ) {
        writer.addUseImports(SmithyGoDependency.FMT);

        writer.write("val, err := decoder.Token()");
        writer.write("if err != nil { return err }");

        BiConsumer<GoWriter, Consumer<GoWriter>> handleJsonFloat = (w, c) -> {
            w.openBlock("if val != nil {", "}", () -> {
                w.write("jtv, ok := val.(json.Number)");
                w.openBlock("if !ok {", "}", () -> {
                    w.write("return fmt.Errorf(\"expected $L to be of type json.Number, got %T instead\", val)",
                            memberShape.getMemberName());
                });
                w.write("f64, err := jtv.Float64()");
                w.write("if err != nil { return err }");
                c.accept(w);
            });
        };

        Shape targetShape = model.expectShape(memberShape.getTarget());
        switch (targetShape.getType()) {
            case FLOAT:
                writer.write("var st *float32");
                handleJsonFloat.accept(writer, w -> {
                    w.write("st = " + CodegenUtils.generatePointerValueIfPointable(w, targetShape, "float32(f64)"));
                });
                break;
            case DOUBLE:
                writer.write("var st *float64");
                handleJsonFloat.accept(writer, w -> {
                    w.write("st = &f64");
                });
                break;
            default:
                throw new CodegenException(
                        "unexpected decimal number type " + targetShape.getType() + ", " + memberShape.getId());
        }

        return "st";
    }

    // Generates deserializer for Timestamp member shape.
    private String generateDocumentBindingTimestampMemberDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            MemberShape memberShape
    ) {
        TimestampFormatTrait.Format format = memberShape.getMemberTrait(model, TimestampFormatTrait.class)
                .map(TimestampFormatTrait::getFormat)
                .orElse(TimestampFormatTrait.Format.EPOCH_SECONDS);

        writer.addUseImports(SmithyGoDependency.SMITHY_TIME);
        writer.addUseImports(SmithyGoDependency.FMT);
        writer.write("val, err := decoder.Token()");
        writer.write("if err != nil { return err }");

        BiConsumer<GoWriter, Consumer<GoWriter>> handleAsJsonString = (w, c) -> {
            w.openBlock("if val != nil {", "}", () -> {
                w.write("jtv, ok := val.(string)");
                w.openBlock("if !ok {", "}", () -> {
                    w.write("return fmt.Errorf(\"expected $L to be of type string, got %T instead\", val)",
                            symbolProvider.toMemberName(memberShape));
                });
                c.accept(w);
            });
        };
        BiConsumer<GoWriter, Consumer<GoWriter>> handleAsJsonNumber = (w, c) -> {
            w.openBlock("if val != nil {", "}", () -> {
                w.write("jtv, ok := val.(json.Number)");
                w.openBlock("if !ok {", "}", () -> {
                    w.write("return fmt.Errorf(\"expected $L to be of type json.Number, got %T instead\", val)",
                            memberShape.getMemberName());
                });
                c.accept(w);
            });
        };

        writer.write("var ts *time.Time");
        switch (format) {
            case DATE_TIME:
                handleAsJsonString.accept(writer, w -> {
                    w.write("t, err := smithytime.ParseDateTime(jtv)");
                    w.write("if err != nil { return err }");
                    w.write("ts = &t");
                });
                break;
            case HTTP_DATE:
                handleAsJsonString.accept(writer, w -> {
                    w.write("t, err := smithytime.ParseHTTPDate(jtv)");
                    w.write("if err != nil { return err }");
                    w.write("ts = &t");
                });
                break;
            case EPOCH_SECONDS:
                handleAsJsonNumber.accept(writer, w -> {
                    w.write("f64, err := jtv.Float64()");
                    w.write("if err != nil { return err }");
                    w.write("t := smithytime.ParseEpochSeconds(f64)");
                    w.write("ts = &t");
                });
                break;
            case UNKNOWN:
                throw new CodegenException("Unknown timestamp format");
        }
        return "ts";
    }

    // Generates deserializer for blob member shape.
    private String generateDocumentBindingBlobMemberDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            MemberShape memberShape
    ) {
        Shape targetShape = model.expectShape(memberShape.getTarget());
        Symbol targetSymbol = symbolProvider.toSymbol(targetShape);

        writer.write("var bs $T", targetSymbol);
        writer.write("err := decoder.Decode(&bs)");
        writer.write("if err != nil { return err }");
        return "bs";
    }

    // Generates deserializer for delegator for structure member shape.
    private String generateDocumentBindingStructureMemberDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            MemberShape memberShape
    ) {
        Shape targetShape = model.expectShape(memberShape.getTarget());
        Symbol targetSymbol = symbolProvider.toSymbol(targetShape);
        String deserFunctionName = ProtocolGenerator
                .getDocumentDeserializerFunctionName(targetShape, getProtocolName());
        writer.write("var val $P", targetSymbol);
        writer.openBlock("if err := $L(&val, decoder); err != nil {", "}", deserFunctionName, () -> {
            writer.write("return err");
        });
        return "val";
    }

    // Generates deserializer for delegator for collection member shape and map member shapes.
    private String generateDocumentBindingCollectionMemberDeserializer(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            MemberShape memberShape
    ) {
        Shape targetShape = model.expectShape(memberShape.getTarget());
        Symbol targetSymbol = symbolProvider.toSymbol(targetShape);

        String deserializerFuncName = ProtocolGenerator
                .getDocumentDeserializerFunctionName(targetShape, getProtocolName());
        writer.write("var col $P", targetSymbol);
        writer.openBlock("if err := $L(&col, decoder); err != nil {", "}", deserializerFuncName, () -> {
            writer.write("return err");
        });
        return "col";
    }

    // generates Json decoder tokenizer start stub wrt to the shape
    private void writeJsonTokenizerStartStub(GoWriter writer, Shape shape) {
        String startToken = shape.isListShape() || shape.isSetShape() ? "[" : "{";
        writer.write("startToken, err := decoder.Token()");
        writer.write("if err == io.EOF { return nil }");
        writer.write("if err != nil { return err }");
        writer.write("if startToken == nil { return nil }");
        writer.openBlock("if t, ok := startToken.(json.Delim); !ok || t != $L {", "}",
                String.format("'%s'", startToken), () -> {
                    writer.addUseImports(SmithyGoDependency.FMT);
                    writer.write("return fmt.Errorf($S)", String.format("expect `%s` as start token", startToken));
                });
        writer.write("");
    }

    // generates Json decoder tokenizer end stub wrt to the shape
    private void writeJsonTokenizerEndStub(GoWriter writer, Shape shape) {
        String endToken = shape.isListShape() || shape.isSetShape() ? "]" : "}";
        writer.write("");
        writer.write("endToken, err := decoder.Token()");
        writer.write("if err != nil { return err }");
        writer.openBlock("if t, ok := endToken.(json.Delim); !ok || t != $L {", "}", String.format("'%s'", endToken),
                () -> {
                    writer.write("return fmt.Errorf($S)", String.format("expect `%s` as end token", endToken));
                });
    }


    @Override
    public void generateSharedDeserializerComponents(GenerationContext context) {
        super.generateSharedDeserializerComponents(context);
    }
}
