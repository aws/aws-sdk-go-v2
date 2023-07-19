package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;

public class DisableEndpointHostPrefix implements GoIntegration {
    @Override
    public void renderPostEndpointResolutionHook(GoSettings settings, GoWriter writer, Model model) {
        writer.write(
            """
            ctx = $T(ctx, true)
            """,
            SymbolUtils.createPointableSymbolBuilder("DisableEndpointHostPrefix", SmithyGoDependency.SMITHY_HTTP_TRANSPORT).build()
        );
    }
    
}
