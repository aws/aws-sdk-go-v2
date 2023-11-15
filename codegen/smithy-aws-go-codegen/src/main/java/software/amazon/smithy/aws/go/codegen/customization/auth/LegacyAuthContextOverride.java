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

package software.amazon.smithy.aws.go.codegen.customization.auth;

import software.amazon.smithy.aws.go.codegen.SdkGoTypes;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.MiddlewareIdentifier;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.go.codegen.auth.SignRequestMiddlewareGenerator;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator.createFinalizeStepMiddleware;
import static software.amazon.smithy.go.codegen.GoWriter.emptyGoTemplate;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * The SDK has historically respected the values for sigv4[a] signing name and region[s] via the context setters
 * SetSigningName and SetSigningRegion. To preserve this we preface the auth resolution middleware with polyfills to
 * update the corresponding values in signing properties for the resolved auth scheme.
 * (post-resolution) the values here accordingly.
 */
public class LegacyAuthContextOverride implements GoIntegration {
    public static final String MIDDLEWARE_NAME = "setLegacyContextSigningOptionsMiddleware";
    public static final String MIDDLEWARE_ID = "setLegacyContextSigningOptions";
    public static final String MIDDLEWARE_ADD_FUNC = "addSetLegacyContextSigningOptionsMiddleware";

    private static final MiddlewareRegistrar MIDDLEWARE_REGISTRAR = MiddlewareRegistrar.builder()
            .resolvedFunction(buildPackageSymbol(MIDDLEWARE_ADD_FUNC))
            .build();

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MIDDLEWARE_REGISTRAR)
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator
    ) {
        goDelegator.useFileWriter("auth.go", settings.getModuleName(), generateMiddleware());
    }

    private GoWriter.Writable generateMiddleware() {
        return GoWriter.ChainWritable.of(
                createFinalizeStepMiddleware(MIDDLEWARE_NAME, MiddlewareIdentifier.string(MIDDLEWARE_ID))
                        .asWritable(generateMiddlewareBody(), emptyGoTemplate()),
                generateAddFunc()
        ).compose();
    }

    private GoWriter.Writable generateMiddlewareBody() {
        return goTemplate("""
                rscheme := getResolvedAuthScheme(ctx)
                schemeID := rscheme.Scheme.SchemeID()
                
                if sn := $getSigningName:T(ctx); sn != "" {
                    if schemeID == "aws.auth#sigv4" {
                        $setSigV4SigningName:T(&rscheme.SignerProperties, sn)
                    } else if schemeID == "aws.auth#sigv4a" {
                        $setSigV4ASigningName:T(&rscheme.SignerProperties, sn)
                    }
                }
                
                if sr := $getSigningRegion:T(ctx); sr != "" {
                    if schemeID == "aws.auth#sigv4" {
                        $setSigV4SigningRegion:T(&rscheme.SignerProperties, sr)
                    } else if schemeID == "aws.auth#sigv4a" {
                        $setSigV4ASigningRegions:T(&rscheme.SignerProperties, []string{sr})
                    }
                }

                return next.HandleFinalize(ctx, in)
                """,
                MapUtils.of(
                        "getSigningName", SdkGoTypes.Aws.Middleware.GetSigningName,
                        "getSigningRegion", SdkGoTypes.Aws.Middleware.GetSigningRegion,
                        "setSigV4SigningName", SmithyGoTypes.Transport.Http.SetSigV4SigningName,
                        "setSigV4ASigningName", SmithyGoTypes.Transport.Http.SetSigV4ASigningName,
                        "setSigV4SigningRegion", SmithyGoTypes.Transport.Http.SetSigV4SigningRegion,
                        "setSigV4ASigningRegions", SmithyGoTypes.Transport.Http.SetSigV4ASigningRegions
                ));
    }

    private GoWriter.Writable generateAddFunc() {
        return goTemplate("""
                func $L(stack $P) error {
                    return stack.Finalize.Insert(&$L{}, $S, $T)
                }
                """,
                MIDDLEWARE_ADD_FUNC,
                SmithyGoTypes.Middleware.Stack,
                MIDDLEWARE_NAME,
                SignRequestMiddlewareGenerator.MIDDLEWARE_ID,
                SmithyGoTypes.Middleware.Before);
    }
}
