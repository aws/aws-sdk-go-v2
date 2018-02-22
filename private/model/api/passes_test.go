// +build 1.6,codegen

package api

import (
	"reflect"
	"testing"
)

func TestUniqueInputAndOutputs(t *testing.T) {
	cases := [][]struct {
		expectedInput  string
		expectedOutput string
		operation      string
		input          string
		inputRef       string
		output         string
		outputRef      string
	}{
		{
			{
				expectedInput:  "FooOperationInput",
				expectedOutput: "FooOperationOutput",
				operation:      "FooOperation",
				input:          "FooInputShape",
				inputRef:       "FooInputShapeRef",
				output:         "FooOutputShape",
				outputRef:      "FooOutputShapeRef",
			},
			{
				expectedInput:  "BarOperationInput",
				expectedOutput: "BarOperationOutput",
				operation:      "BarOperation",
				input:          "FooInputShape",
				inputRef:       "FooInputShapeRef",
				output:         "FooOutputShape",
				outputRef:      "FooOutputShapeRef",
			},
		},
	}

	for _, c := range cases {
		a := &API{
			name:       "FooService",
			Operations: map[string]*Operation{},
		}

		expected := map[string][]string{}
		a.Shapes = map[string]*Shape{}
		for _, op := range c {
			a.Operations[op.operation] = &Operation{
				ExportedName: op.operation,
			}
			a.Operations[op.operation].Name = op.operation
			a.Operations[op.operation].InputRef = ShapeRef{
				API:       a,
				ShapeName: op.inputRef,
				Shape: &Shape{
					API:       a,
					ShapeName: op.input,
				},
			}
			a.Operations[op.operation].OutputRef = ShapeRef{
				API:       a,
				ShapeName: op.outputRef,
				Shape: &Shape{
					API:       a,
					ShapeName: op.output,
				},
			}

			a.Shapes[op.input] = &Shape{
				ShapeName: op.input,
			}
			a.Shapes[op.output] = &Shape{
				ShapeName: op.output,
			}

			expected[op.operation] = append(expected[op.operation], op.expectedInput)
			expected[op.operation] = append(expected[op.operation], op.expectedOutput)
		}

		a.fixStutterNames()
		a.renameToplevelShapes()
		for k, v := range expected {
			if a.Operations[k].InputRef.Shape.ShapeName != v[0] {
				t.Errorf("Error %s case: Expected %q, but received %q", k, v[0], a.Operations[k].InputRef.Shape.ShapeName)
			}
			if a.Operations[k].OutputRef.Shape.ShapeName != v[1] {
				t.Errorf("Error %s case: Expected %q, but received %q", k, v[1], a.Operations[k].OutputRef.Shape.ShapeName)
			}
		}

	}
}

func TestCollidingFields(t *testing.T) {
	cases := []struct {
		api      *API
		expected []*Shapes
	}{
		{
			&API{
				name: "FooService",
				Shapes: []*Shapes{
					{
						MemberRefs: map[string]*ShapeRef{
							"String":   {},
							"GoString": {},
							"Validate": {},
							"Foo":      {},
							"SetFoo":   {},
						},
					},
				},
			},
			[]*Shapes{
				{
					MemberRefs: map[string]*ShapeRef{
						"String_":   {},
						"GoString_": {},
						"Validate_": {},
						"Foo":       {},
						"SetFoo_":   {},
					},
				},
			},
		},
	}

	for _, c := range testCases {
		c.api.renameCollidingFields()
		if !reflect.DeepEqual(c.api.Shapes, c.expected) {
			t.Errorf("expected %v, but received %v", c.expected, c.api.Shapes)
		}
	}
}
