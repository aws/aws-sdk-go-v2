//go:build !linux

package io

import (
	"errors"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	ff := &fakeFile{}
	opens := stubOpenFile(t, ff, nil)

	f, err := Create("/dir/file")
	if err != nil {
		t.Fatal(err)
	}

	want := openCall{name: "/dir/file", flag: os.O_RDWR | os.O_CREATE | os.O_EXCL, perm: 0o666}
	if len(*opens) != 1 || (*opens)[0] != want {
		t.Fatalf("openFile calls = %+v, want [%+v]", *opens, want)
	}

	if _, err := f.WriteAt([]byte("abc"), 7); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	expectCalls(t, ff, "WriteAt", "Close")
	if w := ff.writes[0]; string(w.p) != "abc" || w.off != 7 {
		t.Fatalf("write = %+v, want abc at 7", w)
	}
}

func TestCreateOpenError(t *testing.T) {
	wantErr := errors.New("open failed")
	stubOpenFile(t, nil, wantErr)

	f, err := Create("/dir/file")
	if !errors.Is(err, wantErr) {
		t.Fatalf("Create error = %v, want %v", err, wantErr)
	}
	if f != nil {
		t.Fatalf("Create returned %v on error", f)
	}
}
