package ratelimit

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestTokenRateLimit(t *testing.T) {
	type usage struct {
		Cost      uint
		Release   bool
		Err       string
		AddTokens uint
	}

	cases := map[string]struct {
		Tokens uint
		Usages []usage
	}{
		"retrieve": {
			Tokens: 10,
			Usages: []usage{
				{Cost: 5, Release: true},
				{Cost: 5},
				{Cost: 5},
				{Cost: 5, Err: "retry quota exceeded"},
				{AddTokens: 5},
				{Cost: 5},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			rl := NewTokenRateLimit(c.Tokens)

			for i, u := range c.Usages {
				t.Run(fmt.Sprintf("usage_%d", i), func(t *testing.T) {
					if u.Cost != 0 {
						f, err := rl.GetToken(context.Background(), u.Cost)
						if len(u.Err) != 0 {
							if err == nil {
								t.Fatalf("expect error, got none")
							}
							if e, a := u.Err, err.Error(); !strings.Contains(a, e) {
								t.Fatalf("expect %q error, got %q", e, a)
							}
						} else if err != nil {
							t.Fatalf("expect no error, got %v", err)
						}

						if u.Release {
							if err := f(); err != nil {
								t.Fatalf("expect no error, got %v", err)
							}
						}
					}

					if u.AddTokens != 0 {
						rl.AddTokens(u.AddTokens)
					}
				})
			}
		})
	}
}
