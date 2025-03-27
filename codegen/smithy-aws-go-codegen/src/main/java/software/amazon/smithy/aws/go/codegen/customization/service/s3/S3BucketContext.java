package software.amazon.smithy.aws.go.codegen.customization.service.s3;

import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

/**
 * Adds the input bucket name to the context for S3 operations, which is required for a variety of custom S3 behaviors.
 */
public class S3BucketContext implements GoIntegration {
    private final MiddlewareRegistrar putBucketContextMiddleware =
            MiddlewareRegistrar.builder()
                    .resolvedFunction(SymbolUtils.createValueSymbolBuilder("addPutBucketContextMiddleware").build())
                    .build();

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3ModelUtils::isServiceS3)
                        .registerMiddleware(putBucketContextMiddleware)
                        .build()
        );
    }
}
