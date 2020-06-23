package changes

import (
	"fmt"
)

// VersionBump describes a version increment to a module.
type VersionBump struct {
	From string
	To   string
}

// Release represents a single SDK release, which contains all change metadata and their resulting version bumps.
type Release struct {
	ID            string
	SchemaVersion string
	VersionBumps  map[string]VersionBump
	Changes       []*Change
}

// RenderChangelogForModule returns a new markdown section of a module's CHANGELOG based on the Changes in the Release.
func (r *Release) RenderChangelogForModule(module, headerPrefix string) string {
	sections := map[string]string{}

	for _, c := range r.Changes {
		if c.Module == module {
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
