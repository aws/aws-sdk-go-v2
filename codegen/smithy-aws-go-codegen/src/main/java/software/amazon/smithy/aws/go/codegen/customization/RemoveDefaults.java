package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.AbstractShapeBuilder;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.DefaultTrait;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.utils.SetUtils;
import software.amazon.smithy.utils.ToSmithyBuilder;

import java.util.Map;
import java.util.Set;

public class RemoveDefaults implements GoIntegration {
    private static final Map<ShapeId, Set<ShapeId>> toRemove = MapUtils.of(
            ShapeId.from("com.amazonaws.s3control#AWSS3ControlServiceV20180820"), SetUtils.of(
                    ShapeId.from("com.amazonaws.s3control#PublicAccessBlockConfiguration$BlockPublicAcls"),
                    ShapeId.from("com.amazonaws.s3control#PublicAccessBlockConfiguration$IgnorePublicAcls"),
                    ShapeId.from("com.amazonaws.s3control#PublicAccessBlockConfiguration$BlockPublicPolicy"),
                    ShapeId.from("com.amazonaws.s3control#PublicAccessBlockConfiguration$RestrictPublicBuckets")
            )
    );

    private boolean mustPreprocess(ShapeId service) {
        return toRemove.containsKey(service);
    }

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        var serviceId = settings.getService();
        return mustPreprocess(serviceId)
                ? removeDefaults(model, toRemove.get(serviceId))
                : model;
    }

    private Model removeDefaults(Model model, Set<ShapeId> fromShapes) {
        return ModelTransformer.create().mapShapes(model, it ->
                fromShapes.contains(it.getId())
                        ? withoutDefault(it)
                        : it
        );
    }

    private Shape withoutDefault(Shape shape) {
        var builder = ((ToSmithyBuilder<?>) shape).toBuilder();
        return ((AbstractShapeBuilder<?, ?>) builder).removeTrait(DefaultTrait.ID).build();
    }
}