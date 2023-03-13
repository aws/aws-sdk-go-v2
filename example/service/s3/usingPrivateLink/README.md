# Example

This example demonstrates how you can use the AWS SDK for Go V2's Amazon S3 client
to use AWS PrivateLink for Amazon S3.

# Usage

## How to build vpc endpoint url

To access S3 bucket data using the s3 interface endpoints, prefix the vpc
endpoint with `bucket`. For eg, use endpoint url as `https://bucket.vpce-0xxxxxxx-xxx8xxg.s3.us-west-2.vpce.amazonaws.com`
to access S3 bucket data via the associated vpc endpoint.

To access S3 access point data using the s3 interface endpoints, prefix the vpc
endpoint with `accesspoint`. For eg, use endpoint url as `https://accesspoint.vpce-0xxxxxxx-xxxx8xxg.s3.us-west-2.vpce.amazonaws.com`
to access S3 access point data via the associated vpc endpoint.

To work with S3 control using the s3 interface endpoints, prefix the vpc endpoint
with `control`. For eg, use endpoint url as `https://control.vpce-0xxxxxxx-xxx8xxg.s3.us-west-2.vpce.amazonaws.com`
to use S3 Control operations with the associated vpc endpoint.

## How to use the built vpc endpoint url

Provide the vpc endpoint url (built as per instructions above) to the SDK by providing a custom endpoint resolver on
client config or as a part of request options. The VPC endpoint url should be provided as a custom endpoint,
which means the Endpoint source should be set as EndpointSourceCustom.

You can use client's `EndpointResolverFromURL` utility which by default sets the Endpoint source as Custom source.
Note that the SDK may mutate the vpc endpoint as per the input provided to work with ARNs, unless you explicitly set
HostNameImmutable property for the endpoint.

The example will create s3 client's that use appropriate vpc endpoint url. The example
will then create a bucket of the name provided in code. Replace the value of
the `accountID` const with the account ID for your AWS account. The
`vpcBucketEndpointUrl`, `vpcAccesspointEndpoint`, `vpcControlEndpoint`, `bucket`,
`keyName`, and `accessPoint` const variables need to be updated to match the name
of the appropriate vpc endpoint, Bucket, Object Key, and Access Point that will be
created by the example.

```sh
 AWS_REGION=<region> go run usingPrivateLink.go
 ```
