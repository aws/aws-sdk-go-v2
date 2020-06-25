package defaults

import (
	"net/http"
	"os"
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
)

func TestAddHostExecEnvUserAgentHander(t *testing.T) {
	cases := []struct {
		ExecEnv string
		Expect  string
	}{
		{ExecEnv: "Lambda", Expect: execEnvUAKey + "/Lambda"},
		{ExecEnv: "", Expect: ""},
		{ExecEnv: "someThingCool", Expect: execEnvUAKey + "/someThingCool"},
	}

	for i, c := range cases {
		os.Clearenv()
		os.Setenv(execEnvVar, c.ExecEnv)

		req := &aws.Request{
			HTTPRequest: &http.Request{
				Header: http.Header{},
			},
		}
		AddHostExecEnvUserAgentHander.Fn(req)

		if err := req.Error; err != nil {
			t.Fatalf("%d, expect no error, got %v", i, err)
		}

		if e, a := c.Expect, req.HTTPRequest.Header.Get("User-Agent"); e != a {
			t.Errorf("%d, expect %v user agent, got %v", i, e, a)
		}
	}
}
