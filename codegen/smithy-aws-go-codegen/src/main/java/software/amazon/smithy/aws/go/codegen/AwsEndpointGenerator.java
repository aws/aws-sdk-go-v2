/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

import java.util.List;
import java.util.function.Consumer;
import software.amazon.smithy.aws.go.codegen.customization.S3ModelUtils;
import software.amazon.smithy.aws.traits.ServiceTrait;
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
import software.amazon.smithy.utils.SetUtils;

/**
 * Generates an endpoint resolver from endpoints.json.
 */
public final class AwsEndpointGenerator implements GoIntegration {
    public static final String ENDPOINT_RESOLVER_CONFIG_NAME = "EndpointResolver";
    public static final String ENDPOINT_OPTIONS_CONFIG_NAME = "EndpointOptions";

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {
        String serviceId = settings.getService(model).expectTrait(ServiceTrait.class).getSdkId();
        boolean generateQueryHelpers = serviceId.equalsIgnoreCase("S3")
                                       || serviceId.equalsIgnoreCase("EventBridge");

        EndpointGenerator.builder()
                .settings(settings)
                .model(model)
                .writerFactory(writerFactory)
                .modelQueryHelpers(generateQueryHelpers)
                .build()
                .run();
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .configFields(SetUtils.of(
                                ConfigField.builder()
                                        .name(ENDPOINT_RESOLVER_CONFIG_NAME)
                                        .type(SymbolUtils.createValueSymbolBuilder(EndpointGenerator.RESOLVER_INTERFACE_NAME)
                                                .build())
                                        .documentation(String.format("The service endpoint resolver."))
                                        .deprecated(String.format(
                                            """
                                            %s and With%s are deprecated. See %s and With%s
                                            """,
                                            EndpointGenerator.RESOLVER_INTERFACE_NAME,
                                            EndpointGenerator.RESOLVER_INTERFACE_NAME,
                                            EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME,
                                            EndpointResolutionGenerator.RESOLVER_INTERFACE_NAME
                                        ))
                                        .withHelper(true)
                                        .build(),
                                ConfigField.builder()
                                        .name(ENDPOINT_OPTIONS_CONFIG_NAME)
                                        .type(SymbolUtils.createValueSymbolBuilder(EndpointGenerator.RESOLVER_OPTIONS)
                                                .build())
                                        .documentation("The endpoint options to be used when attempting "
                                                       + "to resolve an endpoint.")
                                        .build()
                        ))
                        .addConfigFieldResolver(ConfigFieldResolver.builder()
                                .location(ConfigFieldResolver.Location.OPERATION)
                                .target(ConfigFieldResolver.Target.FINALIZATION)
                                .resolver(SymbolUtils.createValueSymbolBuilder(
                                        EndpointGenerator.FINALIZE_CLIENT_ENDPOINT_RESOLVER_OPTIONS).build())
                                .build())
                        .build());
    }
}
