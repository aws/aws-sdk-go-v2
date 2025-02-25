package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.traits.DeprecatedTrait;
import software.amazon.smithy.model.transform.ModelTransformer;

import java.util.Set;

/**
 * Deprecates an entire service by applying @deprecated to every shape. This will reflect in IDEs and public package
 * documentation.
 */
public class DeprecateService implements GoIntegration {
    private static final String DEPRECATION_MESSAGE =
            "AWS has deprecated this service. It is no longer available for use.";
    private static final Set<String> DEPRECATED = Set.of(
            "com.amazonaws.nimble#nimble",
            "com.amazonaws.iot1clickdevicesservice#IoT1ClickDevicesService",
            "com.amazonaws.iot1clickprojects#AWSIoT1ClickProjects",
            "com.amazonaws.elasticinference#EC2MatterhornCfSoothsayerApiGatewayLambda"
    );

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        var service = settings.getService(model);
        if (!DEPRECATED.contains(service.getId().toString())) {
            return model;
        }

        return ModelTransformer.create().mapShapes(model, shape -> {
            if (shape.isMemberShape()) {
                return shape;
            }

            var deprecated = DeprecatedTrait.builder()
                    .message(DEPRECATION_MESSAGE)
                    .build();
            return Shape.shapeToBuilder(shape)
                    .addTrait(deprecated)
                    .build();
        });
    }
}
