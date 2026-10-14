package software.amazon.smithy.aws.go.codegen.customization.service.sts;

import java.util.List;
import java.util.Map;

import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

/**
 * Adds IDPCommunicationError as a retryable error code for STS.
 * This error is not modeled with the retryable trait but must be retried as a
 * transient network error for STS clients only.
 */
public class STSRetryableErrors implements GoIntegration {

    private static final String FINALIZE_STS_RETRYER = "finalizeSTSRetryableErrors";

    private static boolean isSTS(Model model, ServiceShape service) {
        return service.getTrait(ServiceTrait.class)
                .map(trait -> "STS".equals(trait.getSdkId()))
                .orElse(false);
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator delegator
    ) {
        ServiceShape service = settings.getService(model);
        if (!isSTS(model, service)) {
            return;
        }

        delegator.useShapeWriter(service, this::writeResolver);
    }

    private void writeResolver(GoWriter writer) {
        writer.write(goTemplate("""
                $retry:D
                func finalizeSTSRetryableErrors(o *Options) {
                    o.Retryer = retry.AddWithErrorCodes(o.Retryer, "IDPCommunicationError")
                }
                """, Map.of(
                "retry", AwsGoDependency.AWS_RETRY
        )));
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(STSRetryableErrors::isSTS)
                        .addConfigFieldResolver(ConfigFieldResolver.builder()
                                .location(ConfigFieldResolver.Location.CLIENT)
                                .target(ConfigFieldResolver.Target.FINALIZATION)
                                .resolver(SymbolUtils.createValueSymbolBuilder(
                                        FINALIZE_STS_RETRYER).build())
                                .build())
                        .build()
        );
    }
}
