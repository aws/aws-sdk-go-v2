package software.amazon.smithy.aws.go.codegen.customization;

import software.amazon.smithy.aws.go.codegen.SdkGoTypes;
import software.amazon.smithy.codegen.core.SymbolProvider;
import software.amazon.smithy.go.codegen.GoDelegator;
import software.amazon.smithy.go.codegen.GoSettings;
import software.amazon.smithy.go.codegen.GoStdlibTypes;
import software.amazon.smithy.go.codegen.GoWriter;
import software.amazon.smithy.go.codegen.integration.GoIntegration;
import software.amazon.smithy.go.codegen.SmithyGoTypes;
import software.amazon.smithy.model.Model;
import software.amazon.smithy.model.shapes.ServiceShape;
import software.amazon.smithy.rulesengine.language.syntax.Identifier;
import software.amazon.smithy.rulesengine.traits.EndpointRuleSetTrait;
import software.amazon.smithy.utils.MapUtils;

import static software.amazon.smithy.go.codegen.GoWriter.goTemplate;

public class AccountIDEndpointRouting implements GoIntegration {
    @Override
    public void renderPreEndpointResolutionHook(GoSettings settings, GoWriter writer, Model model) {
        if (!hasAccountIdEndpoints(model, settings.getService(model))) {
            return;
        }

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
        if (!hasAccountIdEndpoints(model, settings.getService(model))) {
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
                    return $errorf:T("accountID is required but not set")
                } else if ca.Credentials.AccountID == "" {
                    return $errorf:T("accountID is required but not set")
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

    public static boolean hasAccountIdEndpoints(Model model, ServiceShape service) {
        if (!service.hasTrait(EndpointRuleSetTrait.class)) {
            return false;
        }

        var rules = service.expectTrait(EndpointRuleSetTrait.class).getEndpointRuleSet();
        for (var param : rules.getParameters()) {
            if (param.getBuiltIn().orElse("").equals("AWS::Auth::AccountId")) {
                return true;
            }
        }

        return false;
    }
}
