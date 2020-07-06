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
import java.util.function.Consumer;
import software.amazon.smithy.go.codegen.CodegenUtils;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.HttpProtocolTestGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUnitTestGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUnitTestRequestGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUnitTestResponseErrorGenerator;
import software.amazon.smithy.go.codegen.integration.HttpProtocolUnitTestResponseGenerator;
import software.amazon.smithy.go.codegen.integration.IdempotencyTokenMiddlewareGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.traits.EnumTrait;
import software.amazon.smithy.utils.SetUtils;

/**
 * Utility methods for generating AWS protocols.
 */
final class AwsProtocolUtils {
    private AwsProtocolUtils() {}

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
                            writer.openBlock("aws.EndpointResolverFunc("
                                            + "func(service, region string) (e aws.Endpoint, err error) {",
                                    "}),", () -> {
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
                                    writer.write("s.Build.Clear()");
                                    writer.write("s.Finalize.Clear()");
                                    writer.write("return nil");
                                });
                            });
                        })
                        .build()
        ));

        // TODO can this check be replaced with a lookup into the runtime plugins?
        if (IdempotencyTokenMiddlewareGenerator.hasOperationsWithIdempotencyToken(context.getModel(), context.getService())) {
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

    /**
     * Safely accesses a given structure member.
     *
     * @param context The generation context.
     * @param member The member being accessed.
     * @param container The name that the structure is assigned to.
     * @param consumer A string consumer that is given the snippet to access the member value.
     */
    public static void writeSafeMemberAccessor(
            GenerationContext context,
            MemberShape member,
            String container,
            Consumer<String> consumer
    ) {
        Model model = context.getModel();
        Shape target = model.expectShape(member.getTarget());
        String memberName = context.getSymbolProvider().toMemberName(member);
        String operand = container + "." + memberName;

        boolean enumShape = target.hasTrait(EnumTrait.class);

        if (!enumShape && !CodegenUtils.isNilAssignableToShape(model, member)) {
            consumer.accept(operand);
            return;
        }

        String conditionCheck;
        if (enumShape) {
            conditionCheck = "len(" + operand + ") > 0";
        } else {
            conditionCheck = operand + " != nil";
        }

        context.getWriter().openBlock("if $L {", "}", conditionCheck, () -> {
            consumer.accept(operand);
        });
    }
}
