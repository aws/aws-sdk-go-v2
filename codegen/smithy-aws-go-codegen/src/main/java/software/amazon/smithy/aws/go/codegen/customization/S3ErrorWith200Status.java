package software.amazon.smithy.aws.go.codegen.customization;

import java.util.List;
import java.util.Set;
import software.amazon.smithy.go.codegen.SymbolUtils;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.MiddlewareRegistrar;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.SetUtils;

/**
 * Adds middleware to handle S3 response errors with 200 ok status code.
 */
public class S3ErrorWith200Status implements GoIntegration {
    private static String ADD_ERROR_HANDLER_INTERNAL = "HandleResponseErrorWith200Status";

    // list of operations for which this customization is valid.
    private static Set<String> customizedOperations = SetUtils.of(
            "CopyObject", "UploadPartCopy", "CompleteMultipartUpload");


    @Override
    public byte getOrder() {
        // The associated customization ordering is relative to operation deserializers
        // and thus the integration should be added at the end.
        return 127;
    }

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(
                RuntimeClientPlugin.builder()
                        .operationPredicate(S3ErrorWith200Status::supports200Error)
                        .registerMiddleware(MiddlewareRegistrar.builder()
                                .resolvedFunction(SymbolUtils.createValueSymbolBuilder(ADD_ERROR_HANDLER_INTERNAL,
                                        AwsCustomGoDependency.S3_CUSTOMIZATION).build())
                                .build())
                        .build()
        );
    }

    // returns true if the operation supports error response with 200 ok status code
    private static boolean supports200Error(Model model, ServiceShape service, OperationShape operation){
        if (!isS3Service(model, service)) {
            return false;
        }

       return customizedOperations.contains(operation.getId().getName(service));
    }

    // returns true if service is s3
    private static boolean isS3Service(Model model, ServiceShape service) {
        return S3ModelUtils.isServiceS3(model, service);
    }
}
