/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package software.amazon.smithy.aws.go.codegen.customization.service.s3;

import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;

import java.util.Map;

/**
 * Stores the endpoint resolved by EndpointResolverV2
 */
public class StoreResolvedUri implements GoIntegration {
    @Override
    public void renderPostEndpointResolutionHook(GoSettings settings, GoWriter writer, Model model) {
        if (!S3ModelUtils.isServiceS3(model, settings.getService(model))) return;
        writer.writeGoTemplate(
                """
                        ctx = $setFunc:L(ctx, endpt.URI.String())
                        """,
                Map.of(
                        "setFunc", "setS3ResolvedURI")
        );
    }
}
