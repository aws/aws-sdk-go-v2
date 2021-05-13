/*
 * Copyright 2021 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package software.amazon.smithy.aws.go.codegen.customization;

import java.util.Set;
import software.amazon.smithy.codegen.core.Symbol;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.MiddlewareIdentifier;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.ToShapeId;
import software.amazon.smithy.utils.SetUtils;

/**
 * Generates a middleware to fix the AWS REST-JSON Content-Type header from the standard protocol header
 * to the json 1.1 header.
 */
public class AdjustAwsRestJsonContentType implements GoIntegration {
    public static final String RESOLVER_NAME = "addRestJsonContentTypeCustomization";

    private static final String MIDDLEWARE_NAME = "customizeRestJsonContentType";
    private static final GoStackStepMiddlewareGenerator SERIALIZE_STEP_MIDDLEWARE =
            GoStackStepMiddlewareGenerator.createSerializeStepMiddleware(MIDDLEWARE_NAME,
                    MiddlewareIdentifier.builder().name(MIDDLEWARE_NAME).build());
    private static final String CONTENT_TYPE_HEADER = "Content-Type";
    private static final String EXPECTED_CONTENT_TYPE = "application/json";
    private static final String TARGET_CONTENT_TYPE = "application/x-amz-json-1.1";
    private static final String INSERT_AFTER = "OperationSerializer";

    private static final Set<ShapeId> SHAME_SET = SetUtils.of(
            ShapeId.from("com.amazonaws.finspace#AWSHabaneroManagementService"),
            ShapeId.from("com.amazonaws.finspacedata#AWSHabaneroPublicAPI")
    );

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        if (!isServiceOnShameList(settings.getService())) {
            return;
        }

        goDelegator.useShapeWriter(settings.getService(model), writer -> {
            generateMiddleware(writer);
            generateAddMiddleware(writer);
        });
    }

    /**
     * Determine if a service should be customized.
     * @param service the service shape id
     * @return whether the service requires customization
     */
    public static boolean isServiceOnShameList(ToShapeId service) {
        return SHAME_SET.contains(service.toShapeId());
    }

    private void generateAddMiddleware(GoWriter writer) {
        Symbol stackSymbol = SymbolUtils.createPointableSymbolBuilder("Stack",
                SmithyGoDependency.SMITHY_MIDDLEWARE).build();
        writer.write("");
        writer.openBlock("func $L(stack $P) error {", "}", RESOLVER_NAME, stackSymbol, () -> {
            writer.openBlock("return stack.Serialize.Insert(&$T{}, $S, $T)",
                    SERIALIZE_STEP_MIDDLEWARE.getMiddlewareSymbol(), INSERT_AFTER,
                    SymbolUtils.createValueSymbolBuilder("After", SmithyGoDependency.SMITHY_MIDDLEWARE).build());
        });
    }

    private void generateMiddleware(GoWriter writer) {
        SERIALIZE_STEP_MIDDLEWARE.writeMiddleware(writer, (g, w) -> {
            Symbol requestSymbol = SymbolUtils.createPointableSymbolBuilder("Request",
                    SmithyGoDependency.SMITHY_HTTP_TRANSPORT).build();

            w.write("req, ok := in.Request.($P)", requestSymbol);
            w.openBlock("if !ok {", "}", () -> {
                w.addUseImports(SmithyGoDependency.FMT);
                w.write("return out, metadata, fmt.Errorf(\"unknown transport type %T\", in.Request)");
            });
            w.write("");

            w.write("const contentType = $S", CONTENT_TYPE_HEADER);
            w.write("const expectedType = $S", EXPECTED_CONTENT_TYPE);
            w.write("const targetType = $S", TARGET_CONTENT_TYPE);
            w.write("");

            w.addUseImports(SmithyGoDependency.STRINGS);
            w.openBlock("if strings.EqualFold(req.Header.Get(contentType), expectedType) {", "}",
                    () -> {
                        w.write("req.Header.Set(contentType, targetType)");
                    });
            w.write("");

            w.write("return next.$L(ctx, in)", g.getHandleMethodName());
        });
    }
}
