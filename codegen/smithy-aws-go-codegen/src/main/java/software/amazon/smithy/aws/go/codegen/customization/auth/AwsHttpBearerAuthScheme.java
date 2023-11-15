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

import java.util.List;

import software.amazon.smithy.aws.go.codegen.SdkGoTypes;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.integration.auth.HttpBearerDefinition;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.traits.HttpBearerAuthTrait;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

/**
 * Adds full codegen support for smithy.api#httpBearerAuth.
 */
public class AwsHttpBearerAuthScheme implements GoIntegration {
    public static final String TOKEN_PROVIDER_OPTION_NAME = "BearerAuthTokenProvider";
    private static final String SIGNER_OPTION_NAME = "BearerAuthSigner";
    private static final String NEW_DEFAULT_SIGNER_NAME = "newDefault" + SIGNER_OPTION_NAME;
    private static final String SIGNER_RESOLVER_NAME = "resolve" + SIGNER_OPTION_NAME;

    @Override
    public void writeAdditionalFiles(
            GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator
    ) {
        var service = settings.getService(model);
        if (!isHttpBearerService(model, service)) {
            return;
        }

        goDelegator.useShapeWriter(service, writer -> writer.write("""
                $W
                $W
                """, writeSignerConfigFieldResolver(), writeNewSignerFunc()));
    }

    private boolean isHttpBearerService(Model model, ServiceShape service) {
        return service.hasTrait(HttpBearerAuthTrait.class);
    }

    private GoWriter.Writable writeSignerConfigFieldResolver() {
        return goTemplate("""
                func $funcName:L(o *Options) {
                    if o.$signerOption:L != nil {
                        return
                    }
                    o.$signerOption:L = $newDefaultSigner:L(*o)
                }
                """,
                MapUtils.of(
                        "funcName", SIGNER_RESOLVER_NAME,
                        "signerOption", SIGNER_OPTION_NAME,
                        "newDefaultSigner", NEW_DEFAULT_SIGNER_NAME
                ));
    }

    private GoWriter.Writable writeNewSignerFunc() {
        return goTemplate("""
                func $funcName:L(o Options) $signerInterface:T {
                    return $newDefaultSigner:T()
                }
                """,
                MapUtils.of(
                        "funcName", NEW_DEFAULT_SIGNER_NAME,
                        "signerInterface", SmithyGoTypes.Auth.Bearer.Signer,
                        "newDefaultSigner", SmithyGoTypes.Auth.Bearer.NewSignHTTPSMessage
                ));
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(this::isHttpBearerService)
                        .addConfigField(ConfigField.builder()
                                .name(TOKEN_PROVIDER_OPTION_NAME)
                                .type(SmithyGoTypes.Auth.Bearer.TokenProvider)
                                .documentation("Bearer token value provider")
                                .build())
                        .addConfigField(ConfigField.builder()
                                .name(SIGNER_OPTION_NAME)
                                .type(SmithyGoTypes.Auth.Bearer.Signer)
                                .documentation("Signer for authenticating requests with bearer auth")
                                .build())
                        .addConfigFieldResolver(ConfigFieldResolver.builder()
                                .location(ConfigFieldResolver.Location.CLIENT)
                                .target(ConfigFieldResolver.Target.INITIALIZATION)
                                .resolver(SymbolUtils.createValueSymbolBuilder(SIGNER_RESOLVER_NAME).build())
                                .build())
                        .addAuthSchemeDefinition(HttpBearerAuthTrait.ID, new AwsHttpBearer())
                        .build()
        );
    }

    public static class AwsHttpBearer extends HttpBearerDefinition {
        @Override
        public GoWriter.Writable generateDefaultAuthScheme() {
            return goTemplate("$T($S, &$T{Signer: options.BearerAuthSigner})",
                    SdkGoTypes.Internal.Auth.NewHTTPAuthScheme,
                    HttpBearerAuthTrait.ID.toString(),
                    SdkGoTypes.Internal.Auth.Smithy.BearerTokenSignerAdapter);
        }

        @Override
        public GoWriter.Writable generateOptionsIdentityResolver() {
            return goTemplate("&$T{Provider: o.BearerAuthTokenProvider}",
                    SdkGoTypes.Internal.Auth.Smithy.BearerTokenProviderAdapter);
        }
    }
}