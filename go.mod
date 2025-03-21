module github.com/Enflick/aws-sdk-go-v2

require (
	github.com/aws/smithy-go v1.19.0
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ./

replace github.com/aws/aws-sdk-go-v2/aws => ./aws

go 1.20
