module github.com/aws/aws-sdk-go-v2/ec2imds

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.0.0-20200917190052-bb89e83d660c
	github.com/awslabs/smithy-go v0.0.0-20200920191232-15240aa5c76f
	github.com/google/go-cmp v0.5.2
)

replace github.com/aws/aws-sdk-go-v2 => ../

replace github.com/awslabs/smithy-go => ../../smithy-go
