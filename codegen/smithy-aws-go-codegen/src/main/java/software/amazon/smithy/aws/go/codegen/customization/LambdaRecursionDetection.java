package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

/**
 * Add middleware during operation builder step, which detects Lambda environment and sets its X-Ray trace ID to
 * request header if absent to avoid recursion invocation in Lambda
 */
public class LambdaRecursionDetection implements GoIntegration {
    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first.
     *
     * @return Returns the sort order, defaults to -40.
     */
    @Override
    public byte getOrder() {
        return 126;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                                "AddRecursionDetection", AwsGoDependency.AWS_MIDDLEWARE)
                                        .build())
                                .build()
                        )
                        .build()
                );
    }
}
