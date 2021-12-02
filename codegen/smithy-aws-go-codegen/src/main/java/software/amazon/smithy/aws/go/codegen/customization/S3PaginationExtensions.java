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

import java.util.Optional;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.trait.PagingExtensionTrait;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.PaginatedIndex;
import software.amazon.smithy.model.knowledge.PaginationInfo;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;

/**
 * This customization adds support for checking the IsTruncated boolean member for paginated S3 operations to determine
 * if the NextToken should be set for the paginator.
 */
public class S3PaginationExtensions implements GoIntegration {
    @Override
    public Model preprocessModel(
            Model model, GoSettings settings
    ) {
        ServiceShape service = settings.getService(model);
        if (!S3ModelUtils.isServiceS3(model, service)) {
            return model;
        }

        return addMoreResultsKey(model, service);
    }

    private Model addMoreResultsKey(Model model, ServiceShape service) {
        PaginatedIndex paginatedIndex = PaginatedIndex.of(model);

        Model.Builder builder = model.toBuilder();
        for (OperationShape operation : TopDownIndex.of(model).getContainedOperations(service)) {
            Optional<PaginationInfo> optionalPaginationInfo = paginatedIndex.getPaginationInfo(service, operation);
            if (!optionalPaginationInfo.isPresent()) {
                continue;
            }
            StructureShape outputShape = optionalPaginationInfo.get().getOutput();

            Optional<MemberShape> memberShape = Optional.empty();
            for (MemberShape member : outputShape.members()) {
                if (member.getTarget().equals(ShapeId.from("com.amazonaws.s3#IsTruncated"))) {
                    if (memberShape.isPresent()) {
                        throw new CodegenException("IsTruncated member present in output more then once");
                    }
                    memberShape = Optional.of(member);
                }
            }

            if (memberShape.isPresent()) {
                OperationShape.Builder operationBuilder = optionalPaginationInfo.get().getOperation().toBuilder();
                operationBuilder.addTrait(PagingExtensionTrait.builder()
                        .moreResults(memberShape.get())
                        .build());
                builder.addShape(operationBuilder.build());
            }
        }

        return builder.build();
    }
}
