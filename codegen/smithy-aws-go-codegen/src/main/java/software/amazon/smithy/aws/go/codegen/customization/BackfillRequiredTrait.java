package software.amazon.smithy.aws.go.codegen.customization;

import static java.util.stream.Collectors.toSet;

import java.util.Arrays;
import java.util.Map;
import java.util.Set;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.RequiredTrait;
import software.amazon.smithy.model.transform.ModelTransformer;

/**
 * Adds the @required trait to inputs that have @default, but lack of serialization of their zero value causes issues.
 *
 * If a shape is listed here there should generally be an upstreaming effort with the service team to fix. Link issues
 * in comments (or internal ticket IDs) where available.
 */
public class BackfillRequiredTrait implements GoIntegration {
    private static final Map<ShapeId, Set<ShapeId>> toBackfill = Map.ofEntries(
            serviceToShapeIds("com.amazonaws.sqs#AmazonSQS",
                    // https://github.com/aws/aws-sdk/issues/527
                    "com.amazonaws.sqs#ChangeMessageVisibilityBatchRequestEntry$VisibilityTimeout")
    );

    private boolean mustPreprocess(ShapeId service) {
        return toBackfill.containsKey(service);
    }

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        var serviceId = settings.getService();
        return mustPreprocess(serviceId)
                ? backfillRequired(model, toBackfill.get(serviceId))
                : model;
    }

    private Model backfillRequired(Model model, Set<ShapeId> shapes) {
        return ModelTransformer.create()
                .mapShapes(model, (shape) -> shapes.contains(shape.getId()) ? withRequired(shape) : shape);
    }

    private Shape withRequired(Shape shape) {
        return Shape.shapeToBuilder(shape)
                .addTrait(new RequiredTrait())
                .build();
    }

    private static Map.Entry<ShapeId, Set<ShapeId>> serviceToShapeIds(String serviceId, String... shapeIds) {
        return Map.entry(
                ShapeId.from(serviceId),
                Arrays.stream(shapeIds).map(ShapeId::from).collect(toSet()));
    }
}
