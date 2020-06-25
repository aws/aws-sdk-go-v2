package aws

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jviney/aws-sdk-go-v2/internal/sdk"
)

type stubProvider struct {
	creds   Credentials
	expires time.Time
	err     error

	onInvalidate func(*stubProvider)
}

func (s *stubProvider) Retrieve(ctx context.Context) (Credentials, error) {
	creds := s.creds
	creds.Source = "stubProvider"
	creds.CanExpire = !s.expires.IsZero()
	creds.Expires = s.expires

	return creds, s.err
}

func (s *stubProvider) Invalidate() {
	s.onInvalidate(s)
}

func TestSafeCredentialsProvider_Cache(t *testing.T) {
	expect := Credentials{
		AccessKeyID:     "key",
		SecretAccessKey: "secret",
		CanExpire:       true,
		Expires:         time.Now().Add(10 * time.Minute),
	}

	var called bool
	p := &SafeCredentialsProvider{
		RetrieveFn: func() (Credentials, error) {
			if called {
				t.Fatalf("expect RetrieveFn to only be called once")
			}
			called = true
			return expect, nil
		},
	}

	for i := 0; i < 2; i++ {
		creds, err := p.Retrieve(context.Background())
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}
		if e, a := expect, creds; e != a {
			t.Errorf("expect %v creds, got %v", e, a)
		}
	}
}

func TestSafeCredentialsProvider_Expires(t *testing.T) {
	orig := sdk.NowTime
	defer func() { sdk.NowTime = orig }()
	var mockTime time.Time
	sdk.NowTime = func() time.Time { return mockTime }

	cases := []struct {
		Creds  func() Credentials
		Called int
	}{
		{
			Called: 2,
			Creds: func() Credentials {
				return Credentials{
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
					CanExpire:       true,
					Expires:         mockTime.Add(5),
				}
			},
		},
		{
			Called: 1,
			Creds: func() Credentials {
				return Credentials{
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
				}
			},
		},
		{
			Called: 6,
			Creds: func() Credentials {
				return Credentials{
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
					CanExpire:       true,
					Expires:         mockTime,
				}
			},
		},
	}

	for _, c := range cases {
		var called int
		p := &SafeCredentialsProvider{
			RetrieveFn: func() (Credentials, error) {
				called++
				return c.Creds(), nil
			},
		}

		p.Retrieve(context.Background())
		p.Retrieve(context.Background())
		p.Retrieve(context.Background())

		mockTime = mockTime.Add(10)

		p.Retrieve(context.Background())
		p.Retrieve(context.Background())
		p.Retrieve(context.Background())

		if e, a := c.Called, called; e != a {
			t.Errorf("expect %v called, got %v", e, a)
		}
	}
}

func TestSafeCredentialsProvider_Error(t *testing.T) {
	p := &SafeCredentialsProvider{
		RetrieveFn: func() (Credentials, error) {
			return Credentials{}, fmt.Errorf("failed")
		},
	}

	creds, err := p.Retrieve(context.Background())
	if err == nil {
		t.Fatalf("expect error, not none")
	}
	if e, a := "failed", err.Error(); e != a {
		t.Errorf("expect %q, got %q", e, a)
	}
	if e, a := (Credentials{}), creds; e != a {
		t.Errorf("expect empty creds, got %v", a)
	}
}

func TestSafeCredentialsProvider_Race(t *testing.T) {
	expect := Credentials{
		AccessKeyID:     "key",
		SecretAccessKey: "secret",
	}
	var called bool
	p := &SafeCredentialsProvider{
		RetrieveFn: func() (Credentials, error) {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			if called {
				t.Fatalf("expect RetrieveFn only called once")
			}
			called = true
			return expect, nil
		},
	}

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			creds, err := p.Retrieve(context.Background())
			if err != nil {
				t.Errorf("expect no error, got %v", err)
			}
			if e, a := expect, creds; e != a {
				t.Errorf("expect %v, got %v", e, a)
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

type stubSafeProviderConcurrent struct {
	SafeCredentialsProvider
	called uint32
	done   chan struct{}
}

func TestSafeProviderRetrieveConcurrent(t *testing.T) {
	stub := &stubSafeProviderConcurrent{
		done: make(chan struct{}),
	}

	stub.RetrieveFn = func() (Credentials, error) {
		atomic.AddUint32(&stub.called, 1)
		<-stub.done
		return Credentials{
			AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
			SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		}, nil
	}

	done := make(chan struct{})
	for i := 0; i < 2; i++ {
		go func() {
			stub.Retrieve(context.Background())
			done <- struct{}{}
		}()
	}

	// Validates that a single call to Retrieve is shared between two calls to Get
	stub.done <- struct{}{}
	<-done
	<-done

	if e, a := uint32(1), atomic.LoadUint32(&stub.called); e != a {
		t.Errorf("expected %v, got %v", e, a)
	}
}
