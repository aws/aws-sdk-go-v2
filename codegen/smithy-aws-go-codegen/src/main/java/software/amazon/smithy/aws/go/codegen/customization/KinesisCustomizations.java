package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

public class KinesisCustomizations implements GoIntegration {
    private static final String READ_TIMEOUT_ADDER = "AddResponseReadTimeoutMiddleware";
    private static final String READ_TIMEOUT_DURATION = "ReadTimeoutDuration";

    @Override
    public byte getOrder() {
        // We want this to go last so it can be registered at the very end of the list,
        // meaning it will be the first to be called after the actual http request is
        // made.
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate(KinesisCustomizations::isGetRecords)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(READ_TIMEOUT_ADDER,
                                        AwsGoDependency.AWS_HTTP_TRANSPORT).build())
                                .functionArguments(ListUtils.of(
                                        SymbolUtils.createValueSymbolBuilder(READ_TIMEOUT_DURATION,
                                                AwsCustomGoDependency.KINESIS_CUSTOMIZATION).build()))
                                .build())
                        .build()
        );
    }

    private static boolean isGetRecords(Model model, ServiceShape service, OperationShape operation) {
        return isKinesis(model, service) && operation.getId().getName(service).equals("GetRecords");
    }

    private static boolean isKinesis(Model model, ServiceShape service) {
        return service.expectTrait(ServiceTrait.class).getSdkId().equalsIgnoreCase("Kinesis");
    }
}
