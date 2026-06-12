package credentials

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestStaticCredentialsProvider(t *testing.T) {
	s := StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID:     "AKID",
			SecretAccessKey: "SECRET",
			SessionToken:    "",
		},
	}

	creds, err := s.Retrieve(context.Background())
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
	if e, a := "AKID", creds.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "SECRET", creds.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if l := creds.SessionToken; len(l) != 0 {
		t.Errorf("expect no token, got %v", l)
	}
}

func TestStaticCredentialsProviderIsExpired(t *testing.T) {
	s := StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID:     "AKID",
			SecretAccessKey: "SECRET",
			SessionToken:    "",
		},
	}

	creds, err := s.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if creds.Expired() {
		t.Errorf("expect static credentials to never expire")
	}
}

func TestStaticCredentialsProviderValidation(t *testing.T) {
	tests := []struct {
		name         string
		accessKey    string
		secretKey    string
		sessionToken string
		shouldError  bool
		expectedMsg  string
	}{
		{
			name:         "trailing newline in secret (issue #3304)",
			accessKey:    "AKID",
			secretKey:    "SECRET\n",
			sessionToken: "",
			shouldError:  true,
			expectedMsg:  "SecretAccessKey contains invalid whitespace",
		},
		{
			name:         "trailing space in access key",
			accessKey:    "AKID ",
			secretKey:    "SECRET",
			sessionToken: "",
			shouldError:  true,
			expectedMsg:  "AccessKeyID contains invalid whitespace",
		},
		{
			name:         "leading whitespace in token",
			accessKey:    "AKID",
			secretKey:    "SECRET",
			sessionToken: " TOKEN",
			shouldError:  true,
			expectedMsg:  "SessionToken contains invalid whitespace",
		},
		{
			name:         "tabs in secret key",
			accessKey:    "AKID",
			secretKey:    "\tSECRET\r",
			sessionToken: "",
			shouldError:  true,
			expectedMsg:  "SecretAccessKey contains invalid whitespace",
		},
		{
			name:         "valid credentials without whitespace",
			accessKey:    "AKID",
			secretKey:    "SECRET",
			sessionToken: "TOKEN",
			shouldError:  false,
		},
		{
			name:         "empty access key",
			accessKey:    "",
			secretKey:    "SECRET",
			sessionToken: "",
			shouldError:  true,
			expectedMsg:  "static credentials are empty",
		},
		{
			name:         "empty secret key",
			accessKey:    "AKID",
			secretKey:    "",
			sessionToken: "",
			shouldError:  true,
			expectedMsg:  "static credentials are empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStaticCredentialsProvider(tt.accessKey, tt.secretKey, tt.sessionToken)

			creds, err := s.Retrieve(context.Background())

			if tt.shouldError {
				if err == nil {
					t.Fatal("expected error for credentials with whitespace, got nil")
				}

				if err.Error() != tt.expectedMsg {
					t.Errorf("expected error %q, got %q", tt.expectedMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expect no error, got %v", err)
				}

				if e, a := tt.accessKey, creds.AccessKeyID; e != a {
					t.Errorf("expect AccessKeyID %q, got %q", e, a)
				}
				if e, a := tt.secretKey, creds.SecretAccessKey; e != a {
					t.Errorf("expect SecretAccessKey %q, got %q", e, a)
				}
				if e, a := tt.sessionToken, creds.SessionToken; e != a {
					t.Errorf("expect SessionToken %q, got %q", e, a)
				}
			}
		})
	}
}
