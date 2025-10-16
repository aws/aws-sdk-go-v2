package enhancedclient

import (
	"fmt"
	"reflect"

	"github.com/aws/smithy-go/middleware"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Schema[T any] struct {
	options      SchemaOptions
	cachedFields *CachedFields
	enc          *Encoder[T]
	dec          *Decoder[T]
	typ          reflect.Type

	//generators map[string]enhancedclient2.Generator[T]
	extensions   map[ExecutionPhase][]Extension
	interceptors middleware.Stack

	// common
	attributeDefinitions      []types.AttributeDefinition
	keySchema                 []types.KeySchemaElement
	tableName                 *string
	billingMode               types.BillingMode
	deletionProtectionEnabled *bool
	onDemandThroughput        *types.OnDemandThroughput
	provisionedThroughput     *types.ProvisionedThroughput
	sseSpecification          *types.SSESpecification
	streamSpecification       *types.StreamSpecification
	tableClass                types.TableClass
	warmThroughput            *types.WarmThroughput
	// create
	globalSecondaryIndexes []types.GlobalSecondaryIndex
	localSecondaryIndexes  []types.LocalSecondaryIndex
	resourcePolicy         *string
	tags                   []types.Tag
	// update
	multiRegionConsistency types.MultiRegionConsistency
	replicaUpdates         []types.ReplicationGroupUpdate
}

func (s *Schema[T]) createTableInput() (*dynamodb.CreateTableInput, error) {
	return &dynamodb.CreateTableInput{
		TableName:                 s.TableName(),
		KeySchema:                 s.KeySchema(),
		AttributeDefinitions:      s.AttributeDefinitions(),
		BillingMode:               s.BillingMode(),
		DeletionProtectionEnabled: s.DeletionProtectionEnabled(),
		GlobalSecondaryIndexes:    s.GlobalSecondaryIndexes(),
		LocalSecondaryIndexes:     s.LocalSecondaryIndexes(),
		OnDemandThroughput:        s.OnDemandThroughput(),
		ProvisionedThroughput:     s.ProvisionedThroughput(),
		ResourcePolicy:            s.ResourcePolicy(),
		SSESpecification:          s.SSESpecification(),
		StreamSpecification:       s.StreamSpecification(),
		TableClass:                s.TableClass(),
		Tags:                      s.Tags(),
		WarmThroughput:            s.WarmThroughput(),
	}, nil
}

func (s *Schema[T]) describeTableInput() (*dynamodb.DescribeTableInput, error) {
	return &dynamodb.DescribeTableInput{
		TableName: s.TableName(),
	}, nil
}

func (s *Schema[T]) deleteTableInput() (*dynamodb.DeleteTableInput, error) {
	return &dynamodb.DeleteTableInput{
		TableName: s.TableName(),
	}, nil
}

func (s *Schema[T]) createKeyMap(item *T) (Map, error) {
	m, err := s.Encode(item)
	if err != nil {
		return nil, err
	}

	for _, f := range s.cachedFields.fields {
		if !f.Partition && !f.Sort {
			delete(m, f.Name)
		}
	}

	return m, nil
}

func NewSchema[T any](fns ...func(options *SchemaOptions)) (*Schema[T], error) {
	if reflect.TypeFor[T]().Kind() != reflect.Struct {
		return nil, fmt.Errorf("NewClient() can only be created from structs, %T given", *new(T))
	}

	t := new(T)
	cf := unionStructFields(reflect.TypeOf(*t), structFieldOptions{})

	opts := SchemaOptions{}

	for _, fn := range fns {
		fn(&opts)
	}

	s := &Schema[T]{
		options:      opts,
		cachedFields: cf,
		typ:          reflect.TypeFor[T](),
		enc: NewEncoder[T](func(options *EncoderOptions) {
			options.ConverterRegistry = opts.ConverterRegistry
			options.IgnoreNilValueErrors = opts.IgnoreNilValueErrors
		}),
		dec: NewDecoder[T](func(options *DecoderOptions) {
			options.ConverterRegistry = opts.ConverterRegistry
			options.IgnoreNilValueErrors = opts.IgnoreNilValueErrors
		}),
	}

	resolversFns := []func(o *Schema[T]) error{
		(*Schema[T]).defaults,
		(*Schema[T]).resolveTableName,
		(*Schema[T]).resolveKeySchema,
		(*Schema[T]).resolveAttributeDefinitions,
		(*Schema[T]).resolveSecondaryIndexes,
		(*Schema[T]).resolveDefaultExtensions,
	}

	for _, fn := range resolversFns {
		if err := fn(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}
