package testing

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream/eventstreamapi"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/aws/aws-sdk-go-v2/service/internal/eventstreamtesting"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
)

const (
	ModelID = "amazon.nova-sonic-v1:0"
)

func removeValidationMiddleware(stack *middleware.Stack) error {
	_, err := stack.Initialize.Remove("OperationInputValidation")
	return err
}

func toBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func TestStartStreamTranscription_Read(t *testing.T) {
	payload := `{}`
	cfg, cleanupFn, err := setupBasicEventStream(t, payload)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := bedrockruntime.NewFromConfig(cfg)
	resp, err := svc.InvokeModelWithBidirectionalStream(context.Background(),
		&bedrockruntime.InvokeModelWithBidirectionalStreamInput{ModelId: aws.String(ModelID)}, func(options *bedrockruntime.Options) {
			options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
		})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()
	<-resp.GetInitialReply()

	expectEvents := []types.InvokeModelWithBidirectionalStreamOutput{
		&types.InvokeModelWithBidirectionalStreamOutputMemberChunk{
			Value: types.BidirectionalOutputPayloadPart{Bytes: []byte(payload)},
		},
	}

	for i := 0; i < len(expectEvents); i++ {
		event := <-resp.GetStream().Events()
		if event == nil {
			t.Errorf("%d, expect event, got nil", i)
		}
		if diff := cmpDiff(expectEvents[i], event); len(diff) > 0 {
			t.Errorf("expected %v as %d event, got %v", expectEvents[i], i, event)
		}
	}

	if err := resp.GetStream().Err(); err != nil {
		t.Errorf("expect no error, %v", err)
	}
}

func TestStartStreamTranscription_ReadClose(t *testing.T) {
	payload := `{}`
	cfg, cleanupFn, err := setupBasicEventStream(t, payload)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := bedrockruntime.NewFromConfig(cfg)
	resp, err := svc.InvokeModelWithBidirectionalStream(context.Background(),
		&bedrockruntime.InvokeModelWithBidirectionalStreamInput{ModelId: aws.String(ModelID)}, func(options *bedrockruntime.Options) {
			options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
		})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	// Assert calling Err before close does not close the stream.
	resp.GetStream().Err()
	<-resp.GetInitialReply()
	select {
	case _, ok := <-resp.GetStream().Events():
		if !ok {
			t.Fatalf("expect stream not to be closed, but was")
		}
	default:
	}

	resp.GetStream().Close()
	<-resp.GetStream().Events()

	if err := resp.GetStream().Err(); err != nil {
		t.Errorf("expect no error, %v", err)
	}
}

func setupBasicEventStream(t *testing.T, payload string) (aws.Config, func(), error) {
	return eventstreamtesting.SetupEventStream(t,
		&eventstreamtesting.ServeEventStream{
			T: t,
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("chunk"),
						},
						{
							Name:  eventstreamapi.ContentTypeHeader,
							Value: eventstream.StringValue("application/json"),
						},
					},
					Payload: []byte(fmt.Sprintf(`{"bytes":"%s"}`, toBase64(payload))),
				},
			},
			BiDirectional: true,
		},
	)
}

func TestStartStreamTranscription_ReadUnknownEvent(t *testing.T) {
	payload := `{}`
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		&eventstreamtesting.ServeEventStream{
			T: t,
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("chunk"),
						},
						{
							Name:  eventstreamapi.ContentTypeHeader,
							Value: eventstream.StringValue("application/json"),
						},
					},
					Payload: []byte(fmt.Sprintf(`{"bytes":"%s"}`, toBase64(payload))),
				},
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("UnknownEventName"),
						},
					},
					Payload: []byte(`{}`),
				},
			},
			BiDirectional: true,
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := bedrockruntime.NewFromConfig(cfg)
	resp, err := svc.InvokeModelWithBidirectionalStream(context.Background(),
		&bedrockruntime.InvokeModelWithBidirectionalStreamInput{ModelId: aws.String(ModelID)}, func(options *bedrockruntime.Options) {
			options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
		})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()
	<-resp.GetInitialReply()

	expectEvents := []types.InvokeModelWithBidirectionalStreamOutput{
		&types.InvokeModelWithBidirectionalStreamOutputMemberChunk{
			Value: types.BidirectionalOutputPayloadPart{Bytes: []byte(payload)},
		},
		&types.UnknownUnionMember{Tag: "UnknownEventName", Value: func() []byte {
			encoder := eventstream.NewEncoder()
			buff := bytes.NewBuffer(nil)
			encoder.Encode(buff, eventstream.Message{
				Headers: eventstream.Headers{
					eventstreamtesting.EventMessageTypeHeader,
					{
						Name:  eventstreamapi.EventTypeHeader,
						Value: eventstream.StringValue("UnknownEventName"),
					},
				},
				Payload: []byte(`{}`)})
			return buff.Bytes()
		}()},
	}

	for i := 0; i < len(expectEvents); i++ {
		event := <-resp.GetStream().Events()
		if event == nil {
			t.Errorf("%d, expect event, got nil", i)
		}
		if diff := cmpDiff(expectEvents[i], event); len(diff) > 0 {
			t.Errorf("expected %v as %d event, got %v", expectEvents[i], i, event)
		}
	}

	if err := resp.GetStream().Err(); err != nil {
		t.Errorf("expect no error, %v", err)
	}
}

func TestStartStreamTranscription_ReadException(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		&eventstreamtesting.ServeEventStream{
			T: t,
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventExceptionTypeHeader,
						{
							Name:  eventstreamapi.ExceptionTypeHeader,
							Value: eventstream.StringValue("ValidationException"),
						},
					},
					Payload: []byte(`{
 "Message": "Unable to parse input chunk. Please check input format contains correct format."
}`),
				},
			},
			BiDirectional: true,
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := bedrockruntime.NewFromConfig(cfg)
	resp, err := svc.InvokeModelWithBidirectionalStream(context.Background(),
		&bedrockruntime.InvokeModelWithBidirectionalStreamInput{ModelId: aws.String(ModelID)}, func(options *bedrockruntime.Options) {
			options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
		})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()
	<-resp.GetInitialReply()

	<-resp.GetStream().Events()

	err = resp.GetStream().Err()
	if err == nil {
		t.Fatalf("expect err, got none")
	}

	var expectedErr *types.ValidationException
	if !errors.As(err, &expectedErr) {
		t.Errorf("expect err type %T", expectedErr)
	}

	if diff := cmpDiff(
		expectedErr,
		&types.ValidationException{Message: aws.String("Unable to parse input chunk. Please check input format contains correct format.")},
	); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestStartStreamTranscription_ReadUnmodeledException(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		&eventstreamtesting.ServeEventStream{
			T: t,
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventExceptionTypeHeader,
						{
							Name:  eventstreamapi.ExceptionTypeHeader,
							Value: eventstream.StringValue("UnmodeledException"),
						},
					},
					Payload: []byte(`{
 "Message": "this is an unmodeled exception message"
}`),
				},
			},
			BiDirectional: true,
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := bedrockruntime.NewFromConfig(cfg)
	resp, err := svc.InvokeModelWithBidirectionalStream(context.Background(),
		&bedrockruntime.InvokeModelWithBidirectionalStreamInput{ModelId: aws.String(ModelID)}, func(options *bedrockruntime.Options) {
			options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
		})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()
	<-resp.GetInitialReply()

	<-resp.GetStream().Events()

	err = resp.GetStream().Err()
	if err == nil {
		t.Fatalf("expect err, got none")
	}

	var expectedErr *smithy.GenericAPIError
	if !errors.As(err, &expectedErr) {
		t.Errorf("expect err type %T", expectedErr)
	}

	if diff := cmpDiff(
		expectedErr,
		&smithy.GenericAPIError{
			Code:    "UnmodeledException",
			Message: "this is an unmodeled exception message",
		},
	); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestStartStreamTranscription_ReadErrorEvent(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		&eventstreamtesting.ServeEventStream{
			T: t,
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						{
							Name:  eventstreamapi.MessageTypeHeader,
							Value: eventstream.StringValue(eventstreamapi.ErrorMessageType),
						},
						{
							Name:  eventstreamapi.ErrorCodeHeader,
							Value: eventstream.StringValue("AnErrorCode"),
						},
						{
							Name:  eventstreamapi.ErrorMessageHeader,
							Value: eventstream.StringValue("An error message"),
						},
					},
				},
			},
			BiDirectional: true,
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := bedrockruntime.NewFromConfig(cfg)
	resp, err := svc.InvokeModelWithBidirectionalStream(context.Background(),
		&bedrockruntime.InvokeModelWithBidirectionalStreamInput{ModelId: aws.String(ModelID)}, func(options *bedrockruntime.Options) {
			options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
		})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()
	<-resp.GetInitialReply()

	<-resp.GetStream().Events()

	err = resp.GetStream().Err()
	if err == nil {
		t.Fatalf("expect err, got none")
	}

	var expectedErr *smithy.GenericAPIError
	if !errors.As(err, &expectedErr) {
		t.Errorf("expect err type %T", expectedErr)
	}

	if diff := cmpDiff(
		expectedErr,
		&smithy.GenericAPIError{
			Code:    "AnErrorCode",
			Message: "An error message",
		},
	); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestStartStreamTranscription_ReadWrite(t *testing.T) {
	payload := `{}`
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		&eventstreamtesting.ServeEventStream{
			T: t,
			ClientEvents: []eventstream.Message{
				{
					// NOTE deepEqual comparison forces the events to have this specific arrangement,
					// but it's semantically equivalent if type was above say
					Headers: eventstream.Headers{
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("chunk"),
						},
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.ContentTypeHeader,
							Value: eventstream.StringValue("application/json"),
						},
					},
					Payload: []byte(fmt.Sprintf(`{"bytes":"%s"}`, toBase64(payload))),
				},
			},
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("chunk"),
						},
						{
							Name:  eventstreamapi.ContentTypeHeader,
							Value: eventstream.StringValue("application/json"),
						},
					},
					Payload: []byte(fmt.Sprintf(`{"bytes":"%s"}`, toBase64(payload))),
				},
			},
			BiDirectional: true,
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := bedrockruntime.NewFromConfig(cfg)
	resp, err := svc.InvokeModelWithBidirectionalStream(context.Background(),
		&bedrockruntime.InvokeModelWithBidirectionalStreamInput{ModelId: aws.String(ModelID)}, func(options *bedrockruntime.Options) {
			options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
		})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	stream := resp.GetStream()
	defer stream.Close()
	<-resp.GetInitialReply()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		expectedServiceEvents := []types.InvokeModelWithBidirectionalStreamOutput{
			&types.InvokeModelWithBidirectionalStreamOutputMemberChunk{
				Value: types.BidirectionalOutputPayloadPart{Bytes: []byte(payload)},
			},
		}
		for i := 0; i < len(expectedServiceEvents); i++ {
			event := <-resp.GetStream().Events()
			if event == nil {
				t.Errorf("%d, expect event, got nil", i)
			}
			if e, a := expectedServiceEvents[i], event; !reflect.DeepEqual(e, a) {
				t.Errorf("%d, expect %T %v, got %T %v", i, e, e, a, a)
			}
		}
	}()

	clientEvents := []types.InvokeModelWithBidirectionalStreamInput{
		&types.InvokeModelWithBidirectionalStreamInputMemberChunk{Value: types.BidirectionalInputPayloadPart{Bytes: []byte(payload)}},
	}
	for _, event := range clientEvents {
		err = stream.Send(context.Background(), event)
		if err != nil {
			t.Errorf("expect no error, got %v", err)
		}
	}

	wg.Wait()

	stream.Close()

	if err := stream.Err(); err != nil {
		t.Errorf("expect no error, %v", err)
	}
}

func TestInvokeModelWithBidirectionalStream_Write(t *testing.T) {
	payload := `{}`
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		&eventstreamtesting.ServeEventStream{
			T: t,
			ClientEvents: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("chunk"),
						},
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.ContentTypeHeader,
							Value: eventstream.StringValue("application/json"),
						},
					},
					Payload: []byte(fmt.Sprintf(`{"bytes":"%s"}`, toBase64(payload))),
				},
			},
			BiDirectional: true,
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := bedrockruntime.NewFromConfig(cfg)
	resp, err := svc.InvokeModelWithBidirectionalStream(context.Background(), &bedrockruntime.InvokeModelWithBidirectionalStreamInput{ModelId: aws.String(ModelID)}, func(options *bedrockruntime.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	stream := resp.GetStream()

	clientEvents := []types.InvokeModelWithBidirectionalStreamInput{
		&types.InvokeModelWithBidirectionalStreamInputMemberChunk{Value: types.BidirectionalInputPayloadPart{Bytes: []byte(payload)}},
	}
	for _, event := range clientEvents {
		err = stream.Send(context.Background(), event)
		if err != nil {
			t.Errorf("expect no error, got %v", err)
		}
	}

	if err := stream.Close(); err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestInvokeModelWithBidirectionalStream_WriteClose(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		&eventstreamtesting.ServeEventStream{
			T:             t,
			BiDirectional: true,
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := bedrockruntime.NewFromConfig(cfg)
	resp, err := svc.InvokeModelWithBidirectionalStream(context.Background(), &bedrockruntime.InvokeModelWithBidirectionalStreamInput{ModelId: aws.String(ModelID)}, func(options *bedrockruntime.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	stream := resp.GetStream()

	// Assert calling Err before close does not close the stream.
	stream.Err()

	err = stream.Send(context.Background(), &types.InvokeModelWithBidirectionalStreamInputMemberChunk{Value: types.BidirectionalInputPayloadPart{Bytes: []byte{}}})
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	stream.Close()

	if err := stream.Err(); err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestInvokeModelWithBidirectionalStream_WriteError(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		&eventstreamtesting.ServeEventStream{
			T:               t,
			BiDirectional:   true,
			ForceCloseAfter: time.Millisecond * 500,
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := bedrockruntime.NewFromConfig(cfg)
	resp, err := svc.InvokeModelWithBidirectionalStream(context.Background(), &bedrockruntime.InvokeModelWithBidirectionalStreamInput{ModelId: aws.String(ModelID)}, func(options *bedrockruntime.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	stream := resp.GetStream()
	defer stream.Close()

	for {
		err = resp.GetStream().Send(context.Background(), &types.InvokeModelWithBidirectionalStreamInputMemberChunk{Value: types.BidirectionalInputPayloadPart{}})
		if err != nil {
			if strings.Contains("unable to send event", err.Error()) {
				t.Errorf("expected stream closed error, got %v", err)
			}
			break
		}
	}
}

func TestInvokeModelWithBidirectionalStream_ResponseError(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		&eventstreamtesting.ServeEventStream{
			T: t,
			StaticResponse: &eventstreamtesting.StaticResponse{
				StatusCode: 500,
				Body: []byte(`{
					"Message": "this is an exception message"
				}`),
			},
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := bedrockruntime.NewFromConfig(cfg)
	resp, err := svc.InvokeModelWithBidirectionalStream(
		context.Background(),
		&bedrockruntime.InvokeModelWithBidirectionalStreamInput{ModelId: aws.String(ModelID)},
		func(options *bedrockruntime.Options) {
			options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
		},
	)
	if err != nil {
		t.Fatalf("an error on an event stream shouldn't surface on the stream, got %v", err)
	}
	stream := resp.GetStream()
	defer stream.Close()

	<-resp.GetInitialReply()
	err = stream.Err()
	if err == nil {
		t.Fatal("expect an error, got nil")
	}

	var expectedErr *smithy.GenericAPIError
	if !errors.As(err, &expectedErr) {
		t.Errorf("expect err type %T, got %v", expectedErr, err)
	}

	if diff := cmpDiff(
		expectedErr,
		&smithy.GenericAPIError{
			Code:    "UnknownError",
			Message: "this is an exception message",
		},
	); len(diff) > 0 {
		t.Error(diff)
	}
}

func cmpDiff(e, a interface{}) string {
	if !reflect.DeepEqual(e, a) {
		return fmt.Sprintf("%v != %v", e, a)
	}
	return ""
}
