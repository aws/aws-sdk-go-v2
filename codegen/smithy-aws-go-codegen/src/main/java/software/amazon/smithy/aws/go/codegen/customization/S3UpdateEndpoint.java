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
import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.aws.go.codegen.AwsSignatureVersion4;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.Symbol;
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
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.utils.ListUtils;

public class S3UpdateEndpoint implements GoIntegration {
    private static final String USE_PATH_STYLE_OPTION = "UsePathStyle";
    private static final String USE_ACCELERATE_OPTION = "UseAccelerate";
    private static final String UPDATE_ENDPOINT_ADDER = "addUpdateEndpointMiddleware";
    private static final String UPDATE_ENDPOINT_INTERNAL_ADDER = "UpdateEndpoint";

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
        if (!isS3Service(model, settings.getService(model))) {
            return;
        }

        goDelegator.useShapeWriter(settings.getService(model), this::writeMiddlewareHelper);

        // Generate getter's for an operation input that can be used to satisfy interfaces
        for (ShapeId shapeId: settings.getService(model).getOperations()) {
            OperationShape operation = model.expectShape(shapeId, OperationShape.class);
            StructureShape input = model.expectShape(operation.getInput().get(), StructureShape.class);
            goDelegator.useShapeWriter(operation, writer -> {
                writeInputGetter(writer, symbolProvider, input);
            });
        }
    }

    private void writeInputGetter(GoWriter writer, SymbolProvider symbolProvider, StructureShape input) {
        // generate bucketGetter if input has a member named Bucket
        if (input.getMember("Bucket").isPresent()) {
            // generateBucketGetter
            // TODO: Look for alternatives so that these Getters are NOT exported?
            Symbol inputSymbol = symbolProvider.toSymbol(input);
            writer.writeDocs("GetBucket retrieves the Bucket member value if provided");
            writer.openBlock("func (s $P) GetBucket() (v string) {", "}", inputSymbol, () -> {
                writer.write("if s.Bucket == nil { return v }");
                writer.write("return *s.Bucket");
            });
            writer.insertTrailingNewline();
        }
    }

    private void writeMiddlewareHelper(GoWriter writer) {
        writer.openBlock("func $L(stack *middleware.Stack, options Options) {", "}", UPDATE_ENDPOINT_ADDER, () -> {
            writer.write("$T(stack, $T{UsePathStyle: options.$L, UseAccelerate: options.$L, Region: options.Region})",
                    SymbolUtils.createValueSymbolBuilder(UPDATE_ENDPOINT_INTERNAL_ADDER,
                            AwsCustomGoDependency.S3_CUSTOMIZATION).build(),
                    SymbolUtils.createValueSymbolBuilder(UPDATE_ENDPOINT_INTERNAL_ADDER + "Options",
                            AwsCustomGoDependency.S3_CUSTOMIZATION).build(),
                    USE_PATH_STYLE_OPTION,
                    USE_ACCELERATE_OPTION
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
                                                + "i.e., `http://s3.amazonaws.com/BUCKET/KEY`. By default, the S3 client"
                                                + "will use virtual hosted bucket addressing when possible"
                                                + "(`http://BUCKET.s3.amazonaws.com/KEY`).")
                                        .build(),
                                ConfigField.builder()
                                        .name(USE_ACCELERATE_OPTION)
                                        .type(SymbolUtils.createValueSymbolBuilder("bool")
                                                .putProperty(SymbolUtils.GO_UNIVERSE_TYPE, true)
                                                .build())
                                        .documentation("Set this to `true` to enable S3 Accelerate feature. For all operations "
                                                + "compatible with S3 Accelerate will use the accelerate endpoint for "
                                                + "requests. Requests not compatible will fall back to normal S3 requests. "
                                                + ""
                                                + "The bucket must be enable for accelerate to be used with S3 client with "
                                                + "accelerate enabled. If the bucket is not enabled for accelerate an error "
                                                + "will be returned. The bucket name must be DNS compatible to also work "
                                                + "with accelerate."
                                        )
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
