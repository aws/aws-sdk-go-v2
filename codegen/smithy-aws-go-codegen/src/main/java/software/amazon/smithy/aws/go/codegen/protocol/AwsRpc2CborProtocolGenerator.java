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

package software.amazon.smithy.aws.go.codegen.protocol;

import static software.amazon.smithy.go.codegen.protocol.ProtocolUtil.hasEventStream;

import software.amazon.smithy.aws.go.codegen.AwsEventStreamUtils;
import software.amazon.smithy.aws.go.codegen.AwsFnProvider;
import software.amazon.smithy.aws.go.codegen.AwsProtocolUtils;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.endpoints.EndpointResolutionGenerator;
import software.amazon.smithy.go.codegen.protocol.rpc2.cbor.Rpc2CborProtocolGenerator;
import software.amazon.smithy.model.knowledge.TopDownIndex;

/**
 * Extension of the smithy-borne Rpc2CborProtocolGenerator to do protocol tests and event streams since that's currently
 * a 2000+ line dumpster fire in this repo.
 */
public final class AwsRpc2CborProtocolGenerator extends Rpc2CborProtocolGenerator {
    @Override
    public void generateProtocolTests(GenerationContext context) {
        AwsProtocolUtils.generateHttpProtocolTests(context);
    }

    @Override
    public void generateEventStreamComponents(GenerationContext context) {
        // This automagically wires up ALL the framing logic for both directions of streams. All we have to do is fill
        // in the serde elsewhere (it's different signatures than normal request/response), see:
        // * CborEventStreamSerializer
        // * CborEventStreamDeserializer
        AwsEventStreamUtils.generateEventStreamComponents(context);
    }

    @Override
    public void generateEndpointResolution(GenerationContext context) {
        new EndpointResolutionGenerator(new AwsFnProvider()).generate(context);
    }

    @Override
    public void generateSharedSerializerComponents(GenerationContext ctx) {
        super.generateSharedSerializerComponents(ctx);

        var model = ctx.getModel();
        var streamSerializers = TopDownIndex.of(model).getContainedOperations(ctx.getService()).stream()
                .filter(it -> hasEventStream(model, model.expectShape(it.getOutputShape())))
                .map(it -> (GoWriter.Writable) new CborEventStreamSerializer(ctx, it))
                .toList();
        ctx.getWriter().get().write(GoWriter.ChainWritable.of(streamSerializers).compose());
    }

    @Override
    public void generateSharedDeserializerComponents(GenerationContext ctx) {
        super.generateSharedDeserializerComponents(ctx);

        var model = ctx.getModel();
        var streamDeserializers = TopDownIndex.of(model).getContainedOperations(ctx.getService()).stream()
                .filter(it -> hasEventStream(model, model.expectShape(it.getOutputShape())))
                .map(it -> (GoWriter.Writable) new CborEventStreamDeserializer(ctx, it))
                .toList();
        ctx.getWriter().get().write(GoWriter.ChainWritable.of(streamDeserializers).compose());
    }
}
