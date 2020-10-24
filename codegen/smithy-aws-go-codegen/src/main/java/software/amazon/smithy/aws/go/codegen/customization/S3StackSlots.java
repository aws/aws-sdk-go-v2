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

package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import software.amazon.smithy.aws.go.codegen.AwsSlotUtils;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.GoDependency;
import software.amazon.smithy.go.codegen.MiddlewareIdentifier;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.integration.StackSlotRegistrar;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

public class S3StackSlots implements GoIntegration {
    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3StackSlots::isS3OrS3Control)
                        .registerStackSlots(StackSlotRegistrar.builder()
                                .addSerializeSlotMutator(AwsSlotUtils.insertAfter(
                                        MiddlewareIdentifier.symbol(ProtocolUtils.OPERATION_SERIALIZER_MIDDLEWARE_ID),
                                        ListUtils.of(
                                                MiddlewareIdentifier.symbol(
                                                        s3SharedValue("EnableDualstackMiddlewareID"))
                                        )
                                ))
                                .addDeserializeSlotMutators(AwsSlotUtils.insertBefore(
                                        MiddlewareIdentifier.symbol(ProtocolUtils.OPERATION_DESERIALIZER_MIDDLEWARE_ID),
                                        ListUtils.of(
                                                MiddlewareIdentifier.symbol(
                                                        s3SharedValue("MetadataRetrieverMiddlewareID"))
                                        )
                                ))
                                .build())
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3StackSlots::isS3)
                        .registerStackSlots(StackSlotRegistrar.builder()
                                .addSerializeSlotMutator(AwsSlotUtils.insertAfter(
                                        MiddlewareIdentifier.symbol(s3SharedValue("EnableDualstackMiddlewareID")),
                                        ListUtils.of(MiddlewareIdentifier.symbol(
                                                        s3CustomizationValue("UpdateEndpointMiddlewareID")))
                                ))
                                .addDeserializeSlotMutators(AwsSlotUtils.insertAfter(
                                        MiddlewareIdentifier.symbol(ProtocolUtils.OPERATION_DESERIALIZER_MIDDLEWARE_ID),
                                        ListUtils.of(
                                                MiddlewareIdentifier.symbol(
                                                        s3CustomizationValue("Process200ErrorMiddlewareID")),
                                                MiddlewareIdentifier.symbol(
                                                        acceptEncodingValue("DecompressGzipMiddlewareID"))
                                        )
                                ))
                                .addFinalizeSlotMutators(AwsSlotUtils.addBefore(ListUtils.of(
                                        MiddlewareIdentifier.symbol(acceptEncodingValue("EnableGzipMiddlewareID")),
                                        MiddlewareIdentifier.symbol(acceptEncodingValue("DisableGzipMiddlewareID"))
                                )))
                                .build())
                        .build());
    }

    private Symbol acceptEncodingValue(String name) {
        return valueSymbol(name, AwsCustomGoDependency.ACCEPT_ENCODING_CUSTOMIZATION);
    }

    private Symbol s3CustomizationValue(String name) {
        return valueSymbol(name, AwsCustomGoDependency.S3_CUSTOMIZATION);
    }

    private Symbol s3SharedValue(String name) {
        return valueSymbol(name, AwsCustomGoDependency.S3_SHARED_CUSTOMIZATION);
    }

    private Symbol valueSymbol(String name, GoDependency dependency) {
        return SymbolUtils.createValueSymbolBuilder(name, dependency).build();
    }

    // returns true if service is either s3 or s3 control and needs s3 customization
    private static boolean isS3OrS3Control(Model model, ServiceShape service) {
        return isS3(model, service) || isS3Control(model, service);
    }

    // returns true if service is either s3 control and needs s3 customization
    private static boolean isS3(Model model, ServiceShape service) {
        String serviceId = service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("S3");
    }

    // returns true if service is either s3 control and needs s3 customization
    private static boolean isS3Control(Model model, ServiceShape service) {
        String serviceId = service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("S3 Control");
    }
}
