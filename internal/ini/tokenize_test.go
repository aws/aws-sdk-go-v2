package ini

import (
	"testing"
)

func TestTrimComment(t *testing.T) {
	cs := []struct {
		Name   string
		Input  string
		Expect string
	}{
		{
			Name:   "hash",
			Input:  "foo # comment",
			Expect: "foo ",
		},
		{
			Name:   "semi",
			Input:  "foo ; comment",
			Expect: "foo ",
		},
		{
			Name:   "nested",
			Input:  "foo ;# comment",
			Expect: "foo ",
		},
		{
			Name:   "nop",
			Input:  "foo",
			Expect: "foo",
		},
	}
	for _, c := range cs {
		t.Run(c.Name, func(t *testing.T) {
			if a := trimComment(c.Input); c.Expect != a {
				t.Fatalf("%v != %v", c.Expect, a)
			}
		})
	}
}

func TestAsProfile(t *testing.T) {
	cs := []struct {
		Name   string
		Input  string
		Expect *LineTokenProfile
	}{
		{
			Name:  "default",
			Input: "[default]",
			Expect: &LineTokenProfile{
				Name: "default",
			},
		},
		{
			Name:  "profile",
			Input: "[profile foo]",
			Expect: &LineTokenProfile{
				Type: "profile",
				Name: "foo",
			},
		},
		{
			Name:  "rigorous",
			Input: "[ 	profile foo	 ] ; comment",
			Expect: &LineTokenProfile{
				Type: "profile",
				Name: "foo",
			},
		},
		{
			Name:  "rigorous, no type",
			Input: "[ 	 bazzle ] ; comment",
			Expect: &LineTokenProfile{
				Type: "",
				Name: "bazzle",
			},
		},
	}
	for _, c := range cs {
		t.Run(c.Name, func(t *testing.T) {
			a := asProfile(c.Input)
			if a == nil && c.Expect != nil {
				t.Fatalf("shouldn't be nil")
			} else if c.Expect == nil && a != nil {
				t.Fatalf("should be nil")
			}

			if *c.Expect != *a {
				t.Fatalf("%v != %v", *c.Expect, *a)
			}
		})
	}
}
