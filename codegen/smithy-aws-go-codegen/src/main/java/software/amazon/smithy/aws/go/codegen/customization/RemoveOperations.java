package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import java.util.Locale;
import java.util.Map;

import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ShapeId;

public class RemoveOperations implements GoIntegration {

    private Map<String, List<ShapeId>> SHAPES_TO_REMOVE = Map.of(
        "bedrock runtime", List.of(
            ShapeId.from("com.amazonaws.bedrockruntime#InvokeModelWithBidirectionalStream"),
            ShapeId.from("com.amazonaws.bedrockruntime#InvokeModelWithBidirectionalStreamRequest"),
            ShapeId.from("com.amazonaws.bedrockruntime#InvokeModelWithBidirectionalStreamResponse")
        ),
        "sagemaker runtime http2", List.of(
            ShapeId.from("com.amazonaws.sagemakerruntimehttp2#InvokeEndpointWithBidirectionalStream"),
            ShapeId.from("com.amazonaws.sagemakerruntimehttp2#InvokeEndpointWithBidirectionalStreamInput"),
            ShapeId.from("com.amazonaws.sagemakerruntimehttp2#InvokeEndpointWithBidirectionalStreamOutput")
        )
    );

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        String sdkId = settings.getService(model).expectTrait(ServiceTrait.class).getSdkId();
        List<ShapeId> serviceShapes = SHAPES_TO_REMOVE.get(sdkId.toLowerCase(Locale.ROOT));
        if (serviceShapes != null) {
            for (ShapeId shapeId : serviceShapes) {
                model = model.toBuilder().removeShape(shapeId).build();
            }
        }
        return model;
    }
}
