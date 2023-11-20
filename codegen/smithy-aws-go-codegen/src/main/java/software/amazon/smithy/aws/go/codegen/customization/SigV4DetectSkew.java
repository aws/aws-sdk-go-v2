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
 *
 *
 */

package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.aws.go.codegen.SdkGoTypes;
import software.amazon.smithy.aws.go.codegen.customization.util.ServicePredicates;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStdlibTypes;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.go.codegen.auth.SignRequestMiddlewareGenerator;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * Adds clock skew detection middleware to AWS service clients. Detection only applies to instances of v4.Signer.
 */
public class SigV4DetectSkew implements GoIntegration {
    private static final String ADD_FUNC_NAME = "addV4DetectSkewMiddleware";

    private static final MiddlewareRegistrar MIDDLEWARE = MiddlewareRegistrar.builder()
            .resolvedFunction(buildPackageSymbol(ADD_FUNC_NAME))
            .useClientOptions()
            .build();

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(ServicePredicates::isSigV4)
                        .registerMiddleware(MIDDLEWARE)
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator
    ) {
        if (ServicePredicates.isSigV4(model, settings.getService(model))) {
            goDelegator.useFileWriter("api_client.go", settings.getModuleName(), generateAddMiddleware());
        }
    }

    private GoWriter.Writable generateAddMiddleware() {
        return goTemplate("""
                func $func:L(stack $stack:P, options Options) error {
                    signer, ok := options.HTTPSignerV4.($signer:P)
                    if !ok {
                        return nil
                    }

                    m := &$middleware:T{
                        Signer: signer,
                    }
                    if err := stack.Finalize.Insert(m, $id:S, $after:T); err != nil {
                        return $errorf:T("add aws.signer.v4#DetectSkew: %v", err)
                    }
                    return nil
                }
                """,
                MapUtils.of(
                        "func", ADD_FUNC_NAME,
                        "stack", SmithyGoTypes.Middleware.Stack,
                        "signer", SdkGoTypes.Aws.Signer.V4.Signer,
                        "middleware", SdkGoTypes.Aws.Signer.V4.DetectSkewMiddleware,
                        "id", SignRequestMiddlewareGenerator.MIDDLEWARE_ID,
                        "after", SmithyGoTypes.Middleware.After,
                        "errorf", GoStdlibTypes.Fmt.Errorf
                ));
    }
}
