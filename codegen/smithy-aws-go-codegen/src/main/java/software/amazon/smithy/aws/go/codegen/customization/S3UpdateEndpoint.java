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
 *
 *
 */

package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import java.util.stream.Collectors;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ConfigField;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.utils.ListUtils;

/**
 * S3UpdateEndpoint integration serves to apply customizations for S3 service,
 * and modifies the resolved endpoint based on S3 client config or input shape values.
 */
public class S3UpdateEndpoint implements GoIntegration {
    private static final String USE_PATH_STYLE_OPTION = "UsePathStyle";
    private static final String UPDATE_ENDPOINT_ADDER = "addUpdateEndpointMiddleware";
    private static final String UPDATE_ENDPOINT_INTERNAL_ADDER = "UpdateEndpoint";
    private static final String GET_BUCKET_FROM_INPUT = "getBucketFromInput";

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first.
     *
     * @return Returns the sort order, defaults to -40.
     */
    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceShape service = settings.getService(model);
        if (!isS3Service(model, service)) {
            return;
        }

        goDelegator.useShapeWriter(service, this::writeMiddlewareHelper);

        goDelegator.useShapeWriter(service, writer -> {
            writeInputGetter(writer, model, symbolProvider, service);
        });
    }

    private void writeInputGetter(GoWriter writer, Model model, SymbolProvider symbolProvider, ServiceShape service) {
        writer.writeDocs("getBucketFromInput returns a boolean indicating if the input has a modeled bucket name, " +
                " and a pointer to string denoting a provided bucket member value");
        writer.openBlock("func getBucketFromInput(input interface{}) (*string, bool) {","}", ()-> {
            writer.openBlock("switch i:= input.(type) {", "}", () -> {
                service.getAllOperations().forEach((operationId)-> {
                    OperationShape operation = model.expectShape(operationId, OperationShape.class);
                    StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);

                    List<MemberShape> targetBucketShape = input.getAllMembers().values().stream()
                            .filter(m -> m.getTarget().getName().equals("BucketName"))
                            .collect(Collectors.toList());
                    // if model has multiple top level shapes targeting `BucketName`, we throw a codegen exception
                    if (targetBucketShape.size()>1) {
                        throw new CodegenException("BucketName shape should be targeted by only one input member, found " +
                                targetBucketShape.size() +" for Input shape: "+ input.getId());
                    }

                    if (!targetBucketShape.isEmpty()) {
                        writer.write("case $P: return i.$L, true", symbolProvider.toSymbol(input), targetBucketShape.get(0).getMemberName());
                    }
                });
                writer.write("default: return nil, false");
            });
        });
    }

    private void writeMiddlewareHelper(GoWriter writer) {
        writer.openBlock("func $L(stack *middleware.Stack, options Options) {", "}", UPDATE_ENDPOINT_ADDER, () -> {
            writer.write("$T(stack, $T{UsePathStyle: options.$L, GetBucketFromInput: $L})",
                    SymbolUtils.createValueSymbolBuilder(UPDATE_ENDPOINT_INTERNAL_ADDER,
                            AwsCustomGoDependency.S3_CUSTOMIZATION).build(),
                    SymbolUtils.createValueSymbolBuilder(UPDATE_ENDPOINT_INTERNAL_ADDER + "Options",
                            AwsCustomGoDependency.S3_CUSTOMIZATION).build(),
                    USE_PATH_STYLE_OPTION,
                    GET_BUCKET_FROM_INPUT
            );
        });
        writer.insertTrailingNewline();
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                // Add S3 config to use path style host addressing.
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3UpdateEndpoint::isS3Service)
                        .configFields(ListUtils.of(
                                ConfigField.builder()
                                        .name(USE_PATH_STYLE_OPTION)
                                        .type(SymbolUtils.createValueSymbolBuilder("bool")
                                                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                                .build())
                                        .documentation("Allows you to enable the client to use path-style addressing, "
                                                + "i.e., `https://s3.amazonaws.com/BUCKET/KEY`. By default, the S3 client "
                                                + "will use virtual hosted bucket addressing when possible"
                                                + "(`https://BUCKET.s3.amazonaws.com/KEY`).")
                                        .build()
                                ))
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(UPDATE_ENDPOINT_ADDER).build())
                                .useClientOptions()
                                .build())
                        .build()
        );
    }

    private static boolean isS3Service(Model model, ServiceShape service) {
        return service.expectTrait(ServiceTrait.class).getSdkId().equalsIgnoreCase("S3");
    }
}
