package restjson_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/private/protocol/restjson"
	"github.com/aws/aws-sdk-go-v2/service/elastictranscoder"
	"github.com/aws/aws-sdk-go-v2/service/elastictranscoder/types"
)

var (
	elastictranscoderSvc *elastictranscoder.Client
)

func TestMain(m *testing.M) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	cfg, _ := external.LoadDefaultAWSConfig()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)

	cfg.Credentials = aws.NewStaticCredentialsProvider("Key", "Secret", "Token")
	cfg.Region = "us-west-2"

	elastictranscoderSvc = elastictranscoder.New(cfg)

	c := m.Run()
	server.Close()
	os.Exit(c)
}

func BenchmarkRESTJSONBuild_Complex_ETCCreateJob(b *testing.B) {
	params := elastictranscoderCreateJobInput()

	benchRESTJSONBuild(b, func() *aws.Request {
		return elastictranscoderSvc.CreateJobRequest(params).Request
	})
}

func BenchmarkRESTJSONBuild_Simple_ETCListJobsByPipeline(b *testing.B) {
	params := elastictranscoderListJobsByPipeline()

	benchRESTJSONBuild(b, func() *aws.Request {
		return elastictranscoderSvc.ListJobsByPipelineRequest(params).Request
	})
}

func BenchmarkRESTJSONRequest_Complex_CFCreateJob(b *testing.B) {
	benchRESTJSONRequest(b, func() *aws.Request {
		return elastictranscoderSvc.CreateJobRequest(elastictranscoderCreateJobInput()).Request
	})
}

func BenchmarkRESTJSONRequest_Simple_ETCListJobsByPipeline(b *testing.B) {
	benchRESTJSONRequest(b, func() *aws.Request {
		return elastictranscoderSvc.ListJobsByPipelineRequest(elastictranscoderListJobsByPipeline()).Request
	})
}

func benchRESTJSONBuild(b *testing.B, reqFn func() *aws.Request) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := reqFn()
		restjson.Build(req)
		if req.Error != nil {
			b.Fatal("Unexpected error", req.Error)
		}
	}
}

func benchRESTJSONRequest(b *testing.B, reqFn func() *aws.Request) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := reqFn().Send()
		if err != nil {
			b.Fatal("Unexpected error", err)
		}
	}
}

func elastictranscoderListJobsByPipeline() *types.ListJobsByPipelineInput {
	return &types.ListJobsByPipelineInput{
		PipelineId: aws.String("Id"), // Required
		Ascending:  aws.String("Ascending"),
		PageToken:  aws.String("Id"),
	}
}

func elastictranscoderCreateJobInput() *types.CreateJobInput {
	return &types.CreateJobInput{
		Input: &types.JobInput{ // Required
			AspectRatio: aws.String("AspectRatio"),
			Container:   aws.String("JobContainer"),
			DetectedProperties: &types.DetectedProperties{
				DurationMillis: aws.Int64(1),
				FileSize:       aws.Int64(1),
				FrameRate:      aws.String("FloatString"),
				Height:         aws.Int64(1),
				Width:          aws.Int64(1),
			},
			Encryption: &types.Encryption{
				InitializationVector: aws.String("ZeroTo255String"),
				Key:                  aws.String("Base64EncodedString"),
				KeyMd5:               aws.String("Base64EncodedString"),
				Mode:                 aws.String("EncryptionMode"),
			},
			FrameRate:  aws.String("FrameRate"),
			Interlaced: aws.String("Interlaced"),
			Key:        aws.String("Key"),
			Resolution: aws.String("Resolution"),
		},
		PipelineId: aws.String("Id"), // Required
		Output: &types.CreateJobOutputResult{
			AlbumArt: &types.JobAlbumArt{
				Artwork: []types.Artwork{
					{ // Required
						AlbumArtFormat: aws.String("JpgOrPng"),
						Encryption: &types.Encryption{
							InitializationVector: aws.String("ZeroTo255String"),
							Key:                  aws.String("Base64EncodedString"),
							KeyMd5:               aws.String("Base64EncodedString"),
							Mode:                 aws.String("EncryptionMode"),
						},
						InputKey:      aws.String("WatermarkKey"),
						MaxHeight:     aws.String("DigitsOrAuto"),
						MaxWidth:      aws.String("DigitsOrAuto"),
						PaddingPolicy: aws.String("PaddingPolicy"),
						SizingPolicy:  aws.String("SizingPolicy"),
					},
					// More values...
				},
				MergePolicy: aws.String("MergePolicy"),
			},
			Captions: &types.Captions{
				CaptionFormats: []types.CaptionFormat{
					{ // Required
						Encryption: &types.Encryption{
							InitializationVector: aws.String("ZeroTo255String"),
							Key:                  aws.String("Base64EncodedString"),
							KeyMd5:               aws.String("Base64EncodedString"),
							Mode:                 aws.String("EncryptionMode"),
						},
						Format:  aws.String("CaptionFormatFormat"),
						Pattern: aws.String("CaptionFormatPattern"),
					},
					// More values...
				},
				CaptionSources: []types.CaptionSource{
					{ // Required
						Encryption: &types.Encryption{
							InitializationVector: aws.String("ZeroTo255String"),
							Key:                  aws.String("Base64EncodedString"),
							KeyMd5:               aws.String("Base64EncodedString"),
							Mode:                 aws.String("EncryptionMode"),
						},
						Key:        aws.String("Key"),
						Label:      aws.String("Name"),
						Language:   aws.String("Key"),
						TimeOffset: aws.String("TimeOffset"),
					},
					// More values...
				},
				MergePolicy: aws.String("CaptionMergePolicy"),
			},
			Composition: []types.Clip{
				{ // Required
					TimeSpan: &types.TimeSpan{
						Duration:  aws.String("Time"),
						StartTime: aws.String("Time"),
					},
				},
				// More values...
			},
			Encryption: &types.Encryption{
				InitializationVector: aws.String("ZeroTo255String"),
				Key:                  aws.String("Base64EncodedString"),
				KeyMd5:               aws.String("Base64EncodedString"),
				Mode:                 aws.String("EncryptionMode"),
			},
			Key:             aws.String("Key"),
			PresetId:        aws.String("Id"),
			Rotate:          aws.String("Rotate"),
			SegmentDuration: aws.String("FloatString"),
			ThumbnailEncryption: &types.Encryption{
				InitializationVector: aws.String("ZeroTo255String"),
				Key:                  aws.String("Base64EncodedString"),
				KeyMd5:               aws.String("Base64EncodedString"),
				Mode:                 aws.String("EncryptionMode"),
			},
			ThumbnailPattern: aws.String("ThumbnailPattern"),
			Watermarks: []types.JobWatermark{
				{ // Required
					Encryption: &types.Encryption{
						InitializationVector: aws.String("ZeroTo255String"),
						Key:                  aws.String("Base64EncodedString"),
						KeyMd5:               aws.String("Base64EncodedString"),
						Mode:                 aws.String("EncryptionMode"),
					},
					InputKey:          aws.String("WatermarkKey"),
					PresetWatermarkId: aws.String("PresetWatermarkId"),
				},
				// More values...
			},
		},
		OutputKeyPrefix: aws.String("Key"),
		Outputs: []types.CreateJobOutputResult{
			{ // Required
				AlbumArt: &types.JobAlbumArt{
					Artwork: []types.Artwork{
						{ // Required
							AlbumArtFormat: aws.String("JpgOrPng"),
							Encryption: &types.Encryption{
								InitializationVector: aws.String("ZeroTo255String"),
								Key:                  aws.String("Base64EncodedString"),
								KeyMd5:               aws.String("Base64EncodedString"),
								Mode:                 aws.String("EncryptionMode"),
							},
							InputKey:      aws.String("WatermarkKey"),
							MaxHeight:     aws.String("DigitsOrAuto"),
							MaxWidth:      aws.String("DigitsOrAuto"),
							PaddingPolicy: aws.String("PaddingPolicy"),
							SizingPolicy:  aws.String("SizingPolicy"),
						},
						// More values...
					},
					MergePolicy: aws.String("MergePolicy"),
				},
				Captions: &types.Captions{
					CaptionFormats: []types.CaptionFormat{
						{ // Required
							Encryption: &types.Encryption{
								InitializationVector: aws.String("ZeroTo255String"),
								Key:                  aws.String("Base64EncodedString"),
								KeyMd5:               aws.String("Base64EncodedString"),
								Mode:                 aws.String("EncryptionMode"),
							},
							Format:  aws.String("CaptionFormatFormat"),
							Pattern: aws.String("CaptionFormatPattern"),
						},
						// More values...
					},
					CaptionSources: []types.CaptionSource{
						{ // Required
							Encryption: &types.Encryption{
								InitializationVector: aws.String("ZeroTo255String"),
								Key:                  aws.String("Base64EncodedString"),
								KeyMd5:               aws.String("Base64EncodedString"),
								Mode:                 aws.String("EncryptionMode"),
							},
							Key:        aws.String("Key"),
							Label:      aws.String("Name"),
							Language:   aws.String("Key"),
							TimeOffset: aws.String("TimeOffset"),
						},
						// More values...
					},
					MergePolicy: aws.String("CaptionMergePolicy"),
				},
				Composition: []types.Clip{
					{ // Required
						TimeSpan: &types.TimeSpan{
							Duration:  aws.String("Time"),
							StartTime: aws.String("Time"),
						},
					},
					// More values...
				},
				Encryption: &types.Encryption{
					InitializationVector: aws.String("ZeroTo255String"),
					Key:                  aws.String("Base64EncodedString"),
					KeyMd5:               aws.String("Base64EncodedString"),
					Mode:                 aws.String("EncryptionMode"),
				},
				Key:             aws.String("Key"),
				PresetId:        aws.String("Id"),
				Rotate:          aws.String("Rotate"),
				SegmentDuration: aws.String("FloatString"),
				ThumbnailEncryption: &types.Encryption{
					InitializationVector: aws.String("ZeroTo255String"),
					Key:                  aws.String("Base64EncodedString"),
					KeyMd5:               aws.String("Base64EncodedString"),
					Mode:                 aws.String("EncryptionMode"),
				},
				ThumbnailPattern: aws.String("ThumbnailPattern"),
				Watermarks: []types.JobWatermark{
					{ // Required
						Encryption: &types.Encryption{
							InitializationVector: aws.String("ZeroTo255String"),
							Key:                  aws.String("Base64EncodedString"),
							KeyMd5:               aws.String("Base64EncodedString"),
							Mode:                 aws.String("EncryptionMode"),
						},
						InputKey:          aws.String("WatermarkKey"),
						PresetWatermarkId: aws.String("PresetWatermarkId"),
					},
					// More values...
				},
			},
			// More values...
		},
		Playlists: []types.CreateJobPlaylist{
			{ // Required
				Format: aws.String("PlaylistFormat"),
				HlsContentProtection: &types.HlsContentProtection{
					InitializationVector:  aws.String("ZeroTo255String"),
					Key:                   aws.String("Base64EncodedString"),
					KeyMd5:                aws.String("Base64EncodedString"),
					KeyStoragePolicy:      aws.String("KeyStoragePolicy"),
					LicenseAcquisitionUrl: aws.String("ZeroTo512String"),
					Method:                aws.String("HlsContentProtectionMethod"),
				},
				Name: aws.String("Filename"),
				OutputKeys: []string{
					"Key", // Required
					// More values...
				},
				PlayReadyDrm: &types.PlayReadyDrm{
					Format:                aws.String("PlayReadyDrmFormatString"),
					InitializationVector:  aws.String("ZeroTo255String"),
					Key:                   aws.String("NonEmptyBase64EncodedString"),
					KeyId:                 aws.String("KeyIdGuid"),
					KeyMd5:                aws.String("NonEmptyBase64EncodedString"),
					LicenseAcquisitionUrl: aws.String("OneTo512String"),
				},
			},
			// More values...
		},
		UserMetadata: map[string]string{
			"Key": "String", // Required
			// More values...
		},
	}
}
