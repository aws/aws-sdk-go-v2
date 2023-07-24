package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.trait.NoSerializeTrait;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.traits.HttpLabelTrait;
import software.amazon.smithy.model.transform.ModelTransformer;

public class S3HttpLabelBucketFilterIntegration implements GoIntegration {

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        if (!S3ModelUtils.isServiceS3(model, settings.getService(model))) {
            return model;
        }

        return ModelTransformer.create().mapTraits(model, (shape, trait) -> {
            if (!(trait instanceof HttpLabelTrait)) return trait;

            boolean isBucket = shape.asMemberShape()
                    .map(s -> s.getMemberName().equals("Bucket"))
                    .orElse(false);
            return isBucket ? new NoSerializeTrait() : trait;
        });
    }
}
