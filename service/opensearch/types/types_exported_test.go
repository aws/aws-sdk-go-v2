// Code generated by smithy-go-codegen DO NOT EDIT.

package types_test

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/opensearch/types"
)

func ExampleDataSourceType_outputUsage() {
	var union types.DataSourceType
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.DataSourceTypeMemberS3GlueDataCatalog:
		_ = v.Value // Value is types.S3GlueDataCatalog

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.S3GlueDataCatalog

func ExampleDirectQueryDataSourceType_outputUsage() {
	var union types.DirectQueryDataSourceType
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.DirectQueryDataSourceTypeMemberCloudWatchLog:
		_ = v.Value // Value is types.CloudWatchDirectQueryDataSource

	case *types.DirectQueryDataSourceTypeMemberSecurityLake:
		_ = v.Value // Value is types.SecurityLakeDirectQueryDataSource

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.CloudWatchDirectQueryDataSource
var _ *types.SecurityLakeDirectQueryDataSource
