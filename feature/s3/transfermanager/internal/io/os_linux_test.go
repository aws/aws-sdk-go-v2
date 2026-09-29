//go:build linux

package io

import (
	"bytes"
	"errors"
	"os"
	"syscall"
	"testing"
)

const (
	bufferedOpenFlags = os.O_RDWR | os.O_CREATE | os.O_EXCL
	directOpenFlags   = os.O_WRONLY | os.O_CREATE | os.O_EXCL | syscall.O_DIRECT

	testPartSize  = 10 * 1024 * 1024
	testWriteSize = 8 * 1024 * 1024
)

type fallocateCall struct {
	fd   int
	mode uint32
	off  int64
	size int64
}

func stubStatfs(t *testing.T, bsize int64, err error) {
	t.Helper()

	original := statfs
	t.Cleanup(func() { statfs = original })
	statfs = func(_ string, stat *syscall.Statfs_t) error {
		setBsize(&stat.Bsize, bsize)
		return err
	}
}

func stubFallocate(t *testing.T, err error) *[]fallocateCall {
	t.Helper()

	original := fallocate
	t.Cleanup(func() { fallocate = original })

	var calls []fallocateCall
	fallocate = func(fd int, mode uint32, off int64, size int64) error {
		calls = append(calls, fallocateCall{fd: fd, mode: mode, off: off, size: size})
		return err
	}
	return &calls
}

func expectOpen(t *testing.T, opens []openCall, want openCall) {
	t.Helper()

	if len(opens) != 1 || opens[0] != want {
		t.Fatalf("openFile calls = %+v, want [%+v]", opens, want)
	}
}

func TestSupportsDirectIOFilesystemBlockSize(t *testing.T) {
	for _, test := range []struct {
		name    string
		bsize   int64
		allowed bool
	}{
		{name: "512", bsize: 512, allowed: true},
		{name: "1024", bsize: 1024, allowed: true},
		{name: "2048", bsize: 2048, allowed: true},
		{name: "4096", bsize: 4096, allowed: true},
		{name: "zero", bsize: 0},
		{name: "3072", bsize: 3072},
		{name: "8192", bsize: 8192},
	} {
		t.Run(test.name, func(t *testing.T) {
			stubStatfs(t, test.bsize, nil)

			got := supportsDirectIO("/tmp/file", testPartSize, testWriteSize)
			if got != test.allowed {
				t.Fatalf("supportsDirectIO = %v, want %v", got, test.allowed)
			}
		})
	}
}

func TestSupportsDirectIOAlignment(t *testing.T) {
	stubStatfs(t, 4096, nil)

	for _, test := range []struct {
		name      string
		partSize  int64
		writeSize int64
		allowed   bool
	}{
		{name: "part remainder is aligned", partSize: testPartSize, writeSize: testWriteSize, allowed: true},
		{name: "part remainder is not aligned", partSize: testPartSize + 1, writeSize: testWriteSize},
		{name: "write size is not aligned", partSize: testPartSize, writeSize: testWriteSize + 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := supportsDirectIO("/tmp/file", test.partSize, test.writeSize)
			if got != test.allowed {
				t.Fatalf("supportsDirectIO = %v, want %v", got, test.allowed)
			}
		})
	}
}

func TestFileInitUsesBufferedIO(t *testing.T) {
	for _, test := range []struct {
		name      string
		size      int64
		bsize     int64
		statErr   error
		partSize  int64
		directIO  bool
		writeSize int64
	}{
		{name: "size threshold", size: oDirectThreshold - 1, bsize: 4096, partSize: testPartSize, writeSize: testWriteSize, directIO: true},
		{name: "filesystem block size", size: oDirectThreshold, bsize: 8192, partSize: testPartSize, writeSize: testWriteSize, directIO: true},
		{name: "filesystem stat error", size: oDirectThreshold, bsize: 4096, statErr: errors.New("statfs failed"), partSize: testPartSize, writeSize: testWriteSize, directIO: true},
		{name: "part size", size: oDirectThreshold, bsize: 4096, partSize: testPartSize + 1, writeSize: testWriteSize, directIO: true},
		{name: "direct I/O disabled", size: oDirectThreshold, bsize: 4096, partSize: testPartSize, writeSize: testWriteSize},
	} {
		t.Run(test.name, func(t *testing.T) {
			stubStatfs(t, test.bsize, test.statErr)
			fallocates := stubFallocate(t, nil)
			ff := &fakeFile{}
			opens := stubOpenFile(t, ff, nil)

			created, err := Create("/dir/file")
			if err != nil {
				t.Fatal(err)
			}
			f := created.(*file)
			if err := f.Init(test.size, test.partSize, test.writeSize, test.directIO); err != nil {
				t.Fatal(err)
			}

			expectOpen(t, *opens, openCall{name: "/dir/file", flag: bufferedOpenFlags, perm: 0o666})
			if f.direct {
				t.Fatal("file initialized with direct I/O")
			}
			if len(*fallocates) != 0 {
				t.Fatalf("fallocate called for buffered file: %+v", *fallocates)
			}

			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
			expectCalls(t, ff, "Close")
		})
	}
}

func TestFileInitDirectIO(t *testing.T) {
	stubStatfs(t, 4096, nil)
	fallocates := stubFallocate(t, nil)
	ff := &fakeFile{fd: 42}
	opens := stubOpenFile(t, ff, nil)

	created, err := Create("/dir/file")
	if err != nil {
		t.Fatal(err)
	}
	f := created.(*file)
	const size = oDirectThreshold
	if err := f.Init(size, testPartSize, testWriteSize, true); err != nil {
		t.Fatal(err)
	}

	expectOpen(t, *opens, openCall{name: "/dir/file", flag: directOpenFlags, perm: 0o644})
	if !f.direct {
		t.Fatal("file did not initialize with direct I/O")
	}
	want := fallocateCall{fd: 42, mode: 0, off: 0, size: size}
	if len(*fallocates) != 1 || (*fallocates)[0] != want {
		t.Fatalf("fallocate calls = %+v, want [%+v]", *fallocates, want)
	}
	expectCalls(t, ff)
}

func TestFileInitAlreadyInitialized(t *testing.T) {
	opens := stubOpenFile(t, &fakeFile{}, nil)

	created, err := Create("/dir/file")
	if err != nil {
		t.Fatal(err)
	}
	f := created.(*file)
	if err := f.Init(1, 1, 1, false); err != nil {
		t.Fatal(err)
	}
	if err := f.Init(1, 1, 1, false); err == nil {
		t.Fatal("second Init returned nil error")
	}
	if len(*opens) != 1 {
		t.Fatalf("openFile called %d times, want 1", len(*opens))
	}
}

func TestFileInitOpenError(t *testing.T) {
	for _, test := range []struct {
		name  string
		size  int64
		flags int
	}{
		{name: "buffered", size: 1, flags: bufferedOpenFlags},
		{name: "direct", size: oDirectThreshold, flags: directOpenFlags},
	} {
		t.Run(test.name, func(t *testing.T) {
			stubStatfs(t, 4096, nil)
			fallocates := stubFallocate(t, nil)
			wantErr := errors.New("open failed")
			opens := stubOpenFile(t, nil, wantErr)

			created, err := Create("/dir/file")
			if err != nil {
				t.Fatal(err)
			}
			f := created.(*file)
			if err := f.Init(test.size, testPartSize, testWriteSize, true); !errors.Is(err, wantErr) {
				t.Fatalf("Init error = %v, want %v", err, wantErr)
			}

			if len(*opens) != 1 || (*opens)[0].flag != test.flags {
				t.Fatalf("openFile calls = %+v, want flags %#x", *opens, test.flags)
			}
			if f.File != nil {
				t.Fatal("file handle set after open failure")
			}
			if len(*fallocates) != 0 {
				t.Fatalf("fallocate called after open failure: %+v", *fallocates)
			}
		})
	}
}

func TestFileInitFallocateError(t *testing.T) {
	stubStatfs(t, 4096, nil)
	wantErr := errors.New("fallocate failed")
	stubFallocate(t, wantErr)
	ff := &fakeFile{}
	stubOpenFile(t, ff, nil)

	created, err := Create("/dir/file")
	if err != nil {
		t.Fatal(err)
	}
	f := created.(*file)
	if err := f.Init(oDirectThreshold, testPartSize, testWriteSize, true); !errors.Is(err, wantErr) {
		t.Fatalf("Init error = %v, want %v", err, wantErr)
	}

	if f.File != nil {
		t.Fatal("file handle set after fallocate failure")
	}
	if f.direct {
		t.Fatal("file marked direct after fallocate failure")
	}
	expectCalls(t, ff, "Close")
}

func TestFileWriteAt(t *testing.T) {
	const (
		unaligned = oDirectThreshold + 13
		aligned   = oDirectThreshold
	)

	for _, test := range []struct {
		name      string
		direct    bool
		size      int64
		len       int
		off       int64
		expectLen int
	}{
		{name: "buffered", size: unaligned, len: 13, off: unaligned - 13, expectLen: 13},
		{name: "direct not final", direct: true, size: unaligned, len: 4096, off: 0, expectLen: 4096},
		{name: "direct final unaligned", direct: true, size: unaligned, len: 13, off: unaligned - 13, expectLen: 4096},
		{name: "direct final aligned", direct: true, size: aligned, len: 4096, off: aligned - 4096, expectLen: 4096},
	} {
		t.Run(test.name, func(t *testing.T) {
			ff := &fakeFile{}
			f := &file{File: ff, direct: test.direct, size: test.size}

			p := bytes.Repeat([]byte{'x'}, test.len)
			n, err := f.WriteAt(p, test.off)
			if err != nil {
				t.Fatal(err)
			}
			if n != test.len {
				t.Fatalf("WriteAt n = %d, want %d", n, test.len)
			}

			expectCalls(t, ff, "WriteAt")
			w := ff.writes[0]
			if w.off != test.off {
				t.Fatalf("write offset = %d, want %d", w.off, test.off)
			}
			if len(w.p) != test.expectLen {
				t.Fatalf("write length = %d, want %d", len(w.p), test.expectLen)
			}
			if !bytes.Equal(w.p[:test.len], p) {
				t.Fatal("write does not start with the caller's bytes")
			}
			if pad := w.p[test.len:]; !bytes.Equal(pad, make([]byte, len(pad))) {
				t.Fatal("write padding is not zeroed")
			}
		})
	}
}

func TestFileWriteAtError(t *testing.T) {
	const size = oDirectThreshold + 13

	for _, test := range []struct {
		name   string
		direct bool
		off    int64
	}{
		{name: "buffered", off: 0},
		{name: "direct not final", direct: true, off: 0},
		{name: "direct final", direct: true, off: size - 13},
	} {
		t.Run(test.name, func(t *testing.T) {
			wantErr := errors.New("write failed")
			f := &file{File: &fakeFile{writeErr: wantErr}, direct: test.direct, size: size}

			n, err := f.WriteAt(make([]byte, 13), test.off)
			if !errors.Is(err, wantErr) {
				t.Fatalf("WriteAt error = %v, want %v", err, wantErr)
			}
			if n != 0 {
				t.Fatalf("WriteAt n = %d, want 0", n)
			}
		})
	}
}

func TestFileWriteAtUninitialized(t *testing.T) {
	f := &file{}
	if _, err := f.WriteAt([]byte("x"), 0); err == nil {
		t.Fatal("WriteAt on uninitialized file returned nil error")
	}
}

func TestFileClose(t *testing.T) {
	const size = oDirectThreshold + 13
	truncateErr := errors.New("truncate failed")
	closeErr := errors.New("close failed")

	for _, test := range []struct {
		name        string
		direct      bool
		truncateErr error
		closeErr    error
		expectErr   error
		expectCalls []string
	}{
		{name: "buffered", expectCalls: []string{"Close"}},
		{name: "buffered close error", closeErr: closeErr, expectErr: closeErr, expectCalls: []string{"Close"}},
		{name: "direct", direct: true, expectCalls: []string{"Truncate", "Close"}},
		{name: "direct truncate error", direct: true, truncateErr: truncateErr, expectErr: truncateErr, expectCalls: []string{"Truncate", "Close"}},
		{name: "direct close error", direct: true, closeErr: closeErr, expectErr: closeErr, expectCalls: []string{"Truncate", "Close"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			ff := &fakeFile{truncateErr: test.truncateErr, closeErr: test.closeErr}
			f := &file{File: ff, direct: test.direct, size: size}

			if err := f.Close(); !errors.Is(err, test.expectErr) {
				t.Fatalf("Close error = %v, want %v", err, test.expectErr)
			}
			expectCalls(t, ff, test.expectCalls...)
			if test.direct && ff.truncates[0] != size {
				t.Fatalf("Truncate size = %d, want %d", ff.truncates[0], size)
			}
		})
	}
}

func TestFileCloseUninitialized(t *testing.T) {
	f := &file{}
	if err := f.Close(); err != nil {
		t.Fatalf("Close on uninitialized file = %v, want nil", err)
	}
}

// The type of Statfs_t.Bsize varies by GOARCH.
func setBsize[T ~int32 | ~int64 | ~uint32](p *T, v int64) {
	*p = T(v)
}
