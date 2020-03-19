package apigateway

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/internal/sdkio"
	restlegacy "github.com/aws/aws-sdk-go-v2/private/protocol/rest"
)

// protoCreateAPIKeyUnmarshaler defines unmarshaler forProtoCreateAPIKey Operation
type protoCreateAPIKeyUnmarshaler struct {
	output *CreateApiKeyOutput
}

// unmarshalOperation is the top level method used with a handler stack to unmarshal an operation response
// This method calls appropriate unmarshal shape functions as per the output shape and protocol used by the service.
func (u protoCreateAPIKeyUnmarshaler) unmarshalOperation(r *aws.Request) {
	if isRequestError(r) {
		unmarshalError(r)
		return
	}
	restlegacy.UnmarshalMeta(r)
	// jsonlegacy.UnmarshalJSON(u.output, r.Body)
	unmarshalProtoCreateAPIKeyAWSJSON(u.output, r)
}

func unmarshalProtoCreateAPIKeyAWSJSON(output *CreateApiKeyOutput, r *aws.Request) {
	buff := make([]byte, 1024)
	readBuff := make([]byte, 1024)
	ringBuff := sdkio.NewRingBuffer(buff)
	body := io.TeeReader(r.HTTPResponse.Body, ringBuff)
	dec := json.NewDecoder(body)

	// start token
	startToken, err := dec.Token()
	if err == io.EOF {
		fmt.Print("Empty Response")
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	if t, ok := startToken.(json.Delim); !ok {
		if t.String() != "{" {
			log.Fatalf("Invalid JSON : %v %v", t, dec.Buffered())
		}
	}

	for dec.More() {
		// fetch token for key
		t, err := dec.Token()
		if err != nil {
			log.Fatal(err)
		}

		// location name : `value` key with value as `*string`
		if t == "value" {
			val, err := dec.Token()
			if err != nil {
				log.Fatal(err)
			}
			switch v := val.(type) {
			case string:
				output.Value = &v
			case nil:
			}
		}

		// location name : `name` key with value as `*string`
		if t == "name" {
			val, err := dec.Token()
			if err != nil {
				log.Fatal(err)
			}
			switch v := val.(type) {
			case string:
				output.Name = &v
			case nil:
			}
		}

		// location name : `description` key with value as `*string`
		if t == "description" {
			val, err := dec.Token()
			if err != nil {
				log.Fatal(err)
			}
			switch v := val.(type) {
			case string:
				output.Description = &v
			case nil:
			}
		}

		// location name : `customerId` key with value as `*string`
		if t == "customerId" {
			val, err := dec.Token()
			if err != nil {
				log.Fatal(err)
			}
			switch v := val.(type) {
			case string:
				output.CustomerId = &v
			case nil:
			}
		}

		// location name : `id` key with value as `*string`
		if t == "id" {
			val, err := dec.Token()
			if err != nil {
				log.Fatal(err)
			}
			switch v := val.(type) {
			case string:
				output.Id = &v
			case nil:
			}
		}

		// location name : `enabled` key with value as `*boolean`
		if t == "enabled" {
			val, err := dec.Token()
			if err != nil {
				log.Fatal(err)
			}
			switch v := val.(type) {
			case bool:
				output.Enabled = &v
			case nil:
			}
		}

		// location name : `stageKeys` key with value as `[]string`
		if t == "stageKeys" {
			// create []string as modeled for the member shape
			list := make([]string, 0)
			// start of the list
			startToken, err := dec.Token()
			if err == io.EOF {
				fmt.Print("Empty List")
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			if t, ok := startToken.(json.Delim); !ok || t.String() != "[" {
				log.Fatalf("Expected json delimiter at the start of list, found %v instead", t)
			}

			for dec.More() {
				token, err := dec.Token()
				if err != nil {
					log.Fatal(err)
				}
				// based on struct Stage
				switch v := token.(type) {
				case string:
					val := v
					list = append(list, val)
				default:
					log.Fatalf("Expected %v, to be of type string, got %T", v, v)
				}
			}

			// end of the list
			endToken, err := dec.Token()
			if err != nil {
				log.Fatal(err)
			}
			if t, ok := endToken.(json.Delim); !ok || t.String() != "]" {
				log.Fatalf("Invalid JSON : %v %v", t, dec.Buffered())
			}

			// Attach the de-serialized list of string to output shape
			output.StageKeys = list
		}

		// location name : `tags` key with value as `map`
		if t == "tags" {
			// create []string as modeled for the member shape
			m := make(map[string]string, 0)
			// start of the list
			startToken, err := dec.Token()
			if err == io.EOF {
				fmt.Print("Empty List")
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			if t, ok := startToken.(json.Delim); !ok || t.String() != "{" {
				log.Fatalf("Expected json delimiter at the start of list, found %v instead", t)
			}

			// decoder
			for dec.More() {
				token, err := dec.Token()
				if err != nil {
					log.Fatal(err)
				}

				// based on struct Stage
				switch key := token.(type) {
				case string:
					val, err := dec.Token()
					if err != nil {
						log.Fatal(err)
					}
					if v, ok := val.(string); !ok {
						ringBuff.Read(readBuff)
						log.Fatalf("Invalid JSON. Here's an error snapshot %v", readBuff)
					} else {
						m[key] = v
					}
				default:
					log.Fatalf("Expected %v, to be of type string, got %T", key, key)
				}
			}

			// end of the list
			endToken, err := dec.Token()
			if err != nil {
				log.Fatal(err)
			}
			if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
				log.Fatalf("Invalid JSON : %v %v", t, dec.Buffered())
			}

			// Attach de-serialized Tag to output
			output.Tags = m
		}

		// location name : `createdDate` with value as `timestamp`
		if t == "createdDate" {
			val, err := dec.Token()
			if err != nil {
				log.Fatal(err)
			}
			switch v := val.(type) {
			case float64:
				time := time.Unix(int64(v), 0).UTC()
				output.CreatedDate = &time
			case nil:
			}
		}

		// location name : `lastUpdatedDate` with value as `timestamp`
		if t == "lastUpdatedDate" {
			val, err := dec.Token()
			if err != nil {
				log.Fatal(err)
			}
			switch v := val.(type) {
			case float64:
				time := time.Unix(int64(v), 0).UTC()
				output.LastUpdatedDate = &time
			case nil:
			}
		}

	}
	// end of the list
	endToken, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	if t, ok := endToken.(json.Delim); !ok || t.String() != "}" {
		log.Fatalf("Invalid JSON : %v %v", t, dec.Buffered())
	}
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

// unmarshalError unmarshal's an error response.
// some service may have custom error handling
// here we do not handle modelled exceptions.
func unmarshalError(req *aws.Request) {
	defer req.HTTPResponse.Body.Close()
	defer io.Copy(ioutil.Discard, req.HTTPResponse.Body)
	buff := make([]byte, 1024)
	readBuff := make([]byte, 1024)
	ringBuff := sdkio.NewRingBuffer(buff)
	body := io.TeeReader(req.HTTPResponse.Body, ringBuff)
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		ringBuff.Read(readBuff)
		req.Error = awserr.New("SerializationError",
			fmt.Sprintf("failed reading JSON error response, Here's a snapshot %s", readBuff), err)
		return
	}
	if len(bodyBytes) == 0 {
		req.Error = awserr.NewRequestFailure(
			awserr.New("SerializationError", req.HTTPResponse.Status, nil),
			req.HTTPResponse.StatusCode,
			"",
		)
		return
	}
	var jsonErr jsonErrorResponse
	if err := json.Unmarshal(bodyBytes, &jsonErr); err != nil {
		req.Error = awserr.New("SerializationError", "failed decoding JSON RPC error response", err)
		return
	}

	codes := strings.SplitN(jsonErr.Code, "#", 2)
	req.Error = awserr.NewRequestFailure(
		awserr.New(codes[len(codes)-1], jsonErr.Message, nil),
		req.HTTPResponse.StatusCode,
		req.RequestID,
	)
}

type jsonErrorResponse struct {
	Code    string `json:"__type"`
	Message string `json:"message"`
}
