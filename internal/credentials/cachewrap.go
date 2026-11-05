// Package credentials holds internal helpers for credential providers and the credentials cache.
package credentials

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

// inScopeSources is a map of in-scope providers.
var inScopeSources = map[aws.CredentialSource]struct{}{
	aws.CredentialSourceIMDS:                 {},
	aws.CredentialSourceHTTP:                 {},
	aws.CredentialSourceSTSAssumeRole:        {},
	aws.CredentialSourceSTSAssumeRoleWebID:   {},
	aws.CredentialSourceEnvVarsSTSWebIDToken: {},
	aws.CredentialSourceProfileSTSWebIDToken: {},
	aws.CredentialSourceSSO:                  {},
	aws.CredentialSourceProfileSSO:           {},
	aws.CredentialSourceSSOLegacy:            {},
	aws.CredentialSourceProfileSSOLegacy:     {},
	aws.CredentialSourceLogin:                {},
	aws.CredentialSourceProfileLogin:         {},
}

// inScope reports whether the provider is in inScopeSources.
func inScope(provider aws.CredentialsProvider) bool {
	asSource, ok := provider.(aws.CredentialProviderSource)
	if !ok {
		return false
	}
	sources := asSource.ProviderSources()
	if len(sources) == 0 {
		return false
	}
	_, ok = inScopeSources[sources[len(sources)-1]]
	return ok
}

// Wrap returns the provider wrapped in an aws.CredentialsCache if it's in scope and not already wrapped.
// Generated clients call it so a provider set directly on Options.Credentials still gets cached.
func Wrap(provider aws.CredentialsProvider, optFns ...func(options *aws.CredentialsCacheOptions)) aws.CredentialsProvider {
	if provider == nil {
		return nil
	}
	if _, ok := provider.(*aws.CredentialsCache); ok {
		return provider
	}
	if !inScope(provider) {
		return provider
	}
	return aws.NewCredentialsCache(provider, optFns...)
}
