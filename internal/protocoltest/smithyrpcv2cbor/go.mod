module github.com/aws/aws-sdk-go-v2/internal/protocoltest/smithyrpcv2cbor

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.25.2
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.2
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.2
	github.com/aws/smithy-go v1.20.1
)

replace github.com/aws/aws-sdk-go-v2 => ../../../

replace github.com/aws/aws-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/
