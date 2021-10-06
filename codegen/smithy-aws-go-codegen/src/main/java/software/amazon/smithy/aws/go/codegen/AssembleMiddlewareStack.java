package software.amazon.smithy.aws.go.codegen;

import java.util.List;
import software.amazon.smithy.aws.go.codegen.customization.AdjustAwsRestJsonContentType;
import software.amazon.smithy.aws.traits.auth.UnsignedPayloadTrait;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.knowledge.EventStreamIndex;
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
                        .operationPredicate((model, service, operation) ->
                                EventStreamIndex.of(model).getInputInfo(operation).isEmpty())
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                                "AddComputeContentLengthMiddleware",
                                                SmithyGoDependency.SMITHY_HTTP_TRANSPORT)
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

                // Add streaming events payload middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            if (!AwsSignatureVersion4.hasSigV4AuthScheme(
                                    model, service, operation)) {
                                return false;
                            }
                            return EventStreamIndex.of(model).getInputInfo(operation).isPresent();
                        })
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                                "AddStreamingEventsPayload", AwsGoDependency.AWS_SIGNER_V4)
                                        .build())
                                .build())
                        .build(),

                // Add unsigned payload middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            if (!AwsSignatureVersion4.hasSigV4AuthScheme(
                                    model, service, operation)) {
                                return false;
                            }
                            var noEventStream = EventStreamIndex.of(model).getInputInfo(operation).isEmpty();
                            return operation.hasTrait(UnsignedPayloadTrait.class) && noEventStream;
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
                            if (!AwsSignatureVersion4.hasSigV4AuthScheme(
                                    model, service, operation)) {
                                return false;
                            }
                            var noEventStream = EventStreamIndex.of(model).getInputInfo(operation).isEmpty();
                            return !operation.hasTrait(UnsignedPayloadTrait.class) && noEventStream;
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
                            if (!AwsSignatureVersion4.hasSigV4AuthScheme(
                                    model, service, operation)) {
                                return false;
                            }
                            var hasEventStream = EventStreamIndex.of(model).getInputInfo(operation).isPresent();
                            return operation.hasTrait(UnsignedPayloadTrait.class) || hasEventStream;
                        })
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
                // Add middleware to store raw response omn metadata
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                                "AddRawResponseToMetadata", AwsGoDependency.AWS_MIDDLEWARE)
                                        .build())
                                .build())
                        .build(),

                // Add recordResponseTiming middleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                                "AddRecordResponseTiming", AwsGoDependency.AWS_MIDDLEWARE)
                                        .build())
                                .build())
                        .build(),

                // Add Client UserAgent
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createPointableSymbolBuilder(
                                        AwsClientUserAgent.MIDDLEWARE_RESOLVER).build())
                                .build())
                        .build(),

                // Add REST-JSON Content-Type Adjuster
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        AdjustAwsRestJsonContentType.RESOLVER_NAME).build())
                                .build())
                        .servicePredicate((model, serviceShape) ->
                                AdjustAwsRestJsonContentType.isServiceOnShameList(serviceShape))
                        .build(),

                // Add Event Stream Input Writer (must be added AFTER retryer)
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) ->
                                EventStreamIndex.of(model).getInputInfo(operation).isPresent())
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                                "AddInitializeStreamWriter",
                                                AwsGoDependency.SERVICE_INTERNAL_EVENTSTREAMAPI)
                                        .build())
                                .build())
                        .build()
        );
    }
}
