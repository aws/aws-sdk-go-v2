//go:build linux

package io

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
)

const oDirectThreshold = 64 * 1024 * 1024 // 64MiB

var (
	statfs    = syscall.Statfs
	openFile  = os.OpenFile
	fallocate = syscall.Fallocate
)

// Linux files open with O_DIRECT above a size threshold when the filesystem and
// transfer sizes are aligned. This bypasses the page cache and an inode lock,
// which drastically improves performance for writes that are sustained enough.
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

func (f *file) Init(size, partSize, writeSize int64, directIO bool) error {
	if f.File != nil {
		return errors.New("file was already initialized")
	}

	f.size = size
	if size < oDirectThreshold || !directIO || !supportsDirectIO(f.path, partSize, writeSize) {
		ff, err := os.Create(f.path)
		f.File = ff
		return err
	}

	ff, err := openFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|syscall.O_DIRECT, 0o644)
	if err != nil {
		return err
	}

	if err := fallocate(int(ff.Fd()), 0, 0, size); err != nil {
		_ = ff.Close()
		return err
	}

	f.File = ff
	f.direct = true
	return nil
}

func supportsDirectIO(path string, partSize, writeSize int64) bool {
	if writeSize <= 0 || writeSize%alignedBy != 0 {
		return false
	}
	if partSize <= 0 || partSize%writeSize%alignedBy != 0 {
		return false
	}

	var stat syscall.Statfs_t
	if err := statfs(filepath.Dir(path), &stat); err != nil {
		return false
	}

	return stat.Bsize > 0 && alignedBy%stat.Bsize == 0
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
