package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;


public class S3ContentSHA256Header implements  GoIntegration{

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
                        .servicePredicate(S3ContentSHA256Header::requiresS3Customization)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddContentSHA256HeaderMiddleware",
                                        AwsGoDependency.AWS_SIGNER_V4
                                ).build())
                                .build())
                        .build()
        );
    }

    // returns true if service is either s3 or s3 control and needs s3 customization
    private static boolean requiresS3Customization(Model model, ServiceShape service) {
        String serviceId= service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("S3") || serviceId.equalsIgnoreCase("S3 Control");
    }
}
