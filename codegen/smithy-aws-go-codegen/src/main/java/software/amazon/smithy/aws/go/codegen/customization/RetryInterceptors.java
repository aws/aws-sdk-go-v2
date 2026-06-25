package software.amazon.smithy.aws.go.codegen.customization;

import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

import java.util.List;
import software.amazon.smithy.go.codegen.GoCodegenContext;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;

// We add retry interceptors here since Smithy clients don't have a retry loop atm.
public class RetryInterceptors implements GoIntegration {
    private static RuntimeClientPlugin interceptor(String name) {
        return RuntimeClientPlugin.builder()
                .registerMiddleware(
                        MiddlewareRegistrar.builder()
                                .resolvedFunction(buildPackageSymbol(name))
                                .useClientOptions()
                                .build()
                )
                .build();
    }

    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return List.of(
                interceptor("addInterceptBeforeRetryLoop"),
                interceptor("addInterceptAttempt")
        );
    }

    @Override
    public void writeAdditionalFiles(GoCodegenContext ctx) {
        ctx.writerDelegator().useFileWriter("api_client.go", ctx.settings().getModuleName(), writer -> {
            writer.write("""
                    func addInterceptBeforeRetryLoop(stack *middleware.Stack, opts Options) error {
                        return stack.Finalize.Insert(&smithyhttp.InterceptBeforeRetryLoop{
                            Interceptors: opts.Interceptors.BeforeRetryLoop,
                        }, "Retry", middleware.Before)
                    }

                    func addInterceptAttempt(stack *middleware.Stack, opts Options) error {
                        return stack.Finalize.Insert(&smithyhttp.InterceptAttempt{
                            BeforeAttempt: opts.Interceptors.BeforeAttempt,
                            AfterAttempt: opts.Interceptors.AfterAttempt,
                        }, "Retry", middleware.After)
                    }
                    """);
        });
    }
}
