package changes

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ChangeTemplate is an editor friendly template for creating or modifying a Change.
const ChangeTemplate = `modules: %s
# change type is one of: feature, bugfix
change type: %s
description: %s`

// changeTypes maps valid Change Types to the header they are grouped under in CHANGELOGs.
var changeHeaders = map[string]string{
	"feature": "New Features",
	"bugfix":  "Bug Fixes",
}

// Change represents a change to a single Go module.
type Change struct {
	ID            string // ID is a unique identifier for this Change
	SchemaVersion string // SchemaVersion is the version of the library's types used to create this Change.
	Module        string // Module is the Go module affected by this Change.
	Type          string // Type indicates what category of Change was made. Type may be either "feature" or "bugfix".
	Description   string // Description is a human readable description of this Change meant to be included in a CHANGELOG.
}

// NewChanges
func NewChanges(modules []string, changeType, description string) ([]*Change, error) {
	if len(modules) == 0 || changeType == "" || description == "" {
		return nil, errors.New("missing module, type, or description")
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
	lines := strings.Split(filledTemplate, "\n")

	var modules []string
	var changeType string
	var description string

	for _, l := range lines {
		if l != "" && !strings.HasPrefix(l, "#") {
			parts := strings.Split(l, ": ")
			if len(parts) != 2 {
				return nil, fmt.Errorf("template is incorrectly formatted at line: %s", l)
			}

			switch parts[0] {
			case "modules":
				trimmedMods := strings.Trim(parts[1], "[] ")
				modules = strings.Split(trimmedMods, ",")
			case "change type":
				parts[1] = strings.ToLower(parts[1])
				if changeHeaders[parts[1]] == "" {
					return nil, fmt.Errorf("%s is not a valid change type", parts[1])
				}

				changeType = parts[1]
			case "description":
				description = parts[1]
			default:
				return nil, fmt.Errorf("unknown template field: %s", parts[0])
			}
		}
	}

	return NewChanges(modules, changeType, description)
}

// ChangeToTemplate returns a Change template populated with the given Change's data.
func ChangeToTemplate(change *Change) string {
	return fmt.Sprintf(ChangeTemplate, change.Module, change.Type, change.Description)
}

func generateId(module, changeType string) string {
	module = strings.ReplaceAll(module, "/", ".")

	return fmt.Sprintf("%s-%s-%v", module, changeType, time.Now().UnixNano())
}
