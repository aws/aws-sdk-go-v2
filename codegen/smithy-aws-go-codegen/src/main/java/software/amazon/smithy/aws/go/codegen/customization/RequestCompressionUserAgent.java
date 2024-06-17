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
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.traits.RequestCompressionTrait;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * Tracks when the caller uses request compression (smithy.api#requestCompression).
 */
public class RequestCompressionUserAgent implements GoIntegration {
    private static final MiddlewareRegistrar MIDDLEWARE_USERAGENT_RETRY = MiddlewareRegistrar.builder()
            .resolvedFunction(buildPackageSymbol("addIsRequestCompressionUserAgent"))
            .useClientOptions()
            .build();

    private static boolean hasRequestCompression(Model model, ServiceShape service) {
        return TopDownIndex.of(model)
                .getContainedOperations(service).stream()
                .anyMatch(it -> it.hasTrait(RequestCompressionTrait.class));
    }

    private static boolean isRequestCompression(Model model, ServiceShape service, OperationShape operation) {
        return operation.hasTrait(RequestCompressionTrait.class);
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return List.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate(RequestCompressionUserAgent::isRequestCompression)
                        .registerMiddleware(MIDDLEWARE_USERAGENT_RETRY)
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator) {
        if (!hasRequestCompression(model, settings.getService(model))) {
            return;
        }

        goDelegator.useFileWriter("api_client.go", settings.getModuleName(), goTemplate("""
                func addIsRequestCompressionUserAgent(stack $stack:P, options Options) error {
                    ua, err := getOrAddRequestUserAgent(stack)
                    if err != nil {
                        return err
                    }

                    if !options.DisableRequestCompression {
                        ua.AddUserAgentFeature($featureRequestCompression:T)
                    }
                    return nil
                }""",
                Map.of(
                        "stack", SmithyGoDependency.SMITHY_MIDDLEWARE.struct("Stack"),
                        "featureRequestCompression", AwsGoDependency.AWS_MIDDLEWARE
                                .struct("UserAgentFeatureGZIPRequestCompression")
                )));
    }
}
