package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.traits.HttpTrait;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

/**
 * Add middleware, which adds {Expect: 100-continue} header for s3 client HTTP PUT request larger than 2MB
 * or with unknown size streaming bodies, during operation builder step
 */
public class S3100Continue implements GoIntegration {
    private static final String ADD_100Continue_Header = "add100Continue";
    private static final String ADD_100Continue_Header_INTERNAL = "Add100Continue";
    private static final String ADD_100Continue_Header_Option = "AddContinueOption";
    private static final String Continue_Client_Option = "ContinueHeaderThresholdBytes";

    /**
     * Return true if service is Amazon S3.
     *
     * @param model   is the generation model.
     * @param service is the service shape being audited.
     */
    private static boolean isS3Service(Model model, ServiceShape service) {
        return S3ModelUtils.isServiceS3(model, service);
    }

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
        if (!isS3Service(model, service)) {
            return;
        }

        goDelegator.useShapeWriter(service, this::writeMiddlewareHelper);
    }

    private void writeMiddlewareHelper(GoWriter writer) {
        writer.openBlock("func $L(stack *middleware.Stack, options Options) error {", "}", ADD_100Continue_Header, () -> {
            writer.openBlock("return $T(stack, $T{", "})",
                    SymbolUtils.createValueSymbolBuilder(ADD_100Continue_Header_INTERNAL,
                            AwsCustomGoDependency.S3_SHARED_CUSTOMIZATION).build(),
                    SymbolUtils.createValueSymbolBuilder(ADD_100Continue_Header_Option,
                            AwsCustomGoDependency.S3_SHARED_CUSTOMIZATION).build(), () -> {
                        writer.write("ContinueHeaderThresholdBytes: options.ContinueHeaderThresholdBytes,"
                        );
                    }
                    );
        });
        writer.insertTrailingNewline();
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> isS3Service(model, service) && operation.
                                getTrait(HttpTrait.class).get().getMethod().equals("PUT"))
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ADD_100Continue_Header).build())
                                .useClientOptions()
                                .build()
                        )
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3100Continue::isS3Service)
                        .configFields(ListUtils.of(
                            ConfigField.builder()
                                    .name(Continue_Client_Option)
                                    .type(SymbolUtils.createValueSymbolBuilder("int64")
                                            .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                            .build())
                                    .documentation("The threshold ContentLength for HTTP PUT request to receive {Expect: 100-continue} header. " +
                                            "When set to -1, this header will be opt out of the operation request; when set to 0, the threshold" +
                                            "will be set to default 2MB")
                                    .build()
                        ))
                        .build()
        );
    }
}
