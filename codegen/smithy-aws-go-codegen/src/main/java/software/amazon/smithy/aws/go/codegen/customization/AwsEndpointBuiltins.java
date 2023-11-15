package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.aws.go.codegen.SdkGoTypes;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.utils.ListUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

/**
 * Registers bindings for AWS endpoint resolution builtins.
 */
public class AwsEndpointBuiltins implements GoIntegration {
    private static final GoWriter.Writable BindSdkEndpoint =
            goTemplate("options.BaseEndpoint");

    private static final GoWriter.Writable BindAwsRegion =
            goTemplate("$T($T(options.Region))", SdkGoTypes.Aws.String, SdkGoTypes.Internal.Endpoints.MapFIPSRegion);
    private static final GoWriter.Writable BindAwsUseFips =
            goTemplate("$T(options.EndpointOptions.UseFIPSEndpoint == $T)", SdkGoTypes.Aws.Bool, SdkGoTypes.Aws.FIPSEndpointStateEnabled);
    private static final GoWriter.Writable BindAwsUseDualStack =
            goTemplate("$T(options.EndpointOptions.UseDualStackEndpoint == $T)", SdkGoTypes.Aws.Bool, SdkGoTypes.Aws.DualStackEndpointStateEnabled);

    private static final GoWriter.Writable BindAwsS3ForcePathStyle =
            goTemplate("$T(options.UsePathStyle)", SdkGoTypes.Aws.Bool);
    private static final GoWriter.Writable BindAwsS3Accelerate =
            goTemplate("$T(options.UseAccelerate)", SdkGoTypes.Aws.Bool);
    private static final GoWriter.Writable BindAwsS3UseArnRegion =
            goTemplate("$T(options.UseARNRegion)", SdkGoTypes.Aws.Bool);
    private static final GoWriter.Writable BindAwsS3DisableMultiRegionAccessPoints =
            goTemplate("$T(options.DisableMultiRegionAccessPoints)", SdkGoTypes.Aws.Bool);

    @Override
    public List<RuntimeClientPlugin> getClientPlugins() {
        return ListUtils.of(RuntimeClientPlugin.builder()
                .addEndpointBuiltinBinding("SDK::Endpoint", BindSdkEndpoint)
                .addEndpointBuiltinBinding("AWS::Region", BindAwsRegion)
                .addEndpointBuiltinBinding("AWS::UseFIPS", BindAwsUseFips)
                .addEndpointBuiltinBinding("AWS::UseDualStack", BindAwsUseDualStack)
                .addEndpointBuiltinBinding("AWS::S3::ForcePathStyle", BindAwsS3ForcePathStyle)
                .addEndpointBuiltinBinding("AWS::S3::Accelerate", BindAwsS3Accelerate)
                .addEndpointBuiltinBinding("AWS::S3::UseArnRegion", BindAwsS3UseArnRegion)
                .addEndpointBuiltinBinding("AWS::S3::DisableMultiRegionAccessPoints", BindAwsS3DisableMultiRegionAccessPoints)
                .addEndpointBuiltinBinding("AWS::S3Control::UseArnRegion", BindAwsS3UseArnRegion)
                .build());
    }
}
