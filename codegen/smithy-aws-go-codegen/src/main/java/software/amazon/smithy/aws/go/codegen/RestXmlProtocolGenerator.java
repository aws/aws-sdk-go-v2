package software.amazon.smithy.aws.go.codegen;

import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;
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
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.model.traits.XmlAttributeTrait;

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
        HttpBindingIndex bindingIndex = model.getKnowledge(HttpBindingIndex.class);

        Set<MemberShape> documentBindings = bindingIndex.getRequestBindings(operation, HttpBinding.Location.DOCUMENT)
                .stream()
                .map(HttpBinding::getMember)
                .collect(Collectors.toSet());

        if (documentBindings.isEmpty()) {
            return;
        }
        Shape inputShape = ProtocolUtils.expectInput(context.getModel(), operation);
        inputShape.accept(new XmlShapeSerVisitor(context, memberShape -> documentBindings.contains(memberShape) && !memberShape.hasTrait(XmlAttributeTrait.class)));
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

        XMLProtocolUtils.generateXMLStartElement(context, inputShape, "root", "input");
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

        XMLProtocolUtils.generateXMLStartElement(context, payloadShape, "payloadRoot", operand);
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

    /*========================Deserializers==================================*/
    @Override
    protected void writeMiddlewareDocumentDeserializerDelegator(
            GenerationContext context,
            OperationShape operation,
            GoStackStepMiddlewareGenerator generator
    ) {

    }

    @Override
    protected void generateOperationDocumentDeserializer(
            GenerationContext context, OperationShape operation
    ) {

    }

    @Override
    protected void generateDocumentBodyShapeDeserializers(
            GenerationContext context, Set<Shape> shapes
    ) {

    }

    @Override
    protected void writeErrorMessageCodeDeserializer(GenerationContext context) {
        context.getWriter().writeDocs("TODO: implement error message / code deser");
        context.getWriter().write("_ = errorBody");
    }

    @Override
    protected void deserializeError(GenerationContext context, StructureShape shape) {
        GoWriter writer = context.getWriter();

        writer.write("return &smithy.DeserializationError{"
                + "Err: fmt.Errorf(\"TODO: Implement error deserializerfor %v\", errorBody)}");
    }
}
