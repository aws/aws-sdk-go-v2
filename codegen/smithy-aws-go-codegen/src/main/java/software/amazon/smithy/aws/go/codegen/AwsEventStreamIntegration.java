package software.amazon.smithy.aws.go.codegen;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import software.amazon.smithy.go.codegen.GoEventStreamIndex;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ShapeId;

public class AwsEventStreamIntegration implements GoIntegration {
    private final Map<ShapeId, Collection<OperationShape>> serviceOperationMap = new HashMap<>();

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

        Collection<OperationShape> operationShapes = new HashSet<>();

        goEventStreamIndex.getInputEventStreams(service).ifPresent(shapeIdSetMap ->
                shapeIdSetMap.values().forEach(eventStreamInfos ->
                        eventStreamInfos.forEach(info -> operationShapes.add(info.getOperation()))));

        goEventStreamIndex.getOutputEventStreams(service).ifPresent(shapeIdSetMap ->
                shapeIdSetMap.values().forEach(eventStreamInfos ->
                        eventStreamInfos.forEach(info -> operationShapes.add(info.getOperation()))));

        if (!operationShapes.isEmpty()) {
            serviceOperationMap.put(service, operationShapes);
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

        return plugins;
    }
}
