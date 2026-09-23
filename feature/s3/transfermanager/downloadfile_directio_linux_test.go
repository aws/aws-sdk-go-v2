//go:build linux

package transfermanager

import (
	"bytes"
	"errors"
	"os"
	"syscall"
	"testing"
)

func TestDirectFileWriterAtPadsPhysicalWrite(t *testing.T) {
	f := openDirectTestFile(t)
	defer f.Close()

	want := []byte("short direct write")
	buf := newWriteBuffer(directIOAlignment, directIOAlignment)
	copy(buf, want)
	writer := &fileWriterAt{file: f, direct: true}
	if n, err := writer.WriteAt(buf[:len(want)], 0); err != nil {
		t.Fatalf("WriteAt: %v", err)
	} else if n != len(want) {
		t.Fatalf("WriteAt count = %d, want %d", n, len(want))
	}

	got, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(got) != directIOAlignment {
		t.Fatalf("physical file size = %d, want %d", len(got), directIOAlignment)
	}
	if !bytes.Equal(got[:len(want)], want) {
		t.Fatalf("written prefix = %q, want %q", got[:len(want)], want)
	}
	if !bytes.Equal(got[len(want):], make([]byte, directIOAlignment-len(want))) {
		t.Fatal("physical write padding is not zero-filled")
	}
}

func TestDirectFileWriterAtWritesSeparateChunks(t *testing.T) {
	f := openDirectTestFile(t)
	name := f.Name()

	writer := &fileWriterAt{file: f, direct: true}
	want := make([]byte, 4*directIOAlignment)
	for i := 0; i < 4; i++ {
		buf := newWriteBuffer(directIOAlignment, directIOAlignment)
		for j := range buf {
			buf[j] = byte(i + 1)
		}
		copy(want[i*directIOAlignment:], buf)
		if n, err := writer.WriteAt(buf, int64(i*directIOAlignment)); err != nil {
			t.Fatalf("WriteAt %d: %v", i, err)
		} else if n != len(buf) {
			t.Fatalf("WriteAt %d count = %d, want %d", i, n, len(buf))
		}
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	got, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("separate direct writes do not match source buffers")
	}
}

func openDirectTestFile(t *testing.T) *os.File {
	t.Helper()
	path := t.TempDir() + "/direct-write.bin"
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|syscall.O_DIRECT, 0o644)
	if errors.Is(err, syscall.EINVAL) || errors.Is(err, syscall.EOPNOTSUPP) {
		t.Skipf("test filesystem does not support O_DIRECT: %v", err)
	}
	if err != nil {
		t.Fatalf("OpenFile with O_DIRECT: %v", err)
	}
	return f
}

func TestDirectFileWriterAtRejectsUnalignedOffset(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "direct-write-*.bin")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer f.Close()

	writer := &fileWriterAt{file: f, direct: true}
	buf := newWriteBuffer(directIOAlignment, directIOAlignment)
	if _, err := writer.WriteAt(buf, 1); err == nil {
		t.Fatal("WriteAt with unaligned offset returned no error")
	}
}

func TestDirectFileWriterAtPreallocates(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "direct-preallocate-*.bin")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer f.Close()

	writer := &fileWriterAt{file: f, direct: true}
	want := int64(2 * directIOAlignment)
	if err := writer.preallocate(want); err != nil {
		t.Fatalf("preallocate: %v", err)
	}
	info, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Size() != want {
		t.Fatalf("file size = %d, want %d", info.Size(), want)
	}
}
