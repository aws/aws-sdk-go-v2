/*
 * Copyright 2024 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import java.util.Map;
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * Tracks the retry mode being used by the caller.
 */
public class RetryModeUserAgent implements GoIntegration {
    private static final MiddlewareRegistrar MIDDLEWARE = MiddlewareRegistrar.builder()
            .resolvedFunction(buildPackageSymbol("addUserAgentRetryMode"))
            .useClientOptions()
            .build();

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return List.of(
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MIDDLEWARE)
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator) {
        goDelegator.useFileWriter("api_client.go", settings.getModuleName(), goTemplate("""
                func addUserAgentRetryMode(stack $stack:P, options Options) error {
                    ua, err := getOrAddRequestUserAgent(stack)
                    if err != nil {
                        return err
                    }

                    switch options.Retryer.(type) {
                    case $standard:P:
                        ua.AddUserAgentFeature($featureStandard:T)
                    case $adaptive:P:
                        ua.AddUserAgentFeature($featureAdaptive:T)
                    }
                    return nil
                }""",
                Map.of(
                        "stack", SmithyGoDependency.SMITHY_MIDDLEWARE.struct("Stack"),
                        "standard", AwsGoDependency.AWS_RETRY.struct("Standard"),
                        "adaptive", AwsGoDependency.AWS_RETRY.struct("AdaptiveMode"),
                        "featureStandard", AwsGoDependency.AWS_MIDDLEWARE
                                .valueSymbol("UserAgentFeatureRetryModeStandard"),
                        "featureAdaptive", AwsGoDependency.AWS_MIDDLEWARE
                                .valueSymbol("UserAgentFeatureRetryModeAdaptive")
                )));
    }
}
