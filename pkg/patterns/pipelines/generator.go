package pipelines

import "context"

// TODO: seperate out generator into files
type IGenerator[T any] interface {
	GetValues(ctx context.Context) <-chan T
}

type StreamGenerator[T any] struct {
	values []T
}

// TODO: does this need to take in context? if errored? cancel context?
func NewStreamGenerator[T any](values ...T) *StreamGenerator[T] {
	return &StreamGenerator[T]{
		values: values,
	}
}

func (s *StreamGenerator[T]) GetValues(ctx context.Context) <-chan T {
	valuesStream := make(chan T)
	go func() {
		defer close(valuesStream)
		for _, i := range s.values {
			select {
			case <-ctx.Done():
				return
			case valuesStream <- i:
			}
		}
	}()
	return valuesStream
}

type RepeatGenerator[T any] struct {
	repeatFn func(args ...any) T
	cap      int
}

func NewRepeatGenerator[T any](repeatFn func(args ...any) T, cap int) *RepeatGenerator[T] {
	return &RepeatGenerator[T]{
		repeatFn: repeatFn,
		cap:      cap,
	}
}

func (ig *RepeatGenerator[T]) GetValues(ctx context.Context) <-chan T {
	return ig.take(ctx, ig.callRepeatFn(ctx))
}

// TODO: the args ...any isn't getting called
// repeatFn() generator will pass fn() into valueStream forever until the Done channel is closed
func (ig *RepeatGenerator[T]) callRepeatFn(ctx context.Context) <-chan T {
	valueStream := make(chan T)
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-ctx.Done():
				return
			case valueStream <- ig.repeatFn():
			}
		}
	}()

	return valueStream
}

// take() generator will take the first cap items off of its incoming valueStream, then exit
func (ig *RepeatGenerator[T]) take(ctx context.Context, valueStream <-chan T) <-chan T {
	takeStream := make(chan T)
	go func() {
		defer close(takeStream)

		for i := 0; i < ig.cap; i++ {
			select {
			case <-ctx.Done():
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}
