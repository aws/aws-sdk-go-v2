package awsrestxml

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
	restV2 "github.com/aws/aws-sdk-go-v2/private/protocol/rest/v2"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ProtoPutObjectMarshaler defines marshaler for ProtoPutObject operation
type ProtoPutObjectMarshaler struct {
	Input *types.PutObjectInput
}

// MarshalOperation is the top level method used within a handler stack to marshal an operation
// This method calls appropriate marshal shape functions as per the input shape and protocol used by the service.
func (m ProtoPutObjectMarshaler) MarshalOperation(r *aws.Request) {
	var err error
	encoder := restV2.NewEncoder(r.HTTPRequest)

	err = MarshalPutObjectInputShapeAWSREST(m.Input, encoder)
	if err != nil {
		r.Error = err
		return
	}
	encoder.Encode()

	// Todo Instead of passing aws.Request directly to MarshalPutObjectInputShapeAWSXML;
	//  we should pass the payload as an argument
	if err = MarshalPutObjectInputShapeAWSXML(m.Input, r); err != nil {
		r.Error = err
		return
	}
}

// MarshalPutObjectInputShapeAWSREST is a stand alone function used to marshal the HTTP bindings a input shape.
// This method uses the rest encoder utility
func MarshalPutObjectInputShapeAWSREST(input *types.PutObjectInput, encoder *restV2.Encoder) error {
	// Encoding shapes with location `headers`
	marshalShapeMapForHeaders(encoder, "x-amz-meta-", input.Metadata)
	//  Encoding shapes with location `header`
	encoder.AddHeader("Cache-Control").String(*input.CacheControl)
	encoder.AddHeader("Content-Disposition").String(*input.ContentDisposition)
	encoder.AddHeader("Content-Language").String(*input.ContentLanguage)
	encoder.AddHeader("Content-Length").Integer(*input.ContentLength)
	encoder.AddHeader("Content-Md5").String(*input.ContentMD5)
	encoder.AddHeader("Content-Type").String(*input.ContentType)
	encoder.AddHeader("x-amz-acl").String(string(input.ACL))
	encoder.AddHeader("x-amz-grant-full-control").String(*input.GrantFullControl)
	encoder.AddHeader("x-amz-grant-read").String(*input.GrantRead)
	encoder.AddHeader("x-amz-grant-read-acp").String(*input.GrantReadACP)
	encoder.AddHeader("x-amz-grant-write-acp").String(*input.GrantWriteACP)
	encoder.AddHeader("x-amz-object-lock-legal-hold").String(string(input.ObjectLockLegalHoldStatus))
	encoder.AddHeader("x-amz-object-lock-mode").String(string(input.ObjectLockMode))
	encoder.AddHeader("x-amz-tagging").String(*input.Tagging)
	encoder.AddHeader("x-amz-request-payer").String(string(input.RequestPayer))
	encoder.AddHeader("x-amz-server-side-encryption-context").String(*input.SSEKMSEncryptionContext)
	encoder.AddHeader("x-amz-server-side-encryption-aws-kms-key-id").String(*input.SSEKMSKeyId)
	encoder.AddHeader("x-amz-server-side-encryption-customer-key-MD5").String(*input.SSECustomerKeyMD5)
	encoder.AddHeader("x-amz-server-side-encryption-customer-algorithm").String(*input.SSECustomerAlgorithm)
	encoder.AddHeader("x-amz-website-redirect-location").String(*input.WebsiteRedirectLocation)
	encoder.AddHeader("x-amz-storage-class").String(string(input.StorageClass))
	encoder.AddHeader("x-amz-server-side-encryption").String(string(input.ServerSideEncryption))
	if err := encoder.AddHeader("Expires").Time(*input.Expires, protocol.RFC822TimeFormatName); err != nil {
		return fmt.Errorf("failed to encode header for shape Expires: \n \t %v", err)
	}
	if err := encoder.AddHeader("x-amz-object-lock-retain-until-date").Time(*input.ObjectLockRetainUntilDate, protocol.ISO8601TimeFormatName); err != nil {
		return fmt.Errorf("failed to encode header for shape Expires: \n \t %v", err)
	}
	//  Encoding shapes with location `uri`
	if err := encoder.SetURI("Bucket").String(*input.Bucket); err != nil {
		return fmt.Errorf("failed to encode URI, \n\t %v", err)
	}
	if err := encoder.SetURI("Key").String(*input.Key); err != nil {
		return fmt.Errorf("failed to encode URI, \n\t %v", err)
	}
	return nil
}

// MarshalPutObjectInputShapeAWSXML is a stand alone function used to marshal the xml payload
// This should be generated according to the payload type for rest-xml protocol
func MarshalPutObjectInputShapeAWSXML(input *types.PutObjectInput, r *aws.Request) error {
	if input.Body != nil {
		r.SetReaderBody(input.Body)
	}
	return r.Error
}

// marshalShapeMapForHeaders is marshal function that takes in a map[string]string as an input along with an encoder
// and location Name which should be used to marshal the shape with location headers.
func marshalShapeMapForHeaders(encoder *restV2.Encoder, locationName string, input map[string]string) {
	headerObject := encoder.Headers(locationName)
	for k, v := range input {
		headerObject.AddHeader(k).String(v)
	}
}

// GetNamedBuildHandler returns a Named Build Handler for an operation marshal function
func (m ProtoPutObjectMarshaler) GetNamedBuildHandler() aws.NamedHandler {
	const BuildHandler = "ProtoPutBucket.BuildHandler"
	return aws.NamedHandler{
		Name: BuildHandler,
		Fn:   m.MarshalOperation,
	}
}
