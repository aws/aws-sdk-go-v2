package apigateway

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/internal/sdkio"
	restlegacy "github.com/aws/aws-sdk-go-v2/private/protocol/rest"
)

const errorCodeHeader = "X-Amzn-Errortype"

// protoCreateAPIKeyUnmarshaler defines unmarshaler forProtoCreateAPIKey Operation
type protoCreateAPIKeyUnmarshaler struct {
	output *CreateApiKeyOutput
}

// unmarshalOperation is the top level method used with a handler stack to unmarshal an operation response
// This method calls appropriate unmarshal shape functions as per the output shape and protocol used by the service.
func (u protoCreateAPIKeyUnmarshaler) unmarshalOperation(r *aws.Request) {
	// initializes a ring buffer
	buff := make([]byte, 1024)
	ringBuffer := sdkio.NewRingBuffer(buff)
	if isRequestError(r) {
		jsonErr := jsonErrorResponse{}
		// wrap a TeeReader to read from response body & write on a ring buffer
		body := io.TeeReader(r.HTTPResponse.Body, ringBuffer)
		defer r.HTTPResponse.Body.Close()
		// build a json decoder
		decoder := json.NewDecoder(body)
		// call json error unmarshaler
		unmarshalErrorShapeAWSJSON(decoder, &jsonErr, ringBuffer)
		// Get error code from Header error type
		code := r.HTTPResponse.Header.Get(errorCodeHeader)
		if code == "" {
			code = jsonErr.Code
		}

		code = strings.SplitN(code, ":", 2)[0]
		r.Error = awserr.NewRequestFailure(
			awserr.New(code, jsonErr.Message, nil),
			r.HTTPResponse.StatusCode,
			r.RequestID,
		)
		return
	}

	// If not error, then call rest unmarshaler
	restlegacy.UnmarshalMeta(r)

	// unmarshal output shape
	r.Error = unmarshalProtoCreateAPIKeyOutputShapeAWSJSON(u.output, r.HTTPResponse.Body, ringBuffer)
}

func unmarshalProtoCreateAPIKeyOutputShapeAWSJSON(output *CreateApiKeyOutput, responseBody io.ReadCloser, ringBuffer *sdkio.RingBuffer) error {
	// wrap a TeeReader to read from response body & write on snapshot
	body := io.TeeReader(responseBody, ringBuffer)
	defer responseBody.Close()
	// build a json decoder
	decoder := json.NewDecoder(body)

	// start token
	startToken, err := decoder.Token()
	if err == io.EOF {
		// "Empty Response"
		return nil
	}
	if err != nil {
		return awserr.New(aws.ErrCodeSerialization,
			fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
				getSnapshot(ringBuffer)), err)
	}

	if t, ok := startToken.(json.Delim); !ok {
		if t.String() != "{" {
			return awserr.New(aws.ErrCodeSerialization,
				fmt.Sprintf("failed to decode response body with invalid JSON, expected `{` as start Token. "+
					"Here's a snapshot: %s",
					getSnapshot(ringBuffer)), err)
		}
	}

	for decoder.More() {
		// fetch token for key
		t, err := decoder.Token()
		if err != nil {
			return awserr.New(aws.ErrCodeSerialization,
				fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
					getSnapshot(ringBuffer)), err)
		}

		// location name : `value` key with value as `*string`
		if t == "value" {
			val, err := decoder.Token()
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			if v, ok := val.(string); ok {
				output.Value = &v
			} else {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("Expected Value to be of type string. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
		}

		// location name : `name` key with value as `*string`
		if t == "name" {
			val, err := decoder.Token()
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			if v, ok := val.(string); ok {
				output.Name = &v
			} else {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("Expected Name to be of type string. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
		}

		// location name : `description` key with value as `*string`
		if t == "description" {
			val, err := decoder.Token()
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			if v, ok := val.(string); ok {
				output.Description = &v
			} else {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("Expected Description to be of type string. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
		}

		// location name : `customerId` key with value as `*string`
		if t == "customerId" {
			val, err := decoder.Token()
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			if v, ok := val.(string); ok {
				output.CustomerId = &v
			} else {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("Expected CustomerId to be of type string. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
		}

		// location name : `id` key with value as `*string`
		if t == "id" {
			val, err := decoder.Token()
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			if v, ok := val.(string); ok {
				output.Id = &v
			} else {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("Expected Id to be of type string. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
		}

		// location name : `enabled` key with value as `*boolean`
		if t == "enabled" {
			val, err := decoder.Token()
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			if v, ok := val.(bool); ok {
				output.Enabled = &v
			} else {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("Expected Enabled to be of type boolean. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
		}

		// location name : `stageKeys` key with value as `[]string`
		if t == "stageKeys" {
			// create []string as modeled for the member shape
			list := make([]string, 0)
			err = unmarshalStageKeysListShapeAWSJSON(decoder, &list, ringBuffer)
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("Error serializing StageKeys . Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			// Attach the de-serialized list of string to output shape
			output.StageKeys = list
		}

		// location name : `tags` key with value as `map`
		if t == "tags" {
			// create []string as modeled for the member shape
			m := make(map[string]string, 0)
			// start of the list
			err = unmarshalTagsMapShapeAWSJSON(decoder, &m, ringBuffer)
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("Error serializing Tags. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			// Attach de-serialized Tag to output
			output.Tags = m
		}

		// location name : `createdDate` with value as `timestamp`
		if t == "createdDate" {
			val, err := decoder.Token()
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			if v, ok := val.(float64); ok {
				time := time.Unix(int64(v), 0).UTC()
				output.CreatedDate = &time
			} else {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("Expected CreatedDate to be of type timestamp. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
		}

		// location name : `lastUpdatedDate` with value as `timestamp`
		if t == "lastUpdatedDate" {
			val, err := decoder.Token()
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			if v, ok := val.(float64); ok {
				time := time.Unix(int64(v), 0).UTC()
				output.LastUpdatedDate = &time
			} else {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("Expected LastUpdatedDate to be of type timestamp. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
		}

	}
	// end of the json response body
	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return awserr.New(aws.ErrCodeSerialization,
			fmt.Sprintf("failed to decode response body with invalid JSON, expected `}` as end Token. "+
				"Here's a snapshot: %s",
				getSnapshot(ringBuffer)), err)
	}

	return nil
}

func unmarshalStageKeysListShapeAWSJSON(dec *json.Decoder, list *[]string, ringBuffer *sdkio.RingBuffer) error {
	// start of the list
	startToken, err := dec.Token()

	if err == io.EOF {
		// "Empty List"
		return nil
	}
	if err != nil {
		return awserr.New(aws.ErrCodeSerialization,
			fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
				getSnapshot(ringBuffer)), err)
	}
	if t, ok := startToken.(json.Delim); !ok || t.String() != "[" {
		return awserr.New(aws.ErrCodeSerialization,
			fmt.Sprintf("Expected json delimiter at the start of list, found %v instead. "+
				"Here's a snapshot: %s", t, getSnapshot(ringBuffer)), err)
	}

	for dec.More() {
		token, err := dec.Token()
		if err != nil {
			return awserr.New(aws.ErrCodeSerialization,
				fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
					getSnapshot(ringBuffer)), err)
		}
		// based on struct Stage
		if v, ok := token.(string); ok {
			val := v
			*list = append(*list, val)
		} else {
			return awserr.New(aws.ErrCodeSerialization,
				fmt.Sprintf("Expected stageKey to be of type string. Here's a snapshot: %s",
					getSnapshot(ringBuffer)), err)
		}
	}

	// end of the list
	endToken, err := dec.Token()
	if err != nil {
		return awserr.New(aws.ErrCodeSerialization,
			fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
				getSnapshot(ringBuffer)), err)
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "]" {
		return awserr.New(aws.ErrCodeSerialization,
			fmt.Sprintf("failed to decode response body with invalid JSON, expected `]` as end Token. "+
				"Here's a snapshot: %s",
				getSnapshot(ringBuffer)), err)
	}
	return nil
}

func unmarshalTagsMapShapeAWSJSON(dec *json.Decoder, output *map[string]string, ringBuffer *sdkio.RingBuffer) error {
	startToken, err := dec.Token()
	if err == io.EOF {
		// "Empty List"
		return nil
	}
	if err != nil {
		return awserr.New(aws.ErrCodeSerialization,
			fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
				getSnapshot(ringBuffer)), err)
	}
	if t, ok := startToken.(json.Delim); !ok || t.String() != "{" {
		return awserr.New(aws.ErrCodeSerialization,
			fmt.Sprintf("Expected json delimiter at the start of list, found %v instead. "+
				"Here's a snapshot: %s", t,
				getSnapshot(ringBuffer)), err)
	}

	// decoder
	for dec.More() {
		token, err := dec.Token()
		if err != nil {
			return awserr.New(aws.ErrCodeSerialization,
				fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
					getSnapshot(ringBuffer)), err)
		}

		// based on struct Stage
		if key, ok := token.(string); ok {
			val, err := dec.Token()
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			if v, ok := val.(string); !ok {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("invalid json, expected string. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			} else {
				m := *output
				m[key] = v
			}
		} else {
			return awserr.New(aws.ErrCodeSerialization,
				fmt.Sprintf("Expected %v, to be of type string, got %T. "+
					"Here's a snapshot: %s", key, key,
					getSnapshot(ringBuffer)), err)
		}
	}

	// end of the map
	endToken, err := dec.Token()
	if err != nil {
		return awserr.New(aws.ErrCodeSerialization,
			fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
				getSnapshot(ringBuffer)), err)
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return awserr.New(aws.ErrCodeSerialization,
			fmt.Sprintf("failed to decode response body with invalid JSON, expected `}` as end Token. "+
				"Here's a snapshot: %s",
				getSnapshot(ringBuffer)), err)
	}
	return nil
}

// isRequestError would check if a request response was an error
func isRequestError(r *aws.Request) bool {
	if r.HTTPResponse.StatusCode == 0 || r.HTTPResponse.StatusCode >= 300 {
		return true
	}
	return false
}

// namedHandler returns a named handler for an operation unmarshal function
func (u protoCreateAPIKeyUnmarshaler) namedHandler() aws.NamedHandler {
	return aws.NamedHandler{
		Name: "ProtoCreateAPIKey.UnmarshalHandler",
		Fn:   u.unmarshalOperation,
	}
}

// unmarshalErrorShapeAWSJSON unmarshal's a json error response body for the REST JSON protocol..
// some service may have custom error handling
// here we do not handle modelled exceptions.
func unmarshalErrorShapeAWSJSON(decoder *json.Decoder, errorResponse *jsonErrorResponse, ringBuffer *sdkio.RingBuffer) error {
	startToken, err := decoder.Token()
	if err == io.EOF {
		// "Empty error response body"
		return nil
	}
	if err != nil {
		return awserr.New(aws.ErrCodeSerialization,
			fmt.Sprintf("failed to decode error response body with invalid JSON. Here's a snapshot: %s",
				getSnapshot(ringBuffer)), err)
	}

	if t, ok := startToken.(json.Delim); !ok {
		if t.String() != "{" {
			return awserr.New(aws.ErrCodeSerialization,
				fmt.Sprintf("failed to decode response body with invalid JSON, expected `{` as start Token. "+
					"Here's a snapshot: %s",
					getSnapshot(ringBuffer)), err)
		}
	}

	for decoder.More() {
		// fetch token for key
		t, err := decoder.Token()
		if err != nil {
			return awserr.New(aws.ErrCodeSerialization,
				fmt.Sprintf("failed to decode response body with invalid JSON. Here's a snapshot: %s",
					getSnapshot(ringBuffer)), err)
		}

		if t == "code" {
			val, err := decoder.Token()
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("failed to decode response error with invalid JSON. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			if v, ok := val.(string); ok {
				errorResponse.Code = v
			} else {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("Expected error code to be of type string. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}

		}

		if t == "message" {
			val, err := decoder.Token()
			if err != nil {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("failed to decode response error with invalid JSON. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
			if v, ok := val.(string); ok {
				errorResponse.Message = v
			} else {
				return awserr.New(aws.ErrCodeSerialization,
					fmt.Sprintf("Expected error message to be of type string. Here's a snapshot: %s",
						getSnapshot(ringBuffer)), err)
			}
		}
	}

	// end of the json error response body
	endToken, err := decoder.Token()
	if err != nil {
		return err
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return awserr.New(aws.ErrCodeSerialization,
			fmt.Sprintf("failed to decode response body with invalid JSON, expected `}` as end Token. "+
				"Here's a snapshot: %s",
				getSnapshot(ringBuffer)), err)
	}

	return nil
}

type jsonErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// getSnapshot helper utility reads 1024 bytes from ring buffer into the byte slice and returns it
// getSnapshot takes in a ringBuffer to read
func getSnapshot(rb *sdkio.RingBuffer) []byte {
	snapshot := make([]byte, 1024)
	rb.Read(snapshot)
	return snapshot
}
