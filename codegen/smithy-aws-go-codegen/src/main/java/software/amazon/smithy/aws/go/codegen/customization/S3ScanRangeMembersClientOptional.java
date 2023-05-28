package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.ClientOptionalTrait;

public class S3ScanRangeMembersClientOptional implements GoIntegration {

    public static final ShapeId S3_SERVICE_SHAPE = ShapeId.from("com.amazonaws.s3#AmazonS3");
    public static final ShapeId SCAN_RANGE_SHAPE = ShapeId.from("com.amazonaws.s3#ScanRange");
    public static final String[] MEMBERS_TO_MAKE_OPTIONAL = { "Start", "End" };

    /**
     * /**
     * Updates the API model to customize ScanRange member shapes to be
     * clientOptional.
     *
     * @param model    API model
     * @param settings Go codegen settings
     * @return updated API model
     */
    @Override
    public Model preprocessModel(Model model, GoSettings settings) {

        if (!S3_SERVICE_SHAPE.equals(settings.getService())) {
            return model;
        }

        Model.Builder builder = model.toBuilder();

        StructureShape structureShape = model.expectShape(SCAN_RANGE_SHAPE,
                StructureShape.class);
        StructureShape.Builder structureShapeBuilder = structureShape.toBuilder();

        for (String memberName : MEMBERS_TO_MAKE_OPTIONAL) {
            structureShape.getMember(memberName).ifPresent(member -> {
                structureShapeBuilder.addMember(member.toBuilder().addTrait(new ClientOptionalTrait()).build());
            });
        }

        return builder.addShape(structureShapeBuilder.build()).build();
    }
}