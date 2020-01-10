// +build codegen

package api

import (
	"fmt"
	"os"
	"strings"
)

type service struct {
	srcName string
	dstName string

	serviceVersion string
}

func (a *API) EnableSelectGeneratedMarshalers() {
	// Selectivily enable generated marshalers as available
	a.NoGenMarshalers = true
	a.NoGenUnmarshalers = true

	// Enable generated marshalers
	switch a.Metadata.Protocol {
	case "rest-xml", "rest-json":
		a.NoGenMarshalers = false
	}
}

// customizationPasses Executes customization logic for the API by package name.
func (a *API) customizationPasses() {
	var svcCustomizations = map[string]func(*API){
		"s3":         s3Customizations,
		"s3control":  s3ControlCustomizations,
		"cloudfront": cloudfrontCustomizations,
		"rds":        rdsCustomizations,

		// MTurk smoke test is invalid. The service requires AWS account to be
		// linked to Amazon Mechanical Turk Account.
		"mturk": supressSmokeTest,

		// Backfill the authentication type for cognito identity and sts.
		// Removes the need for the customizations in these services.
		"cognitoidentity": backfillAuthType(NoneAuthType,
			"GetId",
			"GetOpenIdToken",
			"UnlinkIdentity",
			"GetCredentialsForIdentity",
		),
		"sts": backfillAuthType(NoneAuthType,
			"AssumeRoleWithSAML",
			"AssumeRoleWithWebIdentity",
		),
	}

	if fn := svcCustomizations[a.PackageName()]; fn != nil {
		fn(a)
	}

	a.EnableSelectGeneratedMarshalers()
}

func supressSmokeTest(a *API) {
	a.SmokeTests.TestCases = []SmokeTestCase{}
}

// s3Customizations customizes the API generation to replace values specific to S3.
func s3Customizations(a *API) {
	var strExpires *Shape

	var keepContentMD5Ref = map[string]struct{}{
		"PutObjectInput":  {},
		"UploadPartInput": {},
	}

	for name, s := range a.Shapes {
		// Remove ContentMD5 members unless specified otherwise.
		if _, keep := keepContentMD5Ref[name]; !keep {
			if _, have := s.MemberRefs["ContentMD5"]; have {
				delete(s.MemberRefs, "ContentMD5")
			}
		}

		// Generate getter methods for API operation fields used by customizations.
		for _, refName := range []string{"Bucket", "SSECustomerKey", "CopySourceSSECustomerKey"} {
			if ref, ok := s.MemberRefs[refName]; ok {
				ref.GenerateGetter = true
			}
		}

		// Expires should be a string not time.Time since the format is not
		// enforced by S3, and any value can be set to this field outside of the SDK.
		if strings.HasSuffix(name, "Output") {
			if ref, ok := s.MemberRefs["Expires"]; ok {
				if strExpires == nil {
					newShape := *ref.Shape
					strExpires = &newShape
					strExpires.Type = "string"
					strExpires.refs = []*ShapeRef{}
				}
				ref.Shape.removeRef(ref)
				ref.Shape = strExpires
				ref.Shape.refs = append(ref.Shape.refs, &s.MemberRef)
			}
		}
	}
	s3CustRemoveHeadObjectModeledErrors(a)
}

// S3 Control service operations with an AccountId need accessors to be
// generated for them so the fields can be dynamically accessed without
// reflection.
func s3ControlCustomizations(a *API) {
	for _, op := range a.Operations {
		// Add moving AccountId into the hostname instead of header.
		if _, ok := op.InputRef.Shape.MemberRefs["AccountId"]; ok {
			op.CustomBuildHandlers = append(op.CustomBuildHandlers,
				`buildPrefixHostHandler("AccountID", aws.StringValue(input.AccountId))`,
				`buildRemoveHeaderHandler("X-Amz-Account-Id")`,
			)
		}
	}
}

// S3 HeadObject API call incorrect models NoSuchKey as valid
// error code that can be returned. This operation does not
// return error codes, all error codes are derived from HTTP
// status codes.
//
// aws/aws-sdk-go#1208
func s3CustRemoveHeadObjectModeledErrors(a *API) {
	op, ok := a.Operations["HeadObject"]
	if !ok {
		return
	}
	op.Documentation += `
//
// See http://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html#RESTErrorResponses
// for more information on returned errors.`
	op.ErrorRefs = []ShapeRef{}
}

// cloudfrontCustomizations customized the API generation to replace values
// specific to CloudFront.
func cloudfrontCustomizations(a *API) {
	// MaxItems members should always be integers
	for _, s := range a.Shapes {
		if ref, ok := s.MemberRefs["MaxItems"]; ok {
			ref.ShapeName = "Integer"
			ref.Shape = a.Shapes["Integer"]
		}
	}
}

// rdsCustomizations are customization for the service/rds. This adds non-modeled fields used for presigning.
func rdsCustomizations(a *API) {
	inputs := []string{
		"CopyDBSnapshotInput",
		"CreateDBInstanceReadReplicaInput",
		"CopyDBClusterSnapshotInput",
		"CreateDBClusterInput",
	}
	for _, input := range inputs {
		if ref, ok := a.Shapes[input]; ok {
			ref.MemberRefs["SourceRegion"] = &ShapeRef{
				Documentation: docstring(`SourceRegion is the source region where the resource exists. This is not sent over the wire and is only used for presigning. This value should always have the same region as the source ARN.`),
				ShapeName:     "String",
				Shape:         a.Shapes["String"],
				Ignore:        true,
			}
			ref.MemberRefs["DestinationRegion"] = &ShapeRef{
				Documentation: docstring(`DestinationRegion is used for presigning the request to a given region.`),
				ShapeName:     "String",
				Shape:         a.Shapes["String"],
			}
		}
	}
}
func backfillAuthType(typ AuthType, opNames ...string) func(*API) {
	return func(a *API) {
		for _, opName := range opNames {
			op, ok := a.Operations[opName]
			if !ok {
				panic("unable to backfill auth-type for unknown operation " + opName)
			}
			if v := op.AuthType; len(v) != 0 {
				fmt.Fprintf(os.Stderr, "unable to backfill auth-type for %s, already set, %s", opName, v)
				continue
			}

			op.AuthType = typ
		}
	}
}
