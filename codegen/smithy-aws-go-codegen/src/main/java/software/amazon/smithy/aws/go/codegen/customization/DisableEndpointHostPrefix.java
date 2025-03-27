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
 *
 *
 */

package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import software.amazon.smithy.aws.go.codegen.SdkGoTypes;
import software.amazon.smithy.aws.go.codegen.customization.service.s3.S3ModelUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.traits.EndpointTrait;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.utils.ListUtils;

/**
 * Adds a middleware to disable host prefix serialization for the AccountId input in S3Control operations if v2 endpoint
 * resolution was used (endpoint rules handle that host prefixing directly).
 */
public class DisableEndpointHostPrefix implements GoIntegration {
    private static final MiddlewareRegistrar DisableHostPrefixMiddleware = MiddlewareRegistrar.builder()
            .resolvedFunction(SdkGoTypes.ServiceCustomizations.S3Control.AddDisableHostPrefixMiddleware)
            .build();

    private static boolean hasAccountIdEndpointPrefix(Model model, ServiceShape service, OperationShape operation) {
        return operation.hasTrait(EndpointTrait.class) &&
                operation.expectTrait(EndpointTrait.class).getHostPrefix().getLabels().stream().anyMatch(it ->
                        it.isLabel() && it.getContent().equals("AccountId"));
    }

    // must be inserted after EndpointResolutionV2
    @Override
    public byte getOrder() { return 127; }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3ModelUtils::isServiceS3Control)
                        .operationPredicate(DisableEndpointHostPrefix::hasAccountIdEndpointPrefix)
                        .registerMiddleware(DisableHostPrefixMiddleware)
                        .build()
        );
    }
}
