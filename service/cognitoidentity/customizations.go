package cognitoidentity

import "github.com/aws/aws-sdk-go-v2/aws"

func init() {
	initRequest = func(c *CognitoIdentity, r *aws.Request) {
		switch r.Operation.Name {
		case opGetOpenIdToken, opGetId, opGetCredentialsForIdentity:
			r.Handlers.Sign.Clear() // these operations are unsigned
		}
	}
}
