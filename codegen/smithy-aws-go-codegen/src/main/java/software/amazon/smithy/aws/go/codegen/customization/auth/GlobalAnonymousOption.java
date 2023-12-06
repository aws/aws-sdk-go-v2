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

import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStdlibTypes;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

import static software.amazon.smithy.aws.go.codegen.customization.auth.EndpointAuthResolution.isEndpointAuthService;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * SDK v2 has always implicitly supported noAuth / anonymous as an option, because all of the legacy signing middlewares
 * would just noop if you didn't have credentials for them. We have to preserve this going forward in a post-SRA
 * landscape.
 */
public class GlobalAnonymousOption implements GoIntegration {
    private static final ConfigFieldResolver wrapWithAnonymousAuth = ConfigFieldResolver.builder()
            .resolver(buildPackageSymbol("wrapWithAnonymousAuth"))
            .location(ConfigFieldResolver.Location.CLIENT)
            .target(ConfigFieldResolver.Target.FINALIZATION)
            .build();

    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public java.util.List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate((model, service) -> !isEndpointAuthService(model, service))
                        .addConfigFieldResolver(wrapWithAnonymousAuth)
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator) {
        goDelegator.useFileWriter("auth.go", settings.getModuleName(), GoWriter.ChainWritable.of(
                generateAuthResolver(),
                generateConfigResolver()
        ).compose());
    }

    private GoWriter.Writable generateAuthResolver() {
        return goTemplate("""
                type withAnonymous struct {
                    resolver AuthSchemeResolver
                }

                var _ AuthSchemeResolver = (*withAnonymous)(nil)

                func (v *withAnonymous) ResolveAuthSchemes(ctx $context:T, params *AuthResolverParameters) ([]$option:P, error) {
                    opts, err := v.resolver.ResolveAuthSchemes(ctx, params)
                    if err != nil {
                        return nil, err
                    }

                    opts = append(opts, &$option:T{
                        SchemeID: $anonymous:T,
                    })
                    return opts, nil
                }
                """,
                MapUtils.of(
                        "context", GoStdlibTypes.Context.Context,
                        "option", SmithyGoTypes.Auth.Option,
                        "anonymous", SmithyGoTypes.Auth.SchemeIDAnonymous
                ));
    }

    private GoWriter.Writable generateConfigResolver() {
        return goTemplate("""
                func wrapWithAnonymousAuth(options *Options) {
                    if _, ok := options.AuthSchemeResolver.(*defaultAuthSchemeResolver); !ok {
                        return
                    }

                    options.AuthSchemeResolver = &withAnonymous{
                        resolver: options.AuthSchemeResolver,
                    }
                }
                """);
    }
}
