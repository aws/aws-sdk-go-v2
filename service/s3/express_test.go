package s3

import (
	"context"
	"net/http"
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
