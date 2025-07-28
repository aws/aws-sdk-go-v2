package kitchensinktest

import (
	"reflect"
	"testing"

	smithyauth "github.com/aws/smithy-go/auth"
)

func toOptions(ids []string) []*smithyauth.Option {
	var options []*smithyauth.Option
	for _, scheme := range ids {
		options = append(options, &smithyauth.Option{
			SchemeID: scheme,
		})
	}
	return options
}

func TestSortAuthOptions(t *testing.T) {
	for name, tt := range map[string]struct {
		ResolvedIDs []string
		Preference  []string
		Expect      []string
	}{
		"no preference": {
			ResolvedIDs: []string{smithyauth.SchemeIDSigV4, smithyauth.SchemeIDSigV4A},
			Preference:  []string{},
			Expect:      []string{smithyauth.SchemeIDSigV4, smithyauth.SchemeIDSigV4A},
		},
		"prefer sigv4a": {
			ResolvedIDs: []string{smithyauth.SchemeIDSigV4, smithyauth.SchemeIDSigV4A},
			Preference:  []string{"sigv4a"},
			Expect:      []string{smithyauth.SchemeIDSigV4A, smithyauth.SchemeIDSigV4},
		},
		"prefer sigv4a,sigv4": {
			ResolvedIDs: []string{smithyauth.SchemeIDSigV4, smithyauth.SchemeIDSigV4A},
			Preference:  []string{"sigv4a", "sigv4"},
			Expect:      []string{smithyauth.SchemeIDSigV4A, smithyauth.SchemeIDSigV4},
		},
		"prefer sigv4": {
			ResolvedIDs: []string{smithyauth.SchemeIDSigV4, smithyauth.SchemeIDSigV4A},
			Preference:  []string{"sigv4"},
			Expect:      []string{smithyauth.SchemeIDSigV4, smithyauth.SchemeIDSigV4A},
		},
		"prefer sigv4,sigv4a": {
			ResolvedIDs: []string{smithyauth.SchemeIDSigV4, smithyauth.SchemeIDSigV4A},
			Preference:  []string{"sigv4", "sigv4a"},
			Expect:      []string{smithyauth.SchemeIDSigV4, smithyauth.SchemeIDSigV4A},
		},
		"prefer nonsense": {
			ResolvedIDs: []string{smithyauth.SchemeIDSigV4, smithyauth.SchemeIDSigV4A},
			Preference:  []string{"nonsense"},
			Expect:      []string{smithyauth.SchemeIDSigV4, smithyauth.SchemeIDSigV4A},
		},
		"preserve remaining order": {
			ResolvedIDs: []string{"a", "b", "c", "d", "e", "f"},
			Preference:  []string{"f", "c", "a"},
			Expect:      []string{"f", "c", "a", "b", "d", "e"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			sorted := sortAuthOptions(toOptions(tt.ResolvedIDs), tt.Preference)
			var actual []string
			for _, option := range sorted {
				actual = append(actual, option.SchemeID)
			}

			if !reflect.DeepEqual(tt.Expect, actual) {
				t.Errorf("resolve auth scheme:\n%#v !=\n%#v", tt.Expect, actual)
			}
		})
	}
}
