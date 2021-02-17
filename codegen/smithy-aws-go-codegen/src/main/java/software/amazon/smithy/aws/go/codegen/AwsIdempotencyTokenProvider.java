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

import java.util.List;
import java.util.Map;
import java.util.Optional;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.IdempotencyTokenMiddlewareGenerator;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.NeighborProviderIndex;
import software.amazon.smithy.model.neighbor.NeighborProvider;
import software.amazon.smithy.model.neighbor.Walker;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.IdempotencyTokenTrait;
import software.amazon.smithy.utils.ListUtils;

/**
 * Registers a client resolver for assigning the AWS default idempotency token provider.
 */
public final class AwsIdempotencyTokenProvider implements GoIntegration {
    private static final String RESOLVER_FUNCTION = "resolveIdempotencyTokenProvider";

    private boolean isTraitUsed = false;

    @Override
    public void writeAdditionalFiles(
            GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator
    ) {
        if (!isTraitUsed) {
            return;
        }
        goDelegator.useShapeWriter(settings.getService(model), this::writeResolver);
    }

    @Override
    public void processFinalizedModel(GoSettings settings, Model model) {
        Map<ShapeId, MemberShape> map = IdempotencyTokenMiddlewareGenerator.getOperationsWithIdempotencyToken(model,
                settings.getService(model));

        isTraitUsed = !map.isEmpty();
    }

    private void writeResolver(GoWriter writer) {
        String idempotencyConfigName = IdempotencyTokenMiddlewareGenerator.IDEMPOTENCY_CONFIG_NAME;
        writer.openBlock("func $L(o *Options) {", "}", RESOLVER_FUNCTION, () -> {
            writer.openBlock("if o.$L != nil {", "}", idempotencyConfigName, () -> writer.write("return"));
            writer.write("o.$L = $T($T)", idempotencyConfigName,
                    SymbolUtils.createValueSymbolBuilder("NewUUIDIdempotencyToken",
                            SmithyGoDependency.SMITHY_RAND).build(),
                    SymbolUtils.createValueSymbolBuilder("Reader", SmithyGoDependency.CRYPTORAND).build());
        });
        writer.write("");
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        if (!isTraitUsed) {
            return ListUtils.of();
        }

        return ListUtils.of(RuntimeClientPlugin.builder()
                .addConfigFieldResolver(ConfigFieldResolver.builder()
                        .location(ConfigFieldResolver.Location.CLIENT)
                        .target(ConfigFieldResolver.Target.INITIALIZATION)
                        .resolver(SymbolUtils.createValueSymbolBuilder(RESOLVER_FUNCTION).build())
                        .build())
                .build());
    }
}
