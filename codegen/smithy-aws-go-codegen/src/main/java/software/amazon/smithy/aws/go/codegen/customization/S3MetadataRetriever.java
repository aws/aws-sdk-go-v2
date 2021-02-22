package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
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
 * Retrieves S3 specific metadata such as host id, request id from Deserialized response, and
 * sets it on middleware's metadata.
 */
public class S3MetadataRetriever implements GoIntegration {
    private static final String ADD_METADATA_RETRIEVER = "addMetadataRetrieverMiddleware";
    private static final String ADD_METADATA_RETRIEVER_INTERNAL = "AddMetadataRetrieverMiddleware";

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
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceShape service = settings.getService(model);
        if (!requiresS3Customization(model, service)) {
            return;
        }

        goDelegator.useShapeWriter(service, this::writeMiddlewareHelper);
    }

    private void writeMiddlewareHelper(GoWriter writer) {
        writer.openBlock("func $L(stack *middleware.Stack) error {", "}", ADD_METADATA_RETRIEVER, () -> {
            writer.write("return $T(stack)",
                    SymbolUtils.createValueSymbolBuilder(ADD_METADATA_RETRIEVER_INTERNAL,
                            AwsCustomGoDependency.S3_SHARED_CUSTOMIZATION).build()
            );
        });
        writer.insertTrailingNewline();
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3MetadataRetriever::requiresS3Customization)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ADD_METADATA_RETRIEVER).build())
                                .build())
                        .build()
        );
    }

    // returns true if service is either s3 or s3 control and needs s3 customization
    private static boolean requiresS3Customization(Model model, ServiceShape service) {
        return S3ModelUtils.isServiceS3(model, service) || S3ModelUtils.isServiceS3Control(model, service);
    }
}
