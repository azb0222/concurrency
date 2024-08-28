package channels

import "context"

// TODO: rename everywhere to use without func? look at reutrn value if thats what the author did?
// TODO: why did the author have this?  func Tee(done <-chan interface{}, in <-chan interface{}) (_,_ <-chan interface{}) { <-chan interface{}} {
func Tee[T any](ctx context.Context, in <-chan T) (_, _ <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)

	go func() {
		defer close(out1)
		defer close(out2)
		//each iteration over in will not continue until both out1 and out2 have been written to
		for val := range OrDone[T](ctx, in) { //orDone is returning <-chan interface{} while performing the logic underneath the hood to check for done status
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ { //two iterations to write to both out1 and out2
				select {
				case <-ctx.Done():
				case out1 <- val:
					out1 = nil //once a channel is written to, set to nil to block further writes
				case out2 <- val:
					out2 = nil //once a channel is written to, set to nil to block further writes
				}
			}
		}
	}()

	return out1, out2
}
