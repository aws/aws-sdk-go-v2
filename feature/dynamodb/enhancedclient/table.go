package enhancedclient

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type TableOptions[T any] struct {
	Schema            *Schema[T]
	DynamoDBOptions   []func(*dynamodb.Options)
	ExtensionRegistry *ExtensionRegistry[T]
}

type Table[T any] struct {
	client  Client
	options TableOptions[T]
}

func NewTable[T any](client Client, fns ...func(options *TableOptions[T])) (*Table[T], error) {
	if reflect.TypeFor[T]().Kind() != reflect.Struct {
		return nil, fmt.Errorf("NewClient() can only be created from structs, %T given", *new(T))
	}

	opts := TableOptions[T]{}

	for _, fn := range fns {
		fn(&opts)
	}

	defaultResolvers := []resolverFn[T]{
		resolveDefaultSchema[T],
		resolveDefaultExtensionRegistry[T],
	}

	for _, fn := range defaultResolvers {
		if err := fn(&opts); err != nil {
			return nil, err
		}
	}

	return &Table[T]{
		client:  client,
		options: opts,
	}, nil
}

func WithSchema[T any](schema *Schema[T]) func(options *TableOptions[T]) {
	return func(options *TableOptions[T]) {
		options.Schema = schema
	}
}

func WithExtensionRegistry[T any](registry *ExtensionRegistry[T]) func(options *TableOptions[T]) {
	return func(options *TableOptions[T]) {
		options.ExtensionRegistry = registry
	}
}

type resolverFn[T any] func(opts *TableOptions[T]) error

func resolveDefaultSchema[T any](opts *TableOptions[T]) error {
	if opts.Schema == nil {
		var err error
		opts.Schema, err = NewSchema[T]()
		if err != nil {
			return err
		}
	}

	return nil
}

func resolveDefaultExtensionRegistry[T any](opts *TableOptions[T]) error {
	if opts.ExtensionRegistry == nil {
		opts.ExtensionRegistry = DefaultExtensionRegistry[T]()
	}

	return nil
}
