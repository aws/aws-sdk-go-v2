/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package software.amazon.smithy.aws.go.codegen.customization.s3;

import software.amazon.smithy.go.codegen.GoCodegenContext;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.shapes.Shape;
import software.amazon.smithy.model.shapes.ShapeId;

import java.util.Map;
import java.util.Optional;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

/**
 *  CodeGen class for a presigner that uses POST
 *  This is a special kind of presigner that instead of generating and signing an HTTP request,
 *  generates a URL and field values to be used on an HTML form.
 * <p>
 *  This method of sending request is only implemented service side for putting objects into an S3 bucket.
 *  However, it could be expanded to other operations if other services start supporting it
 *
 *  @see <a href="https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectPOST.html> the public docs</a>
 */
public class PresignPostCustomization implements GoIntegration {

    // This is the only operation that supports this
    private final ShapeId PUT_OBJECT = ShapeId.from("com.amazonaws.s3#PutObject");

    @Override
    public void writeAdditionalFiles(GoCodegenContext ctx) {
        Optional<Shape> maybeShape = ctx.model().getShape(PUT_OBJECT);
        if (maybeShape.isEmpty()) {
            return;
        }

        Shape shape = maybeShape.get();
        ctx.writerDelegator().useShapeWriter(shape, writer -> {
            generatePreSigner(shape, writer);
        });

    }

    private void generatePreSigner(Shape shape, GoWriter writer) {
        String typeName = shape.getId().getName();
        GoWriter.Writable goTemplate = goTemplate("""
            // PresignPostObject is a special kind of [presigned request] used to send a request using
            // form data, likely from an HTML form on a browser.
            // Unlike other presigned operations, the return values of this function are not meant to be used directly
            // to make an HTTP request but rather to be used as inputs to a form. See [the docs] for more information
            // on how to use these values
            //
            // [presigned request] https://docs.aws.amazon.com/AmazonS3/latest/userguide/ShareObjectPreSignedURL.html
            // [the docs] https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectPOST.html
            func (c *PresignClient) PresignPostObject(ctx context.Context, params *$type:LInput, optFns ...func(*PresignPostOptions)) (*PresignedPostRequest, error) {
                if params == nil {
                    params = &$type:LInput{}
                }
                clientOptions := c.options.copy()
                options := PresignPostOptions{
                    Expires:       clientOptions.Expires,
                    PostPresigner: &postSignAdapter{},
                }
                for _, fn := range optFns {
                    fn(&options)
                }
                clientOptFns := append(clientOptions.ClientOptions, withNopHTTPClientAPIOption)
                cvt := presignPostConverter(options)
                result, _, err := c.client.invokeOperation(ctx, "$type:L", params, clientOptFns,
                    c.client.addOperationPutObjectMiddlewares,
                    cvt.ConvertToPresignMiddleware,
                    func(stack *middleware.Stack, options Options) error {
                        return awshttp.RemoveContentTypeHeader(stack)
                    },
                )
                if err != nil {
                    return nil, err
                }

                out := result.(*PresignedPostRequest)
                return out, nil
            }""", Map.of(
                    "type", typeName
                )
        );
        writer.write(goTemplate);
    }
}
