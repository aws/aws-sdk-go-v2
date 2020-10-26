module github.com/aws/aws-sdk-go-v2/service/ec2

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.27.1-0.20201020212433-5fb7a9ec04bb
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v0.0.0-20201020212433-5fb7a9ec04bb
	github.com/awslabs/smithy-go v0.2.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => ../../service/internal/presigned-url/
