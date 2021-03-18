package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.aws.traits.auth.UnsignedPayloadTrait;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.utils.ListUtils;


public class S3ContentSHA256Header implements GoIntegration {

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first.
     *
     * @return Returns the sort order, defaults to -40.
     */
    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        // Only add the middleware if UnsignedPayloadTrait is not specified, as this middleware
                        // will have already been added.
                        .operationPredicate((model, service, operation) -> {
                            if (!(S3ModelUtils.isServiceS3(model, service)
                                    || S3ModelUtils.isServiceS3Control(model, service))) {
                                return false;
                            }
                            return !operation.hasTrait(UnsignedPayloadTrait.class);
                        })
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddContentSHA256HeaderMiddleware",
                                        AwsGoDependency.AWS_SIGNER_V4
                                ).build())
                                .build())
                        .build()
        );
    }
}
