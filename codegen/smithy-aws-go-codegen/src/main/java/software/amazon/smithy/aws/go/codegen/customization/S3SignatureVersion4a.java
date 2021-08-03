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
    private static final String RESOLVE_CREDENTIAL_PROVIDER = "resolveCredentialProvider";
    private static final String REGISTER_MIDDLEWARE_FUNCTION = "swapWithCustomHTTPSignerMiddleware";
    private static final String V4A_SIGNER_INTERFACE_NAME = "httpSignerV4a";
    private static final String SIGNER_OPTION_FIELD_NAME = V4A_SIGNER_INTERFACE_NAME;
    private static final String NEW_SIGNER_FUNC_NAME = "newDefaultV4aSigner";
    private static final String SIGNER_RESOLVER = "resolveHTTPSignerV4a";

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
                .build(),
                // Add HTTPSigner middleware to operation stack
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3SignatureVersion4a::isS3Service)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        S3SignatureVersion4a.REGISTER_MIDDLEWARE_FUNCTION).build())
                                .useClientOptions()
                                .build())
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3SignatureVersion4a::isS3Service)
                        .addConfigFieldResolver(
                                ConfigFieldResolver.builder()
                                        .location(ConfigFieldResolver.Location.CLIENT)
                                        .target(ConfigFieldResolver.Target.INITIALIZATION)
                                        .resolver(SymbolUtils.createValueSymbolBuilder(SIGNER_RESOLVER).build())
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

    private void writerSignerInterface(GoWriter writer) {
        writer.openBlock("type $L interface {", "}", V4A_SIGNER_INTERFACE_NAME, () -> {
            writer.addUseImports(SmithyGoDependency.CONTEXT);
            writer.addUseImports(AwsGoDependency.AWS_CORE);
            writer.addUseImports(AwsCustomGoDependency.S3_SIGV4A_CUSTOMIZATION);
            writer.addUseImports(SmithyGoDependency.NET_HTTP);
            writer.addUseImports(SmithyGoDependency.TIME);
            writer.write("SignHTTP(ctx context.Context, credentials v4a.Credentials, r *http.Request, "
                    + "payloadHash string, service string, regionSet []string, signingTime time.Time, "
                    + "optFns ...func(*v4a.SignerOptions)) error");
        });
    }

    private void writerConfigFieldResolver(GoWriter writer, ServiceShape serviceShape) {
        writer.openBlock("func $L(o *Options) {", "}", SIGNER_RESOLVER, () -> {
            writer.openBlock("if o.$L != nil {", "}", SIGNER_OPTION_FIELD_NAME, () -> writer.write("return"));
            writer.write("o.$L = $L(*o)", SIGNER_OPTION_FIELD_NAME, NEW_SIGNER_FUNC_NAME);
        });
        writer.write("");
    }

    private void writeNewV4ASignerFunc(GoWriter writer, ServiceShape serviceShape) {
        Symbol signerSymbol = SymbolUtils.createValueSymbolBuilder("Signer",
                AwsCustomGoDependency.S3_SIGV4A_CUSTOMIZATION).build();
        Symbol newSignerSymbol = SymbolUtils.createValueSymbolBuilder("NewSigner",
                AwsCustomGoDependency.S3_SIGV4A_CUSTOMIZATION).build();
        Symbol signerOptionsSymbol = SymbolUtils.createPointableSymbolBuilder("SignerOptions",
                AwsCustomGoDependency.S3_SIGV4A_CUSTOMIZATION).build();

        writer.openBlock("func $L(o Options) *$T {", "}", NEW_SIGNER_FUNC_NAME, signerSymbol, () -> {
            writer.openBlock("return $T(func(so $P) {", "})", newSignerSymbol, signerOptionsSymbol, () -> {
                writer.write("so.Logger = o.$L", AddAwsConfigFields.LOGGER_CONFIG_NAME);
                writer.write("so.LogSigning = o.$L.IsSigning()", AddAwsConfigFields.LOG_MODE_CONFIG_NAME);
                if (DISABLE_URI_PATH_ESCAPE.contains(serviceShape.getId().toString())) {
                    writer.write("so.DisableURIPathEscaping = true");
                }
            });
        });
    }


    private void writeMiddlewareRegister(Model model, GoWriter writer, ServiceShape serviceShape) {
        writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);
        Symbol registerSigningMiddleware = SymbolUtils.createValueSymbolBuilder(
                "RegisterSigningMiddleware", AwsCustomGoDependency.S3_CUSTOMIZATION
        ).build();

        writer.openBlock("func $L(stack $P, o Options) error {", "}", REGISTER_MIDDLEWARE_FUNCTION,
                SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE).build(), () -> {
                    Symbol newMiddlewareSymbol = SymbolUtils.createValueSymbolBuilder(
                            "NewSignHTTPRequestMiddleware", AwsCustomGoDependency.S3_CUSTOMIZATION).build();
                    Symbol middlewareOptionsSymbol = SymbolUtils.createValueSymbolBuilder(
                            "SignHTTPRequestMiddlewareOptions", AwsCustomGoDependency.S3_CUSTOMIZATION).build();

                    writer.openBlock("mw := $T($T{", "})", newMiddlewareSymbol, middlewareOptionsSymbol, () -> {
                                writer.write("CredentialsProvider: o.$L,", AddAwsConfigFields.CREDENTIALS_CONFIG_NAME);
                                writer.write("V4Signer: o.$L,", AwsSignatureVersion4.SIGNER_CONFIG_FIELD_NAME);
                                writer.write("V4aSigner: o.$L,", SIGNER_OPTION_FIELD_NAME);
                                writer.write("LogSigning: o.$L.IsSigning(),", AddAwsConfigFields.LOG_MODE_CONFIG_NAME);
                    });

                    writer.write("return $T(stack, mw)", registerSigningMiddleware);
        });
        writer.write("");
    }

}
