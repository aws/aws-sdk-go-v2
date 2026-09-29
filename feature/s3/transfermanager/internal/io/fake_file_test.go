package io

import (
	"os"
	"testing"
)

type fakeWrite struct {
	p   []byte
	off int64
}

type fakeFile struct {
	fd          uintptr
	writeErr    error
	truncateErr error
	closeErr    error

	calls     []string
	writes    []fakeWrite
	truncates []int64
}

func (f *fakeFile) WriteAt(p []byte, off int64) (int, error) {
	f.calls = append(f.calls, "WriteAt")
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	f.writes = append(f.writes, fakeWrite{p: append([]byte(nil), p...), off: off})
	return len(p), nil
}

func (f *fakeFile) Truncate(size int64) error {
	f.calls = append(f.calls, "Truncate")
	f.truncates = append(f.truncates, size)
	return f.truncateErr
}

func (f *fakeFile) Close() error {
	f.calls = append(f.calls, "Close")
	return f.closeErr
}

func (f *fakeFile) Fd() uintptr {
	return f.fd
}

type openCall struct {
	name string
	flag int
	perm os.FileMode
}

// stubOpenFile replaces openFile for the duration of the test. It returns ff,
// or err if non-nil, and records every call.
func stubOpenFile(t *testing.T, ff *fakeFile, err error) *[]openCall {
	t.Helper()

	original := openFile
	t.Cleanup(func() { openFile = original })

	var calls []openCall
	openFile = func(name string, flag int, perm os.FileMode) (osFile, error) {
		calls = append(calls, openCall{name: name, flag: flag, perm: perm})
		if err != nil {
			return nil, err
		}
		return ff, nil
	}
	return &calls
}

func expectCalls(t *testing.T, ff *fakeFile, want ...string) {
	t.Helper()

	if len(ff.calls) != len(want) {
		t.Fatalf("calls = %v, want %v", ff.calls, want)
	}
	for i := range want {
		if ff.calls[i] != want[i] {
			t.Fatalf("calls = %v, want %v", ff.calls, want)
		}
	}
}
