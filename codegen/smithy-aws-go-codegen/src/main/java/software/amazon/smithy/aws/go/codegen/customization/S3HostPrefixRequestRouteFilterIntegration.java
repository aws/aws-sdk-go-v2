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

package software.amazon.smithy.aws.go.codegen.customization;

import java.util.logging.Logger;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.pattern.SmithyPattern.Segment;
import software.amazon.smithy.model.traits.EndpointTrait;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.utils.SmithyInternalApi;

@SmithyInternalApi
public class S3HostPrefixRequestRouteFilterIntegration implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(S3HostPrefixRequestRouteFilterIntegration.class.getName());

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        if (!S3ModelUtils.isServiceS3(model, settings.getService(model))) {
            return model;
        }
        LOGGER.info("Filtering S3 Endpoint Traits with RequestRoute host prefixes");
        return ModelTransformer.create().removeTraitsIf(model, (shape, trait) -> {
            if (trait instanceof EndpointTrait) {
                EndpointTrait endpointTrait = (EndpointTrait) trait;
                for (Segment segment : endpointTrait.getHostPrefix().getLabels()) {
                    if (segment.isLabel() && segment.getContent().equals("RequestRoute")) {
                        LOGGER.info("Removing Endpoint Trait with RequestRoute host prefix: " + endpointTrait);
                        return true;
                    }
                }
            }
            return false;
        });
    }
}
