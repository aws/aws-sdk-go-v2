package software.amazon.smithy.aws.go.codegen;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.utils.ListUtils;

public final class RegisterServiceMetadataMiddleware implements GoIntegration {
    List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();

    @Override
    public byte getOrder(){
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
        Symbol serviceMetadataProvider = SymbolUtils.createValueSymbolBuilder(
                "RegisterServiceMetadata", AwsGoDependency.AWS_MIDDLEWARE).build();

        for (ShapeId operationId: service.getAllOperations()) {
            String middlewareName = getServiceMetadataMiddlewareName(operationId);
            OperationShape operation = model.expectShape(operationId, OperationShape.class);
            goDelegator.useShapeWriter(operation, writer -> {
                Optional<ServiceTrait> serviceTrait = service.getTrait(ServiceTrait.class);
                Optional<SigV4Trait> sigV4Trait = service.getTrait(SigV4Trait.class);
                writer.openBlock("func $L(region string) $T {", "}",
                        middlewareName, serviceMetadataProvider, () -> {
                    StringBuilder builder = new StringBuilder();
                    builder.append(" return $T{\n");
                    builder.append("Region: region,\n");

                    if (serviceTrait.isPresent()) {
                        ServiceTrait trait = serviceTrait.get();
                        String sdkIdAsSymbol = ServiceIdUtils.toTitleCase(trait.getSdkId()).replace(" ", "");
                        builder.append(String.format("ServiceName: \"%s\",\n", trait.getSdkId()));
                        builder.append(String.format("ServiceID: \"%s\",\n", sdkIdAsSymbol));
                        // TODO: EndpointID can be different but is not modeled in Smithy.
                    }

                    if (sigV4Trait.isPresent()) {
                        SigV4Trait trait = sigV4Trait.get();
                        builder.append(String.format("SigningName: \"%s\",\n", trait.getName()));
                    }
                    builder.append(String.format("OperationName: \"%s\",\n", operation.getId().getName()));
                    builder.append("}");

                    writer.write(builder.toString(), serviceMetadataProvider);
                });
            });
        }
    }

        @Override
    public void processFinalizedModel(GoSettings settings, Model model) {
        ServiceShape service = settings.getService(model);

        for (ShapeId operationId : service.getAllOperations()) {
            String middlewareName = getServiceMetadataMiddlewareName(operationId);
            OperationShape operation = model.expectShape(operationId, OperationShape.class);
            RuntimeClientPlugin runtimeClientPlugin = RuntimeClientPlugin.builder()
                    .operationPredicate((predicateModel, predicateService, predicateOperation) -> {
                        if (operation.equals(predicateOperation)) {
                            return true;
                        }
                        return false;
                    })
                    .registerMiddleware(MiddlewareRegistrar.builder()
                            .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                    middlewareName).build())
                            .registerBefore(MiddlewareRegistrar.StackStep.INITIALIZE)
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

    private String getServiceMetadataMiddlewareName(ShapeId operationID) {
        return "newServiceMetadataMiddleware_op"+ operationID.getName();
    }
}
