package s3shared

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestGetResponseErrorCode(t *testing.T) {
	cases := map[string]struct {
		isS3Service            bool
		status                 int
		errorResponse          io.Reader
		expectedErrorCode      string
		expectedErrorMessage   string
		expectedErrorRequestID string
		expectedErrorHostID    string
	}{
		"standard xml error": {
			isS3Service: true,
			status:      400,
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
		"s3 no response body": {
			isS3Service:          true,
			status:               400,
			errorResponse:        bytes.NewReader([]byte(``)),
			expectedErrorCode:    "BadRequest",
			expectedErrorMessage: "Bad Request",
		},
		"s3control no response body": {
			isS3Service:   false,
			status:        400,
			errorResponse: bytes.NewReader([]byte(``)),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ec, err := GetErrorResponseComponents(c.errorResponse, c.status, c.isS3Service)
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
