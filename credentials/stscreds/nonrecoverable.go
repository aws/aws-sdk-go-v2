package stscreds

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/credentials/nonrecoverable"
	smithy "github.com/aws/smithy-go"
)

// nonRecoverableCodes maps the STS error codes that retrying won't fix to the message shown to the customer.
var nonRecoverableCodes = map[string]string{
	"AccessDenied":            "access denied assuming the role",
	"IDPRejectedClaim":        "the identity provider rejected the web identity token",
	"InvalidIdentityToken":    "the web identity token is invalid or malformed",
	"MalformedPolicyDocument": "the session policy document is malformed",
	"PackedPolicyTooLarge":    "the session policy is too large once packed",
	"RegionDisabled":          "STS is disabled in this region for the account",
}

// classifyError wraps err in a nonrecoverable.Failure if STS returned an error that retrying won't fix.
func classifyError(err error) error {
	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		return err
	}
	if msg, ok := nonRecoverableCodes[strings.TrimSuffix(apiErr.ErrorCode(), "Exception")]; ok {
		return &nonrecoverable.Failure{Err: err, ActionableMessage: msg}
	}
	return err
}
