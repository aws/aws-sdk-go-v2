module github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.30.1
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.13
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.13
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.13
	github.com/aws/smithy-go v1.20.3
)

replace github.com/aws/aws-sdk-go-v2 => ../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../internal/endpoints/v2/

replace github.com/aws/aws-sdk-go-v2/internal/v4a => ../../internal/v4a/
