package enhancedclient

// ExtensionRegistry manages a set of extension hooks for enhanced DynamoDB
// operations on a given type T. It allows registration of pre- and post-processing
// logic for read and write operations, enabling features such as automatic field
// population, versioning, atomic counters, and custom business logic.
//
// Extensions are grouped by operation type:
//   - beforeReaders: Invoked before reading an item (e.g., GetItem).
//   - afterReaders: Invoked after reading an item.
//   - beforeWriters: Invoked before writing an item (e.g., PutItem, UpdateItem).
//   - afterWriters: Invoked after writing an item.
//
// The registry supports method chaining for extension registration.
// DefaultExtensionRegistry provides a registry pre-populated with common extensions.
type ExtensionRegistry[T any] struct {
	// GetItem
	beforeReaders []BeforeReader[T]
	afterReaders  []AfterReader[T]
	// PutItem | UpdateItem
	beforeWriters []BeforeWriter[T]
	afterWriters  []AfterWriter[T]
}

// AddBeforeReader registers a BeforeReader extension to be invoked before
// reading an item. Returns the registry for method chaining.
func (er *ExtensionRegistry[T]) AddBeforeReader(br BeforeReader[T]) *ExtensionRegistry[T] {
	er.beforeReaders = append(er.beforeReaders, br)

	return er
}

// AddAfterReader registers an AfterReader extension to be invoked after
// reading an item. Returns the registry for method chaining.
func (er *ExtensionRegistry[T]) AddAfterReader(ar AfterReader[T]) *ExtensionRegistry[T] {
	er.afterReaders = append(er.afterReaders, ar)

	return er
}

// AddBeforeWriter registers a BeforeWriter extension to be invoked before
// writing an item. Returns the registry for method chaining.
func (er *ExtensionRegistry[T]) AddBeforeWriter(bw BeforeWriter[T]) *ExtensionRegistry[T] {
	er.beforeWriters = append(er.beforeWriters, bw)

	return er
}

// AddAfterWriter registers an AfterWriter extension to be invoked after
// writing an item. Returns the registry for method chaining.
func (er *ExtensionRegistry[T]) AddAfterWriter(aw AfterWriter[T]) *ExtensionRegistry[T] {
	er.afterWriters = append(er.afterWriters, aw)

	return er
}

// Clone creates a new ExtensionRegistry containing copies of all registered
// extensions for type T. The returned registry has independent extension slices,
// so further modifications (adding/removing extensions) do not affect the original.
//
// Note: The extensions themselves are not deep-copied; only the slice references
// are duplicated. If extensions maintain internal state, that state will be shared.
//
// Returns a pointer to the new ExtensionRegistry.
func (er ExtensionRegistry[T]) Clone() *ExtensionRegistry[T] {
	out := &ExtensionRegistry[T]{}

	out.beforeReaders = append(out.beforeReaders, er.beforeReaders...)
	out.afterReaders = append(out.afterReaders, er.afterReaders...)
	out.beforeWriters = append(out.beforeWriters, er.beforeWriters...)
	out.afterWriters = append(out.afterWriters, er.afterWriters...)

	return out
}

// DefaultExtensionRegistry returns a new ExtensionRegistry pre-populated with
// common beforeWriter extensions: AutogenerateExtension, AtomicCounterExtension,
// and VersionExtension. These provide automatic key/timestamp population,
// atomic counter updates, and optimistic versioning for write operations.
func DefaultExtensionRegistry[T any]() *ExtensionRegistry[T] {
	out := &ExtensionRegistry[T]{}

	out.beforeWriters = append(
		out.beforeWriters,
		&AutogenerateExtension[T]{},
		&AtomicCounterExtension[T]{},
		&VersionExtension[T]{},
	)

	return out
}
