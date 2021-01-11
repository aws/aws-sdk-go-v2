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
import java.util.TreeMap;
import java.util.TreeSet;
import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.aws.go.codegen.customization.PresignURLAutoFill;
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
import software.amazon.smithy.model.traits.StreamingTrait;
import software.amazon.smithy.utils.MapUtils;
import software.amazon.smithy.utils.SetUtils;

/**
 * AwsHttpPresignURLClientGenerator class is a runtime plugin integration class
 * that generates code for presign URL clients and associated presign operations.
 * <p>
 * This class pulls in a static list from PresignURLAutofill customization which
 * rely on the generated presigned url client and operation. This is done to
 * deduplicate the listing but make this class dependent on presence of PresignURLAutofill
 * class as a composition.
 */
public class AwsHttpPresignURLClientGenerator implements GoIntegration {
    // constants
    private static final String CONVERT_TO_PRESIGN_MIDDLEWARE_NAME = "convertToPresignMiddleware";
    private static final String CONVERT_TO_PRESIGN_TYPE_NAME = "presignConverter";
    private static final String NOP_HTTP_CLIENT_OPTION_FUNC_NAME = "withNopHTTPClientAPIOption";

    private static final String PRESIGN_CLIENT = "PresignClient";
    private static final Symbol presignClientSymbol = buildSymbol(PRESIGN_CLIENT, true);

    private static final String NEW_CLIENT = "NewPresignClient";
    private static final String PRESIGN_OPTIONS = "PresignOptions";
    private static final Symbol presignOptionsSymbol = buildSymbol(PRESIGN_OPTIONS, true);

    private static final String PRESIGN_OPTIONS_FROM_CLIENT_OPTIONS = "WithPresignClientFromClientOptions";
    private static final String PRESIGN_OPTIONS_FROM_EXPIRES = "WithPresignExpires";

    private static final Symbol presignerInterfaceSymbol = SymbolUtils.createPointableSymbolBuilder(
            "HTTPPresignerV4"
    ).build();
    private static final Symbol v4NewPresignerSymbol = SymbolUtils.createPointableSymbolBuilder(
            "NewSigner", AwsGoDependency.AWS_SIGNER_V4
    ).build();
    private static final Symbol v4PresignedHTTPRequestSymbol = SymbolUtils.createPointableSymbolBuilder(
            "PresignedHTTPRequest", AwsGoDependency.AWS_SIGNER_V4
    ).build();

    // constant map with service to list of operation for which presignedURL client and operation must be generated.
    private static final Map<ShapeId, Set<ShapeId>> presignedClientMap = MapUtils.of(
            ShapeId.from("com.amazonaws.s3#AmazonS3"), SetUtils.of(
                    ShapeId.from("com.amazonaws.s3#GetObject"),
                    ShapeId.from("com.amazonaws.s3#PutObject")
            ),
            ShapeId.from("com.amazonaws.sts#AWSSecurityTokenServiceV20110615"), SetUtils.of(
                    ShapeId.from("com.amazonaws.sts#GetCallerIdentity"))
    );

    // map of service to list of operations for which presignedURL client and operation should
    // be generated.
    private final Map<ShapeId, Set<ShapeId>> PRESIGNER_MAP = new TreeMap<>();

    private static final String addAsUnsignedPayloadName(String operationName) {
        return String.format("add%sPayloadAsUnsigned", operationName);
    }

    // build pointable symbols
    private static Symbol buildSymbol(String name, boolean exported) {
        if (!exported) {
            name = Character.toLowerCase(name.charAt(0)) + name.substring(1);
        }
        return SymbolUtils.createPointableSymbolBuilder(name).build();
    }

    /**
     * generates code to iterate thru func optionals and assign value into the dest variable
     *
     * @param writer GoWriter to write the code to
     * @param src    variable name that denotes functional options
     * @param dest   variable in which result of processed functional options are stored
     */
    private static final void processFunctionalOptions(
            GoWriter writer,
            String src,
            String dest
    ) {
        writer.openBlock("for _, fn := range $L {", "}", src, () -> {
            writer.write("fn(&$L)", dest);
        }).insertTrailingNewline();
    }

    @Override
    public void processFinalizedModel(GoSettings settings, Model model) {
        PRESIGNER_MAP.putAll(presignedClientMap);

        // update map for presign client/operation generation to include
        // service/operations that use PresignURLAutoFill customization class.
        Map<ShapeId, Set<ShapeId>> autofillMap = PresignURLAutoFill.SERVICE_TO_OPERATION_MAP;
        for (ShapeId service : autofillMap.keySet()) {
            if (!PRESIGNER_MAP.containsKey(service)) {
                PRESIGNER_MAP.put(service, autofillMap.get(service));
            } else {
                Set<ShapeId> operations = new TreeSet<>();
                for (ShapeId operation : PRESIGNER_MAP.get(service)) {
                    operations.add(operation);
                }
                for (ShapeId operation : autofillMap.get(service)) {
                    operations.add(operation);
                }
                PRESIGNER_MAP.put(service, operations);
            }
        }
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
            // generate presigner interface
            writePresignInterface(writer, model, symbolProvider, serviceShape);

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

                // generate s3 unsigned payload middleware helper
                writeS3AddAsUnsignedPayloadHelper(writer, model, symbolProvider, serviceShape, operationShape);
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

        writer.writeDocs(
                String.format(
                        "Presign%s is used to generate a presigned HTTP Request which contains presigned URL, signed headers "
                                + "and HTTP method used.", operationSymbol.getName())
        );
        writer.openBlock(
                "func (c *$T) Presign$T(ctx context.Context, params $P, optFns ...func($P)) "
                        + "($P, error) {", "}", presignClientSymbol, operationSymbol,
                operationInputSymbol, presignOptionsSymbol, v4PresignedHTTPRequestSymbol,
                () -> {
                    writer.write("if params == nil { params = &$T{} }", operationInputSymbol).insertTrailingNewline();

                    // process presignerOptions
                    writer.write("options := c.options.copy()");
                    processFunctionalOptions(writer, "optFns", "options");

                    writer.write("clientOptFns := append(options.ClientOptions, $L)", NOP_HTTP_CLIENT_OPTION_FUNC_NAME);
                    writer.write("");

                    writer.openBlock("result, _, err := c.client.invokeOperation(ctx, $S, params, clientOptFns,", ")",
                            operationSymbol.getName(), () -> {
                                writer.write("$L,", OperationGenerator
                                        .getAddOperationMiddlewareFuncName(operationSymbol));
                                writer.write("$L(options).$L,", CONVERT_TO_PRESIGN_TYPE_NAME,
                                        CONVERT_TO_PRESIGN_MIDDLEWARE_NAME);

                                // we should remove Content-Type header if input is a stream and
                                // payload is empty/nil stream.
                                if (operationInputShape.members().stream().anyMatch(memberShape -> {
                                    return memberShape.getMemberTrait(model, StreamingTrait.class).isPresent();
                                })) {
                                    writer.addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE);
                                    writer.addUseImports(AwsGoDependency.AWS_MIDDLEWARE);

                                    Symbol removeContentTypeHeader = SymbolUtils.createValueSymbolBuilder(
                                            "RemoveContentTypeHeader", AwsGoDependency.AWS_HTTP_TRANSPORT
                                    ).build();

                                    writer.openBlock("func(stack *middleware.Stack, options Options) error {", "},",
                                            () -> {
                                                writer.write("return $T(stack)", removeContentTypeHeader);
                                            });
                                }

                                // s3 needs to add a middleware to switch to using unsigned payload .
                                if (isS3ServiceShape(model, serviceShape)) {
                                    writer.write("$L,", addAsUnsignedPayloadName(operationSymbol.getName()));
                                }
                            });
                    writer.write("if err != nil { return nil, err }");
                    writer.write("");

                    writer.write("out := result.($P)", v4PresignedHTTPRequestSymbol);
                    writer.write("return out, nil");
                });
        writer.write("");
    }

    private void writeS3AddAsUnsignedPayloadHelper(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape serviceShape,
            OperationShape operationShape
    ) {
        // if service is not s3, return
        if (!isS3ServiceShape(model, serviceShape)) {
            return;
        }

        Symbol operationSymbol = symbolProvider.toSymbol(operationShape);

        Shape operationInputShape = model.expectShape(operationShape.getInput().get());

        writer.openBlock("func $L(stack $P, options Options) error {", "}",
                addAsUnsignedPayloadName(operationSymbol.getName()),
                SymbolUtils.createPointableSymbolBuilder("Stack", SmithyGoDependency.SMITHY_MIDDLEWARE).build(),
                () -> {
                    writer.addUseImports(AwsGoDependency.AWS_SIGNER_V4);
                    writer.write("v4.RemoveContentSHA256HeaderMiddleware(stack)");
                    writer.write("v4.RemoveComputePayloadSHA256Middleware(stack)");

                    writer.write("return $T(stack)", SymbolUtils.createValueSymbolBuilder(
                            "AddUnsignedPayloadMiddleware", AwsGoDependency.AWS_SIGNER_V4).build());
                });
        writer.write("");
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

        writer.write("type $L $T", CONVERT_TO_PRESIGN_TYPE_NAME, presignOptionsSymbol);
        writer.openBlock("func (c $L) $L(stack $P, options Options) (err error) {", "}",
                CONVERT_TO_PRESIGN_TYPE_NAME,
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

                    Symbol middlewareOptionsSymbol = SymbolUtils.createValueSymbolBuilder(
                            "PresignHTTPRequestMiddlewareOptions", AwsGoDependency.AWS_SIGNER_V4).build();
                    writer.openBlock("pmw := $T($T{", "})", presignMiddleware, middlewareOptionsSymbol, () -> {
                        writer.write("CredentialsProvider: options.$L,", AddAwsConfigFields.CREDENTIALS_CONFIG_NAME);
                        writer.write("Presigner: c.Presigner,");
                        writer.write("LogSigning: options.$L.IsSigning(),", AddAwsConfigFields.LOG_MODE_CONFIG_NAME);
                    });
                    writer.write("err = stack.Finalize.Add(pmw, $T)", smithyAfter);
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

                    // s3 service needs expires and sets unsignedPayload if input is stream
                    if (isS3ServiceShape(model, serviceShape)) {
                        writer.openBlock("if c.Expires < 0 {", "}", () -> {
                            writer.addUseImports(SmithyGoDependency.FMT);
                            writer.write(
                                    "return fmt.Errorf(\"presign URL duration must be 0 or greater, %v\", c.Expires)");
                        });
                        Symbol expiresAsHeaderMiddleware = SymbolUtils.createValueSymbolBuilder(
                                "AddExpiresOnPresignedURL",
                                AwsCustomGoDependency.S3_CUSTOMIZATION).build();
                        writer.writeDocs("add middleware to set expiration for s3 presigned url, "
                                + " if expiration is set to 0, this middleware sets a default expiration of 900 seconds");
                        writer.write("err = stack.Build.Add(&$T{ Expires: c.Expires, }, middleware.After)",
                                expiresAsHeaderMiddleware);
                        writer.write("if err != nil { return err }");
                    }

                    Symbol addAsPresignMiddlewareSymbol = SymbolUtils.createValueSymbolBuilder("AddAsIsPresigingMiddleware",
                            AwsCustomGoDependency.PRESIGNEDURL_CUSTOMIZATION).build();
                    writer.write("err = $T(stack)", addAsPresignMiddlewareSymbol);
                    writer.write("if err != nil { return err }");

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

        writer.writeDocs(String.format("%s represents the presign url client", PRESIGN_CLIENT));
        writer.openBlock("type $T struct {", "}", presignClientSymbol, () -> {
            writer.write("client *Client");
            writer.write("options $T", presignOptionsSymbol);
        });
        writer.write("");

        // generate NewPresignClient
        writer.writeDocs(
                String.format("%s generates a presign client using provided API Client and presign options",
                        NEW_CLIENT)
        );
        writer.openBlock("func $L(c *Client, optFns ...func($P)) $P {", "}",
                NEW_CLIENT, presignOptionsSymbol, presignClientSymbol, () -> {
                    writer.write("var options $T", presignOptionsSymbol);
                    processFunctionalOptions(writer, "optFns", "options");

                    writer.openBlock("if len(options.ClientOptions) != 0 {", "}", () -> {
                        writer.write("c = New(c.options, options.ClientOptions...)");
                    });
                    writer.write("");

                    writer.openBlock("if options.Presigner == nil {", "}", () -> {
                        writer.write("options.Presigner = $L(c.options)", AwsSignatureVersion4.NEW_SIGNER_FUNC_NAME);
                    });
                    writer.write("");

                    writer.openBlock("return &$L{", "}", presignClientSymbol, () -> {
                        writer.write("client: c,");
                        writer.write("options: options,");
                    });
                });
        writer.write("");
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
        // Helper function for NopClient
        writer.openBlock("func $L(o *Options) {", "}", NOP_HTTP_CLIENT_OPTION_FUNC_NAME, () -> {
            Symbol nopClientSymbol = SymbolUtils.createPointableSymbolBuilder("NopClient",
                    SmithyGoDependency.SMITHY_HTTP_TRANSPORT)
                    .build();

            writer.write("o.HTTPClient = $T{}", nopClientSymbol);
        });
        writer.write("");
    }

    /**
     * Writes the presigner interface used by the presign url client
     */
    public void writePresignInterface(
            GoWriter writer,
            Model model,
            SymbolProvider symbolProvider,
            ServiceShape serviceShape
    ) {
        Symbol signerOptionsSymbol = SymbolUtils.createPointableSymbolBuilder("SignerOptions",
                AwsGoDependency.AWS_SIGNER_V4).build();

        writer.writeDocs(
                String.format("%s represents presigner interface used by presign url client",
                        presignerInterfaceSymbol.getName())
        );
        writer.openBlock("type $T interface {", "}", presignerInterfaceSymbol, () -> {
            writer.write("PresignHTTP(");
            writer.write("ctx context.Context, credentials aws.Credentials, r *http.Request,");
            writer.write("payloadHash string, service string, region string, signingTime time.Time,");
            writer.write("optFns ...func($P),", signerOptionsSymbol);
            writer.write(") (url string, signedHeader http.Header, err error)");
        });

        writer.write("");
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

        // generate presign options
        writer.writeDocs(String.format("%s represents the presign client options", PRESIGN_OPTIONS));
        writer.openBlock("type $T struct {", "}", presignOptionsSymbol, () -> {
            writer.write("");
            writer.writeDocs(
                    "ClientOptions are list of functional options to mutate client options used by the presign client."
            );
            writer.write("ClientOptions []func(*Options)");

            writer.write("");
            writer.writeDocs("Presigner is the  presigner used by the presign url client");
            writer.write("Presigner $T", presignerInterfaceSymbol);

            // s3 service has an additional Expires options
            if (isS3ServiceShape(model, serviceShape)) {
                writer.write("");
                writer.writeDocs(
                        String.format("Expires sets the expiration duration for the generated presign url. This should "
                                + "be the duration in seconds the presigned URL should be considered valid for. If "
                                + "not set or set to zero, presign url would default to expire after 900 seconds."
                        )
                );
                writer.write("Expires time.Duration");
            }
        });
        writer.openBlock("func (o $T) copy() $T {", "}", presignOptionsSymbol, presignOptionsSymbol, () -> {
            writer.write("clientOptions := make([]func(*Options), len(o.ClientOptions))");
            writer.write("copy(clientOptions, o.ClientOptions)");
            writer.write("o.ClientOptions = clientOptions");
            writer.write("return o");
        });

        // generate WithPresignClientFromClientOptions Helper
        Symbol presignOptionsFromClientOptionsInternal = buildSymbol(PRESIGN_OPTIONS_FROM_CLIENT_OPTIONS, false);
        writer.writeDocs(
                String.format("%s is a helper utility to retrieve a function that takes PresignOption as input",
                        PRESIGN_OPTIONS_FROM_CLIENT_OPTIONS)
        );
        writer.openBlock("func $L(optFns ...func(*Options)) func($P) {", "}",
                PRESIGN_OPTIONS_FROM_CLIENT_OPTIONS, presignOptionsSymbol, () -> {
                    writer.write("return $L(optFns).options", presignOptionsFromClientOptionsInternal.getName());
                });

        writer.insertTrailingNewline();

        writer.write("type $L []func(*Options)", presignOptionsFromClientOptionsInternal.getName());
        writer.openBlock("func (w $L) options (o $P) {", "}",
                presignOptionsFromClientOptionsInternal.getName(), presignOptionsSymbol, () -> {
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
                    PRESIGN_OPTIONS_FROM_EXPIRES, presignOptionsSymbol, () -> {
                        writer.write("return $L(dur).options", presignOptionsFromExpiresInternal.getName());
                    });

            writer.insertTrailingNewline();

            writer.write("type $L time.Duration", presignOptionsFromExpiresInternal.getName());
            writer.openBlock("func (w $L) options (o $P) {", "}",
                    presignOptionsFromExpiresInternal.getName(), presignOptionsSymbol, () -> {
                        writer.write("o.Expires = time.Duration(w)");
                    }).insertTrailingNewline();
        }
    }

    private final boolean isS3ServiceShape(Model model, ServiceShape service) {
        String serviceId = service.expectTrait(ServiceTrait.class).getSdkId();
        return serviceId.equalsIgnoreCase("S3");
    }
}

