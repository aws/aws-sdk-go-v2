package software.amazon.smithy.aws.go.codegen.customization.s3;

import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.shapes.ShapeId;

import java.util.List;

import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * Adds validation to PutObject to ensure that content-length is derivable _somehow_ (either through the body being
 * seekable or a length being set) - which is required for the operation to function since the service doesn't support
 * chunked transfer-encoding.
 */
public class ValidatePutObjectContentLength implements GoIntegration {
    private static final ShapeId PUT_OBJECT_SHAPE_ID = ShapeId.from("com.amazonaws.s3#PutObject");

    private static final MiddlewareRegistrar MIDDLEWARE = MiddlewareRegistrar.builder()
            .resolvedFunction(buildPackageSymbol("addValidatePutObjectContentLength"))
            .build();

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return List.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate((model, service, operation) -> operation.getId().equals(PUT_OBJECT_SHAPE_ID))
                        .registerMiddleware(MIDDLEWARE)
                        .build()
        );
    }
}
