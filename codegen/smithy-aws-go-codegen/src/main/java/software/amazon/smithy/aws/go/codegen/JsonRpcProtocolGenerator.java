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

import static software.amazon.smithy.aws.go.codegen.AwsProtocolUtils.handleDecodeError;
import static software.amazon.smithy.aws.go.codegen.AwsProtocolUtils.initializeJsonDecoder;
import static software.amazon.smithy.aws.go.codegen.AwsProtocolUtils.writeJsonErrorMessageCodeDeserializer;

import java.util.Set;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.HttpRpcProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.StructureShape;

/**
 * Handles generating the aws.rest-json protocol for services.
 *
 * @inheritDoc
 *
 * @see HttpRpcProtocolGenerator
 */
abstract class JsonRpcProtocolGenerator extends HttpRpcProtocolGenerator {

    /**
     * Creates an AWS JSON RPC protocol generator
     */
    public JsonRpcProtocolGenerator() {
        super();
    }

    @Override
    protected String getOperationPath(GenerationContext context, OperationShape operation) {
        return "/";
    }

    @Override
    protected void writeDefaultHeaders(GenerationContext context, OperationShape operation, GoWriter writer) {
        super.writeDefaultHeaders(context, operation, writer);
        String target = context.getService().getId().getName() + "." + operation.getId().getName();
        writer.write("httpBindingEncoder.SetHeader(\"X-Amz-Target\").String($S)", target);
    }

    @Override
    protected void serializeInputDocument(GenerationContext context, OperationShape operation) {
        GoWriter writer = context.getWriter();
        StructureShape input = ProtocolUtils.expectInput(context.getModel(), operation);
        String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(input, getProtocolName());
        writer.addUseImports(SmithyGoDependency.SMITHY_JSON);

        // If there are no members then there's nothing to serialize
        if (input.members().size() == 0) {
            // Prevent warnings caused by input not being used
            writer.write("_ = input");
            return;
        }

        writer.addUseImports(SmithyGoDependency.SMITHY_JSON);
        writer.write("jsonEncoder := smithyjson.NewEncoder()");
        writer.openBlock("if err := $L(input, jsonEncoder.Value); err != nil {", "}", functionName, () -> {
            writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
        }).write("");

        writer.addUseImports(SmithyGoDependency.BYTES);
        writer.openBlock("if request, err = request.SetStream(bytes.NewReader(jsonEncoder.Bytes())); err != nil {",
                "}", () -> {
                    writer.write("return out, metadata, &smithy.SerializationError{Err: err}");
                });
    }

    @Override
    protected void generateDocumentBodyShapeSerializers(GenerationContext context, Set<Shape> shapes) {
        JsonShapeSerVisitor visitor = new JsonShapeSerVisitor(context);
        shapes.forEach(shape -> shape.accept(visitor));
    }

    @Override
    protected void deserializeOutputDocument(GenerationContext context, OperationShape operation) {
        GoWriter writer = context.getWriter();
        StructureShape output = ProtocolUtils.expectOutput(context.getModel(), operation);
        String functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(output, getProtocolName());
        initializeJsonDecoder(writer, "response.Body");
        AwsProtocolUtils.decodeJsonIntoInterface(writer, "out, metadata, ");
        writer.write("err = $L(&output, shape)", functionName);
        handleDecodeError(writer, "out, metadata, ");
    }


    @Override
    protected void generateDocumentBodyShapeDeserializers(GenerationContext context, Set<Shape> shapes) {
        JsonShapeDeserVisitor visitor = new JsonShapeDeserVisitor(context);
        shapes.forEach(shape -> shape.accept(visitor));
    }


    @Override
    protected void deserializeError(GenerationContext context, StructureShape shape) {
        GoWriter writer = context.getWriter();
        Symbol symbol = context.getSymbolProvider().toSymbol(shape);
        String functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(shape, getProtocolName());

        initializeJsonDecoder(writer, "errorBody");
        AwsProtocolUtils.decodeJsonIntoInterface(writer, "");
        writer.write("output := &$T{}", symbol);
        writer.write("err := $L(&output, shape)", functionName);
        writer.write("");
        handleDecodeError(writer);
        writer.write("errorBody.Seek(0, io.SeekStart)");
        writer.write("return output");
    }

    @Override
    public void generateProtocolTests(GenerationContext context) {
        AwsProtocolUtils.generateHttpProtocolTests(context);
    }

    @Override
    protected void writeErrorMessageCodeDeserializer(GenerationContext context) {
        writeJsonErrorMessageCodeDeserializer(context);
    }
}
