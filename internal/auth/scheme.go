package auth

import (
	"fmt"
	smithy "github.com/aws/smithy-go"
)

// SigV4 is a constant representing
// Authentication Scheme Signature Version 4
const SigV4 = "sigv4"

// SigV4A is a constant representing
// Authentication Scheme Signature Version 4A
const SigV4A = "sigv4a"

// None is a constant representing the
// None Authentication Scheme
const None = "none"

// SupportedSchemes is a data structure
// that indicates the list of supported AWS
// authentication schemes
var SupportedSchemes = map[string]bool{
	SigV4:  true,
	SigV4A: true,
	None:   true,
}

// AuthenticationScheme is a representation of
// AWS authentication schemes
type AuthenticationScheme interface {
	isAuthenticationScheme()
	GetName() string
}

// AuthenticationSchemeV4 is a AWS SigV4 representation
type AuthenticationSchemeV4 struct {
	Name                  string
	SigningName           *string
	SigningRegion         *string
	DisableDoubleEncoding *bool
}

func (a *AuthenticationSchemeV4) isAuthenticationScheme() {}

// GetName provides the name of the AWS Authentication Scheme
func (a *AuthenticationSchemeV4) GetName() string {
	return a.Name
}

// AuthenticationSchemeV4A is a AWS SigV4A representation
type AuthenticationSchemeV4A struct {
	Name                  string
	SigningName           *string
	SigningRegionSet      []string
	DisableDoubleEncoding *bool
}

func (a *AuthenticationSchemeV4A) isAuthenticationScheme() {}

// GetName provides the name of the AWS Authentication Scheme
func (a *AuthenticationSchemeV4A) GetName() string {
	return a.Name
}

// GetAuthenticationSchemes extracts the relevant authentication scheme data
// into a custom strongly typed Go data structure.
func GetAuthenticationSchemes(p *smithy.Properties) ([]AuthenticationScheme, error) {
	var result []AuthenticationScheme
	authSchemes, ok := p.Get("authSchemes").([]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid authSchemes")
	}

	for _, scheme := range authSchemes {
		authScheme, ok := scheme.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Invalid authSchemes")
		}

		if authScheme["name"] == SigV4 {

			v4Scheme := AuthenticationSchemeV4{
				Name:                  SigV4,
				SigningName:           getSigningName(authScheme),
				SigningRegion:         getSigningRegion(authScheme),
				DisableDoubleEncoding: getDisableDoubleEncoding(authScheme),
			}
			result = append(result, AuthenticationScheme(&v4Scheme))

		}

		if authScheme["name"] == SigV4A {
			v4aScheme := AuthenticationSchemeV4A{
				Name:                  SigV4A,
				SigningName:           getSigningName(authScheme),
				SigningRegionSet:      authScheme["signingRegionSet"].([]string),
				DisableDoubleEncoding: getDisableDoubleEncoding(authScheme),
			}
			result = append(result, AuthenticationScheme(&v4aScheme))
		}

	}

	return result, nil
}

func getSigningName(authScheme map[string]interface{}) *string {
	signingName, ok := authScheme["signingName"].(string)
	if !ok || signingName == "" {
		return nil
	}
	return &signingName
}

func getSigningRegion(authScheme map[string]interface{}) *string {
	signingRegion, ok := authScheme["signingRegion"].(string)
	if !ok || signingRegion == "" {
		return nil
	}
	return &signingRegion
}

func getDisableDoubleEncoding(authScheme map[string]interface{}) *bool {
	disableDoubleEncoding, ok := authScheme["disableDoubleEncoding"].(bool)
	if !ok {
		return nil
	}
	return &disableDoubleEncoding
}
