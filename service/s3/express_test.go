package s3

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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
