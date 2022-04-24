package repository

type MaxLenQueue[T comparable] struct {
	slc         *[]T
	maxLen      int
	itemCreator func() T
}

func NewMaxLenQueue[T comparable](maxLen int, nilItemCreator func() T) MaxLenQueue[T] {
	slice := make([]T, 0)
	queue := MaxLenQueue[T]{
		slc:         &slice,
		maxLen:      maxLen,
		itemCreator: nilItemCreator,
	}
	return queue
}

func (q *MaxLenQueue[T]) Enqueue(item T) {
	if len(*q.slc) == q.maxLen {
		newSlice := append((*q.slc)[1:], item)
		q.slc = &newSlice
		return
	}
	newSlice := append((*q.slc), item)
	q.slc = &newSlice
}

func (q *MaxLenQueue[T]) Dequeue() T {
	if len(*q.slc) == 0 {
		return q.itemCreator()
	}
	top := q.Peek()
	newSlice := (*q.slc)[1:]
	q.slc = &newSlice
	return top
}

func (q *MaxLenQueue[T]) Peek() T {
	if len(*q.slc) == 0 {
		return q.itemCreator()
	}
	return (*q.slc)[0]
}

func (q *MaxLenQueue[T]) GetAll() []T {
	return (*q.slc)
}
