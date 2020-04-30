package lexruntimeservice

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdkio"
	"github.com/aws/aws-sdk-go-v2/service/smithyprototype/lexruntimeservice/types"
	"github.com/awslabs/smithy-go/middleware"
	smithyHTTP "github.com/awslabs/smithy-go/transport/http"
)

type awsRestJson1_deserializeOpPutSession struct{}

// ID is the middleware identifier
func (p awsRestJson1_deserializeOpPutSession) ID() string {
	return "awsRestJson1_deserializeOpPutSession"
}

func (p awsRestJson1_deserializeOpPutSession) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err != nil {
		return out, metadata, err
	}
	response, ok := out.RawResponse.(*smithyHTTP.Response)
	if !ok {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("unknown transport type %T", out.RawResponse),
		}
	}

	output, ok := out.Result.(*PutSessionOutput)
	if !ok {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("unknown output result type %T", out.Result),
		}
	}

	if err := awsRestJson1_deserializeRestPutSessionOutput(output, response); err != nil {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("failed to decode response with invalid REST bindings, %w", err),
		}
	}

	// TODO: add error response deserializer function

	// body of type blob
	if err := awsRestJson1_deserializeJsonPutSessionOutput(output, response.Body); err != nil {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("failed to decode response body with invalid JSON, %w", err),
		}
	}

	return out, metadata, err
}

func awsRestJson1_deserializeRestPutSessionOutput(v *PutSessionOutput, response *smithyHTTP.Response) error {
	if v == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", v)
	}

	if v.ContentType != nil {
		val := response.Header.Get("Content-Type")
		v.ContentType = &val
	}

	if len(v.DialogState) > 0 {
		val := response.Header.Get("x-amz-lex-dialog-state")
		v.DialogState = types.DialogState(val)
	}

	if v.IntentName != nil {
		val := response.Header.Get("x-amz-lex-intent-name")
		v.IntentName = &val
	}

	if v.Message != nil {
		val := response.Header.Get("x-amz-lex-message")
		v.Message = &val
	}

	if len(v.MessageFormat) > 0 {
		val := response.Header.Get("x-amz-lex-message-format")
		v.MessageFormat = types.MessageFormatType(val)
	}

	if v.SessionAttributes != nil {
		val := response.Header.Get("x-amz-lex-session-attributes")
		b, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return err
		}
		m := aws.JSONValue{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return err
		}
		v.SessionAttributes = m
	}

	if v.SessionId != nil {
		val := response.Header.Get("x-amz-lex-session-id")
		v.SessionId = &val
	}

	if v.SlotToElicit != nil {
		val := response.Header.Get("x-amz-lex-slot-to-elicit")
		v.SlotToElicit = &val
	}

	if v.Slots != nil {
		val := response.Header.Get("x-amz-lex-slots")
		b, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return err
		}
		m := aws.JSONValue{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return err
		}
		v.Slots = m
	}

	return nil
}

func awsRestJson1_deserializeJsonPutSessionOutput(v *PutSessionOutput, body io.ReadCloser) error {
	if v == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", v)
	}

	v.AudioStream = body
	return nil
}

type awsRestJson1_deserializeOpPostText struct{}

// ID is the middleware identifier
func (p awsRestJson1_deserializeOpPostText) ID() string {
	return "awsRestJson1_deserializeOpPostText"
}

func (p awsRestJson1_deserializeOpPostText) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err != nil {
		return out, metadata, err
	}
	response, ok := out.RawResponse.(*smithyHTTP.Response)
	if !ok {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("unknown transport type %T", out.RawResponse),
		}
	}

	// initialize the out.Result with PostTextOutput
	// This step will overwrite any value that may already
	// be in out.Result
	out.Result = &PostTextOutput{}
	output, ok := out.Result.(*PostTextOutput)
	if !ok {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("unknown output result type %T", out.Result),
		}
	}

	buff := make([]byte, 1024)
	ringBuffer := sdkio.NewRingBuffer(buff)

	// wrap a TeeReader to read from response body & write on snapshot
	body := io.TeeReader(response.Body, ringBuffer)
	defer response.Body.Close()
	decoder := json.NewDecoder(body)
	// UseNumber causes the Decoder to unmarshal a number into an interface{}
	// as a Number instead of as a float64.
	decoder.UseNumber()

	if err := awsRestJson1_deserializeJsonPostTextOutput(output, decoder); err != nil {
		snapshot := make([]byte, 1024)
		ringBuffer.Read(snapshot)
		return out, metadata, &aws.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body with invalid JSON, %w", err),
			Snapshot: snapshot,
		}
	}

	return out, metadata, err
}

func awsRestJson1_deserializeJsonPostTextOutput(output *PostTextOutput, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		// "Empty Response"
		return nil
	}
	if err != nil {
		return err
	}

	if t, ok := startToken.(json.Delim); !ok {
		if t.String() != "{" {
			return fmt.Errorf("expected `{` as start token")
		}
	}

	for decoder.More() {
		// fetch token for key
		t, err := decoder.Token()
		if err != nil {
			return err
		}

		// location name : `dialogState` key with value as `string`
		if t == "dialogState" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.DialogState = types.DialogState(v)
			} else {
				return fmt.Errorf("expected DialogState to be of type String, got %T", val)
			}
		}

		// location name: `intentName` key with value as `string`
		if t == "intentName" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.IntentName = &v
			} else {
				return fmt.Errorf("expected IntentName to be of type *String, got %T", val)
			}
		}

		if t == "message" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.Message = &v
			} else {
				return fmt.Errorf("expected Message to be of type *String, got %T", val)
			}
		}

		if t == "messageFormat" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.MessageFormat = types.MessageFormatType(v)
			} else {
				return fmt.Errorf("expected MessageFormat to be of type *String, got %T", val)
			}
		}

		if t == "responseCard" {
			v := types.ResponseCard{}
			if err = awsRestJson1_deserializeJsonResponseCard(&v, decoder); err != nil {
				return err
			}
			output.ResponseCard = &v
		}

		if t == "sentimentResponse" {
			v := types.SentimentResponse{}
			if err = awsRestJson1_deserializeJsonSentimentResponse(&v, decoder); err != nil {
				return err
			}
			output.SentimentResponse = &v
		}

		if t == "sessionAttributes" {
			v := make(map[string]string, 0)
			if err = awsRestJson1_deserializeJsonSessionAttribute(&v, decoder); err != nil {
				return err
			}
			output.SessionAttributes = v
		}

		if t == "sessionId" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.SessionId = &v
			} else {
				return fmt.Errorf("expected SessionId to be of type *String, got %T", val)
			}
		}

		if t == "slotToElicit" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.SlotToElicit = &v
			} else {
				return fmt.Errorf("expected SlotToElicit to be of type *String, got %T", val)
			}
		}

		if t == "slots" {
			v := make(map[string]string, 0)
			if err = awsRestJson1_deserializeJsonSlots(&v, decoder); err != nil {
				return err
			}
			output.Slots = v
		}
	}

	// end of the json response body
	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return fmt.Errorf("expected `}` as end token")
	}

	return nil
}

type awsRestJson1_deserializeOpPostContent struct{}

// ID is the middleware identifier
func (p awsRestJson1_deserializeOpPostContent) ID() string {
	return "awsRestJson1_deserializeOpPostContent"
}

func (p awsRestJson1_deserializeOpPostContent) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err != nil {
		return out, metadata, err
	}
	response, ok := out.RawResponse.(*smithyHTTP.Response)
	if !ok {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("unknown transport type %T", out.RawResponse),
		}
	}

	output, ok := out.Result.(*PostContentOutput)
	if !ok {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("unknown output result type %T", out.Result),
		}
	}

	if err := awsRestJson1_deserializeRestPostContentOutput(output, response); err != nil {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("failed to decode response with invalid REST bindings, %w", err),
		}
	}

	// body of type blob
	if err := awsRestJson1_deserializeJsonPostContentOutput(output, response.Body); err != nil {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("failed to decode response body with invalid JSON, %w", err),
		}
	}

	return out, metadata, err
}

func awsRestJson1_deserializeRestPostContentOutput(v *PostContentOutput, response *smithyHTTP.Response) error {
	if v == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", v)
	}

	if v.ContentType != nil {
		val := response.Header.Get("Content-Type")
		v.ContentType = &val
	}

	if len(v.DialogState) > 0 {
		val := response.Header.Get("x-amz-lex-dialog-state")
		v.DialogState = types.DialogState(val)
	}

	if v.InputTranscript != nil {
		val := response.Header.Get("x-amz-lex-input-transcript")
		v.InputTranscript = &val
	}

	if v.IntentName != nil {
		val := response.Header.Get("x-amz-lex-intent-name")
		v.IntentName = &val
	}

	if v.Message != nil {
		val := response.Header.Get("x-amz-lex-message")
		v.Message = &val
	}

	if len(v.MessageFormat) > 0 {
		val := response.Header.Get("x-amz-lex-message-format")
		v.MessageFormat = types.MessageFormatType(val)
	}

	if v.SentimentResponse != nil {
		val := response.Header.Get("x-amz-lex-sentiment")
		v.SentimentResponse = &val
	}

	if v.SessionAttributes != nil {
		val := response.Header.Get("x-amz-lex-session-attributes")
		b, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return err
		}
		m := aws.JSONValue{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return err
		}
		v.SessionAttributes = m
	}

	if v.SessionId != nil {
		val := response.Header.Get("x-amz-lex-session-id")
		v.SessionId = &val
	}

	if v.SlotToElicit != nil {
		val := response.Header.Get("x-amz-lex-slot-to-elicit")
		v.SlotToElicit = &val
	}

	if v.Slots != nil {
		val := response.Header.Get("x-amz-lex-slots")
		b, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return err
		}
		m := aws.JSONValue{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return err
		}
		v.Slots = m
	}

	return nil
}

func awsRestJson1_deserializeJsonPostContentOutput(v *PostContentOutput, body io.ReadCloser) error {
	if v == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", v)
	}

	v.AudioStream = body

	return nil
}

type awsRestJson1_deserializeOpGetSession struct{}

// ID is the middleware identifier
func (g awsRestJson1_deserializeOpGetSession) ID() string {
	return "awsRestJson1_deserializeOpGetSession"
}

func (g awsRestJson1_deserializeOpGetSession) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput,
	next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err != nil {
		return out, metadata, err
	}
	response, ok := out.RawResponse.(*smithyHTTP.Response)
	if !ok {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("unknown transport type %T", out.RawResponse),
		}
	}

	output, ok := out.Result.(*GetSessionOutput)
	if !ok {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("unknown output result type %T", out.Result),
		}
	}

	buff := make([]byte, 1024)
	ringBuffer := sdkio.NewRingBuffer(buff)

	// wrap a TeeReader to read from response body
	body := io.TeeReader(response.Body, ringBuffer)
	defer response.Body.Close()
	decoder := json.NewDecoder(body)
	// UseNumber causes the Decoder to unmarshal a number into an interface{}
	// as a Number instead of as a float64.
	decoder.UseNumber()

	if err = awsRestJson1_deserializeJsonGetSessionOutput(output, decoder); err != nil {
		snapshot := make([]byte, 1024)
		ringBuffer.Read(snapshot)
		return out, metadata, &aws.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body with invalid JSON, %w", err),
			Snapshot: snapshot,
		}
	}

	return out, metadata, err
}

func awsRestJson1_deserializeJsonGetSessionOutput(output *GetSessionOutput, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		// "Empty Response"
		return nil
	}
	if err != nil {
		return err
	}

	if t, ok := startToken.(json.Delim); !ok {
		if t.String() != "{" {
			return fmt.Errorf("expected `{` as start token ")
		}
	}

	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			return err
		}

		if t == "dialogAction" {
			v := types.DialogAction{}
			if err := awsRestJson1_deserializeJsondialogAction(&v, decoder); err != nil {
				return err
			}
			output.DialogAction = &v
		}

		if t == "recentIntentSummaryView" {
			v := make([]types.IntentSummary, 0)
			if err = awsRestJson1_deserializeJsonRecentIntentSummaryViewList(&v, decoder); err != nil {
				return err
			}
			output.RecentIntentSummaryView = v
		}

		if t == "sessionAttributes" {
			v := make(map[string]string, 0)
			if err = awsRestJson1_deserializeJsonSessionAttribute(&v, decoder); err != nil {
				return err
			}
			output.SessionAttributes = v
		}

		if t == "sessionId" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}

			v, ok := val.(string)
			if !ok {
				return fmt.Errorf("expected SessionId to be of type *String, got %T", err)
			}
			output.SessionId = &v
		}
	}
	// end of the json response body
	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return fmt.Errorf("expected `}` as end token")
	}

	return nil
}

type awsRestJson1_deserializeOpDeleteSession struct{}

// ID is the middleware identifier
func (d awsRestJson1_deserializeOpDeleteSession) ID() string {
	return "awsRestJson1_deserializeOpDeleteSession"
}

func (d awsRestJson1_deserializeOpDeleteSession) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput,
	next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err != nil {
		return out, metadata, err
	}
	response, ok := out.RawResponse.(*smithyHTTP.Response)
	if !ok {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("unknown transport type %T", out.RawResponse),
		}
	}

	output, ok := out.Result.(*DeleteSessionOutput)
	if !ok {
		return out, metadata, &aws.DeserializationError{
			Err: fmt.Errorf("unknown output result type %T", out.Result),
		}
	}

	buff := make([]byte, 1024)
	ringBuffer := sdkio.NewRingBuffer(buff)

	// wrap a TeeReader to read from response body
	body := io.TeeReader(response.Body, ringBuffer)
	defer response.Body.Close()
	decoder := json.NewDecoder(body)
	// UseNumber causes the Decoder to unmarshal a number into an interface{}
	// as a Number instead of as a float64.
	decoder.UseNumber()

	if err := awsRestJson1_deserializeJsonDeleteSessionOutput(output, decoder); err != nil {
		snapshot := make([]byte, 1024)
		ringBuffer.Read(snapshot)
		return out, metadata, &aws.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body with invalid JSON, %w", err),
			Snapshot: snapshot,
		}
	}

	return out, metadata, err
}

func awsRestJson1_deserializeJsonDeleteSessionOutput(output *DeleteSessionOutput, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		// "Empty Response"
		return nil
	}
	if err != nil {
		return err
	}

	if t, ok := startToken.(json.Delim); !ok {
		if t.String() != "{" {
			return fmt.Errorf("expected `{` as start token ")
		}
	}

	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			return err
		}

		if t == "botAlias" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.BotAlias = &v
			} else {
				return fmt.Errorf("expected BotAlias to be of type *String, got %T", val)
			}
		}

		if t == "botName" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.BotName = &v
			} else {
				return fmt.Errorf("expected BotName to be of type *String, got %T", val)
			}
		}

		if t == "sessionId" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.SessionId = &v
			} else {
				return fmt.Errorf("expected SessionId to be of type *String, got %T", val)
			}
		}

		if t == "userId" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.UserId = &v
			} else {
				return fmt.Errorf("expected UserId to be of type *String, got %T", val)
			}
		}
	}

	// end of the json response body
	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return fmt.Errorf("expected `}` as end token")
	}

	return nil
}

func awsRestJson1_deserializeJsonResponseCard(output *types.ResponseCard, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		// "Empty Response"
		return nil
	}
	if err != nil {
		return err
	}

	if t, ok := startToken.(json.Delim); !ok {
		if t.String() != "{" {
			return fmt.Errorf("expected `{` as start token; ")
		}
	}

	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			return err
		}

		// location name : `contentType` key with value as enum
		if t == "contentType" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.ContentType = types.ContentType(v)
			} else {
				return fmt.Errorf("expected ContentType to be of type Enum, got %T", val)
			}
		}

		if t == "genericAttachments" {
			list := make([]types.GenericAttachment, 0)
			if err := awsRestJson1_deserializeJsonGenericAttachmentsList(&list, decoder); err != nil {
				return err
			}
			output.GenericAttachments = list
		}

		// location name : `version` key with value as `*string`
		if t == "version" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}

			if v, ok := val.(string); ok {
				output.Version = &v
			} else {
				return fmt.Errorf("expected Version to be of type *String, got %T", val)
			}
		}
	}

	// end of the json response body
	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return fmt.Errorf("expected `}` as end token")
	}

	return nil
}

func awsRestJson1_deserializeJsonSentimentResponse(output *types.SentimentResponse, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	if t, ok := startToken.(json.Delim); !ok || t.String() != "{" {
		return fmt.Errorf("expected `{` as start token")
	}

	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			return err
		}

		// location name : `sentimentScore` key with value as `*string`
		if t == "sentimentScore" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.SentimentScore = &v
			} else {
				return fmt.Errorf("expected SentimentScore to be of type *String, got %T", val)
			}
		}

		// location name : `sentimentLabel` key with value as `*string`
		if t == "sentimentLabel" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.SentimentLabel = &v
			} else {
				return fmt.Errorf("expected SentimentLabel to be of type *String, got %T", val)
			}
		}
	}

	// end of the json response body
	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return fmt.Errorf("expected `}` as end token")
	}
	return nil
}

func awsRestJson1_deserializeJsonSessionAttribute(output *map[string]string, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	if t, ok := startToken.(json.Delim); !ok || t.String() != "{" {
		return fmt.Errorf("expected `{` as start token")
	}

	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		// based on struct Stage
		if key, ok := token.(string); ok {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				m := *output
				m[key] = v
			} else {
				return fmt.Errorf("expected session attribute value to be of type String, got %T", val)
			}
		} else {
			return fmt.Errorf("expected session attribute key to be of type String, got %T", key)
		}
	}

	// end of the map
	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return fmt.Errorf("expected `}` as end Token ")
	}

	return nil
}

func awsRestJson1_deserializeJsonSlots(output *map[string]string, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	if t, ok := startToken.(json.Delim); !ok || t.String() != "{" {
		return fmt.Errorf("expected `{` as start token")
	}

	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		if key, ok := token.(string); !ok {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				m := *output
				m[key] = v
			} else {
				return fmt.Errorf("expected slots value to be of type String, got %T", val)
			}
		} else {
			return fmt.Errorf("expected slots key to be of type String, got %T", key)
		}
	}
	// end of the map
	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return fmt.Errorf("expected `}` as end token")
	}

	return nil
}

func awsRestJson1_deserializeJsondialogAction(output *types.DialogAction, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	if t, ok := startToken.(json.Delim); !ok || t.String() != "{" {
		return fmt.Errorf("expected `{` as start token")
	}

	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			return err
		}

		if t == "message" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.Message = &v
			} else {
				return fmt.Errorf("expected Message to be of type *String, got %T", val)
			}
		}

		// location name: `intentName` key with value as `string`
		if t == "intentName" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.IntentName = &v
			} else {
				return fmt.Errorf("expected IntentName to be of type *String, got %T", val)
			}
		}

		if t == "messageFormat" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.MessageFormat = types.MessageFormatType(v)
			} else {
				return fmt.Errorf("expected MessageFormat to be of type *String, got %T", val)
			}
		}

		if t == "slotToElicit" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.SlotToElicit = &v
			} else {
				return fmt.Errorf("expected SlotToElicit to be of type *String, got %T", val)
			}
		}

		if t == "slots" {
			v := make(map[string]string, 0)
			if err = awsRestJson1_deserializeJsonSlots(&v, decoder); err != nil {
				return err
			}
			output.Slots = v
		}

		if t == "type" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.Type = types.DialogActionType(v)
			} else {
				return fmt.Errorf("expected Type to be of type Enum, got %T", val)
			}
		}

		if t == "fulfillmentState" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.FulfillmentState = types.FulfillmentState(v)
			} else {
				return fmt.Errorf("expected FulfillmentState to be of type Enum, got %T", val)
			}
		}
	}

	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return fmt.Errorf("expected `}` as end token")
	}
	return nil
}

func awsRestJson1_deserializeJsonRecentIntentSummaryViewList(output *[]types.IntentSummary, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	if t, ok := startToken.(json.Delim); !ok || t.String() != "[" {
		return fmt.Errorf("expected `[` as start token, found %v instead", t)
	}

	for decoder.More() {
		s := types.IntentSummary{}
		if err = awsRestJson1_deserializeJsonIntentSummary(&s, decoder); err != nil {
			return err
		}
		*output = append(*output, s)
	}

	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "]" {
		return fmt.Errorf("expected `]` as end token, found %v instead", t)

	}
	return nil
}

func awsRestJson1_deserializeJsonIntentSummary(output *types.IntentSummary, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	if t, ok := startToken.(json.Delim); !ok || t.String() != "{" {
		return fmt.Errorf("expected `{` as start token")
	}

	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			return err
		}

		if t == "checkpointLabel" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.CheckpointLabel = &v
			} else {
				return fmt.Errorf("expected CheckpointLabel to be of type *String, got %T", val)
			}
		}

		if t == "confirmationStatus" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.ConfirmationStatus = types.ConfirmationStatus(v)
			} else {
				return fmt.Errorf("expected ConfirmationStatus to be of type Enum, got %T", val)
			}
		}

		if t == "dialogActionType" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.DialogActionType = types.DialogActionType(v)
			} else {
				return fmt.Errorf("expected DialogActionType to be of type Enum, got %T", val)
			}
		}

		if t == "fulfillmentState" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.FulfillmentState = types.FulfillmentState(v)
			} else {
				return fmt.Errorf("expected FulfillmentState to be of type Enum, got %T", val)
			}
		}

		if t == "intentName" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.IntentName = &v
			} else {
				return fmt.Errorf("expected IntentName to be of type *String, got %T", val)
			}
		}

		if t == "slotToElicit" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.SlotToElicit = &v
			} else {
				return fmt.Errorf("expected SlotToElicit to be of type *String, got %T", val)
			}
		}

		if t == "slots" {
			v := make(map[string]string, 0)
			if err = awsRestJson1_deserializeJsonSlots(&v, decoder); err != nil {
				return err
			}
			output.Slots = v
		}
	}

	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return fmt.Errorf("expected `}` as end token")
	}
	return nil
}

func awsRestJson1_deserializeJsonGenericAttachmentsList(output *[]types.GenericAttachment, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	if t, ok := startToken.(json.Delim); !ok || t.String() != "[" {
		return fmt.Errorf("expected `[` as start token, found %v instead", t)
	}

	for decoder.More() {
		s := types.GenericAttachment{}
		if err = awsRestJson1_deserializeJsonGenericAttachment(&s, decoder); err != nil {
			return err
		}
		*output = append(*output, s)
	}

	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "]" {
		return fmt.Errorf("expected `]` as end token, found %v instead", t)

	}
	return nil
}

func awsRestJson1_deserializeJsonGenericAttachment(output *types.GenericAttachment, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	if t, ok := startToken.(json.Delim); !ok || t.String() != "{" {
		return fmt.Errorf("expected `{` as start token")
	}

	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			return err
		}

		if t == "attachmentLinkUrl" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.AttachmentLinkUrl = &v
			} else {
				return fmt.Errorf("expected AttachmentLinkUrl to be of type *String, got %T", val)
			}
		}

		if t == "buttons" {
			list := make([]types.Button, 0)
			if err := awsRestJson1_deserializeJsonButtonsList(&list, decoder); err != nil {
				return err
			}
			output.Buttons = list
		}

		if t == "imageUrl" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.ImageUrl = &v
			} else {
				return fmt.Errorf("expected ImageUrl to be of type *String, got %T", val)
			}
		}

		if t == "subTitle" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.SubTitle = &v
			} else {
				return fmt.Errorf("expected SubTitle to be of type *String, got %T", val)
			}
		}

		if t == "title" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.Title = &v
			} else {
				return fmt.Errorf("expected Title to be of type *String, got %T", val)
			}
		}
	}

	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return fmt.Errorf("expected `}` as end token")
	}
	return nil
}

func awsRestJson1_deserializeJsonButtonsList(output *[]types.Button, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	if t, ok := startToken.(json.Delim); !ok || t.String() != "[" {
		return fmt.Errorf("expected `[` as start token, found %v instead", t)
	}

	for decoder.More() {
		s := types.Button{}
		if err = awsRestJson1_deserializeJsonButton(&s, decoder); err != nil {
			return err
		}
		*output = append(*output, s)
	}

	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "]" {
		return fmt.Errorf("expected `]` as end token, found %v instead", t)

	}
	return nil
}

func awsRestJson1_deserializeJsonButton(output *types.Button, decoder *json.Decoder) error {
	if output == nil {
		return fmt.Errorf("unsupported deserialization of nil %T", output)
	}

	startToken, err := decoder.Token()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	if t, ok := startToken.(json.Delim); !ok || t.String() != "{" {
		return fmt.Errorf("expected `{` as start token")
	}

	for decoder.More() {
		t, err := decoder.Token()
		if err != nil {
			return err
		}

		if t == "text" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.Text = &v
			} else {
				return fmt.Errorf("expected Text to be of type *String, got %T", val)
			}
		}

		if t == "value" {
			val, err := decoder.Token()
			if err != nil {
				return err
			}
			if v, ok := val.(string); ok {
				output.Value = &v
			} else {
				return fmt.Errorf("expected Value to be of type *String, got %T", val)
			}
		}
	}

	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return fmt.Errorf("expected `}` as end token")
	}
	return nil
}
