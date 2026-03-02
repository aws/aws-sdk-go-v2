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

import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;

/**
 * Utility methods for generating AWS protocols.
 */
public final class AwsProtocolUtils {
    private AwsProtocolUtils() {
    }

    public static void initializeJsonDecoder(GoWriter writer, String bodyLocation) {
        // Use a ring buffer and tee reader to help in pinpointing any deserialization errors.
        writer.addUseImports(SmithyGoDependency.SMITHY_IO);
        writer.write("var buff [1024]byte");
        writer.write("ringBuffer := smithyio.NewRingBuffer(buff[:])");
        writer.write("");

        writer.addUseImports(SmithyGoDependency.IO);
        writer.addUseImports(SmithyGoDependency.JSON);
        writer.write("body := io.TeeReader($L, ringBuffer)", bodyLocation);
        writer.write("decoder := json.NewDecoder(body)");
        writer.write("decoder.UseNumber()");
    }

    /**
     * Decodes JSON into {@code shape} with type {@code interface{}} using the encoding/json decoder
     * referenced by {@code decoder}.
     *
     * @param writer            GoWriter to write code to
     * @param errorReturnExtras extra parameters to return if an error occurs
     */
    public static void decodeJsonIntoInterface(GoWriter writer, String errorReturnExtras) {
        writer.write("var shape interface{}");
        writer.addUseImports(SmithyGoDependency.IO);
        writer.openBlock("if err := decoder.Decode(&shape); err != nil && err != io.EOF {", "}", () -> {
            wrapAsDeserializationError(writer);
            writer.write("return $Lerr", errorReturnExtras);
        });
        writer.write("");
    }

    /**
     * Wraps the Go error {@code err} in a {@code DeserializationError} with a snapshot
     *
     * @param writer
     */
    private static void wrapAsDeserializationError(GoWriter writer) {
        writer.write("var snapshot bytes.Buffer");
        writer.write("io.Copy(&snapshot, ringBuffer)");
        writer.openBlock("err = &smithy.DeserializationError {", "}", () -> {
            writer.write("Err: fmt.Errorf(\"failed to decode response body, %w\", err),");
            writer.write("Snapshot: snapshot.Bytes(),");
        });
    }

    public static void handleDecodeError(GoWriter writer, String returnExtras) {
        writer.openBlock("if err != nil {", "}", () -> {
            writer.addUseImports(SmithyGoDependency.BYTES);
            writer.addUseImports(SmithyGoDependency.SMITHY);
            writer.addUseImports(SmithyGoDependency.IO);
            wrapAsDeserializationError(writer);
            writer.write("return $Lerr", returnExtras);
        }).write("");
    }

    public static void handleDecodeError(GoWriter writer) {
        handleDecodeError(writer, "");
    }

    public static void writeJsonEventMessageSerializerDelegator(
            GenerationContext ctx,
            String functionName,
            String operand,
            String contentType
    ) {
        var writer = ctx.getWriter().get();

        var stringValue = SymbolUtils.createValueSymbolBuilder("StringValue",
                AwsGoDependency.SERVICE_INTERNAL_EVENTSTREAM).build();
        var contentTypeHeader = SymbolUtils.createValueSymbolBuilder("ContentTypeHeader",
                AwsGoDependency.SERVICE_INTERNAL_EVENTSTREAMAPI).build();

        writer.write("msg.Headers.Set($T, $T($S))",
                contentTypeHeader, stringValue, contentType);
        var newEncoder = SymbolUtils.createValueSymbolBuilder("NewEncoder",
                SmithyGoDependency.SMITHY_JSON).build();
        writer.write("jsonEncoder := $T()", newEncoder)
                .openBlock("if err := $L($L, jsonEncoder.Value); err != nil {", "}", functionName, operand,
                        () -> writer.write("return err"))
                .write("msg.Payload = jsonEncoder.Bytes()")
                .write("return nil");
    }

    public static void initializeJsonEventMessageDeserializer(GenerationContext ctx) {
        initializeJsonEventMessageDeserializer(ctx, "");
    }

    public static void initializeJsonEventMessageDeserializer(GenerationContext ctx, String errorReturnExtras) {
        var writer = ctx.getWriter().get();
        writer.write("br := $T(msg.Payload)", SymbolUtils.createValueSymbolBuilder(
                "NewReader", SmithyGoDependency.BYTES).build());
        initializeJsonDecoder(writer, "br");
        AwsProtocolUtils.decodeJsonIntoInterface(writer, errorReturnExtras);
    }

    public static void writeJsonEventStreamUnknownExceptionDeserializer(GenerationContext ctx) {
        var writer = ctx.getWriter().get();
        writer.write("br := $T(msg.Payload)", SymbolUtils.createValueSymbolBuilder("NewReader",
                SmithyGoDependency.BYTES).build());
        AwsProtocolUtils.initializeJsonDecoder(writer, "br");
        writer.write("""
                     code, message, err := $T(decoder)
                     if err != nil {
                         return err
                     }
                     errorCode := "UnknownError"
                     errorMessage := errorCode
                     if ev := exceptionType.String(); len(ev) > 0 {
                         errorCode = ev
                     } else if ev := code; len(ev) > 0 {
                         errorCode = ev
                     }
                     if ev := message; len(ev) > 0 {
                         errorMessage = ev
                     }
                     return &$T{
                         Code: errorCode,
                         Message: errorMessage,
                     }
                     """,
                SymbolUtils.createValueSymbolBuilder("GetErrorInfo",
                        AwsGoDependency.AWS_REST_JSON_PROTOCOL).build(),
                SymbolUtils.createValueSymbolBuilder("GenericAPIError", SmithyGoDependency.SMITHY).build());
    }
}
