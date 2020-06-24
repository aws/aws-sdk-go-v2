package changes

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

// changeTypes maps valid Change Types to the header they are grouped under in CHANGELOGs.
var changeHeaders = map[string]string{
	"feature": "New Features",
	"bugfix":  "Bug Fixes",
}

type changeTemplate struct {
	Modules     []string `yaml:",flow"`
	Type        string
	Description string
}

// Change represents a change to a single Go module.
type Change struct {
	ID            string // ID is a unique identifier for this Change
	SchemaVersion string // SchemaVersion is the version of the library's types used to create this Change.
	Module        string // Module is the Go module affected by this Change.
	Type          string // Type indicates what category of Change was made. Type may be either "feature" or "bugfix".
	Description   string // Description is a human readable description of this Change meant to be included in a CHANGELOG.
}

// NewChanges returns a Change slice containing a Change with the given type and description for each of the specified
// modules.
func NewChanges(modules []string, changeType, description string) ([]*Change, error) {
	if len(modules) == 0 || changeType == "" || description == "" {
		return nil, errors.New("missing module, type, or description")
	}

	changeType = strings.ToLower(changeType)
	if _, ok := changeHeaders[changeType]; !ok {
		return nil, fmt.Errorf("change type %s is not valid", changeType)
	}

	changes := make([]*Change, 0, len(modules))

	for _, module := range modules {
		changes = append(changes, &Change{
			ID:          generateId(module, changeType),
			Module:      module,
			Type:        changeType,
			Description: description,
		})
	}

	return changes, nil
}

// TemplateToChanges parses the provided filledTemplate into the provided Change. If Change has no ID, TemplateToChange
// will set the ID.
func TemplateToChanges(filledTemplate string) ([]*Change, error) {
	var template changeTemplate

	err := yaml.Unmarshal([]byte(filledTemplate), &template)
	if err != nil {
		return nil, err
	}

	return NewChanges(template.Modules, template.Type, template.Description)
}

// ChangeToTemplate returns a Change template populated with the given Change's data.
func ChangeToTemplate(change *Change) (string, error) {
	templateBytes, err := yaml.Marshal(changeTemplate{
		Modules:     []string{change.Module},
		Type:        change.Type,
		Description: change.Description,
	})
	if err != nil {
		return "", err
	}

	return string(templateBytes), nil
}

func generateId(module, changeType string) string {
	module = strings.ReplaceAll(module, "/", ".")

	return fmt.Sprintf("%s-%s-%v", module, changeType, time.Now().UnixNano())
}
