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

import java.util.Optional;
import java.util.Set;
import java.util.TreeSet;
import java.util.function.Consumer;
import java.util.stream.Collectors;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.SetUtils;
import software.amazon.smithy.utils.SmithyBuilder;

/**
 * Registers additional AWS specific client configuration fields
 */
public class AddAwsConfigFields implements GoIntegration {
    private static final String REGION_CONFIG_NAME = "Region";
    private static final String CREDENTIALS_CONFIG_NAME = "Credentials";
    private static final String ENDPOINT_RESOLVER_CONFIG_NAME = "EndpointResolver";
    private static final String HTTP_CLIENT_CONFIG_NAME = "HTTPClient";
    private static final String LOGGER_CONFIG_NAME = "Logger";
    private static final String LOG_LEVEL_CONFIG_NAME = "LogLevel";
    private static final String RETRYER_CONFIG_NAME = "Retryer";
    private static final Set<ConfigField> UNIVERSAL_FIELDS = new TreeSet<>(SetUtils.of(
            ConfigField.builder(REGION_CONFIG_NAME, getUniversalSymbol("string"))
                    .documentation("The region to send requests to. (Required)")
                    .build(),
            ConfigField.builder(CREDENTIALS_CONFIG_NAME, getAwsCoreSymbol("CredentialsProvider"))
                    .documentation("The credentials object to use when signing requests.")
                    .build(),
            ConfigField.builder(ENDPOINT_RESOLVER_CONFIG_NAME, getAwsCoreSymbol("EndpointResolver"))
                    .documentation("The resolver to use for looking up endpoints for the service.")
                    .build(),
            ConfigField.builder(HTTP_CLIENT_CONFIG_NAME, getAwsCoreSymbol("HTTPClient"), false)
                    .build(),
            ConfigField.builder(RETRYER_CONFIG_NAME, getAwsCoreSymbol("Retryer"))
                    .documentation("Retryer guides how HTTP requests should be retried in case of\n"
                            + "recoverable failures. When nil the API client will use a default\n"
                            + "retryer.")
                    .build(),
            ConfigField.builder(LOG_LEVEL_CONFIG_NAME, getAwsCoreSymbol("LogLevel"))
                    .documentation("An integer value representing the logging level.")
                    .build(),
            ConfigField.builder(LOGGER_CONFIG_NAME, getAwsCoreSymbol("Logger"))
                    .documentation("The logger writer interface to write logging messages to.")
                    .build()
    ));

    private static Symbol getAwsCoreSymbol(String symbolName) {
        return SymbolUtils.createPointableSymbolBuilder(symbolName,
                AwsGoDependency.AWS_CORE).build();
    }

    private static Symbol getUniversalSymbol(String symbolName) {
        return SymbolUtils.createPointableSymbolBuilder(symbolName)
                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true).build();
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {
        Set<ConfigField> fields = new TreeSet<>(UNIVERSAL_FIELDS);


        writerFactory.accept("api_client.go", settings.getModuleName(), this::writeAwsConfigConstructor);
    }

    private void writeAwsConfigConstructor(GoWriter writer) {
        writer.writeDocs("NewFromConfig returns a new client from the provided config.");
        writer.openBlock("func NewFromConfig(cfg $T, optFns ... func(*Options)) *Client {", "}",
                getAwsCoreSymbol("Config"), () -> {
                    writer.openBlock("opts := Options{", "}", () -> {
                        UNIVERSAL_FIELDS.forEach(configField -> {
                            writer.write("$L: cfg.$L,", configField.getClientConfigField(),
                                    configField.getAwsConfigField());
                        });
                    });
                    writer.write("");
                    writer.openBlock("for _, fn := range optFns {", "}", () -> {
                        writer.write("fn(&opts)");
                    });
                    writer.write("return New(opts)");
                });
        writer.write("");
    }

    @Override
    public void addConfigFields(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoWriter writer
    ) {
        writeUniversalClientConfigFields(writer);
    }

    private void writeUniversalClientConfigFields(GoWriter writer) {
        UNIVERSAL_FIELDS.forEach(configField -> {
            if (!configField.isGenerateOnClient()) {
                return;
            }
            configField.getDocumentation().ifPresent(writer::writeDocs);
            writer.write("$L $T", configField.getClientConfigField(), configField.getTypeSymbol());
            writer.write("");
        });
    }

    /**
     * Represents a 1-1 config field relationship on the AWS SDK Config type and the service client options.
     */
    public static final class ConfigField implements Comparable<ConfigField> {
        private final String awsConfigField;
        private final String clientConfigField;
        private final String documentation;
        private final Symbol typeSymbol;
        private final boolean generateOnClient;

        private ConfigField(Builder builder) {
            this.awsConfigField = SmithyBuilder.requiredState("awsConfigField", builder.awsConfigField);
            this.clientConfigField = SmithyBuilder.requiredState("clientConfigField", builder.clientConfigField);
            this.typeSymbol = SmithyBuilder.requiredState("typeSymbol", builder.typeSymbol);
            this.documentation = builder.documentation;
            this.generateOnClient = SmithyBuilder.requiredState("generateOnClient", builder.generateOnClient);
        }

        /**
         * Get the aws client config field name.
         *
         * @return the aws config field name
         */
        public String getAwsConfigField() {
            return awsConfigField;
        }

        /**
         * Get the client config field name.
         *
         * @return the client config field name
         */
        public String getClientConfigField() {
            return clientConfigField;
        }

        /**
         * Get the documentation string for the client field.
         *
         * @return the documentation string, may return null if no documentation is present
         */
        public Optional<String> getDocumentation() {
            return Optional.ofNullable(documentation);
        }

        /**
         * Get the symbol of the config field type.
         *
         * @return the config field type symbol
         */
        public Symbol getTypeSymbol() {
            return typeSymbol;
        }

        /**
         * Returns whether the config field is present on the client by default.
         *
         * @return whether the field is present on the client
         */
        public boolean isGenerateOnClient() {
            return generateOnClient;
        }

        /**
         * Returns a builder for a {@link ConfigField}
         *
         * @return the builder
         */
        public static Builder builder() {
            return new Builder();
        }

        /**
         * Returns a builder for a {@link ConfigField} using the provided name and type symbol.
         * By default the builder will configure the field to be generated on the client options.
         *
         * @param fieldName  the field name for the aws config and client config
         * @param typeSymbol the type symbol
         * @return the builder
         */
        public static Builder builder(String fieldName, Symbol typeSymbol) {
            return builder(fieldName, typeSymbol, true);
        }

        /**
         * Returns a builder for a {@link ConfigField} using the provided name and type symbol
         * By default the builder will configure the field to be generated on the client options.
         *
         * @param fieldName        the field name for the aws config and client config
         * @param typeSymbol       the type symbol
         * @param generateOnClient whether the field should be generated on the client options
         * @return the builder
         */
        public static Builder builder(String fieldName, Symbol typeSymbol, boolean generateOnClient) {
            return builder(fieldName, fieldName, typeSymbol, generateOnClient);
        }

        private static Builder builder(
                String awsConfigField,
                String clientConfigField,
                Symbol typeSymbol,
                boolean generateOnClient
        ) {
            return ConfigField.builder()
                    .awsConfigField(awsConfigField)
                    .clientConfigField(clientConfigField)
                    .typeSymbol(typeSymbol)
                    .generateOnClient(generateOnClient);
        }

        @Override
        public int compareTo(ConfigField o) {
            return this.getClientConfigField().compareTo(o.getClientConfigField());
        }

        private static class Builder implements SmithyBuilder<ConfigField> {
            private String awsConfigField;
            private String clientConfigField;
            private String documentation;
            private Symbol typeSymbol;
            private boolean generateOnClient;

            private Builder() {
            }

            public Builder awsConfigField(String awsConfigField) {
                this.awsConfigField = awsConfigField;
                return this;
            }

            public Builder clientConfigField(String clientConfigField) {
                this.clientConfigField = clientConfigField;
                return this;
            }

            public Builder typeSymbol(Symbol typeSymbol) {
                this.typeSymbol = typeSymbol;
                return this;
            }

            public Builder documentation(String documentation) {
                this.documentation = documentation;
                return this;
            }

            public Builder generateOnClient(boolean generateOnClient) {
                this.generateOnClient = generateOnClient;
                return this;
            }

            @Override
            public ConfigField build() {
                return new ConfigField(this);
            }
        }
    }
}
