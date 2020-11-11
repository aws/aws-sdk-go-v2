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

package software.amazon.smithy.aws.go.codegen;

import java.util.Map;
import java.util.Set;
import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.aws.traits.protocols.AwsQueryTrait;
import software.amazon.smithy.aws.traits.protocols.Ec2QueryTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.OperationGenerator;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.utils.SetUtils;

public class AwsHttpPresignURLClientGenerator implements GoIntegration {
    // constants
    private static final String CONVERT_TO_PRESIGN_MIDDLEWARE_NAME = "convertToPresignMiddleware";

    private static final String PRESIGN_CLIENT = "PresignClient";
    private static final Symbol presignClientSymbol = buildSymbol(PRESIGN_CLIENT, true);

    private static final String NEW_CLIENT = "NewPresignClient";
    private static final String NEW_CLIENT_FROM_SERVICE = "NewPresignClientWrapper";
    private static final String NEW_CLIENT_FROM_CONFIG = "NewPresignClientFromConfig";

    private static final String PRESIGN_OPTIONS = "PresignOptions";
    private static final Symbol presignOptionsSymbol = buildSymbol(PRESIGN_OPTIONS, true);

    private static final String PRESIGN_OPTIONS_FROM_CLIENT_OPTIONS = "WithPresignClientFromClientOptions";
    private static final String PRESIGN_OPTIONS_FROM_EXPIRES = "WithPresignExpires";

    private static final String PRESIGN_SIGNER = "Presigner";
    private static final Symbol presignerInterfaceSymbol = SymbolUtils.createPointableSymbolBuilder(
            "HTTPPresigner", AwsGoDependency.AWS_SIGNER_V4
    ).build();
    private static final Symbol v4NewPresignerSymbol = SymbolUtils.createPointableSymbolBuilder(
            "NewSigner", AwsGoDependency.AWS_SIGNER_V4
    ).build();

    // map of service to list of operations for which presignedURL client and operation should
    // be generated.
    private static final Map<ShapeId, Set<ShapeId>> PRESIGNER_MAP = MapUtils.of(
            ShapeId.from("com.amazonaws.s3#AmazonS3"), SetUtils.of(
                    ShapeId.from("com.amazonaws.s3#PutObject")),

            ShapeId.from("com.amazonaws.rds#AmazonRDSv19"), SetUtils.of(
                    ShapeId.from("com.amazonaws.rds#CopyDBSnapshot"),
                    ShapeId.from("com.amazonaws.rds#CreateDBInstanceReadReplica"),
                    ShapeId.from("com.amazonaws.rds#CopyDBClusterSnapshot"),
                    ShapeId.from("com.amazonaws.rds#CreateDBCluster")),

            ShapeId.from("com.amazonaws.ec2#AmazonEC2"), SetUtils.of(
                    ShapeId.from("com.amazonaws.ec2#CopySnapshot"))

            // TODO other services
    );

    // build pointable symbols
    private static Symbol buildSymbol(String name, boolean exported) {
        if (!exported) {
            name = Character.toLowerCase(name.charAt(0)) + name.substring(1);
        }
        return SymbolUtils.createPointableSymbolBuilder(name).
                build();
    }

    /**
     * generates code to iterate thru func optionals and assign value into the dest variable
     *
     * @param writer   GoWriter to write the code to
     * @param src      variable name that denotes functional options
     * @param dest     variable in which result of processed functional options are stored
     * @param destType value type used by functional options
     */
    private static final void processFunctionalOptions(
            GoWriter writer,
            String src,
            String dest,
            Symbol destType
    ) {
        writer.write("var $L $T", dest, destType);
        writer.openBlock("for _, fn := range $L {", "}", src, () -> {
            writer.write("fn(&$L)", dest);
        }).insertTrailingNewline();
    }

    /**
     * variables needed in scope:
     * * client
     * * presignOptions
     * <p>
     * generates code to assign client, presigner and return a new presign client
     *
     * @param writer the writer to write to
     */
    private final void returnPresignClientConstructor(
            GoWriter writer,
            Model model,
            ServiceShape serviceShape
    ) {
        writer.write("var presigner $T", presignerInterfaceSymbol);
        writer.openBlock("if presignOptions.Presigner != nil {", "} else {", () -> {
            writer.write("presigner = presignOptions.Presigner");
        }).write("presigner = $T() }", v4NewPresignerSymbol).insertTrailingNewline();

        writer.openBlock("return &$L{", "}", presignClientSymbol, () -> {
            writer.write("client: client, presigner: presigner,");
            //  if s3 assign expires value on client
            if (isS3ServiceShape(model, serviceShape)) {
                writer.write("expires: presignOptions.Expires,");
            }
        });

    }

    @Override
    public byte getOrder() {
        // The associated customization ordering is relative to operation deserializers
        // and thus the integration should be added at the end.
        return 127;
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        ServiceShape serviceShape = settings.getService(model);

        if (!PRESIGNER_MAP.containsKey(serviceShape.getId())) {
            return;
        }

        Set<ShapeId> validOperations = PRESIGNER_MAP.get(serviceShape.getId());
        if (validOperations.isEmpty()) {
            return;
        }

        // delegator for service shape
        goDelegator.useShapeWriter(serviceShape, (writer) -> {
            // generate presign options and helpers per service
            writePresignOptionType(writer, model, symbolProvider, serviceShape);

            // generate Presign client per service
            writePresignClientType(writer, model, symbolProvider, serviceShape);

            // generate client helpers such as copyAPIClient, GetAPIClientOptions()
            writePresignClientHelpers(writer, model, symbolProvider, serviceShape);

            // generate convertToPresignMiddleware per service
            writeConvertToPresignMiddleware(writer, model, symbolProvider, serviceShape);
        });


        for (ShapeId operationId : serviceShape.getAllOperations()) {
            OperationShape operationShape = model.expectShape(operationId, OperationShape.class);
            if (!validOperations.contains(operationShape.getId())) {
                continue;
            }

            goDelegator.useShapeWriter(operationShape, (writer) -> {
                // generate presign operation function for a client operation.
                writePresignOperationFunction(writer, model, symbolProvider, serviceShape, operationShape);
            });
        }
    }

    private void writePresignOperationFunction(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape serviceShape,
            OperationShape operationShape
    ) {
        Symbol operationSymbol = symbolProvider.toSymbol(operationShape);

        Shape operationInputShape = model.expectShape(operationShape.getInput().get());
        Symbol operationInputSymbol = symbolProvider.toSymbol(operationInputShape);

        writer.openBlock(
                "func (c *$T) Presign$T(ctx context.Context, params $P, optFns ...func($P)) "
                        + "(req *v4.PresignedHTTPRequest, err error) {",
                "}",
                presignClientSymbol, operationSymbol, operationInputSymbol, presignOptionsSymbol,
                () -> {
                    Symbol nopClient = SymbolUtils.createPointableSymbolBuilder("NopClient",
                            SmithyGoDependency.SMITHY_HTTP_TRANSPORT)
                            .build();

                    writer.write("if params == nil { params = &$T{} }", operationInputSymbol).insertTrailingNewline();

                    // process presignerOptions
                    processFunctionalOptions(writer, "optFns", "presignOptions", presignOptionsSymbol);

                    // check if presigner was set for presignerOptions
                    writer.openBlock("if presignOptions.Presigner != nil {", "}", () -> {
                        writer.write(
                                "c = NewPresignClientWrapper(c.client, func (o $P) { o.Presigner = presignOptions.Presigner })",
                                presignOptionsSymbol);
                    });

                    writer.write("clientOptFns := presignOptions.ClientOptions");

                    writer.openBlock("clientOptFns = append(clientOptFns, func(o *Options) {", "})", () -> {
                        writer.write("o.HTTPClient = &$T{}", nopClient);
                    });
                    writer.insertTrailingNewline();

                    Symbol withIsPresigning = SymbolUtils.createValueSymbolBuilder("WithIsPresigning",
                            AwsCustomGoDependency.PRESIGNEDURL_CUSTOMIZATION).build();

                    writer.write("ctx = $T(ctx)", withIsPresigning);
                    writer.openBlock("result, _, err := c.client.invokeOperation(ctx, $S, params, clientOptFns,", ")",
                            operationSymbol.getName(), () -> {
                                writer.write("$L,", OperationGenerator
                                        .getAddOperationMiddlewareFuncName(operationSymbol));
                                writer.write("c.$L,", CONVERT_TO_PRESIGN_MIDDLEWARE_NAME);
                            });
                    writer.write("if err != nil { return req, err }");
                    writer.insertTrailingNewline();

                    writer.write("out := result.(*v4.PresignedHTTPRequest)");
                    writer.write("return out, nil");
                });
    }

    /**
     * generates a helper to mutate request middleware stack in favor of generating a presign URL request
     *
     * @param writer         the writer to write to
     * @param model          the service model
     * @param symbolProvider the symbol provider
     * @param serviceShape   the service for which helper is generated
     */
    private void writeConvertToPresignMiddleware(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape serviceShape
    ) {
        Symbol smithyStack = SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE)
                .build();

        writer.openBlock("func (c *$T) $L(stack $P, options Options) (err error) {", "}",
                presignClientSymbol,
                CONVERT_TO_PRESIGN_MIDDLEWARE_NAME,
                smithyStack,
                () -> {
                    Symbol smithyAfter = SymbolUtils.createValueSymbolBuilder("After",
                            SmithyGoDependency.SMITHY_MIDDLEWARE)
                            .build();

                    // Middleware to remove
                    Symbol requestInvocationID = SymbolUtils.createPointableSymbolBuilder(
                            "ClientRequestID",
                            AwsGoDependency.AWS_MIDDLEWARE)
                            .build();

                    Symbol presignMiddleware = SymbolUtils.createValueSymbolBuilder("NewPresignHTTPRequestMiddleware",
                            AwsGoDependency.AWS_SIGNER_V4)
                            .build();


                    // Middleware to add
                    writer.write("stack.Finalize.Clear()");
                    writer.write("stack.Deserialize.Clear()");
                    writer.write("stack.Build.Remove(($P)(nil).ID())", requestInvocationID);

                    writer.write("err = stack.Finalize.Add($T(options.Credentials, c.presigner), $T)",
                            presignMiddleware, smithyAfter);
                    writer.write("if err != nil { return err }");

                    // if protocol used is ec2query or query
                    if (serviceShape.hasTrait(AwsQueryTrait.ID) || serviceShape.hasTrait(Ec2QueryTrait.ID)) {
                        // presigned url should convert to Get request
                        Symbol queryAsGetMiddleware = SymbolUtils.createValueSymbolBuilder("AddAsGetRequestMiddleware",
                                AwsGoDependency.AWS_QUERY_PROTOCOL).build();

                        writer.writeDocs("convert request to a GET request");
                        writer.write("err = $T(stack)", queryAsGetMiddleware);
                        writer.write("if err != nil { return err }");
                    }

                    // s3 service needs expires
                    if (isS3ServiceShape(model, serviceShape)) {
                        Symbol expiresAsHeaderMiddleware = SymbolUtils.createValueSymbolBuilder(
                                "AddExpiresOnPresignedURL",
                                AwsCustomGoDependency.S3_CUSTOMIZATION).build();
                        writer.openBlock("if c.expires != 0 {", "}", () -> {
                            writer.writeDocs("add middleware to set expiration for s3 presigned url");
                            writer.write("err = stack.Build.Add(&$T{ Expires: c.expires, }, middleware.After)",
                                    expiresAsHeaderMiddleware);
                            writer.write("if err != nil { return err }");
                        });
                    }

                    writer.write("return nil");
                }).insertTrailingNewline();
    }


    /**
     * Writes the Presign client's type and methods.
     *
     * @param writer writer to write to
     */
    private void writePresignClientType(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape serviceShape
    ) {
        writer.addUseImports(SmithyGoDependency.CONTEXT);
        writer.addUseImports(AwsGoDependency.AWS_SIGNER_V4);


        Symbol presignerInterfaceSymbol = SymbolUtils.createPointableSymbolBuilder(
                "HTTPPresigner", AwsGoDependency.AWS_SIGNER_V4
        ).build();

        writer.writeDocs(String.format("%s represents the presign url client", PRESIGN_CLIENT));
        writer.openBlock("type $T struct {", "}", presignClientSymbol, () -> {
            writer.write("client *Client");
            writer.write("presigner v4.HTTPPresigner");

            if (isS3ServiceShape(model, serviceShape)) {
                writer.addUseImports(SmithyGoDependency.TIME);
                writer.write("expires time.Duration");
            }

        });

        // generate constructors

        // generate NewPresignClient
        writer.writeDocs(
                String.format("%s generates a presign client using provided Client options and presign options",
                        NEW_CLIENT)
        );
        writer.openBlock("func $L(options Options, optFns ...func($P)) $P {", "}",
                NEW_CLIENT, presignOptionsSymbol, presignClientSymbol, () -> {
                    processFunctionalOptions(writer, "optFns", "presignOptions", presignOptionsSymbol);
                    writer.insertTrailingNewline();

                    writer.write("client := New(options, presignOptions.ClientOptions...)").insertTrailingNewline();
                    writer.insertTrailingNewline();

                    returnPresignClientConstructor(writer, model, serviceShape);
                }).insertTrailingNewline();

        // generate NewPresignClientWrapper
        writer.writeDocs(
                String.format("%s generates a presign client using provided API Client and presign options",
                        NEW_CLIENT_FROM_SERVICE)
        );
        writer.openBlock("func $L(c *Client, optFns ...func($P)) $P {", "}",
                NEW_CLIENT_FROM_SERVICE, presignOptionsSymbol, presignClientSymbol, () -> {
                    processFunctionalOptions(writer, "optFns", "presignOptions", presignOptionsSymbol);
                    writer.insertTrailingNewline();

                    writer.write("client := copyAPIClient(c, presignOptions.ClientOptions...)");
                    writer.insertTrailingNewline();

                    returnPresignClientConstructor(writer, model, serviceShape);
                }).insertTrailingNewline();

        // generate NewPresignClientFromConfig
        writer.writeDocs(
                String.format("%s generates a presign client using provided AWS config and presign options",
                        NEW_CLIENT_FROM_CONFIG)
        );
        writer.openBlock("func $L(cfg aws.Config, optFns ...func($P)) $P {", "}",
                NEW_CLIENT_FROM_CONFIG, presignOptionsSymbol, presignClientSymbol, () -> {
                    processFunctionalOptions(writer, "optFns", "presignOptions", presignOptionsSymbol);
                    writer.insertTrailingNewline();

                    writer.write("client := NewFromConfig(cfg, presignOptions.ClientOptions...)");
                    writer.insertTrailingNewline();

                    returnPresignClientConstructor(writer, model, serviceShape);
                }).insertTrailingNewline();

        writer.insertTrailingNewline();
    }

    /**
     * Writes the Presign client's helper methods.
     *
     * @param writer writer to write to
     */
    private void writePresignClientHelpers(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape serviceShape
    ) {
        // generate copy API client
        final String COPY_API_CLIENT = "copyAPIClient";
        writer.openBlock("func $L(c *Client, optFns ...func(*Options)) *Client {", "}",
                COPY_API_CLIENT, () -> {
                    writer.write("return New(c.options, optFns...)");
                    writer.insertTrailingNewline();
                }).insertTrailingNewline();
        writer.insertTrailingNewline();
    }

    /**
     * Writes the Presign client's type and methods.
     *
     * @param writer writer to write to
     */
    public void writePresignOptionType(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape serviceShape
    ) {
        writer.addUseImports(SmithyGoDependency.CONTEXT);
        writer.addUseImports(AwsGoDependency.AWS_SIGNER_V4);

        Symbol presignOptionSymbol = buildSymbol(PRESIGN_OPTIONS, true);

        // generate presign options
        writer.writeDocs(String.format("%s represents the presign client options", PRESIGN_OPTIONS));
        writer.openBlock("type $T struct {", "}", presignOptionSymbol, () -> {
            writer.writeDocs(
                    "ClientOptions are list of functional options to mutate client options used by presign client");
            writer.write("ClientOptions []func(*Options)");
            writer.insertTrailingNewline();

            writer.writeDocs("Presigner is the  presigner used by the presign url client");
            writer.write("Presigner $T", presignerInterfaceSymbol);

            // s3 service has an additional Expires options
            if (isS3ServiceShape(model, serviceShape)) {
                writer.writeDocs("Expires sets the expiration duration for the generated presign url");
                writer.write("Expires time.Duration");
            }
        });

        // generate WithPresignClientFromClientOptions Helper
        Symbol presignOptionsFromClientOptionsInternal = buildSymbol(PRESIGN_OPTIONS_FROM_CLIENT_OPTIONS, false);
        writer.writeDocs(
                String.format("%s is a helper utility to retrieve a function that takes PresignOption as input",
                        PRESIGN_OPTIONS_FROM_CLIENT_OPTIONS)
        );
        writer.openBlock("func $L(optFns ...func(*Options)) func($P) {", "}",
                PRESIGN_OPTIONS_FROM_CLIENT_OPTIONS, presignOptionSymbol, () -> {
                    writer.write("return $L(optFns).options", presignOptionsFromClientOptionsInternal.getName());
                });

        writer.insertTrailingNewline();

        writer.write("type $L []func(*Options)", presignOptionsFromClientOptionsInternal.getName());
        writer.openBlock("func (w $L) options (o $P) {", "}",
                presignOptionsFromClientOptionsInternal.getName(), presignOptionSymbol, () -> {
                    writer.write("o.ClientOptions = append(o.ClientOptions, w...)");
                }).insertTrailingNewline();


        // s3 specific helpers
        if (isS3ServiceShape(model, serviceShape)) {
            // generate WithPresignExpires Helper
            Symbol presignOptionsFromExpiresInternal = buildSymbol(PRESIGN_OPTIONS_FROM_EXPIRES, false);
            writer.writeDocs(String.format(
                    "%s is a helper utility to append Expires value on presign options optional function",
                    PRESIGN_OPTIONS_FROM_EXPIRES));
            writer.openBlock("func $L(dur time.Duration) func($P) {", "}",
                    PRESIGN_OPTIONS_FROM_EXPIRES, presignOptionSymbol, () -> {
                        writer.write("return $L(dur).options", presignOptionsFromExpiresInternal.getName());
                    });

            writer.insertTrailingNewline();

            writer.write("type $L time.Duration", presignOptionsFromExpiresInternal.getName());
            writer.openBlock("func (w $L) options (o $P) {", "}",
                    presignOptionsFromExpiresInternal.getName(), presignOptionSymbol, () -> {
                        writer.write("o.Expires = time.Duration(w)");
                    }).insertTrailingNewline();
        }
    }

    private final boolean isS3ServiceShape(Model model, ServiceShape service) {
        String serviceId = service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("S3");
    }
}

// TODO: generate tests for presigned urls
