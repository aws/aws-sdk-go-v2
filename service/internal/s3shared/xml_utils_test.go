package s3shared

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestGetResponseErrorCode(t *testing.T) {
	cases := map[string]struct {
		status                 int
		errorResponse          io.Reader
		expectedErrorCode      string
		expectedErrorMessage   string
		expectedErrorRequestID string
		expectedErrorHostID    string
	}{
		"standard xml error": {
			status: 400,
			errorResponse: bytes.NewReader([]byte(`<Error>
    <Type>Sender</Type>
    <Code>InvalidGreeting</Code>
    <Message>Hi</Message>
    <HostId>bar-id</HostId>
    <RequestId>foo-id</RequestId>
</Error>`)),
			expectedErrorCode:      "InvalidGreeting",
			expectedErrorMessage:   "Hi",
			expectedErrorRequestID: "foo-id",
			expectedErrorHostID:    "bar-id",
		},
		"no response body": {
			status:               400,
			errorResponse:        bytes.NewReader([]byte(``)),
			expectedErrorCode:    "BadRequest",
			expectedErrorMessage: "Bad Request",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ec, err := GetErrorResponseComponents(c.errorResponse, c.status)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if e, a := c.expectedErrorCode, ec.Code; !strings.EqualFold(e, a) {
				t.Fatalf("expected %v, got %v", e, a)
			}
			if e, a := c.expectedErrorMessage, ec.Message; !strings.EqualFold(e, a) {
				t.Fatalf("expected %v, got %v", e, a)
			}
			if e, a := c.expectedErrorRequestID, ec.RequestID; !strings.EqualFold(e, a) {
				t.Fatalf("expected %v, got %v", e, a)
			}
			if e, a := c.expectedErrorHostID, ec.HostID; !strings.EqualFold(e, a) {
				t.Fatalf("expected %v, got %v", e, a)
			}
		})
	}
}
