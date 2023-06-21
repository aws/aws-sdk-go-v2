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

import static org.junit.jupiter.api.Assertions.assertTrue;

import org.junit.jupiter.api.Test;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.EndpointTrait;

public class S3ControlHostPrefixAccountIdFilterIntegrationTest {
    @Test
    public void test() {
        Model model = TestUtils.preprocessModelIntegration(
                new S3ControlHostPrefixAccountIdFilterIntegration(),
                S3ModelUtils.SERVICE_S3_CONTROL_MODEL_FILE);
        OperationShape operation = model.expectShape(
                ShapeId.from("com.amazonaws.s3control#PutAccessPointConfigurationForObjectLambda"),
                OperationShape.class);
        assertTrue(operation.getTrait(EndpointTrait.class).isEmpty());
    }
}
