module github.com/aws/aws-sdk-go-v2/internal/auth

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.21.0
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.1.4
	github.com/aws/smithy-go v1.14.2
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/auth => ../../internal/auth

replace github.com/aws/smithy-go => ../../../smithy-go
