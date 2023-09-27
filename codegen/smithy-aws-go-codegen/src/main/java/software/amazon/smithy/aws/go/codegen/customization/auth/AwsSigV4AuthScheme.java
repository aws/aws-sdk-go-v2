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
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.integration.auth.SigV4Definition;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator.generateInitializeMiddlewareFunc;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

/**
 * Adds auth scheme codegen support for aws.auth#sigv4. Region as a config and SDK-specific auth params/resolution are
 * supplied by other integrations. Includes config helpers to set caller overrides for signing name and region.
 */
public class AwsSigV4AuthScheme implements GoIntegration {
    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .addAuthSchemeDefinition(SigV4Trait.ID, new AwsSigV4())
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator
    ) {
        if (settings.getService(model).hasTrait(SigV4Trait.class)) {
            goDelegator.useFileWriter("options.go", settings.getModuleName(), generateAdditionalSource());
        }
    }

    public static class AwsSigV4 extends SigV4Definition {
        @Override
        public GoWriter.Writable generateDefaultAuthScheme() {
            return goTemplate("""
                    $T($S, &$T{
                        Signer: options.HTTPSignerV4,
                        Logger: options.Logger,
                        LogSigning: options.ClientLogMode.IsSigning(),
                    })""",
                    SdkGoTypes.Internal.Auth.NewHTTPAuthScheme,
                    SigV4Trait.ID.toString(),
                    SdkGoTypes.Internal.Auth.Smithy.V4SignerAdapter);
        }

        @Override
        public GoWriter.Writable generateOptionsIdentityResolver() {
            return goTemplate("getSigV4IdentityResolver(o)");
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
                func getSigV4IdentityResolver(o Options) $T {
                    if o.Credentials != nil {
                        return &$T{Provider: o.Credentials}
                    }
                    return nil
                }
                """,
                SmithyGoTypes.Auth.IdentityResolver,
                SdkGoTypes.Internal.Auth.Smithy.CredentialsProviderAdapter);
    }

    private GoWriter.Writable generateHelpers() {
        return goTemplate("""
                // WithSigV4SigningName applies an override to the authentication workflow to
                // use the given signing name for SigV4-authenticated operations.
                //
                // This is an advanced setting. The value here is FINAL, taking precedence over
                // the resolved signing name from both auth scheme resolution and endpoint
                // resolution.
                func WithSigV4SigningName(name string) func(*Options) {
                    fn := $nameMW:W
                    return func(o *Options) {
                        o.APIOptions = append(o.APIOptions, func(s $stack:P) error {
                            return s.Initialize.Add(
                                middleware.InitializeMiddlewareFunc("withSigV4SigningName", fn),
                                $before:T,
                            )
                        })
                    }
                }

                // WithSigV4SigningRegion applies an override to the authentication workflow to
                // use the given signing region for SigV4-authenticated operations.
                //
                // This is an advanced setting. The value here is FINAL, taking precedence over
                // the resolved signing region from both auth scheme resolution and endpoint
                // resolution.
                func WithSigV4SigningRegion(region string) func(*Options) {
                    fn := $regionMW:W
                    return func(o *Options) {
                        o.APIOptions = append(o.APIOptions, func(s $stack:P) error {
                            return s.Initialize.Add(
                                middleware.InitializeMiddlewareFunc("withSigV4SigningRegion", fn),
                                $before:T,
                            )
                        })
                    }
                }
                """,
                MapUtils.of(
                        "stack", SmithyGoTypes.Middleware.Stack,
                        "before", SmithyGoTypes.Middleware.Before,
                        "nameMW", generateInitializeMiddlewareFunc(goTemplate("""
                                return next.HandleInitialize($T(ctx, name), in)
                                """, SdkGoTypes.Aws.Middleware.SetSigningName)),
                        "regionMW", generateInitializeMiddlewareFunc(goTemplate("""
                                return next.HandleInitialize($T(ctx, region), in)
                                """, SdkGoTypes.Aws.Middleware.SetSigningRegion))
                ));
    }
}
