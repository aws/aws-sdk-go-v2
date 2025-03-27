package software.amazon.smithy.aws.go.codegen.customization.service.s3;

import software.amazon.smithy.aws.go.codegen.SdkGoTypes;
import software.amazon.smithy.aws.traits.HttpChecksumTrait;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.SymbolUtils.buildPackageSymbol;

/**
 * Defaults the checksum to CRC32 for requests against the Express storage class when a checksum is required but not set
 * (we'd normally default to MD5). Defaults to CRC32 for transfer manager operations as well.
 */
public class ExpressDefaultChecksum implements GoIntegration {
    private static final MiddlewareRegistrar SET_ALGORITHM_MIDDLEWARE = MiddlewareRegistrar.builder()
            .resolvedFunction(SdkGoTypes.ServiceCustomizations.S3.AddExpressDefaultChecksumMiddleware)
            .build();

    private static final MiddlewareRegistrar SET_MPU_ALGORITHM_MIDDLEWARE = MiddlewareRegistrar.builder()
            .resolvedFunction(buildPackageSymbol("addSetCreateMPUChecksumAlgorithm"))
            .build();

    private static boolean hasRequiredChecksum(Model model, ServiceShape service, OperationShape operation) {
        return operation.hasTrait(HttpChecksumTrait.class)
                && operation.expectTrait(HttpChecksumTrait.class).isRequestChecksumRequired();
    }

    private static boolean isCreateMultipartUpload(Model model, ServiceShape service, OperationShape operation) {
        return operation.getId().getName().equals("CreateMultipartUpload");
    }

    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3ModelUtils::isServiceS3)
                        .operationPredicate(ExpressDefaultChecksum::hasRequiredChecksum)
                        .registerMiddleware(SET_ALGORITHM_MIDDLEWARE)
                        .build(),
                RuntimeClientPlugin.builder()
                        .servicePredicate(S3ModelUtils::isServiceS3)
                        .operationPredicate(ExpressDefaultChecksum::isCreateMultipartUpload)
                        .registerMiddleware(SET_MPU_ALGORITHM_MIDDLEWARE)
                        .build()
        );
    }

    @Override
    public void renderPostEndpointResolutionHook(GoSettings settings, GoWriter writer, Model model) {
        if (!S3ModelUtils.isServiceS3(model, settings.getService(model))) return;

        writer.write("""
                backend := $T(&endpt.Properties)
                ctx = $T(ctx, backend)
                """,
                SdkGoTypes.ServiceCustomizations.S3.GetPropertiesBackend,
                SdkGoTypes.Internal.Context.SetS3Backend);
    }
}
