package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import software.amazon.smithy.aws.go.codegen.AddAwsConfigFields;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.aws.go.codegen.AwsSignatureVersion4;
import software.amazon.smithy.aws.go.codegen.AwsSignatureVersion4aUtils;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

/**
 * This integration configures the CloudFront Key Value Store client for Signature Version 4a
 */
public class CloudFrontKVSSigV4a implements GoIntegration {
    /**
     * Return true if service is CFKVS.
     *
     * @param model   is the generation model.
     * @param service is the service shape being audited.
     */
    private static boolean isCFKVSService(Model model, ServiceShape service) {
        final String sdkId = service.expectTrait(ServiceTrait.class).getSdkId();
        final String serviceId = sdkId.replace("-", "").replace(" ", "").toLowerCase();
        return serviceId.equalsIgnoreCase("cloudfrontkeyvaluestore");
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                    RuntimeClientPlugin.builder()
                         .configFields(
                            ListUtils.of(
                                    ConfigField.builder()
                                            .name(AwsSignatureVersion4aUtils.V4A_SIGNER_INTERFACE_NAME)
                                            .type(SymbolUtils.createValueSymbolBuilder(
                                                            AwsSignatureVersion4aUtils.V4A_SIGNER_INTERFACE_NAME)
                                                    .build())
                                            .documentation("Signature Version 4a (SigV4a) Signer")
                                            .build()
                            )
                        )
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {

        if (!isCFKVSService(model, model.expectShape(settings.getService(), ServiceShape.class))) {
            return;
        }

        ServiceShape serviceShape = settings.getService(model);
        goDelegator.useShapeWriter(serviceShape, writer -> {
            writerSignerInterface(writer);
            writerConfigFieldResolver(writer, serviceShape);
        });

    }


    private void writerSignerInterface(GoWriter writer) {
        AwsSignatureVersion4aUtils.writerSignerInterface(writer);
    }

    private void writerConfigFieldResolver(GoWriter writer, ServiceShape serviceShape) {
        AwsSignatureVersion4aUtils.writerConfigFieldResolver(writer, serviceShape);
    }

}
