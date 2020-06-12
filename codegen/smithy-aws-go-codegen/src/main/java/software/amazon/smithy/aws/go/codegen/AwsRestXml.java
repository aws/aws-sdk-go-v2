package software.amazon.smithy.aws.go.codegen;

import software.amazon.smithy.aws.traits.protocols.RestXmlTrait;
import software.amazon.smithy.model.shapes.ShapeId;

/**
 * Handles generating the aws.rest-xml protocol for services.
 *
 * @inheritDoc
 *
 * @see RestXmlProtocolGenerator
 */
public final class AwsRestXml extends RestXmlProtocolGenerator {

    @Override
    protected String getDocumentContentType() {
        return "application/xml";
    }

    @Override
    public ShapeId getProtocol() {
        return RestXmlTrait.ID;
    }

}
