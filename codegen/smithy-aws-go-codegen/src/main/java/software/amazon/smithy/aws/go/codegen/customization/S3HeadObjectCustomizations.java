package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.knowledge.TopDownIndex;
import software.amazon.smithy.model.shapes.OperationShape;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.model.shapes.ShapeId;

/**
 * This integration removes incorrectly modeled S3 `NoSuchKey` error for HeadObject operation.
 * Related to aws/aws-sdk-go#1208
 */
public class S3HeadObjectCustomizations implements GoIntegration {
    @Override
    public byte getOrder() {
        return 127;
    }

    @Override
    public Model preprocessModel(
            Model model, GoSettings settings
    ) {
        ServiceShape service = model.expectShape(settings.getService(), ServiceShape.class);
        if (!S3ModelUtils.isServiceS3(model, service)) {
            return model;
        }

        for (OperationShape operationShape: TopDownIndex.of(model).getContainedOperations(service)) {
            if (operationShape.getId().getName(service).equalsIgnoreCase("HeadObject")) {
                Model.Builder modelBuilder = model.toBuilder();
                OperationShape.Builder builder = operationShape.toBuilder();
                // clear all errors associated with 'HeadObject' operation aws/aws-sdk-go#1208
                builder.clearErrors();

                // remove old operation shape and add the updated shape to model
                return modelBuilder.removeShape(operationShape.getId()).addShape(builder.build()).build();
            }
        }
        return model;
    }

    // returns true if service is s3
}
