/*
 * Copyright 2023 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *  http://aws.amazon.com/apache2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package software.amazon.smithy.aws.go.codegen;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
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

import java.util.List;
import java.util.function.Consumer;
import java.util.Map;


public class AwsEndpointResolverInitializerGenerator implements GoIntegration {

    public static final String RESOLVE_ENDPOINT_RESOLVER_V2 = "resolve" + EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME;
    public static final String RESOLVE_BASE_ENDPOINT = "resolveBaseEndpoint";

    private Map<String, Object> commonCodegenArgs;


    private static final ConfigField EndpointResolverV2 = ConfigField.builder()
            .name(EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME)
            .type(buildPackageSymbol(EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME))
            .documentation(String.format("""
                    Resolves the endpoint used for a particular service operation.
                    This should be used over the deprecated %s.
                    """, EndpointGenerator.RESOLVER_INTERFACE_NAME)
            )
            .withHelper(true)
            .build();

    private static final ConfigFieldResolver ResolveEndpointResolverV2 = ConfigFieldResolver.builder()
            .resolver(buildPackageSymbol(RESOLVE_ENDPOINT_RESOLVER_V2))
            .location(ConfigFieldResolver.Location.CLIENT)
            .target(ConfigFieldResolver.Target.INITIALIZATION)
            .build();

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {
        String sdkId = settings.getService(model).expectTrait(ServiceTrait.class).getSdkId();
        this.commonCodegenArgs = MapUtils.of(
            "envSdkId", sdkId.toUpperCase().replaceAll(" ", "_"),
            "configSdkId", sdkId.toLowerCase().replaceAll(" ", "_"),
            "urlSdkId", sdkId.toLowerCase().replaceAll(" ", "-"),
            "testing", SymbolUtils.createPointableSymbolBuilder("T", SmithyGoDependency.TESTING).build(),
            "awsString", SymbolUtils.createValueSymbolBuilder("String", AwsGoDependency.AWS_CORE).build()
        );

        writerFactory.accept("endpoints.go", settings.getModuleName(), writer -> {
            writer.write(
                """
                $W

                $W

                """
                ,
                generateResolveMethod(),
                generateConfigResolverMethod(
                    settings.getService(model).expectTrait(ServiceTrait.class).getSdkId())
                );
        });
    }

    private GoWriter.Writable generateConfigResolverMethod(String sdkId) {
        return goTemplate(
            """
                func $resolveMethodName:L(cfg $awsConfig:T, o *Options) {
                    if cfg.BaseEndpoint != nil {
                        o.BaseEndpoint = cfg.BaseEndpoint
                    }

                    _, g := $lookupEnv:T("AWS_ENDPOINT_URL")
                    _, s := $lookupEnv:T("AWS_ENDPOINT_URL_$envSdkId:L")

                    if g && !s  {
                        return
                    }

                    value, found, err := $resolveServiceEndpoint:T(context.Background(), "$sdkId:L", cfg.ConfigSources)
                    if found && err == nil {
                        o.BaseEndpoint = &value
                    }
                }
            """,
            MapUtils.of(
                "resolveMethodName", RESOLVE_BASE_ENDPOINT,
                "awsConfig", SymbolUtils.createValueSymbolBuilder("Config", AwsGoDependency.AWS_CORE).build(),
                "lookupEnv", SymbolUtils.createValueSymbolBuilder("LookupEnv", SmithyGoDependency.OS).build(),
                "resolveServiceEndpoint", SymbolUtils.createValueSymbolBuilder(
                                        "ResolveServiceBaseEndpoint", AwsGoDependency.SERVICE_INTERNAL_CONFIG).build(),
                "sdkId", sdkId
            ),
            this.commonCodegenArgs
        );
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
                        .addConfigField(EndpointResolverV2)
                        .addConfigFieldResolver(ResolveEndpointResolverV2)
                        .build()
        );
    }
}
