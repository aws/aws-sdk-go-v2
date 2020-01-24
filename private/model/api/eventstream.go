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

		err := fmt.Errorf("eventstream support not implemented, %s, %s",
			op.ExportedName)

		if a.IgnoreUnsupportedAPIs {
			fmt.Fprintf(os.Stderr, "removing operation, %s, %v\n", opName, err)
			delete(a.Operations, opName)
			continue
		}
		return UnsupportedAPIModelError{
			Err: err,
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
