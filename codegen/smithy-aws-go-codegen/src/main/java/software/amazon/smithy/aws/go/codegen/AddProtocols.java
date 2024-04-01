package software.amazon.smithy.aws.go.codegen;

import java.util.List;
import software.amazon.smithy.aws.go.codegen.protocol.AwsRpc2CborProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUtils;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.utils.ListUtils;

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
    public List<RuntimeClientPlugin> getClientPlugins() {
        List<RuntimeClientPlugin> plugins = HttpProtocolUtils.getCloseResponseClientPlugins((model, service) -> {
            // All AWS protocols are HTTP based currently. When protocol is added that is not it must be
            // excluded if the service is configured for that protocol.
            return true;
        });

        return plugins;
    }

    @Override
    public List<ProtocolGenerator> getProtocolGenerators() {
        return ListUtils.of(
                new AwsRestJson1(),
                new AwsJsonRpc1_0(),
                new AwsJsonRpc1_1(),
                new AwsRestXml(),
                new AwsQuery(),
                new Ec2Query(),
                new AwsRpc2CborProtocolGenerator()
        );
    }
}
