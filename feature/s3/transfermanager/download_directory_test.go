package transfermanager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting"

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
					Body:          io.NopCloser(bytes.NewReader(c.Data)),
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
					Body:          io.NopCloser(bytes.NewReader(c.Data)),
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
			mgr := New(s3Client)

			dstPath := filepath.Join("testdata", c.destination)
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
					o.Concurrency = c.concurrency
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
			mgr := New(s3Client)

			dstPath := filepath.Join("testdata", c.destination)
			defer os.RemoveAll(dstPath)

			req := &DownloadDirectoryInput{
				Bucket:      aws.String("mock-bucket"),
				Destination: aws.String(dstPath),
			}
			listener := &mockDirectoryListener{}

			_, err := mgr.DownloadDirectory(context.Background(), req, func(o *Options) {
				o.DirectoryProgressListeners.Register(listener)
				o.Concurrency = 1
			})
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			listener.expectObjectsTransferred(t, c.objectsCount...)
		})
	}
}

func TestDownloadDirectoryWithContextCanceled(t *testing.T) {
	dstPath := filepath.Join("testdata", "context-canceled")
	defer os.RemoveAll(dstPath)
	c := s3.New(s3.Options{
		UsePathStyle: true,
		Region:       "mock-region",
	})
	u := New(c)

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

// createdFileCapture records map of file path corresponding to their *os.File handle that
// DownloadDirectory creates via createFileFn, so a test can assert after the call returns that
// each file was either closed or removed if downloaded failed. A
// removed file path returns os.ErrNotExist from os.Stat, and isFileClosed
// reports whether a still-present file's handle was closed, so a file that is
// neither removed nor closed indicates the handle leaked. It is safe for concurrent
// use by the download workers.
type createdFileCapture struct {
	mu    sync.Mutex
	files map[string]*os.File
}

func (c *createdFileCapture) capture(path string, f *os.File) {
	c.mu.Lock()
	c.files[path] = f
	c.mu.Unlock()
}

func (c *createdFileCapture) count() int {
	return len(c.files)
}

// notRemovedOrStillOpen returns the names of captured files that are removed when download
// failed or closed after download.
func (c *createdFileCapture) removedOrClosed() (removed, closed []string) {
	for path, file := range c.files {
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			removed = append(removed, path)
		} else if isFileClosed(file) {
			closed = append(closed, path)
		}
	}
	sort.Strings(closed)
	sort.Strings(removed)
	return
}

// TestDownloadDirectoryClosesCreatedFiles verifies that DownloadDirectory
// closes every file handle it creates before returning - on the success path
// and when individual object downloads fail and are ignored. Regression test
// for a file-handle leak analogous to aws/aws-sdk-go-v2#3512.
func TestDownloadDirectoryClosesCreatedFiles(t *testing.T) {
	cases := map[string]struct {
		destination        string
		objectsLists       [][]s3types.Object
		continuationTokens []string
		getobjectFn        func(*s3testing.TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error)
		failurePolicy      DownloadDirectoryFailurePolicy
		expectCreated      int
		expectRemoved      int
		expectClosed       int
		expectErr          string
	}{
		"single object": {
			destination: "close-single-object",
			objectsLists: [][]s3types.Object{
				{{Key: aws.String("foo/bar")}},
			},
			expectCreated: 1,
			expectClosed:  1,
		},
		"multiple objects with subdirs": {
			destination: "close-multiple-objects",
			objectsLists: [][]s3types.Object{
				{
					{Key: aws.String("foo/bar")},
					{Key: aws.String("baz")},
					{Key: aws.String("foo/zoo/bar")},
					{Key: aws.String("foo/zoo/oii/bababoii")},
				},
			},
			expectCreated: 4,
			expectClosed:  4,
		},
		"multiple objects paginated": {
			destination: "close-multiple-objects-paginated",
			objectsLists: [][]s3types.Object{
				{{Key: aws.String("foo/bar")}, {Key: aws.String("baz")}},
				{{Key: aws.String("foo/zoo/bar")}, {Key: aws.String("foo/zoo/oii/bababoii")}},
				{{Key: aws.String("foo/zoo/baz")}, {Key: aws.String("foo/zoo/oii/yee")}},
			},
			continuationTokens: []string{"token1", "token2"},
			expectCreated:      6,
			expectClosed:       6,
		},
		"created files are closed when some downloads fail and are ignored": {
			destination: "close-error-ignored",
			objectsLists: [][]s3types.Object{
				{
					{Key: aws.String("foo/bar")},
					{Key: aws.String("baz")},
					{Key: aws.String("foo/zoo/bar")},
					{Key: aws.String("foo/zoo/oii/bababoii")},
				},
			},
			getobjectFn: func(c *s3testing.TransferManagerLoggingClient, in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				if key := aws.ToString(in.Key); key == "foo/zoo/bar" || key == "baz" {
					return nil, fmt.Errorf("mocking error")
				}
				return &s3.GetObjectOutput{
					Body:          io.NopCloser(bytes.NewReader(c.Data)),
					ContentLength: aws.Int64(int64(len(c.Data))),
					PartsCount:    aws.Int32(c.PartsCount),
					ETag:          aws.String(etag),
				}, nil
			},
			failurePolicy: IgnoreDownloadFailurePolicy{},
			// The destination file is created before the transfer starts, so
			// os.Create runs for all four objects. foo/zoo/bar and baz fail
			// their GET and are closed then removed; removedOrClosed counts
			// those under expectRemoved rather than expectClosed.
			expectCreated: 4,
			expectClosed:  2,
			expectRemoved: 2,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			capture := &createdFileCapture{
				files: make(map[string]*os.File),
			}
			prev := createFileFn
			createFileFn = func(path string) (*os.File, error) {
				f, err := prev(path)
				if err == nil {
					capture.capture(path, f)
				}
				return f, err
			}
			defer func() { createFileFn = prev }()

			s3Client, _ := s3testing.NewDownloadDirectoryClient()
			s3Client.ListObjectsData = c.objectsLists
			s3Client.ContinuationTokens = c.continuationTokens
			if c.getobjectFn == nil {
				s3Client.GetObjectFn = s3testing.PartGetObjectFn
			} else {
				s3Client.GetObjectFn = c.getobjectFn
			}
			s3Client.Data = make([]byte, 0)
			s3Client.PartsCount = 1
			mgr := New(s3Client)

			dstPath := filepath.Join("testdata", c.destination)
			defer os.RemoveAll(dstPath)

			req := &DownloadDirectoryInput{
				Bucket:        aws.String("mock-bucket"),
				Destination:   aws.String(dstPath),
				FailurePolicy: c.failurePolicy,
			}

			_, err := mgr.DownloadDirectory(context.Background(), req)
			if err != nil {
				if c.expectErr == "" {
					t.Fatalf("expect no error, got %v", err)
				} else if !strings.Contains(err.Error(), c.expectErr) {
					t.Fatalf("expect %s to be contained in %v", c.expectErr, err)
				}
			} else if c.expectErr != "" {
				t.Fatalf("expect error %s, got none", c.expectErr)
			}

			if e, a := c.expectCreated, capture.count(); e != a {
				t.Errorf("expected DownloadDirectory to create %d file(s) under %s, captured %d", e, dstPath, a)
			}
			removed, closed := capture.removedOrClosed()
			if e := c.expectRemoved; e != len(removed) {
				t.Errorf("expected %d files to be removed, but removed %v", e, removed)
			}
			if e := c.expectClosed; e != len(closed) {
				t.Errorf("expected %d created files under %s to be closed after DownloadDirectory returned, but only close those: %v", e, dstPath, closed)
			}
		})
	}
}

// TestDownloadDirectoryNoHeadObject asserts the directory path issues no
// HeadObject: DownloadObject learns each object's size from its first data GET.
func TestDownloadDirectoryNoHeadObject(t *testing.T) {
	s3Client, _ := s3testing.NewDownloadDirectoryClient()
	s3Client.ListObjectsData = [][]s3types.Object{{
		{Key: aws.String("foo/bar")},
		{Key: aws.String("foo/baz")},
	}}
	s3Client.GetObjectFn = s3testing.PartGetObjectFn
	s3Client.Data = []byte("hello world")
	s3Client.PartsCount = 1

	dstPath := filepath.Join("testdata", "no-head-object")
	defer os.RemoveAll(dstPath)

	out, err := New(s3Client).DownloadDirectory(context.Background(), &DownloadDirectoryInput{
		Bucket:      aws.String("mock-bucket"),
		Destination: aws.String(dstPath),
	})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := int64(2), out.ObjectsDownloaded; e != a {
		t.Fatalf("expect %d objects downloaded, got %d", e, a)
	}

	if n := len(s3Client.HeadObjectInputs); n != 0 {
		t.Errorf("expect no HeadObject calls on the directory path, got %d", n)
	}
	// An empty HeadObjectInputs is also consistent with no download happening.
	if s3Client.GetObjectInvocations == 0 {
		t.Error("expect at least one GetObject call, got none")
	}
}

// TestDownloadDirectoryWriteOffsets is a regression test for
// aws/aws-sdk-go-v2#3536. The parts are unequal so the size latched from part 1 is
// wrong for parts 2 and 3, and each part carries a distinct byte so a misplaced
// write shows up as a content mismatch rather than only a length mismatch.
func TestDownloadDirectoryWriteOffsets(t *testing.T) {
	partA := bytes.Repeat([]byte{'A'}, 1000)
	partB := bytes.Repeat([]byte{'B'}, 700)
	partC := bytes.Repeat([]byte{'C'}, 300)
	want := bytes.Join([][]byte{partA, partB, partC}, nil)

	s3Client, _ := s3testing.NewDownloadDirectoryClient()
	s3Client.ListObjectsData = [][]s3types.Object{{
		{Key: aws.String("multipart/object")},
	}}
	s3Client.GetObjectFn = s3testing.UnequalPartGetObjectFn
	s3Client.PartsData = [][]byte{partA, partB, partC}
	s3Client.PartsCount = 3

	dstPath := filepath.Join("testdata", "write-offsets")
	defer os.RemoveAll(dstPath)

	out, err := New(s3Client).DownloadDirectory(context.Background(), &DownloadDirectoryInput{
		Bucket:      aws.String("mock-bucket"),
		Destination: aws.String(dstPath),
	})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if e, a := int64(1), out.ObjectsDownloaded; e != a {
		t.Fatalf("expect %d objects downloaded, got %d", e, a)
	}

	got, err := os.ReadFile(filepath.Join(dstPath, "multipart", "object"))
	if err != nil {
		t.Fatalf("expect to read downloaded file, got %v", err)
	}
	if e, a := len(want), len(got); e != a {
		t.Fatalf("expect downloaded file to be %d bytes, got %d", e, a)
	}
	if !bytes.Equal(want, got) {
		for i := range want {
			if want[i] != got[i] {
				t.Fatalf("downloaded file differs at offset %d: expect %q, got %q", i, want[i], got[i])
			}
		}
	}
}

// TestMapDownloadObjectInputIsTotal guards against field drift: a field added to
// both input types but forgotten in mapDownloadObjectInput would silently drop
// customer input. Every source field is populated and required to arrive.
func TestMapDownloadObjectInputIsTotal(t *testing.T) {
	in := &GetObjectInput{}
	src := reflect.ValueOf(in).Elem()
	for i := 0; i < src.NumField(); i++ {
		f := src.Type().Field(i)
		if !f.IsExported() {
			continue
		}
		v, err := distinctValue(f.Type, i)
		if err != nil {
			// Deliberately fatal rather than skipped: a type this helper cannot
			// populate is a field the guard would silently stop covering.
			t.Fatalf("cannot populate GetObjectInput.%s (%s): %v -- extend distinctValue", f.Name, f.Type, err)
		}
		src.Field(i).Set(v)
	}

	var w nopWriterAt
	got := reflect.ValueOf(mapDownloadObjectInput(in, w)).Elem()

	for i := 0; i < src.NumField(); i++ {
		f := src.Type().Field(i)
		if !f.IsExported() {
			continue
		}
		dst := got.FieldByName(f.Name)
		if !dst.IsValid() {
			t.Errorf("DownloadObjectInput is missing field %s, which GetObjectInput has", f.Name)
			continue
		}
		if e, a := f.Type, dst.Type(); e != a {
			t.Errorf("field %s: GetObjectInput has type %s, DownloadObjectInput has %s", f.Name, e, a)
			continue
		}
		if !reflect.DeepEqual(src.Field(i).Interface(), dst.Interface()) {
			t.Errorf("field %s is not carried across by mapDownloadObjectInput: expect %s, got %s",
				f.Name, deref(src.Field(i)), deref(dst))
		}
	}

	// WriterAt is the only field DownloadObjectInput may add. Anything else new is
	// a field DownloadDirectory would be leaving unset by accident.
	for i := 0; i < got.NumField(); i++ {
		f := got.Type().Field(i)
		if !f.IsExported() || f.Name == "WriterAt" {
			continue
		}
		if _, ok := src.Type().FieldByName(f.Name); !ok {
			t.Errorf("DownloadObjectInput has extra field %s: either map it or document why it stays unset", f.Name)
		}
	}
	if got.FieldByName("WriterAt").IsNil() {
		t.Error("expect WriterAt to be set by mapDownloadObjectInput")
	}
}

// deref renders a pointer field as its value so failures print values, not addresses.
func deref(v reflect.Value) string {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "<nil>"
		}
		return fmt.Sprintf("%v", v.Elem().Interface())
	}
	return fmt.Sprintf("%v", v.Interface())
}

type nopWriterAt struct{}

func (nopWriterAt) WriteAt(p []byte, off int64) (int, error) { return len(p), nil }

// distinctValue builds a non-zero value of t that differs per field index, so a
// mapper that assigns the right type from the wrong source field is still caught.
func distinctValue(t reflect.Type, seed int) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.String:
		return reflect.ValueOf(fmt.Sprintf("value-%d", seed)).Convert(t), nil
	case reflect.Ptr:
		switch t.Elem() {
		case reflect.TypeOf(time.Time{}):
			ts := time.Unix(int64(1700000000+seed), 0).UTC()
			p := reflect.New(t.Elem())
			p.Elem().Set(reflect.ValueOf(ts))
			return p, nil
		}
		switch t.Elem().Kind() {
		case reflect.String:
			p := reflect.New(t.Elem())
			p.Elem().Set(reflect.ValueOf(fmt.Sprintf("value-%d", seed)).Convert(t.Elem()))
			return p, nil
		case reflect.Bool:
			p := reflect.New(t.Elem())
			p.Elem().SetBool(true)
			return p, nil
		case reflect.Int32, reflect.Int64:
			p := reflect.New(t.Elem())
			p.Elem().SetInt(int64(seed + 1))
			return p, nil
		}
	}
	return reflect.Value{}, fmt.Errorf("unsupported kind %s", t.Kind())
}
