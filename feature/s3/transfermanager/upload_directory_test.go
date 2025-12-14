package transfermanager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/internal/awstesting"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type filenameFilter struct {
	keyword string
}

func (ff *filenameFilter) FilterFile(path string) bool {
	if strings.Contains(filepath.Base(path), ff.keyword) {
		return false
	}
	return true
}

type keynameCallback struct {
	keyword string
}

func (kc *keynameCallback) UpdateRequest(in *UploadObjectInput) {
	if k := aws.ToString(in.Key); k == kc.keyword {
		*in.Key = k + "/gotyou"
	}
}

func TestUploadDirectory(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "testdata")

	cases := map[string]struct {
		source               string
		followSymLinks       bool
		recursive            bool
		keyPrefix            string
		filter               FileFilter
		s3Delimiter          string
		callback             PutRequestCallback
		failurePolicy        UploadDirectoryFailurePolicy
		putobjectFunc        func(*s3testing.TransferManagerLoggingClient, *s3.PutObjectInput) (*s3.PutObjectOutput, error)
		preprocessFunc       func(string) (func() error, error)
		expectKeys           []string
		expectErr            string
		expectFilesUploaded  int
		expectFilesFailed    int
		listenerValidationFn func(*testing.T, *mockDirectoryListener, any, any, error)
	}{
		"single file recursively": {
			source:              filepath.Join(root, "single-file-dir"),
			recursive:           true,
			expectKeys:          []string{"foo"},
			expectFilesUploaded: 1,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 1)
			},
		},
		"multi file at root recursively": {
			source:              filepath.Join(root, "multi-file-at-root"),
			recursive:           true,
			expectKeys:          []string{"foo", "bar", "baz"},
			expectFilesUploaded: 3,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 3)
			},
		},
		"multi file with subdir recursively": {
			source:              filepath.Join(root, "multi-file-with-subdir"),
			recursive:           true,
			expectKeys:          []string{"foo", "bar", "zoo/baz", "zoo/oii/yee"},
			expectFilesUploaded: 4,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 4)
			},
		},
		"multi file with subdir non-recursively": {
			source:              filepath.Join(root, "multi-file-with-subdir"),
			expectKeys:          []string{"foo", "bar"},
			expectFilesUploaded: 2,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 2)
			},
		},
		"multi file with subdir and filter recursively": {
			source:              filepath.Join(root, "multi-file-with-subdir"),
			recursive:           true,
			filter:              &filenameFilter{"ar"},
			expectKeys:          []string{"foo", "zoo/baz", "zoo/oii/yee"},
			expectFilesUploaded: 3,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 3)
			},
		},
		"folder with single file and symlink recursively": {
			source:         filepath.Join(root, "single-file-dir"),
			followSymLinks: true,
			recursive:      true,
			preprocessFunc: func(root string) (func() error, error) {
				symlinkPath := filepath.Join(root, "single-file-dir", "symFoo")
				postprocessFunc := func() error {
					os.Remove(symlinkPath)
					return nil
				}
				if err := os.Symlink(filepath.Join(root, "dstFile1"), symlinkPath); err != nil {
					return postprocessFunc, err
				}
				return postprocessFunc, nil
			},
			expectKeys:          []string{"symFoo", "foo"},
			expectFilesUploaded: 2,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 2)
			},
		},
		"folder containing both file and symlink": {
			source:         filepath.Join(root, "multi-file-contain-symlink"),
			followSymLinks: true,
			recursive:      true,
			preprocessFunc: func(root string) (func() error, error) {
				symlinkPath := filepath.Join(root, "multi-file-contain-symlink", "to", "the", "symFoo")
				postprocessFunc := func() error {
					os.Remove(symlinkPath)
					return nil
				}
				if err := os.Symlink(filepath.Join(root, "dstFile1"), symlinkPath); err != nil {
					return postprocessFunc, err
				}
				return postprocessFunc, nil
			},
			expectKeys:          []string{"foo", "bar", "to/baz", "to/the/symFoo", "to/the/yee"},
			expectFilesUploaded: 5,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 5)
			},
		},
		"folder containing multi symlinks": {
			source:         filepath.Join(root, "multi-file-contain-symlink"),
			followSymLinks: true,
			recursive:      true,
			preprocessFunc: func(root string) (func() error, error) {
				symlinkPath1 := filepath.Join(root, "multi-file-contain-symlink", "to", "the", "symFoo")
				symlinkPath2 := filepath.Join(root, "multi-file-contain-symlink", "to", "the", "symBar")
				postprocessFunc := func() error {
					// this cleans up all possible symlinks regardless of
					// whether or not it is successfully created
					os.Remove(symlinkPath1)
					os.Remove(symlinkPath2)
					return nil
				}
				if err := os.Symlink(filepath.Join(root, "dstFile1"), symlinkPath1); err != nil {
					return postprocessFunc, err
				}
				if err := os.Symlink(filepath.Join(root, "dstFile2"), symlinkPath2); err != nil {
					return postprocessFunc, err
				}

				return postprocessFunc, nil
			},
			expectKeys:          []string{"foo", "bar", "to/baz", "to/the/symFoo", "to/the/symBar", "to/the/yee"},
			expectFilesUploaded: 6,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 6)
			},
		},
		"folder containing multi symlinks but not follow": {
			source:    filepath.Join(root, "multi-file-contain-symlink"),
			recursive: true,
			preprocessFunc: func(root string) (func() error, error) {
				symlinkPath1 := filepath.Join(root, "multi-file-contain-symlink", "to", "the", "symFoo")
				symlinkPath2 := filepath.Join(root, "multi-file-contain-symlink", "to", "the", "symBar")
				postprocessFunc := func() error {
					// this cleans up all possible symlinks regardless of
					// whether or not it is successfully created
					os.Remove(symlinkPath1)
					os.Remove(symlinkPath2)
					return nil
				}
				if err := os.Symlink(filepath.Join(root, "dstFile1"), symlinkPath1); err != nil {
					return postprocessFunc, err
				}
				if err := os.Symlink(filepath.Join(root, "dstFile2"), symlinkPath2); err != nil {
					return postprocessFunc, err
				}

				return postprocessFunc, nil
			},
			expectKeys:          []string{"foo", "bar", "to/baz", "to/the/yee"},
			expectFilesUploaded: 4,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 4)
			},
		},
		"folder containing files and symlink referring to folder": {
			source:         filepath.Join(root, "multi-file-contain-symlink"),
			followSymLinks: true,
			recursive:      true,
			preprocessFunc: func(root string) (func() error, error) {
				symlinkPath := filepath.Join(root, "multi-file-contain-symlink", "to", "the", "symFoo")
				postprocessFunc := func() error {
					os.Remove(symlinkPath)
					return nil
				}
				if err := os.Symlink(filepath.Join(root, "dstDir1"), symlinkPath); err != nil {
					return postprocessFunc, err
				}
				return postprocessFunc, nil
			},
			expectKeys:          []string{"foo", "bar", "to/baz", "to/the/symFoo/foo", "to/the/yee"},
			expectFilesUploaded: 5,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 5)
			},
		},
		"folder containing files and empty folder": {
			source:         filepath.Join(root, "multi-file-contain-symlink"),
			followSymLinks: true,
			recursive:      true,
			preprocessFunc: func(root string) (func() error, error) {
				path := filepath.Join(root, "multi-file-contain-symlink", "to", "too")
				postprocessFunc := func() error {
					os.Remove(path)
					return nil
				}
				if err := os.MkdirAll(path, os.ModePerm); err != nil {
					return postprocessFunc, err
				}
				return postprocessFunc, nil
			},
			expectKeys:          []string{"foo", "bar", "to/baz", "to/the/yee"},
			expectFilesUploaded: 4,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 4)
			},
		},
		"error when a file upload fails": {
			source:    filepath.Join(root, "multi-file-with-subdir"),
			recursive: true,
			putobjectFunc: func(svc *s3testing.TransferManagerLoggingClient, param *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
				if aws.ToString(param.Key) == "zoo/oii/yee" {
					return nil, fmt.Errorf("banned key")
				}
				return &s3.PutObjectOutput{}, nil
			},
			expectErr: "banned key",
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectFailed(t, in, err)
			},
		},
		"specified files uploads failure ignored by failure policy for recursive upload": {
			source:    filepath.Join(root, "multi-file-with-subdir"),
			recursive: true,
			putobjectFunc: func(svc *s3testing.TransferManagerLoggingClient, param *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
				if key := aws.ToString(param.Key); key == "zoo/oii/yee" || key == "foo" {
					return nil, fmt.Errorf("banned key")
				}
				return &s3.PutObjectOutput{}, nil
			},
			failurePolicy:       IgnoreUploadFailurePolicy{},
			expectKeys:          []string{"zoo/baz", "bar", "zoo/oii/yee", "foo"},
			expectFilesUploaded: 2,
			expectFilesFailed:   2,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 2)
			},
		},
		"specified files uploads failure ignored by failure policy for non-recursive upload": {
			source: filepath.Join(root, "multi-file-with-subdir"),
			putobjectFunc: func(svc *s3testing.TransferManagerLoggingClient, param *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
				if key := aws.ToString(param.Key); key == "zoo/oii/yee" || key == "foo" {
					return nil, fmt.Errorf("banned key")
				}
				return &s3.PutObjectOutput{}, nil
			},
			failurePolicy:       IgnoreUploadFailurePolicy{},
			expectKeys:          []string{"bar", "foo"},
			expectFilesUploaded: 1,
			expectFilesFailed:   1,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 1)
			},
		},
		"error when a symlink refers to its upper dir": {
			source:         filepath.Join(root, "multi-file-contain-symlink"),
			followSymLinks: true,
			recursive:      true,
			preprocessFunc: func(root string) (func() error, error) {
				symlinkPath := filepath.Join(root, "multi-file-contain-symlink", "to", "the", "symLoop")
				postprocessFunc := func() error {
					os.Remove(symlinkPath)
					return nil
				}
				if err := os.Symlink(filepath.Join(root, "multi-file-contain-symlink", "to"), symlinkPath); err != nil {
					return postprocessFunc, err
				}
				return postprocessFunc, nil
			},
			expectErr: "traversed duplicate path",
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectFailed(t, in, err)
			},
		},
		"error when a symlink refers to another file under source": {
			source:         filepath.Join(root, "multi-file-contain-symlink"),
			followSymLinks: true,
			recursive:      true,
			preprocessFunc: func(root string) (func() error, error) {
				symlinkPath := filepath.Join(root, "multi-file-contain-symlink", "to", "the", "symLoop")
				postprocessFunc := func() error {
					os.Remove(symlinkPath)
					return nil
				}
				if err := os.Symlink(filepath.Join(root, "multi-file-contain-symlink", "foo"), symlinkPath); err != nil {
					return postprocessFunc, err
				}
				return postprocessFunc, nil
			},
			expectErr: "traversed duplicate path",
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectFailed(t, in, err)
			},
		},
		"error when source is not directory": {
			source:    filepath.Join(root, "non-dir-source"),
			expectErr: "doesn't point to a valid directory",
		},
		"single file recursively with keyprefix": {
			source:              filepath.Join(root, "single-file-dir"),
			recursive:           true,
			keyPrefix:           "bla",
			expectKeys:          []string{"bla/foo"},
			expectFilesUploaded: 1,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 1)
			},
		},
		"multi file with subdir and filter recursively with keyprefix": {
			source:              filepath.Join(root, "multi-file-with-subdir"),
			recursive:           true,
			keyPrefix:           "bla",
			filter:              &filenameFilter{"ar"},
			expectKeys:          []string{"bla/foo", "bla/zoo/baz", "bla/zoo/oii/yee"},
			expectFilesUploaded: 3,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 3)
			},
		},
		"folder containing both file and symlink with keyprefix": {
			source:         filepath.Join(root, "multi-file-contain-symlink"),
			followSymLinks: true,
			recursive:      true,
			keyPrefix:      "bla",
			preprocessFunc: func(root string) (func() error, error) {
				symlinkPath1 := filepath.Join(root, "multi-file-contain-symlink", "to", "the", "symFoo")
				symlinkPath2 := filepath.Join(root, "multi-file-contain-symlink", "to", "symBar")
				postprocessFunc := func() error {
					os.Remove(symlinkPath1)
					os.Remove(symlinkPath2)
					return nil
				}
				if err := os.Symlink(filepath.Join(root, "dstFile1"), symlinkPath1); err != nil {
					return postprocessFunc, err
				}
				if err := os.Symlink(filepath.Join(root, "dstDir1"), symlinkPath2); err != nil {
					return postprocessFunc, err
				}
				return postprocessFunc, nil
			},
			expectKeys:          []string{"bla/foo", "bla/bar", "bla/to/baz", "bla/to/the/symFoo", "bla/to/symBar/foo", "bla/to/the/yee"},
			expectFilesUploaded: 6,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 6)
			},
		},
		"folder containing symlink folder with prefix but non-recursive": {
			source:         filepath.Join(root, "multi-file-contain-symlink"),
			followSymLinks: true,
			keyPrefix:      "bla",
			preprocessFunc: func(root string) (func() error, error) {
				symlinkPath := filepath.Join(root, "multi-file-contain-symlink", "symDir")
				postprocessFunc := func() error {
					os.Remove(symlinkPath)
					return nil
				}
				if err := os.Symlink(filepath.Join(root, "dstDir1"), symlinkPath); err != nil {
					return postprocessFunc, err
				}
				return postprocessFunc, nil
			},
			expectKeys:          []string{"bla/foo", "bla/bar"},
			expectFilesUploaded: 2,
			listenerValidationFn: func(t *testing.T, l *mockDirectoryListener, in, out any, err error) {
				l.expectStart(t, in)
				l.expectComplete(t, in, out, 2)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			s3Client, params := s3testing.NewUploadDirectoryClient([]string{"UploadPart", "CompleteMultipartUpload"})
			s3Client.PutObjectFn = c.putobjectFunc
			mgr := New(s3Client, Options{})

			if c.preprocessFunc != nil {
				postprocessFunc, err := c.preprocessFunc(root)
				defer postprocessFunc()
				if err != nil {
					t.Fatalf("error when preprocessing: %v", err)
				}
			}

			req := &UploadDirectoryInput{
				Bucket:              aws.String("mock-bucket"),
				Source:              aws.String(c.source),
				FollowSymbolicLinks: aws.Bool(c.followSymLinks),
				Recursive:           aws.Bool(c.recursive),
				KeyPrefix:           aws.String(c.keyPrefix),
				Filter:              c.filter,
				Callback:            c.callback,
				FailurePolicy:       c.failurePolicy,
			}

			listener := &mockDirectoryListener{}

			resp, err := mgr.UploadDirectory(context.Background(), req, func(o *Options) {
				o.DirectoryProgressListeners.Register(listener)
			})
			if err != nil {
				if c.expectErr == "" {
					t.Fatalf("expect no err, got %v", err)
				} else if e, a := c.expectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %s error message to be in %s", e, a)
				}
			} else if c.expectErr != "" {
				t.Fatalf("expect error %s, got none", c.expectErr)

			}

			if c.listenerValidationFn != nil {
				c.listenerValidationFn(t, listener, req, resp, err)
			}

			if err != nil {
				return
			}

			if e, a := c.expectFilesUploaded, resp.ObjectsUploaded; e != a {
				t.Errorf("expect %d objects uploaded, got %d", e, a)
			}
			if e, a := c.expectFilesFailed, resp.ObjectsFailed; e != a {
				t.Errorf("expect %d objects failed, got %d", e, a)
			}

			var actualKeys []string
			for _, param := range *params {
				if input, ok := param.(*s3.PutObjectInput); ok {
					actualKeys = append(actualKeys, aws.ToString(input.Key))
				} else if input, ok := param.(*s3.CreateMultipartUploadInput); ok {
					actualKeys = append(actualKeys, aws.ToString(input.Key))
				} else {
					t.Fatalf("error when casting captured inputs")
				}
			}

			sort.Strings(actualKeys)
			sort.Strings(c.expectKeys)
			if e, a := c.expectKeys, actualKeys; !reflect.DeepEqual(e, a) {
				t.Errorf("expect upload keys to be %v, got %v", e, a)
			}
		})
	}
}

func TestUploadDirectoryObjectsTransferred(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "testdata")

	cases := map[string]struct {
		source       string
		recursive    bool
		objectsCount []int64
	}{
		"single file": {
			source:       filepath.Join(root, "single-file-dir"),
			recursive:    true,
			objectsCount: []int64{1},
		},
		"multi file with subdir non-recursively": {
			source:       filepath.Join(root, "multi-file-with-subdir"),
			objectsCount: []int64{1, 2},
		},
		"multi file with subdir recursively": {
			source:       filepath.Join(root, "multi-file-with-subdir"),
			recursive:    true,
			objectsCount: []int64{1, 2, 3, 4},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			s3Client, _ := s3testing.NewUploadDirectoryClient([]string{"UploadPart", "CompleteMultipartUpload"})
			mgr := New(s3Client, Options{})

			req := &UploadDirectoryInput{
				Bucket:    aws.String("mock-bucket"),
				Source:    aws.String(c.source),
				Recursive: aws.Bool(c.recursive),
			}

			listener := &mockDirectoryListener{}

			_, err := mgr.UploadDirectory(context.Background(), req, func(o *Options) {
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

func TestUploadDirectoryWithContextCanceled(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "testdata")
	c := s3.New(s3.Options{
		UsePathStyle: true,
		Region:       "mock-region",
	})
	u := New(c, Options{})

	ctx := &awstesting.FakeContext{DoneCh: make(chan struct{})}
	ctx.Error = fmt.Errorf("context canceled")
	close(ctx.DoneCh)

	_, err := u.UploadDirectory(ctx, &UploadDirectoryInput{
		Bucket:    aws.String("mock-bucket"),
		Source:    aws.String(filepath.Join(root, "multi-file-contain-symlink")),
		Recursive: aws.Bool(true),
	})
	if err == nil {
		t.Fatalf("expect error, got nil")
	}

	if e, a := "canceled", err.Error(); !strings.Contains(a, e) {
		t.Errorf("expected error message to contain %q, but did not %q", e, a)
	}
}
