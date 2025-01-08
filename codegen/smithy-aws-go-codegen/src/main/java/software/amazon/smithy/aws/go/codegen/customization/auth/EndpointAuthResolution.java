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

package software.amazon.smithy.aws.go.codegen.customization.auth;

import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

/**
 * Legacy auth resolution for SigV4A services that choose to do so through endpoint resolution.
 */
public class EndpointAuthResolution implements GoIntegration {
    private static final ConfigFieldResolver finalizeServiceEndpointAuthResolver = ConfigFieldResolver.builder()
            .resolver(Symbol.builder().name("finalizeServiceEndpointAuthResolver").build())
            .location(ConfigFieldResolver.Location.CLIENT)
            .target(ConfigFieldResolver.Target.FINALIZATION)
            .build();

    private static final ConfigFieldResolver finalizeOperationEndpointAuthResolver = ConfigFieldResolver.builder()
            .resolver(Symbol.builder().name("finalizeOperationEndpointAuthResolver").build())
            .location(ConfigFieldResolver.Location.OPERATION)
            .target(ConfigFieldResolver.Target.FINALIZATION)
            .build();

    public static boolean isEndpointAuthService(Model model, ServiceShape service) {
        final String sdkId = service.expectTrait(ServiceTrait.class).getSdkId();
        return sdkId.equalsIgnoreCase("s3")
                || sdkId.equalsIgnoreCase("eventbridge")
                || sdkId.equalsIgnoreCase("sesv2");
    };

    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(EndpointAuthResolution::isEndpointAuthService)
                        .addConfigFieldResolver(finalizeServiceEndpointAuthResolver)
                        .addConfigFieldResolver(finalizeOperationEndpointAuthResolver)
                        .build()
        );
    }
}
