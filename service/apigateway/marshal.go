package apigateway

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/json"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/rest"
)

// protoCreateAPIKeyMarshaler defines a marshaler for CreateApiKey operation
type protoCreateAPIKeyMarshaler struct {
	input *CreateApiKeyInput
}

// marshalOperation is the top level method used within a handler stack to marshal an operation
// This method calls appropriate marshal shape functions as per the input shape and protocol used by the service.
func (m protoCreateAPIKeyMarshaler) marshalOperation(r *aws.Request) {
	var err error
	encoder := rest.NewEncoder(r.HTTPRequest)

	// We add Content-Type Header if input shape's is not of type streaming payload
	// The value of Content-Type is decided by following:
	// a. if shape's metadata has JSONVersion and protocol is json.
	//    - application/x-amz-json-%s where %s is ths JSONVersion.
	// b. else if protocol is either json or rest-json
	//     - application/json
	// Here protocol is rest-json, and shape is not of type streaming payload,
	// thus content-type header with value application/json is added.
	encoder.AddHeader("Content-Type").String("application/json")
	err = encoder.Encode()
	if err != nil {
		r.Error = err
		return
	}

	jsonEncoder := json.NewEncoder()
	err = marshalCreateAPIKeyInputShapeAWSJSON(m.input, r, jsonEncoder)
	if err != nil {
		r.Error = err
	}
}

// marshalCreateAPIKeyInputShapeAWSJSON is a stand alone function used to marshal the json body
func marshalCreateAPIKeyInputShapeAWSJSON(v *CreateApiKeyInput, r *aws.Request, enc *json.Encoder) error {
	if r.Metadata.TargetPrefix != "" {
		target := r.Metadata.TargetPrefix + "." + r.Operation.Name
		r.HTTPRequest.Header.Add("X-Amz-Target", target)
	}
	if r.Metadata.JSONVersion != "" {
		jsonVersion := r.Metadata.JSONVersion
		r.HTTPRequest.Header.Add("Content-Type", "application/x-amz-json-"+jsonVersion)
	}

	if v != nil {
		obj := enc.Object()
		if v.CustomerId != nil {
			obj.Key("customerId").String(*v.CustomerId)
		}
		if v.Description != nil {
			obj.Key("description").String(*v.Description)
		}
		if v.Enabled != nil {
			obj.Key("enabled").Boolean(*v.Enabled)
		}
		if v.GenerateDistinctId != nil {
			obj.Key("generateDistinctId").Boolean(*v.GenerateDistinctId)
		}
		if v.Name != nil {
			obj.Key("name").String(*v.Name)
		}

		if v.StageKeys != nil {
			stageKeyArray := obj.Key("stageKeys").Array()
			for _, stageKey := range v.StageKeys {
				stageKeyArrayObj := stageKeyArray.Value().Object()
				if stageKey.RestApiId != nil {
					stageKeyArrayObj.Key("restApiId").String(*stageKey.RestApiId)
				}
				if stageKey.StageName != nil {
					stageKeyArrayObj.Key("stageName").String(*stageKey.StageName)
				}
				stageKeyArrayObj.Close()
			}
			stageKeyArray.Close()
		}

		if v.Tags != nil {
			tagObj := obj.Key("tags").Object()
			for k, val := range v.Tags {
				tagObj.Key(k).String(val)
			}
			tagObj.Close()
		}

		if v.Value != nil {
			obj.Key("value").String(*v.Value)
		}
		obj.Close()
	}
	r.SetBufferBody([]byte(enc.String()))
	return nil
}

// namedHandler returns a Named Build Handler for an operation marshal function
func (m protoCreateAPIKeyMarshaler) namedHandler() aws.NamedHandler {
	const buildHandler = "ProtoGetApiKey.BuildHandler"
	return aws.NamedHandler{
		Name: buildHandler,
		Fn:   m.marshalOperation,
	}
}
