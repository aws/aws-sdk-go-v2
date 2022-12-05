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

import static software.amazon.smithy.aws.go.codegen.AwsProtocolUtils.writeAwsQueryErrorCodeDeserializer;

import software.amazon.smithy.aws.traits.protocols.AwsJson1_0Trait;
import software.amazon.smithy.aws.traits.protocols.AwsQueryCompatibleTrait;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.model.shapes.ShapeId;

/**
 * Handles generating the awsJson1_0 protocol for services.
 *
 * @inheritDoc
 *
 * @see JsonRpcProtocolGenerator
 */
final class AwsJsonRpc1_0 extends JsonRpcProtocolGenerator {

    @Override
    protected String getDocumentContentType() {
        return "application/x-amz-json-1.0";
    }

    @Override
    public ShapeId getProtocol() {
        return AwsJson1_0Trait.ID;
    }

    @Override
    protected void writeErrorMessageCodeDeserializer(GenerationContext context) {
        super.writeErrorMessageCodeDeserializer(context);
        if (context.getService().hasTrait(AwsQueryCompatibleTrait.class)) {
            // writeAwsQueryErrorCodeDeserializer(context);
        }
    }

    @Override
    public void generateSharedDeserializerComponents(GenerationContext context) {
        super.generateSharedDeserializerComponents(context);
        if (context.getService().hasTrait(AwsQueryCompatibleTrait.class)) {
            GoWriter writer = context.getWriter().get();
            writer.openBlock("func getAwsQueryErrorCode(response *smithyhttp.Response) {", "}", () -> {
                writer.write("queryCodeHeader := response.Header.Get(\"x-amzn-query-error\")");
                writer.openBlock("if queryCodeHeader != \"\" {", "}", () -> {
                    writer.write("queryCodeParts := strings.Split(queryCodeHeader, \";\")");
                    writer.openBlock("if queryCodeParts != nil && len(queryCodeParts) == 2 {", "}", () -> {
                        writer.write("return queryCodeParts[0]");
                    });
                    writer.write("return \"\"");
                });
            });
        }
    }
}
