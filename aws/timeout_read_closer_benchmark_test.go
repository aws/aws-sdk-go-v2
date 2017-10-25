package aws_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	"github.com/aws/aws-sdk-go-v2/private/protocol/jsonrpc"
)

func BenchmarkTimeoutReadCloser(b *testing.B) {
	resp := `
	{
		"Bar": "qux"
	}
	`

	handlers := aws.Handlers{}
	handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(resp))),
		}
	})
	handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)

	op := &aws.Operation{
		Name:       "op",
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	meta := aws.Metadata{
		ServiceName:   "fooService",
		SigningName:   "foo",
		SigningRegion: "foo",
		Endpoint:      "localhost",
		APIVersion:    "2001-01-01",
		JSONVersion:   "1.1",
		TargetPrefix:  "Foo",
	}

	cfg := unit.Config()
	cfg.Handlers = handlers

	req := aws.New(
		cfg,
		meta,
		handlers,
		aws.DefaultRetryer{NumMaxRetries: 5},
		op,
		&struct {
			Foo *string
		}{},
		&struct {
			Bar *string
		}{},
	)

	req.ApplyOptions(aws.WithResponseReadTimeout(15 * time.Second))
	for i := 0; i < b.N; i++ {
		err := req.Send()
		if err != nil {
			b.Errorf("Expected no error, but received %v", err)
		}
	}
}
