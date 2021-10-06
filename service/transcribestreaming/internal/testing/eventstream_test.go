package testing

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream/eventstreamapi"
	"github.com/aws/aws-sdk-go-v2/service/internal/eventstreamtesting"
	"github.com/aws/aws-sdk-go-v2/service/transcribestreaming"
	"github.com/aws/aws-sdk-go-v2/service/transcribestreaming/types"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/document"
	"github.com/aws/smithy-go/middleware"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
							Value: eventstream.StringValue("TranscriptEvent"),
						},
					},
					Payload: []byte(`{
  "Transcript": {
    "Results": []
  }
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

	svc := transcribestreaming.NewFromConfig(cfg)
	resp, err := svc.StartStreamTranscription(context.Background(),
		&transcribestreaming.StartStreamTranscriptionInput{}, func(options *transcribestreaming.Options) {
			options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
		})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()

	expectEvents := []types.TranscriptResultStream{
		&types.TranscriptResultStreamMemberTranscriptEvent{
			Value: types.TranscriptEvent{Transcript: &types.Transcript{Results: []types.Result{}}},
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
							Value: eventstream.StringValue("TranscriptEvent"),
						},
					},
					Payload: []byte(`{
  "Transcript": {
    "Results": []
  }
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

	svc := transcribestreaming.NewFromConfig(sess)
	resp, err := svc.StartStreamTranscription(context.Background(), &transcribestreaming.StartStreamTranscriptionInput{}, func(options *transcribestreaming.Options) {
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
							Value: eventstream.StringValue("TranscriptEvent"),
						},
					},
					Payload: []byte(`{
  "Transcript": {
    "Results": []
  }
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
			BiDirectional: true,
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := transcribestreaming.NewFromConfig(cfg)
	resp, err := svc.StartStreamTranscription(context.Background(), &transcribestreaming.StartStreamTranscriptionInput{}, func(options *transcribestreaming.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}
	defer resp.GetStream().Close()

	expectEvents := []types.TranscriptResultStream{
		&types.TranscriptResultStreamMemberTranscriptEvent{Value: types.TranscriptEvent{Transcript: &types.Transcript{Results: []types.Result{}}}},
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
						eventstreamtesting.EventExceptionTypeHeader,
						{
							Name:  eventstreamapi.ExceptionTypeHeader,
							Value: eventstream.StringValue("BadRequestException"),
						},
					},
					Payload: []byte(`{
  "Message": "this is an exception message"
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

	svc := transcribestreaming.NewFromConfig(cfg)
	resp, err := svc.StartStreamTranscription(context.Background(), &transcribestreaming.StartStreamTranscriptionInput{}, func(options *transcribestreaming.Options) {
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

	var expectedErr *types.BadRequestException
	if !errors.As(err, &expectedErr) {
		t.Errorf("expect err type %T", expectedErr)
	}

	if diff := cmp.Diff(
		expectedErr,
		&types.BadRequestException{Message: aws.String("this is an exception message")},
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

	svc := transcribestreaming.NewFromConfig(cfg)
	resp, err := svc.StartStreamTranscription(context.Background(), &transcribestreaming.StartStreamTranscriptionInput{}, func(options *transcribestreaming.Options) {
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

	svc := transcribestreaming.NewFromConfig(cfg)
	resp, err := svc.StartStreamTranscription(context.Background(), &transcribestreaming.StartStreamTranscriptionInput{}, func(options *transcribestreaming.Options) {
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

func TestStartStreamTranscription_ReadWrite(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		eventstreamtesting.ServeEventStream{
			T: t,
			ClientEvents: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("AudioEvent"),
						},
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.ContentTypeHeader,
							Value: eventstream.StringValue("application/octet-stream"),
						},
					},
					Payload: []byte{0x1, 0x2, 0x3, 0x4},
				},
			},
			Events: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("TranscriptEvent"),
						},
					},
					Payload: []byte(`{
  "Transcript": {
    "Results": []
  }
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

	svc := transcribestreaming.NewFromConfig(cfg)
	resp, err := svc.StartStreamTranscription(context.Background(), &transcribestreaming.StartStreamTranscriptionInput{}, func(options *transcribestreaming.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	stream := resp.GetStream()
	defer stream.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		expectedServiceEvents := []types.TranscriptResultStream{
			&types.TranscriptResultStreamMemberTranscriptEvent{Value: types.TranscriptEvent{Transcript: &types.Transcript{Results: []types.Result{}}}},
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

	clientEvents := []types.AudioStream{
		&types.AudioStreamMemberAudioEvent{Value: types.AudioEvent{AudioChunk: []byte{0x1, 0x2, 0x3, 0x4}}},
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

func TestStartStreamTranscription_Write(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		eventstreamtesting.ServeEventStream{
			T: t,
			ClientEvents: []eventstream.Message{
				{
					Headers: eventstream.Headers{
						{
							Name:  eventstreamapi.EventTypeHeader,
							Value: eventstream.StringValue("AudioEvent"),
						},
						eventstreamtesting.EventMessageTypeHeader,
						{
							Name:  eventstreamapi.ContentTypeHeader,
							Value: eventstream.StringValue("application/octet-stream"),
						},
					},
					Payload: []byte{0x1, 0x2, 0x3, 0x4},
				},
			},
			BiDirectional: true,
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := transcribestreaming.NewFromConfig(cfg)
	resp, err := svc.StartStreamTranscription(context.Background(), &transcribestreaming.StartStreamTranscriptionInput{}, func(options *transcribestreaming.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	stream := resp.GetStream()

	clientEvents := []types.AudioStream{
		&types.AudioStreamMemberAudioEvent{Value: types.AudioEvent{AudioChunk: []byte{0x1, 0x2, 0x3, 0x4}}},
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

func TestStartStreamTranscription_WriteClose(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		eventstreamtesting.ServeEventStream{
			T:             t,
			BiDirectional: true,
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := transcribestreaming.NewFromConfig(cfg)
	resp, err := svc.StartStreamTranscription(context.Background(), &transcribestreaming.StartStreamTranscriptionInput{}, func(options *transcribestreaming.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	stream := resp.GetStream()

	// Assert calling Err before close does not close the stream.
	stream.Err()

	err = stream.Send(context.Background(), &types.AudioStreamMemberAudioEvent{Value: types.AudioEvent{AudioChunk: []byte{}}})
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	stream.Close()

	if err := stream.Err(); err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestStartStreamTranscription_WriteError(t *testing.T) {
	cfg, cleanupFn, err := eventstreamtesting.SetupEventStream(t,
		eventstreamtesting.ServeEventStream{
			T:               t,
			BiDirectional:   true,
			ForceCloseAfter: time.Millisecond * 500,
		},
	)
	if err != nil {
		t.Fatalf("expect no error, %v", err)
	}
	defer cleanupFn()

	svc := transcribestreaming.NewFromConfig(cfg)
	resp, err := svc.StartStreamTranscription(context.Background(), &transcribestreaming.StartStreamTranscriptionInput{}, func(options *transcribestreaming.Options) {
		options.APIOptions = append(options.APIOptions, removeValidationMiddleware)
	})
	if err != nil {
		t.Fatalf("expect no error got, %v", err)
	}

	stream := resp.GetStream()
	defer stream.Close()

	for {
		err = resp.GetStream().Send(context.Background(), &types.AudioStreamMemberAudioEvent{Value: types.AudioEvent{}})
		if err != nil {
			if strings.Contains("unable to send event", err.Error()) {
				t.Errorf("expected stream closed error, got %v", err)
			}
			break
		}
	}
}
