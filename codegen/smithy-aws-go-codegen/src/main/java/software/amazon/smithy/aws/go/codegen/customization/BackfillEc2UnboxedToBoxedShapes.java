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
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.ClientOptionalTrait;
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
                updates.add(member.toBuilder().addTrait(new ClientOptionalTrait()).build());
            }
        }
        return ModelTransformer.create().replaceShapes(model, updates);
    }
}
