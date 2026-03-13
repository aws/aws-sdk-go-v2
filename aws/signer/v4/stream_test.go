package v4

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func mustHexDecode(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func TestStreamSigner_GetSignature(t *testing.T) {
	signer := NewStreamSigner(
		aws.Credentials{
			AccessKeyID:     "AKID",
			SecretAccessKey: "SECRET",
		},
		"transcribestreaming",
		"us-east-1",
		[]byte("foobarbazqux"),
	)

	messages := []struct {
		Name    string
		Headers []byte
		Payload []byte
		Time    time.Time
		Expect  []byte
	}{
		{
			Name:    "first message",
			Headers: []byte("foo"),
			Payload: []byte("foo"),
			Time:    time.Unix(10, 0),
			Expect:  mustHexDecode("e3a136405ec1152f62136742efd41f8639cd647a1cf21fab12aedb22cccd8999"),
		},
		{
			Name:    "second message",
			Headers: []byte("bar"),
			Payload: []byte("bar"),
			Time:    time.Unix(20, 0),
			Expect:  mustHexDecode("92a6ad677e839d3003672cb2cefc71ef80d42e8cb300fefce0807e8398c74068"),
		},
		{
			Name:    "third message",
			Headers: []byte("baz"),
			Payload: []byte("baz"),
			Time:    time.Unix(30, 0),
			Expect:  mustHexDecode("46d6660bd40f2ff5d6246a347895cabe7f4e89d9e1f4c6dad6dabdc801923be1"),
		},
		{
			Name:    "end of stream",
			Headers: []byte{},
			Payload: []byte{},
			Time:    time.Unix(30, 0),
			Expect:  mustHexDecode("492bcb3e21de447c85f8a9bef6c0eb833e992e405b9b3584be84edaccc5aad18"),
		},
	}

	for _, tt := range messages {
		t.Run(tt.Name, func(t *testing.T) {
			s, err := signer.GetSignature(context.Background(), tt.Headers, tt.Payload, tt.Time)
			if err != nil {
				t.Fatal(err)
			}

			expect := hex.EncodeToString(tt.Expect)
			actual := hex.EncodeToString(s)
			if expect != actual {
				t.Errorf("%v != %v", expect, actual)
			}
		})
	}
}
