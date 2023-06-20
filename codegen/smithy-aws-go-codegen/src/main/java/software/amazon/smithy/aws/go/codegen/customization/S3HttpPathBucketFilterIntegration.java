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
import java.util.regex.Pattern;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.pattern.UriPattern;
import software.amazon.smithy.model.traits.HttpTrait;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.utils.SmithyInternalApi;

@SmithyInternalApi
public class S3HttpPathBucketFilterIntegration implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(S3HttpPathBucketFilterIntegration.class.getName());

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        if (!S3ModelUtils.isServiceS3(model, settings.getService(model))) {
            return model;
        }
        LOGGER.info("Filtering S3 HTTP Bucket Bindings in URI Paths in Http Traits");
        return ModelTransformer.create().mapTraits(model, (shape, trait) -> {
            if (trait instanceof HttpTrait) {
                HttpTrait httpTrait = (HttpTrait) trait;
                UriPattern uriPattern = UriPattern.parse(Pattern
                        .compile("\\{Bucket}/?")
                        .matcher(httpTrait.getUri().toString())
                        .replaceAll(""));
                LOGGER.info("Replacing URI Path for " + httpTrait + ": "
                        + httpTrait.getUri().toString() + " => " + uriPattern.toString());
                return httpTrait.toBuilder()
                        .uri(uriPattern)
                        .build();
            }
            return trait;
        });
    }
}
