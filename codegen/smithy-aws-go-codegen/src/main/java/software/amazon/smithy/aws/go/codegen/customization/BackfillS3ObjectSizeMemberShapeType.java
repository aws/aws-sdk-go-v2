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
 *
 *
 */

package software.amazon.smithy.aws.go.codegen.customization;

import java.util.Map;
import java.util.Set;
import java.util.logging.Logger;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.IntegerShape;
import software.amazon.smithy.model.shapes.LongShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.ShapeType;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.utils.OptionalUtils;
import software.amazon.smithy.utils.SetUtils;

public class BackfillS3ObjectSizeMemberShapeType  implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(BackfillS3ObjectSizeMemberShapeType.class.getName());

    private static final Map<ShapeId, Set<ShapeId>> SERVICE_TO_SHAPE_MAP = MapUtils.of(
            ShapeId.from("com.amazonaws.s3#AmazonS3"), SetUtils.of(
                    ShapeId.from("com.amazonaws.s3#Size")
            )
    );

    @Override
    public byte getOrder() {
        // This integration should happen before other integrations that rely on the presence of this trait
        return -60;
    }

    @Override
    public Model preprocessModel(
            Model model, GoSettings settings
    ) {
        ShapeId serviceId = settings.getService();
        if (!SERVICE_TO_SHAPE_MAP.containsKey(serviceId)) {
            return model;
        }

        Set<ShapeId> shapeIds = SERVICE_TO_SHAPE_MAP.get(serviceId);

        Model.Builder builder = model.toBuilder();
        for (ShapeId shapeId : shapeIds) {
            OptionalUtils.ifPresentOrElse(
                    model.getShape(shapeId),
                    (shape) -> {
                        if (shape.isLongShape()) {
                            LOGGER.warning("shape is already modeled as an LONG does not require backfill");
                            return;
                        }

                        builder.addShape(LongShape.builder()
                                .id(shape.getId())
                                .addTraits(shape.getAllTraits().values())
                                .build());
                    },
                    () -> LOGGER.warning("shape " + shapeId + " not found in " + serviceId + " model")
            );
        }

        return builder.build();
    }
}
