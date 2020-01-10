// +build codegen

package api

import (
	"fmt"
)

// ExamplesBuilder provides the logic to build modeled examples as Go code.
type ExamplesBuilder struct {
	ShapeValueBuilder

	HasTimestamp bool
}

// NewExamplesBuilder returns an initialized example builder for generating
// example input API shapes from a model.
func NewExamplesBuilder() *ExamplesBuilder {
	b := &ExamplesBuilder{
		ShapeValueBuilder: NewShapeValueBuilder(),
	}
	b.ParseTimeString = b.parseExampleTimeString
	return b
}

// Returns a string which assigns the value of a time member by calling
// parseTime function defined in the file.
func (b *ExamplesBuilder) parseExampleTimeString(ref *ShapeRef, v string) (string, error) {
	b.HasTimestamp = true

	if ref.Location == "header" {
		return fmt.Sprintf("parseTime(%q,%q)", "Mon, 2 Jan 2006 15:04:05 GMT", v), nil
	}

	switch ref.API.Metadata.Protocol {
	case "json", "rest-json", "rest-xml", "ec2", "query":
		return fmt.Sprintf("parseTime(%q,%q)", "2006-01-02T15:04:05Z", v), nil

	default:
		return "", fmt.Errorf("Unsupported time type: %s", ref.API.Metadata.Protocol)
	}
}
