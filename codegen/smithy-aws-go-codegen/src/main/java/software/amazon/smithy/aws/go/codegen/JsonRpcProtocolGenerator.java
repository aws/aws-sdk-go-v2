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

import static software.amazon.smithy.aws.go.codegen.AwsProtocolUtils.handleDecodeError;
import static software.amazon.smithy.aws.go.codegen.AwsProtocolUtils.initializeJsonDecoder;
import static software.amazon.smithy.aws.go.codegen.AwsProtocolUtils.writeJsonErrorMessageCodeDeserializer;

import java.util.HashSet;
import java.util.Set;
import java.util.TreeSet;
import java.util.stream.Collectors;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.CodegenUtils;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.HttpRpcProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.EventStreamInfo;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.shapes.UnionShape;
import software.amazon.smithy.model.traits.ErrorTrait;
import software.amazon.smithy.model.traits.EventHeaderTrait;
import software.amazon.smithy.model.traits.EventPayloadTrait;

/**
 * Handles generating the aws.rest-json protocol for services.
 *
 * @inheritDoc
 * @see HttpRpcProtocolGenerator
 */
abstract class JsonRpcProtocolGenerator extends HttpRpcProtocolGenerator {
    private final Set<ShapeId> generatedDocumentBodyShapeSerializers = new HashSet<>();
    private final Set<ShapeId> generatedEventMessageSerializers = new HashSet<>();
    private final Set<ShapeId> generatedDocumentBodyShapeDeserializers = new HashSet<>();
    private final Set<ShapeId> generatedEventMessageDeserializers = new HashSet<>();

    /**
     * Creates an AWS JSON RPC protocol generator
     */
    public JsonRpcProtocolGenerator() {
        super();
    }

    @Override
    protected String getOperationPath(GenerationContext context, OperationShape operation) {
        return "/";
    }

    @Override
    protected void writeDefaultHeaders(GenerationContext context, OperationShape operation, GoWriter writer) {
        super.writeDefaultHeaders(context, operation, writer);
        ServiceShape service = context.getService();
        String target = service.getId().getName(service) + "." + operation.getId().getName(service);
        writer.write("httpBindingEncoder.SetHeader(\"X-Amz-Target\").String($S)", target);
    }

    @Override
    protected void serializeInputDocument(GenerationContext context, OperationShape operation) {
        GoWriter writer = context.getWriter().get();

        // Stub synthetic clone inputs mean there never was an input modeled, always serialize empty JSON object
        // as place holder.
        if (CodegenUtils.isStubSyntheticClone(ProtocolUtils.expectInput(context.getModel(), operation))) {
            writer.addUseImports(SmithyGoDependency.STRINGS);
            writer.openBlock("if request, err = request.SetStream(strings.NewReader(`{}`)); err != nil {",
                    "}", () -> {
                        writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
                    });
            return;
        }

        StructureShape input = ProtocolUtils.expectInput(context.getModel(), operation);
        String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(input, context.getService(), getProtocolName());

        writer.addUseImports(SmithyGoDependency.SMITHY_JSON);
        writer.write("jsonEncoder := smithyjson.NewEncoder()");
        writer.openBlock("if err := $L(input, jsonEncoder.Value); err != nil {", "}", functionName, () -> {
            writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
        }).write("");

        writer.addUseImports(SmithyGoDependency.BYTES);
        writer.openBlock("if request, err = request.SetStream(bytes.NewReader(jsonEncoder.Bytes())); err != nil {",
                "}", () -> {
                    writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
                });
    }

    @Override
    protected void generateDocumentBodyShapeSerializers(GenerationContext context, Set<Shape> shapes) {
        JsonShapeSerVisitor visitor = new JsonShapeSerVisitor(context);
        shapes.forEach(shape -> {
            if (generatedDocumentBodyShapeSerializers.contains(shape.toShapeId())) {
                return;
            }
            shape.accept(visitor);
            generatedDocumentBodyShapeSerializers.add(shape.toShapeId());
        });
    }

    @Override
    protected void deserializeOutputDocument(GenerationContext context, OperationShape operation) {
        GoWriter writer = context.getWriter().get();
        StructureShape output = ProtocolUtils.expectOutput(context.getModel(), operation);
        String functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(output, context.getService(), getProtocolName());
        initializeJsonDecoder(writer, "response.Body");
        AwsProtocolUtils.decodeJsonIntoInterface(writer, "out, metadata, ");
        writer.write("err = $L(&output, shape)", functionName);
        handleDecodeError(writer, "out, metadata, ");
    }


    @Override
    protected void generateDocumentBodyShapeDeserializers(GenerationContext context, Set<Shape> shapes) {
        JsonShapeDeserVisitor visitor = new JsonShapeDeserVisitor(context);
        shapes.forEach(shape -> {
            if (generatedDocumentBodyShapeDeserializers.contains(shape.toShapeId())) {
                return;
            }
            shape.accept(visitor);
            generatedDocumentBodyShapeDeserializers.add(shape.toShapeId());
        });
    }


    @Override
    protected void deserializeError(GenerationContext context, StructureShape shape) {
        GoWriter writer = context.getWriter().get();
        Symbol symbol = context.getSymbolProvider().toSymbol(shape);
        String functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(shape, context.getService(), getProtocolName());

        initializeJsonDecoder(writer, "errorBody");
        AwsProtocolUtils.decodeJsonIntoInterface(writer, "");
        writer.write("output := &$T{}", symbol);
        writer.write("err := $L(&output, shape)", functionName);
        writer.write("");
        handleDecodeError(writer);
        writer.write("errorBody.Seek(0, io.SeekStart)");
        writer.write("return output");
    }

    @Override
    public void generateProtocolTests(GenerationContext context) {
        AwsProtocolUtils.generateHttpProtocolTests(context);
    }

    @Override
    protected void writeErrorMessageCodeDeserializer(GenerationContext context) {
        writeJsonErrorMessageCodeDeserializer(context);
    }

    @Override
    public void generateProtocolDocumentMarshalerUnmarshalDocument(GenerationContext context) {
        JsonProtocolDocumentUtils.generateProtocolDocumentMarshalerUnmarshalDocument(context);
    }

    @Override
    public void generateProtocolDocumentMarshalerMarshalDocument(GenerationContext context) {
        JsonProtocolDocumentUtils.generateProtocolDocumentMarshalerMarshalDocument(context);
    }

    @Override
    public void generateProtocolDocumentUnmarshalerUnmarshalDocument(GenerationContext context) {
        JsonProtocolDocumentUtils.generateProtocolDocumentUnmarshalerUnmarshalDocument(context);
    }

    @Override
    public void generateProtocolDocumentUnmarshalerMarshalDocument(GenerationContext context) {
        JsonProtocolDocumentUtils.generateProtocolDocumentUnmarshalerMarshalDocument(context);
    }

    @Override
    public void generateEventStreamComponents(GenerationContext context) {
        AwsEventStreamUtils.generateEventStreamComponents(context);
    }

    @Override
    protected void writeOperationSerializerMiddlewareEventStreamSetup(GenerationContext context, EventStreamInfo info) {
        AwsEventStreamUtils.writeOperationSerializerMiddlewareEventStreamSetup(context, info);
    }

    @Override
    protected void generateEventStreamSerializers(
            GenerationContext context,
            UnionShape eventUnion,
            Set<EventStreamInfo> eventStreamInfos
    ) {
        Model model = context.getModel();

        AwsEventStreamUtils.generateEventStreamSerializer(context, eventUnion);
        var memberShapes = eventUnion.members().stream()
                .filter(ms -> ms.getMemberTrait(model, ErrorTrait.class).isEmpty())
                .collect(Collectors.toCollection(TreeSet::new));

        final var eventDocumentShapes = new HashSet<Shape>();
        for (MemberShape member : memberShapes) {
            var targetShape = model.expectShape(member.getTarget());
            if (generatedEventMessageSerializers.contains(targetShape.toShapeId())) {
                continue;
            }

            AwsEventStreamUtils.generateEventMessageSerializer(context, targetShape, (ctx, payloadTarget, operand) -> {
                var functionName = ProtocolGenerator.getDocumentSerializerFunctionName(payloadTarget,
                        ctx.getService(), ctx.getProtocolName());
                AwsProtocolUtils.writeJsonEventMessageSerializerDelegator(ctx, functionName, operand,
                        getDocumentContentType());
            });

            generatedEventMessageSerializers.add(targetShape.toShapeId());

            var hasBindings = targetShape.members().stream()
                    .filter(ms -> ms.getTrait(EventHeaderTrait.class).isPresent()
                                  || ms.getTrait(EventPayloadTrait.class).isPresent())
                    .findAny();
            if (hasBindings.isPresent()) {
                var payload = targetShape.members().stream()
                        .filter(ms -> ms.getTrait(EventPayloadTrait.class).isPresent())
                        .map(ms -> model.expectShape(ms.getTarget()))
                        .filter(ProtocolUtils::requiresDocumentSerdeFunction)
                        .findAny();
                payload.ifPresent(eventDocumentShapes::add);
                continue;
            }
            eventDocumentShapes.add(targetShape);
        }

        eventDocumentShapes.addAll(ProtocolUtils.resolveRequiredDocumentShapeSerde(model, eventDocumentShapes));
        generateDocumentBodyShapeSerializers(context, eventDocumentShapes);

        for (EventStreamInfo streamInfo : eventStreamInfos) {
            var inputShape = model.expectShape(streamInfo.getOperation().getInput().get());
            var functionName = ProtocolGenerator.getDocumentSerializerFunctionName(inputShape,
                    context.getService(), context.getProtocolName());
            AwsEventStreamUtils.generateEventMessageRequestSerializer(context, inputShape,
                    (ctx, payloadTarget, operand) -> {
                        AwsProtocolUtils.writeJsonEventMessageSerializerDelegator(ctx, functionName, operand,
                                getDocumentContentType());
                    });
            var initialMessageMembers = streamInfo.getInitialMessageMembers()
                    .values();
            inputShape.accept(new JsonShapeSerVisitor(context, initialMessageMembers::contains,
                    (shape, serviceShape, proto) -> functionName));
        }
    }

    @Override
    protected void generateEventStreamDeserializers(
            GenerationContext context,
            UnionShape eventUnion,
            Set<EventStreamInfo> eventStreamInfos
    ) {
        var model = context.getModel();

        AwsEventStreamUtils.generateEventStreamDeserializer(context, eventUnion);
        AwsEventStreamUtils.generateEventStreamExceptionDeserializer(context, eventUnion,
                AwsProtocolUtils::writeJsonEventStreamUnknownExceptionDeserializer);

        final var eventDocumentShapes = new HashSet<Shape>();

        for (MemberShape shape : eventUnion.members()) {
            var targetShape = model.expectShape(shape.getTarget());
            if (generatedEventMessageDeserializers.contains(targetShape.toShapeId())) {
                continue;
            }
            generatedEventMessageDeserializers.add(targetShape.toShapeId());
            if (shape.getMemberTrait(model, ErrorTrait.class).isPresent()) {
                AwsEventStreamUtils.generateEventMessageExceptionDeserializer(context, targetShape,
                        (ctx, payloadTarget) -> {
                            AwsProtocolUtils.initializeJsonEventMessageDeserializer(ctx);
                            var functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(
                                    payloadTarget, ctx.getService(), getProtocolName());
                            var ctxWriter = ctx.getWriter().get();
                            ctxWriter.write("v := &$T{}", ctx.getSymbolProvider().toSymbol(payloadTarget))
                                    .openBlock("if err := $L(&v, shape); err != nil {", "}", functionName,
                                            () -> handleDecodeError(ctxWriter))
                                    .write("return v");
                        });

                eventDocumentShapes.add(targetShape);
            } else {
                AwsEventStreamUtils.generateEventMessageDeserializer(context, targetShape,
                        (ctx, payloadTarget, operand) -> {
                            AwsProtocolUtils.initializeJsonEventMessageDeserializer(ctx);
                            var functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(
                                    payloadTarget, ctx.getService(), getProtocolName());
                            var ctxWriter = ctx.getWriter().get();
                            ctxWriter.openBlock("if err := $L(&$L, shape); err != nil {", "}", functionName, operand,
                                            () -> handleDecodeError(ctxWriter))
                                    .write("return nil");
                        });

                var hasBindings = targetShape.members().stream()
                        .filter(ms -> ms.getTrait(EventHeaderTrait.class).isPresent()
                                      || ms.getTrait(EventPayloadTrait.class).isPresent())
                        .findAny();
                if (hasBindings.isPresent()) {
                    var payload = targetShape.members().stream()
                            .filter(ms -> ms.getTrait(EventPayloadTrait.class).isPresent())
                            .map(ms -> model.expectShape(ms.getTarget()))
                            .filter(ProtocolUtils::requiresDocumentSerdeFunction)
                            .findAny();
                    payload.ifPresent(eventDocumentShapes::add);
                    continue;
                }
                eventDocumentShapes.add(targetShape);
            }
        }

        eventDocumentShapes.addAll(ProtocolUtils.resolveRequiredDocumentShapeSerde(model, eventDocumentShapes));
        generateDocumentBodyShapeDeserializers(context, eventDocumentShapes);

        for (EventStreamInfo streamInfo : eventStreamInfos) {
            var outputShape = model.expectShape(streamInfo.getOperation().getOutput().get());
            var functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(outputShape,
                    context.getService(), context.getProtocolName());
            AwsEventStreamUtils.generateEventMessageRequestDeserializer(context, outputShape,
                    (ctx, payloadTarget, operand) -> {
                        AwsProtocolUtils.initializeJsonEventMessageDeserializer(ctx, "nil,");
                        var ctxWriter = ctx.getWriter().get();
                        ctxWriter.openBlock("if err := $L(&$L, shape); err != nil {", "}", functionName, operand,
                                        () -> handleDecodeError(ctxWriter, "nil,"))
                                .write("return v, nil");
                    });
            var initialMessageMembers = streamInfo.getInitialMessageMembers()
                    .values();
            outputShape.accept(new JsonShapeDeserVisitor(context, initialMessageMembers::contains,
                    (shape, serviceShape, proto) -> functionName));
        }
    }

}
