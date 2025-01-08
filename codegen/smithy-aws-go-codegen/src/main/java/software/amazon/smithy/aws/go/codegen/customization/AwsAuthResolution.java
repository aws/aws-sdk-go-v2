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

package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoStdlibTypes;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.auth.AuthParameter;
import software.amazon.smithy.go.codegen.auth.AuthParametersGenerator;
import software.amazon.smithy.go.codegen.auth.AuthParametersResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * Adds customizations for auth resolution in AWS services:
 * 1. Adds a field+resolver for endpoint parameters for the two services (s3, eventbridge) that delegate to endpoint
 *    rules for auth resolution.
 * 2. Adds a resolver for region.
 */
public class AwsAuthResolution implements GoIntegration {
    private final AuthParameter endpointParams = new AuthParameter(
            "endpointParams",
            "The endpoint resolver parameters for this operation. " +
                    "This service's default resolver delegates to endpoint rules.",
            SymbolUtils.createPointableSymbolBuilder("EndpointParameters").build()
    );

    private final AuthParametersResolver regionResolver = new AuthParametersResolver(
            buildPackageSymbol("bindAuthParamsRegion")
    );

    private final AuthParametersResolver endpointParamsResolver = new AuthParametersResolver(
            buildPackageSymbol("bindAuthEndpointParams")
    );

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(this::isEndpointAuthService)
                        .addAuthParameter(endpointParams)
                        .addAuthParameterResolver(endpointParamsResolver)
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(this::isSigV4Service)
                        .addAuthParameterResolver(regionResolver)
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator) {
        if (isSigV4Service(model, settings.getService(model))) {
            goDelegator.useFileWriter("auth.go", settings.getModuleName(), writeRegionResolver());
        }
        if (isEndpointAuthService(model, settings.getService(model))) {
            goDelegator.useFileWriter("auth.go", settings.getModuleName(), writeEndpointParamResolver());
        }
    }

    private boolean isEndpointAuthService(Model model, ServiceShape service) {
        final String sdkId = service.expectTrait(ServiceTrait.class).getSdkId();
        return sdkId.equalsIgnoreCase("s3")
                || sdkId.equalsIgnoreCase("eventbridge")
                || sdkId.equalsIgnoreCase("sesv2");
    };

    private boolean isSigV4Service(Model model, ServiceShape service) {
        return service.hasTrait(SigV4Trait.class);
    };

    private GoWriter.Writable writeRegionResolver() {
        return goTemplate("""
                func bindAuthParamsRegion( _ interface{}, params $P, _ interface{}, options Options) {
                    params.Region = options.Region
                }
                """, AuthParametersGenerator.STRUCT_SYMBOL);
    }

    private GoWriter.Writable writeEndpointParamResolver() {
        return goTemplate("""
                func bindAuthEndpointParams(ctx $P, params $P, input interface{}, options Options) {
                    params.endpointParams = bindEndpointParams(ctx, input, options)
                }
                """, GoStdlibTypes.Context.Context, AuthParametersGenerator.STRUCT_SYMBOL);
    }
}
