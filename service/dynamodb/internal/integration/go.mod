module github.com/aws/aws-sdk-go-v2/service/dynamodb/internal/integration

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v0.24.1-0.20200924225551-a2b886903b8b
	github.com/aws/aws-sdk-go-v2/config v0.0.0-20200924225551-a2b886903b8b
	github.com/aws/aws-sdk-go-v2/credentials v0.0.0-20200924225551-a2b886903b8b // indirect
	github.com/aws/aws-sdk-go-v2/ec2imds v0.0.0-20200924225551-a2b886903b8b // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v0.0.0-20200924225551-a2b886903b8b
	github.com/aws/aws-sdk-go-v2/service/sts v0.0.0-20200924225551-a2b886903b8b // indirect
	github.com/awslabs/smithy-go v0.1.0
)

replace (
	github.com/aws/aws-sdk-go-v2 => ../../../../
	github.com/aws/aws-sdk-go-v2/config => ../../../../config/
	github.com/aws/aws-sdk-go-v2/credentials => ../../../../credentials/
	github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../
)

replace github.com/aws/aws-sdk-go-v2/ec2imds => ../../../../ec2imds/

replace github.com/aws/aws-sdk-go-v2/service/sts => ../../../../service/sts/
