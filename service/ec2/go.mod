module github.com/aws/aws-sdk-go-v2/service/ec2

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.31.1-0.20210105194811-58b543144e2a
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v0.2.0
	github.com/aws/smithy-go v0.5.1-0.20210104190327-c7045c94c1ec
	github.com/jmespath/go-jmespath v0.4.0
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
