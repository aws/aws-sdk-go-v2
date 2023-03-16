/*
 * Copyright 2022 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

/**
 * Records the initial resolve DefaultsMode when the client is constructed.
 */
public class ClientResolvedDefaultsMode implements GoIntegration {
    public static final String RESOLVED_DEFAULTS_MODE_CONFIG_NAME = "resolvedDefaultsMode";

    private static final String RESOLVE_RESOLVED_DEFAULTS_MODE = "setResolvedDefaultsMode";

    private static final ConfigField CONFIG_FIELD = ConfigField.builder()
            .name(RESOLVED_DEFAULTS_MODE_CONFIG_NAME)
            .type(SymbolUtils.createValueSymbolBuilder("DefaultsMode", AwsGoDependency.AWS_CORE)
                    .build())
            .documentation("""
                    The initial DefaultsMode used when the client options were constructed. If the
                    DefaultsMode was set to aws.DefaultsModeAuto this will store what the resolved value
                    was at that point in time.

                    Currently does not support per operation call overrides, may in the future.
                    """)
            .build();

    private void writeSetResolvedDefaultsMode(GoWriter writer) {
        writer.pushState();

        writer.putContext("resolverName", RESOLVE_RESOLVED_DEFAULTS_MODE);
        writer.putContext("resolvedOption", RESOLVED_DEFAULTS_MODE_CONFIG_NAME);
        writer.putContext("modeType", SymbolUtils.createValueSymbolBuilder("DefaultsMode",
                AwsGoDependency.AWS_CORE)
                .build());
        writer.putContext("modeOption", AddAwsConfigFields.DEFAULTS_MODE_CONFIG_NAME);
        writer.putContext("autoResolve", SymbolUtils.createValueSymbolBuilder("ResolveDefaultsModeAuto",
                AwsGoDependency.AWS_DEFAULTS)
                .build());
        writer.putContext("autoMode", SymbolUtils.createValueSymbolBuilder("DefaultsModeAuto",
                AwsGoDependency.AWS_CORE)
                .build());
        writer.putContext("region", "Region");
        writer.putContext("envOption", AddAwsConfigFields.RUNTIME_ENVIRONMENT_CONFIG_NAME);

        writer.write("""
                func $resolverName:L(o *Options) {
                    if len(o.$resolvedOption:L) > 0 {
                        return
                    }

                    var mode $modeType:T
                    mode.SetFromString(string(o.$modeOption:L))

                    if mode == $autoMode:T {
                        mode = $autoResolve:T(o.$region:L, o.$envOption:L)
                    }

                    o.$resolvedOption:L = mode
                }
                """);

        writer.popState();
    }

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first.
     *
     * @return Returns the sort order, defaults to -50.
     */
    @Override
    public byte getOrder() {
        return -50;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceShape serviceShape = settings.getService(model);
        goDelegator.useShapeWriter(serviceShape, w -> {
            writeSetResolvedDefaultsMode(w);
        });
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .addConfigField(CONFIG_FIELD)
                        .addConfigFieldResolver(ConfigFieldResolver.builder()
                                .location(ConfigFieldResolver.Location.CLIENT)
                                .target(ConfigFieldResolver.Target.INITIALIZATION)
                                .resolver(SymbolUtils.createValueSymbolBuilder(RESOLVE_RESOLVED_DEFAULTS_MODE).build())
                                .build())
                        .build()
        );
    }
}
