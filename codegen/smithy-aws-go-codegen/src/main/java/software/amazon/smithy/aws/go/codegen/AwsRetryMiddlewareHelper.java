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

import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;

public class AwsRetryMiddlewareHelper implements GoIntegration {
    public static final String ADD_RETRY_MIDDLEWARES_HELPER = "addRetryMiddlewares";

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator delegator
    ) {
        delegator.useShapeWriter(settings.getService(model), this::generateRetryMiddlewareHelpers);
    }

    private void generateRetryMiddlewareHelpers(GoWriter writer) {
        Symbol stackSymbol = SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE)
                .build();
        Symbol addRetryMiddlewares = SymbolUtils.createValueSymbolBuilder("AddRetryMiddlewares",
                AwsGoDependency.AWS_RETRY).build();
        Symbol addOptions = SymbolUtils.createValueSymbolBuilder("AddRetryMiddlewaresOptions",
                AwsGoDependency.AWS_RETRY).build();

        writer.openBlock("func $L(stack $P, o Options) error {", "}", ADD_RETRY_MIDDLEWARES_HELPER, stackSymbol,
                () -> {
                    writer.openBlock("mo := $T{", "}", addOptions, () -> {
                        writer.write("$L: o.$L,", AddAwsConfigFields.RETRYER_CONFIG_NAME,
                                AddAwsConfigFields.RETRYER_CONFIG_NAME);
                        writer.write("LogRetryAttempts: o.$L.IsRetries(),",
                                AddAwsConfigFields.LOG_MODE_CONFIG_NAME);
                    });

                    writer.write("return $T(stack, mo)", addRetryMiddlewares);
                });
    }
}
