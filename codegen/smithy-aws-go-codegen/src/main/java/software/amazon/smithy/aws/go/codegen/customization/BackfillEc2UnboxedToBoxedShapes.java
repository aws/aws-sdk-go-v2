package software.amazon.smithy.aws.go.codegen.customization;

import java.util.Set;
import java.util.logging.Logger;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.go.codegen.CodegenUtils;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.ShapeType;
import software.amazon.smithy.model.traits.BoxTrait;
import software.amazon.smithy.utils.SetUtils;

public class BackfillEc2UnboxedToBoxedShapes implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(BackfillEc2UnboxedToBoxedShapes.class.getName());

    /**
     * Map of service shape to Set of operation shapes that need to have this
     * presigned url auto fill customization.
     */
    public static final Set<ShapeId> SERVICE_SET = SetUtils.of(
            ShapeId.from("com.amazonaws.ec2#AmazonEC2")
    );

    /**
     * /**
     * Updates the API model to customize all number and boolean unboxed shapes to be boxed.
     *
     * @param model    API model
     * @param settings Go codegen settings
     * @return updated API model
     */
    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        ShapeId serviceId = settings.getService();
        if (!SERVICE_SET.contains(serviceId)) {
            return model;
        }

        Model.Builder builder = model.toBuilder();

        for (Shape shape : model.toSet()) {
            // Only consider number and boolean shapes that are unboxed
            if (shape.isMemberShape()) {
                continue;
            }
            if (!(CodegenUtils.isNumber(shape) || shape.getType() == ShapeType.BOOLEAN)) {
                continue;
            }
            if (shape.hasTrait(BoxTrait.class)) {
                continue;
            }

            switch (shape.getType()) {
                case BYTE:
                    shape = shape.asByteShape().get().toBuilder()
                            .addTrait(new BoxTrait())
                            .build();
                    break;
                case SHORT:
                    shape = shape.asShortShape().get().toBuilder()
                            .addTrait(new BoxTrait())
                            .build();
                    break;
                case INTEGER:
                    shape = shape.asIntegerShape().get().toBuilder()
                            .addTrait(new BoxTrait())
                            .build();
                    break;
                case LONG:
                    shape = shape.asLongShape().get().toBuilder()
                            .addTrait(new BoxTrait())
                            .build();
                    break;
                case FLOAT:
                    shape = shape.asFloatShape().get().toBuilder()
                            .addTrait(new BoxTrait())
                            .build();
                    break;
                case DOUBLE:
                    shape = shape.asDoubleShape().get().toBuilder()
                            .addTrait(new BoxTrait())
                            .build();
                    break;
                case BOOLEAN:
                    shape = shape.asBooleanShape().get().toBuilder()
                            .addTrait(new BoxTrait())
                            .build();
                    break;
                default:
                    throw new CodegenException("unexpected shape type for " + shape.getId() + ", " + shape.getType());
            }

            builder.addShape(shape);
        }

        return builder.build();
    }
}
