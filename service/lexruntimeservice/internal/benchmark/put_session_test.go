package benchmark

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	smithyClient "github.com/aws/aws-sdk-go-v2/service/lexruntimeservice"
	"github.com/aws/aws-sdk-go-v2/service/lexruntimeservice/types"
	v1Aws "github.com/aws/aws-sdk-go/aws"
	v1Request "github.com/aws/aws-sdk-go/aws/request"
	v1Session "github.com/aws/aws-sdk-go/aws/session"
	v1Unit "github.com/aws/aws-sdk-go/awstesting/unit"
	v1Client "github.com/aws/aws-sdk-go/service/lexruntimeservice"
	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/ptr"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

func addPutSessionAPIResponseMiddleware(options *smithyClient.Options) {
	options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
		stack.Deserialize.Add(middleware.DeserializeMiddlewareFunc("StubResponse", func(ctx context.Context, input middleware.DeserializeInput, next middleware.DeserializeHandler,
		) (out middleware.DeserializeOutput, metadata middleware.Metadata, err error) {
			out.RawResponse = &smithyhttp.Response{Response: newPutSessionHTTPResponse()}
			return out, metadata, err
		}), middleware.After)
		return nil
	})
}

func addPutSessionAPIResponseHandler(sess *v1Session.Session) {
	sess.Handlers.Send.Swap("core.SendHandler", v1Request.NamedHandler{
		Name: "StubHandler",
		Fn: func(r *v1Request.Request) {
			r.HTTPResponse = newPutSessionHTTPResponse()
		},
	})
}

var putSessionBody = []byte("fooAudioStream")

func newPutSessionHTTPResponse() *http.Response {
	return &http.Response{
		StatusCode:    200,
		ContentLength: int64(len(putSessionBody)),
		Header: map[string][]string{
			"Content-Type":                 {"application/octet-stream"},
			"x-amz-lex-dialog-state":       {"Fulfilled"},
			"x-amz-lex-intent-name":        {"fooIntent"},
			"x-amz-lex-message":            {"fooMessage"},
			"x-amz-lex-message-format":     {"PlainText"},
			"x-amz-lex-session-attributes": {"eyJmb29LZXkiOiAiZm9vVmFsdWUifQ=="},
			"x-amz-lex-session-id":         {"fooSession"},
			"x-amz-lex-slot-to-elicit":     {"fooSlot"},
			"x-amz-lex-slots":              {"eyJmb29LZXkiOiAiZm9vVmFsdWUifQ=="},
		},
		Body: ioutil.NopCloser(bytes.NewReader(putSessionBody)),
	}
}

func BenchmarkPutSession_Old(b *testing.B) {
	cfg := v1Unit.Session.Copy(&v1Aws.Config{Region: ptr.String("us-west-2")})
	addPutSessionAPIResponseHandler(cfg)
	//cfg.Handlers.Sign.Clear()
	client := v1Client.New(cfg)

	in := &v1Client.PutSessionInput{
		Accept:   ptr.String("text/plain"),
		BotAlias: ptr.String("fooAlias"),
		BotName:  ptr.String("fooName"),
		DialogAction: &v1Client.DialogAction{
			FulfillmentState: ptr.String(v1Client.FulfillmentStateFulfilled),
			IntentName:       ptr.String("fooIntent"),
			Message:          ptr.String("fooMessage"),
			MessageFormat:    ptr.String(v1Client.MessageFormatTypePlainText),
			SlotToElicit:     ptr.String("fooSlot"),
			Slots: ptr.StringMap(map[string]string{
				"fooSlot": "fooValue",
				"barSlot": "barValue",
			}),
			Type: ptr.String(v1Client.DialogActionTypeElicitSlot),
		},
		RecentIntentSummaryView: []*v1Client.IntentSummary{
			{
				CheckpointLabel:    ptr.String("fooLabel"),
				ConfirmationStatus: ptr.String(v1Client.ConfirmationStatusConfirmed),
				DialogActionType:   ptr.String(v1Client.DialogActionTypeElicitSlot),
				FulfillmentState:   ptr.String(v1Client.FulfillmentStateFulfilled),
				IntentName:         ptr.String("fooIntent"),
				SlotToElicit:       ptr.String("fooSlot"),
				Slots: ptr.StringMap(map[string]string{
					"fooSlot": "fooValue",
					"barSlot": "barValue",
				}),
			},
		},
		SessionAttributes: ptr.StringMap(map[string]string{
			"fooAttr": "fooValue",
		}),
		UserId: ptr.String("fooUser"),
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			request, _ := client.PutSessionRequest(in)
			err := request.Send()
			if err != nil {
				b.Errorf("failed to send request: %v", err)
			}
		}
	})
}

func BenchmarkPutSession_Smithy(b *testing.B) {
	cfg := unit.Config()
	cfg.Region = "us-west-2"
	client := smithyClient.NewFromConfig(cfg, addPutSessionAPIResponseMiddleware)

	ctx := context.Background()
	in := &smithyClient.PutSessionInput{
		Accept:   ptr.String("text/plain"),
		BotAlias: ptr.String("fooAlias"),
		BotName:  ptr.String("fooName"),
		DialogAction: &types.DialogAction{
			FulfillmentState: types.FulfillmentStateFulfilled,
			IntentName:       ptr.String("fooIntent"),
			Message:          ptr.String("fooMessage"),
			MessageFormat:    types.MessageFormatTypePlain_text,
			SlotToElicit:     ptr.String("fooSlot"),
			Slots: map[string]*string{
				"fooSlot": ptr.String("fooValue"),
				"barSlot": ptr.String("barValue"),
			},
			Type: types.DialogActionTypeElicit_slot,
		},
		RecentIntentSummaryView: []*types.IntentSummary{
			{
				CheckpointLabel:    ptr.String("fooLabel"),
				ConfirmationStatus: types.ConfirmationStatusConfirmed,
				DialogActionType:   types.DialogActionTypeElicit_slot,
				FulfillmentState:   types.FulfillmentStateFulfilled,
				IntentName:         ptr.String("fooIntent"),
				SlotToElicit:       ptr.String("fooSlot"),
				Slots: map[string]*string{
					"fooSlot": ptr.String("fooValue"),
					"barSlot": ptr.String("barValue"),
				},
			},
		},
		SessionAttributes: map[string]*string{
			"fooAttr": ptr.String("fooValue"),
		},
		UserId: ptr.String("fooUser"),
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.PutSession(ctx, in)
			if err != nil {
				b.Errorf("failed to send request: %v", err)
			}
		}
	})
}
