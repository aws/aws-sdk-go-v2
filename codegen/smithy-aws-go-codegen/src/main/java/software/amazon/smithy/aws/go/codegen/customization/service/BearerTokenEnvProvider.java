package software.amazon.smithy.aws.go.codegen.customization.service;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

import java.util.List;
import java.util.Map;
import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.go.codegen.GoCodegenContext;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.traits.HttpBearerAuthTrait;

public class BearerTokenEnvProvider implements GoIntegration {
    // Set of signing names to which this customization is applied.
    private static final List<String> SIGNING_NAMES = List.of(
            "bedrock"
    );

    private static final ConfigFieldResolver BEARER_TOKEN_RESOLVER =
            ConfigFieldResolver.builder()
                    .location(ConfigFieldResolver.Location.CLIENT)
                    .target(ConfigFieldResolver.Target.INITIALIZATION)
                    .resolver(buildPackageSymbol("resolveEnvBearerToken"))
                    .build();

    private static boolean isApplied(Model model, ServiceShape service) {
        return service.hasTrait(SigV4Trait.class)
                && SIGNING_NAMES.contains(service.expectTrait(SigV4Trait.class).getName())
                && service.hasTrait(HttpBearerAuthTrait.class);
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return List.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(BearerTokenEnvProvider::isApplied)
                        .addConfigFieldResolver(BEARER_TOKEN_RESOLVER)
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(GoCodegenContext ctx) {
        var service = ctx.settings().getService(ctx.model());
        if (isApplied(ctx.model(), service)) {
            ctx.writerDelegator().useFileWriter("api_client.go", ctx.settings().getModuleName(), bearerTokenResolver(service));
        }
    }

    private GoWriter.Writable bearerTokenResolver(ServiceShape service) {
        var envSuffix = service.expectTrait(SigV4Trait.class).getName()
                .toUpperCase()
                .replace(' ', '_')
                .replace('-', '_');
        return goTemplate("""
                $context:D $os:D $bearer:D $awsmiddleware:D
                func resolveEnvBearerToken(options *Options) {
                    token := os.Getenv("AWS_BEARER_TOKEN_$envSuffix:L")
                    if len(token) == 0 { return }

                    options.BearerAuthTokenProvider = bearer.TokenProviderFunc(func(ctx context.Context) (bearer.Token, error) {
                        return bearer.Token{Value: token}, nil
                    })
                    options.AuthSchemePreference = []string{"httpBearerAuth"}
                    options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
                        ua, err := getOrAddRequestUserAgent(stack)
                        if err != nil {
                            return err
                        }

                        ua.AddUserAgentFeature(awsmiddleware.UserAgentFeatureBearerServiceEnvVars)
                        return nil
                    })
                }
                """,
                Map.of(
                        "context", SmithyGoDependency.CONTEXT,
                        "os", SmithyGoDependency.OS,
                        "bearer", SmithyGoDependency.SMITHY_AUTH_BEARER,
                        "awsmiddleware", AwsCustomGoDependency.AWS_MIDDLEWARE,
                        "envSuffix", envSuffix
                ));
    }
}
