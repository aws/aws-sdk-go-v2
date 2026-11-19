package v4

import (
	"bytes"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestSigningKeyDeriver_SecretRotatedSameAccessKey(t *testing.T) {
	deriver := NewSigningKeyDeriver()
	signingTime := NewSigningTime(time.Date(2026, 3, 16, 10, 0, 0, 0, time.UTC))
	laterSameDay := NewSigningTime(time.Date(2026, 3, 16, 14, 0, 0, 0, time.UTC))

	first := aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET-1"}
	rotated := aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET-2"}

	got := deriver.DeriveKey(first, "s3", "us-east-1", signingTime)
	if expect := deriveKey(first.SecretAccessKey, "s3", "us-east-1", signingTime); !bytes.Equal(expect, got) {
		t.Fatalf("first key: expect %x, got %x", expect, got)
	}

	got = deriver.DeriveKey(rotated, "s3", "us-east-1", laterSameDay)
	if expect := deriveKey(rotated.SecretAccessKey, "s3", "us-east-1", laterSameDay); !bytes.Equal(expect, got) {
		t.Errorf("key after secret rotation: expect %x, got %x", expect, got)
	}
}

func TestSigningKeyDeriver_ReusesKeyForSameCredentialsSameDay(t *testing.T) {
	deriver := NewSigningKeyDeriver()
	creds := aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET"}

	first := deriver.DeriveKey(creds, "s3", "us-east-1", NewSigningTime(time.Date(2026, 3, 16, 10, 0, 0, 0, time.UTC)))
	second := deriver.DeriveKey(creds, "s3", "us-east-1", NewSigningTime(time.Date(2026, 3, 16, 14, 0, 0, 0, time.UTC)))

	if &first[0] != &second[0] {
		t.Errorf("expect cached key to be reused for the same credentials on the same day")
	}
}
