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
import java.util.Collection;
import java.util.HashSet;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.function.BiPredicate;
import java.util.logging.Logger;
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
import software.amazon.smithy.go.codegen.integration.auth.HttpBearerAuth;
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
    public static final String BEARER_AUTH_TOKEN_CONFIG_NAME = "TokenProvider";
    public static final String ENDPOINT_RESOLVER_CONFIG_NAME = "EndpointResolver";
    public static final String AWS_ENDPOINT_RESOLVER_WITH_OPTIONS = "EndpointResolverWithOptions";
    public static final String HTTP_CLIENT_CONFIG_NAME = "HTTPClient";
    public static final String RETRY_MAX_ATTEMPTS_CONFIG_NAME = "RetryMaxAttempts";
    public static final String RETRY_MODE_CONFIG_NAME = "RetryMode";
    public static final String RETRYER_CONFIG_NAME = "Retryer";
    public static final String API_OPTIONS_CONFIG_NAME = "APIOptions";
    public static final String LOGGER_CONFIG_NAME = "Logger";
    public static final String LOG_MODE_CONFIG_NAME = "ClientLogMode";
    public static final String DEFAULTS_MODE_CONFIG_NAME = "DefaultsMode";
    public static final String RUNTIME_ENVIRONMENT_CONFIG_NAME = "RuntimeEnvironment";

    private static final String RESOLVE_HTTP_CLIENT = "resolveHTTPClient";
    private static final String RESOLVE_RETRYER = "resolveRetryer";
    private static final String RESOLVE_AWS_CONFIG_ENDPOINT_RESOLVER = "resolveAWSEndpointResolver";
    private static final String RESOLVE_AWS_CONFIG_RETRY_MAX_ATTEMPTS = "resolveAWSRetryMaxAttempts";
    private static final String RESOLVE_AWS_CONFIG_RETRY_MODE = "resolveAWSRetryMode";
    private static final String RESOLVE_AWS_CONFIG_RETRYER_PROVIDER = "resolveAWSRetryerProvider";

    private static final String FINALIZE_RETRY_MAX_ATTEMPTS_OPTIONS = "finalizeRetryMaxAttemptOptions";

    private static final String SDK_APP_ID = "AppID";

    private static final List<AwsConfigField> AWS_CONFIG_FIELDS = ListUtils.of(
            AwsConfigField.builder()
                    .name(REGION_CONFIG_NAME)
                    .type(getUniversalSymbol("string"))
                    .documentation("The region to send requests to. (Required)")
                    .build(),
            AwsConfigField.builder()
                    .name(DEFAULTS_MODE_CONFIG_NAME)
                    .type(getAwsCoreSymbol("DefaultsMode"))
                    .documentation("""
                            The configuration DefaultsMode that the SDK should use when constructing
                            the clients initial default settings.
                            """)
                    .build(),
            AwsConfigField.builder()
                    .name(RUNTIME_ENVIRONMENT_CONFIG_NAME)
                    .type(getAwsCoreSymbol("RuntimeEnvironment"))
                    .documentation("""
                            The RuntimeEnvironment configuration, only populated if the DefaultsMode is set to
                            DefaultsModeAuto and is initialized using `config.LoadDefaultConfig`. You should not
                            populate this structure programmatically, or rely on the values here within your
                            applications.
                            """)
                    .build(),
            AwsConfigField.builder()
                    .name(RETRYER_CONFIG_NAME)
                    .type(getAwsCoreSymbol("Retryer"))
                    .documentation("""
                            Retryer guides how HTTP requests should be retried in case of recoverable failures.
                            When nil the API client will use a default retryer. The kind of default retry created
                            by the API client can be changed with the RetryMode option.
                            """)
                    .addConfigFieldResolvers(getClientInitializationResolver(
                            SymbolUtils.createValueSymbolBuilder(RESOLVE_RETRYER).build())
                            .build()
                    )
                    .awsResolveFunction(SymbolUtils.createValueSymbolBuilder(RESOLVE_AWS_CONFIG_RETRYER_PROVIDER)
                            .build())
                    .build(),
            AwsConfigField.builder()
                    .name(RETRY_MAX_ATTEMPTS_CONFIG_NAME)
                    .type(getUniversalSymbol("int"))
                    .documentation("""
                            RetryMaxAttempts specifies the maximum number attempts an API client
                            will call an operation that fails with a retryable error. A value of 0 is ignored,
                            and will not be used to configure the API client created default retryer, or modify
                            per operation call's retry max attempts.

                            When creating a new API Clients this member will only be used if the
                            Retryer Options member is nil. This value will be ignored if
                            Retryer is not nil.

                            If specified in an operation call's functional options with a value that
                            is different than the constructed client's Options, the Client's Retryer
                            will be wrapped to use the operation's specific RetryMaxAttempts value.
                            """)
                    .awsResolveFunction(SymbolUtils.createValueSymbolBuilder(RESOLVE_AWS_CONFIG_RETRY_MAX_ATTEMPTS)
                            .build())
                    .addConfigFieldResolvers(ConfigFieldResolver.builder()
                            .location(ConfigFieldResolver.Location.OPERATION)
                            .target(ConfigFieldResolver.Target.FINALIZATION)
                            .withClientInput(true)
                            .resolver(SymbolUtils.createValueSymbolBuilder(
                                    FINALIZE_RETRY_MAX_ATTEMPTS_OPTIONS).build())
                            .build())
                    .build(),

            AwsConfigField.builder()
                    .name(RETRY_MODE_CONFIG_NAME)
                    .type(getAwsCoreSymbol("RetryMode"))
                    .documentation("""
                            RetryMode specifies the retry mode the API client will be created with,
                            if Retryer option is not also specified.

                            When creating a new API Clients this member will only be used if the
                            Retryer Options member is nil. This value will be ignored if
                            Retryer is not nil.

                            Currently does not support per operation call overrides, may in the future.
                            """)
                    .awsResolveFunction(SymbolUtils.createValueSymbolBuilder(RESOLVE_AWS_CONFIG_RETRY_MODE)
                            .build())
                    .build(),
            AwsConfigField.builder()
                    .name(HTTP_CLIENT_CONFIG_NAME)
                    .type(SymbolUtils.createValueSymbolBuilder("HTTPClient").build())
                    .generatedOnClient(false)
                    .addConfigFieldResolvers(getClientInitializationResolver(
                            SymbolUtils.createValueSymbolBuilder(RESOLVE_HTTP_CLIENT).build())
                            .build())
                    .build(),
            AwsConfigField.builder()
                    .name(CREDENTIALS_CONFIG_NAME)
                    .type(getAwsCoreSymbol("CredentialsProvider"))
                    .documentation("The credentials object to use when signing requests.")
                    .servicePredicate(AwsSignatureVersion4::isSupportedAuthentication)
                    .build(),
            AwsConfigField.builder()
                    // TOKEN_PROVIDER_OPTION_NAME added API Client's Options by HttpBearerAuth. Only
                    // need to add NewFromConfig resolver from aws#Config type.
                    .name(HttpBearerAuth.TOKEN_PROVIDER_OPTION_NAME)
                    .type(SymbolUtils.createValueSymbolBuilder("TokenProvider",
                            SmithyGoDependency.SMITHY_AUTH_BEARER).build())
                    .documentation("The bearer authentication token provider for authentication requests.")
                    .servicePredicate(HttpBearerAuth::isSupportedAuthentication)
                    .generatedOnClient(false)
                    .build(),
            AwsConfigField.builder()
                    .name(API_OPTIONS_CONFIG_NAME)
                    .type(SymbolUtils.createValueSymbolBuilder("[]func(*middleware.Stack) error")
                            .addDependency(SmithyGoDependency.SMITHY_MIDDLEWARE).build())
                    .documentation("API stack mutators")
                    .generatedOnClient(false)
                    .build(),
            AwsConfigField.builder()
                    .name(ENDPOINT_RESOLVER_CONFIG_NAME)
                    .type(getAwsCoreSymbol("EndpointResolver"))
                    .generatedOnClient(false)
                    .awsResolveFunction(SymbolUtils.createValueSymbolBuilder(RESOLVE_AWS_CONFIG_ENDPOINT_RESOLVER)
                            .build())
                    .build(),
            AwsConfigField.builder()
                    .name(LOGGER_CONFIG_NAME)
                    .type(getAwsCoreSymbol("Logger"))
                    .generatedOnClient(false)
                    .build(),
            AwsConfigField.builder()
                    .name(LOG_MODE_CONFIG_NAME)
                    .type(getAwsCoreSymbol("ClientLogMode"))
                    .documentation("Configures the events that will be sent to the configured logger.")
                    .build(),
            AwsConfigField.builder()
                    .name(SDK_APP_ID)
                    .type(getUniversalSymbol("string"))
                    .documentation("The optional application specific identifier appended to the User-Agent header.")
                    .generatedOnClient(false)
                    .build()
    );

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
        goDelegator.useShapeTestWriter(serviceShape, w -> {
            writerAwsDefaultResolversTests(w);
        });
    }

    private static ConfigFieldResolver.Builder getClientInitializationResolver(Symbol resolver) {
        return ConfigFieldResolver.builder()
                .location(ConfigFieldResolver.Location.CLIENT)
                .target(ConfigFieldResolver.Target.INITIALIZATION)
                .resolver(resolver);
    }

    private void writeAwsDefaultResolvers(GoWriter writer) {
        writeHttpClientResolver(writer);
        writeRetryerResolvers(writer);
        writeRetryMaxAttemptsFinalizeResolver(writer);
        writeAwsConfigEndpointResolver(writer);
    }

    private void writerAwsDefaultResolversTests(GoWriter writer) {
        writeRetryResolverTests(writer);
    }

    private void writeRetryerResolvers(GoWriter writer) {
        writer.pushState();

        writer.putContext("resolvedDefaultsMode",
                ClientResolvedDefaultsMode.RESOLVED_DEFAULTS_MODE_CONFIG_NAME);
        writer.putContext("getConfig", SymbolUtils.createValueSymbolBuilder("GetModeConfiguration",
                AwsGoDependency.AWS_DEFAULTS).build());

        writer.putContext("resolverName", RESOLVE_RETRYER);
        writer.putContext("retryerOption", RETRYER_CONFIG_NAME);
        writer.putContext("retryModeOption", RETRY_MODE_CONFIG_NAME);
        writer.putContext("retryMaxAttemptsOption", RETRY_MAX_ATTEMPTS_CONFIG_NAME);

        writer.putContext("retryerResolveAwsConfig", RESOLVE_AWS_CONFIG_RETRYER_PROVIDER);
        writer.putContext("retryModeResolveAwsConfig", RESOLVE_AWS_CONFIG_RETRY_MODE);
        writer.putContext("retryMaxAttemptsResolveAwsConfig", RESOLVE_AWS_CONFIG_RETRY_MAX_ATTEMPTS);

        writer.putContext("retryModeAdaptive", getAwsCoreSymbol("RetryModeAdaptive"));
        writer.putContext("retryModeStandard", getAwsCoreSymbol("RetryModeStandard"));

        writer.putContext("newStandard", SymbolUtils.createValueSymbolBuilder("NewStandard",
                AwsGoDependency.AWS_RETRY).build());
        writer.putContext("standardOptions", SymbolUtils.createPointableSymbolBuilder("StandardOptions",
                AwsGoDependency.AWS_RETRY).build());
        writer.putContext("newAdaptiveMode", SymbolUtils.createPointableSymbolBuilder("NewAdaptiveMode",
                AwsGoDependency.AWS_RETRY).build());
        writer.putContext("adaptiveModeOptions", SymbolUtils.createValueSymbolBuilder("AdaptiveModeOptions",
                AwsGoDependency.AWS_RETRY).build());

        writer.write("""
                func $resolverName:L(o *Options) {
                    if o.$retryerOption:L != nil {
                        return
                    }

                    if len(o.$retryModeOption:L) == 0 {
                        modeConfig, err := $getConfig:T(o.$resolvedDefaultsMode:L)
                        if err == nil {
                            o.$retryModeOption:L = modeConfig.RetryMode
                        }
                    }
                    if len(o.$retryModeOption:L) == 0 {
                        o.$retryModeOption:L = $retryModeStandard:T
                    }

                    var standardOptions []func(*$standardOptions:T)
                    if v := o.$retryMaxAttemptsOption:L; v != 0 {
                        standardOptions = append(standardOptions, func(so *$standardOptions:T) {
                            so.MaxAttempts = v
                        })
                    }

                    switch o.$retryModeOption:L {
                    case $retryModeAdaptive:T:
                        var adaptiveOptions []func(*$adaptiveModeOptions:T)
                        if len(standardOptions) != 0 {
                            adaptiveOptions = append(adaptiveOptions, func(ao *$adaptiveModeOptions:T) {
                                ao.StandardOptions = append(ao.StandardOptions, standardOptions...)
                            })
                        }
                        o.$retryerOption:L = $newAdaptiveMode:T(adaptiveOptions...)

                    default:
                        o.$retryerOption:L = $newStandard:T(standardOptions...)
                    }
                }

                func $retryerResolveAwsConfig:L(cfg aws.Config, o *Options) {
                    if cfg.$retryerOption:L == nil {
                        return
                    }
                    o.$retryerOption:L = cfg.$retryerOption:L()
                }

                func $retryModeResolveAwsConfig:L(cfg aws.Config, o *Options) {
                    if len(cfg.$retryModeOption:L) == 0 {
                        return
                    }
                    o.$retryModeOption:L = cfg.$retryModeOption:L
                }
                func $retryMaxAttemptsResolveAwsConfig:L(cfg aws.Config, o *Options) {
                    if cfg.$retryMaxAttemptsOption:L == 0 {
                        return
                    }
                    o.$retryMaxAttemptsOption:L = cfg.$retryMaxAttemptsOption:L
                }
                """);

        writer.popState();
    }

    private void writeRetryMaxAttemptsFinalizeResolver(GoWriter writer) {
        writer.pushState();

        writer.putContext("finalizeResolveName", FINALIZE_RETRY_MAX_ATTEMPTS_OPTIONS);
        writer.putContext("withMaxAttempts", SymbolUtils.createValueSymbolBuilder("AddWithMaxAttempts",
                AwsGoDependency.AWS_RETRY).build());

        writer.write("""
                func $finalizeResolveName:L(o *Options, client Client) {
                    if v := o.RetryMaxAttempts; v == 0 || v == client.options.RetryMaxAttempts {
                        return
                    }

                    o.Retryer = $withMaxAttempts:T(o.Retryer, o.RetryMaxAttempts)
                }
                """);

        writer.popState();
    }

    private void writeRetryResolverTests(GoWriter writer) {
        writer.pushState();

        writer.putContext("retryModeOptions", "Client_resolveRetryOptions");
        writer.putContext("retryerType", getAwsCoreSymbol("Retryer"));

        writer.putContext("retryModeType", getAwsCoreSymbol("RetryMode"));
        writer.putContext("retryModeStandard", getAwsCoreSymbol("RetryModeStandard"));
        writer.putContext("retryModeAdaptive", getAwsCoreSymbol("RetryModeAdaptive"));

        writer.putContext("defaultsModeType", getAwsCoreSymbol("DefaultsMode"));
        writer.putContext("defaultsModeStandard", getAwsCoreSymbol("DefaultsModeStandard"));

        writer.putContext("ctxBackground", SymbolUtils.createValueSymbolBuilder("Background",
                SmithyGoDependency.CONTEXT).build());
        writer.putContext("stack", SymbolUtils.createValueSymbolBuilder("Stack",
                SmithyGoDependency.SMITHY_MIDDLEWARE).build());

        writer.putContext("smithyClientDoFunc", SymbolUtils.createValueSymbolBuilder("ClientDoFunc",
                SmithyGoDependency.SMITHY_HTTP_TRANSPORT).build());
        writer.putContext("httpRequest", SymbolUtils.createValueSymbolBuilder("Request",
                SmithyGoDependency.NET_HTTP).build());
        writer.putContext("httpResponse", SymbolUtils.createValueSymbolBuilder("Response",
                SmithyGoDependency.NET_HTTP).build());
        writer.putContext("httpHeader", SymbolUtils.createValueSymbolBuilder("Header",
                SmithyGoDependency.NET_HTTP).build());
        writer.putContext("newStringReader", SymbolUtils.createValueSymbolBuilder("NewReader",
                SmithyGoDependency.STRINGS).build());
        writer.putContext("nopCloser", SymbolUtils.createValueSymbolBuilder("NopCloser",
                SmithyGoDependency.IOUTIL).build());

        writer.addUseImports(SmithyGoDependency.TESTING);
        writer.write("""
                func Test$retryModeOptions:L(t *testing.T) {
                    nopClient := $smithyClientDoFunc:T(func(_ *$httpRequest:T) (*$httpResponse:T, error) {
                        return &$httpResponse:T{
                            StatusCode: 200,
                            Header: $httpHeader:T{},
                            Body: $nopCloser:T($newStringReader:T("")),
                        }, nil
                    })

                    cases := map[string]struct{
                        defaultsMode       $defaultsModeType:T
                        retryer            $retryerType:T
                        retryMaxAttempts   int
                        opRetryMaxAttempts *int
                        retryMode          $retryModeType:T
                        expectClientRetryMode   $retryModeType:T
                        expectClientMaxAttempts int
                        expectOpMaxAttempts     int
                    }{
                        "defaults": {
                            defaultsMode: $defaultsModeStandard:T,
                            expectClientRetryMode: $retryModeStandard:T,
                            expectClientMaxAttempts: 3,
                            expectOpMaxAttempts: 3,
                        },
                        "custom default retry": {
                            retryMode: $retryModeAdaptive:T,
                            retryMaxAttempts: 10,
                            expectClientRetryMode: $retryModeAdaptive:T,
                            expectClientMaxAttempts: 10,
                            expectOpMaxAttempts: 10,
                        },
                        "custom op max attempts": {
                            retryMode: $retryModeAdaptive:T,
                            retryMaxAttempts: 10,
                            opRetryMaxAttempts: aws.Int(2),
                            expectClientRetryMode: $retryModeAdaptive:T,
                            expectClientMaxAttempts: 10,
                            expectOpMaxAttempts: 2,
                        },
                        "custom op no change max attempts": {
                            retryMode: $retryModeAdaptive:T,
                            retryMaxAttempts: 10,
                            opRetryMaxAttempts: aws.Int(10),
                            expectClientRetryMode: $retryModeAdaptive:T,
                            expectClientMaxAttempts: 10,
                            expectOpMaxAttempts: 10,
                        },
                        "custom op 0 max attempts": {
                            retryMode: $retryModeAdaptive:T,
                            retryMaxAttempts: 10,
                            opRetryMaxAttempts: aws.Int(0),
                            expectClientRetryMode: $retryModeAdaptive:T,
                            expectClientMaxAttempts: 10,
                            expectOpMaxAttempts: 10,
                        },
                    }

                    for name, c := range cases {
                        t.Run(name, func(t *testing.T) {
                            client := NewFromConfig(aws.Config{
                                DefaultsMode:     c.defaultsMode,
                                Retryer:          func() func() $retryerType:T {
                                    if c.retryer == nil { return nil }

                                    return func() $retryerType:T { return c.retryer }
                                }(),
                                HTTPClient: nopClient,
                                RetryMaxAttempts: c.retryMaxAttempts,
                                RetryMode:        c.retryMode,
                            })

                            if e, a := c.expectClientRetryMode, client.options.RetryMode; e != a {
                                t.Errorf("expect %v retry mode, got %v", e, a)
                            }
                            if e, a := c.expectClientMaxAttempts, client.options.Retryer.MaxAttempts(); e != a {
                                t.Errorf("expect %v max attempts, got %v", e, a)
                            }

                            _, _, err := client.invokeOperation($ctxBackground:T(), "mockOperation", struct{}{},
                                []func(*Options){
                                    func(o *Options) {
                                        if c.opRetryMaxAttempts == nil {
                                            return
                                        }
                                        o.RetryMaxAttempts = *c.opRetryMaxAttempts
                                    },
                                },
                                func(s *$stack:T, o Options) error {
                                    s.Initialize.Clear()
                                    s.Serialize.Clear()
                                    s.Build.Clear()
                                    s.Finalize.Clear()
                                    s.Deserialize.Clear()

                                    if e, a := c.expectOpMaxAttempts, o.Retryer.MaxAttempts(); e != a {
                                        t.Errorf("expect %v op max attempts, got %v", e, a)
                                    }
                                    return nil
                                })
                            if err != nil {
                                t.Fatalf("expect no operation error, got %v", err)
                            }
                        })
                    }
                }
                """);

        writer.popState();
    }

    private void writeHttpClientResolver(GoWriter writer) {
        writer.pushState();

        writer.putContext("resolverName", RESOLVE_HTTP_CLIENT);
        writer.putContext("resolvedDefaultsMode", ClientResolvedDefaultsMode.RESOLVED_DEFAULTS_MODE_CONFIG_NAME);
        writer.putContext("optionName", HTTP_CLIENT_CONFIG_NAME);
        writer.putContext("newClient", SymbolUtils.createValueSymbolBuilder("NewBuildableClient",
                AwsGoDependency.AWS_HTTP_TRANSPORT).build());
        writer.putContext("buildableType", SymbolUtils.createPointableSymbolBuilder("BuildableClient",
                AwsGoDependency.AWS_HTTP_TRANSPORT).build());
        writer.putContext("legacyModeType", SymbolUtils.createValueSymbolBuilder("DefaultsModeLegacy",
                AwsGoDependency.AWS_CORE).build());
        writer.putContext("getConfig", SymbolUtils.createValueSymbolBuilder("GetModeConfiguration",
                AwsGoDependency.AWS_DEFAULTS).build());
        writer.putContext("dialer", SymbolUtils.createPointableSymbolBuilder("Dialer",
                SmithyGoDependency.NET).build());
        writer.putContext("transport", SymbolUtils.createPointableSymbolBuilder("Transport",
                SmithyGoDependency.NET_HTTP).build());
        writer.putContext("errorf", SymbolUtils.createPointableSymbolBuilder("Errorf",
                SmithyGoDependency.FMT).build());

        writer.write("""
                func $resolverName:L(o *Options) {
                    var buildable $buildableType:P

                    if o.$optionName:L != nil {
                        var ok bool
                        buildable, ok = o.$optionName:L.($buildableType:P)
                        if !ok {
                            return
                        }
                    } else {
                        buildable = $newClient:T()
                    }

                    modeConfig, err := $getConfig:T(o.$resolvedDefaultsMode:L)
                    if err == nil {
                        buildable = buildable.WithDialerOptions(func(dialer $dialer:P) {
                            if dialerTimeout, ok := modeConfig.GetConnectTimeout(); ok {
                                dialer.Timeout = dialerTimeout
                            }
                        })

                        buildable = buildable.WithTransportOptions(func(transport $transport:P) {
                            if tlsHandshakeTimeout, ok := modeConfig.GetTLSNegotiationTimeout(); ok {
                                transport.TLSHandshakeTimeout = tlsHandshakeTimeout
                            }
                        })
                    }

                    o.$optionName:L = buildable
                }
                """);

        writer.popState();
    }

    private void writeAwsConfigEndpointResolver(GoWriter writer) {
        writer.pushState();
        writer.putContext("resolverName", RESOLVE_AWS_CONFIG_ENDPOINT_RESOLVER);
        writer.putContext("clientOption", ENDPOINT_RESOLVER_CONFIG_NAME);
        writer.putContext("wrapperHelper", EndpointGenerator.AWS_ENDPOINT_RESOLVER_HELPER);
        writer.putContext("awsResolver", ENDPOINT_RESOLVER_CONFIG_NAME);
        writer.putContext("awsResolverWithOptions", AWS_ENDPOINT_RESOLVER_WITH_OPTIONS);
        writer.write("""
                func $resolverName:L(cfg aws.Config, o *Options) {
                    if cfg.$awsResolver:L == nil && cfg.$awsResolverWithOptions:L == nil {
                        return
                    }
                    o.$clientOption:L = $wrapperHelper:L(cfg.$awsResolver:L, cfg.$awsResolverWithOptions:L)
                }
                """);
        writer.popState();
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        List<RuntimeClientPlugin> plugins = new ArrayList<>();

        AWS_CONFIG_FIELDS.forEach(awsConfigField -> {
            RuntimeClientPlugin.Builder builder = RuntimeClientPlugin.builder();
            awsConfigField.getServicePredicate().ifPresent(
                    builder::servicePredicate);
            if (awsConfigField.isGeneratedOnClient()) {
                builder.addConfigField(awsConfigField);
            }
            builder.configFieldResolvers(awsConfigField.getConfigFieldResolvers());
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
                            if (field.getAwsResolverFunction().isPresent()) {
                                continue;
                            }
                            writer.write("$L: cfg.$L,", field.getName(), field.getName());
                        }
                    });

                    List<AwsConfigField> configFields = new ArrayList<>(AWS_CONFIG_FIELDS);
                    // add client specific config fields
                    for (AwsConfigField cfgField : ResolveClientConfigFromSources.AWS_CONFIG_FIELDS) {
                        configFields.add(cfgField);
                    }

                    for (AwsConfigField field : configFields) {
                        Optional<Symbol> awsResolverFunction = field.getAwsResolverFunction();
                        if (!awsResolverFunction.isPresent()) {
                            continue;
                        }
                        if (field.getServicePredicate().isPresent()) {
                            if (!field.getServicePredicate().get().test(model, service)) {
                                continue;
                            }
                        }
                        writer.write("$L(cfg, &opts)", awsResolverFunction.get());
                    }

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
        private final Set<ConfigFieldResolver> configFieldResolvers;
        private final Symbol awsResolveFunction;

        private AwsConfigField(Builder builder) {
            super(builder);
            this.generatedOnClient = builder.generatedOnClient;
            this.servicePredicate = builder.servicePredicate;
            this.configFieldResolvers = builder.configFieldResolvers;
            this.awsResolveFunction = builder.awsResolveFunction;
        }

        public boolean isGeneratedOnClient() {
            return generatedOnClient;
        }

        public Optional<BiPredicate<Model, ServiceShape>> getServicePredicate() {
            return Optional.ofNullable(servicePredicate);
        }

        public Set<ConfigFieldResolver> getConfigFieldResolvers() {
            return this.configFieldResolvers;
        }

        public Optional<Symbol> getAwsResolverFunction() {
            return Optional.ofNullable(awsResolveFunction);
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
            private Set<ConfigFieldResolver> configFieldResolvers = new HashSet<>();
            private Symbol awsResolveFunction = null;

            private Builder() {
                super();
            }

            /**
             * This sets the Config field on Client Options structure. By default this is true.
             * If set to false, this field won't be generated on the Client options, but will be used by
             * the NewFromConfig (to copy values from the aws config to client options).
             *
             * @param generatedOnClient bool indicating config field generation on client option structure
             * @return
             */
            public Builder generatedOnClient(boolean generatedOnClient) {
                this.generatedOnClient = generatedOnClient;
                return this;
            }

            public Builder servicePredicate(BiPredicate<Model, ServiceShape> servicePredicate) {
                this.servicePredicate = servicePredicate;
                return this;
            }

            public Builder configFieldResolvers(Collection<ConfigFieldResolver> configFieldResolvers) {
                this.configFieldResolvers = new HashSet<>(configFieldResolvers);
                return this;
            }

            public Builder addConfigFieldResolvers(ConfigFieldResolver configFieldResolver) {
                this.configFieldResolvers.add(configFieldResolver);
                return this;
            }

            public Builder awsResolveFunction(Symbol awsResolveFunction) {
                this.awsResolveFunction = awsResolveFunction;
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

            @Override
            public Builder withHelper(Boolean withHelper) {
                super.withHelper(withHelper);
                return this;
            }

            @Override
            public Builder withHelper() {
                super.withHelper();
                return this;
            }
        }
    }
}
