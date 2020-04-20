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

import software.amazon.smithy.aws.traits.auth.UnsignedPayloadTrait;
import software.amazon.smithy.model.knowledge.HttpBindingIndex;
import software.amazon.smithy.model.knowledge.NeighborProviderIndex;
import software.amazon.smithy.model.neighbor.Walker;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeVisitor;
import software.amazon.smithy.model.traits.IdempotencyTokenTrait;
import software.amazon.smithy.model.traits.TimestampFormatTrait.Format;
import software.amazon.smithy.model.traits.XmlNamespaceTrait;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.HttpProtocolGeneratorUtils;
import software.amazon.smithy.go.codegen.integration.ProtocolGenerator.GenerationContext;
import software.amazon.smithy.utils.IoUtils;

import java.util.Optional;
import java.util.Set;
import java.util.TreeSet;

import static software.amazon.smithy.model.knowledge.HttpBinding.Location.DOCUMENT;

/**
 * Utility methods for generating AWS protocols.
 */
final class AwsProtocolUtils {

    private AwsProtocolUtils() {
    }

    /**
     * Writes an {@code 'x-amz-content-sha256' = 'UNSIGNED_PAYLOAD'} header for an
     * {@code @aws.api#unsignedPayload} trait that specifies the {@code "aws.v4"} auth scheme.
     *
     * @param context   The generation context.
     * @param operation The operation being generated.
     * @see <a href=https://awslabs.github.io/smithy/spec/aws-core.html#aws-api-unsignedpayload-trait>@aws.api#unsignedPayload trait</a>
     */
    static void generateUnsignedPayloadSigV4Header(GenerationContext context, OperationShape operation) {
    }

    /**
     * Writes a serde function for a set of shapes using the passed visitor.
     * This will walk the input set of shapes and invoke the visitor for any
     * members of aggregate shapes in the set.
     *
     * @param context The generation context.
     * @param shapes  A list of shapes to generate serde for, including their members.
     * @param visitor A ShapeVisitor that generates a serde function for shapes.
     * @see software.amazon.smithy.go.codegen.integration.DocumentShapeSerVisitor
     */
    static void generateDocumentBodyShapeSerde(
            GenerationContext context,
            Set<Shape> shapes,
            ShapeVisitor<Void> visitor
    ) {
        // Walk all the shapes within those in the document and generate for them as well.
    }

    /**
     * Writes a response body parser function for JSON protocols. This
     * will parse a present body after converting it to utf-8.
     *
     * @param context The generation context.
     */
    static void generateJsonParseBody(GenerationContext context) {
    }

    /**
     * Writes a response body parser function for XML protocols. This
     * will parse a present body after converting it to utf-8.
     *
     * @param context The generation context.
     */
    static void generateXmlParseBody(GenerationContext context) {
    }

    /**
     * Writes a form urlencoded string builder function for query based protocols.
     * This will escape the keys and values, combine those with an '=', and combine
     * those strings with an '&'.
     *
     * @param context The generation context.
     */
    static void generateBuildFormUrlencodedString(GenerationContext context) {
    }

    /**
     * Writes a default body for query-based operations when the operation doesn't
     * have an input defined.
     *
     * @param context   The generation context.
     * @param operation The operation being generated for.
     * @return That a body variable was generated and should be set on the request.
     */
    static boolean generateUndefinedQueryInputBody(GenerationContext context, OperationShape operation) {
        return true;
    }

    /**
     * Writes an attribute containing information about a Shape's optionally specified
     * XML namespace configuration to an attribute of the passed node name.
     *
     * @param context  The generation context.
     * @param shape    The shape to apply the namespace attribute to, if present on it.
     * @param nodeName The node to apply the namespace attribute to.
     * @return Returns if an XML namespace attribute was written.
     */
    static boolean writeXmlNamespace(GenerationContext context, Shape shape, String nodeName) {
        return true;
    }

    /**
     * Imports a UUID v4 generating function used for auto-filling idempotency tokens.
     *
     * @param context The generation context.
     */
    static void addIdempotencyAutofillImport(GenerationContext context) { }

    /**
     * Writes a statement that auto-fills the value of a member that is an idempotency
     * token if it is undefined at the time of serialization.
     *
     * @param context       The generation context.
     * @param memberShape   The member that may be an idempotency token.
     * @param inputLocation The location of input data for the member.
     */
    static void writeIdempotencyAutofill(GenerationContext context, MemberShape memberShape, String inputLocation) {
        if (memberShape.hasTrait(IdempotencyTokenTrait.class)) {
        }
    }

    /**
     * Gets a value provider for the timestamp member handling proper serialization
     * formatting.
     *
     * @param context       The generation context.
     * @param memberShape   The member that needs timestamp serialization.
     * @param defaultFormat The timestamp format to default to.
     * @param inputLocation The location of input data for the member.
     * @return A string representing the proper value provider for this timestamp.
     */
    static String getInputTimestampValueProvider(
            GenerationContext context,
            MemberShape memberShape,
            Format defaultFormat,
            String inputLocation
    ) {
        return "";
    }
}
