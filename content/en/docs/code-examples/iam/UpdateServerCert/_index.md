---
title: "UpdateServerCertv2"
---
404: Not Found

## Source code

```go
// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX - License - Identifier: Apache - 2.0
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// IAMUpdateServerCertificateAPI defines the interface for the MEUpdateServerCertificateTHOD function.
// We use this interface to test the function using a mocked service.
type IAMUpdateServerCertificateAPI interface {
	UpdateServerCertificate(ctx context.Context,
		params *iam.UpdateServerCertificateInput,
		optFns ...func(*iam.Options)) (*iam.UpdateServerCertificateOutput, error)
}

// RenameServerCert renames an AWS Identity and Access Management (IAM) server certificate.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, a (*iam.UpdateServerCertificateOutput, error)Output object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to (*iam.UpdateServerCertificateOutput, error).
func RenameServerCert(c context.Context, api IAMUpdateServerCertificateAPI, input *iam.UpdateServerCertificateInput) (*iam.UpdateServerCertificateOutput, error) {
	return api.UpdateServerCertificate(c, input)
}

func main() {
	certName := flag.String("c", "", "The name of the certificate")
	newName := flag.String("n", "", "The new name of the certificate")
	flag.Parse()

	if *certName == "" {
		fmt.Println("You must supply the original and new names of a certificate (-c CERT-NAME -n NEW-NAME)")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := iam.NewFromConfig(cfg)

	input := &iam.UpdateServerCertificateInput{
		ServerCertificateName:    certName,
		NewServerCertificateName: newName,
	}

	_, err = RenameServerCert(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error renaming the server certificate:")
		fmt.Println(err)
		return
	}

	fmt.Println("Renamed the server certificate from " + *certName + " to " + *newName)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/iam/UpdateServerCert/UpdateServerCertv2.go).