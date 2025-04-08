package software.amazon.smithy.aws.go.codegen.customization.service.bedrockruntime;

import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ShapeId;

public class RemoveOperations implements GoIntegration {
    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        var sdkId = settings.getService(model).expectTrait(ServiceTrait.class).getSdkId();
        if (!sdkId.equalsIgnoreCase("Bedrock Runtime")) {
            return model;
        }

        return model.toBuilder()
                // remove InvokeModelWithBidirectionalStream because the SDK does not currently support event stream
                // APIs that do not immediately send the initial HTTP response
                .removeShape(ShapeId.from("com.amazonaws.bedrockruntime#InvokeModelWithBidirectionalStream"))
                .removeShape(ShapeId.from("com.amazonaws.bedrockruntime#InvokeModelWithBidirectionalStreamRequest"))
                .removeShape(ShapeId.from("com.amazonaws.bedrockruntime#InvokeModelWithBidirectionalStreamResponse"))
                .build();
    }
}
