//go:build linux

package io

import (
	"errors"
	"os"
	"syscall"
	"testing"
)

func TestSupportsDirectIOFilesystemBlockSize(t *testing.T) {
	originalStatfs := statfs
	defer func() { statfs = originalStatfs }()

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
			statfs = func(_ string, stat *syscall.Statfs_t) error {
				setBsize(&stat.Bsize, test.bsize)
				return nil
			}

			got := supportsDirectIO("/tmp/file", 10*1024*1024, 8*1024*1024)
			if got != test.allowed {
				t.Fatalf("supportsDirectIO = %v, want %v", got, test.allowed)
			}
		})
	}
}

func TestSupportsDirectIOAlignment(t *testing.T) {
	originalStatfs := statfs
	defer func() { statfs = originalStatfs }()
	statfs = func(_ string, stat *syscall.Statfs_t) error {
		stat.Bsize = 4096
		return nil
	}

	for _, test := range []struct {
		name      string
		partSize  int64
		writeSize int64
		allowed   bool
	}{
		{name: "part remainder is aligned", partSize: 10 * 1024 * 1024, writeSize: 8 * 1024 * 1024, allowed: true},
		{name: "part remainder is not aligned", partSize: 10*1024*1024 + 1, writeSize: 8 * 1024 * 1024},
		{name: "write size is not aligned", partSize: 10 * 1024 * 1024, writeSize: 8*1024*1024 + 1},
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
	originalStatfs := statfs
	defer func() { statfs = originalStatfs }()

	for _, test := range []struct {
		name      string
		size      int64
		bsize     int64
		statErr   error
		partSize  int64
		writeSize int64
	}{
		{name: "size threshold", size: oDirectThreshold - 1, bsize: 4096, partSize: 10 * 1024 * 1024, writeSize: 8 * 1024 * 1024},
		{name: "filesystem block size", size: oDirectThreshold, bsize: 1024, partSize: 10 * 1024 * 1024, writeSize: 8 * 1024 * 1024},
		{name: "filesystem stat error", size: oDirectThreshold, bsize: 4096, statErr: errors.New("statfs failed"), partSize: 10 * 1024 * 1024, writeSize: 8 * 1024 * 1024},
		{name: "part size", size: oDirectThreshold, bsize: 4096, partSize: 10*1024*1024 + 1, writeSize: 8 * 1024 * 1024},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := t.TempDir() + "/file"
			created, err := Create(path)
			if err != nil {
				t.Fatal(err)
			}
			f := created.(*file)
			statfs = func(_ string, stat *syscall.Statfs_t) error {
				setBsize(&stat.Bsize, test.bsize)
				return test.statErr
			}

			if err := f.Init(test.size, test.partSize, test.writeSize, true); err != nil {
				t.Fatal(err)
			}
			if f.direct {
				t.Fatal("file initialized with direct I/O")
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(path); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestFileInitDirectIO(t *testing.T) {
	originalStatfs, originalOpenFile, originalFallocate := statfs, openFile, fallocate
	defer func() {
		statfs = originalStatfs
		openFile = originalOpenFile
		fallocate = originalFallocate
	}()

	statfs = func(_ string, stat *syscall.Statfs_t) error {
		stat.Bsize = 4096
		return nil
	}
	var openFlags int
	openFile = func(path string, flags int, perm os.FileMode) (*os.File, error) {
		openFlags = flags
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	}
	var fallocateSize int64
	fallocate = func(_ int, _ uint32, _ int64, size int64) error {
		fallocateSize = size
		return nil
	}

	created, err := Create(t.TempDir() + "/file")
	if err != nil {
		t.Fatal(err)
	}
	f := created.(*file)
	const size = oDirectThreshold
	if err := f.Init(size, 10*1024*1024, 8*1024*1024, true); err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if !f.direct {
		t.Fatal("file did not initialize with direct I/O")
	}
	if openFlags&syscall.O_DIRECT == 0 {
		t.Fatal("file was opened without O_DIRECT")
	}
	if fallocateSize != size {
		t.Fatalf("fallocate size = %d, want %d", fallocateSize, size)
	}
}

func TestFileInitAlreadyInitialized(t *testing.T) {
	created, err := Create(t.TempDir() + "/file")
	if err != nil {
		t.Fatal(err)
	}
	f := created.(*file)
	if err := f.Init(1, 1, 1, false); err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := f.Init(1, 1, 1, false); err == nil {
		t.Fatal("second Init returned nil error")
	}
}

func TestFileInitOpenError(t *testing.T) {
	originalStatfs, originalOpenFile := statfs, openFile
	defer func() {
		statfs = originalStatfs
		openFile = originalOpenFile
	}()

	statfs = func(_ string, stat *syscall.Statfs_t) error {
		stat.Bsize = 4096
		return nil
	}
	wantErr := errors.New("open failed")
	openFile = func(string, int, os.FileMode) (*os.File, error) {
		return nil, wantErr
	}

	created, err := Create(t.TempDir() + "/file")
	if err != nil {
		t.Fatal(err)
	}
	f := created.(*file)
	if err := f.Init(oDirectThreshold, 10*1024*1024, 8*1024*1024, true); !errors.Is(err, wantErr) {
		t.Fatalf("Init error = %v, want %v", err, wantErr)
	}
	if f.File != nil {
		t.Fatal("file handle set after open failure")
	}
}

func TestFileInitFallocateError(t *testing.T) {
	originalStatfs, originalOpenFile, originalFallocate := statfs, openFile, fallocate
	defer func() {
		statfs = originalStatfs
		openFile = originalOpenFile
		fallocate = originalFallocate
	}()

	statfs = func(_ string, stat *syscall.Statfs_t) error {
		stat.Bsize = 4096
		return nil
	}
	openFile = func(path string, _ int, perm os.FileMode) (*os.File, error) {
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	}
	wantErr := errors.New("fallocate failed")
	fallocate = func(int, uint32, int64, int64) error {
		return wantErr
	}

	created, err := Create(t.TempDir() + "/file")
	if err != nil {
		t.Fatal(err)
	}
	f := created.(*file)
	if err := f.Init(oDirectThreshold, 10*1024*1024, 8*1024*1024, true); !errors.Is(err, wantErr) {
		t.Fatalf("Init error = %v, want %v", err, wantErr)
	}
	if f.File != nil {
		t.Fatal("file handle set after fallocate failure")
	}
}

func TestFileInitDirectIODisabled(t *testing.T) {
	originalOpenFile := openFile
	defer func() { openFile = originalOpenFile }()
	openFile = func(string, int, os.FileMode) (*os.File, error) {
		return nil, errors.New("direct open should not be called")
	}

	created, err := Create(t.TempDir() + "/file")
	if err != nil {
		t.Fatal(err)
	}
	f := created.(*file)
	if err := f.Init(oDirectThreshold, 10*1024*1024, 8*1024*1024, false); err != nil {
		t.Fatal(err)
	}
	if f.direct {
		t.Fatal("file initialized with direct I/O")
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

// The type of Statfs_t.Bsize varies by GOARCH.
func setBsize[T ~int32 | ~int64 | ~uint32](p *T, v int64) {
	*p = T(v)
}
