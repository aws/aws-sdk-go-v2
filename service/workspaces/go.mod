module github.com/Enflick/aws-sdk-go-v2/service/workspaces

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.27.0
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.7
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.7
	github.com/aws/smithy-go v1.20.2
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/
