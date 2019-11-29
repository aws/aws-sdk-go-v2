package s3crypto

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	request "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// WrapEntry is builder that return a proper key decrypter and error
type WrapEntry func(Envelope) (CipherDataDecrypter, error)

// CEKEntry is a builder thatn returns a proper content decrypter and error
type CEKEntry func(CipherData) (ContentCipher, error)

// DecryptionClient is an S3 crypto client. The decryption client
// will handle all get object requests from Amazon S3.
// Supported key wrapping algorithms:
//	*AWS KMS
//
// Supported content ciphers:
//	* AES/GCM
//	* AES/CBC
type DecryptionClient struct {
	S3Client s3iface.ClientAPI
	// LoadStrategy is used to load the metadata either from the metadata of the object
	// or from a separate file in s3.
	//
	// Defaults to our default load strategy.
	LoadStrategy LoadStrategy

	WrapRegistry   map[string]WrapEntry
	CEKRegistry    map[string]CEKEntry
	PadderRegistry map[string]Padder
}

// NewDecryptionClient instantiates a new S3 crypto client
//
// Example:
//  cfg, err := external.LoadDefaultAWSConfig()
//	svc := s3crypto.NewDecryptionClient(cfg, func(svc *s3crypto.DecryptionClient{
//		// Custom client options here
//	}))
func NewDecryptionClient(cfg aws.Config, options ...func(*DecryptionClient)) *DecryptionClient {
	s3client := s3.New(cfg)
	client := &DecryptionClient{
		S3Client: s3client,
		LoadStrategy: defaultV2LoadStrategy{
			client: s3client,
		},
		WrapRegistry: map[string]WrapEntry{
			KMSWrap: (kmsKeyHandler{
				kms: kms.New(cfg),
			}).decryptHandler,
		},
		CEKRegistry: map[string]CEKEntry{
			AESGCMNoPadding: newAESGCMContentCipher,
			strings.Join([]string{AESCBC, AESCBCPadder.Name()}, "/"): newAESCBCContentCipher,
		},
		PadderRegistry: map[string]Padder{
			strings.Join([]string{AESCBC, AESCBCPadder.Name()}, "/"): AESCBCPadder,
			"NoPadding": NoPadder,
		},
	}
	for _, option := range options {
		option(client)
	}

	return client
}

// GetObjectRequest will make a request to s3 and retrieve the object. In this process
// decryption will be done. The SDK only supports V2 reads of KMS and GCM.
//
// Example:
//  cfg, err := external.LoadDefaultAWSConfig()
//	svc := s3crypto.NewDecryptionClient(cfg)
//	req, out := svc.GetObjectRequest(&s3.GetObjectInput {
//	  Key: aws.String("testKey"),
//	  Bucket: aws.String("testBucket"),
//	})
//	err := req.Send()
func (c *DecryptionClient) GetObjectRequest(input *types.GetObjectInput) s3.GetObjectRequest {
	req := c.S3Client.GetObjectRequest(input)
	out := req.Data.(*types.GetObjectOutput)

	req.Handlers.Unmarshal.PushBack(func(r *request.Request) {
		env, err := c.LoadStrategy.Load(r)
		if err != nil {
			r.Error = err
			out.Body.Close()
			return
		}

		// If KMS should return the correct CEK algorithm with the proper
		// KMS key provider
		cipher, err := c.contentCipherFromEnvelope(env)
		if err != nil {
			r.Error = err
			out.Body.Close()
			return
		}

		reader, err := cipher.DecryptContents(out.Body)
		if err != nil {
			r.Error = err
			out.Body.Close()
			return
		}
		out.Body = reader
	})
	return req
}

// GetObject is a wrapper for GetObjectRequest
func (c *DecryptionClient) GetObject(input *types.GetObjectInput) (*types.GetObjectOutput, error) {
	req := c.GetObjectRequest(input)
	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return resp.GetObjectOutput, nil
}

// GetObjectWithContext is a wrapper for GetObjectRequest with the additional
// context, and request options support.
//
// GetObjectWithContext is the same as GetObject with the additional support for
// Context input parameters. The Context must not be nil. A nil Context will
// cause a panic. Use the Context to add deadlining, timeouts, ect. In the future
// this may create sub-contexts for individual underlying requests.
func (c *DecryptionClient) GetObjectWithContext(ctx context.Context, input *types.GetObjectInput, opts ...request.Option) (*types.GetObjectOutput, error) {
	req := c.GetObjectRequest(input)
	req.ApplyOptions(opts...)
	resp, err := req.Send(ctx)
	if err != nil {
		return nil, err
	}
	return resp.GetObjectOutput, nil
}
