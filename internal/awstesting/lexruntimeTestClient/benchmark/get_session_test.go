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
	"github.com/awslabs/smithy-go/ptr"
)

var getSessionFullResponse = []byte(`{
  "dialogAction": {
    "fulfillmentState": "ReadyForFulfillment",
    "intentName": "fooIntent",
    "message": "fooMessage",
    "messageFormat": "PlainText",
    "slotToElicit": "fooSlot",
    "slots": {
      "fooSlot": "slotValue",
      "barSlot": "slotValue"
    },
    "type": "ConfirmIntent"
  },
  "recentIntentSummaryView": [
    {
      "checkpointLabel": "fooLabel",
      "confirmationStatus": "Confirmed",
      "dialogActionType": "ConfirmIntent",
      "fulfillmentState": "Fulfilled",
      "intentName": "fooIntent",
      "slotToElicit": "fooSlot",
      "slots": {
        "fooSlot": "slotValue",
        "barSlot": "slotValue"
      }
    }
  ],
  "sessionAttributes": {
    "foo": "fooValue",
    "bar": "barValue"
  },
  "sessionId": "benchmark"
}
`)

func getSessionHandler(tb testing.TB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		tb.Helper()
		_, err := writer.Write(getSessionFullResponse)
		if err != nil {
			tb.Errorf("failed to write http response: %v", err)
		}
	}

}

func BenchmarkGetSession_Old(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(getSessionHandler(b)))
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
		request := client.GetSessionRequest(&oldClient.GetSessionInput{
			BotAlias:              ptr.String("fooAlias"),
			BotName:               ptr.String("fooName"),
			CheckpointLabelFilter: ptr.String("fooFilter"),
			UserId:                ptr.String("fooUser"),
		})
		_, err := request.Send(context.Background())
		if err != nil {
			b.Errorf("failed to send request: %v", err)
		}
		b.StopTimer()
		server.CloseClientConnections()
	}
}

func BenchmarkGetSession_Smithy(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(getSessionHandler(b)))
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
		_, err := client.GetSession(context.Background(), &smithyClient.GetSessionInput{
			BotAlias:              ptr.String("fooAlias"),
			BotName:               ptr.String("fooName"),
			CheckpointLabelFilter: ptr.String("fooFilter"),
			UserId:                ptr.String("fooUser"),
		})
		if err != nil {
			b.Errorf("failed to send request: %v", err)
		}
		b.StopTimer()
		server.CloseClientConnections()
	}
}
