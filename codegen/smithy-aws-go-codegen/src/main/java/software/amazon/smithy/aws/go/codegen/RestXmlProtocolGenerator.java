package software.amazon.smithy.aws.go.codegen;


import static software.amazon.smithy.go.codegen.integration.HttpProtocolGeneratorUtils.isShapeWithResponseBindings;
import static software.amazon.smithy.aws.go.codegen.XmlProtocolUtils.initializeXmlDecoder;
import static software.amazon.smithy.aws.go.codegen.XmlProtocolUtils.handleDecodeError;
import static software.amazon.smithy.aws.go.codegen.XmlProtocolUtils.writeXmlErrorMessageCodeDeserializer;
import static software.amazon.smithy.aws.go.codegen.XmlProtocolUtils.generateXMLStartElement;
import static software.amazon.smithy.aws.go.codegen.XmlProtocolUtils.generatePayloadAsDocumentXMLStartElement;

import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;
import java.util.TreeSet;
import java.util.function.Predicate;
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
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.MediaTypeTrait;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.model.traits.StreamingTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.model.traits.XmlAttributeTrait;
import software.amazon.smithy.model.traits.XmlNamespaceTrait;

abstract class RestXmlProtocolGenerator extends HttpBindingProtocolGenerator {
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
        GoWriter writer = context.getWriter();
        writer.addUseImports(SmithyGoDependency.SMITHY);
        writer.addUseImports(SmithyGoDependency.SMITHY_XML);

        writer.write("restEncoder.SetHeader(\"Content-Type\").String($S)", getDocumentContentType());
        writer.write("");

        Shape inputShape = ProtocolUtils.expectInput(context.getModel(), operation);
        String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(inputShape, getProtocolName());

        writer.addUseImports(SmithyGoDependency.SMITHY_XML);
        writer.addUseImports(SmithyGoDependency.BYTES);
        writer.write("xmlEncoder := smithyxml.NewEncoder(bytes.NewBuffer(nil))");

        generateXMLStartElement(context, inputShape, "root", "input");

        // check if service shape is bound by xmlNameSpace Trait
        Optional<XmlNamespaceTrait> xmlNamespaceTrait = context.getService().getTrait(XmlNamespaceTrait.class);
        if (xmlNamespaceTrait.isPresent()) {
            XmlNamespaceTrait namespace = xmlNamespaceTrait.get();
            writer.write("root.Attr = append(root.Attr, smithyxml.NewNamespaceAttribute($S, $S))",
                    namespace.getPrefix().isPresent() ? namespace.getPrefix().get() : "", namespace.getUri()
            );
        }

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

    @Override
    protected void writeMiddlewarePayloadAsDocumentSerializerDelegator(
            GenerationContext context,
            MemberShape memberShape,
            String operand
    ) {
        GoWriter writer = context.getWriter();
        Model model = context.getModel();
        Shape payloadShape = model.expectShape(memberShape.getTarget());

        String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(payloadShape,
                getProtocolName());
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
        shapes.forEach(shape -> shape.accept(visitor));
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
        GoWriter writer = context.getWriter();
        Symbol symbol = context.getSymbolProvider().toSymbol(shape);

        writer.write("output := &$T{}", symbol);
        writer.insertTrailingNewline();

        if (isShapeWithRestResponseBindings(context.getModel(), shape)) {
            String bindingDeserFunctionName = ProtocolGenerator.getOperationHttpBindingsDeserFunctionName(
                    shape, getProtocolName());
            writer.openBlock("if err := $L(output, response); err != nil {", "}", bindingDeserFunctionName, () -> {
                writer.addUseImports(SmithyGoDependency.SMITHY);
                writer.write(String.format("return &smithy.DeserializationError{Err: %s}",
                        "fmt.Errorf(\"failed to decode response error with invalid HTTP bindings, %w\", err)"));
            });
            writer.insertTrailingNewline();
        }

        if (isShapeWithResponseBindings(context.getModel(), shape, HttpBinding.Location.DOCUMENT)) {
            String documentDeserFunctionName = ProtocolGenerator.getDocumentDeserializerFunctionName(
                    shape, getProtocolName());
            writer.addUseImports(SmithyGoDependency.IO);
            initializeXmlDecoder(writer, "errorBody", "output");
            writer.write("err = $L(&output, decoder)", documentDeserFunctionName);
            handleDecodeError(writer, "");
            writer.insertTrailingNewline();
        }

        writer.write("return output");
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
        GoWriter writer = context.getWriter();
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
                writeMiddlewarePayloadBindingDeserializerDelegator(writer, targetShape);
                return;
            }
            // for other payload target types we should deserialize using the appropriate document deserializer
            targetShape = payloadShape;
            operand += "." + context.getSymbolProvider().toMemberName(memberShape);
        }

        writeMiddlewareDocumentBindingDeserializerDelegator(writer, targetShape, operand);
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
        GoWriter writer = context.getWriter();

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
        for (Shape shape : shapes) {
            shape.accept(visitor);
        }
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
                            writer.write("v.$L = string(bs)", memberName);
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
    private void writeMiddlewarePayloadBindingDeserializerDelegator(GoWriter writer, Shape shape) {
        String deserFuncName = ProtocolGenerator.getDocumentDeserializerFunctionName(shape, getProtocolName());
        writer.write("err = $L(output, response.Body)", deserFuncName);
        writer.openBlock("if err != nil {", "}", () -> {
            writer.addUseImports(SmithyGoDependency.SMITHY);
            writer.write(String.format("return out, metadata, &smithy.DeserializationError{Err:%s}",
                    "fmt.Errorf(\"failed to deserialize response payload, %w\", err)"));
        });
    }

    // Write middleware that delegates to deserializers for shapes that have implicit payload
    private void writeMiddlewareDocumentBindingDeserializerDelegator(
            GoWriter writer,
            Shape shape,
            String operand
    ) {
        XmlProtocolUtils.initializeXmlDecoder(writer, "response.Body", "out, metadata,", "nil");
        String deserFuncName = ProtocolGenerator.getDocumentDeserializerFunctionName(shape, getProtocolName());
        writer.addUseImports(SmithyGoDependency.IO);
        writer.write("err = $L(&$L, decoder)", deserFuncName, operand);
        XmlProtocolUtils.handleDecodeError(writer, "out, metadata,");
    }
}
