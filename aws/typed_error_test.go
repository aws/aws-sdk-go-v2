package aws_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdkio"
)

type mockOP struct {
	Message *string
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

			op := mockOP{}
			u := mockUnmarshaler{
				output: &op,
			}
			// unmarshal request response
			u.unmarshalOperation(r)

			if r.Error != nil && !c.errorIsSet {
				t.Fatalf("Expected no error, got %v", r.Error)
			}

			if c.errorIsSet {
				var e *aws.HTTPDeserializationError
				if errors.As(r.Error, &e) {
					if e, a := c.requestID, e.RequestID; e != a {
						t.Fatalf("Expected request ID to be %v, got %v", e, a)
					}
					if e, a := c.responseStatus, e.ErrorStatusCode(); e != a {
						t.Fatalf("Expected request ID to be %v, got %v", e, a)
					}
					if e, a := c.responseBody, e.Reason; !strings.Contains(a, e) {
						t.Fatalf("Expected request ID to be %v, got %v", e, a)
					}
				} else {
					t.Fatalf("Expected error to be of type %T", e)
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
	decoder := json.NewDecoder(body)

	// If unmarshaling function returns an error, it is a deserialization error
	if err := unmarshalJSON(decoder, u.output); err != nil {
		snapshot := make([]byte, 1024)
		ringbuffer.Read(snapshot)
		r.Error = &aws.HTTPDeserializationError{
			Response:  r.HTTPResponse,
			RequestID: r.RequestID,
			Reason:    fmt.Sprintf("Here's a snapshot of response being deserialized: %s", snapshot), // Additional context
			Err:       err,
		}
	}
}

// Call unmarshal output shape of JSON doc
func unmarshalJSON(dec *json.Decoder, output *mockOP) error {
	// start token
	startToken, err := dec.Token()
	if err == io.EOF {
		// Empty Response
		return nil
	}
	if err != nil {
		return fmt.Errorf("Failed  to decode JSON: %v", err)
	}
	if t, ok := startToken.(json.Delim); !ok {
		if t.String() != "{" {
			return fmt.Errorf("invalid JSON , expected { at start of json")
		}
	}

	for dec.More() {
		// fetch token for key
		t, err := dec.Token()
		if err != nil {
			return fmt.Errorf("Failed  to decode JSON: %v", err)
		}

		// key with value as `*string`
		if t == "message" {
			val, err := dec.Token()
			if err != nil {
				log.Fatal(err)
			}
			if v, ok := val.(string); ok {
				output.Message = &v
			} else {
				return fmt.Errorf("expected Message to be of type string, got %T", val)
			}
		}
	}

	// end token
	endToken, err := dec.Token()
	if err != nil {
		return fmt.Errorf("Failed  to decode JSON: %v ", err)
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		return fmt.Errorf("invalid JSON , expected } at start of json")
	}
	return nil
}
