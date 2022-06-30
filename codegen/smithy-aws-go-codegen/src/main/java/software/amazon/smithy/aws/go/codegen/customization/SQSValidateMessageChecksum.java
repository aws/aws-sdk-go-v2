package software.amazon.smithy.aws.go.codegen.customization;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.logging.Logger;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoCodegenPlugin;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.utils.SetUtils;

public class SQSValidateMessageChecksum implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(SQSValidateMessageChecksum.class.getName());

    /**
     * Map of service shape to Set of operation shapes that need to have this
     * customization.
     */
    public static final Map<ShapeId, Set<ShapeId>> SERVICE_TO_OPERATION_MAP = MapUtils.of(
            ShapeId.from("com.amazonaws.sqs#AmazonSQS"), SetUtils.of(
                    ShapeId.from("com.amazonaws.sqs#SendMessage"),
                    ShapeId.from("com.amazonaws.sqs#SendMessageBatch"),
                    ShapeId.from("com.amazonaws.sqs#ReceiveMessage")
            )
    );
    static final String DISABLE_MESSAGE_CHECKSUM_VALIDATION_OPTION_NAME = "DisableMessageChecksumValidation";

    private final List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();

    /**
     * Builds the set of runtime plugs used by the customization.
     *
     * @param settings codegen settings
     * @param model    api model
     */
    @Override
    public void processFinalizedModel(GoSettings settings, Model model) {
        ShapeId serviceId = settings.getService();
        if (!SERVICE_TO_OPERATION_MAP.containsKey(serviceId)) {
            return;
        }

        ServiceShape service = settings.getService(model);

        // Add option to disable message checksum validation
        runtimeClientPlugins.add(RuntimeClientPlugin.builder()
                .servicePredicate((m, s) -> s.equals(service))
                .addConfigField(ConfigField.builder()
                        .name(DISABLE_MESSAGE_CHECKSUM_VALIDATION_OPTION_NAME)
                        .type(SymbolUtils.createValueSymbolBuilder("bool")
                                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true).build())
                        .documentation("Allows you to disable the client's validation of "
                                + "response message checksums. Enabled by default. "
                                + "Used by SendMessage, SendMessageBatch, and ReceiveMessage.")
                        .build())
                .build());

        for (ShapeId operationId : SERVICE_TO_OPERATION_MAP.get(serviceId)) {
            final OperationShape operation = model.expectShape(operationId, OperationShape.class);

            // Create a symbol provider because one is not available in this call.
            SymbolProvider symbolProvider = GoCodegenPlugin.createSymbolProvider(model, settings);

            String helperFuncName = addMiddlewareFuncName(symbolProvider.toSymbol(operation).getName());

            runtimeClientPlugins.add(RuntimeClientPlugin.builder()
                    .servicePredicate((m, s) -> s.equals(service))
                    .operationPredicate((m, s, o) -> o.equals(operation))
                    .registerMiddleware(MiddlewareRegistrar.builder()
                            .resolvedFunction(SymbolUtils.createValueSymbolBuilder(helperFuncName)
                                    .build())
                            .useClientOptions()
                            .build())
                    .build());
        }
    }

    String addMiddlewareFuncName(String operationName) {
        return "addValidate" + operationName + "Checksum";
    }

    /**
     * Returns the list of runtime client plugins added by this customization
     *
     * @return runtime client plugins
     */
    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return runtimeClientPlugins;
    }
}
