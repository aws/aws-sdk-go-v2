/*
 * Copyright 2021 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

import java.util.Map;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.ShapeType;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.utils.MapUtils;

public class ShapeBackwardsCompatabilityTransforms implements GoIntegration {
    private static final Map<ShapeId, ShapeType> SHAPE_ID_SHAPE_TYPE_MAP = MapUtils.of(
            ShapeId.from("com.amazonaws.nimble#StudioComponentConfiguration"), ShapeType.STRUCTURE
    );

    @Override
    public byte getOrder() {
        return -127;
    }

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        return ModelTransformer.create()
                .changeShapeType(model, SHAPE_ID_SHAPE_TYPE_MAP);
    }
}
