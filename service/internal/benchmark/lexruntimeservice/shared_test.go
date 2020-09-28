package benchmark

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"testing"

	smithyClient "github.com/aws/aws-sdk-go-v2/service/lexruntimeservice"
	"github.com/awslabs/smithy-go/middleware"
)

var (
	disableSmithySigning bool
)

func init() {
	flag.BoolVar(&disableSmithySigning, "disable-smithy-signing", false,
		"Instructs the test to be run with smithy signing turned off.")
}

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func loadTestData(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var b bytes.Buffer
	if _, err = io.Copy(&b, f); err != nil {
		return nil, fmt.Errorf("failed to read test data, %v", err)
	}

	return b.Bytes(), nil
}

func removeSmithySigner(options *smithyClient.Options) {
	options.APIOptions = append(options.APIOptions, func(stack *middleware.Stack) error {
		stack.Finalize.Remove("SigV4SignHTTPRequestMiddleware")
		stack.Finalize.Remove("SigV4ContentSHA256HeaderMiddleware")
		stack.Finalize.Remove("ComputePayloadSHA256Middleware")
		stack.Finalize.Remove("SigV4UnsignedPayloadMiddleware")
		return nil
	})
}
