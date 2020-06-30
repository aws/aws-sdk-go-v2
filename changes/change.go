package changes

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

const FeatureType = "feature" // FeatureType is a constant change type for a new feature.
const BugFixType = "bugfix"   // BugFixType is a constant change type for a bug fix.

// changeTypes maps valid Change Types to the header they are grouped under in CHANGELOGs.
var changeHeaders = map[string]string{
	FeatureType: "New Features",
	BugFixType:  "Bug Fixes",
}

const changeTemplateSuffix = `
# type may be one of "feature" or "bugfix".
# multiple modules may be listed. A change metadata file will be created for each module.`

type changeTemplate struct {
	Modules     []string
	Type        string
	Description string
}

// Change represents a change to a single Go module.
type Change struct {
	ID            string // ID is a unique identifier for this Change.
	SchemaVersion string // SchemaVersion is the version of the library's types used to create this Change.
	Module        string // Module is a shortened Go module path for the module affected by this Change. Module is the path from the root of the repository to the module.
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
func TemplateToChanges(filledTemplate []byte) ([]*Change, error) {
	var template changeTemplate

	err := yaml.Unmarshal(filledTemplate, &template)
	if err != nil {
		return nil, err
	}

	return NewChanges(template.Modules, template.Type, template.Description)
}

// ChangeToTemplate returns a Change template populated with the given Change's data.
func ChangeToTemplate(change *Change) ([]byte, error) {
	templateBytes, err := yaml.Marshal(changeTemplate{
		Modules:     []string{change.Module},
		Type:        change.Type,
		Description: change.Description,
	})
	if err != nil {
		return nil, err
	}

	return append(templateBytes, []byte(changeTemplateSuffix)...), nil
}

func generateId(module, changeType string) string {
	module = strings.ReplaceAll(module, "/", ".")

	return fmt.Sprintf("%s-%s-%v", module, changeType, time.Now().UnixNano())
}
