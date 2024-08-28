package waitGroup

import (
	"fmt"
	"sync"
	"time"
)

/*
	sync.WaitGroup will block until a set of concurrent operations complete.
	use sync.WaitGroup if we do not care about the result of the concurrent operations or plan to collect them in some other way.
		otherwise use channels and select.
*/

func Main() {
	var wg sync.WaitGroup

	wg.Add(1) //wg counter++
	go func() {
		defer wg.Done()
		fmt.Println("1st goroutine sleeping ...")
		time.Sleep(1 * time.Second)
	}()

	wg.Add(1) //wg counter++
	go func() {
		defer wg.Done()
		fmt.Println("2nd goroutine sleeping ...")
		time.Sleep(1 * time.Second)
	}()

	wg.Wait() //blocks until wg counter == 0
	fmt.Println("All goroutines completed.")
}
