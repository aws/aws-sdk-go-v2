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
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.BoxTrait;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.utils.SetUtils;

/**
 * Integration that back fills the `boxed` traits to API shapes that were not decorated with the trait in the model.
 */
public class BackfillBoxTrait implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(BackfillBoxTrait.class.getName());

    /**
     * Map of service shape to Set of operation shapes that need to have this
     * presigned url auto fill customization.
     */
    public static final Map<ShapeId, Set<ShapeId>> SERVICE_TO_MEMBER_MAP = MapUtils.of(
            ShapeId.from("com.amazonaws.s3control#AWSS3ControlServiceV20180820"), SetUtils.of(
                    ShapeId.from("com.amazonaws.s3control#S3ExpirationInDays")
            ));

    /**
     * /**
     * Updates the API model to add additional members to the operation input shape that are needed for presign url
     * customization.
     *
     * @param model    API model
     * @param settings Go codegen settings
     * @return updated API model
     */
    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        ShapeId serviceId = settings.getService();
        if (!SERVICE_TO_MEMBER_MAP.containsKey(serviceId)) {
            return model;
        }
        Model.Builder builder = model.toBuilder();

        Set<ShapeId> shapeIds = SERVICE_TO_MEMBER_MAP.get(serviceId);
        for (ShapeId shapeId : shapeIds) {
            IntegerShape shape = model.expectShape(shapeId, IntegerShape.class);
            if (shape.getTrait(BoxTrait.class).isPresent()) {
                LOGGER.warning("BoxTrait is present in model and does not require backfill");
                continue;
            }
            builder.addShape(shape.toBuilder()
                    .addTrait(new BoxTrait())
                    .build());
        }

        return builder.build();
    }
}
