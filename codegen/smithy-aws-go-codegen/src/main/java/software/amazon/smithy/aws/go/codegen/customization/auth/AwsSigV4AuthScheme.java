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
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.integration.auth.SigV4Definition;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

/**
 * Adds auth scheme codegen support for aws.auth#sigv4. Region as a config and SDK-specific auth params/resolution are
 * supplied by other integrations.
 */
public class AwsSigV4AuthScheme implements GoIntegration {
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .addAuthSchemeDefinition(SigV4Trait.ID, new AwsSigV4())
                        .build()
        );
    }

    public static class AwsSigV4 extends SigV4Definition {
        @Override
        public GoWriter.Writable generateDefaultAuthScheme() {
            return goTemplate("$T($S, &$T{Signer: options.HTTPSignerV4})",
                    SdkGoTypes.Internal.Auth.NewHTTPAuthScheme,
                    SigV4Trait.ID.toString(),
                    SdkGoTypes.Internal.Auth.Smithy.V4SignerAdapter);
        }

        @Override
        public GoWriter.Writable generateOptionsIdentityResolver() {
            return goTemplate("&$T{Provider: o.Credentials}",
                    SdkGoTypes.Internal.Auth.Smithy.CredentialsProviderAdapter);
        }
    }


}
