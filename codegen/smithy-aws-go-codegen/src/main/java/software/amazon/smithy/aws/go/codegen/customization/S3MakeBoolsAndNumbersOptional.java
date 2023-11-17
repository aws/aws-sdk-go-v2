package software.amazon.smithy.aws.go.codegen.customization;

import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import java.util.logging.Logger;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.NumberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.DefaultTrait;
import software.amazon.smithy.utils.SetUtils;

public class S3MakeBoolsAndNumbersOptional implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(S3MakeBoolsAndNumbersOptional.class.getName());

    /**
     * Map of service shape to Set of operation shapes that need to have this
     * presigned url auto fill customization.
     */
    public static final Set<ShapeId> SERVICE_SET = SetUtils.of(
            ShapeId.from("com.amazonaws.s3#AmazonS3")
    );

    /**
     * /**
     * Updates the API model to customize all structured members to be nullable.
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

        List<Shape> updates = new ArrayList<>();
        for (StructureShape struct : model.getStructureShapes()) {
            for (MemberShape member : struct.getAllMembers().values()) {
                Shape target = model.expectShape(member.getTarget());
                if (isNumberShape(target) || target.isBooleanShape()) {
                    updates.add(member.toBuilder().removeTrait(DefaultTrait.ID).build());
                    updates.add(Shape.shapeToBuilder(target).removeTrait(DefaultTrait.ID).build());
                }
            }
        }
        return ModelTransformer.create().replaceShapes(model, updates);
    }

    private static boolean isNumberShape(Shape shape) {
        return shape.isByteShape()
                || shape.isShortShape()
                || shape.isIntegerShape()
                || shape.isLongShape()
                || shape.isFloatShape()
                || shape.isDoubleShape()
                || shape.isBigDecimalShape()
                || shape.isBigIntegerShape();
    }
}