package software.amazon.smithy.aws.go.codegen;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.aws.go.codegen.EndpointGenerator.RESOLVER_OPTIONS;
import static software.amazon.smithy.aws.go.codegen.EndpointGenerator.USE_FIPS_ENDPOINT_OPTION;
import static software.amazon.smithy.aws.go.codegen.EndpointGenerator.DUAL_STACK_ENDPOINT_OPTION;




import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoCodegenPlugin;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.MiddlewareIdentifier;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.OperationIndex;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.shapes.ToShapeId;
import software.amazon.smithy.rulesengine.language.EndpointRuleSet;
import software.amazon.smithy.rulesengine.language.syntax.parameters.Parameter;
import software.amazon.smithy.rulesengine.language.syntax.parameters.Parameters;
import software.amazon.smithy.utils.StringUtils;
import software.amazon.smithy.rulesengine.traits.ClientContextParamDefinition;
import software.amazon.smithy.rulesengine.traits.ClientContextParamsTrait;
import software.amazon.smithy.rulesengine.traits.ContextParamTrait;
import software.amazon.smithy.rulesengine.traits.EndpointRuleSetTrait;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.utils.MapUtils;




import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.function.Consumer;

public class AwsEndpointResolverMiddlewareGenerator implements GoIntegration {



    private final List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();

    private static String getAddEndpointMiddlewareFuncName(String operationName) {
        return String.format("add%sResolveEndpointMiddleware", operationName);
    }

    private static String getMiddlewareObjectName(String operationName) {
        return String.format("op%sResolveEndpointMiddlware", operationName);
    }



    @Override
    public void processFinalizedModel(GoSettings settings, Model model) {
        ServiceShape service = settings.getService(model);
        for (ShapeId operationId : service.getAllOperations()) {
            final OperationShape operation = model.expectShape(operationId, OperationShape.class);

            // Create a symbol provider because one is not available in this call.
            SymbolProvider symbolProvider = GoCodegenPlugin.createSymbolProvider(model, settings);

            // Input helper
            String inputHelperFuncName = getAddEndpointMiddlewareFuncName(
                    symbolProvider.toSymbol(operation).getName()
            );
            runtimeClientPlugins.add(RuntimeClientPlugin.builder()
                    .operationPredicate((m, s, o) -> {
                        return o.equals(operation);
                    })
                    .registerMiddleware(MiddlewareRegistrar.builder()
                            .resolvedFunction(SymbolUtils.createValueSymbolBuilder(inputHelperFuncName)
                                    .build())
                            .useClientOptions()
                            .build())
                    .build());
        }

    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return runtimeClientPlugins;
    }





    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {

        var serviceShape = settings.getService(model);

        var rulesetTrait = serviceShape.getTrait(EndpointRuleSetTrait.class);
        var clientContextParamsTrait = serviceShape.getTrait(ClientContextParamsTrait.class);

        Optional<EndpointRuleSet> rulesetOpt = (rulesetTrait.isPresent()) 
                ? Optional.of(EndpointRuleSet.fromNode(rulesetTrait.get().getRuleSet()))
                : Optional.empty();
        
        TopDownIndex topDownIndex = TopDownIndex.of(model);

        for (ToShapeId operation : topDownIndex.getContainedOperations(serviceShape)) {
            OperationShape operationShape = model.expectShape(operation.toShapeId(), OperationShape.class);
            goDelegator.useShapeWriter(operationShape, writer -> {
                if (rulesetOpt.isPresent()) {
                    var parameters = rulesetOpt.get().getParameters();
                    Symbol operationSymbol = symbolProvider.toSymbol(operationShape);
                    String operationName = operationSymbol.getName();


                    writer.write(
                        """
                        $W

                        $W
  
                        """,
                        generateMiddlewareType(parameters, clientContextParamsTrait, operationName),
                        generateMiddlewareMethods(parameters, clientContextParamsTrait, symbolProvider, operationShape, model)

                        // generateMiddlewareAdder();
                    );


                }

            });
        }
    }






    // private void generateMiddlewareAdder(GoWriter writer) {
    //     writer.write(
    //         """
    //         func $L(stack $P, options Options) error {
    //             return $L{
    //                 BuiltInResolver: $T{
    //                     $W
    //                 },
    //             }
    //         }
    //         """,
    //         SymbolUtils.createValueSymbolBuilder(getAddEndpointMiddlewareFuncName(operationName)).build(),
    //         SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE).build(),
    //         SymbolUtils.createValueSymbolBuilder(getMiddlewareObjectName(operationName)).build(),
    //         SymbolUtils.createValueSymbolBuilder("BuiltInResolver", AwsGoDependency.INTERNAL_ENDPOINTS).build(),
    //         generateBuiltInInitializeFieldMembers(parameters)

    //     );
    // }

    private GoWriter.Writable generateMiddlewareType(Parameters parameters, Optional<ClientContextParamsTrait> clientContextParamsTrait, String operationName) {
        return (GoWriter w) -> {
            w.openBlock("type $L struct {", "}", getMiddlewareObjectName(operationName), () -> {
                w.write("EndpointResolver $T", SymbolUtils.createValueSymbolBuilder("EndpointResolverV2").build());
                w.write("BuiltInResolver $T", SymbolUtils.createValueSymbolBuilder("BuiltInResolver", AwsGoDependency.INTERNAL_ENDPOINTS).build());
                if (clientContextParamsTrait.isPresent()) {
                    var clientContextParams = clientContextParamsTrait.get();
                    parameters.toList().stream().forEach(param -> {
                        if (clientContextParams.getParameters().containsKey(param.getName().asString())) {
                            w.write("$L $T", getExportedParameterName(param), parameterAsSymbol(param));
                        }
                    });
                }

            });
        };
    }

    private GoWriter.Writable generateMiddlewareMethods(Parameters parameters, Optional<ClientContextParamsTrait> clientContextParamsTrait, SymbolProvider symbolProvider, OperationShape operationShape, Model model) {
        Symbol operationSymbol = symbolProvider.toSymbol(operationShape);
        String operationName = operationSymbol.getName();


        String middlewareName = getMiddlewareObjectName(operationName);
        Symbol middlewareSymbol = SymbolUtils.createPointableSymbolBuilder(middlewareName).build();



        return (GoWriter writer) -> {
            writer.openBlock("func ($P) ID() string {", "}", middlewareSymbol, () -> {
                writer.writeInline("return ");
                MiddlewareIdentifier.string(middlewareName).writeInline(writer);
                writer.write("");
            });
    
            writer.write("");
    
            // each middleware must implement their given handlerMethodName in order to satisfy the interface for
            // their respective step.

            String handleMethodName = "HandleSerialize";
            Symbol contextType = SymbolUtils.createValueSymbolBuilder("Context", SmithyGoDependency.CONTEXT).build();
            Symbol metadataType = SymbolUtils.createValueSymbolBuilder("Metadata", SmithyGoDependency.SMITHY_MIDDLEWARE).build();
            var inputType = SymbolUtils.createValueSymbolBuilder("SerializeInput", SmithyGoDependency.SMITHY_MIDDLEWARE).build();
            var outputType = SymbolUtils.createValueSymbolBuilder("SerializeOutput", SmithyGoDependency.SMITHY_MIDDLEWARE).build();
            var handlerType = SymbolUtils.createValueSymbolBuilder("SerializeHandler", SmithyGoDependency.SMITHY_MIDDLEWARE).build();


            writer.openBlock("func (m $P) $L(ctx $T, in $T, next $T) (\n"
                            + "\tout $T, metadata $T, err error,\n"
                            + ") {", "}",
                    new Object[]{
                            middlewareSymbol, handleMethodName, contextType, inputType, handlerType, outputType,
                            metadataType,
                    },
                    () -> {
                        // TODO: isaiah fill in HandleSerialize method body
                        writer.write("$W", generateMiddlewareResolverBody(operationShape, model, parameters, clientContextParamsTrait));
                    });
        };
    }

    private GoWriter.Writable generateMiddlewareResolverBody(OperationShape operationShape, Model model, Parameters parameters, Optional<ClientContextParamsTrait> clientContextParamsTrait) {
        return (GoWriter writer) -> {
            writer.write(
                """
                    req, ok := in.Request.($P)
                    if !ok {
                        return out, metadata, fmt.Errorf(\"unknown transport type %T\", in.Request)
                    }
                
                    input, ok := in.Parameters.($P)
                    if !ok {
                        return out, metadata, fmt.Errorf(\"unknown transport type %T\", in.Request)
                    }
                
                    if m.Resolver == nil {
                        return out, metadata, fmt.Errorf(\"expected endpoint resolver to not be nil\")
                    }
                
                    if m.BuiltInResolver == nil {
                        m.BuiltInResolver = endpoint.NopBuiltInResolver{}
                    }
                
                    params := EndpointParameters{}

                    resolveBuiltIns(&params, m.BuiltInResolver)

                    $W

                    $W

                """,
                SymbolUtils.createPointableSymbolBuilder("Request", SmithyGoDependency.SMITHY_HTTP_TRANSPORT).build(),
                SymbolUtils.createPointableSymbolBuilder("PutObjectInput").build(),
                generateClientContextParamBinding(parameters, clientContextParamsTrait),
                generateContextParamBinding(operationShape, model)
            );
        };
    }

    private GoWriter.Writable generateClientContextParamBinding(Parameters parameters, Optional<ClientContextParamsTrait> clientContextParamsTrait) {
        return (GoWriter writer) -> {
            if (clientContextParamsTrait.isPresent()) {
                var clientContextParams = clientContextParamsTrait.get();
                parameters.toList().stream().forEach(param -> {
                    if (clientContextParams.getParameters().containsKey(param.getName().asString())) {
                        var name = getExportedParameterName(param);
                        writer.write("params.$L = &m.$L", name, name);
                    }
                });
            }
        };
    }

    private GoWriter.Writable generateContextParamBinding(OperationShape operationShape, Model model) {

        return (GoWriter writer) -> {
            var opIndex = OperationIndex.of(model);
            var inputOpt = opIndex.getInput(operationShape);
            if (inputOpt.isPresent()) {
                var input = inputOpt.get();
                input.getAllMembers().values().forEach(inputMember -> {
                    var contextParamTraitOpt = inputMember.getTrait(ContextParamTrait.class);
                    if (contextParamTraitOpt.isPresent()) {
                        var contextParamTrait = contextParamTraitOpt.get();
                        writer.write(
                            """
                            params.$L = input.$L     
                            """,
                            contextParamTrait.getName(),
                            contextParamTrait.getName()
                        );
                        writer.write("");
                    }
                });
            }
            writer.write("");

        };
    }


    private GoWriter.Writable generateBuiltInInitializeFieldMembers(Parameters parameters) {
        return (GoWriter writer) -> {
            parameters.toList().stream().filter(
                p -> p.getBuiltIn().isPresent())
                .forEach(parameter -> {
                    if (parameter.getBuiltIn().get().equals("AWS::Region")) {
                        writer.write("$L: options.Region,", getExportedParameterName(parameter));
                        writer.insertTrailingNewline();
                    }
                    if (parameter.getBuiltIn().get().equals("AWS::UseFIPS")){
                        writer.write("$L: options.$L.$L,", getExportedParameterName(parameter), RESOLVER_OPTIONS, USE_FIPS_ENDPOINT_OPTION);
                        writer.insertTrailingNewline();
                    }
                    if (parameter.getBuiltIn().get().equals("AWS::UseDualStack")) {
                        writer.write("$L: options.$L.$L,", getExportedParameterName(parameter), RESOLVER_OPTIONS, DUAL_STACK_ENDPOINT_OPTION);
                        writer.insertTrailingNewline();
                    }
                    if (parameter.getBuiltIn().get().equals("AWS::S3::ForcePathStyle")) {
                        writer.write("$L: options.UsePathStyle,", getExportedParameterName(parameter));
                        writer.insertTrailingNewline();
                    }
                    if (parameter.getBuiltIn().get().equals("AWS::S3::Accelerate")) {
                        writer.write("$L: options.UseAccelerate,", getExportedParameterName(parameter));
                        writer.insertTrailingNewline();
                    }
                    if (parameter.getBuiltIn().get().equals("AWS::S3::UseArnRegion")) {
                        writer.write("S3$L: options.UseARNRegion,", getExportedParameterName(parameter));
                        writer.insertTrailingNewline();
                    }
                    if (parameter.getBuiltIn().get().equals("AWS::S3::DisableMultiRegionAccessPoints")) {
                        writer.write("$L: options.DisableMultiRegionAccessPoints,", getExportedParameterName(parameter));
                        writer.insertTrailingNewline();
                    }
                    if (parameter.getBuiltIn().get().equals("AWS::S3Control::UseArnRegion")) {
                        writer.write("S3Control$L: options.UseARNRegion,", getExportedParameterName(parameter));
                        writer.insertTrailingNewline();
                    }
                }
            );
        };
    }


    public static String getExportedParameterName(Parameter parameter) {
        return StringUtils.capitalize(parameter.getName().asString());
    }

    public static Symbol parameterAsSymbol(Parameter parameter) {
        return switch (parameter.getType()) {
            case STRING -> SymbolUtils.createPointableSymbolBuilder("string")
                    .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true).build();

            case BOOLEAN -> SymbolUtils.createPointableSymbolBuilder("bool")
                    .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true).build();
        };
    }
}
