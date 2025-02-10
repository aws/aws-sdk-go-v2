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
                $aws:D $awsMiddleware:D

                type SetCredentialSourceMiddleware struct {
                       ua *awsmiddleware.RequestUserAgent
                       options Options
                }

                func (m SetCredentialSourceMiddleware) ID() string { return "SetCredentialSourceMiddleware" }

                func (m SetCredentialSourceMiddleware) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (
                       out middleware.BuildOutput, metadata middleware.Metadata, err error,
                ) {
                       asChain, ok := m.options.Credentials.(aws.CredentialProviderChain)
                       if !ok {
                               return next.HandleBuild(ctx, in)
                       }
                       credChain := asChain.CredentialChain()
                       for _, source := range credChain {
                               m.ua.AddCredentialsSource(source)
                       }
                       return next.HandleBuild(ctx, in)
                }

                func addCredentialSource(stack *middleware.Stack, options Options) error {
                       ua, err := getOrAddRequestUserAgent(stack)
                       if err != nil {
                               return err
                       }

                       mw := SetCredentialSourceMiddleware{ua: ua, options: options}
                       return stack.Build.Insert(&mw, "UserAgent", middleware.Before)
                }
                """,
                Map.of(
                        "aws", AwsGoDependency.AWS_CORE,
                        "awsMiddleware", AwsGoDependency.AWS_MIDDLEWARE,
                        "stack", SmithyGoDependency.SMITHY_MIDDLEWARE.struct("Stack")
                )));
    }



}
