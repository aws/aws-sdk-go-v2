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
func (a *API) customizationPasses() error {
	var svcCustomizations = map[string]func(*API) error{
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
		err := fn(a)
		if err != nil {
			return fmt.Errorf("service customization pass failure for %s: %v", a.PackageName(), err)
		}
	}

	a.EnableSelectGeneratedMarshalers()

	return nil
}

func supressSmokeTest(a *API) error {
	a.SmokeTests.TestCases = []SmokeTestCase{}
	return nil
}

// s3Customizations customizes the API generation to replace values specific to S3.
func s3Customizations(a *API) error {
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

	// Generate an endpointARN method for the BucketName shape operation inputs
	for _, o := range a.Operations {
		if o.Name == "CreateBucket" {
			// For all operations but CreateBucket the BucketName shape
			// needs to be decorated.
			continue
		}
		var endpointARNShape *ShapeRef
		for _, ref := range o.InputRef.Shape.MemberRefs {
			if ref.OrigShapeName != "BucketName" || ref.Shape.Type != "string" {
				continue
			}
			if endpointARNShape != nil {
				return fmt.Errorf("more then one BucketName shape present on shape")
			}
			ref.EndpointARN = true
			endpointARNShape = ref
		}
		if endpointARNShape != nil {
			o.InputRef.Shape.HasEndpointARNMember = true
			o.HasEndpointARN = true
		}
	}

	return nil
}

// S3 Control service operations with an AccountId need accessors to be
// generated for them so the fields can be dynamically accessed without
// reflection.
func s3ControlCustomizations(a *API) error {
	for _, op := range a.Operations {
		// Add moving AccountId into the hostname instead of header.
		if _, ok := op.InputRef.Shape.MemberRefs["AccountId"]; ok {
			op.CustomBuildHandlers = append(op.CustomBuildHandlers,
				`buildPrefixHostHandler("AccountID", aws.StringValue(input.AccountId))`,
				`buildRemoveHeaderHandler("X-Amz-Account-Id")`,
			)
		}
	}

	return nil
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
func cloudfrontCustomizations(a *API) error {
	// MaxItems members should always be integers
	for _, s := range a.Shapes {
		if ref, ok := s.MemberRefs["MaxItems"]; ok {
			ref.ShapeName = "Integer"
			ref.Shape = a.Shapes["Integer"]
		}
	}

	return nil
}

// rdsCustomizations are customization for the service/rds. This adds non-modeled fields used for presigning.
func rdsCustomizations(a *API) error {
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

	return nil
}

func backfillAuthType(typ AuthType, opNames ...string) func(*API) error {
	return func(a *API) error {
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
		return nil
	}
}
