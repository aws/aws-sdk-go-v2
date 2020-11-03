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
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;

public class RequestResponseLogging implements GoIntegration {
    private final static String MIDDLEWARE_HELPER = "addRequestResponseLogging";

    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        goDelegator.useShapeWriter(settings.getService(model), writer -> {
            Symbol stackSymbol = SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE)
                    .build();
            Symbol middlewareSymbol = SymbolUtils.createValueSymbolBuilder("RequestResponseLogger",
                    SmithyGoDependency.SMITHY_HTTP_TRANSPORT).build();

            writer.openBlock("func $L(stack $P, o Options) error {", "}", MIDDLEWARE_HELPER, stackSymbol, () -> {
                writer.openBlock("return stack.Deserialize.Add(&$T{", "}, middleware.After)", middlewareSymbol, () -> {
                    writer.write("LogRequest: o.$L.IsRequest(),", AddAwsConfigFields.LOG_MODE_CONFIG_NAME);
                    writer.write("LogRequestWithBody: o.$L.IsRequestWithBody(),",
                            AddAwsConfigFields.LOG_MODE_CONFIG_NAME);
                    writer.write("LogResponse: o.$L.IsResponse(),", AddAwsConfigFields.LOG_MODE_CONFIG_NAME);
                    writer.write("LogResponseWithBody: o.$L.IsResponseWithBody(),",
                            AddAwsConfigFields.LOG_MODE_CONFIG_NAME);
                });
            });
        });
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(RuntimeClientPlugin.builder()
                .registerMiddleware(MiddlewareRegistrar.builder()
                        .resolvedFunction(SymbolUtils.createValueSymbolBuilder(MIDDLEWARE_HELPER).build())
                        .useClientOptions()
                        .build())
                .build());
    }
}
