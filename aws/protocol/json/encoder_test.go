package json_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	v2Encoder "github.com/aws/aws-sdk-go-v2/aws/protocol/json"
	v1Encoder "github.com/aws/aws-sdk-go-v2/private/protocol/json"
	reflectEncoder "github.com/aws/aws-sdk-go-v2/private/protocol/json/jsonutil"
)

func TestEncoder(t *testing.T) {
	encoder := v2Encoder.NewEncoder()

	object := encoder.Object()

	object.Key("stringKey").String("stringValue")
	object.Key("integerKey").Integer(1024)
	object.Key("floatKey").Float(3.14)

	subObj := object.Key("foo").Object()

	subObj.Key("byteSlice").ByteSlice([]byte("foo bar"))
	subObj.Close()

	object.Close()

	e := []byte(`{"stringKey":"stringValue","integerKey":1024,"floatKey":3.14,"foo":{"byteSlice":"Zm9vIGJhcg=="}}`)
	if a := encoder.Bytes(); bytes.Compare(e, a) != 0 {
		t.Errorf("expected %+q, but got %+q", e, a)
	}

	if a := encoder.String(); string(e) != a {
		t.Errorf("expected %s, but got %s", e, a)
	}
}

func TestEncoderComparability(t *testing.T) {
	for i, operationCase := range testOperationCases {
		t.Run(fmt.Sprintf("Case%d", i), func(t *testing.T) {
			v2 := v2Encoder.NewEncoder()
			_ = MarshalTestOperationInputAWSJSON(operationCase, v2)
			v2Bytes := v2.Bytes()

			v1 := v1Encoder.NewEncoder()
			err := operationCase.MarshalFields(v1)
			if err != nil {
				t.Fatal(err)
			}
			v1Reader, err := v1.Encode()
			if err != nil {
				t.Fatal(err)
			}

			if v1Reader == nil {
				t.Logf("v1 encoder returns no reader, and v2 returned %+q", v2Bytes)
			} else {
				v1Bytes, err := ioutil.ReadAll(v1Reader)
				if err != nil {
					t.Fatal(err)
				}

				if bytes.Compare(v1Bytes, v2Bytes) != 0 {
					t.Fatalf("expected %+q, but got %+q", v1Bytes, v2Bytes)
				}
			}

			reflectBytes, err := reflectEncoder.BuildJSON(operationCase)
			if err != nil {
				t.Fatal(err)
			}

			if bytes.Compare(reflectBytes, v2Bytes) != 0 {
				t.Fatalf("expected %+q, but got %+q", reflectBytes, v2Bytes)
			}
		})
	}
}
