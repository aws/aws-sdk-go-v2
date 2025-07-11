// Code generated by smithy-go-codegen DO NOT EDIT.

package types

type AccessAdvisorUsageGranularityType string

// Enum values for AccessAdvisorUsageGranularityType
const (
	AccessAdvisorUsageGranularityTypeServiceLevel AccessAdvisorUsageGranularityType = "SERVICE_LEVEL"
	AccessAdvisorUsageGranularityTypeActionLevel  AccessAdvisorUsageGranularityType = "ACTION_LEVEL"
)

// Values returns all known values for AccessAdvisorUsageGranularityType. Note
// that this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AccessAdvisorUsageGranularityType) Values() []AccessAdvisorUsageGranularityType {
	return []AccessAdvisorUsageGranularityType{
		"SERVICE_LEVEL",
		"ACTION_LEVEL",
	}
}

type AssertionEncryptionModeType string

// Enum values for AssertionEncryptionModeType
const (
	AssertionEncryptionModeTypeRequired AssertionEncryptionModeType = "Required"
	AssertionEncryptionModeTypeAllowed  AssertionEncryptionModeType = "Allowed"
)

// Values returns all known values for AssertionEncryptionModeType. Note that this
// can be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AssertionEncryptionModeType) Values() []AssertionEncryptionModeType {
	return []AssertionEncryptionModeType{
		"Required",
		"Allowed",
	}
}

type AssignmentStatusType string

// Enum values for AssignmentStatusType
const (
	AssignmentStatusTypeAssigned   AssignmentStatusType = "Assigned"
	AssignmentStatusTypeUnassigned AssignmentStatusType = "Unassigned"
	AssignmentStatusTypeAny        AssignmentStatusType = "Any"
)

// Values returns all known values for AssignmentStatusType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (AssignmentStatusType) Values() []AssignmentStatusType {
	return []AssignmentStatusType{
		"Assigned",
		"Unassigned",
		"Any",
	}
}

type ContextKeyTypeEnum string

// Enum values for ContextKeyTypeEnum
const (
	ContextKeyTypeEnumString      ContextKeyTypeEnum = "string"
	ContextKeyTypeEnumStringList  ContextKeyTypeEnum = "stringList"
	ContextKeyTypeEnumNumeric     ContextKeyTypeEnum = "numeric"
	ContextKeyTypeEnumNumericList ContextKeyTypeEnum = "numericList"
	ContextKeyTypeEnumBoolean     ContextKeyTypeEnum = "boolean"
	ContextKeyTypeEnumBooleanList ContextKeyTypeEnum = "booleanList"
	ContextKeyTypeEnumIp          ContextKeyTypeEnum = "ip"
	ContextKeyTypeEnumIpList      ContextKeyTypeEnum = "ipList"
	ContextKeyTypeEnumBinary      ContextKeyTypeEnum = "binary"
	ContextKeyTypeEnumBinaryList  ContextKeyTypeEnum = "binaryList"
	ContextKeyTypeEnumDate        ContextKeyTypeEnum = "date"
	ContextKeyTypeEnumDateList    ContextKeyTypeEnum = "dateList"
)

// Values returns all known values for ContextKeyTypeEnum. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ContextKeyTypeEnum) Values() []ContextKeyTypeEnum {
	return []ContextKeyTypeEnum{
		"string",
		"stringList",
		"numeric",
		"numericList",
		"boolean",
		"booleanList",
		"ip",
		"ipList",
		"binary",
		"binaryList",
		"date",
		"dateList",
	}
}

type DeletionTaskStatusType string

// Enum values for DeletionTaskStatusType
const (
	DeletionTaskStatusTypeSucceeded  DeletionTaskStatusType = "SUCCEEDED"
	DeletionTaskStatusTypeInProgress DeletionTaskStatusType = "IN_PROGRESS"
	DeletionTaskStatusTypeFailed     DeletionTaskStatusType = "FAILED"
	DeletionTaskStatusTypeNotStarted DeletionTaskStatusType = "NOT_STARTED"
)

// Values returns all known values for DeletionTaskStatusType. Note that this can
// be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (DeletionTaskStatusType) Values() []DeletionTaskStatusType {
	return []DeletionTaskStatusType{
		"SUCCEEDED",
		"IN_PROGRESS",
		"FAILED",
		"NOT_STARTED",
	}
}

type EncodingType string

// Enum values for EncodingType
const (
	EncodingTypeSsh EncodingType = "SSH"
	EncodingTypePem EncodingType = "PEM"
)

// Values returns all known values for EncodingType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (EncodingType) Values() []EncodingType {
	return []EncodingType{
		"SSH",
		"PEM",
	}
}

type EntityType string

// Enum values for EntityType
const (
	EntityTypeUser               EntityType = "User"
	EntityTypeRole               EntityType = "Role"
	EntityTypeGroup              EntityType = "Group"
	EntityTypeLocalManagedPolicy EntityType = "LocalManagedPolicy"
	EntityTypeAWSManagedPolicy   EntityType = "AWSManagedPolicy"
)

// Values returns all known values for EntityType. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (EntityType) Values() []EntityType {
	return []EntityType{
		"User",
		"Role",
		"Group",
		"LocalManagedPolicy",
		"AWSManagedPolicy",
	}
}

type FeatureType string

// Enum values for FeatureType
const (
	FeatureTypeRootCredentialsManagement FeatureType = "RootCredentialsManagement"
	FeatureTypeRootSessions              FeatureType = "RootSessions"
)

// Values returns all known values for FeatureType. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (FeatureType) Values() []FeatureType {
	return []FeatureType{
		"RootCredentialsManagement",
		"RootSessions",
	}
}

type GlobalEndpointTokenVersion string

// Enum values for GlobalEndpointTokenVersion
const (
	GlobalEndpointTokenVersionV1Token GlobalEndpointTokenVersion = "v1Token"
	GlobalEndpointTokenVersionV2Token GlobalEndpointTokenVersion = "v2Token"
)

// Values returns all known values for GlobalEndpointTokenVersion. Note that this
// can be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (GlobalEndpointTokenVersion) Values() []GlobalEndpointTokenVersion {
	return []GlobalEndpointTokenVersion{
		"v1Token",
		"v2Token",
	}
}

type JobStatusType string

// Enum values for JobStatusType
const (
	JobStatusTypeInProgress JobStatusType = "IN_PROGRESS"
	JobStatusTypeCompleted  JobStatusType = "COMPLETED"
	JobStatusTypeFailed     JobStatusType = "FAILED"
)

// Values returns all known values for JobStatusType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (JobStatusType) Values() []JobStatusType {
	return []JobStatusType{
		"IN_PROGRESS",
		"COMPLETED",
		"FAILED",
	}
}

type PermissionsBoundaryAttachmentType string

// Enum values for PermissionsBoundaryAttachmentType
const (
	PermissionsBoundaryAttachmentTypePolicy PermissionsBoundaryAttachmentType = "PermissionsBoundaryPolicy"
)

// Values returns all known values for PermissionsBoundaryAttachmentType. Note
// that this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PermissionsBoundaryAttachmentType) Values() []PermissionsBoundaryAttachmentType {
	return []PermissionsBoundaryAttachmentType{
		"PermissionsBoundaryPolicy",
	}
}

type PolicyEvaluationDecisionType string

// Enum values for PolicyEvaluationDecisionType
const (
	PolicyEvaluationDecisionTypeAllowed      PolicyEvaluationDecisionType = "allowed"
	PolicyEvaluationDecisionTypeExplicitDeny PolicyEvaluationDecisionType = "explicitDeny"
	PolicyEvaluationDecisionTypeImplicitDeny PolicyEvaluationDecisionType = "implicitDeny"
)

// Values returns all known values for PolicyEvaluationDecisionType. Note that
// this can be expanded in the future, and so it is only as up to date as the
// client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PolicyEvaluationDecisionType) Values() []PolicyEvaluationDecisionType {
	return []PolicyEvaluationDecisionType{
		"allowed",
		"explicitDeny",
		"implicitDeny",
	}
}

type PolicyOwnerEntityType string

// Enum values for PolicyOwnerEntityType
const (
	PolicyOwnerEntityTypeUser  PolicyOwnerEntityType = "USER"
	PolicyOwnerEntityTypeRole  PolicyOwnerEntityType = "ROLE"
	PolicyOwnerEntityTypeGroup PolicyOwnerEntityType = "GROUP"
)

// Values returns all known values for PolicyOwnerEntityType. Note that this can
// be expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PolicyOwnerEntityType) Values() []PolicyOwnerEntityType {
	return []PolicyOwnerEntityType{
		"USER",
		"ROLE",
		"GROUP",
	}
}

type PolicyScopeType string

// Enum values for PolicyScopeType
const (
	PolicyScopeTypeAll   PolicyScopeType = "All"
	PolicyScopeTypeAws   PolicyScopeType = "AWS"
	PolicyScopeTypeLocal PolicyScopeType = "Local"
)

// Values returns all known values for PolicyScopeType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PolicyScopeType) Values() []PolicyScopeType {
	return []PolicyScopeType{
		"All",
		"AWS",
		"Local",
	}
}

type PolicySourceType string

// Enum values for PolicySourceType
const (
	PolicySourceTypeUser        PolicySourceType = "user"
	PolicySourceTypeGroup       PolicySourceType = "group"
	PolicySourceTypeRole        PolicySourceType = "role"
	PolicySourceTypeAwsManaged  PolicySourceType = "aws-managed"
	PolicySourceTypeUserManaged PolicySourceType = "user-managed"
	PolicySourceTypeResource    PolicySourceType = "resource"
	PolicySourceTypeNone        PolicySourceType = "none"
)

// Values returns all known values for PolicySourceType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PolicySourceType) Values() []PolicySourceType {
	return []PolicySourceType{
		"user",
		"group",
		"role",
		"aws-managed",
		"user-managed",
		"resource",
		"none",
	}
}

type PolicyType string

// Enum values for PolicyType
const (
	PolicyTypeInline  PolicyType = "INLINE"
	PolicyTypeManaged PolicyType = "MANAGED"
)

// Values returns all known values for PolicyType. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PolicyType) Values() []PolicyType {
	return []PolicyType{
		"INLINE",
		"MANAGED",
	}
}

type PolicyUsageType string

// Enum values for PolicyUsageType
const (
	PolicyUsageTypePermissionsPolicy   PolicyUsageType = "PermissionsPolicy"
	PolicyUsageTypePermissionsBoundary PolicyUsageType = "PermissionsBoundary"
)

// Values returns all known values for PolicyUsageType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (PolicyUsageType) Values() []PolicyUsageType {
	return []PolicyUsageType{
		"PermissionsPolicy",
		"PermissionsBoundary",
	}
}

type ReportFormatType string

// Enum values for ReportFormatType
const (
	ReportFormatTypeTextCsv ReportFormatType = "text/csv"
)

// Values returns all known values for ReportFormatType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ReportFormatType) Values() []ReportFormatType {
	return []ReportFormatType{
		"text/csv",
	}
}

type ReportStateType string

// Enum values for ReportStateType
const (
	ReportStateTypeStarted    ReportStateType = "STARTED"
	ReportStateTypeInprogress ReportStateType = "INPROGRESS"
	ReportStateTypeComplete   ReportStateType = "COMPLETE"
)

// Values returns all known values for ReportStateType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (ReportStateType) Values() []ReportStateType {
	return []ReportStateType{
		"STARTED",
		"INPROGRESS",
		"COMPLETE",
	}
}

type SortKeyType string

// Enum values for SortKeyType
const (
	SortKeyTypeServiceNamespaceAscending       SortKeyType = "SERVICE_NAMESPACE_ASCENDING"
	SortKeyTypeServiceNamespaceDescending      SortKeyType = "SERVICE_NAMESPACE_DESCENDING"
	SortKeyTypeLastAuthenticatedTimeAscending  SortKeyType = "LAST_AUTHENTICATED_TIME_ASCENDING"
	SortKeyTypeLastAuthenticatedTimeDescending SortKeyType = "LAST_AUTHENTICATED_TIME_DESCENDING"
)

// Values returns all known values for SortKeyType. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (SortKeyType) Values() []SortKeyType {
	return []SortKeyType{
		"SERVICE_NAMESPACE_ASCENDING",
		"SERVICE_NAMESPACE_DESCENDING",
		"LAST_AUTHENTICATED_TIME_ASCENDING",
		"LAST_AUTHENTICATED_TIME_DESCENDING",
	}
}

type StatusType string

// Enum values for StatusType
const (
	StatusTypeActive   StatusType = "Active"
	StatusTypeInactive StatusType = "Inactive"
	StatusTypeExpired  StatusType = "Expired"
)

// Values returns all known values for StatusType. Note that this can be expanded
// in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (StatusType) Values() []StatusType {
	return []StatusType{
		"Active",
		"Inactive",
		"Expired",
	}
}

type SummaryKeyType string

// Enum values for SummaryKeyType
const (
	SummaryKeyTypeUsers                             SummaryKeyType = "Users"
	SummaryKeyTypeUsersQuota                        SummaryKeyType = "UsersQuota"
	SummaryKeyTypeGroups                            SummaryKeyType = "Groups"
	SummaryKeyTypeGroupsQuota                       SummaryKeyType = "GroupsQuota"
	SummaryKeyTypeServerCertificates                SummaryKeyType = "ServerCertificates"
	SummaryKeyTypeServerCertificatesQuota           SummaryKeyType = "ServerCertificatesQuota"
	SummaryKeyTypeUserPolicySizeQuota               SummaryKeyType = "UserPolicySizeQuota"
	SummaryKeyTypeGroupPolicySizeQuota              SummaryKeyType = "GroupPolicySizeQuota"
	SummaryKeyTypeGroupsPerUserQuota                SummaryKeyType = "GroupsPerUserQuota"
	SummaryKeyTypeSigningCertificatesPerUserQuota   SummaryKeyType = "SigningCertificatesPerUserQuota"
	SummaryKeyTypeAccessKeysPerUserQuota            SummaryKeyType = "AccessKeysPerUserQuota"
	SummaryKeyTypeMFADevices                        SummaryKeyType = "MFADevices"
	SummaryKeyTypeMFADevicesInUse                   SummaryKeyType = "MFADevicesInUse"
	SummaryKeyTypeAccountMFAEnabled                 SummaryKeyType = "AccountMFAEnabled"
	SummaryKeyTypeAccountAccessKeysPresent          SummaryKeyType = "AccountAccessKeysPresent"
	SummaryKeyTypeAccountPasswordPresent            SummaryKeyType = "AccountPasswordPresent"
	SummaryKeyTypeAccountSigningCertificatesPresent SummaryKeyType = "AccountSigningCertificatesPresent"
	SummaryKeyTypeAttachedPoliciesPerGroupQuota     SummaryKeyType = "AttachedPoliciesPerGroupQuota"
	SummaryKeyTypeAttachedPoliciesPerRoleQuota      SummaryKeyType = "AttachedPoliciesPerRoleQuota"
	SummaryKeyTypeAttachedPoliciesPerUserQuota      SummaryKeyType = "AttachedPoliciesPerUserQuota"
	SummaryKeyTypePolicies                          SummaryKeyType = "Policies"
	SummaryKeyTypePoliciesQuota                     SummaryKeyType = "PoliciesQuota"
	SummaryKeyTypePolicySizeQuota                   SummaryKeyType = "PolicySizeQuota"
	SummaryKeyTypePolicyVersionsInUse               SummaryKeyType = "PolicyVersionsInUse"
	SummaryKeyTypePolicyVersionsInUseQuota          SummaryKeyType = "PolicyVersionsInUseQuota"
	SummaryKeyTypeVersionsPerPolicyQuota            SummaryKeyType = "VersionsPerPolicyQuota"
	SummaryKeyTypeGlobalEndpointTokenVersion        SummaryKeyType = "GlobalEndpointTokenVersion"
)

// Values returns all known values for SummaryKeyType. Note that this can be
// expanded in the future, and so it is only as up to date as the client.
//
// The ordering of this slice is not guaranteed to be stable across updates.
func (SummaryKeyType) Values() []SummaryKeyType {
	return []SummaryKeyType{
		"Users",
		"UsersQuota",
		"Groups",
		"GroupsQuota",
		"ServerCertificates",
		"ServerCertificatesQuota",
		"UserPolicySizeQuota",
		"GroupPolicySizeQuota",
		"GroupsPerUserQuota",
		"SigningCertificatesPerUserQuota",
		"AccessKeysPerUserQuota",
		"MFADevices",
		"MFADevicesInUse",
		"AccountMFAEnabled",
		"AccountAccessKeysPresent",
		"AccountPasswordPresent",
		"AccountSigningCertificatesPresent",
		"AttachedPoliciesPerGroupQuota",
		"AttachedPoliciesPerRoleQuota",
		"AttachedPoliciesPerUserQuota",
		"Policies",
		"PoliciesQuota",
		"PolicySizeQuota",
		"PolicyVersionsInUse",
		"PolicyVersionsInUseQuota",
		"VersionsPerPolicyQuota",
		"GlobalEndpointTokenVersion",
	}
}
