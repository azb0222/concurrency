package channels

import "context"

// TODO: can I use this anywhere in the ./pkg/patterns/pipelines package?
// used in cases where its unclear if the channel we are reading from is cancelled when the goroutine is cancelled
// TODO: set breakpoints to make sure I understand how this actually works
func OrDone[T any](ctx context.Context, c <-chan T) <-chan T {
	valStream := make(chan T)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-ctx.Done():
				}
			}
		}
	}()
	return valStream
}
