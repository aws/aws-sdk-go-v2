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
import software.amazon.smithy.aws.go.codegen.AddAwsConfigFields;
import software.amazon.smithy.aws.go.codegen.AddAwsConfigFields.AwsConfigField;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

public class DynamoDBValidateResponseChecksum implements GoIntegration {
    private static final String DDB_CHECKSUM_NAME = "DisableResponseChecksumValidation";

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
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                // Add DynamoDB Checksum customization middleware to deserialize.
                RuntimeClientPlugin.builder()
                        .servicePredicate(DynamoDBValidateResponseChecksum::isDynamoDBService)
                        .configFields(ListUtils.of(
                                AwsConfigField.builder()
                                        .name(DDB_CHECKSUM_NAME)
                                        .type(SymbolUtils.createValueSymbolBuilder("bool")
                                                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                                .build())
                                        .documentation(
                                                DDB_CHECKSUM_NAME + " allows you to disable the client's validation of "
                                                        + "response integrity using CRC32 checksum.")
                                        .build()
                        ))
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddChecksumMiddleware", AwsCustomGoDependency.DYNAMODB_CUSTOMIZATION)
                                        .build())
                                .useClientOptions()
                                .build()
                        )
                        .build(),

                // Add DynamoDB explicit control over accept-encoding: gzip.
                RuntimeClientPlugin.builder()
                        .servicePredicate(DynamoDBValidateResponseChecksum::isDynamoDBService)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        "AddAcceptEncodingGzip", AwsCustomGoDependency.DYNAMODB_CUSTOMIZATION)
                                        .build())
                                .build()
                        )
                        .build()
        );
    }

    private static boolean isDynamoDBService(Model model, ServiceShape service) {
        return service.expectTrait(ServiceTrait.class).getSdkId().equalsIgnoreCase("DynamoDB");
    }
}
