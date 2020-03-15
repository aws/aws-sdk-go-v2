package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/private/protocol/restxml"
)

const protoOpGetObject = "GetObject"

// ProtoGetObjectRequest returns a request value for making API operation for
// Amazon Simple Storage Service.
//
// Retrieves objects from Amazon S3. To use GET, you must have READ access to
// the object. If you grant READ access to the anonymous user, you can return
// the object without using an authorization header.
//
// An Amazon S3 bucket has no directory hierarchy such as you would find in
// a typical computer file system. You can, however, create a logical hierarchy
// by using object key names that imply a folder structure. For example, instead
// of naming an object sample.jpg, you can name it photos/2006/February/sample.jpg.
//
// To get an object from such a logical hierarchy, specify the full key name
// for the object in the GET operation. For a virtual hosted-style request example,
// if you have the object photos/2006/February/sample.jpg, specify the resource
// as /photos/2006/February/sample.jpg. For a path-style request example, if
// you have the object photos/2006/February/sample.jpg in the bucket named examplebucket,
// specify the resource as /examplebucket/photos/2006/February/sample.jpg. For
// more information about request types, see HTTP Host Header Bucket Specification
// (https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html#VirtualHostingSpecifyBucket).
//
// To distribute large files to many people, you can save bandwidth costs by
// using BitTorrent. For more information, see Amazon S3 Torrent (https://docs.aws.amazon.com/AmazonS3/latest/dev/S3Torrent.html).
// For more information about returning the ACL of an object, see GetObjectAcl.
//
// If the object you are retrieving is stored in the GLACIER or DEEP_ARCHIVE
// storage classes, before you can retrieve the object you must first restore
// a copy using . Otherwise, this operation returns an InvalidObjectStateError
// error. For information about restoring archived objects, see Restoring Archived
// Objects (https://docs.aws.amazon.com/AmazonS3/latest/dev/restoring-objects.html).
//
// Encryption request headers, like x-amz-server-side-encryption, should not
// be sent for GET requests if your object uses server-side encryption with
// CMKs stored in AWS KMS (SSE-KMS) or server-side encryption with Amazon S3–managed
// encryption keys (SSE-S3). If your object does use these types of keys, you’ll
// get an HTTP 400 BadRequest error.
//
// If you encrypt an object by using server-side encryption with customer-provided
// encryption keys (SSE-C) when you store the object in Amazon S3, then when
// you GET the object, you must use the following headers:
//
//    * x-amz-server-side​-encryption​-customer-algorithm
//
//    * x-amz-server-side​-encryption​-customer-key
//
//    * x-amz-server-side​-encryption​-customer-key-MD5
//
// For more information about SSE-C, see Server-Side Encryption (Using Customer-Provided
// Encryption Keys) (https://docs.aws.amazon.com/AmazonS3/latest/dev/ServerSideEncryptionCustomerKeys.html).
//
// Assuming you have permission to read object tags (permission for the s3:GetObjectVersionTagging
// action), the response also returns the x-amz-tagging-count header that provides
// the count of number of tags associated with the object. You can use GetObjectTagging
// to retrieve the tag set associated with an object.
//
// Permissions
//
// You need the s3:GetObject permission for this operation. For more information,
// see Specifying Permissions in a Policy (https://docs.aws.amazon.com/AmazonS3/latest/dev/using-with-s3-actions.html).
// If the object you request does not exist, the error Amazon S3 returns depends
// on whether you also have the s3:ListBucket permission.
//
//    * If you have the s3:ListBucket permission on the bucket, Amazon S3 will
//    return an HTTP status code 404 ("no such key") error.
//
//    * If you don’t have the s3:ListBucket permission, Amazon S3 will return
//    an HTTP status code 403 ("access denied") error.
//
// Versioning
//
// By default, the GET operation returns the current version of an object. To
// return a different version, use the versionId subresource.
//
// If the current version of the object is a delete marker, Amazon S3 behaves
// as if the object was deleted and includes x-amz-delete-marker: true in the
// response.
//
// For more information about versioning, see PutBucketVersioning.
//
// Overriding Response Header Values
//
// There are times when you want to override certain response header values
// in a GET response. For example, you might override the Content-Disposition
// response header value in your GET request.
//
// You can override values for a set of response headers using the following
// query parameters. These response header values are sent only on a successful
// request, that is, when status code 200 OK is returned. The set of headers
// you can override using these parameters is a subset of the headers that Amazon
// S3 accepts when you create an object. The response headers that you can override
// for the GET response are Content-Type, Content-Language, Expires, Cache-Control,
// Content-Disposition, and Content-Encoding. To override these header values
// in the GET response, you use the following request parameters.
//
// You must sign the request, either using an Authorization header or a presigned
// URL, when using these parameters. They cannot be used with an unsigned (anonymous)
// request.
//
//    * response-content-type
//
//    * response-content-language
//
//    * response-expires
//
//    * response-cache-control
//
//    * response-content-disposition
//
//    * response-content-encoding
//
// Additional Considerations about Request Headers
//
// If both of the If-Match and If-Unmodified-Since headers are present in the
// request as follows: If-Match condition evaluates to true, and; If-Unmodified-Since
// condition evaluates to false; then, S3 returns 200 OK and the data requested.
//
// If both of the If-None-Match and If-Modified-Since headers are present in
// the request as follows:If-None-Match condition evaluates to false, and; If-Modified-Since
// condition evaluates to true; then, S3 returns 304 Not Modified response code.
//
// For more information about conditional requests, see RFC 7232 (https://tools.ietf.org/html/rfc7232).
//
// The following operations are related to GetObject:
//
//    * ListBuckets
//
//    * GetObjectAcl
//
//    // Example sending a request using GetObjectRequest.
//    req := client.GetObjectRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/s3-2006-03-01/GetObject
func (c *Client) ProtoGetObjectRequest(input *GetObjectInput) ProtoGetObjectRequest {
	op := &aws.Operation{
		Name:       protoOpGetObject,
		HTTPMethod: "GET",
		HTTPPath:   "/{Bucket}/{Key+}",
	}

	if input == nil {
		input = &GetObjectInput{}
	}

	req := c.newRequest(op, input, &GetObjectOutput{})
	// unmarshalMeta handler is used for metadata headers
	req.Handlers.UnmarshalMeta.Swap(restxml.UnmarshalMetaHandler.Name,
		protoGetObjectUnmarshaler{output: &GetObjectOutput{}}.namedHandler(),
	)
	return ProtoGetObjectRequest{Request: req, Input: input, Copy: c.ProtoGetObjectRequest}
}

// ProtoGetObjectRequest is the request type for the
// GetObject API operation.
type ProtoGetObjectRequest struct {
	*aws.Request
	Input *GetObjectInput
	Copy  func(*GetObjectInput) ProtoGetObjectRequest
}

// Send marshals and sends the GetObject API request.
func (r ProtoGetObjectRequest) Send(ctx context.Context) (*ProtoGetObjectResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &ProtoGetObjectResponse{
		GetObjectOutput: r.Request.Data.(*GetObjectOutput),
		response:        &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// ProtoGetObjectResponse is the response type for the
// GetObject API operation.
type ProtoGetObjectResponse struct {
	*GetObjectOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// GetObject request.
func (r *ProtoGetObjectResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
