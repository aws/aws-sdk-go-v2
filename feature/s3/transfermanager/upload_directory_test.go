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

func (kc *keynameCallback) UpdateRequest(in *PutObjectInput) {
	if in.Key == kc.keyword {
		in.Key = in.Key + "/gotyou"
	}
}

func TestUploadDirectory(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "testdata")

	cases := map[string]struct {
		source              string
		followSymLinks      bool
		recursive           bool
		keyPrefix           string
		filter              FileFilter
		s3Delimiter         string
		callback            PutRequestCallback
		putobjectFunc       func(*s3testing.TransferManagerLoggingClient, *s3.PutObjectInput) (*s3.PutObjectOutput, error)
		preprocessFunc      func(string) (func() error, error)
		expectKeys          []string
		expectErr           string
		expectFilesUploaded int
	}{
		"single file recursively": {
			source:              filepath.Join(root, "single-file-dir"),
			recursive:           true,
			expectKeys:          []string{"foo"},
			expectFilesUploaded: 1,
		},
		"multi file at root recursively": {
			source:              filepath.Join(root, "multi-file-at-root"),
			recursive:           true,
			expectKeys:          []string{"foo", "bar", "baz"},
			expectFilesUploaded: 3,
		},
		"multi file with subdir recursively": {
			source:              filepath.Join(root, "multi-file-with-subdir"),
			recursive:           true,
			expectKeys:          []string{"foo", "bar", "zoo/baz", "zoo/oii/yee"},
			expectFilesUploaded: 4,
		},
		"multi file with subdir non-recursively": {
			source:              filepath.Join(root, "multi-file-with-subdir"),
			expectKeys:          []string{"foo", "bar"},
			expectFilesUploaded: 2,
		},
		"multi file with subdir and filter recursively": {
			source:              filepath.Join(root, "multi-file-with-subdir"),
			recursive:           true,
			filter:              &filenameFilter{"ar"},
			expectKeys:          []string{"foo", "zoo/baz", "zoo/oii/yee"},
			expectFilesUploaded: 3,
		},
		"single symlink recursively": {
			source:         filepath.Join(root, "single-symlink"),
			followSymLinks: true,
			recursive:      true,
			preprocessFunc: func(root string) (func() error, error) {
				symlinkPath := filepath.Join(root, "single-symlink", "symFoo")
				postprocessFunc := func() error {
					os.Remove(symlinkPath)
					return nil
				}
				if err := os.Symlink(filepath.Join(root, "dstFile1"), symlinkPath); err != nil {
					return postprocessFunc, err
				}
				return postprocessFunc, nil
			},
			expectKeys:          []string{"symFoo"},
			expectFilesUploaded: 1,
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
			expectKeys:          []string{"foo", "bar", "to/baz", "to/the/symFoo"},
			expectFilesUploaded: 4,
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
			expectKeys:          []string{"foo", "bar", "to/baz", "to/the/symFoo", "to/the/symBar"},
			expectFilesUploaded: 5,
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
			expectKeys:          []string{"foo", "bar", "to/baz"},
			expectFilesUploaded: 3,
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
			expectKeys:          []string{"foo", "bar", "to/baz", "to/the/symFoo/foo"},
			expectFilesUploaded: 4,
		},
		"folder containing files and empty folder": {
			source:              filepath.Join(root, "multi-file-contain-symlink"),
			followSymLinks:      true,
			recursive:           true,
			expectKeys:          []string{"foo", "bar", "to/baz"},
			expectFilesUploaded: 3,
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
		},
		"multi file with subdir and filter recursively with keyprefix": {
			source:              filepath.Join(root, "multi-file-with-subdir"),
			recursive:           true,
			keyPrefix:           "bla",
			filter:              &filenameFilter{"ar"},
			expectKeys:          []string{"bla/foo", "bla/zoo/baz", "bla/zoo/oii/yee"},
			expectFilesUploaded: 3,
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
			expectKeys:          []string{"bla/foo", "bla/bar", "bla/to/baz", "bla/to/the/symFoo", "bla/to/symBar/foo"},
			expectFilesUploaded: 5,
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
		},
		"folder containing both file and symlink with keyprefix and custome delimiter": {
			source:         filepath.Join(root, "multi-file-contain-symlink"),
			followSymLinks: true,
			recursive:      true,
			keyPrefix:      "bla",
			s3Delimiter:    "#",
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
			expectKeys:          []string{"bla#foo", "bla#bar", "bla#to#baz", "bla#to#the#symFoo", "bla#to#symBar#foo"},
			expectFilesUploaded: 5,
		},
		"folder containing both file and symlink with keyprefix, custome delimiter and request callback": {
			source:         filepath.Join(root, "multi-file-contain-symlink"),
			followSymLinks: true,
			recursive:      true,
			keyPrefix:      "bla",
			s3Delimiter:    "#",
			callback:       &keynameCallback{"bla#to#baz"},
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
			expectKeys:          []string{"bla#foo", "bla#bar", "bla#to#baz/gotyou", "bla#to#the#symFoo", "bla#to#symBar#foo"},
			expectFilesUploaded: 5,
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

			resp, err := mgr.UploadDirectory(context.Background(), &UploadDirectoryInput{
				Bucket:              "mock-bucket",
				Source:              c.source,
				FollowSymbolicLinks: c.followSymLinks,
				Recursive:           c.recursive,
				KeyPrefix:           c.keyPrefix,
				Filter:              c.filter,
				Callback:            c.callback,
				S3Delimiter:         c.s3Delimiter,
			})
			if err != nil {
				if c.expectErr == "" {
					t.Fatalf("expect no err, got %v", err)
				} else if e, a := c.expectErr, err.Error(); !strings.Contains(a, e) {
					t.Fatalf("expect %s error message to be in %s", e, a)
				}
			} else {
				if c.expectErr != "" {
					t.Fatalf("expect error %s, got none", c.expectErr)
				}
			}

			if err != nil {
				return
			}

			if e, a := c.expectFilesUploaded, resp.ObjectsUploaded; e != a {
				t.Errorf("expect %d objects uploaded, got %d", e, a)
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
