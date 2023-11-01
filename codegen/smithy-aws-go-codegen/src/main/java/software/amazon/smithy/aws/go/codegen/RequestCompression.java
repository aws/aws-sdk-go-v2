package software.amazon.smithy.aws.go.codegen;


import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ToShapeId;
import software.amazon.smithy.model.traits.RequestCompressionTrait;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

import static software.amazon.smithy.aws.go.codegen.AwsGoDependency.AWS_MIDDLEWARE;

public final class RequestCompression implements GoIntegration {
    private static final String ADD_Request_Compression = "addRequestCompression";

    private static final String ADD_Request_Compression_Internal = "AddRequestCompression";

    private static final String Disable_Request_Compression = "DisableRequestCompression";

    private static final String Request_Min_Compression_Size_Bytes = "RequestMinCompressSizeBytes";

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceShape service = settings.getService(model);
        if (!isRequestCompressionService(model, service)) {
            return;
        }

        goDelegator.useShapeWriter(service, this::writeMiddlewareHelper);
    }


    public static boolean isRequestCompressionService(Model model, ServiceShape service) {
        TopDownIndex topDownIndex = TopDownIndex.of(model);
        for (ToShapeId operation : topDownIndex.getContainedOperations(service)) {
            OperationShape operationShape = model.expectShape(operation.toShapeId(), OperationShape.class);
            if (operationShape.hasTrait(RequestCompressionTrait.class)) {
                return true;
            }
        }
        return false;
}

    private void writeMiddlewareHelper(GoWriter writer) {
        writer.openBlock("func $L(stack *middleware.Stack, options Options) error {", "}",
        ADD_Request_Compression, () -> {
            writer.write("return $T(stack, options.DisableRequestCompression, options.RequestMinCompressSizeBytes)",
                    SymbolUtils.createValueSymbolBuilder(ADD_Request_Compression_Internal,
                            AwsGoDependency.AWS_MIDDLEWARE).build()
            );
        });
        writer.insertTrailingNewline();
    }


    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> operation.hasTrait(RequestCompressionTrait.class))
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ADD_Request_Compression).build())
                                .useClientOptions()
                                .build())
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(RequestCompression::isRequestCompressionService)
                        .configFields(ListUtils.of(
                                ConfigField.builder()
                                        .name(Disable_Request_Compression)
                                        .type(SymbolUtils.createValueSymbolBuilder("bool")
                                                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                                .build())
                                        .documentation("Determine if request compression is allowed, default to false")
                                        .build(),
                                ConfigField.builder()
                                        .name(Request_Min_Compression_Size_Bytes)
                                        .type(SymbolUtils.createValueSymbolBuilder("int64")
                                                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                                .build())
                                        .documentation("Inclusive threshold request body size to trigger compression, " +
                                         "default to 10240 and must be within 0 and 10485760 bytes inclusively")
                                        .build()
                        ))
                        .build()
        );
    }
}