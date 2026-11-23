package eventstream

import (
	"errors"
	"testing"
)

func TestValidateMessageLen(t *testing.T) {
	// maxBody is the most headers and payload bytes, together, that fit.
	const maxBody = maxMsgLen - minMsgLen

	cases := map[string]struct {
		headersLen, payloadLen uint64
		wantErr                bool
	}{
		"empty":                       {0, 0, false},
		"max payload":                 {0, maxBody, false},
		"max headers and payload":     {testHeadersLen, maxBody - testHeadersLen, false},
		"payload one too long":        {0, maxBody + 1, true},
		"headers push payload over":   {1, maxBody, true},
		"payload wraps 32-bit length": {0, maxMsgLen + 1, true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := validateMessageLen(c.headersLen, c.payloadLen)
			if c.wantErr {
				if !errors.As(err, &LengthError{}) {
					t.Fatalf("expect LengthError, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
		})
	}
}
