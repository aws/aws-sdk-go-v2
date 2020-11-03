package software.amazon.smithy.aws.go.codegen;

import java.util.List;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

/**
 * Retrieves request id from response header of deserialized response shapes
 */
public class AWSRequestIDRetriever implements GoIntegration {
    private static final String ADD_REQUEST_ID_RETRIEVER = "addRequestIDRetrieverMiddleware";
    private static final String ADD_REQUEST_ID_RETRIEVER_INTERNAL = "AddRequestIDRetrieverMiddleware";

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
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceShape service = settings.getService(model);

        // S3 doesn't need aws specific wrapper
        if (requiresS3Customization(model, service)) {return;}

        goDelegator.useShapeWriter(service, this::writeMiddlewareHelper);
    }

    private void writeMiddlewareHelper(GoWriter writer) {
        writer.openBlock("func $L(stack *middleware.Stack) error {", "}", ADD_REQUEST_ID_RETRIEVER, () -> {
            writer.write("return $T(stack)",
                    SymbolUtils.createValueSymbolBuilder(ADD_REQUEST_ID_RETRIEVER_INTERNAL,
                            AwsGoDependency.AWS_MIDDLEWARE).build()
            );
        });
        writer.insertTrailingNewline();
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(((model, serviceShape) -> !requiresS3Customization(model,serviceShape)))
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ADD_REQUEST_ID_RETRIEVER).build())
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
