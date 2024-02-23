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

package software.amazon.smithy.aws.go.codegen;

import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

import java.util.List;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

/**
 * wraps a smithy request error with an AWS Request error.
 */
public class AWSResponseErrorWrapper implements GoIntegration {
    private static final String ADD_ERROR_MIDDLEWARE = "addResponseErrorMiddleware";

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first.
     *
     * @return Returns the sort order, defaults to -40.
     */
    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceShape service = settings.getService(model);

        // S3 doesn't need aws specific wrapper
        if (requiresS3Customization(model, service)) {return;}

        goDelegator.useShapeWriter(service, this::writeMiddlewareHelper);
    }

    private void writeMiddlewareHelper(GoWriter writer) {
        writer.openBlock("func $L(stack *middleware.Stack) error {", "}", ADD_ERROR_MIDDLEWARE, () -> {
            writer.write("""
                    return stack.Deserialize.Insert(&$T{}, "RequestIDRetriever", middleware.Before)
                    """, AwsGoDependency.AWS_HTTP_TRANSPORT.valueSymbol("ResponseErrorWrapper"));
        });
        writer.insertTrailingNewline();
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(((model, serviceShape) -> !requiresS3Customization(model,serviceShape)))
                        .registerMiddleware(
                                MiddlewareRegistrar.builder()
                                        .resolvedFunction(buildPackageSymbol(ADD_ERROR_MIDDLEWARE))
                                        .build()
                        ).build()
        );
    }

    // returns true if service is either s3 or s3 control and needs s3 customization
    private static boolean requiresS3Customization(Model model, ServiceShape service) {
        String serviceId= service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("S3") || serviceId.equalsIgnoreCase("S3 Control");
    }
}
