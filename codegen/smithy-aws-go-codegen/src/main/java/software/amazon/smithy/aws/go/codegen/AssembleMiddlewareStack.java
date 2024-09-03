/*
 * Copyright 2024 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

import static software.amazon.smithy.aws.go.codegen.AwsSignatureVersion4.hasSigV4X;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

import java.util.List;
import java.util.Map;

import software.amazon.smithy.aws.go.codegen.customization.AdjustAwsRestJsonContentType;
import software.amazon.smithy.aws.traits.auth.UnsignedPayloadTrait;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.go.codegen.middleware.FinalizeStepMiddleware;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.EventStreamIndex;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

public class AssembleMiddlewareStack implements GoIntegration {

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first.
     *
     * @return Returns the sort order, defaults to -40.
     */
    @Override
    public byte getOrder() {
        return -40;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                // Add RequestInvocationIDMiddleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(
                                MiddlewareRegistrar.builder()
                                        .resolvedFunction(buildPackageSymbol("addClientRequestID"))
                                        .build()
                        ).build(),

                // Add ContentLengthMiddleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) ->
                                EventStreamIndex.of(model).getInputInfo(operation).isEmpty())
                        .registerMiddleware(
                                MiddlewareRegistrar.builder()
                                        .resolvedFunction(buildPackageSymbol("addComputeContentLength"))
                                        .build()
                        ).build(),

                // Add endpoint serialize middleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        EndpointGenerator.ADD_MIDDLEWARE_HELPER_NAME).build())
                                .useClientOptions()
                                .build())
                        .build(),

                // Add streaming events payload middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            if (!hasSigV4X(
                                    model, service, operation)) {
                                return false;
                            }
                            return EventStreamIndex.of(model).getInputInfo(operation).isPresent();
                        })
                        .registerMiddleware(
                                MiddlewareRegistrar.builder()
                                        .resolvedFunction(buildPackageSymbol("addStreamingEventsPayload"))
                                        .build()
                        ).build(),

                // Add unsigned payload middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            if (!hasSigV4X(
                                    model, service, operation)) {
                                return false;
                            }
                            var noEventStream = EventStreamIndex.of(model).getInputInfo(operation).isEmpty();
                            return operation.hasTrait(UnsignedPayloadTrait.class) && noEventStream;
                        })
                        .registerMiddleware(
                                MiddlewareRegistrar.builder()
                                        .resolvedFunction(buildPackageSymbol("addUnsignedPayload"))
                                        .build()
                        ).build(),

                // Add signed payload middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            if (!hasSigV4X(
                                    model, service, operation)) {
                                return false;
                            }
                            var noEventStream = EventStreamIndex.of(model).getInputInfo(operation).isEmpty();
                            return !operation.hasTrait(UnsignedPayloadTrait.class) && noEventStream;
                        })
                        .registerMiddleware(
                                MiddlewareRegistrar.builder()
                                        .resolvedFunction(buildPackageSymbol("addComputePayloadSHA256"))
                                        .build()
                        ).build(),

                // Add content-sha256 payload header middleware to operation stack
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> {
                            if (!hasSigV4X(
                                    model, service, operation)) {
                                return false;
                            }
                            var hasEventStream = EventStreamIndex.of(model).getInputInfo(operation).isPresent();
                            return operation.hasTrait(UnsignedPayloadTrait.class) || hasEventStream;
                        })
                        .registerMiddleware(
                                MiddlewareRegistrar.builder()
                                        .resolvedFunction(buildPackageSymbol("addContentSHA256Header"))
                                        .build()
                        ).build(),

                // Add retryer middleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                                AwsRetryMiddlewareHelper.ADD_RETRY_MIDDLEWARES_HELPER)
                                        .build())
                                .useClientOptions()
                                .build())
                        .build(),

                // Add middleware to store raw response omn metadata
                RuntimeClientPlugin.builder()
                        .registerMiddleware(
                                MiddlewareRegistrar.builder()
                                        .resolvedFunction(buildPackageSymbol("addRawResponseToMetadata"))
                                        .build()
                        ).build(),

                // Add recordResponseTiming middleware to operation stack
                RuntimeClientPlugin.builder()
                        .registerMiddleware(
                                MiddlewareRegistrar.builder()
                                        .resolvedFunction(buildPackageSymbol("addRecordResponseTiming"))
                                        .build()
                        ).build(),

                // wrap the retry loop in a span
                RuntimeClientPlugin.builder()
                        .registerMiddleware(
                                MiddlewareRegistrar.builder()
                                        .resolvedFunction(buildPackageSymbol("addSpanRetryLoop"))
                                        .useClientOptions()
                                        .build()
                        ).build(),

                // Add Client UserAgent
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createPointableSymbolBuilder(
                                        AwsClientUserAgent.MIDDLEWARE_RESOLVER).build())
                                .useClientOptions()
                                .build())
                        .build(),

                // Add REST-JSON Content-Type Adjuster
                RuntimeClientPlugin.builder()
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        AdjustAwsRestJsonContentType.RESOLVER_NAME).build())
                                .build())
                        .servicePredicate((model, serviceShape) ->
                                AdjustAwsRestJsonContentType.isServiceOnShameList(serviceShape))
                        .build(),

                // Add Event Stream Input Writer (must be added AFTER retryer)
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) ->
                                EventStreamIndex.of(model).getInputInfo(operation).isPresent())
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                                "AddInitializeStreamWriter",
                                                AwsGoDependency.SERVICE_INTERNAL_EVENTSTREAMAPI)
                                        .build())
                                .build())
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator) {
        goDelegator.useFileWriter("api_client.go", settings.getModuleName(), writer -> {
            writer.write(addMiddleware());
            writer.write(spanRetryLoopMiddleware());
            if (hasSigV4X(model, settings.getService(model))) {
                writer.write(addSigV4XMiddleware());
            }
        });
    }

    private GoWriter.Writable spanRetryLoopMiddleware() {
        return new FinalizeStepMiddleware() {
            public String getStructName() {
                return "spanRetryLoop";
            }

            public Map<String, Symbol> getFields() {
                return Map.of("options", buildPackageSymbol("Options"));
            }

            public GoWriter.Writable getFuncBody() {
                return goTemplate("""
                        tracer := operationTracer(m.options.TracerProvider)
                        ctx, span := tracer.StartSpan(ctx, "RetryLoop")
                        defer span.End()

                        return next.HandleFinalize(ctx, in)
                        """);
            }
        };
    }

    private GoWriter.Writable addMiddleware() {
        return goTemplate("""
                $D $D $D
                func addClientRequestID(stack *middleware.Stack) error {
                    return stack.Build.Add(&awsmiddleware.ClientRequestID{}, middleware.After)
                }

                func addComputeContentLength(stack *middleware.Stack) error {
                    return stack.Build.Add(&smithyhttp.ComputeContentLength{}, middleware.After)
                }

                func addRawResponseToMetadata(stack *middleware.Stack) error {
                    return stack.Deserialize.Add(&awsmiddleware.AddRawResponse{}, middleware.Before)
                }

                func addRecordResponseTiming(stack *middleware.Stack) error {
                    return stack.Deserialize.Add(&awsmiddleware.RecordResponseTiming{}, middleware.After)
                }

                func addSpanRetryLoop(stack *middleware.Stack, options Options) error {
                    return stack.Finalize.Insert(&spanRetryLoop{options: options}, "Retry", middleware.Before)
                }
                """, SmithyGoDependency.SMITHY_MIDDLEWARE, AwsGoDependency.AWS_MIDDLEWARE, SmithyGoDependency.SMITHY_HTTP_TRANSPORT);
    }

    private GoWriter.Writable addSigV4XMiddleware() {
        return goTemplate("""
                $D
                func addStreamingEventsPayload(stack *middleware.Stack) error {
                    return stack.Finalize.Add(&v4.StreamingEventsPayload{}, middleware.Before)
                }

                func addUnsignedPayload(stack *middleware.Stack) error {
                    return stack.Finalize.Insert(&v4.UnsignedPayload{}, "ResolveEndpointV2", middleware.After)
                }

                func addComputePayloadSHA256(stack *middleware.Stack) error {
                    return stack.Finalize.Insert(&v4.ComputePayloadSHA256{}, "ResolveEndpointV2", middleware.After)
                }

                func addContentSHA256Header(stack *middleware.Stack) error {
                    return stack.Finalize.Insert(&v4.ContentSHA256Header{}, (*v4.ComputePayloadSHA256)(nil).ID(), middleware.After)
                }
                """, AwsGoDependency.AWS_SIGNER_V4);
    }
}
