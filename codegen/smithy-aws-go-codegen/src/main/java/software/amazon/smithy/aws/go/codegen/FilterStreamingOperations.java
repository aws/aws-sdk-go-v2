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
 */

package software.amazon.smithy.aws.go.codegen;

import java.util.Set;
import java.util.logging.Logger;
import java.util.stream.Collectors;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.EventStreamIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.transform.ModelTransformer;

/**
 * Filters out event stream operations.
 * TODO: implement event streams
 */
public final class FilterStreamingOperations implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(FilterStreamingOperations.class.getName());

    public FilterStreamingOperations() {}

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        EventStreamIndex index = EventStreamIndex.of(model);
        Set<Shape> streamingOperations = model.shapes(OperationShape.class)
                .filter(op -> index.getOutputInfo(op).isPresent() || index.getInputInfo(op).isPresent())
                .peek(op -> LOGGER.warning(String.format(
                        "Filtering out unsupported event stream operation: %s", op.getId().toString())))
                .collect(Collectors.toSet());
        ModelTransformer transformer = ModelTransformer.create();
        return transformer.removeUnreferencedShapes(transformer.removeShapes(model, streamingOperations));
    }
}
