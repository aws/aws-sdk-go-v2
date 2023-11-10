package ini

import (
	"fmt"
)

func parse(tokens []LineToken, path string) (Sections, error) {
	cfg := Sections{
		container: map[string]Section{},
	}

	var currSection, currKey string

	for _, otok := range tokens {
		switch tok := otok.(type) {
		case *LineTokenProfile:
			name := tok.Name
			if tok.Type != "" {
				name = fmt.Sprintf("%s %s", tok.Type, tok.Name)
			}
			currKey = ""
			currSection = name
			if _, ok := cfg.container[name]; !ok {
				cfg.container[name] = NewSection(name)
			}
		case *LineTokenProperty:
			if currSection == "" {
				continue // LEGACY: don't error on "global" properties
			}

			currKey = tok.Key
			if _, ok := cfg.container[currSection].values[tok.Key]; ok {
				section := cfg.container[currSection]
				section.Logs = append(cfg.container[currSection].Logs,
					fmt.Sprintf(
						"For profile: %v, overriding %v value, with a %v value found in a duplicate profile defined later in the same file %v. \n",
						currSection, tok.Key, tok.Key, path,
					),
				)
				cfg.container[currSection] = section
			}

			cfg.container[currSection].values[tok.Key] = Value{
				str: tok.Value,
			}
			cfg.container[currSection].SourceFile[tok.Key] = path

		case *LineTokenSubProperty:
			if currKey == "" {
				continue
			}

			value, ok := cfg.container[currSection].values[currKey]
			if !ok {
				return cfg, fmt.Errorf("something went wrong")
			}

			if value.mp == nil && value.str == "" {
				value.mp = map[string]string{
					tok.Key: tok.Value,
				}
			}

			cfg.container[currSection].values[currKey] = value
		case *LineTokenContinuation:
			if currKey == "" {
				continue
			}

			if isBracketed(tok.Value) {
				continue // LEGACY: ignore this for some reason
			}

			value, ok := cfg.container[currSection].values[currKey]
			if !ok {
				return cfg, fmt.Errorf("something went wrong")
			}

			if value.str != "" && value.mp == nil {
				value.str = fmt.Sprintf("%s\n%s", value.str, tok.Value)
			}

			cfg.container[currSection].values[currKey] = value
		}
	}
	return cfg, nil
}
