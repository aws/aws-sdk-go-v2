// Code generated by smithy-go-codegen DO NOT EDIT.

package types

import (
	smithydocument "github.com/aws/smithy-go/document"
	"time"
)

// A provider representing an Amazon Cognito user pool and its client ID.
type CognitoIdentityProvider struct {

	// The client ID for the Amazon Cognito user pool.
	ClientId *string

	// The provider name for an Amazon Cognito user pool. For example,
	// cognito-idp.us-east-1.amazonaws.com/us-east-1_123456789 .
	ProviderName *string

	// TRUE if server-side token validation is enabled for the identity provider’s
	// token.
	//
	// Once you set ServerSideTokenCheck to TRUE for an identity pool, that identity
	// pool will check with the integrated user pools to make sure that the user has
	// not been globally signed out or deleted before the identity pool provides an
	// OIDC token or Amazon Web Services credentials for the user.
	//
	// If the user is signed out or deleted, the identity pool will return a 400 Not
	// Authorized error.
	ServerSideTokenCheck *bool

	noSmithyDocumentSerde
}

// Credentials for the provided identity ID.
type Credentials struct {

	// The Access Key portion of the credentials.
	AccessKeyId *string

	// The date at which these credentials will expire.
	Expiration *time.Time

	// The Secret Access Key portion of the credentials
	SecretKey *string

	// The Session Token portion of the credentials
	SessionToken *string

	noSmithyDocumentSerde
}

// A description of the identity.
type IdentityDescription struct {

	// Date on which the identity was created.
	CreationDate *time.Time

	// A unique identifier in the format REGION:GUID.
	IdentityId *string

	// Date on which the identity was last modified.
	LastModifiedDate *time.Time

	// The provider names.
	Logins []string

	noSmithyDocumentSerde
}

// A description of the identity pool.
type IdentityPoolShortDescription struct {

	// An identity pool ID in the format REGION:GUID.
	IdentityPoolId *string

	// A string that you provide.
	IdentityPoolName *string

	noSmithyDocumentSerde
}

// A rule that maps a claim name, a claim value, and a match type to a role ARN.
type MappingRule struct {

	// The claim name that must be present in the token, for example, "isAdmin" or
	// "paid".
	//
	// This member is required.
	Claim *string

	// The match condition that specifies how closely the claim value in the IdP token
	// must match Value .
	//
	// This member is required.
	MatchType MappingRuleMatchType

	// The role ARN.
	//
	// This member is required.
	RoleARN *string

	// A brief string that the claim must match, for example, "paid" or "yes".
	//
	// This member is required.
	Value *string

	noSmithyDocumentSerde
}

// A role mapping.
type RoleMapping struct {

	// The role mapping type. Token will use cognito:roles and cognito:preferred_role
	// claims from the Cognito identity provider token to map groups to roles. Rules
	// will attempt to match claims from the token to map to a role.
	//
	// This member is required.
	Type RoleMappingType

	// If you specify Token or Rules as the Type , AmbiguousRoleResolution is required.
	//
	// Specifies the action to be taken if either no rules match the claim value for
	// the Rules type, or there is no cognito:preferred_role claim and there are
	// multiple cognito:roles matches for the Token type.
	AmbiguousRoleResolution AmbiguousRoleResolutionType

	// The rules to be used for mapping users to roles.
	//
	// If you specify Rules as the role mapping type, RulesConfiguration is required.
	RulesConfiguration *RulesConfigurationType

	noSmithyDocumentSerde
}

// A container for rules.
type RulesConfigurationType struct {

	// An array of rules. You can specify up to 25 rules per identity provider.
	//
	// Rules are evaluated in order. The first one to match specifies the role.
	//
	// This member is required.
	Rules []MappingRule

	noSmithyDocumentSerde
}

// An array of UnprocessedIdentityId objects, each of which contains an ErrorCode
// and IdentityId.
type UnprocessedIdentityId struct {

	// The error code indicating the type of error that occurred.
	ErrorCode ErrorCode

	// A unique identifier in the format REGION:GUID.
	IdentityId *string

	noSmithyDocumentSerde
}

type noSmithyDocumentSerde = smithydocument.NoSerde
