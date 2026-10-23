package testing

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream/eventstreamapi"
	"github.com/aws/aws-sdk-go-v2/service/internal/eventstreamtesting"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
)

func removeValidationMiddleware(stack *middleware.Stack) error {
	_, err := stack.Initialize.Remove("OperationInputValidation")
	return err
}

func TestSubscribeToShard_Read(t *testing.T) {
	stream := &eventstreamtesting.ServeEventStream{T: t}
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t, stream)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(cfg)

	// Encode payloads in whatever wire format the client's resolved protocol decodes
	f := clientWireFormat(svc)
	stream.Events = []eventstream.Message{
		{
			Headers: eventstream.Headers{
				eventstreamtesting.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("initial-response"),
				},
			},
			Payload: emptyPayload(f),
		},
		{
			Headers: eventstream.Headers{
				eventstreamtesting.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("SubscribeToShardEvent"),
				},
			},
			Payload: objectPayload(f, map[string]string{
				"ContinuationSequenceNumber": "01234",
			}),
		},
	}

	resp, err := svc.SubscribeToShard(context.Background(),
		&kinesis.SubscribeToShardInput{}, func(options *kinesis.Options) {
			options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
		})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()

	expectEvents := []types.SubscribeToShardEventStream{
		&types.SubscribeToShardEventStreamMemberSubscribeToShardEvent{
			Value: types.SubscribeToShardEvent{ContinuationSequenceNumber: aws.String("01234")},
		},
	}

	for i := range expectEvents {
		event := <-resp.GetStream().Events()
		if event == nil {
			t.Errorf("%d, expect event, got nil", i)
		}
		if diff := cmpDiff(expectEvents[i], event); len(diff) > 0 {
			t.Errorf("%d, %v", i, diff)
		}
	}

	if err := resp.GetStream().Err(); err != nil {
		t.Errorf("expect no error, %v", err)
	}
}

func TestSubscribeToShard_ReadClose(t *testing.T) {
	stream := &eventstreamtesting.ServeEventStream{T: t}
	sess, cleanupFn, err := eventstreamtesting.SetupEventStream(t, stream)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(sess)

	f := clientWireFormat(svc)
	stream.Events = []eventstream.Message{
		{
			Headers: eventstream.Headers{
				eventstreamtesting.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("initial-response"),
				},
			},
			Payload: emptyPayload(f),
		},
		{
			Headers: eventstream.Headers{
				eventstreamtesting.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("SubscribeToShardEvent"),
				},
			},
			Payload: objectPayload(f, map[string]string{
				"ContinuationSequenceNumber": "01234",
			}),
		},
	}

	resp, err := svc.SubscribeToShard(context.Background(), &kinesis.SubscribeToShardInput{}, func(options *kinesis.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	// Assert calling Err before close does not close the stream.
	resp.GetStream().Err()
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

func TestSubscribeToShard_ReadUnknownEvent(t *testing.T) {
	stream := &eventstreamtesting.ServeEventStream{T: t}
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t, stream)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(cfg)

	f := clientWireFormat(svc)
	stream.Events = []eventstream.Message{
		{
			Headers: eventstream.Headers{
				eventstreamtesting.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("initial-response"),
				},
			},
			Payload: emptyPayload(f),
		},
		{
			Headers: eventstream.Headers{
				eventstreamtesting.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("SubscribeToShardEvent"),
				},
			},
			Payload: objectPayload(f, map[string]string{
				"ContinuationSequenceNumber": "01234",
			}),
		},
		{
			Headers: eventstream.Headers{
				eventstreamtesting.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("UnknownEventName"),
				},
			},
			Payload: emptyPayload(f),
		},
	}

	resp, err := svc.SubscribeToShard(context.Background(), &kinesis.SubscribeToShardInput{}, func(options *kinesis.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()

	expectEvents := []types.SubscribeToShardEventStream{
		&types.SubscribeToShardEventStreamMemberSubscribeToShardEvent{
			Value: types.SubscribeToShardEvent{ContinuationSequenceNumber: aws.String("01234")},
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
				Payload: emptyPayload(f)})
			return buff.Bytes()
		}()},
	}

	for i := range expectEvents {
		event := <-resp.GetStream().Events()
		if event == nil {
			t.Errorf("%d, expect event, got nil", i)
		}
		if diff := cmpDiff(expectEvents[i], event); len(diff) > 0 {
			t.Errorf("%d, %v", i, diff)
		}
	}

	if err := resp.GetStream().Err(); err != nil {
		t.Errorf("expect no error, %v", err)
	}
}

func TestSubscribeToShard_ReadException(t *testing.T) {
	stream := &eventstreamtesting.ServeEventStream{T: t}
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t, stream)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(cfg)

	f := clientWireFormat(svc)
	stream.Events = []eventstream.Message{
		{
			Headers: eventstream.Headers{
				eventstreamtesting.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("initial-response"),
				},
			},
			Payload: emptyPayload(f),
		},
		{
			Headers: eventstream.Headers{
				eventstreamtesting.EventExceptionTypeHeader,
				{
					Name:  eventstreamapi.ExceptionTypeHeader,
					Value: eventstream.StringValue("InternalFailureException"),
				},
			},
			Payload: objectPayload(f, map[string]string{
				"message": "this is an exception message",
			}),
		},
	}

	resp, err := svc.SubscribeToShard(context.Background(), &kinesis.SubscribeToShardInput{}, func(options *kinesis.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()

	<-resp.GetStream().Events()

	err = resp.GetStream().Err()
	if err == nil {
		t.Fatalf("expect err, got none")
	}

	var expectedErr *types.InternalFailureException
	if !errors.As(err, &expectedErr) {
		t.Errorf("expect err type %T", expectedErr)
	}

	if diff := cmpDiff(
		expectedErr,
		&types.InternalFailureException{Message: aws.String("this is an exception message")},
	); len(diff) > 0 {
		t.Error(diff)
	}
}

func TestSubscribeToShard_ReadUnmodeledException(t *testing.T) {
	stream := &eventstreamtesting.ServeEventStream{T: t}
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t, stream)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(cfg)

	f := clientWireFormat(svc)
	stream.Events = []eventstream.Message{
		{
			Headers: eventstream.Headers{
				eventstreamtesting.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("initial-response"),
				},
			},
			Payload: emptyPayload(f),
		},
		{
			Headers: eventstream.Headers{
				eventstreamtesting.EventExceptionTypeHeader,
				{
					Name:  eventstreamapi.ExceptionTypeHeader,
					Value: eventstream.StringValue("UnmodeledException"),
				},
			},
			Payload: objectPayload(f, map[string]string{
				"message": "this is an unmodeled exception message",
			}),
		},
	}

	resp, err := svc.SubscribeToShard(context.Background(), &kinesis.SubscribeToShardInput{}, func(options *kinesis.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()

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

func TestSubscribeToShard_ReadErrorEvent(t *testing.T) {
	stream := &eventstreamtesting.ServeEventStream{T: t}
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t, stream)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(cfg)

	f := clientWireFormat(svc)
	stream.Events = []eventstream.Message{
		{
			Headers: eventstream.Headers{
				eventstreamtesting.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("initial-response"),
				},
			},
			Payload: emptyPayload(f),
		},
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
	}

	resp, err := svc.SubscribeToShard(context.Background(), &kinesis.SubscribeToShardInput{}, func(options *kinesis.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()

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

func TestSubscribeToShard_ResponseError(t *testing.T) {
	stream := &eventstreamtesting.ServeEventStream{
		T:              t,
		StaticResponse: &eventstreamtesting.StaticResponse{StatusCode: 500},
	}
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t, stream)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(cfg)

	f := clientWireFormat(svc)
	stream.StaticResponse.Body = objectPayload(f, map[string]string{
		"message": "this is an exception message",
	})

	_, err = svc.SubscribeToShard(context.Background(), &kinesis.SubscribeToShardInput{}, func(options *kinesis.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err == nil {
		t.Fatal("expect error got nil")
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

func cmpDiff(e, a any) string {
	if !reflect.DeepEqual(e, a) {
		return fmt.Sprintf("%v != %v", e, a)
	}
	return ""
}
