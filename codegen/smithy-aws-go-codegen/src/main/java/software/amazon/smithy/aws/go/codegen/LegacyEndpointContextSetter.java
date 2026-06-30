package software.amazon.smithy.aws.go.codegen;

import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;

public class LegacyEndpointContextSetter implements GoIntegration {

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first. Needs to execute after Rules Engine endpoint
     * resolution middleware insertion.
     *
     * @return Returns the sort order, defaults to 127.
     */
    @Override
    public byte getOrder() {
            return -128;
    }

    @Override
    public void renderPreEndpointResolutionHook(GoSettings settings, GoWriter writer, Model model) {
        writer.write("""
                if $T(ctx) {
                    return next.HandleFinalize(ctx, in)
                }
                """, AwsGoDependency.AWS_MIDDLEWARE.func("GetRequiresLegacyEndpoints"));
    }

}
