// +build go1.6

package s3

import (
	request "github.com/aws/aws-sdk-go-v2/aws"
)

func platformRequestHandlers(r *request.Request) {
	if r.Operation.HTTPMethod == "PUT" {
		// 100-Continue should only be used on put requests.
		r.Handlers.Sign.PushBack(add100Continue)
	}
}

func add100Continue(r *request.Request) {
	if r.Config.S3Disable100Continue {
		return
	}
	if r.HTTPRequest.ContentLength < 1024*1024*2 {
		// Ignore requests smaller than 2MB. This helps prevent delaying
		// requests unnecessarily.
		return
	}

	r.HTTPRequest.Header.Set("Expect", "100-Continue")
}
