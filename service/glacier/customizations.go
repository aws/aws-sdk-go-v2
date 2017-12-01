package glacier

import (
	"encoding/hex"
	"reflect"

	request "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

var (
	defaultAccountID = "-"
)

func init() {
	initRequest = func(c *Glacier, r *request.Request) {
		r.Handlers.Validate.PushFront(addAccountID)
		r.Handlers.Validate.PushFront(copyParams) // this happens first
		r.Handlers.Build.PushBack(addChecksum)
		r.Handlers.Build.PushBack(addAPIVersion)
	}
}

func copyParams(r *request.Request) {
	r.Params = awsutil.CopyOf(r.Params)
}

func addAccountID(r *request.Request) {
	if !r.ParamsFilled() {
		return
	}

	v := reflect.Indirect(reflect.ValueOf(r.Params))
	if f := v.FieldByName("AccountId"); f.IsNil() {
		f.Set(reflect.ValueOf(&defaultAccountID))
	}
}

func addChecksum(r *request.Request) {
	if r.Body == nil || r.HTTPRequest.Header.Get("X-Amz-Sha256-Tree-Hash") != "" {
		return
	}

	h := ComputeHashes(r.Body)
	hstr := hex.EncodeToString(h.TreeHash)
	r.HTTPRequest.Header.Set("X-Amz-Sha256-Tree-Hash", hstr)

	hLstr := hex.EncodeToString(h.LinearHash)
	r.HTTPRequest.Header.Set("X-Amz-Content-Sha256", hLstr)
}

func addAPIVersion(r *request.Request) {
	r.HTTPRequest.Header.Set("X-Amz-Glacier-Version", r.Metadata.APIVersion)
}
