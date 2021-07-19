package software.amazon.smithy.aws.go.codegen;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
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
import software.amazon.smithy.model.knowledge.ServiceIndex;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.ToShapeId;
import software.amazon.smithy.model.traits.Trait;
import software.amazon.smithy.utils.ListUtils;

public final class RegisterServiceMetadataMiddleware implements GoIntegration {
    List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();

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
        ServiceIndex serviceIndex = ServiceIndex.of(model);

        TopDownIndex topDownIndex = TopDownIndex.of(model);

        for (ToShapeId operation : topDownIndex.getContainedOperations(service)) {
            String middlewareName = getServiceMetadataMiddlewareName(service, operation.toShapeId());
            OperationShape operationShape = model.expectShape(operation.toShapeId(), OperationShape.class);
            goDelegator.useShapeWriter(operationShape, writer -> {
                writer.openBlock("func $L(region string) $P {", "}",
                        middlewareName, serviceMetadataProvider, () -> {
                            StringBuilder builder = new StringBuilder();
                            builder.append(" return &$T{\n");
                            builder.append("Region: region,\n");
                            builder.append("ServiceID: ServiceID,\n");

                            Map<ShapeId, Trait> authSchemes = serviceIndex.getEffectiveAuthSchemes(service, operationShape);
                            if (authSchemes.containsKey(SigV4Trait.ID)) {
                                SigV4Trait trait = (SigV4Trait) authSchemes.get(SigV4Trait.ID);
                                builder.append(String.format("SigningName: \"%s\",\n", trait.getName()));
                            }
                            builder.append(String.format("OperationName: \"%s\",\n",
                                    operationShape.getId().getName(service)));
                            builder.append("}");

                            writer.write(builder.toString(), serviceMetadataProvider);
                        });
            });
        }
    }

    @Override
    public void processFinalizedModel(GoSettings settings, Model model) {
        ServiceShape service = settings.getService(model);
        TopDownIndex index = TopDownIndex.of(model);

        for (ToShapeId operation : index.getContainedOperations(service)) {
            String middlewareName = getServiceMetadataMiddlewareName(service, operation.toShapeId());
            OperationShape operationShape = model.expectShape(operation.toShapeId(), OperationShape.class);
            RuntimeClientPlugin runtimeClientPlugin = RuntimeClientPlugin.builder()
                    .operationPredicate((m, s, o) -> {
                        if (!s.equals(service)) {
                            return false;
                        }
                        return operationShape.equals(o);
                    })
                    .registerMiddleware(MiddlewareRegistrar.builder()
                            .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                    middlewareName).build())
                            .registerBefore(MiddlewareStackStep.INITIALIZE)
                            .functionArguments(ListUtils.of(
                                    SymbolUtils.createValueSymbolBuilder("options.Region").build()
                            ))
                            .build())
                    .build();
            runtimeClientPlugins.add(runtimeClientPlugin);
        }
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return runtimeClientPlugins;
    }

    private String getServiceMetadataMiddlewareName(ServiceShape service, ShapeId operationID) {
        return "newServiceMetadataMiddleware_op" + operationID.getName(service);
    }
}
