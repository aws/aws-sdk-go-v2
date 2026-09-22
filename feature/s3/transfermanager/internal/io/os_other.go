//go:build !linux

package io

import "os"

type file struct {
	*os.File
}

func (*file) Init(_, _ int64) error {
	return nil
}

// Create creates the named file.
//
// Create in non-linux contexts just delegates to os.Create for now.
func Create(path string) (File, error) {
	f, err := os.Create(path)
	return &file{f}, err
}
