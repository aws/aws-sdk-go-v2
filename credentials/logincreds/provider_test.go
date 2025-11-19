package logincreds

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/signin"
	"github.com/aws/aws-sdk-go-v2/service/signin/types"
)

func mockNowTime(t time.Time) func() {
	sdk.NowTime = func() time.Time { return t }
	return func() { sdk.NowTime = time.Now }
}

type writer struct {
	p []byte
}

func (w *writer) Write(p []byte) (int, error) {
	w.p = append(w.p, p...)
	return len(p), nil
}

func (w *writer) Close() error {
	return nil
}

type tokenAPIClient struct {
	in  *signin.CreateOAuth2TokenInput
	out *signin.CreateOAuth2TokenOutput
	err error
}

func mockTokenAPIClient(out *signin.CreateOAuth2TokenOutput, err error) *tokenAPIClient {
	return &tokenAPIClient{out: out, err: err}
}

func (m *tokenAPIClient) CreateOAuth2Token(ctx context.Context, in *signin.CreateOAuth2TokenInput, opts ...func(*signin.Options)) (*signin.CreateOAuth2TokenOutput, error) {
	m.in = in
	return m.out, m.err
}

func mockOpenFile(t *loginToken, terr error) func() {
	orig := openFile
	openFile = func(name string) (io.ReadCloser, error) {
		if terr != nil {
			return nil, terr
		}

		j, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		return io.NopCloser(bytes.NewReader(j)), nil
	}
	return func() { openFile = orig }
}

func mockCreateFile(terr error) (*writer, func()) {
	w := &writer{}

	orig := createFile
	createFile = func(name string) (io.WriteCloser, error) {
		if terr != nil {
			return nil, terr
		}

		return w, nil
	}
	return w, func() { createFile = orig }
}

func mockKey() (*ecdsa.PrivateKey, string) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), cryptorand.Reader)
	if err != nil {
		panic(err)
	}
	der, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		panic(err)
	}

	return key, string(pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: der,
	}))
}

func mockToken() (*loginToken, *ecdsa.PrivateKey) {
	key, kpem := mockKey()
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
		DPOPKey:       kpem,
	}, key
}

// Success - Valid credentials are returned immediately
func TestRetrieve_OK_NotExpired(t *testing.T) {
	token, _ := mockToken()
	restoreOpen := mockOpenFile(token, nil)
	restoreNowTime := mockNowTime(time.Unix(60, 0).UTC())
	defer restoreOpen()
	defer restoreNowTime()

	token.AccessToken.ExpiresAt = time.Unix(120, 0).UTC()
	svc := mockTokenAPIClient(nil, nil) // won't be called

	p := New(svc, "mocktokenpath")
	creds, err := p.Retrieve(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if svc.in != nil {
		t.Fatal("CreateOAuth2Token shouldn't be called")
	}

	// should match what we "loaded"
	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("akid: %q != %q", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("secret: %q != %q", e, a)
	}
	if e, a := "SESSION", creds.SessionToken; e != a {
		t.Errorf("session: %q != %q", e, a)
	}
	if e, a := "ACCOUNTID", creds.AccountID; e != a {
		t.Errorf("account id: %q != %q", e, a)
	}
	if e, a := time.Unix(120, 0).UTC(), creds.Expires; e != a {
		t.Errorf("expires: %v != %v", e, a)
	}
}

// Failure - No cache file
func TestRetrieve_Failure_NoCacheFile(t *testing.T) {
	restoreOpen := mockOpenFile(nil, os.ErrNotExist)
	defer restoreOpen()

	svc := mockTokenAPIClient(nil, nil)
	p := New(svc, "mocktokenpath")
	_, err := p.Retrieve(context.Background())
	if err == nil {
		t.Fatal("expect err, got none")
	}
	if !strings.Contains(err.Error(), "token file not found, please reauthenticate") {
		t.Errorf("unexpected error: %v", err)
	}
}

// Failure - Missing accessToken
func TestRetrieve_Failure_MissingAccessToken(t *testing.T) {
	token, _ := mockToken()
	token.AccessToken = nil

	restoreOpen := mockOpenFile(token, nil)
	defer restoreOpen()

	svc := mockTokenAPIClient(nil, nil)
	p := New(svc, "mocktokenpath")
	_, err := p.Retrieve(context.Background())
	if err == nil {
		t.Fatal("expect err, got none")
	}
	if !strings.Contains(err.Error(), "validate login token") {
		t.Errorf("unexpected error: %v", err)
	}
}

// Failure - Missing refreshToken
func TestRetrieve_Failure_MissingRefreshToken(t *testing.T) {
	token, _ := mockToken()
	token.RefreshToken = ""

	restoreOpen := mockOpenFile(token, nil)
	defer restoreOpen()

	svc := mockTokenAPIClient(nil, nil)
	p := New(svc, "mocktokenpath")
	_, err := p.Retrieve(context.Background())
	if err == nil {
		t.Fatal("expect err, got none")
	}
	if !strings.Contains(err.Error(), "validate login token") {
		t.Errorf("unexpected error: %v", err)
	}
}

// Failure - Missing clientId in cache
func TestRetrieve_Failure_MissingClientID(t *testing.T) {
	token, _ := mockToken()
	token.ClientID = ""

	restoreOpen := mockOpenFile(token, nil)
	defer restoreOpen()

	svc := mockTokenAPIClient(nil, nil)
	p := New(svc, "mocktokenpath")
	_, err := p.Retrieve(context.Background())
	if err == nil {
		t.Fatal("expect err, got none")
	}
	if !strings.Contains(err.Error(), "validate login token") {
		t.Errorf("unexpected error: %v", err)
	}
}

// Failure - Missing dpopKey
func TestRetrieve_Failure_MissingDPOPKey(t *testing.T) {
	token, _ := mockToken()
	token.DPOPKey = ""

	restoreOpen := mockOpenFile(token, nil)
	defer restoreOpen()

	svc := mockTokenAPIClient(nil, nil)
	p := New(svc, "mocktokenpath")
	_, err := p.Retrieve(context.Background())
	if err == nil {
		t.Fatal("expect err, got none")
	}
	if !strings.Contains(err.Error(), "validate login token") {
		t.Errorf("unexpected error: %v", err)
	}
}

// Success - Expired token triggers successful refresh
func TestRetrieve_OK_Refresh(t *testing.T) {
	token, _ := mockToken()
	restoreOpen := mockOpenFile(token, nil)
	written, restoreCreate := mockCreateFile(nil)
	restoreNowTime := mockNowTime(time.Unix(60, 0).UTC())
	defer restoreOpen()
	defer restoreCreate()
	defer restoreNowTime()

	svc := mockTokenAPIClient(&signin.CreateOAuth2TokenOutput{
		TokenOutput: &types.CreateOAuth2TokenResponseBody{
			AccessToken: &types.AccessToken{
				AccessKeyId:     aws.String("NEW_AKID"),
				SecretAccessKey: aws.String("NEW_SECRET"),
				SessionToken:    aws.String("NEW_SESSION"),
			},
			ExpiresIn:    aws.Int32(900), // actual service returns 15-min creds
			RefreshToken: aws.String("NewRefreshToken"),
			TokenType:    aws.String("NewTokenType"),
			IdToken:      aws.String("NewIdToken"),
		},
	}, nil)

	p := New(svc, "mocktokenpath")
	creds, err := p.Retrieve(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	var savedToken *loginToken
	if err := json.Unmarshal(written.p, &savedToken); err != nil {
		t.Fatal(err)
	}

	if e, a := "NEW_AKID", creds.AccessKeyID; e != a {
		t.Errorf("akid: %q != %q", e, a)
	}
	if e, a := "NEW_SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("secret: %q != %q", e, a)
	}
	if e, a := "NEW_SESSION", creds.SessionToken; e != a {
		t.Errorf("session: %q != %q", e, a)
	}
	if e, a := "ACCOUNTID", creds.AccountID; e != a {
		t.Errorf("account id: %q != %q", e, a)
	}
	// we mocked time.Now() to return 60 and creds expire in 900 seconds
	if e, a := time.Unix(960, 0).UTC(), creds.Expires; e != a {
		t.Errorf("expires: %v != %v", e, a)
	}

	if e, a := "NEW_AKID", savedToken.AccessToken.AccessKeyID; e != a {
		t.Errorf("akid: %q != %q", e, a)
	}
	if e, a := "NEW_SECRET", savedToken.AccessToken.SecretAccessKey; e != a {
		t.Errorf("secret: %q != %q", e, a)
	}
	if e, a := "NEW_SESSION", savedToken.AccessToken.SessionToken; e != a {
		t.Errorf("session: %q != %q", e, a)
	}
	if e, a := time.Unix(960, 0).UTC(), savedToken.AccessToken.ExpiresAt; e != a {
		t.Errorf("expires: %v != %v", e, a)
	}
	if e, a := "NewRefreshToken", savedToken.RefreshToken; e != a {
		t.Errorf("saved refresh token: %q != %q", e, a)
	}
}

// Failure - Expired token triggers failed refresh
func TestRetrieve_Failure_Refresh(t *testing.T) {
	token, _ := mockToken()
	restoreOpen := mockOpenFile(token, nil)
	restoreNowTime := mockNowTime(time.Unix(60, 0).UTC())
	defer restoreOpen()
	defer restoreNowTime()

	svc := mockTokenAPIClient(nil, &types.AccessDeniedException{
		Error_: types.OAuth2ErrorCodeTokenExpired,
	})

	p := New(svc, "mocktokenpath")
	_, err := p.Retrieve(context.Background())
	if err == nil {
		t.Fatal("expect err, got none")
	}
	if !strings.Contains(err.Error(), "create oauth2 token") {
		t.Errorf("unexpected error: %v", err)
	}
}
