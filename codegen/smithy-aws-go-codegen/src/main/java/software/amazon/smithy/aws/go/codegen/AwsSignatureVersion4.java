/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

import java.util.List;
import java.util.Map;

import software.amazon.smithy.aws.traits.auth.SigV4ATrait;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.ConfigFieldResolver;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.ServiceIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.OptionalAuthTrait;
import software.amazon.smithy.model.traits.Trait;
import software.amazon.smithy.utils.ListUtils;

/**
 * Generates Client Configuration and Config Resolvers for AWS Signature Version 4 support.
 */
public final class AwsSignatureVersion4 implements GoIntegration {
    public static final String SIGNER_INTERFACE_NAME = "HTTPSignerV4";
    public static final String SIGNER_CONFIG_FIELD_NAME = SIGNER_INTERFACE_NAME;
    public static final String NEW_SIGNER_FUNC_NAME = "newDefaultV4Signer";
    public static final String NEW_SIGNER_V4A_FUNC_NAME = "newDefaultV4aSigner";
    public static final String SIGNER_RESOLVER = "resolve" + SIGNER_CONFIG_FIELD_NAME;

    private static final List<String> DISABLE_URI_PATH_ESCAPE = ListUtils.of("com.amazonaws.s3#AmazonS3");

    @Override
    public byte getOrder() {
        return -48;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceShape serviceShape = settings.getService(model);
        if (isSupportedAuthentication(model, serviceShape)) {
            goDelegator.useShapeWriter(serviceShape, writer -> {
                writerSignerInterface(writer);
                writerConfigFieldResolver(writer, serviceShape);
                writeNewV4SignerFunc(writer, serviceShape);
            });
        }
    }

    private void writerSignerInterface(GoWriter writer) {
        writer.openBlock("type $L interface {", "}", SIGNER_INTERFACE_NAME, () -> {
            writer.addUseImports(SmithyGoDependency.CONTEXT);
            writer.addUseImports(AwsGoDependency.AWS_CORE);
            writer.addUseImports(AwsGoDependency.AWS_SIGNER_V4);
            writer.addUseImports(SmithyGoDependency.NET_HTTP);
            writer.addUseImports(SmithyGoDependency.TIME);
            writer.write("SignHTTP(ctx context.Context, credentials aws.Credentials, r *http.Request, "
                    + "payloadHash string, service string, region string, signingTime time.Time, "
                    + "optFns ...func(*v4.SignerOptions)) error");
        });
    }

    private void writerConfigFieldResolver(GoWriter writer, ServiceShape serviceShape) {
        writer.openBlock("func $L(o *Options) {", "}", SIGNER_RESOLVER, () -> {
            writer.openBlock("if o.$L != nil {", "}", SIGNER_CONFIG_FIELD_NAME, () -> writer.write("return"));
            writer.write("o.$L = $L(*o)", SIGNER_CONFIG_FIELD_NAME, NEW_SIGNER_FUNC_NAME);
        });
        writer.write("");
    }

    private void writeNewV4SignerFunc(GoWriter writer, ServiceShape serviceShape) {
        Symbol signerSymbol = SymbolUtils.createValueSymbolBuilder("Signer",
                AwsGoDependency.AWS_SIGNER_V4).build();
        Symbol newSignerSymbol = SymbolUtils.createValueSymbolBuilder("NewSigner",
                AwsGoDependency.AWS_SIGNER_V4).build();
        Symbol signerOptionsSymbol = SymbolUtils.createPointableSymbolBuilder("SignerOptions",
                AwsGoDependency.AWS_SIGNER_V4).build();

        writer.openBlock("func $L(o Options) *$T {", "}", NEW_SIGNER_FUNC_NAME, signerSymbol, () -> {
            writer.openBlock("return $T(func(so $P) {", "})", newSignerSymbol, signerOptionsSymbol, () -> {
                writer.write("so.Logger = o.$L", AddAwsConfigFields.LOGGER_CONFIG_NAME);
                writer.write("so.LogSigning = o.$L.IsSigning()", AddAwsConfigFields.LOG_MODE_CONFIG_NAME);
                if (DISABLE_URI_PATH_ESCAPE.contains(serviceShape.getId().toString())) {
                    writer.write("so.DisableURIPathEscaping = true");
                }
            });
        });
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(RuntimeClientPlugin.builder()
                .servicePredicate(AwsSignatureVersion4::isSupportedAuthentication)
                .addConfigField(ConfigField.builder()
                        .name(SIGNER_INTERFACE_NAME)
                        .type(SymbolUtils.createValueSymbolBuilder(SIGNER_INTERFACE_NAME).build())
                        .documentation("Signature Version 4 (SigV4) Signer")
                        .build())
                .addConfigFieldResolver(
                        ConfigFieldResolver.builder()
                                .location(ConfigFieldResolver.Location.CLIENT)
                                .target(ConfigFieldResolver.Target.INITIALIZATION)
                                .resolver(SymbolUtils.createValueSymbolBuilder(SIGNER_RESOLVER).build())
                                .build())
                .build());
    }

    /**
     * Returns if the SigV4Trait is a auth scheme supported by the service.
     *
     * @param model        model definition
     * @param serviceShape service shape for the API
     * @return if the SigV4 trait is used by the service.
     */
    public static boolean isSupportedAuthentication(Model model, ServiceShape serviceShape) {
        return ServiceIndex.of(model).getAuthSchemes(serviceShape).values().stream().anyMatch(trait -> trait.getClass()
                .equals(SigV4Trait.class));
    }

    /**
     * Returns if the SigV4Trait is a auth scheme for the service and operation.
     *
     * @param model     model definition
     * @param service   service shape for the API
     * @param operation operation shape
     * @return if SigV4Trait is an auth scheme for the operation and service.
     */
    public static boolean hasSigV4AuthScheme(Model model, ServiceShape service, OperationShape operation) {
        Map<ShapeId, Trait> auth = ServiceIndex.of(model).getEffectiveAuthSchemes(service.getId(), operation.getId());
        return auth.containsKey(SigV4Trait.ID) && !operation.hasTrait(OptionalAuthTrait.class);
    }

    public static boolean hasSigV4X(Model model, ServiceShape service) {
        var auth = ServiceIndex.of(model)
                .getEffectiveAuthSchemes(service.getId());
        return auth.containsKey(SigV4Trait.ID) || auth.containsKey(SigV4ATrait.ID);
    }

    public static boolean hasSigV4X(Model model, ServiceShape service, OperationShape operation) {
        var auth = ServiceIndex.of(model)
                .getEffectiveAuthSchemes(service.getId(), operation.getId());
        return auth.containsKey(SigV4Trait.ID) || auth.containsKey(SigV4ATrait.ID);
    }
}
