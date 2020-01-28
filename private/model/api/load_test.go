// +build codegen

package api

import (
	"testing"
)

func TestResolvedReferences(t *testing.T) {
	json := `{
		"operations": {
			"OperationName": {
				"input": { "shape": "TestName" }
			}
		},
		"shapes": {
			"TestName": {
				"type": "structure",
				"members": {
					"memberName1": { "shape": "OtherTest" },
					"memberName2": { "shape": "OtherTest" }
				}
			},
			"OtherTest": { "type": "string" }
		}
	}`
	a := API{}
	err := a.AttachString(json)
	if err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}
	if len(a.Shapes["OtherTest"].refs) != 2 {
		t.Errorf("Expected %d, but received %d", 2, len(a.Shapes["OtherTest"].refs))
	}
}
