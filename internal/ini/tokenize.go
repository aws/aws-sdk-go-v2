package ini

import (
	"fmt"
	"strings"
)

func tokenize(lines []string) ([]LineToken, error) {
	tokens := make([]LineToken, 0, len(lines))
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 || isLineComment(line) {
			continue
		}

		if tok := asProfile(line); tok != nil {
			tokens = append(tokens, tok)
		} else if tok := asProperty(line); tok != nil {
			tokens = append(tokens, tok)
		} else if tok := asSubProperty(line); tok != nil {
			tokens = append(tokens, tok)
		} else if tok := asContinuation(line); tok != nil {
			tokens = append(tokens, tok)
		} else {
			return nil, fmt.Errorf("unrecognized token '%s'", line)
		}
	}
	return tokens, nil
}

func isLineComment(line string) bool {
	trimmed := strings.TrimLeft(line, " \t")
	return strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, ";")
}

func asProfile(line string) *LineTokenProfile { // "[ type name ] ; comment"
	trimmed := strings.TrimRight(trimComment(line), " \t") // "[ type name ]"
	if !strings.HasPrefix(trimmed, "[") || !strings.HasSuffix(trimmed, "]") {
		return nil
	}
	trimmed = trimmed[1 : len(trimmed)-1] // " type name " (or just " name ")
	trimmed = strings.TrimSpace(trimmed)  // "type name" / "name"
	typ, name := splitProfile(trimmed)
	return &LineTokenProfile{
		Type: typ,
		Name: name,
	}
}

func asProperty(line string) *LineTokenProperty {
	if isLineSpace(rune(line[0])) {
		return nil
	}

	trimmed := strings.TrimRight(trimComment(line), " \t")
	k, v, ok := splitProperty(trimmed)
	if !ok {
		return nil
	}

	return &LineTokenProperty{
		Key:   strings.ToLower(k), // LEGACY: normalize key case
		Value: legacyStrconv(v),   // LEGACY: see func docs
	}
}

func asSubProperty(line string) *LineTokenSubProperty {
	if !isLineSpace(rune(line[0])) {
		return nil
	}

	// comments on sub-properties are included in the value
	trimmed := strings.TrimLeft(line, " \t")
	k, v, ok := splitProperty(trimmed)
	if !ok {
		return nil
	}

	return &LineTokenSubProperty{ // same LEGACY constraints as in normal property
		Key:   strings.ToLower(k),
		Value: legacyStrconv(v),
	}
}

func asContinuation(line string) *LineTokenContinuation {
	if !isLineSpace(rune(line[0])) {
		return nil
	}

	// includes comments like sub-properties
	trimmed := strings.TrimLeft(line, " \t")
	return &LineTokenContinuation{
		Value: trimmed,
	}
}
