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
import software.amazon.smithy.go.codegen.GoDependency;
import software.amazon.smithy.go.codegen.MiddlewareIdentifier;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.integration.StackSlotRegistrar;
import software.amazon.smithy.utils.ListUtils;

public class AwsStackSlots implements GoIntegration {
    @Override
    public byte getOrder() {
        return -126;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .registerStackSlots(StackSlotRegistrar.builder()
                                .addInitializeSlotMutator(AwsSlotUtils.addBefore(ListUtils.of(
                                        AwsSlotUtils.awsSymbolId("RegisterServiceMetadata")
                                )))
                                .addSerializeSlotMutator(AwsSlotUtils.insertBefore(
                                        MiddlewareIdentifier.symbol(ProtocolUtils.OPERATION_SERIALIZER_MIDDLEWARE_ID),
                                        ListUtils.of(
                                                AwsSlotUtils.awsSymbolId("ResolveEndpoint")
                                        )))
                                .addBuildSlotMutator(AwsSlotUtils.addAfter(ListUtils.of(
                                        AwsSlotUtils.awsSymbolId("ClientRequestID"),
                                        AwsSlotUtils.awsSymbolId("ComputePayloadHash"),
                                        AwsSlotUtils.awsSymbolId("UserAgent")
                                )))
                                .addFinalizeSlotMutators(AwsSlotUtils.addAfter(ListUtils.of(
                                        AwsSlotUtils.awsSymbolId("Retry"),
                                        AwsSlotUtils.awsSymbolId("Signing")
                                )))
                                .addDeserializeSlotMutators(AwsSlotUtils.insertBefore(
                                        MiddlewareIdentifier.symbol(ProtocolUtils.OPERATION_DESERIALIZER_MIDDLEWARE_ID),
                                        ListUtils.of(
                                                AwsSlotUtils.awsSymbolId("ResponseErrorWrapper"),
                                                AwsSlotUtils.awsSymbolId("RequestIDRetriever")
                                        )))
                                .build())
                        .build());
    }
}
