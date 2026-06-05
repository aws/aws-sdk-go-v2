package software.amazon.smithy.aws.go.codegen;

import java.util.ArrayList;
import java.util.List;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.MiddlewareStackStep;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ToShapeId;
import software.amazon.smithy.utils.ListUtils;

/**
 * Emits a single shared {@code newServiceMetadataMiddleware(region, operation)}
 * helper per service and registers a per-operation plugin that calls it with
 * the appropriate operation name literal.
 *
 * <p>The shared function eliminates N redundant per-operation function
 * definitions. The plugin registration stays per-operation because the
 * middleware uses Initialize.Add(Before) which requires correct prepend
 * ordering relative to other per-operation Initialize middlewares.
 */
public final class RegisterServiceMetadataMiddleware implements GoIntegration {
    private final List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();

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
        ServiceShape service = settings.getService(model);
        Symbol serviceMetadataProvider = SymbolUtils.createPointableSymbolBuilder(
                "RegisterServiceMetadata", AwsGoDependency.AWS_MIDDLEWARE).build();

        goDelegator.useFileWriter("api_client.go", settings.getModuleName(), writer -> {
            writer.openBlock("func newServiceMetadataMiddleware(region, operation string) $P {", "}",
                    serviceMetadataProvider, () -> {
                        writer.write("return &$T{", serviceMetadataProvider);
                        writer.write("Region: region,");
                        writer.write("ServiceID: ServiceID,");
                        writer.write("OperationName: operation,");
                        writer.write("}");
                    });
        });
    }

    @Override
    public void processFinalizedModel(GoSettings settings, Model model) {
        ServiceShape service = settings.getService(model);
        TopDownIndex index = TopDownIndex.of(model);

        for (ToShapeId operation : index.getContainedOperations(service)) {
            OperationShape operationShape = model.expectShape(operation.toShapeId(), OperationShape.class);
            String opName = operationShape.getId().getName(service);
            runtimeClientPlugins.add(RuntimeClientPlugin.builder()
                    .operationPredicate((m, s, o) -> s.equals(service) && o.equals(operationShape))
                    .registerMiddleware(MiddlewareRegistrar.builder()
                            .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                    "newServiceMetadataMiddleware").build())
                            .registerBefore(MiddlewareStackStep.INITIALIZE)
                            .functionArguments(ListUtils.of(
                                    SymbolUtils.createValueSymbolBuilder("options.Region").build(),
                                    SymbolUtils.createValueSymbolBuilder("\"" + opName + "\"").build()
                            ))
                            .build())
                    .build());
        }
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return runtimeClientPlugins;
    }
}
