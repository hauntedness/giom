package lists

import "sync"

type Buffer[T any] struct {
	buffer chan T
	sync.WaitGroup
}

func (b *Buffer[T]) Send(t T) {
	b.buffer <- t
}

func (b *Buffer[T]) Receive() T {
	return <-b.buffer
}
