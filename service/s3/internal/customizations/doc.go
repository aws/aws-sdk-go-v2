/*
Package customizations provides customizations for the Amazon S3 API client.

This package provides support for following S3 customizations

    UpdateEndpoint Middleware: Virtual Host style url addressing

    processResponseWith200Error Middleware: Deserializing response error with 200 status code


Virtual Host style url addressing

Since serializers serialize by default as path style url, we use customization
to modify the endpoint url when `UsePathStyle` option on S3Client is unset or
false.

UpdateEndpoint middleware handler for virtual host url addressing needs to be
executed after request serialization.

    Middleware layering:

    HTTP Request -> operation serializer -> Update-Endpoint customization -> next middleware

Customization option:
    UsePathStyle (Disabled by Default)

Handle Error response with 200 status code

S3 operations: CopyObject, CompleteMultipartUpload, UploadPartCopy can have an
error Response with status code 2xx. The processResponseWith200Error middleware
customizations enables SDK to check for an error within response body prior to
deserialization.

As the check for 2xx response containing an error needs to be performed earlier
than response deserialization. Since the behavior of Deserialization is in
reverse order to the other stack steps its easier to consider that "after" means
"before".

    Middleware layering:

	HTTP Response -> handle 200 error customization -> deserialize

*/
package customizations
