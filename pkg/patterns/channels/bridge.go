package channels

import "context"

// TODO: see if my piplines package would even work with the scanner
// TODO: idk if I properly understand how this integrates with the piplines
// Bridge() is used to consume values while maintaining order from a sequence of channels (<-chan <-chan interface{})/ "difference sources"
func Bridge[T any](ctx context.Context, chanStream <-chan <-chan T) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			var stream <-chan T
			select {
			case maybeStream, ok := <-chanStream: //TODO: do i undestand correctly, ok returns false if sends a closed signal?
				if ok == false {
					return
				}

				stream = maybeStream
			case <-ctx.Done():
				return
			}

			//TODO: i don't understand the orDone properly
			for val := range OrDone[T](ctx, stream) { //break out of loop if current stream is closed, to continue with next iteration of parent for loop and select a channel to read from
				select {
				case valStream <- val:
				case <-ctx.Done():
				}
			}
		}
	}()
	return valStream
}
