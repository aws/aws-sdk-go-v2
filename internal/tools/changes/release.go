package changes

import (
	"bytes"
	"fmt"
	"text/template"
)

// VersionBump describes a version increment to a module.
type VersionBump struct {
	From string
	To   string
}

// Release represents a single SDK release, which contains all change metadata and their resulting version bumps.
type Release struct {
	ID            string
	SchemaVersion int
	VersionBumps  map[string]VersionBump
	Changes       []Change
}

type changelogModuleEntry struct {
	Prefix   string
	Module   string
	Version  string
	Sections map[ChangeType][]Change
}

const changelogModule = `{{.Prefix}}# ` + "`" + `{{.Module}}` + "`" + `{{with .Version}} - {{.}}{{end}}
{{range $key, $section := .Sections}}{{$.Prefix}}## {{ $key.HeaderTitle }}
{{range $section}}* {{.Description}}
{{end}}
{{end}}`

var changelogTemplate *template.Template

func init() {
	var err error

	changelogTemplate, err = template.New("changelog-entry").Parse(changelogModule)
	if err != nil {
		panic(err)
	}
}

// RenderChangelogForModule returns a new markdown section of a module's CHANGELOG based on the Changes in the Release.
func (r *Release) RenderChangelogForModule(module, headerPrefix string) (string, error) {
	sections := map[ChangeType][]Change{}

	for _, c := range r.Changes {
		if c.Module == module {
			sections[c.Type] = append(sections[c.Type], c)
		}
	}

	var version string
	if bump, ok := r.VersionBumps[module]; ok {
		version = bump.To
	}

	buff := new(bytes.Buffer)

	err := changelogTemplate.Execute(buff, changelogModuleEntry{
		Prefix:   headerPrefix,
		Module:   module,
		Version:  version,
		Sections: sections,
	})
	if err != nil {
		return "", fmt.Errorf("failed to render module %s's changelog entry: %v", module, err)
	}

	return buff.String(), nil
}

// AffectedModules returns a sorted list of all modules affected by this Release. A module is considered affected if
// it is the Module of one or more Changes in the Release.
func (r *Release) AffectedModules() []string {
	return AffectedModules(r.Changes)
}
