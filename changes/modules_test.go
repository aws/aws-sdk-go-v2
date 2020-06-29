package changes

import "testing"

func TestGetCurrentModule(t *testing.T) {
	mod, err := GetCurrentModule()
	if err != nil {
		t.Errorf("expected nil err, got %v", err)
	}

	// TODO: update to reflect path in SDK repo.
	if mod != "changes" {
		t.Errorf("expected mod to be \"changes\", got %s", mod)
	}
}
