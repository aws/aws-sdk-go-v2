package enhancedclient

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (s *Schema[T]) defaults() error {
	if s.billingMode == "" {
		s.billingMode = types.BillingModePayPerRequest
	}

	return nil
}

func (s *Schema[T]) resolveTableName() error {
	if s.typ == nil {
		s.typ = reflect.TypeFor[T]()
	}

	r := s.typ
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}

	if r.Kind() != reflect.Struct {
		return fmt.Errorf("resolveTableName() expected the type to be a struct or struct pointer, got: %V", reflect.New(s.typ).Interface())
	}

	s.tableName = pointer(r.Name())

	return nil
}

func (s *Schema[T]) resolveKeySchema() error {
	if s.keySchema != nil && len(s.keySchema) > 0 {
		return nil
	}

	var primary []string
	var sort []string

	for _, f := range s.cachedFields.fields {
		if f.Tag.Partition && f.Tag.Sort {
			return fmt.Errorf("Field %s is both primary and sort", f.Name)
		}

		if f.Tag.Partition {
			primary = append(primary, f.Name)
		}

		if f.Tag.Sort {
			sort = append(sort, f.Name)
		}
	}

	cp := len(primary)
	if cp != 1 {
		return fmt.Errorf("exactly 1 primary Field is expected, %d given, fields: %s", len(primary), strings.Join(primary, ", "))
	}

	cs := len(sort)
	if cs > 1 {
		return fmt.Errorf("exactly 0 or 1 sort Field is expected, %d given, fields: %s", len(sort), strings.Join(sort, ", "))
	}

	s.keySchema = make([]types.KeySchemaElement, cp+cs)
	s.keySchema[0].AttributeName = &primary[0]
	s.keySchema[0].KeyType = types.KeyTypeHash

	if cs > 0 {
		s.keySchema[1].AttributeName = &sort[0]
		s.keySchema[1].KeyType = types.KeyTypeRange
	}

	return nil
}

// []types.AttributeDefinition
func (s *Schema[T]) resolveAttributeDefinitions() error {
	for _, f := range s.cachedFields.fields {
		isKey := f.Tag.Partition || f.Tag.Sort
		for _, i := range f.Tag.Indexes {
			isKey = isKey || i.Partition || i.Sort
		}

		if !isKey {
			continue
		}

		at, ok := typeToScalarAttributeType(f.Type)
		if ok != true {
			continue
		}

		s.attributeDefinitions = append(s.attributeDefinitions, types.AttributeDefinition{
			AttributeName: &f.Name,
			AttributeType: at,
		})
	}

	return nil
}

func extractIndexes(fields []Field) (map[string][][]int, map[string][][]int, error) {
	globals := make(map[string][][]int)
	locals := make(map[string][][]int)
	unknowns := make(map[string][][]int)

	// collect index data
	for f, fld := range fields {
		for i, idx := range fld.Indexes {
			if idx.Global && idx.Local {
				return nil, nil, fmt.Errorf(`Field "%s" for index "%s" is configured to be both local and global`, fld.Name, idx.Name)
			}

			if idx.Partition && idx.Sort {
				return nil, nil, fmt.Errorf(`Field "%s" for index "%s" is configured to be both primarty and sort`, fld.Name, idx.Name)
			}

			if idx.Partition && idx.Local {
				return nil, nil, fmt.Errorf(`Field "%s" for index "%s" is configured to be the primarty key for a local index, local indexes inherit the primary from the table`, fld.Name, idx.Name)
			}

			pos := []int{f, i}

			switch {
			case idx.Global:
				globals[idx.Name] = append(globals[idx.Name], pos)
				break
			case idx.Local:
				locals[idx.Name] = append(locals[idx.Name], pos)
				break
			case !idx.Global && !idx.Local:
				unknowns[idx.Name] = append(unknowns[idx.Name], pos)
				break
			}
		}
	}

	for name, positions := range unknowns {
		_, gOk := globals[name]
		_, lOk := locals[name]

		if gOk && lOk {
			return nil, nil, fmt.Errorf(`index "%s" is configured both as global and local secondary index`, name)
		}
		if !gOk && !lOk {
			return nil, nil, fmt.Errorf(`index "%s" type cannot be determined`, name)
		}

		if gOk {
			globals[name] = append(globals[name], positions...)
		}
		if lOk {
			locals[name] = append(locals[name], positions...)
		}
	}

	return globals, locals, nil
}

func (s *Schema[T]) resolveSecondaryIndexes() error {
	globals, locals, err := extractIndexes(s.cachedFields.fields)
	if err != nil {
		return err
	}

	var tablePrimary *types.KeySchemaElement
	if len(s.keySchema) == 0 {
		if err := s.resolveKeySchema(); err != nil {
			return err
		}
	}
	for _, ks := range s.keySchema {
		if ks.KeyType == types.KeyTypeHash {
			tablePrimary = &ks
		}
	}

	if tablePrimary == nil {
		return fmt.Errorf("unable to determine the table primary key %v", s.TableName())
	}

	s.localSecondaryIndexes, err = processLSIs(s.cachedFields.fields, *tablePrimary, locals)
	if err != nil {
		return err
	}

	s.globalSecondaryIndexes, err = processGSIs(s.cachedFields.fields, globals)
	if err != nil {
		return err
	}

	return nil
}

func (s *Schema[T]) resolveDefaultExtensions() error {
	if s.extensions == nil {
		s.extensions = map[ExecutionPhase][]Extension{}
	} else {
		return nil
	}

	// register expression builder extensions first
	s.WithExtension(BeforeWrite, &VersionExtension[T]{})
	s.WithExtension(BeforeWrite, &AtomicCounterExtension[T]{})
	s.WithExtension(BeforeWrite, &AutogenerateExtension[T]{})

	return nil
}

func processGSIs(fields []Field, globals map[string][][]int) ([]types.GlobalSecondaryIndex, error) {
	gs := make([]types.GlobalSecondaryIndex, 0, len(globals))

	// build globals
	for name, positions := range globals {
		numPrimaries := 0
		numSorts := 0

		g := types.GlobalSecondaryIndex{
			IndexName: pointer(name),
			Projection: &types.Projection{
				NonKeyAttributes: nil,
				ProjectionType:   types.ProjectionTypeAll,
			},
		}

		for _, pos := range positions {
			f := fields[pos[0]]
			i := f.Indexes[pos[1]]

			switch {
			case i.Partition:
				g.KeySchema = append(g.KeySchema, types.KeySchemaElement{
					AttributeName: pointer(f.Name),
					KeyType:       types.KeyTypeHash,
				})
				numPrimaries++
				break
			case i.Sort:
				g.KeySchema = append(g.KeySchema, types.KeySchemaElement{
					AttributeName: pointer(f.Name),
					KeyType:       types.KeyTypeRange,
				})
				numSorts++
				break
			default:
				g.Projection.NonKeyAttributes = append(g.Projection.NonKeyAttributes, f.Name)
				g.Projection.ProjectionType = types.ProjectionTypeInclude
			}

			// the hash must be first
			if len(g.KeySchema) == 2 {
				slices.SortStableFunc(g.KeySchema, ksSortFunc)
			}
		}

		if numPrimaries != 1 {
			return nil, fmt.Errorf(`index "%s" has %d primary keys, it must have exactly 1`, name, numPrimaries)
		}

		if numSorts > 1 {
			return nil, fmt.Errorf(`index "%s" has %d sort keys, it must have exactly 0 or 1`, name, numSorts)
		}

		gs = append(gs, g)
	}

	return gs, nil
}

func ksSortFunc(a, b types.KeySchemaElement) int {
	switch types.KeyTypeHash {
	case a.KeyType:
		return -1
	case b.KeyType:
		return 1
	default:
		return 0
	}
}

func processLSIs(fields []Field, tablePrimary types.KeySchemaElement, locals map[string][][]int) ([]types.LocalSecondaryIndex, error) {
	ls := make([]types.LocalSecondaryIndex, 0, len(locals))
	numSorts := 0

	for name, positions := range locals {
		l := types.LocalSecondaryIndex{
			IndexName: pointer(name),
			KeySchema: []types.KeySchemaElement{
				tablePrimary,
			},
			Projection: &types.Projection{
				NonKeyAttributes: nil,
				ProjectionType:   types.ProjectionTypeAll,
			},
		}

		for _, pos := range positions {
			f := fields[pos[0]]
			i := f.Indexes[pos[1]]

			switch {
			case i.Partition:
				return nil, fmt.Errorf(`index "%s" has Field "%s" is configured as the primary key for a secondary index, secondary indexes inherit the primary key of their table`, name, f.Name)
			case i.Sort:
				l.KeySchema = append(l.KeySchema, types.KeySchemaElement{
					AttributeName: pointer(f.Name),
					KeyType:       types.KeyTypeRange,
				})
				numSorts++
			default:
				l.Projection.NonKeyAttributes = append(l.Projection.NonKeyAttributes, f.Name)
				l.Projection.ProjectionType = types.ProjectionTypeInclude
			}
		}

		ls = append(ls, l)
	}

	return ls, nil
}
