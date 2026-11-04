package s3

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/internal/sync/singleflight"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
)

type mockCreateSession struct {
	wg sync.WaitGroup

	calls []mockCreateSessionCall
	times int
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

func (m *mockCreateSession) CreateSession(context.Context, *CreateSessionInput, ...func(*Options)) (*CreateSessionOutput, error) {
	defer m.wg.Done()
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
}

// blockingCreateSession blocks every CreateSession call until release is
// closed or the call's context is canceled, and records what happened.
type blockingCreateSession struct {
	started chan struct{}
	release chan struct{}

	mu       sync.Mutex
	calls    int
	canceled int
}

func newBlockingCreateSession() *blockingCreateSession {
	return &blockingCreateSession{
		started: make(chan struct{}, 100),
		release: make(chan struct{}),
	}
}

func (m *blockingCreateSession) CreateSession(ctx context.Context, _ *CreateSessionInput, _ ...func(*Options)) (*CreateSessionOutput, error) {
	m.mu.Lock()
	m.calls++
	m.mu.Unlock()
	m.started <- struct{}{}

	select {
	case <-m.release:
	case <-ctx.Done():
		m.mu.Lock()
		m.canceled++
		m.mu.Unlock()
		return nil, ctx.Err()
	}

	return &CreateSessionOutput{
		Credentials: &types.SessionCredentials{
			AccessKeyId:     aws.String("AccessKeyId-1"),
			Expiration:      aws.Time(time.Unix(3600, 0).UTC()),
			SecretAccessKey: aws.String("SecretAccessKey-1"),
			SessionToken:    aws.String("SessionToken-1"),
		},
	}, nil
}

func (m *blockingCreateSession) counts() (calls, canceled int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.calls, m.canceled
}

func TestS3Express_CloneSharesSingleflight(t *testing.T) {
	sdk.NowTime = func() time.Time {
		return time.Unix(0, 0)
	}
	mockClient := newBlockingCreateSession()

	base := newDefaultS3ExpressCredentialsProvider()
	base.client = mockClient
	v4creds := newMockCreds("AKID", "SECRET", "SESSION")
	base.v4creds = v4creds

	key := cacheKey{
		CredentialsHash: gethmac("AKID", "SECRET"),
		Bucket:          "bucket-0",
	}

	// Every operation clones the provider (finalizeOperationExpressCredentials).
	// Start one cold-cache retrieve and wait for its CreateSession to be in
	// flight.
	first := base.CloneWithBaseCredentials(v4creds).doChanRetrieve(context.Background(), key)
	<-mockClient.started

	// Retrieves on other clones for the same key must join that call instead
	// of starting their own.
	var chans []<-chan singleflight.Result
	for i := 0; i < 10; i++ {
		chans = append(chans, base.CloneWithBaseCredentials(v4creds).doChanRetrieve(context.Background(), key))
	}

	close(mockClient.release)
	for _, ch := range append(chans, first) {
		if r := <-ch; r.Err != nil {
			t.Fatalf("expected no error, got %v", r.Err)
		}
	}

	if calls, _ := mockClient.counts(); calls != 1 {
		t.Errorf("expected 1 call to CreateSession across clones, got %d", calls)
	}
}

func TestS3Express_AsyncRefreshNotCanceledByCaller(t *testing.T) {
	sdk.NowTime = func() time.Time {
		return time.Unix(0, 0)
	}
	mockClient := newBlockingCreateSession()

	c := newDefaultS3ExpressCredentialsProvider()
	c.client = mockClient
	c.v4creds = newMockCreds("AKID", "SECRET", "SESSION")

	// Seed the cache with a session that is valid but inside the refresh
	// window, so Retrieve returns it and starts an async refresh.
	key := cacheKey{
		CredentialsHash: gethmac("AKID", "SECRET"),
		Bucket:          "bucket-0",
	}
	c.cache.Put(key, &aws.Credentials{
		AccessKeyID:     "AccessKeyId-0",
		SecretAccessKey: "SecretAccessKey-0",
		SessionToken:    "SessionToken-0",
		CanExpire:       true,
		Expires:         time.Unix(30, 0).UTC(),
	})

	ctx, cancel := context.WithCancel(context.Background())
	creds, err := c.Retrieve(ctx, "bucket-0")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if creds.AccessKeyID != "AccessKeyId-0" {
		t.Errorf("expected cached credentials, got %v", creds.AccessKeyID)
	}

	// The operation that triggered the refresh finishes and cancels its
	// context before CreateSession returns.
	<-mockClient.started
	cancel()
	close(mockClient.release)

	// Wait for the refresh to finish by joining it.
	r := <-c.doChanRetrieve(context.Background(), key)
	if r.Err != nil {
		t.Fatalf("expected refresh to succeed, got %v", r.Err)
	}

	calls, canceled := mockClient.counts()
	if canceled != 0 {
		t.Errorf("expected async refresh not to be canceled, got %d canceled", canceled)
	}
	if calls != 1 {
		t.Errorf("expected 1 call to CreateSession, got %d", calls)
	}
	if cached, _ := c.cache.Get(key); cached.AccessKeyID != "AccessKeyId-1" {
		t.Errorf("expected refreshed credentials in cache, got %v", cached.AccessKeyID)
	}
}

func TestS3Express_CanceledWaiterDoesNotCancelSharedRetrieve(t *testing.T) {
	sdk.NowTime = func() time.Time {
		return time.Unix(0, 0)
	}
	mockClient := newBlockingCreateSession()

	base := newDefaultS3ExpressCredentialsProvider()
	base.client = mockClient
	v4creds := newMockCreds("AKID", "SECRET", "SESSION")
	base.v4creds = v4creds

	// The first operation starts the cold-cache CreateSession and then gives
	// up (for example its timeout fires).
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		_, err := base.CloneWithBaseCredentials(v4creds).Retrieve(ctx, "bucket-0")
		errCh <- err
	}()
	<-mockClient.started

	// A second operation joins the same in-flight call.
	key := cacheKey{
		CredentialsHash: gethmac("AKID", "SECRET"),
		Bucket:          "bucket-0",
	}
	second := base.CloneWithBaseCredentials(v4creds).doChanRetrieve(context.Background(), key)

	cancel()
	if err := <-errCh; err == nil {
		t.Errorf("expected canceled caller to get an error")
	}

	close(mockClient.release)
	r := <-second
	if r.Err != nil {
		t.Fatalf("expected the other caller to get credentials, got %v", r.Err)
	}
	if got := r.Val.(aws.Credentials).AccessKeyID; got != "AccessKeyId-1" {
		t.Errorf("expected AccessKeyId-1, got %v", got)
	}

	calls, canceled := mockClient.counts()
	if canceled != 0 {
		t.Errorf("expected shared CreateSession not to be canceled, got %d canceled", canceled)
	}
	if calls != 1 {
		t.Errorf("expected 1 call to CreateSession, got %d", calls)
	}
}
