package v4

import "strings"

// BuildCredentialScope builds the Signature Version 4 (SigV4) signing scope
func BuildCredentialScope(signingTime SigningTime, region, service string) string {
	const suffix = "aws4_request"
	t := signingTime.ShortTimeFormat()

	var sb strings.Builder
	sb.Grow(len(t) + 1 + len(region) + 1 + len(service) + 1 + len(suffix))
	sb.WriteString(t)
	sb.WriteByte('/')
	sb.WriteString(region)
	sb.WriteByte('/')
	sb.WriteString(service)
	sb.WriteByte('/')
	sb.WriteString(suffix)
	return sb.String()
}
