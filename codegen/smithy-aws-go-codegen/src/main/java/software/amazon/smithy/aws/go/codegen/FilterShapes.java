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

package software.amazon.smithy.aws.go.codegen;

import java.util.Optional;
import java.util.Set;
import java.util.logging.Logger;
import java.util.stream.Collectors;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.utils.SetUtils;

/**
 * Filters out certain shapes such as an operation.
 */
public final class FilterShapes implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(FilterShapes.class.getName());

    private static final Set<ShapeId> SHAPE_IDS = SetUtils.of();

    public FilterShapes() {
    }

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        var toRemove = SHAPE_IDS.stream()
                .map(model::getShape)
                .filter(Optional::isPresent)
                .map(Optional::get)
                .collect(Collectors.toSet());

        if (toRemove.size() == 0) {
            return model;
        }

        ModelTransformer transformer = ModelTransformer.create();

        return transformer.removeUnreferencedShapes(transformer.removeShapes(model, toRemove));
    }
}
