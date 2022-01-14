package v4a

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/v4a/internal/crypto"
	"github.com/aws/smithy-go/logging"
	"github.com/google/go-cmp/cmp"
)

const (
	accessKey = "AKISORANDOMAASORANDOM"
	secretKey = "q+jcrXGc+0zWN6uzclKVhvMmUsIfRPa4rlRandom"
)

func TestDeriveECDSAKeyPairFromSecret(t *testing.T) {
	privateKey, err := deriveKeyFromAccessKeyPair(accessKey, secretKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedX := func() *big.Int {
		t.Helper()
		b, ok := new(big.Int).SetString("15D242CEEBF8D8169FD6A8B5A746C41140414C3B07579038DA06AF89190FFFCB", 16)
		if !ok {
			t.Fatalf("failed to parse big integer")
		}
		return b
	}()
	expectedY := func() *big.Int {
		t.Helper()
		b, ok := new(big.Int).SetString("515242CEDD82E94799482E4C0514B505AFCCF2C0C98D6A553BF539F424C5EC0", 16)
		if !ok {
			t.Fatalf("failed to parse big integer")
		}
		return b
	}()

	if privateKey.X.Cmp(expectedX) != 0 {
		t.Errorf("expected % X, got % X", expectedX, privateKey.X)
	}
	if privateKey.Y.Cmp(expectedY) != 0 {
		t.Errorf("expected % X, got % X", expectedY, privateKey.Y)
	}
}

func TestSignHTTP(t *testing.T) {
	req := buildRequest("dynamodb", "us-east-1")

	signer, credProvider := buildSigner(t, true)

	key, err := credProvider.RetrievePrivateKey(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	err = signer.SignHTTP(context.Background(), key, req, EmptyStringSHA256, "dynamodb", []string{"us-east-1"}, time.Unix(0, 0))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedDate := "19700101T000000Z"
	expectedAlg := "AWS4-ECDSA-P256-SHA256"
	expectedCredential := "AKISORANDOMAASORANDOM/19700101/dynamodb/aws4_request"
	expectedSignedHeaders := "content-length;content-type;host;x-amz-date;x-amz-meta-other-header;x-amz-meta-other-header_with_underscore;x-amz-region-set;x-amz-security-token;x-amz-target"
	expectedStrToSignHash := "4ba7d0482cf4d5450cefdc067a00de1a4a715e444856fa3e1d85c35fb34d9730"

	q := req.Header

	validateAuthorization(t, q.Get("Authorization"), expectedAlg, expectedCredential, expectedSignedHeaders, expectedStrToSignHash)

	if e, a := expectedDate, q.Get("X-Amz-Date"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestSignHTTP_NoSessionToken(t *testing.T) {
	req := buildRequest("dynamodb", "us-east-1")

	signer, credProvider := buildSigner(t, false)

	key, err := credProvider.RetrievePrivateKey(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	err = signer.SignHTTP(context.Background(), key, req, EmptyStringSHA256, "dynamodb", []string{"us-east-1"}, time.Unix(0, 0))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedAlg := "AWS4-ECDSA-P256-SHA256"
	expectedCredential := "AKISORANDOMAASORANDOM/19700101/dynamodb/aws4_request"
	expectedSignedHeaders := "content-length;content-type;host;x-amz-date;x-amz-meta-other-header;x-amz-meta-other-header_with_underscore;x-amz-region-set;x-amz-target"
	expectedStrToSignHash := "1aeefb422ae6aa0de7aec829da813e55cff35553cac212dffd5f9474c71e47ee"

	q := req.Header

	validateAuthorization(t, q.Get("Authorization"), expectedAlg, expectedCredential, expectedSignedHeaders, expectedStrToSignHash)
}

func TestPresignHTTP(t *testing.T) {
	req := buildRequest("dynamodb", "us-east-1")

	signer, credProvider := buildSigner(t, false)

	key, err := credProvider.RetrievePrivateKey(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	query := req.URL.Query()
	query.Set("X-Amz-Expires", "18000")
	req.URL.RawQuery = query.Encode()

	signedURL, _, err := signer.PresignHTTP(context.Background(), key, req, EmptyStringSHA256, "dynamodb", []string{"us-east-1"}, time.Unix(0, 0))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedDate := "19700101T000000Z"
	expectedAlg := "AWS4-ECDSA-P256-SHA256"
	expectedHeaders := "content-length;content-type;host;x-amz-meta-other-header;x-amz-meta-other-header_with_underscore"
	expectedCredential := "AKISORANDOMAASORANDOM/19700101/dynamodb/aws4_request"
	expectedStrToSignHash := "d7ffbd2fab644384c056957e6ac38de4ae68246764b5f5df171b3824153b6397"
	expectedTarget := "prefix.Operation"

	signedReq, err := url.Parse(signedURL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	q := signedReq.Query()

	validateSignature(t, expectedStrToSignHash, q.Get("X-Amz-Signature"))

	if e, a := expectedAlg, q.Get("X-Amz-Algorithm"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedCredential, q.Get("X-Amz-Credential"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedHeaders, q.Get("X-Amz-SignedHeaders"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedDate, q.Get("X-Amz-Date"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if a := q.Get("X-Amz-Meta-Other-Header"); len(a) != 0 {
		t.Errorf("expect %v to be empty", a)
	}
	if e, a := expectedTarget, q.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "us-east-1", q.Get("X-Amz-Region-Set"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}

func TestPresignHTTP_BodyWithArrayRequest(t *testing.T) {
	req := buildRequest("dynamodb", "us-east-1")
	req.URL.RawQuery = "Foo=z&Foo=o&Foo=m&Foo=a"

	signer, credProvider := buildSigner(t, true)

	key, err := credProvider.RetrievePrivateKey(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	query := req.URL.Query()
	query.Set("X-Amz-Expires", "300")
	req.URL.RawQuery = query.Encode()

	signedURI, _, err := signer.PresignHTTP(context.Background(), key, req, EmptyStringSHA256, "dynamodb", []string{"us-east-1"}, time.Unix(0, 0))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	signedReq, err := url.Parse(signedURI)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedAlg := "AWS4-ECDSA-P256-SHA256"
	expectedDate := "19700101T000000Z"
	expectedHeaders := "content-length;content-type;host;x-amz-meta-other-header;x-amz-meta-other-header_with_underscore"
	expectedStrToSignHash := "acff64fd3689be96259d4112c3742ff79f4da0d813bc58a285dc1c4449760bec"
	expectedCred := "AKISORANDOMAASORANDOM/19700101/dynamodb/aws4_request"
	expectedTarget := "prefix.Operation"

	q := signedReq.Query()

	validateSignature(t, expectedStrToSignHash, q.Get("X-Amz-Signature"))

	if e, a := expectedAlg, q.Get("X-Amz-Algorithm"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedCred, q.Get("X-Amz-Credential"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedHeaders, q.Get("X-Amz-SignedHeaders"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := expectedDate, q.Get("X-Amz-Date"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if a := q.Get("X-Amz-Meta-Other-Header"); len(a) != 0 {
		t.Errorf("expect %v to be empty, was not", a)
	}
	if e, a := expectedTarget, q.Get("X-Amz-Target"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "us-east-1", q.Get("X-Amz-Region-Set"); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}
func TestSign_buildCanonicalHeaders(t *testing.T) {
	serviceName := "mockAPI"
	region := "mock-region"
	endpoint := "https://" + serviceName + "." + region + ".amazonaws.com"

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		t.Fatalf("failed to create request, %v", err)
	}

	req.Header.Set("FooInnerSpace", "   inner      space    ")
	req.Header.Set("FooLeadingSpace", "    leading-space")
	req.Header.Add("FooMultipleSpace", "no-space")
	req.Header.Add("FooMultipleSpace", "\ttab-space")
	req.Header.Add("FooMultipleSpace", "trailing-space    ")
	req.Header.Set("FooNoSpace", "no-space")
	req.Header.Set("FooTabSpace", "\ttab-space\t")
	req.Header.Set("FooTrailingSpace", "trailing-space    ")
	req.Header.Set("FooWrappedSpace", "   wrapped-space    ")

	credProvider := &SymmetricCredentialAdaptor{
		SymmetricProvider: staticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			},
		},
	}
	key, err := credProvider.RetrievePrivateKey(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	ctx := &httpSigner{
		Request:     req,
		ServiceName: serviceName,
		RegionSet:   []string{region},
		Credentials: key,
		Time:        time.Date(2021, 10, 20, 12, 42, 0, 0, time.UTC),
	}

	build, err := ctx.Build()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectCanonicalString := strings.Join([]string{
		`POST`,
		`/`,
		``,
		`fooinnerspace:inner space`,
		`fooleadingspace:leading-space`,
		`foomultiplespace:no-space,tab-space,trailing-space`,
		`foonospace:no-space`,
		`footabspace:tab-space`,
		`footrailingspace:trailing-space`,
		`foowrappedspace:wrapped-space`,
		`host:mockAPI.mock-region.amazonaws.com`,
		`x-amz-date:20211020T124200Z`,
		`x-amz-region-set:mock-region`,
		``,
		`fooinnerspace;fooleadingspace;foomultiplespace;foonospace;footabspace;footrailingspace;foowrappedspace;host;x-amz-date;x-amz-region-set`,
		``,
	}, "\n")
	if diff := cmp.Diff(expectCanonicalString, build.CanonicalString); diff != "" {
		t.Errorf("expect match, got\n%s", diff)
	}
}

func validateAuthorization(t *testing.T, authorization, expectedAlg, expectedCredential, expectedSignedHeaders, expectedStrToSignHash string) {
	t.Helper()
	split := strings.SplitN(authorization, " ", 2)

	if len(split) != 2 {
		t.Fatal("unexpected authorization header format")
	}

	if e, a := split[0], expectedAlg; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	keyValues := strings.Split(split[1], ", ")
	seen := make(map[string]string)

	for _, kv := range keyValues {
		idx := strings.Index(kv, "=")
		if idx == -1 {
			continue
		}
		key, value := kv[:idx], kv[idx+1:]
		seen[key] = value
	}

	if a, ok := seen["Credential"]; ok {
		if expectedCredential != a {
			t.Errorf("expected credential %v, got %v", expectedCredential, a)
		}
	} else {
		t.Errorf("Credential not found in authorization string")
	}

	if a, ok := seen["SignedHeaders"]; ok {
		if expectedSignedHeaders != a {
			t.Errorf("expected signed headers %v, got %v", expectedSignedHeaders, a)
		}
	} else {
		t.Errorf("SignedHeaders not found in authorization string")
	}

	if a, ok := seen["Signature"]; ok {
		validateSignature(t, expectedStrToSignHash, a)
	} else {
		t.Errorf("signature not found in authorization string")
	}
}

func validateSignature(t *testing.T, expectedHash, signature string) {
	t.Helper()
	pair, err := deriveKeyFromAccessKeyPair(accessKey, secretKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	hash, _ := hex.DecodeString(expectedHash)
	sig, _ := hex.DecodeString(signature)

	ok, err := crypto.VerifySignature(&pair.PublicKey, hash, sig)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Errorf("failed to verify signing singature")
	}
}

func buildRequest(serviceName, region string) *http.Request {
	endpoint := "https://" + serviceName + "." + region + ".amazonaws.com"
	req, _ := http.NewRequest("POST", endpoint, nil)
	req.URL.Opaque = "//example.org/bucket/key-._~,!@%23$%25^&*()"
	req.Header.Set("X-Amz-Target", "prefix.Operation")
	req.Header.Set("Content-Type", "application/x-amz-json-1.0")

	req.Header.Set("Content-Length", strconv.Itoa(1024))

	req.Header.Set("X-Amz-Meta-Other-Header", "some-value=!@#$%^&* (+)")
	req.Header.Add("X-Amz-Meta-Other-Header_With_Underscore", "some-value=!@#$%^&* (+)")
	req.Header.Add("X-amz-Meta-Other-Header_With_Underscore", "some-value=!@#$%^&* (+)")
	return req
}

func buildSigner(t *testing.T, withToken bool) (*Signer, CredentialsProvider) {
	creds := aws.Credentials{
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
	}

	if withToken {
		creds.SessionToken = "TOKEN"
	}

	return NewSigner(func(options *SignerOptions) {
			options.Logger = loggerFunc(func(format string, v ...interface{}) {
				t.Logf(format, v...)
			})
		}), &SymmetricCredentialAdaptor{
			SymmetricProvider: staticCredentialsProvider{
				Value: creds,
			},
		}
}

type loggerFunc func(format string, v ...interface{})

func (l loggerFunc) Logf(_ logging.Classification, format string, v ...interface{}) {
	l(format, v...)
}

type staticCredentialsProvider struct {
	Value aws.Credentials
}

func (s staticCredentialsProvider) Retrieve(_ context.Context) (aws.Credentials, error) {
	v := s.Value
	if v.AccessKeyID == "" || v.SecretAccessKey == "" {
		return aws.Credentials{
			Source: "Source Name",
		}, fmt.Errorf("static credentials are empty")
	}

	if len(v.Source) == 0 {
		v.Source = "Source Name"
	}

	return v, nil
}
