package software.amazon.smithy.aws.go.codegen;

import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

public class AddProtocols implements GoIntegration {
    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first.
     *
     * @return Returns the sort order, defaults to -10.
     */
    @Override
    public byte getOrder() {
        return -10;
    }

    @Override
    public List<ProtocolGenerator> getProtocolGenerators() {
        return ListUtils.of(new AwsRestJson1(), new AwsJsonRpc1_0(), new AwsJsonRpc1_1(), new AwsRestXml());
    }
}
