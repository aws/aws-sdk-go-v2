package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import java.util.ArrayList;
import software.amazon.smithy.aws.go.codegen.AddAwsConfigFields;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.aws.go.codegen.AwsSignatureVersion4;
import software.amazon.smithy.aws.go.codegen.AwsSignatureVersion4aUtils;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.aws.traits.auth.SigV4ATrait;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
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
import software.amazon.smithy.model.traits.AuthTrait;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.SetUtils;

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

    private final List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();


    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return runtimeClientPlugins;
    }

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        ServiceShape service = settings.getService(model);
        if (!isCFKVSService(model, service)) {
            return model;
        }

        if (settings.getService(model).hasTrait(SigV4ATrait.class)) {
            return model;
        }

        var v4a = SigV4ATrait.builder()
                .name(service.expectTrait(SigV4Trait.class).getName())
                .build();

        return model.toBuilder()
                .addShape(
                        service.toBuilder()
                                .addTrait(v4a)
                                // FUTURE: https://github.com/aws/smithy-go/issues/493
                                // we are keeping sigv4 at the end of this list (it will never be selected)
                                // as a stopgap to drive codegen of payload checksum routines
                                .addTrait(new AuthTrait(SetUtils.of(SigV4ATrait.ID, SigV4Trait.ID)))
                                .build()
                )
                .build();
    }

    @Override
    public void processFinalizedModel(GoSettings settings, Model model) {
        if (!isCFKVSService(model, model.expectShape(settings.getService(), ServiceShape.class))) {
            return;
        }
        runtimeClientPlugins.add(
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
                    .build());
        runtimeClientPlugins.add(
                RuntimeClientPlugin.builder()
                    .servicePredicate(CloudFrontKVSSigV4a::isCFKVSService)
                    .addConfigFieldResolver(
                            ConfigFieldResolver.builder()
                                    .location(ConfigFieldResolver.Location.CLIENT)
                                    .target(ConfigFieldResolver.Target.INITIALIZATION)
                                    .resolver(SymbolUtils.createValueSymbolBuilder(
                                            AwsSignatureVersion4aUtils.SIGNER_RESOLVER).build())
                                    .build())
                    .build());
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
            writeNewV4ASignerFunc(writer, serviceShape);
        });

    }


    private void writerSignerInterface(GoWriter writer) {
        AwsSignatureVersion4aUtils.writerSignerInterface(writer);
    }

    private void writeNewV4ASignerFunc(GoWriter writer, ServiceShape serviceShape) {
        AwsSignatureVersion4aUtils.writeNewV4ASignerFunc(writer, serviceShape);
    }

    private void writerConfigFieldResolver(GoWriter writer, ServiceShape serviceShape) {
        AwsSignatureVersion4aUtils.writerConfigFieldResolver(writer, serviceShape);
    }

}
