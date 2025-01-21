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

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;

public class AwsRetryMiddlewareHelper implements GoIntegration {
    public static final String ADD_RETRY_MIDDLEWARES_HELPER = "addRetry";

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator delegator
    ) {
        delegator.useShapeWriter(settings.getService(model), writer ->
                generateRetryMiddlewareHelpers(writer, settings.getModuleName()));
    }

    private void generateRetryMiddlewareHelpers(GoWriter writer, String moduleName) {
        writer
                .addUseImports(SmithyGoDependency.SMITHY_MIDDLEWARE)
                .addUseImports(SmithyGoDependency.SMITHY_HTTP_TRANSPORT)
                .addUseImports(AwsGoDependency.AWS_RETRY)
                .write(goTemplate("""
                        func addRetry(stack *middleware.Stack, o Options) error {
                        attempt := retry.NewAttemptMiddleware(o.Retryer, smithyhttp.RequestCloner, func(m *retry.Attempt) {
                            m.LogAttempts = o.ClientLogMode.IsRetries()
                            m.OperationMeter = o.MeterProvider.Meter($S)
                        })
                        if err := stack.Finalize.Insert(attempt, "ResolveAuthScheme", middleware.Before); err != nil {
                            return err
                        }
                        if err := stack.Finalize.Insert(&retry.MetricsHeader{}, attempt.ID(), middleware.After); err != nil {
                            return err
                        }
                        return nil
                    }""", moduleName));
    }
}
