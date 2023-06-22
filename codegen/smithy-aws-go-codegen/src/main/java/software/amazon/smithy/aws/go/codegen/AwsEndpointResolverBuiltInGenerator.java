package software.amazon.smithy.aws.go.codegen;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.aws.go.codegen.EndpointGenerator.USE_FIPS_ENDPOINT_OPTION;
import static software.amazon.smithy.aws.go.codegen.EndpointGenerator.DUAL_STACK_ENDPOINT_OPTION;

import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.go.codegen.endpoints.EndpointResolutionGenerator;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.MapUtils;

import software.amazon.smithy.rulesengine.language.EndpointRuleSet;
import software.amazon.smithy.rulesengine.traits.EndpointRuleSetTrait;
import software.amazon.smithy.rulesengine.language.syntax.parameters.Parameter;
import software.amazon.smithy.rulesengine.language.syntax.parameters.ParameterType;
import software.amazon.smithy.rulesengine.language.syntax.parameters.Parameters;
import software.amazon.smithy.utils.StringUtils;

import java.util.Map;
import java.util.Optional;
import java.util.function.Consumer;

public class AwsEndpointResolverBuiltInGenerator implements GoIntegration {

    public static final String BUILTIN_RESOLVER_INTERFACE_TYPE = "BuiltInParameterResolver";
    public static final String BUILTIN_RESOLVER_IMPLEMENTATION_TYPE = "BuiltInResolver";


    private Map<String, Object> commonCodegenArgs;

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


        var content = new GoWriter.ChainWritable()
            .add(generateBuiltInInterfaceType())
            .add(generateBuiltInImplementationType(parameters))
            .add(generateBuiltInResolverMethod(parameters))
            .compose();

        for (Parameter parameter : parameters.toList()) {
            if (parameter.getBuiltIn().isPresent()) {
                writerFactory.accept("endpoints.go", settings.getModuleName(), writer -> {
                    writer.write("$W", content);
                });
                break;
            }
        }

    }

    private GoWriter.Writable generateBuiltInInterfaceType() {
        return (GoWriter writer) -> {
            writer.write(
                """
                    type $T interface {
                        ResolveBuiltIns(*EndpointParameters) error
                    }
                """,
                SymbolUtils.createValueSymbolBuilder(BUILTIN_RESOLVER_INTERFACE_TYPE).build()
            );
        };
    }

    // private GoWriter.Writable generateBuiltInResolverEntryPoint(Parameters parameters) {
    //     return (GoWriter w) -> {
    //         w.write(
    //             """
    //             func resolveBuiltIns(parameters $P, resolver $T) {
    //                 var value interface{}; var present bool
    //                 $W
    //             }
    //             """,
    //             SymbolUtils.createValueSymbolBuilder(EndpointResolutionGenerator.PARAMETERS_TYPE_NAME).build(),
    //             SymbolUtils.createValueSymbolBuilder(BUILTIN_RESOLVER_INTERFACE_TYPE).build(),
    //             generateBuiltInResolutionInvocation(parameters)
    //         );
    //     };
    // }

    // private GoWriter.Writable generateBuiltInResolutionInvocation(Parameters parameters) {
    //     return (GoWriter w) -> {
    //         parameters.toList().stream().filter(
    //             p -> p.getBuiltIn().isPresent())
    //             .forEach(parameter -> {
    //                 w.write(
    //                     """
    //                         value, present = resolver.ResolveBuiltIn(\"$L\")
    //                         if v, ok := value.($T); present && ok {
    //                             parameters.$L = &v
    //                         }
    //                     """,
    //                     parameter.getBuiltIn(),
    //                     parameterAsSymbol(parameter),
    //                     getExportedParameterName(parameter)
    //                 );
    //                 w.write("");
    //         });
    //     };
    // }

    private GoWriter.Writable generateBuiltInImplementationType(Parameters parameters) {
        return (GoWriter writer) -> {
            writer.write(
                """
                type $T struct {
                    $W
                }
                """,
                SymbolUtils.createValueSymbolBuilder(BUILTIN_RESOLVER_IMPLEMENTATION_TYPE).build(),
                generateBuiltInResolverMembers(parameters)
            );
        };
    }

    private GoWriter.Writable generateBuiltInResolverMembers(Parameters parameters) {
        return (GoWriter w) -> {
            parameters.toList().stream().filter(
                p -> p.getBuiltIn().isPresent())
                .forEach(parameter -> {
                    if (parameter.getBuiltIn().get().equals("SDK::Endpoint")) {
                        w.write("$L $P",
                            getExportedParameterName(parameter),
                            parameterAsSymbol(parameter));
                    } else if (parameter.getBuiltIn().get().equals("AWS::UseFIPS")) {
                        w.write("$L $T",
                            getExportedParameterName(parameter),
                            SymbolUtils.createValueSymbolBuilder("FIPSEndpointState", AwsGoDependency.AWS_CORE).build());
                    } else if (parameter.getBuiltIn().get().equals("AWS::UseDualStack")) {
                        w.write("$L $T",
                            getExportedParameterName(parameter),
                            SymbolUtils.createValueSymbolBuilder("DualStackEndpointState", AwsGoDependency.AWS_CORE).build());
                    } else {
                        w.write("$L $T",
                            getExportedParameterName(parameter),
                            parameterAsSymbol(parameter));
                    }
                    w.insertTrailingNewline();
            });
        };
    }


    private GoWriter.Writable generateBuiltInResolverMethod(Parameters parameters) {
        return (GoWriter writer) -> {
            writer.write(
                """
                func (b *BuiltInResolver) ResolveBuiltIns(params *EndpointParameters) error {
                """
            );

            parameters.toList().stream().filter(
                p -> p.getBuiltIn().isPresent())
                .forEach(parameter -> {
                    if (parameter.getBuiltIn().get().equals("SDK::Endpoint")) {
                        writer.write("$W", generateSdkEndpointBuiltInResolver(parameter));
                    } else if (parameter.getBuiltIn().get().equals("AWS::Region")) {
                        writer.write("$W", generateAwsRegionBuiltInResolver(parameter));
                    } else if (parameter.getBuiltIn().get().equals("AWS::UseFIPS")) {
                        writer.write("$W", generateAwsFipsBuiltInResolver(parameter));
                    } else if (parameter.getBuiltIn().get().equals("AWS::UseDualStack")) {
                        writer.write("$W", generateAwsDualStackBuiltInResolver(parameter));
                    } else if (parameter.getType() == ParameterType.STRING) {
                        writer.write(
                            "params.$L = $T(b.$L)",
                            parameter.getName().asString(),
                            SymbolUtils.createValueSymbolBuilder("String", AwsGoDependency.AWS_CORE).build(),
                            getExportedParameterName(parameter)
                        );
                    } else if (parameter.getType() == ParameterType.BOOLEAN) {
                        writer.write(
                            "params.$L = $T(b.$L)",
                            parameter.getName().asString(),
                            SymbolUtils.createValueSymbolBuilder("Bool", AwsGoDependency.AWS_CORE).build(),
                            getExportedParameterName(parameter)
                        );
                    } else {
                        throw new CodegenException(String.format("Invalid Builtin %s", parameter.getBuiltIn().get()));
                    }
                }
            );

            writer.write(
                """
                    return nil
                }
                """
            );
        };
    }

    private GoWriter.Writable generateSdkEndpointBuiltInResolver(Parameter parameter) {
        return (GoWriter writer) -> {
            writer.write(
                """

                params.$L = b.Endpoint

                """,
                parameter.getName().asString()
            );
        };
    }
    

    private GoWriter.Writable generateAwsRegionBuiltInResolver(Parameter parameter) {
        return (GoWriter writer) -> {
            writer.write(
                """

                    region, _ := mapPseudoRegion(b.Region)
                    if len(region) == 0 {
                        return fmt.Errorf(\"Could not resolve AWS::Region\")
                    } else {
                        params.$L = $T(region)
                    }
                    

                """,
                parameter.getName().asString(),
                SymbolUtils.createValueSymbolBuilder("String", AwsGoDependency.AWS_CORE).build()
            );
        };
    }

    private GoWriter.Writable generateAwsFipsBuiltInResolver(Parameter parameter) {
        return (GoWriter writer) -> {
            String paramName = parameter.getName().asString();
            writer.write(
                """
                    if b.UseFIPS == $T {
                        params.$L = $T(true)
                    } else {
                        params.$L = $T(false)
                    }
                    
                """,
                SymbolUtils.createValueSymbolBuilder("FIPSEndpointStateEnabled", AwsGoDependency.AWS_CORE).build(),
                paramName,
                SymbolUtils.createValueSymbolBuilder("Bool", AwsGoDependency.AWS_CORE).build(),
                paramName,
                SymbolUtils.createValueSymbolBuilder("Bool", AwsGoDependency.AWS_CORE).build()
            );
        };
    }

    private GoWriter.Writable generateAwsDualStackBuiltInResolver(Parameter parameter) {
        return (GoWriter writer) -> {
            String paramName = parameter.getName().asString();
            writer.write(
                """
                    if b.UseDualStack == $T {
                        params.$L = $T(true)
                    } else {
                        params.$L = $T(false)
                    }
                    
                """,
                SymbolUtils.createValueSymbolBuilder("DualStackEndpointStateEnabled", AwsGoDependency.AWS_CORE).build(),
                paramName,
                SymbolUtils.createValueSymbolBuilder("Bool", AwsGoDependency.AWS_CORE).build(),
                paramName,
                SymbolUtils.createValueSymbolBuilder("Bool", AwsGoDependency.AWS_CORE).build()
            );
        };
    }




    // private GoWriter.Writable generatePseudoRegionUtility() {
    //     return (GoWriter writer) -> {
    //         writer.write(
    //             """
    //             func mapPseudoRegion(pr string) (region string, fips aws.FIPSEndpointState) {
    //                 const fipsInfix = \"-fips-\"
    //                 const fipsPrefix = \"fips-\"
    //                 const fipsSuffix = \"-fips\"
                
    //                 if strings.Contains(pr, fipsInfix) ||
    //                     strings.Contains(pr, fipsPrefix) ||
    //                     strings.Contains(pr, fipsSuffix) {
    //                     region = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
    //                         pr, fipsInfix, "-"), fipsPrefix, ""), fipsSuffix, "")
    //                     fips = aws.FIPSEndpointStateEnabled
    //                 } else {
    //                     region = pr
    //                 }
                
    //                 return region, fips
    //             }
    //             """
    //         );
    //     };
    // }


    // private GoWriter.Writable generateForcePathStyleMethod(Parameters parameters) {
    //     return (GoWriter w) -> {
    //         parameters.toList().stream().filter(
    //             p -> p.getBuiltIn() == "AWS:S3::ForcePathStyle")
    //             .forEach(parameter -> {
    //                 w.write(
    //                     """
    //                     func (b *builtInResolver) resolveForcePathStyle() (value interface{}, ok bool) {

    //                     }
    //                     """

    //             }

    //         )
    //     }
    // }

    @Override
    public void renderEndpointBuiltInField(GoWriter writer) {
        writer.write("BuiltInResolver $T", SymbolUtils.createValueSymbolBuilder("BuiltInParameterResolver").build());
    }

    @Override
    public void renderEndpointBuiltInInvocation(GoWriter writer) {
        writer.write("m.BuiltInResolver.ResolveBuiltIns(&params)");
    }

    @Override 
    public void renderEndpointBuiltInInitialization(GoWriter writer, Parameters parameters) {
        writer.write(
            """
                BuiltInResolver: &$T{
                    $W
                },
            """,
            SymbolUtils.createValueSymbolBuilder("BuiltInResolver").build(),
            generateBuiltInInitializeFieldMembers(parameters)
        );
    }

    private GoWriter.Writable generateBuiltInInitializeFieldMembers(Parameters parameters) {
        return (GoWriter writer) -> {
            parameters.toList().stream().filter(
                p -> p.getBuiltIn().isPresent())
                .forEach(parameter -> {
                    if (parameter.getBuiltIn().get().equals("SDK::Endpoint")) {
                        writer.write("$L: options.BaseEndpoint,", getExportedParameterName(parameter));
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
                        writer.write("$L: options.UseARNRegion,", getExportedParameterName(parameter));
                        writer.insertTrailingNewline();
                    }
                    if (parameter.getBuiltIn().get().equals("AWS::S3::DisableMultiRegionAccessPoints")) {
                        writer.write("$L: options.DisableMultiRegionAccessPoints,", getExportedParameterName(parameter));
                        writer.insertTrailingNewline();
                    }
                    if (parameter.getBuiltIn().get().equals("AWS::S3Control::UseArnRegion")) {
                        writer.write("$L: options.UseARNRegion,", getExportedParameterName(parameter));
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
