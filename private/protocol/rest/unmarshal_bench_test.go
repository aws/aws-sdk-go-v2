package rest

import (
	"bytes"
	"github.com/jviney/aws-sdk-go-v2/aws"
	"io/ioutil"
	"net/http"
	"testing"
)

type DataOutput struct {
	_ struct{} `type:"structure"`

	HeaderEnum string `location:"header" locationName:"x-amz-enum" type:"string" enum:"true"`

	StatusCode int64 `location:"statusCode"`
}

type DataOutput2 struct {
	_ struct{} `type:"structure" payload:"TestString"`

	TestString *string `type:"string"`
}

func BenchmarkRESTUnmarshalMeta(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnmarshalMeta(getRESTMeta())
	}
}

func BenchmarkRESTUnmarshal_Short(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unmarshal(getRESTResponseShortREST())
	}
}

func BenchmarkRESTUnmarshal_Long(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unmarshal(getRESTResponseLongREST())
	}
}

func getRESTMeta() *aws.Request {
	output := DataOutput{}
	req := aws.Request{Data: &output, HTTPResponse: &http.Response{StatusCode: 200, Header: http.Header{}}}
	req.HTTPResponse.Header.Set("x-amz-enum", "baz")
	return &req
}

func getRESTResponseShortREST() *aws.Request {
	buf := bytes.NewReader([]byte("test_string"))
	req := aws.Request{Data: &DataOutput2{}, HTTPResponse: &http.Response{Body: ioutil.NopCloser(buf)}}
	return &req
}

func getRESTResponseLongREST() *aws.Request {
	buf := bytes.NewReader([]byte("O8sxCe0h66aCwPY2gLAhmKPFVP9c5b0o6iViHf4G58JT5ugSPX8VvFj7KAk5ykjR4Xlq7MAzGFcaFRFMiSwDuSMHNCO7KiiE3Eud38zkq62uwlkLteORA1EZLfbnTO2QigizYAcbJQkTf3mDpJdxDOkpinAMOUTA3uGq1VqG2lrVu4D0vnCXUh6L8HzQ7vGRM0ydWpCdhefedJSDOuNr4PWvESEPFS3gAZYFFYhdJGyeF9NfyqrmAApog25pERNrd7cY5wkahx1bCh4wqkUchlwg795rbKpYIO4feYChnp94I0rtX7RNUfEksRuamEgvcdJODV8yWdI2TiBobkCUwB1avVixmo5frWIXPvQVPls1sfxlOKpgsP9dm3dKaxalUWcQVG7Nv5b3TXbh6IeXKivYtwEaHo2T2pFt39cfty6k2GEFAjxjDXBrw4ljsZv0HcA36XNvWOCoRblISa7CMV9J0QjQ3W1D0MfXqkygojpxcxgKNbyNR28oZmX3H0ZY4kz4NCj34tg0jCHzpj9KExZGLvinWF4IuR4gKi8usXp6j7q7ZqX2qf7bx9tWs3Ci3N7lnq5NPrBVSeHGyNmBHdVzFEPBFiFPyiFkLtSQFCS054WaDmR5m2DJK9gmY9gMxmnMDk6q2Vzdp1j8mgzFsjbhQWiKMuwIel9YKqDzijlxaNSJKd9UWvNs4CFewlzq4JP7VC7vfSkOmI0RmeTirey3KrCLzHCB1bizMhSxwkZZz745vvQMfYXZ14vw4KEsqQVKTldg74EQuLiEHjpxcJ5PPh31FxkkEPvq4AD9JcfB8b4Uead3ij20dCZ7qhDmKrJ2eTbP5mYpZdiymETpRtCtA6mbcTxeJQrtlJbUgJjVOixhTZOS6Qe3GTsLES8lKasGQPFWbPPE44Ot2vZJtmeFZfLyCTtzIIHNx00Zzv6aHh1I8rstQef3rKQ9yExELk4J6ckCmaPnBGsDX1U91hVrx2LsJ1h8dUNXEQzE"))
	req := aws.Request{Data: &DataOutput2{}, HTTPResponse: &http.Response{Body: ioutil.NopCloser(buf)}}
	return &req
}
