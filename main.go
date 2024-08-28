package main

import (
	"asritha.dev/concurrency/pkg/patterns/channels"
	"asritha.dev/concurrency/pkg/syncExamples/signals"
	"context"
	"fmt"
	"sync"
)

func main() {

	/*
		thinking through concurrent scanner:

		seems like bridge is alternative to FanIn?

		ex pipeline:

		input[] order and output[] order should match, which won't happen if we use FanOutFanIn
		stage 1: add by 1
		stage 2: multiply by 2
		stage 3: subtract by 1

	*/
	signals.CondBroadcastEx()

}

func IDKWTFISGOINGONHEREDEBUGTHISTODO() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	test := []int{1, 2, 3, 4}
	c := make(chan int)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for num := range test {
			select {
			case <-ctx.Done():
				return
			case c <- num:
			}
		}
	}()

	out1, out2 := channels.Tee[int](ctx, c)
	for val1 := range out1 {
		fmt.Printf("out1: %v, out2: %v", val1, <-out2)
	}
	wg.Wait()
}
