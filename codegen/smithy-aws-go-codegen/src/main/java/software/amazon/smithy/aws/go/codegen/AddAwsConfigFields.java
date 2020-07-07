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
import java.util.function.Consumer;
import java.util.logging.Logger;
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


    private static final List<ConfigField> UNIVERSAL_FIELDS = new ArrayList<>(SetUtils.of(
            ConfigField.builder()
                    .name(REGION_CONFIG_NAME)
                    .type(getUniversalSymbol("string"))
                    .documentation("The region to send requests to. (Required)")
                    .build(),
            ConfigField.builder()
                    .name(ENDPOINT_RESOLVER_CONFIG_NAME)
                    .type(getAwsCoreSymbol("EndpointResolver"))
                    .documentation("The resolver to use for looking up endpoints for the service.")
                    .build(),
            ConfigField.builder()
                    .name(RETRYER_CONFIG_NAME)
                    .type(getAwsRetrySymbol("Retryer"))
                    .documentation("Retryer guides how HTTP requests should be retried in case of\n"
                            + "recoverable failures. When nil the API client will use a default\n"
                            + "retryer.")
                    .build(),
            ConfigField.builder()
                    .name(HTTP_SIGNER_CONFIG_NAME)
                    .type(getAwsSignerV4Symbol("HTTPSigner"))
                    .documentation("HTTPSigner provides AWS request signing for HTTP requests made\n"
                            + "from the client. When nil the API client will use a default signer.")
                    .build(),
            ConfigField.builder()
                    .name(LOG_LEVEL_CONFIG_NAME)
                    .type(getAwsCoreSymbol("LogLevel"))
                    .documentation("An integer value representing the logging level.")
                    .build(),
            ConfigField.builder()
                    .name(LOGGER_CONFIG_NAME)
                    .type(getAwsCoreSymbol("Logger"))
                    .documentation("The logger writer interface to write logging messages to.")
                    .build()
    ));

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
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .configFields(UNIVERSAL_FIELDS)
                        .build(),
                RuntimeClientPlugin.builder()
                        .configFields(ListUtils.of(ConfigField.builder()
                                .name(CREDENTIALS_CONFIG_NAME)
                                .type(getAwsCoreSymbol("CredentialsProvider"))
                                .documentation("The credentials object to use when signing requests.")
                                .build()))
                        .servicePredicate((model, serviceShape) -> model.getKnowledge(ServiceIndex.class)
                                .getAuthSchemes(serviceShape).values().stream().anyMatch(trait -> trait.getClass()
                                        .equals(SigV4Trait.class)))
                        .build()
        );
    }

    private void writeAwsConfigConstructor(Model model, ServiceShape service, GoWriter writer) {
        writer.writeDocs("NewFromConfig returns a new client from the provided config.");
        writer.openBlock("func NewFromConfig(cfg $T, optFns ... func(*Options)) *Client {", "}",
                getAwsCoreSymbol("Config"), () -> {
                    writer.openBlock("opts := Options{", "}", () -> {
                        writer.write("$L: cfg.$L,", HTTP_CLIENT_CONFIG_NAME, HTTP_CLIENT_CONFIG_NAME);
                        for (RuntimeClientPlugin plugin : getClientPlugins()) {
                            if (!plugin.matchesService(model, service)) {
                                continue;
                            }
                            plugin.getConfigFields().forEach(configField -> {
                                if (configField.getName().equals(HTTP_SIGNER_CONFIG_NAME)) {
                                    // TODO signer does not exist in the aws.Config.
                                    return;
                                }
                                writer.write("$L: cfg.$L,", configField.getName(), configField.getName());
                            });
                        }
                    });
                    writer.write("");
                    writer.openBlock("for _, fn := range optFns {", "}", () -> {
                        writer.write("fn(&opts)");
                    });
                    writer.write("return New(opts)");
                });
        writer.write("");
    }
}
