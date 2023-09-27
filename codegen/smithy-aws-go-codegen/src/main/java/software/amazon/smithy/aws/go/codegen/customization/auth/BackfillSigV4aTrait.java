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
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.auth.SigV4aTrait;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.traits.AuthTrait;
import software.amazon.smithy.utils.SetUtils;

/**
 * Throws the aws.auth#sigv4a trait onto the service such that auth codegen picks it up.
 */
public class BackfillSigV4aTrait implements GoIntegration {
    private boolean isBackfillSigV4aService(ServiceShape service) {
        final String sdkId = service.expectTrait(ServiceTrait.class).getSdkId();
        return sdkId.equalsIgnoreCase("s3") || sdkId.equalsIgnoreCase("eventbridge");
    };

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        ServiceShape service = settings.getService(model);
        if (!isBackfillSigV4aService(service)) {
            return model;
        }

        return model.toBuilder()
                .addShape(
                        service.toBuilder()
                                .addTrait(new SigV4aTrait())
                                .addTrait(new AuthTrait(SetUtils.of(SigV4Trait.ID, SigV4aTrait.ID)))
                                .build()
                )
                .build();
    }
}
