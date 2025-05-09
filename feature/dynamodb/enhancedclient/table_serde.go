package enhancedclient

import (
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (s *Schema[T]) Encode(t *T) (map[string]types.AttributeValue, error) {
	v := reflect.ValueOf(t)
	out := map[string]types.AttributeValue{}

	for _, f := range s.cachedFields.All() {
		var fv reflect.Value
		var err error

		if f.tag.Getter != "" {
			m := v.MethodByName(f.tag.Getter)
			fv = m.Call([]reflect.Value{})[0]
		} else {
			fv, err = v.Elem().FieldByIndexErr(f.Index)
			if err != nil {
				if unwrap(s.options.ErrorOnMissingField) {
					return nil, err
				} else {
					continue
				}
			}
		}

		av, err := s.enc.encode(fv, f.tag)
		if err != nil && unwrap(s.options.ErrorOnMissingField) {
			return nil, err
		}

		out[f.Name] = av
	}

	return out, nil
}

func (s *Schema[T]) Decode(m map[string]types.AttributeValue) (*T, error) {
	t := new(T)
	v := reflect.ValueOf(t)

	for _, f := range s.cachedFields.All() {
		av, ok := m[f.Name]
		if !ok {
			continue
		}

		if f.tag.Setter != "" && f.tag.Getter != "" {
			current := v.MethodByName(f.tag.Getter).
				Call([]reflect.Value{})[0]

			if current.Kind() != reflect.Ptr {
				current = reflect.New(current.Type())
			}

			if err := s.dec.decode(av, current, f.tag); err != nil {
				return nil, err
			}

			v.MethodByName(f.tag.Setter).
				Call([]reflect.Value{
					current.Elem(),
				})

			continue
		}

		fv, err := v.Elem().FieldByIndexErr(f.Index)
		if err != nil {
			if unwrap(s.options.ErrorOnMissingField) {
				return nil, err
			} else {
				continue
			}
		}

		err = s.dec.decode(av, fv, f.tag)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}
