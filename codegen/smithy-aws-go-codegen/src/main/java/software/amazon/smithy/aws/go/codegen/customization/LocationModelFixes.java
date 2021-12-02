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

import java.util.List;
import java.util.logging.Logger;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.AuthTrait;
import software.amazon.smithy.utils.ListUtils;

public class LocationModelFixes implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(LocationModelFixes.class.getName());

    private static final List<ShapeId> SHAPE_ID_EMPTY_AUTH_TRAIT_REMOVAL = ListUtils.of(
            ShapeId.from("com.amazonaws.location#BatchEvaluateGeofences"),
            ShapeId.from("com.amazonaws.location#DescribeGeofenceCollection"),
            ShapeId.from("com.amazonaws.location#DescribeMap"),
            ShapeId.from("com.amazonaws.location#DescribePlaceIndex"),
            ShapeId.from("com.amazonaws.location#DescribeRouteCalculator"),
            ShapeId.from("com.amazonaws.location#DescribeTracker")
    );

    @Override
    public Model preprocessModel(
            Model model,
            GoSettings settings
    ) {
        if (SHAPE_ID_EMPTY_AUTH_TRAIT_REMOVAL.size() == 0) {
            return model;
        }

        var builder = model.toBuilder();

        for (ShapeId shapeId : SHAPE_ID_EMPTY_AUTH_TRAIT_REMOVAL) {
            var optionalShape = model.getShape(shapeId);

            if (optionalShape.isEmpty()) {
                continue;
            }

            var shape = optionalShape.get().asOperationShape().get();

            var optionalAuthTrait = shape.getTrait(AuthTrait.class);

            if (optionalAuthTrait.isEmpty()) {
                LOGGER.warning(() -> String.format("%s no longer has an AuthTrait", shapeId));
                continue;
            }

            var authTrait = optionalAuthTrait.get();

            if (authTrait.getValueSet().size() != 0) {
                LOGGER.warning(() -> String.format("%s has a non-empty AuthTrait list and will not be removed",
                        shapeId));
                continue;
            }

            builder.addShape(shape.toBuilder()
                    .removeTrait(AuthTrait.ID)
                    .build());
        }

        return builder.build();
    }
}
