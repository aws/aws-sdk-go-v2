package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.DefaultTrait;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.utils.MapUtils;

import java.util.Arrays;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

public class RemoveDefaults implements GoIntegration {
    private static final Map<ShapeId, Set<ShapeId>> toRemove = MapUtils.ofEntries(
        serviceToShapeIds("com.amazonaws.s3control#AWSS3ControlServiceV20180820",
                "com.amazonaws.s3control#PublicAccessBlockConfiguration$BlockPublicAcls",
                "com.amazonaws.s3control#PublicAccessBlockConfiguration$IgnorePublicAcls",
                "com.amazonaws.s3control#PublicAccessBlockConfiguration$BlockPublicPolicy",
                "com.amazonaws.s3control#PublicAccessBlockConfiguration$RestrictPublicBuckets"),
        serviceToShapeIds("com.amazonaws.evidently#Evidently",
                "com.amazonaws.evidently#ResultsPeriod"),
        serviceToShapeIds("com.amazonaws.amplifyuibuilder#AmplifyUIBuilder",
                "smithy.go.synthetic#ListPlaceIndexesInput$MaxResults",
                "smithy.go.synthetic#SearchPlaceIndexForSuggestionsInput$MaxResults",
                "com.amazonaws.location#PlaceIndexSearchResultLimit"),
        serviceToShapeIds("com.amazonaws.paymentcryptographydata#PaymentCryptographyDataPlane",
                "com.amazonaws.paymentcryptographydata#IntegerRangeBetween4And12"),
        serviceToShapeIds("com.amazonaws.emrserverless#AwsToledoWebService",
                "com.amazonaws.emrserverless#WorkerCounts"),
        serviceToShapeIds("com.amazonaws.imagebuilder#imagebuilder",
                    // https://github.com/aws/aws-sdk-go-v2/issues/2734
                    // V1479153907
                    "com.amazonaws.imagebuilder#LaunchTemplateConfiguration$setDefaultVersion"));

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
        Set<ShapeId> removedRootDefaults = new HashSet<>();
        Model removedRootDefaultsModel = ModelTransformer.create().mapShapes(model, (shape) -> {
            if (shouldRemoveRootDefault(shape, fromShapes)) {
                removedRootDefaults.add(shape.getId());
                return withoutDefault(shape);
            } else {
                return shape;
            }
        });
        return ModelTransformer.create().mapShapes(removedRootDefaultsModel, (shape) -> {
            if (shouldRemoveMemberDefault(shape, removedRootDefaults, fromShapes)) {
                return withoutDefault(shape);
            } else {
                return shape;
            }
        });
    }

    private boolean shouldRemoveRootDefault(Shape shape, Set<ShapeId> removeDefaultsFrom) {
        return !shape.isMemberShape()
                && removeDefaultsFrom.contains(shape.getId())
                && shape.hasTrait(DefaultTrait.class);
    }

    private boolean shouldRemoveMemberDefault(
            Shape shape,
            Set<ShapeId> removedRootDefaults,
            Set<ShapeId> removeDefaultsFrom
    ) {
        if (!shape.isMemberShape()) {
            return false;
        }
        MemberShape member = shape.asMemberShape().get();
        return (removedRootDefaults.contains(member.getTarget()) || removeDefaultsFrom.contains(member.getId()))
                && member.hasTrait(DefaultTrait.class);
    }

    private Shape withoutDefault(Shape shape) {
        return Shape.shapeToBuilder(shape)
                .removeTrait(DefaultTrait.ID)
                .build();
    }

    private static Map.Entry<ShapeId, Set<ShapeId>> serviceToShapeIds(String serviceId, String... shapeIds) {
        return Map.entry(
                ShapeId.from(serviceId),
                Arrays.stream(shapeIds).map(ShapeId::from).collect(Collectors.toSet()));
    }
}
