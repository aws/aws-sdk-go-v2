package logincreds

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
	"testing"
)

// the signature is random so the only way to test that is to actually
// cryptographically verify the signature
func TestDPOP(t *testing.T) {
	token, key := mockToken()
	dpop, err := mkdpop(token, "htu")
	if err != nil {
		t.Fatal(err)
	}

	parts := strings.Split(dpop, ".")
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		t.Fatal(err)
	}

	msg := fmt.Sprintf("%s.%s", parts[0], parts[1])
	h := sha256.New()
	h.Write([]byte(msg))

	r := new(big.Int).SetBytes(sig[0:curvelen])
	s := new(big.Int).SetBytes(sig[curvelen:])
	if !ecdsa.Verify(&key.PublicKey, h.Sum(nil), r, s) {
		t.Error("ecdsa signature verify failed")
	}
}
