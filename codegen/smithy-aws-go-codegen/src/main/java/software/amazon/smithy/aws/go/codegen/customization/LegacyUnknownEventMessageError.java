package software.amazon.smithy.aws.go.codegen.customization;

import java.util.Set;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.codegen.core.SymbolProvider;

/**
 * Generates the per-service UnknownEventMessageError type for services that had event streams before the
 * useExperimentalSerde migration. This preserves backwards compatibility for those services without generating the type
 * for any new services that add event streams in the future.
 */
public class LegacyUnknownEventMessageError implements GoIntegration {
    // Services that have event streams as of the experimental serde migration. This list is intentionally static -- new
    // services that add event streams will NOT get this type.
    private static final Set<String> LEGACY_SERVICES = Set.of(
            "com.amazonaws.kinesis",
            "com.amazonaws.s3",
            "com.amazonaws.lambda",
            "com.amazonaws.polly",
            "com.amazonaws.iotsitewise",
            "com.amazonaws.sagemakerruntime",
            "com.amazonaws.bedrockruntime",
            "com.amazonaws.bedrockagentruntime",
            "com.amazonaws.bedrockagentcore",
            "com.amazonaws.transcribestreaming",
            "com.amazonaws.lexruntimev2",
            "com.amazonaws.qbusiness",
            "com.amazonaws.cloudwatchlogs",
            "com.amazonaws.devopsagent",
            "com.amazonaws.connecthealth"
    );

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        if (!settings.useExperimentalSerde()) {
            return;
        }

        var serviceNamespace = settings.getService().getNamespace();
        if (!LEGACY_SERVICES.contains(serviceNamespace)) {
            return;
        }

        goDelegator.useShapeWriter(settings.getService(model), this::generateType);
    }

    private void generateType(GoWriter writer) {
        var message = SymbolUtils.createPointableSymbolBuilder("Message",
                AwsGoDependency.SERVICE_INTERNAL_EVENTSTREAM).build();

        writer.write("""
                // UnknownEventMessageError provides an error when a message is received from the stream,
                // but the reader is unable to determine what kind of message it is.
                type UnknownEventMessageError struct {
                    Type    string
                    Message $T
                }

                // Error retruns the error message string.
                func (e *UnknownEventMessageError) Error() string {
                    return "unknown event stream message type, " + e.Type
                }
                """, message);
    }
}
