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

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.function.BiPredicate;
import java.util.logging.Logger;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

/**
 * Registers additional AWS specific client configuration fields
 */
public class AddAwsConfigFields implements GoIntegration {
    private static final Logger LOGGER = Logger.getLogger(AddAwsConfigFields.class.getName());

    public static final String REGION_CONFIG_NAME = "Region";
    public static final String CREDENTIALS_CONFIG_NAME = "Credentials";
    public static final String ENDPOINT_RESOLVER_CONFIG_NAME = "EndpointResolver";
    public static final String HTTP_CLIENT_CONFIG_NAME = "HTTPClient";
    public static final String LOGGER_CONFIG_NAME = "Logger";
    public static final String LOG_LEVEL_CONFIG_NAME = "LogLevel";
    public static final String RETRYER_CONFIG_NAME = "Retryer";

    private static final String RESOLVE_HTTP_CLIENT = "resolveHTTPClient";
    private static final String RESOLVE_RETRYER = "resolveRetryer";

    private static final List<AwsConfigField> AWS_CONFIG_FIELDS = ListUtils.of(
            AwsConfigField.builder()
                    .name(REGION_CONFIG_NAME)
                    .type(getUniversalSymbol("string"))
                    .documentation("The region to send requests to. (Required)")
                    .build(),
            AwsConfigField.builder()
                    .name(RETRYER_CONFIG_NAME)
                    .type(getAwsRetrySymbol("Retryer"))
                    .documentation("Retryer guides how HTTP requests should be retried in case of\n"
                            + "recoverable failures. When nil the API client will use a default\n"
                            + "retryer.")
                    .resolveFunction(SymbolUtils.createValueSymbolBuilder(RESOLVE_RETRYER).build())
                    .build(),
            AwsConfigField.builder()
                    .name(LOG_LEVEL_CONFIG_NAME)
                    .type(getAwsCoreSymbol("LogLevel"))
                    .documentation("An integer value representing the logging level.")
                    .build(),
            AwsConfigField.builder()
                    .name(LOGGER_CONFIG_NAME)
                    .type(getAwsCoreSymbol("Logger"))
                    .documentation("The logger writer interface to write logging messages to.")
                    .build(),
            AwsConfigField.builder()
                    .name(HTTP_CLIENT_CONFIG_NAME)
                    .type(SymbolUtils.createValueSymbolBuilder("HTTPClient").build())
                    .generatedOnClient(false)
                    .resolveFunction(SymbolUtils.createValueSymbolBuilder(RESOLVE_HTTP_CLIENT).build())
                    .build(),
            AwsConfigField.builder()
                    .name(CREDENTIALS_CONFIG_NAME)
                    .type(getAwsCoreSymbol("CredentialsProvider"))
                    .documentation("The credentials object to use when signing requests.")
                    .servicePredicate(AwsSignatureVersion4::isSupportedAuthentication)
                    .build()
    );

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first.
     *
     * @return Returns the sort order, defaults to -50.
     */
    @Override
    public byte getOrder() {
        return -50;
    }

    private static Symbol getAwsCoreSymbol(String symbolName) {
        return SymbolUtils.createValueSymbolBuilder(symbolName,
                AwsGoDependency.AWS_CORE).build();
    }

    private static Symbol getAwsSignerV4Symbol(String symbolName) {
        return SymbolUtils.createValueSymbolBuilder(symbolName,
                AwsGoDependency.AWS_SIGNER_V4).build();
    }

    private static Symbol getUniversalSymbol(String symbolName) {
        return SymbolUtils.createValueSymbolBuilder(symbolName)
                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true).build();
    }

    private static Symbol getAwsRetrySymbol(String symbolName) {
        return SymbolUtils.createValueSymbolBuilder(symbolName,
                AwsGoDependency.AWS_RETRY).build();
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        LOGGER.info("generating aws.Config based client constructor");
        ServiceShape serviceShape = settings.getService(model);
        goDelegator.useShapeWriter(serviceShape, w -> {
            writeAwsConfigConstructor(model, serviceShape, w);
            writeAwsDefaultResolvers(w);
        });
    }

    private void writeAwsDefaultResolvers(GoWriter writer) {
        writeHttpClientResolver(writer);
        writeRetryerResolver(writer);
    }

    private void writeRetryerResolver(GoWriter writer) {
        writer.openBlock("func $L(o *Options) {", "}", RESOLVE_RETRYER, () -> {
            writer.openBlock("if o.$L != nil {", "}", RETRYER_CONFIG_NAME, () -> writer.write("return"));
            writer.write("o.$L = $T()", RETRYER_CONFIG_NAME, SymbolUtils.createValueSymbolBuilder("NewStandard",
                    AwsGoDependency.AWS_RETRY).build());
        });
        writer.write("");
    }

    private void writeHttpClientResolver(GoWriter writer) {
        writer.openBlock("func $L(o *Options) {", "}", RESOLVE_HTTP_CLIENT, () -> {
            writer.openBlock("if o.$L != nil {", "}", HTTP_CLIENT_CONFIG_NAME, () -> writer.write("return"));
            writer.write("o.$L = $T()", HTTP_CLIENT_CONFIG_NAME,
                    SymbolUtils.createValueSymbolBuilder("NewBuildableHTTPClient", AwsGoDependency.AWS_CORE).build());
        });
        writer.write("");
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        List<RuntimeClientPlugin> plugins = new ArrayList<>();

        AWS_CONFIG_FIELDS.forEach(awsConfigField -> {
            RuntimeClientPlugin.Builder builder = RuntimeClientPlugin.builder();
            awsConfigField.getServicePredicate().ifPresent(
                    modelServiceShapeBiPredicate -> builder.servicePredicate(modelServiceShapeBiPredicate));
            if (awsConfigField.isGeneratedOnClient()) {
                builder.addConfigField(awsConfigField);
            }
            awsConfigField.getResolverFunction().ifPresent(symbol -> {
                builder.resolveFunction(symbol);
            });
            plugins.add(builder.build());
        });

        return plugins;
    }

    private void writeAwsConfigConstructor(Model model, ServiceShape service, GoWriter writer) {
        writer.writeDocs("NewFromConfig returns a new client from the provided config.");
        writer.openBlock("func NewFromConfig(cfg $T, optFns ... func(*Options)) *Client {", "}",
                getAwsCoreSymbol("Config"), () -> {
                    writer.openBlock("opts := Options{", "}", () -> {
                        for (AwsConfigField field : AWS_CONFIG_FIELDS) {
                            if (field.getServicePredicate().isPresent()) {
                                if (!field.getServicePredicate().get().test(model, service)) {
                                    continue;
                                }
                            }
                            writer.write("$L: cfg.$L,", field.getName(), field.getName());
                        }
                    });
                    writer.write("return New(opts, optFns...)");
                });
        writer.write("");
    }

    /**
     * Provides configuration field for AWS client.
     */
    public static class AwsConfigField extends ConfigField {
        private final boolean generatedOnClient;
        private final BiPredicate<Model, ServiceShape> servicePredicate;
        private final Symbol resolveFunction;

        private AwsConfigField(Builder builder) {
            super(builder);
            this.generatedOnClient = builder.generatedOnClient;
            this.servicePredicate = builder.servicePredicate;
            this.resolveFunction = builder.resolveFunction;
        }

        public boolean isGeneratedOnClient() {
            return generatedOnClient;
        }

        public Optional<BiPredicate<Model, ServiceShape>> getServicePredicate() {
            return Optional.ofNullable(servicePredicate);
        }

        public Optional<Symbol> getResolverFunction() {
            return Optional.ofNullable(resolveFunction);
        }

        public static Builder builder() {
            return new Builder();
        }

        /**
         * Provides builder for AWSConfigFile values.
         */
        public static class Builder extends ConfigField.Builder {
            private boolean generatedOnClient = true;
            private BiPredicate<Model, ServiceShape> servicePredicate = null;
            private Symbol resolveFunction = null;

            private Builder() {
                super();
            }

            public Builder generatedOnClient(boolean generatedOnClient) {
                this.generatedOnClient = generatedOnClient;
                return this;
            }

            public Builder servicePredicate(BiPredicate<Model, ServiceShape> servicePredicate) {
                this.servicePredicate = servicePredicate;
                return this;
            }

            public Builder resolveFunction(Symbol resolveFunction) {
                this.resolveFunction = resolveFunction;
                return this;
            }

            @Override
            public AwsConfigField build() {
                return new AwsConfigField(this);
            }

            @Override
            public Builder name(String name) {
                super.name(name);
                return this;
            }

            @Override
            public Builder type(Symbol type) {
                super.type(type);
                return this;
            }

            @Override
            public Builder documentation(String documentation) {
                super.documentation(documentation);
                return this;
            }
        }
    }
}
