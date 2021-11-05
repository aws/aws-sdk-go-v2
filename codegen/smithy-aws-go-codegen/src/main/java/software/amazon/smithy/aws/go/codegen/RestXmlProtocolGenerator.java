package software.amazon.smithy.aws.go.codegen;


import static software.amazon.smithy.go.codegen.integration.HttpProtocolGeneratorUtils.isShapeWithResponseBindings;
import static software.amazon.smithy.aws.go.codegen.XmlProtocolUtils.initializeXmlDecoder;
import static software.amazon.smithy.aws.go.codegen.XmlProtocolUtils.handleDecodeError;
import static software.amazon.smithy.aws.go.codegen.XmlProtocolUtils.writeXmlErrorMessageCodeDeserializer;
import static software.amazon.smithy.aws.go.codegen.XmlProtocolUtils.generateXMLStartElement;
import static software.amazon.smithy.aws.go.codegen.XmlProtocolUtils.generatePayloadAsDocumentXMLStartElement;

import java.util.HashSet;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;
import java.util.TreeSet;
import java.util.function.Predicate;
import software.amazon.smithy.aws.traits.protocols.RestXmlTrait;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.HttpBindingProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.EventStreamInfo;
import software.amazon.smithy.model.knowledge.HttpBinding;
import software.amazon.smithy.model.knowledge.HttpBindingIndex;
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
import software.amazon.smithy.model.traits.MediaTypeTrait;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.model.traits.StreamingTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.model.traits.XmlAttributeTrait;
import software.amazon.smithy.model.traits.XmlNamespaceTrait;

abstract class RestXmlProtocolGenerator extends HttpBindingProtocolGenerator {
    private final Set<ShapeId> generatedDocumentBodyShapeSerializers = new HashSet<>();
    private final Set<ShapeId> generatedEventMessageSerializers = new HashSet<>();
    private final Set<ShapeId> generatedDocumentBodyShapeDeserializers = new HashSet<>();
    private final Set<ShapeId> generatedEventMessageDeserializers = new HashSet<>();

    /**
     * Creates a AWS REST XML protocol generator.
     */
    RestXmlProtocolGenerator() {
        super(true);
    }

    @Override
    public void generateProtocolTests(GenerationContext context) {
        AwsProtocolUtils.generateHttpProtocolTests(context);
    }

    @Override
    protected TimestampFormatTrait.Format getDocumentTimestampFormat() {
        return TimestampFormatTrait.Format.DATE_TIME;
    }


    @Override
    protected void generateOperationDocumentSerializer(GenerationContext context, OperationShape operation) {
        Model model = context.getModel();
        HttpBindingIndex bindingIndex = HttpBindingIndex.of(model);

        Set<MemberShape> documentBindings = bindingIndex.getRequestBindings(operation, HttpBinding.Location.DOCUMENT)
                .stream()
                .map(HttpBinding::getMember)
                .collect(Collectors.toSet());

        if (documentBindings.isEmpty()) {
            return;
        }
        Shape inputShape = ProtocolUtils.expectInput(context.getModel(), operation);
        inputShape.accept(new XmlShapeSerVisitor(context,
                memberShape -> documentBindings.contains(memberShape) && !memberShape.hasTrait(
                        XmlAttributeTrait.class)));
    }

    @Override
    protected void writeMiddlewareDocumentSerializerDelegator(
            GenerationContext context,
            OperationShape operation,
            GoStackStepMiddlewareGenerator generator
    ) {
        GoWriter writer = context.getWriter().get();
        writer.addUseImports(SmithyGoDependency.SMITHY);
        writer.addUseImports(SmithyGoDependency.SMITHY_XML);

        writer.write("restEncoder.SetHeader(\"Content-Type\").String($S)", getDocumentContentType());
        writer.write("");

        Shape inputShape = ProtocolUtils.expectInput(context.getModel(), operation);
        String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(inputShape, context.getService(), getProtocolName());

        initalizeXmlEncoder(context, writer, inputShape, "root", "input");

        writer.openBlock("if err := $L(input, xmlEncoder.RootElement(root)); err != nil {", "}",
                functionName, () -> {
                    writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
                });
        writer.insertTrailingNewline();

        writer.openBlock("if request, err = request.SetStream(bytes.NewReader(xmlEncoder.Bytes())); "
                + "err != nil {", "}", () -> {
            writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
        });
    }

    private void initalizeXmlEncoder(
            GenerationContext context,
            GoWriter writer,
            Shape inputShape,
            String nodeDst,
            String inputSrc
    ) {
        writer.addUseImports(SmithyGoDependency.SMITHY_XML);
        writer.addUseImports(SmithyGoDependency.BYTES);
        writer.write("xmlEncoder := smithyxml.NewEncoder(bytes.NewBuffer(nil))");

        generateXMLStartElement(context, inputShape, nodeDst, inputSrc);

        // check if service shape is bound by xmlNameSpace Trait
        Optional<XmlNamespaceTrait> xmlNamespaceTrait = context.getService().getTrait(XmlNamespaceTrait.class);
        if (xmlNamespaceTrait.isPresent()) {
            XmlNamespaceTrait namespace = xmlNamespaceTrait.get();
            writer.write("$L.Attr = append($L.Attr, smithyxml.NewNamespaceAttribute($S, $S))", nodeDst, nodeDst,
                    namespace.getPrefix().isPresent() ? namespace.getPrefix().get() : "", namespace.getUri()
            );
        }
    }

    @Override
    protected void writeMiddlewarePayloadAsDocumentSerializerDelegator(
            GenerationContext context,
            MemberShape memberShape,
            String operand
    ) {
        GoWriter writer = context.getWriter().get();
        Model model = context.getModel();
        Shape payloadShape = model.expectShape(memberShape.getTarget());

        writer.addUseImports(SmithyGoDependency.SMITHY_XML);
        writer.addUseImports(SmithyGoDependency.BYTES);
        writer.write("xmlEncoder := smithyxml.NewEncoder(bytes.NewBuffer(nil))");

        generatePayloadAsDocumentXMLStartElement(context, memberShape, "payloadRoot", operand);

        // check if service shape is bound by xmlNameSpace Trait
        Optional<XmlNamespaceTrait> xmlNamespaceTrait = context.getService().getTrait(XmlNamespaceTrait.class);
        if (xmlNamespaceTrait.isPresent()) {
            XmlNamespaceTrait namespace = xmlNamespaceTrait.get();
            writer.write("payloadRoot.Attr = append(payloadRoot.Attr, smithyxml.NewNamespaceAttribute($S, $S))",
                    namespace.getPrefix().isPresent() ? namespace.getPrefix().get() : "", namespace.getUri()
            );
        }

        String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(
                payloadShape, context.getService(), getProtocolName());
        writer.openBlock("if err := $L($L, xmlEncoder.RootElement(payloadRoot)); err != nil {", "}", functionName,
                operand, () -> {
                    writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
                });
        writer.write("payload := bytes.NewReader(xmlEncoder.Bytes())");
    }

    @Override
    protected void generateDocumentBodyShapeSerializers(GenerationContext context, Set<Shape> shapes) {
        // filter shapes marked as attributes
        XmlShapeSerVisitor visitor = new XmlShapeSerVisitor(context, memberShape -> !memberShape.hasTrait(XmlAttributeTrait.class));
        shapes.forEach(shape -> {
            if (generatedDocumentBodyShapeSerializers.contains(shape.toShapeId())) {
                return;
            }
            shape.accept(visitor);
            generatedDocumentBodyShapeSerializers.add(shape.toShapeId());
        });
    }

    /**
     * Returns the MediaType for the payload shape derived from the MediaTypeTrait, shape type, or document content type.
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

    /*     ================Deserializer===========================     */

    @Override
    protected void deserializeError(GenerationContext context, StructureShape shape) {
        GoWriter writer = context.getWriter().get();
        Symbol symbol = context.getSymbolProvider().toSymbol(shape);

        writer.write("output := &$T{}", symbol);
        writer.insertTrailingNewline();

        if (isShapeWithRestResponseBindings(context.getModel(), shape)) {
            String bindingDeserFunctionName = ProtocolGenerator.getOperationHttpBindingsDeserFunctionName(
                    shape, context.getService(), getProtocolName());
            writer.openBlock("if err := $L(output, response); err != nil {", "}", bindingDeserFunctionName, () -> {
                writer.addUseImports(SmithyGoDependency.SMITHY);
                writer.write(String.format("return &smithy.DeserializationError{Err: %s}",
                        "fmt.Errorf(\"failed to decode response error with invalid HTTP bindings, %w\", err)"));
            });
            writer.insertTrailingNewline();
        }

        if (isShapeWithResponseBindings(context.getModel(), shape, HttpBinding.Location.DOCUMENT)) {
            String documentDeserFunctionName = ProtocolGenerator.getDocumentDeserializerFunctionName(
                    shape, context.getService(), getProtocolName());

            initializeXmlDecoder(writer, "errorBody", "output");
            boolean isNoErrorWrapping = isNoErrorWrapping(context);
            Runnable writeErrorDelegator = () -> {
                writer.write("err = $L(&output, decoder)", documentDeserFunctionName);
                handleDecodeError(writer, "");
                writer.insertTrailingNewline();
            };

            if (isNoErrorWrapping) {
                writeErrorDelegator.run();
            } else {
                writer.write("t, err = decoder.GetElement(\"Error\")");
                XmlProtocolUtils.handleDecodeError(writer, "");
                Symbol wrapNodeDecoder = SymbolUtils.createValueSymbolBuilder("WrapNodeDecoder",
                        SmithyGoDependency.SMITHY_XML).build();
                writer.write("decoder = $T(decoder.Decoder, t)", wrapNodeDecoder);
                writeErrorDelegator.run();
            }
        }

        writer.write("return output");
    }

    private Boolean isNoErrorWrapping(GenerationContext context) {
        return context.getService().getTrait(RestXmlTrait.class).map(
                RestXmlTrait::isNoErrorWrapping).orElse(false);
    }

    @Override
    protected void writeErrorMessageCodeDeserializer(GenerationContext context) {
        writeXmlErrorMessageCodeDeserializer(context);
    }

    @Override
    protected void writeMiddlewareDocumentDeserializerDelegator(
            GenerationContext context,
            OperationShape operation,
            GoStackStepMiddlewareGenerator generator
    ) {
        Model model = context.getModel();
        GoWriter writer = context.getWriter().get();
        Shape targetShape = ProtocolUtils.expectOutput(model, operation);
        String operand = "output";

        if (isShapeWithResponseBindings(model, operation, HttpBinding.Location.PAYLOAD)) {
            // since payload trait can only be applied to a single member in a output shape
            MemberShape memberShape = HttpBindingIndex.of(model)
                    .getResponseBindings(operation, HttpBinding.Location.PAYLOAD).stream()
                    .findFirst()
                    .orElseThrow(() -> new CodegenException("Expected payload binding member"))
                    .getMember();

            Shape payloadShape = model.expectShape(memberShape.getTarget());

            // if target shape is of type String or type Blob, then delegate deserializers for explicit payload shapes
            if (payloadShape.isStringShape() || payloadShape.isBlobShape()) {
                writeMiddlewarePayloadBindingDeserializerDelegator(writer, context.getService(), targetShape);
                return;
            }
            // for other payload target types we should deserialize using the appropriate document deserializer
            targetShape = payloadShape;
            operand += "." + context.getSymbolProvider().toMemberName(memberShape);
        }

        writeMiddlewareDocumentBindingDeserializerDelegator(context, writer, targetShape, operand);
    }

    @Override
    protected void generateOperationDocumentDeserializer(
            GenerationContext context, OperationShape operation
    ) {
        Model model = context.getModel();
        HttpBindingIndex bindingIndex = HttpBindingIndex.of(model);
        Set<MemberShape> documentBindings = bindingIndex.getResponseBindings(operation, HttpBinding.Location.DOCUMENT)
                .stream()
                .map(HttpBinding::getMember)
                .collect(Collectors.toSet());

        Shape outputShape = ProtocolUtils.expectOutput(model, operation);
        GoWriter writer = context.getWriter().get();

        if (documentBindings.size() != 0) {
            outputShape.accept(new XmlShapeDeserVisitor(context, documentBindings::contains));
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
    protected void generateDocumentBodyShapeDeserializers(GenerationContext context, Set<Shape> shapes) {
        XmlShapeDeserVisitor visitor = new XmlShapeDeserVisitor(context);
        shapes.forEach(shape -> {
            if (generatedDocumentBodyShapeDeserializers.contains(shape.toShapeId())) {
                return;
            }
            shape.accept(visitor);
            generatedDocumentBodyShapeDeserializers.add(shape.toShapeId());
        });
    }

    // Generate deserializers for shapes with payload binding
    private void writePayloadBindingDeserializer(
            GenerationContext context,
            Shape shape,
            Predicate<MemberShape> filterMemberShapes
    ) {
        GoWriter writer = context.getWriter().get();
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol shapeSymbol = symbolProvider.toSymbol(shape);
        String funcName = ProtocolGenerator.getDocumentDeserializerFunctionName(shape, context.getService(), getProtocolName());

        for (MemberShape memberShape : new TreeSet<>(shape.members())) {
            if (!filterMemberShapes.test(memberShape)) {
                continue;
            }

            String memberName = symbolProvider.toMemberName(memberShape);
            Shape targetShape = context.getModel().expectShape(memberShape.getTarget());
            if (!targetShape.isStringShape() && !targetShape.isBlobShape()) {
                shape.accept(new XmlShapeDeserVisitor(context, filterMemberShapes));
                return;
            }
            writer.openBlock("func $L(v $P, body io.ReadCloser) error {", "}", funcName, shapeSymbol, () -> {
                writer.openBlock("if v == nil {", "}", () -> {
                    writer.write("return fmt.Errorf(\"unsupported deserialization of nil %T\", v)");
                });
                writer.insertTrailingNewline();

                if (targetShape.hasTrait(StreamingTrait.class)) {
                    writer.write("v.$L = body", memberName);
                    writer.write("return nil");
                    return;
                }

                writer.addUseImports(SmithyGoDependency.IOUTIL);
                writer.write("bs, err := ioutil.ReadAll(body)");
                writer.write("if err != nil { return err }");
                writer.openBlock("if len(bs) > 0 {", "}", () -> {
                    if (targetShape.isBlobShape()) {
                        writer.write("v.$L = bs", memberName);
                    } else { // string
                        writer.addUseImports(SmithyGoDependency.SMITHY_PTR);
                        if (targetShape.hasTrait(EnumTrait.class)) {
                            writer.write("v.$L = $T(bs)", memberName, symbolProvider.toSymbol(targetShape));
                        } else {
                            writer.write("v.$L = ptr.String(string(bs))", memberName);
                        }
                    }
                });
                writer.write("return nil");
            });
        }
    }

    // Writes middleware that delegates to deserializers for shapes that have explicit payload.
    private void writeMiddlewarePayloadBindingDeserializerDelegator(
            GoWriter writer, ServiceShape service, Shape shape
    ) {
        String deserFuncName = ProtocolGenerator.getDocumentDeserializerFunctionName(shape, service, getProtocolName());
        writer.write("err = $L(output, response.Body)", deserFuncName);
        writer.openBlock("if err != nil {", "}", () -> {
            writer.addUseImports(SmithyGoDependency.SMITHY);
            writer.write(String.format("return out, metadata, &smithy.DeserializationError{Err:%s}",
                    "fmt.Errorf(\"failed to deserialize response payload, %w\", err)"));
        });
    }

    // Writes middleware that delegates to deserializers for shapes that have implicit payload.
    private void writeMiddlewareDocumentBindingDeserializerDelegator(
            GenerationContext context,
            GoWriter writer,
            Shape shape,
            String operand
    ) {
        XmlProtocolUtils.initializeXmlDecoder(writer, "response.Body", "out, metadata,", "nil");

        String functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(
                shape, context.getService(), context.getProtocolName());

        writer.write("err = $L(&$L, decoder)", functionName, operand);
        XmlProtocolUtils.handleDecodeError(writer, "out, metadata,");
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
                var ctxWriter = ctx.getWriter().get();
                var stringValue = SymbolUtils.createValueSymbolBuilder("StringValue",
                        AwsGoDependency.SERVICE_INTERNAL_EVENTSTREAM).build();
                var contentTypeHeader = SymbolUtils.createValueSymbolBuilder("ContentTypeHeader",
                        AwsGoDependency.SERVICE_INTERNAL_EVENTSTREAMAPI).build();

                ctxWriter.write("msg.Headers.Set($T, $T($S))",
                        contentTypeHeader, stringValue, getDocumentContentType());

                String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(payloadTarget,
                        context.getService(), getProtocolName());

                initalizeXmlEncoder(context, ctxWriter, payloadTarget, "root", operand);

                ctxWriter.openBlock("if err := $L(input, xmlEncoder.RootElement(root)); err != nil {", "}",
                                functionName, () -> {
                                    ctxWriter.write("return &$T{Err: err}",
                                            SymbolUtils.createValueSymbolBuilder("SerializationError",
                                                    SmithyGoDependency.SMITHY).build());
                                })
                        .write("msg.Payload = xmlEncoder.Bytes()")
                        .write("return nil");
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
    }

    @Override
    protected void generateEventStreamDeserializers(
            GenerationContext context,
            UnionShape eventUnion,
            Set<EventStreamInfo> eventStreamInfos
    ) {
        var model = context.getModel();

        AwsEventStreamUtils.generateEventStreamDeserializer(context, eventUnion);
        AwsEventStreamUtils.generateEventStreamExceptionDeserializer(context, eventUnion, ctx -> {
            var ctxWriter = ctx.getWriter().get();
            ctxWriter.write("br := $T(msg.Payload)", SymbolUtils.createValueSymbolBuilder("NewReader",
                    SmithyGoDependency.BYTES).build());
            AwsProtocolUtils.initializeJsonDecoder(ctxWriter, "br");
            ctxWriter.addUseImports(AwsGoDependency.AWS_XML);
            ctxWriter.write("""
                            errorComponents, err := $T(br, $L)
                            if err != nil {
                                return err
                            }
                            errorCode := "UnknownError"
                            errorMessage := errorCode
                            if ev := exceptionType.String(); len(ev) > 0 {
                                errorCode = ev
                            } else if ev := errorComponents.Code; len(ev) > 0 {
                                errorCode = ev
                            }
                            if ev := errorComponents.Message; len(ev) > 0 {
                                errorMessage = ev
                            }
                            return &$T{
                                Code: errorCode,
                                Message: errorMessage,
                            }
                            """,
                    SymbolUtils.createValueSymbolBuilder("GetErrorResponseComponents", AwsGoDependency.AWS_XML).build(),
                    isNoErrorWrapping(context),
                    SymbolUtils.createValueSymbolBuilder("GenericAPIError", SmithyGoDependency.SMITHY).build()
            );
        });

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
                            var ctxWriter = ctx.getWriter().get();
                            ctxWriter.write("br := $T(msg.Payload)", SymbolUtils.createValueSymbolBuilder("NewReader",
                                            SmithyGoDependency.BYTES).build())
                                    .write("output := &$T{}", context.getSymbolProvider().toSymbol(payloadTarget));

                            String functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(
                                    payloadTarget, context.getService(), getProtocolName());

                            initializeXmlDecoder(ctxWriter, "br", "output");
                            boolean isNoErrorWrapping = isNoErrorWrapping(context);
                            Runnable writeErrorDelegator = () -> {
                                ctxWriter.write("err = $L(&output, decoder)", functionName);
                                handleDecodeError(ctxWriter, "");
                            };

                            if (isNoErrorWrapping) {
                                writeErrorDelegator.run();
                            } else {
                                ctxWriter.write("t, err = decoder.GetElement(\"Error\")");
                                XmlProtocolUtils.handleDecodeError(ctxWriter, "");
                                Symbol wrapNodeDecoder = SymbolUtils.createValueSymbolBuilder("WrapNodeDecoder",
                                        SmithyGoDependency.SMITHY_XML).build();
                                ctxWriter.write("decoder = $T(decoder.Decoder, t)", wrapNodeDecoder);
                                writeErrorDelegator.run();
                            }
                        });

                eventDocumentShapes.add(targetShape);
            } else {
                AwsEventStreamUtils.generateEventMessageDeserializer(context, targetShape,
                        (ctx, payloadTarget, operand) -> {
                            var ctxWriter = ctx.getWriter().get();

                            ctxWriter.write("br := $T(msg.Payload)", SymbolUtils.createValueSymbolBuilder(
                                    "NewReader", SmithyGoDependency.BYTES).build());

                            XmlProtocolUtils.initializeXmlDecoder(ctxWriter, "br", "", "nil");

                            String functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(
                                    payloadTarget, context.getService(), context.getProtocolName());

                            ctxWriter.write("err = $L(&$L, decoder)", functionName, operand);
                            XmlProtocolUtils.handleDecodeError(ctxWriter, "");

                            ctxWriter.write("return nil");
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
    }
}
