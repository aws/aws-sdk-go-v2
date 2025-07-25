// Code generated by smithy-go-codegen DO NOT EDIT.

package types_test

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/datazone/types"
)

func ExampleActionParameters_outputUsage() {
	var union types.ActionParameters
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.ActionParametersMemberAwsConsoleLink:
		_ = v.Value // Value is types.AwsConsoleLinkParameters

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.AwsConsoleLinkParameters

func ExampleAssetFilterConfiguration_outputUsage() {
	var union types.AssetFilterConfiguration
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.AssetFilterConfigurationMemberColumnConfiguration:
		_ = v.Value // Value is types.ColumnFilterConfiguration

	case *types.AssetFilterConfigurationMemberRowConfiguration:
		_ = v.Value // Value is types.RowFilterConfiguration

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.ColumnFilterConfiguration
var _ *types.RowFilterConfiguration

func ExampleAwsAccount_outputUsage() {
	var union types.AwsAccount
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.AwsAccountMemberAwsAccountId:
		_ = v.Value // Value is string

	case *types.AwsAccountMemberAwsAccountIdPath:
		_ = v.Value // Value is string

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *string
var _ *string

func ExampleConnectionPropertiesInput_outputUsage() {
	var union types.ConnectionPropertiesInput
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.ConnectionPropertiesInputMemberAthenaProperties:
		_ = v.Value // Value is types.AthenaPropertiesInput

	case *types.ConnectionPropertiesInputMemberGlueProperties:
		_ = v.Value // Value is types.GluePropertiesInput

	case *types.ConnectionPropertiesInputMemberHyperPodProperties:
		_ = v.Value // Value is types.HyperPodPropertiesInput

	case *types.ConnectionPropertiesInputMemberIamProperties:
		_ = v.Value // Value is types.IamPropertiesInput

	case *types.ConnectionPropertiesInputMemberRedshiftProperties:
		_ = v.Value // Value is types.RedshiftPropertiesInput

	case *types.ConnectionPropertiesInputMemberS3Properties:
		_ = v.Value // Value is types.S3PropertiesInput

	case *types.ConnectionPropertiesInputMemberSparkEmrProperties:
		_ = v.Value // Value is types.SparkEmrPropertiesInput

	case *types.ConnectionPropertiesInputMemberSparkGlueProperties:
		_ = v.Value // Value is types.SparkGluePropertiesInput

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.SparkEmrPropertiesInput
var _ *types.GluePropertiesInput
var _ *types.S3PropertiesInput
var _ *types.AthenaPropertiesInput
var _ *types.IamPropertiesInput
var _ *types.SparkGluePropertiesInput
var _ *types.HyperPodPropertiesInput
var _ *types.RedshiftPropertiesInput

func ExampleConnectionPropertiesOutput_outputUsage() {
	var union types.ConnectionPropertiesOutput
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.ConnectionPropertiesOutputMemberAthenaProperties:
		_ = v.Value // Value is types.AthenaPropertiesOutput

	case *types.ConnectionPropertiesOutputMemberGlueProperties:
		_ = v.Value // Value is types.GluePropertiesOutput

	case *types.ConnectionPropertiesOutputMemberHyperPodProperties:
		_ = v.Value // Value is types.HyperPodPropertiesOutput

	case *types.ConnectionPropertiesOutputMemberIamProperties:
		_ = v.Value // Value is types.IamPropertiesOutput

	case *types.ConnectionPropertiesOutputMemberRedshiftProperties:
		_ = v.Value // Value is types.RedshiftPropertiesOutput

	case *types.ConnectionPropertiesOutputMemberS3Properties:
		_ = v.Value // Value is types.S3PropertiesOutput

	case *types.ConnectionPropertiesOutputMemberSparkEmrProperties:
		_ = v.Value // Value is types.SparkEmrPropertiesOutput

	case *types.ConnectionPropertiesOutputMemberSparkGlueProperties:
		_ = v.Value // Value is types.SparkGluePropertiesOutput

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.S3PropertiesOutput
var _ *types.AthenaPropertiesOutput
var _ *types.SparkGluePropertiesOutput
var _ *types.IamPropertiesOutput
var _ *types.RedshiftPropertiesOutput
var _ *types.HyperPodPropertiesOutput
var _ *types.GluePropertiesOutput
var _ *types.SparkEmrPropertiesOutput

func ExampleConnectionPropertiesPatch_outputUsage() {
	var union types.ConnectionPropertiesPatch
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.ConnectionPropertiesPatchMemberAthenaProperties:
		_ = v.Value // Value is types.AthenaPropertiesPatch

	case *types.ConnectionPropertiesPatchMemberGlueProperties:
		_ = v.Value // Value is types.GluePropertiesPatch

	case *types.ConnectionPropertiesPatchMemberIamProperties:
		_ = v.Value // Value is types.IamPropertiesPatch

	case *types.ConnectionPropertiesPatchMemberRedshiftProperties:
		_ = v.Value // Value is types.RedshiftPropertiesPatch

	case *types.ConnectionPropertiesPatchMemberS3Properties:
		_ = v.Value // Value is types.S3PropertiesPatch

	case *types.ConnectionPropertiesPatchMemberSparkEmrProperties:
		_ = v.Value // Value is types.SparkEmrPropertiesPatch

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.SparkEmrPropertiesPatch
var _ *types.IamPropertiesPatch
var _ *types.RedshiftPropertiesPatch
var _ *types.AthenaPropertiesPatch
var _ *types.GluePropertiesPatch
var _ *types.S3PropertiesPatch

func ExampleDataSourceConfigurationInput_outputUsage() {
	var union types.DataSourceConfigurationInput
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.DataSourceConfigurationInputMemberGlueRunConfiguration:
		_ = v.Value // Value is types.GlueRunConfigurationInput

	case *types.DataSourceConfigurationInputMemberRedshiftRunConfiguration:
		_ = v.Value // Value is types.RedshiftRunConfigurationInput

	case *types.DataSourceConfigurationInputMemberSageMakerRunConfiguration:
		_ = v.Value // Value is types.SageMakerRunConfigurationInput

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.SageMakerRunConfigurationInput
var _ *types.RedshiftRunConfigurationInput
var _ *types.GlueRunConfigurationInput

func ExampleDataSourceConfigurationOutput_outputUsage() {
	var union types.DataSourceConfigurationOutput
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.DataSourceConfigurationOutputMemberGlueRunConfiguration:
		_ = v.Value // Value is types.GlueRunConfigurationOutput

	case *types.DataSourceConfigurationOutputMemberRedshiftRunConfiguration:
		_ = v.Value // Value is types.RedshiftRunConfigurationOutput

	case *types.DataSourceConfigurationOutputMemberSageMakerRunConfiguration:
		_ = v.Value // Value is types.SageMakerRunConfigurationOutput

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.GlueRunConfigurationOutput
var _ *types.SageMakerRunConfigurationOutput
var _ *types.RedshiftRunConfigurationOutput

func ExampleDomainUnitGrantFilter_outputUsage() {
	var union types.DomainUnitGrantFilter
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.DomainUnitGrantFilterMemberAllDomainUnitsGrantFilter:
		_ = v.Value // Value is types.AllDomainUnitsGrantFilter

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.AllDomainUnitsGrantFilter

func ExampleDomainUnitOwnerProperties_outputUsage() {
	var union types.DomainUnitOwnerProperties
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.DomainUnitOwnerPropertiesMemberGroup:
		_ = v.Value // Value is types.DomainUnitGroupProperties

	case *types.DomainUnitOwnerPropertiesMemberUser:
		_ = v.Value // Value is types.DomainUnitUserProperties

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.DomainUnitGroupProperties
var _ *types.DomainUnitUserProperties

func ExampleEventSummary_outputUsage() {
	var union types.EventSummary
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.EventSummaryMemberOpenLineageRunEventSummary:
		_ = v.Value // Value is types.OpenLineageRunEventSummary

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.OpenLineageRunEventSummary

func ExampleFilterClause_outputUsage() {
	var union types.FilterClause
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.FilterClauseMemberAnd:
		_ = v.Value // Value is []types.FilterClause

	case *types.FilterClauseMemberFilter:
		_ = v.Value // Value is types.Filter

	case *types.FilterClauseMemberOr:
		_ = v.Value // Value is []types.FilterClause

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.Filter
var _ []types.FilterClause

func ExampleGrantedEntity_outputUsage() {
	var union types.GrantedEntity
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.GrantedEntityMemberListing:
		_ = v.Value // Value is types.ListingRevision

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.ListingRevision

func ExampleGrantedEntityInput_outputUsage() {
	var union types.GrantedEntityInput
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.GrantedEntityInputMemberListing:
		_ = v.Value // Value is types.ListingRevisionInput

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.ListingRevisionInput

func ExampleGroupPolicyGrantPrincipal_outputUsage() {
	var union types.GroupPolicyGrantPrincipal
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.GroupPolicyGrantPrincipalMemberGroupIdentifier:
		_ = v.Value // Value is string

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *string

func ExampleJobRunDetails_outputUsage() {
	var union types.JobRunDetails
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.JobRunDetailsMemberLineageRunDetails:
		_ = v.Value // Value is types.LineageRunDetails

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.LineageRunDetails

func ExampleListingItem_outputUsage() {
	var union types.ListingItem
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.ListingItemMemberAssetListing:
		_ = v.Value // Value is types.AssetListing

	case *types.ListingItemMemberDataProductListing:
		_ = v.Value // Value is types.DataProductListing

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.AssetListing
var _ *types.DataProductListing

func ExampleMatchRationaleItem_outputUsage() {
	var union types.MatchRationaleItem
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.MatchRationaleItemMemberTextMatches:
		_ = v.Value // Value is []types.TextMatchItem

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ []types.TextMatchItem

func ExampleMember_outputUsage() {
	var union types.Member
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.MemberMemberGroupIdentifier:
		_ = v.Value // Value is string

	case *types.MemberMemberUserIdentifier:
		_ = v.Value // Value is string

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *string

func ExampleMemberDetails_outputUsage() {
	var union types.MemberDetails
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.MemberDetailsMemberGroup:
		_ = v.Value // Value is types.GroupDetails

	case *types.MemberDetailsMemberUser:
		_ = v.Value // Value is types.UserDetails

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.GroupDetails
var _ *types.UserDetails

func ExampleModel_outputUsage() {
	var union types.Model
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.ModelMemberSmithy:
		_ = v.Value // Value is string

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *string

func ExampleOwnerProperties_outputUsage() {
	var union types.OwnerProperties
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.OwnerPropertiesMemberGroup:
		_ = v.Value // Value is types.OwnerGroupProperties

	case *types.OwnerPropertiesMemberUser:
		_ = v.Value // Value is types.OwnerUserProperties

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.OwnerGroupProperties
var _ *types.OwnerUserProperties

func ExampleOwnerPropertiesOutput_outputUsage() {
	var union types.OwnerPropertiesOutput
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.OwnerPropertiesOutputMemberGroup:
		_ = v.Value // Value is types.OwnerGroupPropertiesOutput

	case *types.OwnerPropertiesOutputMemberUser:
		_ = v.Value // Value is types.OwnerUserPropertiesOutput

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.OwnerGroupPropertiesOutput
var _ *types.OwnerUserPropertiesOutput

func ExamplePolicyGrantDetail_outputUsage() {
	var union types.PolicyGrantDetail
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.PolicyGrantDetailMemberAddToProjectMemberPool:
		_ = v.Value // Value is types.AddToProjectMemberPoolPolicyGrantDetail

	case *types.PolicyGrantDetailMemberCreateAssetType:
		_ = v.Value // Value is types.CreateAssetTypePolicyGrantDetail

	case *types.PolicyGrantDetailMemberCreateDomainUnit:
		_ = v.Value // Value is types.CreateDomainUnitPolicyGrantDetail

	case *types.PolicyGrantDetailMemberCreateEnvironment:
		_ = v.Value // Value is types.Unit

	case *types.PolicyGrantDetailMemberCreateEnvironmentFromBlueprint:
		_ = v.Value // Value is types.Unit

	case *types.PolicyGrantDetailMemberCreateEnvironmentProfile:
		_ = v.Value // Value is types.CreateEnvironmentProfilePolicyGrantDetail

	case *types.PolicyGrantDetailMemberCreateFormType:
		_ = v.Value // Value is types.CreateFormTypePolicyGrantDetail

	case *types.PolicyGrantDetailMemberCreateGlossary:
		_ = v.Value // Value is types.CreateGlossaryPolicyGrantDetail

	case *types.PolicyGrantDetailMemberCreateProject:
		_ = v.Value // Value is types.CreateProjectPolicyGrantDetail

	case *types.PolicyGrantDetailMemberCreateProjectFromProjectProfile:
		_ = v.Value // Value is types.CreateProjectFromProjectProfilePolicyGrantDetail

	case *types.PolicyGrantDetailMemberDelegateCreateEnvironmentProfile:
		_ = v.Value // Value is types.Unit

	case *types.PolicyGrantDetailMemberOverrideDomainUnitOwners:
		_ = v.Value // Value is types.OverrideDomainUnitOwnersPolicyGrantDetail

	case *types.PolicyGrantDetailMemberOverrideProjectOwners:
		_ = v.Value // Value is types.OverrideProjectOwnersPolicyGrantDetail

	case *types.PolicyGrantDetailMemberUseAssetType:
		_ = v.Value // Value is types.UseAssetTypePolicyGrantDetail

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.CreateProjectFromProjectProfilePolicyGrantDetail
var _ *types.UseAssetTypePolicyGrantDetail
var _ *types.CreateDomainUnitPolicyGrantDetail
var _ *types.OverrideProjectOwnersPolicyGrantDetail
var _ *types.CreateEnvironmentProfilePolicyGrantDetail
var _ *types.CreateGlossaryPolicyGrantDetail
var _ *types.AddToProjectMemberPoolPolicyGrantDetail
var _ *types.CreateProjectPolicyGrantDetail
var _ *types.OverrideDomainUnitOwnersPolicyGrantDetail
var _ *types.CreateAssetTypePolicyGrantDetail
var _ *types.Unit
var _ *types.CreateFormTypePolicyGrantDetail

func ExamplePolicyGrantPrincipal_outputUsage() {
	var union types.PolicyGrantPrincipal
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.PolicyGrantPrincipalMemberDomainUnit:
		_ = v.Value // Value is types.DomainUnitPolicyGrantPrincipal

	case *types.PolicyGrantPrincipalMemberGroup:
		_ = v.Value // Value is types.GroupPolicyGrantPrincipal

	case *types.PolicyGrantPrincipalMemberProject:
		_ = v.Value // Value is types.ProjectPolicyGrantPrincipal

	case *types.PolicyGrantPrincipalMemberUser:
		_ = v.Value // Value is types.UserPolicyGrantPrincipal

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.DomainUnitPolicyGrantPrincipal
var _ types.GroupPolicyGrantPrincipal
var _ types.UserPolicyGrantPrincipal
var _ *types.ProjectPolicyGrantPrincipal

func ExampleProjectGrantFilter_outputUsage() {
	var union types.ProjectGrantFilter
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.ProjectGrantFilterMemberDomainUnitFilter:
		_ = v.Value // Value is types.DomainUnitFilterForProject

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.DomainUnitFilterForProject

func ExampleProvisioningConfiguration_outputUsage() {
	var union types.ProvisioningConfiguration
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.ProvisioningConfigurationMemberLakeFormationConfiguration:
		_ = v.Value // Value is types.LakeFormationConfiguration

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.LakeFormationConfiguration

func ExampleProvisioningProperties_outputUsage() {
	var union types.ProvisioningProperties
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.ProvisioningPropertiesMemberCloudFormation:
		_ = v.Value // Value is types.CloudFormationProperties

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.CloudFormationProperties

func ExampleRedshiftCredentials_outputUsage() {
	var union types.RedshiftCredentials
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.RedshiftCredentialsMemberSecretArn:
		_ = v.Value // Value is string

	case *types.RedshiftCredentialsMemberUsernamePassword:
		_ = v.Value // Value is types.UsernamePassword

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *string
var _ *types.UsernamePassword

func ExampleRedshiftStorage_outputUsage() {
	var union types.RedshiftStorage
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.RedshiftStorageMemberRedshiftClusterSource:
		_ = v.Value // Value is types.RedshiftClusterStorage

	case *types.RedshiftStorageMemberRedshiftServerlessSource:
		_ = v.Value // Value is types.RedshiftServerlessStorage

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.RedshiftClusterStorage
var _ *types.RedshiftServerlessStorage

func ExampleRedshiftStorageProperties_outputUsage() {
	var union types.RedshiftStorageProperties
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.RedshiftStoragePropertiesMemberClusterName:
		_ = v.Value // Value is string

	case *types.RedshiftStoragePropertiesMemberWorkgroupName:
		_ = v.Value // Value is string

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *string

func ExampleRegion_outputUsage() {
	var union types.Region
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.RegionMemberRegionName:
		_ = v.Value // Value is string

	case *types.RegionMemberRegionNamePath:
		_ = v.Value // Value is string

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *string
var _ *string

func ExampleRowFilter_outputUsage() {
	var union types.RowFilter
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.RowFilterMemberAnd:
		_ = v.Value // Value is []types.RowFilter

	case *types.RowFilterMemberExpression:
		_ = v.Value // Value is types.RowFilterExpression

	case *types.RowFilterMemberOr:
		_ = v.Value // Value is []types.RowFilter

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ types.RowFilterExpression
var _ []types.RowFilter

func ExampleRowFilterExpression_outputUsage() {
	var union types.RowFilterExpression
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.RowFilterExpressionMemberEqualTo:
		_ = v.Value // Value is types.EqualToExpression

	case *types.RowFilterExpressionMemberGreaterThan:
		_ = v.Value // Value is types.GreaterThanExpression

	case *types.RowFilterExpressionMemberGreaterThanOrEqualTo:
		_ = v.Value // Value is types.GreaterThanOrEqualToExpression

	case *types.RowFilterExpressionMemberIn:
		_ = v.Value // Value is types.InExpression

	case *types.RowFilterExpressionMemberIsNotNull:
		_ = v.Value // Value is types.IsNotNullExpression

	case *types.RowFilterExpressionMemberIsNull:
		_ = v.Value // Value is types.IsNullExpression

	case *types.RowFilterExpressionMemberLessThan:
		_ = v.Value // Value is types.LessThanExpression

	case *types.RowFilterExpressionMemberLessThanOrEqualTo:
		_ = v.Value // Value is types.LessThanOrEqualToExpression

	case *types.RowFilterExpressionMemberLike:
		_ = v.Value // Value is types.LikeExpression

	case *types.RowFilterExpressionMemberNotEqualTo:
		_ = v.Value // Value is types.NotEqualToExpression

	case *types.RowFilterExpressionMemberNotIn:
		_ = v.Value // Value is types.NotInExpression

	case *types.RowFilterExpressionMemberNotLike:
		_ = v.Value // Value is types.NotLikeExpression

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.NotLikeExpression
var _ *types.GreaterThanExpression
var _ *types.LessThanExpression
var _ *types.IsNotNullExpression
var _ *types.NotEqualToExpression
var _ *types.GreaterThanOrEqualToExpression
var _ *types.IsNullExpression
var _ *types.LessThanOrEqualToExpression
var _ *types.LikeExpression
var _ *types.NotInExpression
var _ *types.InExpression
var _ *types.EqualToExpression

func ExampleRuleDetail_outputUsage() {
	var union types.RuleDetail
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.RuleDetailMemberMetadataFormEnforcementDetail:
		_ = v.Value // Value is types.MetadataFormEnforcementDetail

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.MetadataFormEnforcementDetail

func ExampleRuleTarget_outputUsage() {
	var union types.RuleTarget
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.RuleTargetMemberDomainUnitTarget:
		_ = v.Value // Value is types.DomainUnitTarget

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.DomainUnitTarget

func ExampleSearchInventoryResultItem_outputUsage() {
	var union types.SearchInventoryResultItem
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.SearchInventoryResultItemMemberAssetItem:
		_ = v.Value // Value is types.AssetItem

	case *types.SearchInventoryResultItemMemberDataProductItem:
		_ = v.Value // Value is types.DataProductResultItem

	case *types.SearchInventoryResultItemMemberGlossaryItem:
		_ = v.Value // Value is types.GlossaryItem

	case *types.SearchInventoryResultItemMemberGlossaryTermItem:
		_ = v.Value // Value is types.GlossaryTermItem

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.GlossaryItem
var _ *types.DataProductResultItem
var _ *types.AssetItem
var _ *types.GlossaryTermItem

func ExampleSearchResultItem_outputUsage() {
	var union types.SearchResultItem
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.SearchResultItemMemberAssetListing:
		_ = v.Value // Value is types.AssetListingItem

	case *types.SearchResultItemMemberDataProductListing:
		_ = v.Value // Value is types.DataProductListingItem

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.DataProductListingItem
var _ *types.AssetListingItem

func ExampleSearchTypesResultItem_outputUsage() {
	var union types.SearchTypesResultItem
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.SearchTypesResultItemMemberAssetTypeItem:
		_ = v.Value // Value is types.AssetTypeItem

	case *types.SearchTypesResultItemMemberFormTypeItem:
		_ = v.Value // Value is types.FormTypeData

	case *types.SearchTypesResultItemMemberLineageNodeTypeItem:
		_ = v.Value // Value is types.LineageNodeTypeItem

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.LineageNodeTypeItem
var _ *types.FormTypeData
var _ *types.AssetTypeItem

func ExampleSelfGrantStatusOutput_outputUsage() {
	var union types.SelfGrantStatusOutput
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.SelfGrantStatusOutputMemberGlueSelfGrantStatus:
		_ = v.Value // Value is types.GlueSelfGrantStatusOutput

	case *types.SelfGrantStatusOutputMemberRedshiftSelfGrantStatus:
		_ = v.Value // Value is types.RedshiftSelfGrantStatusOutput

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.RedshiftSelfGrantStatusOutput
var _ *types.GlueSelfGrantStatusOutput

func ExampleSubscribedListingItem_outputUsage() {
	var union types.SubscribedListingItem
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.SubscribedListingItemMemberAssetListing:
		_ = v.Value // Value is types.SubscribedAssetListing

	case *types.SubscribedListingItemMemberProductListing:
		_ = v.Value // Value is types.SubscribedProductListing

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.SubscribedAssetListing
var _ *types.SubscribedProductListing

func ExampleSubscribedPrincipal_outputUsage() {
	var union types.SubscribedPrincipal
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.SubscribedPrincipalMemberProject:
		_ = v.Value // Value is types.SubscribedProject

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.SubscribedProject

func ExampleSubscribedPrincipalInput_outputUsage() {
	var union types.SubscribedPrincipalInput
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.SubscribedPrincipalInputMemberProject:
		_ = v.Value // Value is types.SubscribedProjectInput

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.SubscribedProjectInput

func ExampleUserPolicyGrantPrincipal_outputUsage() {
	var union types.UserPolicyGrantPrincipal
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.UserPolicyGrantPrincipalMemberAllUsersGrantFilter:
		_ = v.Value // Value is types.AllUsersGrantFilter

	case *types.UserPolicyGrantPrincipalMemberUserIdentifier:
		_ = v.Value // Value is string

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *string
var _ *types.AllUsersGrantFilter

func ExampleUserProfileDetails_outputUsage() {
	var union types.UserProfileDetails
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.UserProfileDetailsMemberIam:
		_ = v.Value // Value is types.IamUserProfileDetails

	case *types.UserProfileDetailsMemberSso:
		_ = v.Value // Value is types.SsoUserProfileDetails

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.SsoUserProfileDetails
var _ *types.IamUserProfileDetails
