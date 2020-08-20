// +build disabled

package external

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
)

func TestResolveEC2Region(t *testing.T) {
	configs := Configs{}

	cfg := unit.Config()

	err := ResolveEC2Region(&cfg, configs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if e, a := "mock-region", cfg.Region; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	resetOrig := swapEC2MetadataNew(func(config aws.Config) ec2MetadataRegionClient {
		return mockEC2MetadataClient{
			retRegion: "foo-region",
		}
	})
	defer resetOrig()

	cfg.Region = ""
	err = ResolveEC2Region(&cfg, configs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if e, a := "foo-region", cfg.Region; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	_ = swapEC2MetadataNew(func(config aws.Config) ec2MetadataRegionClient {
		return mockEC2MetadataClient{
			retErr: fmt.Errorf("some error"),
		}
	})

	cfg.Region = ""
	err = ResolveEC2Region(&cfg, configs)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(cfg.Region) != 0 {
		t.Errorf("expected region to remain unset")
	}
}

type mockEC2MetadataClient struct {
	retRegion string
	retErr    error
}

func (m mockEC2MetadataClient) Region(ctx context.Context) (string, error) {
	if m.retErr != nil {
		return "", m.retErr
	}

	return m.retRegion, nil
}

func swapEC2MetadataNew(f func(config aws.Config) ec2MetadataRegionClient) func() {
	orig := newEC2MetadataClient
	newEC2MetadataClient = f
	return func() {
		newEC2MetadataClient = orig
	}
}
