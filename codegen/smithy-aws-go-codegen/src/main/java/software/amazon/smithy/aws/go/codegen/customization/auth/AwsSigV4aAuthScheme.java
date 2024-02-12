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

import software.amazon.smithy.aws.go.codegen.AwsSignatureVersion4aUtils;
import software.amazon.smithy.aws.go.codegen.SdkGoTypes;
import software.amazon.smithy.aws.traits.auth.SigV4ATrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStdlibTypes;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.integration.auth.SigV4ADefinition;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator.generateFinalizeMiddlewareFunc;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * Adds auth scheme codegen support for aws.auth#sigv4a.
 */
public class AwsSigV4aAuthScheme implements GoIntegration {
    private static final ConfigField Signer = ConfigField.builder()
            .name(AwsSignatureVersion4aUtils.V4A_SIGNER_INTERFACE_NAME)
            .type(buildPackageSymbol(AwsSignatureVersion4aUtils.V4A_SIGNER_INTERFACE_NAME))
            .documentation("Signature Version 4a (SigV4a) Signer")
            .build();

    private static final ConfigFieldResolver SignerResolver = ConfigFieldResolver.builder()
            .location(ConfigFieldResolver.Location.CLIENT)
            .target(ConfigFieldResolver.Target.INITIALIZATION)
            .resolver(buildPackageSymbol(AwsSignatureVersion4aUtils.SIGNER_RESOLVER))
            .build();

    private static boolean isSigV4A(Model model, ServiceShape service) {
        return service.hasTrait(SigV4ATrait.class);
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(AwsSigV4aAuthScheme::isSigV4A)
                        .addAuthSchemeDefinition(SigV4ATrait.ID, new AwsSigV4A())
                        .addConfigField(Signer)
                        .addConfigFieldResolver(SignerResolver)
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator
    ) {
        var service = settings.getService(model);
        if (!isSigV4A(model, service)) {
            return;
        }

        goDelegator.useFileWriter("options.go", settings.getModuleName(), generateAdditionalSource());
        goDelegator.useShapeWriter(service, writer -> {
            AwsSignatureVersion4aUtils.writerSignerInterface(writer);
            AwsSignatureVersion4aUtils.writerConfigFieldResolver(writer, service);
            AwsSignatureVersion4aUtils.writeNewV4ASignerFunc(writer, service);
        });
    }

    public static class AwsSigV4A extends SigV4ADefinition {
        @Override
        public GoWriter.Writable generateDefaultAuthScheme() {
            return goTemplate("""
                    $T($S, &$T{
                        Signer: options.httpSignerV4a,
                        Logger: options.Logger,
                        LogSigning: options.ClientLogMode.IsSigning(),
                    })""",
                    SdkGoTypes.Internal.Auth.NewHTTPAuthScheme,
                    SigV4ATrait.ID.toString(),
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
