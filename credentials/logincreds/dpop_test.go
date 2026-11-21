package logincreds

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/signin"
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

func p256KeyFromScalar(d int64) *ecdsa.PrivateKey {
	curve := elliptic.P256()
	key := new(ecdsa.PrivateKey)
	key.PublicKey.Curve = curve
	key.D = big.NewInt(d)
	key.X, key.Y = curve.ScalarBaseMult(key.D.Bytes())
	return key
}

func tokenFromKey(key *ecdsa.PrivateKey) *loginToken {
	der, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		panic(err)
	}
	return &loginToken{
		AccessToken: &loginTokenAccessToken{
			AccessKeyID:     "AKID",
			SecretAccessKey: "SECRET",
			SessionToken:    "SESSION",
			AccountID:       "ACCOUNTID",
			ExpiresAt:       time.Unix(0, 0).UTC(),
		},
		TokenType:     "TokenType",
		RefreshToken:  "RefreshToken",
		IdentityToken: "IdentityToken",
		ClientID:      "ClientID",
		DPOPKey: string(pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: der,
		})),
	}
}

func TestDPOP_JWKCoordinateWidth(t *testing.T) {
	cases := []struct {
		name       string
		scalar     int64
		shortCoord string
	}{
		{name: "short X", scalar: 379, shortCoord: "x"},
		{name: "short Y", scalar: 43, shortCoord: "y"},
		{name: "ordinary width", scalar: 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			key := p256KeyFromScalar(c.scalar)
			switch c.shortCoord {
			case "x":
				if got := len(key.X.Bytes()); got != curvelen-1 {
					t.Fatalf("fixture X width: %d, want %d (leading-zero case)", got, curvelen-1)
				}
			case "y":
				if got := len(key.Y.Bytes()); got != curvelen-1 {
					t.Fatalf("fixture Y width: %d, want %d (leading-zero case)", got, curvelen-1)
				}
			default:
				if got := len(key.X.Bytes()); got != curvelen {
					t.Fatalf("control X width: %d, want %d", got, curvelen)
				}
				if got := len(key.Y.Bytes()); got != curvelen {
					t.Fatalf("control Y width: %d, want %d", got, curvelen)
				}
			}

			dpop, err := mkdpop(tokenFromKey(key), "htu")
			if err != nil {
				t.Fatal(err)
			}
			assertDPOPJWK(t, dpop, &key.PublicKey)
		})
	}
}

func TestRetrieve_DPOPJWKCoordinateWidth(t *testing.T) {
	cases := []struct {
		name   string
		scalar int64
	}{
		{name: "short X", scalar: 379},
		{name: "short Y", scalar: 43},
		{name: "ordinary width", scalar: 1},
	}

	now := time.Unix(60, 0).UTC()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			key := p256KeyFromScalar(c.scalar)
			token := tokenFromKey(key)
			restoreOpen := mockOpenFile(token, nil)
			_, restoreCreate := mockCreateFile(nil)
			restoreNow := mockNowTime(now)
			defer restoreOpen()
			defer restoreCreate()
			defer restoreNow()

			var (
				requests int
				dpop     string
			)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				dpop = r.Header.Get("DPoP")
				if r.Method != http.MethodPost {
					http.Error(w, "method", http.StatusMethodNotAllowed)
					return
				}
				if r.URL.Path != "/v1/token" {
					http.Error(w, "path", http.StatusNotFound)
					return
				}
				if dpop == "" {
					http.Error(w, "missing DPoP", http.StatusBadRequest)
					return
				}
				if !dpopProofAcceptable(dpop, now) {
					http.Error(w, "invalid DPoP", http.StatusBadRequest)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"accessToken":{"accessKeyId":"NEW_AKID","secretAccessKey":"NEW_SECRET","sessionToken":"NEW_SESSION"},"expiresIn":900,"refreshToken":"NewRefreshToken","tokenType":"aws_sigv4"}`)
			}))
			defer srv.Close()

			client := signin.New(signin.Options{
				Region:       "us-east-1",
				BaseEndpoint: aws.String(srv.URL),
				Credentials:  aws.AnonymousCredentials{},
				Retryer:      aws.NopRetryer{},
			})

			_, err := New(client, "mocktokenpath").Retrieve(context.Background())
			if requests != 1 {
				t.Errorf("token endpoint requests: %d, want 1", requests)
			}
			if dpop == "" {
				t.Fatal("missing DPoP header")
			}
			assertDPOPJWK(t, dpop, &key.PublicKey)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func assertDPOPJWK(t *testing.T, dpop string, orig *ecdsa.PublicKey) {
	t.Helper()

	parts := strings.Split(dpop, ".")
	if len(parts) != 3 {
		t.Fatalf("dpop jwt parts: %d, want 3", len(parts))
	}

	hb, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("decode header: %v", err)
	}
	var hdr dpopHeader
	if err := json.Unmarshal(hb, &hdr); err != nil {
		t.Fatalf("unmarshal header: %v", err)
	}
	if hdr.Jwk == nil {
		t.Fatal("missing jwk")
	}

	xb, err := base64.RawURLEncoding.DecodeString(hdr.Jwk.X)
	if err != nil {
		t.Fatalf("decode jwk.x: %v", err)
	}
	yb, err := base64.RawURLEncoding.DecodeString(hdr.Jwk.Y)
	if err != nil {
		t.Fatalf("decode jwk.y: %v", err)
	}
	if len(xb) != curvelen {
		t.Errorf("jwk.x decoded length %d, want %d", len(xb), curvelen)
	}
	if len(yb) != curvelen {
		t.Errorf("jwk.y decoded length %d, want %d", len(yb), curvelen)
	}

	x := new(big.Int).SetBytes(xb)
	y := new(big.Int).SetBytes(yb)
	curve := elliptic.P256()
	if !curve.IsOnCurve(x, y) {
		t.Fatal("jwk point is not on P-256")
	}
	if orig.X.Cmp(x) != 0 || orig.Y.Cmp(y) != 0 {
		t.Fatal("reconstructed jwk point does not match original public key")
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		t.Fatalf("decode signature: %v", err)
	}
	if len(sig) != curvelen*2 {
		t.Fatalf("signature length %d, want %d", len(sig), curvelen*2)
	}

	msg := fmt.Sprintf("%s.%s", parts[0], parts[1])
	h := sha256.New()
	h.Write([]byte(msg))
	r := new(big.Int).SetBytes(sig[0:curvelen])
	s := new(big.Int).SetBytes(sig[curvelen:])
	pub := ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	if !ecdsa.Verify(&pub, h.Sum(nil), r, s) {
		t.Error("ecdsa signature verify against JWK-derived key failed")
	}
}

func dpopProofAcceptable(dpop string, now time.Time) bool {
	parts := strings.Split(dpop, ".")
	if len(parts) != 3 {
		return false
	}

	hb, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	var hdr dpopHeader
	if json.Unmarshal(hb, &hdr) != nil || hdr.Jwk == nil {
		return false
	}

	xb, err := base64.RawURLEncoding.DecodeString(hdr.Jwk.X)
	if err != nil || len(xb) != curvelen {
		return false
	}
	yb, err := base64.RawURLEncoding.DecodeString(hdr.Jwk.Y)
	if err != nil || len(yb) != curvelen {
		return false
	}

	pb, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	var payload dpopPayload
	if json.Unmarshal(pb, &payload) != nil {
		return false
	}
	if payload.Htm != "POST" || !strings.Contains(payload.Htu, "/v1/token") {
		return false
	}
	if payload.Iat != now.Unix() {
		return false
	}

	x := new(big.Int).SetBytes(xb)
	y := new(big.Int).SetBytes(yb)
	curve := elliptic.P256()
	if !curve.IsOnCurve(x, y) {
		return false
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || len(sig) != curvelen*2 {
		return false
	}

	msg := fmt.Sprintf("%s.%s", parts[0], parts[1])
	h := sha256.New()
	h.Write([]byte(msg))
	r := new(big.Int).SetBytes(sig[0:curvelen])
	s := new(big.Int).SetBytes(sig[curvelen:])
	pub := ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	return ecdsa.Verify(&pub, h.Sum(nil), r, s)
}
