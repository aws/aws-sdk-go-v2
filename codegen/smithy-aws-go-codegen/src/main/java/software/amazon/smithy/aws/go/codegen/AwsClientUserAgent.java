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

import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.CodegenUtils;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;

public class AwsClientUserAgent implements GoIntegration {
    public static final String MIDDLEWARE_RESOLVER = "addClientUserAgent";

    @Override
    public byte getOrder() {
        return -49;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceTrait serviceTrait = settings.getService(model).expectTrait(ServiceTrait.class);
        String serviceId = serviceTrait.getSdkId().replace("-", "").replace(" ", "").toLowerCase();

        goDelegator.useShapeWriter(settings.getService(model), writer -> {
            writer.openBlock("func $L(stack $P) error {", "}", MIDDLEWARE_RESOLVER, SymbolUtils.createPointableSymbolBuilder("Stack",
                    SmithyGoDependency.SMITHY_MIDDLEWARE).build(), () -> {
                writer.write("return $T($T, $S, $T)(stack)",
                        SymbolUtils.createValueSymbolBuilder("AddSDKAgentKeyValue", AwsGoDependency.AWS_MIDDLEWARE)
                                .build(),
                        SymbolUtils.createValueSymbolBuilder("APIMetadata",
                                AwsGoDependency.AWS_MIDDLEWARE).build(),
                        serviceId,
                        SymbolUtils.createValueSymbolBuilder("goModuleVersion").build());
            });
            writer.write("");
        });
    }
}
