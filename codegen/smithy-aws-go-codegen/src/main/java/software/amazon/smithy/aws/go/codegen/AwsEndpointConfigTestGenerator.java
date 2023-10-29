package software.amazon.smithy.aws.go.codegen;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.TriConsumer;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;

import java.util.Map;
import java.util.function.Consumer;


/*
 * This class generates service specific tests for resolving configured endpoints.
 */
public class AwsEndpointConfigTestGenerator implements GoIntegration {

    private Map<String, Object> commonCodegenArgs;


    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            TriConsumer<String, String, Consumer<GoWriter>> writerFactory
    ) {
        String sdkId = settings.getService(model).expectTrait(ServiceTrait.class).getSdkId();
        this.commonCodegenArgs = MapUtils.of(
            "envSdkId", sdkId.toUpperCase().replaceAll(" ", "_"),
            "configSdkId", sdkId.toLowerCase().replaceAll(" ", "_"),
            "urlSdkId", sdkId.toLowerCase().replaceAll(" ", "-"),
            "testing", SymbolUtils.createPointableSymbolBuilder("T", SmithyGoDependency.TESTING).build(),
            "awsString", SymbolUtils.createValueSymbolBuilder("String", AwsGoDependency.AWS_CORE).build(),
            "context", SymbolUtils.createValueSymbolBuilder("Context", SmithyGoDependency.CONTEXT).build()
        );

        writerFactory.accept("endpoints_config_test.go", settings.getModuleName(), writer -> {
            writer.write("$W", generate());
        });

    }

    private GoWriter.Writable generate() {
        return (GoWriter w) -> {
            w.write(
                """
                $W

                $W
                """,
                generateMockProviders(),
                generateTestFunction()
            );
        };
    }

    private GoWriter.Writable generateMockProviders() {
        return goTemplate(
            """
                type mockConfigSource struct {
                    global string
                    service string
                    ignore bool
                }

                // GetIgnoreConfiguredEndpoints is used in knowing when to disable configured
                // endpoints feature.
                func (m mockConfigSource) GetIgnoreConfiguredEndpoints($context:T) (bool, bool, error) {
                    return m.ignore, m.ignore, nil
                }
                
                // GetServiceBaseEndpoint is used to retrieve a normalized SDK ID for use
                // with configured endpoints.
                func (m mockConfigSource) GetServiceBaseEndpoint(ctx $context:T, sdkID string) (string, bool, error) {
                    if m.service != "" {
                        return m.service, true, nil
                    }
                    return "", false, nil
                }
                
            """,
            this.commonCodegenArgs
        );
    }


    private GoWriter.Writable generateTestFunction() {
        return (GoWriter w) -> {
            w.write(
                """
                func TestResolveBaseEndpoint(t $P) {
                    $W

                    $W
                }
                """,
                SymbolUtils.createPointableSymbolBuilder("T", SmithyGoDependency.TESTING).build(),
                generateCases(),
                generateTests()
            );
        };
    }

    private GoWriter.Writable generateCases() {
        return goTemplate(
            """
                cases := map[string]struct {
                    envGlobal        string
                    envService       string
                    envIgnore        bool
                    configGlobal     string
                    configService    string
                    configIgnore     bool
                    clientEndpoint   *string
                    expectURL        *string
                }{
                    "env ignore": {
                        envGlobal: "https://env-global.dev",
                        envService: "https://env-$urlSdkId:L.dev",
                        envIgnore: true,
                        configGlobal: "http://config-global.dev",
                        configService: "http://config-$urlSdkId:L.dev",
                        expectURL: nil,
                    },
                    "env global": {
                        envGlobal: "https://env-global.dev",
                        configGlobal: "http://config-global.dev",
                        configService: "http://config-$urlSdkId:L.dev",
                        expectURL: aws.String("https://env-global.dev"),
                    },
                    "env service": {
                        envGlobal: "https://env-global.dev",
                        envService: "https://env-$urlSdkId:L.dev",
                        configGlobal: "http://config-global.dev",
                        configService: "http://config-$urlSdkId:L.dev",
                        expectURL: aws.String("https://env-$urlSdkId:L.dev"),
                    },
                    "config ignore": {
                        envGlobal:     "https://env-global.dev",
                        envService:    "https://env-$urlSdkId:L.dev",
                        configGlobal:  "http://config-global.dev",
                        configService: "http://config-$urlSdkId:L.dev",
                        configIgnore: true,
                        expectURL:     nil,
                    },
                    "config global": {
                        configGlobal:  "http://config-global.dev",
                        expectURL:     aws.String("http://config-global.dev"),
                    },
                    "config service": {
                        configGlobal:  "http://config-global.dev",
                        configService: "http://config-$urlSdkId:L.dev",
                        expectURL:     aws.String("http://config-$urlSdkId:L.dev"),
                    },
                    "client": {
                        envGlobal:     "https://env-global.dev",
                        envService:    "https://env-$urlSdkId:L.dev",
                        configGlobal:  "http://config-global.dev",
                        configService: "http://config-$urlSdkId:L.dev",
                        clientEndpoint: aws.String("https://client-$urlSdkId:L.dev"),
                        expectURL:     aws.String("https://client-$urlSdkId:L.dev"),
                    },
                }
            """,
            this.commonCodegenArgs
        );
    }

    private GoWriter.Writable generateTests() {
        return goTemplate(
            """
                for name, c := range cases {
                    t.Run(name, func(t $testing:P) {
                        $clearEnv:T()

                        awsConfig := $awsConfig:T{}
                        ignore := c.envIgnore || c.configIgnore
            
                        if c.configGlobal != "" && !ignore {
                            awsConfig.BaseEndpoint = $awsString:T(c.configGlobal)
                        }
            
                        if c.envGlobal != "" {
                            t.Setenv("AWS_ENDPOINT_URL", c.envGlobal)
                            if !ignore {
                                awsConfig.BaseEndpoint = $awsString:T(c.envGlobal)
                            }
                        }
            
                        if c.envService != "" {
                            t.Setenv("AWS_ENDPOINT_URL_$envSdkId:L", c.envService)
                        }
            
                        awsConfig.ConfigSources = []interface{}{
                            mockConfigSource{
                                global: c.envGlobal,
                                service: c.envService,
                                ignore: c.envIgnore,
                            },
                            mockConfigSource{
                                global: c.configGlobal,
                                service: c.configService,
                                ignore: c.configIgnore,
                            },
                        }
            
                        client := NewFromConfig(awsConfig, func (o *Options) {
                            if c.clientEndpoint != nil {
                                o.BaseEndpoint = c.clientEndpoint
                            }
                        })

                        if e, a := c.expectURL, client.options.BaseEndpoint; !$deepEqual:T(e, a) {
                            t.Errorf("expect endpoint %v , got %v", e, a)
                        }
                    })
                }
            """,
            this.commonCodegenArgs,
            MapUtils.of(
                "clearEnv", SymbolUtils.createValueSymbolBuilder("Clearenv", SmithyGoDependency.OS).build(),
                "awsConfig", SymbolUtils.createValueSymbolBuilder("Config", AwsGoDependency.AWS_CORE).build(),
                "deepEqual", SymbolUtils.createValueSymbolBuilder("DeepEqual", SmithyGoDependency.REFLECT).build()
            )
        );
    }
}
