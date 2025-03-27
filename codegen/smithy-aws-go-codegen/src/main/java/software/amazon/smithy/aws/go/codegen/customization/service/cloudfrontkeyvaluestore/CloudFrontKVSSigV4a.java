/*
 * Copyright 2024 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package software.amazon.smithy.aws.go.codegen.customization.service.cloudfrontkeyvaluestore;

import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.aws.traits.auth.SigV4ATrait;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.traits.AuthTrait;
import software.amazon.smithy.utils.SetUtils;

/**
 * This integration configures the CloudFront Key Value Store client for Signature Version 4a
 */
public class CloudFrontKVSSigV4a implements GoIntegration {
    // hardcoded from model so we don't have to extract it from whatever auth trait
    private static final String SIGNING_NAME = "cloudfront-keyvaluestore";

    /**
     * Return true if service is CFKVS.
     *
     * @param model   is the generation model.
     * @param service is the service shape being audited.
     */
    private static boolean isCFKVSService(Model model, ServiceShape service) {
        final String sdkId = service.expectTrait(ServiceTrait.class).getSdkId();
        final String serviceId = sdkId.replace("-", "").replace(" ", "").toLowerCase();
        return serviceId.equalsIgnoreCase("cloudfrontkeyvaluestore");
    }

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        ServiceShape service = settings.getService(model);
        if (!isCFKVSService(model, service)) {
            return model;
        }

        // we MUST preserve the sigv4 trait as released since it affects the exported API
        // (signer interface and config field)
        return model.toBuilder()
                .addShape(
                        service.toBuilder()
                                .addTrait(SigV4ATrait.builder().name(SIGNING_NAME).build())
                                .addTrait(SigV4Trait.builder().name(SIGNING_NAME).build())
                                .addTrait(new AuthTrait(SetUtils.of(SigV4ATrait.ID, SigV4Trait.ID)))
                                .build()
                )
                .build();
    }
}
