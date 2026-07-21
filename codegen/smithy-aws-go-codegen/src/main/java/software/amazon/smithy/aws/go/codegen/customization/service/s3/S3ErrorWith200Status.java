/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package software.amazon.smithy.aws.go.codegen.customization.service.s3;

import java.util.List;
import java.util.Optional;

import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.traits.HttpPayloadTrait;
import software.amazon.smithy.model.traits.StreamingTrait;

/**
 * Adds middleware to handle S3 response errors with 200 ok status code.
 *
 * Per internal specification, this customization MUST be applied to all S3
 * operations with a structured XML response. The response MUST NOT have a
 * streaming binary payload or event-stream. In Smithy terms, this means all
 * S3 operations whose output does not contain a member with @httpPayload
 * targeting a @streaming blob, @streaming union, or string.
 */
public class S3ErrorWith200Status implements GoIntegration {
    private static String ADD_ERROR_HANDLER_INTERNAL = "HandleResponseErrorWith200Status";

    @Override
    public byte getOrder() {
        // The associated customization ordering is relative to operation deserializers
        // and thus the integration should be added at the end.
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return List.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate(S3ErrorWith200Status::supports200Error)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ADD_ERROR_HANDLER_INTERNAL,
                                        AwsCustomGoDependency.S3_CUSTOMIZATION).build())
                                .build())
                        .build()
        );
    }

    /**
     * Returns true if the operation supports error response with 200 ok status code.
     *
     * Per internal specification, this applies to all S3 operations whose output
     * does NOT contain an @httpPayload member targeting either:
     *   - a shape with the @streaming trait (blob or event stream)
     *   - a string shape
     */
    private static boolean supports200Error(Model model, ServiceShape service, OperationShape operation) {
        if (!isS3Service(model, service)) {
            return false;
        }

        Optional<ShapeId> output = operation.getOutput();
        if (!output.isPresent()) {
            // No output structure means a structured (empty) XML response, apply the check.
            return true;
        }

        StructureShape outputShape = model.expectShape(output.get(), StructureShape.class);

        // Exclude if any @httpPayload member targets a @streaming shape or a string.
        return outputShape.getAllMembers().values().stream()
                .filter(member -> member.hasTrait(HttpPayloadTrait.class))
                .map(member -> model.expectShape(member.getTarget()))
                .noneMatch(target -> target.hasTrait(StreamingTrait.class) || target.isStringShape());
    }

    // returns true if service is s3
    private static boolean isS3Service(Model model, ServiceShape service) {
        return S3ModelUtils.isServiceS3(model, service);
    }
}
