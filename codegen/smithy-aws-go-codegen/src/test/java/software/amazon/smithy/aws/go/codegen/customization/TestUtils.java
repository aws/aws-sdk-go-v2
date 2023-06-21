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

import java.nio.file.Path;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.loader.ModelAssembler;
import software.amazon.smithy.model.node.Node;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.utils.IoUtils;

public class TestUtils {
    public static final String AWS_MODELS_PATH_PREFIX = "../sdk-codegen/aws-models/";

    public static Node getAwsModel(String modelFile) {
        try {
            return Node.parseJsonWithComments(IoUtils.readUtf8File(Path.of(AWS_MODELS_PATH_PREFIX, modelFile)));
        } catch (Exception e) {
            throw new CodegenException(e);
        }
    }

    public static Model preprocessModelIntegration(GoIntegration integration, String modelFile) {
        GoSettings settings = new GoSettings();
        Model model = new ModelAssembler()
                .addDocumentNode(getAwsModel(modelFile))
                .disableValidation()
                .putProperty(ModelAssembler.ALLOW_UNKNOWN_TRAITS, true)
                .assemble()
                .unwrap();
        ShapeId service = model.getServiceShapes().stream().findFirst().get().getId();
        settings.setService(service);
        return integration.preprocessModel(model, settings);
    }
}
