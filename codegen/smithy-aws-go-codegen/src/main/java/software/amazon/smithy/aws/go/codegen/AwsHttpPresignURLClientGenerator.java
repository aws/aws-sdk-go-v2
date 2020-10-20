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

import java.util.ArrayList;
import java.util.List;
import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.OperationGenerator;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.utils.SmithyBuilder;

public class AwsHttpPresignURLClientGenerator {
    private static final String CONVERT_TO_PRESIGN_MIDDLEWARE_NAME = "convertToPresignMiddleware";

    private final Model model;
    private final SymbolProvider symbolProvider;

    private final Symbol presignClientSymbol;
    private final Symbol newPresignClientSymbol;

    private final OperationShape operation;
    private final Symbol operationSymbol;
    private final Shape operationInput;
    private final Symbol operationInputSymbol;

    private final boolean exported;

    private final List<Symbol> convertToPresignMiddlewareHelpers;

    private AwsHttpPresignURLClientGenerator(Builder builder) {
        this.exported = builder.exported;

        this.model = SmithyBuilder.requiredState("model", builder.model);
        this.symbolProvider = SmithyBuilder.requiredState("symbolProvider", builder.symbolProvider);
        this.convertToPresignMiddlewareHelpers = builder.convertToPresignMiddlewareHelpers;

        this.operation = SmithyBuilder.requiredState("operation", builder.operation);
        this.operationSymbol = symbolProvider.toSymbol(operation);

        this.operationInput = ProtocolUtils.expectInput(model, operation);
        this.operationInputSymbol = symbolProvider.toSymbol(operationInput);

        this.presignClientSymbol = buildPresignClientSymbol(operationSymbol, exported);
        this.newPresignClientSymbol = buildNewPresignClientSymbol(operationSymbol, exported);
    }

    /**
     * Writes the Presign client's type and methods.
     *
     * @param writer writer to write to
     */
    public void writePresignClientType(GoWriter writer) {
        writer.addUseImports(SmithyGoDependency.CONTEXT);
        writer.addUseImports(AwsGoDependency.AWS_SIGNER_V4);

        writer.openBlock("type $T struct {", "}", presignClientSymbol, () -> {
            writer.write("client *Client");
            writer.write("presigner *v4.Signer");
        });

        writer.openBlock("func $L(options Options, optFns ...func(*Options)) *$T {", "}",
                newPresignClientSymbol.getName(),
                presignClientSymbol,
                () -> {
                    writer.openBlock("return &$T{", "}", presignClientSymbol, () -> {
                        writer.write("client: New(options, optFns...),");

                        writer.addUseImports(AwsGoDependency.AWS_SIGNER_V4);
                        writer.write("presigner: v4.NewSigner(),");
                    });
                });

        writer.addUseImports(SmithyGoDependency.NET_HTTP);
        writer.openBlock(
                // TODO presign with expire can be supported with a builder param that adds an additional expires param to presign signature.

                // TODO Should this return a v4.PresignedHTTPRequest type instead of individual fields?
                "func (c *$T) Presign$T(ctx context.Context, params $P, optFns ...func(*Options)) "
                        + "(string, http.Header, error) {",
                "}",
                presignClientSymbol, operationSymbol, operationInputSymbol,
                () -> {
                    Symbol nopClient = SymbolUtils.createPointableSymbolBuilder("NopClient",
                            SmithyGoDependency.SMITHY_HTTP_TRANSPORT)
                            .build();

                    writer.write("if params == nil { params = &$T{} }", operationInputSymbol);
                    writer.write("");

                    // TODO could be replaced with a `WithAPIOptions` client option helper.
                    // TODO could be replaced with a `WithHTTPClient` client option helper.
                    writer.openBlock("optFns = append(optFns, func(o *Options) {", "})", () -> {
                        writer.write("o.HTTPClient = &$T{}", nopClient);
                    });
                    writer.write("");

                    Symbol withIsPresigning = SymbolUtils.createValueSymbolBuilder("WithIsPresigning",
                            AwsCustomGoDependency.PRESIGNEDURL_CUSTOMIZATION).build();

                    writer.write("ctx = $T(ctx)", withIsPresigning);
                    writer.openBlock("result, _, err := c.client.invokeOperation(ctx, $S, params, optFns,", ")",
                            operationSymbol.getName(),
                            () -> {
                        writer.write("$L,", OperationGenerator.getAddOperationMiddlewareFuncName(operationSymbol));
                        writer.write("c.$L,", CONVERT_TO_PRESIGN_MIDDLEWARE_NAME);
                    });
                    writer.write("if err != nil { return ``, nil, err }");
                    writer.write("");

                    writer.write("out := result.(*v4.PresignedHTTPRequest)");
                    writer.write("return out.URL, out.SignedHeader, nil");
                });

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
                    Symbol requestInvocationIDMiddleware = SymbolUtils.createValueSymbolBuilder(
                            "RequestInvocationIDMiddleware",
                            AwsGoDependency.AWS_MIDDLEWARE)
                            .build();

                    Symbol presignMiddleware = SymbolUtils.createValueSymbolBuilder("NewPresignHTTPRequestMiddleware",
                            AwsGoDependency.AWS_SIGNER_V4)
                            .build();


                    // Middleware to add
                    writer.write("stack.Finalize.Clear()");
                    writer.write("stack.Deserialize.Clear()");
                    writer.write("stack.Build.Remove($T{}.ID())", requestInvocationIDMiddleware);

                    writer.write("err = stack.Finalize.Add($T(options.Credentials, c.presigner), $T)", presignMiddleware,
                            smithyAfter);
                    writer.write("if err != nil { return err }");

                    convertToPresignMiddlewareHelpers.forEach((symbol) -> {
                        writer.write("err = $T(stack)", symbol);
                        writer.write("if err != nil { return err }");
                    });

                    writer.write("return nil");
                });
    }

    public Symbol getPresignClientSymbol() {
        return presignClientSymbol;
    }

    public Symbol getNewPresignClientSymbol() {
        return newPresignClientSymbol;
    }

    private static Symbol buildNewPresignClientSymbol(Symbol operation, boolean exported) {
        String name = String.format("New%sHTTPPresignURLClient", operation.getName());
        return buildSymbol(name, exported);
    }

    private static Symbol buildPresignClientSymbol(Symbol operation, boolean exported) {
        String name = String.format("%sHTTPPresignURLClient", operation.getName());
        return buildSymbol(name, exported);
    }

    private static Symbol buildAPIClientSymbol(Symbol operation, boolean exported) {
        String name = String.format("%sAPIClient", operation.getName());
        return buildSymbol(name, exported);
    }

    private static Symbol buildSymbol(String name, boolean exported) {
        if (!exported) {
            name = Character.toLowerCase(name.charAt(0)) + name.substring(1);
        }
        return SymbolUtils.createValueSymbolBuilder(name).
                build();
    }

    /**
     * Builder for the HTTP Presign URL client client generator.
     */
    public static class Builder implements SmithyBuilder<AwsHttpPresignURLClientGenerator> {
        private Model model;
        private SymbolProvider symbolProvider;
        private OperationShape operation;
        private boolean exported;
        private List<Symbol> convertToPresignMiddlewareHelpers = new ArrayList<>();

        /**
         * Sets the model for the builder
         * @param model API model
         * @return builder
         */
        public Builder model(Model model) {
            this.model = model;
            return this;
        }

        /**
         * Sets the symbol provider for the builder
         * @param symbolProvider the symbol provider
         * @return buidler
         */
        public Builder symbolProvider(SymbolProvider symbolProvider) {
            this.symbolProvider = symbolProvider;
            return this;
        }

        /**
         * Sets the operation for the builder
         * @param operation api operation
         * @return builder
         */
        public Builder operation(OperationShape operation) {
            this.operation = operation;
            return this;
        }

        /**
         * Sets that the generated client type should be exported, defaults to false.
         * @return builder
         */
        public Builder exported() {
            return this.exported(true);
        }

        /**
         * Sets if the generate client type should be exported or not.
         * @param exported if exported
         * @return builder
         */
        public Builder exported(boolean exported) {
            this.exported = exported;
            return this;
        }

        /**
         * Sets additional middleware mutator that will be generated into the client's convert to presign URL operation.
         * Used by the client to convert a API operation to a presign URL.
         * @param middlewareHelpers list of middleware helpers to set
         * @return builder
         */
        public Builder convertToPresignMiddlewareHelpers(List<Symbol> middlewareHelpers) {
            this.convertToPresignMiddlewareHelpers.clear();
            this.convertToPresignMiddlewareHelpers.addAll(middlewareHelpers);
            return this;
        }

        /**
         * Adds a single middleware mutator that will be generated into the client's convert to presign URL operation.
         * Used by the client to convert API operation to a presigned URL.
         * @param middlewareHelper the middleware helper to add
         * @return builder
         */
        public Builder addConvertToPresignMiddlewareHelpers(Symbol middlewareHelper) {
            this.convertToPresignMiddlewareHelpers.add(middlewareHelper);
            return this;
        }

        // TODO presign with expire can be supported with a builder param that enables expires param behavior.

        public AwsHttpPresignURLClientGenerator build() {
            return new AwsHttpPresignURLClientGenerator(this);
        }
    }
}
