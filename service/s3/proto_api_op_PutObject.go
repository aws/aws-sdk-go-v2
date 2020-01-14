package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol/restxml"
)

const protoOpPutObject = "PutObject"

// ProtoPutObjectRequest returns a request value for making API operation for
// Amazon Simple Storage Service.
//
// Adds an object to a bucket. You must have WRITE permissions on a bucket to
// add an object to it.
//
// Amazon S3 never adds partial objects; if you receive a success response,
// Amazon S3 added the entire object to the bucket.
//
// Amazon S3 is a distributed system. If it receives multiple write requests
// for the same object simultaneously, it overwrites all but the last object
// written. Amazon S3 does not provide object locking; if you need this, make
// sure to build it into your application layer or use versioning instead.
//
// To ensure that data is not corrupted traversing the network, use the Content-MD5
// header. When you use this header, Amazon S3 checks the object against the
// provided MD5 value and, if they do not match, returns an error. Additionally,
// you can calculate the MD5 while putting an object to Amazon S3 and compare
// the returned ETag to the calculated MD5 value.
//
// To configure your application to send the request headers before sending
// the request body, use the 100-continue HTTP status code. For PUT operations,
// this helps you avoid sending the message body if the message is rejected
// based on the headers (for example, because authentication fails or a redirect
// occurs). For more information on the 100-continue HTTP status code, see Section
// 8.2.3 of http://www.ietf.org/rfc/rfc2616.txt (http://www.ietf.org/rfc/rfc2616.txt).
//
// You can optionally request server-side encryption. With server-side encryption,
// Amazon S3 encrypts your data as it writes it to disks in its data centers
// and decrypts the data when you access it. You have the option to provide
// your own encryption key or use AWS-managed encryption keys. For more information,
// see Using Server-Side Encryption (https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingServerSideEncryption.html).
//
// Access Permissions
//
// You can optionally specify the accounts or groups that should be granted
// specific permissions on the new object. There are two ways to grant the permissions
// using the request headers:
//
//    * Specify a canned ACL with the x-amz-acl request header. For more information,
//    see Canned ACL (https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html#CannedACL).
//
//    * Specify access permissions explicitly with the x-amz-grant-read, x-amz-grant-read-acp,
//    x-amz-grant-write-acp, and x-amz-grant-full-control headers. These parameters
//    map to the set of permissions that Amazon S3 supports in an ACL. For more
//    information, see Access Control List (ACL) Overview (https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html).
//
// You can use either a canned ACL or specify access permissions explicitly.
// You cannot do both.
//
// Server-Side- Encryption-Specific Request Headers
//
// You can optionally tell Amazon S3 to encrypt data at rest using server-side
// encryption. Server-side encryption is for data encryption at rest. Amazon
// S3 encrypts your data as it writes it to disks in its data centers and decrypts
// it when you access it. The option you use depends on whether you want to
// use AWS-managed encryption keys or provide your own encryption key.
//
//    * Use encryption keys managed Amazon S3 or customer master keys (CMKs)
//    stored in AWS Key Management Service (KMS) – If you want AWS to manage
//    the keys used to encrypt data, specify the following headers in the request.
//    x-amz-server-side​-encryption x-amz-server-side-encryption-aws-kms-key-id
//    x-amz-server-side-encryption-context If you specify x-amz-server-side-encryption:aws:kms,
//    but don't provide x-amz-server-side- encryption-aws-kms-key-id, Amazon
//    S3 uses the AWS managed CMK in AWS KMS to protect the data. All GET and
//    PUT requests for an object protected by AWS KMS fail if you don't make
//    them with SSL or by using SigV4. For more information on Server-Side Encryption
//    with CMKs stored in AWS KMS (SSE-KMS), see Protecting Data Using Server-Side
//    Encryption with CMKs stored in AWS (https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingKMSEncryption.html).
//
//    * Use customer-provided encryption keys – If you want to manage your
//    own encryption keys, provide all the following headers in the request.
//    x-amz-server-side​-encryption​-customer-algorithm x-amz-server-side​-encryption​-customer-key
//    x-amz-server-side​-encryption​-customer-key-MD5 For more information
//    on Server-Side Encryption with CMKs stored in KMS (SSE-KMS), see Protecting
//    Data Using Server-Side Encryption with CMKs stored in AWS KMS (https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingKMSEncryption.html).
//
// Access-Control-List (ACL)-Specific Request Headers
//
// You also can use the following access control–related headers with this
// operation. By default, all objects are private. Only the owner has full access
// control. When adding a new object, you can grant permissions to individual
// AWS accounts or to predefined groups defined by Amazon S3. These permissions
// are then added to the Access Control List (ACL) on the object. For more information,
// see Using ACLs (https://docs.aws.amazon.com/AmazonS3/latest/dev/S3_ACLs_UsingACLs.html).
// With this operation, you can grant access permissions using one of the following
// two methods:
//
//    * Specify a canned ACL (x-amz-acl) — Amazon S3 supports a set of predefined
//    ACLs, known as canned ACLs. Each canned ACL has a predefined set of grantees
//    and permissions. For more information, see Canned ACL (https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html#CannedACL).
//
//    * Specify access permissions explicitly — To explicitly grant access
//    permissions to specific AWS accounts or groups, use the following headers.
//    Each header maps to specific permissions that Amazon S3 supports in an
//    ACL. For more information, see Access Control List (ACL) Overview (https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html).
//    In the header, you specify a list of grantees who get the specific permission.
//    To grant permissions explicitly use: x-amz-grant-read x-amz-grant-write
//    x-amz-grant-read-acp x-amz-grant-write-acp x-amz-grant-full-control You
//    specify each grantee as a type=value pair, where the type is one of the
//    following: emailAddress – if the value specified is the email address
//    of an AWS account Using email addresses to specify a grantee is only supported
//    in the following AWS Regions: US East (N. Virginia) US West (N. California)
//    US West (Oregon) Asia Pacific (Singapore) Asia Pacific (Sydney) Asia Pacific
//    (Tokyo) EU (Ireland) South America (São Paulo) For a list of all the
//    Amazon S3 supported regions and endpoints, see Regions and Endpoints (https://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region)
//    in the AWS General Reference id – if the value specified is the canonical
//    user ID of an AWS account uri – if you are granting permissions to a
//    predefined group For example, the following x-amz-grant-read header grants
//    the AWS accounts identified by email addresses permissions to read object
//    data and its metadata: x-amz-grant-read: emailAddress="xyz@amazon.com",
//    emailAddress="abc@amazon.com"
//
// Server-Side- Encryption-Specific Request Headers
//
// You can optionally tell Amazon S3 to encrypt data at rest using server-side
// encryption. Server-side encryption is for data encryption at rest. Amazon
// S3 encrypts your data as it writes it to disks in its data centers and decrypts
// it when you access it. The option you use depends on whether you want to
// use AWS-managed encryption keys or provide your own encryption key.
//
//    * Use encryption keys managed by Amazon S3 or customer master keys (CMKs)
//    stored in AWS Key Management Service (KMS) – If you want AWS to manage
//    the keys used to encrypt data, specify the following headers in the request.
//    x-amz-server-side​-encryption x-amz-server-side-encryption-aws-kms-key-id
//    x-amz-server-side-encryption-context If you specify x-amz-server-side-encryption:aws:kms,
//    but don't provide x-amz-server-side- encryption-aws-kms-key-id, Amazon
//    S3 uses the default AWS KMS CMK to protect the data. All GET and PUT requests
//    for an object protected by AWS KMS fail if you don't make them with SSL
//    or by using SigV4. For more information on Server-Side Encryption with
//    CMKs stored in AWS KMS (SSE-KMS), see Protecting Data Using Server-Side
//    Encryption with CMKs stored in AWS KMS (https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingKMSEncryption.html).
//
//    * Use customer-provided encryption keys – If you want to manage your
//    own encryption keys, provide all the following headers in the request.
//    If you use this feature, the ETag value that Amazon S3 returns in the
//    response is not the MD5 of the object. x-amz-server-side​-encryption​-customer-algorithm
//    x-amz-server-side​-encryption​-customer-key x-amz-server-side​-encryption​-customer-key-MD5
//    For more information on Server-Side Encryption with CMKs stored in AWS
//    KMS (SSE-KMS), see Protecting Data Using Server-Side Encryption with CMKs
//    stored in AWS KMS (https://docs.aws.amazon.com/AmazonS3/latest/dev/UsingKMSEncryption.html).
//
// Storage Class Options
//
// By default, Amazon S3 uses the Standard storage class to store newly created
// objects. The Standard storage class provides high durability and high availability.
// You can specify other storage classes depending on the performance needs.
// For more information, see Storage Classes (https://docs.aws.amazon.com/AmazonS3/latest/dev/storage-class-intro.html)
// in the Amazon Simple Storage Service Developer Guide.
//
// Versioning
//
// If you enable versioning for a bucket, Amazon S3 automatically generates
// a unique version ID for the object being stored. Amazon S3 returns this ID
// in the response using the x-amz-version-id response header. If versioning
// is suspended, Amazon S3 always uses null as the version ID for the object
// stored. For more information about returning the versioning state of a bucket,
// see GetBucketVersioning. If you enable versioning for a bucket, when Amazon
// S3 receives multiple write requests for the same object simultaneously, it
// stores all of the objects.
//
// Related Resources
//
//    * CopyObject
//
//    * DeleteObject
//
//    // Example sending a request using PutObjectRequest.
//    req := client.ProtoPutObjectRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/s3-2006-03-01/PutObject
func (c *Client) ProtoPutObjectRequest(input *PutObjectInput) ProtoPutObjectRequest {
	op := &aws.Operation{
		Name:       protoOpPutObject,
		HTTPMethod: "PUT",
		HTTPPath:   "/{Bucket}/{Key+}",
	}

	if input == nil {
		input = &PutObjectInput{}
	}

	req := c.newRequest(op, input, &PutObjectOutput{})

	// swap existing build handler on svc, with a new named build handler
	req.Handlers.Build.Swap(restxml.BuildHandler.Name,
		protoPutObjectMarshaler{
			input: input,
		}.NamedHandler(),
	)

	return ProtoPutObjectRequest{Request: req, Input: input, Copy: c.ProtoPutObjectRequest}
}

// ProtoPutObjectRequest is the request type for the
// ProtoPutObject API operation.
type ProtoPutObjectRequest struct {
	*aws.Request
	Input *PutObjectInput
	Copy  func(*PutObjectInput) ProtoPutObjectRequest
}

// Send marshals and sends the PutObject API request.
func (r ProtoPutObjectRequest) Send(ctx context.Context) (*ProtoPutObjectResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &ProtoPutObjectResponse{
		PutObjectOutput: r.Request.Data.(*PutObjectOutput),
		response:        &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// ProtoPutObjectResponse is the response type for the
// ProtoPutObject API operation.
type ProtoPutObjectResponse struct {
	*PutObjectOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// PutObject request.
func (r *ProtoPutObjectResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
