// Code generated by smithy-go-codegen DO NOT EDIT.

package types

type BatchGetPolicyErrorCode string

// Enum values for BatchGetPolicyErrorCode
const (
	BatchGetPolicyErrorCodePolicyStoreNotFound BatchGetPolicyErrorCode = "POLICY_STORE_NOT_FOUND"
	BatchGetPolicyErrorCodePolicyNotFound      BatchGetPolicyErrorCode = "POLICY_NOT_FOUND"
)

// Values returns all known values for BatchGetPolicyErrorCode. Note that this can
// be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (BatchGetPolicyErrorCode) Values() []BatchGetPolicyErrorCode {
	return []BatchGetPolicyErrorCode{
		"POLICY_STORE_NOT_FOUND",
		"POLICY_NOT_FOUND",
	}
}

type CedarVersion string

// Enum values for CedarVersion
const (
	CedarVersionCedar2 CedarVersion = "CEDAR_2"
	CedarVersionCedar4 CedarVersion = "CEDAR_4"
)

// Values returns all known values for CedarVersion. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (CedarVersion) Values() []CedarVersion {
	return []CedarVersion{
		"CEDAR_2",
		"CEDAR_4",
	}
}

type Decision string

// Enum values for Decision
const (
	DecisionAllow Decision = "ALLOW"
	DecisionDeny  Decision = "DENY"
)

// Values returns all known values for Decision. Note that this can be expanded in
// the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (Decision) Values() []Decision {
	return []Decision{
		"ALLOW",
		"DENY",
	}
}

type DeletionProtection string

// Enum values for DeletionProtection
const (
	DeletionProtectionEnabled  DeletionProtection = "ENABLED"
	DeletionProtectionDisabled DeletionProtection = "DISABLED"
)

// Values returns all known values for DeletionProtection. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (DeletionProtection) Values() []DeletionProtection {
	return []DeletionProtection{
		"ENABLED",
		"DISABLED",
	}
}

type OpenIdIssuer string

// Enum values for OpenIdIssuer
const (
	OpenIdIssuerCognito OpenIdIssuer = "COGNITO"
)

// Values returns all known values for OpenIdIssuer. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (OpenIdIssuer) Values() []OpenIdIssuer {
	return []OpenIdIssuer{
		"COGNITO",
	}
}

type PolicyEffect string

// Enum values for PolicyEffect
const (
	PolicyEffectPermit PolicyEffect = "Permit"
	PolicyEffectForbid PolicyEffect = "Forbid"
)

// Values returns all known values for PolicyEffect. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PolicyEffect) Values() []PolicyEffect {
	return []PolicyEffect{
		"Permit",
		"Forbid",
	}
}

type PolicyType string

// Enum values for PolicyType
const (
	PolicyTypeStatic         PolicyType = "STATIC"
	PolicyTypeTemplateLinked PolicyType = "TEMPLATE_LINKED"
)

// Values returns all known values for PolicyType. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PolicyType) Values() []PolicyType {
	return []PolicyType{
		"STATIC",
		"TEMPLATE_LINKED",
	}
}

type ResourceType string

// Enum values for ResourceType
const (
	ResourceTypeIdentitySource ResourceType = "IDENTITY_SOURCE"
	ResourceTypePolicyStore    ResourceType = "POLICY_STORE"
	ResourceTypePolicy         ResourceType = "POLICY"
	ResourceTypePolicyTemplate ResourceType = "POLICY_TEMPLATE"
	ResourceTypeSchema         ResourceType = "SCHEMA"
)

// Values returns all known values for ResourceType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ResourceType) Values() []ResourceType {
	return []ResourceType{
		"IDENTITY_SOURCE",
		"POLICY_STORE",
		"POLICY",
		"POLICY_TEMPLATE",
		"SCHEMA",
	}
}

type ValidationMode string

// Enum values for ValidationMode
const (
	ValidationModeOff    ValidationMode = "OFF"
	ValidationModeStrict ValidationMode = "STRICT"
)

// Values returns all known values for ValidationMode. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ValidationMode) Values() []ValidationMode {
	return []ValidationMode{
		"OFF",
		"STRICT",
	}
}
