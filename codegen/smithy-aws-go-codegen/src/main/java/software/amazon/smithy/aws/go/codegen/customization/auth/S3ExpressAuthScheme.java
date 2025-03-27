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
import software.amazon.smithy.aws.go.codegen.customization.service.s3.S3ModelUtils;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.AuthSchemeDefinition;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.traits.AuthTrait;
import software.amazon.smithy.utils.ListUtils;

import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.function.Consumer;

import static software.amazon.smithy.go.codegen.GoWriter.emptyGoTemplate;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

public class S3ExpressAuthScheme implements GoIntegration {
    private static final ConfigField s3ExpressCredentials =
            ConfigField.builder()
                    .name("ExpressCredentials")
                    .type(SymbolUtils.createValueSymbolBuilder("ExpressCredentialsProvider").build())
                    .documentation("The credentials provider for S3Express requests.")
                    .build();

    private static final ConfigFieldResolver s3ExpressCredentialsResolver =
            ConfigFieldResolver.builder()
                    .location(ConfigFieldResolver.Location.CLIENT)
                    .target(ConfigFieldResolver.Target.FINALIZATION)
                    .resolver(SymbolUtils.createValueSymbolBuilder("resolveExpressCredentials").build())
                    .build();

    private static final ConfigFieldResolver s3ExpressCredentialsClientFinalizer =
            ConfigFieldResolver.builder()
                    .location(ConfigFieldResolver.Location.CLIENT)
                    .target(ConfigFieldResolver.Target.FINALIZATION_WITH_CLIENT)
                    .resolver(SymbolUtils.createValueSymbolBuilder("finalizeExpressCredentials").build())
                    .withClientInput(true)
                    .build();

    private static final ConfigFieldResolver s3ExpressCredentialsOperationFinalizer =
            ConfigFieldResolver.builder()
                    .location(ConfigFieldResolver.Location.OPERATION)
                    .target(ConfigFieldResolver.Target.FINALIZATION)
                    .resolver(buildPackageSymbol("finalizeOperationExpressCredentials"))
                    .withClientInput(true)
                    .build();

    @Override
    public void writeAdditionalFiles(
            GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator
    ) {
        if (S3ModelUtils.isServiceS3(model, settings.getService(model))) {
            goDelegator.useFileWriter("options.go", settings.getModuleName(), generateGetIdentityResolver());
        }
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3ModelUtils::isServiceS3)
                        .addConfigField(s3ExpressCredentials)
                        .addConfigFieldResolver(s3ExpressCredentialsResolver)
                        .addConfigFieldResolver(s3ExpressCredentialsClientFinalizer)
                        .addConfigFieldResolver(s3ExpressCredentialsOperationFinalizer)
                        .addAuthSchemeDefinition(SigV4S3ExpressTrait.ID, new SigV4S3Express())
                        .build()
        );
    }

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        ServiceShape service = settings.getService(model);
        if (!S3ModelUtils.isServiceS3(model, service)) {
            return model;
        }

        AuthTrait newAuth = new AuthTrait(buildSet(s -> {
            s.addAll(service.expectTrait(AuthTrait.class).getValueSet());
            s.add(SigV4S3ExpressTrait.ID);
        }));
        return model.toBuilder()
                .addShape(
                        service.toBuilder()
                                .addTrait(new SigV4S3ExpressTrait())
                                .addTrait(newAuth)
                                .build()
                )
                .build();
    }

    private static <V> Set<V> buildSet(Consumer<Set<V>> builder) {
        Set<V> s = new HashSet<>();
        builder.accept(s);
        return Collections.unmodifiableSet(s);
    }

    private static class SigV4S3Express implements AuthSchemeDefinition {
        @Override
        public GoWriter.Writable generateServiceOption(ProtocolGenerator.GenerationContext context, ServiceShape service) {
            return emptyGoTemplate(); // not modeled
        }

        @Override
        public GoWriter.Writable generateOperationOption(ProtocolGenerator.GenerationContext context, OperationShape operation) {
            return emptyGoTemplate(); // not modeled
        }

        @Override
        public GoWriter.Writable generateDefaultAuthScheme() {
            return goTemplate("""
                    $T($S, &$T{
                        Signer: options.HTTPSignerV4,
                        Logger: options.Logger,
                        LogSigning: options.ClientLogMode.IsSigning(),
                    })""",
                    SdkGoTypes.Internal.Auth.NewHTTPAuthScheme,
                    SigV4S3ExpressTrait.ID.toString(),
                    SdkGoTypes.ServiceCustomizations.S3.ExpressSigner);
        }

        @Override
        public GoWriter.Writable generateOptionsIdentityResolver() {
            return goTemplate("getExpressIdentityResolver(o)");
        }
    }

    private GoWriter.Writable generateGetIdentityResolver() {
        return goTemplate("""
                func getExpressIdentityResolver(o Options) $T {
                    if o.ExpressCredentials != nil {
                        return &$T{Provider: o.ExpressCredentials}
                    }
                    return nil
                }
                """,
                SmithyGoTypes.Auth.IdentityResolver,
                SdkGoTypes.ServiceCustomizations.S3.ExpressIdentityResolver);
    }
}
