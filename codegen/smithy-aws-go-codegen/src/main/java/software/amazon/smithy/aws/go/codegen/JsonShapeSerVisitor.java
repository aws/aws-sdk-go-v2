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

import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.*;
import software.amazon.smithy.model.traits.JsonNameTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait.Format;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.DocumentMemberSerVisitor;
import software.amazon.smithy.go.codegen.integration.DocumentShapeSerVisitor;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;

import java.util.Map;
import java.util.TreeMap;

/**
 * Visitor to generate serialization functions for shapes in AWS JSON protocol
 * document bodies.
 *
 * This class handles function body generation for all types expected by the {@code
 * DocumentShapeSerVisitor}. No other shape type serialization is overridden.
 *
 * Timestamps are serialized to {@link Format}.EPOCH_SECONDS by default.
 */
final class JsonShapeSerVisitor extends DocumentShapeSerVisitor {
    private static final Format TIMESTAMP_FORMAT = Format.EPOCH_SECONDS;

    JsonShapeSerVisitor(GenerationContext context) {
        super(context);
    }

    private DocumentMemberSerVisitor getMemberVisitor(String dataSource) {
        return new JsonMemberSerVisitor(getContext(), dataSource, TIMESTAMP_FORMAT);
    }

    @Override
    public void serializeCollection(GenerationContext context, CollectionShape shape) {
    }

    @Override
    public void serializeDocument(GenerationContext context, DocumentShape shape) {
    }

    @Override
    public void serializeMap(GenerationContext context, MapShape shape) {
    }

    @Override
    public void serializeStructure(GenerationContext context, StructureShape shape) {
    }

    @Override
    public void serializeUnion(GenerationContext context, UnionShape shape) {
    }
}
