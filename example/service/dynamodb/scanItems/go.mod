module github.com/aws/aws-sdk-go-v2/example/service/dynamodb/listItems

go 1.16

require (
	github.com/aws/aws-sdk-go-v2 v1.13.0
	github.com/aws/aws-sdk-go-v2/config v1.13.1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.13.0
)

replace github.com/aws/aws-sdk-go-v2/config => ../../../../config/

replace github.com/aws/aws-sdk-go-v2 => ../../../../

replace github.com/aws/aws-sdk-go-v2/service/dynamodb => ../../../../service/dynamodb/
