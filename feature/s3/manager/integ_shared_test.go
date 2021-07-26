//go:build integration
// +build integration

package manager_test

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager/internal/integration"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var integConfig aws.Config

func init() {
	var err error

	integConfig, err = config.LoadDefaultConfig(context.Background(), config.WithDefaultRegion("us-west-2"))
	if err != nil {
		panic(err)
	}
}

var bucketName *string
var client *s3.Client

func TestMain(m *testing.M) {
	flag.Parse()
	flag.CommandLine.Visit(func(f *flag.Flag) {
		if !(f.Name == "run" || f.Name == "test.run") {
			return
		}
		value := f.Value.String()
		if value == `NONE` {
			os.Exit(0)
		}
	})

	client = s3.NewFromConfig(integConfig)
	bucketName = aws.String(integration.GenerateBucketName())
	if err := integration.SetupBucket(client, *bucketName, integConfig.Region); err != nil {
		panic(err)
	}

	var result int
	defer func() {
		if err := integration.CleanupBucket(client, *bucketName); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "S3 integration tests panicked,", r)
			result = 1
		}
		os.Exit(result)
	}()

	result = m.Run()
}

type dlwriter struct {
	buf []byte
}

func newDLWriter(size int) *dlwriter {
	return &dlwriter{buf: make([]byte, size)}
}

func (d dlwriter) WriteAt(p []byte, pos int64) (n int, err error) {
	if pos > int64(len(d.buf)) {
		return 0, io.EOF
	}

	written := 0
	for i, b := range p {
		if i >= len(d.buf) {
			break
		}
		d.buf[pos+int64(i)] = b
		written++
	}
	return written, nil
}

func validate(t *testing.T, key string, md5value string) {
	mgr := manager.NewDownloader(client)
	params := &s3.GetObjectInput{
		Bucket:       bucketName,
		Key:          &key,
		ChecksumMode: s3types.ChecksumModeEnabled,
	}

	w := newDLWriter(1024 * 1024 * 20)
	n, err := mgr.Download(context.Background(), w, params)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := md5value, fmt.Sprintf("%x", md5.Sum(w.buf[0:n])); e != a {
		t.Errorf("expect %s md5 value, got %s", e, a)
	}
}
