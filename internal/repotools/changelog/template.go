package changelog

import (
	"bytes"
	"fmt"

	"github.com/pelletier/go-toml"
)

// TemplateToAnnotation parses the provided filledTemplate into the provided Change. If Change has no ID, TemplateToChange
// will set the ID.
func TemplateToAnnotation(filledTemplate []byte) (annotation Annotation, err error) {
	err = toml.Unmarshal(filledTemplate, &annotation)
	if err != nil {
		return Annotation{}, err
	}

	if len(annotation.Modules) == 0 {
		return Annotation{}, fmt.Errorf("annotation should include at least one module")
	}

	return annotation, nil
}

// AnnotationToTemplate returns a Change template populated with the given Change's data.
func AnnotationToTemplate(annotation Annotation) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	encoder := toml.NewEncoder(buffer)

	err := encoder.Order(toml.OrderPreserve).
		ArraysWithOneElementPerLine(true).
		Indentation("    ").
		Encode(annotation)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
