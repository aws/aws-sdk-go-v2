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

package software.amazon.smithy.aws.go.codegen.customization;

import java.util.Map;
import java.util.Set;

import software.amazon.smithy.aws.go.codegen.AwsGoDependency;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.GoCodegenContext;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.Waiters2;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * Extends the base smithy Waiters integration to track in the User-Agent string.
 */
public class AwsWaiters extends Waiters2 {
    @Override
    public Set<Symbol> getAdditionalClientOptions() {
        return Set.of(buildPackageSymbol("addIsWaiterUserAgent"));
    }

    @Override
    public void writeAdditionalFiles(GoCodegenContext ctx) {
        super.writeAdditionalFiles(ctx);

        ctx.writerDelegator().useFileWriter("api_client.go", ctx.settings().getModuleName(), goTemplate("""
                func addIsWaiterUserAgent(o *Options) {
                    o.APIOptions = append(o.APIOptions, func(stack $stack:P) error {
                        ua, err := getOrAddRequestUserAgent(stack)
                        if err != nil {
                            return err
                        }

                        ua.AddUserAgentFeature($featureWaiter:T)
                        return nil
                    })
                }""",
                Map.of(
                        "stack", SmithyGoDependency.SMITHY_MIDDLEWARE.struct("Stack"),
                        "featureWaiter", AwsGoDependency.AWS_MIDDLEWARE.valueSymbol("UserAgentFeatureWaiter")
                )));
    }
}
