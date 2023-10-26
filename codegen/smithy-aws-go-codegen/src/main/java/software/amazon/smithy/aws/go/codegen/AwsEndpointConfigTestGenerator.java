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
            "awsString", SymbolUtils.createValueSymbolBuilder("String", AwsGoDependency.AWS_CORE).build()
        );

        writerFactory.accept("endpoints_config_test.go", settings.getModuleName(), writer -> {
            writer.write("$W", generate());
        });

    }


    private GoWriter.Writable generate() {
        return (GoWriter w) -> {
            w.write(
                """
                func TestConfiguredEndpoints(t $P) {
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
                    Env                  map[string]string
                    SharedConfigFile     string
                    ClientEndpoint       *string
                    ExpectURL            *string
                }{
                    "env ignore": {
                        Env: map[string]string{
                            "AWS_ENDPOINT_URL":         "https://env-global.dev",
                            "AWS_ENDPOINT_URL_$envSdkId:L":  "https://env-$urlSdkId:L.dev",
                            "AWS_IGNORE_CONFIGURED_ENDPOINT_URLS": "true",
                        },
                        SharedConfigFile: `[profile dev]
endpoint_url = http://config-global.dev
services = testing-$urlSdkId:L

[services testing-$urlSdkId:L]
$configSdkId:L =
    endpoint_url = http://config-$urlSdkId:L.dev
`,
                        ExpectURL:            nil,
                    },
                    "env global": {
                        Env: map[string]string{
                            "AWS_ENDPOINT_URL": "https://env-global.dev",
                        },
                        SharedConfigFile: `[profile dev]
endpoint_url = http://config-global.dev
`,
                        ExpectURL:            $awsString:T("https://env-global.dev"),
                    },
                    "env service": {
                        Env: map[string]string{
                            "AWS_ENDPOINT_URL":                    "https://env-global.dev",
                            "AWS_ENDPOINT_URL_$envSdkId:L":  "https://env-$urlSdkId:L.dev",
                        },
                        SharedConfigFile: `[profile dev]
endpoint_url = http://config-global.dev
services = testing-$urlSdkId:L

[services testing-$urlSdkId:L]
$configSdkId:L =
    endpoint_url = http://config-$urlSdkId:L.dev
`,
                        ExpectURL:            $awsString:T("https://env-$urlSdkId:L.dev"),
                    },
                    "config ignore": {
                        Env: map[string]string{
                            "AWS_ENDPOINT_URL":                    "https://env-global.dev",
                            "AWS_ENDPOINT_URL_$envSdkId:L":  "https://env-$urlSdkId:L.dev",
                        },
                        SharedConfigFile: `[profile dev]
endpoint_url = http://config-global.dev
services = testing-$urlSdkId:L
ignore_configured_endpoint_urls = true

[services testing-$urlSdkId:L]
$configSdkId:L =
    endpoint_url = http://config-$urlSdkId:L.dev
`,
                        ExpectURL:            nil,
                    },
                    "config global": {
                        SharedConfigFile: `[profile dev]
endpoint_url = http://config-global.dev
`,
                        ExpectURL:            $awsString:T("http://config-global.dev"),
                    },
                    "config service": {
                        Env: map[string]string{
                            "AWS_ENDPOINT_URL": "https://env-global.dev",
                        },
                        SharedConfigFile: `[profile dev]
endpoint_url = http://config-global.dev
services = testing-$urlSdkId:L

[services testing-$urlSdkId:L]
$configSdkId:L = 
    endpoint_url = http://config-$urlSdkId:L.dev
`,
                        ExpectURL:            $awsString:T("http://config-$urlSdkId:L.dev"),
                    },
                    "client": {
                        Env: map[string]string{
                            "AWS_ENDPOINT_URL":                    "https://env-global.dev",
                            "AWS_ENDPOINT_URL_$envSdkId:L":  "https://env-$urlSdkId:L.dev",
                            "AWS_IGNORE_CONFIGURED_ENDPOINT_URLS": "true",
                        },
                        SharedConfigFile: `[profile dev]
endpoint_url = http://config-global.dev
services = testing-$urlSdkId:L

[services testing-$urlSdkId:L]
$configSdkId:L = 
    endpoint_url = http://config-$urlSdkId:L.dev
`,
                        ClientEndpoint:       $awsString:T("https://client-$urlSdkId:L.dev"),
                        ExpectURL:            $awsString:T("https://client-$urlSdkId:L.dev"),
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
                        for k, v := range c.Env {
                            t.Setenv(k, v)
                        }

                        tmpDir := t.TempDir()
                        $writeFile:T($joinFile:T(tmpDir, "test_shared_config"), []byte(c.SharedConfigFile), $fileMode:T(int(0777)))

                        awsConfig, err := $loadDefaultConfig:T(
                            $contextTodo:T(),
                            $withSharedConfig:T([]string{$joinFile:T(tmpDir, "test_shared_config")}),
                            $withSharedConfigProfile:T("dev"),
                        )
                        if err != nil {
                            t.Fatalf("error loading default config: %v", err)
                        }

                        client := NewFromConfig(awsConfig, func (o *Options) {
                            if c.ClientEndpoint != nil {
                                o.BaseEndpoint = c.ClientEndpoint
                            }
                        })

                        if e, a := c.ExpectURL, client.options.BaseEndpoint; !$deepEqual:T(e, a) {
                            t.Errorf("expect endpoint %v , got %v", e, a)
                        }
                    })
                }
            """,
            this.commonCodegenArgs,
            MapUtils.of(
                "clearEnv", SymbolUtils.createValueSymbolBuilder("Clearenv", SmithyGoDependency.OS).build(),
                "writeFile", SymbolUtils.createValueSymbolBuilder("WriteFile", SmithyGoDependency.OS).build(),
                "joinFile", SymbolUtils.createValueSymbolBuilder("Join", SmithyGoDependency.PATH_FILEPATH).build(),
                "fileMode", SymbolUtils.createValueSymbolBuilder("FileMode", SmithyGoDependency.OS).build(),
                "loadDefaultConfig", SymbolUtils.createValueSymbolBuilder("LoadDefaultConfig", AwsGoDependency.CONFIG).build(),
                "contextTodo", SymbolUtils.createValueSymbolBuilder("TODO", SmithyGoDependency.CONTEXT).build(),
                "withSharedConfig", SymbolUtils.createValueSymbolBuilder("WithSharedConfigFiles", AwsGoDependency.CONFIG).build(),
                "withSharedConfigProfile", SymbolUtils.createValueSymbolBuilder("WithSharedConfigProfile", AwsGoDependency.CONFIG).build(),
                "deepEqual", SymbolUtils.createValueSymbolBuilder("DeepEqual", SmithyGoDependency.REFLECT).build()
            )
        );
    }
}
