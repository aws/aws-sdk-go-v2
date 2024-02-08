package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.aws.go.codegen.SdkGoTypes;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.*;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.rulesengine.traits.EndpointRuleSetTrait;
import software.amazon.smithy.utils.MapUtils;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

public class AccountIDEndpointRouting implements GoIntegration {
    @Override
    public void renderPreEndpointResolutionHook(GoSettings settings, GoWriter writer, Model model) {
        writer.write("""
                if err := checkAccountID(getIdentity(ctx), m.options.AccountIDEndpointMode); err != nil {
                    return out, metadata, $T("invalid accountID set: %w", err)
                }
                """,
                GoStdlibTypes.Fmt.Errorf);
    }

    @Override
    public void writeAdditionalFiles(
            GoSettings settings,
            Model model,
            SymbolProvider symbolProvider,
            GoDelegator goDelegator
    ) {
        if (!settings.getService(model).hasTrait(EndpointRuleSetTrait.class)) {
            return;
        }
        goDelegator.useShapeWriter(settings.getService(model), goTemplate("""
        func checkAccountID(identity $auth:T, mode $accountIDEndpointMode:T) error {
            switch mode {
            case $aidModeUnset:T:
            case $aidModePreferred:T:
            case $aidModeDisabled:T:
            case $aidModeRequired:T:
                if ca, ok := identity.(*$credentialsAdapter:T); !ok {
                    return $errorf:T("the accountID is configured to be required, but the " +
                        "identity provider could not be converted to a valid credentials adapter " +
                        "and provide an accountID, should switch to a valid credentials provider " +
                        "or change to another account id endpoint mode")
                } else if ca.Credentials.AccountID == "" {
                    return $errorf:T("the required accountID could not be empty, should switch to a valid " +
                        "credentials provider or change to another account id endpoint mode")
                }
            // default check in case invalid mode is configured through request config
            default:
                return $errorf:T("invalid accountID endpoint mode %s, must be preferred/required/disabled", mode)
            }
        
            return nil
        }
        """,
        MapUtils.of(
        "auth", SmithyGoTypes.Auth.Identity,
        "accountIDEndpointMode", SdkGoTypes.Aws.AccountIDEndpointMode,
        "credentialsAdapter", SdkGoTypes.Internal.Auth.Smithy.CredentialsAdapter,
        "aidModePreferred", SdkGoTypes.Aws.AccountIDEndpointModePreferred,
        "aidModeRequired", SdkGoTypes.Aws.AccountIDEndpointModeRequired,
        "aidModeUnset", SdkGoTypes.Aws.AccountIDEndpointModeUnset,
        "aidModeDisabled", SdkGoTypes.Aws.AccountIDEndpointModeDisabled,
        "errorf", GoStdlibTypes.Fmt.Errorf
        )
        ));
    }
}
