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
	"github.com/aws/aws-sdk-go/aws/corehandlers"
	v1Creds "github.com/aws/aws-sdk-go/aws/credentials"
	v1Request "github.com/aws/aws-sdk-go/aws/request"
	v1Unit "github.com/aws/aws-sdk-go/awstesting/unit"
	v1Client "github.com/aws/aws-sdk-go/service/lexruntimeservice"
	"github.com/awslabs/smithy-go/ptr"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

func BenchmarkPutSession(b *testing.B) {
	b.Run("old", func(b *testing.B) {
		benchPutSessionOld(b)
	})

	b.Run("smithy", func(b *testing.B) {
		benchPutSessionSmithy(b)
	})
}

func benchPutSessionOld(b *testing.B) {
	sess := v1Unit.Session.Copy(&v1Aws.Config{
		Credentials: v1Creds.NewStaticCredentials("AKID", "SECRET", ""),
		Region:      ptr.String("us-west-2"),
	})
	sess.Handlers.Send.SwapNamed(v1Request.NamedHandler{
		Name: corehandlers.SendHandler.Name,
		Fn: func(r *v1Request.Request) {
			r.HTTPResponse = newPutSessionHTTPResponse()
		},
	})

	client := v1Client.New(sess)

	ctx := context.Background()
	params := v1Client.PutSessionInput{
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
			_, err := client.PutSessionWithContext(ctx, &params)
			if err != nil {
				b.Errorf("failed to send request: %v", err)
			}
		}
	})
}

func benchPutSessionSmithy(b *testing.B) {
	var args []func(*smithyClient.Options)
	if disableSmithySigning {
		args = append(args, removeSmithySigner)
	}

	client := smithyClient.New(smithyClient.Options{
		Region:      "us-west-2",
		Credentials: unit.StubCredentialsProvider{},
		HTTPClient: smithyhttp.ClientDoFunc(
			func(r *http.Request) (*http.Response, error) {
				return newPutSessionHTTPResponse(), nil
			}),
	}, args...)

	ctx := context.Background()
	params := smithyClient.PutSessionInput{
		Accept:   ptr.String("text/plain"),
		BotAlias: ptr.String("fooAlias"),
		BotName:  ptr.String("fooName"),
		DialogAction: &types.DialogAction{
			FulfillmentState: types.FulfillmentStateFulfilled,
			IntentName:       ptr.String("fooIntent"),
			Message:          ptr.String("fooMessage"),
			MessageFormat:    types.MessageFormatTypePlainText,
			SlotToElicit:     ptr.String("fooSlot"),
			Slots: map[string]string{
				"fooSlot": "fooValue",
				"barSlot": "barValue",
			},
			Type: types.DialogActionTypeElicitSlot,
		},
		RecentIntentSummaryView: []types.IntentSummary{
			{
				CheckpointLabel:    ptr.String("fooLabel"),
				ConfirmationStatus: types.ConfirmationStatusConfirmed,
				DialogActionType:   types.DialogActionTypeElicitSlot,
				FulfillmentState:   types.FulfillmentStateFulfilled,
				IntentName:         ptr.String("fooIntent"),
				SlotToElicit:       ptr.String("fooSlot"),
				Slots: map[string]string{
					"fooSlot": "fooValue",
					"barSlot": "barValue",
				},
			},
		},
		SessionAttributes: map[string]string{
			"fooAttr": "fooValue",
		},
		UserId: ptr.String("fooUser"),
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.PutSession(ctx, &params)
			if err != nil {
				b.Errorf("failed to send request: %v", err)
			}
		}
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
