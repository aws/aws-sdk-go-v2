package testing

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream/eventstreamapi"
	"github.com/aws/aws-sdk-go-v2/service/internal/eventstreamtesting"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/document"
	"github.com/aws/smithy-go/middleware"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
)

func removeValidationMiddleware(stack *middleware.Stack) error {
	_, err := stack.Initialize.Remove("OperationInputValidation")
	return err
}

func TestStartStreamTranscription_Read(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		eventstreamtesting.ServeEventStream{
			T: t,
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("initial-response"),
						},
					},
					Payload: []byte(`{}`),
				},
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("SubscribeToShardEvent"),
						},
					},
					Payload: []byte(`{
  "ContinuationSequenceNumber": "01234"
}`),
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(cfg)
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

	for i := 0; i < len(expectEvents); i++ {
		event := <-resp.GetStream().Events()
		if event == nil {
			t.Errorf("%d, expect event, got nil", i)
		}
		if diff := cmp.Diff(expectEvents[i], event, cmpopts.IgnoreTypes(document.NoSerde{})); len(diff) > 0 {
			t.Errorf("%d, %v", i, diff)
		}
	}

	if err := resp.GetStream().Err(); err != nil {
		t.Errorf("expect no error, %v", err)
	}
}

func TestStartStreamTranscription_ReadClose(t *testing.T) {
	sess, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		eventstreamtesting.ServeEventStream{
			T: t,
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("initial-response"),
						},
					},
					Payload: []byte(`{}`),
				},
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("SubscribeToShardEvent"),
						},
					},
					Payload: []byte(`{
  "ContinuationSequenceNumber": "01234"
}`),
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(sess)
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

func TestStartStreamTranscription_ReadUnknownEvent(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		eventstreamtesting.ServeEventStream{
			T: t,
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("initial-response"),
						},
					},
					Payload: []byte(`{}`),
				},
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("SubscribeToShardEvent"),
						},
					},
					Payload: []byte(`{
  "ContinuationSequenceNumber": "01234"
}`),
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
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(cfg)
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
				Payload: []byte(`{}`)})
			return buff.Bytes()
		}()},
	}

	for i := 0; i < len(expectEvents); i++ {
		event := <-resp.GetStream().Events()
		if event == nil {
			t.Errorf("%d, expect event, got nil", i)
		}
		if diff := cmp.Diff(expectEvents[i], event, cmpopts.IgnoreTypes(document.NoSerde{})); len(diff) > 0 {
			t.Errorf("%d, %v", i, diff)
		}
	}

	if err := resp.GetStream().Err(); err != nil {
		t.Errorf("expect no error, %v", err)
	}
}

func TestStartStreamTranscription_ReadException(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		eventstreamtesting.ServeEventStream{
			T: t,
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("initial-response"),
						},
					},
					Payload: []byte(`{}`),
				},
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventExceptionTypeHeader,
						{
							Name:  eventstreamapi.ExceptionTypeHeader,
							Value: eventstream.StringValue("InternalFailureException"),
						},
					},
					Payload: []byte(`{
  "message": "this is an exception message"
}`),
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(cfg)
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

	if diff := cmp.Diff(
		expectedErr,
		&types.InternalFailureException{Message: aws.String("this is an exception message")},
		cmpopts.IgnoreTypes(document.NoSerde{}),
	); len(diff) > 0 {
		t.Errorf(diff)
	}
}

func TestStartStreamTranscription_ReadUnmodeledException(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		eventstreamtesting.ServeEventStream{
			T: t,
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("initial-response"),
						},
					},
					Payload: []byte(`{}`),
				},
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
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(cfg)
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

	if diff := cmp.Diff(
		expectedErr,
		&smithy.GenericAPIError{
			Code:    "UnmodeledException",
			Message: "this is an unmodeled exception message",
		},
		cmpopts.IgnoreTypes(document.NoSerde{}),
	); len(diff) > 0 {
		t.Errorf(diff)
	}
}

func TestStartStreamTranscription_ReadErrorEvent(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		eventstreamtesting.ServeEventStream{
			T: t,
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("initial-response"),
						},
					},
					Payload: []byte(`{}`),
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
			},
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := kinesis.NewFromConfig(cfg)
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

	if diff := cmp.Diff(
		expectedErr,
		&smithy.GenericAPIError{
			Code:    "AnErrorCode",
			Message: "An error message",
		},
		cmpopts.IgnoreTypes(document.NoSerde{}),
	); len(diff) > 0 {
		t.Errorf(diff)
	}
}
