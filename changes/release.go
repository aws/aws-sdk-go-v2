package changes

import (
	"fmt"
	"sort"
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
	var sectionsList []string

	for _, c := range r.Changes {
		if c.Module == module {
			if _, ok := sections[c.Type]; !ok {
				sectionsList = append(sectionsList, c.Type)
			}

			sections[c.Type] += fmt.Sprintf("* %s\n", c.Description)
		}
	}

	entry := fmt.Sprintf("%s# `%s`", headerPrefix, module)
	if bump, ok := r.VersionBumps[module]; ok {
		entry += " - " + bump.To + "\n"
	} else {
		entry += "\n"
	}

	sort.Strings(sectionsList)

	for _, section := range sectionsList {
		entry += headerPrefix + "## " + changeHeaders[section] + "\n"
		entry += sections[section] + "\n"
	}

	return entry
}

// AffectedModules returns a sorted list of all modules affected by this Release. A module is considered affected if
// it is the Module of one or more Changes in the Release.
func (r *Release) AffectedModules() []string {
	var modules []string
	seen := make(map[string]struct{})

	for _, c := range r.Changes {
		if _, ok := seen[c.Module]; !ok {
			seen[c.Module] = struct{}{}
			modules = append(modules, c.Module)
		}
	}

	sort.Strings(modules)
	return modules
}
