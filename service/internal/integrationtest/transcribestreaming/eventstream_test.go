//go:build integration
// +build integration

package transcribestreaming

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/internal/integrationtest"
	"github.com/aws/aws-sdk-go-v2/service/transcribestreaming"
	"github.com/aws/aws-sdk-go-v2/service/transcribestreaming/types"
)

var (
	audioFilename   string
	audioFormat     string
	audioLang       string
	audioSampleRate int
	audioFrameSize  int
	withDebug       bool
)

func init() {
	flag.BoolVar(&withDebug, "debug", false, "Include debug logging with test.")
	flag.StringVar(&audioFilename, "audio-file", "", "Audio file filename to perform test with.")
	flag.StringVar(&audioLang, "audio-lang", string(types.LanguageCodeEnUs), "Language of audio speech.")
	flag.StringVar(&audioFormat, "audio-format", string(types.MediaEncodingPcm), "Format of audio.")
	flag.IntVar(&audioSampleRate, "audio-sample", 16000, "Sample rate of the audio.")
	flag.IntVar(&audioFrameSize, "audio-frame", 15*1024, "Size of frames of audio uploaded.")
}

func TestInteg_StartStreamTranscription(t *testing.T) {
	var audio io.Reader
	if len(audioFilename) != 0 {
		audioFile, err := os.Open(audioFilename)
		if err != nil {
			t.Fatalf("expect to open file, %v", err)
		}
		defer audioFile.Close()
		audio = audioFile
	} else {
		b, err := base64.StdEncoding.DecodeString(
			`UklGRjzxPQBXQVZFZm10IBAAAAABAAEAgD4AAAB9AAACABAAZGF0YVTwPQAAAAAAAAAAAAAAAAD//wIA/f8EAA==`,
		)
		if err != nil {
			t.Fatalf("expect decode audio bytes, %v", err)
		}
		audio = bytes.NewReader(b)
	}

	cfg, _ := integrationtest.LoadConfigWithDefaultRegion("us-west-2")

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	client := transcribestreaming.NewFromConfig(cfg)
	resp, err := client.StartStreamTranscription(ctx, &transcribestreaming.StartStreamTranscriptionInput{
		LanguageCode:         types.LanguageCode(audioLang),
		MediaEncoding:        types.MediaEncoding(audioFormat),
		MediaSampleRateHertz: aws.Int32(int32(audioSampleRate)),
	})
	if err != nil {
		t.Fatalf("failed to start streaming, %v", err)
	}
	stream := resp.GetStream()
	defer stream.Close()

	go streamAudioFromReader(context.Background(), stream.Writer, audioFrameSize, audio)

	for event := range stream.Events() {
		switch e := event.(type) {
		case *types.TranscriptResultStreamMemberTranscriptEvent:
			t.Logf("got event, %v results", len(e.Value.Transcript.Results))
			for _, res := range e.Value.Transcript.Results {
				for _, alt := range res.Alternatives {
					t.Logf("* %s", aws.ToString(alt.Transcript))
				}
			}
		default:
			t.Fatalf("unexpected event, %T", event)
		}
	}

	if err := stream.Err(); err != nil {
		t.Fatalf("expect no error from stream, got %v", err)
	}
}

func TestInteg_StartStreamTranscription_contextClose(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString(
		`UklGRjzxPQBXQVZFZm10IBAAAAABAAEAgD4AAAB9AAACABAAZGF0YVTwPQAAAAAAAAAAAAAAAAD//wIA/f8EAA==`,
	)
	if err != nil {
		t.Fatalf("expect decode audio bytes, %v", err)
	}
	audio := bytes.NewReader(b)

	cfg, _ := integrationtest.LoadConfigWithDefaultRegion("us-west-2")

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	client := transcribestreaming.NewFromConfig(cfg)
	resp, err := client.StartStreamTranscription(ctx, &transcribestreaming.StartStreamTranscriptionInput{
		LanguageCode:         types.LanguageCodeEnUs,
		MediaEncoding:        types.MediaEncodingPcm,
		MediaSampleRateHertz: aws.Int32(16000),
	})
	if err != nil {
		t.Fatalf("failed to start streaming, %v", err)
	}

	stream := resp.GetStream()
	defer stream.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := streamAudioFromReader(ctx, stream.Writer, audioFrameSize, audio)
		if err == nil {
			t.Errorf("expect error")
		}
		if e, a := "context canceled", err.Error(); !strings.Contains(a, e) {
			t.Errorf("expect %q error in %q", e, a)
		}
		wg.Done()
	}()

	cancelFn()

Loop:
	for {
		select {
		case <-ctx.Done():
			break Loop
		case event, ok := <-stream.Events():
			if !ok {
				break Loop
			}
			switch e := event.(type) {
			case *types.TranscriptResultStreamMemberTranscriptEvent:
				t.Logf("got event, %v results", len(e.Value.Transcript.Results))
				for _, res := range e.Value.Transcript.Results {
					for _, alt := range res.Alternatives {
						t.Logf("* %s", aws.ToString(alt.Transcript))
					}
				}
			default:
				t.Fatalf("unexpected event, %T", event)
			}
		}
	}

	wg.Wait()

	if err := stream.Err(); err != nil {
		t.Fatalf("expect no error from stream, got %v", err)
	}
}

func streamAudioFromReader(ctx context.Context, stream transcribestreaming.AudioStreamWriter, frameSize int, input io.Reader) (err error) {
	defer func() {
		if closeErr := stream.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close stream, %v", closeErr)
		}
	}()

	frame := make([]byte, frameSize)
	for {
		var n int
		n, err = input.Read(frame)
		if n > 0 {
			err = stream.Send(ctx, &types.AudioStreamMemberAudioEvent{Value: types.AudioEvent{
				AudioChunk: frame[:n],
			}})
			if err != nil {
				return fmt.Errorf("failed to send audio event, %v", err)
			}
		}

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to read audio, %v", err)
		}
	}
}
