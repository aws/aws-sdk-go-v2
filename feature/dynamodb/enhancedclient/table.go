package enhancedclient

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// TableOptions provides configuration options for a DynamoDB Table.
//
// T is the type of the items stored in the table.
type TableOptions[T any] struct {
	// Schema defines the schema for the table, including attribute mapping and validation.
	Schema *Schema[T]

	// DynamoDBOptions is a list of functions to customize the underlying DynamoDB client options.
	DynamoDBOptions []func(*dynamodb.Options)

	// ExtensionRegistry holds the registry of extensions to be used with the table.
	ExtensionRegistry *ExtensionRegistry[T]

	// MaxConsecutiveErrors sets the maximum number of consecutive errors allowed during batch, query, or scan operations.
	// If this threshold is exceeded, the operation will stop and return.
	// If set to 0, the default value of DefaultMaxConsecutiveErrors will be used.
	MaxConsecutiveErrors uint
}

// DefaultMaxConsecutiveErrors is the fallback value used for MaxConsecutiveErrors when it is set to 0.
// A value of 1 means the operation will stop after the first error.
const DefaultMaxConsecutiveErrors uint = 1

// Table represents a strongly-typed DynamoDB table for items of type T.
//
// It provides methods for interacting with DynamoDB using the provided client and options.
type Table[T any] struct {
	// client is the DynamoDB client used to perform operations on the table.
	client Client

	// options holds the configuration options for the table.
	options TableOptions[T]
}

// NewTable creates a new Table for items of type T using the provided client and configuration functions.
//
// The configuration functions can be used to customize the TableOptions before the table is created.
// Returns an error if T is not a struct type or if required options cannot be resolved.
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
		resolveDefaultMaxConsecutiveErrors[T],
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

// WithSchema returns a configuration function that sets the Schema for TableOptions.
//
// Use this to specify a custom schema when creating a Table.
func WithSchema[T any](schema *Schema[T]) func(options *TableOptions[T]) {
	return func(options *TableOptions[T]) {
		options.Schema = schema
	}
}

// WithExtensionRegistry returns a configuration function that sets the ExtensionRegistry for TableOptions.
//
// Use this to specify a custom extension registry when creating a Table.
func WithExtensionRegistry[T any](registry *ExtensionRegistry[T]) func(options *TableOptions[T]) {
	return func(options *TableOptions[T]) {
		options.ExtensionRegistry = registry
	}
}

// WithMaxConsecutiveErrors returns a configuration function that sets the MaxConsecutiveErrors option for TableOptions.
//
// Use this to specify the maximum number of consecutive errors allowed during batch, query, or scan operations.
// A value of 0 means no limit is enforced.
// WithMaxConsecutiveErrors returns a configuration function that sets the MaxConsecutiveErrors option for TableOptions.
//
// Use this to specify the maximum number of consecutive errors allowed during batch, query, or scan operations.
// If set to 0, the default value of DefaultMaxConsecutiveErrors will be used.
func WithMaxConsecutiveErrors[T any](maxConsecutiveErrors uint) func(options *TableOptions[T]) {
	return func(options *TableOptions[T]) {
		options.MaxConsecutiveErrors = maxConsecutiveErrors
	}
}

// resolverFn defines a function type for resolving or setting default options on TableOptions.
type resolverFn[T any] func(opts *TableOptions[T]) error

// resolveDefaultSchema sets a default schema on TableOptions if none is provided.
//
// Returns an error if the schema cannot be created.
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

// resolveDefaultExtensionRegistry sets a default extension registry on TableOptions if none is provided.
func resolveDefaultExtensionRegistry[T any](opts *TableOptions[T]) error {
	if opts.ExtensionRegistry == nil {
		opts.ExtensionRegistry = DefaultExtensionRegistry[T]()
	}

	return nil
}

// resolveDefaultMaxConsecutiveErrors sets MaxConsecutiveErrors to DefaultMaxConsecutiveErrors
// if it is not explicitly set (i.e., if the value is 0).
// This ensures a sensible default for error handling in batch, query, or scan operations.
func resolveDefaultMaxConsecutiveErrors[T any](opts *TableOptions[T]) error {
	if opts.MaxConsecutiveErrors == 0 {
		opts.MaxConsecutiveErrors = DefaultMaxConsecutiveErrors
	}
	return nil
}
