package main

import (
	"fmt"
	"strings"
	"text/template"
)

var repoChangeLogTemplate = template.Must(template.New("repoChangeLog").
	Funcs(map[string]interface{}{
		"inlineCodeBlock": inlineCodeBlock,
		"modulesForHighlight": func(v map[string]moduleSummary) (f map[string]moduleSummary) {
			f = make(map[string]moduleSummary)
			for modDir, summary := range v {
				for i := range summary.Annotations {
					if !summary.Annotations[i].Collapse {
						f[modDir] = summary
					}
				}
			}
			return f
		},
		"changeLogLink": func(relModDir, version string, releaseID string) string {
			if relModDir == "." {
				return version
			}

			lv := strings.ReplaceAll(version, ".", "")
			lr := strings.ReplaceAll(releaseID, ".", "")

			return fmt.Sprintf("[%s](%s/CHANGELOG.md#%s)", version, relModDir, strings.ToLower(lv+"-"+lr))
		},
	}).
	Parse(`{{ define "entry" -}}
# Release ({{ .ReleaseID }})

{{ template "summary" . -}}
{{ end }}{{/* template */}}
{{ define "summary" -}}
{{ if (not .IsEmptyReleaseSummary) -}}
{{ if (gt (len .General) 0) -}}
## General Highlights
{{ range $_, $a := .General -}}
* **{{ $a.Type.ChangelogPrefix }}**: {{ $a.Description }}
{{ end }}{{/* range */}}
{{ end -}}{{/* if */ -}}
{{ $mh := (modulesForHighlight .Modules) -}}
{{ if (gt (len $mh) 0) -}}
## Module Highlights
{{ range $name, $mod := $mh -}}
* {{ inlineCodeBlock $mod.ModulePath }}: {{ changeLogLink $name $mod.Version $mod.ReleaseID }}
{{ range $_, $a := $mod.Annotations -}}
{{- if (not $a.Collapse) }}  * **{{ $a.Type.ChangelogPrefix }}**: {{ $a.Description }}
{{ end -}}{{/* if */ -}}
{{ end -}}{{/* range */ -}}
{{ end }}{{/* range */}}
{{ end -}}{{/* if */ -}}
{{ else -}}{{/* if */ -}}
* No change notes available for this release.

{{ end -}}{{/* if * */ -}}
{{ end -}}{{/* define */ -}}
`))

var moduleChangeLogTemplate = template.Must(template.New("moduleChangeLog").
	Funcs(map[string]interface{}{
		"inlineCodeBlock": inlineCodeBlock,
	}).
	Parse(`# {{ .Version }} ({{ .ReleaseID }})

{{ if (len .Annotations) -}}
{{- range $_, $a := .Annotations -}}
* **{{ $a.Type.ChangelogPrefix }}**: {{ $a.Description }}
{{ end -}}
{{- else -}}
* No change notes available for this release.
{{ end }}
`))
