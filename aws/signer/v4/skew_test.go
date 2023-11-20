package v4

import (
	"context"
	"net/http"
	"testing"
	"time"

	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/smithy-go/middleware"
	smithytime "github.com/aws/smithy-go/time"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func TestDetectSkew(t *testing.T) {
	sdk.NowTime = func() time.Time {
		return time.Unix(0, 0).UTC()
	}
	m := &DetectSkewMiddleware{
		Signer: NewSigner(),
	}

	stack := middleware.NewStack("TestDetectSkew", smithyhttp.NewStackRequest)
	stack.Finalize.Add(m, middleware.After)
	stack.Deserialize.Add(&awsmiddleware.RecordResponseTiming{}, middleware.After)

	hfn := func(ctx context.Context, input interface{}) (
		out interface{}, metadata middleware.Metadata, err error,
	) {
		resp := &smithyhttp.Response{
			Response: &http.Response{
				Header: http.Header{},
			},
		}
		resp.Header.Set("Date", smithytime.FormatHTTPDate(time.Unix(7, 0).UTC()))

		out = resp
		return out, metadata, err
	}
	handler := middleware.DecorateHandler(middleware.HandlerFunc(hfn), stack)
	handler.Handle(context.Background(), 1)

	expected := int64(7 * time.Second)
	actual := m.Signer.clockSkew.Load()
	if expected != actual {
		t.Fatalf("%v != %v", expected, actual)
	}
}
