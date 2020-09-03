/*
Package customizations provides customizations for the Amazon DynamoDB API client.

The DynamoDB API client uses two customizations, response checksum validation,
and manual content-encoding: gzip support.

Middleware layering

Checksum validation needs to be performed first in deserialization chain
on top of gzip decompression. Since the behavior of Deserialization is
in reverse order to the other stack steps its easier to consider that
"after" means "before".

    HTTP Response -> Checksum ->  gzip decompress -> deserialize

Response checksum validation

DynamoDB responses can include a X-Amz-Crc32 header with the CRC32 checksum
value of the response body. If the response body is content-encoding: gzip, the
checksum is of the gzipped response content.

If the header is present, the SDK should validate that the response payload
computed CRC32 checksum matches the value provided in the header. The checksum
header is based on the original payload provided returned by the service. Which
means that if the response is gzipped the checksum is of the gzipped response,
not the decompressed response bytes.

Customization option:
    DisableValidateResponseChecksum (Enabled by Default)

Accept encoding gzip

The Go HTTP client automatically supports accept-encoding and content-encoding
gzip by default. This default behavior is not desired by the SDK, and prevents
validating the response body's checksum. To prevent this the SDK must manually
control usage of content-encoding gzip.

To control content-encoding, the SDK must always set the `Accept-Encoding`
header to a value. This prevents the HTTP client from using gzip automatically.
When gzip is enabled on the API client, the SDK's customization will control
decompressing the gzip data in order to not break the checksum validation. When
gzip is disabled, the API client will disable gzip, preventing the HTTP
client's default behavior.

Customization option:
    EnableAcceptEncodingGzip (Disabled by Default)

*/
package customizations
