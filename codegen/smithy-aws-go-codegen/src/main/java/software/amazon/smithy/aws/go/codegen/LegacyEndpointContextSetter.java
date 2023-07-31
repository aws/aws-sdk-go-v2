package software.amazon.smithy.aws.go.codegen;

import java.util.ArrayList;
import java.util.List;

import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.MiddlewareIdentifier;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;

public class LegacyEndpointContextSetter implements GoIntegration {
    
    private final List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();

    public static final String MIDDLEWARE_ID = "legacyEndpointContextSetter";
    public static final String MIDDLEWARE_ADDER = String.format("add%s", MIDDLEWARE_ID);

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first. Needs to execute after Rules Engine endpoint
     * resolution middleware insertion.
     *
     * @return Returns the sort order, defaults to 127.
     */
    @Override
    public byte getOrder() {
            return -128;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
            return runtimeClientPlugins;
    }

    @Override
    public void processFinalizedModel(GoSettings settings, Model model) {

            var serviceShape = settings.getService(model);

            runtimeClientPlugins.add(RuntimeClientPlugin.builder()
                            .servicePredicate((m, s) -> s.equals(serviceShape))
                            .registerMiddleware(MiddlewareRegistrar.builder()
                                            .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                                            MIDDLEWARE_ADDER)
                                                            .build())
                                            .useClientOptions()
                                            .build())
                            .build());

    }

    @Override
    public void renderPreEndpointResolutionHook(GoSettings settings, GoWriter writer, Model model) {
        writer.write(
                """
                        if $T(ctx) {
                            return next.HandleSerialize(ctx, in)
                        }
                """,
                SymbolUtils.createValueSymbolBuilder("GetRequiresLegacyEndpoints", AwsGoDependency.AWS_MIDDLEWARE).build()
        );
    }

    @Override
    public void writeAdditionalFiles(
                    GoSettings settings,
                    Model model,
                    SymbolProvider symbolProvider,
                    GoDelegator goDelegator) {

            var serviceShape = settings.getService(model);
            goDelegator.useShapeWriter(serviceShape, writer -> {


                    GoStackStepMiddlewareGenerator middleware = GoStackStepMiddlewareGenerator
                                    .createInitializeStepMiddleware(
                                                    MIDDLEWARE_ID,
                                                    MiddlewareIdentifier.string(MIDDLEWARE_ID));
                    middleware.writeMiddleware(writer, this::generateMiddlewareResolverBody,
                                    this::generateMiddlewareStructureMembers);

                    writer.write(
                                    """
                                                            func $L(stack $P, o Options) error {
                                                                    return stack.Initialize.Add(&$L{
                                                                            LegacyResolver: o.EndpointResolver,
                                                                    }, middleware.Before)
                                                            }
                                                    """,
                                    MIDDLEWARE_ADDER,
                                    SymbolUtils.createPointableSymbolBuilder("Stack",
                                                    SmithyGoDependency.SMITHY_MIDDLEWARE).build(),
                                    MIDDLEWARE_ID);
                    writer.write("");
            });
    }

    private void generateMiddlewareResolverBody(GoStackStepMiddlewareGenerator g, GoWriter writer) {
            writer.write(
                            """
                                                    if m.LegacyResolver != nil {
                                                        ctx = $T(ctx, true)
                                                    }

                                                    return next.HandleInitialize(ctx, in)
                                            """,
                            SymbolUtils.createValueSymbolBuilder("SetRequiresLegacyEndpoints", AwsGoDependency.AWS_MIDDLEWARE).build()
                            );
    }

    private void generateMiddlewareStructureMembers(GoStackStepMiddlewareGenerator g, GoWriter writer) {
            writer.write("LegacyResolver $L", "EndpointResolver");
    }
}
