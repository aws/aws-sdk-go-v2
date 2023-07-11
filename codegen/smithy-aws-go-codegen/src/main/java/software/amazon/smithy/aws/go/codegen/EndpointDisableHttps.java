/*
 * Copyright 2023 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
import java.util.List;

import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.MiddlewareIdentifier;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.endpoints.EndpointMiddlewareGenerator;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ToShapeId;

/*
 * Adds support for non-SSL endpoints during endpoint resolution.
 * The new Rules Engine endpoint resolution doesnt support non-SSL endpoints.
 * So this middleware exists for backwards compatibility with legacy
 * endpoint resolution. It is operation specific because it is being inserted
 * directly after the operation-specific endpoint resolution middleware.
 */
public class EndpointDisableHttps implements GoIntegration {

        private final List<RuntimeClientPlugin> runtimeClientPlugins = new ArrayList<>();

        public static String getAddMiddlewareFuncName(String operationName) {
                return String.format("add%sEndpointDisableHTTPSMiddleware", operationName);
        }

        public static String getMiddlewareObjectName(String operationName) {
                return String.format("op%sEndpointDisableHTTPSMiddleware", operationName);
        }

        /**
         * Gets the sort order of the customization from -128 to 127, with lowest
         * executed first. Needs to execute after Rules Engine endpoint
         * resolution middleware insertion.
         *
         * @return Returns the sort order, defaults to 127.
         */
        @Override
        public byte getOrder() {
                return 127;
        }

        @Override
        public List<RuntimeClientPlugin> getClientPlugins() {
                return runtimeClientPlugins;
        }

        @Override
        public void processFinalizedModel(GoSettings settings, Model model) {

                TopDownIndex topDownIndex = TopDownIndex.of(model);
                var serviceShape = settings.getService(model);

                for (ToShapeId operation : topDownIndex.getContainedOperations(serviceShape)) {
                        OperationShape operationShape = model.expectShape(operation.toShapeId(), OperationShape.class);
                        String operationName = operationShape.getId().getName();

                        runtimeClientPlugins.add(RuntimeClientPlugin.builder()
                                        .servicePredicate((m, s) -> s.equals(serviceShape))
                                        .operationPredicate((m, s, o) -> o.equals(operationShape))
                                        .registerMiddleware(MiddlewareRegistrar.builder()
                                                        .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                                                        getAddMiddlewareFuncName(operationName))
                                                                        .build())
                                                        .useClientOptions()
                                                        .build())
                                        .build());

                }
        }

        @Override
        public void writeAdditionalFiles(
                        GoSettings settings,
                        Model model,
                        SymbolProvider symbolProvider,
                        GoDelegator goDelegator) {

                TopDownIndex topDownIndex = TopDownIndex.of(model);
                var serviceShape = settings.getService(model);

                for (ToShapeId operation : topDownIndex.getContainedOperations(serviceShape)) {
                        OperationShape operationShape = model.expectShape(operation.toShapeId(), OperationShape.class);

                        goDelegator.useShapeWriter(operationShape, writer -> {
                                Symbol operationSymbol = symbolProvider.toSymbol(operationShape);
                                String operationName = operationSymbol.getName();

                                GoStackStepMiddlewareGenerator middleware = GoStackStepMiddlewareGenerator
                                                .createSerializeStepMiddleware(
                                                                getMiddlewareObjectName(operationName),
                                                                MiddlewareIdentifier.string(getMiddlewareObjectName(
                                                                                operationName)));
                                middleware.writeMiddleware(writer, this::generateMiddlewareResolverBody,
                                                this::generateMiddlewareStructureMembers);

                                writer.write(
                                                """
                                                                        func $L(stack $P, o Options) error {
                                                                                return stack.Serialize.Insert(&$L{
                                                                                        EndpointDisableHTTPS: o.EndpointOptions.DisableHTTPS,
                                                                                }, \"$L\", middleware.After)
                                                                        }
                                                                """,
                                                getAddMiddlewareFuncName(operationName),
                                                SymbolUtils.createPointableSymbolBuilder("Stack",
                                                                SmithyGoDependency.SMITHY_MIDDLEWARE).build(),
                                                getMiddlewareObjectName(operationName),
                                                EndpointMiddlewareGenerator.MIDDLEWARE_ID);
                                writer.write("");

                        });
                }
        }

        private void generateMiddlewareResolverBody(GoStackStepMiddlewareGenerator g, GoWriter writer) {
                writer.write(
                                """
                                                        req, ok := in.Request.($P)
                                                        if !ok {
                                                                return out, metadata, $T(\"unknown transport type %T\", in.Request)
                                                        }

                                                        if m.EndpointDisableHTTPS {
                                                                req.URL.Scheme = \"http\"
                                                        }

                                                        return next.HandleSerialize(ctx, in)
                                                """,
                                SymbolUtils.createPointableSymbolBuilder("Request",
                                                SmithyGoDependency.SMITHY_HTTP_TRANSPORT).build(),
                                SymbolUtils.createValueSymbolBuilder("Errorf", SmithyGoDependency.FMT).build());
        }

        private void generateMiddlewareStructureMembers(GoStackStepMiddlewareGenerator g, GoWriter writer) {
                writer.write("EndpointDisableHTTPS $L", "bool");
        }

}
