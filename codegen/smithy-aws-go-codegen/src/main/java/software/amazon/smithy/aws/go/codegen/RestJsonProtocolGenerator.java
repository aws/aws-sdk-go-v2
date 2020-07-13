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

import static software.amazon.smithy.go.codegen.integration.ProtocolUtils.writeSafeMemberAccessor;

import java.util.Collection;
import java.util.Optional;
import java.util.Set;
import java.util.function.Predicate;
import java.util.stream.Collectors;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.HttpBindingProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.HttpBinding;
import software.amazon.smithy.model.knowledge.HttpBindingIndex;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.model.traits.HttpErrorTrait;
import software.amazon.smithy.model.traits.MediaTypeTrait;
import software.amazon.smithy.model.traits.StreamingTrait;
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

        Shape inputShape = ProtocolUtils.expectInput(model, operation);
        inputShape.accept(new JsonShapeSerVisitor(context, documentBindings::contains));
    }

    @Override
    protected void writeMiddlewarePayloadSerializerDelegator(
            GenerationContext context,
            OperationShape operation,
            MemberShape memberShape,
            GoStackStepMiddlewareGenerator generator
    ) {
        GoWriter writer = context.getWriter();
        Model model = context.getModel();
        Shape payloadShape = model.expectShape(memberShape.getTarget());

        writeSafeMemberAccessor(context, memberShape, "input", s -> {
            writer.openBlock("if !restEncoder.HasHeader(\"Content-Type\") {", "}", () -> {
                writer.write("restEncoder.SetHeader(\"Content-Type\").String($S)", getPayloadShapeMediaType(payloadShape));
            });
            writer.write("");

            if (payloadShape.hasTrait(StreamingTrait.class)) {
                writer.write("payload := $L", s);

            } else if (payloadShape.isBlobShape()) {
                writer.addUseImports(SmithyGoDependency.BYTES);
                writer.write("payload := bytes.NewReader($L)", s);

            } else if (payloadShape.isStringShape()) {
                writer.addUseImports(SmithyGoDependency.STRINGS);
                writer.write("payload := strings.NewReader(*$L)", s);

            } else {
                String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(payloadShape,
                        getProtocolName());
                writer.addUseImports(SmithyGoDependency.SMITHY_JSON);
                writer.write("jsonEncoder := smithyjson.NewEncoder()");
                writer.openBlock("if err := $L($L, jsonEncoder.Value); err != nil {", "}", functionName,
                        s, () -> {
                            writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
                        });
                writer.write("payload := bytes.NewReader(jsonEncoder.Bytes())");
            }

            writer.openBlock("if request, err = request.SetStream(payload); err != nil {", "}",
                    () -> {
                        writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
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
            return "text/plain";
        }

        return getDocumentContentType();
    }


    @Override
    protected void writeMiddlewareDocumentSerializerDelegator(
            GenerationContext context,
            OperationShape operation,
            GoStackStepMiddlewareGenerator generator
    ) {
        GoWriter writer = context.getWriter();
        writer.addUseImports(SmithyGoDependency.SMITHY);
        writer.addUseImports(SmithyGoDependency.SMITHY_JSON);

        writer.write("restEncoder.SetHeader(\"Content-Type\").String($S)", getDocumentContentType());
        writer.write("");

        Shape inputShape = ProtocolUtils.expectInput(context.getModel(), operation);
        String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(inputShape, getProtocolName());

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
        JsonShapeSerVisitor visitor = new JsonShapeSerVisitor(context);
        shapes.forEach(shape -> shape.accept(visitor));
    }

    @Override
    protected void writeMiddlewareDocumentDeserializerDelegator(
            GenerationContext context,
            OperationShape operation,
            GoStackStepMiddlewareGenerator generator
    ) {
        Model model = context.getModel();
        GoWriter writer = context.getWriter();
        Shape targetShape = ProtocolUtils.expectOutput(model, operation);
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
                writeMiddlewarePayloadBindingDeserializerDelegator(writer, targetShape);
                return;
            }
            // for other payload target types we should deserialize using the appropriate document deserializer
            targetShape = payloadShape;
            operand += "." + context.getSymbolProvider().toMemberName(memberShape);
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

        writeMiddlewareDocumentBindingDeserializerDelegator(writer, targetShape, operand);
    }

    // Writes middleware that delegates to deserializers for shapes that have explicit payload.
    private void writeMiddlewarePayloadBindingDeserializerDelegator(GoWriter writer, Shape shape) {
        String deserFuncName = ProtocolGenerator.getDocumentDeserializerFunctionName(shape, getProtocolName());
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
            String operand
    ) {
        String deserFuncName = ProtocolGenerator.getDocumentDeserializerFunctionName(shape, getProtocolName());
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
            GenerationContext context,
            OperationShape operationShape,
            GoStackStepMiddlewareGenerator generator
    ) {
        GoWriter writer = context.getWriter();
        Model model = context.getModel();
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
            writeErrorShapeDeserializerDelegator(writer, model, context.getSymbolProvider(), ErrorShapeIds);
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
                    writeMiddlewareDocumentBindingDeserializerDelegator(writer, errorShape, "output");
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

        Shape outputShape = ProtocolUtils.expectOutput(model, operation);
        GoWriter writer = context.getWriter();

        if (documentBindings.size() != 0) {
            outputShape.accept(new JsonShapeDeserVisitor(context, documentBindings::contains));
        }

        Set<MemberShape> payloadBindings = bindingIndex.getResponseBindings(operation, HttpBinding.Location.PAYLOAD)
                .stream()
                .map(HttpBinding::getMember)
                .collect(Collectors.toSet());

        if (payloadBindings.size() == 0) {
            return;
        }

        writePayloadBindingDeserializer(context, outputShape, payloadBindings::contains);
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
            shape.accept(new JsonShapeDeserVisitor(context, errorDocumentBinding::contains));
        }

        Set<MemberShape> errorPayloadBinding = bindingIndex.getResponseBindings(shapeId, HttpBinding.Location.PAYLOAD)
                .stream()
                .map(HttpBinding::getMember)
                .collect(Collectors.toSet());

        // do not generate if no payload binding deserializer for Error Binding
        if (errorPayloadBinding.size() == 0) {
            return;
        }

        writePayloadBindingDeserializer(context, shape, errorPayloadBinding::contains);
        writer.write("");
    }


    @Override
    protected void generateDocumentBodyShapeDeserializers(GenerationContext context, Set<Shape> shapes) {
        JsonShapeDeserVisitor visitor = new JsonShapeDeserVisitor(context);
        shapes.forEach(shape -> shape.accept(visitor));
    }


    // Generate deserializers for shapes with payload binding
    private void writePayloadBindingDeserializer(
            GenerationContext context,
            Shape shape,
            Predicate<MemberShape> filterMemberShapes
    ) {
        GoWriter writer = context.getWriter();
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol shapeSymbol = symbolProvider.toSymbol(shape);
        String funcName = ProtocolGenerator.getDocumentDeserializerFunctionName(shape, getProtocolName());

        for (MemberShape memberShape : shape.members()) {
            if (!filterMemberShapes.test(memberShape)) {
                continue;
            }

            String memberName = symbolProvider.toMemberName(memberShape);
            Shape targetShape = context.getModel().expectShape(memberShape.getTarget());
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
                shape.accept(new JsonShapeDeserVisitor(context, filterMemberShapes));
            }
        }
    }

    @Override
    public void generateSharedDeserializerComponents(GenerationContext context) {
        super.generateSharedDeserializerComponents(context);
    }

    @Override
    public void generateProtocolTests(GenerationContext context) {
        AwsProtocolUtils.generateHttpProtocolTests(context);
    }
}
