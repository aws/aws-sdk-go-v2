package s3

import (
	"encoding/xml"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
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

	var decoder *xml.Decoder
	decoder = xml.NewDecoder(txt)
	startTag, err := decoder.Token()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err.Error())
	}
	err = unmarshalErrorPrototype(&mockReq, decoder, startTag)

	if err == nil {
		t.Fatalf("Expected an error, got none")
	}

	if respErr, ok := err.(awserr.RequestFailure); ok {
		if e, a := "AccessDenied", respErr.Code(); e != a {
			t.Fatalf(" Response Error Code: Expected %v, got %v", e, a)
		}
		if e, a := "Access Denied", respErr.Message(); e != a {
			t.Fatalf(" Response Error Message: Expected %v, got %v", e, a)
		}
		if e, a := "0380AF11C6689A57", respErr.RequestID(); e != a {
			t.Fatalf(" Response Error Request Id: Expected %v, got %v", e, a)
		}
	} else {
		t.Fatalf("Expected error to be of type Request Failure")
	}
}
