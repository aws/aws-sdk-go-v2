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

package software.amazon.smithy.aws.go.codegen.customization;

import java.util.Map;
import java.util.Set;
import java.util.logging.Logger;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.OptionalAuthTrait;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.utils.SetUtils;

/**
 * Backfill missing Smithy OptionalAuth traits to AWS models.
 */
public class BackfillOptionalAuthTrait implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(BackfillOptionalAuthTrait.class.getName());

    private static final Map<ShapeId, Set<ShapeId>> SERVICE_TO_OPERATION_MAP = MapUtils.of(
            ShapeId.from("com.amazonaws.sts#AWSSecurityTokenServiceV20110615"), SetUtils.of(
                    ShapeId.from("com.amazonaws.sts#AssumeRoleWithSAML"),
                    ShapeId.from("com.amazonaws.sts#AssumeRoleWithWebIdentity")),
            ShapeId.from("com.amazonaws.cognitoidentity#AWSCognitoIdentityService"), SetUtils.of(
                    ShapeId.from("com.amazonaws.cognitoidentity#GetId"),
                    ShapeId.from("com.amazonaws.cognitoidentity#GetOpenIdToken"),
                    ShapeId.from("com.amazonaws.cognitoidentity#UnlinkIdentity"),
                    ShapeId.from("com.amazonaws.cognitoidentity#GetCredentialsForIdentity")));

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
        if (!SERVICE_TO_OPERATION_MAP.containsKey(serviceId)) {
            return model;
        }

        Set<ShapeId> operationIds = SERVICE_TO_OPERATION_MAP.get(serviceId);

        Model.Builder builder = model.toBuilder();
        for (ShapeId operationId : operationIds) {
            OperationShape operationShape = model.expectShape(operationId).asOperationShape().get();
            if (operationShape.getTrait(OptionalAuthTrait.class).isPresent()) {
                LOGGER.warning("optionalAuth trait is present in model and does not require backfill");
                continue;
            }
            builder.addShape(operationShape.toBuilder()
                    .addTrait(new OptionalAuthTrait())
                    .build());
        }

        return builder.build();
    }
}
