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
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.SetUtils;

/**
 * Generates an endpoint resolver from endpoints.json.
 */
public final class AwsEndpointGenerator implements GoIntegration {

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {
        new EndpointGenerator(settings, model, writerFactory).run();
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(RuntimeClientPlugin.builder()
                .configFields(SetUtils.of(
                        ConfigField.builder()
                                .name("EndpointResolver")
                                .type(SymbolUtils.createValueSymbolBuilder(EndpointGenerator.RESOLVER_INTERFACE_NAME)
                                        .build())
                                .documentation("The service endpoint resolver.")
                                .build(),
                        ConfigField.builder()
                                .name("EndpointOptions")
                                .type(SymbolUtils.createValueSymbolBuilder(EndpointGenerator.RESOLVER_OPTIONS)
                                        .build())
                                .documentation("The endpoint options to be used when attempting "
                                        + "to resolve an endpoint.")
                                .build()
                ))
                .resolveFunction(SymbolUtils.createValueSymbolBuilder(EndpointGenerator.CLIENT_CONFIG_RESOLVER)
                        .build())
                .build());
    }
}
