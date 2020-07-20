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
import java.util.function.BiConsumer;
import java.util.function.BiPredicate;
import java.util.function.Consumer;
import java.util.function.Predicate;
import java.util.logging.Logger;
import java.util.stream.Collectors;
import software.amazon.smithy.aws.traits.auth.SigV4Trait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.ServiceIndex;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.SetUtils;

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
    public static final String HTTP_SIGNER_CONFIG_NAME = "HTTPSigner";

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
                    .build(),
            AwsConfigField.builder()
                    .name(HTTP_SIGNER_CONFIG_NAME)
                    .type(getAwsSignerV4Symbol("HTTPSigner"))
                    .documentation("HTTPSigner provides AWS request signing for HTTP requests made\n"
                            + "from the client. When nil the API client will use a default signer.")
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
                    .build(),
            AwsConfigField.builder()
                    .name(CREDENTIALS_CONFIG_NAME)
                    .type(getAwsCoreSymbol("CredentialsProvider"))
                    .documentation("The credentials object to use when signing requests.")
                    .servicePredicate((model, serviceShape) -> model.getKnowledge(ServiceIndex.class)
                            .getAuthSchemes(serviceShape).values().stream().anyMatch(trait -> trait.getClass()
                                    .equals(SigV4Trait.class)))
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
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {
        LOGGER.info("generating aws.Config based client constructor");
        writerFactory.accept("api_client.go", settings.getModuleName(), w -> {
            writeAwsConfigConstructor(model, model.expectShape(settings.getService()).asServiceShape().get(), w);
        });
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        List<RuntimeClientPlugin> plugins = new ArrayList<>();

        // Collect fields that have no service predicate into a single runtime client plugin
        List<ConfigField> allClients = AWS_CONFIG_FIELDS.stream().filter(AwsConfigField::isGeneratedOnClient)
                .filter(field -> !field.getServicePredicate().isPresent())
                .collect(Collectors.toList());
        plugins.add(RuntimeClientPlugin.builder().configFields(allClients).build());

        // For each service predicate construct runtime client plugins for the field
        AWS_CONFIG_FIELDS.stream().filter(AwsConfigField::isGeneratedOnClient)
                .filter(field -> field.getServicePredicate().isPresent())
                .forEach(field -> {
                    RuntimeClientPlugin.Builder builder = RuntimeClientPlugin.builder()
                            .configFields(ListUtils.of(field))
                            .servicePredicate(field.getServicePredicate().get());
                    field.getResolverFunction().ifPresent(builder::resolveFunction);
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
                            if (field.getName().equals(HTTP_SIGNER_CONFIG_NAME)) {
                                // TODO signer does not exist in the aws.Config.
                                continue;
                            }
                            writer.write("$L: cfg.$L,", field.getName(), field.getName());
                        }
                    });
                    writer.write("return New(opts, optFns...)");
                });
        writer.write("");
    }

    private static class AwsConfigField extends ConfigField {
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

        private static class Builder extends ConfigField.Builder {
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
