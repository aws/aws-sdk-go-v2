package lexruntimeservice

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/json"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/rest"
	"github.com/aws/aws-sdk-go-v2/service/smithyprototype/lexruntimeservice/types"
	"github.com/awslabs/smithy-go/middleware"
	smithyHTTP "github.com/awslabs/smithy-go/transport/http"
)

type awsRestJson1_serializeOpDeleteSession struct{}

// ID is the identifier for the middleware
func (g *awsRestJson1_serializeOpDeleteSession) ID() string {
	return "awsRestJson1_serializeOpDeleteSession"
}

// HandleSerialize will serialize the middleware input parameters to the provided input http request
func (g *awsRestJson1_serializeOpDeleteSession) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	request, ok := in.Request.(*smithyHTTP.Request)
	if !ok {
		return out, metadata, &aws.SerializationError{Err: fmt.Errorf("unknown transport type %T", in.Request)}
	}

	input, ok := in.Parameters.(*DeleteSessionInput)
	if !ok {
		return out, metadata, &aws.SerializationError{Err: fmt.Errorf("unknown input parameters type %T", in.Parameters)}
	}

	restEncoder := rest.NewEncoder(request.Request)

	if err := awsRestJson1_serializeRestDeleteSessionInput(input, restEncoder); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}

	if err := restEncoder.Encode(); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}

	return next.HandleSerialize(ctx, in)
}

type awsRestJson1_serializeOpGetSession struct{}

// ID is the identifier for the middleware
func (g *awsRestJson1_serializeOpGetSession) ID() string {
	return "awsRestJson1_serializeOpGetSession"
}

// HandleSerialize will serialize the middleware input parameters to the provided input http request
func (g *awsRestJson1_serializeOpGetSession) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	request, ok := in.Request.(*smithyHTTP.Request)
	if !ok {
		return out, metadata, &aws.SerializationError{Err: fmt.Errorf("unknown transport type %T", in.Request)}
	}

	input, ok := in.Parameters.(*GetSessionInput)
	if !ok {
		return out, metadata, &aws.SerializationError{Err: fmt.Errorf("unknown input parameters type %T", in.Parameters)}
	}

	restEncoder := rest.NewEncoder(request.Request)

	if err := awsRestJson1_serializeRestGetSessionInput(input, restEncoder); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}

	if err := restEncoder.Encode(); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}

	return next.HandleSerialize(ctx, in)
}

type awsRestJson1_serializeOpPostContent struct{}

func (p awsRestJson1_serializeOpPostContent) ID() string {
	return "awsRestJson1_serializeOpPostContent"
}

func (p awsRestJson1_serializeOpPostContent) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	request, ok := in.Request.(*smithyHTTP.Request)
	if !ok {
		return out, metadata, &aws.SerializationError{Err: fmt.Errorf("unknown transport type %T", in.Request)}
	}

	params, ok := in.Parameters.(*PostContentInput)
	if !ok {
		return out, metadata, &aws.SerializationError{Err: fmt.Errorf("unknown params parameters type %T", in.Parameters)}
	}

	restEncoder := rest.NewEncoder(request.Request)

	if err := awsRestJson1_serializeRestPostContentInput(params, restEncoder); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}

	if err := restEncoder.Encode(); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}

	// TODO: If MediaType trait was applied in the Smithy model the content-type should get reflected as such
	if request, err = request.SetStream(params.InputStream); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}
	in.Request = request

	return next.HandleSerialize(ctx, in)
}

type awsRestJson1_serializeOpPostText struct{}

func (p awsRestJson1_serializeOpPostText) ID() string {
	return "awsRestJson1_serializeOpPostText"
}

func (p awsRestJson1_serializeOpPostText) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	request, ok := in.Request.(*smithyHTTP.Request)
	if !ok {
		return out, metadata, &aws.SerializationError{Err: fmt.Errorf("unknown transport type %T", in.Request)}
	}

	input, ok := in.Parameters.(*PostTextInput)
	if !ok {
		return out, metadata, &aws.SerializationError{Err: fmt.Errorf("unknown input parameters type %T", in.Parameters)}
	}

	restEncoder := rest.NewEncoder(request.Request)

	restEncoder.AddHeader("Content-Type").String("application/json")

	if err := awsRestJson1_serializeRestPostTextInput(input, restEncoder); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}

	if err := restEncoder.Encode(); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}

	jsonEncoder := json.NewEncoder()
	if err := awsRestJson1_serializeJsonPostTextInput(input, jsonEncoder.Value); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}

	if request, err = request.SetStream(bytes.NewReader(jsonEncoder.Bytes())); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}
	in.Request = request

	return next.HandleSerialize(ctx, in)
}

type awsRestJson1_serializeOpPutSession struct{}

// ID is the middleware identifier
func (p *awsRestJson1_serializeOpPutSession) ID() string {
	return "awsRestJson1_serializeOpPutSession"
}

func (p *awsRestJson1_serializeOpPutSession) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	request, ok := in.Request.(*smithyHTTP.Request)
	if !ok {
		return out, metadata, &aws.SerializationError{Err: fmt.Errorf("unknown transport type %T", in.Request)}
	}

	input, ok := in.Parameters.(*PutSessionInput)
	if !ok {
		return out, metadata, &aws.SerializationError{Err: fmt.Errorf("unknown input parameters type %T", in.Parameters)}
	}

	restEncoder := rest.NewEncoder(request.Request)

	restEncoder.AddHeader("Content-Type").String("application/json")

	if err := awsRestJson1_serializeRestPutSessionInput(input, restEncoder); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}

	jsonEncoder := json.NewEncoder()
	if err := awsRestJson1_serializeJsonPutSessionInput(input, jsonEncoder.Value); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}

	if request, err = request.SetStream(bytes.NewReader(jsonEncoder.Bytes())); err != nil {
		return out, metadata, &aws.SerializationError{Err: err}
	}
	in.Request = request

	return next.HandleSerialize(ctx, in)
}

func awsRestJson1_serializeJsonPutSessionInput(v *PutSessionInput, encoder json.Value) error {
	if v == nil {
		return fmt.Errorf("unsupported serialization of nil %T", v)
	}

	object := encoder.Object()
	defer object.Close()

	if v.DialogAction != nil {
		value := object.Key("dialogAction")
		if err := awsRestJson1_serializeJsonDialogAction(v.DialogAction, value); err != nil {
			return err
		}
	}

	if v.RecentIntentSummaryView != nil {
		value := object.Key("recentIntentSummaryView")
		if err := awsRestJson1_serializeJsonIntentSummaryList(v.RecentIntentSummaryView, value); err != nil {
			return err
		}
	}

	if v.SessionAttributes != nil {
		value := object.Key("sessionAttributes")
		if err := awsRestJson1_serializeJsonStringMap(v.SessionAttributes, value); err != nil {
			return err
		}
	}

	return nil
}

func awsRestJson1_serializeRestPutSessionInput(v *PutSessionInput, encoder *rest.Encoder) error {
	if v == nil {
		return fmt.Errorf("unsupported serialization of nil %T", v)
	}

	if v.Accept != nil {
		encoder.SetHeader("Accept").String(*v.Accept)
	}

	if v.BotAlias != nil {
		if err := encoder.SetURI("botAlias").String(*v.BotAlias); err != nil {
			return err
		}
	}

	if v.BotName != nil {
		if err := encoder.SetURI("botName").String(*v.BotName); err != nil {
			return err
		}
	}

	if v.UserId != nil {
		if err := encoder.SetURI("userId").String(*v.UserId); err != nil {
			return err
		}
	}

	return nil
}

func awsRestJson1_serializeRestPostContentInput(v *PostContentInput, encoder *rest.Encoder) error {
	if v == nil {
		return fmt.Errorf("unsupported serialization of nil %T", v)
	}

	if v.Accept != nil {
		encoder.SetHeader("Accept").String(*v.Accept)
	}

	if v.BotAlias != nil {
		if err := encoder.SetURI("botAlias").String(*v.BotAlias); err != nil {
			return err
		}
	}

	if v.BotName != nil {
		if err := encoder.SetURI("botName").String(*v.BotName); err != nil {
			return err
		}
	}

	if v.ContentType != nil {
		encoder.SetHeader("Content-Type").String(*v.ContentType)
	}

	if v.UserId != nil {
		if err := encoder.SetURI("userId").String(*v.UserId); err != nil {
			return err
		}
	}

	if v.RequestAttributes != nil {
		if err := encoder.SetHeader("x-amz-lex-session-attributes").JSONValue(v.RequestAttributes); err != nil {
			return err
		}
	}

	if v.SessionAttributes != nil {
		if err := encoder.SetHeader("x-amz-lex-session-attributes").JSONValue(v.SessionAttributes); err != nil {
			return err
		}
	}

	return nil
}

// awsRestJson1_serializeRestGetSessionInput marshals the top level members of GetSessionInput that have HTTP bindings
func awsRestJson1_serializeRestGetSessionInput(v *GetSessionInput, encoder *rest.Encoder) error {
	if v == nil {
		return fmt.Errorf("unsupported serialization of nil %T", v)
	}

	if v.BotAlias != nil {
		if err := encoder.SetURI("botAlias").String(*v.BotAlias); err != nil {
			return err
		}
	}

	if v.BotName != nil {
		if err := encoder.SetURI("botName").String(*v.BotName); err != nil {
			return err
		}
	}

	if v.CheckpointLabelFilter != nil {
		encoder.SetQuery("checkpointLabelFilter").String(*v.CheckpointLabelFilter)
	}

	if v.UserId != nil {
		if err := encoder.SetURI("userId").String(*v.UserId); err != nil {
			return err
		}
	}

	return nil
}

// awsRestJson1_serializeRestDeleteSessionInput marshals the top level members of GetSessionInput that have HTTP bindings
func awsRestJson1_serializeRestDeleteSessionInput(v *DeleteSessionInput, encoder *rest.Encoder) error {
	if v == nil {
		return fmt.Errorf("unsupported serialization of nil %T", v)
	}

	if v.BotAlias != nil {
		if err := encoder.SetURI("botAlias").String(*v.BotAlias); err != nil {
			return err
		}
	}

	if v.BotName != nil {
		if err := encoder.SetURI("botName").String(*v.BotName); err != nil {
			return err
		}
	}

	if v.UserId != nil {
		if err := encoder.SetURI("userId").String(*v.UserId); err != nil {
			return err
		}
	}

	return nil
}

func awsRestJson1_serializeJsonPostTextInput(v *PostTextInput, value json.Value) error {
	if v == nil {
		return fmt.Errorf("unsupported serialization of nil %T", v)
	}

	object := value.Object()
	defer object.Close()

	if v.InputText != nil {
		object.Key("inputText").String(*v.InputText)
	}

	if v.RequestAttributes != nil {
		value := object.Key("requestAttributes")
		if err := awsRestJson1_serializeJsonStringMap(v.RequestAttributes, value); err != nil {
			return err
		}
	}

	if v.SessionAttributes != nil {
		value := object.Key("sessionAttributes")
		if err := awsRestJson1_serializeJsonStringMap(v.SessionAttributes, value); err != nil {
			return err
		}
	}

	return nil
}

func awsRestJson1_serializeRestPostTextInput(v *PostTextInput, encoder *rest.Encoder) error {
	if v == nil {
		return fmt.Errorf("unsupported serialization of nil %T", v)
	}

	if v.BotAlias != nil {
		if err := encoder.SetURI("botAlias").String(*v.BotAlias); err != nil {
			return err
		}
	}

	if v.BotName != nil {
		if err := encoder.SetURI("botName").String(*v.BotName); err != nil {
			return err
		}
	}

	if v.UserId != nil {
		if err := encoder.SetURI("userId").String(*v.UserId); err != nil {
			return err
		}
	}

	return nil
}

func awsRestJson1_serializeJsonIntentSummaryList(v []types.IntentSummary, value json.Value) error {
	if v == nil {
		return fmt.Errorf("unsupported serialization of nil %T", v)
	}

	array := value.Array()

	for _, arrayValue := range v {
		av := array.Value()
		if err := awsRestJson1_serializeIntentSummary(&arrayValue, av); err != nil {
			return err
		}
	}

	return nil
}

func awsRestJson1_serializeIntentSummary(v *types.IntentSummary, value json.Value) error {
	if v == nil {
		return fmt.Errorf("unsupported serialization of nil %T", v)
	}

	object := value.Object()
	defer object.Close()

	if v.CheckpointLabel != nil {
		object.Key("checkpointLabel").String(*v.CheckpointLabel)
	}

	if len(v.ConfirmationStatus) > 0 {
		object.Key("confirmationStatus").String(string(v.ConfirmationStatus))
	}

	if len(v.DialogActionType) > 0 {
		object.Key("dialogActionType").String(string(v.DialogActionType))
	}

	if len(v.FulfillmentState) > 0 {
		object.Key("fulfillmentState").String(string(v.FulfillmentState))
	}

	if v.IntentName != nil {
		object.Key("intentName").String(*v.IntentName)
	}

	if v.SlotToElicit != nil {
		object.Key("slotToElicit").String(*v.SlotToElicit)
	}

	if v.Slots != nil {
		value := object.Key("slots")
		if err := awsRestJson1_serializeJsonStringMap(v.Slots, value); err != nil {
			return err
		}
	}

	return nil
}

func awsRestJson1_serializeJsonDialogAction(v *types.DialogAction, value json.Value) error {
	if v == nil {
		return fmt.Errorf("unsupported serialization of nil %T", v)
	}

	object := value.Object()
	defer object.Close()

	if len(v.FulfillmentState) > 0 {
		object.Key("fulfillmentState").String(string(v.FulfillmentState))
	}

	if v.IntentName != nil {
		object.Key("intentName").String(*v.IntentName)
	}

	if v.Message != nil {
		object.Key("message").String(*v.Message)
	}

	if len(v.MessageFormat) > 0 {
		object.Key("messageFormat").String(string(v.MessageFormat))
	}

	if v.SlotToElicit != nil {
		object.Key("slotToElicit").String(*v.SlotToElicit)
	}

	if v.Slots != nil {
		value := object.Key("slots")
		if err := awsRestJson1_serializeJsonStringMap(v.Slots, value); err != nil {
			return err
		}
	}

	if len(v.Type) > 0 {
		object.Key("type").String(string(v.Type))
	}

	return nil
}

func awsRestJson1_serializeJsonStringMap(v map[string]string, value json.Value) error {
	if v == nil {
		return fmt.Errorf("unsupported serialization of nil %T", v)
	}

	object := value.Object()
	defer object.Close()

	for k, kv := range v {
		object.Key(k).String(kv)
	}

	return nil
}
