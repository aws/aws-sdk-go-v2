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

package software.amazon.smithy.aws.go.codegen;

import java.util.Optional;
import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;


/**
 * Used by integrations to generate an AWS
 * authentication scheme resolution.
 *
 */
public class AwsEndpointAuthSchemeGenerator implements GoIntegration {

    @Override
    public void renderPostEndpointResolutionHook(GoSettings settings, GoWriter writer, Model model) {
        ServiceShape serviceShape = settings.getService(model);
        writer.write(
            """
            $W

            for _, authScheme := range authSchemes {
                switch authScheme.(type) {
                    case $P:
                        $W
                        break
                    case $P:
                        $W
                        break
                    case $P:
                        break
                }
            }
            """,
            generateAuthSchemeDetection(serviceShape),
            SymbolUtils.createPointableSymbolBuilder("AuthenticationSchemeV4", AwsGoDependency.INTERNAL_AUTH).build(),
            generateSigV4Resolution(serviceShape),
            SymbolUtils.createPointableSymbolBuilder("AuthenticationSchemeV4A", AwsGoDependency.INTERNAL_AUTH).build(),
            generateSigV4AResolution(serviceShape),
            SymbolUtils.createPointableSymbolBuilder("AuthenticationSchemeNone", AwsGoDependency.INTERNAL_AUTH).build()
        );
    }

    private GoWriter.Writable generateAuthSchemeDetection(ServiceShape serviceShape) {
        GoWriter.Writable signerVersion = (GoWriter writer) -> {
            String serviceId = serviceShape.expectTrait(ServiceTrait.class).getSdkId();
            if (serviceId.equalsIgnoreCase("S3")) {
                writer.write(
                    """
                        ctx = $T(ctx, $T)
                    """,
                    SymbolUtils.createValueSymbolBuilder("SetSignerVersion", AwsCustomGoDependency.S3_CUSTOMIZATION).build(),
                    SymbolUtils.createValueSymbolBuilder("SigV4", AwsGoDependency.INTERNAL_AUTH).build()
                );
            } else if (serviceId.equalsIgnoreCase("EventBridge")) {
                writer.write(
                    """
                        ctx = $T(ctx, $T)
                    """,
                    SymbolUtils.createValueSymbolBuilder("SetSignerVersion", AwsCustomGoDependency.EVENTBRIDGE_CUSTOMIZATION).build(),
                    SymbolUtils.createValueSymbolBuilder("SigV4", AwsGoDependency.INTERNAL_AUTH).build()
                );            } else {
                writer.write("");
            } 
        };
        return (GoWriter writer) -> {
            var signingNameDefaultOpt = getDefaultSigningName(serviceShape);
            var signingNameDefault = signingNameDefaultOpt.isPresent() ? signingNameDefaultOpt.get() : "";
            writer.write(
                """
                    authSchemes, err := $T(&resolvedEndpoint.Properties)
                    if err != nil {
                        var nfe $P
                        if $T(err, &nfe) {
                            // if no auth scheme is found, default to sigv4
                            signingName := \"$L\"
                            signingRegion := m.BuiltInResolver.(*BuiltInResolver).Region
                            ctx = $T(ctx, signingName)
                            ctx = $T(ctx, signingRegion)
                            $W
                        }
                        var ue $P
                        if errors.As(err, &ue) {
                            return out, metadata, $T(
                                \"This operation requests signer version(s) %v but the client only supports %v\",
                                ue.UnsupportedSchemes,
                                $T,
                            )
                        }
                    }
                """,
                SymbolUtils.createValueSymbolBuilder("GetAuthenticationSchemes", AwsGoDependency.INTERNAL_AUTH).build(),
                SymbolUtils.createPointableSymbolBuilder("NoAuthenticationSchemesFoundError", AwsGoDependency.INTERNAL_AUTH).build(),
                SymbolUtils.createValueSymbolBuilder("As", SmithyGoDependency.ERRORS).build(),
                signingNameDefault,
                SymbolUtils.createValueSymbolBuilder("SetSigningName", AwsGoDependency.AWS_MIDDLEWARE).build(),
                SymbolUtils.createValueSymbolBuilder("SetSigningRegion", AwsGoDependency.AWS_MIDDLEWARE).build(),
                signerVersion,
                SymbolUtils.createPointableSymbolBuilder("UnSupportedAuthenticationSchemeSpecifiedError", AwsGoDependency.INTERNAL_AUTH).build(),
                SymbolUtils.createValueSymbolBuilder("Errorf", SmithyGoDependency.FMT).build(),
                SymbolUtils.createValueSymbolBuilder("SupportedSchemes", AwsGoDependency.INTERNAL_AUTH).build()

            );
        };
    }

    private GoWriter.Writable generateSigV4Resolution(ServiceShape serviceShape) {
        GoWriter.Writable signerVersion = (GoWriter writer) -> {
            String serviceId = serviceShape.expectTrait(ServiceTrait.class).getSdkId();
            if (serviceId.equalsIgnoreCase("S3")) {
                writer.write(
                    """
                        ctx = $T(ctx, v4Scheme.Name)
                    """,
                    SymbolUtils.createValueSymbolBuilder("SetSignerVersion", AwsCustomGoDependency.S3_CUSTOMIZATION).build()
                );
            } else if (serviceId.equalsIgnoreCase("EventBridge")) {
                writer.write(
                    """
                        ctx = $T(ctx, v4Scheme.Name)
                    """,
                    SymbolUtils.createValueSymbolBuilder("SetSignerVersion", AwsCustomGoDependency.EVENTBRIDGE_CUSTOMIZATION).build()
                );            } else {
                writer.write("");
            } 
        };

        return (GoWriter writer) -> {
            var signingNameDefaultOpt = getDefaultSigningName(serviceShape);
            var signingNameDefault = signingNameDefaultOpt.isPresent() ? signingNameDefaultOpt.get() : "";
            writer.write(
                """
                    v4Scheme, _ := authScheme.($P)
                    var signingName, signingRegion string
                    if v4Scheme.SigningName == nil {
                        signingName = \"$L\"
                    }
                    if v4Scheme.SigningRegion == nil {
                        signingRegion = m.BuiltInResolver.(*BuiltInResolver).Region
                    }
                    if v4Scheme.DisableDoubleEncoding != nil {
                        // The signer sets an equivalent value at client initialization time.
                        // Setting this context value will cause the signer to extract it
                        // and override the value set at client initialization time.
                        ctx = $T(ctx, *v4Scheme.DisableDoubleEncoding)
                    }
                    ctx = $T(ctx, signingName)
                    ctx = $T(ctx, signingRegion)
                    $W
                """,
                SymbolUtils.createPointableSymbolBuilder("AuthenticationSchemeV4", AwsGoDependency.INTERNAL_AUTH).build(),
                signingNameDefault,
                SymbolUtils.createValueSymbolBuilder("SetDisableDoubleEncoding", AwsGoDependency.INTERNAL_AUTH).build(),
                SymbolUtils.createValueSymbolBuilder("SetSigningName", AwsGoDependency.AWS_MIDDLEWARE).build(),
                SymbolUtils.createValueSymbolBuilder("SetSigningRegion", AwsGoDependency.AWS_MIDDLEWARE).build(),
                signerVersion
            );
        };
    }

    private GoWriter.Writable generateSigV4AResolution(ServiceShape serviceShape) {
        GoWriter.Writable signerVersion = (GoWriter writer) -> {
            String serviceId = serviceShape.expectTrait(ServiceTrait.class).getSdkId();
            if (serviceId.equalsIgnoreCase("S3")) {
                writer.write(
                    """
                    ctx = $T(ctx, v4aScheme.Name)
                    """,
                    SymbolUtils.createValueSymbolBuilder("SetSignerVersion", AwsCustomGoDependency.S3_CUSTOMIZATION).build()
                );
            } else if (serviceId.equalsIgnoreCase("EventBridge")) {
                writer.write(
                    """
                    ctx = $T(ctx, v4aScheme.Name)
                    """,
                    SymbolUtils.createValueSymbolBuilder("SetSignerVersion", AwsCustomGoDependency.EVENTBRIDGE_CUSTOMIZATION).build()
                );            } else {
                writer.write("");
            }
        };
        
        return (GoWriter writer) -> {
            var signingNameDefaultOpt = getDefaultSigningName(serviceShape);
            var signingNameDefault = signingNameDefaultOpt.isPresent() ? signingNameDefaultOpt.get() : "";
            writer.write(
                """
                    v4aScheme, _ := authScheme.($P)
                    if v4aScheme.SigningName == nil {
                        v4aScheme.SigningName = $T(\"$L\")
                    }
                    if v4aScheme.DisableDoubleEncoding != nil {
                        // The signer sets an equivalent value at client initialization time.
                        // Setting this context value will cause the signer to extract it
                        // and override the value set at client initialization time.
                        ctx = $T(ctx, *v4aScheme.DisableDoubleEncoding)
                    }
                    ctx = $T(ctx, *v4aScheme.SigningName)
                    ctx = $T(ctx, v4aScheme.SigningRegionSet[0]) 
                    $W
                """,
                SymbolUtils.createPointableSymbolBuilder("AuthenticationSchemeV4A", AwsGoDependency.INTERNAL_AUTH).build(),
                SymbolUtils.createValueSymbolBuilder("String", AwsGoDependency.AWS_CORE).build(),
                SymbolUtils.createValueSymbolBuilder(signingNameDefault).build(),
                SymbolUtils.createValueSymbolBuilder("SetDisableDoubleEncoding", AwsGoDependency.INTERNAL_AUTH).build(),
                SymbolUtils.createValueSymbolBuilder("SetSigningName", AwsGoDependency.AWS_MIDDLEWARE).build(),
                SymbolUtils.createValueSymbolBuilder("SetSigningRegion", AwsGoDependency.AWS_MIDDLEWARE).build(),
                signerVersion

            );
        };
    }

    private Optional<String> getDefaultSigningName(ServiceShape serviceShape) {
        var sigV4TraitOpt = serviceShape.getTrait(SigV4Trait.class);
        String signingNameDefault = "";
        if (sigV4TraitOpt.isPresent()) {
            signingNameDefault = sigV4TraitOpt.get().getName();
            return Optional.of(signingNameDefault);
        }
        return Optional.empty();
    }
    
}
