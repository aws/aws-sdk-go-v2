package benchmark

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	smithyClient "github.com/aws/aws-sdk-go-v2/service/lexruntimeservice"
	v1Aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/corehandlers"
	v1Creds "github.com/aws/aws-sdk-go/aws/credentials"
	v1Request "github.com/aws/aws-sdk-go/aws/request"
	v1Unit "github.com/aws/aws-sdk-go/awstesting/unit"
	v1Client "github.com/aws/aws-sdk-go/service/lexruntimeservice"
	"github.com/awslabs/smithy-go/ptr"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

func BenchmarkGetSession(b *testing.B) {
	filename := filepath.Join("testdata", "get_session_resp.json")
	bodyBytes, err := loadTestData(filename)
	if err != nil {
		b.Fatalf("failed to load test data, %s, %v", filename, err)
	}

	b.Run("old", func(b *testing.B) {
		benchGetSessionOld(b, bodyBytes)
	})

	b.Run("smithy", func(b *testing.B) {
		benchGetSessionSmithy(b, bodyBytes)
	})
}

func benchGetSessionOld(b *testing.B, respBytes []byte) {
	sess := v1Unit.Session.Copy(&v1Aws.Config{
		Credentials: v1Creds.NewStaticCredentials("AKID", "SECRET", ""),
		Region:      ptr.String("us-west-2"),
	})
	sess.Handlers.Send.SwapNamed(v1Request.NamedHandler{
		Name: corehandlers.SendHandler.Name,
		Fn: func(r *v1Request.Request) {
			r.HTTPResponse = newGetSessionHTTPResponse(respBytes)
		},
	})

	client := v1Client.New(sess)
	params := v1Client.GetSessionInput{
		BotAlias:              ptr.String("fooAlias"),
		BotName:               ptr.String("fooName"),
		CheckpointLabelFilter: ptr.String("fooFilter"),
		UserId:                ptr.String("fooUser"),
	}

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.GetSessionWithContext(ctx, &params)
			if err != nil {
				b.Fatalf("failed to send request: %v", err)
			}
		}
	})
}

func benchGetSessionSmithy(b *testing.B, respBytes []byte) {
	var args []func(*smithyClient.Options)
	if disableSmithySigning {
		args = append(args, removeSmithySigner)
	}

	client := smithyClient.New(smithyClient.Options{
		Region:      "us-west-2",
		Credentials: unit.StubCredentialsProvider{},
		HTTPClient: smithyhttp.ClientDoFunc(
			func(r *http.Request) (*http.Response, error) {
				return newGetSessionHTTPResponse(respBytes), nil
			}),
	}, args...)

	ctx := context.Background()
	params := smithyClient.GetSessionInput{
		BotAlias:              ptr.String("fooAlias"),
		BotName:               ptr.String("fooName"),
		CheckpointLabelFilter: ptr.String("fooFilter"),
		UserId:                ptr.String("fooUser"),
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.GetSession(ctx, &params)
			if err != nil {
				b.Fatalf("failed to send: %v", err)
			}
		}
	})
}

func newGetSessionHTTPResponse(body []byte) *http.Response {
	return &http.Response{
		StatusCode: 200,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		ContentLength: int64(len(body)),
		Body:          ioutil.NopCloser(bytes.NewReader(body)),
	}
}
