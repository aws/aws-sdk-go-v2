package golist

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseGoList(t *testing.T) {
	out, err := ioutil.ReadFile(filepath.Join("testdata", "golist.json"))
	if err != nil {
		t.Fatal(err)
	}

	packages, err := parseGoList(out)
	if err != nil {
		t.Fatal(err)
	}

	wantPackages := []string{
		"github.com/aws/aws-sdk-go-v2/internal/tools/changes",
		"github.com/aws/aws-sdk-go-v2/internal/tools/changes/cmd/changetool",
	}

	if diff := cmp.Diff(wantPackages, packages); diff != "" {
		t.Errorf("expect packages to match:\n%v", diff)
	}
}
