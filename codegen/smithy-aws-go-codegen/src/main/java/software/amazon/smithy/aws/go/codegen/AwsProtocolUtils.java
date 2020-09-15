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

import java.util.Set;
import java.util.TreeSet;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.HttpProtocolTestGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUnitTestGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUnitTestRequestGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUnitTestResponseErrorGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUnitTestResponseGenerator;
import software.amazon.smithy.go.codegen.integration.IdempotencyTokenMiddlewareGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;
import software.amazon.smithy.utils.SetUtils;

/**
 * Utility methods for generating AWS protocols.
 */
final class AwsProtocolUtils {
    private AwsProtocolUtils() {
    }

    /**
     * Generates HTTP protocol tests with all required AWS-specific configuration set.
     *
     * @param context The generation context.
     */
    static void generateHttpProtocolTests(GenerationContext context) {
        Set<HttpProtocolUnitTestGenerator.ConfigValue> configValues = new TreeSet<>(SetUtils.of(
                HttpProtocolUnitTestGenerator.ConfigValue.builder()
                        .name(AddAwsConfigFields.REGION_CONFIG_NAME)
                        .value(writer -> writer.write("$S,", "us-west-2"))
                        .build(),
                HttpProtocolUnitTestGenerator.ConfigValue.builder()
                        .name(AddAwsConfigFields.HTTP_CLIENT_CONFIG_NAME)
                        .value(writer -> {
                            writer.addUseImports(AwsGoDependency.AWS_CORE);
                            writer.write("aws.NewBuildableHTTPClient(),");
                        })
                        .build(),
                HttpProtocolUnitTestGenerator.ConfigValue.builder()
                        .name(AddAwsConfigFields.ENDPOINT_RESOLVER_CONFIG_NAME)
                        .value(writer -> {
                            writer.addUseImports(AwsGoDependency.AWS_CORE);
                            writer.openBlock("$L(func(region string, options $L) (e aws.Endpoint, err error) {", "}),",
                                    EndpointGenerator.RESOLVER_FUNC_NAME, EndpointGenerator.RESOLVER_OPTIONS, () -> {
                                        writer.write("e.URL = server.URL");
                                        writer.write("e.SigningRegion = \"us-west-2\"");
                                        writer.write("return e, err");
                                    });
                        })
                        .build(),
                HttpProtocolUnitTestGenerator.ConfigValue.builder()
                        .name("APIOptions")
                        .value(writer -> {
                            writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);
                            writer.openBlock("[]APIOptionFunc{", "},", () -> {
                                writer.openBlock("func(s *middleware.Stack) error {", "},", () -> {
                                    writer.write("s.Finalize.Clear()");
                                    writer.write("return nil");
                                });
                            });
                        })
                        .build()
        ));

        // TODO can this check be replaced with a lookup into the runtime plugins?
        if (IdempotencyTokenMiddlewareGenerator.hasOperationsWithIdempotencyToken(context.getModel(),
                context.getService())) {
            configValues.add(
                    HttpProtocolUnitTestGenerator.ConfigValue.builder()
                            .name(IdempotencyTokenMiddlewareGenerator.IDEMPOTENCY_CONFIG_NAME)
                            .value(writer -> {
                                writer.addUseImports(SmithyGoDependency.SMITHY_RAND);
                                writer.addUseImports(SmithyGoDependency.SMITHY_TESTING);
                                writer.write("smithyrand.NewUUIDIdempotencyToken(&smithytesting.ByteLoop{}),");
                            })
                            .build()
            );
        }

        new HttpProtocolTestGenerator(context,
                (HttpProtocolUnitTestRequestGenerator.Builder) new HttpProtocolUnitTestRequestGenerator
                        .Builder()
                        .addClientConfigValues(configValues),
                (HttpProtocolUnitTestResponseGenerator.Builder) new HttpProtocolUnitTestResponseGenerator
                        .Builder()
                        .addClientConfigValues(configValues),
                (HttpProtocolUnitTestResponseErrorGenerator.Builder) new HttpProtocolUnitTestResponseErrorGenerator
                        .Builder()
                        .addClientConfigValues(configValues)
        ).generateProtocolTests();
    }

    public static void writeJsonErrorMessageCodeDeserializer(GenerationContext context) {
        GoWriter writer = context.getWriter();
        // The error code could be in the headers, even though for this protocol it should be in the body.
        writer.write("code := response.Header.Get(\"X-Amzn-ErrorType\")");
        writer.write("if len(code) != 0 { errorCode = restjson.SanitizeErrorCode(code) }");
        writer.write("");

        initializeJsonDecoder(writer, "errorBody");
        writer.addUseImports(AwsGoDependency.AWS_REST_JSON_PROTOCOL);
        // This will check various body locations for the error code and error message
        writer.write("code, message, err := restjson.GetErrorInfo(decoder)");
        handleDecodeError(writer);

        writer.addUseImports(SmithyGoDependency.IO);
        // Reset the body in case it needs to be used for anything else.
        writer.write("errorBody.Seek(0, io.SeekStart)");

        // Only set the values if something was found so that we keep the default values.
        writer.write("if len(code) != 0 { errorCode = restjson.SanitizeErrorCode(code) }");
        writer.write("if len(message) != 0 { errorMessage = message }");
        writer.write("");
    }

    public static void initializeJsonDecoder(GoWriter writer, String bodyLocation) {
        // Use a ring buffer and tee reader to help in pinpointing any deserialization errors.
        writer.addUseImports(SmithyGoDependency.SMITHY_IO);
        writer.write("var buff [1024]byte");
        writer.write("ringBuffer := smithyio.NewRingBuffer(buff[:])");
        writer.write("");

        writer.addUseImports(SmithyGoDependency.IO);
        writer.addUseImports(SmithyGoDependency.JSON);
        writer.write("body := io.TeeReader($L, ringBuffer)", bodyLocation);
        writer.write("decoder := json.NewDecoder(body)");
        writer.write("decoder.UseNumber()");
        writer.write("");
    }

    public static void handleDecodeError(GoWriter writer, String returnExtras) {
        writer.openBlock("if err != nil {", "}", () -> {
            writer.addUseImports(SmithyGoDependency.BYTES);
            writer.addUseImports(SmithyGoDependency.SMITHY);
            writer.addUseImports(SmithyGoDependency.IO);
            writer.write("var snapshot bytes.Buffer");
            writer.write("io.Copy(&snapshot, ringBuffer)");
            writer.openBlock("return $L&smithy.DeserializationError {", "}", returnExtras, () -> {
                writer.write("Err: fmt.Errorf(\"failed to decode response body, %w\", err),");
                writer.write("Snapshot: snapshot.Bytes(),");
            });
        }).write("");
    }

    public static void handleDecodeError(GoWriter writer) {
        handleDecodeError(writer, "");
    }
}
