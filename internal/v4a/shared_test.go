package v4a

import (
	"bytes"
	"context"
	"crypto/ecdsa"
)

var stubCredentials = stubCredentialsProviderFunc(func(ctx context.Context) (Credentials, error) {
	stubKey, err := ecdsa.GenerateKey(p256, bytes.NewReader(bytes.Repeat([]byte{1}, 40)))
	if err != nil {
		return Credentials{}, err
	}
	return Credentials{
		Context:    "STUB",
		PrivateKey: stubKey,
	}, nil
})
