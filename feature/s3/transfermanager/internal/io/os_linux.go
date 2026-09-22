//go:build linux

package io

import (
	"errors"
	"os"
	"syscall"
)

const oDirectThreshold = 64 * 1024 * 1024 // 64MiB

// Linux files open with O_DIRECT above a size threshold when the write size is
// aligned. This bypasses the page cache and an inode lock, which drastically
// improves performance for writes that are sustained enough.
type file struct {
	*os.File
	direct bool
	path   string
	size   int64
}

func (f *file) WriteAt(p []byte, off int64) (int, error) {
	if f.File == nil {
		return 0, errors.New("file was not initialized")
	}
	if !f.direct || off+int64(len(p)) != f.size {
		return f.File.WriteAt(p, off)
	}

	// last write needs pad
	padded := makealigned(len(p) + int(align(f.size)-f.size))
	copy(padded, p) // yes it's a copy but it's only the last write

	_, err := f.File.WriteAt(padded, off)
	if err != nil {
		return 0, err
	}

	return len(p), err
}

func (f *file) Init(size, writeSize int64) error {
	if f.File != nil {
		return errors.New("file was already initialized")
	}

	f.size = size
	if size < oDirectThreshold || writeSize%alignedBy != 0 {
		ff, err := os.Create(f.path)
		f.File = ff
		return err
	}

	ff, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|syscall.O_DIRECT, 0o644)
	if err != nil {
		return err
	}

	if err := syscall.Fallocate(int(ff.Fd()), 0, 0, size); err != nil {
		return err
	}

	f.File = ff
	f.direct = true
	return nil
}

func (f *file) Close() error {
	if f.File == nil {
		return nil
	}

	if f.direct {
		if err := f.File.Truncate(f.size); err != nil {
			return err
		}
	}
	return f.File.Close()
}

// Create creates the named file.
//
// Create on Linux returns a lazy wrapper. Actual file creation is delayed until
// the file size and write chunk size are known.
func Create(path string) (File, error) {
	return &file{path: path}, nil
}
