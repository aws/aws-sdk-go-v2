package awsutil_test

import (
	"testing"

	"github.com/Enflick/aws-sdk-go-v2/internal/awsutil"
	"github.com/Enflick/smithy-go/ptr"
)

type testStruct struct {
	Field1 string
	Field2 *string
	Field3 []byte `sensitive:"true"`
	Value  []string
}

func TestStringValue(t *testing.T) {
	cases := map[string]struct {
		Value  interface{}
		Expect string
	}{
		"general": {
			Value: testStruct{
				Field1: "abc123",
				Field2: ptr.String("abc123"),
				Field3: []byte("don't show me"),
				Value: []string{
					"first",
					"second",
				},
			},
			Expect: `{
  Field1: "abc123",
  Field2: "abc123",
  Field3: <sensitive>,
  Value: ["first","second"],

}`,
		},
	}

	for d, c := range cases {
		t.Run(d, func(t *testing.T) {
			actual := awsutil.StringValue(c.Value)
			if e, a := c.Expect, actual; e != a {
				t.Errorf("expect:\n%v\nactual:\n%v\n", e, a)
			}
		})
	}
}
