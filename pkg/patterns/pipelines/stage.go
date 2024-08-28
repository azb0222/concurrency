package pipelines

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
)

// Or could have interface called TeeStage, Stage, BridgeStage
// TODO: make sure pointers are returned everywhere they are supposed to be
// TODO: create a type alias for fn func(value T, args ...any) T
// TODO: private or public?
type Stage[T any] struct {
	Fn func(valueStream <-chan T, args ...any) <-chan T
}

//TODO: the only way to have this is to have a field on Stage: s.TeeFn

// newStage() will take in a function, and return a stage with context managed //TODO write a better comment
func newStage[T any](ctx context.Context, fn func(value T, args ...any) T) *Stage[T] {
	stageFn := func(valueStream <-chan T, args ...any) <-chan T {
		newValueStream := make(chan T)
		go func() {
			defer close(newValueStream)

			for i := range valueStream {
				select {
				case <-ctx.Done():
					return
				case newValueStream <- fn(i, args...):
				}
			}
		}()
		return newValueStream
	}

	return &Stage[T]{
		Fn: stageFn,
	}
}

//func (s *Stage[T]) Tee(ctx context.Context) {
//	out1, out2 := channels.Tee(ctx, s.Fn())
//}

// TODO: debug, crashing rn
func (s *Stage[T]) FanOutFanIn(ctx context.Context) Stage[T] {
	// fanInWorkers() will multiplex streams of data onto a single stream
	fanInWorkers := func(ctx context.Context, channels ...<-chan T) <-chan T {
		fmt.Println("Fanning In ...")
		var wg sync.WaitGroup
		multiplexedStream := make(chan T)

		multiplex := func(c <-chan T) {
			defer wg.Done()
			for i := range c {
				select {
				case <-ctx.Done():
					return
				case multiplexedStream <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	log.Println("Fanning Out ...")
	numWorkers := runtime.NumCPU()
	newStageFn := func(valueStream <-chan T, args ...any) <-chan T {
		workers := make([]<-chan T, numWorkers)
		for i := 0; i < numWorkers; i++ {
			workers[i] = s.Fn(valueStream, args...)
		}

		return fanInWorkers(ctx, workers...)
	}

	s.Fn = newStageFn
	return *s
}
