package ini

import (
	"fmt"
	"strings"
)

func parse(tokens []lineToken, path string) (Sections, error) {
	cfg := Sections{
		container: map[string]Section{},
	}

	var currSection, currKey string

	for _, otok := range tokens {
		switch tok := otok.(type) {
		case *lineTokenProfile:
			name := tok.Name
			if tok.Type != "" {
				name = fmt.Sprintf("%s %s", tok.Type, tok.Name)
			}
			currKey = ""
			currSection = name
			if _, ok := cfg.container[name]; !ok {
				cfg.container[name] = NewSection(name)
			}
		case *lineTokenProperty:
			handleProperty(&currSection, &currKey, &cfg, path, tok)
		case *lineTokenSubProperty:
			if currSection == "" {
				continue
			}

			if currKey == "" || cfg.container[currSection].values[currKey].str != "" {
				// This is an "orphaned" subproperty, either because it's at
				// the beginning of a section or because the last property's
				// value isn't empty. Either way we're lenient here and
				// "promote" this to a normal property.
				handleProperty(&currSection, &currKey, &cfg, path, &lineTokenProperty{
					Key:   tok.Key,
					Value: strings.TrimSpace(trimComment(tok.Value)),
				})
				continue
			}

			if cfg.container[currSection].values[currKey].mp == nil {
				cfg.container[currSection].values[currKey] = Value{
					mp: map[string]string{},
				}
			}
			cfg.container[currSection].values[currKey].mp[tok.Key] = tok.Value
		case *lineTokenContinuation:
			if currKey == "" {
				continue
			}

			value, _ := cfg.container[currSection].values[currKey]
			if value.str != "" && value.mp == nil {
				value.str = fmt.Sprintf("%s\n%s", value.str, tok.Value)
			}

			cfg.container[currSection].values[currKey] = value
		}
	}
	return cfg, nil
}

func handleProperty(currSection, currKey *string, cfg *Sections, path string, tok *lineTokenProperty) {
	if *currSection == "" {
		return // LEGACY: don't error on "global" properties
	}

	*currKey = tok.Key
	if _, ok := cfg.container[*currSection].values[tok.Key]; ok {
		section := cfg.container[*currSection]
		section.Logs = append(cfg.container[*currSection].Logs,
			fmt.Sprintf(
				"For profile: %v, overriding %v value, with a %v value found in a duplicate profile defined later in the same file %v. \n",
				*currSection, tok.Key, tok.Key, path,
			),
		)
		cfg.container[*currSection] = section
	}

	cfg.container[*currSection].values[tok.Key] = Value{
		str: tok.Value,
	}
	cfg.container[*currSection].SourceFile[tok.Key] = path
}
