package changes

import "testing"

func TestGetCurrentModule(t *testing.T) {
	mod, err := GetCurrentModule()
	if err != nil {
		t.Errorf("expected nil err, got %v", err)
	}

	if mod != "internal/tools/changes" {
		t.Errorf("expected mod to be \"internal/tools/changes\", got %s", mod)
	}
}
