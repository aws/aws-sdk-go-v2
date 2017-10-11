// +build !go1.6

package s3

import "github.com/aws/aws-sdk-go-v2/aws/request"

func platformRequestHandlers(r *request.Request) {
}
