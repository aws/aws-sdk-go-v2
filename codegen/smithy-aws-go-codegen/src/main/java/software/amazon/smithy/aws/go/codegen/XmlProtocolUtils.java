package software.amazon.smithy.aws.go.codegen;

import java.util.Collection;
import java.util.Optional;
import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.aws.traits.protocols.RestXmlTrait;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoValueAccessUtils;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.SyntheticClone;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.knowledge.GoPointableIndex;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.model.traits.XmlAttributeTrait;
import software.amazon.smithy.model.traits.XmlNameTrait;
import software.amazon.smithy.model.traits.XmlNamespaceTrait;

public final class XmlProtocolUtils {
    private XmlProtocolUtils() {

    }

    /**
     * generateXMLStartElement generates the XML start element for a shape. It is used to generate smithy xml's startElement.
     *
     * @param context  is the generation context.
     * @param shape    is the Shape for which xml start element is to be generated.
     * @param dst      is the operand name which holds the generated start element.
     * @param inputSrc is the input variable for the shape with values to be serialized.
     */
    public static void generateXMLStartElement(
            ProtocolGenerator.GenerationContext context, Shape shape, String dst, String inputSrc
    ) {
        GoWriter writer = context.getWriter().get();
        String attrName = dst + "Attr";
        generateXmlNamespaceAndAttributes(context, shape, attrName, inputSrc);

        writer.openBlock("$L := smithyxml.StartElement{ ", "}", dst, () -> {
            writer.openBlock("Name:smithyxml.Name{", "},", () -> {
                writer.write("Local: $S,", getSerializedXMLShapeName(context, shape));
            });
            writer.write("Attr : $L,", attrName);
        });
    }

    /**
     * Generates XML Start element for a document shape marked as a payload.
     *
     * @param context     is the generation context.
     * @param memberShape is the payload as document member shape
     * @param dst         is the operand name which holds the generated start element.
     * @param inputSrc    is the input variable for the shape with values to be serialized.
     */
    public static void generatePayloadAsDocumentXMLStartElement(
            ProtocolGenerator.GenerationContext context, MemberShape memberShape, String dst, String inputSrc
    ) {
        GoWriter writer = context.getWriter().get();
        String attrName = dst + "Attr";
        Shape targetShape = context.getModel().expectShape(memberShape.getTarget());

        generateXmlNamespaceAndAttributes(context, targetShape, attrName, inputSrc);

        writer.openBlock("$L := smithyxml.StartElement{ ", "}", dst, () -> {
            writer.openBlock("Name:smithyxml.Name{", "},", () -> {
                String name = memberShape.getMemberName();
                if (targetShape.isStructureShape()) {
                    if (memberShape.hasTrait(XmlNameTrait.class)) {
                        name = getSerializedXMLMemberName(memberShape);
                    } else {
                        name = getSerializedXMLShapeName(context, targetShape);
                    }
                }

                writer.write("Local: $S,", name);

            });
            writer.write("Attr : $L,", attrName);
        });
    }


    /**
     * Generates XML Attributes as per xmlNamespace and xmlAttribute traits.
     *
     * @param context  is the generation context.
     * @param shape    is the shape that is decorated with XmlNamespace, XmlAttribute trait.
     * @param dst      is the operand name which holds the generated xml Attribute value.
     * @param inputSrc is the input variable for the shape with values to be put as xml attributes.
     */
    private static void generateXmlNamespaceAndAttributes(
            ProtocolGenerator.GenerationContext context, Shape shape, String dst, String inputSrc
    ) {
        GoWriter writer = context.getWriter().get();
        writer.write("$L := []smithyxml.Attr{}", dst);

        Optional<XmlNamespaceTrait> xmlNamespaceTrait = shape.getTrait(XmlNamespaceTrait.class);
        if (xmlNamespaceTrait.isPresent()) {
            XmlNamespaceTrait namespace = xmlNamespaceTrait.get();
            writer.write("$L = append($L, smithyxml.NewNamespaceAttribute($S, $S))",
                    dst, dst,
                    namespace.getPrefix().isPresent() ? namespace.getPrefix().get() : "", namespace.getUri()
            );
        }

        // Traverse member shapes to get attributes
        if (shape.isMemberShape()) {
            MemberShape memberShape = shape.asMemberShape().get();
            Shape target = context.getModel().expectShape(memberShape.getTarget());
            String memberName = context.getSymbolProvider().toMemberName(memberShape);
            String operand = inputSrc + "." + memberName;
            generateXmlAttributes(context, target.members(), operand, dst);
        } else {
            generateXmlAttributes(context, shape.members(), inputSrc, dst);
        }
    }

    private static void generateXmlAttributes(
            ProtocolGenerator.GenerationContext context,
            Collection<MemberShape> members,
            String inputSrc,
            String dst
    ) {
        GoWriter writer = context.getWriter().get();
        members.forEach(memberShape -> {
            if (memberShape.hasTrait(XmlAttributeTrait.class)) {
                GoValueAccessUtils.writeIfNonZeroValueMember(context.getModel(), context.getSymbolProvider(),
                        writer, memberShape, inputSrc, true, memberShape.isRequired(), (operand) -> {
                            // xml attributes should always be string
                            String dest = "av";
                            formatXmlAttributeValueAsString(context, memberShape, operand, dest);
                            writer.write("$L = append($L, smithyxml.NewAttribute($S, $L))",
                                    dst, dst, getSerializedXMLMemberName(memberShape), dest);
                        });
            }
        });
    }

    // generates code to format xml attributes. If a shape type is timestamp, number, or boolean
    // it will be formatted into a string.
    private static void formatXmlAttributeValueAsString(
            ProtocolGenerator.GenerationContext context,
            MemberShape member, String src, String dest
    ) {
        GoWriter writer = context.getWriter().get();
        Shape target = context.getModel().expectShape(member.getTarget());

        // declare destination variable
        writer.write("var $L string", dest);

        // Pointable value references need to be dereferenced before being used.
        String derefSource = src;
        if (GoPointableIndex.of(context.getModel()).isPointable(member)) {
            derefSource = "*" + src;
        }

        if (target.hasTrait(EnumTrait.class)) {
            writer.write("$L = string($L)", dest, derefSource);
            return;
        } else if (target.isStringShape()) {
            // create dereferenced copy of pointed to value.
            writer.write("$L = $L", dest, derefSource);
            return;
        }

        if (target.isTimestampShape() || target.hasTrait(TimestampFormatTrait.class)) {
            TimestampFormatTrait.Format format = member.getMemberTrait(context.getModel(), TimestampFormatTrait.class)
                    .map(TimestampFormatTrait::getFormat).orElse(TimestampFormatTrait.Format.DATE_TIME);
            writer.addUseImports(SmithyGoDependency.SMITHY_TIME);
            switch (format) {
                case DATE_TIME:
                    writer.write("$L = smithytime.FormatDateTime($L)", dest, derefSource);
                    break;
                case HTTP_DATE:
                    writer.write("$L = smithytime.FormatHTTPDate($L)", dest, derefSource);
                    break;
                case EPOCH_SECONDS:
                    writer.addUseImports(SmithyGoDependency.STRCONV);
                    writer.write("$L = strconv.FormatFloat(smithytime.FormatEpochSeconds($L), 'f', -1, 64)",
                            dest, derefSource);
                    break;
                case UNKNOWN:
                    throw new CodegenException("Unknown timestamp format");
            }
            return;
        }

        if (target.isBooleanShape()) {
            writer.write(SmithyGoDependency.STRCONV);
            writer.write("$L = strconv.FormatBool($L)", dest, derefSource);
            return;
        }

        if (target.isByteShape() || target.isShortShape() || target.isIntegerShape() || target.isLongShape()) {
            writer.write(SmithyGoDependency.STRCONV);
            writer.write("$L = strconv.FormatInt(int64($L), 10)", dest, derefSource);
            return;
        }

        if (target.isFloatShape()) {
            writer.write(SmithyGoDependency.STRCONV);
            writer.write("$L = strconv.FormatFloat(float64($L),'f', -1, 32)", dest, derefSource);
            return;
        }

        if (target.isDoubleShape()) {
            writer.write(SmithyGoDependency.STRCONV);
            writer.write("$L = strconv.FormatFloat($L,'f', -1, 64)", dest, derefSource);
            return;
        }

        if (target.isBigIntegerShape() || target.isBigDecimalShape()) {
            throw new CodegenException(String.format("Cannot serialize shape type %s on protocol, shape: %s.",
                    target.getType(), target.getId()));
        }

        throw new CodegenException(
                "Members serialized as XML attributes can only be of string, number, boolean or timestamp format");
    }

    /**
     * getSerializedXMLMemberName returns a xml member name used for serializing. If a member shape has
     * XML name trait, xml name would be given precedence over member name.
     *
     * @param memberShape is the member shape for which serializer name is queried.
     * @return name of a xml member shape used by serializers
     */
    private static String getSerializedXMLMemberName(MemberShape memberShape) {
        Optional<XmlNameTrait> xmlNameTrait = memberShape.getTrait(XmlNameTrait.class);
        return xmlNameTrait.isPresent() ? xmlNameTrait.get().getValue() : memberShape.getMemberName();
    }

    /**
     * getSerializedXMLShapeName returns a xml shape name used for serializing. If a member shape
     * has xml name trait, xml name would be given precedence over member name.
     * This correctly handles renamed shapes, and returns the original shape name.
     *
     * @param context is the generation context for which
     * @param shape   is the Shape for which serializer name is queried.
     * @return name of a xml member shape used by serializers.
     */
    private static String getSerializedXMLShapeName(ProtocolGenerator.GenerationContext context, Shape shape) {
        SymbolProvider symbolProvider = context.getSymbolProvider();
        Symbol shapeSymbol = symbolProvider.toSymbol(shape);
        String shapeName = shapeSymbol.getName();
        ServiceShape service = context.getService();

        // check if synthetic cloned shape
        Optional<SyntheticClone> clone = shape.getTrait(SyntheticClone.class);
        if (clone.isPresent()) {
            SyntheticClone cl = clone.get();
            if (cl.getArchetype().isPresent()) {
                shapeName = cl.getArchetype().get().getName(service);
            }
        }

        // check if shape is member shape
        Optional<MemberShape> member = shape.asMemberShape();
        if (member.isPresent()) {
            return getSerializedXMLMemberName(member.get());
        }

        return shape.getTrait(XmlNameTrait.class).map(XmlNameTrait::getValue).orElse(shapeName);
    }

    /**
     * initializeXmlDecoder generates stub code to initialize xml decoder.
     * Returns nil in case EOF occurs while initializing xml decoder.
     *
     * @param writer       the go writer used to write
     * @param bodyLocation the variable used to represent response body
     */
    public static void initializeXmlDecoder(GoWriter writer, String bodyLocation) {
        initializeXmlDecoder(writer, bodyLocation, "", "nil");
    }

    /**
     * initializeXmlDecoder generates stub code to initialize xml decoder
     *
     * @param writer       the go writer used to write
     * @param bodyLocation the variable used to represent response body
     * @param returnOnEOF  the variable to return in case an EOF error occurs while initializing xml decoder
     */
    public static void initializeXmlDecoder(GoWriter writer, String bodyLocation, String returnOnEOF) {
        initializeXmlDecoder(writer, bodyLocation, "", returnOnEOF);
    }

    /**
     * initializeXmlDecoder generates stub code to initialize xml decoder
     *
     * @param writer       the go writer used to write
     * @param bodyLocation the variable used to represent response body
     * @param returnExtras the extra variables to be returned with the wrapped error check statement
     * @param returnOnEOF  the variable to return in case an EOF error occurs while initializing xml decoder
     */
    public static void initializeXmlDecoder(
            GoWriter writer, String bodyLocation, String returnExtras, String returnOnEOF
    ) {
        // Use a ring buffer and tee reader to help in pinpointing any deserialization errors.
        writer.addUseImports(SmithyGoDependency.SMITHY_IO);
        writer.write("var buff [1024]byte");
        writer.write("ringBuffer := smithyio.NewRingBuffer(buff[:])");
        writer.insertTrailingNewline();

        writer.addUseImports(SmithyGoDependency.IO);
        writer.addUseImports(SmithyGoDependency.XML);
        writer.addUseImports(SmithyGoDependency.SMITHY_XML);
        writer.write("body := io.TeeReader($L, ringBuffer)", bodyLocation);
        writer.write("rootDecoder := xml.NewDecoder(body)");
        writer.write("t, err := smithyxml.FetchRootElement(rootDecoder)");
        writer.write("if err == io.EOF { return $L$L}", returnExtras, returnOnEOF);
        handleDecodeError(writer, returnExtras);

        writer.insertTrailingNewline();
        writer.write("decoder := smithyxml.WrapNodeDecoder(rootDecoder, t)");
        writer.insertTrailingNewline();
    }

    /**
     * handleDecodeError handles the xml deserialization error wrapping.
     *
     * @param writer       the go writer used to write
     * @param returnExtras extra variables to be returned with the wrapped error statement
     */
    public static void handleDecodeError(GoWriter writer, String returnExtras) {
        writer.addUseImports(SmithyGoDependency.IO);
        writer.openBlock("if err != nil {", "}", () -> {
            writer.addUseImports(SmithyGoDependency.BYTES);
            writer.addUseImports(SmithyGoDependency.SMITHY);
            writer.write("var snapshot bytes.Buffer");
            writer.write("io.Copy(&snapshot, ringBuffer)");
            writer.openBlock("return $L&smithy.DeserializationError {", "}", returnExtras, () -> {
                writer.write("Err : fmt.Errorf(\"failed to decode response body, %w\", err),");
                writer.write("Snapshot: snapshot.Bytes(),");
            });
        }).write("");
    }

    /**
     * Generates code to retrieve error code or error message from the error response body
     * This method is used indirectly by generateErrorDispatcher to generate operation specific error handling functions
     *
     * @param context the generation context
     * @see <a href="https://awslabs.github.io/smithy/1.0/spec/aws/aws-restxml-protocol.html#operation-error-serialization">Rest-XML operation error serialization.</a>
     */
    public static void writeXmlErrorMessageCodeDeserializer(ProtocolGenerator.GenerationContext context) {
        GoWriter writer = context.getWriter().get();

        // Check if service uses isNoErrorWrapping setting
        boolean isNoErrorWrapping = context.getService().getTrait(RestXmlTrait.class).map(
                RestXmlTrait::isNoErrorWrapping).orElse(false);

        ServiceShape service = context.getService();

        if (requiresS3Customization(service)) {
            Symbol getErrorComponentFunction = SymbolUtils.createValueSymbolBuilder(
                    "GetErrorResponseComponents",
                    AwsCustomGoDependency.S3_SHARED_CUSTOMIZATION
            ).build();

            Symbol errorOptions = SymbolUtils.createValueSymbolBuilder(
                    "ErrorResponseDeserializerOptions",
                    AwsCustomGoDependency.S3_SHARED_CUSTOMIZATION
            ).build();

            if (isS3Service(service)) {
                // s3 service
                writer.openBlock("errorComponents, err := $T(errorBody, $T{",
                        "})", getErrorComponentFunction, errorOptions, () -> {
                            writer.write("UseStatusCode : true, StatusCode : response.StatusCode,");
                        });
            } else {
                // s3 control
                writer.openBlock("errorComponents, err := $T(errorBody, $T{",
                        "})", getErrorComponentFunction, errorOptions, () -> {
                            writer.write("IsWrappedWithErrorTag: true,");
                        });
            }

            writer.write("if err != nil { return err }");

            writer.insertTrailingNewline();
            writer.openBlock("if hostID := errorComponents.HostID; len(hostID)!=0 {", "}", () -> {
                writer.write("s3shared.SetHostIDMetadata(metadata, hostID)");
            });
        } else {
            writer.addUseImports(AwsGoDependency.AWS_XML);
            writer.write("errorComponents, err := awsxml.GetErrorResponseComponents(errorBody, $L)", isNoErrorWrapping);
            writer.write("if err != nil { return err }");
            writer.insertTrailingNewline();
        }

        writer.addUseImports(AwsGoDependency.AWS_MIDDLEWARE);
        writer.openBlock("if reqID := errorComponents.RequestID; len(reqID)!=0 {", "}", () -> {
            writer.write("awsmiddleware.SetRequestIDMetadata(metadata, reqID)");
        });
        writer.insertTrailingNewline();

        writer.write("if len(errorComponents.Code) != 0 { errorCode = errorComponents.Code}");
        writer.write("if len(errorComponents.Message) != 0 { errorMessage = errorComponents.Message}");
        writer.insertTrailingNewline();

        writer.write("errorBody.Seek(0, io.SeekStart)");
        writer.insertTrailingNewline();
    }

    // returns true if service is either s3 or s3 control and needs s3 customization
    private static boolean requiresS3Customization(ServiceShape service) {
        String serviceId = service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("S3") || serviceId.equalsIgnoreCase("S3 Control");
    }

    private static boolean isS3Service(ServiceShape service) {
        String serviceId = service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("S3");
    }
}
