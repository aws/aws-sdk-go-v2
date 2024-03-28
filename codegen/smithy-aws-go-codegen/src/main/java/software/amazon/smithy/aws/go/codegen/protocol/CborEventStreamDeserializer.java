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

import static software.amazon.smithy.aws.go.codegen.AwsEventStreamUtils.getEventStreamDeserializerName;
import static software.amazon.smithy.aws.go.codegen.AwsEventStreamUtils.getEventStreamExceptionDeserializerName;
import static software.amazon.smithy.aws.go.codegen.AwsEventStreamUtils.getEventStreamMessageResponseDeserializerName;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildSymbol;
import static software.amazon.smithy.go.codegen.serde.cbor.CborDeserializerGenerator.getDeserializerName;

import java.util.Map;
import software.amazon.smithy.aws.go.codegen.customization.AwsCustomGoDependency;
import software.amazon.smithy.codegen.core.CodegenException;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.StructureShape;
import software.amazon.smithy.model.shapes.UnionShape;
import software.amazon.smithy.model.traits.ErrorTrait;
import software.amazon.smithy.model.traits.StreamingTrait;
import software.amazon.smithy.utils.MapUtils;

final class CborEventStreamDeserializer implements GoWriter.Writable {
    private static final Map<String, Object> templateEnv = MapUtils.of(
            "bytesNewBuffer", SmithyGoDependency.BYTES.func("NewBuffer"),
            "cborDecode", SmithyGoDependency.SMITHY_CBOR.func("Decode"),
            "cborMap", SmithyGoDependency.SMITHY_CBOR.func("Map"),
            "eventstreamMessage", AwsCustomGoDependency.SERVICE_INTERNAL_EVENTSTREAM.pointableSymbol("Message"),
            "eventstreamNewEncoder", AwsCustomGoDependency.SERVICE_INTERNAL_EVENTSTREAM.func("NewEncoder"),
            "eventstreamapiEventTypeHeader", AwsCustomGoDependency.SERVICE_INTERNAL_EVENTSTREAMAPI.valueSymbol("EventTypeHeader"),
            "eventstreamapiExceptionTypeHeader", AwsCustomGoDependency.SERVICE_INTERNAL_EVENTSTREAMAPI.valueSymbol("ExceptionTypeHeader"),
            "fmtErrorf", SmithyGoDependency.FMT.func("Errorf"),
            "smithyGenericAPIError", SmithyGoDependency.SMITHY.pointableSymbol("GenericAPIError"),
            "stringsEqualFold", SmithyGoDependency.STRINGS.valueSymbol("EqualFold")
    );

    private final ProtocolGenerator.GenerationContext ctx;
    private final StructureShape output;
    private final UnionShape stream;

    public CborEventStreamDeserializer(ProtocolGenerator.GenerationContext ctx, OperationShape operation) {
        this.ctx = ctx;

        this.output = ctx.getModel().expectShape(operation.getOutputShape(), StructureShape.class);
        this.stream = output.members().stream()
                .filter(it -> StreamingTrait.isEventStream(ctx.getModel(), it))
                .map(it -> ctx.getModel().expectShape(it.getTarget(), UnionShape.class))
                .findFirst().orElseThrow(() -> new CodegenException("operation must have an output event stream"));
    }

    @Override
    public void accept(GoWriter writer) {
        writer.write(deserializeInitialResponse());
        writer.write(deserializeEvent());
        writer.write(deserializeError());
    }

    private GoWriter.Writable deserializeInitialResponse() {
        return goTemplate("""
                func $fn:L(msg $eventstreamMessage:P) (interface{}, error) {
                    cv, err := $cborDecode:T(msg.Payload)
                    if err != nil {
                        return nil, err
                    }

                    return $deserialize:L(cv)
                }
                """,
                templateEnv,
                MapUtils.of(
                        "fn", getEventStreamMessageResponseDeserializerName(output, ctx.getService(), ctx.getProtocolName()),
                        "deserialize", getDeserializerName(output)
                ));
    }

    private GoWriter.Writable deserializeEvent() {
        return goTemplate("""
                func $fn:L(v *$union:T, msg $eventstreamMessage:P) error {
                    cv, err := $cborDecode:T(msg.Payload)
                    if err != nil {
                        return err
                    }
                    
                    typ := msg.Headers.Get($eventstreamapiEventTypeHeader:T)
                    if typ == nil {
                        return $fmtErrorf:T("%s event header not present", $eventstreamapiEventTypeHeader:T)
                    }
                    
                    switch {
                    $variants:W
                    default:
                        buffer := $bytesNewBuffer:T(nil)
                        $eventstreamNewEncoder:T().Encode(buffer, *msg)
                        *v = &types.UnknownUnionMember{
                            Tag:   typ.String(),
                            Value: buffer.Bytes(),
                        }
                        return nil
                    }
                }
                """,
                templateEnv,
                MapUtils.of(
                        "fn", getEventStreamDeserializerName(stream, ctx.getService(), ctx.getProtocolName()),
                        "union", ctx.getSymbolProvider().toSymbol(stream),
                        "variants", GoWriter.ChainWritable.of(
                                stream.members().stream()
                                        .filter(it -> !ctx.getModel().expectShape(it.getTarget()).hasTrait(ErrorTrait.class))
                                        .map(this::deserializeEventVariant)
                                        .toList()
                        ).compose(false)
                ));
    }

    private GoWriter.Writable deserializeError() {
        return goTemplate("""
                func $fn:L(msg $eventstreamMessage:P) error {
                    cv, err := $cborDecode:T(msg.Payload)
                    if err != nil {
                        return err
                    }

                    typ := msg.Headers.Get($eventstreamapiExceptionTypeHeader:T)
                    if typ == nil {
                        return $fmtErrorf:T("%s event header not present", $eventstreamapiExceptionTypeHeader:T)
                    }

                    switch {
                    $variants:W
                    default:
                        code, msg, _, err := getProtocolErrorInfo(msg.Payload)
                        if err != nil {
                            return err
                        }
                        
                        if len(code) == 0 {
                            code = typ.String()
                        }
                        if len(code) == 0 {
                            code = "UnknownError"
                        }
                        if len(msg) == 0 {
                            msg = "UnknownError"
                        }
                        return &$smithyGenericAPIError:T{Code: code, Message: msg}
                    }
                }
                """,
                templateEnv,
                MapUtils.of(
                        "fn", getEventStreamExceptionDeserializerName(stream, ctx.getService(), ctx.getProtocolName()),
                        "variants", GoWriter.ChainWritable.of(
                                stream.members().stream()
                                        .filter(it -> ctx.getModel().expectShape(it.getTarget()).hasTrait(ErrorTrait.class))
                                        .map(this::deserializeErrorVariant)
                                        .toList()
                        ).compose(false)
                ));
    }

    private GoWriter.Writable deserializeEventVariant(MemberShape variant) {
        var variantSymbol = buildSymbol(
                ctx.getSymbolProvider().toMemberName(variant),
                ctx.getSymbolProvider().toSymbol(variant).getNamespace()
        );
        return goTemplate("""
                case $stringsEqualFold:T(typ.String(), $variantName:S):
                     vv, err := $deserialize:L(cv)
                     if err != nil {
                         return err
                     }
                     *v = &$variantStruct:T{Value: *vv}
                     return nil
                """,
                templateEnv,
                MapUtils.of(
                        "variantName", variant.getMemberName(),
                        "variantStruct", variantSymbol,
                        "deserialize", getDeserializerName(ctx.getModel().expectShape(variant.getTarget()))
                ));
    }

    private GoWriter.Writable deserializeErrorVariant(MemberShape variant) {
        return goTemplate("""
                case $stringsEqualFold:T(typ.String(), $variantName:S):
                     verr, err := $deserialize:L(cv)
                     if err != nil {
                         return err
                     }
                     return verr
                """,
                templateEnv,
                MapUtils.of(
                        "variantName", variant.getMemberName(),
                        "deserialize", getDeserializerName(ctx.getModel().expectShape(variant.getTarget()))
                ));
    }
}
