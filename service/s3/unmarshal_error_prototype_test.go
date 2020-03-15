package s3

import (
	"encoding/xml"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/s3err"
	"github.com/aws/aws-sdk-go-v2/internal/sdkio"
)

func TestUnmarshalError_Prototype(t *testing.T) {
	msg :=
		`<Error>
   <Code>AccessDenied</Code>
   <Message>Access Denied</Message>
   <RequestId>0380AF11C6689A57</RequestId>
   <HostId>c8bv/z7AAkLXxI8qsf/SYXTGmW0RHYI3o4hS+b7nVRKnGwhyrMsC3Hyf/3/3dNiQ3zJYF/ZYHXg=</HostId>
</Error>`

	txt := strings.NewReader(msg)
	mockReq := aws.Request{
		HTTPResponse: &http.Response{StatusCode: 400},
		RequestID:    "",
	}

	buff := make([]byte, 1024)
	ringBuff := sdkio.NewRingBuffer(buff)
	var decoder *xml.Decoder
	decoder = xml.NewDecoder(txt)
	startTag, err := decoder.Token()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err.Error())
	}

	if start, ok := startTag.(xml.StartElement); ok {
		err = unmarshalErrorPrototype(&mockReq, decoder, start, ringBuff)
		if err == nil {
			t.Fatalf("Expected an error, got none")
		}
	} else {
		t.Fatalf("Failed to get start element from invalid xml")
	}

	if respErr, ok := err.(*s3err.RequestFailure); ok {
		if e, a := "AccessDenied", respErr.Code(); e != a {
			t.Fatalf(" Response Error Code: Expected %v, got %v", e, a)
		}
		if e, a := "Access Denied", respErr.Message(); e != a {
			t.Fatalf(" Response Error Message: Expected %v, got %v", e, a)
		}
		if e, a := "0380AF11C6689A57", respErr.RequestID(); e != a {
			t.Fatalf(" Response Error Request Id: Expected %v, got %v", e, a)
		}
		if e, a := "c8bv/z7AAkLXxI8qsf/SYXTGmW0RHYI3o4hS+b7nVRKnGwhyrMsC3Hyf/3/3dNiQ3zJYF/ZYHXg=",
			respErr.HostID(); e != a {
			t.Fatalf("Response Error Host Id: Expected %v, got %v", e, a)
		}
	} else {
		t.Fatalf("Expected error to be of type s3 Request Failure, got %T", err)
	}
}
