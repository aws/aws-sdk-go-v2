package software.amazon.smithy.aws.go.codegen;


import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.go.codegen.endpoints.EndpointResolutionGenerator;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.utils.SetUtils;

import java.util.List;
import java.util.function.Consumer;

public class AwsEndpointResolverInitializerGenerator implements GoIntegration {

    public static final String RESOLVE_ENDPOINT_RESOLVER_V2 = "resolve" + EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME;

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {

        writerFactory.accept("endpoints.go", settings.getModuleName(), writer -> {
            writer.write("$W", generateResolveMethod());
        });
    }


    private GoWriter.Writable generateResolveMethod() {
        return goTemplate(
            """
                func $resolveMethodName:L(options *Options) {
                    if options.EndpointResolverV2 == nil {
                        options.EndpointResolverV2 = $newResolverFuncName:L()
                    }
                }

            """,
            MapUtils.of(
                       "resolveMethodName", RESOLVE_ENDPOINT_RESOLVER_V2,
                    "newResolverFuncName", EndpointResolutionGenerator.NEW_RESOLVER_FUNC_NAME
            ));

    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .configFields(SetUtils.of(
                                ConfigField.builder()
                                        .name(EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME)
                                        .type(SymbolUtils.createValueSymbolBuilder(EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME)
                                                .build())
                                        .documentation(String.format(
                                            """
                                            Resolves the endpoint used for a particular service. This should be used over the
                                            deprecated %s
                                            """,
                                            EndpointGenerator.RESOLVER_INTERFACE_NAME
                                        ))
                                        .withHelper(true)
                                        .build()
                        ))
                        .addConfigFieldResolver(ConfigFieldResolver.builder()
                                        .location(ConfigFieldResolver.Location.OPERATION)
                                        .target(ConfigFieldResolver.Target.INITIALIZATION)
                                        .resolver(SymbolUtils.createValueSymbolBuilder(
                                                RESOLVE_ENDPOINT_RESOLVER_V2).build())
                                        .build())
                        .build());
    }
}
