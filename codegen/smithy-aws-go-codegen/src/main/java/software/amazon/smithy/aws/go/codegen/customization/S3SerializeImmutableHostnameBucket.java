package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

/**
 * Adds the input bucket name back to the path when an immutable endpoint is returned through v1 EndpointResolution.
 */
public class S3SerializeImmutableHostnameBucket implements GoIntegration {
    @Override
    public byte getOrder() { return 127; } // depends on ResolveEndpointV2 middleware for stack insert
    
    private final MiddlewareRegistrar serializeImmutableHostnameBucketMiddleware =
            MiddlewareRegistrar.builder()
                    .resolvedFunction(SymbolUtils.createValueSymbolBuilder("addSerializeImmutableHostnameBucketMiddleware").build())
                    .build();

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3ModelUtils::isServiceS3)
                        .registerMiddleware(serializeImmutableHostnameBucketMiddleware)
                        .build()
        );
    }
}
