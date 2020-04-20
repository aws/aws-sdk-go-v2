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
import software.amazon.smithy.model.traits.TimestampFormatTrait.Format;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.DocumentMemberDeserVisitor;
import software.amazon.smithy.go.codegen.integration.DocumentShapeDeserVisitor;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;

import java.util.Map;
import java.util.TreeMap;

/**
 * Visitor to generate deserialization functions for shapes in AWS JSON protocol
 * document bodies.
 *
 * No standard visitation methods are overridden; function body generation for all
 * expected deserializers is handled by this class.
 *
 * Timestamps are deserialized from {@link Format}.EPOCH_SECONDS by default.
 */
final class JsonShapeDeserVisitor extends DocumentShapeDeserVisitor {

    JsonShapeDeserVisitor(GenerationContext context) {
        super(context);
    }

    private DocumentMemberDeserVisitor getMemberVisitor(String dataSource) {
        return new JsonMemberDeserVisitor(getContext(), dataSource, Format.EPOCH_SECONDS);
    }

    @Override
    protected void deserializeCollection(GenerationContext context, CollectionShape shape) {
    }

    @Override
    protected void deserializeDocument(GenerationContext context, DocumentShape shape) {
    }

    @Override
    protected void deserializeMap(GenerationContext context, MapShape shape) {
    }

    @Override
    protected void deserializeStructure(GenerationContext context, StructureShape shape) {
    }

    @Override
    protected void deserializeUnion(GenerationContext context, UnionShape shape) {
    }
}
