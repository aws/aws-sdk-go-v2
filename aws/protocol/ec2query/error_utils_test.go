package ec2query

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestGetResponseErrorCode(t *testing.T) {
	cases := map[string]struct {
		errorResponse          io.Reader
		expectedErrorCode      string
	}{
		"Invalid Greeting": {
			errorResponse: bytes.NewReader([]byte(`<Response>
			    <Errors>
			        <Error>
			            <Code>InvalidGreeting</Code>
			            <Message>Hi</Message>
			        </Error>
			    </Errors>
			    <RequestId>foo-id</RequestId>
			</Response>`)),
			expectedErrorCode:      "InvalidGreeting",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			errorcode, err := GetResponseErrorCode(c.errorResponse)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if e, a := c.expectedErrorCode, errorcode; !strings.EqualFold(e, a) {
				t.Fatalf("expected %v, got %v", e, a)
			}
		})
	}
}
