package accountid

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	accountidmode "github.com/aws/aws-sdk-go-v2/aws/accountid/mode"
	"github.com/aws/aws-sdk-go-v2/internal/auth/smithy"
	"github.com/aws/smithy-go/auth"
)

func AccountID(identity auth.Identity, mode accountidmode.AIDMode) *string {
	if ca, ok := identity.(*smithy.CredentialsAdapter); ok && (mode == accountidmode.Preferred || mode == accountidmode.Required) {
		return aws.String(ca.Credentials.AccountID)
	}

	return nil
}

func CheckAccountID(identity auth.Identity, mode accountidmode.AIDMode) error {
	switch mode {
	case "":
	case accountidmode.Preferred:
	case accountidmode.Disabled:
	case accountidmode.Required:
		if ca, ok := identity.(*smithy.CredentialsAdapter); !ok {
			return fmt.Errorf("the accountID is configured to be required, but the " +
				"identity provider could not be converted to a valid credentials adapter " +
				"and provide an accountID, should try to configure a valid credentials provider")
		} else if ca.Credentials.AccountID == "" {
			return fmt.Errorf("the required accountID could not be empty")
		}
	// default check in case invalid mode is configured through request config
	default:
		return fmt.Errorf("invalid accountID endpoint mode %s, must be preferred/required/disabled", mode)
	}

	return nil
}
