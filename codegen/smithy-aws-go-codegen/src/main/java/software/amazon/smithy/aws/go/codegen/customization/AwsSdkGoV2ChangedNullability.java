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

import static software.amazon.smithy.aws.go.codegen.customization.ApiGatewayExportsNullabilityExceptionIntegration.getNullabilityExceptions;

import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.diff.Differences;
import software.amazon.smithy.diff.evaluators.AbstractDiffEvaluator;
import software.amazon.smithy.go.codegen.knowledge.GoPointableIndex;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.NullableIndex;
import software.amazon.smithy.model.shapes.NumberShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.DefaultTrait;
import software.amazon.smithy.model.validation.ValidationEvent;

/**
 * An aws-sdk-go-v2 specific ChangedNullability diff evaluator.
 *
 * The following conditions will create events:
 * - Added shapes with a default trait that are not in the nullability
 *   exceptions.
 *   - TODO(APIGW): remove after model fixes
 * - Any changed shape that differs in nullability in GoPointableIndex and
 *   NullableIndex
 */
public class AwsSdkGoV2ChangedNullability extends AbstractDiffEvaluator {
    @Override
    public List<ValidationEvent> evaluate(Differences differences) {
        Model oldModel = differences.getOldModel();
        GoPointableIndex oldGoPointableIndex = GoPointableIndex.of(oldModel);
        NullableIndex oldNullableIndex = NullableIndex.of(oldModel);

        Model newModel = differences.getNewModel();
        GoPointableIndex newGoPointableIndex = GoPointableIndex.of(newModel);
        NullableIndex newNullableIndex = NullableIndex.of(newModel);
        Set<NumberShape> newNumberShapes = newModel.toSet(NumberShape.class);

        List<ValidationEvent> events = new ArrayList<ValidationEvent>();

        ShapeId service = newModel
                .getServiceShapesWithTrait(ServiceTrait.class)
                .iterator().next()
                .toShapeId();
        Set<ShapeId> nullabilityExceptions = getNullabilityExceptions(service);

        differences.addedShapes().forEach(shape -> {
            // TODO(APIGW): remove after model fixes
            if (shape.hasTrait(DefaultTrait.class)
                    && (shape.isBooleanShape() || newNumberShapes.contains(shape))
                    && !nullabilityExceptions.contains(shape.toShapeId())) {
                events.add(error(shape, "Shape must be added to the APIGW exports nullability exceptions"));
            }
        });

        differences.changedShapes().forEach(changedShape -> {
            ShapeId shape = changedShape.getShapeId();
            if (oldGoPointableIndex.isDereferencable(shape) != newGoPointableIndex.isDereferencable(shape)) {
                events.add(error(changedShape.getNewShape(), "Shape changed GoPointableIndex::isDereferencable()"));
            }
            if (oldGoPointableIndex.isNillable(shape) != newGoPointableIndex.isNillable(shape)) {
                events.add(error(changedShape.getNewShape(), "Shape changed GoPointableIndex::isNillable()"));
            }
            if (oldGoPointableIndex.isPointable(shape) != newGoPointableIndex.isPointable(shape)) {
                events.add(error(changedShape.getNewShape(), "Shape changed GoPointableIndex::isPointable()"));
            }
            if (oldNullableIndex.isNullable(shape) != newNullableIndex.isNullable(shape)) {
                events.add(error(changedShape.getNewShape(), "Shape changed NullableIndex::isNullable()"));
            }
            if (changedShape.getOldShape().isMemberShape() && changedShape.getNewShape().isMemberShape()) {
                boolean isOldMemberNullable = oldNullableIndex.isMemberNullable(
                        changedShape.getOldShape().asMemberShape().get(),
                        NullableIndex.CheckMode.CLIENT_ZERO_VALUE_V1_NO_INPUT);
                boolean isNewMemberNullable = newNullableIndex.isMemberNullable(
                        changedShape.getNewShape().asMemberShape().get(),
                        NullableIndex.CheckMode.CLIENT_ZERO_VALUE_V1_NO_INPUT);
                if (isOldMemberNullable != isNewMemberNullable) {
                    events.add(error(changedShape.getNewShape(), "Shape changed NullableIndex::isMemberNullable()"));
                }
            }
            if (changedShape.getOldShape().isMemberShape() && !changedShape.getNewShape().isMemberShape()) {
                events.add(error(changedShape.getNewShape(), "Shape is not a member shape"));
            }
            if (!changedShape.getOldShape().isMemberShape() && changedShape.getNewShape().isMemberShape()) {
                events.add(error(changedShape.getNewShape(), "Shape should not be a member shape"));
            }
        });
        return events;
    }
}
