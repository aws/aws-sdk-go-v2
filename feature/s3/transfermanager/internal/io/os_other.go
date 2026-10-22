//go:build !linux

package io

import "os"

type file struct {
	osFile
}

func (*file) Init(_, _, _ int64, _ bool) error {
	return nil
}

// Create creates the named file. The file must not already exist.
//
// Create in non-linux contexts just delegates to os.OpenFile for now.
func Create(path string) (File, error) {
	f, err := openFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o666)
	if err != nil {
		return nil, err
	}
	return &file{f}, nil
}
