package polly

import (
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/modeledendpoints"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/unit"
)

func TestRestGETStrategy(t *testing.T) {
	cfg := unit.Config()
	cfg.Region = "us-west-2"

	svc := New(cfg)
	r, _ := svc.SynthesizeSpeechRequest(nil)

	if err := restGETPresignStrategy(r); err != nil {
		t.Error(err)
	}
	if "GET" != r.HTTPRequest.Method {
		t.Errorf("Expected 'GET', but received %s", r.HTTPRequest.Method)
	}
	if r.Operation.BeforePresignFn == nil {
		t.Error("Expected non-nil value for 'BeforePresignFn'")
	}
}

func TestPresign(t *testing.T) {
	cfg := unit.Config()
	cfg.Region = "us-west-2"
	cfg.EndpointResolver = modeledendpoints.NewDefaultResolver()

	svc := New(cfg)
	r, _ := svc.SynthesizeSpeechRequest(&SynthesizeSpeechInput{
		Text:         aws.String("Moo"),
		OutputFormat: OutputFormatMp3,
		VoiceId:      VoiceIdGeraint,
	})
	url, err := r.Presign(time.Second)

	if err != nil {
		t.Error(err)
	}
	expectedURL := `^https://polly.us-west-2.amazonaws.com/v1/speech\?.*?OutputFormat=mp3.*?Text=Moo.*?VoiceId=Geraint.*`
	if matched, err := regexp.MatchString(expectedURL, url); !matched || err != nil {
		t.Errorf("Expected:\n%q\nReceived:\n%q\nError:\n%v\n", expectedURL, url, err)
	}
}
