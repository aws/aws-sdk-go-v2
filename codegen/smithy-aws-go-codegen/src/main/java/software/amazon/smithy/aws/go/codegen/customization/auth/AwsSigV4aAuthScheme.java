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
import software.amazon.smithy.go.codegen.GoStdlibTypes;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.go.codegen.auth.SigV4aTrait;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.integration.auth.SigV4aDefinition;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator.generateFinalizeMiddlewareFunc;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

/**
 * Adds auth scheme codegen support for aws.auth#sigv4a.
 */
public class AwsSigV4aAuthScheme implements GoIntegration {
    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .addAuthSchemeDefinition(SigV4aTrait.ID, new AwsSigV4a())
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator
    ) {
        if (settings.getService(model).hasTrait(SigV4aTrait.class)) {
            goDelegator.useFileWriter("options.go", settings.getModuleName(), generateAdditionalSource());
        }
    }

    public static class AwsSigV4a extends SigV4aDefinition {
        @Override
        public GoWriter.Writable generateDefaultAuthScheme() {
            return goTemplate("$T($S, &$T{Signer: options.httpSignerV4a})",
                    SdkGoTypes.Internal.Auth.NewHTTPAuthScheme,
                    SigV4aTrait.ID.toString(),
                    SdkGoTypes.Internal.V4A.SignerAdapter);
        }

        @Override
        public GoWriter.Writable generateOptionsIdentityResolver() {
            return goTemplate("getSigV4AIdentityResolver(o)");
        }
    }

    private GoWriter.Writable generateAdditionalSource() {
        return GoWriter.ChainWritable.of(
                generateGetIdentityResolver(),
                generateHelpers()
        ).compose();
    }

    private GoWriter.Writable generateGetIdentityResolver() {
        return goTemplate("""
                func getSigV4AIdentityResolver(o Options) $T {
                    if o.Credentials != nil {
                        return &$T{
                            Provider: &$T{
                                SymmetricProvider: o.Credentials,
                            },
                        }
                    }
                    return nil
                }
                """,
                SmithyGoTypes.Auth.IdentityResolver,
                SdkGoTypes.Internal.V4A.CredentialsProviderAdapter,
                SdkGoTypes.Internal.V4A.SymmetricCredentialAdaptor);
    }

    private GoWriter.Writable generateHelpers() {
        return goTemplate("""
                // WithSigV4ASigningRegions applies an override to the authentication workflow to
                // use the given signing region set for SigV4A-authenticated operations.
                //
                // This is an advanced setting. The value here is FINAL, taking precedence over
                // the resolved signing region set from both auth scheme resolution and endpoint
                // resolution.
                func WithSigV4ASigningRegions(regions []string) func(*Options) {
                    fn := $mw:W
                    return func(o *Options) {
                        o.APIOptions = append(o.APIOptions, func(s $stack:P) error {
                            return s.Finalize.Insert(
                                middleware.FinalizeMiddlewareFunc("withSigV4ASigningRegions", fn),
                                "Signing",
                                $before:T,
                            )
                        })
                    }
                }
                """,
                MapUtils.of(
                        "stack", SmithyGoTypes.Middleware.Stack,
                        "before", SmithyGoTypes.Middleware.Before,
                        "mw", generateFinalizeMiddlewareFunc(goTemplate("""
                                rscheme := getResolvedAuthScheme(ctx)
                                if rscheme == nil {
                                    return out, metadata, $T("no resolved auth scheme")
                                }

                                $T(&rscheme.SignerProperties, regions)
                                return next.HandleFinalize(ctx, in)
                                """,
                                GoStdlibTypes.Fmt.Errorf,
                                SmithyGoTypes.Transport.Http.SetSigV4ASigningRegions))
                ));
    }
}
