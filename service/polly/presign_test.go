package polly

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func TestPresignOpSynthesizeSpeechInput(t *testing.T) {
	cases := map[string]struct {
		LexiconNames    []string
		OutputFormat    types.OutputFormat
		SampleRate      *string
		Text            *string
		TextType        types.TextType
		VoiceID         types.VoiceId
		Engine          types.Engine
		LanguageCode    types.LanguageCode
		SpeechMarkTypes []types.SpeechMarkType
		ExpectStream    string
		Error           error
		ExpectError     bool
	}{
		"Single LexiconNames": {
			LexiconNames: []string{"abc"},
			OutputFormat: types.OutputFormatMp3,
			SampleRate:   aws.String("128"),
			Text:         aws.String("Test"),
			TextType:     types.TextTypeText,
			VoiceID:      types.VoiceIdAmy,
			ExpectStream: "LexiconNames=abc&OutputFormat=mp3&SampleRate=128&Text=Test&TextType=text&VoiceId=Amy",
		},
		"Multiple LexiconNames": {
			LexiconNames: []string{"abc", "mno"},
			OutputFormat: types.OutputFormatMp3,
			SampleRate:   aws.String("128"),
			Text:         aws.String("Test"),
			TextType:     types.TextTypeText,
			VoiceID:      types.VoiceIdAmy,
			ExpectStream: "LexiconNames=abc&LexiconNames=mno&OutputFormat=mp3&SampleRate=128&Text=Test&TextType=text&VoiceId=Amy",
		},
		"Text needs parsing": {
			LexiconNames: []string{"abc", "mno"},
			OutputFormat: types.OutputFormatMp3,
			SampleRate:   aws.String("128"),
			Text:         aws.String("Test /Text"),
			TextType:     types.TextTypeText,
			VoiceID:      types.VoiceIdAmy,
			ExpectStream: "LexiconNames=abc&LexiconNames=mno&OutputFormat=mp3&SampleRate=128&Text=Test+%2FText&TextType=text&VoiceId=Amy",
		},
		"With Engine and LanguageCode": {
			LexiconNames: []string{"abc"},
			OutputFormat: types.OutputFormatMp3,
			SampleRate:   aws.String("128"),
			Text:         aws.String("Test"),
			TextType:     types.TextTypeText,
			VoiceID:      types.VoiceIdAmy,
			Engine:       types.EngineNeural,
			LanguageCode: types.LanguageCodeArb,
			ExpectStream: "Engine=neural&LanguageCode=arb&LexiconNames=abc&OutputFormat=mp3&SampleRate=128&Text=Test&TextType=text&VoiceId=Amy",
		},
		"With single SpeechMarkType": {
			OutputFormat:    types.OutputFormatJson,
			Text:            aws.String("Test"),
			VoiceID:         types.VoiceIdJoanna,
			SpeechMarkTypes: []types.SpeechMarkType{types.SpeechMarkTypeSsml},
			ExpectStream:    "OutputFormat=json&SpeechMarkTypes=ssml&Text=Test&VoiceId=Joanna",
		},
		"With multiple SpeechMarkTypes": {
			OutputFormat:    types.OutputFormatJson,
			Text:            aws.String("Test"),
			VoiceID:         types.VoiceIdJoanna,
			SpeechMarkTypes: []types.SpeechMarkType{types.SpeechMarkTypeSsml, types.SpeechMarkTypeWord, types.SpeechMarkTypeSentence},
			ExpectStream:    "OutputFormat=json&SpeechMarkTypes=ssml&SpeechMarkTypes=word&SpeechMarkTypes=sentence&Text=Test&VoiceId=Joanna",
		},
		"With Engine, LanguageCode and SpeechMarkTypes": {
			LexiconNames:    []string{"abc", "mno"},
			OutputFormat:    types.OutputFormatJson,
			Text:            aws.String("Test"),
			VoiceID:         types.VoiceIdMatthew,
			Engine:          types.EngineStandard,
			LanguageCode:    types.LanguageCodeDeAt,
			SpeechMarkTypes: []types.SpeechMarkType{types.SpeechMarkTypeWord, types.SpeechMarkTypeViseme},
			ExpectStream:    "Engine=standard&LanguageCode=de-AT&LexiconNames=abc&LexiconNames=mno&OutputFormat=json&SpeechMarkTypes=word&SpeechMarkTypes=viseme&Text=Test&VoiceId=Matthew",
		},
		"Next serializer return error": {
			Error:       fmt.Errorf("next handler return error"),
			ExpectError: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			req := smithyhttp.NewStackRequest().(*smithyhttp.Request)
			var updatedRequest *smithyhttp.Request
			param := &SynthesizeSpeechInput{
				LexiconNames:    c.LexiconNames,
				OutputFormat:    c.OutputFormat,
				SampleRate:      c.SampleRate,
				Text:            c.Text,
				TextType:        c.TextType,
				VoiceId:         c.VoiceID,
				Engine:          c.Engine,
				LanguageCode:    c.LanguageCode,
				SpeechMarkTypes: c.SpeechMarkTypes,
			}

			m := presignOpSynthesizeSpeechInput{}
			_, _, err := m.HandleSerialize(context.Background(),
				middleware.SerializeInput{
					Request:    req,
					Parameters: param,
				},
				middleware.SerializeHandlerFunc(func(ctx context.Context, input middleware.SerializeInput) (
					out middleware.SerializeOutput, metadata middleware.Metadata, err error) {
					updatedRequest = input.Request.(*smithyhttp.Request)
					return out, metadata, c.Error
				}),
			)

			if err != nil && !c.ExpectError {
				t.Fatalf("expect no error, got %v", err)
			} else if err != nil != c.ExpectError {
				t.Fatalf("expect error but got nil")
			}

			stream := updatedRequest.GetStream()
			b, _ := io.ReadAll(stream)
			if e, a := c.ExpectStream, string(b); e != a {
				t.Errorf("expect request stream value %v, got %v", e, a)
			}
		})
	}
}
