package enhancedclient

// ItemResult represents the result of a DynamoDB operation that returns an item or an error.
// Used in iterators for batch, scan, and query operations to convey either a successfully decoded item or an error.
type ItemResult[T any] struct {
	item *T    // The decoded item, if successful
	err  error // The error encountered, if any
}

// Item returns the decoded item, or nil if an error occurred.
func (it *ItemResult[T]) Item() *T {
	return it.item
}

// Error returns the error encountered during the operation, or nil if successful.
func (it *ItemResult[T]) Error() error {
	return it.err
}
