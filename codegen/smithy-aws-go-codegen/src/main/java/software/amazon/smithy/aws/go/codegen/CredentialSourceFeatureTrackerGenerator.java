package software.amazon.smithy.aws.go.codegen;

import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.GoCodegenContext;

import java.util.List;
import java.util.Map;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * Generates code to track which credential provider was used on
 * the User Agent
 */
public class CredentialSourceFeatureTrackerGenerator implements GoIntegration {

    private static final MiddlewareRegistrar MIDDLEWARE = MiddlewareRegistrar.builder()
            .resolvedFunction(buildPackageSymbol("addCredentialSource"))
            .useClientOptions()
            .build();

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return List.of(
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MIDDLEWARE)
                        .servicePredicate(AwsSignatureVersion4::hasSigV4X)
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(GoCodegenContext ctx) {
        if (!AwsSignatureVersion4.hasSigV4X(ctx.model(), ctx.settings().getService(ctx.model()))) {
            return;
        }

        ctx.writerDelegator().useFileWriter("api_client.go", ctx.settings().getModuleName(), goTemplate("""
                $aws:D

                func addCredentialSource(stack *middleware.Stack, options Options) error {
                       ua, err := getOrAddRequestUserAgent(stack)
                       if err != nil {
                               return err
                       }

                       asProviderSource, ok := options.Credentials.(aws.CredentialProviderSource)
                       if !ok {
                               return nil
                       }

                       for _, source := range asProviderSource.ProviderSources() {
                               ua.AddCredentialsSource(source)
                       }
                       return nil
                }
                """,
                Map.of(
                        "aws", AwsGoDependency.AWS_CORE,
                        "stack", SmithyGoDependency.SMITHY_MIDDLEWARE.struct("Stack")
                )));
    }



}
