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

import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;

public class AwsSdkServiceId implements GoIntegration {
    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public String processServiceId(GoSettings settings, Model model, String serviceId) {
        ServiceShape serviceShape = settings.getService(model);
        return serviceShape.expectTrait(ServiceTrait.class).getSdkId();
    }
}
