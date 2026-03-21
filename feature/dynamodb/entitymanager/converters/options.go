package converters

import (
	"fmt"
	"strings"
)

func getOpt(opts []string, name string) string {
	p := fmt.Sprintf("%s=", name)

	for _, opt := range opts {
		if strings.HasPrefix(opt, p) {
			return opt[len(p):]
		}
	}

	return ""
}
