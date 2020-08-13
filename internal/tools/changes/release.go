package changes

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Release represents a single SDK release, which contains all change metadata and their resulting version bumps.
type Release struct {
	ID            string
	SchemaVersion int
	VersionBumps  map[string]VersionBump
	Changes       []Change
}

type changelogModuleEntry struct {
	Module    string
	Version   string
	Sections  map[ChangeType][]Change
	TopLevel  bool
	ReleaseID string
}

func (e changelogModuleEntry) Link() string {
	anchor := "Release-" + strings.ReplaceAll(e.ReleaseID, " ", "-")
	return fmt.Sprintf("[%s](%s/CHANGELOG.md#%s)", e.Module, e.Module, anchor)
}

const changelogModule = `{{- if .TopLevel -}}
* {{.Link}}{{with .Version}} - {{.}}{{end}}{{else -}}
## Release {{.ReleaseID}}
* ` + "`" + `{{.Module}}` + "`" + `{{with .Version}} - {{.}}{{end}}{{end -}}
{{range $key, $section := .Sections}}{{range $section}}
  * {{ $key.ChangelogPrefix }}{{.IndentedDescription "  "}}{{end}}{{end}}{{if not .TopLevel}}
{{end}}`

var changelogTemplate *template.Template
var rootChangelogTemplate *template.Template

func init() {
	var err error

	changelogTemplate, err = template.New("changelog-entry").Parse(changelogModule)
	if err != nil {
		panic(err)
	}

	rootChangelogTemplate, err = template.New("root-changelog").Parse(rootChangelogTemplateContents)
	if err != nil {
		panic(err)
	}
}

// RenderChangelogForModule returns a new markdown section of a module's CHANGELOG based on the Changes in the Release.
func (r *Release) RenderChangelogForModule(module string, topLevel bool) (string, error) {
	sections := map[ChangeType][]Change{}

	for _, c := range r.Changes {
		if topLevel && c.Module == module {
			sections[c.Type] = append(sections[c.Type], c)
		} else if !topLevel && c.matches(module) {
			sections[c.Type] = append(sections[c.Type], c)
		}
	}

	if len(sections) == 0 {
		return "", nil
	}

	var version string
	if bump, ok := r.VersionBumps[module]; ok {
		version = bump.To
	}

	buff := new(bytes.Buffer)

	err := changelogTemplate.Execute(buff, changelogModuleEntry{
		Module:    module,
		Version:   version,
		Sections:  sections,
		ReleaseID: r.ID,
		TopLevel:  topLevel,
	})
	if err != nil {
		return "", fmt.Errorf("failed to render module %s's changelog entry: %v", module, err)
	}

	return buff.String(), nil
}

const rootChangelogTemplateContents = `# Release {{.ID}}
{{with .AnnouncementsSection}}## Announcements
{{range .}}{{.}}
{{end -}}
{{end}}{{with .ServiceSection}}## Service Client Highlights
{{range .}}{{.}}
{{end -}}
{{end}}{{with .CoreSection}}## Core SDK Highlights
{{range .}}{{.}}
{{end -}}
{{end}}`

// RenderChangelog generates a top level CHANGELOG.md for the Release r.
func (r *Release) RenderChangelog() (string, error) {
	buff := new(bytes.Buffer)

	err := rootChangelogTemplate.Execute(buff, r)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

// AffectedModules returns a sorted list of all modules affected by this Release. A module is considered affected if
// it is the Module of one or more Changes in the Release.
func (r *Release) AffectedModules() []string {
	return AffectedModules(r.Changes)
}

// wildcards returns a sorted list of wildcards Changes whose Module begin with the given prefix.
func (r *Release) wildcards() []Change {
	var changes []Change

	for _, c := range r.Changes {
		if c.isWildcard() {
			changes = append(changes, c)
		}
	}

	return changes
}

// splitSections groups entries (including wildcard Changes and module Changelog entries) into three groups: Announcements,
// Services, and Core SDK modules.
func (r *Release) splitSections() ([]string, []string, []string, error) {
	const servicePrefix = "service/"

	var announcements []string
	var services []string
	var core []string

	for _, c := range r.wildcards() {
		if c.Type == AnnouncementChangeType {
			announcements = append(announcements, c.String())
		} else if strings.HasPrefix(c.Module, servicePrefix) {
			services = append(services, c.String())
		} else {
			core = append(core, c.String())
		}
	}

	mods := r.AffectedModules()

	for _, m := range mods {
		entry, err := r.RenderChangelogForModule(m, true)
		if err != nil {
			return nil, nil, nil, err
		}
		if entry == "" {
			continue
		}

		if strings.HasPrefix(m, servicePrefix) {
			services = append(services, entry)
		} else {
			core = append(core, entry)
		}
	}

	return announcements, services, core, nil
}

// AnnouncementsSection returns a list of Changelog bullet entries that should be included under the Announcements header.
func (r *Release) AnnouncementsSection() ([]string, error) {
	announcements, _, _, err := r.splitSections()
	return announcements, err
}

// ServiceSection returns a list of Changelog bullet entries that should be included under the Service Clients header.
func (r *Release) ServiceSection() ([]string, error) {
	_, services, _, err := r.splitSections()
	return services, err
}

// CoreSection returns a list of Changelog bullet entries that should be included under the Core SDK header.
func (r *Release) CoreSection() ([]string, error) {
	_, _, core, err := r.splitSections()
	return core, err
}
