//go:build linux

package transfermanager

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func openDownloadFile(path string) (*os.File, io.WriterAt, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|syscall.O_DIRECT, 0o644)
	if err == nil {
		writer := &fileVectorWriterAt{fd: int(f.Fd()), direct: true}
		return f, newGroupedVectorWriterAt(writer, downloadFileVectorChunkSize, downloadFileVectorCount, true), nil
	}
	if !errors.Is(err, syscall.EINVAL) && !errors.Is(err, syscall.EOPNOTSUPP) {
		return nil, nil, err
	}

	f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, nil, err
	}
	writer := &fileVectorWriterAt{fd: int(f.Fd())}
	return f, newGroupedVectorWriterAt(writer, downloadFileVectorChunkSize, downloadFileVectorCount, false), nil
}

type fileVectorWriterAt struct {
	fd     int
	direct bool
}

func (w *fileVectorWriterAt) writeBufferAlignment() int64 {
	if w.direct {
		return directIOAlignment
	}
	return 1
}

func (w *fileVectorWriterAt) preallocate(size int64) error {
	if !w.direct || size <= 0 {
		return nil
	}
	for {
		err := syscall.Fallocate(w.fd, 0, 0, size)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}

func (w *fileVectorWriterAt) writeVectorAt(vectors [][]byte, off int64) (int, error) {
	if w.direct {
		if err := validateDirectIOVectors(vectors, off); err != nil {
			return 0, err
		}
	}

	for {
		n, err := unix.Pwritev(w.fd, vectors, off)
		if errors.Is(err, syscall.EINTR) {
			continue
		}
		return n, err
	}
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
