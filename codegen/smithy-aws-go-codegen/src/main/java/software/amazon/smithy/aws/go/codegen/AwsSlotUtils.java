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

import java.util.List;
import software.amazon.smithy.go.codegen.GoDependency;
import software.amazon.smithy.go.codegen.MiddlewareIdentifier;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.StackSlotRegistrar;

public class AwsSlotUtils {
    public static MiddlewareIdentifier awsSymbolId(String name) {
        return symbolId(name, AwsGoDependency.AWS_MIDDLEWARE_ID);
    }

    public static MiddlewareIdentifier smithySymbolId(String name) {
        return symbolId(name, SmithyGoDependency.SMITHY_MIDDLEWARE_ID);
    }

    public static MiddlewareIdentifier symbolId(String name, GoDependency dependency) {
        return MiddlewareIdentifier.symbol(SymbolUtils.createValueSymbolBuilder(name, dependency).build());
    }

    public static StackSlotRegistrar.SlotMutator addBefore(List<MiddlewareIdentifier> identifiers) {
        return StackSlotRegistrar.SlotMutator.addBefore().identifiers(identifiers).build();
    }

    public static StackSlotRegistrar.SlotMutator addAfter(List<MiddlewareIdentifier> identifiers) {
        return StackSlotRegistrar.SlotMutator.addAfter().identifiers(identifiers).build();
    }

    public static StackSlotRegistrar.SlotMutator insertBefore(
            MiddlewareIdentifier relativeTo,
            List<MiddlewareIdentifier> identifiers
    ) {
        return StackSlotRegistrar.SlotMutator.insertBefore(relativeTo).identifiers(identifiers).build();
    }

    public static StackSlotRegistrar.SlotMutator insertAfter(
            MiddlewareIdentifier relativeTo,
            List<MiddlewareIdentifier> identifiers
    ) {
        return StackSlotRegistrar.SlotMutator.insertAfter(relativeTo).identifiers(identifiers).build();
    }
}
