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
 */

package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;

/**
 * Backfills the AWS service trait onto smithy-namespaced protocol tests. Long-term these protocol tests should live in
 * smithy-go instead but they're currently stuck here due to conflated codegen nonsense.
 */
public class BackfillProtocolTestServiceTrait implements GoIntegration {
    @Override
    public byte getOrder() {
        return -128;
    }

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        var service = settings.getService(model);
        if (!service.getId().getNamespace().startsWith("smithy.protocoltests")) {
            return model;
        }

        return model.toBuilder()
                .addShape(
                        service.toBuilder()
                                .addTrait(
                                        ServiceTrait.builder()
                                                .sdkId("")
                                                .arnNamespace("")
                                                .cloudFormationName("")
                                                .cloudTrailEventSource("")
                                                .endpointPrefix("")
                                                .build()
                                )
                                .build()
                )
                .build();
    }
}
