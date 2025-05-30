// Code generated by smithy-go-codegen DO NOT EDIT.

package types_test

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func ExampleExecutionInputs_outputUsage() {
	var union types.ExecutionInputs
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.ExecutionInputsMemberAutomation:
		_ = v.Value // Value is types.AutomationExecutionInputs

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.AutomationExecutionInputs

func ExampleExecutionPreview_outputUsage() {
	var union types.ExecutionPreview
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.ExecutionPreviewMemberAutomation:
		_ = v.Value // Value is types.AutomationExecutionPreview

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.AutomationExecutionPreview

func ExampleNodeType_outputUsage() {
	var union types.NodeType
	// type switches can be used to check the union value
	switch v := union.(type) {
	case *types.NodeTypeMemberInstance:
		_ = v.Value // Value is types.InstanceInfo

	case *types.UnknownUnionMember:
		fmt.Println("unknown tag:", v.Tag)

	default:
		fmt.Println("union is nil or unknown type")

	}
}

var _ *types.InstanceInfo
