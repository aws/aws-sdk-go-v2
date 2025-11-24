package software.amazon.smithy.aws.go.codegen.customization.service.s3;

import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.ErrorTrait;
import software.amazon.smithy.model.transform.ModelTransformer;

public class S3AddErrorMessage implements GoIntegration {
    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        if (!S3ModelUtils.isServiceS3(model, settings.getService(model))) {
            return model;
        }

        return ModelTransformer.create()
                .mapShapes(model, (shape) ->
                        shape.hasTrait(ErrorTrait.class)
                                ? withMessage(shape)
                                : shape
                );
    }

    private static Shape withMessage(Shape shape) {
        return ((StructureShape) shape).toBuilder()
                .addMember("Message", ShapeId.from("com.amazonaws.s3#Message"))
                .build();
    }
}
