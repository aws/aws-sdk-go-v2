package enhancedclient

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestCreateExtensionContext(t *testing.T) {
	type extensionContextCreator interface {
		createExtensionContext() context.Context
	}

	makeContext := func(args ...any) context.Context {
		ctx := context.Background()

		for c := range len(args) / 2 {
			ctx = context.WithValue(ctx, args[c*2], args[c*2+1])
		}

		return ctx
	}

	keys := []any{
		CachedFieldsKey{},
		TableSchemaKey{},
	}

	cases := []struct {
		source   extensionContextCreator
		expected context.Context
	}{
		{
			source: &Schema[any]{},
			expected: makeContext(
				CachedFieldsKey{},
				(*CachedFields)(nil),
				TableSchemaKey{},
				&Schema[any]{},
			),
		},
		{
			source: &Schema[any]{
				cachedFields: &CachedFields{
					fields: []Field{
						{
							Tag:         Tag{},
							Name:        "",
							NameFromTag: false,
							Index:       nil,
							Type:        nil,
						},
					},
					fieldsByName: map[string]int{
						"": 0,
					},
				},
			},
			expected: makeContext(
				CachedFieldsKey{},
				&CachedFields{
					fields: []Field{
						{
							Tag:         Tag{},
							Name:        "",
							NameFromTag: false,
							Index:       nil,
							Type:        nil,
						},
					},
					fieldsByName: map[string]int{
						"": 0,
					},
				},
				TableSchemaKey{},
				&Schema[any]{
					cachedFields: &CachedFields{
						fields: []Field{
							{
								Tag:         Tag{},
								Name:        "",
								NameFromTag: false,
								Index:       nil,
								Type:        nil,
							},
						},
						fieldsByName: map[string]int{
							"": 0,
						},
					},
				},
			),
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := c.source.createExtensionContext()

			for _, k := range keys {
				if diff := cmpDiff(c.expected.Value(k), actual.Value(k)); len(diff) != 0 {
					t.Fatalf("unexpected diff: %s", diff)
				}
			}
		})
	}
}

type testStruct struct {
	beforeWrite int
	afterRead   int
}

func TestApplyBeforeWriteExtensions(t *testing.T) {
	cases := []struct {
		when       ExecutionPhase
		extensions []Extension
		input      any
		expected   any
		error      bool
	}{
		{
			extensions: []Extension{},
			error:      false,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			str := testStruct{}
			sch := &Schema[testStruct]{
				extensions: map[ExecutionPhase][]Extension{
					BeforeWrite: c.extensions,
				},
			}
			err := sch.applyAfterReadExtensions(&str)

			if !c.error && err != nil {
				t.Errorf("unexpected error: %v", err)

				return
			}

			if c.error && err == nil {
				t.Error("expected error")

				return
			}

			//if diff := cmpDiff(c.expected, av); len(diff) != 0 {
			//	t.Errorf("unexpected diff: %s", diff)
			//}
			_ = c
		})
	}
	_ = cases
}

func TestSchemaApplyExtension(t *testing.T) {
	if true {
		return
	}
	cases := []struct {
		when     ExecutionPhase
		actual   map[string]types.AttributeValue
		expected map[string]types.AttributeValue
		error    bool
	}{
		{
			when:   BeforeWrite,
			actual: map[string]types.AttributeValue{},
			expected: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{
					Value: "",
				},
			},
		},
	}

	buff := bytes.Buffer{}
	cryptorand.Reader = io.TeeReader(cryptorand.Reader, &buff)

	s, _ := NewSchema[order]()

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Logf("buffer too big: %d", len(buff.Bytes()))
			buff.Reset()
			t.Logf("buffer too big: %d", len(buff.Bytes()))

			actual, _ := s.Decode(c.actual)

			var err error
			switch c.when {
			case BeforeWrite:
				err = s.applyBeforeWriteExtensions(actual)
			case AfterRead:
				err = s.applyAfterReadExtensions(actual)
			//case BeforeQuery:
			//	err = s.apply(actual)
			//case BeforeScan:
			//	err = s.applyBeforeWriteExtensions(actual)
			default:
				t.Fatalf("i don't know how to handle: %s", c.when)
			}
			//err := s.applyExtensions(BeforeWrite, actual)

			t.Logf("buffer too big: %d", len(buff.Bytes()))

			fmt.Printf("%#+v\n", actual)

			b := buff.Bytes()[:]
			if len(b) != 16 {
				t.Fatalf("buffer too big: %d", len(b))
			}

			//b[6] = (b[6] & 0x0f) | 0x40
			//b[8] = (b[8] & 0x3f) | 0x80
			//c.expected["id"] = &types.AttributeValueMemberS{
			//	Value: fmt.Sprintf(
			//		"%x-%x-%x-%x-%x",
			//		b[0:4],
			//		b[4:6],
			//		b[6:8],
			//		b[8:10],
			//		b[10:16],
			//	),
			//}
			//fmt.Println(c.expected["id"].(*types.AttributeValueMemberS).Value)

			if !c.error && err != nil {
				t.Errorf("unexpected error: %v", err)

				return
			}

			if c.error && err == nil {
				t.Error("expected error")

				return
			}

			if diff := cmpDiff(c.expected, c.actual); len(diff) != 0 {
				e := json.NewEncoder(os.Stdout)
				e.SetIndent("", "  ")
				_ = e.Encode(c.actual)
				_ = e.Encode(c.expected)
				t.Errorf("unexpected diff: %s", diff)
			}
		})
	}
}
