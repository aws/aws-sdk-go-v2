package enhancedclient

type ItemResult[T any] struct {
	item T
	err  error
}

func (it *ItemResult[T]) Item() T {
	return it.item
}

func (it *ItemResult[T]) Error() error {
	return it.err
}
