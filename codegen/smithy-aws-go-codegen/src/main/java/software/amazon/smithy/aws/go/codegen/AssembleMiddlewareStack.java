package software.amazon.smithy.aws.go.codegen;

import java.util.List;
import software.amazon.smithy.aws.traits.auth.UnsignedPayloadTrait;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
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
                                        "AddClientRequestIDMiddleware", AwsGoDependency.AWS_MIDDLEWARE)
                                        .build())
                                .build()
                        )
                        .build(),

                // Add ContentLengthMiddleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddComputeContentLengthMiddleware", SmithyGoDependency.SMITHY_HTTP_TRANSPORT)
                                        .build())
                                .build()
                        )
                        .build(),

                // Add endpoint serialize middleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        EndpointGenerator.ADD_MIDDLEWARE_HELPER_NAME).build())
                                .useClientOptions()
                                .build())
                        .build(),

                // Add unsigned payload middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> AwsSignatureVersion4.hasSigV4AuthScheme(
                                model, service, operation) && operation.hasTrait(UnsignedPayloadTrait.class))
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddUnsignedPayloadMiddleware", AwsGoDependency.AWS_SIGNER_V4)
                                        .build())
                                .build())
                        .build(),

                // Add signed payload middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> AwsSignatureVersion4.hasSigV4AuthScheme(
                                model, service, operation) && !operation.hasTrait(UnsignedPayloadTrait.class))
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddComputePayloadSHA256Middleware", AwsGoDependency.AWS_SIGNER_V4)
                                        .build())
                                .build())
                        .build(),

                // Add content-sha256 payload header middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> AwsSignatureVersion4.hasSigV4AuthScheme(
                                model, service, operation) && operation.hasTrait(UnsignedPayloadTrait.class))
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddContentSHA256HeaderMiddleware", AwsGoDependency.AWS_SIGNER_V4)
                                        .build())
                                .build())
                        .build(),

                // Add retryer middleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        AwsRetryMiddlewareHelper.ADD_RETRY_MIDDLEWARES_HELPER)
                                        .build())
                                .useClientOptions()
                                .build())
                        .build(),

                // Add HTTPSigner middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate(AwsSignatureVersion4::hasSigV4AuthScheme)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        AwsSignatureVersion4.REGISTER_MIDDLEWARE_FUNCTION).build())
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

                // Add Client UserAgent
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createPointableSymbolBuilder(
                                        AwsClientUserAgent.MIDDLEWARE_RESOLVER).build())
                                .build())
                        .build()
        );
    }
}
