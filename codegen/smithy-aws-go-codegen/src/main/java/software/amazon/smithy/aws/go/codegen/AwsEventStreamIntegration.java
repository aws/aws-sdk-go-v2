package software.amazon.smithy.aws.go.codegen;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;
import java.util.stream.Collectors;
import java.util.stream.Stream;
import software.amazon.smithy.go.codegen.GoEventStreamIndex;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.utils.ListUtils;

public class AwsEventStreamIntegration implements GoIntegration {
    private final Map<ShapeId, Collection<OperationShape>> serviceOperationMap = new HashMap<>();
    private final Map<ShapeId, Collection<OperationShape>> minHttp2 = new HashMap<>();

    @Override
    public byte getOrder() {
        return -127;
    }

    @Override
    public void processFinalizedModel(
            GoSettings settings,
            Model model
    ) {
        var goEventStreamIndex = GoEventStreamIndex.of(model);
        var service = settings.getService();

        Collection<OperationShape> inputOperations = new TreeSet<>();
        Collection<OperationShape> outputOperation = new TreeSet<>();

        goEventStreamIndex.getInputEventStreams(service).ifPresent(shapeIdSetMap ->
                shapeIdSetMap.values().forEach(eventStreamInfos ->
                        eventStreamInfos.forEach(info -> inputOperations.add(info.getOperation()))));

        goEventStreamIndex.getOutputEventStreams(service).ifPresent(shapeIdSetMap ->
                shapeIdSetMap.values().forEach(eventStreamInfos ->
                        eventStreamInfos.forEach(info -> outputOperation.add(info.getOperation()))));

        var operationShapes = Stream.concat(inputOperations.stream(), outputOperation.stream())
                .collect(Collectors.toSet());

        if (operationShapes.size() > 0) {
            serviceOperationMap.put(service, operationShapes);
        }

        var biDirectional = new TreeSet<>(inputOperations);
        biDirectional.retainAll(outputOperation);
        if (biDirectional.size() > 0) {
            minHttp2.put(service, biDirectional);
        }
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        final List<RuntimeClientPlugin> plugins = new ArrayList<>();

        plugins.add(RuntimeClientPlugin.builder()
                .servicePredicate((model, serviceShape) -> serviceOperationMap.containsKey(serviceShape.toShapeId()))
                .addConfigFieldResolver(ConfigFieldResolver.builder()
                        .location(ConfigFieldResolver.Location.OPERATION)
                        .resolver(AwsEventStreamUtils.getEventStreamClientLogModeFinalizerSymbol())
                        .target(ConfigFieldResolver.Target.FINALIZATION)
                        .withOperationName(true)
                        .build())
                .build());

        serviceOperationMap.entrySet().stream()
                .map(entry -> entry.getValue().stream().map(operationShape ->
                                RuntimeClientPlugin.builder()
                                        .operationPredicate((model, service, operation) ->
                                                service.getId() == entry.getKey() && operation.equals(operationShape))
                                        .registerMiddleware(MiddlewareRegistrar.builder()
                                                .resolvedFunction(
                                                        AwsEventStreamUtils.getAddEventStreamOperationMiddlewareSymbol(
                                                                operationShape))
                                                .useClientOptions()
                                                .build())
                                        .build())
                        .collect(Collectors.toSet()))
                .forEach(plugins::addAll);


        minHttp2.entrySet().stream()
                .map(entry -> entry.getValue().stream().map(operationShape ->
                                RuntimeClientPlugin.builder()
                                        .operationPredicate((model, service, operation) ->
                                                service.getId() == entry.getKey() && operation.equals(operationShape))
                                        .registerMiddleware(MiddlewareRegistrar.builder()
                                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                                        "AddRequireMinimumProtocol",
                                                        SmithyGoDependency.SMITHY_HTTP_TRANSPORT).build())
                                                .functionArguments(ListUtils.of(
                                                        SymbolUtils.createValueSymbolBuilder("2").build(),
                                                        SymbolUtils.createValueSymbolBuilder("0").build()))
                                                .build())
                                        .build())
                        .collect(Collectors.toSet()))
                .forEach(plugins::addAll);

        return plugins;
    }
}
