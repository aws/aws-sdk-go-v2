// Package v4 implements signing for AWS V4 signer
//
// Provides request signing for request that need to be signed with
// AWS V4 Signatures.
//
// Standalone Signer
//
// Generally using the signer outside of the SDK should not require any additional
//  The signer does this by taking advantageof the URL.EscapedPath method. If your request URI requires
// additional escaping you many need to use the URL.Opaque to define what the raw URI should be sent
// to the service as.
//
// The signer will first check the URL.Opaque field, and use its value if set.
// The signer does require the URL.Opaque field to be set in the form of:
//
//     "//<hostname>/<path>"
//
//     // e.g.
//     "//example.com/some/path"
//
// The leading "//" and hostname are required or the URL.Opaque escaping will
// not work correctly.
//
// If URL.Opaque is not set the signer will fallback to the URL.EscapedPath()
// method and using the returned value.
//
// AWS v4 signature validation requires that the canonical string's URI path
// element must be the URI escaped form of the HTTP request's path.
// http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
//
// The Go HTTP client will perform escaping automatically on the request. Some
// of these escaping may cause signature validation errors because the HTTP
// request differs from the URI path or query that the signature was generated.
// https://golang.org/pkg/net/url/#URL.EscapedPath
//
// Because of this, it is recommended that when using the signer outside of the
// SDK that explicitly escaping the request prior to being signed is preferable,
// and will help prevent signature validation errors. This can be done by setting
// the URL.Opaque or URL.RawPath. The SDK will use URL.Opaque first and then
// call URL.EscapedPath() if Opaque is not set.
//
// Test `TestStandaloneSign` provides a complete example of using the signer
// outside of the SDK and pre-escaping the URI path.
package v4

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4Internal "github.com/aws/aws-sdk-go-v2/aws/signer/internal/v4"
	"github.com/awslabs/smithy-go/httpbinding"
)

const (
	signingAlgorithm = "AWS4-HMAC-SHA256"
)

// HTTPSigner is an interface to a SigV4 signer that can sign HTTP requests
type HTTPSigner interface {
	SignHTTP(ctx context.Context, r *http.Request, payloadHash string, service string, region string, signingTime time.Time) error
}

// Signer applies AWS v4 signing to given request. Use this to sign requests
// that need to be signed with AWS V4 Signatures.
type Signer struct {
	// The authentication credentials the request will be signed against.
	// This value must be set to sign requests.
	Credentials aws.CredentialsProvider

	// Sets the log level the signer should use when reporting information to
	// the logger. If the logger is nil nothing will be logged. See
	// aws.LogLevel for more information on available logging levels
	//
	// By default nothing will be logged.
	Debug aws.LogLevel

	// The logger loging information will be written to. If there the logger
	// is nil, nothing will be logged.
	Logger aws.Logger

	// Disables the Signer's moving HTTP header key/value pairs from the HTTP
	// request header to the request's query string. This is most commonly used
	// with pre-signed requests preventing headers from being added to the
	// request's query string.
	DisableHeaderHoisting bool

	// Disables the automatic escaping of the URI path of the request for the
	// siganture's canonical string's path. For services that do not need additional
	// escaping then use this to disable the signer escaping the path.
	//
	// S3 is an example of a service that does not need additional escaping.
	//
	// http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
	DisableURIPathEscaping bool
}

// NewSigner returns a Signer pointer configured with the credentials and optional
// option values provided. If not options are provided the Signer will use its
// default configuration.
func NewSigner(credsProvider aws.CredentialsProvider, options ...func(*Signer)) *Signer {
	v4 := &Signer{
		Credentials: credsProvider,
	}

	for _, option := range options {
		option(v4)
	}

	return v4
}

type httpSigner struct {
	Request     *http.Request
	ServiceName string
	Region      string
	Time        time.Time
	ExpireTime  time.Duration
	Credentials aws.Credentials
	IsPreSign   bool

	// PayloadHash is the hex encoded SHA-256 hash of the request payload
	// If len(PayloadHash) == 0 the signer will attempt to send the request
	// as an unsigned payload. Note: Unsigned payloads only work for a subset of services.
	PayloadHash string

	DisableHeaderHoisting  bool
	DisableURIPathEscaping bool
}

func (s *httpSigner) Build() (signedRequest, error) {
	req := s.Request

	query := req.URL.Query()
	headers := req.Header

	s.setRequiredSigningFields(headers, query)

	// Sort Each Query Key's Values
	for key := range query {
		sort.Strings(query[key])
	}

	v4Internal.SanitizeHostForHeader(req)

	credentialScope := s.buildCredentialScope()
	credentialStr := s.Credentials.AccessKeyID + "/" + credentialScope
	if s.IsPreSign {
		query.Set(v4Internal.AmzCredentialKey, credentialStr)
	}

	unsignedHeaders := headers
	if s.IsPreSign && !s.DisableHeaderHoisting {
		urlValues := url.Values{}
		urlValues, unsignedHeaders = buildQuery(v4Internal.AllowedQueryHoisting, unsignedHeaders)
		for k := range urlValues {
			query[k] = urlValues[k]
		}
	}

	host := req.URL.Host
	if len(req.Host) > 0 {
		host = req.Host
	}

	signedHeaders, signedHeadersStr, canonicalHeaderStr := s.buildCanonicalHeaders(host, v4Internal.IgnoredHeaders, unsignedHeaders)

	if s.IsPreSign {
		query.Set(v4Internal.AmzSignedHeadersKey, signedHeadersStr)
	}

	rawQuery := strings.Replace(query.Encode(), "+", "%20", -1)

	canonicalURI := v4Internal.GetURIPath(req.URL)
	if !s.DisableURIPathEscaping {
		canonicalURI = httpbinding.EscapePath(canonicalURI, false)
	}

	canonicalString := s.buildCanonicalString(
		req.Method,
		canonicalURI,
		rawQuery,
		signedHeadersStr,
		canonicalHeaderStr,
	)

	strToSign := s.buildStringToSign(credentialScope, canonicalString)
	signingSignature := s.buildSignature(strToSign)

	if s.IsPreSign {
		rawQuery += "&X-Amz-Signature=" + signingSignature
	} else {
		parts := []string{
			"Credential=" + credentialStr,
			"SignedHeaders=" + signedHeadersStr,
			"Signature=" + signingSignature,
		}
		headers.Set("Authorization", signingAlgorithm+" "+strings.Join(parts, ", "))
	}

	req.URL.RawQuery = rawQuery

	return signedRequest{
		Request:         req,
		SignedHeaders:   signedHeaders,
		CanonicalString: canonicalString,
		StringToSign:    strToSign,
		PreSigned:       s.IsPreSign,
	}, nil
}

// SignHTTP signs AWS v4 requests with the provided payload hash, service name, region the
// request is made to, and time the request is signed at. The signTime allows
// you to specify that a request is signed for the future, and cannot be
// used until then.
//
// Sign differs from Presign in that it will sign the request using HTTP
// header values. This type of signing is intended for http.Request values that
// will not be shared, or are shared in a way the header values on the request
// will not be lost.
//
// The passed in request will be modified in place.
func (v4 Signer) SignHTTP(ctx context.Context, r *http.Request, payloadHash string, service string, region string, signingTime time.Time) error {
	credentials, err := v4.Credentials.Retrieve(ctx)
	if err != nil {
		return err
	}

	signer := &httpSigner{
		Request:                r,
		PayloadHash:            payloadHash,
		ServiceName:            service,
		Region:                 region,
		Credentials:            credentials,
		Time:                   signingTime.UTC(),
		DisableHeaderHoisting:  v4.DisableHeaderHoisting,
		DisableURIPathEscaping: v4.DisableURIPathEscaping,
	}

	signedRequest, err := signer.Build()
	if err != nil {
		return err
	}

	v4.logHTTPSigningInfo(signedRequest)

	return nil
}

// PresignHTTP signs AWS v4 requests with the payload hash, service name, region
// the request is made to, and time the request is signed at. The signTime
// allows you to specify that a request is signed for the future, and cannot
// be used until then.
//
// Returns the signed URL and the map of HTTP headers that were included in the signature or an
// error if signing the request failed. For presigned requests these headers
// and their values must be included on the HTTP request when it is made. This
// is helpful to know what header values need to be shared with the party the
// presigned request will be distributed to.
//
// PresignHTTP differs from SignHTTP in that it will sign the request using query string
// instead of header values. This allows you to share the Presigned Request's
// URL with third parties, or distribute it throughout your system with minimal
// dependencies.
//
// PresignHTTP also takes an exp value which is the duration the
// signed request will be valid after the signing time. This is allows you to
// set when the request will expire.
//
// This method does not modify the provided request.
func (v4 *Signer) PresignHTTP(ctx context.Context, r *http.Request, payloadHash string, service string, region string, expireTime time.Duration, signingTime time.Time) (signedURI string, signedHeaders http.Header, err error) {
	credentials, err := v4.Credentials.Retrieve(ctx)
	if err != nil {
		return "", nil, err
	}

	signer := &httpSigner{
		Request:                r.Clone(r.Context()),
		PayloadHash:            payloadHash,
		ServiceName:            service,
		Region:                 region,
		Credentials:            credentials,
		Time:                   signingTime.UTC(),
		IsPreSign:              true,
		ExpireTime:             expireTime,
		DisableHeaderHoisting:  v4.DisableHeaderHoisting,
		DisableURIPathEscaping: v4.DisableURIPathEscaping,
	}

	signedRequest, err := signer.Build()
	if err != nil {
		return "", nil, err
	}

	v4.logHTTPSigningInfo(signedRequest)

	return signedRequest.Request.URL.String(), signedRequest.SignedHeaders, nil
}

const logSignInfoMsg = `DEBUG: Request Signature:
---[ CANONICAL STRING  ]-----------------------------
%s
---[ STRING TO SIGN ]--------------------------------
%s%s
-----------------------------------------------------`
const logSignedURLMsg = `
---[ SIGNED URL ]------------------------------------
%s`

func (v4 Signer) logHTTPSigningInfo(r signedRequest) {
	if !v4.Debug.Matches(aws.LogDebugWithSigning) || v4.Logger == nil {
		return
	}

	signedURLMsg := ""
	if r.PreSigned {
		signedURLMsg = fmt.Sprintf(logSignedURLMsg, r.Request.URL.String())
	}
	msg := fmt.Sprintf(logSignInfoMsg, r.CanonicalString, r.StringToSign, signedURLMsg)
	v4.Logger.Log(msg)
}

func (s *httpSigner) buildCredentialScope() string {
	return strings.Join([]string{
		s.Time.Format(v4Internal.ShortTimeFormat),
		s.Region,
		s.ServiceName,
		"aws4_request",
	}, "/")
}

func buildQuery(r v4Internal.Rule, header http.Header) (url.Values, http.Header) {
	query := url.Values{}
	unsignedHeaders := http.Header{}
	for k, h := range header {
		if r.IsValid(k) {
			query[k] = h
		} else {
			unsignedHeaders[k] = h
		}
	}

	return query, unsignedHeaders
}

func (s *httpSigner) buildCanonicalHeaders(host string, rule v4Internal.Rule, header http.Header) (signed http.Header, signedHeaders, canonicalHeaders string) {
	signed = make(http.Header)

	var headers []string
	headers = append(headers, "host")
	signed["host"] = append(signed["host"], host)

	for k, v := range header {
		canonicalKey := http.CanonicalHeaderKey(k)
		if !rule.IsValid(canonicalKey) {
			continue // ignored header
		}

		lowerCaseKey := strings.ToLower(k)
		if _, ok := signed[lowerCaseKey]; ok {
			// include additional values
			signed[lowerCaseKey] = append(signed[lowerCaseKey], v...)
			continue
		}

		headers = append(headers, lowerCaseKey)
		signed[lowerCaseKey] = v
	}
	sort.Strings(headers)

	signedHeaders = strings.Join(headers, ";")

	headerValues := make([]string, len(headers))
	for i, k := range headers {
		if k == "host" {
			headerValues[i] = "host:" + host
		} else {
			headerValues[i] = k + ":" + strings.Join(signed[k], ",")
		}
	}
	v4Internal.StripExcessSpaces(headerValues)
	canonicalHeaders = strings.Join(headerValues, "\n")

	return signed, signedHeaders, canonicalHeaders
}

func (s *httpSigner) buildCanonicalString(method, uri, query, signedHeaders, canonicalHeaders string) string {
	return strings.Join([]string{
		method,
		uri,
		query,
		canonicalHeaders + "\n",
		signedHeaders,
		s.PayloadHash,
	}, "\n")
}

func (s *httpSigner) buildStringToSign(credentialScope, canonicalRequestString string) string {
	return strings.Join([]string{
		signingAlgorithm,
		s.Time.Format(v4Internal.TimeFormat),
		credentialScope,
		hex.EncodeToString(makeHash(sha256.New(), []byte(canonicalRequestString))),
	}, "\n")
}

func makeHash(hash hash.Hash, b []byte) []byte {
	hash.Reset()
	hash.Write(b)
	return hash.Sum(nil)
}

func (s *httpSigner) buildSignature(strToSign string) string {
	secret := s.Credentials.SecretAccessKey
	date := makeHmacSha256([]byte("AWS4"+secret), []byte(s.Time.Format(v4Internal.ShortTimeFormat)))
	region := makeHmacSha256(date, []byte(s.Region))
	service := makeHmacSha256(region, []byte(s.ServiceName))
	credentials := makeHmacSha256(service, []byte("aws4_request"))
	signature := makeHmacSha256(credentials, []byte(strToSign))
	return hex.EncodeToString(signature)
}

func (s *httpSigner) setRequiredSigningFields(headers http.Header, query url.Values) {
	amzDate := s.Time.Format(v4Internal.TimeFormat)

	if s.IsPreSign {
		query.Set(v4Internal.AmzAlgorithmKey, signingAlgorithm)
		if sessionToken := s.Credentials.SessionToken; len(sessionToken) > 0 {
			query.Set("X-Amz-Security-Token", sessionToken)
		}

		duration := int64(s.ExpireTime / time.Second)
		query.Set(v4Internal.AmzDateKey, amzDate)
		query.Set(v4Internal.AmzExpiresKey, strconv.FormatInt(duration, 10))
		return
	}

	headers.Set(v4Internal.AmzDateKey, amzDate)

	if len(s.Credentials.SessionToken) > 0 {
		headers.Set(v4Internal.AmzSecurityTokenKey, s.Credentials.SessionToken)
	}
}

func makeHmacSha256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

type signedRequest struct {
	Request         *http.Request
	SignedHeaders   http.Header
	CanonicalString string
	StringToSign    string
	PreSigned       bool
}
