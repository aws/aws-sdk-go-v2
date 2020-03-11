package sdk

import "testing"

func TestUUIDVersion4(t *testing.T) {
	uuid := uuidVersion4(make([]byte, 16))
	if e, a := `00000000-0000-4000-8000-000000000000`, uuid; e != a {
		t.Errorf("expect %v uuid, got %v", e, a)
	}

	b := make([]byte, 16)
	for i := 0; i < len(b); i++ {
		b[i] = 1
	}
	uuid = uuidVersion4(b)
	if e, a := `01010101-0101-4101-8101-010101010101`, uuid; e != a {
		t.Errorf("expect %v uuid, got %v", e, a)
	}
}
