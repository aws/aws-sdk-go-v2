package sqs

import request "github.com/aws/aws-sdk-go-v2/aws"

func init() {
	initRequest = func(r *request.Request) {
		setupChecksumValidation(r)
	}
}
