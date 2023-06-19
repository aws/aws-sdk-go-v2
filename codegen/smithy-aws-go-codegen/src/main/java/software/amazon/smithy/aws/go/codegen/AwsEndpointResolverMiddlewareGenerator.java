package software.amazon.smithy.aws.go.codegen;

import static software.amazon.smithy.aws.go.codegen.EndpointGenerator.USE_FIPS_ENDPOINT_OPTION;
import static software.amazon.smithy.aws.go.codegen.EndpointGenerator.DUAL_STACK_ENDPOINT_OPTION;

import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoCodegenPlugin;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.MiddlewareIdentifier;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.OperationIndex;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.ToShapeId;
import software.amazon.smithy.rulesengine.language.EndpointRuleSet;
import software.amazon.smithy.rulesengine.language.syntax.parameters.Parameter;
import software.amazon.smithy.rulesengine.language.syntax.parameters.ParameterType;
import software.amazon.smithy.rulesengine.language.syntax.parameters.Parameters;
import software.amazon.smithy.utils.StringUtils;
import software.amazon.smithy.rulesengine.traits.ClientContextParamsTrait;
import software.amazon.smithy.rulesengine.traits.ContextParamTrait;
import software.amazon.smithy.rulesengine.traits.EndpointRuleSetTrait;
import software.amazon.smithy.rulesengine.traits.StaticContextParamsTrait;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.go.codegen.endpoints.EndpointResolutionGenerator;
import software.amazon.smithy.utils.ListUtils;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.function.Consumer;

public class AwsEndpointResolverMiddlewareGenerator implements GoIntegration {

    private final List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();


    private static String getAddEndpointMiddlewareFuncName(String operationName) {
        return String.format("add%sResolveEndpointMiddleware", operationName);
    }

    private static String getMiddlewareObjectName(String operationName) {
        return String.format("op%sResolveEndpointMiddleware", operationName);
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        runtimeClientPlugins.add(RuntimeClientPlugin.builder()
        .configFields(ListUtils.of(
            ConfigField.builder()
                    .name("MutableBaseEndpoint")
                    .type(SymbolUtils.createPointableSymbolBuilder("URL", SmithyGoDependency.NET_URL).build())
                    .documentation(
                        """
                        This endpoint will be given as input to an EndpointResolverV2.
                        It is used for providing a custom base endpoint that is subject 
                        to modifications by the processing EndpointResolverV2.        
                        """
                    )
                    .build()
        ))
        .build());
        return runtimeClientPlugins;
    }

    @Override
    public void processFinalizedModel(GoSettings settings, Model model) {
        ServiceShape service = settings.getService(model);
        var rulesetTrait = service.getTrait(EndpointRuleSetTrait.class);
        Optional<EndpointRuleSet> rulesetOpt = (rulesetTrait.isPresent()) 
        ? Optional.of(EndpointRuleSet.fromNode(rulesetTrait.get().getRuleSet()))
        : Optional.empty();
        var clientContextParamsTrait = service.getTrait(ClientContextParamsTrait.class);

        var topDownIndex = TopDownIndex.of(model);

        for (ToShapeId operationId : topDownIndex.getContainedOperations(service)) {
            OperationShape operationShape = model.expectShape(operationId.toShapeId(), OperationShape.class);

            // Create a symbol provider because one is not available in this call.
            SymbolProvider symbolProvider = GoCodegenPlugin.createSymbolProvider(model, settings);

            // Input helper
            String inputHelperFuncName = getAddEndpointMiddlewareFuncName(
                    symbolProvider.toSymbol(operationShape).getName()
            );
            runtimeClientPlugins.add(RuntimeClientPlugin.builder()
                    .operationPredicate((m, s, o) -> {
                        return o.equals(operationShape);
                    })
                    .registerMiddleware(MiddlewareRegistrar.builder()
                            .resolvedFunction(SymbolUtils.createValueSymbolBuilder(inputHelperFuncName)
                                    .build())
                            .useClientOptions()
                            .build())
                    .build());

            if (clientContextParamsTrait.isPresent()) {
                if (rulesetOpt.isPresent()) {
                    var clientContextParams = clientContextParamsTrait.get();
                    var parameters = rulesetOpt.get().getParameters();
                    parameters.toList().stream().forEach(param -> {
                        if (
                            clientContextParams.getParameters().containsKey(param.getName().asString()) &&
                            !param.getBuiltIn().isPresent()
                        ) {
                            var documentation = param.getDocumentation().isPresent() ?  param.getDocumentation().get() : "";
                            runtimeClientPlugins.add(RuntimeClientPlugin.builder()
                            .configFields(ListUtils.of(
                                ConfigField.builder()
                                        .name(getExportedParameterName(param))
                                        .type(parameterAsSymbol(param))
                                        .documentation(documentation)
                                        .build()
                            ))
                            .build());
                        }
                    });
                }
            }
        }
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {

        var serviceShape = settings.getService(model);

        Optional<EndpointRuleSet> ruleset = Optional.empty();
        var rulesetTrait = serviceShape.getTrait(EndpointRuleSetTrait.class);
        if (rulesetTrait.isPresent()) {
            ruleset = Optional.of(EndpointRuleSet.fromNode(rulesetTrait.get().getRuleSet()));
        }

        Parameters parameters = ruleset.get().getParameters();

        writerFactory.accept("endpoints.go", settings.getModuleName(), writer -> {
            writer.write("$W", generateBuiltInResolverEntryPoint(parameters));
        });
    }

    private GoWriter.Writable generateBuiltInResolverEntryPoint(Parameters parameters) {
        return (GoWriter w) -> {
            w.write(
                """
                func resolveBuiltIns(parameters $P, resolver $T) {
                    var value interface{}; var present bool
                    $W
                }
                """,
                SymbolUtils.createValueSymbolBuilder(EndpointResolutionGenerator.PARAMETERS_TYPE_NAME).build(),
                SymbolUtils.createValueSymbolBuilder("BuiltInParameterResolver", AwsGoDependency.INTERNAL_ENDPOINTS).build(),
                generateBuiltInResolutionInvocation(parameters)
            );
        };
    }

    private GoWriter.Writable generateBuiltInResolutionInvocation(Parameters parameters) {
        return (GoWriter w) -> {
            parameters.toList().stream().filter(
                p -> p.getBuiltIn().isPresent())
                .forEach(parameter -> {
                    w.write(
                        """
                            value, present = resolver.ResolveBuiltIn(\"$L\")
                            if v, ok := value.($T); present && ok {
                                parameters.$L = &v
                            }
                        """,
                        parameter.getBuiltIn(),
                        parameterAsSymbol(parameter),
                        getExportedParameterName(parameter)
                    );
                    w.write("");
            });
        };
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
  
                        $W
                        """,
                        generateMiddlewareType(parameters, clientContextParamsTrait, operationName),
                        generateMiddlewareMethods(parameters, clientContextParamsTrait, symbolProvider, operationShape, model),
                        generateMiddlewareAdder(parameters, operationName, clientContextParamsTrait)
                    );
                }
            });
        }
    }

    private GoWriter.Writable generateMiddlewareAdder(Parameters parameters, String operationName, Optional<ClientContextParamsTrait> clientContextParamsTrait) {
        return (GoWriter writer) -> {
            writer.write(
                """
                func $L(stack $P, options Options) error {
                    return stack.Serialize.Insert(&$L{
                        BuiltInResolver: &$T{
                            $W
                        },
                        EndpointResolver: options.EndpointResolverV2,
                        $W
                    }, \"ResolveEndpoint\", middleware.After)
                }
                """,
                SymbolUtils.createValueSymbolBuilder(getAddEndpointMiddlewareFuncName(operationName)).build(),
                SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE).build(),
                SymbolUtils.createValueSymbolBuilder(getMiddlewareObjectName(operationName)).build(),
                SymbolUtils.createValueSymbolBuilder("BuiltInResolver", AwsGoDependency.INTERNAL_ENDPOINTS).build(),
                generateBuiltInInitializeFieldMembers(parameters),
                generateClientContextParamInitialization(parameters, clientContextParamsTrait)

    
            );
        };
    }

    private GoWriter.Writable generateClientContextParamInitialization(Parameters parameters, Optional<ClientContextParamsTrait> clientContextParamsTrait) {
        return (GoWriter writer) -> {
            if (clientContextParamsTrait.isPresent()) {
                var clientContextParams = clientContextParamsTrait.get();
                parameters.toList().stream().forEach(param -> {
                    if (
                        clientContextParams.getParameters().containsKey(param.getName().asString()) &&
                        !param.getBuiltIn().isPresent()
                    ) {
                        var name = getExportedParameterName(param);
                        writer.write("$L: options.$L,", name, name);
                    }
                });
            }
        };
    }

    private GoWriter.Writable generateMiddlewareType(Parameters parameters, Optional<ClientContextParamsTrait> clientContextParamsTrait, String operationName) {
        return (GoWriter w) -> {
            w.openBlock("type $L struct {", "}", getMiddlewareObjectName(operationName), () -> {
                w.write("EndpointResolver $T", SymbolUtils.createValueSymbolBuilder("EndpointResolverV2").build());
                w.write("BuiltInResolver $T", SymbolUtils.createValueSymbolBuilder("BuiltInParameterResolver", AwsGoDependency.INTERNAL_ENDPOINTS).build());
                if (clientContextParamsTrait.isPresent()) {
                    var clientContextParams = clientContextParamsTrait.get();
                    parameters.toList().stream().forEach(param -> {
                        if (clientContextParams.getParameters().containsKey(param.getName().asString())) {
                            w.write("$L $P", getExportedParameterName(param), parameterAsSymbol(param));
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
                        writer.write("$W", generateMiddlewareResolverBody(operationShape, model, parameters, clientContextParamsTrait));
                    });
        };
    }

    private GoWriter.Writable generateMiddlewareResolverBody(OperationShape operationShape, Model model, Parameters parameters, Optional<ClientContextParamsTrait> clientContextParamsTrait) {

        return (GoWriter writer) -> {
            var fmtErrorSymbol = SymbolUtils.createValueSymbolBuilder("Errorf", SmithyGoDependency.FMT).build();
            writer.write(
                """
                    $W
                
                    $W
                
                    if m.EndpointResolver == nil {
                        return out, metadata, $T(\"expected endpoint resolver to not be nil\")
                    }
                
                    if m.BuiltInResolver == nil {
                        m.BuiltInResolver = &$T{}
                    }
                
                    params := EndpointParameters{}

                    resolveBuiltIns(params, m.BuiltInResolver)

                    $W

                    $W

                    $W

                    var resolvedEndpoint $T
                    resolvedEndpoint, err = m.EndpointResolver.ResolveEndpoint(ctx, params)
                    if err != nil {
                        return out, metadata, $T(\"failed to resolve service endpoint, %w\", err)
                    }

                    req.URL = &resolvedEndpoint.URI

                    $W

                    return next.HandleSerialize(ctx, in)
                """,
                generateRequestValidator(),
                generateInputValidator(model, operationShape),
                fmtErrorSymbol,
                SymbolUtils.createValueSymbolBuilder("NopBuiltInResolver", AwsGoDependency.INTERNAL_ENDPOINTS).build(),
                generateClientContextParamBinding(parameters, clientContextParamsTrait),
                generateContextParamBinding(operationShape, model),
                generateStaticContextParamBinding(parameters, operationShape),
                SymbolUtils.createValueSymbolBuilder("Endpoint", SmithyGoDependency.SMITHY_ENDPOINTS).build(),
                fmtErrorSymbol,
                generateAuthSchemeResolution()
            );
        };
    }

    private GoWriter.Writable generateRequestValidator() {
        return (GoWriter writer) -> {
            writer.write(
                """
                    req, ok := in.Request.($P)
                    if !ok {
                        return out, metadata, $T(\"unknown transport type %T\", in.Request)
                    }
                """,
                SymbolUtils.createPointableSymbolBuilder("Request", SmithyGoDependency.SMITHY_HTTP_TRANSPORT).build(),
                SymbolUtils.createValueSymbolBuilder("Errorf", SmithyGoDependency.FMT).build()
            );
        };
    }

    private GoWriter.Writable generateInputValidator(Model model, OperationShape operationShape) {
        var opIndex = OperationIndex.of(model);
        var inputOpt = opIndex.getInput(operationShape);
        GoWriter.Writable inputValidator = (GoWriter writer) -> {writer.write("");};

        if (inputOpt.isPresent()) {
            var input = inputOpt.get();
            for (var inputMember : input.getAllMembers().values()) {
                var contextParamTraitOpt = inputMember.getTrait(ContextParamTrait.class);
                if (contextParamTraitOpt.isPresent()) {
                    inputValidator = (GoWriter writer) -> {writer.write(
                        """
                            input, ok := in.Parameters.($P)
                            if !ok {
                                return out, metadata, $T(\"unknown transport type %T\", in.Request)
                            }     
                        """,
                        SymbolUtils.createPointableSymbolBuilder(operationShape.getInput().get().getName()).build(),
                        SymbolUtils.createValueSymbolBuilder("Errorf", SmithyGoDependency.FMT).build()
                    );};
                }
            }
        }
        return inputValidator;
    }

    private GoWriter.Writable generateAuthSchemeResolution() {
        return (GoWriter writer) -> {
            writer.write(
                """
                    auth, ok := resolvedEndpoint.Properties.Get(\"authSchemes\").([]interface{})
                    if ok {
                        for _, schemes := range auth {
                            scheme, ok := schemes.(map[string]interface{})
                            if !ok {
                                continue
                            }
                            if len($T(ctx)) == 0 {
                                signingName := scheme[\"signingName\"].(string)
                                if len(signingName) == 0 {
                                    signingName = \"s3\"
                                }
                                ctx = $T(ctx, signingName)
                            }
                        }
                    }
                """,
                SymbolUtils.createValueSymbolBuilder("GetSigningName", AwsGoDependency.AWS_MIDDLEWARE).build(),
                SymbolUtils.createValueSymbolBuilder("SetSigningName", AwsGoDependency.AWS_MIDDLEWARE).build()
            );
        };
    }

    private GoWriter.Writable generateClientContextParamBinding(Parameters parameters, Optional<ClientContextParamsTrait> clientContextParamsTrait) {
        return (GoWriter writer) -> {
            if (clientContextParamsTrait.isPresent()) {
                var clientContextParams = clientContextParamsTrait.get();
                parameters.toList().stream().forEach(param -> {
                    if (clientContextParams.getParameters().containsKey(param.getName().asString())
                    && !param.getBuiltIn().isPresent()
                    ) {
                        var name = getExportedParameterName(param);
                        writer.write("params.$L = m.$L", name, name);
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

    private GoWriter.Writable generateStaticContextParamBinding(Parameters parameters, OperationShape operationShape) {
        var staticContextParamTraitOpt = operationShape.getTrait(StaticContextParamsTrait.class);
        return (GoWriter writer) -> {
            parameters.toList().stream().forEach( param -> {
                if (staticContextParamTraitOpt.isPresent()) {
                    var paramName = param.getName().asString();

                    var staticParam = staticContextParamTraitOpt
                                            .get()
                                            .getParameters()
                                            .get(paramName);
                    if (staticParam != null) {
                        Symbol valueWrapper;
                        if (param.getType() == ParameterType.BOOLEAN) {
                            valueWrapper = SymbolUtils.createValueSymbolBuilder("Bool", AwsGoDependency.AWS_CORE).build();
                        } else if (param.getType() == ParameterType.STRING) {
                            valueWrapper = SymbolUtils.createValueSymbolBuilder("String", AwsGoDependency.AWS_CORE).build();
                        } else {
                            throw new CodegenException(String.format("unexpected static context param type: %s", param.getType()));
                        }
                        writer.write("params.$L = $T($L)", paramName, valueWrapper, staticParam.getValue());
                    }
                }
            });
            writer.write("");
        };
    }


    private GoWriter.Writable generateBuiltInInitializeFieldMembers(Parameters parameters) {
        return (GoWriter writer) -> {
            parameters.toList().stream().filter(
                p -> p.getBuiltIn().isPresent())
                .forEach(parameter -> {
                    if (parameter.getBuiltIn().get().equals("SDK::Endpoint")) {
                        writer.write("$L: options.MutableBaseEndpoint,", getExportedParameterName(parameter));
                        writer.insertTrailingNewline();
                    }
                    if (parameter.getBuiltIn().get().equals("AWS::Region")) {
                        writer.write("$L: options.Region,", getExportedParameterName(parameter));
                        writer.insertTrailingNewline();
                    }
                    if (parameter.getBuiltIn().get().equals("AWS::UseFIPS")){
                        writer.write("$L: options.$L.$L,", getExportedParameterName(parameter), "EndpointOptions", USE_FIPS_ENDPOINT_OPTION);
                        writer.insertTrailingNewline();
                    }
                    if (parameter.getBuiltIn().get().equals("AWS::UseDualStack")) {
                        writer.write("$L: options.$L.$L,", getExportedParameterName(parameter), "EndpointOptions", DUAL_STACK_ENDPOINT_OPTION);
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
