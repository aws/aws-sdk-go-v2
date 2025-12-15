package transfermanager

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type objectkeyFilter struct {
	keyword string
}

func (of *objectkeyFilter) FilterObject(object s3types.Object) bool {
	if strings.Contains(aws.ToString(object.Key), of.keyword) {
		return false
	}
	return true
}

type objectkeyCallback struct {
	keyword string
}

func (oc *objectkeyCallback) UpdateRequest(in *GetObjectInput) {
	if key := aws.ToString(in.Key); key == oc.keyword {
		in.Key = aws.String(key + "gotyou")
	}
}

func TestDownloadDirectory(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "testdata")

	cases := map[string]struct {
		destination             string
		keyPrefix               string
		objectsLists            [][]s3types.Object
		continuationTokens      []string
		filter                  ObjectFilter
		concurrency             int
		callback                GetRequestCallback
		failurePolicy           DownloadDirectoryFailurePolicy
		getobjectFn             func(*s3testing.TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error)
		expectTokens            []string
		expectKeys              []string
		expectFiles             []string
		expectErr               string
		expectObjectsDownloaded int64
		expectObjectsFailed     int64
		listenerValidationFn    func(*testing.T, *mockDirectoryListener, any, any, error)
	}{
		"single object": {
			destination: "single-object",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("foo/bar"),
					},
				},
			},
			expectTokens:            []string{""},
			expectKeys:              []string{"foo/bar"},
			expectFiles:             []string{"foo/bar"},
			expectObjectsDownloaded: 1,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 1)
			},
		},
		"multiple objects": {
			destination: "multiple-objects",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("foo/bar"),
					},
					{
						Key: aws.String("baz"),
					},
					{
						Key: aws.String("foo/zoo/bar"),
					},
					{
						Key: aws.String("foo/zoo/oii/bababoii"),
					},
				},
			},
			expectTokens:            []string{""},
			expectKeys:              []string{"foo/bar", "baz", "foo/zoo/bar", "foo/zoo/oii/bababoii"},
			expectFiles:             []string{"foo/bar", "baz", "foo/zoo/bar", "foo/zoo/oii/bababoii"},
			expectObjectsDownloaded: 4,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 4)
			},
		},
		"multiple objects paginated": {
			destination: "multiple-objects-paginated",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("foo/bar"),
					},
					{
						Key: aws.String("baz"),
					},
				},
				{
					{
						Key: aws.String("foo/zoo/bar"),
					},
					{
						Key: aws.String("foo/zoo/oii/bababoii"),
					},
				},
				{
					{
						Key: aws.String("foo/zoo/baz"),
					},
					{
						Key: aws.String("foo/zoo/oii/yee"),
					},
				},
			},
			continuationTokens:      []string{"token1", "token2"},
			expectTokens:            []string{"", "token1", "token2"},
			expectKeys:              []string{"foo/bar", "baz", "foo/zoo/bar", "foo/zoo/oii/bababoii", "foo/zoo/baz", "foo/zoo/oii/yee"},
			expectObjectsDownloaded: 6,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 6)
			},
		},
		"multiple objects containing folder object": {
			destination: "multiple-objects-with-folder-object",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("foo/bar"),
					},
					{
						Key: aws.String("baz"),
					},
					{
						Key: aws.String("foo/zoo/"),
					},
				},
			},
			expectTokens:            []string{""},
			expectKeys:              []string{"foo/bar", "baz"},
			expectFiles:             []string{"foo/bar", "baz"},
			expectObjectsDownloaded: 2,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 2)
			},
		},
		"single object named with keyprefix": {
			destination: "single-object-named-with-keyprefix",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("a"),
					},
				},
			},
			keyPrefix:               "a",
			expectTokens:            []string{""},
			expectKeys:              []string{"a"},
			expectFiles:             []string{"a"},
			expectObjectsDownloaded: 1,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 1)
			},
		},
		"multiple objects with keyprefix without delimiter suffix": {
			destination: "multiple-objects-with-keyprefix-no-delimiter",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("a/"),
					},
					{
						Key: aws.String("a/b"),
					},
					{
						Key: aws.String("ad"),
					},
					{
						Key: aws.String("ab/c"),
					},
					{
						Key: aws.String("ae"),
					},
				},
			},
			keyPrefix:               "a",
			expectTokens:            []string{""},
			expectKeys:              []string{"a/b", "ad", "ab/c", "ae"},
			expectFiles:             []string{"b", "ad", "ab/c", "ae"},
			expectObjectsDownloaded: 4,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 4)
			},
		},
		"multiple objects with keyprefix with default delimiter suffix": {
			destination: "multiple-objects-with-keyprefix-default-delimiter",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("a/"),
					},
					{
						Key: aws.String("a/b"),
					},
					{
						Key: aws.String("a/c"),
					},
					{
						Key: aws.String("ad"),
					},
					{
						Key: aws.String("ab/c/d"),
					},
					{
						Key: aws.String("ab/c/e"),
					},
				},
			},
			keyPrefix:               "a/",
			expectTokens:            []string{""},
			expectKeys:              []string{"a/b", "a/c", "ad", "ab/c/d", "ab/c/e"},
			expectFiles:             []string{"b", "c", "ad", "ab/c/d", "ab/c/e"},
			expectObjectsDownloaded: 5,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 5)
			},
		},
		"error when path resolved from objects key out of destination scope": {
			destination: "error-bucket",
			concurrency: 1,
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("a/"),
					},
					{
						Key: aws.String("a/b"),
					},
					{
						Key: aws.String(filepath.Join("a", "..", "..", "d")),
					},
					{
						Key: aws.String("a/c"),
					},
				},
			},
			expectErr: "outside of destination",
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				// only validate failure listener since start listener
				// might never be triggerred if the error response is returned first
				l.expectFailed(t, in, err)
			},
		},
		"multiple objects with filter applied": {
			destination: "multiple-objects-with-filter-applied",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("foo/bar"),
					},
					{
						Key: aws.String("baz"),
					},
					{
						Key: aws.String("foo/zoo/bar"),
					},
					{
						Key: aws.String("foo/zoo/oii/bababoii"),
					},
				},
			},
			filter:                  &objectkeyFilter{"bababoii"},
			expectTokens:            []string{""},
			expectKeys:              []string{"foo/bar", "baz", "foo/zoo/bar"},
			expectFiles:             []string{"foo/bar", "baz", "foo/zoo/bar"},
			expectObjectsDownloaded: 3,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 3)
			},
		},
		"multiple objects with keyprefix and filter": {
			destination: "multiple-objects-with-keyprefix-and-filter",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("a/"),
					},
					{
						Key: aws.String("a/b"),
					},
					{
						Key: aws.String("ad"),
					},
					{
						Key: aws.String("ab/c"),
					},
					{
						Key: aws.String("ae"),
					},
				},
			},
			keyPrefix:               "a",
			filter:                  &objectkeyFilter{"e"},
			expectTokens:            []string{""},
			expectKeys:              []string{"a/b", "ad", "ab/c"},
			expectFiles:             []string{"b", "ad", "ab/c"},
			expectObjectsDownloaded: 3,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 3)
			},
		},
		"multiple objects with keyprefix and request callback": {
			destination: "multiple-objects-with-keyprefix-and-callback",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("a/"),
					},
					{
						Key: aws.String("a/b"),
					},
					{
						Key: aws.String("ad"),
					},
					{
						Key: aws.String("ab/c"),
					},
					{
						Key: aws.String("ae"),
					},
				},
			},
			keyPrefix:               "a",
			callback:                &objectkeyCallback{"ad"},
			expectTokens:            []string{""},
			expectKeys:              []string{"a/b", "adgotyou", "ab/c", "ae"},
			expectFiles:             []string{"b", "ad", "ab/c", "ae"},
			expectObjectsDownloaded: 4,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 4)
			},
		},
		"error when getting object": {
			destination: "error-bucket",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("foo/bar"),
					},
					{
						Key: aws.String("baz"),
					},
				},
				{
					{
						Key: aws.String("foo/zoo/bar"),
					},
					{
						Key: aws.String("foo/zoo/oii/bababoii"),
					},
				},
				{
					{
						Key: aws.String("foo/zoo/baz"),
					},
					{
						Key: aws.String("foo/zoo/oii/yee"),
					},
				},
			},
			concurrency:        1,
			continuationTokens: []string{"token1", "token2"},
			getobjectFn: func(c *s3testing.TransferManagerLoggingClient, in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				if aws.ToString(in.Key) == "foo/zoo/bar" {
					return nil, fmt.Errorf("mocking error")
				}
				return &s3.GetObjectOutput{
					Body:          ioutil.NopCloser(bytes.NewReader(c.Data)),
					ContentLength: aws.Int64(int64(len(c.Data))),
					PartsCount:    aws.Int32(c.PartsCount),
					ETag:          aws.String(etag),
				}, nil
			},
			expectErr: "mocking error",
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectFailed(t, in, err)
			},
		},
		"specified getting object failure ignored by failure policy": {
			destination: "error-ignored",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("fo/"),
					},
					{
						Key: aws.String("baz"),
					},
				},
				{
					{
						Key: aws.String("foo/zoo/bar"),
					},
					{
						Key: aws.String("foo/zoo/oii/bababoii"),
					},
				},
				{
					{
						Key: aws.String("foo/zoo/baz"),
					},
					{
						Key: aws.String("foo/zoo/oii/yee"),
					},
				},
			},
			concurrency:        1,
			failurePolicy:      IgnoreDownloadFailurePolicy{},
			continuationTokens: []string{"token1", "token2"},
			getobjectFn: func(c *s3testing.TransferManagerLoggingClient, in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				if key := aws.ToString(in.Key); key == "foo/zoo/bar" || key == "baz" {
					return nil, fmt.Errorf("mocking error")
				}
				return &s3.GetObjectOutput{
					Body:          ioutil.NopCloser(bytes.NewReader(c.Data)),
					ContentLength: aws.Int64(int64(len(c.Data))),
					PartsCount:    aws.Int32(c.PartsCount),
					ETag:          aws.String(etag),
				}, nil
			},
			expectTokens:            []string{"", "token1", "token2"},
			expectKeys:              []string{"baz", "foo/zoo/bar", "foo/zoo/oii/bababoii", "foo/zoo/baz", "foo/zoo/oii/yee"},
			expectFiles:             []string{"foo/zoo/oii/bababoii", "foo/zoo/baz", "foo/zoo/oii/yee"},
			expectObjectsDownloaded: 3,
			expectObjectsFailed:     2,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 3)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			s3Client, params := s3testing.NewDownloadDirectoryClient()
			s3Client.ListObjectsData = c.objectsLists
			s3Client.ContinuationTokens = c.continuationTokens
			if c.getobjectFn == nil {
				s3Client.GetObjectFn = s3testing.PartGetObjectFn
			} else {
				s3Client.GetObjectFn = c.getobjectFn
			}
			s3Client.Data = make([]byte, 0)
			s3Client.PartsCount = 1
			mgr := New(s3Client, Options{})

			dstPath := filepath.Join(root, c.destination)
			defer os.RemoveAll(dstPath)

			req := &DownloadDirectoryInput{
				Bucket:        aws.String("mock-bucket"),
				Destination:   aws.String(dstPath),
				KeyPrefix:     nzstring(c.keyPrefix),
				Filter:        c.filter,
				Callback:      c.callback,
				FailurePolicy: c.failurePolicy,
			}
			listener := &mockDirectoryListener{}

			resp, err := mgr.DownloadDirectory(context.Background(), req, func(o *Options) {
				o.DirectoryProgressListeners.Register(listener)
				if c.concurrency > 0 {
					o.DirectoryConcurrency = c.concurrency
				}
			})

			if err != nil {
				if c.expectErr == "" {
					t.Fatalf("expect not error, got %v", err)
				} else if e, a := c.expectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %s error message to be in %s", e, a)
				}
			} else if c.expectErr != "" {
				t.Fatalf("expect error %s, got none", c.expectErr)
			}

			if c.listenerValidationFn != nil {
				c.listenerValidationFn(t, listener, req, resp, err)
			}

			if c.expectErr != "" {
				return
			}

			if e, a := c.expectObjectsDownloaded, resp.ObjectsDownloaded; e != a {
				t.Errorf("expect %d objects downloaded, got %d", e, a)
			}
			if e, a := c.expectObjectsFailed, resp.ObjectsFailed; e != a {
				t.Errorf("expect %d objects failed, got %d", e, a)
			}

			var actualTokens []string
			var actualKeys []string
			for _, param := range *params {
				if input, ok := param.(*s3.ListObjectsV2Input); ok {
					actualTokens = append(actualTokens, aws.ToString(input.ContinuationToken))
				} else if input, ok := param.(*s3.GetObjectInput); ok {
					actualKeys = append(actualKeys, aws.ToString(input.Key))
				} else {
					t.Fatalf("error when casting captured inputs")
				}
			}

			if e, a := c.expectTokens, actualTokens; !reflect.DeepEqual(e, a) {
				t.Errorf("expect continuation tokens to be %v, got %v", e, a)
			}

			sort.Strings(actualKeys)
			sort.Strings(c.expectKeys)
			if e, a := c.expectKeys, actualKeys; !reflect.DeepEqual(e, a) {
				t.Errorf("expect downloaded keys to be %v, got %v", e, a)
			}

			for _, file := range c.expectFiles {
				path := filepath.Join(dstPath, strings.ReplaceAll(file, "/", string(os.PathSeparator)))
				_, err := os.Stat(path)
				if os.IsNotExist(err) {
					t.Errorf("expect %s to be downloaded, got none", path)
				}
			}
		})
	}
}

func TestDownloadDirectoryObjectsTransferred(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "testdata")
	cases := map[string]struct {
		destination        string
		objectsLists       [][]s3types.Object
		continuationTokens []string
		objectsCount       []int64
	}{
		"single object": {
			destination: "single-object",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("foo/bar"),
					},
				},
			},
			objectsCount: []int64{1},
		},
		"multiple objects": {
			destination: "multiple-objects",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("foo/bar"),
					},
					{
						Key: aws.String("baz"),
					},
					{
						Key: aws.String("foo/zoo/bar"),
					},
					{
						Key: aws.String("foo/zoo/oii/bababoii"),
					},
				},
			},
			objectsCount: []int64{1, 2, 3, 4},
		},
		"multiple objects paginated": {
			destination: "multiple-objects-with-keyprefix-delimiter-filter-callback",
			objectsLists: [][]s3types.Object{
				{
					{
						Key: aws.String("a/"),
					},
					{
						Key: aws.String("a/b"),
					},
					{
						Key: aws.String("a/b"),
					},
				},
				{
					{
						Key: aws.String("a/foo/bar"),
					},
					{
						Key: aws.String("ac"),
					},
					{
						Key: aws.String("ac@d/e"),
					},
				},
				{
					{
						Key: aws.String("a/k.b"),
					},
				},
			},
			continuationTokens: []string{"token1", "token2"},
			objectsCount:       []int64{1, 2, 3, 4, 5, 6},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			s3Client, _ := s3testing.NewDownloadDirectoryClient()
			s3Client.ListObjectsData = c.objectsLists
			s3Client.ContinuationTokens = c.continuationTokens
			s3Client.GetObjectFn = s3testing.PartGetObjectFn

			s3Client.Data = make([]byte, 0)
			s3Client.PartsCount = 1
			mgr := New(s3Client, Options{})

			dstPath := filepath.Join(root, c.destination)
			defer os.RemoveAll(dstPath)

			req := &DownloadDirectoryInput{
				Bucket:      aws.String("mock-bucket"),
				Destination: aws.String(dstPath),
			}
			listener := &mockDirectoryListener{}

			_, err := mgr.DownloadDirectory(context.Background(), req, func(o *Options) {
				o.DirectoryProgressListeners.Register(listener)
				o.DirectoryConcurrency = 1
			})
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			listener.expectObjectsTransferred(t, c.objectsCount...)
		})
	}
}

func TestDownloadDirectoryWithContextCanceled(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "testdata")
	dstPath := filepath.Join(root, "context-canceled")
	defer os.RemoveAll(dstPath)
	c := s3.New(s3.Options{
		UsePathStyle: true,
		Region:       "mock-region",
	})
	u := New(c, Options{})

	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{})}
	ctx.Error = fmt.Errorf("context canceled")
	close(ctx.DoneCh)

	_, err := u.DownloadDirectory(ctx, &DownloadDirectoryInput{
		Bucket:      aws.String("mock-bucket"),
		Destination: aws.String(dstPath),
	})
	if err == nil {
		t.Fatalf("expect error, got nil")
	}

	if e, a := "canceled", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expected error message to contain %q, but did not %q", e, a)
	}
}
