package io

import (
	"errors"
	"os"
	"syscall"
)

const oDirectThreshold = 64 * 1024 * 1024 // 64MiB

// Linux files will open w/ O_DIRECT above a certain size threshold, which
// bypasses page cache as well as a kernal inode lock, drastically improving
// write performance.
type file struct {
	*os.File
	direct    bool
	path      string
	size      int64
	writeSize int64
}

func (f *file) WriteAt(p []byte, off int64) (n int, err error) {
	if f.File == nil {
		return 0, errors.New("file wasn't created")
	}

	if !f.direct {
		return f.File.WriteAt(p, off)
	}

	// TODO: if it's the last write then pad
	return f.File.WriteAt(p, off)
}

func (f *file) Init(size, writeSize int64) error {
	if size < oDirectThreshold || writeSize%alignedBy != 0 {
		ff, err := os.Create(f.path)
		f.File = ff
		return err
	}

	ff, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|syscall.O_DIRECT, 0o644)
	if err := syscall.Fallocate(ff.Fd(), size /*TODO*/); err != nil {
		return err
	}

	f.direct = true
	return nil
}

func (f *file) Close() error {
	if f.File == nil {
		return nil
	}
	return f.Close()
}

// Create creates the named file.
//
// Create on Linux will return a lazy wrapper, actual file creation is delayed
// until the point at which the file and write chunk sizes are known.
func Create(path string) (File, error) {
	return &file{
		path: path,
	}, err
}
