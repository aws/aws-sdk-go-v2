package smithy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/internal/v4a"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/auth"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// V4ACredentialsAdapter adapts v4a.Credentials to auth.Identity.
type V4ACredentialsAdapter struct {
	creds v4a.Credentials
}

var _ (smithyhttp.Signer) = (*V4ASignerAdapter)(nil)

// Expiration returns the time of expiration for the credentials.
func (v *V4ACredentialsAdapter) Expiration() time.Time {
	return v.creds.Expires
}

// V4SignerAdapter adapts v4a.HTTPSigner to smithy http.Signer.
type V4ASignerAdapter struct {
	signer v4a.HTTPSigner
}

var _ (smithyhttp.Signer) = (*V4ASignerAdapter)(nil)

// SignRequest signs the request with the provided identity.
func (v *V4ASignerAdapter) SignRequest(ctx context.Context, r *http.Request, identity auth.Identity, props *smithy.Properties) error {
	ca, ok := identity.(*V4ACredentialsAdapter)
	if !ok {
		return fmt.Errorf("unexpected identity type: %T", identity)
	}

	name, ok := smithyhttp.GetSigV4ASigningName(props)
	if !ok {
		return fmt.Errorf("sigv4a signing name is required")
	}

	regions, ok := smithyhttp.GetSigV4ASigningRegions(props)
	if !ok {
		return fmt.Errorf("sigv4a signing region set is required")
	}

	hash := v4.GetPayloadHash(ctx)
	err := v.signer.SignHTTP(ctx, ca.creds, r, hash, name, regions, sdk.NowTime())
	if err != nil {
		return fmt.Errorf("sign http: %v", err)
	}

	return nil
}
