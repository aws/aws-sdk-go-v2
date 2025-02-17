package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.aws.go.codegen.SdkGoTypes;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.integration.RuntimeClientPlugin;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.rulesengine.traits.EndpointRuleSetTrait;
import software.amazon.smithy.utils.ListUtils;
import software.amazon.smithy.utils.MapUtils;

import java.util.List;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

/**
 * Registers bindings for AWS endpoint resolution builtins.
 */
public class AwsEndpointBuiltins implements GoIntegration {
    private static final GoWriter.Writable BindSdkEndpoint =
            goTemplate("options.BaseEndpoint");

    private static final GoWriter.Writable BindAwsRegion =
            goTemplate("bindRegion(options.Region)");
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
    private static final GoWriter.Writable BindAccountID =
            goTemplate("resolveAccountID(getIdentity(ctx), options.AccountIDEndpointMode)");
    private static final GoWriter.Writable BindAccountIDEndpointMode =
            goTemplate("$T(string(options.AccountIDEndpointMode))", SdkGoTypes.Aws.String);

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
                .addEndpointBuiltinBinding("AWS::Auth::AccountId", BindAccountID)
                .addEndpointBuiltinBinding("AWS::Auth::AccountIdEndpointMode", BindAccountIDEndpointMode)
                .build());
    }

    @Override
    public void writeAdditionalFiles(GoSettings settings, Model model, SymbolProvider symbolProvider, GoDelegator goDelegator) {
        goDelegator.useFileWriter("endpoints.go", settings.getModuleName(), builtinBindingSource());
        if (!settings.getService(model).hasTrait(EndpointRuleSetTrait.class)) {
            return;
        }
        goDelegator.useShapeWriter(settings.getService(model), goTemplate("""
        func resolveAccountID(identity $auth:T, mode $accountIDEndpointMode:T) *string {
            if mode == $aidModeDisabled:T {
                return nil
            }
            
            if ca, ok := identity.(*$credentialsAdapter:T); ok && ca.Credentials.AccountID != "" {
                return $string:T(ca.Credentials.AccountID)
            }
            
            return nil
        }
        """,
                MapUtils.of(
                        "auth", SmithyGoTypes.Auth.Identity,
                        "accountIDEndpointMode", SdkGoTypes.Aws.AccountIDEndpointMode,
                        "aidModeUnset", SdkGoTypes.Aws.AccountIDEndpointModeUnset,
                        "aidModeDisabled", SdkGoTypes.Aws.AccountIDEndpointModeDisabled,
                        "credentialsAdapter", SdkGoTypes.Internal.Auth.Smithy.CredentialsAdapter,
                        "string", SdkGoTypes.Aws.String
                )
        ));
    }

    private GoWriter.Writable builtinBindingSource() {
        return goTemplate("""
                func bindRegion(region string) *string {
                    if region == "" {
                        return nil
                    }
                    return $T($T(region))
                }
                """, SdkGoTypes.Aws.String, SdkGoTypes.Internal.Endpoints.MapFIPSRegion);
    }
}
