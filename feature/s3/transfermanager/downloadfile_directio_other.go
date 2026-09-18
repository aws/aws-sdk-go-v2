//go:build !linux

package transfermanager

import (
	"io"
	"os"
)

func openDownloadFile(path string) (*os.File, io.WriterAt, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, nil, err
	}
	return f, f, nil
}
