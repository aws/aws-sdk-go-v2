package changes

import (
	"fmt"
)

type VersionBump struct {
	From string
	To   string
}

type Release struct {
	Id            string
	SchemaVersion string
	VersionBumps  map[string]VersionBump
	Changes       []*Change
}

// RenderChangelogForModule returns a new markdown section of a module's CHANGELOG based on the Changes in the Release.
func (r *Release) RenderChangelogForModule(module, headerPrefix string) string {
	sections := map[string]string{}

	for _, c := range r.Changes {
		if c.AffectsModule(module) {
			sections[c.Type] += fmt.Sprintf("* %s\n", c.Description)
		}
	}

	entry := fmt.Sprintf("%s# `%s` - %s\n", headerPrefix, module, r.VersionBumps[module].To)
	for section, content := range sections {
		entry += headerPrefix + "## " + changeHeaders[section] + "\n"
		entry += content + "\n"
	}

	return entry
}

func (r *Release) SetSchemaVersion(version string) {
	r.SchemaVersion = version
}
