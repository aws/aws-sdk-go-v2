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

// Change represents a change to one or more Go modules.
type Change struct {
	Id            string
	SchemaVersion string
	Modules       []string
	Type          string
	Description   string
	SetVersion    string
}

func NewChange(modules []string, changeType, description string) *Change {
	return &Change{
		Id:          generateId(modules[0], changeType),
		Modules:     modules,
		Type:        changeType,
		Description: description,
	}
}

// AffectsModule returns whether the Change contains the given module in its Modules. AffectsModule does not resolve
// wildcards.
func (c *Change) AffectsModule(module string) bool {
	for _, m := range c.Modules {
		if m == module {
			return true
		}
	}

	return false
}

func (c *Change) SetSchemaVersion(version string) {
	c.SchemaVersion = version
}

func (c *Change) isValid() bool {
	return len(c.Modules) > 0 && c.Type != "" && c.Description != ""
}

// TemplateToChange parses the provided filledTemplate into the provided Change. If Change has no Id, TemplateToChange
// will set the Id.
func TemplateToChange(filledTemplate string, change *Change) error {
	if change == nil {
		return errors.New("change must be non-nil")
	}

	lines := strings.Split(filledTemplate, "\n")

	for _, l := range lines {
		if l != "" && !strings.HasPrefix(l, "#") {
			parts := strings.Split(l, ": ")
			if len(parts) != 2 {
				return fmt.Errorf("template is incorrectly formatted at line: %s", l)
			}

			switch parts[0] {
			case "modules":
				trimmedMods := strings.Trim(parts[1], "[] ")
				change.Modules = strings.Split(trimmedMods, ",")
			case "change type":
				change.Type = parts[1]
			case "description":
				change.Description = parts[1]
			default:
				return fmt.Errorf("unknown template field: %s", parts[0])
			}
		}
	}

	if !change.isValid() {
		return errors.New("change template is missing a type, module, or description")
	}

	if change.Id == "" {
		change.Id = generateId(change.Modules[0], change.Type)
	}

	return nil
}

// ChangeToTemplate returns a Change template populated with the given Change's data.
func ChangeToTemplate(change *Change) string {
	modules := strings.Join(change.Modules, ", ")
	return fmt.Sprintf(ChangeTemplate, modules, change.Type, change.Description)
}

func generateId(module, changeType string) string {
	module = strings.ReplaceAll(module, "/", ".")

	id := fmt.Sprintf("%s-%s-%v", module, changeType, time.Now().UnixNano())

	return id
}
