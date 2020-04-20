/*
 * Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

import software.amazon.smithy.go.codegen.integration.DocumentMemberDeserVisitor;
import software.amazon.smithy.model.knowledge.HttpBinding;
import software.amazon.smithy.model.shapes.*;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.go.codegen.integration.DocumentMemberSerVisitor;
import software.amazon.smithy.go.codegen.integration.HttpBindingProtocolGenerator;

import java.util.List;
import java.util.Set;

/**
 * Handles general components across the AWS JSON protocols that have HTTP bindings.
 * It handles reading and writing from document bodies, including generating any
 * functions needed for performing serde.
 *
 * @see <a href="https://awslabs.github.io/smithy/spec/http.html">Smithy HTTP protocol bindings.</a>
 */
abstract class RestJsonProtocolGenerator extends HttpBindingProtocolGenerator {
    /**
     * Creates a AWS JSON RPC protocol generator.
     */
    RestJsonProtocolGenerator() {
        super(true);
    }

    @Override
    protected TimestampFormatTrait.Format getDocumentTimestampFormat() {
        return TimestampFormatTrait.Format.EPOCH_SECONDS;
    }

    @Override
    protected void generateDocumentBodyShapeSerializers(GenerationContext context, Set<Shape> shapes) {
        AwsProtocolUtils.generateDocumentBodyShapeSerde(context, shapes, new JsonShapeSerVisitor(context));
    }

    @Override
    protected void generateDocumentBodyShapeDeserializers(GenerationContext context, Set<Shape> shapes) {
        AwsProtocolUtils.generateDocumentBodyShapeSerde(context, shapes, new JsonShapeDeserVisitor(context));
    }

    @Override
    public void generateSharedComponents(GenerationContext context) {
        super.generateSharedComponents(context);
        AwsProtocolUtils.generateJsonParseBody(context);
        AwsProtocolUtils.addIdempotencyAutofillImport(context);

//        GoWriter writer = context.getWriter();
//        writer.addUseImports(getApplicationProtocol().getResponseType());
//        writer.write(IoUtils.readUtf8Resource(getClass(), "load-rest-json-error-code-stub.ts"));
    }

    @Override
    protected void writeDefaultHeaders(GenerationContext context, OperationShape operation) {
        super.writeDefaultHeaders(context, operation);
        AwsProtocolUtils.generateUnsignedPayloadSigV4Header(context, operation);
    }

    @Override
    public void serializeInputDocument(
            GenerationContext context,
            OperationShape operation,
            List<HttpBinding> documentBindings
    ) {
        // Short circuit when we have no bindings.
//        GoWriter writer = context.getWriter();
    }

    @Override
    protected void serializeInputPayload(
            GenerationContext context,
            OperationShape operation,
            HttpBinding payloadBinding
    ) {
        // We want the standard serialization, but need to alter it to JSON.
        super.serializeInputPayload(context, operation, payloadBinding);

//        GoWriter writer = context.getWriter();
//        Shape target = context.getModel().expectShape(payloadBinding.getMember().getTarget());

        // Default to an empty JSON body instead of an undefined body if
        // the target is a structure or union, and make sure any structure
        // or union content ends up as a JSON string.
//        if (target instanceof StructureShape || target instanceof UnionShape) {
//        }
    }

    private DocumentMemberSerVisitor getMemberSerVisitor(GenerationContext context, String dataSource) {
        return new JsonMemberSerVisitor(context, dataSource, getDocumentTimestampFormat());
    }

    @Override
    protected void writeErrorCodeParser(GenerationContext context) {
//        GoWriter writer = context.getWriter();

        // Outsource error code parsing since it's complex for this protocol.
    }

    @Override
    public void deserializeOutputDocument(
            GenerationContext context,
            Shape operationOrError,
            List<HttpBinding> documentBindings
    ) {
//        GoWriter writer = context.getWriter();
//        SymbolProvider symbolProvider = context.getSymbolProvider();
//
//        for (HttpBinding binding : documentBindings) {
//            Shape target = context.getModel().expectShape(binding.getMember().getTarget());
//            // The name of the member to get from the input shape.
//            String memberName = symbolProvider.toMemberName(binding.getMember());
//            // Use the jsonName trait value if present, otherwise use the member name.
//            String locationName = binding.getMember().getTrait(JsonNameTrait.class)
//                    .map(JsonNameTrait::getValue)
//                    .orElseGet(binding::getLocationName);
//        }
    }

    private DocumentMemberDeserVisitor getMemberDeserVisitor(GenerationContext context, String dataSource) {
        return new JsonMemberDeserVisitor(context, dataSource, getDocumentTimestampFormat());
    }
}
