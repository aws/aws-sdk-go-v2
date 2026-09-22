/*
 * Copyright 2026 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

import java.util.List;
import java.util.Map;
import java.util.Set;

import software.amazon.smithy.aws.traits.ServiceTrait;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.traits.LongPollTrait;
import software.amazon.smithy.utils.ListUtils;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

/**
 * Marks long-polling operations on the request context so the retry middleware
 * can apply backoff even when the retry quota is exhausted.
 */
public class LongPollingRetryIntegration implements GoIntegration {

    private static final String SET_LONG_POLLING_FUNC = "addSetLongPollingContext";

    // Hardcoded long-polling operations until the aws.api#longPoll trait is
    // applied to service models.
    private static final Map<String, Set<String>> LONG_POLLING_OPERATIONS = Map.of(
            "SQS", Set.of("ReceiveMessage"),
            "SFN", Set.of("GetActivityTask"),
            "SWF", Set.of("PollForActivityTask", "PollForDecisionTask")
    );

    private static boolean isLongPollingOperation(Model model, ServiceShape service, OperationShape operation) {
        if (operation.hasTrait(LongPollTrait.class)) {
            return true;
        }
        return service.getTrait(ServiceTrait.class)
                .map(trait -> {
                    Set<String> ops = LONG_POLLING_OPERATIONS.get(trait.getSdkId());
                    return ops != null && ops.contains(operation.getId().getName());
                })
                .orElse(false);
    }

    private static boolean serviceHasLongPollingOps(Model model, ServiceShape service) {
        return service.getAllOperations().stream()
                .flatMap(id -> model.getShape(id).stream())
                .filter(shape -> shape.isOperationShape())
                .map(shape -> shape.asOperationShape().get())
                .anyMatch(op -> isLongPollingOperation(model, service, op));
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator delegator
    ) {
        ServiceShape service = settings.getService(model);
        if (!serviceHasLongPollingOps(model, service)) {
            return;
        }

        delegator.useShapeWriter(service, this::writeSetLongPollingContext);
    }

    private void writeSetLongPollingContext(GoWriter writer) {
        writer.write(goTemplate("""
                $os:D $internalcontext:D $middleware:D
                func addSetLongPollingContext(stack *middleware.Stack, options Options) error {
                    if os.Getenv("AWS_NEW_RETRIES_2026") != "true" {
                        return nil
                    }
                    return stack.Initialize.Add(&setLongPollingContextMiddleware{}, middleware.Before)
                }

                type setLongPollingContextMiddleware struct{}

                func (*setLongPollingContextMiddleware) ID() string { return "SetLongPollingContext" }

                func (*setLongPollingContextMiddleware) HandleInitialize(
                    ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
                ) (middleware.InitializeOutput, middleware.Metadata, error) {
                    ctx = internalcontext.SetIsLongPolling(ctx, true)
                    return next.HandleInitialize(ctx, in)
                }
                """, Map.of(
                "internalcontext", AwsGoDependency.INTERNAL_CONTEXT,
                "middleware", SmithyGoDependency.SMITHY_MIDDLEWARE,
                "os", SmithyGoDependency.OS
        )));
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate(LongPollingRetryIntegration::isLongPollingOperation)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(
                                        SET_LONG_POLLING_FUNC).build())
                                .useClientOptions()
                                .build())
                        .build()
        );
    }
}
