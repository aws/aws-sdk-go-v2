package attributevalue

import (
	"fmt"
	"reflect"
	"testing"
)

type testUnionValues struct {
	Name  string
	Value interface{}
}

type unionSimple struct {
	A int
	B string
	C []string
}

type unionComplex struct {
	unionSimple
	A int
}

type unionTagged struct {
	A int `dynamodbav:"ddbav" json:"A" taga:"TagA" tagb:"TagB"`
}

type unionTaggedComplex struct {
	unionSimple
	unionTagged
	B string
}

func TestUnionStructFields(t *testing.T) {
	origFieldCache := fieldCache
	defer func() { fieldCache = origFieldCache }()

	fieldCache = &fieldCacher{}

	var cases = map[string]struct {
		in     interface{}
		opts   structFieldOptions
		expect []testUnionValues
	}{
		"simple input": {
			in:   unionSimple{1, "2", []string{"abc"}},
			opts: structFieldOptions{TagKey: "json"},
			expect: []testUnionValues{
				{"A", 1},
				{"B", "2"},
				{"C", []string{"abc"}},
			},
		},
		"nested struct": {
			in: unionComplex{
				unionSimple: unionSimple{1, "2", []string{"abc"}},
				A:           2,
			},
			opts: structFieldOptions{TagKey: "json"},
			expect: []testUnionValues{
				{"B", "2"},
				{"C", []string{"abc"}},
				{"A", 2},
			},
		},
		"with TagKey unset": {
			in: unionTaggedComplex{
				unionSimple: unionSimple{1, "2", []string{"abc"}},
				unionTagged: unionTagged{3},
				B:           "3",
			},
			expect: []testUnionValues{
				{"A", 1},
				{"C", []string{"abc"}},
				{"ddbav", 3},
				{"B", "3"},
			},
		},
		"with TagKey json": {
			in: unionTaggedComplex{
				unionSimple: unionSimple{1, "2", []string{"abc"}},
				unionTagged: unionTagged{3},
				B:           "3",
			},
			opts: structFieldOptions{TagKey: "json"},
			expect: []testUnionValues{
				{"C", []string{"abc"}},
				{"A", 3},
				{"B", "3"},
			},
		},
		"with TagKey taga": {
			in: unionTaggedComplex{
				unionSimple: unionSimple{1, "2", []string{"abc"}},
				unionTagged: unionTagged{3},
				B:           "3",
			},
			opts: structFieldOptions{TagKey: "taga"},
			expect: []testUnionValues{
				{"A", 1},
				{"C", []string{"abc"}},
				{"TagA", 3},
				{"B", "3"},
			},
		},
		"with TagKey tagb": {
			in: unionTaggedComplex{
				unionSimple: unionSimple{1, "2", []string{"abc"}},
				unionTagged: unionTagged{3},
				B:           "3",
			},
			opts: structFieldOptions{TagKey: "tagb"},
			expect: []testUnionValues{
				{"A", 1},
				{"C", []string{"abc"}},
				{"TagB", 3},
				{"B", "3"},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			v := reflect.ValueOf(c.in)

			fields := unionStructFields(v.Type(), c.opts)
			for i, f := range fields.All() {
				expected := c.expect[i]
				if e, a := expected.Name, f.Name; e != a {
					t.Errorf("%d expect %v, got %v, %v", i, e, a, f)
				}
				actual := v.FieldByIndex(f.Index).Interface()
				if e, a := expected.Value, actual; !reflect.DeepEqual(e, a) {
					t.Errorf("%d expect %v, got %v, %v", i, e, a, f)
				}
			}
		})
	}
}

func TestCachedFields(t *testing.T) {
	type myStruct struct {
		Dog  int `tag1:"rabbit" tag2:"cow" tag3:"horse"`
		CAT  string
		bird bool
	}

	cases := map[string][]struct {
		Name      string
		FieldName string
		Found     bool
	}{
		"": {
			{"Dog", "Dog", true},
			{"dog", "Dog", true},
			{"DOG", "Dog", true},
			{"Yorkie", "", false},
			{"Cat", "CAT", true},
			{"cat", "CAT", true},
			{"CAT", "CAT", true},
			{"tiger", "", false},
			{"bird", "", false},
		},
		"tag1": {
			{"rabbit", "rabbit", true},
			{"Rabbit", "rabbit", true},
			{"cow", "", false},
			{"Cow", "", false},
			{"horse", "", false},
			{"Horse", "", false},
			{"Dog", "", false},
			{"dog", "", false},
			{"DOG", "", false},
			{"Cat", "CAT", true},
			{"cat", "CAT", true},
			{"CAT", "CAT", true},
			{"tiger", "", false},
			{"bird", "", false},
		},
		"tag2": {
			{"rabbit", "", false},
			{"Rabbit", "", false},
			{"cow", "cow", true},
			{"Cow", "cow", true},
			{"horse", "", false},
			{"Horse", "", false},
			{"Dog", "", false},
			{"dog", "", false},
			{"DOG", "", false},
			{"Cat", "CAT", true},
			{"cat", "CAT", true},
			{"CAT", "CAT", true},
			{"tiger", "", false},
			{"bird", "", false},
		},
		"tag3": {
			{"rabbit", "", false},
			{"Rabbit", "", false},
			{"cow", "", false},
			{"Cow", "", false},
			{"horse", "horse", true},
			{"Horse", "horse", true},
			{"Dog", "", false},
			{"dog", "", false},
			{"DOG", "", false},
			{"Cat", "CAT", true},
			{"cat", "CAT", true},
			{"CAT", "CAT", true},
			{"tiger", "", false},
			{"bird", "", false},
		},
	}

	for tagKey, cs := range cases {
		for _, c := range cs {
			name := tagKey
			if name == "" {
				name = "none"
			}
			t.Run(fmt.Sprintf("%s/%s", name, c.Name), func(t *testing.T) {
				t.Parallel()

				fields := unionStructFields(reflect.TypeOf(myStruct{}), structFieldOptions{
					TagKey: tagKey,
				})

				const expectedNumFields = 2
				if numFields := len(fields.All()); numFields != expectedNumFields {
					t.Errorf("expect %v fields, got %d", expectedNumFields, numFields)
				}

				f, found := fields.FieldByName(c.Name)
				if found != c.Found {
					t.Errorf("expect %v found, got %v", c.Found, found)
				}
				if found && f.Name != c.FieldName {
					t.Errorf("expect %v field name, got %s", c.FieldName, f.Name)
				}
			})
		}
	}
}
