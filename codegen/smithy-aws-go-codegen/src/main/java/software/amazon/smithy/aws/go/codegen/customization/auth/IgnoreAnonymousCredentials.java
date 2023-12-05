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
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * aws.AnonymousCredentials as a sentinel is redundant with the SRA auth refactor and prevents the client from resolving
 * to anonymous if necessary, this integration nils them out if set on client options.
 */
public class IgnoreAnonymousCredentials implements GoIntegration {
    public static final ConfigFieldResolver IGNORE_ANONYMOUS_AUTH = ConfigFieldResolver.builder()
            .location(ConfigFieldResolver.Location.CLIENT)
            .target(ConfigFieldResolver.Target.FINALIZATION)
            .resolver(buildPackageSymbol("ignoreAnonymousAuth"))
            .build();

    private static boolean hasCredentials(Model model, ServiceShape service) {
        return service.hasTrait(SigV4Trait.class) || service.hasTrait(SigV4Trait.class);
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(IgnoreAnonymousCredentials::hasCredentials)
                        .addConfigFieldResolver(IGNORE_ANONYMOUS_AUTH)
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

    private GoWriter.Writable generateResolver() {
        return goTemplate("""
                func ignoreAnonymousAuth(options *Options) {
                    if $T(options.Credentials, ($P)(nil)) {
                        options.Credentials = nil
                    }
                }
                """, SdkGoTypes.Aws.IsCredentialsProvider, SdkGoTypes.Aws.AnonymousCredentials);
    }
}
