package release

import (
	"github.com/aws/aws-sdk-go-v2/internal/repotools"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"testing"
	"time"
)

type mockFinder struct {
	RootPath string
	Modules  map[string][]string
}

func (m *mockFinder) Root() string {
	return m.RootPath
}

func (m *mockFinder) ModulesRel() (map[string][]string, error) {
	return m.Modules, nil
}

func TestCalculateNextVersion(t *testing.T) {
	type args struct {
		modulePath  string
		latest      string
		config      repotools.ModuleConfig
		annotations []changelog.Annotation
	}
	tests := map[string]struct {
		args     args
		wantNext string
		wantErr  bool
	}{
		"new module v1 major": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/shinynew",
			},
			wantNext: "v1.0.0-preview",
		},
		"new module v1 major with release annotation": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/shinynew",
				annotations: []changelog.Annotation{{
					Type: changelog.ReleaseChangeType,
				}},
			},
			wantNext: "v1.0.0",
		},
		"new module v2 or higher major": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/shinynew/v2",
			},
			wantNext: "v2.0.0-preview",
		},
		"new module v2 or higher with release annotation": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/shinynew/v2",
				annotations: []changelog.Annotation{{
					Type: changelog.ReleaseChangeType,
				}},
			},
			wantNext: "v2.0.0",
		},
		"existing module version, not pre-release, no annotation": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/existing",
				latest:     "v1.0.0",
			},
			wantNext: "v1.0.1",
		},
		"existing module version, not pre-release, with patch semver annotation": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/existing",
				latest:     "v1.0.0",
				annotations: []changelog.Annotation{
					{Type: changelog.BugFixChangeType},
				},
			},
			wantNext: "v1.0.1",
		},
		"existing module version, not pre-release, with minor semver annotation": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/existing",
				latest:     "v1.0.1",
				annotations: []changelog.Annotation{
					{Type: changelog.FeatureChangeType},
				},
			},
			wantNext: "v1.1.0",
		},
		"existing module version, set for pre-release": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/existing",
				latest:     "v1.0.1",
				config:     repotools.ModuleConfig{PreRelease: "rc"},
			},
			wantNext: "v1.1.0-rc",
		},
		"existing module preview version": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/existing",
				latest:     "v1.1.0-preview",
				config:     repotools.ModuleConfig{PreRelease: "preview"},
			},
			wantNext: "v1.1.0-preview.1",
		},
		"existing module preview version, with non-release annotation types": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/existing",
				latest:     "v1.1.0-preview.1",
				config:     repotools.ModuleConfig{PreRelease: "preview"},
				annotations: []changelog.Annotation{{
					Type: changelog.FeatureChangeType,
				}},
			},
			wantNext: "v1.1.0-preview.2",
		},
		"existing module preview version, with new pre-release tag": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/existing",
				latest:     "v1.1.0-preview.2",
				config:     repotools.ModuleConfig{PreRelease: "rc"},
				annotations: []changelog.Annotation{{
					Type: changelog.FeatureChangeType,
				}},
			},
			wantNext: "v1.1.0-rc",
		},
		"existing module preview version, with new invalid pre-release tag": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/existing",
				latest:     "v1.1.0-rc.5",
				config:     repotools.ModuleConfig{PreRelease: "alpha"},
				annotations: []changelog.Annotation{{
					Type: changelog.FeatureChangeType,
				}},
			},
			wantErr: true,
		},
		"existing module preview version, with release annotation": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/existing",
				latest:     "v1.1.0-rc.5",
				annotations: []changelog.Annotation{{
					Type: changelog.ReleaseChangeType,
				}},
			},
			wantNext: "v1.1.0",
		},
		"invalid latest tag": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/existing",
				latest:     "1.1.0",
			},
			wantErr: true,
		},
		"module tag with build metadata": {
			args: args{
				modulePath: "github.com/aws/aws-sdk-go-v2/service/existing",
				latest:     "v1.1.0+build.12345",
			},
			wantNext: "v1.1.1",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotNext, err := CalculateNextVersion(tt.args.modulePath, tt.args.latest, tt.args.config, tt.args.annotations)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateNextVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNext != tt.wantNext {
				t.Errorf("CalculateNextVersion() gotNext = %v, want %v", gotNext, tt.wantNext)
			}
		})
	}
}

func TestNextReleaseID(t *testing.T) {
	origNowTime := nowTime
	defer func() {
		nowTime = origNowTime
	}()

	type args struct {
		tags []string
	}
	tests := map[string]struct {
		args     args
		nowTime  func() time.Time
		wantNext string
	}{
		"no tags": {
			wantNext: "2021-05-06",
		},
		"other tags": {
			args:     args{tags: []string{"v1.2.0", "release/foo/v2"}},
			wantNext: "2021-05-06",
		},
		"older tags": {
			args:     args{tags: []string{"release-2021-05-04", "release-2021-05-04.2"}},
			wantNext: "2021-05-06",
		},
		"second release": {
			args:     args{tags: []string{"release-2021-05-06"}},
			wantNext: "2021-05-06.2",
		},
		"third release": {
			args:     args{tags: []string{"release-2021-05-06", "release-2021-05-06.2"}},
			wantNext: "2021-05-06.3",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.nowTime == nil {
				nowTime = func() time.Time {
					return time.Date(2021, 5, 6, 7, 8, 9, 10, time.UTC)
				}
			} else {
				nowTime = tt.nowTime
			}

			if gotNext := NextReleaseID(tt.args.tags); gotNext != tt.wantNext {
				t.Errorf("NextReleaseID() = %v, want %v", gotNext, tt.wantNext)
			}
		})
	}
}
