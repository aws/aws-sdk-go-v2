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

package software.amazon.smithy.aws.go.codegen.customization.service.s3;

import static software.amazon.smithy.aws.go.codegen.customization.service.s3.S3ModelUtils.isServiceS3;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

import java.util.List;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoDependency;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.MemberShape;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;
import software.amazon.smithy.model.shapes.StringShape;
import software.amazon.smithy.model.traits.DeprecatedTrait;
import software.amazon.smithy.model.traits.DocumentationTrait;
import software.amazon.smithy.model.traits.HttpHeaderTrait;
import software.amazon.smithy.model.traits.OutputTrait;
import software.amazon.smithy.model.transform.ModelTransformer;
import software.amazon.smithy.utils.MapUtils;

/**
 * Restrictions around timestamp formatting for the 'Expires' value in some S3 responses has never been standardized and
 * thus many non-conforming values for the field (unsupported formats, arbitrary strings, etc.) exist in the wild. This
 * customization makes the response parsing forgiving for this field in responses and adds an ExpiresString field that
 * contains the unparsed value.
 */
public class S3ExpiresShapeCustomization implements GoIntegration {
    private static final ShapeId S3_EXPIRES = ShapeId.from("com.amazonaws.s3#Expires");
    private static final ShapeId S3_EXPIRES_STRING = ShapeId.from("com.amazonaws.s3#ExpiresString");
    private static final String DESERIALIZE_S3_EXPIRES = "deserializeS3Expires";

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return List.of(RuntimeClientPlugin.builder()
                .addShapeDeserializer(S3_EXPIRES, buildPackageSymbol(DESERIALIZE_S3_EXPIRES))
                .build());
    }

    @Override
    public Model preprocessModel(Model model, GoSettings settings) {
        if (!isServiceS3(model, settings.getService(model))) {
            return model;
        }

        var withExpiresString = model.toBuilder()
                .addShape(StringShape.builder()
                        .id(S3_EXPIRES_STRING)
                        .build())
                .build();
        return ModelTransformer.create().mapShapes(withExpiresString, this::addExpiresString);
    }

    @Override
    public void writeAdditionalFiles(GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator) {
        goDelegator.useFileWriter("deserializers.go", settings.getModuleName(), deserializeS3Expires());
    }

    private Shape addExpiresString(Shape shape) {
        if (!shape.hasTrait(OutputTrait.class)) {
            return shape;
        }

        var expires = shape.getMember(S3_EXPIRES.getName());
        if (expires.isEmpty()) {
            return shape;
        }

        if (!expires.get().getTarget().equals(S3_EXPIRES)) {
            return shape;
        }

        var deprecated = DeprecatedTrait.builder()
                .message("This field is handled inconsistently across AWS SDKs. Prefer using the ExpiresString field " +
                        "which contains the unparsed value from the service response.")
                .build();
        var stringDocs = new DocumentationTrait("The unparsed value of the Expires field from the service " +
                "response. Prefer use of this value over the normal Expires response field where possible.");
        return Shape.shapeToBuilder(shape)
                .addMember(expires.get().toBuilder()
                        .addTrait(deprecated)
                        .build())
                .addMember(MemberShape.builder()
                        .id(shape.getId().withMember(S3_EXPIRES_STRING.getName()))
                        .target(S3_EXPIRES_STRING)
                        .addTrait(expires.get().expectTrait(HttpHeaderTrait.class)) // copies header name
                        .addTrait(stringDocs)
                        .build())
                .build();
    }

    private GoWriter.Writable deserializeS3Expires() {
        return goTemplate("""
                func $name:L(v string) ($time:P, error) {
                    t, err := $parseHTTPDate:T(v)
                    if err != nil {
                        return nil, nil
                    }
                    return &t, nil
                }
                """,
                MapUtils.of(
                        "name", DESERIALIZE_S3_EXPIRES,
                        "time", SmithyGoDependency.TIME.struct("Time"),
                        "parseHTTPDate", SmithyGoDependency.SMITHY_TIME.func("ParseHTTPDate")
                ));
    }
}
