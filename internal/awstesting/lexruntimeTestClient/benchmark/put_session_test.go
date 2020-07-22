package benchmark

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	oldClient "github.com/aws/aws-sdk-go-v2/internal/awstesting/lexruntimeTestClient/benchmark/lexruntimeservice"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	smithyClient "github.com/aws/aws-sdk-go-v2/service/lexruntimeservice"
	"github.com/aws/aws-sdk-go-v2/service/lexruntimeservice/types"
	"github.com/awslabs/smithy-go/ptr"
)

var putSessionFullResponse = []byte("fooAudioStream")

func putSessionHandler(tb testing.TB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		tb.Helper()
		header := writer.Header()
		header.Set("Content-Type", "application/octet-stream")
		header.Set("x-amz-lex-dialog-state", "Fulfilled")
		header.Set("x-amz-lex-intent-name", "fooIntent")
		header.Set("x-amz-lex-message", "fooMessage")
		header.Set("x-amz-lex-message-format", "PlainText")
		header.Set("x-amz-lex-session-attributes", "eyJmb29LZXkiOiAiZm9vVmFsdWUifQ==")
		header.Set("x-amz-lex-session-id", "fooSession")
		header.Set("x-amz-lex-slot-to-elicit", "fooSlot")
		header.Set("x-amz-lex-slots", "eyJmb29LZXkiOiAiZm9vVmFsdWUifQ==")
		writer.WriteHeader(http.StatusOK)
		_, err := writer.Write(putSessionFullResponse)
		if err != nil {
			tb.Errorf("failed to write http response: %v", err)
		}
	}
}

func BenchmarkPutSession_Old(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(putSessionHandler(b)))
	defer server.Close()

	cfg := unit.Config()
	cfg.Retryer = &aws.NoOpRetryer{}
	cfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:         server.URL,
			SigningName: "foo",
		}, nil
	})
	client := oldClient.New(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		request := client.PutSessionRequest(&oldClient.PutSessionInput{
			Accept:   ptr.String("text/plain"),
			BotAlias: ptr.String("fooAlias"),
			BotName:  ptr.String("fooName"),
			DialogAction: &oldClient.DialogAction{
				FulfillmentState: oldClient.FulfillmentStateFulfilled,
				IntentName:       ptr.String("fooIntent"),
				Message:          ptr.String("fooMessage"),
				MessageFormat:    oldClient.MessageFormatTypePlainText,
				SlotToElicit:     ptr.String("fooSlot"),
				Slots: map[string]string{
					"fooSlot": "fooValue",
					"barSlot": "barValue",
				},
				Type: oldClient.DialogActionTypeElicitSlot,
			},
			RecentIntentSummaryView: []oldClient.IntentSummary{
				{
					CheckpointLabel:    ptr.String("fooLabel"),
					ConfirmationStatus: oldClient.ConfirmationStatusConfirmed,
					DialogActionType:   oldClient.DialogActionTypeElicitSlot,
					FulfillmentState:   oldClient.FulfillmentStateFulfilled,
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
		})
		_, err := request.Send(context.Background())
		if err != nil {
			b.Errorf("failed to send request: %v", err)
		}
		b.StopTimer()
		server.CloseClientConnections()
	}
}

func BenchmarkPutSession_Smithy(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(putSessionHandler(b)))
	defer server.Close()

	cfg := unit.Config()
	cfg.Retryer = &aws.NoOpRetryer{}
	cfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:         server.URL,
			SigningName: "foo",
		}, nil
	})
	client := smithyClient.NewFromConfig(cfg, func(o *smithyClient.Options) {
		o.HTTPSigner = v4.NewSigner(o.Credentials)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, err := client.PutSession(context.Background(), &smithyClient.PutSessionInput{
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
		})
		if err != nil {
			b.Errorf("failed to send request: %v", err)
		}
		b.StopTimer()
		server.CloseClientConnections()
	}
}
