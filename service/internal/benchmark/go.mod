module github.com/Enflick/aws-sdk-go-v2/service/internal/benchmark

go 1.20

require (
	github.com/Enflick/aws-sdk-go-v2 v1.3.0
	github.com/Enflick/aws-sdk-go-v2/service/dynamodb v0.0.0-20250325221836-b4195dfa2eb5
	github.com/Enflick/aws-sdk-go-v2/service/lexruntimeservice v1.20.8
	github.com/Enflick/aws-sdk-go-v2/service/s3 v1.54.3
	github.com/Enflick/aws-sdk-go-v2/service/schemas v1.24.8
	github.com/Enflick/smithy-go v1.3.0
	github.com/aws/aws-sdk-go v1.44.28
)

require (
	github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/configsources v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 v2.0.0-20250325155711-0a4bf6fdbeb3 // indirect
	github.com/Enflick/aws-sdk-go-v2/internal/v4a v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/checksum v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/endpoint-discovery v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/Enflick/aws-sdk-go-v2/service/internal/s3shared v0.0.0-20250325221836-b4195dfa2eb5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace github.com/Enflick/aws-sdk-go-v2 => ../../../

replace github.com/Enflick/aws-sdk-go-v2/aws/protocol/eventstream => ../../../aws/protocol/eventstream/

replace github.com/Enflick/aws-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/Enflick/aws-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/Enflick/aws-sdk-go-v2/internal/v4a => ../../../internal/v4a/

replace github.com/Enflick/aws-sdk-go-v2/service/dynamodb => ../../../service/dynamodb/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/checksum => ../../../service/internal/checksum/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/endpoint-discovery => ../../../service/internal/endpoint-discovery/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/

replace github.com/Enflick/aws-sdk-go-v2/service/internal/s3shared => ../../../service/internal/s3shared/

replace github.com/Enflick/aws-sdk-go-v2/service/lexruntimeservice => ../../../service/lexruntimeservice/

replace github.com/Enflick/aws-sdk-go-v2/service/s3 => ../../../service/s3/

replace github.com/Enflick/aws-sdk-go-v2/service/schemas => ../../../service/schemas/
