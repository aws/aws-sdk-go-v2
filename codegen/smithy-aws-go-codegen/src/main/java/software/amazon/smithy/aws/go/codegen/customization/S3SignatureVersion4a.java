package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import software.amazon.smithy.aws.go.codegen.AddAwsConfigFields;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.aws.go.codegen.AwsSignatureVersion4;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

/**
 * This integration configures the S3 client for Signature Version 4a
 */
public class S3SignatureVersion4a implements GoIntegration {
    private static final String RESOLVE_CREDENTIAL_PROVIDER = "resolveCredentialProvider";

    /**
     * Return true if service is Amazon S3.
     *
     * @param model   is the generation model.
     * @param service is the service shape being audited.
     */
    private static boolean isS3Service(Model model, ServiceShape service) {
        String serviceId = service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("S3");
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        Symbol resolver = SymbolUtils.createValueSymbolBuilder(RESOLVE_CREDENTIAL_PROVIDER)
                .build();

        return ListUtils.of(RuntimeClientPlugin.builder()
                .addConfigFieldResolver(ConfigFieldResolver.builder()
                        .location(ConfigFieldResolver.Location.CLIENT)
                        .target(ConfigFieldResolver.Target.FINALIZATION)
                        .resolver(resolver)
                        .build())
                .addConfigFieldResolver(ConfigFieldResolver.builder()
                        .location(ConfigFieldResolver.Location.OPERATION)
                        .target(ConfigFieldResolver.Target.FINALIZATION)
                        .resolver(resolver)
                        .build())
                .servicePredicate((model, serviceShape) -> {
                    if (!S3SignatureVersion4a.isS3Service(model, serviceShape)) {
                        return false;
                    }
                    return AwsSignatureVersion4.isSupportedAuthentication(model, serviceShape);
                })
                .build());
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        if (!isS3Service(model, model.expectShape(settings.getService(), ServiceShape.class))) {
            return;
        }

        goDelegator.useShapeWriter(settings.getService(model), this::writeCredentialProviderResolver);
    }

    private void writeCredentialProviderResolver(GoWriter writer) {
        final String fieldName = AddAwsConfigFields.CREDENTIALS_CONFIG_NAME;

        writer.openBlock("func $L(o *Options) {", "}", RESOLVE_CREDENTIAL_PROVIDER, () -> {
            writer.openBlock("if o.$L == nil {", "}", fieldName, () -> writer.write("return"));

            Symbol adaptorSymbol = SymbolUtils.createPointableSymbolBuilder("SymmetricCredentialAdaptor",
                    AwsCustomGoDependency.S3_SIGV4A_CUSTOMIZATION).build();
            Symbol credentialProvider = SymbolUtils.createPointableSymbolBuilder("CredentialsProvider",
                    AwsCustomGoDependency.S3_SIGV4A_CUSTOMIZATION).build();

            writer.openBlock("if _, ok := o.$L.($T); ok {", "}", fieldName, credentialProvider,
                    () -> writer.write("return"));
            writer.write("");

            Symbol anonymousCredentials = SymbolUtils.createPointableSymbolBuilder("AnonymousCredentials",
                    AwsGoDependency.AWS_CORE).build();
            writer.openBlock("switch o.$L.(type) {", "}", fieldName, () -> {
                writer.openBlock("case $T, $P:", "", anonymousCredentials, anonymousCredentials, () -> {
                    writer.write("return");
                });
            });
            writer.write("");

            writer.write("o.$L = &$T{SymmetricProvider: o.$L}", fieldName, adaptorSymbol, fieldName);
        });
    }
}
