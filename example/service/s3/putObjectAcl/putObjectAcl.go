// +build example

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/enums"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Put an ACL on an S3 object
//
// Usage:
// putBucketAcl <params>
//	-region <region> // required
//	-bucket <bucket> // required
//	-key <key> // required
//	-owner-name <owner-name>
//	-owner-id <owner-id>
//	-grantee-type <some type> // required
//	-uri <uri to group>
//	-email <email address>
//	-user-id <user-id>
func main() {
	regionPtr := flag.String("region", "", "region of your request")
	bucketPtr := flag.String("bucket", "", "name of your bucket")
	keyPtr := flag.String("key", "", "of your object")
	ownerNamePtr := flag.String("owner-name", "", "of your request")
	ownerIDPtr := flag.String("owner-id", "", "of your request")
	granteeTypePtr := flag.String("grantee-type", "", "of your request")
	uriPtr := flag.String("uri", "", "of your grantee type")
	emailPtr := flag.String("email", "", "of your grantee type")
	userPtr := flag.String("user-id", "", "of your grantee type")
	displayNamePtr := flag.String("display-name", "", "of your grantee type")
	flag.Parse()

	// Based off the type, fields must be excluded.
	switch enums.Type(*granteeTypePtr) {
	case enums.TypeCanonicalUser:
		emailPtr, uriPtr = nil, nil
		if *displayNamePtr == "" {
			displayNamePtr = nil
		}

		if *userPtr == "" {
			userPtr = nil
		}
	case enums.TypeAmazonCustomerByEmail:
		uriPtr, userPtr = nil, nil
	case enums.TypeGroup:
		emailPtr, userPtr = nil, nil
	}

	// Create a AWS SDK for Go config that will load credentials using the SDK's
	// default credential change.
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		exitErrorf("failed to load config, %v", err)
	}
	cfg.Region = *regionPtr

	svc := s3.New(cfg)

	req := svc.PutObjectAclRequest(&types.PutObjectAclInput{
		Bucket: bucketPtr,
		Key:    keyPtr,
		AccessControlPolicy: &types.AccessControlPolicy{
			Owner: &types.Owner{
				DisplayName: ownerNamePtr,
				ID:          ownerIDPtr,
			},
			Grants: []types.Grant{
				{
					Grantee: &types.Grantee{
						Type:         enums.Type(*granteeTypePtr),
						DisplayName:  displayNamePtr,
						URI:          uriPtr,
						EmailAddress: emailPtr,
						ID:           userPtr,
					},
					Permission: enums.PermissionFullControl,
				},
			},
		},
	})
	resp, err := req.Send(context.Background())
	if err != nil {
		exitErrorf("failed, %v", err)
	}

	fmt.Println("success", resp)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
