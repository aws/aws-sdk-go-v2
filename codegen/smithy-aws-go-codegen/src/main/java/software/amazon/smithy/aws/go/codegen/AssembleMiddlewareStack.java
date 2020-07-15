package software.amazon.smithy.aws.go.codegen;

import java.util.List;
import java.util.Map;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.aws.traits.auth.UnsignedPayloadTrait;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.ServiceIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.AuthTrait;
import software.amazon.smithy.model.traits.Trait;
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

                // Add ContentLengthMiddleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddContentLengthMiddleware", SmithyGoDependency.SMITHY_HTTP_TRANSPORT)
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

                // Add unsigned payload middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            return hasSigV4AuthScheme(model, service, operation)
                                    && operation.hasTrait(UnsignedPayloadTrait.class);
                        })
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddUnsignedPayloadMiddleware", AwsGoDependency.AWS_SIGNER_V4)
                                        .build())
                                .build())
                        .build(),

                // Add signed payload middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            return hasSigV4AuthScheme(model, service, operation)
                                    && !operation.hasTrait(UnsignedPayloadTrait.class);
                        })
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddComputePayloadSHA256Middleware", AwsGoDependency.AWS_SIGNER_V4)
                                        .build())
                                .build())
                        .build(),

                // Add content-sha256 payload header middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            return hasSigV4AuthScheme(model, service, operation)
                                    && operation.hasTrait(UnsignedPayloadTrait.class);
                        })
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddContentSHA256HeaderMiddleware", AwsGoDependency.AWS_SIGNER_V4)
                                        .build())
                                .build())
                        .build(),

                // Add newAttempt middleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddRetryMiddlewares", AwsGoDependency.AWS_RETRY)
                                        .build())
                                .useClientOptions()
                                .build())
                        .build(),

                // Add HTTPSigner middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate(this::hasSigV4AuthScheme)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddHTTPSignerMiddleware", AwsGoDependency.AWS_SIGNER_V4)
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
                        .build()
        );
    }

    /**
     * Returns if the SigV4Trait is a auth scheme for the service and operation.
     *
     * @param model     model definition
     * @param service   service shape for the API
     * @param operation operation shape
     * @return if SigV4Trait is an auth scheme for the operation and service.
     */
    private boolean hasSigV4AuthScheme(Model model, ServiceShape service, OperationShape operation) {
        ServiceIndex serviceIndex = model.getKnowledge(ServiceIndex.class);
        Map<ShapeId, Trait> auth = serviceIndex.getEffectiveAuthSchemes(
                service.getId(),
                operation.getId()
        );

        return auth.containsKey(SigV4Trait.ID);
    }
}
