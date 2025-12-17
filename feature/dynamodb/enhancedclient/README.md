# AWS SDK Go V2 - High Level Client

This package provides an high level DynamoDB client for [AWS SDK Go v2](https://github.com/aws/aws-sdk-go-v2), featuring a flexible data mapper layer. It simplifies object mapping, schema management, and table operations, enabling idiomatic Go struct-to-table mapping, lifecycle hooks, and extension support for DynamoDB applications.

## Features

- [Automated struct to schema mapping](#automated-struct-to-schema-mapping)
- [Table management](#table-management)
- [Item operations](#item-operations)
- [Extensions](#extensions)
- [Custom converters](#custom-converters)

## Automated Struct to Schema Mapping

The `Schema[T]` type supports advanced table configuration options, allowing you to fine-tune your DynamoDB tables:

- **Provisioned and On-Demand Throughput:**
	- Use `WithProvisionedThroughput` or `WithOnDemandThroughput` to set capacity modes.
- **Table Class:**
	- Use `WithTableClass` to set the table class (e.g., Standard, StandardInfrequentAccess).
- **Tags:**
	- Use `WithTags` to add tags to your table.
- **Resource Policy, Encryption, Streams, and More:**
	- Methods like `WithResourcePolicy`, `WithSSESpecification`, `WithStreamSpecification`, and `WithWarmThroughput` allow further customization.

**Example:**

```go
type Product struct {
		ProductID string  `dynamodbav:"product_id,partition"`
		Category  string  `dynamodbav:",sort"`
		Price     float64
		InStock   bool    `dynamodbav:",omitempty"`
}

// default usage
table := enhancedclient.NewTable[Product](client)

// customized schema options
schema := enhancedclient.NewSchema[Product]()
schema = schema.WithProvisionedThroughput(&types.ProvisionedThroughput{
		ReadCapacityUnits:  5,
		WriteCapacityUnits: 5,
})
schema = schema.WithTableClass(types.TableClassStandardInfrequentAccess)
schema = schema.WithTags([]types.Tag{{Key: aws.String("env"), Value: aws.String("prod")}})
table := enhancedclient.NewTable(client, enhancedclient.WithSchema(schema))
```

The `Table[T]` type provides a high-level, type-safe interface for managing DynamoDB tables. It abstracts away much of the boilerplate required for table lifecycle management, making it easier to work with DynamoDB in Go.

**Key table management methods:**

- `Create(ctx context.Context) (*dynamodb.CreateTableOutput, error)`: Creates the DynamoDB table based on the inferred or provided schema.
- `CreateWithWait(ctx context.Context, maxWaitDur time.Duration) error`: Creates the table and waits until it becomes active, or until the specified timeout is reached.
- `Describe(ctx context.Context) (*dynamodb.DescribeTableOutput, error)`: Retrieves metadata and status information about the table.
- `Delete(ctx context.Context) (*dynamodb.DeleteTableOutput, error)`: Deletes the DynamoDB table.
- `DeleteWithWait(ctx context.Context, maxWaitDur time.Duration) error`: Deletes the table and waits until it is fully removed or the timeout is reached.
- `Exists(ctx context.Context) (bool, error)`: Checks if the table exists and is accessible.

These features allow you to manage the full lifecycle of your DynamoDB tables in a concise, idiomatic Go style, while leveraging the full power of DynamoDB's management capabilities.

> **Note:** Only table level scans are supported at the moment.

> **Note:** Table schema updates (such as adding or modifying attributes, indexes, or throughput settings after table creation) are not supported at this time. Only the table management functions listed above are available. Support for table updates is planned for a future release.

## Item Operations

The `Table[T]` type provides a set of strongly-typed methods for common item-level operations, making it easy to interact with DynamoDB records as native Go structs. These methods handle marshaling and unmarshaling, key construction, and error handling, so you can focus on your application logic.

**Key item operations:**

- `GetItem(ctx, key, ...) (*T, error)`: Retrieve a single item by its key.
- `GetItemWithProjection(ctx, key, projection, ...) (*T, error)`: Retrieve a single item by key, returning only the specified attributes.
- `PutItem(ctx, item, ...) (*T, error)`: Insert or replace an item in the table.
- `UpdateItem(ctx, item, ...) (*T, error)`: Update an existing item, using the struct as the source of changes.
- `DeleteItem(ctx, item, ...) error`: Delete an item by providing the struct value.
- `DeleteItemByKey(ctx, key, ...) error`: Delete an item by its key.
- `Scan(ctx, expr, ...) iter.Seq[ItemResult[T]]`: Scan the table with a filter expression, returning an iterator over results.
- `ScanIndex(ctx, indexName, expr, ...) iter.Seq[ItemResult[T]]`: Scan the index with a filter expression, returning an iterator over results.
- `Query(ctx, expr, ...) iter.Seq[ItemResult[T]]`: Query the table or an index using a key condition expression, returning an iterator over results.
- `QueryIndex(ctx, indexName, expr, ...) iter.Seq[ItemResult[T]]`: Query the index or an index using a key condition expression, returning an iterator over results.

**Batch operations:**

- `CreateBatchWriteOperation() *BatchWriteOperation[T]`: Returns a new batch write operation, allowing you to queue multiple put and delete requests and execute them efficiently in batches. Handles chunking, retries for unprocessed items, and respects DynamoDB's batch size limits.
    - Use `AddPut(item *T)` or `AddRawPut(map[string]types.AttributeValue)` to queue items for writing.
    - Use `AddDelete(item *T)` or `AddRawDelete(map[string]types.AttributeValue)` to queue items for deletion.
    - Call `Execute(ctx, ...)` to perform the batch write.

- `CreateBatchGetOperation() *BatchGetOperation[T]`: Returns a new batch get operation, allowing you to queue multiple keys for retrieval and execute them in a single batch request. Handles chunking, retries for unprocessed keys, and respects DynamoDB's batch size limits.
    - Use `AddReadItem(item *T)` or `AddReadItemByMap(map[string]types.AttributeValue)` to queue keys for retrieval.
    - Call `Execute(ctx, ...)` to perform the batch get, which yields results as an iterator.

Batch operations are useful for efficiently processing large numbers of items, minimizing network calls, and handling DynamoDB's batch constraints automatically.

These methods are designed to be ergonomic and safe, leveraging Go's type system to reduce boilerplate and runtime errors when working with DynamoDB items.

**Iterators and ItemResult:**

Many methods, such as `Scan`, `Query`, and `BatchGetOperation.Execute`, return an iterator in the form of an `iter.Seq[ItemResult[T]]`, which is a function that accepts a callback. Each callback invocation receives an `ItemResult[T]` containing either a successfully decoded item or an error encountered during retrieval or decoding.

When consuming these iterators, use the callback or range pattern and always check the `Error()` method on each result before using the item:

```go
// Callback-based iteration (idiomatic for iter.Seq):
table.Scan(ctx, expr, ...)(func(result ItemResult[T]) bool {
		if err := result.Error(); err != nil {
				// handle error, e.g. log or collect
				return true // continue to next result
		}
		item := result.Item()
		// process item, e.g. append to a slice or print
		return true // continue, or return false to stop early
})
// Alternative: idiomatic Go for-range over the iterator:
for res := range table.Scan(ctx, expr, ...) {
		if err := res.Error(); err != nil {
				// handle error
				continue
		}
		item := res.Item()
		// process item
}
```

This pattern ensures robust error handling and makes it easy to process large result sets efficiently and safely.

## Extensions

The enhanced client supports an extension system that allows you to inject custom logic at key points in the item lifecycle. Extensions can be used for auditing, validation, automatic field population, versioning, atomic counters, and more.

### Extension Registry and Lifecycle Hooks

Extensions are registered using the `ExtensionRegistry`, which manages hooks for different operation phases:

- **BeforeReader / AfterReader:** Invoked before/after reading an item (e.g., `GetItem`).
- **BeforeWriter / AfterWriter:** Invoked before/after writing an item (e.g., `PutItem`, `UpdateItem`).

You can register multiple extensions for each phase. The registry supports method chaining for easy configuration.

**Example: Registering extensions**

```go
reg := &enhancedclient.ExtensionRegistry[Product]{}
reg.AddBeforeReader(&MyAuditExtension{}).
    AddAfterReader(&MyAuditExtension{}).
    AddBeforeWriter(&MyValidationExtension{})

table := enhancedclient.NewTable[Product](client, enhancedclient.WithExtensionRegistry(reg))
```

### Built-in Extensions

The default registry includes useful extensions for common patterns:

- **AutogenerateExtension:** Automatically populates fields such as UUIDs or timestamps.
- **AtomicCounterExtension:** Handles atomic increment/decrement fields.
- **VersionExtension:** Implements optimistic versioning for concurrency control.

You can use the default registry or customize it as needed:

```go
table := enhancedclient.NewTable[Product](
    client,
    enhancedclient.WithExtensionRegistry(
        enhancedclient.DefaultExtensionRegistry[Product](),
    ),
)
```

### Writing a Custom Extension

To create your own extension, implement one or more of the extension interfaces (e.g., `BeforeWriter`, `AfterReader`). Each hook receives the context and the item, and can return an error to abort the operation (for "before" hooks).

```go
type MyAuditExtension struct{}

func (a *MyAuditExtension) BeforeWrite(ctx context.Context, v *Product) error {
    log.Printf("Audit: about to write item: %+v", v)
    return nil
}

func (a *MyAuditExtension) AfterRead(ctx context.Context, v *Product) error {
    log.Printf("Audit: read item: %+v", v)
    return nil
}


ext := &MyAuditExtension{}
registry := DefaultExtensionRegistry[order]().Clone()
registry.AddBeforeWriter(ext)
registry.AddAfterReader(ext)

table, err := NewTable[order](
    c,
    WithSchema(sch),
    WithExtensionRegistry(registry),
)
```

### Advanced: Expression Builders

Extensions can also participate in building DynamoDB expressions (conditions, filters, projections, updates) by implementing the relevant builder interfaces. This allows for powerful customization of query and update logic.

See the source and tests for more advanced extension usage patterns.

### Built-in extensions available in the default registry

The default extension registry includes several built-in extensions that provide common DynamoDB patterns out of the box:


- **AutogenerateExtension**
    - Automatically populates struct fields marked for auto-generation. This includes generating UUIDs for primary keys, setting timestamps for created/updated fields, or populating other values at write time.
    - **Use cases:**
        - Automatically generate unique IDs for new items:
            ```go
            type Order struct {
                ID        string    `dynamodbav:"id,partition,autogenerated|key"`
                CreatedAt string    `dynamodbav:"created_at,autogenerated|timestamp"`
            }
            // On PutItem, ID and CreatedAt will be set if empty.
            ```
        - Set audit fields (created/updated timestamps) without manual code.
    - **How it works:**
        - Fields with the `autogenerated|key` or `autogenerated|timestamp` tag option are detected and set by the extension before writing.

- **AtomicCounterExtension**
    - Enables atomic increment or decrement of numeric fields marked as atomic counters. This is useful for counters, sequence numbers, or version fields that must be updated safely in concurrent environments.
    - **Use cases:**
        - Track the number of times an item is accessed or updated:
            ```go
            type PageView struct {
                URL     string `dynamodbav:"url,partition"`
                Counter int64  `dynamodbav:"counter,atomiccounter"`
            }
            // On UpdateItem, Counter can be atomically incremented.
            ```
        - Maintain a version or sequence number for items.
    - **How it works:**
        - Fields with the `atomiccounter` tag option are updated using DynamoDB's atomic update expressions, ensuring thread-safe increments/decrements.

- **VersionExtension**
    - Implements optimistic concurrency control by managing a version field on your items. This helps prevent lost updates and ensures that concurrent writes do not overwrite each other unintentionally.
    - **Use cases:**
        - Add a version field to your struct to enable safe concurrent updates:
            ```go
            type Document struct {
                DocID   string `dynamodbav:"doc_id,partition"`
                Version int64  `dynamodbav:"version,version"`
            }
            // On each update, Version is checked and incremented.
            ```
        - Prevent accidental overwrites in collaborative or distributed systems.
    - **How it works:**
        - Fields with the `version` tag option are checked and incremented on each write. If the version in the database does not match the expected value, the write fails, preventing lost updates.

These extensions are automatically included when you use `DefaultExtensionRegistry`:

```go
table := enhancedclient.NewTable[Product](
    client,
    enhancedclient.WithExtensionRegistry(
        enhancedclient.DefaultExtensionRegistry[Product](),
    ),
)
```

You can also clone and customize the registry to add or remove extensions as needed for your application.

**Note:** Because extensions can modify how your data is processed, the extension registry is not enabled by default. Enable it explicitly if you want to use these features.

## Custom converters

The enhanced client includes a set of built-in converters for common Go types (booleans, numbers, strings, time, JSON, byte arrays, pointers, etc.), making most struct fields work out of the box.

For advanced scenarios, you can define custom converters to handle complex or non-standard data types in your structs. By implementing the `AttributeConverter` interface, you control how a field is marshaled to and unmarshaled from DynamoDB attribute values.

**Implementing a custom converter:**

```go
type MyCustomType struct {
    Value string
}

type MyCustomConverter struct{}

func (c MyCustomConverter) ToAttributeValue(v any) (types.AttributeValue, error) {
    t, ok := v.(MyCustomType)
    if !ok {
        return nil, fmt.Errorf("expected MyCustomType")
    }
    return &types.AttributeValueMemberS{Value: "custom:" + t.Value}, nil
}

func (c MyCustomConverter) FromAttributeValue(av types.AttributeValue) (any, error) {
    s, ok := av.(*types.AttributeValueMemberS)
    if !ok {
        return nil, fmt.Errorf("expected string attribute")
    }
    return MyCustomType{Value: strings.TrimPrefix(s.Value, "custom:")}, nil
}
```

**Registering a custom converter:**

```go
my_registry := converters.DefaultRegistry.Clone()
// or
my_registry := converters.NewRegistry()

// register converter
my_registry.Add("my_custom_converter", &MyCustomConverter{})

schema := enhancedclient.NewSchema[MyStruct]func(options *SchemaOptions) {
    options.ConverterRegistry = my_registry
})

// add it to the struct
type MyStruct struct {
    ID            string        `dynamodbav:"id,partition"`
    CustomField   MyCustomType  `dynamodbav:"custom_field,converter|my_custom_converter"`
}
```

For more details and built-in examples, see the `converters/` directory in the source tree.

## Example Usage

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/enhancedclient"
)

type Product struct {
    ProductID string  `dynamodbav:"product_id,partition"`
    Category  string  `dynamodbav:",sort"`
    Price     float64
    InStock   bool    `dynamodbav:",omitempty"`
}

func main() {
    cfg, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }
    client := dynamodb.NewFromConfig(cfg)

    // Create the table with schema inference
    table, err := enhancedclient.NewTable[Product](client)
    if err != nil {
        log.Fatalf("failed to create table: %v", err)
    }
    if err := table.CreateWithWait(context.Background(), 2*time.Minute); err != nil {
        log.Fatalf("failed to create table: %v", err)
    }

    // Put an item
    prod := &Product{ProductID: "p1", Category: "books", Price: 19.99, InStock: true}
    _, err = table.PutItem(context.Background(), prod)
    if err != nil {
        log.Fatalf("failed to put item: %v", err)
    }

    // Get the item
    key := enhancedclient.Map{}.With("product_id", "p1").With("category", "books")
    got, err := table.GetItem(context.Background(), key)
    if err != nil {
        log.Fatalf("failed to get item: %v", err)
    }
    log.Printf("Got item: %+v", got)

    // Scan all items
    for res := range table.Scan(context.Background(), enhancedclient.ScanExpression{}) {
        if err := res.Error(); err != nil {
            log.Printf("scan error: %v", err)
            continue
        }
        item := res.Item()
        log.Printf("Scanned item: %+v", item)
    }
}
```
