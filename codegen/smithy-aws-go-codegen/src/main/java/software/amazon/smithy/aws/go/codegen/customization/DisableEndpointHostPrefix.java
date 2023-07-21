package software.amazon.smithy.aws.go.codegen.customization;

import java.util.Optional;


import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.pattern.SmithyPattern.Segment;
import software.amazon.smithy.model.traits.EndpointTrait;
import software.amazon.smithy.model.shapes.OperationShape;


public class DisableEndpointHostPrefix implements GoIntegration {
    @Override
    public void renderPostEndpointResolutionHook(
        GoSettings settings, GoWriter writer, Model model, Optional<OperationShape> operation
    ) {
        if (!S3ModelUtils.isServiceS3Control(model, settings.getService(model))) {
            return;
        }

        if (operation.isPresent() && operation.get().getTrait(EndpointTrait.class).isPresent()) {
            EndpointTrait endpointTrait = operation.get().getTrait(EndpointTrait.class).get();
            boolean written = false;
            for (Segment segment : endpointTrait.getHostPrefix().getLabels()) {
                if (segment.isLabel() && segment.getContent().equals("AccountId") && !written) {
                    writer.write(
                        """
                        ctx = $T(ctx, true)
                        """,
                        SymbolUtils.createPointableSymbolBuilder(
                            "DisableEndpointHostPrefix", SmithyGoDependency.SMITHY_HTTP_TRANSPORT).build()
                    );
                    // we only want the content above written ONCE per operation.
                    written = true;
                }
            }
        }
    }
    
}
