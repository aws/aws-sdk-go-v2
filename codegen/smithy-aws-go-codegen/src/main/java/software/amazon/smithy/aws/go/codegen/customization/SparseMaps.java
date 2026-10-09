/*
 * Copyright 2026 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.MapShape;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.SparseTrait;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.utils.MapUtils;

/**
 * Marks map members as sparse for services that return explicit null map values despite modeling the map as
 * non-sparse.
 *
 * <p>Because {@code @sparse} is a trait on the map shape, each affected member is retargeted to a sparse copy of its
 * map so that other members sharing the original map are unaffected.
 */
public class SparseMaps implements GoIntegration {
    private static final String SPARSE_SUFFIX = "Sparse";

    private static final Map<ShapeId, Set<ShapeId>> toSparse = MapUtils.ofEntries(
        // https://github.com/aws/aws-sdk-go-v2/issues/3575
        serviceToMembers("com.amazonaws.apigateway#BackplaneControlService",
                "com.amazonaws.apigateway#IntegrationResponse$responseTemplates"),
        // https://github.com/aws/aws-sdk-go-v2/issues/3501
        serviceToMembers("com.amazonaws.appstream#PhotonAdminProxyService",
                "com.amazonaws.appstream#Application$Metadata")
    );

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        var members = toSparse.get(settings.getService());
        if (members == null) {
            return model;
        }

        var builder = model.toBuilder();
        var updates = new ArrayList<Shape>();
        var added = new HashSet<ShapeId>();
        for (var memberId : members) {
            var member = model.expectShape(memberId, MemberShape.class);
            var map = model.expectShape(member.getTarget(), MapShape.class);
            if (map.hasTrait(SparseTrait.class)) {
                continue;
            }

            var sparseId = ShapeId.fromParts(map.getId().getNamespace(), map.getId().getName() + SPARSE_SUFFIX);
            if (added.add(sparseId)) {
                builder.addShape(sparseCopy(map, sparseId));
            }
            updates.add(member.toBuilder().target(sparseId).build());
        }

        return ModelTransformer.create().replaceShapes(builder.build(), updates);
    }

    private static MapShape sparseCopy(MapShape map, ShapeId id) {
        return map.toBuilder()
                .id(id)
                .key(map.getKey().toBuilder().id(id.withMember("key")).build())
                .value(map.getValue().toBuilder().id(id.withMember("value")).build())
                .addTrait(new SparseTrait())
                .build();
    }

    private static Map.Entry<ShapeId, Set<ShapeId>> serviceToMembers(String serviceId, String... memberIds) {
        return Map.entry(
                ShapeId.from(serviceId),
                Arrays.stream(memberIds).map(ShapeId::from).collect(Collectors.toSet()));
    }
}
