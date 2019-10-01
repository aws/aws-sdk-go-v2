package s3manager_test

import (
	"math/rand"

	"github.com/aws/aws-sdk-go-v2/internal/sdkio"
)

var buf12MB = make([]byte, 1024*1024*12)
var buf2MB = make([]byte, 1024*1024*2)

var randBytes = func() []byte {
	b := make([]byte, 10*sdkio.MebiByte)

	// always returns len(b) and nil error
	_, _ = rand.Read(b)

	return b
}()

func getTestBytes(size int) []byte {
	if len(randBytes) >= size {
		return randBytes[:size]
	}

	b := append(randBytes, getTestBytes(size-len(randBytes))...)
	return b
}
