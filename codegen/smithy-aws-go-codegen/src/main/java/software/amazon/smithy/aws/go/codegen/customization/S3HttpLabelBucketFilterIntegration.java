package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.trait.NoSerializeTrait;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.traits.HttpLabelTrait;
import software.amazon.smithy.model.transform.ModelTransformer;

public class S3HttpLabelBucketFilterIntegration implements GoIntegration {

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        if (!S3ModelUtils.isServiceS3(model, settings.getService(model))) {
            return model;
        }

        return ModelTransformer.create().mapShapes(model, (shape) -> {
            if (!shape.hasTrait(HttpLabelTrait.class)) return shape;

            boolean isBucket = shape.asMemberShape()
                    .map(s -> s.getMemberName().equals("Bucket"))
                    .orElse(false);
            if (isBucket) {
                shape = Shape.shapeToBuilder(shape).addTrait(new NoSerializeTrait()).build();
            }

            return shape;
        });
    }
}
