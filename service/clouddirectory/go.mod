module github.com/aws/aws-sdk-go-v2/service/clouddirectory

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.13.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.4
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.2.0
	github.com/aws/smithy-go v1.10.1-0.20220208165225-5adb4b73ede9
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
