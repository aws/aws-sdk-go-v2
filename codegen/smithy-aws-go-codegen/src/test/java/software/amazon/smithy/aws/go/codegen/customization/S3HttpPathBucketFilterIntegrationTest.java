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

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;

import org.junit.jupiter.api.Test;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.HttpTrait;

public class S3HttpPathBucketFilterIntegrationTest {
    @Test
    public void test() {
        Model model = TestUtils.preprocessModelIntegration(
                new S3HttpPathBucketFilterIntegration(),
                S3ModelUtils.SERVICE_S3_MODEL_FILE);
        OperationShape operation = model.expectShape(
            ShapeId.from("com.amazonaws.s3#DeleteBucketWebsite"),
            OperationShape.class);
        String uri = operation.expectTrait(HttpTrait.class)
            .getUri().toString();
        // URI is originally: /{Bucket}?website
        assertFalse(uri.contains("{Bucket}"));
        assertEquals(uri, "/?website");
    }
}
