package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import java.util.Map;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoCodegenContext;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;

import static software.amazon.smithy.aws.go.codegen.customization.AccountIDEndpointRouting.hasAccountIdEndpoints;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * Tracks the retry mode being used by the caller.
 */
public class AccountIdEndpointModeUserAgent implements GoIntegration {
    private static final MiddlewareRegistrar MIDDLEWARE = MiddlewareRegistrar.builder()
            .resolvedFunction(buildPackageSymbol("addUserAgentAccountIDEndpointMode"))
            .useClientOptions()
            .build();

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return List.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(AccountIDEndpointRouting::hasAccountIdEndpoints)
                        .registerMiddleware(MIDDLEWARE)
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(GoCodegenContext ctx) {
        if (!hasAccountIdEndpoints(ctx.model(), ctx.settings().getService(ctx.model()))) {
            return;
        }

        ctx.writerDelegator().useFileWriter("api_client.go", ctx.settings().getModuleName(), goTemplate("""
                $aws:D $awsMiddleware:D
                func addUserAgentAccountIDEndpointMode(stack $stack:P, options Options) error {
                    ua, err := getOrAddRequestUserAgent(stack)
                    if err != nil {
                        return err
                    }

                    switch options.AccountIDEndpointMode {
                    case aws.AccountIDEndpointModePreferred:
                        ua.AddUserAgentFeature(awsmiddleware.UserAgentFeatureAccountIDModePreferred)
                    case aws.AccountIDEndpointModeRequired:
                        ua.AddUserAgentFeature(awsmiddleware.UserAgentFeatureAccountIDModeRequired)
                    case aws.AccountIDEndpointModeDisabled:
                        ua.AddUserAgentFeature(awsmiddleware.UserAgentFeatureAccountIDModeDisabled)
                    }
                    return nil
                }""",
                Map.of(
                        "aws", AwsGoDependency.AWS_CORE,
                        "awsMiddleware", AwsGoDependency.AWS_MIDDLEWARE,
                        "stack", SmithyGoDependency.SMITHY_MIDDLEWARE.struct("Stack")
                )));
    }
}
