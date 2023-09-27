/*
 * Copyright 2023 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

import java.util.List;

import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStdlibTypes;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.MiddlewareIdentifier;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.go.codegen.endpoints.EndpointMiddlewareGenerator;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.utils.ListUtils;

import static software.amazon.smithy.go.codegen.GoStackStepMiddlewareGenerator.createFinalizeStepMiddleware;
import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;
import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * Adds support for non-SSL endpoints during endpoint resolution.
 * The new Rules Engine endpoint resolution doesnt support non-SSL endpoints.
 * So this middleware exists for backwards compatibility with legacy
 * endpoint resolution.
 */
public class EndpointDisableHttps implements GoIntegration {
    public static final String MIDDLEWARE_NAME = "disableHTTPSMiddleware";
    public static final String MIDDLEWARE_ID = "disableHTTPS";
    public static final String MIDDLEWARE_ADDER = "addDisableHTTPSMiddleware";

    private static final MiddlewareRegistrar DISABLE_HTTPS_MIDDLEWARE =
            MiddlewareRegistrar.builder()
                    .resolvedFunction(buildPackageSymbol(MIDDLEWARE_ADDER))
                    .useClientOptions()
                    .build();

    /**
     * Gets the sort order of the customization from -128 to 127, with lowest
     * executed first. Needs to execute after Rules Engine endpoint
     * resolution middleware insertion.
     *
     * @return Returns the sort order, defaults to 127.
     */
    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .registerMiddleware(DISABLE_HTTPS_MIDDLEWARE)
                        .build()
        );
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator
    ) {
        var serviceShape = settings.getService(model);
        goDelegator.useShapeWriter(serviceShape, writer -> {
            writer.write(generateMiddleware());
            writer.write("");
            writer.write("""
                    func $L(stack $P, o Options) error {
                        return stack.Finalize.Insert(&$L{
                            DisableHTTPS: o.EndpointOptions.DisableHTTPS,
                        }, $S, $T)
                    }
                    """,
                    MIDDLEWARE_ADDER,
                    SmithyGoTypes.Middleware.Stack,
                    MIDDLEWARE_NAME,
                    EndpointMiddlewareGenerator.MIDDLEWARE_ID,
                    SmithyGoTypes.Middleware.After);
        });
    }

    private GoWriter.Writable generateMiddleware() {
        return createFinalizeStepMiddleware(MIDDLEWARE_NAME, MiddlewareIdentifier.string(MIDDLEWARE_ID))
                .asWritable(generateBody(), generateFields());
    }

    private GoWriter.Writable generateFields() {
        return goTemplate("""
                DisableHTTPS bool
                """);
    }

    private GoWriter.Writable generateBody() {
        return goTemplate("""
                req, ok := in.Request.($P)
                if !ok {
                    return out, metadata, $T("unknown transport type %T", in.Request)
                }

                if m.DisableHTTPS && !$T(ctx) {
                    req.URL.Scheme = "http"
                }

                return next.HandleFinalize(ctx, in)
                """,
                SmithyGoTypes.Transport.Http.Request,
                GoStdlibTypes.Fmt.Errorf,
                SmithyGoTypes.Transport.Http.GetHostnameImmutable);
    }
}
