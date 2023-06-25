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
import software.amazon.smithy.go.codegen.AuthenticationSchemeGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.shapes.ServiceShape;


/**
 * Used by integrations to generate an AWS
 * authentication scheme resolution.
 *
 */
public class AwsAuthenticationSchemeGenerator implements GoIntegration {

    @Override
    public Optional<AuthenticationSchemeGenerator> getAuthenticationSchemeGenerator() {
        return Optional.of(new AwsEndpointAuthSchemeGenerator());

    }

    private class AwsEndpointAuthSchemeGenerator implements AuthenticationSchemeGenerator{

        AwsEndpointAuthSchemeGenerator() {
        }

        @Override
        public void renderEndpointBasedAuthSchemeResolution(GoWriter writer, ServiceShape serviceShape) {
            writer.write(
                """
                $W

                supportedAuthSchemeFound := false
                for _, authScheme := range authSchemes {
                    name := authScheme.GetName()
                    _, supportedAuthSchemeFound = $T[name]
                    if supportedAuthSchemeFound {
                        if name == $T {
                            $W
                            break
                        }
                        if name == $T {
                            $W
                            break
                        }
                        if name == $T {
                            break
                        }
                    }
                }
                $W
                """,
                generateAuthSchemeDetection(),
                SymbolUtils.createValueSymbolBuilder("SupportedSchemes", AwsGoDependency.INTERNAL_AUTH).build(),
                SymbolUtils.createValueSymbolBuilder("SigV4", AwsGoDependency.INTERNAL_AUTH).build(),
                generateSigV4Resolution(serviceShape),
                SymbolUtils.createValueSymbolBuilder("SigV4A", AwsGoDependency.INTERNAL_AUTH).build(),
                generateSigV4AResolution(serviceShape),
                SymbolUtils.createValueSymbolBuilder("None", AwsGoDependency.INTERNAL_AUTH).build(),
                generateAuthSchemeValidation()
            );
        }

        private GoWriter.Writable generateAuthSchemeDetection() {
            return (GoWriter writer) -> {
                writer.write(
                    """
                        authSchemes, err := $T(&resolvedEndpoint.Properties)
                        if err != nil || len(authSchemes) == 0 {
                            return out, metadata, $T(\"Failed to resolve authentication scheme\")
                        }  
                    """,
                    SymbolUtils.createValueSymbolBuilder("GetAuthenticationSchemes", AwsGoDependency.INTERNAL_AUTH).build(),
                    SymbolUtils.createValueSymbolBuilder("Errorf", SmithyGoDependency.FMT).build()

                );
            };
        }

        private GoWriter.Writable generateSigV4Resolution(ServiceShape serviceShape) {
            GoWriter.Writable signerVersion = (GoWriter writer) -> {
                if (isS3ServiceShape(serviceShape)) {
                    writer.write(
                        """
                        ctx = $T(ctx, v4Scheme.Name)
                        """,
                        SymbolUtils.createValueSymbolBuilder("SetSignerVersion", AwsCustomGoDependency.S3_CUSTOMIZATION).build()
                    );
                } else {
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
                        ctx = $T(ctx, signingName)
                        ctx = $T(ctx, signingRegion)
                        $W
                    """,
                    SymbolUtils.createPointableSymbolBuilder("AuthenticationSchemeV4", AwsGoDependency.INTERNAL_AUTH).build(),
                    SymbolUtils.createValueSymbolBuilder(signingNameDefault).build(),
                    SymbolUtils.createValueSymbolBuilder("SetSigningName", AwsGoDependency.AWS_MIDDLEWARE).build(),
                    SymbolUtils.createValueSymbolBuilder("SetSigningRegion", AwsGoDependency.AWS_MIDDLEWARE).build(),
                    signerVersion
                );
            };
        }

        private GoWriter.Writable generateSigV4AResolution(ServiceShape serviceShape) {
            GoWriter.Writable signerVersion = (GoWriter writer) -> {
                if (isS3ServiceShape(serviceShape)) {
                    writer.write(
                        """
                        ctx = $T(ctx, v4aScheme.Name)
                        """,
                        SymbolUtils.createValueSymbolBuilder("SetSignerVersion", AwsCustomGoDependency.S3_CUSTOMIZATION).build()
                    );
                } else {
                    writer.write("");
                }   
            };
            
            return (GoWriter writer) -> {
                var signingNameDefaultOpt = getDefaultSigningName(serviceShape);
                var signingNameDefault = signingNameDefaultOpt.isPresent() ? signingNameDefaultOpt.get() : "";
                writer.write(
                    """
                        v4aScheme, _ := authScheme.($P)
                        var signingName string
                        if v4aScheme.SigningName == nil {
                            signingName = \"$L\"
                        }
                        ctx = $T(ctx, signingName)
                        ctx = $T(ctx, v4aScheme.SigningRegionSet[0]) 
                        $W
                    """,
                    SymbolUtils.createPointableSymbolBuilder("AuthenticationSchemeV4A", AwsGoDependency.INTERNAL_AUTH).build(),
                    SymbolUtils.createValueSymbolBuilder(signingNameDefault).build(),
                    SymbolUtils.createValueSymbolBuilder("SetSigningName", AwsGoDependency.AWS_MIDDLEWARE).build(),
                    SymbolUtils.createValueSymbolBuilder("SetSigningRegion", AwsGoDependency.AWS_MIDDLEWARE).build(),
                    signerVersion

                );
            };
        }

        private GoWriter.Writable generateAuthSchemeValidation() {
            return (GoWriter writer) -> {
                writer.write(
                    """
                        if !supportedAuthSchemeFound {
                            return out, metadata, $T(
                                \"This operation requests signer version %s but the client only supports %v\",
                                authSchemes[0].GetName(),
                                $T,
                            )
                        }     
                    """,
                    SymbolUtils.createValueSymbolBuilder("Errorf", SmithyGoDependency.FMT).build(),
                    SymbolUtils.createValueSymbolBuilder("SupportedSchemes", AwsGoDependency.INTERNAL_AUTH).build()
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

        private final boolean isS3ServiceShape(ServiceShape serviceShape) {
            String serviceId = serviceShape.expectTrait(ServiceTrait.class).getSdkId();
            return serviceId.equalsIgnoreCase("S3");
        }
    }
    
}
