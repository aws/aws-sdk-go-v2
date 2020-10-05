package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

public class Route53ErrorCustomizations implements GoIntegration {
    private static String ADD_ERROR_HANDLER_INTERNAL = "HandleCustomErrorDeserialization";

    @Override
    public byte getOrder() {
        // The associated customization ordering is relative to operation deserializers
        // and thus the integration should be added at the end.
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate(Route53ErrorCustomizations::supportsCustomError)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ADD_ERROR_HANDLER_INTERNAL,
                                        AwsCustomGoDependency.ROUTE53_CUSTOMIZATION).build())
                                .build())
                        .build()
        );
    }

    // returns true if the operation supports custom route53 error response
    private static boolean supportsCustomError(Model model, ServiceShape service, OperationShape operation){
        if (!isRoute53Service(model, service)) {
            return false;
        }

        return operation.getId().getName().equalsIgnoreCase("ChangeResourceRecordSets");
    }

    // returns true if service is route53
    private static boolean isRoute53Service(Model model, ServiceShape service) {
        String serviceId= service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("Route 53");
    }
}
