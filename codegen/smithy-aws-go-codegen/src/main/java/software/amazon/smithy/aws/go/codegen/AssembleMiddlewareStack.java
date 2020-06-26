package software.amazon.smithy.aws.go.codegen;

import java.util.List;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.aws.traits.auth.UnsignedPayloadTrait;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.traits.AuthTrait;
import software.amazon.smithy.utils.ListUtils;

public class AssembleMiddlewareStack implements GoIntegration {

    /**
 	 * Gets the sort order of the customization from -128 to 127, with lowest
	 * executed first.
     *
     * @return Returns the sort order, defaults to -40.
     */
    @Override
    public byte getOrder() {
        return -40;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                // Add RequestInvocationIDMiddleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddRequestInvocationIDMiddleware", AwsGoDependency.AWS_MIDDLEWARE)
                                        .build())
                                .build()
                        )
                        .build(),

                // Add endpoint serialize middleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddResolveServiceEndpointMiddleware", AwsGoDependency.AWS_MIDDLEWARE)
                                .build())
                                .useClientOptions()
                                .build())
                        .build(),

                // Add attemptClockSkew middleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddAttemptClockSkewMiddleware", AwsGoDependency.AWS_MIDDLEWARE)
                                        .build())
                                .build())
                        .build(),

                // Add newAttempt middleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction( SymbolUtils.createValueSymbolBuilder(
                                        "AddRetryMiddlewares", AwsGoDependency.AWS_RETRY)
                                        .build())
                                .useClientOptions()
                                .build())
                        .build(),

                // Add unsigned payload middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            if (operation.hasTrait(UnsignedPayloadTrait.class)) {
                                return true;
                            }
                            return false;
                        })
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddUnsignedPayloadMiddleware", AwsGoDependency.AWS_SIGNER_V4)
                                        .build())
                                .build())
                        .build(),

                // Add HTTPSigner middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            if (service.hasTrait(SigV4Trait.class) && (!operation.hasTrait(UnsignedPayloadTrait.class))
                                    && (operation.hasTrait(SigV4Trait.class) || !operation.hasTrait(AuthTrait.class))
                            ){
                                return true;
                            }
                            return false;
                        })
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddHTTPSignerMiddlewares", AwsGoDependency.AWS_SIGNER_V4)
                                        .build())
                                .useClientOptions()
                                .build())
                        .build()
        );
    }
}
