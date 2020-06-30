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

import java.util.Set;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.HttpRpcProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.go.codegen.integration.ProtocolUtils;
import software.amazon.smithy.model.Model;
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
        super(true);
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
    protected void serializeInputDocument(
            Model model, SymbolProvider symbolProvider, OperationShape operation,
            GoStackStepMiddlewareGenerator generator, GoWriter writer
    ) {
        StructureShape input = ProtocolUtils.expectInput(model, operation);
        String functionName = ProtocolGenerator.getDocumentSerializerFunctionName(input, getProtocolName());

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
    protected void deserializeOutputDocument(
            Model model,
            SymbolProvider symbolProvider,
            OperationShape operation,
            GoStackStepMiddlewareGenerator generator,
            GoWriter writer
    ) {
        StructureShape output = ProtocolUtils.expectOutput(model, operation);
        String functionName = ProtocolGenerator.getDocumentDeserializerFunctionName(output, getProtocolName());

        writer.addUseImports(SmithyGoDependency.SMITHY_IO);
        writer.write("buff := make([]byte, 1024)");
        writer.write("ringBuffer := smithyio.NewRingBuffer(buff)");
        writer.write("");

        writer.addUseImports(SmithyGoDependency.IO);
        writer.write("body := io.TeeReader(response.Body, ringBuffer)");
        writer.write("defer response.Body.Close()");
        writer.write("");

        writer.addUseImports(SmithyGoDependency.JSON);
        writer.write("decoder := json.NewDecoder(body)");
        writer.write("decoder.UseNumber()");
        writer.write("");

        writer.write("err = $L(&output, decoder)", functionName);
        writer.openBlock("if err != nil {", "}", () -> {
            writer.addUseImports(SmithyGoDependency.BYTES);
            writer.addUseImports(SmithyGoDependency.SMITHY);
            writer.write("var snapshot bytes.Buffer");
            writer.write("io.Copy(&snapshot, ringBuffer)");
            writer.openBlock("return out, metadata, &smithy.DeserializationError {", "}", () -> {
                writer.write("Err: fmt.Errorf(\"failed to decode response body with invalid JSON, %w\", err),");
                writer.write("Snapshot: snapshot.Bytes(),");
            });
        }).write("");
    }

    @Override
    protected void generateDocumentBodyShapeDeserializers(GenerationContext context, Set<Shape> shapes) {
        JsonShapeDeserVisitor visitor = new JsonShapeDeserVisitor(context);
        shapes.forEach(shape -> shape.accept(visitor));
    }

    @Override
    public void generateProtocolTests(GenerationContext context) {
        AwsProtocolUtils.generateHttpProtocolTests(context);
    }
}
