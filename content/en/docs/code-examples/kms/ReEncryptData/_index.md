---
title: "ReEncryptDatav2"
---
## ReEncryptDatav2.go

This example reencrypts some text using an AWS Key Management Service (AWS KMS) customer master key (CMK).

`go run ReEncryptDatav2.go -k KeyID -d DATA`

- _KeyID_ is the ID of the AWS KMS key to use for reencrypting the data.
- _DATA_ is the data to reencrypt, as a string.

The unit test accepts similar values in _config.json_.

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
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

// KMSReEncryptAPI defines the interface for the ReEncrypt function.
// We use this interface to test the function using a mocked service.
type KMSReEncryptAPI interface {
	ReEncrypt(ctx context.Context,
		params *kms.ReEncryptInput,
		optFns ...func(*kms.Options)) (*kms.ReEncryptOutput, error)
}

// ReEncryptText reencrypts some text using a new AWS Key Management Service (AWS KMS) customer master key (CMK).
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a ReEncryptOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ReEncrypt.
func ReEncryptText(c context.Context, api KMSReEncryptAPI, input *kms.ReEncryptInput) (*kms.ReEncryptOutput, error) {
	return api.ReEncrypt(c, input)
}

func main() {
	keyID := flag.String("k", "", "The ID of a KMS key")
	data := flag.String("d", "", "The data to reencrypt, as a string")
	flag.Parse()

	if *keyID == "" || *data == "" {
		fmt.Println("You must supply the ID of a KMS key and data")
		fmt.Println("-k KEY-ID -d DATA")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := kms.NewFromConfig(cfg)

	blob := []byte(*data)

	input := &kms.ReEncryptInput{
		CiphertextBlob:   blob,
		DestinationKeyId: keyID,
	}

	result, err := ReEncryptText(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error reencrypting data:")
		fmt.Println(err)
		return
	}

	fmt.Println("Blob (base-64 byte array):")
	fmt.Println(result.CiphertextBlob)
}

```

See the [complete example in GitHub](https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/kms/ReEncryptData/ReEncryptDatav2.go).