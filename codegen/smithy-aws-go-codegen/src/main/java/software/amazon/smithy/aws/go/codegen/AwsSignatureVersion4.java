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
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.ServiceIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.traits.Trait;
import software.amazon.smithy.utils.ListUtils;

/**
 * Generates Client Configuration, Middleware, and Config Resolvers for AWS Signature Version 4 support.
 */
public final class AwsSignatureVersion4 implements GoIntegration {
    public static final String SIGNER_INTERFACE_TYPE = "HTTPSignerV4";
    public static final String SIGNER_CONFIG_NAME = SIGNER_INTERFACE_TYPE;
    public static final String DEFAULT_RESOLVER_FUNCTION = "resolveHTTPSignerV4";
    public static final String REGISTER_MIDDLEWARE_FUNCTION = "registerHTTPSignerV4Middleware";

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
                writeSignerInterface(writer);
                writeDefaultResolver(writer);
                writeMiddlewareRegister(writer);
            });
        }
    }

    private void writeMiddlewareRegister(GoWriter writer) {
        writer.openBlock("func $L(stack $P, o Options) {", "}", REGISTER_MIDDLEWARE_FUNCTION,
                SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE).build(),
                () -> writer.write("stack.Finalize.Add($T(o.$L, o.$L), middleware.After)",
                        SymbolUtils.createValueSymbolBuilder("NewSignHTTPRequestMiddleware",
                                AwsGoDependency.AWS_SIGNER_V4).build(),
                        AddAwsConfigFields.CREDENTIALS_CONFIG_NAME,
                        SIGNER_CONFIG_NAME));
        writer.write("");
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(RuntimeClientPlugin.builder()
                .servicePredicate(AwsSignatureVersion4::isSupportedAuthentication)
                .addConfigField(ConfigField.builder()
                        .name(SIGNER_CONFIG_NAME)
                        .type(SymbolUtils.createValueSymbolBuilder(SIGNER_INTERFACE_TYPE).build())
                        .documentation("Provides the AWS Signature Version 4 Implementation")
                        .build())
                .resolveFunction(SymbolUtils.createValueSymbolBuilder(DEFAULT_RESOLVER_FUNCTION).build())
                .build());
    }

    private void writeDefaultResolver(GoWriter writer) {
        writer.openBlock("func $L(o *Options) {", "}", DEFAULT_RESOLVER_FUNCTION, () -> {
            writer.openBlock("if o.$L != nil {", "}", SIGNER_CONFIG_NAME, () -> writer.write("return"));
            writer.openBlock("o.$L = $T(func(s $P) {", "})", SIGNER_CONFIG_NAME,
                    SymbolUtils.createValueSymbolBuilder("NewSigner", AwsGoDependency.AWS_SIGNER_V4).build(),
                    SymbolUtils.createPointableSymbolBuilder("Signer", AwsGoDependency.AWS_SIGNER_V4).build(),
                    () -> {
                        writer.write("o.Logger = o.Logger");
                    });
        });
        writer.write("");
    }

    private void writeSignerInterface(GoWriter writer) {
        writer.openBlock("type $L interface {", "}", SIGNER_INTERFACE_TYPE, () -> {
            writer.addUseImports(SmithyGoDependency.CONTEXT);
            writer.addUseImports(AwsGoDependency.AWS_CORE);
            writer.addUseImports(SmithyGoDependency.NET_HTTP);
            writer.addUseImports(SmithyGoDependency.TIME);
            writer.write("SignHTTP(ctx context.Context, credentials aws.Credentials, r *http.Request, "
                    + "payloadHash string, service string, region string, signingTime time.Time) error");
        });
        writer.write("");
    }

    /**
     * Returns if the SigV4Trait is a auth scheme supported by the service.
     *
     * @param model model definition
     * @param serviceShape service shape for the API
     * @return if the SigV4 trait is used by the service.
     */
    public static boolean isSupportedAuthentication(Model model, ServiceShape serviceShape) {
        return model.getKnowledge(ServiceIndex.class)
                .getAuthSchemes(serviceShape).values().stream().anyMatch(trait -> trait.getClass()
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
        ServiceIndex serviceIndex = model.getKnowledge(ServiceIndex.class);
        Map<ShapeId, Trait> auth = serviceIndex.getEffectiveAuthSchemes(service.getId(), operation.getId());
        return auth.containsKey(SigV4Trait.ID);
    }
}
