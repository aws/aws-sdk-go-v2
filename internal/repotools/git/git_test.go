package git_test

import (
	"github.com/aws/aws-sdk-go-v2/internal/repotools/git"
	"github.com/google/go-cmp/cmp"
	"strconv"
	"testing"
)

func TestToTag(t *testing.T) {
	tests := []struct {
		Path     string
		Version  string
		Expected string
		WantErr  bool
	}{
		{
			Path:     ".",
			Version:  "v1.2.3",
			Expected: "v1.2.3",
		},
		{
			Path:     "service/s3",
			Version:  "v0.1.2",
			Expected: "service/s3/v0.1.2",
		},
		{
			Path:     "service/s3",
			Version:  "v1.2.3",
			Expected: "service/s3/v1.2.3",
		},
		{
			Path:     "service/s3/v2",
			Version:  "v2.3.4",
			Expected: "service/s3/v2.3.4",
		},
		{
			Path:     "service/s3/volumetric",
			Version:  "v1.3.4",
			Expected: "service/s3/volumetric/v1.3.4",
		},
		{
			Path:    "service/s3/v2",
			Version: "v1.3.4",
			WantErr: true,
		},
		{
			Path:    "service/s3/v0",
			Version: "v1.3.4",
			WantErr: true,
		},
		{
			Path:    "service/s3/v1",
			Version: "v1.3.4",
			WantErr: true,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got, err := git.ToModuleTag(tt.Path, tt.Version)
			if (err != nil) != tt.WantErr {
				t.Errorf("WantErr=%v, got err=%v", tt.WantErr, err)
			}
			if tt.Expected != got {
				t.Errorf("expect %v, got %v", tt.Expected, got)
			}
		})
	}
}

func TestParseTags(t *testing.T) {
	tests := map[string]struct {
		args []string
		want git.ModuleTags
	}{
		"standard tags": {
			args: []string{
				"v0.2.1",
				"v1.0.0",
				"v1.3.0",
				"v2.0.0",
				"feature/ec2/imds/v0.1.0",
				"feature/ec2/imds/v1.0.1",
				"feature/ec2/imds/v1.0.6",
				"feature/ec2/imds/v2.0.0",
			},
			want: map[string][]string{
				".":                   {"v1.3.0", "v1.0.0", "v0.2.1"},
				"v2":                  {"v2.0.0"},
				"feature/ec2/imds":    {"v1.0.6", "v1.0.1", "v0.1.0"},
				"feature/ec2/imds/v2": {"v2.0.0"},
			},
		},
		"invalid tags": {
			args: []string{
				"v0.2.1",
				"v1.0.0",
				"release-1-2021-04-09",
				"v1.3.0",
				"1.4.0",
				"feature/ec2/imds/v0.1.0",
				"feature/ec2/imds/v1.0.1",
				"feature/ec2/imds/v1.0.6",
				"feature/ec2/imds@v1.2.0",
			},
			want: map[string][]string{
				".":                {"v1.3.0", "v1.0.0", "v0.2.1"},
				"feature/ec2/imds": {"v1.0.6", "v1.0.1", "v0.1.0"},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := git.ParseModuleTags(tt.args)
			if diff := cmp.Diff(tt.want, got); len(diff) > 0 {
				t.Error(diff)
			}
		})
	}
}

func TestModuleTags_Add(t *testing.T) {
	tags := []string{
		"v0.2.1",
		"v1.0.0",
		"v1.3.0",
	}

	moduleTags := git.ParseModuleTags(tags)

	for _, tag := range []string{"v0.2.1", "v1.3.0", "v0.2.0", "v1.0.0", "v1.4.0", "v2.0.0", "feature/ec2/imds/v1.0.0"} {
		if ok := moduleTags.Add(tag); !ok {
			t.Errorf("expect tag to have been inserted")
		}
	}

	if ok := moduleTags.Add("invalid-tag"); ok {
		t.Errorf("expected tag to not be inserted")
	}

	want := git.ModuleTags{
		".":                {"v1.4.0", "v1.3.0", "v1.0.0", "v0.2.1", "v0.2.0"},
		"v2":               {"v2.0.0"},
		"feature/ec2/imds": {"v1.0.0"},
	}

	if diff := cmp.Diff(want, moduleTags); len(diff) > 0 {
		t.Errorf(diff)
	}
}
