package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.utils.ListUtils;

public class MachineLearningCustomizations implements GoIntegration {
    private static final String ADD_PREDICT_ENDPOINT = "AddPredictEndpointMiddleware";
    private static final String ENDPOINT_ACCESSOR = "getPredictEndpoint";

    @Override
    public byte getOrder() {
        // This needs to be run after the generic endpoint resolver gets added
        return 50;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate(MachineLearningCustomizations::isPredict)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ADD_PREDICT_ENDPOINT,
                                        AwsCustomGoDependency.MACHINE_LEARNING_CUSTOMIZATION).build())
                                .functionArguments(ListUtils.of(
                                        SymbolUtils.createValueSymbolBuilder(ENDPOINT_ACCESSOR).build()
                                ))
                                .build())
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceShape service = settings.getService(model);
        if (!isMachineLearning(model, service)) {
            return;
        }

        TopDownIndex.of(model).getContainedOperations(service).stream()
                .filter(shape -> shape.getId().getName(service).equalsIgnoreCase("Predict"))
                .findAny()
                .ifPresent(operation -> {
                    goDelegator.useShapeWriter(operation, writer -> writeEndpointAccessor(
                            writer, model, symbolProvider, operation));
                });
    }

    private void writeEndpointAccessor(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            OperationShape operation
    ) {
        StructureShape input = ProtocolUtils.expectInput(model, operation);
        writer.openBlock("func $L(input interface{}) (*string, error) {", "}", ENDPOINT_ACCESSOR, () -> {
            writer.write("in, ok := input.($P)", symbolProvider.toSymbol(input));
            writer.openBlock("if !ok {", "}", () -> {
                writer.addUseImports(SmithyGoDependency.SMITHY);
                writer.addUseImports(SmithyGoDependency.FMT);
                writer.write("return nil, &smithy.SerializationError{Err: fmt.Errorf("
                        + "\"expected $P, but was %T\", input)}", symbolProvider.toSymbol(input));
            });
            writer.write("return in.PredictEndpoint, nil");
        });
    }

    private static boolean isPredict(Model model, ServiceShape service, OperationShape operation) {
        return isMachineLearning(model, service) && operation.getId().getName(service).equalsIgnoreCase("Predict");
    }

    private static boolean isMachineLearning(Model model, ServiceShape service) {
        return service.expectTrait(ServiceTrait.class).getSdkId().equalsIgnoreCase("Machine Learning");
    }
}
