package enhancedclient

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestResolveTableName(t *testing.T) {
	cases := []struct {
		input    any
		expected *string
		error    bool
	}{
		{
			input:    &Schema[order]{},
			expected: pointer("order"),
		},
		{
			input:    &Schema[*order]{},
			expected: pointer("order"),
		},
		{
			input:    &Schema[address]{},
			expected: pointer("address"),
		},
		{
			input:    &Schema[*address]{},
			expected: pointer("address"),
		},
		{
			input:    &Schema[reflect.Value]{},
			expected: pointer("Value"),
		},
		{
			input:    &Schema[*reflect.Value]{},
			expected: pointer("Value"),
		},
		{
			input:    &Schema[any]{},
			expected: nil,
			error:    true,
		},
		{
			input:    &Schema[string]{},
			expected: nil,
			error:    true,
		},
		{
			input:    &Schema[[]byte]{},
			expected: nil,
			error:    true,
		},
		{
			input:    &Schema[[]order]{},
			expected: nil,
			error:    true,
		},
	}

	type tableNameResolver interface {
		TableName() *string
		resolveTableName() error
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var actual tableNameResolver
			var ok bool

			if actual, ok = c.input.(tableNameResolver); !ok && !c.error {
				t.Fatalf("unable to check the presence of the resolveTableName() error method")
			}

			err := actual.resolveTableName()

			if c.error && err == nil {
				t.Fatalf("expected error")
			}

			if !c.error && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if c.error && err != nil {
				return
			}

			if diff := cmpDiff(c.expected, actual.TableName()); len(diff) > 0 {
				t.Errorf(`failed to resolve table name, expected: %s"`, diff)
			}
		})
	}
}

func TestResolveKeySchema(t *testing.T) {
	pk := "pk"
	sk := "sk"

	cases := []struct {
		input    []Field
		expected any
		error    bool
	}{
		{
			input: []Field{
				{
					Name: "pk",
					Tag:  Tag{Partition: true},
				},
				{
					Name: "sk",
					Tag:  Tag{Sort: true},
				},
				{},
				{},
				{},
			},
			expected: []types.KeySchemaElement{
				{
					AttributeName: &pk,
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: &sk,
					KeyType:       types.KeyTypeRange,
				},
			},
			error: false,
		},
		{
			input: []Field{
				{
					Name: "pk",
					Tag:  Tag{Partition: true},
				},
				{},
				{},
				{},
			},
			expected: []types.KeySchemaElement{
				{
					AttributeName: &pk,
					KeyType:       types.KeyTypeHash,
				},
			},
			error: false,
		},
		{
			input: []Field{
				{
					Name: "pk",
					Tag:  Tag{Partition: true},
				},
				{
					Name: "sk",
					Tag:  Tag{Partition: true},
				},
				{},
				{},
				{},
			},
			expected: []types.KeySchemaElement{},
			error:    true,
		},
		{
			input: []Field{
				{
					Name: "pk",
					Tag:  Tag{Partition: true},
				},
				{
					Name: "sk",
					Tag:  Tag{Sort: true},
				},
				{
					Name: "sk",
					Tag:  Tag{Sort: true},
				},
				{},
				{},
			},
			expected: []types.KeySchemaElement{},
			error:    true,
		},
		{
			input: []Field{
				{
					Name: "pk",
				},
				{
					Name: "sk",
				},
				{},
				{},
			},
			expected: []types.KeySchemaElement{},
			error:    true,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			o := &Schema[order]{
				cachedFields: &CachedFields{
					fields: c.input,
				},
			}

			err := o.resolveKeySchema()
			if c.error && err == nil {
				t.Fatalf("expected error")
			}

			if !c.error && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if c.error && err != nil {
				return
			}

			if diff := cmpDiff(c.expected, o.KeySchema()); len(diff) != 0 {
				t.Fatalf("unexpected diff: %s", diff)
			}
		})
	}
}

func TestResolveAttributeDefinitions(t *testing.T) {
	cases := []struct {
		input    []Field
		expected []types.AttributeDefinition
	}{
		{
			input:    []Field{},
			expected: []types.AttributeDefinition(nil),
		},
		{
			input: []Field{
				{
					Name: "pk",
					Tag: Tag{
						Partition: true,
					},
					Type: reflect.TypeFor[string](),
				},
			},
			expected: []types.AttributeDefinition{
				{
					AttributeName: pointer("pk"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
		},
		{
			input: []Field{
				{
					Name: "pk",
					Tag: Tag{
						Partition: true,
					},
					Type: reflect.TypeFor[int32](),
				},
			},
			expected: []types.AttributeDefinition{
				{
					AttributeName: pointer("pk"),
					AttributeType: types.ScalarAttributeTypeN,
				},
			},
		},
		{
			input: []Field{
				{
					Name: "pk",
					Tag: Tag{
						Partition: true,
					},
					Type: reflect.TypeFor[[]byte](),
				},
			},
			expected: []types.AttributeDefinition{
				{
					AttributeName: pointer("pk"),
					AttributeType: types.ScalarAttributeTypeB,
				},
			},
		},
		{
			input: []Field{
				{
					Name: "sk",
					Tag: Tag{
						Sort: true,
					},
					Type: reflect.TypeFor[[]byte](),
				},
			},
			expected: []types.AttributeDefinition{
				{
					AttributeName: pointer("sk"),
					AttributeType: types.ScalarAttributeTypeB,
				},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			s := &Schema[any]{
				cachedFields: &CachedFields{
					fields: c.input,
				},
			}

			_ = s.resolveAttributeDefinitions()

			if diff := cmpDiff(c.expected, s.AttributeDefinitions()); len(diff) != 0 {
				fmt.Printf("%#+v\n", c.expected)
				fmt.Printf("%#+v\n", s.AttributeDefinitions())
				t.Fatalf("unexpected diff: %s", diff)
			}
		})
	}
}

func TestResolveSecondaryIndexes(t *testing.T) {
	cases := []struct {
		input        []Field
		expectedLSIs []types.LocalSecondaryIndex
		expectedGSIs []types.GlobalSecondaryIndex
		error        bool
	}{
		{
			error: true,
		},
		{
			input: []Field{},
			error: true,
		},
		{
			input: []Field{
				{
					Name: "pk",
					Tag: Tag{
						Partition: true,
						Indexes: []Index{
							{
								Name:      "gsi1",
								Global:    true,
								Partition: true,
							},
							{
								Name: "gsi2",
								Sort: true,
							},
						},
					},
				},
				{
					Name: "sk",
					Tag: Tag{
						Indexes: []Index{
							{
								Name: "gsi1",
								Sort: true,
							},
							{
								Name:      "gsi2",
								Global:    true,
								Partition: true,
							},
						},
					},
				},
			},
			expectedLSIs: []types.LocalSecondaryIndex{},
			expectedGSIs: []types.GlobalSecondaryIndex{
				{
					IndexName: pointer("gsi1"),
					Projection: &types.Projection{
						NonKeyAttributes: nil,
						ProjectionType:   types.ProjectionTypeAll,
					},
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: pointer("pk"),
							KeyType:       types.KeyTypeHash,
						},
						{
							AttributeName: pointer("sk"),
							KeyType:       types.KeyTypeRange,
						},
					},
				},
				{
					IndexName: pointer("gsi2"),
					Projection: &types.Projection{
						NonKeyAttributes: nil,
						ProjectionType:   types.ProjectionTypeAll,
					},
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: pointer("sk"),
							KeyType:       types.KeyTypeHash,
						},
						{
							AttributeName: pointer("pk"),
							KeyType:       types.KeyTypeRange,
						},
					},
				},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			o := &Schema[order]{
				cachedFields: &CachedFields{
					fields: c.input,
				},
			}

			err := o.resolveSecondaryIndexes()
			if c.error && err == nil {
				t.Fatalf("expected error")
			}

			if !c.error && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if c.error && err != nil {
				return
			}

			if diff := cmpDiff(c.expectedGSIs, o.GlobalSecondaryIndexes()); len(diff) != 0 {
				t.Fatalf("unexpected diff in GSIs: %s", diff)
			}

			if diff := cmpDiff(c.expectedLSIs, o.LocalSecondaryIndexes()); len(diff) != 0 {
				t.Fatalf("unexpected diff in LSIs: %s", diff)
			}
		})
	}
}

func TestResolveGlobalSecondaryIndexUpdates(t *testing.T) {
}
