/*
 * Copyright 2026 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.Writable;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * Wraps the client's credentials provider in a credentials cache, so a provider set directly on
 * Options.Credentials, without going through config.LoadDefaultConfig, still gets cached.
 * internal/credentials/cachewrap decides which providers get wrapped.
 */
public class WrapCredentialsCache implements GoIntegration {
    public static final ConfigFieldResolver WRAP_CREDENTIALS_CACHE = ConfigFieldResolver.builder()
            .location(ConfigFieldResolver.Location.CLIENT)
            .target(ConfigFieldResolver.Target.FINALIZATION)
            .resolver(buildPackageSymbol("wrapCredentialsCache"))
            .build();

    private static boolean hasCredentials(Model model, ServiceShape service) {
        return service.hasTrait(SigV4Trait.class);
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(WrapCredentialsCache::hasCredentials)
                        .addConfigFieldResolver(WRAP_CREDENTIALS_CACHE)
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator
    ) {
        if (hasCredentials(model, settings.getService(model))) {
            goDelegator.useFileWriter("options.go", settings.getModuleName(), generateResolver());
        }
    }

    private Writable generateResolver() {
        return goTemplate("""
                func wrapCredentialsCache(options *Options) {
                    options.Credentials = $wrap:T(options.Credentials)
                }
                """,
                MapUtils.of(
                        "wrap", AwsGoDependency.INTERNAL_CREDENTIALS_CACHEWRAP.func("Wrap")
                ));
    }
}
