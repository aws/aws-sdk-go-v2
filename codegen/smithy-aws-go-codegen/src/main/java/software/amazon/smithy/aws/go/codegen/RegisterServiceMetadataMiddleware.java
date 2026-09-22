package software.amazon.smithy.aws.go.codegen;

import java.util.List;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;

/**
 * Sets service and operation metadata (ServiceID, Region, OperationName, and the
 * legacy-endpoints flag) on the request context.
 *
 * <p>Rather than registering a per-request Initialize-step middleware
 * (RegisterServiceMetadata), the values are set directly in invokeOperation via
 * an operation context resolver. The generated resolveServiceMetadata function
 * writes each value using the exported setters in aws/middleware. The public
 * RegisterServiceMetadata type is unchanged and still available for external
 * callers.
 */
public final class RegisterServiceMetadataMiddleware implements GoIntegration {
    private static final String RESOLVER = "resolveServiceMetadata";

    @Override
    public byte getOrder() {
        return 30;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        Symbol context = SymbolUtils.createValueSymbolBuilder(
                "Context", SmithyGoDependency.CONTEXT).build();
        Symbol setServiceID = SymbolUtils.createValueSymbolBuilder(
                "SetServiceID", AwsGoDependency.AWS_MIDDLEWARE).build();
        Symbol setRegion = SymbolUtils.createValueSymbolBuilder(
                "SetRegion", AwsGoDependency.AWS_MIDDLEWARE).build();
        Symbol setOperationName = SymbolUtils.createValueSymbolBuilder(
                "SetOperationName", AwsGoDependency.AWS_MIDDLEWARE).build();
        Symbol setRequiresLegacyEndpoints = SymbolUtils.createValueSymbolBuilder(
                "SetRequiresLegacyEndpoints", AwsGoDependency.AWS_MIDDLEWARE).build();

        goDelegator.useFileWriter("api_client.go", settings.getModuleName(), writer -> {
            writer.openBlock("func $L(ctx $T, options Options, operation string) $T {", "}",
                    RESOLVER, context, context, () -> {
                        writer.write("ctx = $T(ctx, ServiceID)", setServiceID);
                        writer.openBlock("if options.Region != \"\" {", "}", () -> {
                            writer.write("ctx = $T(ctx, options.Region)", setRegion);
                        });
                        writer.write("ctx = $T(ctx, operation)", setOperationName);
                        writer.openBlock("if options.EndpointResolver != nil {", "}", () -> {
                            writer.write("ctx = $T(ctx, true)", setRequiresLegacyEndpoints);
                        });
                        writer.write("return ctx");
                    });
        });
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .addOperationContextResolver(
                                SymbolUtils.createValueSymbolBuilder(RESOLVER).build())
                        .build());
    }
}
