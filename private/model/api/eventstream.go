// +build codegen

package api

import (
	"fmt"
	"os"
)

func (a *API) suppressEventStreams() error {
	for opName, op := range a.Operations {
		inputRef := getEventStreamMember(op.InputRef.Shape)
		outputRef := getEventStreamMember(op.OutputRef.Shape)

		if inputRef == nil && outputRef == nil {
			continue
		}

		if !a.KeepUnsupportedAPIs {
			fmt.Fprintf(os.Stderr,
				"removing unsupported eventstream operation, %s\n",
				opName)
			a.removeOperation(opName)
			continue
		}
		return UnsupportedAPIModelError{
			Err: fmt.Errorf("eventstream support not implemented, %s",
				op.ExportedName),
		}
	}

	return nil
}

func getEventStreamMember(topShape *Shape) *ShapeRef {
	for _, ref := range topShape.MemberRefs {
		if !ref.Shape.IsEventStream {
			continue
		}
		return ref
	}

	return nil
}
