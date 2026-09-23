//go:build linux

package transfermanager

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

func openDownloadFile(path string) (*os.File, io.WriterAt, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|syscall.O_DIRECT, 0o644)
	if err == nil {
		return f, &fileWriterAt{file: f, direct: true}, nil
	}
	if !errors.Is(err, syscall.EINVAL) && !errors.Is(err, syscall.EOPNOTSUPP) {
		return nil, nil, err
	}

	f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, nil, err
	}
	return f, &fileWriterAt{file: f}, nil
}

type fileWriterAt struct {
	file   *os.File
	direct bool
}

func (w *fileWriterAt) writeBufferAlignment() int64 {
	if w.direct {
		return directIOAlignment
	}
	return 1
}

func (w *fileWriterAt) preallocate(size int64) error {
	if !w.direct || size <= 0 {
		return nil
	}
	for {
		err := syscall.Fallocate(int(w.file.Fd()), 0, 0, size)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (w *fileWriterAt) WriteAt(p []byte, off int64) (int, error) {
	logicalLen := len(p)
	if w.direct {
		if err := validateDirectIOVectors([][]byte{p}, off); err != nil {
			return 0, err
		}

		physicalLen := logicalLen
		if remainder := physicalLen % directIOAlignment; remainder != 0 {
			physicalLen += directIOAlignment - remainder
			if physicalLen > cap(p) {
				return 0, fmt.Errorf("write buffer capacity %d is smaller than padded length %d", cap(p), physicalLen)
			}
			p = p[:physicalLen]
			clear(p[logicalLen:])
		}
	}

	n, err := w.file.WriteAt(p, off)
	return min(n, logicalLen), err
}

func validateDirectIOVectors(vectors [][]byte, off int64) error {
	if off%directIOAlignment != 0 {
		return fmt.Errorf("O_DIRECT write offset %d is not aligned to %d bytes", off, directIOAlignment)
	}
	for _, vector := range vectors {
		if len(vector)%directIOAlignment != 0 {
			return fmt.Errorf("O_DIRECT write length %d is not aligned to %d bytes", len(vector), directIOAlignment)
		}
		if len(vector) > 0 && uintptr(unsafe.Pointer(unsafe.SliceData(vector)))%directIOAlignment != 0 {
			return fmt.Errorf("O_DIRECT write buffer is not aligned to %d bytes", directIOAlignment)
		}
	}
	return nil
}
