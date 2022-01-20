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
 * This integration configures the S3 client for Signature Version 4a
 */
public class S3SignatureVersion4a implements GoIntegration {
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

    private static final List<String> DISABLE_URI_PATH_ESCAPE = ListUtils.of("com.amazonaws.s3#AmazonS3");

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        Symbol resolver = SymbolUtils.createValueSymbolBuilder(AwsSignatureVersion4aUtils.RESOLVE_CREDENTIAL_PROVIDER)
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
                        .build(),
                // Add HTTPSigner middleware to operation stack
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3SignatureVersion4a::isS3Service)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        AwsSignatureVersion4aUtils.REGISTER_MIDDLEWARE_FUNCTION).build())
                                .useClientOptions()
                                .build())
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3SignatureVersion4a::isS3Service)
                        .addConfigFieldResolver(
                                ConfigFieldResolver.builder()
                                        .location(ConfigFieldResolver.Location.CLIENT)
                                        .target(ConfigFieldResolver.Target.INITIALIZATION)
                                        .resolver(SymbolUtils.createValueSymbolBuilder(
                                                AwsSignatureVersion4aUtils.SIGNER_RESOLVER).build())
                                        .build())
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

        if (!isS3Service(model, model.expectShape(settings.getService(), ServiceShape.class))) {
            return;
        }

        ServiceShape serviceShape = settings.getService(model);
        goDelegator.useShapeWriter(serviceShape, writer -> {
            writeCredentialProviderResolver(writer);
            writeMiddlewareRegister(model, writer, serviceShape);
            writerSignerInterface(writer);
            writerConfigFieldResolver(writer, serviceShape);
            writeNewV4ASignerFunc(writer, serviceShape);
        });

    }

    private void writeCredentialProviderResolver(GoWriter writer) {
        AwsSignatureVersion4aUtils.writeCredentialProviderResolver(writer);
    }

    private void writerSignerInterface(GoWriter writer) {
        AwsSignatureVersion4aUtils.writerSignerInterface(writer);
    }

    private void writerConfigFieldResolver(GoWriter writer, ServiceShape serviceShape) {
        AwsSignatureVersion4aUtils.writerConfigFieldResolver(writer, serviceShape);
    }

    private void writeNewV4ASignerFunc(GoWriter writer, ServiceShape serviceShape) {
        AwsSignatureVersion4aUtils.writeNewV4ASignerFunc(writer, serviceShape,
                DISABLE_URI_PATH_ESCAPE.contains(serviceShape.getId().toString()));
    }

    private void writeMiddlewareRegister(Model model, GoWriter writer, ServiceShape serviceShape) {
        AwsSignatureVersion4aUtils.writeMiddlewareRegister(model, writer, serviceShape,
                AwsCustomGoDependency.S3_CUSTOMIZATION);
    }
}
