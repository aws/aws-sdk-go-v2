package stscreds

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/credentials/nonrecoverable"
	smithy "github.com/aws/smithy-go"
)

// classifyError wraps err in a nonrecoverable.Failure if STS returned an error that retrying won't fix.
func classifyError(err error) error {
	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		return err
	}

	code := strings.TrimSuffix(apiErr.ErrorCode(), "Exception")
	switch code {
	case "AccessDenied":
		return &nonrecoverable.Failure{
			Err:               err,
			ActionableMessage: "the role's trust policy does not allow this principal to assume it",
		}
	case "IDPRejectedClaim":
		return &nonrecoverable.Failure{
			Err:               err,
			ActionableMessage: "the identity provider rejected the web identity token",
		}
	case "InvalidIdentityToken":
		return &nonrecoverable.Failure{
			Err:               err,
			ActionableMessage: "the web identity token is invalid or malformed",
		}
	case "MalformedPolicyDocument":
		return &nonrecoverable.Failure{
			Err:               err,
			ActionableMessage: "the session policy document is malformed",
		}
	case "PackedPolicyTooLarge":
		return &nonrecoverable.Failure{
			Err:               err,
			ActionableMessage: "the session policy is too large once packed",
		}
	case "RegionDisabled":
		return &nonrecoverable.Failure{
			Err:               err,
			ActionableMessage: "STS is disabled in this region for the account",
		}
	default:
		return err
	}
}
