package software.amazon.smithy.aws.go.codegen.customization;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import software.amazon.smithy.aws.go.codegen.AddAwsConfigFields;
import software.amazon.smithy.aws.go.codegen.AwsEndpointGenerator;
import software.amazon.smithy.aws.go.codegen.AwsSignatureVersion4;
import software.amazon.smithy.aws.go.codegen.AwsSignatureVersion4aUtils;
import software.amazon.smithy.aws.go.codegen.EndpointGenerator;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.OperationIndex;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.ShapeType;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

/**
 * This integration configures the EventBridge client for Signature Version 4a
 */
public class EventBridgeMultiRegionEndpoint implements GoIntegration {
    private static Map<ShapeId, String> SUPPORTED_OPERATIONS = MapUtils.of(
            ShapeId.from("com.amazonaws.eventbridge#PutEvents"), "EndpointId"
    );

    /**
     * Return true if service is .
     *
     * @param model   is the generation model.
     * @param service is the service shape being audited.
     */
    private static boolean isEventBridgeService(Model model, ServiceShape service) {
        var serviceId = service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("EventBridge");
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        var resolver = SymbolUtils.createValueSymbolBuilder(
                AwsSignatureVersion4aUtils.RESOLVE_CREDENTIAL_PROVIDER).build();

        List<RuntimeClientPlugin> plugins = new ArrayList<>();

        for (var operationId : SUPPORTED_OPERATIONS.keySet()) {
            plugins.add(RuntimeClientPlugin.builder()
                    .operationPredicate((m, s, o) -> {
                        if (!isEventBridgeService(m, s)) {
                            return false;
                        }
                        return o.toShapeId().equals(operationId);
                    })
                    .registerMiddleware(MiddlewareRegistrar.builder()
                            .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                            getEndpointMiddlewareHelperName(operationId.getName()))
                                    .build())
                            .useClientOptions()
                            .build())
                    .build());
        }

        plugins.addAll(ListUtils.of(RuntimeClientPlugin.builder()
                        .addConfigFieldResolver(ConfigFieldResolver.builder()
                                .location(ConfigFieldResolver.Location.CLIENT)
                                .target(ConfigFieldResolver.Target.FINALIZATION)
                                .resolver(resolver)
                                .build())
                        .addConfigFieldResolver(ConfigFieldResolver.builder()
                                .location(ConfigFieldResolver.Location.OPERATION)
                                .target(ConfigFieldResolver.Target.FINALIZATION)
                                .resolver(resolver)
                                .build())
                        .servicePredicate((model, serviceShape) -> {
                            if (!EventBridgeMultiRegionEndpoint.isEventBridgeService(model, serviceShape)) {
                                return false;
                            }
                            return AwsSignatureVersion4.isSupportedAuthentication(model, serviceShape);
                        })
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(EventBridgeMultiRegionEndpoint::isEventBridgeService)
                        .addConfigField(ConfigField.builder()
                                .name(AwsSignatureVersion4aUtils.V4A_SIGNER_INTERFACE_NAME)
                                .type(SymbolUtils.createValueSymbolBuilder(
                                                AwsSignatureVersion4aUtils.V4A_SIGNER_INTERFACE_NAME)
                                        .build())
                                .documentation("Signature Version 4a (SigV4a) Signer")
                                .build())
                        .build(),
                // Add HTTPSigner middleware to operation stack
                RuntimeClientPlugin.builder()
                        .servicePredicate(EventBridgeMultiRegionEndpoint::isEventBridgeService)
                        .operationPredicate((model, service, operation) -> {
                            if (!isEventBridgeService(model, service)) {
                                return false;
                            }
                            return SUPPORTED_OPERATIONS.containsKey(operation.toShapeId());
                        })
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        AwsSignatureVersion4aUtils.REGISTER_MIDDLEWARE_FUNCTION).build())
                                .useClientOptions()
                                .build())
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(EventBridgeMultiRegionEndpoint::isEventBridgeService)
                        .addConfigFieldResolver(
                                ConfigFieldResolver.builder()
                                        .location(ConfigFieldResolver.Location.CLIENT)
                                        .target(ConfigFieldResolver.Target.INITIALIZATION)
                                        .resolver(SymbolUtils.createValueSymbolBuilder(
                                                AwsSignatureVersion4aUtils.SIGNER_RESOLVER).build())
                                        .build())
                        .build()
        ));

        return plugins;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        if (!isEventBridgeService(model, model.expectShape(settings.getService(), ServiceShape.class))) {
            return;
        }

        if (SUPPORTED_OPERATIONS.size() == 0) {
            return;
        }

        var serviceShape = settings.getService(model);
        goDelegator.useShapeWriter(serviceShape, writer -> {
            writeCredentialProviderResolver(writer);
            writeSigningMiddlewareRegister(model, writer, serviceShape);
            writerSignerInterface(writer);
            writerSignerConfigFieldResolver(writer, serviceShape);
            writeNewV4ASignerFunc(writer, serviceShape);
        });

        for (var operationShape : TopDownIndex.of(model).getContainedOperations(serviceShape)) {
            if (!SUPPORTED_OPERATIONS.containsKey(operationShape.toShapeId())) {
                continue;
            }
            goDelegator.useShapeWriter(operationShape, writer -> {
                writeEndpointIdAccessorHelper(writer, model, symbolProvider, operationShape);
                writeEndpointMiddlewareHelper(writer, model, symbolProvider, operationShape);
            });
        }
    }

    private void writeCredentialProviderResolver(GoWriter writer) {
        AwsSignatureVersion4aUtils.writeCredentialProviderResolver(writer);
    }

    private void writerSignerInterface(GoWriter writer) {
        AwsSignatureVersion4aUtils.writerSignerInterface(writer);
    }

    private void writerSignerConfigFieldResolver(GoWriter writer, ServiceShape serviceShape) {
        AwsSignatureVersion4aUtils.writerConfigFieldResolver(writer, serviceShape);
    }

    private void writeNewV4ASignerFunc(GoWriter writer, ServiceShape serviceShape) {
        AwsSignatureVersion4aUtils.writeNewV4ASignerFunc(writer, serviceShape);
    }

    private void writeSigningMiddlewareRegister(Model model, GoWriter writer, ServiceShape serviceShape) {
        AwsSignatureVersion4aUtils.writeMiddlewareRegister(model, writer, serviceShape,
                AwsCustomGoDependency.EVENTBRIDGE_CUSTOMIZATION);
    }

    // retrieves function name for get bucket accessor function
    private String getEndpointIdAccessorName(String operationName, String memberName) {
        return String.format("get%s%s", operationName, memberName);
    }

    private String getEndpointMiddlewareHelperName(String operationName) {
        return String.format("add%s%s", operationName, "UpdateEndpoint");
    }

    private void writeEndpointMiddlewareHelper(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            OperationShape operationShape
    ) {
        var opShapeId = operationShape.toShapeId();

        var input = OperationIndex.of(model).getInput(opShapeId).get();

        var memberShape = input.getMember(SUPPORTED_OPERATIONS.get(opShapeId)).get();

        writer.pushState();
        writer.putContext("middlewareAdder", getEndpointMiddlewareHelperName(opShapeId.getName()));
        writer.putContext("stackType", SymbolUtils.createPointableSymbolBuilder("Stack",
                SmithyGoDependency.SMITHY_MIDDLEWARE).build());
        writer.putContext("regHelper", SymbolUtils.createPointableSymbolBuilder("UpdateEndpoint",
                AwsCustomGoDependency.EVENTBRIDGE_CUSTOMIZATION).build());
        writer.putContext("helperOptions", SymbolUtils.createPointableSymbolBuilder("UpdateEndpointOptions",
                AwsCustomGoDependency.EVENTBRIDGE_CUSTOMIZATION).build());
        writer.putContext("memberHelper", getEndpointIdAccessorName(operationShape.getId().getName(),
                memberShape.getMemberName()));
        writer.putContext("resolver", AwsEndpointGenerator.ENDPOINT_RESOLVER_CONFIG_NAME);
        writer.putContext("resolverOptions", AwsEndpointGenerator.ENDPOINT_OPTIONS_CONFIG_NAME);
        writer.write("""
                     func $middlewareAdder:L(stack $stackType:P, o Options) error {
                         return $regHelper:T(stack, $helperOptions:T{
                             GetEndpointIDFromInput: $memberHelper:L,
                             EndpointResolver: o.$resolver:L,
                             EndpointResolverOptions: o.$resolverOptions:L,
                         })
                     }
                     """);
        writer.popState();
    }

    private void writeEndpointIdAccessorHelper(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            OperationShape operationShape
    ) {
        var opShapeId = operationShape.toShapeId();

        var input = OperationIndex.of(model).getInput(opShapeId).get();

        var memberShape = input.getMember(SUPPORTED_OPERATIONS.get(opShapeId)).get();

        // limit to string for now
        if (model.expectShape(memberShape.getTarget()).getType() != ShapeType.STRING) {
            throw new CodegenException("expect EventBridge EndpointId shape type to be string");
        }

        var accessorName = getEndpointIdAccessorName(operationShape.getId().getName(),
                memberShape.getMemberName());

        writer.pushState();
        writer.putContext("funcName", accessorName);
        writer.putContext("inputShape", symbolProvider.toSymbol(input));
        writer.putContext("memberName", memberShape.getMemberName());
        writer.writeDocs(String.format("""
                                       %s returns a pointer to string denoting a provided member value and a boolean indicating if the
                                       value is not nil
                                       """, accessorName));
        writer.write("""
                     func $funcName:L(input interface{}) (*string, bool) {
                         in := input.($inputShape:P)
                         if in.$memberName:L == nil {
                             return nil, false
                         }
                         return in.$memberName:L, true
                     }
                     """);
        writer.popState();
    }
}
