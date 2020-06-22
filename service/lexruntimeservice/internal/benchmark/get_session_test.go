package benchmark

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
	smithyClient "github.com/aws/aws-sdk-go-v2/service/lexruntimeservice"
	v1Aws "github.com/aws/aws-sdk-go/aws"
	v1Request "github.com/aws/aws-sdk-go/aws/request"
	v1Session "github.com/aws/aws-sdk-go/aws/session"
	v1Unit "github.com/aws/aws-sdk-go/awstesting/unit"
	v1Client "github.com/aws/aws-sdk-go/service/lexruntimeservice"
	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/ptr"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

func addGetSessionAPIResponseMiddleware(options *smithyClient.Options) {
	options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
		stack.Deserialize.Add(middleware.DeserializeMiddlewareFunc("StubResponse", func(ctx context.Context, input middleware.DeserializeInput, next middleware.DeserializeHandler,
		) (out middleware.DeserializeOutput, metadata middleware.Metadata, err error) {
			out.RawResponse = &smithyhttp.Response{Response: newGetSessionHTTPResponse()}
			return out, metadata, err
		}), middleware.After)
		return nil
	})
}

func addGetSessionAPIResponseHandler(sess *v1Session.Session) {
	sess.Handlers.Send.Swap("core.SendHandler", v1Request.NamedHandler{
		Name: "StubHandler",
		Fn: func(r *v1Request.Request) {
			r.HTTPResponse = newGetSessionHTTPResponse()
		},
	})
}

func newGetSessionHTTPResponse() *http.Response {
	return &http.Response{
		StatusCode:    200,
		ContentLength: int64(len(getSessionResponseBody)),
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: ioutil.NopCloser(bytes.NewReader(getSessionResponseBody)),
	}
}

var getSessionResponseBody = []byte(`{
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
      "checkpointLabel": "fooLabel1",
      "confirmationStatus": "Confirmed",
      "dialogActionType": "ConfirmIntent",
      "fulfillmentState": "Fulfilled",
      "intentName": "fooIntent",
      "slotToElicit": "fooSlot",
      "slots": {
        "fooSlot": "slotValue",
        "barSlot": "slotValue"
      }
    },
    {
      "checkpointLabel": "fooLabel2",
      "confirmationStatus": "Confirmed",
      "dialogActionType": "ConfirmIntent",
      "fulfillmentState": "Fulfilled",
      "intentName": "fooIntent",
      "slotToElicit": "fooSlot",
      "slots": {
        "fooSlot": "slotValue",
        "barSlot": "slotValue"
      }
    },
    {
      "checkpointLabel": "fooLabel3",
      "confirmationStatus": "Confirmed",
      "dialogActionType": "ConfirmIntent",
      "fulfillmentState": "Fulfilled",
      "intentName": "fooIntent",
      "slotToElicit": "fooSlot",
      "slots": {
        "fooSlot": "slotValue",
        "barSlot": "slotValue"
      }
    },
    {
      "checkpointLabel": "fooLabel4",
      "confirmationStatus": "Confirmed",
      "dialogActionType": "ConfirmIntent",
      "fulfillmentState": "Fulfilled",
      "intentName": "fooIntent",
      "slotToElicit": "fooSlot",
      "slots": {
        "fooSlot": "slotValue",
        "barSlot": "slotValue"
      }
    },
    {
      "checkpointLabel": "fooLabel5",
      "confirmationStatus": "Confirmed",
      "dialogActionType": "ConfirmIntent",
      "fulfillmentState": "Fulfilled",
      "intentName": "fooIntent",
      "slotToElicit": "fooSlot",
      "slots": {
        "fooSlot": "slotValue",
        "barSlot": "slotValue"
      }
    },
    {
      "checkpointLabel": "fooLabel6",
      "confirmationStatus": "Confirmed",
      "dialogActionType": "ConfirmIntent",
      "fulfillmentState": "Fulfilled",
      "intentName": "fooIntent",
      "slotToElicit": "fooSlot",
      "slots": {
        "fooSlot": "slotValue",
        "barSlot": "slotValue"
      }
    },
    {
      "checkpointLabel": "fooLabel7",
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
}`)

func BenchmarkGetSession_Old(b *testing.B) {
	cfg := v1Unit.Session.Copy(&v1Aws.Config{Region: ptr.String("us-west-2")})
	addGetSessionAPIResponseHandler(cfg)
	//cfg.Handlers.Sign.Clear()
	client := v1Client.New(cfg)
	in := &v1Client.GetSessionInput{
		BotAlias:              ptr.String("fooAlias"),
		BotName:               ptr.String("fooName"),
		CheckpointLabelFilter: ptr.String("fooFilter"),
		UserId:                ptr.String("fooUser"),
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			request, _ := client.GetSessionRequest(in)
			err := request.Send()
			if err != nil {
				b.Errorf("failed to send v1Request: %v", err)
			}
		}
	})
}

func BenchmarkGetSession_Smithy(b *testing.B) {
	cfg := unit.Config()
	cfg.Region = "us-west-2"

	client := smithyClient.NewFromConfig(cfg, addGetSessionAPIResponseMiddleware)
	ctx := context.Background()
	in := &smithyClient.GetSessionInput{
		BotAlias:              ptr.String("fooAlias"),
		BotName:               ptr.String("fooName"),
		CheckpointLabelFilter: ptr.String("fooFilter"),
		UserId:                ptr.String("fooUser"),
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.GetSession(ctx, in)
			if err != nil {
				b.Errorf("failed to send: %v", err)
			}
		}
	})
}

func removeSmithySigner(options *smithyClient.Options) {
	options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
		stack.Finalize.Remove("SigV4SignHTTPRequestMiddleware")
		stack.Finalize.Remove("SigV4ContentSHA256HeaderMiddleware")
		stack.Finalize.Remove("ComputePayloadSHA256Middleware")
		stack.Finalize.Remove("SigV4UnsignedPayloadMiddleware")
		return nil
	})
}
