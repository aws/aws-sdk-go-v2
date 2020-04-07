package aws_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdkio"
	"github.com/aws/aws-sdk-go-v2/private/protocol/json/jsonutil"
)

type mockOP struct {
	Message *string `locationName:"message" type:"string"`
}

type mockUnmarshaler struct {
	output *mockOP
}

const correctResponse = `{"message": "Hello"}`
const responseWithMissingDelimiter = `"message": "Hello"}`
const invalidResponse = `{"message": true}`

func TestHTTPDeserializationError(t *testing.T) {
	cases := map[string]struct {
		responseBody   string
		responseStatus int
		requestID      string
		errorIsSet     bool
	}{
		"No error": {
			responseBody:   correctResponse,
			requestID:      "mockReqID",
			responseStatus: 200,
		},
		"Missing delimiter": {
			responseBody:   responseWithMissingDelimiter,
			errorIsSet:     true,
			requestID:      "mockReqID",
			responseStatus: 200,
		},
		"Invalid response": {
			responseBody:   invalidResponse,
			errorIsSet:     true,
			requestID:      "mockReqID",
			responseStatus: 200,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			r := &aws.Request{
				HTTPResponse: &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(c.responseBody))),
				},
				RequestID: "mockReqID",
			}

			u := mockUnmarshaler{
				output: &mockOP{},
			}
			// unmarshal request response
			u.unmarshalOperation(r)

			if c.errorIsSet {
				if r.Error == nil {
					t.Fatal("Expected error, got none")
				}
				if r.Error == nil {
					t.Fatal("Expected error, got none")
				}
				var e *aws.DeserializationError
				if errors.As(r.Error, &e) {
					if e, a := c.responseBody, e.Reason; !strings.Contains(a, e) {
						t.Fatalf("Expected response body to contain %v, got %v", e, a)
					}
				} else {
					t.Fatalf("Expected error to be of type %T, got %T", e, r.Error)
				}
			} else {
				if r.Error != nil {
					t.Fatalf("Expected no error, got %v", r.Error)
				}
			}
		})
	}
}

// unmarshal operation unmarshal's request response
func (u *mockUnmarshaler) unmarshalOperation(r *aws.Request) {
	b := make([]byte, 1024)
	ringbuffer := sdkio.NewRingBuffer(b)
	// wraps ring buffer around the response body
	body := io.TeeReader(r.HTTPResponse.Body, ringbuffer)

	// If unmarshaling function returns an error, it is a deserialization error
	if err := jsonutil.UnmarshalJSON(u.output, body); err != nil {
		snapshot := make([]byte, 1024)
		ringbuffer.Read(snapshot)
		r.Error = &aws.DeserializationError{
			Reason: fmt.Sprintf("Here's a snapshot of response being deserialized: %s", snapshot), // Additional context
			Err:    err,
		}
	}
}

// // Call unmarshal output shape of JSON doc
// func unmarshalJSON(dec *json.Decoder, output *mockOP) error {
// 	// start token
// 	startToken, err := dec.Token()
// 	if err == io.EOF {
// 		// Empty Response
// 		return nil
// 	}
// 	if err != nil {
// 		return fmt.Errorf("Failed  to decode JSON: %v", err)
// 	}
// 	if t, ok := startToken.(json.Delim); !ok {
// 		if t.String() != "{" {
// 			return fmt.Errorf("invalid JSON , expected { at start of json")
// 		}
// 	}
//
// 	for dec.More() {
// 		// fetch token for key
// 		t, err := dec.Token()
// 		if err != nil {
// 			return fmt.Errorf("Failed  to decode JSON: %v", err)
// 		}
//
// 		// key with value as `*string`
// 		if t == "message" {
// 			val, err := dec.Token()
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			if v, ok := val.(string); ok {
// 				output.Message = &v
// 			} else {
// 				return fmt.Errorf("expected Message to be of type string, got %T", val)
// 			}
// 		}
// 	}
//
// 	// end token
// 	endToken, err := dec.Token()
// 	if err != nil {
// 		return fmt.Errorf("Failed  to decode JSON: %v ", err)
// 	}
// 	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
// 		return fmt.Errorf("invalid JSON , expected } at start of json")
// 	}
// 	return nil
// }
