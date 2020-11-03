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
 *
 *
 */

package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

public class DynamoDBValidateResponseChecksum implements GoIntegration {
    private static final String CHECKSUM_CLIENT_OPTION = "DisableValidateResponseChecksum";
    private static final String CHECKSUM_ADDER = "addValidateResponseChecksum";
    private static final String CHECKSUM_INTERNAL_ADDER = "AddValidateResponseChecksum";

    private static final String GZIP_CLIENT_OPTION = "EnableAcceptEncodingGzip";
    private static final String GZIP_ADDER = "addAcceptEncodingGzip";
    private static final String GZIP_INTERNAL_ADDER = "AddAcceptEncodingGzip";

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first.
     *
     * @return Returns the sort order, defaults to -40.
     */
    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        if (!isDynamoDBService(model, settings.getService(model))) {
            return;
        }

        goDelegator.useShapeWriter(settings.getService(model), this::writeMiddlewareHelper);
    }

    private void writeMiddlewareHelper(GoWriter writer) {
        writer.openBlock("func $L(stack *middleware.Stack, options Options) error {", "}", CHECKSUM_ADDER, () -> {
            writer.write("return $T(stack, $T{Disable: options.$L})",
                    SymbolUtils.createValueSymbolBuilder(CHECKSUM_INTERNAL_ADDER,
                            AwsCustomGoDependency.DYNAMODB_CUSTOMIZATION).build(),
                    SymbolUtils.createValueSymbolBuilder(CHECKSUM_INTERNAL_ADDER + "Options",
                            AwsCustomGoDependency.DYNAMODB_CUSTOMIZATION).build(),
                    CHECKSUM_CLIENT_OPTION
            );
        });
        writer.write("");

        writer.openBlock("func $L(stack *middleware.Stack, options Options) error {", "}", GZIP_ADDER, () -> {
            writer.write("return $T(stack, $T{Enable: options.$L})",
                    SymbolUtils.createValueSymbolBuilder(GZIP_INTERNAL_ADDER,
                            AwsCustomGoDependency.ACCEPT_ENCODING_CUSTOMIZATION).build(),
                    SymbolUtils.createValueSymbolBuilder(GZIP_INTERNAL_ADDER + "Options",
                            AwsCustomGoDependency.ACCEPT_ENCODING_CUSTOMIZATION).build(),
                    GZIP_CLIENT_OPTION
            );
        });
        writer.write("");
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                // Add DynamoDB Checksum customization middleware to deserialize.
                RuntimeClientPlugin.builder()
                        .servicePredicate(DynamoDBValidateResponseChecksum::isDynamoDBService)
                        .configFields(ListUtils.of(
                                ConfigField.builder()
                                        .name(CHECKSUM_CLIENT_OPTION)
                                        .type(SymbolUtils.createValueSymbolBuilder("bool")
                                                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                                .build())
                                        .documentation("Allows you to disable the client's validation of "
                                                + "response integrity using CRC32 checksum. Enabled by default.")
                                        .build()
                        ))
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(CHECKSUM_ADDER)
                                        .build())
                                .useClientOptions()
                                .build()
                        )
                        .build(),

                // Add DynamoDB explicit control over accept-encoding: gzip.
                RuntimeClientPlugin.builder()
                        .servicePredicate(DynamoDBValidateResponseChecksum::isDynamoDBService)
                        .configFields(ListUtils.of(
                                ConfigField.builder()
                                        .name(GZIP_CLIENT_OPTION)
                                        .type(SymbolUtils.createValueSymbolBuilder("bool")
                                                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                                .build())
                                        .documentation("Allows you to enable the client's support for "
                                                + "compressed gzip responses. Disabled by default.")
                                        .build()
                        ))
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(GZIP_ADDER)
                                        .build())
                                .useClientOptions()
                                .build()
                        )
                        .build()
        );
    }

    private static boolean isDynamoDBService(Model model, ServiceShape service) {
        return service.expectTrait(ServiceTrait.class).getSdkId().equalsIgnoreCase("DynamoDB");
    }
}
