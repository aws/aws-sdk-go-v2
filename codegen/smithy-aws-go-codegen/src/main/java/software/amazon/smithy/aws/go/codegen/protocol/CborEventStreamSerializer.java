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

import static software.amazon.smithy.aws.go.codegen.AwsEventStreamUtils.getEventStreamMessageRequestSerializerName;
import static software.amazon.smithy.aws.go.codegen.AwsEventStreamUtils.getEventStreamMessageResponseDeserializerName;
import static software.amazon.smithy.aws.go.codegen.AwsEventStreamUtils.getEventStreamSerializerName;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildSymbol;
import static software.amazon.smithy.go.codegen.serde.cbor.CborDeserializerGenerator.getDeserializerName;
import static software.amazon.smithy.go.codegen.serde.cbor.CborSerializerGenerator.getSerializerName;

import java.util.Map;
import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.shapes.UnionShape;
import software.amazon.smithy.model.traits.StreamingTrait;
import software.amazon.smithy.utils.MapUtils;

final class CborEventStreamSerializer implements GoWriter.Writable {
    private static final Map<String, Object> templateEnv = MapUtils.of(
            "cborEncode", SmithyGoDependency.SMITHY_CBOR.func("Encode"),
            "eventstreamMessage", AwsCustomGoDependency.SERVICE_INTERNAL_EVENTSTREAM.struct("Message"),
            "eventstreamStringValue", AwsCustomGoDependency.SERVICE_INTERNAL_EVENTSTREAM.valueSymbol("StringValue"),
            "eventstreamapiEventTypeHeader", AwsCustomGoDependency.SERVICE_INTERNAL_EVENTSTREAMAPI.valueSymbol("EventTypeHeader"),
            "fmtErrorf", SmithyGoDependency.FMT.func("Errorf")
    );

    private final ProtocolGenerator.GenerationContext ctx;

    private final StructureShape input;
    private final UnionShape stream;

    public CborEventStreamSerializer(ProtocolGenerator.GenerationContext ctx, OperationShape operation) {
        this.ctx = ctx;

        this.input = ctx.getModel().expectShape(operation.getInputShape(), StructureShape.class);
        this.stream = input.members().stream()
                .filter(it -> StreamingTrait.isEventStream(ctx.getModel(), it))
                .map(it -> ctx.getModel().expectShape(it.getTarget(), UnionShape.class))
                .findFirst().orElseThrow(() -> new CodegenException("operation must have an input event stream"));
    }

    @Override
    public void accept(GoWriter writer) {
        writer.write(serializeInitialRequest());
        writer.write(serializeMessage());
    }

    private GoWriter.Writable serializeInitialRequest() {
        return goTemplate("""
                func $fn:L(input interface{}, msg $eventstreamMessage:P) error {
                    in, ok := input.($input:P)
                    if !ok {
                        return $fmtErrorf:T("unexpected input type %T", input)
                    }
                    cv, err := $serialize:L(in)
                    if err != nil {
                        return err
                    }
                    msg.Payload = $cborEncode:T(cv)
                    return nil
                }
                """,
                templateEnv,
                MapUtils.of(
                        "fn", getEventStreamMessageRequestSerializerName(input, ctx.getService(), ctx.getProtocolName()),
                        "input", ctx.getSymbolProvider().toSymbol(input),
                        "serialize", getSerializerName(input)
                ));
    }

    private GoWriter.Writable serializeMessage() {
        return goTemplate("""
                func $fn:L(v $union:T, msg $eventstreamMessage:P) error {
                    switch vv := v.(type) {
                    $variants:W
                    default:
                        return $fmtErrorf:T("unexpected event message type: %T", v)
                    }
                }
                """,
                templateEnv,
                MapUtils.of(
                        "fn", getEventStreamSerializerName(stream, ctx.getService(), ctx.getProtocolName()),
                        "union", ctx.getSymbolProvider().toSymbol(stream),
                        "variants", GoWriter.ChainWritable.of(
                                stream.members().stream()
                                        .map(this::serializeEventVariant)
                                        .toList()
                        ).compose(false)
                ));
    }

    private GoWriter.Writable serializeEventVariant(MemberShape variant) {
        var variantSymbol = buildSymbol(
                ctx.getSymbolProvider().toMemberName(variant),
                ctx.getSymbolProvider().toSymbol(variant).getNamespace()
        ).toBuilder()
                .putProperty(SymbolUtils.POINTABLE, true)
                .build();
        return goTemplate("""
                case $variant:P:
                    msg.Headers.Set($eventstreamapiEventTypeHeader:T, $eventstreamStringValue:T($variantName:S))
                    cv, err := $serialize:L(&vv.Value)
                    if err != nil {
                        return err
                    }
                    msg.Payload = $cborEncode:T(cv)
                    return nil
                """,
                templateEnv,
                MapUtils.of(
                        "variant", variantSymbol,
                        "variantName", variant.getMemberName(),
                        "serialize", getSerializerName(ctx.getModel().expectShape(variant.getTarget()))
                ));
    }
}
