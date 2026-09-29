package transfermanager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	internalio "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/io"
	s3testing "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/internal/testing"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const testFilePath = "/dir/object"

type fakeInit struct {
	size      int64
	partSize  int64
	writeSize int64
	directIO  bool
}

// fakeDownloadFile is an internalio.File that assembles writes in memory.
type fakeDownloadFile struct {
	initErr  error
	writeErr error
	closeErr error

	h     *downloadFileHarness
	mu    sync.Mutex
	inits []fakeInit
	data  []byte
}

func (f *fakeDownloadFile) Init(size, partSize, writeSize int64, directIO bool) error {
	f.h.record("Init")
	f.mu.Lock()
	defer f.mu.Unlock()

	f.inits = append(f.inits, fakeInit{size: size, partSize: partSize, writeSize: writeSize, directIO: directIO})
	return f.initErr
}

func (f *fakeDownloadFile) WriteAt(p []byte, off int64) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.writeErr != nil {
		return 0, f.writeErr
	}
	if end := int(off) + len(p); end > len(f.data) {
		f.data = append(f.data, make([]byte, end-len(f.data))...)
	}
	copy(f.data[off:], p)
	return len(p), nil
}

func (f *fakeDownloadFile) Close() error {
	f.h.record("Close")
	return f.closeErr
}

// downloadFileHarness records every filesystem operation DownloadFile makes.
// Events excludes WriteAt, whose count depends on part size.
type downloadFileHarness struct {
	mu      sync.Mutex
	events  []string
	created []string
	renames [][2]string
	removes []string
}

func (h *downloadFileHarness) record(event string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.events = append(h.events, event)
}

func stubDownloadFile(t *testing.T, file *fakeDownloadFile, createErr, renameErr error) *downloadFileHarness {
	t.Helper()

	h := &downloadFileHarness{}
	if file != nil {
		file.h = h
	}

	prevCreate, prevRename, prevRemove := createDownloadFileFn, renameFileFn, removeFileFn
	t.Cleanup(func() {
		createDownloadFileFn, renameFileFn, removeFileFn = prevCreate, prevRename, prevRemove
	})

	createDownloadFileFn = func(path string) (internalio.File, error) {
		h.mu.Lock()
		h.created = append(h.created, path)
		h.mu.Unlock()
		if createErr != nil {
			return nil, createErr
		}
		return file, nil
	}
	renameFileFn = func(from, to string) error {
		h.mu.Lock()
		h.renames = append(h.renames, [2]string{from, to})
		h.mu.Unlock()
		h.record("Rename")
		return renameErr
	}
	removeFileFn = func(path string) error {
		h.mu.Lock()
		h.removes = append(h.removes, path)
		h.mu.Unlock()
		h.record("Remove")
		return nil
	}
	return h
}

// expect asserts DownloadFile created exactly one temp file next to
// testFilePath, performed events in order, and renamed/removed that temp file.
func (h *downloadFileHarness) expect(t *testing.T, events ...string) {
	t.Helper()

	if len(h.created) != 1 {
		t.Fatalf("created = %v, want one temp file", h.created)
	}
	tmp := h.created[0]
	if filepath.Dir(tmp) != filepath.Dir(testFilePath) || !strings.HasPrefix(tmp, testFilePath+".") || !strings.HasSuffix(tmp, ".tmp") {
		t.Fatalf("temp file %q is not <FilePath>.<id>.tmp", tmp)
	}

	if fmt.Sprint(h.events) != fmt.Sprint(events) {
		t.Fatalf("events = %v, want %v", h.events, events)
	}

	var wantRenames [][2]string
	var wantRemoves []string
	for _, e := range events {
		switch e {
		case "Rename":
			wantRenames = append(wantRenames, [2]string{tmp, testFilePath})
		case "Remove":
			wantRemoves = append(wantRemoves, tmp)
		}
	}
	if fmt.Sprint(h.renames) != fmt.Sprint(wantRenames) {
		t.Fatalf("renames = %v, want %v", h.renames, wantRenames)
	}
	if fmt.Sprint(h.removes) != fmt.Sprint(wantRemoves) {
		t.Fatalf("removes = %v, want %v", h.removes, wantRemoves)
	}
}

type rangeNotSatisfiableError struct{}

func (rangeNotSatisfiableError) Error() string       { return "InvalidRange" }
func (rangeNotSatisfiableError) HTTPStatusCode() int { return http.StatusRequestedRangeNotSatisfiable }

func failingGetObjectFn(*s3testing.TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return nil, errors.New("first request failed")
}

func TestDownloadFile(t *testing.T) {
	data := randomBytes(1, 20*megabyte+13)
	partSizes := []int{8 * megabyte, 5 * megabyte, 7*megabyte + 13}
	ranges := func(o *Options) { o.GetObjectType = types.GetObjectRanges }
	parts := func(o *Options) { o.GetObjectType = types.GetObjectParts }

	cases := map[string]struct {
		data        []byte
		getObjectFn func(*s3testing.TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error)
		optFn       func(*Options)
		rng         string
		initErr     error
		writeErr    error
		closeErr    error
		createErr   error
		renameErr   error

		expectErr    string
		expectEvents []string
		expectData   []byte
		expectInit   *fakeInit
	}{
		"ranges": {
			data:         data,
			getObjectFn:  s3testing.RangeGetObjectFn,
			optFn:        ranges,
			expectEvents: []string{"Init", "Close", "Rename"},
			expectData:   data,
			expectInit:   &fakeInit{size: int64(len(data)), directIO: true},
		},
		"parts": {
			data:         data,
			getObjectFn:  partsGetObjectFn(data, partSizes),
			optFn:        parts,
			expectEvents: []string{"Init", "Close", "Rename"},
			expectData:   data,
			expectInit:   &fakeInit{size: int64(len(data))},
		},
		"explicit range": {
			data:         data,
			getObjectFn:  s3testing.RangeGetObjectFn,
			optFn:        ranges,
			rng:          "bytes=2-16777218",
			expectEvents: []string{"Init", "Close", "Rename"},
			expectData:   data[2:16777219],
			expectInit:   &fakeInit{size: 16777217, directIO: true},
		},
		"empty object": {
			data:         []byte{},
			getObjectFn:  s3testing.RangeGetObjectFn,
			optFn:        ranges,
			expectEvents: []string{"Init", "Close", "Rename"},
			expectData:   []byte{},
			expectInit:   &fakeInit{size: 0, directIO: true},
		},
		"empty object, range not satisfiable": {
			getObjectFn: func(*s3testing.TransferManagerLoggingClient, *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return nil, rangeNotSatisfiableError{}
			},
			optFn:        ranges,
			expectEvents: []string{"Init", "Close", "Rename"},
			expectData:   []byte{},
			expectInit:   &fakeInit{size: 0},
		},
		"create fails": {
			data:        data,
			getObjectFn: s3testing.RangeGetObjectFn,
			optFn:       ranges,
			createErr:   errors.New("create failed"),
			expectErr:   "create failed",
		},
		"first request fails": {
			data:         data,
			getObjectFn:  failingGetObjectFn,
			optFn:        ranges,
			expectErr:    "first request failed",
			expectEvents: []string{"Close", "Remove"},
		},
		"later request fails": {
			data:         data,
			getObjectFn:  s3testing.ErrRangeGetObjectFn,
			optFn:        ranges,
			expectErr:    "s3 service error",
			expectEvents: []string{"Init", "Close", "Remove"},
		},
		"init fails": {
			data:         data,
			getObjectFn:  s3testing.RangeGetObjectFn,
			optFn:        ranges,
			initErr:      errors.New("init failed"),
			expectErr:    "init failed",
			expectEvents: []string{"Init", "Close", "Remove"},
		},
		"write fails": {
			data:         data,
			getObjectFn:  s3testing.RangeGetObjectFn,
			optFn:        ranges,
			writeErr:     errors.New("write failed"),
			expectErr:    "write failed",
			expectEvents: []string{"Init", "Close", "Remove"},
		},
		"close fails": {
			data:         data,
			getObjectFn:  s3testing.RangeGetObjectFn,
			optFn:        ranges,
			closeErr:     errors.New("close failed"),
			expectErr:    "close failed",
			expectEvents: []string{"Init", "Close", "Remove"},
		},
		"rename fails": {
			data:         data,
			getObjectFn:  s3testing.RangeGetObjectFn,
			optFn:        ranges,
			renameErr:    errors.New("rename failed"),
			expectErr:    "rename failed",
			expectEvents: []string{"Init", "Close", "Rename", "Remove"},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			file := &fakeDownloadFile{initErr: c.initErr, writeErr: c.writeErr, closeErr: c.closeErr}
			h := stubDownloadFile(t, file, c.createErr, c.renameErr)

			client := New(&s3testing.TransferManagerLoggingClient{
				Data:        c.data,
				GetObjectFn: c.getObjectFn,
			}, func(o *Options) {
				o.Concurrency = 1
				c.optFn(o)
			})

			in := &DownloadFileInput{
				Bucket:   aws.String("bucket"),
				Key:      aws.String("key"),
				FilePath: testFilePath,
			}
			if c.rng != "" {
				in.Range = aws.String(c.rng)
			}
			out, err := client.DownloadFile(context.Background(), in)

			if c.expectErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.expectErr) {
					t.Fatalf("expect error containing %q, got %v", c.expectErr, err)
				}
			} else if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			h.expect(t, c.expectEvents...)

			if c.expectInit != nil {
				want := *c.expectInit
				want.partSize = defaultPartSizeBytes
				want.writeSize = getWriteSize(defaultPartSizeBytes)
				if len(file.inits) != 1 || file.inits[0] != want {
					t.Fatalf("Init calls = %+v, want [%+v]", file.inits, want)
				}
			}
			if c.expectData != nil {
				checkBytes(t, file.data, c.expectData)
				if got := aws.ToInt64(out.ContentLength); got != int64(len(c.expectData)) {
					t.Errorf("ContentLength = %d, want %d", got, len(c.expectData))
				}
			}
		})
	}
}

func TestDownloadFileCanceled(t *testing.T) {
	file := &fakeDownloadFile{}
	h := stubDownloadFile(t, file, nil, nil)

	// The mock client ignores ctx, so simulate the real client failing requests
	// issued after cancellation.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := New(&s3testing.TransferManagerLoggingClient{
		Data: randomBytes(2, 20*megabyte),
		GetObjectFn: func(c *s3testing.TransferManagerLoggingClient, in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			out, err := s3testing.RangeGetObjectFn(c, in)
			cancel()
			return out, err
		},
	}, func(o *Options) {
		o.Concurrency = 1
		o.GetObjectType = types.GetObjectRanges
	})

	_, err := client.DownloadFile(ctx, &DownloadFileInput{
		Bucket:   aws.String("bucket"),
		Key:      aws.String("key"),
		FilePath: testFilePath,
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expect context.Canceled, got %v", err)
	}
	h.expect(t, "Init", "Close", "Remove")
}

func TestDownloadFileUniqueTempPath(t *testing.T) {
	h := stubDownloadFile(t, nil, errors.New("create failed"), nil)
	client := New(&s3testing.TransferManagerLoggingClient{})

	in := &DownloadFileInput{Bucket: aws.String("bucket"), Key: aws.String("key"), FilePath: testFilePath}
	for range 2 {
		if _, err := client.DownloadFile(context.Background(), in); err == nil {
			t.Fatal("expect create error")
		}
	}

	if len(h.created) != 2 || h.created[0] == h.created[1] {
		t.Fatalf("expect two distinct temp paths, got %v", h.created)
	}
}

func TestDownloadFileRequiresPath(t *testing.T) {
	h := stubDownloadFile(t, &fakeDownloadFile{}, nil, nil)
	client := New(&s3testing.TransferManagerLoggingClient{})

	for _, in := range []*DownloadFileInput{nil, {}} {
		if _, err := client.DownloadFile(context.Background(), in); err == nil {
			t.Errorf("expect error for input %v", in)
		}
	}
	if len(h.created) != 0 || len(h.events) != 0 {
		t.Fatalf("expect no file operations, got created=%v events=%v", h.created, h.events)
	}
}
