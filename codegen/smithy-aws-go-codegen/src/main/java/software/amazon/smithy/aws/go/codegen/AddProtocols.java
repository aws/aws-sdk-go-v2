package software.amazon.smithy.aws.go.codegen;

import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;

import java.util.Collections;
import java.util.List;

public class AddProtocols implements GoIntegration {
    @Override
    public List<ProtocolGenerator> getProtocolGenerators() {
        return Collections.emptyList();
    }
}
