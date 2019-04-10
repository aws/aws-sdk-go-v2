package s3control

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/s3err"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

type accountIDGetter interface {
	getAccountId() string
}

func init() {
	initClient = defaultInitClientFn
}
func defaultInitClientFn(c *S3Control) {
	c.Handlers.UnmarshalError.PushBackNamed(s3err.RequestFailureWrapperHandler())
}
func buildPrefixHostHandler(fieldName, value string) aws.NamedHandler {
	return aws.NamedHandler{
		Name: "awssdk.s3control.prefixhost",
		Fn: func(r *aws.Request) {
			paramErrs := aws.ErrInvalidParams{Context: r.Operation.Name}
			if !protocol.ValidHostLabel(value) {
				paramErrs.Add(aws.NewErrParamFormat(fieldName, "[a-zA-Z0-9-]{1,63}", value))
				r.Error = paramErrs
				return
			}
			r.HTTPRequest.URL.Host = value + "." + r.HTTPRequest.URL.Host
		},
	}
}
func buildRemoveHeaderHandler(key string) aws.NamedHandler {
	return aws.NamedHandler{
		Name: "awssdk.s3control.removeHeader",
		Fn: func(r *aws.Request) {
			r.HTTPRequest.Header.Del(key)
		},
	}
}
