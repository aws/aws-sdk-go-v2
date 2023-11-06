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

import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
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

public class AwsEndpointResolverInitializerGenerator implements GoIntegration {

    public static final String RESOLVE_ENDPOINT_RESOLVER_V2 = "resolve" + EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME;

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
                        .addConfigField(EndpointResolverV2)
                        .addConfigFieldResolver(ResolveEndpointResolverV2)
                        .build()
        );
    }
}
