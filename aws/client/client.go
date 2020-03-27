package client

import (
	"context"

	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/transport/http"
)

var _ http.Header
var _ middleware.MiddlewareHandler

type signingMetadataKey struct{}

type SigningMetadata struct {
	SigningName   string
	SigningRegion string
}

func GetSigningMetadata(ctx context.Context) SigningMetadata {
	value := ctx.Value(signingMetadataKey{})

	sm, ok := value.(SigningMetadata)
	if !ok {
		return SigningMetadata{}
	}

	return sm
}

func SetSigningMetadata(ctx context.Context, metadata SigningMetadata) context.Context {
	ctx = context.WithValue(ctx, signingMetadataKey{}, metadata)
	return ctx
}
