package software.amazon.smithy.aws.go.codegen;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

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
import software.amazon.smithy.rulesengine.language.syntax.parameters.Parameters;
import software.amazon.smithy.utils.StringUtils;

import java.util.Map;
import java.util.Optional;
import java.util.function.Consumer;

public class AwsEndpointResolverBuiltInGenerator implements GoIntegration {

    public static final String BUILTIN_RESOLVER_INTERFACE_TYPE = "BuiltInParameterResolver";
    public static final String BUILTIN_RESOLVER_IMPLEMENTATION_TYPE = "builtInResolver";


    private Map<String, Object> commonCodegenArgs;

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {
        var serviceShape = settings.getService(model);

        this.commonCodegenArgs = MapUtils.of(
            "resolverInterfaceType", SymbolUtils.createValueSymbolBuilder(BUILTIN_RESOLVER_INTERFACE_TYPE).build(),
            "resolverImplType", SymbolUtils.createValueSymbolBuilder(BUILTIN_RESOLVER_IMPLEMENTATION_TYPE).build()

        );


        Optional<EndpointRuleSet> ruleset = Optional.empty();
        var rulesetTrait = serviceShape.getTrait(EndpointRuleSetTrait.class);
        if (rulesetTrait.isPresent()) {
            ruleset = Optional.of(EndpointRuleSet.fromNode(rulesetTrait.get().getRuleSet()));
        }

        Parameters parameters = ruleset.get().getParameters();
        var content = new GoWriter.ChainWritable()
            .add(generateBuiltInInterfaceType())
            .add(generateBuiltInResolverEntryPoint(parameters))
            .add(generateBuiltInImplementationType(parameters))
            .compose();

        writerFactory.accept("endpoints.go", settings.getModuleName(), writer -> {
            writer.write("$W", content);
        });


    }

    private GoWriter.Writable generateBuiltInInterfaceType() {
        return goTemplate(
            """
            type $resolverInterfaceType:T interface {
                ResolveBuiltIn(name string) (value interface{}, ok bool)
            }
            """,
            commonCodegenArgs);
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
                SymbolUtils.createValueSymbolBuilder(BUILTIN_RESOLVER_INTERFACE_TYPE).build(),
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

    private GoWriter.Writable generateBuiltInImplementationType(Parameters parameters) {
        return goTemplate(
            """
            type $resolverImplType:T struct {
                $builtInMembers:W
            }
            """,
            commonCodegenArgs,
            MapUtils.of(
                "builtInMembers", generateBuiltInResolverMembers(parameters)));
    }

    private GoWriter.Writable generateBuiltInResolverMembers(Parameters parameters) {
        return (GoWriter w) -> {
            parameters.toList().stream().filter(
                p -> p.getBuiltIn().isPresent())
                .forEach(parameter -> {
                    w.write(
                        """
                            $L $T
                        """,
                        getExportedParameterName(parameter),
                        parameterAsSymbol(parameter)
                    );
                    w.write("");
            });
        };
    }

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
