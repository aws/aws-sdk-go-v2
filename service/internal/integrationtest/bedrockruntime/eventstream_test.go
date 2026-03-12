//go:build integration
// +build integration

package bedrockruntime

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

const (
	AudioFile = "japan16k.raw"
	ModelID   = "amazon.nova-sonic-v1:0"
)

var PreambleEvents = []string{
	`{"event":{"sessionStart":{"inferenceConfiguration":{"maxTokens":10000,"topP":0.95,"temperature":0.9}}}}`,
	`{"event":{"promptStart":{"promptName":"126680f5-5859-4d15-ae70-488de4146484","textOutputConfiguration":{"mediaType":"text/plain"},"audioOutputConfiguration":{"mediaType":"audio/lpcm","sampleRateHertz":24000,"sampleSizeBits":16,"channelCount":1,"voiceId":"en_us_matthew","encoding":"base64","audioType":"SPEECH"},"toolUseOutputConfiguration":{"mediaType":"application/json"},"toolConfiguration":{"tools":[]}}}}`,
	`{"event":{"contentStart":{"promptName":"126680f5-5859-4d15-ae70-488de4146484","contentName":"a6431ef2-e23c-4f8c-a552-3f308629d3c3","type":"TEXT","interactive":true,"textInputConfiguration":{"mediaType":"text/plain"}}}}`,
	`{"event":{"textInput":{"promptName":"126680f5-5859-4d15-ae70-488de4146484","contentName":"a6431ef2-e23c-4f8c-a552-3f308629d3c3","content":"You are a friend. The user and you will engage in a spoken dialog exchanging the transcripts of a natural real-time conversation. Keep your responses short, generally two or three sentences for chatty scenarios.","role":"SYSTEM"}}}`,
	`{"event":{"contentEnd":{"promptName":"126680f5-5859-4d15-ae70-488de4146484","contentName":"a6431ef2-e23c-4f8c-a552-3f308629d3c3"}}}`,
	`{"event":{"contentStart":{"promptName":"126680f5-5859-4d15-ae70-488de4146484","contentName":"b3917935-2398-4889-94a8-e677f6c3e351","type":"AUDIO","interactive":true,"audioInputConfiguration":{"mediaType":"audio/lpcm","sampleRateHertz":16000,"sampleSizeBits":16,"channelCount":1,"audioType":"SPEECH","encoding":"base64"}}}}`,
}

const ContentEnd = `{"event": {"contentEnd": {"promptName": "126680f5-5859-4d15-ae70-488de4146484", "contentName": "b3917935-2398-4889-94a8-e677f6c3e351" } } }`
const PromptEnd = `{"event": {"promptEnd": {"promptName": "126680f5-5859-4d15-ae70-488de4146484", "contentName": "b3917935-2398-4889-94a8-e677f6c3e351" } } }`
const SessionEnd = `{"event":{"sessionEnd":{}}}`

const ContentName2ndEvent = "562cec92-9c44-4363-9605-428eb860335c"

var (
	contentStart2 = fmt.Sprintf(`{"event":{"contentStart":{"promptName":"126680f5-5859-4d15-ae70-488de4146484","contentName":"%s","type":"AUDIO","interactive":true,"audioInputConfiguration":{"mediaType":"audio/lpcm","sampleRateHertz":16000,"sampleSizeBits":16,"channelCount":1,"audioType":"SPEECH","encoding":"base64"}}}}`, ContentName2ndEvent)
	audioEvent2   = fmt.Sprintf(`{"event":{"audioInput":{"promptName":"126680f5-5859-4d15-ae70-488de4146484","contentName":"%s","content":"%%s","role":"USER"}}}`, ContentName2ndEvent)
	contentEnd2   = fmt.Sprintf(`{"event": {"contentEnd": {"promptName": "126680f5-5859-4d15-ae70-488de4146484", "contentName": "%s" } } }`, ContentName2ndEvent)
)

const AudioEvent = `{"event":{"audioInput":{"promptName":"126680f5-5859-4d15-ae70-488de4146484","contentName":"b3917935-2398-4889-94a8-e677f6c3e351","content":"%s","role":"USER"}}}`

func sendAudioContent(ctx context.Context, t *testing.T, stream *bedrockruntime.InvokeModelWithBidirectionalStreamEventStream, file *os.File, contentStart, audioEventTemplate, contentEnd string, contentEndReceived chan bool) error {
	if contentStart != "" {
		if err := sendEvent(ctx, stream, contentStart); err != nil {
			return fmt.Errorf("sending content start: %w", err)
		}
	}

	file.Seek(0, io.SeekStart)
	buffer := make([]byte, 1024)
	chunkCount := 0
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading audio file: %w", err)
		}

		blob := base64.StdEncoding.EncodeToString(buffer[:n])
		audioEvent := fmt.Sprintf(audioEventTemplate, blob)
		if err := sendEvent(ctx, stream, audioEvent); err != nil {
			return fmt.Errorf("sending audio event: %w", err)
		}
		chunkCount++
	}

	t.Logf("Sent %d audio chunks", chunkCount)

	if err := sendEvent(ctx, stream, contentEnd); err != nil {
		return fmt.Errorf("sending content end: %w", err)
	}

	select {
	case <-contentEndReceived:
		t.Log("Received response signal")
	case <-ctx.Done():
		t.Error("Unexpected context cancel, failing")
		return ctx.Err()
	}

	return nil
}

func TestBedrockBidirectionalStream(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if _, err := os.Stat(AudioFile); os.IsNotExist(err) {
		t.Fatalf("Audio file %s not found", AudioFile)
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		// Nova Sonic is only available on us-east-1 (02/2026)
		config.WithRegion("us-east-1"),
	)

	client := bedrockruntime.NewFromConfig(cfg)

	input := &bedrockruntime.InvokeModelWithBidirectionalStreamInput{
		ModelId: aws.String(ModelID),
	}

	response, err := client.InvokeModelWithBidirectionalStream(ctx, input)
	if err != nil {
		t.Fatalf("failed to start bidirectional response: %v", err)
	}
	stream := response.GetStream()

	var wg sync.WaitGroup
	contentEndReceived := make(chan bool, 1)

	// Start a goroutine to send events
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer t.Log("Sender goroutine completed")

		// Send preamble events
		for i, event := range PreambleEvents {
			err = sendEvent(ctx, stream, event)
			if err != nil {
				t.Errorf("Error sending event: %v", err)
				return
			}
			t.Logf("Sent preamble event %d", i)
		}
		t.Log("Sent all preamble events")

		// Open and read audio file
		file, err := os.Open(AudioFile)
		if err != nil {
			t.Errorf("Error opening audio file: %v", err)
			return
		}
		defer file.Close()

		if err := sendAudioContent(ctx, t, stream, file, "", AudioEvent, ContentEnd, contentEndReceived); err != nil {
			t.Errorf("Error sending first audio content: %v", err)
			return
		}

		if err := sendAudioContent(ctx, t, stream, file, contentStart2, audioEvent2, contentEnd2, contentEndReceived); err != nil {
			t.Errorf("Error sending second audio content: %v", err)
			return
		}

		if err := sendEvent(ctx, stream, PromptEnd); err != nil {
			t.Errorf("Error sending prompt end: %v", err)
		}

		if err := sendEvent(ctx, stream, SessionEnd); err != nil {
			t.Errorf("Error sending session end: %v", err)
		}

		if err := stream.Writer.Close(); err != nil {
			t.Errorf("Error closing response: %v", err)
			return
		}
		t.Log("Closed writer successfully")
	}()

	// Wait until there's a reply from the server
	<-response.GetInitialReply()
	responseCount := readEvents(ctx, t, stream, contentEndReceived)

	// Check for response errors
	if streamErr := stream.Err(); streamErr != nil {
		t.Errorf("Stream error: %v", streamErr)
	}

	if responseCount == 0 {
		t.Error("No responses received from server")
	}
	// from all events, ensuring we got all the different kind of events
	expectedEvents := 419
	if responseCount != expectedEvents {
		t.Errorf("Expected %d responses, got %d", expectedEvents, responseCount)
	}

	stream.Close()
	wg.Wait()
}

func TestWrongModel(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
	t.Cleanup(cancel)

	cfg, err := config.LoadDefaultConfig(
		ctx,
		// Nova Sonic is only available on us-east-1 (02/2026)
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		t.Fatalf("unable to load SDK config: %v", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	input := &bedrockruntime.InvokeModelWithBidirectionalStreamInput{
		ModelId: aws.String("wrong"),
	}

	resp, err := client.InvokeModelWithBidirectionalStream(ctx, input)
	if err != nil {
		t.Fatalf("expected no error even with wrong input since the error will be communicated on the stream, instead got: %v", err)
	}
	<-resp.GetInitialReply()
	stream := resp.GetStream()
	select {
	case <-ctx.Done():
		t.Error("Context cancelled, closing writer")
	case event, ok := <-stream.Events():
		if !ok {
			t.Log("Event stream closed")
			break
		}
		if event != nil {
			t.Errorf("Expected no event, got %v with type %T", event, event)
		}
	}
	err = resp.GetStream().Err()
	if err == nil {
		t.Fatal("Expected resp error on wrong model, got nil")
	}
	target := &types.ValidationException{}
	if !errors.As(err, &target) {
		t.Fatalf("Expected target error to be of type %T, but was %T", target, err)
	}
}

func TestMalformedPreambleEvent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(
		ctx,
		// Nova Sonic is only available on us-east-1 (02/2026)
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		t.Fatalf("unable to load SDK config: %v", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	input := &bedrockruntime.InvokeModelWithBidirectionalStreamInput{
		ModelId: aws.String(ModelID),
	}
	response, err := client.InvokeModelWithBidirectionalStream(ctx, input)
	if err != nil {
		t.Fatalf("failed to start bidirectional response: %v", err)
	}

	stream := response.GetStream()
	defer stream.Close()

	err = sendEvent(ctx, stream, PreambleEvents[0])
	if err != nil {
		t.Fatalf("Error sending first preamble event: %v", err)
	}

	// Create a malformed version by adding an extra character to make invalid JSON
	// Take the second preamble event and add an extra 'X' character at the end
	malformedEvent := PreambleEvents[1] + "X"

	err = sendEvent(ctx, stream, malformedEvent)
	if err != nil {
		t.Logf("Got expected error when sending malformed event: %v", err)
	}

	// Wait for server to process the malformed event and respond with error
	t.Log("Waiting for server response to malformed event...")

	// Wait for initial reply
	t.Log("Waiting for initial reply")
	<-response.GetInitialReply()

	select {
	case event, ok := <-stream.Events():
		if !ok {
			t.Log("Event stream closed")
		} else {
			t.Fatalf("Received event expected none: %v", event)
		}
	case <-time.After(1 * time.Second):
		t.Log("No events received within timeout")
	}

	// Check for stream error
	streamErr := stream.Err()
	if streamErr == nil {
		t.Fatal("Expected stream error due to malformed preamble event, but got nil")
	}

	t.Logf("Successfully caught expected stream error: %v", streamErr)

	target := &types.ValidationException{}
	if !errors.As(streamErr, &target) {
		t.Fatalf("Expected target error type %T but got %T", target, err)
	}
}

func sendEvent(ctx context.Context, stream *bedrockruntime.InvokeModelWithBidirectionalStreamEventStream, event string) error {
	chunk := &types.InvokeModelWithBidirectionalStreamInputMemberChunk{
		Value: types.BidirectionalInputPayloadPart{
			Bytes: []byte(event),
		},
	}

	if err := stream.Writer.Send(ctx, chunk); err != nil {
		return err
	}
	return nil
}

func readEvents(ctx context.Context, t *testing.T, stream *bedrockruntime.InvokeModelWithBidirectionalStreamEventStream, done chan<- bool) int {
	responseCount := 0
	for {
		select {
		case event, ok := <-stream.Events():
			if !ok {
				t.Log("Event stream closed")
				return responseCount
			}
			switch v := event.(type) {
			case *types.InvokeModelWithBidirectionalStreamOutputMemberChunk:
				chunk := v
				responseCount++
				bytes := string(chunk.Value.Bytes)
				t.Logf("Received chunk %d: %s", responseCount, bytes)
				if strings.Contains(bytes, "contentEnd") && strings.Contains(bytes, "END_TURN") {
					done <- true
				}
			default:
				// really never happens, only to make the compiler/linter happy
				t.Logf("Other kind of message: %v", v)
			}
		case <-ctx.Done():
			t.Log("Context cancelled while waiting for responses")
			return responseCount
		}
	}
}
