package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import java.util.logging.Logger;
import software.amazon.smithy.aws.go.codegen.AddAwsConfigFields;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;

/**
 * S3AcceptEncodingGzip adds a customization for s3 client to disable
 * auto decoding of GZip content by Golang HTTP Client.
 *
 * This customization provides an option on the S3 client options to enable
 * AcceptEncoding for GZIP. The flag if set, will enable auto decompression of
 * GZIP by the S3 Client.
 *
 * By default, the client's auto decompression of GZIP content is turned off.
 */
public class S3AcceptEncodingGzip implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(AddAwsConfigFields.class.getName());

    private static final String GZIP_DISABLE = "disableAcceptEncodingGzip";
    private static final String GZIP_INTERNAL_ADDER = "AddAcceptEncodingGzip";

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first.
     *
     * @return Returns the sort order, defaults to -50.
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
        if (!S3ModelUtils.isServiceS3(model, settings.getService(model))) {
            return;
        }

        goDelegator.useShapeWriter(settings.getService(model), this::writeMiddlewareHelper);
    }

    private void writeMiddlewareHelper(GoWriter writer) {
        writer.openBlock("func $L(stack *middleware.Stack) error {", "}", GZIP_DISABLE, () -> {
            writer.write("return $T(stack, $T{})",
                    SymbolUtils.createValueSymbolBuilder(GZIP_INTERNAL_ADDER,
                            AwsCustomGoDependency.ACCEPT_ENCODING_CUSTOMIZATION).build(),
                    SymbolUtils.createValueSymbolBuilder(GZIP_INTERNAL_ADDER + "Options",
                            AwsCustomGoDependency.ACCEPT_ENCODING_CUSTOMIZATION).build()
            );
        });
        writer.insertTrailingNewline();
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                // register disableAcceptEncodingGzip middleware
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3ModelUtils::isServiceS3)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(GZIP_DISABLE)
                                        .build())
                                .build()
                        )
                        .build()
        );
    }


}
