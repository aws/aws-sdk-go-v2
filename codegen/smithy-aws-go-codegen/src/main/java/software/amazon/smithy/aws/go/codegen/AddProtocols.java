package software.amazon.smithy.aws.go.codegen;

import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;

import java.util.List;
import software.amazon.smithy.utils.ListUtils;

public class AddProtocols implements GoIntegration {
    @Override
    public List<ProtocolGenerator> getProtocolGenerators() {
        return ListUtils.of(new AwsRestJson1());
    }
}
