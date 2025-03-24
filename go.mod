module github.com/Enflick/aws-sdk-go-v2

require (
	github.com/Enflick/smithy-go v1.3.0
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/Enflick/aws-sdk-go-v2 => ./

replace github.com/Enflick/aws-sdk-go-v2/aws => ./aws

replace github.com/aws/smithy-go => github.com/Enflick/smithy-go v1.3.0

go 1.20
