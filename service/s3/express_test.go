package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
)

type mockCreateSession struct {
	wg sync.WaitGroup

	calls []mockCreateSessionCall
	times int

	// AKIDs of the sigv4 credentials each CreateSession call was made with,
	// as set by the per-call options
	akids []string
}

type mockCreateSessionCall struct {
	output *CreateSessionOutput
	err    error
}

func (m *mockCreateSession) expectCalled(t *testing.T, times int) {
	if m.times != times {
		t.Errorf("expected %d calls to CreateSession, got %d", times, m.times)
	}
}

func (m *mockCreateSession) CreateSession(ctx context.Context, _ *CreateSessionInput, optFns ...func(*Options)) (*CreateSessionOutput, error) {
	defer m.wg.Done()
	var opts Options
	for _, fn := range optFns {
		fn(&opts)
	}
	akid := ""
	if opts.Credentials != nil {
		if creds, err := opts.Credentials.Retrieve(ctx); err == nil {
			akid = creds.AccessKeyID
		}
	}
	m.akids = append(m.akids, akid)

	o := m.calls[m.times]
	m.times++
	return o.output, o.err
}

type mockCreds struct {
	akid, secret, session string
}

func newMockCreds(akid, secret, session string) *mockCreds {
	return &mockCreds{akid: akid, secret: secret, session: session}
}

func (m *mockCreds) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     m.akid,
		SecretAccessKey: m.secret,
		SessionToken:    m.session,
	}, nil
}

func TestS3Express_Retrieve(t *testing.T) {
	sdk.NowTime = func() time.Time {
		return time.Unix(0, 0)
	}
	mockClient := &mockCreateSession{
		calls: []mockCreateSessionCall{
			{
				output: &CreateSessionOutput{
					Credentials: &types.SessionCredentials{
						AccessKeyId:     aws.String("AccessKeyId-0"),
						Expiration:      aws.Time(time.Unix(3600, 0).UTC()),
						SecretAccessKey: aws.String("SecretAccessKey-0"),
						SessionToken:    aws.String("SessionToken-0"),
					},
				},
			},
			{
				output: &CreateSessionOutput{
					Credentials: &types.SessionCredentials{
						AccessKeyId:     aws.String("AccessKeyId-1"),
						Expiration:      aws.Time(time.Unix(7200, 0).UTC()),
						SecretAccessKey: aws.String("SecretAccessKey-1"),
						SessionToken:    aws.String("SessionToken-1"),
					},
				},
			},
		},
	}

	c := newDefaultS3ExpressCredentialsProvider()
	c.client = mockClient
	c.v4creds = newMockCreds("AKID", "SECRET", "SESSION")

	mockClient.wg.Add(3)
	c0, err := c.Retrieve(context.Background(), "bucket-0")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	c1, err := c.Retrieve(context.Background(), "bucket-1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	c2, err := c.Retrieve(context.Background(), "bucket-0")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected0 := aws.Credentials{
		AccessKeyID:     "AccessKeyId-0",
		SecretAccessKey: "SecretAccessKey-0",
		SessionToken:    "SessionToken-0",
		CanExpire:       true,
		Expires:         time.Unix(3600, 0).UTC(),
	}
	expected1 := aws.Credentials{
		AccessKeyID:     "AccessKeyId-1",
		SecretAccessKey: "SecretAccessKey-1",
		SessionToken:    "SessionToken-1",
		CanExpire:       true,
		Expires:         time.Unix(7200, 0).UTC(),
	}

	// one should have been a cache hit
	mockClient.expectCalled(t, 2)
	if expected0 != c0 {
		t.Errorf("expected credentials %v, got %v", expected0, c0)
	}
	if expected0 != c2 {
		t.Errorf("expected credentials %v, got %v", expected0, c2)
	}
	if expected1 != c1 {
		t.Errorf("expected credentials %v, got %v", expected1, c1)
	}
}

func TestS3Express_AsyncRefresh(t *testing.T) {
	sdk.NowTime = func() time.Time {
		return time.Unix(0, 0)
	}
	mockClient := &mockCreateSession{
		calls: []mockCreateSessionCall{
			{
				output: &CreateSessionOutput{
					Credentials: &types.SessionCredentials{
						AccessKeyId:     aws.String("AccessKeyId-0"),
						Expiration:      aws.Time(time.Unix(30, 0).UTC()),
						SecretAccessKey: aws.String("SecretAccessKey-0"),
						SessionToken:    aws.String("SessionToken-0"),
					},
				},
			},
			{
				output: &CreateSessionOutput{
					Credentials: &types.SessionCredentials{
						AccessKeyId:     aws.String("AccessKeyId-0"),
						Expiration:      aws.Time(time.Unix(3600, 0).UTC()),
						SecretAccessKey: aws.String("SecretAccessKey-0"),
						SessionToken:    aws.String("SessionToken-0"),
					},
				},
			},
		},
	}

	c := newDefaultS3ExpressCredentialsProvider()
	c.client = mockClient
	c.v4creds = newMockCreds("AKID", "SECRET", "SESSION")

	mockClient.wg.Add(2)
	c0, err := c.Retrieve(context.Background(), "bucket-0")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	c1, err := c.Retrieve(context.Background(), "bucket-0")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	mockClient.wg.Wait() // block on the async retrieve that we've now triggered

	// the first set should still be returned the 2nd time since it's still valid
	expected := aws.Credentials{
		AccessKeyID:     "AccessKeyId-0",
		SecretAccessKey: "SecretAccessKey-0",
		SessionToken:    "SessionToken-0",
		CanExpire:       true,
		Expires:         time.Unix(30, 0).UTC(),
	}

	// 2nd call should happen due to refresh window
	mockClient.expectCalled(t, 2)

	if expected != c0 {
		t.Errorf("expected credentials %v, got %v", expected, c0)
	}
	if expected != c1 {
		t.Errorf("expected credentials %v, got %v", expected, c1)
	}
}

type mockHTTP struct{}

func (*mockHTTP) Do(*http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}

func TestS3Express_OperationCredentialOverride(t *testing.T) {
	sdk.NowTime = func() time.Time {
		return time.Unix(0, 0)
	}

	createSessionClient := &mockCreateSession{
		calls: []mockCreateSessionCall{
			{
				output: &CreateSessionOutput{
					Credentials: &types.SessionCredentials{
						AccessKeyId:     aws.String("EXPRESS_AKID0"),
						SecretAccessKey: aws.String("EXPRESS_SECRET0"),
						SessionToken:    aws.String("EXPRESS_TOKEN0"),
						Expiration:      aws.Time(time.Unix(3600, 0).UTC()),
					},
				},
			},
			{
				output: &CreateSessionOutput{
					Credentials: &types.SessionCredentials{
						AccessKeyId:     aws.String("EXPRESS_AKID1"),
						SecretAccessKey: aws.String("EXPRESS_SECRET1"),
						SessionToken:    aws.String("EXPRESS_TOKEN1"),
						Expiration:      aws.Time(time.Unix(3600, 0).UTC()),
					},
				},
			},
		},
	}
	createSessionClient.wg.Add(2)

	svc := New(Options{
		Region:      "us-west-2",
		Credentials: newMockCreds("AKID0", "SECRET0", "SESSION0"),
		HTTPClient:  &mockHTTP{},
		APIOptions: []func(*middleware.Stack) error{
			func(stack *middleware.Stack) error {
				stack.Deserialize.Clear()
				return stack.Deserialize.Add(
					middleware.DeserializeMiddlewareFunc(
						"mockResponse",
						func(context.Context, middleware.DeserializeInput, middleware.DeserializeHandler) (middleware.DeserializeOutput, middleware.Metadata, error) {
							out := middleware.DeserializeOutput{
								Result: &GetObjectOutput{},
							}
							return out, middleware.Metadata{}, nil
						},
					),
					middleware.After,
				)
			},
		},
	})

	expressProvider, _ := svc.options.ExpressCredentials.(*defaultS3ExpressCredentialsProvider)
	expressProvider.client = createSessionClient

	_, err := svc.GetObject(context.Background(), &GetObjectInput{
		Bucket: aws.String("bucket--usw2-az1--x-s3"),
		Key:    aws.String("key"),
	})
	if err != nil {
		t.Errorf("get object: %v", err)
	}

	// there should be one set of credentials in the cache
	key0 := cacheKey{
		CredentialsHash: gethmac("AKID0", "SECRET0"),
		Bucket:          "bucket--usw2-az1--x-s3",
	}
	_, ok := expressProvider.cache.Get(key0)
	if !ok {
		t.Errorf("creds for AKID0/SECRET0 are missing")
	}

	_, err = svc.GetObject(context.Background(), &GetObjectInput{
		Bucket: aws.String("bucket--usw2-az1--x-s3"),
		Key:    aws.String("key"),
	}, func(o *Options) {
		o.Credentials = newMockCreds("AKID1", "SECRET1", "SESSION1")
	})
	if err != nil {
		t.Errorf("get object: %v", err)
	}

	// checking two things here:
	//   - we have a new cache entry since creds changed
	//   - note we're still using the original pointer, the operation finalizer
	//     should have copied it and passed the cache along
	key1 := cacheKey{
		CredentialsHash: gethmac("AKID1", "SECRET1"),
		Bucket:          "bucket--usw2-az1--x-s3",
	}
	_, ok = expressProvider.cache.Get(key1)
	if !ok {
		t.Errorf("creds for AKID1/SECRET1 are missing")
	}

	// repeat of 1st call, should be a cache hit
	_, err = svc.GetObject(context.Background(), &GetObjectInput{
		Bucket: aws.String("bucket--usw2-az1--x-s3"),
		Key:    aws.String("key"),
	})
	if err != nil {
		t.Errorf("get object: %v", err)
	}

	createSessionClient.expectCalled(t, 2)

	// each CreateSession must be signed with the credentials that key its
	// cache entry
	expectAKIDs := []string{"AKID0", "AKID1"}
	if len(createSessionClient.akids) != len(expectAKIDs) {
		t.Fatalf("expected CreateSession akids %v, got %v", expectAKIDs, createSessionClient.akids)
	}
	for i, v := range expectAKIDs {
		if createSessionClient.akids[i] != v {
			t.Errorf("CreateSession %d: expected akid %q, got %q", i, v, createSessionClient.akids[i])
		}
	}
}

var expressStubCredentialRe = regexp.MustCompile(`Credential=([^/]+)/`)

// expressStubHTTP answers CreateSession with session credentials named after
// the AKID that signed it, and records the session token of data-plane calls.
type expressStubHTTP struct {
	mu                 sync.Mutex
	createSessionAKIDs []string
	dataPlaneTokens    []string
}

func (s *expressStubHTTP) Do(r *http.Request) (*http.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	akid := ""
	if m := expressStubCredentialRe.FindStringSubmatch(r.Header.Get("Authorization")); m != nil {
		akid = m[1]
	}
	if _, ok := r.URL.Query()["session"]; ok {
		s.createSessionAKIDs = append(s.createSessionAKIDs, akid)
		body := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<CreateSessionResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Credentials><SessionToken>SESSION_TOKEN_OF_%s</SessionToken><SecretAccessKey>SESSION_SECRET_OF_%s</SecretAccessKey><AccessKeyId>SESSION_AKID_OF_%s</AccessKeyId><Expiration>2099-01-01T00:00:00Z</Expiration></Credentials></CreateSessionResult>`, akid, akid, akid)
		return &http.Response{
			StatusCode: 200,
			Header:     http.Header{"Content-Type": {"application/xml"}},
			Body:       io.NopCloser(bytes.NewBufferString(body)),
			Request:    r,
		}, nil
	}

	s.dataPlaneTokens = append(s.dataPlaneTokens, r.Header.Get("X-Amz-S3session-Token"))
	return &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
		Body:       io.NopCloser(bytes.NewBufferString("data")),
		Request:    r,
	}, nil
}

func expectStrings(t *testing.T, name string, expect, actual []string) {
	t.Helper()
	if len(expect) != len(actual) {
		t.Errorf("%s: expected %v, got %v", name, expect, actual)
		return
	}
	for i := range expect {
		if expect[i] != actual[i] {
			t.Errorf("%s: expected %v, got %v", name, expect, actual)
			return
		}
	}
}

const expressTestBucket = "bucket--usw2-az1--x-s3"

func TestS3Express_OperationCredentialOverrideSignsCreateSession(t *testing.T) {
	stub := &expressStubHTTP{}
	svc := New(Options{
		Region:      "us-west-2",
		Credentials: newMockCreds("AKID_A", "SECRET_A", ""),
		HTTPClient:  stub,
	})

	out, err := svc.GetObject(context.Background(), &GetObjectInput{
		Bucket: aws.String(expressTestBucket),
		Key:    aws.String("key"),
	}, func(o *Options) {
		o.Credentials = newMockCreds("AKID_B", "SECRET_B", "")
	})
	if err != nil {
		t.Fatalf("get object: %v", err)
	}
	out.Body.Close()

	out, err = svc.GetObject(context.Background(), &GetObjectInput{
		Bucket: aws.String(expressTestBucket),
		Key:    aws.String("key"),
	})
	if err != nil {
		t.Fatalf("get object: %v", err)
	}
	out.Body.Close()

	expectStrings(t, "CreateSession signers", []string{"AKID_B", "AKID_A"}, stub.createSessionAKIDs)
	expectStrings(t, "data-plane session tokens",
		[]string{"SESSION_TOKEN_OF_AKID_B", "SESSION_TOKEN_OF_AKID_A"}, stub.dataPlaneTokens)
}

func TestS3Express_PresignCredentialOverrideSignsCreateSession(t *testing.T) {
	stub := &expressStubHTTP{}
	svc := New(Options{
		Region:      "us-west-2",
		Credentials: newMockCreds("AKID_A", "SECRET_A", ""),
		HTTPClient:  stub,
	})

	req, err := NewPresignClient(svc).PresignGetObject(context.Background(), &GetObjectInput{
		Bucket: aws.String(expressTestBucket),
		Key:    aws.String("key"),
	}, func(po *PresignOptions) {
		po.ClientOptions = append(po.ClientOptions, func(o *Options) {
			o.Credentials = newMockCreds("AKID_B", "SECRET_B", "")
		})
	})
	if err != nil {
		t.Fatalf("presign get object: %v", err)
	}

	expectStrings(t, "CreateSession signers", []string{"AKID_B"}, stub.createSessionAKIDs)
	if !strings.Contains(req.URL, "SESSION_TOKEN_OF_AKID_B") {
		t.Errorf("expected presigned URL to carry session of AKID_B, got %s", req.URL)
	}
}

func TestS3Express_DerivedClientDoesNotRebindProvider(t *testing.T) {
	stub1 := &expressStubHTTP{}
	svc1 := New(Options{
		Region:      "us-west-2",
		Credentials: newMockCreds("AKID_A", "SECRET_A", ""),
		HTTPClient:  stub1,
	})

	stub2 := &expressStubHTTP{}
	svc2 := New(svc1.Options(), func(o *Options) {
		o.Credentials = newMockCreds("AKID_B", "SECRET_B", "")
		o.HTTPClient = stub2
	})

	if svc1.options.ExpressCredentials == svc2.options.ExpressCredentials {
		t.Errorf("expected derived client to have its own express credentials provider")
	}

	for _, svc := range []*Client{svc1, svc2} {
		out, err := svc.GetObject(context.Background(), &GetObjectInput{
			Bucket: aws.String(expressTestBucket),
			Key:    aws.String("key"),
		})
		if err != nil {
			t.Fatalf("get object: %v", err)
		}
		out.Body.Close()
	}

	expectStrings(t, "client 1 CreateSession signers", []string{"AKID_A"}, stub1.createSessionAKIDs)
	expectStrings(t, "client 1 data-plane session tokens", []string{"SESSION_TOKEN_OF_AKID_A"}, stub1.dataPlaneTokens)
	expectStrings(t, "client 2 CreateSession signers", []string{"AKID_B"}, stub2.createSessionAKIDs)
	expectStrings(t, "client 2 data-plane session tokens", []string{"SESSION_TOKEN_OF_AKID_B"}, stub2.dataPlaneTokens)
}

// rotatingCreds returns the first identity on its first call and the second
// identity on every call after that.
type rotatingCreds struct {
	mu    sync.Mutex
	calls int
}

func (r *rotatingCreds) Retrieve(context.Context) (aws.Credentials, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls++
	if r.calls == 1 {
		return aws.Credentials{AccessKeyID: "AKID_ROT1", SecretAccessKey: "SECRET_ROT1"}, nil
	}
	return aws.Credentials{AccessKeyID: "AKID_ROT2", SecretAccessKey: "SECRET_ROT2"}, nil
}

func TestS3Express_RotatingCredentialsSignCreateSessionWithCacheKeyIdentity(t *testing.T) {
	mockClient := &mockCreateSession{
		calls: []mockCreateSessionCall{
			{
				output: &CreateSessionOutput{
					Credentials: &types.SessionCredentials{
						AccessKeyId:     aws.String("AccessKeyId-0"),
						Expiration:      aws.Time(time.Now().Add(time.Hour).UTC()),
						SecretAccessKey: aws.String("SecretAccessKey-0"),
						SessionToken:    aws.String("SessionToken-0"),
					},
				},
			},
		},
	}
	mockClient.wg.Add(1)

	c := newDefaultS3ExpressCredentialsProvider()
	c.client = mockClient
	c.v4creds = &rotatingCreds{}

	if _, err := c.Retrieve(context.Background(), "bucket"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mockClient.expectCalled(t, 1)
	expectStrings(t, "CreateSession signers", []string{"AKID_ROT1"}, mockClient.akids)

	key := cacheKey{
		CredentialsHash: gethmac("AKID_ROT1", "SECRET_ROT1"),
		Bucket:          "bucket",
	}
	if _, ok := c.cache.Get(key); !ok {
		t.Errorf("expected session cached under the AKID_ROT1 key")
	}
}
