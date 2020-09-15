module github.com/aws/aws-sdk-go-v2/config

go 1.15

require github.com/aws/aws-sdk-go-v2 v0.24.0

require (
	github.com/aws/aws-sdk-go-v2/credentials v0.0.0-20200915195926-9dd18af694c4
	github.com/awslabs/smithy-go v0.0.0-20200828214850-b1c39f43623b
)

replace github.com/aws/aws-sdk-go-v2 => ../

replace github.com/aws/aws-sdk-go-v2/credentials => ../credentials
